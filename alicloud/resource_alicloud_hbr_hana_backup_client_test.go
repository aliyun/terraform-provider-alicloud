package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func SkipTestAccAlicloudHbrHanaBackupClient_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_hbr_hana_backup_client.default"
	ra := resourceAttrInit(resourceId, AlicloudHbrHanaBackupClientMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrHanaBackupClient")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sHbrHanaBackupClient%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHbrHanaBackupClientBasicDependence0)
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
					"vault_id":      "${data.alicloud_hbr_vaults.default.vaults.0.id}",
					"client_info":   `[ { \"instanceId\": \"i-bp1dpl8hfbkh5rvvcmsg\", \"clusterId\": \"cl-000cnu7ti2rmj23dhp77\", \"sourceTypes\": [ \"HANA\" ]  }]`,
					"alert_setting": "INHERITED",
					"use_https":     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_id":      CHECKSET,
						"alert_setting": "INHERITED",
						"use_https":     "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_info"},
			},
		},
	})
}

var AlicloudHbrHanaBackupClientMap = map[string]string{
	"client_id":   CHECKSET,
	"instance_id": CHECKSET,
	"cluster_id":  CHECKSET,
	"status":      CHECKSET,
}

func AlicloudHbrHanaBackupClientBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_hbr_vaults" "default" {
  		name_regex = "tf-test-hbr-hana-client"
	}
`, name)
}
