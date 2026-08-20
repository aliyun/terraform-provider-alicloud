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
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"member": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				// The Flink member API is case-insensitive for role and normalizes
				// the value to lowercase on read (e.g. config "VIEWER" is returned
				// as "viewer"). Suppress pure-case diffs so the state stored from
				// the API response does not produce a spurious update on the next
				// plan; a genuine role change still produces a diff.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRealtimeComputeMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceId := d.Get("resource_id").(string)
	namespace := d.Get("namespace").(string)
	member := d.Get("member").(string)
	action := fmt.Sprintf("/gateway/v2/namespaces/%s/members", namespace)
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	var body map[string]interface{}
	var err error
	header["workspace"] = StringPointer(resourceId)

	bodyMap := make(map[string]interface{})
	bodyMap["member"] = member
	if v, ok := d.GetOk("role"); ok {
		bodyMap["role"] = v
	}
	body = bodyMap

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("ververica", "2022-07-18", action, query, header, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, body)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", resourceId, namespace, member))

	return resourceAliCloudRealtimeComputeMemberRead(d, meta)
}

func resourceAliCloudRealtimeComputeMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}

	objectRaw, err := realtimeComputeServiceV2.DescribeRealtimeComputeFlinkMember(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_realtime_compute_member DescribeRealtimeComputeFlinkMember Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("resource_id", parts[0])
	d.Set("namespace", parts[1])
	d.Set("member", parts[2])
	d.Set("role", objectRaw["role"])
	d.Set("region_id", client.RegionId)

	return nil
}

func resourceAliCloudRealtimeComputeMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	resourceId := parts[0]
	namespace := parts[1]
	member := parts[2]
	action := fmt.Sprintf("/gateway/v2/namespaces/%s/members", namespace)
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	var body map[string]interface{}
	var err error
	header["workspace"] = StringPointer(resourceId)

	if d.HasChange("role") {
		bodyMap := make(map[string]interface{})
		bodyMap["member"] = member
		if v, ok := d.GetOk("role"); ok {
			bodyMap["role"] = v
		}
		body = bodyMap
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("ververica", "2022-07-18", action, query, header, body, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, body)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudRealtimeComputeMemberRead(d, meta)
}

func resourceAliCloudRealtimeComputeMemberDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	resourceId := parts[0]
	namespace := parts[1]
	member := parts[2]
	action := fmt.Sprintf("/gateway/v2/namespaces/%s/members/%s", namespace, member)
	var request map[string]interface{}
	var response map[string]interface{}
	header := make(map[string]*string)
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	header["workspace"] = StringPointer(resourceId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("ververica", "2022-07-18", action, query, header, nil, true)
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
		if IsExpectedErrors(err, []string{"990301"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
