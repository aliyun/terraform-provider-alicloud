package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	waf_openapi "github.com/aliyun/alibaba-cloud-sdk-go/services/waf-openapi"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_waf_domain", &resource.Sweeper{
		Name: "alicloud_waf_domain",
		F:    testSweepWafDomains,
	})
}
var instance_ids []string
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
	var v waf_openapi.DescribeDomainResponse

	resourceId := "alicloud_waf_domain.domain"
	ra := resourceAttrInit(resourceId, wafDomainBasicMap)

	serviceFunc := func() interface{} {
		return &Waf_openapiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%s%d.waf.com", defaultRegionToTest, rand)
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
					"domain":            "waf.test123.abc",
					"instance_id":       "waf_elasticity-cn-0pp1ko78f00d", // waf_openapi.DescribeInstanceInfoRequest{}.InstanceId,
					"is_access_product": "0",
					"source_ips":        "1.1.1.1",
					"cluster_type":      "0",
					"http2_port":        "433",
					"http_port":         "80",
					"https_port":        "433",
					"http_to_user_ip":   "0",
					"https_redirect":    "0",
					"load_balancing":    "0",
					"resource_group_id": "rg-atstuj3rtoptyui",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":            "www.example.com",
						"instance_id":       REMOVEKEY,
						"is_access_product": "0",
						"source_ips":        "1.1.1.1",
						"cluster_type":      "0",
						"http2_port":        "433",
						"http_port":         "80",
						"https_port":        "433",
						"http_to_user_ip":   "0",
						"https_redirect":    "0",
						"load_balancing":    "0",
						"resource_group_id": "rg-atstuj3rtoptyui",
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
					"source_ips": "2.2.2.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_ips": "2.2.2.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http2_port": "433",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http2_port": "433",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"http_port": "80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_port": "80",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"https_port": "433",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"https_port": "433",
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
					"is_access_product": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_access_product": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http_to_user_ip": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_to_user_ip": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"https_redirect": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"https_redirect": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain":            "www.example.com",
					"instance_id":       "waf_elasticity-cn-0pp1ko78f00d",
					"is_access_product": "0",
					"source_ips":        "1.1.1.1",
					"cluster_type":      "0",
					"http2_port":        "433",
					"http_port":         "80",
					"https_port":        "433",
					"http_to_user_ip":   "0",
					"https_redirect":    "0",
					"load_balancing":    "0",
					"resource_group_id": "rg-atstuj3rtoptyui",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":            "www.example.com",
						"instance_id":       REMOVEKEY,
						"is_access_product": "0",
						"source_ips":        "1.1.1.1",
						"cluster_type":      "0",
						"http2_port":        "433",
						"http_port":         "80",
						"https_port":        "433",
						"http_to_user_ip":   "0",
						"https_redirect":    "0",
						"load_balancing":    "0",
						"resource_group_id": "rg-atstuj3rtoptyui",
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
