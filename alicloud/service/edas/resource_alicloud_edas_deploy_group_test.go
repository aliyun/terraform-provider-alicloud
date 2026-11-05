package edas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	edassvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/edas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudEdasDeployGroup_basic(t *testing.T) {
	var v *edas.DeployGroup
	resourceId := "alicloud_edas_deploy_group.default"

	ra := acctest.ResourceAttrInit(resourceId, edasDeployGroupBasicMap)
	serviceFunc := func() interface{} {
		return edassvc.NewEdasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}

	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(1000, 9999)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasdeploygroupbasic%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEdasDeployGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEdasDeployGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":     "${alicloud_edas_application.default.id}",
					"group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name,
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
					"group_name": fmt.Sprintf("tf-testacc-edasdeploygroupchange%v", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": fmt.Sprintf("tf-testacc-edasdeploygroupchange%v", rand)}),
				),
			},
		},
	})
}

func TestAccAliCloudEdasDeployGroup_multi(t *testing.T) {
	var v *edas.DeployGroup
	resourceId := "alicloud_edas_deploy_group.default.1"

	ra := acctest.ResourceAttrInit(resourceId, edasDeployGroupBasicMap)
	serviceFunc := func() interface{} {
		return edassvc.NewEdasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}

	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(1000, 9999)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasdeploygroupmulti%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEdasDeployGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEdasDeployGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":      "2",
					"app_id":     "${alicloud_edas_application.default.id}",
					"group_name": "${var.name}-${count.index}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckEdasDeployGroupDestroy(s *terraform.State) error {
	return nil
}

var edasDeployGroupBasicMap = map[string]string{
	"app_id":     acctest.CHECKSET,
	"group_name": acctest.CHECKSET,
	"group_type": acctest.CHECKSET,
}

func resourceEdasDeployGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		data "alicloud_vpcs" "default" {
			name_regex = "^default-NODELETING$"
		}

		resource "alicloud_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = data.alicloud_vpcs.default.ids.0
		}

		resource "alicloud_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = "${alicloud_edas_cluster.default.id}"
		  package_type = "JAR"
		}
		`, name)
}
