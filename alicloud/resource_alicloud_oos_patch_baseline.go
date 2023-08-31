// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOosPatchBaseline() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOosPatchBaselineCreate,
		Read:   resourceAliCloudOosPatchBaselineRead,
		Update: resourceAliCloudOosPatchBaselineUpdate,
		Delete: resourceAliCloudOosPatchBaselineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"approval_rules": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operation_system": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Windows", "Ubuntu", "Debian", "AliyunLinux", "RedhatEnterpriseLinux", "Anolis", "CentOS", "AlmaLinux"}, false),
			},
			"patch_baseline_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rejected_patches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rejected_patches_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ALLOW_AS_DEPENDENCY", "BLOCK"}, false),
			},
		},
	}
}

func resourceAliCloudOosPatchBaselineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePatchBaseline"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["Name"] = d.Get("patch_baseline_name")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["OperationSystem"] = d.Get("operation_system")
	request["ApprovalRules"] = d.Get("approval_rules")
	if v, ok := d.GetOk("rejected_patches"); ok {
		rejectedPatchesMaps := v.([]interface{})
		rejectedPatchesMapsJson, err := json.Marshal(rejectedPatchesMaps)
		if err != nil {
			return WrapError(err)
		}
		request["RejectedPatches"] = string(rejectedPatchesMapsJson)
	}

	if v, ok := d.GetOk("rejected_patches_action"); ok {
		request["RejectedPatchesAction"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_patch_baseline", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.PatchBaseline.Name", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudOosPatchBaselineUpdate(d, meta)
}

func resourceAliCloudOosPatchBaselineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosServiceV2 := OosServiceV2{client}

	objectRaw, err := oosServiceV2.DescribeOosPatchBaseline(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_patch_baseline DescribeOosPatchBaseline Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	jsonMarshalResult, err := json.Marshal(objectRaw["ApprovalRules"])
	if err != nil {
		return WrapError(err)
	}
	d.Set("approval_rules", string(jsonMarshalResult))
	d.Set("create_time", objectRaw["CreatedDate"])
	d.Set("description", objectRaw["Description"])
	d.Set("operation_system", objectRaw["OperationSystem"])
	d.Set("rejected_patches_action", objectRaw["RejectedPatchesAction"])
	d.Set("patch_baseline_name", objectRaw["Name"])
	rejectedPatches1Raw := make([]interface{}, 0)
	if objectRaw["RejectedPatches"] != nil {
		rejectedPatches1Raw = objectRaw["RejectedPatches"].([]interface{})
	}

	d.Set("rejected_patches", rejectedPatches1Raw)

	return nil
}

func resourceAliCloudOosPatchBaselineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "UpdatePatchBaseline"
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["Name"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("approval_rules") {
		update = true
	}
	request["ApprovalRules"] = d.Get("approval_rules")
	if !d.IsNewResource() && d.HasChange("rejected_patches") {
		update = true
		if v, ok := d.GetOk("rejected_patches"); ok {
			rejectedPatchesMaps := v.([]interface{})
			rejectedPatchesMapsJson, err := json.Marshal(rejectedPatchesMaps)
			if err != nil {
				return WrapError(err)
			}
			request["RejectedPatches"] = string(rejectedPatchesMapsJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("rejected_patches_action") {
		update = true
		request["RejectedPatchesAction"] = d.Get("rejected_patches_action")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

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

	return resourceAliCloudOosPatchBaselineRead(d, meta)
}

func resourceAliCloudOosPatchBaselineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeletePatchBaseline"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["Name"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
