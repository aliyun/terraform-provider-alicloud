// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEhpcUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEhpcUserCreate,
		Read:   resourceAliCloudEhpcUserRead,
		Update: resourceAliCloudEhpcUserUpdate,
		Delete: resourceAliCloudEhpcUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"async": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEhpcUserCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddUsers"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["ClusterId"] = d.Get("cluster_id")

	userMap := map[string]interface{}{
		"Name":  d.Get("user_name"),
		"Group": d.Get("group"),
	}
	if v, ok := d.GetOk("password"); ok {
		userMap["Password"] = v
	}

	request["User"] = []map[string]interface{}{userMap}

	if v, ok := d.GetOk("async"); ok {
		request["Async"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcGet("EHPC", "2018-04-12", action, request, nil)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ehpc_user", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ClusterId"], d.Get("user_name")))

	return resourceAliCloudEhpcUserRead(d, meta)
}

func resourceAliCloudEhpcUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ehpcServiceV2 := EhpcServiceV2{client}

	objectRaw, err := ehpcServiceV2.DescribeEhpcUser(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ehpc_user DescribeEhpcUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("group", objectRaw["Group"])
	d.Set("user_name", objectRaw["Name"])
	d.Set("region_id", client.RegionId)

	parts := strings.Split(d.Id(), ":")
	d.Set("cluster_id", parts[0])

	return nil
}

func resourceAliCloudEhpcUserUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyUserGroups"
	request = make(map[string]interface{})
	request["ClusterId"] = parts[0]

	userMap := map[string]interface{}{
		"Name":  parts[1],
		"Group": d.Get("group"),
	}

	request["User"] = []map[string]interface{}{userMap}

	if v, ok := d.GetOk("async"); ok {
		request["Async"] = v
	}

	if d.HasChange("group") {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcGet("EHPC", "2018-04-12", action, request, nil)
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

	return resourceAliCloudEhpcUserRead(d, meta)
}

func resourceAliCloudEhpcUserDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteUsers"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["ClusterId"] = parts[0]

	userMap := map[string]interface{}{
		"Name": parts[1],
	}

	request["User"] = []map[string]interface{}{userMap}

	if v, ok := d.GetOk("async"); ok {
		request["Async"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcGet("EHPC", "2018-04-12", action, request, nil)
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
