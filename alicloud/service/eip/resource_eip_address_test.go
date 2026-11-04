package eip_test

// collection-order: bypass security_protection_types -- pre-existing ForceNew TypeList without DiffSuppressFunc; reordering members triggers resource recreation (non-convergent), and the attribute is unchanged by this migration PR.

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eip"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// eipTestClientForRegion delegates to sweep.SharedClientForRegion so the EIP
// test suite reuses the shared, region-isolated sweeper client (account/STS
// fallback and signature handled there) instead of a package-local copy.
func eipTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	raw, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return nil, err
	}
	return raw.(*connectivity.AliyunClient), nil
}

func init() {
	resource.AddTestSweepers("alicloud_eip_address", &resource.Sweeper{
		Name: "alicloud_eip_address",
		F:    testSweepEipAddress,
	})
}

func testSweepEipAddress(region string) error {
	rawClient, err := eipTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeEipAddresses"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
	}
	addressIds := make([]string, 0)
	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Println("List Eip Address Failed!", err)
			return nil
		}
		resp, err := jsonpath.Get("$.EipAddresses.EipAddress", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.EipAddresses.EipAddress", response)
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
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	for _, addressId := range addressIds {
		log.Printf("[INFO] Deleting Eip Address: (%s)", addressId)
		action := "ReleaseEipAddress"
		request := map[string]interface{}{
			"AllocationId": addressId,
		}
		request["RegionId"] = client.RegionId
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(time.Minute*9, func() *retry.RetryError {
			_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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

func TestAccAliCloudEIPAddress_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"resource_group_id": acctest.CHECKSET,
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
						"resource_group_id":   acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

var AlicloudEIPAddressMap0 = map[string]string{
	"description":          "",
	"status":               acctest.CHECKSET,
	"activity_id":          acctest.NOSET,
	"auto_pay":             acctest.NOSET,
	"netmode":              acctest.CHECKSET,
	"period":               acctest.NOSET,
	"pricing_cycle":        acctest.NOSET,
	"bandwidth":            acctest.CHECKSET,
	"resource_group_id":    acctest.CHECKSET,
	"address_name":         acctest.CHECKSET,
	"isp":                  acctest.CHECKSET,
	"internet_charge_type": acctest.CHECKSET,
	"payment_type":         acctest.CHECKSET,
	"ip_address":           acctest.CHECKSET,
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

func TestAccAliCloudEIPAddress_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.VPCPublicIpAddressPoolCidrBlockSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"isp":                       "BGP",
					"address_name":              "${var.name}",
					"internet_charge_type":      "PayByTraffic",
					"payment_type":              "PayAsYouGo",
					"public_ip_address_pool_id": "${alicloud_vpc_public_ip_address_pool_cidr_block.default.public_ip_address_pool_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"isp":                       "BGP",
						"address_name":              name,
						"internet_charge_type":      "PayByTraffic",
						"payment_type":              "PayAsYouGo",
						"public_ip_address_pool_id": acctest.CHECKSET,
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
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

var AlicloudEIPAddressMap1 = map[string]string{
	"auto_pay":             acctest.NOSET,
	"period":               acctest.NOSET,
	"status":               acctest.CHECKSET,
	"bandwidth":            acctest.CHECKSET,
	"description":          "",
	"netmode":              acctest.CHECKSET,
	"resource_group_id":    acctest.CHECKSET,
	"activity_id":          acctest.NOSET,
	"pricing_cycle":        acctest.NOSET,
	"payment_type":         "PayAsYouGo",
	"isp":                  "BGP",
	"address_name":         acctest.CHECKSET,
	"internet_charge_type": "PayByTraffic",
	"ip_address":           acctest.CHECKSET,
}

func AlicloudEIPAddressBasicDependence1(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
		default = "%s"
	}

	resource "alicloud_vpc_public_ip_address_pool" "default" {
		public_ip_address_pool_name = var.name
	}

	resource "alicloud_vpc_public_ip_address_pool_cidr_block" "default" {
  		public_ip_address_pool_id = alicloud_vpc_public_ip_address_pool.default.id
  		cidr_block                = "47.118.126.0/25"
	}
`, name)
}

func TestAccAliCloudEIPAddress_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap2)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EipAddressBGPProSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"isp":  "BGP_PRO",
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"isp":  "BGP_PRO",
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "Update",
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
					"name":        name,
					"bandwidth":   "5",
					"description": "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"bandwidth":   "5",
						"description": "update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

var AlicloudEIPAddressMap2 = map[string]string{
	"bandwidth":            acctest.CHECKSET,
	"description":          "",
	"payment_type":         acctest.CHECKSET,
	"pricing_cycle":        acctest.NOSET,
	"auto_pay":             acctest.NOSET,
	"internet_charge_type": acctest.CHECKSET,
	"period":               acctest.NOSET,
	"status":               acctest.CHECKSET,
	"activity_id":          acctest.NOSET,
	"resource_group_id":    acctest.CHECKSET,
	"netmode":              acctest.CHECKSET,
	"isp":                  "BGP_PRO",
	"address_name":         acctest.CHECKSET,
	"ip_address":           acctest.CHECKSET,
}

func AlicloudEIPAddressBasicDependence2(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func SkipTestAccAliCloudEIPAddress_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap3)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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

func TestAccAliCloudEIPAddress_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap4)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithTime(t, []int{1})
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

func SkipTestAccAliCloudEIPAddress_basic5(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

func TestAccAliCloudEIPAddress_basic6(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

func TestAccAliCloudEIPAddress_basic7(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap4)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithTime(t, []int{1})
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

func TestAccAliCloudEIPAddress_basic8(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
					"high_definition_monitor_log_status": "ON",
					"log_project":                        "hdmonitor-cn-shenzhen",
					"log_store":                          "hdmonitor",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"high_definition_monitor_log_status": "ON",
						"log_project":                        "hdmonitor-cn-shenzhen",
						"log_store":                          "hdmonitor",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

func TestAccAliCloudEIPAddress_basic9(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.VpcDhcpOptionsSetSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name,
					"activity_id":               "12345",
					"bandwidth":                 "10",
					"description":               name,
					"instance_charge_type":      "PostPaid",
					"isp":                       "BGP",
					"netmode":                   "public",
					"zone":                      "cn-hangzhou-i",
					"security_protection_types": []string{"AntiDDoS_Enhanced"},
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
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
						"resource_group_id":    acctest.CHECKSET,
						"zone":                 "cn-hangzhou-i",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mode": "BINDED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mode": "BINDED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mode": "NAT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mode": "NAT",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

func TestAccAliCloudEIPAddress_basic13(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.VpcDhcpOptionsSetSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name,
					"bandwidth":                 "10",
					"description":               name,
					"instance_charge_type":      "PostPaid",
					"allocation_id":             "eip-bp1pktw6wwavjqudm4mon",
					"isp":                       "BGP",
					"netmode":                   "public",
					"security_protection_types": []string{"AntiDDoS_Enhanced"},
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"bandwidth":            "10",
						"description":          name,
						"instance_charge_type": "PostPaid",
						"isp":                  "BGP",
						"netmode":              "public",
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

func TestAccAliCloudEIPAddress_basic9_ip(t *testing.T) {
	t.Skipped()
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EipAddressZoneSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
					"ip_address":           "47.91.57.17",
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
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay"},
			},
		},
	})
}

func TestAccAliCloudEIPAddress_basic10(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEIPAddressMap4)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eip.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEipAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seipaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEIPAddressBasicDependence4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithTime(t, []int{1})
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":  "Subscription",
					"auto_pay":      "true",
					"period":        "1",
					"pricing_cycle": "Month",
					"address_name":  name,
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
				ImportStateVerifyIgnore: []string{"activity_id", "period", "auto_pay", "pricing_cycle"},
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

// lintignore: R001
func TestUnitAlicloudEIPAddress(t *testing.T) {
	p := acctest.Provider.ResourcesMap
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
	rawClient, err := eipTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
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
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_eip_address", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = p["alicloud_eip_address"].Create(dInit, rawClient)
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
		err := p["alicloud_eip_address"].Create(dInit, rawClient)
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
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = p["alicloud_eip_address"].Update(dExisted, rawClient)
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
	diff, err := acctest.NewInstanceDiff("alicloud_eip_address", attributes, attributesDiff, dInit.State())
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
		err := p["alicloud_eip_address"].Update(dExisted, rawClient)
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
	diff, err = acctest.NewInstanceDiff("alicloud_eip_address", attributes, attributesDiff, dExisted.State())
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
		err := p["alicloud_eip_address"].Update(dExisted, rawClient)
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
	diff, err = acctest.NewInstanceDiff("alicloud_eip_address", attributes, attributesDiff, dExisted.State())
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
		err := p["alicloud_eip_address"].Update(dExisted, rawClient)
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
		err := p["alicloud_eip_address"].Read(dExisted, rawClient)
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
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = p["alicloud_eip_address"].Delete(dExisted, rawClient)
	patches.Reset()
	assert.Nil(t, err)
	attributesDiff = map[string]interface{}{
		"payment_type": "PayAsYouGo",
	}
	diff, err = acctest.NewInstanceDiff("alicloud_eip_address", attributes, attributesDiff, dInit.State())
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
		err := p["alicloud_eip_address"].Delete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidAllocationId.NotFound":
			assert.Nil(t, err)
		}
	}

}
