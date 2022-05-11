package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDMSEnterpriseUsersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_dms_enterprise_users.default"
	name := fmt.Sprintf("tf_testAccDmsEnterpriseUsersDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDmsEnterpriseUsersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dms_enterprise_user.default.uid}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dms_enterprise_user.default.uid}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dms_enterprise_user.default.user_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dms_enterprise_user.default.user_name}-fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_dms_enterprise_user.default.uid}"},
			"status": "NORMAL",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_dms_enterprise_user.default.uid}"},
			"status": "DISABLE",
		}),
	}

	roleConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alicloud_dms_enterprise_user.default.uid}"},
			"role": "DBA",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alicloud_dms_enterprise_user.default.uid}"},
			"role": "USER",
		}),
	}

	var existDmsEnterpriseUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"ids.0":                CHECKSET,
			"users.#":              "1",
			"users.0.mobile":       "15910799999",
			"users.0.nick_name":    name,
			"users.0.parent_uid":   CHECKSET,
			"users.0.role_ids.#":   "1",
			"users.0.role_names.#": "1",
			"users.0.status":       "NORMAL",
			"users.0.id":           CHECKSET,
			"users.0.user_id":      CHECKSET,
		}
	}

	var fakeDmsEnterpriseUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"users.#": "0",
		}
	}

	var kmsKeysCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDmsEnterpriseUsersMapFunc,
		fakeMapFunc:  fakeDmsEnterpriseUsersMapFunc,
	}

	kmsKeysCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, roleConf)
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
