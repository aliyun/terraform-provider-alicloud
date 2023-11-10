package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	cs "github.com/alibabacloud-go/cs-20151215/v3/client"
)

func TestAccAliCloudCSKubernetesPermissions_basic(t *testing.T) {
	var v []*cs.GrantPermissionsRequestBody
	resourceId := "alicloud_cs_kubernetes_permissions.default"
	serviceFunc := func() interface{} {
		client, _ := testAccProvider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return &CsClient{client}
	}

	ra := resourceAttrInit(resourceId, permissionMap)
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rc.describeMethod = "DescribeUserPermission"
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSKubernetesPermission-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSPermissionsConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"uid": "${alicloud_ram_user.user.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"uid":           CHECKSET,
						"permissions.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"uid": "${alicloud_ram_user.user.id}",
					"permissions": []map[string]string{
						{
							"cluster":     "${alicloud_cs_managed_kubernetes.default.id}",
							"role_type":   "cluster",
							"role_name":   "dev",
							"namespace":   "",
							"is_custom":   "false",
							"is_ram_role": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"uid":           CHECKSET,
						"permissions.#": "1",
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
					"permissions": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"uid":           CHECKSET,
						"permissions.#": "0",
					}),
				),
			},
		},
	})
}

var permissionMap = map[string]string{
	"uid": CHECKSET,
}

func resourceCSPermissionsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation  = "VSwitch"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

# Create a management cluster
resource "alicloud_cs_managed_kubernetes" "default" {
  name                        = var.name
  cluster_spec                = "ack.pro.small"
  worker_vswitch_ids          = [local.vswitch_id]
  new_nat_gateway             = false
  pod_cidr                    = cidrsubnet("10.0.0.0/8", 8, 41)
  service_cidr                = cidrsubnet("172.17.0.0/16", 4, 3)
  slb_internet_enabled        = false
  depends_on                  = ["alicloud_ram_user_policy_attachment.attach"]
}

# Create a new RAM user.
resource "alicloud_ram_user" "user" {
  name         = var.name
  display_name = var.name
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
}

# Create a new RAM Policy, .
resource "alicloud_ram_policy" "policy" {
  policy_name     = var.name
  policy_document = <<EOF
  {
    "Statement":[
      {
        "Action":"cs:Get*",
        "Effect":"Allow",
        "Resource":[
            "*"
        ]
      }
    ],
    "Version":"1"
  }
  EOF
  description = "this is a policy test by tf"
}

# Authorize the RAM user
resource "alicloud_ram_user_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.name
  policy_type = alicloud_ram_policy.policy.type
  user_name   = alicloud_ram_user.user.name
}
`, name)
}
