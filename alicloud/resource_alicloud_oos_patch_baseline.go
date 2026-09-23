// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
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
				Type:     schema.TypeString,
				Required: true,
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
				ValidateFunc: StringInSlice([]string{"Windows", "CentOS", "AliyunLinux", "Ubuntu", "Debian", "Fedora", "Suse", "RockyLinux", "RedhatEnterpriseLinux", "Anolis", "AlmaLinux"}, false),
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
	if v, ok := d.GetOk("patch_baseline_name"); ok {
		request["Name"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["ApprovalRules"] = d.Get("approval_rules")
	if v, ok := d.GetOk("sources"); ok {
		sourcesMapsArray := convertToInterfaceArray(v)

		sourcesMapsJson, err := json.Marshal(sourcesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Sources"] = string(sourcesMapsJson)
	}

	if v, ok := d.GetOk("rejected_patches"); ok {
		rejectedPatchesMapsArray := convertToInterfaceArray(v)

		rejectedPatchesMapsJson, err := json.Marshal(rejectedPatchesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["RejectedPatches"] = string(rejectedPatchesMapsJson)
	}

	if v, ok := d.GetOk("approved_patches"); ok {
		approvedPatchesMapsArray := convertToInterfaceArray(v)

		approvedPatchesMapsJson, err := json.Marshal(approvedPatchesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["ApprovedPatches"] = string(approvedPatchesMapsJson)
	}

	if v, ok := d.GetOkExists("approved_patches_enable_non_security"); ok {
		request["ApprovedPatchesEnableNonSecurity"] = v
	}
	request["OperationSystem"] = d.Get("operation_system")
	if v, ok := d.GetOk("rejected_patches_action"); ok {
		request["RejectedPatchesAction"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)
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
	d.Set("approved_patches_enable_non_security", objectRaw["ApprovedPatchesEnableNonSecurity"])
	d.Set("create_time", objectRaw["CreatedDate"])
	d.Set("description", objectRaw["Description"])
	d.Set("operation_system", objectRaw["OperationSystem"])
	d.Set("rejected_patches_action", objectRaw["RejectedPatchesAction"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("patch_baseline_name", objectRaw["Name"])

	approvedPatchesRaw := make([]interface{}, 0)
	if objectRaw["ApprovedPatches"] != nil {
		approvedPatchesRaw = convertToInterfaceArray(objectRaw["ApprovedPatches"])
	}

	d.Set("approved_patches", approvedPatchesRaw)
	rejectedPatchesRaw := make([]interface{}, 0)
	if objectRaw["RejectedPatches"] != nil {
		rejectedPatchesRaw = convertToInterfaceArray(objectRaw["RejectedPatches"])
	}

	d.Set("rejected_patches", rejectedPatchesRaw)
	sourcesRaw := make([]interface{}, 0)
	if objectRaw["Sources"] != nil {
		sourcesRaw = convertToInterfaceArray(objectRaw["Sources"])
	}

	d.Set("sources", sourcesRaw)
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

	var err error
	action := "UpdatePatchBaseline"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Name"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok || d.HasChange("resource_group_id") {
		request["ResourceGroupId"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["Description"] = v
	}
	if d.HasChange("approval_rules") {
		update = true
	}
	request["ApprovalRules"] = d.Get("approval_rules")
	if d.HasChange("sources") {
		update = true
	}
	if v, ok := d.GetOk("sources"); ok || d.HasChange("sources") {
		sourcesMapsArray := convertToInterfaceArray(v)

		sourcesMapsJson, err := json.Marshal(sourcesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Sources"] = string(sourcesMapsJson)
	}

	if d.HasChange("rejected_patches") {
		update = true
	}
	if v, ok := d.GetOk("rejected_patches"); ok || d.HasChange("rejected_patches") {
		rejectedPatchesMapsArray := convertToInterfaceArray(v)

		rejectedPatchesMapsJson, err := json.Marshal(rejectedPatchesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["RejectedPatches"] = string(rejectedPatchesMapsJson)
	}

	if d.HasChange("approved_patches") {
		update = true
	}
	if v, ok := d.GetOk("approved_patches"); ok || d.HasChange("approved_patches") {
		approvedPatchesMapsArray := convertToInterfaceArray(v)

		approvedPatchesMapsJson, err := json.Marshal(approvedPatchesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["ApprovedPatches"] = string(approvedPatchesMapsJson)
	}

	if d.HasChange("approved_patches_enable_non_security") {
		update = true
	}
	if v, ok := d.GetOkExists("approved_patches_enable_non_security"); ok || d.HasChange("approved_patches_enable_non_security") {
		request["ApprovedPatchesEnableNonSecurity"] = v
	}
	if d.HasChange("rejected_patches_action") {
		update = true
	}
	if v, ok := d.GetOk("rejected_patches_action"); ok || d.HasChange("rejected_patches_action") {
		request["RejectedPatchesAction"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("oos", "2019-06-01", action, query, request, true)
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

	if d.HasChange("tags") {
		oosServiceV2 := OosServiceV2{client}
		if err := oosServiceV2.SetOssResourceTags(d, "patchbaseline"); err != nil {
			return WrapError(err)
		}
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
	request["Name"] = d.Id()
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
