// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketResponseHeader() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketResponseHeaderCreate,
		Read:   resourceAliCloudOssBucketResponseHeaderRead,
		Update: resourceAliCloudOssBucketResponseHeaderUpdate,
		Delete: resourceAliCloudOssBucketResponseHeaderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operation": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"hide_headers": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudOssBucketResponseHeaderCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?responseHeader")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	responseHeaderConfiguration := make(map[string]interface{})

	if v := d.Get("rule"); !IsNil(v) {
		if v, ok := d.GetOk("rule"); ok {
			localData, err := jsonpath.Get("$", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Name"] = dataLoopTmp["name"]
				localData1 := make(map[string]interface{})
				operation1, _ := jsonpath.Get("$[0].operation", dataLoopTmp["filters"])
				if operation1 != nil && operation1 != "" {
					localData1["Operation"] = operation1
				}
				if len(localData1) > 0 {
					dataLoopMap["Filters"] = localData1
				}
				localData2 := make(map[string]interface{})
				header1, _ := jsonpath.Get("$[0].header", dataLoopTmp["hide_headers"])
				if header1 != nil && header1 != "" {
					localData2["Header"] = header1
				}
				if len(localData2) > 0 {
					dataLoopMap["HideHeaders"] = localData2
				}
				localMaps = append(localMaps, dataLoopMap)
			}
			responseHeaderConfiguration["Rule"] = localMaps
		}

		request["ResponseHeaderConfiguration"] = responseHeaderConfiguration
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketResponseHeader", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_response_header", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	return resourceAliCloudOssBucketResponseHeaderRead(d, meta)
}

func resourceAliCloudOssBucketResponseHeaderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketResponseHeader(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_response_header DescribeOssBucketResponseHeader Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	ruleRaw, _ := jsonpath.Get("$.Rule", objectRaw)

	ruleMaps := make([]map[string]interface{}, 0)
	if ruleRaw != nil {
		for _, ruleChildRaw := range convertToInterfaceArray(ruleRaw) {
			ruleMap := make(map[string]interface{})
			ruleChildRaw := ruleChildRaw.(map[string]interface{})
			ruleMap["name"] = ruleChildRaw["Name"]

			filtersMaps := make([]map[string]interface{}, 0)
			filtersMap := make(map[string]interface{})
			operationRaw, _ := jsonpath.Get("$.Filters.Operation", ruleChildRaw)

			filtersMap["operation"] = operationRaw
			filtersMaps = append(filtersMaps, filtersMap)
			ruleMap["filters"] = filtersMaps
			hideHeadersMaps := make([]map[string]interface{}, 0)
			hideHeadersMap := make(map[string]interface{})
			headerRaw, _ := jsonpath.Get("$.HideHeaders.Header", ruleChildRaw)

			hideHeadersMap["header"] = headerRaw
			hideHeadersMaps = append(hideHeadersMaps, hideHeadersMap)
			ruleMap["hide_headers"] = hideHeadersMaps
			ruleMaps = append(ruleMaps, ruleMap)
		}
	}
	if err := d.Set("rule", ruleMaps); err != nil {
		return err
	}

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketResponseHeaderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/?responseHeader")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())

	if d.HasChange("rule") {
		update = true
	}
	responseHeaderConfiguration := make(map[string]interface{})

	if v := d.Get("rule"); !IsNil(v) || d.HasChange("rule") {
		if v, ok := d.GetOk("rule"); ok {
			localData, err := jsonpath.Get("$", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Name"] = dataLoopTmp["name"]
				if !IsNil(dataLoopTmp["filters"]) {
					localData1 := make(map[string]interface{})
					operation1, _ := jsonpath.Get("$[0].operation", dataLoopTmp["filters"])
					if operation1 != nil && operation1 != "" {
						localData1["Operation"] = operation1
					}
					if len(localData1) > 0 {
						dataLoopMap["Filters"] = localData1
					}
				}
				if !IsNil(dataLoopTmp["hide_headers"]) {
					localData2 := make(map[string]interface{})
					header1, _ := jsonpath.Get("$[0].header", dataLoopTmp["hide_headers"])
					if header1 != nil && header1 != "" {
						localData2["Header"] = header1
					}
					if len(localData2) > 0 {
						dataLoopMap["HideHeaders"] = localData2
					}
				}
				localMaps = append(localMaps, dataLoopMap)
			}
			responseHeaderConfiguration["Rule"] = localMaps
		}

		request["ResponseHeaderConfiguration"] = responseHeaderConfiguration
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketResponseHeader", action), query, body, nil, hostMap, false)
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
	}

	return resourceAliCloudOssBucketResponseHeaderRead(d, meta)
}

func resourceAliCloudOssBucketResponseHeaderDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?responseHeader")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Id())

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketResponseHeader", action), query, nil, nil, hostMap, false)
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
		if IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchResponseHeaderConfiguration"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
