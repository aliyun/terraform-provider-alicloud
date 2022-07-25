package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	cs "github.com/alibabacloud-go/cs-20151215/v3/client"
)

func TestAccAlicloudCSKubernetesPermissions_basic(t *testing.T) {
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
					"permissions": []map[string]string{
						{
							"cluster":     "${alicloud_cs_managed_kubernetes.default.0.id}",
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
						"uid": CHECKSET,
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

data "alicloud_instance_types" "default" {
	availability_zone          = data.alicloud_zones.default.zones.0.id
	cpu_core_count             = 2
	memory_size                = 4
	kubernetes_node_role       = "Worker"
}

resource "alicloud_vpc" "vpc" {
	vpc_name   = "tf_unittest_cs"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vswitch" {
	count             = 1
	vpc_id            = alicloud_vpc.vpc.id
	cidr_block        = cidrsubnet(alicloud_vpc.vpc.cidr_block, 8, 8)
	zone_id           = data.alicloud_zones.default.zones.0.id
	vswitch_name      = var.name
}

# Create a management cluster
resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  worker_number                = 2
  deletion_protection          = false
  password                     = "Hello1234"
  pod_cidr                     = "10.11.0.0/16"
  service_cidr                 = "192.168.0.0/16"
  worker_vswitch_ids           = [alicloud_vswitch.vswitch.id]
  worker_instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
  depends_on                   = ["alicloud_ram_user_policy_attachment.attach"]
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
