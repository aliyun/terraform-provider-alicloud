// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudSlsOssExportSink() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsOssExportSinkCreate,
		Read:   resourceAliCloudSlsOssExportSinkRead,
		Update: resourceAliCloudSlsOssExportSinkUpdate,
		Delete: resourceAliCloudSlsOssExportSinkDelete,
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
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"to_time": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"sink": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"buffer_interval": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"path_format": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"content_detail": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.ValidateJsonString,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											equal, _ := compareJsonTemplateAreEquivalent(old, new)
											return equal
										},
									},
									"prefix": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"path_format_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"time"}, false),
									},
									"role_arn": {
										Type:     schema.TypeString,
										Required: true,
									},
									"buffer_size": {
										Type:     schema.TypeString,
										Required: true,
									},
									"time_zone": {
										Type:     schema.TypeString,
										Required: true,
									},
									"suffix": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"endpoint": {
										Type:     schema.TypeString,
										Required: true,
									},
									"delay_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"bucket": {
										Type:     schema.TypeString,
										Required: true,
									},
									"compression_type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"logstore": {
							Type:     schema.TypeString,
							Required: true,
						},
						"from_time": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
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
			"job_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudSlsOssExportSinkCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/ossexports")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project").(string))
	request["name"] = d.Get("job_name")

	request["displayName"] = d.Get("display_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("configuration"); v != nil {
		logstore1, _ := jsonpath.Get("$[0].logstore", d.Get("configuration"))
		if logstore1 != nil && logstore1 != "" {
			objectDataLocalMap["logstore"] = logstore1
		}
		sink := make(map[string]interface{})
		endpoint1, _ := jsonpath.Get("$[0].sink[0].endpoint", d.Get("configuration"))
		if endpoint1 != nil && endpoint1 != "" {
			sink["endpoint"] = endpoint1
		}
		bucket1, _ := jsonpath.Get("$[0].sink[0].bucket", d.Get("configuration"))
		if bucket1 != nil && bucket1 != "" {
			sink["bucket"] = bucket1
		}
		prefix1, _ := jsonpath.Get("$[0].sink[0].prefix", d.Get("configuration"))
		if prefix1 != nil && prefix1 != "" {
			sink["prefix"] = prefix1
		}
		suffix1, _ := jsonpath.Get("$[0].sink[0].suffix", d.Get("configuration"))
		if suffix1 != nil && suffix1 != "" {
			sink["suffix"] = suffix1
		}
		pathFormat1, _ := jsonpath.Get("$[0].sink[0].path_format", d.Get("configuration"))
		if pathFormat1 != nil && pathFormat1 != "" {
			sink["pathFormat"] = pathFormat1
		}
		pathFormatType1, _ := jsonpath.Get("$[0].sink[0].path_format_type", d.Get("configuration"))
		if pathFormatType1 != nil && pathFormatType1 != "" {
			sink["pathFormatType"] = pathFormatType1
		}
		timeZone1, _ := jsonpath.Get("$[0].sink[0].time_zone", d.Get("configuration"))
		if timeZone1 != nil && timeZone1 != "" {
			sink["timeZone"] = timeZone1
		}
		contentType1, _ := jsonpath.Get("$[0].sink[0].content_type", d.Get("configuration"))
		if contentType1 != nil && contentType1 != "" {
			sink["contentType"] = contentType1
		}
		compressionType1, _ := jsonpath.Get("$[0].sink[0].compression_type", d.Get("configuration"))
		if compressionType1 != nil && compressionType1 != "" {
			sink["compressionType"] = compressionType1
		}
		bufferInterval1, _ := jsonpath.Get("$[0].sink[0].buffer_interval", d.Get("configuration"))
		if bufferInterval1 != nil && bufferInterval1 != "" {
			sink["bufferInterval"] = bufferInterval1
		}
		bufferSize1, _ := jsonpath.Get("$[0].sink[0].buffer_size", d.Get("configuration"))
		if bufferSize1 != nil && bufferSize1 != "" {
			sink["bufferSize"] = bufferSize1
		}
		delaySeconds1, _ := jsonpath.Get("$[0].sink[0].delay_seconds", d.Get("configuration"))
		if delaySeconds1 != nil && delaySeconds1 != "" {
			sink["delaySeconds"] = delaySeconds1
		}
		roleArn1, _ := jsonpath.Get("$[0].sink[0].role_arn", d.Get("configuration"))
		if roleArn1 != nil && roleArn1 != "" {
			sink["roleArn"] = roleArn1
		}
		contentDetail1, _ := jsonpath.Get("$[0].sink[0].content_detail", d.Get("configuration"))
		if contentDetail1 != nil && contentDetail1 != "" {
			sink["contentDetail"] = contentDetail1
		}

		objectDataLocalMap["sink"] = sink
		fromTime1, _ := jsonpath.Get("$[0].from_time", d.Get("configuration"))
		if fromTime1 != nil && fromTime1 != "" {
			objectDataLocalMap["fromTime"] = fromTime1
		}
		toTime1, _ := jsonpath.Get("$[0].to_time", d.Get("configuration"))
		if toTime1 != nil && toTime1 != "" {
			objectDataLocalMap["toTime"] = toTime1
		}
		roleArn3, _ := jsonpath.Get("$[0].role_arn", d.Get("configuration"))
		if roleArn3 != nil && roleArn3 != "" {
			objectDataLocalMap["roleArn"] = roleArn3
		}

		request["configuration"] = objectDataLocalMap
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateOSSExport", action), query, body, nil, hostMap, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"403"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_oss_export_sink", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["name"]))

	slsServiceV2 := SlsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, slsServiceV2.SlsOssExportSinkStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudSlsOssExportSinkRead(d, meta)
}

func resourceAliCloudSlsOssExportSinkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsOssExportSink(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_oss_export_sink DescribeSlsOssExportSink Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["createTime"] != nil {
		d.Set("create_time", objectRaw["createTime"])
	}
	if objectRaw["description"] != nil {
		d.Set("description", objectRaw["description"])
	}
	if objectRaw["displayName"] != nil {
		d.Set("display_name", objectRaw["displayName"])
	}
	if objectRaw["status"] != nil {
		d.Set("status", objectRaw["status"])
	}

	configurationMaps := make([]map[string]interface{}, 0)
	configurationMap := make(map[string]interface{})
	configuration1Raw := make(map[string]interface{})
	if objectRaw["configuration"] != nil {
		configuration1Raw = objectRaw["configuration"].(map[string]interface{})
	}
	if len(configuration1Raw) > 0 {
		configurationMap["from_time"] = configuration1Raw["fromTime"]
		configurationMap["logstore"] = configuration1Raw["logstore"]
		configurationMap["role_arn"] = configuration1Raw["roleArn"]
		configurationMap["to_time"] = configuration1Raw["toTime"]

		sinkMaps := make([]map[string]interface{}, 0)
		sinkMap := make(map[string]interface{})
		sink1Raw := make(map[string]interface{})
		if configuration1Raw["sink"] != nil {
			sink1Raw = configuration1Raw["sink"].(map[string]interface{})
		}
		if len(sink1Raw) > 0 {
			sinkMap["bucket"] = sink1Raw["bucket"]
			sinkMap["buffer_interval"] = sink1Raw["bufferInterval"]
			sinkMap["buffer_size"] = sink1Raw["bufferSize"]
			sinkMap["compression_type"] = sink1Raw["compressionType"]
			sinkMap["content_detail"] = convertObjectToJsonString(sink1Raw["contentDetail"])
			sinkMap["content_type"] = sink1Raw["contentType"]
			sinkMap["delay_seconds"] = sink1Raw["delaySeconds"]
			sinkMap["endpoint"] = sink1Raw["endpoint"]
			sinkMap["path_format"] = sink1Raw["pathFormat"]
			sinkMap["path_format_type"] = sink1Raw["pathFormatType"]
			sinkMap["prefix"] = sink1Raw["prefix"]
			sinkMap["role_arn"] = sink1Raw["roleArn"]
			sinkMap["suffix"] = sink1Raw["suffix"]
			sinkMap["time_zone"] = sink1Raw["timeZone"]

			sinkMaps = append(sinkMaps, sinkMap)
		}
		configurationMap["sink"] = sinkMaps
		configurationMaps = append(configurationMaps, configurationMap)
	}
	if objectRaw["configuration"] != nil {
		if err := d.Set("configuration", configurationMaps); err != nil {
			return err
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("project", parts[0])
	d.Set("job_name", parts[1])

	return nil
}

func resourceAliCloudSlsOssExportSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	parts := strings.Split(d.Id(), ":")
	ossExportName := parts[1]
	action := fmt.Sprintf("/ossexports/%s", ossExportName)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])

	if d.HasChange("display_name") {
		update = true
	}
	request["displayName"] = d.Get("display_name")
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("configuration") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("configuration"); v != nil {
		logstore1, _ := jsonpath.Get("$[0].logstore", v)
		if logstore1 != nil && (d.HasChange("configuration.0.logstore") || logstore1 != "") {
			objectDataLocalMap["logstore"] = logstore1
		}
		sink := make(map[string]interface{})
		endpoint1, _ := jsonpath.Get("$[0].sink[0].endpoint", v)
		if endpoint1 != nil && (d.HasChange("configuration.0.sink.0.endpoint") || endpoint1 != "") {
			sink["endpoint"] = endpoint1
		}
		bucket1, _ := jsonpath.Get("$[0].sink[0].bucket", v)
		if bucket1 != nil && (d.HasChange("configuration.0.sink.0.bucket") || bucket1 != "") {
			sink["bucket"] = bucket1
		}
		prefix1, _ := jsonpath.Get("$[0].sink[0].prefix", v)
		if prefix1 != nil && (d.HasChange("configuration.0.sink.0.prefix") || prefix1 != "") {
			sink["prefix"] = prefix1
		}
		suffix1, _ := jsonpath.Get("$[0].sink[0].suffix", v)
		if suffix1 != nil && (d.HasChange("configuration.0.sink.0.suffix") || suffix1 != "") {
			sink["suffix"] = suffix1
		}
		pathFormat1, _ := jsonpath.Get("$[0].sink[0].path_format", v)
		if pathFormat1 != nil && (d.HasChange("configuration.0.sink.0.path_format") || pathFormat1 != "") {
			sink["pathFormat"] = pathFormat1
		}
		pathFormatType1, _ := jsonpath.Get("$[0].sink[0].path_format_type", v)
		if pathFormatType1 != nil && (d.HasChange("configuration.0.sink.0.path_format_type") || pathFormatType1 != "") {
			sink["pathFormatType"] = pathFormatType1
		}
		timeZone1, _ := jsonpath.Get("$[0].sink[0].time_zone", v)
		if timeZone1 != nil && (d.HasChange("configuration.0.sink.0.time_zone") || timeZone1 != "") {
			sink["timeZone"] = timeZone1
		}
		contentType1, _ := jsonpath.Get("$[0].sink[0].content_type", v)
		if contentType1 != nil && (d.HasChange("configuration.0.sink.0.content_type") || contentType1 != "") {
			sink["contentType"] = contentType1
		}
		compressionType1, _ := jsonpath.Get("$[0].sink[0].compression_type", v)
		if compressionType1 != nil && (d.HasChange("configuration.0.sink.0.compression_type") || compressionType1 != "") {
			sink["compressionType"] = compressionType1
		}
		bufferInterval1, _ := jsonpath.Get("$[0].sink[0].buffer_interval", v)
		if bufferInterval1 != nil && (d.HasChange("configuration.0.sink.0.buffer_interval") || bufferInterval1 != "") {
			sink["bufferInterval"] = bufferInterval1
		}
		bufferSize1, _ := jsonpath.Get("$[0].sink[0].buffer_size", v)
		if bufferSize1 != nil && (d.HasChange("configuration.0.sink.0.buffer_size") || bufferSize1 != "") {
			sink["bufferSize"] = bufferSize1
		}
		delaySeconds1, _ := jsonpath.Get("$[0].sink[0].delay_seconds", v)
		if delaySeconds1 != nil && (d.HasChange("configuration.0.sink.0.delay_seconds") || delaySeconds1 != "") {
			sink["delaySeconds"] = delaySeconds1
		}
		roleArn1, _ := jsonpath.Get("$[0].sink[0].role_arn", v)
		if roleArn1 != nil && (d.HasChange("configuration.0.sink.0.role_arn") || roleArn1 != "") {
			sink["roleArn"] = roleArn1
		}
		contentDetail1, _ := jsonpath.Get("$[0].sink[0].content_detail", v)
		if contentDetail1 != nil && (d.HasChange("configuration.0.sink.0.content_detail") || contentDetail1 != "") {
			sink["contentDetail"] = contentDetail1
		}

		objectDataLocalMap["sink"] = sink
		fromTime1, _ := jsonpath.Get("$[0].from_time", v)
		if fromTime1 != nil && (d.HasChange("configuration.0.from_time") || fromTime1 != "") {
			objectDataLocalMap["fromTime"] = fromTime1
		}
		toTime1, _ := jsonpath.Get("$[0].to_time", v)
		if toTime1 != nil && (d.HasChange("configuration.0.to_time") || toTime1 != "") {
			objectDataLocalMap["toTime"] = toTime1
		}
		roleArn3, _ := jsonpath.Get("$[0].role_arn", v)
		if roleArn3 != nil && (d.HasChange("configuration.0.role_arn") || roleArn3 != "") {
			objectDataLocalMap["roleArn"] = roleArn3
		}

		request["configuration"] = objectDataLocalMap
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateOSSExport", action), query, body, nil, hostMap, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"403"}) || NeedRetry(err) {
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
		slsServiceV2 := SlsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, slsServiceV2.SlsOssExportSinkStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudSlsOssExportSinkRead(d, meta)
}

func resourceAliCloudSlsOssExportSinkDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	ossExportName := parts[1]
	action := fmt.Sprintf("/ossexports/%s", ossExportName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteOSSExport", action), query, nil, nil, hostMap, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"403"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
