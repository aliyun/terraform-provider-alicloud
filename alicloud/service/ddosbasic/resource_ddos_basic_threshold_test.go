package ddosbasic_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ddosbasic"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudDdosBasicThreshold_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddos_basic_threshold.default"
	// ap-southeast-1 (CheckoutSupportedRegions' default backup region when
	// ALICLOUD_REGION is unset) lacks sellable ECS for the domestic test
	// account, so the alicloud_instance dependency returns an empty
	// instance_types list there (the EIP-based defense threshold ACC, with no
	// ECS dependency, passes in the same region/account). Pin this test to
	// cn-hangzhou, which supports both DdosBasic and ECS for the account.
	// Region adjustment per migration-standard §5.
	os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
	acctest.DefaultRegionToTest = string(connectivity.Hangzhou)
	acctest.CheckoutSupportedRegions(t, true, connectivity.DdosBasicSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDdosBasicThresholdMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddosbasic.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosBasicThreshold")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sddosbasicthreshold%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDdosBasicThresholdBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "ecs",
					"instance_id":   "${alicloud_instance.default.id}",
					"internet_ip":   "${alicloud_instance.default.public_ip}",
					"bps":           "100",
					"pps":           "60000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": "ecs",
						"instance_id":   acctest.CHECKSET,
						"internet_ip":   acctest.CHECKSET,
						"bps":           "100",
						"pps":           "60000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bps": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bps": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pps": "70000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pps": "70000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bps": "100",
					"pps": "60000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bps": "100",
						"pps": "60000",
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

var AliCloudDdosBasicThresholdMap0 = map[string]string{
	"max_bps": acctest.CHECKSET,
	"max_pps": acctest.CHECKSET,
}

func AliCloudDdosBasicThresholdBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
}
`, name)
}
