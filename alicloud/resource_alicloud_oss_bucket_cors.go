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
						"expose_header": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_headers": {
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

	if v, ok := d.GetOkExists("response_vary"); ok {
		objectDataLocalMap["ResponseVary"] = v
	}

	if v := d.Get("cors_rule"); !IsNil(v) {
		if v, ok := d.GetOk("cors_rule"); ok {
			localData, err := jsonpath.Get("$", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.(*schema.Set).List() {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["AllowedOrigin"] = dataLoopTmp["allowed_origins"].(*schema.Set).List()
				dataLoopMap["AllowedMethod"] = dataLoopTmp["allowed_methods"].(*schema.Set).List()
				dataLoopMap["AllowedHeader"] = dataLoopTmp["allowed_headers"].(*schema.Set).List()
				dataLoopMap["MaxAgeSeconds"] = dataLoopTmp["max_age_seconds"]
				dataLoopMap["ExposeHeader"] = dataLoopTmp["expose_header"].(*schema.Set).List()
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["CORSRule"] = localMaps
		}

	}

	request["CORSConfiguration"] = objectDataLocalMap

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketCors", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_cors", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	ossServiceV2 := OssServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketCorsStateRefreshFunc(d.Id(), "#CORSRule", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

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

	cORSRuleRaw := objectRaw["CORSRule"]
	corsRuleMaps := make([]map[string]interface{}, 0)
	if cORSRuleRaw != nil {
		for _, cORSRuleChildRaw := range cORSRuleRaw.([]interface{}) {
			corsRuleMap := make(map[string]interface{})
			cORSRuleChildRaw := cORSRuleChildRaw.(map[string]interface{})
			corsRuleMap["max_age_seconds"] = cORSRuleChildRaw["MaxAgeSeconds"]

			allowedHeaderRaw := make([]interface{}, 0)
			if cORSRuleChildRaw["AllowedHeader"] != nil {
				allowedHeaderRaw = cORSRuleChildRaw["AllowedHeader"].([]interface{})
			}

			corsRuleMap["allowed_headers"] = allowedHeaderRaw
			allowedMethodRaw := make([]interface{}, 0)
			if cORSRuleChildRaw["AllowedMethod"] != nil {
				allowedMethodRaw = cORSRuleChildRaw["AllowedMethod"].([]interface{})
			}

			corsRuleMap["allowed_methods"] = allowedMethodRaw
			allowedOriginRaw := make([]interface{}, 0)
			if cORSRuleChildRaw["AllowedOrigin"] != nil {
				allowedOriginRaw = cORSRuleChildRaw["AllowedOrigin"].([]interface{})
			}

			corsRuleMap["allowed_origins"] = allowedOriginRaw
			exposeHeaderRaw := make([]interface{}, 0)
			if cORSRuleChildRaw["ExposeHeader"] != nil {
				exposeHeaderRaw = cORSRuleChildRaw["ExposeHeader"].([]interface{})
			}

			corsRuleMap["expose_header"] = exposeHeaderRaw
			corsRuleMaps = append(corsRuleMaps, corsRuleMap)
		}
	}
	if err := d.Set("cors_rule", corsRuleMaps); err != nil {
		return err
	}

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

	var err error
	action := fmt.Sprintf("/?cors")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())

	objectDataLocalMap := make(map[string]interface{})

	if d.HasChange("response_vary") {
		update = true
	}
	if v, ok := d.GetOk("response_vary"); ok {
		objectDataLocalMap["ResponseVary"] = v
	}

	if d.HasChange("cors_rule") {
		update = true
	}
	if v := d.Get("cors_rule"); v != nil {
		if v, ok := d.GetOk("cors_rule"); ok {
			localData, err := jsonpath.Get("$", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.(*schema.Set).List() {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["MaxAgeSeconds"] = dataLoopTmp["max_age_seconds"]
				dataLoopMap["AllowedHeader"] = dataLoopTmp["allowed_headers"].(*schema.Set).List()
				dataLoopMap["AllowedMethod"] = dataLoopTmp["allowed_methods"].(*schema.Set).List()
				dataLoopMap["AllowedOrigin"] = dataLoopTmp["allowed_origins"].(*schema.Set).List()
				dataLoopMap["ExposeHeader"] = dataLoopTmp["expose_header"].(*schema.Set).List()
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
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketCors", action), query, body, nil, hostMap, false)
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
		ossServiceV2 := OssServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ossServiceV2.OssBucketCorsStateRefreshFunc(d.Id(), "#CORSRule", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Id())

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketCors", action), query, nil, nil, hostMap, false)

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
