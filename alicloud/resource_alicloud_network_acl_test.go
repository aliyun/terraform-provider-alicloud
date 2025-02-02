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
	resource.AddTestSweepers("alicloud_network_acl", &resource.Sweeper{
		Name: "alicloud_network_acl",
		F:    testSweepNetworkAcl,
	})
}

func testSweepNetworkAcl(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeNetworkAcls"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	var response map[string]interface{}
	networkAclIds := make([]string, 0)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			log.Printf("Error retrieving network acl: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.NetworkAcls.NetworkAcl", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NetworkAcls.NetworkAcl", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["NetworkAclName"])
			id := fmt.Sprint(item["NetworkAclId"])
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Network Acl: %s (%s)", name, id)
					continue
				}
			}
			networkAclIds = append(networkAclIds, id)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	vpcService := VpcService{client}
	for _, id := range networkAclIds {
		//	Delete attach resources
		object, err := vpcService.DescribeNetworkAcl(id)
		if err != nil {
			log.Println("DescribeNetworkAcl failed", err)
		}
		deleteResources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		if len(deleteResources) > 0 {
			request := map[string]interface{}{
				"NetworkAclId": id,
			}
			resourcesMaps := make([]map[string]interface{}, 0)
			for _, resources := range deleteResources {
				resourcesArg := resources.(map[string]interface{})
				resourcesMap := map[string]interface{}{
					"ResourceId":   resourcesArg["ResourceId"],
					"ResourceType": resourcesArg["ResourceType"],
				}
				resourcesMaps = append(resourcesMaps, resourcesMap)
			}
			request["Resource"] = resourcesMaps
			request["RegionId"] = client.RegionId
			action := "UnassociateNetworkAcl"
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
				log.Println("UnassociateNetworkAcl failed", err)
			}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, 5*time.Minute, 5*time.Second, vpcService.NetworkAclStateRefreshFunc(id, []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				log.Println("UnassociateNetworkAcl failed", err)
			}
		}

		log.Printf("[INFO] Deleting Network Acl: (%s)", id)
		request := map[string]interface{}{
			"NetworkAclId": id,
		}
		action := "DeleteNetworkAcl"
		request["RegionId"] = client.RegionId
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
		if err != nil {
			log.Printf("[ERROR] Failed to delete Network Acl (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccAliCloudVpcNetworkAcl_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_network_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudNetworkAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snetworkacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNetworkAclBasicDependence0)
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
					"vpc_id":           "${alicloud_vpc.default.id}",
					"network_acl_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":           CHECKSET,
						"network_acl_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"egress_acl_entries": []map[string]interface{}{
						{
							"description":            "engress test",
							"destination_cidr_ip":    "10.0.0.0/24",
							"network_acl_entry_name": "tf-testacc78924",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"egress_acl_entries.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ingress_acl_entries": []map[string]interface{}{
						{
							"description":            "ingress test",
							"network_acl_entry_name": "tf-testacc78999",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
							"source_cidr_ip":         "10.0.0.0/24",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ingress_acl_entries.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_acl_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_acl_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${alicloud_vswitch.default0.id}",
							"resource_type": "VSwitch",
						},
						{
							"resource_id":   "${alicloud_vswitch.default1.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${alicloud_vswitch.default0.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      name,
					"network_acl_name": name,
					"ingress_acl_entries": []map[string]interface{}{
						{
							"description":            "ingress test change",
							"network_acl_entry_name": "tf-testacc78999",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
							"source_cidr_ip":         "10.0.0.0/24",
						},
					},
					"egress_acl_entries": []map[string]interface{}{
						{
							"description":            "engress test change",
							"destination_cidr_ip":    "10.0.0.0/24",
							"network_acl_entry_name": "tf-testacc78924",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           name,
						"network_acl_name":      name,
						"ingress_acl_entries.#": "1",
						"egress_acl_entries.#":  "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudVpcNetworkAcl_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_network_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudNetworkAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snetworkacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNetworkAclBasicDependence0)
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
					"vpc_id":      "${alicloud_vpc.default.id}",
					"name":        name,
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":      CHECKSET,
						"name":        name,
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "1",
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

func TestAccAliCloudVpcNetworkAcl_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_network_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudNetworkAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snetworkacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNetworkAclBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-6"})
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":           "${alicloud_vpc.default.id}",
					"network_acl_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":           CHECKSET,
						"network_acl_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"egress_acl_entries": []map[string]interface{}{
						{
							"description":            "engress test",
							"destination_cidr_ip":    "2408:4004:cc:400::/56",
							"network_acl_entry_name": "tf-testacc78924",
							"policy":                 "accept",
							"protocol":               "icmpv6",
							"ip_version":             "IPV6",
							"entry_type":             "custom",
							"port":                   "-1/-1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"egress_acl_entries.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ingress_acl_entries": []map[string]interface{}{
						{
							"description":            "ingress test",
							"network_acl_entry_name": "tf-testacc78999",
							"policy":                 "accept",
							"protocol":               "icmpv6",
							"source_cidr_ip":         "2408:4004:cc:400::/56",
							"ip_version":             "IPV6",
							"entry_type":             "custom",
							"port":                   "-1/-1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ingress_acl_entries.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_acl_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_acl_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${alicloud_vswitch.default0.id}",
							"resource_type": "VSwitch",
						},
						{
							"resource_id":   "${alicloud_vswitch.default1.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${alicloud_vswitch.default0.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudVpcNetworkAcl_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_network_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudNetworkAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snetworkacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNetworkAclBasicDependence0)
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
					"vpc_id":      "${alicloud_vpc.default.id}",
					"name":        name,
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":      CHECKSET,
						"name":        name,
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_network_acl_id": "${alicloud_network_acl.copyed.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_network_acl_id"},
			},
		},
	})
}

var AlicloudNetworkAclMap0 = map[string]string{}

func AlicloudNetworkAclBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%[1]s"
		}
variable "name_change" {
			default = "%[1]s_change"
		}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name = var.name
}
resource "alicloud_vswitch" "default0" {
  vpc_id            = alicloud_vpc.default.id
  vswitch_name      = var.name
  cidr_block        = cidrsubnets(alicloud_vpc.default.cidr_block, 4, 4)[0]
  zone_id           = data.alicloud_zones.default.ids.0
}
resource "alicloud_vswitch" "default1" {
  vpc_id            = alicloud_vpc.default.id
  vswitch_name      = var.name_change
  cidr_block        = cidrsubnets(alicloud_vpc.default.cidr_block, 4, 4)[1]
  zone_id           = data.alicloud_zones.default.ids.0
}

resource "alicloud_network_acl" "copyed" {
  vpc_id           = alicloud_vpc.default.id
  network_acl_name = var.name
  description      = var.name
  ingress_acl_entries {
    description            = "${var.name}-ingress"
    network_acl_entry_name = "${var.name}-ingress"
    source_cidr_ip         = "10.0.0.0/24"
    policy                 = "accept"
    port                   = "20/80"
    protocol               = "tcp"
  }
  egress_acl_entries {
    description            = "${var.name}-egress"
    network_acl_entry_name = "${var.name}-egress"
    destination_cidr_ip    = "10.0.0.0/24"
    policy                 = "accept"
    port                   = "20/80"
    protocol               = "tcp"
  }
}
`, name)
}

func AlicloudNetworkAclBasicDependence2(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%[1]s"
		}
variable "name_change" {
			default = "%[1]s_change"
		}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name = var.name
  enable_ipv6 = "true"
}
resource "alicloud_vswitch" "default0" {
  vpc_id            = alicloud_vpc.default.id
  vswitch_name      = var.name
  cidr_block        = cidrsubnets(alicloud_vpc.default.cidr_block, 4, 4)[0]
  zone_id           = data.alicloud_zones.default.ids.0
}
resource "alicloud_vswitch" "default1" {
  vpc_id            = alicloud_vpc.default.id
  vswitch_name      = var.name_change
  cidr_block        = cidrsubnets(alicloud_vpc.default.cidr_block, 4, 4)[1]
  zone_id           = data.alicloud_zones.default.ids.0
}

`, name)
}

// Test Vpc NetworkAcl. >>> Resource test cases, automatically generated.
// Case 2583
func TestAccAliCloudVpcNetworkAcl_basic2583(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_network_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcNetworkAclMap2583)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcnetworkacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcNetworkAclBasicDependence2583)
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
					"vpc_id":           "${alicloud_vpc.defaultVpc.id}",
					"network_acl_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":           CHECKSET,
						"network_acl_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_acl_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_acl_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testacc-acl-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testacc-acl-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_acl_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_acl_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "test",
					"vpc_id":           "${alicloud_vpc.defaultVpc.id}",
					"network_acl_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "test",
						"vpc_id":           CHECKSET,
						"network_acl_name": name + "_update",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudVpcNetworkAclMap2583 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcNetworkAclBasicDependence2583(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultVpc" {
  description = "tf-testacc-acl-vpc"
}


