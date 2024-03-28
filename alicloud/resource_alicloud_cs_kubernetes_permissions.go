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

//var modeApply = "apply"
var modePatch = "patch"
var modeDelete = "delete"

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
	if perms, ok := d.GetOk("permissions"); ok {
		return manageUserPermissions(modePatch, d, meta, perms.(*schema.Set).List())
	}

	return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
}

func resourceAlicloudCSKubernetesPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("uid").(string))
	d.Set("uid", d.Id())
	return nil
}

func resourceAlicloudCSKubernetesPermissionsUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	defer d.Partial(false)

	// Update the permissions of the specified cluster.
	// If other permissions of the cluster already exist, they will replace the existing permissions, and they will be added if they do not exist.
	// Keep other existing cluster permissions.
	if d.HasChange("permissions") {
		oldValue, newValue := d.GetChange("permissions")
		// delete old permissions
		err := manageUserPermissions(modeDelete, d, meta, oldValue.(*schema.Set).List())
		if err != nil {
			return err
		}
		// create new permissions
		return manageUserPermissions(modePatch, d, meta, newValue.(*schema.Set).List())
	}

	return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
}

func resourceAlicloudCSKubernetesPermissionsDelete(d *schema.ResourceData, meta interface{}) error {
	// Delete old permissions of cluster for the specified user.
	if perms, ok := d.GetOk("permissions"); ok {
		return manageUserPermissions(modeDelete, d, meta, perms.(*schema.Set).List())
	}

	return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
}

func manageUserPermissions(mode string, d *schema.ResourceData, meta interface{}, permissions []interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}
	uid := d.Get("uid").(string)

	// convert terraform permission resources to sdk request format
	updateUserPermissionsRequest := buildUpdateUserPermissionsArgs(permissions)
	updateUserPermissionsRequest.Mode = &mode
	addDebug("UpdateUserPermissions", updateUserPermissionsRequest)

	// call sdk update cluster permissions for user
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.UpdateUserPermissions(&uid, updateUserPermissionsRequest)
		if err == nil {
			return resource.NonRetryableError(err)
		}
		time.Sleep(5 * time.Second)
		return resource.RetryableError(Error("[ERROR] Update user permission failed %s error %v", d.Id(), err.Error()))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "UpdatePermissions", AliyunTablestoreGoSdk)
	}

	return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
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
