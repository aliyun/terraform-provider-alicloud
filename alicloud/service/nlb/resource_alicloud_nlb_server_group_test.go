package nlb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nlb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudNlbServerGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNLBServerGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNLBServerGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name,
						"vpc_id":            acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": []map[string]interface{}{
						{
							"health_check_interval":        "2",
							"health_check_url":             "/tf-testAcc",
							"unhealthy_threshold":          "3",
							"health_check_connect_timeout": "1",
							"health_check_enabled":         "true",
							"health_check_connect_port":    "46325",
							"healthy_threshold":            "3",
							"http_check_method":            "HEAD",
							"health_check_domain":          "tf-testAcc.com",
							"health_check_http_code":       []string{"http_2xx", "http_3xx", "http_4xx"},
							"health_check_type":            "HTTP",
							"health_check_http_version":    "HTTP1.0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": []map[string]interface{}{
						{
							"health_check_interval":        "2",
							"health_check_url":             "/tf-testAcc",
							"unhealthy_threshold":          "3",
							"health_check_connect_timeout": "1",
							"health_check_enabled":         "true",
							"health_check_connect_port":    "46325",
							"healthy_threshold":            "3",
							"http_check_method":            "HEAD",
							"health_check_domain":          "tf-testAcc.com",
							"health_check_http_code":       []string{"http_2xx", "http_3xx", "http_4xx"},
							"health_check_type":            "HTTP",
							"health_check_http_version":    "HTTP1.0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": []map[string]interface{}{
						{
							"health_check_interval":        "2",
							"health_check_url":             "/tf-testAcc",
							"unhealthy_threshold":          "3",
							"health_check_connect_timeout": "1",
							"health_check_enabled":         "true",
							"health_check_connect_port":    "46325",
							"healthy_threshold":            "3",
							"http_check_method":            "HEAD",
							"health_check_domain":          "tf-testAcc.com",
							"health_check_http_code":       []string{"http_2xx", "http_3xx", "http_4xx"},
							"health_check_type":            "HTTP",
							"health_check_http_version":    "HTTP1.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
						"health_check.0.health_check_http_version": "HTTP1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preserve_client_ip_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preserve_client_ip_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Tch",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Tch",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

func TestAccAliCloudNlbServerGroup_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNLBServerGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNLBServerGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"address_ip_version":       "Ipv4",
					"any_port_enabled":         "true",
					"connection_drain_enabled": "true",
					"connection_drain_timeout": "60",
					"health_check": []map[string]interface{}{
						{
							"health_check_interval":        "2",
							"health_check_url":             "/tf-testAcc",
							"unhealthy_threshold":          "3",
							"health_check_connect_timeout": "1",
							"health_check_enabled":         "true",
							"health_check_connect_port":    "46325",
							"healthy_threshold":            "3",
							"http_check_method":            "HEAD",
							"health_check_domain":          "tf-testAcc.com",
							"health_check_http_code":       []string{"http_2xx", "http_3xx", "http_4xx"},
							"health_check_type":            "HTTP",
							"health_check_http_version":    "HTTP1.0",
						},
					},
					"preserve_client_ip_enabled": "false",
					"protocol":                   "UDP",
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"scheduler":                  "Tch",
					"server_group_name":          name,
					"server_group_type":          "Ip",
					"vpc_id":                     "${data.alicloud_vpcs.default.vpcs.0.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_ip_version":         "Ipv4",
						"any_port_enabled":           "true",
						"connection_drain_enabled":   "true",
						"connection_drain_timeout":   "60",
						"health_check.#":             "1",
						"preserve_client_ip_enabled": "false",
						"protocol":                   "UDP",
						"resource_group_id":          acctest.CHECKSET,
						"scheduler":                  "Tch",
						"server_group_name":          name,
						"server_group_type":          "Ip",
						"vpc_id":                     acctest.CHECKSET,
						"tags.%":                     "2",
						"tags.Created":               "TF",
						"tags.For":                   "Test",
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

var AliCloudNLBServerGroupMap0 = map[string]string{
	"address_ip_version":         acctest.CHECKSET,
	"health_check.#":             acctest.CHECKSET,
	"preserve_client_ip_enabled": acctest.CHECKSET,
	"protocol":                   acctest.CHECKSET,
	"region_id":                  acctest.CHECKSET,
	"resource_group_id":          acctest.CHECKSET,
	"scheduler":                  acctest.CHECKSET,
	"server_group_type":          acctest.CHECKSET,
	"status":                     acctest.CHECKSET,
}

func AliCloudNLBServerGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}

func TestAccAliCloudNlbServerGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNLBServerGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNLBServerGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"server_group_name": "${var.name}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":            acctest.CHECKSET,
						"server_group_name": name,
						"health_check.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Wrr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Wrr",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "TCP",
							"health_check_connect_port":    "0",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "5",
							"health_check_interval":        "10",
							"http_check_method":            "GET",
							"health_check_url":             "/test/index.html",
							"health_check_domain":          "tf-testAcc.com",
							"health_check_http_code":       []string{"http_2xx", "http_3xx", "http_4xx"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain":         "true",
					"connection_drain_timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain":         "true",
						"connection_drain_timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc5",
						"For":     "Tftestacc 5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc5",
						"tags.For":     "Tftestacc 5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name":        "${var.name}",
					"scheduler":                "Rr",
					"connection_drain_timeout": "100",
					"tags": map[string]string{
						"Created": "tfTestAcc6",
						"For":     "Tftestacc 6",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name":        name,
						"scheduler":                "Rr",
						"connection_drain_timeout": "100",
						"tags.%":                   "2",
						"tags.Created":             "tfTestAcc6",
						"tags.For":                 "Tftestacc 6",
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

// lintignore: R001

func TestAccAliCloudNlbServerGroup_basic3734(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMap3734)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependence3734)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${alicloud_vpc.defaultxDhVj6.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name,
						"vpc_id":            acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preserve_client_ip_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preserve_client_ip_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":          "TCP",
					"server_group_type": "Instance",
					"server_group_name": name + "_update",
					"vpc_id":            "${alicloud_vpc.defaultxDhVj6.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":          "TCP",
						"server_group_type": "Instance",
						"server_group_name": name + "_update",
						"vpc_id":            acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AliCloudNlbServerGroupMap3734 = map[string]string{
	"status":                     acctest.CHECKSET,
	"address_ip_version":         acctest.CHECKSET,
	"scheduler":                  acctest.CHECKSET,
	"health_check.#":             acctest.CHECKSET,
	"server_group_type":          acctest.CHECKSET,
	"preserve_client_ip_enabled": acctest.CHECKSET,
	"connection_drain_enabled":   acctest.CHECKSET,
	"connection_drain_timeout":   acctest.CHECKSET,
	"any_port_enabled":           acctest.CHECKSET,
	"protocol":                   acctest.CHECKSET,
}

func AliCloudNlbServerGroupBasicDependence3734(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultxDhVj6" {
  description = "ceshi"
  vpc_name    = var.name

  cidr_block = "10.0.0.0/8"
}


`, name)
}

// Case 4668

func TestAccAliCloudNlbServerGroup_basic4668(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMap4668)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependence4668)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"any_port_enabled":  "false",
					"protocol":          "UDP",
					"server_group_type": "Instance",
					"server_group_name": name,
					"vpc_id":            "${alicloud_vpc.defaultEyoN1M.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_type":            "HTTP",
							"health_check_connect_port":    "0",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_enabled":         "true",
							"health_check_connect_timeout": "20",
							"health_check_interval":        "10",
							"health_check_domain":          "example.com",
							"health_check_url":             "/rdk",
							"health_check_http_code": []string{
								"http_2xx", "http_3xx", "http_4xx"},
							"http_check_method":         "HEAD",
							"health_check_http_version": "HTTP1.0",
						},
					},
					"scheduler":                "Qch",
					"connection_drain_timeout": "20",
					"connection_drain_enabled": "true",
					"address_ip_version":       "Ipv4",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"any_port_enabled":         "false",
						"protocol":                 "UDP",
						"server_group_type":        "Instance",
						"server_group_name":        name,
						"vpc_id":                   acctest.CHECKSET,
						"scheduler":                "Qch",
						"connection_drain_timeout": "20",
						"connection_drain_enabled": "true",
						"address_ip_version":       "Ipv4",
						"tags.%":                   "2",
						"tags.Created":             "TF",
						"tags.For":                 "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preserve_client_ip_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preserve_client_ip_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Qch",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Qch",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_timeout": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_timeout": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_connect_port":    "0",
							"health_check_connect_timeout": "20",
							"health_check_interval":        "10",
							"health_check_domain":          "example.com",
							"health_check_http_code":       []string{"http_4xx"},
							"http_check_method":            "HEAD",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preserve_client_ip_enabled": "true",
					"connection_drain_enabled":   "true",
					"any_port_enabled":           "true",
					"protocol":                   "UDP",
					"server_group_type":          "Instance",
					"server_group_name":          name + "_update",
					"vpc_id":                     "${alicloud_vpc.defaultEyoN1M.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "TCP",
							"health_check_connect_port":    "1",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "1",
							"health_check_interval":        "5",
							"health_check_domain":          "$SERVER_IP",
							"health_check_url":             "/rdktest",
							"health_check_http_code": []string{
								"http_2xx"},
							"http_check_method": "HEAD",
						},
					},
					"scheduler":                "Qch",
					"connection_drain_timeout": "10",
					"address_ip_version":       "Ipv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preserve_client_ip_enabled": "true",
						"connection_drain_enabled":   "true",
						"any_port_enabled":           "true",
						"protocol":                   "UDP",
						"server_group_type":          "Instance",
						"server_group_name":          name + "_update",
						"vpc_id":                     acctest.CHECKSET,
						"scheduler":                  "Qch",
						"connection_drain_timeout":   "10",
						"address_ip_version":         "Ipv4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AliCloudNlbServerGroupMap4668 = map[string]string{
	"status":                     acctest.CHECKSET,
	"address_ip_version":         acctest.CHECKSET,
	"scheduler":                  acctest.CHECKSET,
	"health_check.#":             acctest.CHECKSET,
	"server_group_type":          acctest.CHECKSET,
	"preserve_client_ip_enabled": acctest.CHECKSET,
	"connection_drain_enabled":   acctest.CHECKSET,
	"connection_drain_timeout":   acctest.CHECKSET,
	"any_port_enabled":           acctest.CHECKSET,
	"protocol":                   acctest.CHECKSET,
}

func AliCloudNlbServerGroupBasicDependence4668(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultEyoN1M" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}


`, name)
}

// Case 4639

func TestAccAliCloudNlbServerGroup_basic4639(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMap4639)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependence4639)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${alicloud_vpc.testsr.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name,
						"vpc_id":            acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preserve_client_ip_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preserve_client_ip_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Wrr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Wrr",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Rr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Rr",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_timeout": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_timeout": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Wrr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Wrr",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Tch",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Tch",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preserve_client_ip_enabled": "true",
					"connection_drain_enabled":   "true",
					"any_port_enabled":           "true",
					"protocol":                   "TCP",
					"server_group_type":          "Instance",
					"server_group_name":          name + "_update",
					"vpc_id":                     "${alicloud_vpc.testsr.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "TCP",
							"health_check_connect_port":    "1",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "1",
							"health_check_interval":        "5",
							"health_check_domain":          "$SERVER_IP",
							"health_check_url":             "/rdktest",
							"health_check_http_code": []string{
								"http_2xx"},
							"http_check_method": "HEAD",
						},
					},
					"scheduler":                "Wrr",
					"connection_drain_timeout": "10",
					"address_ip_version":       "Ipv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preserve_client_ip_enabled": "true",
						"connection_drain_enabled":   "true",
						"any_port_enabled":           "true",
						"protocol":                   "TCP",
						"server_group_type":          "Instance",
						"server_group_name":          name + "_update",
						"vpc_id":                     acctest.CHECKSET,
						"scheduler":                  "Wrr",
						"connection_drain_timeout":   "10",
						"address_ip_version":         "Ipv4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AliCloudNlbServerGroupMap4639 = map[string]string{
	"status":                     acctest.CHECKSET,
	"address_ip_version":         acctest.CHECKSET,
	"scheduler":                  acctest.CHECKSET,
	"health_check.#":             acctest.CHECKSET,
	"server_group_type":          acctest.CHECKSET,
	"preserve_client_ip_enabled": acctest.CHECKSET,
	"connection_drain_enabled":   acctest.CHECKSET,
	"connection_drain_timeout":   acctest.CHECKSET,
	"any_port_enabled":           acctest.CHECKSET,
	"protocol":                   acctest.CHECKSET,
}

func AliCloudNlbServerGroupBasicDependence4639(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "testsr" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}


`, name)
}

// Case 4662

func TestAccAliCloudNlbServerGroup_basic4662(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMap4662)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependence4662)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${alicloud_vpc.defaultEyoN1M.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name,
						"vpc_id":            acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preserve_client_ip_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preserve_client_ip_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Wrr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Wrr",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_timeout": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_timeout": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preserve_client_ip_enabled": "true",
					"connection_drain_enabled":   "true",
					"any_port_enabled":           "true",
					"protocol":                   "TCP_UDP",
					"server_group_type":          "Instance",
					"server_group_name":          name + "_update",
					"vpc_id":                     "${alicloud_vpc.defaultEyoN1M.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "TCP",
							"health_check_connect_port":    "1",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "1",
							"health_check_interval":        "5",
							"health_check_domain":          "$SERVER_IP",
							"health_check_url":             "/rdktest",
							"health_check_http_code": []string{
								"http_2xx", "http_4xx"},
							"http_check_method": "HEAD",
						},
					},
					"scheduler":                "Wrr",
					"connection_drain_timeout": "10",
					"address_ip_version":       "DualStack",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preserve_client_ip_enabled": "true",
						"connection_drain_enabled":   "true",
						"any_port_enabled":           "true",
						"protocol":                   "TCP_UDP",
						"server_group_type":          "Instance",
						"server_group_name":          name + "_update",
						"vpc_id":                     acctest.CHECKSET,
						"scheduler":                  "Wrr",
						"connection_drain_timeout":   "10",
						"address_ip_version":         "DualStack",
						"resource_group_id":          acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AliCloudNlbServerGroupMap4662 = map[string]string{
	"status":                     acctest.CHECKSET,
	"address_ip_version":         acctest.CHECKSET,
	"scheduler":                  acctest.CHECKSET,
	"health_check.#":             acctest.CHECKSET,
	"server_group_type":          acctest.CHECKSET,
	"preserve_client_ip_enabled": acctest.CHECKSET,
	"connection_drain_enabled":   acctest.CHECKSET,
	"connection_drain_timeout":   acctest.CHECKSET,
	"any_port_enabled":           acctest.CHECKSET,
	"protocol":                   acctest.CHECKSET,
}

func AliCloudNlbServerGroupBasicDependence4662(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultEyoN1M" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

  enable_ipv6 = true
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}

`, name)
}

