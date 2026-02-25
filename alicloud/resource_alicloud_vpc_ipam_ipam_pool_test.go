package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_vpc_ipam_ipam_pool", &resource.Sweeper{
		Name: "alicloud_vpc_ipam_ipam_pool",
		F:    testSweepVpcIpamIpamPool,
	})
}

func testSweepVpcIpamIpamPool(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test",
		"tf-test",
		"tfacc",
		"testAcc",
		"default",
	}

	ipamPoolIds := make([]string, 0)
	action := "ListIpamPools"
	var response map[string]interface{}
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
		"RegionId":   client.RegionId,
	}
	for {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve VPC IPAM Pool in service list: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.IpamPools", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.IpamPools", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			skip := true
			item := v.(map[string]interface{})
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["IpamPoolName"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping VPC IPAM Pool: %v (%v)", item["IpamPoolName"], item["IpamPoolId"])
					continue
				}
			}
			ipamPoolIds = append(ipamPoolIds, fmt.Sprint(item["IpamPoolId"]))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	for _, id := range ipamPoolIds {
		log.Printf("[INFO] Deleting VPC IPAM Pool: (%s)", id)
		action := "DeleteIpamPool"
		request := map[string]interface{}{
			"IpamPoolId": id,
			"RegionId":   client.RegionId,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(time.Minute*10, func() *resource.RetryError {
			response, err = client.RpcPost("VpcIpam", "2023-02-28", action, nil, request, false)
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
			log.Printf("[ERROR] Failed to delete VPC IPAM Pool (%s): %v", id, err)
			continue
		}
	}
	return nil
}

// Test VpcIpam IpamPool. >>> Resource test cases, automatically generated.
// Case test_parent_ipam_pool 8374
func TestAccAliCloudVpcIpamIpamPool_basic8374(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamPoolMap8374)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpamPool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipampool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamPoolBasicDependence8374)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_scope_id":       "${alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id}",
					"pool_region_id":      "${alicloud_vpc_ipam_ipam_pool.parentIpamPool.pool_region_id}",
					"ipam_pool_name":      name,
					"source_ipam_pool_id": "${alicloud_vpc_ipam_ipam_pool.parentIpamPool.id}",
					"ip_version":          "IPv4",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_scope_id":       CHECKSET,
						"pool_region_id":      CHECKSET,
						"ipam_pool_name":      name,
						"source_ipam_pool_id": CHECKSET,
						"ip_version":          "IPv4",
						"resource_group_id":   CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"clear_allocation_default_cidr_mask", "region_id"},
			},
		},
	})
}

var AlicloudVpcIpamIpamPoolMap8374 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudVpcIpamIpamPoolBasicDependence8374(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "parentIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = "cn-hangzhou"
}


`, name)
}

// Case test_ipam_pool 8026
func TestAccAliCloudVpcIpamIpamPool_basic8026(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamPoolMap8026)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpamPool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamPoolBasicDependence8026)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_scope_id":                "${alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id}",
					"ipam_pool_description":        "This is my ipam pool.",
					"ipam_pool_name":               name,
					"ip_version":                   "IPv4",
					"allocation_default_cidr_mask": "20",
					"allocation_min_cidr_mask":     "16",
					"allocation_max_cidr_mask":     "24",
					"pool_region_id":               "cn-hangzhou",
					"auto_import":                  "true",
					"resource_group_id":            "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_scope_id":                CHECKSET,
						"ipam_pool_description":        "This is my ipam pool.",
						"ipam_pool_name":               name,
						"ip_version":                   "IPv4",
						"allocation_default_cidr_mask": "20",
						"allocation_min_cidr_mask":     "16",
						"allocation_max_cidr_mask":     "24",
						"pool_region_id":               "cn-hangzhou",
						"auto_import":                  "true",
						"resource_group_id":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_pool_description":              "This is my new ipam pool description.",
					"ipam_pool_name":                     name + "_update",
					"allocation_default_cidr_mask":       "24",
					"allocation_min_cidr_mask":           "12",
					"allocation_max_cidr_mask":           "26",
					"auto_import":                        "false",
					"resource_group_id":                  "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"clear_allocation_default_cidr_mask": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_pool_description":              "This is my new ipam pool description.",
						"ipam_pool_name":                     name + "_update",
						"allocation_default_cidr_mask":       "24",
						"allocation_min_cidr_mask":           "12",
						"allocation_max_cidr_mask":           "26",
						"auto_import":                        "false",
						"resource_group_id":                  CHECKSET,
						"clear_allocation_default_cidr_mask": "false",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"clear_allocation_default_cidr_mask", "region_id"},
			},
		},
	})
}

var AlicloudVpcIpamIpamPoolMap8026 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudVpcIpamIpamPoolBasicDependence8026(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
  ipam_name             = var.name
}


`, name)
}

// Test VpcIpam IpamPool. <<< Resource test cases, automatically generated.
