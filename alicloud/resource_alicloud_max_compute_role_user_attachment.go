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
	"github.com/tidwall/sjson"
)

func resourceAliCloudMaxComputeRoleUserAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMaxComputeRoleUserAttachmentCreate,
		Read:   resourceAliCloudMaxComputeRoleUserAttachmentRead,
		Delete: resourceAliCloudMaxComputeRoleUserAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudMaxComputeRoleUserAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	projectName := d.Get("project_name")
	roleName := d.Get("role_name")
	action := fmt.Sprintf("/api/v1/projects/%s/roles/%s/users", projectName, roleName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "add.0", d.Get("user"))
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_max_compute_role_user_attachment", action, AlibabaCloudSdkGoERROR)
	}

	addVar, _ := jsonpath.Get("add[0]", request)
	d.SetId(fmt.Sprintf("%v&%v&%v", projectName, roleName, addVar))

	return resourceAliCloudMaxComputeRoleUserAttachmentRead(d, meta)
}

func resourceAliCloudMaxComputeRoleUserAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxComputeServiceV2 := MaxComputeServiceV2{client}

	objectRaw, err := maxComputeServiceV2.DescribeMaxComputeRoleUserAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_max_compute_role_user_attachment DescribeMaxComputeRoleUserAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["name"] != nil {
		d.Set("user", objectRaw["name"])
	}

	parts := strings.Split(d.Id(), "&")
	d.Set("project_name", parts[0])
	d.Set("role_name", parts[1])

	return nil
}

func resourceAliCloudMaxComputeRoleUserAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), "&")
	projectName := parts[0]
	roleName := parts[1]
	action := fmt.Sprintf("/api/v1/projects/%s/roles/%s/users", projectName, roleName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "remove.0", parts[2])
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)

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
