package alicloud

import (
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	cs "github.com/alibabacloud-go/cs-20151215/v4/client"
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
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	// Query existing permissions
	uid := d.Get("uid").(string)
	permissions := make([]interface{}, 0)
	if perms, ok := d.GetOk("permissions"); ok {
		permissions = perms.(*schema.Set).List()
	}
	// Grant Permissions
	// If other permissions with this right already exist, the existing permissions will be merged
	grantPermissionsRequest := buildPermissionsArgs(permissions)
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		err := grantPermissionsForAddPerm(client, uid, grantPermissionsRequest)
		if err == nil {
			return resource.NonRetryableError(err)
		}
		time.Sleep(5 * time.Second)
		return resource.RetryableError(Error("[ERROR] Grant user permission failed %s error %v", d.Id(), err.Error()))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "GrantPermissions", AliyunTablestoreGoSdk)
	}

	addDebug("GrantPermissions", grantPermissionsRequest, err)
	d.SetId(uid)
	return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
}

func resourceAlicloudCSKubernetesPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	d.Set("uid", d.Id())
	return nil
}

func resourceAlicloudCSKubernetesPermissionsUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	uid := d.Get("uid").(string)

	// Update the permissions of the specified cluster.
	// If other permissions of the cluster already exist, they will replace the existing permissions, and they will be added if they do not exist.
	// Keep other existing cluster permissions.
	if d.HasChange("permissions") {
		oldValue, newValue := d.GetChange("permissions")
		o := oldValue.(*schema.Set).List()
		n := newValue.(*schema.Set).List()
		oldPermissions := buildPermissionsArgs(o)
		newPermissions := buildPermissionsArgs(n)

		// Remove old permissions owned by the user:q:q
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			err := grantPermissionsForDeleteSomePerms(client, uid, oldPermissions)
			if err == nil {
				return resource.NonRetryableError(err)
			}
			time.Sleep(5 * time.Second)
			return resource.RetryableError(Error("[ERROR] Grant user permission failed %s error %v", d.Id(), err.Error()))
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceName, "GrantPermissions", AliyunTablestoreGoSdk)
		}

		// Add new permissions owned by the user
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			err := grantPermissionsForAddPerm(client, uid, newPermissions)
			if err == nil {
				return resource.NonRetryableError(err)
			}
			time.Sleep(5 * time.Second)
			return resource.RetryableError(Error("[ERROR] Grant user permission failed %s error %v", d.Id(), err.Error()))
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceName, "GrantPermissions", AliyunTablestoreGoSdk)
		}
	}
	d.Partial(false)
	return resourceAlicloudCSKubernetesPermissionsRead(d, meta)
}

func resourceAlicloudCSKubernetesPermissionsDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	uid := d.Id()
	permissions := make([]interface{}, 0)
	if perms, ok := d.GetOk("permissions"); ok {
		permissions = perms.(*schema.Set).List()
	}
	grantPermissionsRequest := buildPermissionsArgs(permissions)
	// Remove up some clusters permissions owned by the user
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		err := grantPermissionsForDeleteSomePerms(client, uid, grantPermissionsRequest)
		if err == nil {
			return resource.NonRetryableError(err)
		}
		time.Sleep(5 * time.Second)
		return resource.RetryableError(Error("[ERROR] Grant user permission failed %s error %v", d.Id(), err.Error()))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "GrantPermissions", AliyunTablestoreGoSdk)
	}
	return nil
}

func buildPermissionsArgs(permissions []interface{}) []*cs.GrantPermissionsRequestBody {
	grantPermissionsRequest := make([]*cs.GrantPermissionsRequestBody, 0)
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
	return grantPermissionsRequest
}

func convertDescribePermissionsToGrantPermissionsRequestBody(perms []*cs.DescribeUserPermissionResponseBody) []*cs.GrantPermissionsRequestBody {
	permReqs := make([]*cs.GrantPermissionsRequestBody, 0)
	for _, p := range perms {
		p := p
		req := &cs.GrantPermissionsRequestBody{
			Cluster:   tea.String(""),
			IsCustom:  tea.Bool(false),
			RoleName:  tea.String(""),
			RoleType:  tea.String("cluster"),
			Namespace: tea.String(""),
			IsRamRole: tea.Bool(false),
		}
		resourceId := ""
		resourceType := tea.StringValue(p.ResourceType)

		req.IsRamRole = tea.Bool(tea.Int64Value(p.IsRamRole) == 1)
		if tea.StringValue(p.RoleType) == "custom" {
			req.IsCustom = tea.Bool(true)
			req.RoleName = tea.String(tea.StringValue(p.RoleName))
		} else {
			req.RoleName = tea.String(tea.StringValue(p.RoleType))
		}
		resourceId = tea.StringValue(p.ResourceId)
		if strings.Contains(resourceId, "/") {
			parts := strings.Split(resourceId, "/")
			cluster := parts[0]
			namespace := parts[1]
			req.Cluster = tea.String(cluster)
			req.Namespace = tea.String(namespace)
			req.RoleType = tea.String("namespace")
		} else if resourceType == "cluster" {
			cluster := resourceId
			req.Cluster = tea.String(cluster)
			req.RoleType = tea.String("cluster")
		}
		if resourceType == "console" && resourceId == "all-clusters" {
			req.RoleType = tea.String("all-clusters")
		}

		permReqs = append(permReqs, req)
	}

	return permReqs
}

func describeUserPermission(client *cs.Client, uid string) ([]*cs.DescribeUserPermissionResponseBody, error) {
	resp, err := client.DescribeUserPermission(tea.String(uid))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func grantPermissionsForAddPerm(client *cs.Client, uid string, body []*cs.GrantPermissionsRequestBody) error {
	existPerms, err := describeUserPermission(client, uid)
	if err != nil {
		return err
	}
	perms := convertDescribePermissionsToGrantPermissionsRequestBody(existPerms)
	perms = append(perms, body...)
	req := &cs.GrantPermissionsRequest{
		Body: perms,
	}
	_, err = client.GrantPermissions(tea.String(uid), req)
	if err != nil {
		return err
	}
	return nil
}

func grantPermissionsForDeleteSomePerms(client *cs.Client, uid string, deletedPerms []*cs.GrantPermissionsRequestBody) error {
	describePerms, err := describeUserPermission(client, uid)
	if err != nil {
		return err
	}
	existPerms := convertDescribePermissionsToGrantPermissionsRequestBody(describePerms)
	newPerms := make([]*cs.GrantPermissionsRequestBody, 0)
	for _, existPerm := range existPerms {
		isDeleted := false
		for _, deletedPerm := range deletedPerms {
			if reflect.DeepEqual(existPerm, deletedPerm) {
				isDeleted = true
				break
			}
		}
		if !isDeleted {
			newPerms = append(newPerms, existPerm)
		}
	}
	req := &cs.GrantPermissionsRequest{
		Body: newPerms,
	}

	_, err = client.GrantPermissions(tea.String(uid), req)
	if err != nil {
		return err
	}

	return nil
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
