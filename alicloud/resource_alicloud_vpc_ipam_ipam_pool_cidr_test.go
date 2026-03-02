package alicloud

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_vpc_ipam_ipam_pool_cidr", &resource.Sweeper{
		Name: "alicloud_vpc_ipam_ipam_pool_cidr",
		F:    testSweepVpcIpamIpamPoolCidr,
	})
}

func testSweepVpcIpamIpamPoolCidr(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// First, get all IPAM Pools
	poolAction := "ListIpamPools"
	var poolResponse map[string]interface{}
	poolRequest := map[string]interface{}{
		"MaxResults": PageSizeLarge,
		"RegionId":   client.RegionId,
	}

	poolIds := make([]string, 0)
	for {
		poolResponse, err = client.RpcPost("VpcIpam", "2023-02-28", poolAction, nil, poolRequest, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve VPC IPAM Pool list: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.IpamPools", poolResponse)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, poolAction, "$.IpamPools", poolResponse)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			poolIds = append(poolIds, fmt.Sprint(item["IpamPoolId"]))
		}
		if nextToken, ok := poolResponse["NextToken"].(string); ok && nextToken != "" {
			poolRequest["NextToken"] = nextToken
		} else {
			break
		}
	}

	// For each pool, list and delete all CIDRs
	for _, poolId := range poolIds {
		action := "ListIpamPoolCidrs"
		var response map[string]interface{}
		request := map[string]interface{}{
			"MaxResults": PageSizeLarge,
			"IpamPoolId": poolId,
			"RegionId":   client.RegionId,
		}

		for {
			response, err = client.RpcPost("VpcIpam", "2023-02-28", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to retrieve VPC IPAM Pool CIDR list for pool %s: %s", poolId, err)
				break
			}
			resp, err := jsonpath.Get("$.IpamPoolCidrs", response)
			if err != nil {
				log.Printf("[WARN] Failed to get IpamPoolCidrs for pool %s: %v", poolId, err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				cidr := fmt.Sprint(item["Cidr"])
				log.Printf("[INFO] Deleting VPC IPAM Pool CIDR: %s (Pool: %s)", cidr, poolId)

				deleteAction := "DeleteIpamPoolCidr"
				deleteRequest := map[string]interface{}{
					"IpamPoolId": poolId,
					"Cidr":       cidr,
					"RegionId":   client.RegionId,
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(time.Minute*10, func() *resource.RetryError {
					response, err = client.RpcPost("VpcIpam", "2023-02-28", deleteAction, nil, deleteRequest, false)
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
					log.Printf("[ERROR] Failed to delete VPC IPAM Pool CIDR %s (Pool: %s): %v", cidr, poolId, err)
					continue
				}
			}
			if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
				request["NextToken"] = nextToken
			} else {
				break
			}
		}
	}
	return nil
}

// Test VpcIpam IpamPoolCidr. >>> Resource test cases, automatically generated.
// Case test_ipam_sub_pool_cidr 10812
func TestAccAliCloudVpcIpamIpamPoolCidr_basic10812(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_pool_cidr.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamPoolCidrMap10812)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpamPoolCidr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpcipam%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamPoolCidrBasicDependence10812)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"netmask_length": "24",
					"ipam_pool_id":   "${alicloud_vpc_ipam_ipam_pool.subIpamPool.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"netmask_length": "24",
						"ipam_pool_id":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"netmask_length"},
			},
		},
	})
}

var AlicloudVpcIpamIpamPoolCidrMap10812 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVpcIpamIpamPoolCidrBasicDependence10812(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ipam_scope_id = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  ip_version    = "IPv4"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "defaultIpamPoolCidr" {
  cidr         = "10.0.0.0/8"
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id
}

resource "alicloud_vpc_ipam_ipam_pool" "subIpamPool" {
  depends_on = ["alicloud_vpc_ipam_ipam_pool_cidr.defaultIpamPoolCidr"]
  ipam_scope_id       = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id      = alicloud_vpc_ipam_ipam.defaultIpam.region_id
  ip_version          = "IPv4"
  source_ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id
}


`, name)
}

// Case test_ipam_pool_cidr_ipv6 11310
func TestAccAliCloudVpcIpamIpamPoolCidr_basic11310(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_pool_cidr.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamPoolCidrMap11310)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpamPoolCidr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpcipam%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamPoolCidrBasicDependence11310)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_pool_id":   "${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}",
					"netmask_length": "56",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_pool_id":   CHECKSET,
						"netmask_length": "56",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"netmask_length"},
			},
		},
	})
}

var AlicloudVpcIpamIpamPoolCidrMap11310 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVpcIpamIpamPoolCidrBasicDependence11310(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ip_version     = "IPv6"
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.public_default_scope_id
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
}


`, name)
}

// Case test_ipam_pool_cidr 8028
func TestAccAliCloudVpcIpamIpamPoolCidr_basic8028(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_pool_cidr.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamPoolCidrMap8028)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpamPoolCidr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpcipam%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamPoolCidrBasicDependence8028)
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
					"cidr":         "10.0.0.0/8",
					"ipam_pool_id": "${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr":         "10.0.0.0/8",
						"ipam_pool_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"netmask_length"},
			},
		},
	})
}

var AlicloudVpcIpamIpamPoolCidrMap8028 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVpcIpamIpamPoolCidrBasicDependence8028(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
  ip_version     = "IPv4"
}


`, name)
}

// Test VpcIpam IpamPoolCidr. <<< Resource test cases, automatically generated.
