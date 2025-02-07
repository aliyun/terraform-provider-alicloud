package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudServiceMeshUserPermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudServiceMeshUserPermissionCreate,
		Read:   resourceAliCloudServiceMeshUserPermissionRead,
		Update: resourceAliCloudServiceMeshUserPermissionUpdate,
		Delete: resourceAliCloudServiceMeshUserPermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"sub_account_user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_name": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"istio-admin", "istio-ops", "istio-readonly"}, false),
						},
						"service_mesh_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"role_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"custom"}, false),
							Computed:     true,
						},
						"is_custom": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_ram_role": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudServiceMeshUserPermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "GrantUserPermissions"
	request := make(map[string]interface{})
	var err error
	request["SubAccountUserId"] = d.Get("sub_account_user_id")
	userSlice := make([]interface{}, 0)
	if v, ok := d.GetOk("permissions"); ok {
		for _, raw := range v.(*schema.Set).List() {
			obj := make(map[string]interface{}, 0)
			rawMap := raw.(map[string]interface{})
			if v, ok := rawMap["role_name"]; ok {
				obj["RoleName"] = v
			}
			if v, ok := rawMap["service_mesh_id"]; ok {
				obj["Cluster"] = v
			}
			if v, ok := rawMap["role_type"]; ok {
				obj["RoleType"] = v
			}
			if v, ok := rawMap["is_ram_role"]; ok {
				obj["IsRamRole"] = v
			}
			if v, ok := rawMap["is_custom"]; ok {
				obj["IsCustom"] = v
			}
			userSlice = append(userSlice, obj)
		}
	}
	raw, err := json.Marshal(userSlice)
	if err != nil {
		return WrapError(err)
	}
	request["Permissions"] = string(raw)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidOperation.Grant.NotRunning"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_service_mesh_user_permission", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(d.Get("sub_account_user_id")))
	return resourceAliCloudServiceMeshUserPermissionRead(d, meta)
}
func resourceAliCloudServiceMeshUserPermissionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	servicemeshService := ServicemeshService{client}
	object, err := servicemeshService.DescribeUserPermissions(d.Id())
	if err != nil {
		return WrapError(err)
	}
	permissionSli := make([]interface{}, 0)
	for _, raw := range object["Permissions"].([]interface{}) {
		rawMap := raw.(map[string]interface{})
		obj := make(map[string]interface{}, 0)
		if v, ok := rawMap["ResourceId"]; ok {
			obj["service_mesh_id"] = v
		}
		if v, ok := rawMap["IsRamRole"]; ok {
			obj["is_ram_role"] = v
			obj["is_custom"] = !v.(bool)
		} else {
			obj["is_ram_role"] = false
			obj["is_custom"] = true
		}
		if v, ok := rawMap["RoleType"]; ok {
			obj["role_type"] = v
		}
		if v, ok := rawMap["RoleName"]; ok {
			obj["role_name"] = v
		}

		permissionSli = append(permissionSli, obj)
	}
	d.Set("permissions", permissionSli)
	d.Set("sub_account_user_id", d.Id())
	return nil
}
func resourceAliCloudServiceMeshUserPermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)
	var err error
	update := false
	request := map[string]interface{}{
		"SubAccountUserId": d.Id(),
	}
	if d.HasChange("permissions") {
		update = true
	}
	userSlice := make([]interface{}, 0)
	if v, ok := d.GetOk("permissions"); ok {
		for _, raw := range v.(*schema.Set).List() {
			obj := make(map[string]interface{}, 0)
			rawMap := raw.(map[string]interface{})
			if v, ok := rawMap["role_name"]; ok {
				obj["RoleName"] = v
			}
			if v, ok := rawMap["service_mesh_id"]; ok {
				obj["Cluster"] = v
			}
			if v, ok := rawMap["role_type"]; ok {
				obj["RoleType"] = v
			}
			if v, ok := rawMap["is_ram_role"]; ok {
				obj["IsRamRole"] = v
			}
			if v, ok := rawMap["is_custom"]; ok {
				obj["IsCustom"] = v
			}
			userSlice = append(userSlice, obj)
		}
	}
	raw, err := json.Marshal(userSlice)
	if err != nil {
		return WrapError(err)
	}
	request["Permissions"] = string(raw)
	if update {
		action := "GrantUserPermissions"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidOperation.Grant.NotRunning"}) {
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
	d.Partial(false)

	return resourceAliCloudServiceMeshUserPermissionRead(d, meta)
}
func resourceAliCloudServiceMeshUserPermissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)
	var err error

	request := map[string]interface{}{
		"SubAccountUserId": d.Id(),
	}
	userSlice := make([]interface{}, 0)
	raw, err := json.Marshal(userSlice)
	if err != nil {
		return WrapError(err)
	}
	request["Permissions"] = string(raw)

	action := "GrantUserPermissions"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidOperation.Grant.NotRunning"}) {
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
	return nil
}
