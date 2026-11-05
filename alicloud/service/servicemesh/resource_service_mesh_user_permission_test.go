package servicemesh_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/servicemesh"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudServiceMeshUserPermission_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_mesh_user_permission.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudServiceMeshUserPermissionMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return servicemesh.NewServicemeshService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeUserPermissions")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sservicemesh%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudServiceMeshUserPermissionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.ServiceMeshStandardUnsupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sub_account_user_id": "${alicloud_ram_user.default.id}",
					"permissions": []map[string]interface{}{
						{
							"role_name":       "istio-ops",
							"service_mesh_id": "${alicloud_service_mesh_service_mesh.default1.id}",
							"role_type":       "custom",
							"is_custom":       "true",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permissions": []map[string]interface{}{
						{
							"role_name":       "istio-admin",
							"service_mesh_id": "${alicloud_service_mesh_service_mesh.default1.id}",
							"role_type":       "custom",
							"is_custom":       "true",
							"is_ram_role":     "false",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"permissions": []map[string]interface{}{
						{
							"role_name":       "istio-readonly",
							"service_mesh_id": "${alicloud_service_mesh_service_mesh.default2.id}",
							"role_type":       "custom",
							"is_custom":       "true",
							"is_ram_role":     "false",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"permissions.role_type", "permissions.is_custom"},
			},
		},
	})
}

var AlicloudServiceMeshUserPermissionMap0 = map[string]string{}

func AlicloudServiceMeshUserPermissionBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
variable "service_mesh_name"{
  default = "%s"
}
data "alicloud_service_mesh_versions" "default" {
	edition = "Default"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

resource "alicloud_vpc" "default" {
	count = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
}

data "alicloud_vswitches" "default" {
	vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
}

resource "alicloud_vswitch" "default" {
	count      = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	vpc_id     = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
	cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
	zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

resource "alicloud_service_mesh_service_mesh" "default1" {
	service_mesh_name = var.name
    edition = "Default"
    cluster_spec = "standard"
    version = data.alicloud_service_mesh_versions.default.versions.0.version
	network {
		vpc_id        = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
		vswitche_list = [length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id]
	}
	load_balancer {
		pilot_public_eip = false
		api_server_public_eip = false
	}
}

resource "alicloud_service_mesh_service_mesh" "default2" {
	service_mesh_name = "${var.name}-2"
    edition = "Default"
    cluster_spec = "standard"
    version = data.alicloud_service_mesh_versions.default.versions.0.version
	network {
		vpc_id        = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
		vswitche_list = [length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id]
	}
	load_balancer {
		pilot_public_eip = false
		api_server_public_eip = false
	}
}



`, name, fmt.Sprintf("%s%d", name, tfacctest.RandIntRange(1, 999)))
}
