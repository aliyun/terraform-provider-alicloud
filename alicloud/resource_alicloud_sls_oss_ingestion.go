// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsOssIngestion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsOssIngestionCreate,
		Read:   resourceAliCloudSlsOssIngestionRead,
		Update: resourceAliCloudSlsOssIngestionUpdate,
		Delete: resourceAliCloudSlsOssIngestionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logstore": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pattern": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"time_format": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"end_time": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"start_time": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"use_meta_index": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"encoding": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"prefix": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"role_arn": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"time_field": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"time_zone": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"format": {
										Type:     schema.TypeMap,
										Optional: true,
									},
									"endpoint": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"bucket": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"time_pattern": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"compression_codec": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"restore_object_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"interval": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oss_ingestion_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"run_immediately": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"cron_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"interval": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudSlsOssIngestionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/ossingestions")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	conn, err := client.NewSlsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project").(string))
	request["name"] = d.Get("oss_ingestion_name")

	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("configuration"); v != nil {
		source := make(map[string]interface{})
		nodeNative, _ := jsonpath.Get("$[0].source[0].endpoint", d.Get("configuration"))
		if nodeNative != nil && nodeNative != "" {
			source["endpoint"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].source[0].time_format", d.Get("configuration"))
		if nodeNative1 != nil && nodeNative1 != "" {
			source["timeFormat"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].source[0].time_pattern", d.Get("configuration"))
		if nodeNative2 != nil && nodeNative2 != "" {
			source["timePattern"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].source[0].interval", d.Get("configuration"))
		if nodeNative3 != nil && nodeNative3 != "" {
			source["interval"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].source[0].end_time", d.Get("configuration"))
		if nodeNative4 != nil && nodeNative4 != "" {
			source["endTime"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].source[0].compression_codec", d.Get("configuration"))
		if nodeNative5 != nil && nodeNative5 != "" {
			source["compressionCodec"] = nodeNative5
		}
		nodeNative6, _ := jsonpath.Get("$[0].source[0].bucket", d.Get("configuration"))
		if nodeNative6 != nil && nodeNative6 != "" {
			source["bucket"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].source[0].pattern", d.Get("configuration"))
		if nodeNative7 != nil && nodeNative7 != "" {
			source["pattern"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].source[0].time_field", d.Get("configuration"))
		if nodeNative8 != nil && nodeNative8 != "" {
			source["timeField"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].source[0].restore_object_enabled", d.Get("configuration"))
		if nodeNative9 != nil && nodeNative9 != "" {
			source["restoreObjectEnabled"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].source[0].use_meta_index", d.Get("configuration"))
		if nodeNative10 != nil && nodeNative10 != "" {
			source["useMetaIndex"] = nodeNative10
		}
		nodeNative11, _ := jsonpath.Get("$[0].source[0].start_time", d.Get("configuration"))
		if nodeNative11 != nil && nodeNative11 != "" {
			source["startTime"] = nodeNative11
		}
		nodeNative12, _ := jsonpath.Get("$[0].source[0].prefix", d.Get("configuration"))
		if nodeNative12 != nil && nodeNative12 != "" {
			source["prefix"] = nodeNative12
		}
		nodeNative13, _ := jsonpath.Get("$[0].source[0].encoding", d.Get("configuration"))
		if nodeNative13 != nil && nodeNative13 != "" {
			source["encoding"] = nodeNative13
		}
		nodeNative14, _ := jsonpath.Get("$[0].source[0].time_zone", d.Get("configuration"))
		if nodeNative14 != nil && nodeNative14 != "" {
			source["timeZone"] = nodeNative14
		}
		nodeNative15, _ := jsonpath.Get("$[0].source[0].role_arn", d.Get("configuration"))
		if nodeNative15 != nil && nodeNative15 != "" {
			source["roleARN"] = nodeNative15
		}
		nodeNative16, _ := jsonpath.Get("$[0].source[0].format", d.Get("configuration"))
		if nodeNative16 != nil && nodeNative16 != "" {
			source["format"] = nodeNative16
		}
		objectDataLocalMap["source"] = source
		nodeNative17, _ := jsonpath.Get("$[0].logstore", d.Get("configuration"))
		if nodeNative17 != nil && nodeNative17 != "" {
			objectDataLocalMap["logstore"] = nodeNative17
		}
		request["configuration"] = objectDataLocalMap
	}

	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("schedule"); !IsNil(v) {
		nodeNative18, _ := jsonpath.Get("$[0].cron_expression", d.Get("schedule"))
		if nodeNative18 != nil && nodeNative18 != "" {
			objectDataLocalMap1["cronExpression"] = nodeNative18
		}
		nodeNative19, _ := jsonpath.Get("$[0].delay", d.Get("schedule"))
		if nodeNative19 != nil && nodeNative19 != "" {
			objectDataLocalMap1["delay"] = nodeNative19
		}
		nodeNative20, _ := jsonpath.Get("$[0].time_zone", d.Get("schedule"))
		if nodeNative20 != nil && nodeNative20 != "" {
			objectDataLocalMap1["timeZone"] = nodeNative20
		}
		nodeNative21, _ := jsonpath.Get("$[0].type", d.Get("schedule"))
		if nodeNative21 != nil && nodeNative21 != "" {
			objectDataLocalMap1["type"] = nodeNative21
		}
		nodeNative22, _ := jsonpath.Get("$[0].interval", d.Get("schedule"))
		if nodeNative22 != nil && nodeNative22 != "" {
			objectDataLocalMap1["interval"] = nodeNative22
		}
		nodeNative23, _ := jsonpath.Get("$[0].run_immediately", d.Get("schedule"))
		if nodeNative23 != nil && nodeNative23 != "" {
			objectDataLocalMap1["runImmediately"] = nodeNative23
		}
		request["schedule"] = objectDataLocalMap1
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["displayName"] = d.Get("display_name")
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.Execute(genRoaParam("CreateOssIngestion", "POST", "2020-12-30", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_oss_ingestion", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["name"]))

	return resourceAliCloudSlsOssIngestionRead(d, meta)
}

func resourceAliCloudSlsOssIngestionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsOssIngestion(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_oss_ingestion DescribeSlsOssIngestion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("display_name", objectRaw["displayName"])
	d.Set("status", objectRaw["status"])
	d.Set("oss_ingestion_name", objectRaw["name"])

	configurationMaps := make([]map[string]interface{}, 0)
	configurationMap := make(map[string]interface{})
	configuration1Raw := make(map[string]interface{})
	if objectRaw["configuration"] != nil {
		configuration1Raw = objectRaw["configuration"].(map[string]interface{})
	}
	if len(configuration1Raw) > 0 {
		configurationMap["logstore"] = configuration1Raw["logstore"]

		sourceMaps := make([]map[string]interface{}, 0)
		sourceMap := make(map[string]interface{})
		source1Raw := make(map[string]interface{})
		if configuration1Raw["source"] != nil {
			source1Raw = configuration1Raw["source"].(map[string]interface{})
		}
		if len(source1Raw) > 0 {
			sourceMap["bucket"] = source1Raw["bucket"]
			sourceMap["compression_codec"] = source1Raw["compressionCodec"]
			sourceMap["encoding"] = source1Raw["encoding"]
			sourceMap["end_time"] = source1Raw["endTime"]
			sourceMap["endpoint"] = source1Raw["endpoint"]
			sourceMap["format"] = source1Raw["format"]
			sourceMap["interval"] = source1Raw["interval"]
			sourceMap["pattern"] = source1Raw["pattern"]
			sourceMap["prefix"] = source1Raw["prefix"]
			sourceMap["restore_object_enabled"] = source1Raw["restoreObjectEnabled"]
			sourceMap["role_arn"] = source1Raw["roleARN"]
			sourceMap["start_time"] = source1Raw["startTime"]
			sourceMap["time_field"] = source1Raw["timeField"]
			sourceMap["time_format"] = source1Raw["timeFormat"]
			sourceMap["time_pattern"] = source1Raw["timePattern"]
			sourceMap["time_zone"] = source1Raw["timeZone"]
			sourceMap["use_meta_index"] = source1Raw["useMetaIndex"]

			sourceMaps = append(sourceMaps, sourceMap)
		}
		configurationMap["source"] = sourceMaps
		configurationMaps = append(configurationMaps, configurationMap)
	}
	d.Set("configuration", configurationMaps)
	scheduleMaps := make([]map[string]interface{}, 0)
	scheduleMap := make(map[string]interface{})
	schedule1Raw := make(map[string]interface{})
	if objectRaw["schedule"] != nil {
		schedule1Raw = objectRaw["schedule"].(map[string]interface{})
	}
	if len(schedule1Raw) > 0 {
		scheduleMap["cron_expression"] = schedule1Raw["cronExpression"]
		scheduleMap["delay"] = schedule1Raw["delay"]
		scheduleMap["interval"] = schedule1Raw["interval"]
		scheduleMap["run_immediately"] = schedule1Raw["runImmediately"]
		scheduleMap["time_zone"] = schedule1Raw["timeZone"]
		scheduleMap["type"] = schedule1Raw["type"]

		scheduleMaps = append(scheduleMaps, scheduleMap)
	}
	d.Set("schedule", scheduleMaps)

	parts := strings.Split(d.Id(), ":")
	d.Set("project", parts[0])
	d.Set("oss_ingestion_name", parts[1])

	return nil
}

func resourceAliCloudSlsOssIngestionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	var hostMap map[string]*string
	update := false
	parts := strings.Split(d.Id(), ":")
	ossIngestionName := parts[1]
	action := fmt.Sprintf("/ossingestions/%s", ossIngestionName)
	conn, err := client.NewSlsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap = make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])
	request["OssIngestionName"] = parts[1]
	if d.HasChange("configuration") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("configuration"); v != nil {
		source := make(map[string]interface{})
		nodeNative, _ := jsonpath.Get("$[0].source[0].endpoint", v)
		if nodeNative != nil && nodeNative != "" {
			source["endpoint"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].source[0].time_format", v)
		if nodeNative1 != nil && nodeNative1 != "" {
			source["timeFormat"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].source[0].time_pattern", v)
		if nodeNative2 != nil && nodeNative2 != "" {
			source["timePattern"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].source[0].interval", v)
		if nodeNative3 != nil && nodeNative3 != "" {
			source["interval"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].source[0].end_time", v)
		if nodeNative4 != nil && nodeNative4 != "" {
			source["endTime"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].source[0].compression_codec", v)
		if nodeNative5 != nil && nodeNative5 != "" {
			source["compressionCodec"] = nodeNative5
		}
		nodeNative6, _ := jsonpath.Get("$[0].source[0].bucket", v)
		if nodeNative6 != nil && nodeNative6 != "" {
			source["bucket"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].source[0].pattern", v)
		if nodeNative7 != nil && nodeNative7 != "" {
			source["pattern"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].source[0].time_field", v)
		if nodeNative8 != nil && nodeNative8 != "" {
			source["timeField"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].source[0].restore_object_enabled", v)
		if nodeNative9 != nil && nodeNative9 != "" {
			source["restoreObjectEnabled"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].source[0].use_meta_index", v)
		if nodeNative10 != nil && nodeNative10 != "" {
			source["useMetaIndex"] = nodeNative10
		}
		nodeNative11, _ := jsonpath.Get("$[0].source[0].start_time", v)
		if nodeNative11 != nil && nodeNative11 != "" {
			source["startTime"] = nodeNative11
		}
		nodeNative12, _ := jsonpath.Get("$[0].source[0].prefix", v)
		if nodeNative12 != nil && nodeNative12 != "" {
			source["prefix"] = nodeNative12
		}
		nodeNative13, _ := jsonpath.Get("$[0].source[0].encoding", v)
		if nodeNative13 != nil && nodeNative13 != "" {
			source["encoding"] = nodeNative13
		}
		nodeNative14, _ := jsonpath.Get("$[0].source[0].time_zone", v)
		if nodeNative14 != nil && nodeNative14 != "" {
			source["timeZone"] = nodeNative14
		}
		nodeNative15, _ := jsonpath.Get("$[0].source[0].role_arn", v)
		if nodeNative15 != nil && nodeNative15 != "" {
			source["roleARN"] = nodeNative15
		}
		nodeNative16, _ := jsonpath.Get("$[0].source[0].format", v)
		if nodeNative16 != nil && nodeNative16 != "" {
			source["format"] = nodeNative16
		}
		objectDataLocalMap["source"] = source
		nodeNative17, _ := jsonpath.Get("$[0].logstore", v)
		if nodeNative17 != nil && nodeNative17 != "" {
			objectDataLocalMap["logstore"] = nodeNative17
		}
		request["configuration"] = objectDataLocalMap
	}

	if d.HasChange("schedule") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("schedule"); v != nil {
		nodeNative18, _ := jsonpath.Get("$[0].cron_expression", v)
		if nodeNative18 != nil && nodeNative18 != "" {
			objectDataLocalMap1["cronExpression"] = nodeNative18
		}
		nodeNative19, _ := jsonpath.Get("$[0].delay", v)
		if nodeNative19 != nil && nodeNative19 != "" {
			objectDataLocalMap1["delay"] = nodeNative19
		}
		nodeNative20, _ := jsonpath.Get("$[0].time_zone", v)
		if nodeNative20 != nil && nodeNative20 != "" {
			objectDataLocalMap1["timeZone"] = nodeNative20
		}
		nodeNative21, _ := jsonpath.Get("$[0].type", v)
		if nodeNative21 != nil && nodeNative21 != "" {
			objectDataLocalMap1["type"] = nodeNative21
		}
		nodeNative22, _ := jsonpath.Get("$[0].interval", v)
		if nodeNative22 != nil && nodeNative22 != "" {
			objectDataLocalMap1["interval"] = nodeNative22
		}
		nodeNative23, _ := jsonpath.Get("$[0].run_immediately", v)
		if nodeNative23 != nil && nodeNative23 != "" {
			objectDataLocalMap1["runImmediately"] = nodeNative23
		}
		request["schedule"] = objectDataLocalMap1
	}

	if d.HasChange("description") {
		update = true
	}
	request["description"] = d.Get("description")
	if d.HasChange("display_name") {
		update = true
	}
	request["displayName"] = d.Get("display_name")
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.Execute(genRoaParam("UpdateOssIngestion", "PUT", "2020-12-30", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudSlsOssIngestionRead(d, meta)
}

func resourceAliCloudSlsOssIngestionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	ossIngestionName := parts[1]
	action := fmt.Sprintf("/ossingestions/%s", ossIngestionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	conn, err := client.NewSlsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])
	request["OssIngestionName"] = parts[1]

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.Execute(genRoaParam("DeleteOssIngestion", "DELETE", "2020-12-30", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
