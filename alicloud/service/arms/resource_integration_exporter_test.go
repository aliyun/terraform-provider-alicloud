package arms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudArmsIntegrationExporter_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_integration_exporter.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsIntegrationExporterMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsIntegrationExporter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10, 99)
	name := fmt.Sprintf("tf-testacc-ArmsIE%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsIntegrationExporterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":       "${alicloud_arms_prometheus.default.id}",
					"integration_type": "kafka",
					"param":            `{\"tls_insecure-skip-tls-verify\":\"none=tls.insecure-skip-tls-verify\",\"tls_enabled\":\"none=tls.enabled\",\"sasl_mechanism\":\"\",\"name\":\"kafka1\",\"sasl_enabled\":\"none=sasl.enabled\",\"ip_ports\":\"abc:666\",\"scrape_interval\":30,\"version\":\"0.10.1.0\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":       acctest.CHECKSET,
						"integration_type": "kafka",
						"param":            "{\"tls_insecure-skip-tls-verify\":\"none=tls.insecure-skip-tls-verify\",\"tls_enabled\":\"none=tls.enabled\",\"sasl_mechanism\":\"\",\"name\":\"kafka1\",\"sasl_enabled\":\"none=sasl.enabled\",\"ip_ports\":\"abc:666\",\"scrape_interval\":30,\"version\":\"0.10.1.0\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param": `{\"tls_insecure-skip-tls-verify\":\"none=tls.insecure-skip-tls-verify\",\"tls_enabled\":\"none=tls.enabled\",\"sasl_mechanism\":\"\",\"name\":\"kafka1\",\"sasl_enabled\":\"none=sasl.enabled\",\"ip_ports\":\"abc:888\",\"scrape_interval\":30,\"version\":\"0.10.1.0\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param": "{\"tls_insecure-skip-tls-verify\":\"none=tls.insecure-skip-tls-verify\",\"tls_enabled\":\"none=tls.enabled\",\"sasl_mechanism\":\"\",\"name\":\"kafka1\",\"sasl_enabled\":\"none=sasl.enabled\",\"ip_ports\":\"abc:888\",\"scrape_interval\":30,\"version\":\"0.10.1.0\"}",
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

var AlicloudArmsIntegrationExporterMap = map[string]string{
	"instance_id": acctest.CHECKSET,
}

func AlicloudArmsIntegrationExporterBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_arms_prometheus" "default" {
  		cluster_type        = "ecs"
  		grafana_instance_id = "free"
  		vpc_id              = data.alicloud_vpcs.default.ids.0
  		vswitch_id          = data.alicloud_vswitches.default.ids.0
  		security_group_id   = alicloud_security_group.default.id
  		cluster_name        = "${var.name}-${data.alicloud_vpcs.default.ids.0}"
  		resource_group_id   = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	}
`, name)
}
