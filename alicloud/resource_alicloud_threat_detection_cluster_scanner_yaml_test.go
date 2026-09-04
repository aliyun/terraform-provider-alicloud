package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudThreatDetectionClusterScannerYaml_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_cluster_scanner_yaml.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClusterScannerYamlMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClusterScannerYaml")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionClusterScannerYaml%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClusterScannerYamlBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":   "${alicloud_cs_serverless_kubernetes.default.id}",
					"webhook_open": 1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":       CHECKSET,
						"webhook_open":     "1",
						"region_id":        CHECKSET,
						"ca_cert_base64":   CHECKSET,
						"tls_key_base64":   CHECKSET,
						"tls_cert_base64":  CHECKSET,
						"cluster_env_info": CHECKSET,
						"image":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":   "${alicloud_cs_serverless_kubernetes.default.id}",
					"webhook_open": 0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   CHECKSET,
						"webhook_open": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudThreatDetectionClusterScannerYamlMap = map[string]string{
	"cluster_id":       CHECKSET,
	"webhook_open":     CHECKSET,
	"region_id":        CHECKSET,
	"ca_cert_base64":   CHECKSET,
	"tls_key_base64":   CHECKSET,
	"tls_cert_base64":  CHECKSET,
	"cluster_env_info": CHECKSET,
	"image":            CHECKSET,
}

func AlicloudThreatDetectionClusterScannerYamlBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

data "alicloud_cs_kubernetes_version" "version" {
  cluster_type       = "Kubernetes"
  kubernetes_version = "1.28"
  profile            = "Serverless"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_cs_serverless_kubernetes" "default" {
  name                = var.name
  vpc_id              = alicloud_vpc.default.id
  vswitch_ids         = [alicloud_vswitch.default.id]
  security_group_id   = alicloud_security_group.default.id
  new_nat_gateway     = true
  deletion_protection = false
  version             = data.alicloud_cs_kubernetes_version.version.metadata.0.version
}
`, name)
}
