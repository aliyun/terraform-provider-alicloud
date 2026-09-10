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
	resource.AddTestSweepers(
		"alicloud_vpc_dhcp_options_set",
		&resource.Sweeper{
			Name: "alicloud_vpc_dhcp_options_set",
			F:    testSweepVpcDhcpOptionsSet,
		})
}

func testSweepVpcDhcpOptionsSet(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListDhcpOptionsSets"
	request := map[string]interface{}{}
	request["MaxResults"] = PageSizeLarge
	request["RegionId"] = client.RegionId
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.DhcpOptionsSets", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DhcpOptionsSets", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			_, ok := item["DhcpOptionsSetName"]
			if !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DhcpOptionsSetName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Vpc Dhcp Options Set: %s", item["DhcpOptionsSetName"].(string))
				continue
			}
			action := "DeleteDhcpOptionsSet"
			request := map[string]interface{}{
				"DhcpOptionsSetId": item["DhcpOptionsSetId"],
			}
			request["RegionId"] = client.RegionId
			_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Vpc Dhcp Options Set (%s): %s", item["DhcpOptionsSetName"].(string), err)
			}
			log.Printf("[INFO] Delete Vpc Dhcp Options Set success: %s ", item["DhcpOptionsSetName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudVPCDhcpOptionsSet_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_dhcp_options_set.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCDhcpOptionsSetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcDhcpOptionsSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcdhcpoptionsset%d", defaultRegionToTest, rand)
	domainName := fmt.Sprintf("tftestacc%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCDhcpOptionsSetBasicDependence0)
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
					"dhcp_options_set_name":        "${var.name}",
					"dhcp_options_set_description": "${var.name}",
					"domain_name":                  domainName,
					"domain_name_servers":          "100.100.2.136",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_name":        name,
						"dhcp_options_set_description": name,
						"domain_name":                  domainName,
						"domain_name_servers":          "100.100.2.136",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dhcp_options_set_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dhcp_options_set_description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "update" + domainName,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": "update" + domainName,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name_servers": "100.100.2.138",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name_servers": "100.100.2.138",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_lease_time": "24h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_lease_time": "24h",
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
				Config: testAccConfig(map[string]interface{}{
					"dhcp_options_set_name":        "${var.name}",
					"dhcp_options_set_description": "${var.name}",
					"domain_name":                  domainName,
					"domain_name_servers":          "100.100.2.136",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_name":        name,
						"dhcp_options_set_description": name,
						"domain_name":                  domainName,
						"domain_name_servers":          "100.100.2.136",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAliCloudVPCDhcpOptionsSet_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_dhcp_options_set.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCDhcpOptionsSetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcDhcpOptionsSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcdhcpoptionsset%d", defaultRegionToTest, rand)
	domainName := fmt.Sprintf("tftestacc%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCDhcpOptionsSetBasicDependence0)
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
					"dhcp_options_set_name":        "${var.name}",
					"dhcp_options_set_description": "${var.name}",
					"domain_name":                  domainName,
					"domain_name_servers":          "100.100.2.136",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_name":        name,
						"dhcp_options_set_description": name,
						"domain_name":                  domainName,
						"domain_name_servers":          "100.100.2.136",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudVPCDhcpOptionsSet_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_dhcp_options_set.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCDhcpOptionsSetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcDhcpOptionsSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcdhcpoptionsset%d", defaultRegionToTest, rand)
	domainName := fmt.Sprintf("tftestacc%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCDhcpOptionsSetBasicDependence0)
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
					"dhcp_options_set_name":        "${var.name}",
					"dhcp_options_set_description": "${var.name}",
					"domain_name":                  domainName,
					"domain_name_servers":          "100.100.2.136",
					"dry_run":                      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_name":        name,
						"dhcp_options_set_description": name,
						"domain_name":                  domainName,
						"domain_name_servers":          "100.100.2.136",
						"dry_run":                      "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAliCloudVPCDhcpOptionsSet_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_dhcp_options_set.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCDhcpOptionsSetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcDhcpOptionsSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcdhcpoptionsset%d", defaultRegionToTest, rand)
	domainName := fmt.Sprintf("tftestacc%d.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCDhcpOptionsSetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcDhcpOptionsSetSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dhcp_options_set_name":        "${var.name}",
					"dhcp_options_set_description": "${var.name}",
					"domain_name":                  domainName,
					"domain_name_servers":          "100.100.2.136",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_name":        name,
						"dhcp_options_set_description": name,
						"domain_name":                  domainName,
						"domain_name_servers":          "100.100.2.136",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dhcp_options_set_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dhcp_options_set_description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "update" + domainName,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": "update" + domainName,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name_servers": "100.100.2.138",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name_servers": "100.100.2.138",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_lease_time": "24h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_lease_time": "24h",
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_lease_time": "48h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_lease_time": "48h",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lease_time": "24h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lease_time": "24h",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lease_time": "48h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lease_time": "48h",
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
				Config: testAccConfig(map[string]interface{}{
					"dhcp_options_set_name":        "${var.name}",
					"dhcp_options_set_description": "${var.name}",
					"domain_name":                  domainName,
					"domain_name_servers":          "100.100.2.136",
					"ipv6_lease_time":              "24h",
					"lease_time":                   "24h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dhcp_options_set_name":        name,
						"dhcp_options_set_description": name,
						"domain_name":                  domainName,
						"domain_name_servers":          "100.100.2.136",
						"ipv6_lease_time":              "24h",
						"lease_time":                   "24h",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudVPCDhcpOptionsSetMap0 = map[string]string{
	"dry_run": NOSET,
	"status":  CHECKSET,
}

func AlicloudVPCDhcpOptionsSetBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}
