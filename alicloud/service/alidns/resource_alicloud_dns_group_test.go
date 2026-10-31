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
		"alicloud_dns_group",
		&resource.Sweeper{
			Name: "alicloud_dns_group",
			F:    testSweepDnsGroup,
		})
}

func testSweepDnsGroup(region string) error {
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
			if strings.HasPrefix(domainGroup.GroupName, "tf-testacc") {
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

func TestAccAlicloudAlidnsGroup_basic(t *testing.T) {
	resourceId := "alicloud_dns_group.default"
	var v alidnssdk.DomainGroup
	serviceFunc := func() interface{} {
		return alidns.NewDnsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(10000, 99999)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testaccdns%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDnsGroupConfigDependence)

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
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": fmt.Sprintf("tf-testaccdns%d", rand-1),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testaccdns%d", rand-1),
					}),
				),
			},
		},
	})
}

func resourceDnsGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}
