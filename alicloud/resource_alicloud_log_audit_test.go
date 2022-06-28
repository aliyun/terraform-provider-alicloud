package alicloud

import (
	"fmt"
	"testing"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogAudit_basic(t *testing.T) {
	var v *slsPop.DescribeAppResponse
	resourceId := "alicloud_log_audit.foo"
	ra := resourceAttrInit(resourceId, logAuditMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogaudit-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogAuditConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
					"aliuid":       "${data.alicloud_account.default.id}",
					"variable_map": map[string]string{
						"actiontrail_enabled": "false",
						"actiontrail_ttl":     "10",
						"oss_access_enabled":  "true",
						"oss_access_ttl":      "155",
						"oss_sync_enabled":    "true",
						"oss_sync_ttl":        "180",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name":                     name,
						"aliuid":                           CHECKSET,
						"variable_map.%":                   "6",
						"variable_map.actiontrail_enabled": "false",
						"variable_map.actiontrail_ttl":     "10",
						"variable_map.oss_access_enabled":  "true",
						"variable_map.oss_access_ttl":      "155",
						"variable_map.oss_sync_enabled":    "true",
						"variable_map.oss_sync_ttl":        "180",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"multi_account": []string{"1234567", "123123123213", "123141412"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"multi_account.#": "3",
					}),
				),
			},
			// TODO: only when center account is resource directory master or resource directory admin need to check resource type config，otherwise pass it directly
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"resource_directory_type": "custom",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"resource_directory_type": "custom",
			// 		}),
			// 	),
			// },
		},
	})
}

func resourceLogAuditConfigDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_account" "default" {}
`)
}

var logAuditMap = map[string]string{
	"display_name": CHECKSET,
}
