// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudRealtimeComputeMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRealtimeComputeMemberCreate,
		Read:   resourceAliCloudRealtimeComputeMemberRead,
		Update: resourceAliCloudRealtimeComputeMemberUpdate,
		Delete: resourceAliCloudRealtimeComputeMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"member": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"viewer", "owner", "editor"}, false),
			},
		},
	}
}

func resourceAliCloudRealtimeComputeMemberCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	namespace := d.Get("namespace")
	action := fmt.Sprintf("/gateway/v2/namespaces/%s/members", namespace)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(d.Get("resource_id").(string))
	if v, ok := d.GetOk("member"); ok {
		request["member"] = v
	}

	if v, ok := d.GetOk("role"); ok {
		request["role"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("ververica", "2022-07-18", action, query, header, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, " alicloud_realtime_compute_member", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	dataMemberVar, _ := jsonpath.Get("$.data.member", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", d.Get("resource_id").(string), namespace, dataMemberVar))

	return resourceAliCloudRealtimeComputeMemberRead(d, meta)
}

func resourceAliCloudRealtimeComputeMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}

	objectRaw, err := realtimeComputeServiceV2.DescribeRealtimeComputeMember(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource  alicloud_realtime_compute_member DescribeRealtimecomputemember Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("role", objectRaw["role"])
	d.Set("member", objectRaw["member"])

	parts := strings.Split(d.Id(), ":")
	d.Set("resource_id", parts[0])
	d.Set("namespace", parts[1])

	return nil
}

func resourceAliCloudRealtimeComputeMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	action := fmt.Sprintf("/gateway/v2/namespaces/%s/members", namespace)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	body = make(map[string]interface{})
	header["workspace"] = StringPointer(parts[0])
	request["member"] = parts[2]

	if d.HasChange("role") {
		update = true
	}
	if v, ok := d.GetOk("role"); ok || d.HasChange("role") {
		request["role"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("ververica", "2022-07-18", action, query, header, body, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudRealtimeComputeMemberRead(d, meta)
}

func resourceAliCloudRealtimeComputeMemberDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	namespace := parts[1]
	member := parts[2]
	action := fmt.Sprintf("/gateway/v2/namespaces/%s/members/%s", namespace, member)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("ververica", "2022-07-18", action, query, header, nil, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"990301"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
