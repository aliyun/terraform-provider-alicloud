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

func resourceAliCloudAligreenImageLib() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAligreenImageLibCreate,
		Read:   resourceAliCloudAligreenImageLibRead,
		Update: resourceAliCloudAligreenImageLibUpdate,
		Delete: resourceAliCloudAligreenImageLibDelete,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"image_lib_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scene": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAligreenImageLibCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateImageLib"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ServiceModule"] = "open_api"
	request["Scene"] = d.Get("scene")
	request["Category"] = d.Get("category")
	if v, ok := d.GetOk("enable"); ok {
		request["Enable"] = v
	}
	request["Name"] = d.Get("image_lib_name")
	if v, ok := d.GetOk("biz_types"); ok {
		jsonPathResult4, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult4 != "" {
			request["BizTypes"] = convertListToJsonString(jsonPathResult4.([]interface{}))
		}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_aligreen_image_lib", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudAligreenImageLibUpdate(d, meta)
}

func resourceAliCloudAligreenImageLibRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	aligreenServiceV2 := AligreenServiceV2{client}

	objectRaw, err := aligreenServiceV2.DescribeAligreenImageLib(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_aligreen_image_lib DescribeAligreenImageLib Failed!!! %s", err)
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
		d.Set("image_lib_name", objectRaw["Name"])
	}
	if objectRaw["Scene"] != nil {
		d.Set("scene", objectRaw["Scene"])
	}

	bizTypes1Raw := make([]interface{}, 0)
	if objectRaw["BizTypes"] != nil {
		bizTypes1Raw = objectRaw["BizTypes"].([]interface{})
	}

	d.Set("biz_types", bizTypes1Raw)

	return nil
}

func resourceAliCloudAligreenImageLibUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateImageLib"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = d.Id()

	if !d.IsNewResource() && d.HasChange("scene") {
		update = true
	}
	request["Scene"] = d.Get("scene")
	if !d.IsNewResource() && d.HasChange("category") {
		update = true
	}
	request["Category"] = d.Get("category")
	if !d.IsNewResource() && d.HasChange("enable") {
		update = true
		request["Enable"] = d.Get("enable")
	}

	if !d.IsNewResource() && d.HasChange("image_lib_name") {
		update = true
	}
	request["Name"] = d.Get("image_lib_name")
	if !d.IsNewResource() && d.HasChange("biz_types") {
		update = true
	}
	jsonPathResult4, err := jsonpath.Get("$", d.Get("biz_types"))
	if err == nil {
		request["BizTypes"] = convertListToJsonString(jsonPathResult4.([]interface{}))
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

	return resourceAliCloudAligreenImageLibRead(d, meta)
}

func resourceAliCloudAligreenImageLibDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteImageLib"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Id"] = d.Id()

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
