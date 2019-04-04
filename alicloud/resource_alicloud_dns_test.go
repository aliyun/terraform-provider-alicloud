package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers(
		"alicloud.dns",
		&resource.Sweeper{
			Name: "alicloud.dns",
			F:    testSweepDns,
		})
}

func testSweepDns(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	queryRequest := alidns.CreateDescribeDomainsRequest()
	var allDomains []alidns.Domain
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
			if strings.HasPrefix(domain.DomainName, "tf-testaccdns") {
				allDomains = append(allDomains, domain)
			} else {
				log.Printf("Skip %#v", domain)
			}
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
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DeleteDomain(removeRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error %s", removeRequest.GetActionName(), err)
		}
		addDebug(removeRequest.GetActionName(), raw)

	}

	return nil

}

func TestAccAlicloudDns_basic(t *testing.T) {
	randInt := acctest.RandInt()
	var v *alidns.DescribeDomainInfoResponse
	ra := resourceAttrInit("alicloud_dns.dns", map[string]string{})
	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit("alicloud_dns.dns", &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns.dns",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsConfig_create(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         fmt.Sprintf("tf-testaccdnsbasic%v.abc", randInt),
						"dns_server.#": CHECKSET,
					}),
				),
			},
			{
				Config: testAccDnsConfig_group_id(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": CHECKSET,
					}),
				),
			},
		},
	})

}

func testAccCheckDnsDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dns" {
			continue
		}

		// Try to find the domain
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		dnsService := &DnsService{client: client}
		_, err := dnsService.DescribeDns(rs.Primary.Attributes["name"])

		if err != nil && !IsExceptedErrors(err, []string{InvalidDomainNameNoExist}) {
			return WrapError(err)
		}
	}

	return nil
}

func testAccDnsConfig_create(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsbasic%v.abc"
}
`, randInt)
}

func testAccDnsConfig_group_id(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "group" {
  name = "test-dns-group"
}

resource "alicloud_dns" "dns" {
  name = "tf-testaccdnsbasic%v.abc"
  group_id = "${alicloud_dns_group.group.id}"
}
`, randInt)
}
