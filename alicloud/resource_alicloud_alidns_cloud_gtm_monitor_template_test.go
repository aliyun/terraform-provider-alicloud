package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudAlidnsCloudGtmMonitorTemplate_basic12684(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_monitor_template.default"
	ra := resourceAttrInit(resourceId, AlicloudAlidnsCloudGtmMonitorTemplateMap12684)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsCloudGtmMonitorTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsCloudGtmMonitorTemplateBasicDependence12684)
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
					"ip_version": "IPv4",
					"timeout":    "2000",
					"isp_city_nodes": []map[string]interface{}{
						{
							"city_code": "357",
							"isp_code":  "465",
						},
						{
							"city_code": "738",
							"isp_code":  "465",
						},
					},
					"evaluation_count": "2",
					"protocol":         "http",
					"failure_rate":     "50",
					"extend_info":      "{\\\"code\\\":500,\\\"followRedirect\\\":true,\\\"path\\\":\\\"/\\\"}",
					"name":             name,
					"interval":         "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version":       "IPv4",
						"timeout":          CHECKSET,
						"isp_city_nodes.#": "2",
						"evaluation_count": "2",
						"protocol":         "http",
						"failure_rate":     "50",
						"extend_info":      CHECKSET,
						"name":             name,
						"interval":         CHECKSET,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudAlidnsCloudGtmMonitorTemplateMap12684 = map[string]string{}

func AlicloudAlidnsCloudGtmMonitorTemplateBasicDependence12684(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resourceCase_20260325_z3jf7s 12685
func TestAccAliCloudAlidnsCloudGtmMonitorTemplate_basic12685(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_monitor_template.default"
	ra := resourceAttrInit(resourceId, AlicloudAlidnsCloudGtmMonitorTemplateMap12685)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsCloudGtmMonitorTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsCloudGtmMonitorTemplateBasicDependence12685)
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
					"ip_version": "IPv6",
					"timeout":    "3000",
					"isp_city_nodes": []map[string]interface{}{
						{
							"city_code": "357",
							"isp_code":  "465",
						},
						{
							"city_code": "738",
							"isp_code":  "465",
						},
					},
					"evaluation_count": "2",
					"protocol":         "tcp",
					"failure_rate":     "50",
					"name":             name,
					"interval":         "60",
					"remark":           "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version":       "IPv6",
						"timeout":          CHECKSET,
						"isp_city_nodes.#": "2",
						"evaluation_count": "2",
						"protocol":         "tcp",
						"failure_rate":     "50",
						"name":             name,
						"interval":         CHECKSET,
						"remark":           "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"interval": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"interval": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"failure_rate": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"failure_rate": "20",
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

var AlicloudAlidnsCloudGtmMonitorTemplateMap12685 = map[string]string{}

func AlicloudAlidnsCloudGtmMonitorTemplateBasicDependence12685(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resourceCase_20260325_u2rOov 12690
func TestAccAliCloudAlidnsCloudGtmMonitorTemplate_basic12690(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_monitor_template.default"
	ra := resourceAttrInit(resourceId, AlicloudAlidnsCloudGtmMonitorTemplateMap12690)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsCloudGtmMonitorTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsCloudGtmMonitorTemplateBasicDependence12690)
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
					"ip_version": "IPv4",
					"timeout":    "2000",
					"isp_city_nodes": []map[string]interface{}{
						{
							"city_code": "357",
							"isp_code":  "465",
						},
						{
							"city_code": "738",
							"isp_code":  "465",
						},
					},
					"evaluation_count": "1",
					"protocol":         "ping",
					"failure_rate":     "20",
					"extend_info":      "{\\\"packetLossRate\\\":10,\\\"packetNum\\\":20}",
					"name":             name,
					"interval":         "60",
					"remark":           "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version":       "IPv4",
						"timeout":          CHECKSET,
						"isp_city_nodes.#": "2",
						"evaluation_count": "1",
						"protocol":         "ping",
						"failure_rate":     "20",
						"extend_info":      CHECKSET,
						"name":             name,
						"interval":         CHECKSET,
						"remark":           "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"isp_city_nodes": []map[string]interface{}{
						{
							"city_code": "304",
							"isp_code":  "5",
						},
						{
							"city_code": "738",
							"isp_code":  "465",
						},
						{
							"city_code": "304",
							"isp_code":  "465",
						},
					},
					"failure_rate": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"isp_city_nodes.#": "3",
						"failure_rate":     "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"extend_info": "{\\\"packetLossRate\\\":80,\\\"packetNum\\\":20}",
					"timeout":     "3000",
					"isp_city_nodes": []map[string]interface{}{
						{
							"city_code": "357",
							"isp_code":  "465",
						},
						{
							"city_code": "738",
							"isp_code":  "465",
						},
					},
					"evaluation_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout":          CHECKSET,
						"isp_city_nodes.#": "2",
						"evaluation_count": "2",
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

var AlicloudAlidnsCloudGtmMonitorTemplateMap12690 = map[string]string{}

func AlicloudAlidnsCloudGtmMonitorTemplateBasicDependence12690(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resourceCase_20260325_9q5Tgr 12691
func TestAccAliCloudAlidnsCloudGtmMonitorTemplate_basic12691(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_monitor_template.default"
	ra := resourceAttrInit(resourceId, AlicloudAlidnsCloudGtmMonitorTemplateMap12691)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsCloudGtmMonitorTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsCloudGtmMonitorTemplateBasicDependence12691)
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
					"ip_version": "IPv4",
					"timeout":    "3000",
					"isp_city_nodes": []map[string]interface{}{
						{
							"city_code": "357",
							"isp_code":  "465",
						},
						{
							"city_code": "738",
							"isp_code":  "465",
						},
					},
					"evaluation_count": "2",
					"protocol":         "https",
					"failure_rate":     "80",
					"extend_info":      "{\\\"code\\\":500,\\\"followRedirect\\\":true,\\\"path\\\":\\\"/\\\",\\\"sni\\\":true}",
					"name":             name,
					"interval":         "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_version":       "IPv4",
						"timeout":          CHECKSET,
						"isp_city_nodes.#": "2",
						"evaluation_count": "2",
						"protocol":         "https",
						"failure_rate":     "80",
						"extend_info":      CHECKSET,
						"name":             name,
						"interval":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"evaluation_count": "1",
					"failure_rate":     "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"evaluation_count": "1",
						"failure_rate":     "50",
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

var AlicloudAlidnsCloudGtmMonitorTemplateMap12691 = map[string]string{}

func AlicloudAlidnsCloudGtmMonitorTemplateBasicDependence12691(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
