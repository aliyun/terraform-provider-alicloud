package imm_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/imm"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudIMMProject_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_imm_project.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudIMMProjectMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return imm.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeImmProject")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%simmproject%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudIMMProjectBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.IMMSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_role": "${alicloud_ram_role.role.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_role": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_role": "AliyunIMMDefaultRole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_role": "AliyunIMMDefaultRole",
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

var AlicloudIMMProjectMap0 = map[string]string{}

func AlicloudIMMProjectBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_ram_role" "role" {
  name     = var.name
  document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "imm.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
  description = "this is a role test."
  force       = true
}
`, name)
}
