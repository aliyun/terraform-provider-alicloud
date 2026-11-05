package maxcompute

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	jsonString := convert.ConvertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "grant.0", d.Get("tenant_role"))
	_ = json.Unmarshal([]byte(jsonString), &request)

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_max_compute_tenant_role_user_attachment", action, errs.AlibabaCloudSdkGoERROR)
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
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_max_compute_tenant_role_user_attachment DescribeMaxComputeTenantRoleUserAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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

	jsonString := convert.ConvertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "revoke.0", parts[1])
	_ = json.Unmarshal([]byte(jsonString), &request)

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaPut("MaxCompute", "2022-01-04", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
