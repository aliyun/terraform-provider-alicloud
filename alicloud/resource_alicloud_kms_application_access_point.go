// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/blues/jsonata-go"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudKmsApplicationAccessPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudKmsApplicationAccessPointCreate,
		Read:   resourceAliCloudKmsApplicationAccessPointRead,
		Update: resourceAliCloudKmsApplicationAccessPointUpdate,
		Delete: resourceAliCloudKmsApplicationAccessPointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"application_access_point_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudKmsApplicationAccessPointCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateApplicationAccessPoint"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["Name"] = d.Get("application_access_point_name")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	jsonPathResult1, err := jsonpath.Get("$", d.Get("policies"))
	if err != nil {
		return WrapError(err)
	}
	request["Policies"], _ = convertArrayObjectToJsonString(jsonPathResult1)

	request["AuthenticationMethod"] = "ClientKey"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_application_access_point", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Name"]))

	return resourceAliCloudKmsApplicationAccessPointRead(d, meta)
}

func resourceAliCloudKmsApplicationAccessPointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsServiceV2 := KmsServiceV2{client}

	objectRaw, err := kmsServiceV2.DescribeKmsApplicationAccessPoint(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_application_access_point DescribeKmsApplicationAccessPoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("application_access_point_name", objectRaw["Name"])

	e := jsonata.MustCompile("$map($split($substring($.Policies, 1, $length($.Policies)-2), \",\"), function($v, $i, $a) {$substring($v, 1, $length($v)-2)})")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("policies", evaluation)

	return nil
}

func resourceAliCloudKmsApplicationAccessPointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "UpdateApplicationAccessPoint"
	var err error
	request = make(map[string]interface{})
	request["Name"] = d.Id()
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("policies") {
		update = true
	}
	jsonPathResult1, err := jsonpath.Get("$", d.Get("policies"))
	if err != nil {
		return WrapError(err)
	}
	request["Policies"], _ = convertArrayObjectToJsonString(jsonPathResult1)

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)

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

	return resourceAliCloudKmsApplicationAccessPointRead(d, meta)
}

func resourceAliCloudKmsApplicationAccessPointDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteApplicationAccessPoint"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["Name"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)

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
		if IsExpectedErrors(err, []string{"Forbidden.ResourceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
