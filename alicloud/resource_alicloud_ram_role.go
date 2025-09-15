// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudRamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamRoleCreate,
		Read:   resourceAliCloudRamRoleRead,
		Update: resourceAliCloudRamRoleUpdate,
		Delete: resourceAliCloudRamRoleDelete,
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
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
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
			"tags": tagsSchema(),
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
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
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
				ValidateFunc:  StringInSlice([]string{"1"}, false),
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

func resourceAliCloudRamRoleCreate(d *schema.ResourceData, meta interface{}) error {

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
		tagsMap := ConvertTags(v.(map[string]interface{}))
		tagsMapJSON, err := convertListMapToJsonString(tagsMap)
		if err != nil {
			return WrapError(err)
		}
		request["Tag"] = tagsMapJSON
	}

	assumeRolePolicyDocument, assumeRolePolicyDocumentOk := d.GetOk("assume_role_policy_document")
	document, documentOk := d.GetOk("document")
	ramUsers, usersOk := d.GetOk("ram_users")
	services, servicesOk := d.GetOk("services")

	if !assumeRolePolicyDocumentOk && !documentOk && !usersOk && !servicesOk {
		return WrapError(Error("At least one of 'assume_role_policy_document', 'document', 'ram_users' or 'services' must be set."))
	}

	if assumeRolePolicyDocumentOk {
		request["AssumeRolePolicyDocument"] = assumeRolePolicyDocument
	} else if documentOk {
		request["AssumeRolePolicyDocument"] = document
	} else {
		ramServiceV2 := RamServiceV2{client}

		rolePolicyDocument, err := ramServiceV2.AssembleRolePolicyDocument(ramUsers.(*schema.Set).List(), services.(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return WrapError(err)
		}

		request["AssumeRolePolicyDocument"] = rolePolicyDocument
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_role", action, AlibabaCloudSdkGoERROR)
	}

	responseRole := response["Role"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseRole["RoleName"]))

	return resourceAliCloudRamRoleRead(d, meta)
}

func resourceAliCloudRamRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	objectRaw, err := ramServiceV2.DescribeRamRole(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_role DescribeRamRole Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		assumeRolePolicyDocumentArg, err := convertJsonStringToMap(v)
		if err != nil {
			return WrapError(err)
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
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudRamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
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
			return WrapError(err)
		}

		request["NewAssumeRolePolicyDocument"] = rolePolicyDocument
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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
		ramServiceV2 := RamServiceV2{client}
		if err := ramServiceV2.SetResourceTags(d, "role"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudRamRoleRead(d, meta)
}

func resourceAliCloudRamRoleDelete(d *schema.ResourceData, meta interface{}) error {

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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", listAction, nil, listRequest, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(listAction, response, listRequest)

		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) || NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		policyResp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies.Policy", response)
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

				err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
					response, err = client.RpcPost("Ram", "2015-05-01", policyAction, nil, policyRequest, true)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, policyRequest)

				if err != nil {
					if IsExpectedErrors(err, []string{"EntityNotExist"}) || NotFoundError(err) {
						return nil
					}
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}

		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"DeleteConflict.Role.Policy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
