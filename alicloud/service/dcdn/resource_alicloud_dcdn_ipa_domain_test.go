package dcdn_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dcdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_dcdn_ipa_domain",
		&resource.Sweeper{
			Name: "alicloud_dcdn_ipa_domain",
			F:    testSweepDcdnIpaDomain,
		})
}
func testSweepDcdnIpaDomain(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DCDNSupportRegions) {
		log.Printf("[INFO] Skipping Dcdn Ipa Domain unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDcdnIpaUserDomains"
	request := map[string]interface{}{}

	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Domains.PageData", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Domains.PageData", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DomainName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Dcdn Ipa Domain: %s", item["DomainName"].(string))
				continue
			}
			action := "DeleteDcdnIpaDomain"
			request := map[string]interface{}{
				"DomainName": item["DomainName"],
			}
			_, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Dcdn Ipa Domain (%s): %s", item["DomainName"].(string), err)
			}
			log.Printf("[INFO] Delete Dcdn Ipa Domain success: %s ", item["DomainName"].(string))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}
func TestAccAliCloudDcdnIpaDomain_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_ipa_domain.default"
	checkoutSupportedRegions(t, true, connectivity.DCDNSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDCDNIpaDomainMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnIpaDomain")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDCDNIpaDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${var.domain_name}",
					"scope":       "overseas",
					"sources": []map[string]interface{}{
						{
							"content":  "www.alicloud-provider.cn",
							"port":     "8898",
							"priority": "20",
							"type":     "domain",
							"weight":   "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
						"domain_name":       name,
						"scope":             "overseas",
						"sources.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "offline",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "offline",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "online",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "online",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []map[string]interface{}{
						{
							"content":  "3.3.3.3",
							"port":     "7261",
							"priority": "30",
							"type":     "ipaddr",
							"weight":   "7",
						},
						{
							"content":  "5.3.3.3",
							"port":     "7221",
							"priority": "20",
							"type":     "ipaddr",
							"weight":   "9",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sources.#": "2",
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
func TestAccAliCloudDcdnIpaDomain_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_ipa_domain.default"
	checkoutSupportedRegions(t, true, connectivity.DCDNSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDCDNIpaDomainMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnIpaDomain")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDCDNIpaDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${var.domain_name}",
					"sources": []map[string]interface{}{
						{
							"content":  "www.alicloud-provider.cn",
							"port":     "8898",
							"priority": "20",
							"type":     "domain",
							"weight":   "10",
						},
					},
					"scope":             "overseas",
					"status":            "online",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
						"domain_name":       name,
						"sources.#":         "1",
						"status":            "online",
						"scope":             "overseas",
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

var AlicloudDCDNIpaDomainMap0 = map[string]string{
	"domain_name":       acctest.CHECKSET,
	"resource_group_id": acctest.CHECKSET,
	"status":            acctest.CHECKSET,
	"sources.#":         acctest.CHECKSET,
}

func AlicloudDCDNIpaDomainBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "domain_name" {	
	default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
	name_regex = "default"
}
`, name)
}
