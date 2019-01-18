package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"time"
)

func resourceAlicoudLogtailConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicoudLogtaiConfiglCreate,
		Read:   resourceAlicoudLogtailConfigRead,
		Update: resourceAlicoudLogtaiConfiglUpdate,
		Delete: resourceAlicoudLogtailConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"input_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"log_sample": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"input_detail": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					yaml, _ := normalizeJsonString(v)
					return yaml
				},
				ValidateFunc: validateJsonString,
			},
		},
	}
}

func resourceAlicoudLogtaiConfiglCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var inputConfigInputDetail = make(map[string]interface{})
	data := d.Get("input_detail").(string)
	if err := json.Unmarshal([]byte(data), &inputConfigInputDetail); err != nil {
		return fmt.Errorf("Input detail covert to string get an error: %#v.", err)
	}
	_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		logconfig := &sls.LogConfig{
			Name:       d.Get("name").(string),
			LogSample:  d.Get("log_sample").(string),
			InputType:  d.Get("input_type").(string),
			OutputType: d.Get("output_type").(string),
			OutputDetail: sls.OutputDetail{
				ProjectName:  d.Get("project").(string),
				LogStoreName: d.Get("logstore").(string),
			},
		}
		if covert_input, covert_err := assertInputDetailType(inputConfigInputDetail, logconfig); covert_err != nil {
			return nil, covert_err
		} else {
			logconfig.InputDetail = covert_input
		}
		return nil, slsClient.CreateConfig(d.Get("project").(string), logconfig)
	})
	if err != nil {
		return fmt.Errorf("CreateLogtailConfig got an error: %#v.", err)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("logstore").(string), COLON_SEPARATED, d.Get("name").(string)))
	return resourceAlicoudLogtailConfigRead(d, meta)
}

func resourceAlicoudLogtailConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)
	config, err := logService.DescribeLogLogtailConfig(split[0], split[2])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		fmt.Println(fmt.Errorf("DescribeLogLogtailConfig got an error: %#v.", err))
		return fmt.Errorf("DescribeLogLogtailConfig got an error: %#v.", err)
	}

	d.Set("project", split[0])
	d.Set("config_name", config.Name)
	d.Set("logstore", split[1])
	d.Set("input_type", config.InputType)
	d.Set("input_detail", config.InputDetail)
	d.Set("log_sample", config.LogSample)
	d.Set("output_type", config.OutputType)
	return nil
}

func resourceAlicoudLogtaiConfiglUpdate(d *schema.ResourceData, meta interface{}) error {
	split := strings.Split(d.Id(), COLON_SEPARATED)

	update := false
	if d.HasChange("log_sample") {
		update = true
		d.SetPartial("log_sample")
	}
	if d.HasChange("input_detail") {
		update = true
		d.SetPartial("input_detail")
	}
	if d.HasChange("input_type") {
		update = true
		d.SetPartial("input_type")
	}
	if update {
		logconfig := &sls.LogConfig{}
		inputConfigInputDetail := make(map[string]interface{})
		data := d.Get("input_detail").(string)
		if err := json.Unmarshal([]byte(data), &inputConfigInputDetail); err != nil {
			return fmt.Errorf("InputDetail convert got an error: %#v.", err)
		}
		if covert_input, covert_err := assertInputDetailType(inputConfigInputDetail, logconfig); covert_err != nil {
			return covert_err
		} else {
			logconfig.InputDetail = covert_input
		}

		client := meta.(*connectivity.AliyunClient)
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.UpdateConfig(split[0], &sls.LogConfig{
				Name:        split[2],
				LogSample:   d.Get("log_sample").(string),
				InputType:   d.Get("input_type").(string),
				OutputType:  d.Get("output_type").(string),
				InputDetail: logconfig.InputDetail,
				OutputDetail: sls.OutputDetail{
					ProjectName:  d.Get("project").(string),
					LogStoreName: d.Get("logstore").(string),
				},
			})
		})
		if err != nil {
			return fmt.Errorf("UpdateLogTailConfig %s got an error: %#v.", split[2], err)
		}
	}
	return resourceAlicoudLogtaiConfiglUpdate(d, meta)
}

func resourceAlicoudLogtailConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteConfig(split[0], split[2])
		})
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting logtail config %s got an error: %#v", split[2], err))
		}
		if _, err := logService.DescribeLogLogtailConfig(split[0], split[2]); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("Deleting logtail config %s timeout.", split[2]))
	})
}

// This function is used to assert and convert the type to sls.LogConfig
func assertInputDetailType(inputConfigInputDetail map[string]interface{}, logconfig *sls.LogConfig) (sls.InputDetailInterface, error) {
	if inputConfigInputDetail["logType"] == "json_log" {
		JSONConfigInputDetail, ok := sls.ConvertToJSONConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, fmt.Errorf("covert to JSONConfigInputDetail false ")
		}
		logconfig.InputDetail = JSONConfigInputDetail
	}
	if inputConfigInputDetail["logType"] == "apsara_log" {
		ApsaraLogConfigInputDetail, ok := sls.ConvertToApsaraLogConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, fmt.Errorf("covert to JSONConfigInputDetail false ")
		}
		logconfig.InputDetail = ApsaraLogConfigInputDetail
	}

	if inputConfigInputDetail["logType"] == "common_reg_log" {
		RegexConfigInputDetail, ok := sls.ConvertToRegexConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, fmt.Errorf("covert to JSONConfigInputDetail false ")
		}
		logconfig.InputDetail = RegexConfigInputDetail
	}
	if inputConfigInputDetail["logType"] == "delimiter_log" {
		DelimiterConfigInputDetail, ok := sls.ConvertToDelimiterConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, fmt.Errorf("covert to JSONConfigInputDetail false ")
		}
		logconfig.InputDetail = DelimiterConfigInputDetail
	}
	if logconfig.InputType == "plugin" {
		PluginLogConfigInputDetail, ok := sls.ConvertToPluginLogConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, fmt.Errorf("covert to JSONConfigInputDetail false ")
		}
		logconfig.InputDetail = PluginLogConfigInputDetail
	}
	return logconfig.InputDetail, nil
}

