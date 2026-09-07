package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ExpressConnect TrafficQosQueue. >>> Resource test cases, automatically generated.
// Case QoS队列-高优先级-线上 6831
func TestAccAliCloudExpressConnectTrafficQosQueue_basic6831(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_queue.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosQueueMap6831)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosQueue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosqueue%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosQueueBasicDependence6831)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":     "${alicloud_express_connect_traffic_qos.QoSCreate.id}",
					"queue_type": "High",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":     CHECKSET,
						"queue_type": "High",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"queue_description": "meijian-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue_description": "meijian-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"queue_name": "meijian-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue_name": "meijian-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            "${alicloud_express_connect_traffic_qos.QoSCreate.id}",
					"queue_description": "meijian-test",
					"queue_name":        "meijian-test",
					"queue_type":        "High",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            CHECKSET,
						"queue_description": "meijian-test",
						"queue_name":        "meijian-test",
						"queue_type":        "High",
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

var AlicloudExpressConnectTrafficQosQueueMap6831 = map[string]string{
	"status":   CHECKSET,
	"queue_id": CHECKSET,
}

func AlicloudExpressConnectTrafficQosQueueBasicDependence6831(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_express_connect_traffic_qos" "QoSCreate" {
  qos_name        = "meijian-test"
  qos_description = "meijian-test"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "preserved-NODELETING"
}

resource "alicloud_express_connect_traffic_qos_association" "Ass" {
  instance_id   = data.alicloud_express_connect_physical_connections.default.ids.0
  qos_id        = alicloud_express_connect_traffic_qos.QoSCreate.id
  instance_type = "PHYSICALCONNECTION"
}


`, name)
}

// Case QoS队列-线上 6832
func TestAccAliCloudExpressConnectTrafficQosQueue_basic6832(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_queue.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosQueueMap6832)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosQueue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosqueue%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosQueueBasicDependence6832)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            "${alicloud_express_connect_traffic_qos.QoSCreate.id}",
					"queue_type":        "Medium",
					"bandwidth_percent": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":     CHECKSET,
						"queue_type": "Medium",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_percent": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_percent": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"queue_description": "meijian-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue_description": "meijian-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"queue_name": "meijian-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue_name": "meijian-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_percent": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_percent": "40",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"queue_description": "meijian-test-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue_description": "meijian-test-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"queue_name": "meijian-test-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue_name": "meijian-test-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            "${alicloud_express_connect_traffic_qos.QoSCreate.id}",
					"bandwidth_percent": "60",
					"queue_description": "meijian-test",
					"queue_name":        "meijian-test",
					"queue_type":        "Medium",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            CHECKSET,
						"bandwidth_percent": "60",
						"queue_description": "meijian-test",
						"queue_name":        "meijian-test",
						"queue_type":        "Medium",
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

var AlicloudExpressConnectTrafficQosQueueMap6832 = map[string]string{
	"status":   CHECKSET,
	"queue_id": CHECKSET,
}

func AlicloudExpressConnectTrafficQosQueueBasicDependence6832(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_express_connect_traffic_qos" "QoSCreate" {
  qos_name        = "meijian-test"
  qos_description = "meijian-test"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "preserved-NODELETING"
}

resource "alicloud_express_connect_traffic_qos_association" "Ass" {
  instance_id   = data.alicloud_express_connect_physical_connections.default.ids.1
  qos_id        = alicloud_express_connect_traffic_qos.QoSCreate.id
  instance_type = "PHYSICALCONNECTION"
}


`, name)
}

// Case QoS队列-高优先级-线上 6831  twin
func TestAccAliCloudExpressConnectTrafficQosQueue_basic6831_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_queue.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosQueueMap6831)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosQueue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosqueue%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosQueueBasicDependence6831)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            "${alicloud_express_connect_traffic_qos.QoSCreate.id}",
					"queue_description": "meijian-test",
					"queue_name":        "meijian-test",
					"queue_type":        "High",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            CHECKSET,
						"queue_description": "meijian-test",
						"queue_name":        "meijian-test",
						"queue_type":        "High",
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

// Case QoS队列-线上 6832  twin
func TestAccAliCloudExpressConnectTrafficQosQueue_basic6832_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_queue.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosQueueMap6832)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosQueue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosqueue%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosQueueBasicDependence6832)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            "${alicloud_express_connect_traffic_qos.QoSCreate.id}",
					"bandwidth_percent": "60",
					"queue_description": "meijian-test",
					"queue_name":        "meijian-test",
					"queue_type":        "Medium",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            CHECKSET,
						"bandwidth_percent": "60",
						"queue_description": "meijian-test",
						"queue_name":        "meijian-test",
						"queue_type":        "Medium",
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

// Case QoS队列-高优先级-线上 6831  raw
func TestAccAliCloudExpressConnectTrafficQosQueue_basic6831_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_queue.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosQueueMap6831)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosQueue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosqueue%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosQueueBasicDependence6831)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            "${alicloud_express_connect_traffic_qos_association.Ass.qos_id}",
					"queue_description": "meijian-test",
					"queue_name":        "meijian-test",
					"queue_type":        "High",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            CHECKSET,
						"queue_description": "meijian-test",
						"queue_name":        "meijian-test",
						"queue_type":        "High",
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

// Case QoS队列-线上 6832  raw
func TestAccAliCloudExpressConnectTrafficQosQueue_basic6832_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_traffic_qos_queue.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectTrafficQosQueueMap6832)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectTrafficQosQueue")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnecttrafficqosqueue%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectTrafficQosQueueBasicDependence6832)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"qos_id":            "${alicloud_express_connect_traffic_qos.QoSCreate.id}",
					"bandwidth_percent": "60",
					"queue_description": "meijian-test",
					"queue_name":        "meijian-test",
					"queue_type":        "Medium",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            CHECKSET,
						"bandwidth_percent": "60",
						"queue_description": "meijian-test",
						"queue_name":        "meijian-test",
						"queue_type":        "Medium",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_percent": "40",
					"queue_description": "meijian-test-1",
					"queue_name":        "meijian-test-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_percent": "40",
						"queue_description": "meijian-test-1",
						"queue_name":        "meijian-test-1",
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

// Test ExpressConnect TrafficQosQueue. <<< Resource test cases, automatically generated.
