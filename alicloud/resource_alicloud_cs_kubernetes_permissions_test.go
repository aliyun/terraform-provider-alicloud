package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	cs "github.com/alibabacloud-go/cs-20151215/v4/client"
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
			{
				Config: testAccConfig(map[string]interface{}{
					"uid": "${alicloud_ram_user.user.id}",
					"permissions": []map[string]string{
						{
							"cluster":     "${alicloud_cs_managed_kubernetes.default.0.id}",
							"role_type":   "cluster",
							"role_name":   "admin",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"uid": "${alicloud_ram_user.user.id}",
					"permissions": []map[string]string{
						{
							"cluster":     "${alicloud_cs_managed_kubernetes.default.0.id}",
							"role_type":   "cluster",
							"role_name":   "admin",
							"namespace":   "",
							"is_custom":   "false",
							"is_ram_role": "false",
						},
						{
							"cluster":     "${alicloud_cs_managed_kubernetes.default.0.id}",
							"role_type":   "cluster",
							"role_name":   "admin",
							"namespace":   "",
							"is_custom":   "true",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"uid": "${alicloud_ram_user.user.id}",
					"permissions": []map[string]string{
						{
							"cluster":     "${alicloud_cs_managed_kubernetes.default.1.id}",
							"role_type":   "namespace",
							"role_name":   "dev",
							"namespace":   "default",
							"is_custom":   "false",
							"is_ram_role": "true",
						},
						{
							"cluster":     "${alicloud_cs_managed_kubernetes.default.1.id}",
							"role_type":   "namespace",
							"role_name":   "dev",
							"namespace":   "kube-system",
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

variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "10.0.0.0/8"
}

variable "vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'vswitch_ids' is not specified."
  type        = list(string)
  default     = ["10.1.0.0/16", "10.2.0.0/16"]
}

variable "pod_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or service's and cannot be in them."
  type        = list(string)
  default     = ["172.16.0.0/16", "172.17.0.0/16"]
}

variable "service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  type        = list(string)
  default     = ["172.18.0.0/16", "172.19.0.0/16"]
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

data "alicloud_cs_kubernetes_version" "default" {
  cluster_type = "ManagedKubernetes"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = var.vpc_cidr
}

resource "alicloud_vswitch" "default" {
  count      = length(var.vswitch_cidrs)
  vpc_id     = alicloud_vpc.vpc.id
  cidr_block = element(var.vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

# Create a new RAM cluster.
resource "alicloud_cs_managed_kubernetes" "default" {
  count                = 2
  name                 = join("-", [var.name, count.index])
  cluster_spec         = "ack.pro.small"
  version              = data.alicloud_cs_kubernetes_version.default.metadata.0.version
  worker_vswitch_ids   = split(",", join(",", alicloud_vswitch.default.*.id))
  new_nat_gateway      = false
  pod_cidr             = element(var.pod_cidr, count.index)
  service_cidr         = element(var.service_cidr, count.index)
  slb_internet_enabled = false
}

# Create a new RAM user.
resource "alicloud_ram_user" "user" {
  name = var.name
}
`, name)
}
