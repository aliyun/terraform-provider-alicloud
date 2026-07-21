package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection ImageEventOperation. >>> Resource test cases, automatically generated.
// Case 4600
func TestAccAliCloudThreatDetectionImageEventOperation_basic4600(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_image_event_operation.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionImageEventOperationMap4600)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionImageEventOperation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionimageeventoperation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionImageEventOperationBasicDependence4600)
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
					"event_type":     "maliciousFile",
					"operation_code": "whitelist",
					"conditions":     "[{\\\"condition\\\":\\\"MD5\\\",\\\"type\\\":\\\"equals\\\",\\\"value\\\":\\\"0083a31cc0083a31ccf7c10367a6e783e\\\"}]",
					"scenarios":      "{\\\"type\\\":\\\"repo\\\",\\\"value\\\":\\\"test/repo-01\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_type":     "maliciousFile",
						"operation_code": "whitelist",
						"conditions":     "[{\"condition\":\"MD5\",\"type\":\"equals\",\"value\":\"0083a31cc0083a31ccf7c10367a6e783e\"}]",
						"scenarios":      "{\"type\":\"repo\",\"value\":\"test/repo-01\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"conditions": "[{\\\"condition\\\":\\\"MD5\\\",\\\"type\\\":\\\"equals\\\",\\\"value\\\":\\\"cb1ddb97bad0c443b438e6d013a6de6f\\\"}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"conditions": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scenarios": "{\\\"type\\\":\\\"default\\\",\\\"value\\\":\\\"\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scenarios": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"note": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"note": "test",
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

var AliCloudThreatDetectionImageEventOperationMap4600 = map[string]string{
	"source": CHECKSET,
}

func AliCloudThreatDetectionImageEventOperationBasicDependence4600(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4588
func TestAccAliCloudThreatDetectionImageEventOperation_basic4588(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_image_event_operation.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionImageEventOperationMap4588)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionImageEventOperation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionimageeventoperation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionImageEventOperationBasicDependence4588)
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
					"event_type":     "sensitiveFile",
					"operation_code": "whitelist",
					"conditions":     "[{\\\"condition\\\":\\\"MD5\\\",\\\"type\\\":\\\"equals\\\",\\\"value\\\":\\\"0083a31cc0083a31ccf7c10367a6e783e\\\"}]",
					"source":         "agentless",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_type":     "sensitiveFile",
						"operation_code": "whitelist",
						"conditions":     "[{\"condition\":\"MD5\",\"type\":\"equals\",\"value\":\"0083a31cc0083a31ccf7c10367a6e783e\"}]",
						"source":         "agentless",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"conditions": "[{\\\"condition\\\":\\\"MD5\\\",\\\"type\\\":\\\"equals\\\",\\\"value\\\":\\\"cb1ddb97bad0c443b438e6d013a6de6f\\\"}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"conditions": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scenarios": "{\\\"type\\\":\\\"repo\\\",\\\"value\\\":\\\"test/repo-01\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scenarios": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"note": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"note": "test",
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

var AliCloudThreatDetectionImageEventOperationMap4588 = map[string]string{
	"scenarios": CHECKSET,
}

func AliCloudThreatDetectionImageEventOperationBasicDependence4588(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4600  twin
func TestAccAliCloudThreatDetectionImageEventOperation_basic4600_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_image_event_operation.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionImageEventOperationMap4600)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionImageEventOperation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionimageeventoperation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionImageEventOperationBasicDependence4600)
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
					"event_type":     "maliciousFile",
					"operation_code": "whitelist",
					"event_key":      "huaweicloud_ak",
					"scenarios":      "{\\\"type\\\":\\\"repo\\\",\\\"value\\\":\\\"test/repo-01\\\"}",
					"event_name":     "华为AK",
					"conditions":     "[{\\\"condition\\\":\\\"MD5\\\",\\\"type\\\":\\\"equals\\\",\\\"value\\\":\\\"0083a31cc0083a31ccf7c10367a6e783e\\\"}]",
					"source":         "agentless",
					"note":           "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_type":     "maliciousFile",
						"operation_code": "whitelist",
						"event_key":      "huaweicloud_ak",
						"scenarios":      "{\"type\":\"repo\",\"value\":\"test/repo-01\"}",
						"event_name":     "华为AK",
						"conditions":     "[{\"condition\":\"MD5\",\"type\":\"equals\",\"value\":\"0083a31cc0083a31ccf7c10367a6e783e\"}]",
						"source":         "agentless",
						"note":           "test",
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

// Case 4588  twin
func TestAccAliCloudThreatDetectionImageEventOperation_basic4588_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_image_event_operation.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionImageEventOperationMap4588)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionImageEventOperation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sthreatdetectionimageeventoperation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionImageEventOperationBasicDependence4588)
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
					"event_type":     "sensitiveFile",
					"operation_code": "whitelist",
					"event_key":      "alibabacloud_ak",
					"scenarios":      "{\\\"type\\\":\\\"repo\\\",\\\"value\\\":\\\"test/repo-01\\\"}",
					"event_name":     "阿里云AK",
					"conditions":     "[{\\\"condition\\\":\\\"MD5\\\",\\\"type\\\":\\\"equals\\\",\\\"value\\\":\\\"0083a31cc0083a31ccf7c10367a6e783e\\\"}]",
					"source":         "agentless",
					"note":           "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_type":     "sensitiveFile",
						"operation_code": "whitelist",
						"event_key":      "alibabacloud_ak",
						"scenarios":      "{\"type\":\"repo\",\"value\":\"test/repo-01\"}",
						"event_name":     "阿里云AK",
						"conditions":     "[{\"condition\":\"MD5\",\"type\":\"equals\",\"value\":\"0083a31cc0083a31ccf7c10367a6e783e\"}]",
						"source":         "agentless",
						"note":           "test",
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

// Test ThreatDetection ImageEventOperation. <<< Resource test cases, automatically generated.
