package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"sync"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func init() {
	resource.AddTestSweepers("alicloud_waf_domain", &resource.Sweeper{
		Name: "alicloud_waf_domain",
		F:    testSweepWafDomains,
	})
}

func testSweepWafDomains(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		fmt.Sprintf("tf-testacc%s", region),
		fmt.Sprintf("tf_testacc%s", region),
	}

	wafInstanceIds := make([]string, 0)
	domainIds := make([]string, 0)
	request := make(map[string]interface{})
	conn, err := client.NewWafClient()
	if err != nil {
		return WrapError(err)
	}
	action := "DescribeInstanceInfos"

	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve waf instance in service list: %s", err)
	}
	resp, err := jsonpath.Get("$.InstanceInfos", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceInfos", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if item["InstanceId"] == nil {
			continue
		}
		wafInstanceIds = append(wafInstanceIds, item["InstanceId"].(string))
	}

	for _, instanceId := range wafInstanceIds {
		action = "DescribeDomainNames"
		request = make(map[string]interface{})
		request["InstanceId"] = instanceId
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve waf domain in service list: %s", err)
		}
		if response["DomainNames"] != nil {
			for _, item := range response["DomainNames"].([]interface{}) {
				domainIds = append(domainIds, fmt.Sprintf(`%s:%s`, instanceId, item.(string)))
			}
		}
	}

	for _, id := range domainIds {
		part := strings.Split(id, ":")
		name := part[1]
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping WAF domain: %s", id)
			continue
		}
		log.Printf("[INFO] Deleting WAF domain: %s", id)

		request = make(map[string]interface{})
		action = "DeleteDomain"
		request["InstanceId"] = part[0]
		request["Domain"] = name
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete WAF domain (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccAlicloudWAFDomain(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_waf_domain.domain"
	ra := resourceAttrInit(resourceId, wafDomainBasicMap)

	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Waf_openapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.wafqa3.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceWafDomainDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain":            name,
					"instance_id":       "${data.alicloud_waf_instances.default.ids.0}",
					"is_access_product": "Off",
					"source_ips":        []string{"1.1.1.1"},
					"cluster_type":      "PhysicalCluster",
					"http2_port":        []string{"443"},
					"http_port":         []string{"80"},
					"https_port":        []string{"443"},
					"http_to_user_ip":   "Off",
					"https_redirect":    "Off",
					"load_balancing":    "IpHash",
					"log_headers": []map[string]interface{}{
						{
							"key":   "kkk",
							"value": "vvv",
						},
						{
							"key":   "test",
							"value": "ddd",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":            name,
						"instance_id":       CHECKSET,
						"is_access_product": "Off",
						"source_ips.#":      "1",
						"cluster_type":      "PhysicalCluster",
						"http2_port.#":      "1",
						"http_port.#":       "1",
						"https_port.#":      "1",
						"http_to_user_ip":   "Off",
						"https_redirect":    "Off",
						"load_balancing":    "IpHash",
						"log_headers.#":     "2",
						"resource_group_id": CHECKSET,
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Only the exclusive version of WAF instance supports modification
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"cluster_type": "VirtualCluster",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"cluster_type": "VirtualCluster",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_time": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_time": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http2_port": []string{"443", "8443"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http2_port.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http_to_user_ip": "On",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_to_user_ip": "On",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http_port": []string{"8080", "80"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_port.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_port": []string{"8443", "443"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"https_port.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_redirect": "On",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"https_redirect": "On",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_ips": []string{"1.1.1.1", "2.2.2.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_ips.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancing": "RoundRobin",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancing": "RoundRobin",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_headers": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_headers.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_time": "120",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_time": "120",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_access_product": "On",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_access_product": "On",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"write_time": "150",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"write_time": "150",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain":            name,
					"instance_id":       "${data.alicloud_waf_instances.default.ids.0}",
					"cluster_type":      "PhysicalCluster",
					"connection_time":   "30",
					"http2_port":        []string{"8443"},
					"http_port":         []string{"80"},
					"http_to_user_ip":   "Off",
					"https_port":        []string{"443"},
					"https_redirect":    "Off",
					"is_access_product": "Off",
					"load_balancing":    "IpHash",
					"log_headers": []map[string]interface{}{
						{
							"key":   "kkk1",
							"value": "vvv1",
						},
						{
							"key":   "Test2",
							"value": "Sdd2",
						},
					},
					"read_time":         "140",
					"source_ips":        []string{"1.1.1.1"},
					"write_time":        "800",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":            name,
						"instance_id":       CHECKSET,
						"cluster_type":      "PhysicalCluster",
						"connection_time":   "30",
						"http2_port.#":      "1",
						"http_port.#":       "1",
						"http_to_user_ip":   "Off",
						"https_port.#":      "1",
						"https_redirect":    "Off",
						"is_access_product": "Off",
						"load_balancing":    "IpHash",
						"log_headers.#":     "2",
						"read_time":         "140",
						"source_ips.#":      "1",
						"write_time":        "800",
						"resource_group_id": CHECKSET,
					}),
				),
			},
		},
	})
}

var wafDomainBasicMap = map[string]string{}

func resourceWafDomainDependence(name string) string {
	return `
data "alicloud_resource_manager_resource_groups" "default" {
 name_regex = "^default$"
 }
data "alicloud_waf_instances" "default" {}
`
}

var wg sync.WaitGroup