// Case 5353

func TestAccAliCloudNlbServerGroup_basic5353(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMap5353)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependence5353)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${alicloud_vpc.sgtestnlb.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name,
						"vpc_id":            acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_timeout": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_timeout": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_enabled": "true",
					"any_port_enabled":         "true",
					"protocol":                 "TCP",
					"server_group_type":        "Ip",
					"server_group_name":        name + "_update",
					"vpc_id":                   "${alicloud_vpc.sgtestnlb.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "HTTP",
							"health_check_connect_port":    "1",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "1",
							"health_check_interval":        "5",
							"health_check_domain":          "$SERVER_IP",
							"health_check_url":             "/rdktest",
							"health_check_http_code": []string{
								"http_2xx"},
							"http_check_method":         "HEAD",
							"health_check_http_version": "HTTP1.0",
						},
					},
					"connection_drain_timeout": "10",
					"address_ip_version":       "Ipv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_enabled": "true",
						"any_port_enabled":         "true",
						"protocol":                 "TCP",
						"server_group_type":        "Ip",
						"server_group_name":        name + "_update",
						"vpc_id":                   acctest.CHECKSET,
						"connection_drain_timeout": "10",
						"address_ip_version":       "Ipv4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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

var AliCloudNlbServerGroupMap5353 = map[string]string{
	"status":                     acctest.CHECKSET,
	"address_ip_version":         acctest.CHECKSET,
	"scheduler":                  acctest.CHECKSET,
	"health_check.#":             acctest.CHECKSET,
	"server_group_type":          acctest.CHECKSET,
	"preserve_client_ip_enabled": acctest.CHECKSET,
	"connection_drain_enabled":   acctest.CHECKSET,
	"connection_drain_timeout":   acctest.CHECKSET,
	"any_port_enabled":           acctest.CHECKSET,
	"protocol":                   acctest.CHECKSET,
}

