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

func resourceAliCloudOssBucketOverwriteConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketOverwriteConfigCreate,
		Read:   resourceAliCloudOssBucketOverwriteConfigRead,
		Update: resourceAliCloudOssBucketOverwriteConfigUpdate,
		Delete: resourceAliCloudOssBucketOverwriteConfigDelete,
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
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"suffix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"principals": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"principal": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudOssBucketOverwriteConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?overwriteConfig")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	overwriteConfiguration := make(map[string]interface{})

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
				localData1 := make(map[string]interface{})
				principal1, _ := jsonpath.Get("$[0].principal", dataLoopTmp["principals"])
				if principal1 != nil && principal1 != "" {
					localData1["Principal"] = principal1
				}
				if len(localData1) > 0 {
					dataLoopMap["Principals"] = localData1
				}
				dataLoopMap["Action"] = dataLoopTmp["action"]
				dataLoopMap["ID"] = dataLoopTmp["id"]
				dataLoopMap["Suffix"] = dataLoopTmp["suffix"]
				dataLoopMap["Prefix"] = dataLoopTmp["prefix"]
				localMaps = append(localMaps, dataLoopMap)
			}
			overwriteConfiguration["Rule"] = localMaps
		}

		request["OverwriteConfiguration"] = overwriteConfiguration
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketOverwriteConfig", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_overwrite_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	ossServiceV2 := OssServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketOverwriteConfigStateRefreshFunc(d.Id(), "#$.Rule[0].ID", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudOssBucketOverwriteConfigRead(d, meta)
}

func resourceAliCloudOssBucketOverwriteConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketOverwriteConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_overwrite_config DescribeOssBucketOverwriteConfig Failed!!! %s", err)
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
			ruleMap["action"] = ruleChildRaw["Action"]
			ruleMap["id"] = ruleChildRaw["ID"]
			ruleMap["prefix"] = ruleChildRaw["Prefix"]
			ruleMap["suffix"] = ruleChildRaw["Suffix"]

			principalsMaps := make([]map[string]interface{}, 0)
			principalsMap := make(map[string]interface{})
			principalRaw, _ := jsonpath.Get("$.Principals.Principal", ruleChildRaw)

			principalsMap["principal"] = principalRaw
			principalsMaps = append(principalsMaps, principalsMap)
			ruleMap["principals"] = principalsMaps
			ruleMaps = append(ruleMaps, ruleMap)
		}
	}
	if err := d.Set("rule", ruleMaps); err != nil {
		return err
	}

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketOverwriteConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/?overwriteConfig")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())

	if d.HasChange("rule") {
		update = true
	}
	overwriteConfiguration := make(map[string]interface{})

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
				if !IsNil(dataLoopTmp["principals"]) {
					localData1 := make(map[string]interface{})
					principal1, _ := jsonpath.Get("$[0].principal", dataLoopTmp["principals"])
					if principal1 != nil && principal1 != "" {
						localData1["Principal"] = principal1
					}
					if len(localData1) > 0 {
						dataLoopMap["Principals"] = localData1
					}
				}
				dataLoopMap["Action"] = dataLoopTmp["action"]
				dataLoopMap["ID"] = dataLoopTmp["id"]
				dataLoopMap["Suffix"] = dataLoopTmp["suffix"]
				dataLoopMap["Prefix"] = dataLoopTmp["prefix"]
				localMaps = append(localMaps, dataLoopMap)
			}
			overwriteConfiguration["Rule"] = localMaps
		}

		request["OverwriteConfiguration"] = overwriteConfiguration
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketOverwriteConfig", action), query, body, nil, hostMap, false)
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
		stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ossServiceV2.OssBucketOverwriteConfigStateRefreshFunc(d.Id(), "#$.Rule[0].ID", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudOssBucketOverwriteConfigRead(d, meta)
}

func resourceAliCloudOssBucketOverwriteConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?overwriteConfig")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Id())

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketOverwriteConfig", action), query, nil, nil, hostMap, false)
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
		if IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchBucketOverwriteConfig"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
