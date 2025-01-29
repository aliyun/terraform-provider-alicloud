// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAligreenKeywordLib() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAligreenKeywordLibCreate,
		Read:   resourceAliCloudAligreenKeywordLibRead,
		Update: resourceAliCloudAligreenKeywordLibUpdate,
		Delete: resourceAliCloudAligreenKeywordLibDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"biz_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"keyword_lib_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"language": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"lib_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"match_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAligreenKeywordLibCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateKeywordLib"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Name"] = d.Get("keyword_lib_name")
	request["ResourceType"] = d.Get("resource_type")
	if v, ok := d.GetOk("lib_type"); ok {
		request["LibType"] = v
	}
	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}
	if v, ok := d.GetOk("match_mode"); ok {
		request["MatchMode"] = v
	}
	if v, ok := d.GetOk("enable"); ok {
		request["Enable"] = v
	}
	request["ServiceModule"] = "open_api"
	if v, ok := d.GetOk("language"); ok {
		request["Language"] = v
	}
	if v, ok := d.GetOk("biz_types"); ok {
		jsonPathResult7, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult7 != "" {
			request["BizTypes"] = convertListToJsonString(jsonPathResult7.([]interface{}))
		}
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_aligreen_keyword_lib", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudAligreenKeywordLibRead(d, meta)
}

func resourceAliCloudAligreenKeywordLibRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	aligreenServiceV2 := AligreenServiceV2{client}

	objectRaw, err := aligreenServiceV2.DescribeAligreenKeywordLib(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_aligreen_keyword_lib DescribeAligreenKeywordLib Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Category"] != nil {
		d.Set("category", objectRaw["Category"])
	}
	if objectRaw["Enable"] != nil {
		d.Set("enable", objectRaw["Enable"])
	}
	if objectRaw["Name"] != nil {
		d.Set("keyword_lib_name", objectRaw["Name"])
	}
	if objectRaw["Language"] != nil {
		d.Set("language", objectRaw["Language"])
	}
	if objectRaw["LibType"] != nil {
		d.Set("lib_type", objectRaw["LibType"])
	}
	if objectRaw["MatchMode"] != nil {
		d.Set("match_mode", objectRaw["MatchMode"])
	}
	if objectRaw["ResourceType"] != nil {
		d.Set("resource_type", objectRaw["ResourceType"])
	}

	bizTypes1Raw := make([]interface{}, 0)
	if objectRaw["BizTypes"] != nil {
		bizTypes1Raw = objectRaw["BizTypes"].([]interface{})
	}

	d.Set("biz_types", bizTypes1Raw)

	return nil
}

func resourceAliCloudAligreenKeywordLibUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateKeywordLib"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = d.Id()

	if d.HasChange("keyword_lib_name") {
		update = true
	}
	request["Name"] = d.Get("keyword_lib_name")
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if d.HasChange("biz_types") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$", d.Get("biz_types"))
		if err == nil {
			request["BizTypes"] = convertListToJsonString(jsonPathResult2.([]interface{}))
		}
	}

	if d.HasChange("match_mode") {
		update = true
		request["MatchMode"] = d.Get("match_mode")
	}

	if d.HasChange("enable") {
		update = true
		request["Enable"] = d.Get("enable")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Green", "2017-08-23", action, query, request, false)
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

	return resourceAliCloudAligreenKeywordLibRead(d, meta)
}

func resourceAliCloudAligreenKeywordLibDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteKeywordLib"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Id"] = d.Id()

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, false)

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
