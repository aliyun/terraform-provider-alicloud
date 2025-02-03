// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"approved_patches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"approved_patches_enable_non_security": {
				Type:     schema.TypeBool,
				Optional: true,
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
				ValidateFunc: StringInSlice([]string{"Windows", "CentOS", "AliyunLinux", "Ubuntu", "Debian", "RedhatEnterpriseLinux", "Anolis", "AlmaLinux"}, true),
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
				ValidateFunc: StringInSlice([]string{"ALLOW_AS_DEPENDENCY", "BLOCK"}, true),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudOosPatchBaselineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePatchBaseline"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Name"] = d.Get("patch_baseline_name")
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
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("approved_patches_enable_non_security"); ok {
		request["ApprovedPatchesEnableNonSecurity"] = v
	}
	if v, ok := d.GetOk("approved_patches"); ok {
		approvedPatchesMaps := v.([]interface{})
		request["ApprovedPatches"] = approvedPatchesMaps
	}

	if v, ok := d.GetOk("sources"); ok {
		sourcesMaps := v.([]interface{})
		request["Sources"] = sourcesMaps
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)
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

	return resourceAliCloudOosPatchBaselineRead(d, meta)
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
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("patch_baseline_name", objectRaw["Name"])

	approvedPatches1Raw := make([]interface{}, 0)
	if objectRaw["ApprovedPatches"] != nil {
		approvedPatches1Raw = objectRaw["ApprovedPatches"].([]interface{})
	}

	d.Set("approved_patches", approvedPatches1Raw)
	rejectedPatches1Raw := make([]interface{}, 0)
	if objectRaw["RejectedPatches"] != nil {
		rejectedPatches1Raw = objectRaw["RejectedPatches"].([]interface{})
	}

	d.Set("rejected_patches", rejectedPatches1Raw)
	sources1Raw := make([]interface{}, 0)
	if objectRaw["Sources"] != nil {
		sources1Raw = objectRaw["Sources"].([]interface{})
	}

	d.Set("sources", sources1Raw)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("patch_baseline_name", d.Id())

	return nil
}

func resourceAliCloudOosPatchBaselineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdatePatchBaseline"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Name"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("approval_rules") {
		update = true
	}
	request["ApprovalRules"] = d.Get("approval_rules")
	if d.HasChange("rejected_patches") {
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

	if d.HasChange("rejected_patches_action") {
		update = true
		request["RejectedPatchesAction"] = d.Get("rejected_patches_action")
	}

	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if d.HasChange("approved_patches_enable_non_security") {
		update = true
		request["ApprovedPatchesEnableNonSecurity"] = d.Get("approved_patches_enable_non_security")
	}

	if d.HasChange("approved_patches") {
		update = true
		if v, ok := d.GetOk("approved_patches"); ok {
			approvedPatchesMaps := v.([]interface{})
			approvedPatchesMapsJson, err := json.Marshal(approvedPatchesMaps)
			if err != nil {
				return WrapError(err)
			}
			request["ApprovedPatches"] = string(approvedPatchesMapsJson)
		}
	}

	if d.HasChange("sources") {
		update = true
		if v, ok := d.GetOk("sources"); ok {
			sourcesMaps := v.([]interface{})
			sourcesMapsJson, err := json.Marshal(sourcesMaps)
			if err != nil {
				return WrapError(err)
			}
			request["Sources"] = string(sourcesMapsJson)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)
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

	if d.HasChange("tags") {
		oosServiceV2 := OosServiceV2{client}
		if err := oosServiceV2.SetResourceTags(d, "patchbaseline"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	return resourceAliCloudOosPatchBaselineRead(d, meta)
}

func resourceAliCloudOosPatchBaselineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeletePatchBaseline"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Name"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)

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
