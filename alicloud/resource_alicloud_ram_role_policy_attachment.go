// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var ramRolePolicyAttachmentScopePattern = regexp.MustCompile(`^rg-[A-Za-z0-9]+$`)

func resourceAliCloudRamRolePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamRolePolicyAttachmentCreate,
		Read:   resourceAliCloudRamRolePolicyAttachmentRead,
		Delete: resourceAliCloudRamRolePolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"System", "Custom"}, false),
			},
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^$|^rg-[A-Za-z0-9]+$`), "must be empty for account scope or a resource group ID"),
			},
		},
	}
}

type ramRolePolicyAttachmentIdentity struct {
	policyName, policyType, roleName, scope string
}

func parseRamRolePolicyAttachmentID(id string) (ramRolePolicyAttachmentIdentity, error) {
	parts := strings.Split(id, ":")
	identity := ramRolePolicyAttachmentIdentity{}
	if (len(parts) != 4 && len(parts) != 5) || parts[0] != "role" {
		return identity, fmt.Errorf("invalid RAM role policy attachment ID %q: expected role:<policy_name>:<policy_type>:<role_name>[:<resource_group_id>]", id)
	}
	if parts[1] == "" || parts[3] == "" || (parts[2] != "Custom" && parts[2] != "System") {
		return identity, fmt.Errorf("invalid RAM role policy attachment identity in %q", id)
	}
	identity.policyName, identity.policyType, identity.roleName = parts[1], parts[2], parts[3]
	if len(parts) == 5 {
		if !ramRolePolicyAttachmentScopePattern.MatchString(parts[4]) {
			return identity, fmt.Errorf("invalid resource group scope in RAM role policy attachment ID %q", id)
		}
		identity.scope = parts[4]
	}
	return identity, nil
}

func (identity ramRolePolicyAttachmentIdentity) request() map[string]interface{} {
	request := map[string]interface{}{"PolicyName": identity.policyName, "PolicyType": identity.policyType, "RoleName": identity.roleName}
	if identity.scope != "" {
		request["ResourceGroupId"] = identity.scope
	}
	return request
}

// Only the scoped query's explicit absence result may clear state, not an arbitrary HTTP 404.
func ramRolePolicyAttachmentNotFound(err error) bool {
	wrapped, ok := err.(*ComplexError)
	return ok && wrapped.Err != nil && strings.HasPrefix(wrapped.Err.Error(), ResourceNotfound)
}

func setRamRolePolicyAttachmentData(d *schema.ResourceData, identity ramRolePolicyAttachmentIdentity) error {
	if err := d.Set("policy_name", identity.policyName); err != nil {
		return WrapError(err)
	}
	if err := d.Set("policy_type", identity.policyType); err != nil {
		return WrapError(err)
	}
	if err := d.Set("role_name", identity.roleName); err != nil {
		return WrapError(err)
	}
	if err := d.Set("resource_group_id", identity.scope); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAliCloudRamRolePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	id := strings.Join([]string{"role", d.Get("policy_name").(string), d.Get("policy_type").(string), d.Get("role_name").(string)}, ":")
	if scope := d.Get("resource_group_id").(string); scope != "" {
		id += ":" + scope
	}
	identity, err := parseRamRolePolicyAttachmentID(id)
	if err != nil {
		return WrapError(err)
	}
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	request := identity.request()
	action := "AttachPolicyToRole"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err := client.RpcPost("Ram", "2015-05-01", action, nil, request, true)
		addDebug(action, response, request)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"EntityNotExist.Role", "EntityNotExist.Policy", "EntityNotExists.ResourceGroup"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_role_policy_attachment", action, AlibabaCloudSdkGoERROR)
	}

	// Preserve the legacy four-part account ID and retain the identity if read-back fails.
	d.SetId(id)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := ramServiceV2.DescribeRamRolePolicyAttachment(id)
		if ramRolePolicyAttachmentNotFound(err) {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapError(err)
	}

	return setRamRolePolicyAttachmentData(d, identity)
}

func resourceAliCloudRamRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	identity, err := parseRamRolePolicyAttachmentID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}
	_, err = ramServiceV2.DescribeRamRolePolicyAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && ramRolePolicyAttachmentNotFound(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	return setRamRolePolicyAttachmentData(d, identity)
}

func resourceAliCloudRamRolePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	identity, err := parseRamRolePolicyAttachmentID(d.Id())
	if err != nil {
		return WrapError(err)
	}
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}
	request := identity.request()
	action := "DetachPolicyFromRole"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err := client.RpcPost("Ram", "2015-05-01", action, nil, request, true)
		addDebug(action, response, request)

		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExist.Role", "EntityNotExist.Role.Policy", "EntityNotExists.ResourceGroup"}) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := ramServiceV2.DescribeRamRolePolicyAttachment(d.Id())
		if ramRolePolicyAttachmentNotFound(err) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("RAM role policy attachment %s is still visible", d.Id()))
	})
}
