package alidns_test

import (
	"fmt"
	"log"
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
		"alicloud_alidns_domain_group",
		&resource.Sweeper{
			Name: "alicloud_alidns_domain_group",
			F:    testSweepAlidnsDomainGroup,
		})
}

func testSweepAlidnsDomainGroup(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := alidnssdk.CreateDescribeDomainGroupsRequest()

	var allGroups []alidnssdk.DomainGroup
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", request.GetActionName(), err)
		}
		logs.AddDebug(request.GetActionName(), raw)
		response, _ := raw.(*alidnssdk.DescribeDomainGroupsResponse)
		groups := response.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if strings.HasPrefix(strings.ToLower(domainGroup.GroupName), "tf-testacc") {
				allGroups = append(allGroups, domainGroup)
			} else {
				log.Printf("Skip %#v.", domainGroup)
			}
		}
		if len(groups) < helper.PageSizeLarge {
			break
		}
		if page, err := helper.GetNextpageNumber(request.PageNumber); err != nil {
			return errs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	removeRequest := alidnssdk.CreateDeleteDomainGroupRequest()

	for _, group := range allGroups {
		removeRequest.GroupId = group.GroupId
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.DeleteDomainGroup(removeRequest)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", request.GetActionName(), err)
		}
		logs.AddDebug(request.GetActionName(), raw)
	}
	return nil
}

func TestAccAliCloudAlidnsDomainGroup_basic(t *testing.T) {
	var v alidnssdk.DomainGroup
	resourceId := "alicloud_alidns_domain_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlidnsDomainGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alidns.NewAlidnsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlidnsDomainGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDG%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlidnsDomainGroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_group_name": name,
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
					"domain_group_name": fmt.Sprintf("tf-testaccdns%d", rand-1),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_group_name": fmt.Sprintf("tf-testaccdns%d", rand-1),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_group_name": fmt.Sprintf("tf-testaccdns%d", rand+1),
					"lang":              "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_group_name": fmt.Sprintf("tf-testaccdns%d", rand+1),
						"lang":              "en",
					}),
				),
			},
		},
	})
}

var AlidnsDomainGroupMap = map[string]string{}

func AlidnsDomainGroupBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
