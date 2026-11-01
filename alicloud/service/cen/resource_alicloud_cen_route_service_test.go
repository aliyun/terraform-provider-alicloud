package cen_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cen"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cen_route_service", &resource.Sweeper{
		Name: "alicloud_cen_route_service",
		F:    testSweepCenRouteService,
	})
}

func testSweepCenRouteService(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_example",
		"tf-example",
		"terraform-example",
	}

	request := cbn.CreateDescribeCensRequest()
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var cenIds []string
	for {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCens(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve cen instance in service list: %s", err)
		}

		response, _ := raw.(*cbn.DescribeCensResponse)

		for _, v := range response.Cens.Cen {
			skip := true
			if !sweep.SweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(v.Name), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping cen instance: %s ", v.Name)
				}
			}
			cenIds = append(cenIds, v.CenId)
		}
		if len(response.Cens.Cen) < helper.PageSizeLarge {
			break
		}
		page, err := helper.GetNextpageNumber(request.PageNumber)
		if err != nil {
			return errs.WrapError(err)
		}
		request.PageNumber = page
	}

	for _, cenId := range cenIds {
		request := cbn.CreateDescribeRouteServicesInCenRequest()
		request.CenId = cenId
		request.PageSize = requests.NewInteger(helper.PageSizeLarge)
		request.PageNumber = requests.NewInteger(1)

		for {
			raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.DescribeRouteServicesInCen(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete cen route service(%s): %s.", cenId, err)
			}
			response, _ := raw.(*cbn.DescribeRouteServicesInCenResponse)

			for _, item := range response.RouteServiceEntries.RouteServiceEntry {
				id := fmt.Sprintf("%v:%v:%v:%v:%v", item.CenId, item.HostRegionId, item.Host, item.AccessRegionId, item.HostVpcId)
				if err := deleteCenRouteService(id, client); err != nil {
					log.Printf("[ERROR] Failed to delete cen route service (%s): %s.", cenId, err)
				} else {
					log.Printf("[INFO] Delete cen route service success: %s.", id)
				}
			}
			if len(response.RouteServiceEntries.RouteServiceEntry) < helper.PageSizeLarge {
				break
			}
			page, err := helper.GetNextpageNumber(request.PageNumber)
			if err != nil {
				return errs.WrapError(err)
			}
			request.PageNumber = page
		}
	}
	return nil
}

func TestAccAliCloudCenRouteService_basic(t *testing.T) {
	var v cbn.RouteServiceEntry
	resourceId := "alicloud_cen_route_service.default"
	ra := acctest.ResourceAttrInit(resourceId, CenRouteServiceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCbnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenRouteService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenRouteService%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, CenRouteServiceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_region_id": "${alicloud_cen_instance_attachment.vpc.child_instance_region_id}",
					"cen_id":           "${alicloud_cen_instance_attachment.vpc.instance_id}",
					"host":             "100.118.28.52/32",
					"host_region_id":   "${alicloud_cen_instance_attachment.vpc.child_instance_region_id}",
					"host_vpc_id":      "${alicloud_cen_instance_attachment.vpc.child_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_region_id": acctest.DefaultRegionToTest,
						"cen_id":           acctest.CHECKSET,
						"host":             "100.118.28.52/32",
						"host_region_id":   acctest.DefaultRegionToTest,
						"host_vpc_id":      acctest.CHECKSET,
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

func TestAccAliCloudCenRouteService_basic1(t *testing.T) {
	var v cbn.RouteServiceEntry
	resourceId := "alicloud_cen_route_service.default"
	ra := acctest.ResourceAttrInit(resourceId, CenRouteServiceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCbnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenRouteService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenRouteService%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, CenRouteServiceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_region_id": "${alicloud_cen_instance_attachment.vpc.child_instance_region_id}",
					"cen_id":           "${alicloud_cen_instance_attachment.vpc.instance_id}",
					"host":             "100.118.28.52/32",
					"host_region_id":   "${alicloud_cen_instance_attachment.vpc.child_instance_region_id}",
					"host_vpc_id":      "${alicloud_cen_instance_attachment.vpc.child_instance_id}",
					"description":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_region_id": acctest.DefaultRegionToTest,
						"cen_id":           acctest.CHECKSET,
						"host":             "100.118.28.52/32",
						"host_region_id":   acctest.DefaultRegionToTest,
						"host_vpc_id":      acctest.CHECKSET,
						"description":      name,
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

var CenRouteServiceMap = map[string]string{
	"status": acctest.CHECKSET,
}

func CenRouteServiceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
resource "alicloud_cen_instance" "default" {
    cen_instance_name = var.name
}
resource "alicloud_cen_instance_attachment" "vpc" {
    instance_id = alicloud_cen_instance.default.id
    child_instance_id = data.alicloud_vpcs.default.ids.0
	child_instance_type = "VPC"
    child_instance_region_id = "%s"
}
`, name, acctest.DefaultRegionToTest)
}

func deleteCenRouteService(id string, client *connectivity.AliyunClient) error {
	parts, err := helper.ParseResourceId(id, 5)
	if err != nil {
		return errs.WrapError(err)
	}
	request := cbn.CreateDeleteRouteServiceInCenRequest()
	request.AccessRegionId = parts[3]
	request.CenId = parts[0]
	request.Host = parts[2]
	request.HostRegionId = parts[1]
	request.HostVpcId = parts[4]
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DeleteRouteServiceInCen(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "InvalidOperation.CenInstanceStatus", "InvalidOperation.CloudRouteStatusNotAllow"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
