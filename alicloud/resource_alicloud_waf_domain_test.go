package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	//"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"os"

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

	var domains []waf_openapi.Domain
	args := waf_openapi.CreateDescribeDomainRequest()
	args.Port = requests.DefaultHttpPort

	for {

		raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {

			return waf_openapiClient.DescribeDomain(args)
		})
		if err != nil {
			log.Printf("Error retrieving WAF Domain: %s", err)
		}
		addDebug(args.GetActionName(), raw)
		resp, _ := raw.(*waf_openapi.DescribeDomainResponse)
		if resp == nil || len(resp.Domain.Cname) < 1 {
			break
		}
		domains = append(domains, resp.Domain)

	}

	for _, v := range domains {
		name := v.Cname
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping WAF domain: %s", name)
			continue
		}
		log.Printf("[INFO] Deleting WAF domain: %s", name)
		request := waf_openapi.CreateDeleteDomainRequest()
		request.Domain = name
		raw, err := client.WithWafOpenapiClient(func(waf_openapiClient *waf_openapi.Client) (interface{}, error) {
			return waf_openapiClient.DeleteDomain(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete WAF domain (%s): %s", name, err)
		}
		addDebug(request.GetActionName(), raw)
	}
	return nil
}

func TestAccAlicloudWafDomain(t *testing.T) {
	var v waf_openapi.Domain

	resourceId := "alicloud_waf_domain.domain"
	ra := resourceAttrInit(resourceId, wafDomainBasicMap)

	serviceFunc := func() interface{} {
		return &Waf_openapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.wafqa3.com", defaultRegionToTest, rand)
	//name := "tf-testacctest.wafqa3.com"
	instanceId := os.Getenv("ALICLOUD_WAF_INSTANCE_ID")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceWafDomainDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithWafInstanceSetting(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain":            name,
					"instance_id":       instanceId,
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
					"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":            name,
						"instance_id":       instanceId,
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
						"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				//ImportStateVerifyIgnore: []string{"resource_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_type": "VirtualCluster",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_type": "VirtualCluster",
					}),
				),
			},
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
					"instance_id":       instanceId,
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
					"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":            name,
						"instance_id":       instanceId,
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
						"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
					}),
				),
			},
		},
	})
}

var wafDomainBasicMap = map[string]string{}

func resourceWafDomainDependence(name string) string {
	return ``
}

var wg sync.WaitGroup
