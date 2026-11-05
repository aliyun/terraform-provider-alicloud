package swasopen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	swasopen "github.com/aliyun/terraform-provider-alicloud/alicloud/service/swasopen"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSimpleApplicationServerCustomImage_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_simple_application_server_custom_image.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSimpleApplicationServerCustomImageMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return swasopen.NewSwasOpenService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSimpleApplicationServerCustomImage")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-swascustomimage%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSimpleApplicationServerCustomImageBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
			acctest.PreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_image_name":  "${var.name}",
					"system_snapshot_id": "${alicloud_simple_application_server_snapshot.default.id}",
					"instance_id":        "${data.alicloud_simple_application_server_disks.default.disks.0.instance_id}",
					"description":        "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_image_name": name,
						"description":       name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Share",
				}),
				Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{})),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "UnShare",
				}),
				Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{})),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "system_snapshot_id", "status"},
			},
		},
	})
}

var AlicloudSimpleApplicationServerCustomImageMap0 = map[string]string{}

func AlicloudSimpleApplicationServerCustomImageBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_simple_application_server_images" "default" {
	platform = "Linux"
}
data "alicloud_simple_application_server_plans" "default" {
	platform = "Linux"
}

resource "alicloud_simple_application_server_instance" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

data "alicloud_simple_application_server_disks" "default" {
  disk_type = "System"
  instance_id = alicloud_simple_application_server_instance.default.id
}

resource "alicloud_simple_application_server_snapshot" "default" {
  disk_id       = data.alicloud_simple_application_server_disks.default.ids.0
  snapshot_name = var.name
}
`, name)
}