func AliCloudNlbServerGroupBasicDependence5353(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "sgtestnlb" {
  description = "nlb tf用例测试"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}


`, name)
}

// Case 3734  twin

func TestAccAliCloudNlbServerGroup_basic3734_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMap3734)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependence3734)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"preserve_client_ip_enabled": "true",
					"protocol":                   "TCP",
					"server_group_type":          "Instance",
					"server_group_name":          name,
					"vpc_id":                     "${alicloud_vpc.defaultxDhVj6.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preserve_client_ip_enabled": "true",
						"protocol":                   "TCP",
						"server_group_type":          "Instance",
						"server_group_name":          name,
						"vpc_id":                     acctest.CHECKSET,
						"tags.%":                     "2",
						"tags.Created":               "TF",
						"tags.For":                   "Test",
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

// Case 4668  twin

func TestAccAliCloudNlbServerGroup_basic4668_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMap4668)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependence4668)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"any_port_enabled":  "false",
					"protocol":          "UDP",
					"server_group_type": "Instance",
					"server_group_name": name,
					"vpc_id":            "${alicloud_vpc.defaultEyoN1M.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_type":            "HTTP",
							"health_check_connect_port":    "0",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_enabled":         "true",
							"health_check_connect_timeout": "20",
							"health_check_interval":        "10",
							"health_check_domain":          "example.com",
							"health_check_url":             "/rdk",
							"health_check_http_code": []string{
								"http_2xx", "http_3xx", "http_4xx"},
							"http_check_method":         "HEAD",
							"health_check_http_version": "HTTP1.0",
						},
					},
					"scheduler":                "Qch",
					"connection_drain_timeout": "20",
					"connection_drain_enabled": "true",
					"address_ip_version":       "Ipv4",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"any_port_enabled":         "false",
						"protocol":                 "UDP",
						"server_group_type":        "Instance",
						"server_group_name":        name,
						"vpc_id":                   acctest.CHECKSET,
						"scheduler":                "Qch",
						"connection_drain_timeout": "20",
						"connection_drain_enabled": "true",
						"address_ip_version":       "Ipv4",
						"tags.%":                   "2",
						"tags.Created":             "TF",
						"tags.For":                 "Test",
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

// Case 4639  twin

func TestAccAliCloudNlbServerGroup_basic4639_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMap4639)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependence4639)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"any_port_enabled":  "false",
					"protocol":          "TCP",
					"server_group_type": "Instance",
					"server_group_name": name,
					"vpc_id":            "${alicloud_vpc.testsr.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "HTTP",
							"health_check_connect_port":    "0",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "20",
							"health_check_interval":        "10",
							"health_check_domain":          "example.com",
							"health_check_url":             "/rdk",
							"health_check_http_code": []string{
								"http_2xx", "http_3xx", "http_4xx"},
							"http_check_method":         "HEAD",
							"health_check_http_version": "HTTP1.0",
						},
					},
					"scheduler":                "Tch",
					"connection_drain_timeout": "20",
					"connection_drain_enabled": "true",
					"address_ip_version":       "Ipv4",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"any_port_enabled":         "false",
						"protocol":                 "TCP",
						"server_group_type":        "Instance",
						"server_group_name":        name,
						"vpc_id":                   acctest.CHECKSET,
						"scheduler":                "Tch",
						"connection_drain_timeout": "20",
						"address_ip_version":       "Ipv4",
						"tags.%":                   "2",
						"tags.Created":             "TF",
						"tags.For":                 "Test",
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

// Case 4662  twin

func TestAccAliCloudNlbServerGroup_basic4662_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMap4662)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependence4662)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"any_port_enabled":  "false",
					"protocol":          "TCP_UDP",
					"server_group_type": "Instance",
					"server_group_name": name,
					"vpc_id":            "${alicloud_vpc.defaultEyoN1M.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "HTTP",
							"health_check_connect_port":    "0",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "20",
							"health_check_interval":        "10",
							"health_check_domain":          "example.com",
							"health_check_url":             "/rdk",
							"health_check_http_code": []string{
								"http_2xx", "http_3xx", "http_4xx"},
							"http_check_method":         "HEAD",
							"health_check_http_version": "HTTP1.0",
						},
					},
					"scheduler":                "Wrr",
					"connection_drain_enabled": "true",
					"connection_drain_timeout": "20",
					"address_ip_version":       "DualStack",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"any_port_enabled":         "false",
						"protocol":                 "TCP_UDP",
						"server_group_type":        "Instance",
						"server_group_name":        name,
						"vpc_id":                   acctest.CHECKSET,
						"scheduler":                "Wrr",
						"connection_drain_timeout": "20",
						"address_ip_version":       "DualStack",
						"resource_group_id":        acctest.CHECKSET,
						"tags.%":                   "2",
						"tags.Created":             "TF",
						"tags.For":                 "Test",
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

// Case 5353  twin

func TestAccAliCloudNlbServerGroup_basic5353_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMap5353)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependence5353)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_enabled": "true",
					"any_port_enabled":         "true",
					"protocol":                 "UDP",
					"server_group_type":        "Ip",
					"server_group_name":        name,
					"vpc_id":                   "${alicloud_vpc.sgtestnlb.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "UDP",
							"health_check_connect_port":    "1",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "1",
							"health_check_interval":        "5",
							"health_check_domain":          "$SERVER_IP",
							"health_check_url":             "/rdktest",
							"health_check_http_code": []string{
								"http_2xx"},
							"http_check_method": "HEAD",
							"health_check_req":  "example",
							"health_check_exp":  "example",
						},
					},
					"connection_drain_timeout": "10",
					"address_ip_version":       "Ipv4",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain_enabled": "true",
						"any_port_enabled":         "true",
						"protocol":                 "UDP",
						"server_group_type":        "Ip",
						"server_group_name":        name,
						"vpc_id":                   acctest.CHECKSET,
						"connection_drain_timeout": "10",
						"address_ip_version":       "Ipv4",
						"tags.%":                   "2",
						"tags.Created":             "TF",
						"tags.For":                 "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "UDP",
							"health_check_connect_port":    "1",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "1",
							"health_check_interval":        "5",
							"health_check_domain":          "$SERVER_IP",
							"health_check_url":             "/rdktest",
							"health_check_http_code": []string{
								"http_2xx"},
							"http_check_method": "HEAD",
							"health_check_req":  "example-change",
							"health_check_exp":  "example-change",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

// Case IpVersionAffinityMode

func TestAccAliCloudNlbServerGroup_basicIpVersionAffinityMode(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMapIpVersionAffinityMode)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependenceIpVersionAffinityMode)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name":        name,
					"vpc_id":                   "${alicloud_vpc.default.id}",
					"address_ip_version":       "DualStack",
					"ip_version_affinity_mode": "NonAffinity",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name":        name,
						"vpc_id":                   acctest.CHECKSET,
						"address_ip_version":       "DualStack",
						"ip_version_affinity_mode": "NonAffinity",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_version_affinity_mode": "Affinity",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version_affinity_mode": "Affinity",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_version_affinity_mode": "NonAffinity",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version_affinity_mode": "NonAffinity",
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

// Case IpVersionAffinityMode twin

func TestAccAliCloudNlbServerGroup_basicIpVersionAffinityMode_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNlbServerGroupMapIpVersionAffinityMode)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNlbServerGroupBasicDependenceIpVersionAffinityMode)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"any_port_enabled":  "false",
					"protocol":          "TCP_UDP",
					"server_group_type": "Instance",
					"server_group_name": name,
					"vpc_id":            "${alicloud_vpc.default.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "HTTP",
							"health_check_connect_port":    "0",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "20",
							"health_check_interval":        "10",
							"health_check_domain":          "example.com",
							"health_check_url":             "/rdk",
							"health_check_http_code": []string{
								"http_2xx", "http_3xx", "http_4xx"},
							"http_check_method":         "HEAD",
							"health_check_http_version": "HTTP1.0",
						},
					},
					"scheduler":                "Wrr",
					"connection_drain_enabled": "true",
					"connection_drain_timeout": "20",
					"address_ip_version":       "DualStack",
					"ip_version_affinity_mode": "Affinity",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"any_port_enabled":         "false",
						"protocol":                 "TCP_UDP",
						"server_group_type":        "Instance",
						"server_group_name":        name,
						"vpc_id":                   acctest.CHECKSET,
						"scheduler":                "Wrr",
						"connection_drain_timeout": "20",
						"address_ip_version":       "DualStack",
						"ip_version_affinity_mode": "Affinity",
						"resource_group_id":        acctest.CHECKSET,
						"tags.%":                   "2",
						"tags.Created":             "TF",
						"tags.For":                 "Test",
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

var AliCloudNlbServerGroupMapIpVersionAffinityMode = map[string]string{
	"status":                     acctest.CHECKSET,
	"address_ip_version":         acctest.CHECKSET,
	"ip_version_affinity_mode":   acctest.CHECKSET,
	"scheduler":                  acctest.CHECKSET,
	"health_check.#":             acctest.CHECKSET,
	"server_group_type":          acctest.CHECKSET,
	"preserve_client_ip_enabled": acctest.CHECKSET,
	"connection_drain_enabled":   acctest.CHECKSET,
	"connection_drain_timeout":   acctest.CHECKSET,
	"any_port_enabled":           acctest.CHECKSET,
	"protocol":                   acctest.CHECKSET,
}

func AliCloudNlbServerGroupBasicDependenceIpVersionAffinityMode(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

  enable_ipv6 = true
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}

`, name)
}

