package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"time"
)

func resourceAlicloudRamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamPolicyCreate,
		Read:   resourceAlicloudRamPolicyRead,
		Update: resourceAlicloudRamPolicyUpdate,
		Delete: resourceAlicloudRamPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"statement": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := Effect(v.(string))
								if value != Allow && value != Deny {
									errors = append(errors, fmt.Errorf(
										"%q must be '%s' or '%s'.", k, Allow, Deny))
								}
								return
							},
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
			"document": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"statement", "version"},
				ValidateFunc:  validateJsonString,
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
			"version": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "1",
				ConflictsWith: []string{"document"},
				ValidateFunc:  validatePolicyDocVersion,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request, err := buildAlicloudRamPolicyCreateArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreatePolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_policy", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ram.CreatePolicyResponse)
	d.SetId(response.Policy.PolicyName)
	return resourceAlicloudRamPolicyRead(d, meta)
}

func resourceAlicloudRamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	request, err := buildAlicloudRamPolicyUpdateArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	//check the quantity of version ,reserved 5 at most ,remove oldest version
	err = ramPolicyPruneVersions(d.Id(), "Custom", meta)
	if err != nil {
		return WrapError(err)
	}
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreatePolicyVersion(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	d.Partial(false)

	return resourceAlicloudRamPolicyRead(d, meta)
}

func resourceAlicloudRamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	object, err := ramService.DescribeRamPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	policy := object.Policy

	getPolicyRequest := ram.CreateGetPolicyVersionRequest()
	getPolicyRequest.VersionId = policy.DefaultVersion
	getPolicyRequest.PolicyType = policy.PolicyType
	getPolicyRequest.PolicyName = policy.PolicyName

	var statement []map[string]interface{}
	var version string
	var policyVersionResp *ram.GetPolicyVersionResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetPolicyVersion(getPolicyRequest)
		})
		if err != nil {
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), getPolicyRequest.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		policyVersionResp, _ = raw.(*ram.GetPolicyVersionResponse)
		statement, version, err = ramService.ParsePolicyDocument(policyVersionResp.PolicyVersion.PolicyDocument)
		if err != nil {
			if IsExceptedError(err, RoleAttachmentUnExpectedJson) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(WrapError(err))
		}
		addDebug(getPolicyRequest.GetActionName(), raw)
		return nil
	})
	if err != nil {
		return err
	}

	d.Set("name", policy.PolicyName)
	d.Set("type", policy.PolicyType)
	d.Set("description", policy.Description)
	d.Set("attachment_count", policy.AttachmentCount)
	d.Set("version", version)
	d.Set("statement", statement)
	d.Set("document", policyVersionResp.PolicyVersion.PolicyDocument)

	return nil
}

func resourceAlicloudRamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	listEntitiesForPolicyRequest := ram.CreateListEntitiesForPolicyRequest()
	listEntitiesForPolicyRequest.PolicyName = d.Id()

	if d.Get("force").(bool) {
		listEntitiesForPolicyRequest.PolicyType = "Custom"

		// list and detach entities for this policy
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListEntitiesForPolicy(listEntitiesForPolicyRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), listEntitiesForPolicyRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(listEntitiesForPolicyRequest.GetActionName(), raw)
		response := raw.(*ram.ListEntitiesForPolicyResponse)

		if len(response.Users.User) > 0 {
			for _, v := range response.Users.User {
				request := ram.CreateDetachPolicyFromUserRequest()
				request.UserName = v.UserName
				request.PolicyName = d.Id()
				request.PolicyType = "Custom"
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromUser(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw)
			}
		}

		if len(response.Groups.Group) > 0 {
			for _, v := range response.Groups.Group {
				request := ram.CreateDetachPolicyFromGroupRequest()
				request.GroupName = v.GroupName
				request.PolicyName = d.Id()
				request.PolicyType = "Custom"
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromGroup(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw)
			}
		}

		if len(response.Roles.Role) > 0 {
			for _, v := range response.Roles.Role {
				request := ram.CreateDetachPolicyFromRoleRequest()
				request.RoleName = v.RoleName
				request.PolicyName = d.Id()
				request.PolicyType = "Custom"
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromRole(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw)
			}
		}

		// list and delete policy version which are not default
		versions, err := ramPolicyListVersions(d.Id(), "Custom", meta)
		if len(versions) > 1 {
			for _, v := range versions {
				if !v.IsDefaultVersion {
					err = ramPolicyDeleteVersion(v.VersionId, d.Id(), meta)
					if err != nil {
						return WrapError(err)
					}
				}
			}
		}
	}

	deletePolicyRequest := ram.CreateDeletePolicyRequest()
	deletePolicyRequest.PolicyName = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeletePolicy(deletePolicyRequest)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{DeleteConflictPolicyUser, DeleteConflictPolicyGroup, DeleteConflictRolePolicy, DeleteConflictPolicyVersion}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(deletePolicyRequest.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), deletePolicyRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ramService.WaitForRamPolicy(d.Id(), Deleted, DefaultTimeout))
}

