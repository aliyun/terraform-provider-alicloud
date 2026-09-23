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

func resourceAliCloudMaxComputeTenantRoleUserAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMaxComputeTenantRoleUserAttachmentCreate,
		Read:   resourceAliCloudMaxComputeTenantRoleUserAttachmentRead,
		Delete: resourceAliCloudMaxComputeTenantRoleUserAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tenant_role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudMaxComputeTenantRoleUserAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/tenants/user/roles/new")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("account_id"); ok {
		request["accountId"] = v
	}

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "grant.0", d.Get("tenant_role"))
	_ = json.Unmarshal([]byte(jsonString), &request)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_max_compute_tenant_role_user_attachment", action, AlibabaCloudSdkGoERROR)
	}

	grantVar, _ := jsonpath.Get("grant[0]", request)
	d.SetId(fmt.Sprintf("%v:%v", request["accountId"], grantVar))

	return resourceAliCloudMaxComputeTenantRoleUserAttachmentRead(d, meta)
}

func resourceAliCloudMaxComputeTenantRoleUserAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxComputeServiceV2 := MaxComputeServiceV2{client}

	_, err := maxComputeServiceV2.DescribeMaxComputeTenantRoleUserAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_max_compute_tenant_role_user_attachment DescribeMaxComputeTenantRoleUserAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("account_id", parts[0])
	d.Set("tenant_role", parts[1])

	return nil
}

func resourceAliCloudMaxComputeTenantRoleUserAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/api/v1/tenants/user/roles/new")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["accountId"] = parts[0]

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "revoke.0", parts[1])
	_ = json.Unmarshal([]byte(jsonString), &request)

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
