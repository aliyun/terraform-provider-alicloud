package mongodb

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mongoDBAccountGovernanceSchema adds the account-level Dds APIs to the
// existing alicloud_mongodb_account type.
func mongoDBAccountGovernanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_roles": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"role":    {Type: schema.TypeString, Computed: true},
					"db_name": {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"privilege": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"role":    {Type: schema.TypeString, Required: true},
					"db_name": {Type: schema.TypeString, Required: true},
				},
			},
		},
		"privilege_reset_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "update",
			ValidateFunc: validators.StringInSlice([]string{"update", "grant", "revoke"}, false),
		},
		"force_destroy": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func mongoDBGovernanceReadAccount(d *schema.ResourceData, client *connectivity.AliyunClient) {
	d.Set("account_roles", []map[string]interface{}{})
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return
	}
	request := map[string]interface{}{"DBInstanceId": parts[0]}
	if response, err := mongoDBGovernanceRpcPost(client, 2*time.Minute, "DescribeAccountRoles", request); err != nil {
		log.Printf("[DEBUG] DescribeAccountRoles for %s: %s", d.Id(), err)
	} else {
		roles := make([]map[string]interface{}, 0)
		for _, item := range mongoDBGovernanceMapList(mongoDBGovernancePick(response, "Roles", "roles")) {
			roles = append(roles, map[string]interface{}{
				"role":    item["Role"],
				"db_name": item["DBName"],
			})
		}
		d.Set("account_roles", roles)
	}
}

func mongoDBGovernanceUpdateAccount(d *schema.ResourceData, client *connectivity.AliyunClient) error {
	privileges, _ := d.Get("privilege").([]interface{})
	if !d.HasChange("privilege") && !d.HasChange("privilege_reset_type") {
		return nil
	}
	if len(privileges) == 0 {
		return nil
	}
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return errs.WrapErrorf(errs.Error("invalid account id %s", d.Id()), errs.IdMsg, d.Id())
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"AccountName":  parts[1],
		"ResetType":    d.Get("privilege_reset_type").(string),
	}
	if v, ok := d.GetOk("character_type"); ok {
		request["CharacterType"] = v
	}
	for index, item := range privileges {
		privilege, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		request[fmt.Sprintf("Roles.%d.Role", index+1)] = privilege["role"]
		request[fmt.Sprintf("Roles.%d.DBName", index+1)] = privilege["db_name"]
	}
	if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutUpdate), "ModifyAccountPrivilege", request); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "ModifyAccountPrivilege", errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

func mongoDBGovernanceDeleteAccount(d *schema.ResourceData, client *connectivity.AliyunClient) error {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return errs.WrapErrorf(errs.Error("invalid account id %s", d.Id()), errs.IdMsg, d.Id())
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"AccountName":  parts[1],
	}
	if v, ok := d.GetOk("character_type"); ok {
		request["CharacterType"] = v
	}
	if _, err := mongoDBGovernanceRpcPost(client, d.Timeout(schema.TimeoutDelete), "DeleteAccount", request); err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteAccount", errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
