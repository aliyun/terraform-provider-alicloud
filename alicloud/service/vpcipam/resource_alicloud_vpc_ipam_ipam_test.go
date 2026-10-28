package vpcipam_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpcipam"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

func init() {
	resource.AddTestSweepers("alicloud_vpc_ipam_ipam", &resource.Sweeper{
		Name: "alicloud_vpc_ipam_ipam",
		F:    testSweepVpcIpamIpam,
		Dependencies: []string{
			"alicloud_vpc_ipam_ipam_pool",
		},
	})
}

func testSweepVpcIpamIpam(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tfacc",
		"testAcc",
		"default",
	}

	ipamIds := make([]string, 0)
	action := "ListIpams"
	var response map[string]interface{}
	request := map[string]interface{}{
		"MaxResults": helper.PageSizeLarge,
		"RegionId":   client.RegionId,
	}
	for {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve VPC IPAM in service list: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Ipams", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Ipams", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			skip := true
			item := v.(map[string]interface{})
			if !sweep.SweepAll() {
				// Get ipam name, handle nil case
				ipamName := ""
				if item["IpamName"] != nil {
					ipamName = fmt.Sprint(item["IpamName"])
				}
				// Skip if name is empty or doesn't match prefix
				if ipamName == "" {
					log.Printf("[INFO] Skipping VPC IPAM with empty name: (%v)", item["IpamId"])
					continue
				}
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(ipamName), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping VPC IPAM: %v (%v)", ipamName, item["IpamId"])
					continue
				}
			}
			ipamIds = append(ipamIds, fmt.Sprint(item["IpamId"]))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	for _, id := range ipamIds {
		log.Printf("[INFO] Deleting VPC IPAM: (%s)", id)
		action := "DeleteIpam"
		request := map[string]interface{}{
			"IpamId":   id,
			"RegionId": client.RegionId,
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(time.Minute*10, func() *retry.RetryError {
			response, err = client.RpcPost("VpcIpam", "2023-02-28", action, nil, request, false)
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
			log.Printf("[ERROR] Failed to delete VPC IPAM (%s): %v", id, err)
			continue
		}
	}
	return nil
}

// Test VpcIpam Ipam. >>> Resource test cases, automatically generated.
// Case test_ipam_20250115 10035
func TestAccAliCloudVpcIpamIpam_basic10035(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpcIpamIpamMap10035)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vpcipam.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcIpamIpam")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpcipam%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpcIpamIpamBasicDependence10035)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_description":  "This is my first Ipam.",
					"ipam_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my first Ipam.",
						"ipam_name":               name,
						"resource_group_id":       acctest.CHECKSET,
						"operating_region_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_description":  "This is my new ipam.",
					"ipam_name":         name + "_update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"operating_region_list": []string{
						"cn-hangzhou", "cn-beijing", "cn-qingdao"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my new ipam.",
						"ipam_name":               name + "_update",
						"resource_group_id":       acctest.CHECKSET,
						"operating_region_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":       acctest.CHECKSET,
						"operating_region_list.#": "1",
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

var AlicloudVpcIpamIpamMap10035 = map[string]string{
	"status":                   acctest.CHECKSET,
	"create_time":              acctest.CHECKSET,
	"region_id":                acctest.CHECKSET,
	"private_default_scope_id": acctest.CHECKSET,
	"public_default_scope_id":  acctest.CHECKSET,
}

func AlicloudVpcIpamIpamBasicDependence10035(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case test_ipam 7856
func TestAccAliCloudVpcIpamIpam_basic7856(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpcIpamIpamMap7856)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vpcipam.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcIpamIpam")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpcipam%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpcIpamIpamBasicDependence7856)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_description":  "This is my first Ipam.",
					"ipam_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my first Ipam.",
						"ipam_name":               name,
						"resource_group_id":       acctest.CHECKSET,
						"operating_region_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_description":  "This is my new ipam.",
					"ipam_name":         name + "_update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"operating_region_list": []string{
						"cn-hangzhou", "cn-beijing", "cn-qingdao"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my new ipam.",
						"ipam_name":               name + "_update",
						"resource_group_id":       acctest.CHECKSET,
						"operating_region_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":       acctest.CHECKSET,
						"operating_region_list.#": "1",
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

var AlicloudVpcIpamIpamMap7856 = map[string]string{
	"status":                   acctest.CHECKSET,
	"create_time":              acctest.CHECKSET,
	"region_id":                acctest.CHECKSET,
	"private_default_scope_id": acctest.CHECKSET,
	"public_default_scope_id":  acctest.CHECKSET,
}

func AlicloudVpcIpamIpamBasicDependence7856(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test VpcIpam Ipam. <<< Resource test cases, automatically generated.
