package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamPolicyCreate,
		Read:   resourceAliCloudRamPolicyRead,
		Update: resourceAliCloudRamPolicyUpdate,
		Delete: resourceAliCloudRamPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(26 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"policy_document": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validatePolicyDocument(true),
				ConflictsWith: []string{"document", "version", "statement"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rotate_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"None", "DeleteOldestNonDefaultVersionWhenLimitExceeded"}, false),
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
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
			"attachment_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"policy_name"},
				Deprecated:    "Field `name` has been deprecated from provider version 1.114.0. New field `policy_name` instead.",
			},
			"document": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validatePolicyDocument(true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				ConflictsWith: []string{"policy_document", "version", "statement"},
				Deprecated:    "Field `document` has been deprecated from provider version 1.114.0. New field `policy_document` instead.",
			},
			"version": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "1",
				ConflictsWith: []string{"document"},
				ValidateFunc:  StringInSlice([]string{"1"}, false),
				Deprecated:    "Field `version` has been deprecated from version 1.49.0, and use field `document` to replace.",
			},
			"statement": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"document"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"Allow", "Deny"}, false),
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
				Deprecated: "Field `statement` has been deprecated from version 1.49.0, and use field `document` to replace.",
			},
		},
	}
}

func resourceAliCloudRamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePolicy"
	request := make(map[string]interface{})
	conn, err := client.NewRamClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("policy_name"); ok {
		request["PolicyName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["PolicyName"] = v
	} else {
		return WrapError(Error("One of `policy_name`, `name` must be specified."))

	}

	if v, ok := d.GetOk("policy_document"); ok {
		request["PolicyDocument"] = v
	} else if v, ok := d.GetOk("document"); ok {
		request["PolicyDocument"] = v
	} else if v, ok := d.GetOk("statement"); ok {
		ramService := RamService{client}
		doc, err := ramService.AssemblePolicyDocument(v.(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return WrapError(err)
		}

		request["PolicyDocument"] = doc
	} else {
		return WrapError(Error("One of `policy_document`, `document`, `statement` must be specified."))
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_policy", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.Policy", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ram_policy")
	} else {
		policyName := resp.(map[string]interface{})["PolicyName"]
		d.SetId(fmt.Sprint(policyName))
	}

	return resourceAliCloudRamPolicyRead(d, meta)
}

func resourceAliCloudRamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	object, err := ramService.DescribeRamPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_policy ramService.DescribeRamPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if policy, ok := object["Policy"]; ok {
		policyArg := policy.(map[string]interface{})

		d.Set("policy_name", policyArg["PolicyName"])
		d.Set("description", policyArg["Description"])
		d.Set("type", policyArg["PolicyType"])
		d.Set("default_version", policyArg["DefaultVersion"])
		d.Set("attachment_count", policyArg["AttachmentCount"])
		d.Set("name", policyArg["PolicyName"])
	}

	if defaultPolicyVersion, ok := object["DefaultPolicyVersion"]; ok {
		defaultPolicyVersionArg := defaultPolicyVersion.(map[string]interface{})

		d.Set("policy_document", defaultPolicyVersionArg["PolicyDocument"])
		d.Set("version_id", defaultPolicyVersionArg["VersionId"])
		d.Set("document", defaultPolicyVersionArg["PolicyDocument"])

		statement, version, err := ramService.ParsePolicyDocument(defaultPolicyVersionArg["PolicyDocument"].(string))
		if err != nil {
			return WrapError(err)
		}

		d.Set("version", version)
		d.Set("statement", statement)
	}

	return nil
}

func resourceAliCloudRamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"SetAsDefault": true,
		"PolicyName":   d.Id(),
	}

	if d.HasChange("policy_document") {
		update = true

		if v, ok := d.GetOk("policy_document"); ok {
			request["PolicyDocument"] = v
		}
	}

	if d.HasChange("document") {
		update = true

		if v, ok := d.GetOk("document"); ok {
			request["PolicyDocument"] = v
		}
	}

	if d.HasChange("version") || d.HasChange("statement") {
		ramService := RamService{client}
		document, err := ramService.AssemblePolicyDocument(d.Get("statement").(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return WrapError(err)
		}

		request["PolicyDocument"] = document
	}

	if v, ok := d.GetOk("rotate_strategy"); ok {
		request["RotateStrategy"] = v
	}

	if update {
		action := "CreatePolicyVersion"
		conn, err := client.NewRamClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &runtime)
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

	return resourceAliCloudRamPolicyRead(d, meta)
}

func resourceAliCloudRamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicy"
	var response map[string]interface{}

	conn, err := client.NewRamClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"PolicyName": d.Id(),
	}

	if d.Get("force").(bool) {
		listRequest := map[string]interface{}{
			"PolicyName": d.Id(),
			"PolicyType": "Custom",
		}

		listAction := "ListEntitiesForPolicy"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(listAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, listRequest, &runtime)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		userResp, err := jsonpath.Get("$.Users.User", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Users.User", response)
		}

		if userResp != nil && len(userResp.([]interface{})) > 0 {
			for _, v := range userResp.([]interface{}) {
				userAction := "DetachPolicyFromUser"
				userRequest := map[string]interface{}{
					"PolicyName": d.Id(),
					"UserName":   v.(map[string]interface{})["UserName"],
					"PolicyType": "Custom",
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(userAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, userRequest, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, userRequest)

				if err != nil {
					if IsExpectedErrors(err, []string{"EntityNotExist"}) {
						return nil
					}
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		groupResp, err := jsonpath.Get("$.Groups.Group", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Groups.Group", response)
		}

		if groupResp != nil && len(groupResp.([]interface{})) > 0 {
			for _, v := range groupResp.([]interface{}) {
				groupAction := "DetachPolicyFromGroup"
				groupRequest := map[string]interface{}{
					"PolicyName": d.Id(),
					"GroupName":  v.(map[string]interface{})["GroupName"],
					"PolicyType": "Custom",
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(groupAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, groupRequest, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, groupRequest)

				if err != nil {
					if IsExpectedErrors(err, []string{"EntityNotExist"}) {
						return nil
					}
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		roleResp, err := jsonpath.Get("$.Roles.Role", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Roles.Role", response)
		}
		if roleResp != nil && len(roleResp.([]interface{})) > 0 {
			for _, v := range roleResp.([]interface{}) {
				roleAction := "DetachPolicyFromRole"
				roleRequest := map[string]interface{}{
					"PolicyName": d.Id(),
					"RoleName":   v.(map[string]interface{})["RoleName"],
					"PolicyType": "Custom",
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(roleAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, roleRequest, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, roleRequest)

				if err != nil {
					if IsExpectedErrors(err, []string{"EntityNotExist"}) {
						return nil
					}
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		listVersionsRequest := map[string]interface{}{
			"PolicyName": d.Id(),
			"PolicyType": "Custom",
		}
		listVersionsAction := "ListPolicyVersions"
		runtime = util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(listVersionsAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, listVersionsRequest, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, listVersionsRequest)

		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		versionsResp, er := jsonpath.Get("$.PolicyVersions.PolicyVersion", response)
		if er != nil {
			return WrapErrorf(er, FailedGetAttributeMsg, action, "$.PolicyVersions.PolicyVersion", response)
		}
		// More than one means there are other versions besides the default version
		if versionsResp != nil && len(versionsResp.([]interface{})) > 1 {
			for _, v := range versionsResp.([]interface{}) {
				if !v.(map[string]interface{})["IsDefaultVersion"].(bool) {
					versionAction := "DeletePolicyVersion"
					versionRequest := map[string]interface{}{
						"PolicyName": d.Id(),
						"VersionId":  v.(map[string]interface{})["VersionId"],
					}

					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(versionAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, versionRequest, &util.RuntimeOptions{})
						if err != nil {
							if NeedRetry(err) {
								wait()
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						return nil
					})
					addDebug(versionAction, response, versionRequest)

					if err != nil {
						return WrapErrorf(err, DefaultErrorMsg, d.Id(), versionAction, AlibabaCloudSdkGoERROR)
					}
				}
			}
		}
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"DeleteConflict.Policy.Group", "DeleteConflict.Policy.User", "DeleteConflict.Policy.Version", "DeleteConflict.Role.Policy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
