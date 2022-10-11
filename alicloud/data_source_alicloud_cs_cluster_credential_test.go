package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCSClusterCredentialDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_cluster_credential.default"

	testAccConfigInternet := dataSourceTestAccConfigFunc(
		resourceId,
		fmt.Sprintf("tf-testaccinternetk8s-%d", rand),
		dataSourceCSClusterCredentialConfigDependence_Internet,
	)

	idConfig := dataSourceTestAccConfig{
		existConfig: testAccConfigInternet(map[string]interface{}{
			"cluster_id":                 "${alicloud_cs_managed_kubernetes.default.id}",
			"temporary_duration_minutes": "60",
		}),
		fakeConfig: testAccConfigInternet(map[string]interface{}{
			"cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
		}),
	}

	var existCSClusterCredentialMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"cluster_id":                         CHECKSET,
			"cluster_name":                       REGEXMATCH + fmt.Sprintf("tf-testaccinternetk8s-%d", rand),
			"kube_config":                        CHECKSET,
			"certificate_authority.cluster_cert": CHECKSET,
			"certificate_authority.client_cert":  CHECKSET,
			"certificate_authority.client_key":   CHECKSET,
			"expiration":                         CHECKSET,
		}
	}

	var fakeCSClusterCredentialMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"cluster_id":                         CHECKSET,
			"cluster_name":                       REGEXMATCH + fmt.Sprintf("tf-testaccinternetk8s-%d", rand),
			"kube_config":                        CHECKSET,
			"certificate_authority.cluster_cert": CHECKSET,
			"certificate_authority.client_cert":  CHECKSET,
			"certificate_authority.client_key":   CHECKSET,
			"expiration":                         CHECKSET,
		}
	}

	var csClusterAuthCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCSClusterCredentialMapFunc,
		fakeMapFunc:  fakeCSClusterCredentialMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
	}
	csClusterAuthCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idConfig)
}

func dataSourceCSClusterCredentialConfigDependence_Internet(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
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

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix                 = "${var.name}"
  cluster_spec                = "ack.pro.small"
  worker_vswitch_ids          = [local.vswitch_id]
  new_nat_gateway             = true
  node_port_range             = "30000-32767"
  password                    = "Hello1234"
  pod_cidr                    = cidrsubnet("10.0.0.0/8", 8, 35)
  service_cidr                = cidrsubnet("172.16.0.0/16", 4, 6)
  install_cloud_monitor       = true
  slb_internet_enabled        = true
  worker_disk_category        = "cloud_efficiency"
  worker_data_disk_category   = "cloud_ssd"
  worker_data_disk_size       = 200
  worker_disk_size            = 40
  worker_instance_charge_type = "PostPaid"
}
`, name)
}
