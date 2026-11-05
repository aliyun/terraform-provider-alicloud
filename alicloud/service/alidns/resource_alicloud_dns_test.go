package alidns_test

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	queryRequest := alidnssdk.CreateDescribeDomainsRequest()
	var allDomains []alidnssdk.DomainInDescribeDomains
	queryRequest.PageSize = requests.NewInteger(helper.PageSizeLarge)
	queryRequest.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.DescribeDomains(queryRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error %#v", queryRequest.GetActionName(), err)
		}
		logs.AddDebug(queryRequest.GetActionName(), raw)
		response, _ := raw.(*alidnssdk.DescribeDomainsResponse)
		domains := response.Domains.Domain
		for _, domain := range domains {
			if strings.HasPrefix(domain.DomainName, "tf-testacc"+acctest.DefaultRegionToTest) {
				allDomains = append(allDomains, domain)
			} else {
				log.Printf("Skip %#v", domain)
			}
		}

		if len(domains) < helper.PageSizeLarge {
			break
		}
		if page, err := helper.GetNextpageNumber(queryRequest.PageNumber); err != nil {
			return errs.WrapError(err)
		} else {
			queryRequest.PageNumber = page
		}
	}

	removeRequest := alidnssdk.CreateDeleteDomainRequest()
	removeRequest.DomainName = ""
	for _, domain := range allDomains {
		removeRequest.DomainName = domain.DomainName
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.DeleteDomain(removeRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error %s", removeRequest.GetActionName(), err)
		}
		logs.AddDebug(removeRequest.GetActionName(), raw)

	}

	return nil

}

func TestAccAlicloudAlidns_basic(t *testing.T) {
	resourceId := "alicloud_dns.default"
	randInt := tfacctest.RandIntRange(10000, 99999)
	var v *alidnssdk.DescribeDomainInfoResponse
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return alidns.NewDnsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, strconv.FormatInt(int64(randInt), 10), resourceDnsConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              "${var.dnsName}",
					"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         fmt.Sprintf("tf-testacc%sdnsbasic%d.abc", acctest.DefaultRegionToTest, randInt),
						"dns_server.#": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id": "${alicloud_dns_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceDnsConfigDependence(name string) string {
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
`, acctest.DefaultRegionToTest, name, name)
}
