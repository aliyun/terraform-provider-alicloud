package dmsenterprise_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDMSEnterpriseUsersDataSource(t *testing.T) {
	tfacctest.CheckoutSupportedRegions(t, true, connectivity.DMSEnterpriseSupportRegions)
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_dms_enterprise_users.default"
	name := fmt.Sprintf("tf_testAccDmsEnterpriseUsersDataSource_%d", rand)
	testAccConfig := tfacctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceDmsEnterpriseUsersConfigDependence)

	idsConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dms_enterprise_user.default.uid}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dms_enterprise_user.default.uid}-fake"},
		}),
	}
	nameRegexConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dms_enterprise_user.default.user_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dms_enterprise_user.default.user_name}-fake",
		}),
	}
	statusConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_dms_enterprise_user.default.uid}"},
			"status": "NORMAL",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_dms_enterprise_user.default.uid}"},
			"status": "DISABLE",
		}),
	}

	roleConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alicloud_dms_enterprise_user.default.uid}"},
			"role": "DBA",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alicloud_dms_enterprise_user.default.uid}"},
			"role": "USER",
		}),
	}

	var existDmsEnterpriseUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"ids.0":                tfacctest.CHECKSET,
			"users.#":              "1",
			"users.0.nick_name":    name,
			"users.0.parent_uid":   tfacctest.CHECKSET,
			"users.0.role_ids.#":   "1",
			"users.0.role_names.#": "1",
			"users.0.status":       "NORMAL",
			"users.0.id":           tfacctest.CHECKSET,
			"users.0.user_id":      tfacctest.CHECKSET,
		}
	}

	var fakeDmsEnterpriseUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"users.#": "0",
		}
	}

	var kmsKeysCheckInfo = tfacctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existDmsEnterpriseUsersMapFunc,
		FakeMapFunc:  fakeDmsEnterpriseUsersMapFunc,
	}

	kmsKeysCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, roleConf)
}

func dataSourceDmsEnterpriseUsersConfigDependence(name string) string {
	return fmt.Sprintf(`                                            
		resource "alicloud_ram_user" "user" {                           
		  name         = "%[1]s"                                           
		  display_name = "user_display_name"                            
		  mobile       = "86-18688888888"                               
		  email        = "hello.uuu@aaa.com"                            
		  comments     = "yoyoyo"                                       
		}
		resource "alicloud_dms_enterprise_user" "default" {
		  uid = alicloud_ram_user.user.id
		  user_name = "%[1]s"
		  mobile = "15910799999"
		  role_names = ["DBA"]
	}`, name)
}