func buildAlicloudRamPolicyCreateArgs(d *schema.ResourceData, meta interface{}) (*ram.CreatePolicyRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	var document string

	doc, docOk := d.GetOk("document")
	statement, statementOk := d.GetOk("statement")

	if !docOk && !statementOk {
		return &ram.CreatePolicyRequest{}, WrapError(Error("One of 'document' and 'statement' must be specified."))
	}

	if docOk {
		document = doc.(string)
	} else {
		doc, err := ramService.AssemblePolicyDocument(statement.(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return &ram.CreatePolicyRequest{}, WrapError(err)
		}
		document = doc
	}

	request := ram.CreateCreatePolicyRequest()
	request.PolicyDocument = document
	request.PolicyName = d.Get("name").(string)

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}

	return request, nil
}

func buildAlicloudRamPolicyUpdateArgs(d *schema.ResourceData, meta interface{}) (*ram.CreatePolicyVersionRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	request := ram.CreateCreatePolicyVersionRequest()
	request.SetAsDefault = "true"
	request.PolicyName = d.Id()

	if d.HasChange("document") {
		d.SetPartial("document")
		request.PolicyDocument = d.Get("document").(string)

	} else if d.HasChange("statement") || d.HasChange("version") {

		if d.HasChange("statement") {
			d.SetPartial("statement")
		}
		if d.HasChange("version") {
			d.SetPartial("version")
		}

		document, err := ramService.AssemblePolicyDocument(d.Get("statement").(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return &ram.CreatePolicyVersionRequest{}, err
		}
		request.PolicyDocument = document
	}

	return request, nil
}

func ramPolicyPruneVersions(policyName, policyType string, meta interface{}) error {
	versions, err := ramPolicyListVersions(policyName, policyType, meta)
	if err != nil {
		return WrapError(err)
	}
	if len(versions) < 5 {
		return nil
	}
	var oldestVersion ram.PolicyVersion

	for _, version := range versions {
		if version.IsDefaultVersion {
			continue
		}
		if oldestVersion.CreateDate == "" ||
			version.CreateDate < oldestVersion.CreateDate {
			oldestVersion = version
		}
	}
	return ramPolicyDeleteVersion(oldestVersion.VersionId, policyName, meta)
}

func ramPolicyDeleteVersion(versionId, policyName string, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ram.CreateDeletePolicyVersionRequest()
	request.VersionId = versionId
	request.PolicyName = policyName
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.DeletePolicyVersion(request)
	})
	if err != nil && !RamEntityNotExist(err) {
		return WrapErrorf(err, DefaultErrorMsg, policyName, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return nil
}

func ramPolicyListVersions(policyName, policyType string, meta interface{}) ([]ram.PolicyVersion, error) {
	client := meta.(*connectivity.AliyunClient)
	request := ram.CreateListPolicyVersionsRequest()
	request.PolicyName = policyName
	request.PolicyType = policyType
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPolicyVersions(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, policyName, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ram.ListPolicyVersionsResponse)

	return response.PolicyVersions.PolicyVersion, nil
}
