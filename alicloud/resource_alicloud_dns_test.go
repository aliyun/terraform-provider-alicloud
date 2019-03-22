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
	var v *alidns.DescribeDomainInfoResponse

	randInt := acctest.RandInt()
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
					testAccCheckDnsExists("alicloud_dns.dns", v),
					resource.TestCheckResourceAttr("alicloud_dns.dns", "name", fmt.Sprintf("tf-testaccdnsbasic%v.abc", randInt)),
					resource.TestCheckResourceAttrSet("alicloud_dns.dns", "dns_server.#"),
				),
			},
			{
				Config: testAccDnsConfig_group_id(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsExists("alicloud_dns.dns", v),
					resource.TestCheckResourceAttr("alicloud_dns.dns", "name", fmt.Sprintf("tf-testaccdnsbasic%v.abc", randInt)),
					resource.TestCheckResourceAttrSet("alicloud_dns.dns", "group_id"),
					resource.TestCheckResourceAttrSet("alicloud_dns.dns", "dns_server.#"),
				),
			},
		},
	})

}

func testAccCheckDnsExists(n string, domain *alidns.DescribeDomainInfoResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Domain ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		dnsService := &DnsService{client: client}
		domainInfo, err := dnsService.DescribeDns(rs.Primary.Attributes["name"])
		log.Printf("[WARN] Domain id %#v", rs.Primary.ID)

		if err == nil {
			domain = domainInfo
			return nil
		}
		return WrapError(err)
	}
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
