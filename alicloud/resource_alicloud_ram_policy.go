package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"statement": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": &schema.Schema{
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
						"action": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource": &schema.Schema{
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
			"document": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"statement", "version"},
				ValidateFunc:  validateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamDesc,
			},
			"version": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "1",
				ConflictsWith: []string{"document"},
				ValidateFunc:  validatePolicyDocVersion,
			},
			"force": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args, err := buildAlicloudRamPolicyCreateArgs(d, meta)
	if err != nil {
		return err
	}

	response, err := conn.CreatePolicy(args)
	if err != nil {
		return fmt.Errorf("CreatePolicy got an error: %#v", err)
	}

	d.SetId(response.Policy.PolicyName)
	return resourceAlicloudRamPolicyUpdate(d, meta)
}

func resourceAlicloudRamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	d.Partial(true)

	args, attributeUpdate, err := buildAlicloudRamPolicyUpdateArgs(d, meta)
	if err != nil {
		return err
	}

	if !d.IsNewResource() && attributeUpdate {
		if _, err := conn.CreatePolicyVersion(args); err != nil {
			return fmt.Errorf("Error updating policy %s: %#v", d.Id(), err)
		}
	}

	d.Partial(false)

	return resourceAlicloudRamPolicyRead(d, meta)
}

func resourceAlicloudRamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		PolicyName: d.Id(),
		PolicyType: ram.Custom,
	}

	policyResp, err := conn.GetPolicy(args)
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
		}
		return fmt.Errorf("GetPolicy got an error: %#v", err)
	}
	policy := policyResp.Policy

	args.VersionId = policy.DefaultVersion
	policyVersionResp, err := conn.GetPolicyVersionNew(args)
	if err != nil {
		return fmt.Errorf("GetPolicyVersion got an error: %#v", err)
	}

	statement, version, err := ParsePolicyDocument(policyVersionResp.PolicyVersion.PolicyDocument)
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
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		PolicyName: d.Id(),
	}

	if d.Get("force").(bool) {
		args.PolicyType = ram.Custom

		// list and detach entities for this policy
		response, err := conn.ListEntitiesForPolicy(args)
		if err != nil {
			return fmt.Errorf("Error listing entities for policy %s when trying to delete: %#v", d.Id(), err)
		}

		if len(response.Users.User) > 0 {
			for _, v := range response.Users.User {
				_, err := conn.DetachPolicyFromUser(ram.AttachPolicyRequest{
					PolicyRequest: args,
					UserName:      v.UserName,
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error detaching policy %s from user %s:%#v", d.Id(), v.UserId, err)
				}
			}
		}

		if len(response.Groups.Group) > 0 {
			for _, v := range response.Groups.Group {
				_, err := conn.DetachPolicyFromGroup(ram.AttachPolicyToGroupRequest{
					PolicyRequest: args,
					GroupName:     v.GroupName,
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error detaching policy %s from group %s:%#v", d.Id(), v.GroupName, err)
				}
			}
		}

		if len(response.Roles.Role) > 0 {
			for _, v := range response.Roles.Role {
				_, err := conn.DetachPolicyFromRole(ram.AttachPolicyToRoleRequest{
					PolicyRequest: args,
					RoleName:      v.RoleName,
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error detaching policy %s from role %s:%#v", d.Id(), v.RoleId, err)
				}
			}
		}

		// list and delete policy version which are not default
		pvResp, err := conn.ListPolicyVersionsNew(args)
		if err != nil {
			return fmt.Errorf("Error listing policy versions for policy %s:%#v", d.Id(), err)
		}
		if len(pvResp.PolicyVersions.PolicyVersion) > 1 {
			for _, v := range pvResp.PolicyVersions.PolicyVersion {
				if !v.IsDefaultVersion {
					args.VersionId = v.VersionId
					if _, err = conn.DeletePolicyVersion(args); err != nil && !RamEntityNotExist(err) {
						return fmt.Errorf("Error delete policy version %s for policy %s:%#v", v.VersionId, d.Id(), err)
					}
				}
			}
		}
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DeletePolicy(args); err != nil {
			if IsExceptedError(err, DeleteConflictPolicyUser) || IsExceptedError(err, DeleteConflictPolicyGroup) || IsExceptedError(err, DeleteConflictRolePolicy) {
				return resource.RetryableError(fmt.Errorf("The policy can not been attached to any user or group or role while deleting the policy. - you can set force with true to force delete the policy."))
			}
			if IsExceptedError(err, DeleteConflictPolicyVersion) {
				return resource.RetryableError(fmt.Errorf("The policy can not has any version except the defaul version. - you can set force with true to force delete the policy."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting policy %s: %#v", d.Id(), err))
		}
		return nil
	})
}

func buildAlicloudRamPolicyCreateArgs(d *schema.ResourceData, meta interface{}) (ram.PolicyRequest, error) {
	var document string

	doc, docOk := d.GetOk("document")
	statement, statementOk := d.GetOk("statement")

	if !docOk && !statementOk {
		return ram.PolicyRequest{}, fmt.Errorf("One of 'document' and 'statement' must be specified.")
	}

	if docOk {
		document = doc.(string)
	} else {
		doc, err := AssemblePolicyDocument(statement.(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return ram.PolicyRequest{}, err
		}
		document = doc
	}

	args := ram.PolicyRequest{
		PolicyName:     d.Get("name").(string),
		PolicyDocument: document,
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.Description = v.(string)
	}

	return args, nil
}

func buildAlicloudRamPolicyUpdateArgs(d *schema.ResourceData, meta interface{}) (ram.PolicyRequest, bool, error) {
	args := ram.PolicyRequest{
		PolicyName:   d.Id(),
		SetAsDefault: "true",
	}

	attributeUpdate := false
	if d.HasChange("document") {
		d.SetPartial("document")
		attributeUpdate = true
		args.PolicyDocument = d.Get("document").(string)

	} else if d.HasChange("statement") || d.HasChange("version") {
		attributeUpdate = true

		if d.HasChange("statement") {
			d.SetPartial("statement")
		}
		if d.HasChange("version") {
			d.SetPartial("version")
		}

		document, err := AssemblePolicyDocument(d.Get("statement").(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return ram.PolicyRequest{}, attributeUpdate, err
		}
		args.PolicyDocument = document
	}

	return args, attributeUpdate, nil
}
