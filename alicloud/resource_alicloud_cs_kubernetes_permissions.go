package alicloud

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	cs "github.com/alibabacloud-go/cs-20151215/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const ResourceName = "resource_alicloud_cs_kubernetes_permissions"

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
	uid := d.Get("uid").(string)
	// Grant Permissions
	if err := grantPermissions(d, meta, false); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "GrantPermissions", AlibabaCloudSdkGoERROR)
	}
	d.SetId(uid)
	return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
}

func resourceAlicloudCSKubernetesPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	d.Set("uid", d.Id())
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}
	describePerms, err := describeUserPermission(client, d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "DescribeUserPermission", AlibabaCloudSdkGoERROR)
	}
	var permissions []map[string]interface{}
	for _, perm := range describePerms {
		permission := map[string]interface{}{}
		if perm.IsRamRole != nil {
			permission["is_ram_role"] = tea.Int64Value(perm.IsRamRole) == 1
		}
		resourceType := tea.StringValue(perm.ResourceType)
		if tea.StringValue(perm.RoleType) == "custom" {
			permission["is_custom"] = true
			permission["role_name"] = tea.StringValue(perm.RoleName)
		} else {
			permission["role_name"] = tea.StringValue(perm.RoleType)
		}
		resourceId := tea.StringValue(perm.ResourceId)
		if strings.Contains(resourceId, "/") {
			parts := strings.Split(resourceId, "/")
			permission["cluster"] = parts[0]
			permission["namespace"] = parts[1]
			permission["role_type"] = "namespace"
		} else if resourceType == "cluster" {
			permission["cluster"] = resourceId
			permission["role_type"] = "cluster"
		}
		if resourceType == "console" && resourceId == "all-clusters" {
			permission["role_type"] = "all-clusters"
		}
		permissions = append(permissions, permission)
	}
	d.Set("permissions", permissions)
	return nil
}

func resourceAlicloudCSKubernetesPermissionsUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	// Update the permissions of the specified cluster, override
	if d.HasChange("permissions") {
		if err := grantPermissions(d, meta, false); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceName, "UpdatePermissions", AlibabaCloudSdkGoERROR)
		}
		d.Partial(false)
		return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
	}

	// Update all-clusters level permissions, if not exist, add new ones
	// TODO

	d.Partial(false)
	return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
}

func resourceAlicloudCSKubernetesPermissionsDelete(d *schema.ResourceData, meta interface{}) error {
	// Remove all permissions owned by the user
	err := grantPermissions(d, meta, true)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "RemoveAllPermissions", AlibabaCloudSdkGoERROR)
	}

	return nil
}

func buildPermissionArgs(d *schema.ResourceData) []*cs.GrantPermissionsRequestBody {
	var grantPermissionsRequest []*cs.GrantPermissionsRequestBody
	if perms, ok := d.GetOk("permissions"); ok {
		permissions := perms.(*schema.Set).List()
		var perms *cs.GrantPermissionsRequestBody
		for _, v := range permissions {
			pack := v.(map[string]interface{})
			perms = &cs.GrantPermissionsRequestBody{
				Cluster:   tea.String(pack["cluster"].(string)),
				RoleName:  tea.String(pack["role_name"].(string)),
				RoleType:  tea.String(pack["role_type"].(string)),
				Namespace: tea.String(pack["namespace"].(string)),
				IsCustom:  tea.Bool(pack["is_custom"].(bool)),
				IsRamRole: tea.Bool(pack["is_ram_role"].(bool)),
			}
			grantPermissionsRequest = append(grantPermissionsRequest, perms)
		}
	}

	return grantPermissionsRequest
}

func describeUserPermission(client *cs.Client, uid string) ([]*cs.DescribeUserPermissionResponseBody, error) {
	resp, err := client.DescribeUserPermission(tea.String(uid))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func grantPermissions(d *schema.ResourceData, meta interface{}, isDelete bool) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return err
	}

	uid := d.Get("uid").(string)
	body := buildPermissionArgs(d)
	if isDelete || body == nil {
		body = []*cs.GrantPermissionsRequestBody{}
	}
	req := &cs.GrantPermissionsRequest{
		Body: body,
	}
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.GrantPermissions(tea.String(uid), req)
		if err != nil {
			if NeedRetry(err) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug("GrantPermissions", body, err)
	if err != nil {
		return err
	}

	return nil
}

func parseClusterIds(perms []interface{}) []string {
	var clusters []string
	for _, v := range perms {
		m := v.(map[string]interface{})
		clusters = append(clusters, m["cluster"].(string))
	}
	return clusters
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
