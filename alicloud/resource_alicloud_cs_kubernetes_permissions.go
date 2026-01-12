package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const ResourceName = "resource_alicloud_cs_kubernetes_permissions"

const (
	//	ModeApply       = "apply"
	ModePatch       = "patch"
	ModeDelete      = "delete"
	ConflictError   = "ErrAuthorizationConflict"
	ThrottlingError = "Throttling.User"
)

func resourceAlicloudCSKubernetesPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesPermissionsCreate,
		Read:   resourceAlicloudCSKubernetesPermissionsRead,
		Update: resourceAlicloudCSKubernetesPermissionsUpdate,
		Delete: resourceAlicloudCSKubernetesPermissionsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permissions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"cluster", "namespace", "all-clusters"}, false),
						},
						"role_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cluster": {
							Type:     schema.TypeString,
							Required: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_custom": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"is_ram_role": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudCSKubernetesPermissionsCreate(d *schema.ResourceData, meta interface{}) error {
	// Create new permissions of cluster for the specified user.
	var err error
	uid := d.Get("uid").(string)
	if perms, ok := d.GetOk("permissions"); ok {
		err = manageUserPermissions(ModePatch, uid, meta, perms.(*schema.Set).List())
	}

	if err != nil {
		return WrapError(err)
	}

	d.SetId(uid)
	return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
}

func resourceAlicloudCSKubernetesPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}
	// Query existing permissions, DescribeUserPermission
	uid := d.Id()

	var perms []*cs.DescribeUserPermissionResponseBody
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		perms, err = describeUserPermissions(client, uid)
		if isRetryforThrottling(err) {
			time.Sleep(1 * time.Minute)
		} else if tea.BoolValue(tea.Retryable(err)) {
			time.Sleep(5 * time.Second)
		} else {
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(Error("[ERROR] Describe user permission failed %s error %v", uid, err.Error()))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "DescribeUserPermission", err)
	}
	if len(perms) == 0 {
		err = d.Set("permissions", nil)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceName, "Read set permissions", err)
		}
	}
	err = d.Set("uid", uid)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "Read set uid", err)
	}

	return nil
}

func resourceAlicloudCSKubernetesPermissionsUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	defer d.Partial(false)

	// Update the permissions of the specified cluster.
	// If other permissions of the cluster already exist, they will replace the existing permissions, and they will be added if they do not exist.
	// Keep other existing cluster permissions.
	uid := d.Id()
	if d.HasChange("permissions") {
		oldPermissions, newPermissions := d.GetChange("permissions")
		// if resource no change, return
		oldPermissionsList := oldPermissions.(*schema.Set).List()
		newPermissionsList := newPermissions.(*schema.Set).List()

		if len(oldPermissionsList) > 0 && len(newPermissionsList) > 0 {
			oldPermissionsList, newPermissionsList = diffPermissions(oldPermissionsList, newPermissionsList)
		}

		// delete old permissions
		if len(oldPermissionsList) > 0 {
			err := manageUserPermissions(ModeDelete, uid, meta, oldPermissionsList)
			if err != nil {
				return WrapError(err)
			}
		}

		// create new permissions
		if len(newPermissionsList) > 0 {
			err := manageUserPermissions(ModePatch, uid, meta, newPermissionsList)
			if err != nil {
				return WrapError(err)
			}
		}
	}

	return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
}

func resourceAlicloudCSKubernetesPermissionsDelete(d *schema.ResourceData, meta interface{}) error {
	// Delete old permissions of cluster for the specified user.
	if perms, ok := d.GetOk("permissions"); ok {
		return manageUserPermissions(ModeDelete, d.Get("uid").(string), meta, perms.(*schema.Set).List())
	}

	return nil
}

func manageUserPermissions(mode, uid string, meta interface{}, permissions []interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	// convert terraform permission resources to sdk request format
	updateUserPermissionsRequest := buildUpdateUserPermissionsArgs(permissions)
	updateUserPermissionsRequest.Mode = &mode
	addDebug("UpdateUserPermissions", updateUserPermissionsRequest)

	// call sdk update cluster permissions for user
	err = resource.Retry(60*time.Minute, func() *resource.RetryError {
		_, err := client.UpdateUserPermissions(&uid, updateUserPermissionsRequest)
		if isRetryforThrottling(err) {
			time.Sleep(1 * time.Minute)
		} else if isRetryforConflict(err) || tea.BoolValue(tea.Retryable(err)) {
			time.Sleep(5 * time.Second)
		} else {
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(Error("[ERROR] Update user permission failed %s error %v", uid, err.Error()))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "UpdatePermissions", AliyunTablestoreGoSdk)
	}

	_ = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.DescribeUserPermission(&uid)
		if isRetryforThrottling(err) {
			time.Sleep(1 * time.Minute)
		} else if tea.BoolValue(tea.Retryable(err)) {
			time.Sleep(5 * time.Second)
		} else {
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(Error("[ERROR] Describe user permission failed %s error %v", uid, err.Error()))
	})

	return nil
}

func buildUpdateUserPermissionsArgs(permissions []interface{}) *cs.UpdateUserPermissionsRequest {
	updateUserPermissions := make([]*cs.UpdateUserPermissionsRequestBody, 0)
	var perms *cs.UpdateUserPermissionsRequestBody
	for _, v := range permissions {
		pack := v.(map[string]interface{})
		perms = &cs.UpdateUserPermissionsRequestBody{
			Cluster:   tea.String(pack["cluster"].(string)),
			RoleName:  tea.String(pack["role_name"].(string)),
			RoleType:  tea.String(pack["role_type"].(string)),
			Namespace: tea.String(pack["namespace"].(string)),
			IsCustom:  tea.Bool(pack["is_custom"].(bool)),
			IsRamRole: tea.Bool(pack["is_ram_role"].(bool)),
		}
		updateUserPermissions = append(updateUserPermissions, perms)
	}

	return &cs.UpdateUserPermissionsRequest{Body: updateUserPermissions}
}

func diffPermissions(oldPermissionList, newPermissionList []interface{}) ([]interface{}, []interface{}) {
	for i := 0; i < len(oldPermissionList); {
		oldMap := oldPermissionList[i].(map[string]interface{})
		i += 1
		for j := 0; j < len(newPermissionList); j++ {
			newMap := newPermissionList[j].(map[string]interface{})
			if oldMap["role_type"] == newMap["role_type"] && oldMap["role_name"] == newMap["role_name"] && oldMap["cluster"] == newMap["cluster"] && oldMap["namespace"] == newMap["namespace"] && oldMap["is_custom"] == newMap["is_custom"] && oldMap["is_ram_role"] == newMap["is_ram_role"] {
				i -= 1
				oldPermissionList = append(oldPermissionList[:i], oldPermissionList[i+1:]...)
				newPermissionList = append(newPermissionList[:j], newPermissionList[j+1:]...)
				// Jump out of the first cycle
				break
			}
		}
	}

	return oldPermissionList, newPermissionList
}

func isRetryforConflict(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && tea.StringValue(e.Code) == ConflictError {
		return true
	}

	return false
}

func isRetryforThrottling(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && tea.StringValue(e.Code) == ThrottlingError {
		return true
	}

	return false
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}
