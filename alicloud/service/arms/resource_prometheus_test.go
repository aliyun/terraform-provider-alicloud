package arms_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudArmsPrometheus_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsPrometheusMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsPrometheus")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%sArmsPrometheus%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsPrometheusBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sub_clusters_json":   "60",
					"cluster_type":        "remote-write",
					"grafana_instance_id": "free",
					"cluster_name":        name,
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"payment_type":        "POSTPAY",
					"duration":            "90",
					"archive_duration":    "60",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Prometheus",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_type":        "remote-write",
						"grafana_instance_id": "free",
						"cluster_name":        name,
						"resource_group_id":   acctest.CHECKSET,
						"payment_type":        "POSTPAY",
						"duration":            "90",
						"archive_duration":    "60",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "Prometheus",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sub_clusters_json": "61",
					"duration":          "180",
					"archive_duration":  "90",
					"payment_type":      "POSTPAY_GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"duration":         "180",
						"archive_duration": "90",
						"payment_type":     "POSTPAY_GB",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "Prometheus_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "Prometheus_Update",
					}),
				),
			},
			{
				// cluster_id only applies to cluster_type "aliyun-cs" (an existing ACK
				// cluster) and rebinding grafana_instance_id needs a dedicated Grafana
				// workspace; neither is available to this test. Removing the required
				// grafana_instance_id stops the plan before the sentinel reaches the API.
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":          "",
					"grafana_instance_id": acctest.REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudArmsPrometheus_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsPrometheusMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsPrometheus")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc%s", acctest.DefaultRegionToTest)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsPrometheusBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_type":        "ecs",
					"grafana_instance_id": "free",
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":          "${data.alicloud_vswitches.default.ids.0}",
					"security_group_id":   "${alicloud_security_group.default.id}",
					"cluster_name":        name + "-" + "${data.alicloud_vpcs.default.ids.0}",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Prometheus",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_type":        "ecs",
						"grafana_instance_id": "free",
						"vpc_id":              acctest.CHECKSET,
						"vswitch_id":          acctest.CHECKSET,
						"security_group_id":   acctest.CHECKSET,
						"cluster_name":        acctest.CHECKSET,
						"resource_group_id":   acctest.CHECKSET,
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "Prometheus",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "Prometheus_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "Prometheus_Update",
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

var AlicloudArmsPrometheusMap = map[string]string{}

func AlicloudArmsPrometheusBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}
	
	data "alicloud_account" "default" {
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name)
}

func AlicloudArmsPrometheusBasicDependence1(name string) string {
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
  		vpc_id      = data.alicloud_vpcs.default.ids.0
	}
`, name)
}
