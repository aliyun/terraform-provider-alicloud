package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudHBRVault_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_vault.default"
	ra := resourceAttrInit(resourceId, AlicloudHBRVaultMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrVault")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrvault%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRVaultBasicDependence0)
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
					"vault_type":          "STANDARD",
					"vault_storage_class": "STANDARD",
					"vault_name":          name,
					"description":         "接入测试描述",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_type":          "STANDARD",
						"vault_storage_class": "STANDARD",
						"vault_name":          name,
						"description":         "接入测试描述",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vault_name": name + "_update1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_name": name + "_update1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "接入测试描述1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "接入测试描述1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vault_name":  name + "_update2",
					"description": "接入测试描述2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_name":  name + "_update2",
						"description": "接入测试描述2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudHBRVaultMap0 = map[string]string{
	"status":              CHECKSET,
	"vault_type":          "STANDARD",
	"vault_storage_class": "STANDARD",
	"vault_name":          CHECKSET,
	"description":         "接入测试描述",
}

func AlicloudHBRVaultBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
