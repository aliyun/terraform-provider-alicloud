package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

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
				ImportStateVerifyIgnore: []string{"activity_id", "netmode", "period"},
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
				ImportStateVerifyIgnore: []string{"activity_id", "netmode", "period"},
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
				ImportStateVerifyIgnore: []string{"activity_id", "netmode", "period"},
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
