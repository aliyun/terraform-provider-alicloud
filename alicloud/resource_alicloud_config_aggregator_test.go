package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudConfigAggregator_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregator.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigAggregatorMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigAggregatorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckEnterpriseAccountEnabled(t)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_accounts": []map[string]interface{}{
						{
							"account_id":   `123968452689****`,
							"account_name": "tf-testacc1",
							"account_type": "ResourceDirectory",
						},
					},
					"aggregator_name": "${var.name}",
					"aggregator_type": "CUSTOM",
					"description":     "tf-create-aggregator",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_accounts.#": "1",
						"aggregator_name":       name,
						"aggregator_type":       "CUSTOM",
						"description":           "tf-create-aggregator",
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
					"aggregator_accounts": []map[string]interface{}{
						{
							"account_id":   "118258279412****",
							"account_name": "tf-testacc2",
							"account_type": "ResourceDirectory",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_accounts.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-modify-aggregator",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-modify-aggregator",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_accounts": []map[string]interface{}{
						{
							"account_id":   `123968452689****`,
							"account_name": "tf-testacc1",
							"account_type": "ResourceDirectory",
						},
					},
					"aggregator_name": "${var.name}",
					"description":     "tf-create-aggregator",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_accounts.#": "1",
						"aggregator_name":       name,
						"description":           "tf-create-aggregator",
					}),
				),
			},
		},
	})
}

var AlicloudConfigAggregatorMap0 = map[string]string{
	"aggregator_type": "CUSTOM",
}

func AlicloudConfigAggregatorBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}
`, name)
}
