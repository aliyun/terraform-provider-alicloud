package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudExpressConnectPhysicalConnection_domesic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectPhysicalConnectionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou, connectivity.Beijing, connectivity.Shanghai})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// currently， not all access points are available
					//"access_point_id":          "${data.alicloud_express_connect_access_points.default.ids.0}",
					"access_point_id":          getAccessPointId(),
					"type":                     "VPC",
					"peer_location":            "testacc12345",
					"physical_connection_name": "${var.name}",
					"description":              "${var.name}",
					"line_operator":            "CU",
					"port_type":                "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id":          CHECKSET,
						"type":                     "VPC",
						"peer_location":            "testacc12345",
						"physical_connection_name": name,
						"description":              name,
						"line_operator":            "CU",
						"port_type":                "1000Base-LX",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name + "_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"circuit_code": "longtel001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"circuit_code": "longtel001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line_operator": "CU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"line_operator": "CU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line_operator": "CM",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"line_operator": "CM",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line_operator": "CO",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"line_operator": "CO",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_location": "浙江省---vfjdbg_21e",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_location": "浙江省---vfjdbg_21e",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"port_type": "10GBase-LR",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"port_type": "10GBase-LR",
			//		}),
			//	),
			//},
			// Only confirmed connection can be enabled.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "Enabled",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"status": "Enabled",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name,
					"status":                   "Canceled",
					"bandwidth":                "15",
					"circuit_code":             "longtel002",
					"description":              name,
					"line_operator":            "CT",
					"peer_location":            "testacc12345",
					"port_type":                "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name,
						"status":                   "Canceled",
						"bandwidth":                "15",
						"circuit_code":             "longtel002",
						"description":              name,
						"line_operator":            "CT",
						"peer_location":            "testacc12345",
						"port_type":                "1000Base-LX",
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

func TestAccAlicloudExpressConnectPhysicalConnection_intl(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectPhysicalConnectionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.EUCentral1, connectivity.APSouthEast1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// currently， not all access points are available
					//"access_point_id": "${data.alicloud_express_connect_access_points.default.ids.0}",
					"access_point_id":          getAccessPointId(),
					"type":                     "VPC",
					"peer_location":            "testacc12345",
					"physical_connection_name": "${var.name}",
					"description":              "${var.name}",
					"line_operator":            "Other",
					"port_type":                "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id":          CHECKSET,
						"type":                     "VPC",
						"peer_location":            "testacc12345",
						"physical_connection_name": name,
						"description":              name,
						"line_operator":            "Other",
						"port_type":                "1000Base-LX",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name + "_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"circuit_code": "longtel001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"circuit_code": "longtel001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line_operator": "Equinix",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"line_operator": "Equinix",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_location": "国际---vfjdbg_21e",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_location": "国际---vfjdbg_21e",
					}),
				),
			},
			// Currently, the internal region does not support 10G
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"port_type": "10GBase-LR",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"port_type": "10GBase-LR",
			//		}),
			//	),
			//},
			// Only confirmed connection can be enabled.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "Enabled",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"status": "Enabled",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name,
					"status":                   "Canceled",
					"bandwidth":                "15",
					"circuit_code":             "longtel002",
					"description":              name,
					"line_operator":            "Other",
					"peer_location":            "testacc12345",
					"port_type":                "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name,
						"status":                   "Canceled",
						"bandwidth":                "15",
						"circuit_code":             "longtel002",
						"description":              name,
						"line_operator":            "Other",
						"peer_location":            "testacc12345",
						"port_type":                "1000Base-LX",
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

func TestAccAlicloudExpressConnectPhysicalConnection_domesic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectPhysicalConnectionBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou, connectivity.Beijing, connectivity.Shanghai})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// currently， not all access points are available
					//"access_point_id":          "${data.alicloud_express_connect_access_points.default.ids.0}",
					"access_point_id":                  getAccessPointId(),
					"redundant_physical_connection_id": "${data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id}",
					"type":                             "VPC",
					"peer_location":                    "testacc12345",
					"physical_connection_name":         name,
					"description":                      "${var.name}",
					"line_operator":                    "CU",
					"port_type":                        "10GBase-LR",
					"bandwidth":                        "10",
					"circuit_code":                     "longtel001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id":                  CHECKSET,
						"redundant_physical_connection_id": CHECKSET,
						"type":                             "VPC",
						"peer_location":                    "testacc12345",
						"physical_connection_name":         name,
						"description":                      name,
						"line_operator":                    "CU",
						"port_type":                        "10GBase-LR",
						"bandwidth":                        "10",
						"circuit_code":                     "longtel001",
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

var AlicloudExpressConnectPhysicalConnectionMap0 = map[string]string{
	"status":                           CHECKSET,
	"redundant_physical_connection_id": "",
	"bandwidth":                        CHECKSET,
}

func AlicloudExpressConnectPhysicalConnectionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_express_connect_access_points" "default" {
	status = "recommended"
}
`, name)
}

func AlicloudExpressConnectPhysicalConnectionBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_express_connect_access_points" "default" {
	status = "recommended"
}

data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

`, name)
}
