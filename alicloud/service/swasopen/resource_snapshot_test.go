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

func TestAccAliCloudSimpleApplicationServerSnapshot_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_simple_application_server_snapshot.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSimpleApplicationServerSnapshotMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return swasopen.NewSwasOpenService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSimpleApplicationServerSnapshot")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc-swassnapshot%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSimpleApplicationServerSnapshotBasicDependence0)
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
					"disk_id":       "${data.alicloud_simple_application_server_disks.default.ids.0}",
					"snapshot_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_id":       acctest.CHECKSET,
						"snapshot_name": name,
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

var AlicloudSimpleApplicationServerSnapshotMap0 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudSimpleApplicationServerSnapshotBasicDependence0(name string) string {
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
  instance_id = alicloud_simple_application_server_instance.default.id
}
`, name)
}
