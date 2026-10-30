package hbr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/hbr"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func SkipTestAccAlicloudHbrHanaBackupClient_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_hbr_hana_backup_client.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHbrHanaBackupClientMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrHanaBackupClient")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sHbrHanaBackupClient%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHbrHanaBackupClientBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"vault_id":      acctest.CHECKSET,
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
	"client_id":   acctest.CHECKSET,
	"instance_id": acctest.CHECKSET,
	"cluster_id":  acctest.CHECKSET,
	"status":      acctest.CHECKSET,
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