func TestAccAliCloudNlbServerGroup_healthCheckHttpCodeOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNLBServerGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNLBServerGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":   "true",
							"health_check_type":      "HTTP",
							"health_check_http_code": []string{"http_2xx", "http_3xx"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
						"health_check.0.health_check_http_code.#": "2",
						"health_check.0.health_check_http_code.0": "http_2xx",
						"health_check.0.health_check_http_code.1": "http_3xx",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":   "true",
							"health_check_type":      "HTTP",
							"health_check_http_code": []string{"http_3xx", "http_2xx"},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":   "true",
							"health_check_type":      "HTTP",
							"health_check_http_code": []string{"http_3xx", "http_2xx"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
						"health_check.0.health_check_http_code.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":   "true",
							"health_check_type":      "HTTP",
							"health_check_http_code": []string{"http_2xx", "http_3xx", "http_4xx"},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":   "true",
							"health_check_type":      "HTTP",
							"health_check_http_code": []string{"http_2xx", "http_3xx", "http_4xx"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
						"health_check.0.health_check_http_code.#": "3",
						"health_check.0.health_check_http_code.0": "http_2xx",
						"health_check.0.health_check_http_code.1": "http_3xx",
						"health_check.0.health_check_http_code.2": "http_4xx",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudNlbServerGroup_healthCheckHttpCodeOrderCreation(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNLBServerGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNLBServerGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":   "true",
							"health_check_type":      "HTTP",
							"health_check_http_code": []string{"http_3xx", "http_2xx"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
						"health_check.0.health_check_http_code.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":   "true",
							"health_check_type":      "HTTP",
							"health_check_http_code": []string{"http_3xx", "http_2xx"},
						},
					},
				}),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
					}),
				),
			},
		},
	})
}
