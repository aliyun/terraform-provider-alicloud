package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamRoleCreate,
		Read:   resourceAlicloudRamRoleRead,
		Update: resourceAlicloudRamRoleUpdate,
		Delete: resourceAlicloudRamRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"ram_users": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"document"},
			},
			"services": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"document"},
			},
			"document": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"ram_users", "services", "version"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				ValidateFunc: validateJsonString,
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
			"arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args, err := buildAlicloudRamRoleCreateArgs(d, meta)
	if err != nil {
		return err
	}

	response, err := conn.CreateRole(args)
	if err != nil {
		return fmt.Errorf("CreateRole got an error: %#v", err)
	}

	d.SetId(response.Role.RoleName)
	return resourceAlicloudRamRoleUpdate(d, meta)
}

func resourceAlicloudRamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	d.Partial(true)

	args, attributeUpdate, err := buildAlicloudRamRoleUpdateArgs(d, meta)
	if err != nil {
		return err
	}

	if !d.IsNewResource() && attributeUpdate {
		if _, err := conn.UpdateRole(args); err != nil {
			return fmt.Errorf("UpdateRole got an error: %v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamRoleRead(d, meta)
}

func resourceAlicloudRamRoleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.RoleQueryRequest{
		RoleName: d.Id(),
	}

	response, err := conn.GetRole(args)
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
		}
		return fmt.Errorf("GetRole got an error: %v", err)
	}

	role := response.Role
	rolePolicy, err := ParseRolePolicyDocument(role.AssumeRolePolicyDocument)
	if err != nil {
		return err
	}
	if len(rolePolicy.Statement) > 0 {
		principal := rolePolicy.Statement[0].Principal
		d.Set("services", principal.Service)
		d.Set("ram_users", principal.RAM)
	}

	d.Set("name", role.RoleName)
	d.Set("arn", role.Arn)
	d.Set("description", role.Description)
	d.Set("version", rolePolicy.Version)
	d.Set("document", role.AssumeRolePolicyDocument)
	return nil
}

func resourceAlicloudRamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.RoleQueryRequest{
		RoleName: d.Id(),
	}

	if d.Get("force").(bool) {
		resp, err := conn.ListPoliciesForRole(args)
		if err != nil {
			return fmt.Errorf("Error listing Policies for Role (%s) when trying to delete: %#v", d.Id(), err)
		}

		// Loop and remove the Policies from the Role
		if len(resp.Policies.Policy) > 0 {
			for _, v := range resp.Policies.Policy {
				_, err = conn.DetachPolicyFromRole(ram.AttachPolicyToRoleRequest{
					PolicyRequest: ram.PolicyRequest{
						PolicyName: v.PolicyName,
						PolicyType: ram.Type(v.PolicyType),
					},
					RoleName: d.Id(),
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error detach Policy from Role %s: %#v", d.Id(), err)
				}
			}
		}
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DeleteRole(args); err != nil {
			if IsExceptedError(err, DeleteConflictRolePolicy) {
				return resource.RetryableError(fmt.Errorf("The role can not has any attached policy while deleting the role. - you can set force with true to force delete the role."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting role %s: %#v, you can set force with true to force delete the role.", d.Id(), err))
		}
		return nil
	})
}

func buildAlicloudRamRoleCreateArgs(d *schema.ResourceData, meta interface{}) (ram.RoleRequest, error) {

	args := ram.RoleRequest{
		RoleName: d.Get("name").(string),
	}

	ramUsers, usersOk := d.GetOk("ram_users")
	services, servicesOk := d.GetOk("services")
	document, documentOk := d.GetOk("document")

	if !usersOk && !servicesOk && !documentOk {
		return ram.RoleRequest{}, fmt.Errorf("At least one of 'ram_users', 'services' or 'document' must be set.")
	}

	if documentOk {
		args.AssumeRolePolicyDocument = document.(string)
	} else {
		rolePolicyDocument, err := AssembleRolePolicyDocument(ramUsers.(*schema.Set).List(), services.(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return ram.RoleRequest{}, err
		}
		args.AssumeRolePolicyDocument = rolePolicyDocument
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.Description = v.(string)
	}

	return args, nil
}

func buildAlicloudRamRoleUpdateArgs(d *schema.ResourceData, meta interface{}) (ram.UpdateRoleRequest, bool, error) {
	args := ram.UpdateRoleRequest{
		RoleName: d.Id(),
	}

	attributeUpdate := false

	if d.HasChange("document") {
		d.SetPartial("document")
		attributeUpdate = true
		args.NewAssumeRolePolicyDocument = d.Get("document").(string)

	} else if d.HasChange("ram_users") || d.HasChange("services") || d.HasChange("version") {
		attributeUpdate = true

		if d.HasChange("ram_users") {
			d.SetPartial("ram_users")
		}
		if d.HasChange("services") {
			d.SetPartial("services")
		}
		if d.HasChange("version") {
			d.SetPartial("version")
		}

		document, err := AssembleRolePolicyDocument(d.Get("ram_users").(*schema.Set).List(), d.Get("services").(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return ram.UpdateRoleRequest{}, attributeUpdate, err
		}
		args.NewAssumeRolePolicyDocument = document
	}

	return args, attributeUpdate, nil
}
