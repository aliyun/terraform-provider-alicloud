package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsLogtailConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsLogtailConfigCreate,
		Read:   resourceAliCloudSlsLogtailConfigRead,
		Delete: resourceAliCloudSlsLogtailConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"input_detail": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareArrayJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"input_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"plugin", "file"}, false),
			},
			"last_modify_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"log_sample": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"logtail_config_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"output_detail": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logstore_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"output_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudSlsLogtailConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/configs")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project_name").(string))
	if v, ok := d.GetOk("logtail_config_name"); ok {
		request["configName"] = v
	}

	if v, ok := d.GetOk("log_sample"); ok {
		request["logSample"] = v
	}
	request["inputType"] = d.Get("input_type")
	inputDetail := d.Get("input_detail").(string)
	request["inputDetail"] = NormalizeMap(convertJsonStringToObject(inputDetail))
	request["outputType"] = d.Get("output_type")
	if v, ok := d.GetOk("create_time"); ok {
		request["createTime"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("output_detail"); v != nil {
		endpoint1, _ := jsonpath.Get("$[0].endpoint", v)
		if endpoint1 != nil && endpoint1 != "" {
			objectDataLocalMap["endpoint"] = endpoint1
		}
		region1, _ := jsonpath.Get("$[0].region", v)
		if region1 != nil && region1 != "" {
			objectDataLocalMap["region"] = region1
		}
		logstoreName1, _ := jsonpath.Get("$[0].logstore_name", v)
		if logstoreName1 != nil && logstoreName1 != "" {
			objectDataLocalMap["logstoreName"] = logstoreName1
		}

		request["outputDetail"] = objectDataLocalMap
	}

	if v, ok := d.GetOkExists("last_modify_time"); ok {
		request["lastModifyTime"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateConfig", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_logtail_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["configName"]))

	return resourceAliCloudSlsLogtailConfigRead(d, meta)
}

func resourceAliCloudSlsLogtailConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsLogtailConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_logtail_config DescribeSlsLogtailConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	if objectRaw["inputDetail"] != nil {
		d.Set("input_detail", convertMapToJsonStringIgnoreError(objectRaw["inputDetail"].(map[string]interface{})))
	}
	d.Set("input_type", objectRaw["inputType"])
	d.Set("last_modify_time", objectRaw["lastModifyTime"])
	d.Set("log_sample", objectRaw["logSample"])
	d.Set("output_type", objectRaw["outputType"])
	d.Set("logtail_config_name", objectRaw["configName"])

	outputDetailMaps := make([]map[string]interface{}, 0)
	outputDetailMap := make(map[string]interface{})
	outputDetailRaw := make(map[string]interface{})
	if objectRaw["outputDetail"] != nil {
		outputDetailRaw = objectRaw["outputDetail"].(map[string]interface{})
	}
	if len(outputDetailRaw) > 0 {
		outputDetailMap["endpoint"] = outputDetailRaw["endpoint"]
		outputDetailMap["logstore_name"] = outputDetailRaw["logstoreName"]
		outputDetailMap["region"] = outputDetailRaw["region"]

		outputDetailMaps = append(outputDetailMaps, outputDetailMap)
	}
	if err := d.Set("output_detail", outputDetailMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("project_name", parts[0])

	return nil
}

func resourceAliCloudSlsLogtailConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	configName := parts[1]
	action := fmt.Sprintf("/configs/%s", configName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteConfig", action), query, nil, nil, hostMap, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
