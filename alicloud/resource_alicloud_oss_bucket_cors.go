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

func resourceAliCloudOssBucketCors() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketCorsCreate,
		Read:   resourceAliCloudOssBucketCorsRead,
		Update: resourceAliCloudOssBucketCorsUpdate,
		Delete: resourceAliCloudOssBucketCorsDelete,
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
			"cors_rule": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_methods": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_origins": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_age_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"allowed_headers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expose_header": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"response_vary": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudOssBucketCorsCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?cors")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOk("response_vary"); ok {
		objectDataLocalMap["ResponseVary"] = v
	}

	if v := d.Get("cors_rule"); !IsNil(v) {
		if v, ok := d.GetOk("cors_rule"); ok {
			localData, err := jsonpath.Get("$", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.(*schema.Set).List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["AllowedOrigin"] = dataLoopTmp["allowed_origins"].(*schema.Set).List()
				dataLoopMap["AllowedMethod"] = dataLoopTmp["allowed_methods"].(*schema.Set).List()
				dataLoopMap["AllowedHeader"] = dataLoopTmp["allowed_headers"].(*schema.Set).List()
				dataLoopMap["ExposeHeader"] = dataLoopTmp["expose_header"].(*schema.Set).List()
				dataLoopMap["MaxAgeSeconds"] = dataLoopTmp["max_age_seconds"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["CORSRule"] = localMaps
		}
	}

	request["CORSConfiguration"] = objectDataLocalMap
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", genXmlParam("PUT", "2019-05-17", "PutBucketCors", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_cors", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	return resourceAliCloudOssBucketCorsRead(d, meta)
}

func resourceAliCloudOssBucketCorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketCors(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_cors DescribeOssBucketCors Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("response_vary", objectRaw["ResponseVary"])

	cORSRule1Raw := objectRaw["CORSRule"]
	corsRuleMaps := make([]map[string]interface{}, 0)
	if cORSRule1Raw != nil {
		for _, cORSRuleChild1Raw := range cORSRule1Raw.([]interface{}) {
			corsRuleMap := make(map[string]interface{})
			cORSRuleChild1Raw := cORSRuleChild1Raw.(map[string]interface{})
			corsRuleMap["max_age_seconds"] = cORSRuleChild1Raw["MaxAgeSeconds"]

			allowedHeader1Raw := make([]interface{}, 0)
			if cORSRuleChild1Raw["AllowedHeader"] != nil {
				allowedHeader1Raw = cORSRuleChild1Raw["AllowedHeader"].([]interface{})
			}

			corsRuleMap["allowed_headers"] = allowedHeader1Raw
			allowedMethod1Raw := make([]interface{}, 0)
			if cORSRuleChild1Raw["AllowedMethod"] != nil {
				allowedMethod1Raw = cORSRuleChild1Raw["AllowedMethod"].([]interface{})
			}

			corsRuleMap["allowed_methods"] = allowedMethod1Raw
			allowedOrigin1Raw := make([]interface{}, 0)
			if cORSRuleChild1Raw["AllowedOrigin"] != nil {
				allowedOrigin1Raw = cORSRuleChild1Raw["AllowedOrigin"].([]interface{})
			}

			corsRuleMap["allowed_origins"] = allowedOrigin1Raw
			exposeHeader1Raw := make([]interface{}, 0)
			if cORSRuleChild1Raw["ExposeHeader"] != nil {
				exposeHeader1Raw = cORSRuleChild1Raw["ExposeHeader"].([]interface{})
			}

			corsRuleMap["expose_header"] = exposeHeader1Raw
			corsRuleMaps = append(corsRuleMaps, corsRuleMap)
		}
	}
	d.Set("cors_rule", corsRuleMaps)

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketCorsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	action := fmt.Sprintf("/?cors")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())
	objectDataLocalMap := make(map[string]interface{})

	if d.HasChange("response_vary") {
		update = true
	}
	objectDataLocalMap["ResponseVary"] = d.Get("response_vary")
	if d.HasChange("cors_rule") {
		update = true
	}
	if v := d.Get("cors_rule"); v != nil {
		if v, ok := d.GetOk("cors_rule"); ok {
			localData, err := jsonpath.Get("$", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.(*schema.Set).List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["MaxAgeSeconds"] = dataLoopTmp["max_age_seconds"]
				dataLoopMap["ExposeHeader"] = dataLoopTmp["expose_header"].(*schema.Set).List()
				dataLoopMap["AllowedHeader"] = dataLoopTmp["allowed_headers"].(*schema.Set).List()
				dataLoopMap["AllowedMethod"] = dataLoopTmp["allowed_methods"].(*schema.Set).List()
				dataLoopMap["AllowedOrigin"] = dataLoopTmp["allowed_origins"].(*schema.Set).List()
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["CORSRule"] = localMaps
		}
	}

	request["CORSConfiguration"] = objectDataLocalMap
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", genXmlParam("PUT", "2019-05-17", "PutBucketCors", action), query, body, nil, hostMap, false)
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

	return resourceAliCloudOssBucketCorsRead(d, meta)
}

func resourceAliCloudOssBucketCorsDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?cors")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Id())

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", genXmlParam("DELETE", "2019-05-17", "DeleteBucketCors", action), query, body, nil, hostMap, false)
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
