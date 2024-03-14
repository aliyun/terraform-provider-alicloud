package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strconv"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudLogtailConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogtailConfigCreate,
		Read:   resourceAlicloudLogtailConfigRead,
		Update: resourceAlicloudLogtailConfiglUpdate,
		Delete: resourceAlicloudLogtailConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"input_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"file", "plugin"}, false),
			},
			"log_sample": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"LogService"}, false),
			},
			"input_detail": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					yaml, _ := normalizeJsonString(v)
					return yaml
				},
				ValidateFunc: validation.ValidateJsonString,
			},
		},
	}
}

func resourceAlicloudLogtailConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var inputConfigInputDetail = make(map[string]interface{})
	data := d.Get("input_detail").(string)
	if jsonErr := json.Unmarshal([]byte(data), &inputConfigInputDetail); jsonErr != nil {
		return WrapError(jsonErr)
	}
	var requestInfo *sls.Client
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
	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		requestInfo = slsClient
		sls.AddNecessaryInputConfigField(inputConfigInputDetail)
		covertInput, covertErr := assertInputDetailType(inputConfigInputDetail, logconfig)
		if covertErr != nil {
			return nil, WrapError(covertErr)
		}
		logconfig.InputDetail = covertInput
		return nil, slsClient.CreateConfig(d.Get("project").(string), logconfig)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_logtail_config", "CreateConfig", AliyunLogGoSdkERROR)
	}
	if debugOn() {
		addDebug("CreateConfig", raw, requestInfo, map[string]interface{}{
			"project": d.Get("project").(string),
			"config":  logconfig,
		})
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("logstore").(string), COLON_SEPARATED, d.Get("name").(string)))
	return resourceAlicloudLogtailConfigSave(d, meta)
}

func getConfig(d *schema.ResourceData, meta interface{}) (*sls.LogConfig, error) {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	config, err := logService.DescribeLogtailConfig(d.Id())
	return config, err
}

func isLogtailConfigUnconvertable(d *schema.ResourceData, err error, split []string) bool {
	if IsExpectedErrors(err, []string{"unconvertable", "Unconvertable"}) {
		d.Set("project", split[0])
		d.Set("logstore", split[1])
		d.Set("name", split[2])
		d.Set("input_detail", "[Warning!] Server configuration is an unconvertable pipeline config.")
		d.Set("last_modify_time", "")
		return true
	}
	return false
}

func resourceAlicloudLogtailConfigSave(d *schema.ResourceData, meta interface{}) error {
	split := strings.Split(d.Id(), COLON_SEPARATED)
	config, err := getConfig(d, meta)
	if err != nil {
		if isLogtailConfigUnconvertable(d, err, split) {
			return nil
		}
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		d.Set("last_modify_time", "")
		return WrapError(err)
	}

	// Save the input_detail in tf to state, instead of saving the server-side config's input_detail.
	lastModifyTimeNew := strconv.Itoa(int(config.LastModifyTime))
	d.Set("last_modify_time", lastModifyTimeNew)
	d.Set("project", split[0])
	d.Set("name", config.Name)
	d.Set("logstore", split[1])
	d.Set("input_type", config.InputType)
	d.Set("log_sample", config.LogSample)
	d.Set("output_type", config.OutputType)
	return nil
}

// The return value of this function will be saved to the state during import and refresh, and used for diff during plan and apply.
func resourceAlicloudLogtailConfigRead(d *schema.ResourceData, meta interface{}) error {
	split := strings.Split(d.Id(), COLON_SEPARATED)
	config, err := getConfig(d, meta)
	if err != nil {
		if isLogtailConfigUnconvertable(d, err, split) {
			return nil
		}
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		d.Set("last_modify_time", "")
		return WrapError(err)
	}

	// Determine if there is "last_modify_time" in the state, if yes, then retrieve it; otherwise, leave it empty.
	var lastModifyTimeOld string
	if val, ok := d.GetOk("last_modify_time"); ok {
		lastModifyTimeOld = val.(string)
	}
	// Because the server will return redundant parameters, we filter here
	inputDetail := d.Get("input_detail").(string)
	var oMap map[string]interface{}
	json.Unmarshal([]byte(inputDetail), &oMap)
	nMap := config.InputDetail.(map[string]interface{})
	if inputDetail != "" && !strings.HasPrefix(inputDetail, "[Warning!]") {
		for nk := range nMap {
			if _, ok := oMap[nk]; !ok {
				delete(nMap, nk)
			}
		}
	}
	nMapJson, err := json.Marshal(nMap)
	if err != nil {
		d.Set("last_modify_time", "")
		return WrapError(err)
	}
	lastModifyTimeNew := strconv.Itoa(int(config.LastModifyTime))
	if len(lastModifyTimeOld) == 0 {
		// When Terraform is upgraded from a lower version or when executing terraform import, if lastModifyTimeOld is empty, set input_detail to the server-side config's input_detail.
		d.Set("input_detail", string(nMapJson))
	} else if lastModifyTimeNew != lastModifyTimeOld {
		// If the last_modify_time in state is different from that on the server, it means that the server configuration has been changed. The input_detail of input needs to be set as the input_detail of the server configuration.
		d.Set("input_detail", string(nMapJson))
	} else {
		// If the last_modify_time in state is the same as that on the server, it means that there is no change in the server configuration. The input_detail and state remain consistent.
		d.Set("input_detail", inputDetail)
	}
	d.Set("last_modify_time", lastModifyTimeNew)
	d.Set("project", split[0])
	d.Set("name", config.Name)
	d.Set("logstore", split[1])
	d.Set("input_type", config.InputType)
	d.Set("log_sample", config.LogSample)
	d.Set("output_type", config.OutputType)
	return nil
}

func resourceAlicloudLogtailConfiglUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	update := false
	if d.HasChange("log_sample") {
		update = true
	}
	if d.HasChange("input_detail") {
		update = true
	}
	if d.HasChange("input_type") {
		update = true
	}
	if update {
		logconfig := &sls.LogConfig{
			InputType: d.Get("input_type").(string),
		}
		inputConfigInputDetail := make(map[string]interface{})
		data := d.Get("input_detail").(string)
		conver_err := json.Unmarshal([]byte(data), &inputConfigInputDetail)
		if conver_err != nil {
			old, _ := d.GetChange("input_detail")
			d.Set("input_detail", old)
			return WrapError(conver_err)
		}
		sls.AddNecessaryInputConfigField(inputConfigInputDetail)
		covertInput, covertErr := assertInputDetailType(inputConfigInputDetail, logconfig)
		if covertErr != nil {
			old, _ := d.GetChange("input_detail")
			d.Set("input_detail", old)
			return WrapError(covertErr)
		}
		logconfig.InputDetail = covertInput

		client := meta.(*connectivity.AliyunClient)
		var requestInfo *sls.Client
		params := &sls.LogConfig{
			Name:        parts[2],
			LogSample:   d.Get("log_sample").(string),
			InputType:   d.Get("input_type").(string),
			OutputType:  d.Get("output_type").(string),
			InputDetail: logconfig.InputDetail,
			OutputDetail: sls.OutputDetail{
				ProjectName:  d.Get("project").(string),
				LogStoreName: d.Get("logstore").(string),
			},
		}
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.UpdateConfig(parts[0], params)
		})
		if err != nil {
			old, _ := d.GetChange("input_detail")
			d.Set("input_detail", old)
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateConfig", AliyunLogGoSdkERROR)
		}
		if debugOn() {
			addDebug("UpdateConfig", raw, requestInfo, map[string]interface{}{
				"project": parts[0],
				"config":  params,
			})
		}
	}
	return resourceAlicloudLogtailConfigSave(d, meta)
}

func resourceAlicloudLogtailConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *sls.Client
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteConfig(parts[0], parts[2])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("DeleteConfig", raw, requestInfo, map[string]string{
				"project": parts[0],
				"config":  parts[2],
			})
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ProjectNotExist", "LogStoreNotExist", "ConfigNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteConfig", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogtailConfig(d.Id(), Deleted, DefaultTimeout))
}

// This function is used to assert and convert the type to sls.LogConfig
func assertInputDetailType(inputConfigInputDetail map[string]interface{}, logconfig *sls.LogConfig) (sls.InputDetailInterface, error) {
	if inputConfigInputDetail["logType"] == "json_log" {
		JSONConfigInputDetail, ok := sls.ConvertToJSONConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, WrapError(Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = JSONConfigInputDetail
	}
	if inputConfigInputDetail["logType"] == "apsara_log" {
		ApsaraLogConfigInputDetail, ok := sls.ConvertToApsaraLogConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, WrapError(Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = ApsaraLogConfigInputDetail
	}

	if inputConfigInputDetail["logType"] == "common_reg_log" {
		RegexConfigInputDetail, ok := sls.ConvertToRegexConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, WrapError(Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = RegexConfigInputDetail
	}
	if inputConfigInputDetail["logType"] == "delimiter_log" {
		DelimiterConfigInputDetail, ok := sls.ConvertToDelimiterConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, WrapError(Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = DelimiterConfigInputDetail
	}
	if logconfig.InputType == "plugin" {
		PluginLogConfigInputDetail, ok := sls.ConvertToPluginLogConfigInputDetail(inputConfigInputDetail)
		if ok != true {
			return nil, WrapError(Error("covert to JSONConfigInputDetail false "))
		}
		logconfig.InputDetail = PluginLogConfigInputDetail
	}
	return logconfig.InputDetail, nil
}
