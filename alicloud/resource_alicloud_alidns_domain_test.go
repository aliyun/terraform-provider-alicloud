package alicloud

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_alidns_domain",
		&resource.Sweeper{
			Name: "alicloud_alidns_domain",
			F:    testSweepAlidnsDomain,
		})
}

func testSweepAlidnsDomain(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	queryRequest := alidns.CreateDescribeDomainsRequest()
	var allDomains []alidns.DomainInDescribeDomains
	queryRequest.PageSize = requests.NewInteger(PageSizeLarge)
	queryRequest.PageNumber = requests.NewInteger(1)

	for {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DescribeDomains(queryRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error %#v", queryRequest.GetActionName(), err)
		}
		addDebug(queryRequest.GetActionName(), raw)
		response, _ := raw.(*alidns.DescribeDomainsResponse)
		domains := response.Domains.Domain
		for _, domain := range domains {
			// RED LINE: alidns may hold real存 / ICP-backed production domains.
			// Always prefix-guard (even when SWEEP_ALL=1) so production domains are
			// never deleted. Skip any domain whose name carries a do-not-delete marker.
			if protectedSweepVpcName(domain.DomainName) {
				log.Printf("[INFO] Skipping DNS domain with do-not-delete marker: %s", domain.DomainName)
				continue
			}
			if !isTestDnsName(domain.DomainName) {
				log.Printf("[INFO] Skipping non-test DNS domain (real存 preserve): %s", domain.DomainName)
				continue
			}
			allDomains = append(allDomains, domain)
		}

		if len(domains) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(queryRequest.PageNumber); err != nil {
			return WrapError(err)
		} else {
			queryRequest.PageNumber = page
		}
	}
	removeRequest := alidns.CreateDeleteDomainRequest()
	removeRequest.DomainName = ""
	for _, domain := range allDomains {
		removeRequest.DomainName = domain.DomainName
		raw, err := client.WithDnsClient(func(dnsClietn *alidns.Client) (interface{}, error) {
			return dnsClietn.DeleteDomain(removeRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error %s", removeRequest.GetActionName(), err)
		}
		addDebug(removeRequest.GetActionName(), raw)
	}
	return nil
}

// isTestDnsName reports whether a DNS domain/record name looks like it was
// created by the provider's acceptance tests. alidns may hold real存 / ICP-backed
// production domains and records, so sweepers MUST only match test-prefixed
// names EVEN when SWEEP_ALL=1 (sweepAll() bypasses the framework prefix filter).
func isTestDnsName(name string) bool {
	lower := strings.ToLower(name)
	// Standard acceptance-test prefixes (see testSweepDBInstances).
	for _, prefix := range []string{"tf-test", "tf_test", "test-"} {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	// Marker emitted by the SSL certificates service acceptance tests, which
	// create validation records such as "_pki-validation.tfaccsslcasXXX" under
	// real存-style domains (e.g. aliterraform.com). The domain itself is preserved,
	// but these stray test records must be swept.
	if strings.Contains(lower, "tfaccssl") || strings.Contains(lower, "tfacc") {
		return true
	}
	return false
}

func TestAccAliCloudAlidnsDomain_basic(t *testing.T) {
	resourceId := "alicloud_alidns_domain.default"
	randInt := acctest.RandIntRange(10000, 99999)
	var v alidns.DescribeDomainInfoResponse
	ra := resourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return &AlidnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, strconv.FormatInt(int64(randInt), 10), resourceDnsDomainConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":       "${var.dnsName}",
					"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
					"remark":            "test new domain",
					"tags": map[string]string{
						"Created": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":  fmt.Sprintf("tf-testacc%sdnsbasic%d.abc", defaultRegionToTest, randInt),
						"remark":       "test new domain",
						"tags.%":       "1",
						"tags.Created": "TF",
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
					"group_id": "${alicloud_dns_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TFM",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "1",
						"tags.Created": "TFM",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "test new domain again",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "test new domain again",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${var.dnsName}",
					"remark":      "test new domain",
					"tags": map[string]string{
						"Created": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":  fmt.Sprintf("tf-testacc%sdnsbasic%d.abc", defaultRegionToTest, randInt),
						"remark":       "test new domain",
						"tags.%":       "1",
						"tags.Created": "TF",
					}),
				),
			},
		},
	})
}

func resourceDnsDomainConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "dnsName"{
	default = "tf-testacc%sdnsbasic%s.abc"
}

variable "dnsGroupName"{
	default = "tf-testaccdns%s"
}

resource "alicloud_dns_group" "default" {
  name = "${var.dnsGroupName}"
}
`, defaultRegionToTest, name, name)
}