`, name)
}

// Case 2583  twin
func TestAccAliCloudVpcNetworkAcl_basic2583_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_network_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcNetworkAclMap2583)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcnetworkacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcNetworkAclBasicDependence2583)
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
					"description":      "tf-testacc-acl-update",
					"vpc_id":           "${alicloud_vpc.defaultVpc.id}",
					"network_acl_name": name,
					"ingress_acl_entries": []map[string]interface{}{
						{
							"description":            "ingress test change",
							"network_acl_entry_name": "tf-testacc78999",
							"policy":                 "accept",
							"ip_version":             "IPV4",
							"entry_type":             "custom",
							"port":                   "20/80",
							"protocol":               "tcp",
							"source_cidr_ip":         "10.0.0.0/24",
						},
					},
					"egress_acl_entries": []map[string]interface{}{
						{
							"description":            "engress test change",
							"destination_cidr_ip":    "10.0.0.0/24",
							"network_acl_entry_name": "tf-testacc78924",
							"ip_version":             "IPV4",
							"entry_type":             "custom",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           "tf-testacc-acl-update",
						"vpc_id":                CHECKSET,
						"network_acl_name":      name,
						"ingress_acl_entries.#": "1",
						"egress_acl_entries.#":  "1",
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "Test",
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

// Test Vpc NetworkAcl. <<< Resource test cases, automatically generated.
