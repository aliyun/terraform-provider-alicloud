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
		"alicloud_dns_group",
		&resource.Sweeper{
			Name: "alicloud_dns_group",
			F:    testSweepDnsGroup,
		})
}

func testSweepDnsGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := alidns.CreateDescribeDomainGroupsRequest()

	var allGroups []alidns.DomainGroup
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", request.GetActionName(), err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*alidns.DescribeDomainGroupsResponse)
		groups := response.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if strings.HasPrefix(domainGroup.GroupName, "tf-testacc") {
				allGroups = append(allGroups, domainGroup)
			} else {
				log.Printf("Skip %#v.", domainGroup)
			}
		}
		if len(groups) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	removeRequest := alidns.CreateDeleteDomainGroupRequest()

	for _, group := range allGroups {
		removeRequest.GroupId = group.GroupId
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DeleteDomainGroup(removeRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", request.GetActionName(), err)
		}
		addDebug(request.GetActionName(), raw)
	}
	return nil
}

func TestAccAlicloudDnsGroup_basic(t *testing.T) {
	resourceId := "alicloud_dns_group.group"
	var v alidns.DomainGroup
	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(100000, 999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns_group.group",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsGroupConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testacc%d", rand),
					}),
				),
			},
			{
				Config: testAccDnsGroupConfig(rand - 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testacc%d", rand-1),
					}),
				),
			},
		},
	})
}

func testAccDnsGroupConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_dns_group" "group" {
	  name = "tf-testacc%d"
	}
	`, rand)
}

func testAccCheckDnsGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dns_group" {
			continue
		}

		// Try to find the domain group
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		dnsService := &DnsService{client: client}
		_, err := dnsService.DescribeDnsGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}

		return WrapError(Error("Error groups still exist"))
	}
	return nil
}
