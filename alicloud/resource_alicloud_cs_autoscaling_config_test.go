package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	cs "github.com/alibabacloud-go/cs-20151215/v3/client"
)

func TestAccAlicloudCSAutoscalingConfig_basic(t *testing.T) {
	var v *cs.CreateAutoscalingConfigRequest
	resourceId := "alicloud_cs_autoscaling_config.default"
	serviceFunc := func() interface{} {
		client, _ := testAccProvider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return &CsClient{client}
	}

	ra := resourceAttrInit(resourceId, baseMap)
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rc.describeMethod = "DescribeCsAutoscalingConfig"
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSAutoscalingConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSAuoscalingConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":                "${alicloud_cs_managed_kubernetes.default.0.id}",
					"cool_down_duration":        "10m",
					"unneeded_duration":         "10m",
					"utilization_threshold":     "0.5",
					"gpu_utilization_threshold": "0.5",
					"scan_interval":             "30s",
					"scale_down_enabled":        "true",
					"expander":                  "least-waste",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                CHECKSET,
						"cool_down_duration":        CHECKSET,
						"unneeded_duration":         CHECKSET,
						"utilization_threshold":     CHECKSET,
						"gpu_utilization_threshold": CHECKSET,
						"scan_interval":             CHECKSET,
						"scale_down_enabled":        CHECKSET,
						"expander":                  CHECKSET,
					}),
				),
			},
		},
	})
}

var baseMap = map[string]string{
	"cluster_id": CHECKSET,
}

func resourceCSAuoscalingConfigDependence(name string) string {
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
	vpc_name   = var.name
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
  worker_number                = 0
  worker_disk_size             = 40
  deletion_protection          = false
  node_port_range              = "30000-32767"
  password                     = "Hello1234"
  pod_cidr                     = "10.11.0.0/16"
  service_cidr                 = "192.168.0.0/16"
  worker_vswitch_ids           = [alicloud_vswitch.vswitch.id]
  worker_instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
}
`, name)
}
