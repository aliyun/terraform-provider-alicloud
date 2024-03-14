package alicloud

import (
	"time"

	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCSKubernetesPermissions() *schema.Resource {
	return &schema.Resource{
		Read: dataAlicloudCSKubernetesPermissionsRead,

		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_owner": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_ram_role": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataAlicloudCSKubernetesPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	// Query existing permissions, DescribeUserPermission
	var perms []*cs.DescribeUserPermissionResponseBody
	uid := d.Get("uid").(string)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		perms, err = describeUserPermissions(client, uid)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "DescribeUserPermission", err)
	}

	_ = d.Set("permissions", flattenPermissionsConfig(perms))
	_ = d.Set("uid", uid)

	d.SetId(tea.ToString(hashcode.String(uid)))
	return nil
}

func describeUserPermissions(client *cs.Client, uid string) ([]*cs.DescribeUserPermissionResponseBody, error) {
	resp, err := client.DescribeUserPermission(tea.String(uid))
	if err != nil {
		return nil, err
	}

	addDebug("DescribeUserPermission", resp)
	return resp.Body, nil
}

func flattenPermissionsConfig(permissions []*cs.DescribeUserPermissionResponseBody) (m []map[string]interface{}) {
	if permissions == nil {
		return []map[string]interface{}{}
	}
	for _, permission := range permissions {
		m = append(m, map[string]interface{}{
			"resource_id":   permission.ResourceId,
			"resource_type": permission.ResourceType,
			"role_name":     permission.RoleName,
			"role_type":     permission.RoleType,
			"is_owner":      convertToBool(permission.IsOwner),
			"is_ram_role":   convertToBool(permission.IsRamRole),
		})
	}

	return m
}

func convertToBool(i *int64) bool {
	in := tea.Int64Value(i)
	if in != 1 {
		return false
	}

	return true
}
