package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_eip_address", &resource.Sweeper{
		Name: "alicloud_eip_address",
		F:    testSweepEipAddress,
	})
}

func testSweepEipAddress(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	action := "DescribeEipAddresses"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	addressIds := make([]string, 0)
	var response map[string]interface{}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Println("List Eip Address Failed!", err)
			return nil
		}
		resp, err := jsonpath.Get("$.EipAddresses.EipAddress", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.EipAddresses.EipAddress", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["Name"])), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Eip Address: %v (%v)", item["Name"], item["AllocationId"])
				continue
			}
			addressIds = append(addressIds, fmt.Sprint(item["AllocationId"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	for _, addressId := range addressIds {
		log.Printf("[INFO] Deleting Eip Address: (%s)", addressId)
		action := "ReleaseEipAddress"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		request := map[string]interface{}{
			"AllocationId": addressId,
		}
		request["RegionId"] = client.RegionId
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(time.Minute*9, func() *resource.RetryError {
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed To Delete Eip Address : %s", err)
		}
		log.Printf("[INFO] Delete Eip Address Success : %s", addressId)
	}
	return nil
}

func TestAccAlicloudEIPAddress_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := resourceAttrInit(resourceId, AlicloudEIPAddressMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEIPAddressBasicDependence0)
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
					"address_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]interface{}{
						"Create": "tfTest",
						"For":    "tfTest 123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":      "2",
						"tags.Create": "tfTest",
						"tags.For":    "tfTest 123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test for terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test for terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_name":        "${var.name}",
					"bandwidth":           "3",
					"description":         "testForTerraform1234",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"deletion_protection": "false",
					"tags": map[string]interface{}{
						"Create": "tfTest Update",
						"For":    "tfTest 123 Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_name":        name,
						"bandwidth":           "3",
						"description":         "testForTerraform1234",
						"resource_group_id":   CHECKSET,
						"deletion_protection": "false",
						"tags.%":              "2",
						"tags.Create":         "tfTest Update",
						"tags.For":            "tfTest 123 Update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "netmode", "period", "auto_pay"},
			},
		},
	})
}

var AlicloudEIPAddressMap0 = map[string]string{
	"description":          "",
	"status":               CHECKSET,
	"activity_id":          NOSET,
	"auto_pay":             NOSET,
	"netmode":              NOSET,
	"period":               NOSET,
	"pricing_cycle":        NOSET,
	"bandwidth":            CHECKSET,
	"resource_group_id":    CHECKSET,
	"address_name":         CHECKSET,
	"isp":                  CHECKSET,
	"internet_charge_type": CHECKSET,
	"payment_type":         CHECKSET,
	"ip_address":           CHECKSET,
}

func AlicloudEIPAddressBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}
`, name)
}

func TestAccAlicloudEIPAddress_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := resourceAttrInit(resourceId, AlicloudEIPAddressMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEIPAddressBasicDependence1)
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
					"isp":                  "BGP",
					"address_name":         "${var.name}",
					"internet_charge_type": "PayByTraffic",
					"payment_type":         "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"isp":                  "BGP",
						"address_name":         name,
						"internet_charge_type": "PayByTraffic",
						"payment_type":         "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_name": name + "update",
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
					"description": "test for terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test for terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_name": name,
					"bandwidth":    "5",
					"description":  "test for terraform update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_name": name,
						"bandwidth":    "5",
						"description":  "test for terraform update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "netmode", "period", "auto_pay"},
			},
		},
	})
}

var AlicloudEIPAddressMap1 = map[string]string{
	"auto_pay":             NOSET,
	"period":               NOSET,
	"status":               CHECKSET,
	"bandwidth":            CHECKSET,
	"description":          "",
	"netmode":              NOSET,
	"resource_group_id":    CHECKSET,
	"activity_id":          NOSET,
	"pricing_cycle":        NOSET,
	"payment_type":         "PayAsYouGo",
	"isp":                  "BGP",
	"address_name":         CHECKSET,
	"internet_charge_type": "PayByTraffic",
	"ip_address":           CHECKSET,
}

func AlicloudEIPAddressBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestAccAlicloudEIPAddress_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := resourceAttrInit(resourceId, AlicloudEIPAddressMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEIPAddressBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EipAddressBGPProSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"isp":          "BGP_PRO",
					"address_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"isp":          "BGP_PRO",
						"address_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_name": name + "Update",
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
					"description": "test for terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test for terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_name": "${var.name}",
					"bandwidth":    "5",
					"description":  "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_name": name,
						"bandwidth":    "5",
						"description":  "update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "netmode", "period", "auto_pay"},
			},
		},
	})
}

var AlicloudEIPAddressMap2 = map[string]string{
	"bandwidth":            CHECKSET,
	"description":          "",
	"payment_type":         CHECKSET,
	"pricing_cycle":        NOSET,
	"auto_pay":             NOSET,
	"internet_charge_type": CHECKSET,
	"period":               NOSET,
	"status":               CHECKSET,
	"activity_id":          NOSET,
	"resource_group_id":    CHECKSET,
	"netmode":              NOSET,
	"isp":                  "BGP_PRO",
	"address_name":         CHECKSET,
	"ip_address":           CHECKSET,
}

func AlicloudEIPAddressBasicDependence2(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestAccAlicloudEIPAddress_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := resourceAttrInit(resourceId, AlicloudEIPAddressMap3)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEIPAddressBasicDependence3)
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
					"isp":                  "BGP",
					"address_name":         "${var.name}",
					"internet_charge_type": "PayByDominantTraffic",
					"payment_type":         "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"isp":                  "BGP",
						"address_name":         name,
						"internet_charge_type": "PayByDominantTraffic",
						"payment_type":         "PayAsYouGo",
					}),
				),
			},
		},
	})
}

var AlicloudEIPAddressMap3 = map[string]string{}

func AlicloudEIPAddressBasicDependence3(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestAccAlicloudEIPAddress_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := resourceAttrInit(resourceId, AlicloudEIPAddressMap4)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEIPAddressBasicDependence4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"auto_pay":     "true",
					"period":       "1",
					"address_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_name": name,
						"payment_type": "Subscription",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "netmode", "period", "auto_pay"},
			},
		},
	})
}

func TestAccAlicloudEIPAddress_basic5(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := resourceAttrInit(resourceId, AlicloudEIPAddressMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEIPAddressBasicDependence0)
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
					"address_name":         name,
					"activity_id":          "12345",
					"bandwidth":            "10",
					"description":          name,
					"internet_charge_type": "PayByDominantTraffic",
					"payment_type":         "PayAsYouGo",
					"isp":                  "BGP",
					"netmode":              "public",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_name":         name,
						"activity_id":          "12345",
						"bandwidth":            "10",
						"description":          name,
						"internet_charge_type": "PayByDominantTraffic",
						"payment_type":         "PayAsYouGo",
						"isp":                  "BGP",
						"netmode":              "public",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "netmode", "period", "auto_pay"},
			},
		},
	})
}

func TestAccAlicloudEIPAddress_basic6(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := resourceAttrInit(resourceId, AlicloudEIPAddressMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEIPAddressBasicDependence0)
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
					"name":                 name,
					"activity_id":          "12345",
					"bandwidth":            "10",
					"description":          name,
					"instance_charge_type": "PostPaid",
					"isp":                  "BGP",
					"netmode":              "public",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"activity_id":          "12345",
						"bandwidth":            "10",
						"description":          name,
						"instance_charge_type": "PostPaid",
						"isp":                  "BGP",
						"netmode":              "public",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "netmode", "period", "auto_pay"},
			},
		},
	})
}

func TestAccAlicloudEIPAddress_basic7(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := resourceAttrInit(resourceId, AlicloudEIPAddressMap4)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEIPAddressBasicDependence4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"auto_pay":     "true",
					"period":       "12",
					"address_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_name": name,
						"payment_type": "Subscription",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "netmode", "period", "auto_pay"},
			},
		},
	})
}

var AlicloudEIPAddressMap4 = map[string]string{}

func AlicloudEIPAddressBasicDependence4(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestUnitAlicloudEIPAddress(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_eip_address"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_eip_address"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"address_name":         "AllocateEipAddressValue",
		"activity_id":          "AllocateEipAddressValue",
		"bandwidth":            "10",
		"description":          "AllocateEipAddressValue",
		"internet_charge_type": "AllocateEipAddressValue",
		"payment_type":         "Subscription",
		"isp":                  "AllocateEipAddressValue",
		"netmode":              "AllocateEipAddressValue",
		"resource_group_id":    "AllocateEipAddressValue",
		"period":               12,
		"auto_pay":             false,
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeEipAddresses
		"EipAddresses": map[string]interface{}{
			"EipAddress": []interface{}{
				map[string]interface{}{
					"Name":               "AllocateEipAddressValue",
					"Bandwidth":          "10",
					"DeletionProtection": "AllocateEipAddressValue",
					"Descritpion":        "AllocateEipAddressValue",
					"ISP":                "AllocateEipAddressValue",
					"InternetChargeType": "AllocateEipAddressValue",
					"ChargeType":         "PrePaid",
					"IpAddress":          "AllocateEipAddressValue",
					"ResourceGroupId":    "AllocateEipAddressValue",
					"Status":             "Available",
					"AllocationId":       "AllocateEipAddressValue",
				},
			},
		},
		"TagResources": map[string]interface{}{
			"TagResource": []interface{}{
				map[string]interface{}{
					"TagKey":   "Create",
					"TagValue": "tfTest Update",
				},
				map[string]interface{}{
					"TagKey":   "For",
					"TagValue": "tfTest Update",
				},
			},
		},
		"AllocationId": "AllocateEipAddressValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateJobTemplate
		"AllocationId": "AllocateEipAddressValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_eip_address", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEipAddressCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeEipAddresses Response
		"EipAddresses": map[string]interface{}{
			"EipAddress": []interface{}{
				map[string]interface{}{
					"AllocationId": "AllocateEipAddressValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AllocateEipAddress" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEipAddressCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_eip_address"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEipAddressUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// DeletionProtection
	attributesDiff := map[string]interface{}{
		"tags": map[string]interface{}{
			"Create": "tfTest Update",
			"For":    "tfTest Update",
			"Test":   "tfTest Update",
		},
		"deletion_protection": true,
	}
	diff, err := newInstanceDiff("alicloud_eip_address", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_eip_address"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeEipAddresses Response
		"EipAddresses": map[string]interface{}{
			"EipAddress": []interface{}{
				map[string]interface{}{
					"DeletionProtection": true,
				},
			},
		},
		"TagResources": map[string]interface{}{
			"TagResource": []interface{}{
				map[string]interface{}{
					"TagKey":   "Create",
					"TagValue": "tfTest Update",
				},
				map[string]interface{}{
					"TagKey":   "For",
					"TagValue": "tfTest Update",
				},
				map[string]interface{}{
					"TagKey":   "Test",
					"TagValue": "tfTest Update",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeletionProtection" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEipAddressUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_eip_address"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// MoveResourceGroup
	attributesDiff = map[string]interface{}{
		"resource_group_id": "MoveResourceGroupValue",
	}
	diff, err = newInstanceDiff("alicloud_eip_address", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_eip_address"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeEipAddresses Response
		"EipAddresses": map[string]interface{}{
			"EipAddress": []interface{}{
				map[string]interface{}{
					"ResourceGroupId": "MoveResourceGroupValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "MoveResourceGroup" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEipAddressUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_eip_address"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyEipAddressAttribute
	attributesDiff = map[string]interface{}{
		"address_name": "ModifyEipAddressAttributeValue",
		"bandwidth":    "20",
		"description":  "ModifyEipAddressAttributeValue",
	}
	diff, err = newInstanceDiff("alicloud_eip_address", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_eip_address"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeEipAddresses Response
		"EipAddresses": map[string]interface{}{
			"EipAddress": []interface{}{
				map[string]interface{}{
					"Name":        "ModifyEipAddressAttributeValue",
					"Bandwidth":   "20",
					"Descritpion": "ModifyEipAddressAttributeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyEipAddressAttribute" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEipAddressUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_eip_address"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeEipAddresses" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEipAddressRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEipAddressDelete(dExisted, rawClient)
	patches.Reset()
	assert.Nil(t, err)
	attributesDiff = map[string]interface{}{
		"payment_type": "PayAsYouGo",
	}
	diff, err = newInstanceDiff("alicloud_eip_address", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_eip_address"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "IncorrectEipStatus", "TaskConflict.AssociateGlobalAccelerationInstance", "nil", "InvalidAllocationId.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ReleaseEipAddress" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"EipAddresses": map[string]interface{}{
								"EipAddress": []interface{}{},
							},
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEipAddressDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidAllocationId.NotFound":
			assert.Nil(t, err)
		}
	}

}
