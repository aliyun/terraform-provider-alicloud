package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsSubscription_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_subscription.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsSubscriptionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsSubscription")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cms-subscription%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsSubscriptionBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"subscription_name":  "${var.name}",
					"description":        "${var.name}",
					"notify_strategy_id": "",
					"workspace":          "",
					"enable":             true,
					"filter_setting": []map[string]interface{}{
						{
							"conditions": []map[string]interface{}{
								{
									"field": "product",
									"value": "ECS",
									"op":    "EQ",
								},
							},
							"expression": "product",
							"relation":   "AND",
						},
					},
					"pushing_setting": []map[string]interface{}{
						{
							"alert_action_ids":   []interface{}{},
							"restore_action_ids": []interface{}{},
							"template_uuid":      "",
							"response_plan_id":   "",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"subscription_name": name,
						"description":       name,
						"enable":            "true",
						"filter_setting.#":  "1",
						"pushing_setting.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"subscription_name": "${var.name}",
					"description":       "${var.name}_update",
					"enable":            true,
					"filter_setting": []map[string]interface{}{
						{
							"conditions": []map[string]interface{}{
								{
									"field": "product",
									"value": "RDS",
									"op":    "IN",
								},
							},
							"expression": "product",
							"relation":   "OR",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      name + "_update",
						"filter_setting.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"subscription_name": "${var.name}",
					"description":       "${var.name}_update",
					"enable":            false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "update_time", "user_id", "region_id", "subscription_type", "subscription_id"},
			},
		},
	})
}

var AlicloudCmsSubscriptionMap = map[string]string{
	"subscription_name": CHECKSET,
	"description":       CHECKSET,
	"enable":            CHECKSET,
}

func AlicloudCmsSubscriptionBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
