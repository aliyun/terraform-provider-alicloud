// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package ram

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceRole() *schema.Resource {
	return &schema.Resource{
		Create: ResourceRoleCreate,
		Read:   ResourceRoleRead,
		Update: ResourceRoleUpdate,
		Delete: ResourceRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assume_role_policy_document": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				ConflictsWith: []string{"document", "version", "ram_users", "services"},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_session_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"tags": tags.TagsSchema(),
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.252.0. New field 'role_name' instead.",
				ConflictsWith: []string{"role_name"},
			},
			"document": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				Deprecated:    "Field 'document' has been deprecated from provider version 1.252.0. New field 'assume_role_policy_document' instead.",
				ConflictsWith: []string{"assume_role_policy_document", "version", "ram_users", "services"},
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1",
				// can only be '1' so far.
				ValidateFunc:  validators.StringInSlice([]string{"1"}, false),
				Deprecated:    "Field 'version' has been deprecated from provider version 1.49.0. New field 'document' instead.",
				ConflictsWith: []string{"assume_role_policy_document", "document"},
			},
			"ram_users": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Set:           schema.HashString,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Deprecated:    "Field 'ram_users' has been deprecated from provider version 1.49.0. New field 'document' instead.",
				ConflictsWith: []string{"assume_role_policy_document", "document"},
			},
			"services": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Set:           schema.HashString,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Deprecated:    "Field 'services' has been deprecated from provider version 1.49.0. New field 'document' instead.",
				ConflictsWith: []string{"assume_role_policy_document", "document"},
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func ResourceRoleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRole"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("role_name"); ok {
		request["RoleName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["RoleName"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("max_session_duration"); ok {
		request["MaxSessionDuration"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		tagsMapJSON, err := convert.ConvertListMapToJsonString(tagsMap)
		if err != nil {
			return errs.WrapError(err)
		}
		request["Tag"] = tagsMapJSON
	}

	assumeRolePolicyDocument, assumeRolePolicyDocumentOk := d.GetOk("assume_role_policy_document")
	document, documentOk := d.GetOk("document")
	ramUsers, usersOk := d.GetOk("ram_users")
	services, servicesOk := d.GetOk("services")

	if !assumeRolePolicyDocumentOk && !documentOk && !usersOk && !servicesOk {
		return errs.WrapError(errs.Error("At least one of 'assume_role_policy_document', 'document', 'ram_users' or 'services' must be set."))
	}

	if assumeRolePolicyDocumentOk {
		request["AssumeRolePolicyDocument"] = assumeRolePolicyDocument
	} else if documentOk {
		request["AssumeRolePolicyDocument"] = document
	} else {
		ramServiceV2 := RamServiceV2{client}

		rolePolicyDocument, err := ramServiceV2.AssembleRolePolicyDocument(ramUsers.(*schema.Set).List(), services.(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return errs.WrapError(err)
		}

		request["AssumeRolePolicyDocument"] = rolePolicyDocument
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ram_role", action, errs.AlibabaCloudSdkGoERROR)
	}

	responseRole := response["Role"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseRole["RoleName"]))

	return ResourceRoleRead(d, meta)
}

func ResourceRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	objectRaw, err := ramServiceV2.DescribeRamRole(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_role DescribeRamRole Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("arn", objectRaw["Arn"])
	d.Set("assume_role_policy_document", objectRaw["AssumeRolePolicyDocument"])
	d.Set("create_time", objectRaw["CreateDate"])
	d.Set("description", objectRaw["Description"])
	d.Set("max_session_duration", objectRaw["MaxSessionDuration"])
	d.Set("role_id", objectRaw["RoleId"])
	d.Set("role_name", objectRaw["RoleName"])
	d.Set("name", objectRaw["RoleName"])
	d.Set("document", objectRaw["AssumeRolePolicyDocument"])

	if v, ok := objectRaw["AssumeRolePolicyDocument"].(string); ok && v != "" {
		assumeRolePolicyDocumentArg, err := convert.ConvertJsonStringToMap(v)
		if err != nil {
			return errs.WrapError(err)
		}

		if version, ok := assumeRolePolicyDocumentArg["Version"]; ok {
			d.Set("version", version)
		}

		if statement, ok := assumeRolePolicyDocumentArg["Statement"]; ok {
			statementList := statement.([]interface{})
			if len(statementList) > 0 {
				if principal, ok := statementList[0].(map[string]interface{})["Principal"]; ok {
					principalArg := principal.(map[string]interface{})
					d.Set("ram_users", principalArg["RAM"])
					d.Set("services", principalArg["Service"])
				}
			}
		}
	}

	objectRaw, err = ramServiceV2.DescribeRoleListTagResources(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func ResourceRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateRole"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RoleName"] = d.Id()

	if d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
	}

	if d.HasChange("max_session_duration") {
		update = true

		if v, ok := d.GetOkExists("max_session_duration"); ok {
			request["NewMaxSessionDuration"] = v
		}
	}

	if d.HasChange("assume_role_policy_document") {
		update = true
		request["NewAssumeRolePolicyDocument"] = d.Get("assume_role_policy_document")
	}

	if d.HasChange("document") {
		update = true
		request["NewAssumeRolePolicyDocument"] = d.Get("document")
	}

	if d.HasChange("ram_users") || d.HasChange("services") || d.HasChange("version") {
		update = true
		ramServiceV2 := RamServiceV2{client}

		rolePolicyDocument, err := ramServiceV2.AssembleRolePolicyDocument(d.Get("ram_users").(*schema.Set).List(), d.Get("services").(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return errs.WrapError(err)
		}

		request["NewAssumeRolePolicyDocument"] = rolePolicyDocument
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		ramServiceV2 := RamServiceV2{client}
		if err := ramServiceV2.SetResourceTags(d, "role"); err != nil {
			return errs.WrapError(err)
		}
	}
	return ResourceRoleRead(d, meta)
}

func ResourceRoleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRole"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RoleName"] = d.Id()

	if d.Get("force").(bool) {
		listRequest := map[string]interface{}{
			"RoleName": d.Id(),
		}
		listAction := "ListPoliciesForRole"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", listAction, nil, listRequest, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(listAction, response, listRequest)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) || errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), listAction, errs.AlibabaCloudSdkGoERROR)
		}

		policyResp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, listAction, "$.Policies.Policy", response)
		}

		// Loop and remove the Policies from the Role
		if policyResp != nil && len(policyResp.([]interface{})) > 0 {
			for _, v := range policyResp.([]interface{}) {
				policyAction := "DetachPolicyFromRole"
				policyRequest := map[string]interface{}{
					"RoleName":   d.Id(),
					"PolicyName": v.(map[string]interface{})["PolicyName"],
					"PolicyType": v.(map[string]interface{})["PolicyType"],
				}

				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
					response, err = client.RpcPost("Ram", "2015-05-01", policyAction, nil, policyRequest, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(policyAction, response, policyRequest)

				if err != nil {
					if errs.IsExpectedErrors(err, []string{"EntityNotExist", "EntityNotExist.Role", "EntityNotExist.Role.Policy", "EntityNotExist.Policy"}) || errs.NotFoundError(err) {
						continue
					}
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), policyAction, errs.AlibabaCloudSdkGoERROR)
				}
			}

		}
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"DeleteConflict.Role.Policy"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
