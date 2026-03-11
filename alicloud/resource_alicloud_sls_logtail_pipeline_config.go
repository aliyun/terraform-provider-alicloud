// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsLogtailPipelineConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsLogtailPipelineConfigCreate,
		Read:   resourceAliCloudSlsLogtailPipelineConfigRead,
		Update: resourceAliCloudSlsLogtailPipelineConfigUpdate,
		Delete: resourceAliCloudSlsLogtailPipelineConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_sample": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"global": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"filepath", "machine_group_topic", "custom"}, false),
						},
						"topic_format": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable_timestamp_nanosecond": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"pipeline_meta_tag_key": {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := compareArrayJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
					},
				},
			},
			"inputs": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareArrayJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"processors": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareArrayJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"aggregators": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareArrayJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"flushers": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareArrayJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudSlsLogtailPipelineConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/pipelineconfigs")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project_name").(string))
	
	if v, ok := d.GetOk("config_name"); ok {
		request["configName"] = v
	}

	if v, ok := d.GetOk("log_sample"); ok {
		request["logSample"] = v
	}

	// Global configuration
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("global"); v != nil {
		topicType1, _ := jsonpath.Get("$[0].topic_type", v)
		if topicType1 != nil && topicType1 != "" {
			objectDataLocalMap["TopicType"] = topicType1
		}
		topicFormat1, _ := jsonpath.Get("$[0].topic_format", v)
		if topicFormat1 != nil && topicFormat1 != "" {
			objectDataLocalMap["TopicFormat"] = topicFormat1
		}
		enableTimestampNanosecond1, _ := jsonpath.Get("$[0].enable_timestamp_nanosecond", v)
		if enableTimestampNanosecond1 != nil {
			objectDataLocalMap["EnableTimestampNanosecond"] = enableTimestampNanosecond1
		}
		pipelineMetaTagKey1, _ := jsonpath.Get("$[0].pipeline_meta_tag_key", v)
		if pipelineMetaTagKey1 != nil && pipelineMetaTagKey1 != "" {
			objectDataLocalMap["PipelineMetaTagKey"] = NormalizeMap(convertJsonStringToObject(pipelineMetaTagKey1.(string)))
		}

		if len(objectDataLocalMap) > 0 {
			request["global"] = objectDataLocalMap
		}
	}

	// Inputs (required, JSON array string)
	if v, ok := d.GetOk("inputs"); ok {
		inputs := v.(string)
		inputsList, err := convertJsonStringToList(inputs)
		if err != nil {
			return WrapError(err)
		}
		request["inputs"] = inputsList
	}

	// Processors (optional, JSON array string)
	if v, ok := d.GetOk("processors"); ok {
		processors := v.(string)
		processorsList, err := convertJsonStringToList(processors)
		if err != nil {
			return WrapError(err)
		}
		request["processors"] = processorsList
	}

	// Aggregators (optional, JSON array string)
	if v, ok := d.GetOk("aggregators"); ok {
		aggregators := v.(string)
		aggregatorsList, err := convertJsonStringToList(aggregators)
		if err != nil {
			return WrapError(err)
		}
		request["aggregators"] = aggregatorsList
	}

	// Flushers (required, JSON array string)
	if v, ok := d.GetOk("flushers"); ok {
		flushers := v.(string)
		flushersList, err := convertJsonStringToList(flushers)
		if err != nil {
			return WrapError(err)
		}
		request["flushers"] = flushersList
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateLogtailPipelineConfig", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_logtail_pipeline_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["configName"]))

	return resourceAliCloudSlsLogtailPipelineConfigRead(d, meta)
}

func resourceAliCloudSlsLogtailPipelineConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsLogtailPipelineConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_logtail_pipeline_config DescribeSlsLogtailPipelineConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config_name", objectRaw["configName"])
	d.Set("log_sample", objectRaw["logSample"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("last_modify_time", objectRaw["lastModifyTime"])

	// Set global configuration
	if objectRaw["global"] != nil {
		globalMaps := make([]map[string]interface{}, 0)
		globalMap := make(map[string]interface{})
		globalRaw := objectRaw["global"].(map[string]interface{})
		
		if topicType, ok := globalRaw["TopicType"]; ok {
			globalMap["topic_type"] = topicType
		}
		if topicFormat, ok := globalRaw["TopicFormat"]; ok {
			globalMap["topic_format"] = topicFormat
		}
		if enableTimestampNanosecond, ok := globalRaw["EnableTimestampNanosecond"]; ok {
			globalMap["enable_timestamp_nanosecond"] = enableTimestampNanosecond
		}
		if pipelineMetaTagKey, ok := globalRaw["PipelineMetaTagKey"]; ok {
			globalMap["pipeline_meta_tag_key"] = convertMapToJsonStringIgnoreError(pipelineMetaTagKey.(map[string]interface{}))
		}

		if len(globalMap) > 0 {
			globalMaps = append(globalMaps, globalMap)
		}
		if err := d.Set("global", globalMaps); err != nil {
			return err
		}
	}

	// Set inputs
	if objectRaw["inputs"] != nil {
		if inputsJson, err := json.Marshal(objectRaw["inputs"]); err == nil {
			d.Set("inputs", string(inputsJson))
		}
	}

	// Set processors
	if objectRaw["processors"] != nil {
		if processorsJson, err := json.Marshal(objectRaw["processors"]); err == nil {
			d.Set("processors", string(processorsJson))
		}
	}

	// Set aggregators
	if objectRaw["aggregators"] != nil {
		if aggregatorsJson, err := json.Marshal(objectRaw["aggregators"]); err == nil {
			d.Set("aggregators", string(aggregatorsJson))
		}
	}

	// Set flushers
	if objectRaw["flushers"] != nil {
		if flushersJson, err := json.Marshal(objectRaw["flushers"]); err == nil {
			d.Set("flushers", string(flushersJson))
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("project_name", parts[0])

	return nil
}

func resourceAliCloudSlsLogtailPipelineConfigUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	parts := strings.Split(d.Id(), ":")
	configName := parts[1]
	action := fmt.Sprintf("/pipelineconfigs/%s", configName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	if d.HasChange("log_sample") {
		if v, ok := d.GetOk("log_sample"); ok {
			request["logSample"] = v
		}
	}

	if d.HasChange("global") {
		objectDataLocalMap := make(map[string]interface{})
		if v := d.Get("global"); v != nil {
			topicType1, _ := jsonpath.Get("$[0].topic_type", v)
			if topicType1 != nil && topicType1 != "" {
				objectDataLocalMap["TopicType"] = topicType1
			}
			topicFormat1, _ := jsonpath.Get("$[0].topic_format", v)
			if topicFormat1 != nil && topicFormat1 != "" {
				objectDataLocalMap["TopicFormat"] = topicFormat1
			}
			enableTimestampNanosecond1, _ := jsonpath.Get("$[0].enable_timestamp_nanosecond", v)
			if enableTimestampNanosecond1 != nil {
				objectDataLocalMap["EnableTimestampNanosecond"] = enableTimestampNanosecond1
			}
			pipelineMetaTagKey1, _ := jsonpath.Get("$[0].pipeline_meta_tag_key", v)
			if pipelineMetaTagKey1 != nil && pipelineMetaTagKey1 != "" {
				objectDataLocalMap["PipelineMetaTagKey"] = NormalizeMap(convertJsonStringToObject(pipelineMetaTagKey1.(string)))
			}

			if len(objectDataLocalMap) > 0 {
				request["global"] = objectDataLocalMap
			}
		}
	}

	if d.HasChange("inputs") {
		if v, ok := d.GetOk("inputs"); ok {
			inputs := v.(string)
			inputsList, err := convertJsonStringToList(inputs)
			if err != nil {
				return WrapError(err)
			}
			request["inputs"] = inputsList
		}
	}

	if d.HasChange("processors") {
		if v, ok := d.GetOk("processors"); ok {
			processors := v.(string)
			processorsList, err := convertJsonStringToList(processors)
			if err != nil {
				return WrapError(err)
			}
			request["processors"] = processorsList
		}
	}

	if d.HasChange("aggregators") {
		if v, ok := d.GetOk("aggregators"); ok {
			aggregators := v.(string)
			aggregatorsList, err := convertJsonStringToList(aggregators)
			if err != nil {
				return WrapError(err)
			}
			request["aggregators"] = aggregatorsList
		}
	}

	if d.HasChange("flushers") {
		if v, ok := d.GetOk("flushers"); ok {
			flushers := v.(string)
			flushersList, err := convertJsonStringToList(flushers)
			if err != nil {
				return WrapError(err)
			}
			request["flushers"] = flushersList
		}
	}

	// ConfigName must be set in body for update
	request["configName"] = configName

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateLogtailPipelineConfig", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudSlsLogtailPipelineConfigRead(d, meta)
}

func resourceAliCloudSlsLogtailPipelineConfigDelete(d *schema.ResourceData, meta interface{}) error {

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
