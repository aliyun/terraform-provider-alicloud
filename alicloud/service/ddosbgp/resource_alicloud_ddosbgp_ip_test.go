package ddosbgp_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	ddosbgp "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ddosbgp"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test DdosBgp Ip. >>> Resource test cases, automatically generated.
// Case IP用例 11365
func TestAccAliCloudDdosBgpIp_basic11365(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddosbgp_ip.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDdosBgpIpMap11365)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddosbgp.NewDdosBgpServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosBgpIp")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccddosbgp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDdosBgpIpBasicDependence11365)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_ddosbgp_instance.default.id}",
					"ip":          "${alicloud_eip_address.default.ip_address}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
						"ip":          acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudDdosBgpIp_basic11365_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddosbgp_ip.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDdosBgpIpMap11365)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddosbgp.NewDdosBgpServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosBgpIp")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccddosbgp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDdosBgpIpBasicDependence11365)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_ddosbgp_instance.default.id}",
					"ip":          "${alicloud_eip_address.default.ip_address}",
					"member_uid":  "${data.alicloud_account.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
						"ip":          acctest.CHECKSET,
						"member_uid":  acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudDdosBgpIpMap11365 = map[string]string{
	"status":     acctest.CHECKSET,
	"member_uid": acctest.CHECKSET,
}

func AliCloudDdosBgpIpBasicDependence11365(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {
}

resource "alicloud_ddosbgp_instance" "default" {
  normal_bandwidth = "300"
  ip_count         = "100"
  type             = "Enterprise"
  base_bandwidth   = "20"
  ip_type          = "IPv4"
  bandwidth        = "-1"
  period           = "12"
}

resource "alicloud_eip_address" "default" {
	address_name = var.name
}


`, name)
}

// Test DdosBgp Ip. <<< Resource test cases, automatically generated.
