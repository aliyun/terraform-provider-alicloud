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

func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		Create: ResourcePolicyCreate,
		Read:   ResourcePolicyRead,
		Update: ResourcePolicyUpdate,
		Delete: ResourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(26 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_document": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				ConflictsWith: []string{"document", "version", "statement"},
			},
			"policy_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"rotate_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"None", "DeleteOldestNonDefaultVersionWhenLimitExceeded"}, false),
			},
			"tags": tags.TagsSchema(),
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.114.0. New field 'policy_name' instead.",
				ConflictsWith: []string{"policy_name"},
			},
			"version": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "1",
				ConflictsWith: []string{"document"},
				// can only be '1' so far.
				ValidateFunc: validators.StringInSlice([]string{"1"}, false),
				Deprecated:   "Field 'version' has been deprecated from version 1.49.0, and use field 'document' to replace. ",
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
				Deprecated:    "Field 'document' has been deprecated from provider version 1.114.0. New field 'policy_document' instead.",
				ConflictsWith: []string{"policy_document", "version", "statement"},
			},
			"statement": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'statement' has been deprecated from version 1.49.0, and use field 'document' to replace. ",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validators.StringInSlice([]string{"Allow", "Deny"}, false),
						},
						"action": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				ConflictsWith: []string{"document"},
			},
		},
	}
}

func ResourcePolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("policy_name"); ok {
		request["PolicyName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["PolicyName"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("policy_document"); ok {
		request["PolicyDocument"] = v
	} else if v, ok := d.GetOk("document"); ok {
		request["PolicyDocument"] = v
	} else if v, ok := d.GetOk("statement"); ok {
		ramServiceV2 := RamServiceV2{client}
		doc, err := ramServiceV2.AssemblePolicyDocument(v.(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return errs.WrapError(err)
		}
		request["PolicyDocument"] = doc
	} else {
		return errs.WrapError(errs.Error("One of 'policy_document', 'document', 'statement'  must be specified."))

	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		tagsMapJSON, err := convert.ConvertListMapToJsonString(tagsMap)
		if err != nil {
			return errs.WrapError(err)
		}
		request["Tag"] = tagsMapJSON
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ram_policy", action, errs.AlibabaCloudSdkGoERROR)
	}

	responsePolicy := response["Policy"].(map[string]interface{})
	d.SetId(fmt.Sprint(responsePolicy["PolicyName"]))

	return ResourcePolicyRead(d, meta)
}

func ResourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	objectRaw, err := ramServiceV2.DescribeRamPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_policy DescribeRamPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	defaultPolicyVersionRawObj, _ := jsonpath.Get("$.DefaultPolicyVersion", objectRaw)
	defaultPolicyVersionRaw := make(map[string]interface{})
	if defaultPolicyVersionRawObj != nil {
		defaultPolicyVersionRaw = defaultPolicyVersionRawObj.(map[string]interface{})
	}
	d.Set("policy_document", defaultPolicyVersionRaw["PolicyDocument"])
	d.Set("document", defaultPolicyVersionRaw["PolicyDocument"])
	d.Set("version_id", defaultPolicyVersionRaw["VersionId"])

	statement, version, err := ramServiceV2.ParsePolicyDocument(defaultPolicyVersionRaw["PolicyDocument"].(string))
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("version", version)
	d.Set("statement", statement)

	policyRawObj, _ := jsonpath.Get("$.Policy", objectRaw)
	policyRaw := make(map[string]interface{})
	if policyRawObj != nil {
		policyRaw = policyRawObj.(map[string]interface{})
	}
	d.Set("create_time", policyRaw["CreateDate"])
	d.Set("description", policyRaw["Description"])
	d.Set("policy_name", policyRaw["PolicyName"])
	d.Set("name", policyRaw["PolicyName"])
	d.Set("default_version", policyRaw["DefaultVersion"])
	d.Set("attachment_count", policyRaw["AttachmentCount"])
	d.Set("type", policyRaw["PolicyType"])

	objectRaw, err = ramServiceV2.DescribePolicyListTagResources(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func ResourcePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "CreatePolicyVersion"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PolicyName"] = d.Id()

	request["SetAsDefault"] = true

	if d.HasChange("policy_document") {
		update = true
		request["PolicyDocument"] = d.Get("policy_document")
	}
	if d.HasChange("document") {
		update = true
		request["PolicyDocument"] = d.Get("document")
	}
	if d.HasChange("statement") || d.HasChange("version") {
		update = true
		ramServiceV2 := RamServiceV2{client}
		document, err := ramServiceV2.AssemblePolicyDocument(d.Get("statement").(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return errs.WrapError(err)
		}
		request["PolicyDocument"] = document
	}

	if v, ok := d.GetOk("rotate_strategy"); ok {
		request["RotateStrategy"] = v
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
	update = false
	action = "UpdatePolicyDescription"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PolicyName"] = d.Id()

	if d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
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
		if err := ramServiceV2.SetResourceTags(d, "policy"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return ResourcePolicyRead(d, meta)
}

func ResourcePolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["PolicyName"] = d.Id()

	request["CascadingDelete"] = true

	if d.Get("force").(bool) {
		listRequest := map[string]interface{}{
			"PolicyName": d.Id(),
			"PolicyType": "Custom",
		}
		listAction := "ListEntitiesForPolicy"
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
			if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) || errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), listAction, errs.AlibabaCloudSdkGoERROR)
		}

		listResponse := response
		userResp, err := jsonpath.Get("$.Users.User", listResponse)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, listAction, "$.Users.User", listResponse)
		}
		if userResp != nil && len(userResp.([]interface{})) > 0 {
			for _, v := range userResp.([]interface{}) {
				userAction := "DetachPolicyFromUser"
				userRequest := map[string]interface{}{
					"PolicyName": d.Id(),
					"UserName":   v.(map[string]interface{})["UserName"],
					"PolicyType": "Custom",
				}
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
					response, err = client.RpcPost("Ram", "2015-05-01", userAction, nil, userRequest, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(userAction, response, userRequest)

				if err != nil {
					if errs.IsExpectedErrors(err, []string{"EntityNotExist", "EntityNotExist.User", "EntityNotExist.User.Policy", "EntityNotExist.Policy"}) || errs.NotFoundError(err) {
						continue
					}
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), userAction, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}

		groupResp, err := jsonpath.Get("$.Groups.Group", listResponse)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, listAction, "$.Groups.Group", listResponse)
		}
		if groupResp != nil && len(groupResp.([]interface{})) > 0 {
			for _, v := range groupResp.([]interface{}) {
				groupAction := "DetachPolicyFromGroup"
				groupRequest := map[string]interface{}{
					"PolicyName": d.Id(),
					"GroupName":  v.(map[string]interface{})["GroupName"],
					"PolicyType": "Custom",
				}
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
					response, err = client.RpcPost("Ram", "2015-05-01", groupAction, nil, groupRequest, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(groupAction, response, groupRequest)

				if err != nil {
					if errs.IsExpectedErrors(err, []string{"EntityNotExist", "EntityNotExist.Group", "EntityNotExist.Group.Policy", "EntityNotExist.Policy"}) || errs.NotFoundError(err) {
						continue
					}
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), groupAction, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}

		roleResp, err := jsonpath.Get("$.Roles.Role", listResponse)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, listAction, "$.Roles.Role", listResponse)
		}
		if roleResp != nil && len(roleResp.([]interface{})) > 0 {
			for _, v := range roleResp.([]interface{}) {
				roleAction := "DetachPolicyFromRole"
				roleRequest := map[string]interface{}{
					"PolicyName": d.Id(),
					"RoleName":   v.(map[string]interface{})["RoleName"],
					"PolicyType": "Custom",
				}
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
					response, err = client.RpcPost("Ram", "2015-05-01", roleAction, nil, roleRequest, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(roleAction, response, roleRequest)

				if err != nil {
					if errs.IsExpectedErrors(err, []string{"EntityNotExist", "EntityNotExist.Role", "EntityNotExist.Role.Policy", "EntityNotExist.Policy"}) || errs.NotFoundError(err) {
						continue
					}
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), roleAction, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}

		listVersionsRequest := map[string]interface{}{
			"PolicyName": d.Id(),
			"PolicyType": "Custom",
		}
		listVersionsAction := "ListPolicyVersions"
		wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", listVersionsAction, nil, listVersionsRequest, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(listVersionsAction, response, listVersionsRequest)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) || errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), listVersionsAction, errs.AlibabaCloudSdkGoERROR)
		}

		versionsResp, er := jsonpath.Get("$.PolicyVersions.PolicyVersion", response)
		if er != nil {
			return errs.WrapErrorf(er, errs.FailedGetAttributeMsg, listVersionsAction, "$.PolicyVersions.PolicyVersion", response)
		}
		// More than one means there are other versions besides the default version
		if versionsResp != nil && len(versionsResp.([]interface{})) > 1 {
			for _, v := range versionsResp.([]interface{}) {
				isDefaultVersion, _ := v.(map[string]interface{})["IsDefaultVersion"].(bool)
				if !isDefaultVersion {
					versionAction := "DeletePolicyVersion"
					versionRequest := map[string]interface{}{
						"PolicyName": d.Id(),
						"VersionId":  v.(map[string]interface{})["VersionId"],
					}

					wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
					err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
						response, err = client.RpcPost("Ram", "2015-05-01", versionAction, nil, versionRequest, false)
						if err != nil {
							if helper.NeedRetry(err) {
								wait()
								return retry.RetryableError(err)
							}
							return retry.NonRetryableError(err)
						}
						return nil
					})
					logs.AddDebug(versionAction, response, versionRequest)

					if err != nil {
						if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) || errs.NotFoundError(err) {
							continue
						}
						return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), versionAction, errs.AlibabaCloudSdkGoERROR)
					}
				}
			}
		}
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"DeleteConflict.Policy.Group", "DeleteConflict.Policy.User", "DeleteConflict.Policy.Version", "DeleteConflict.Role.Policy"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
