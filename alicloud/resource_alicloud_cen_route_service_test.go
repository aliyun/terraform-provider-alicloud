package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cen_route_service", &resource.Sweeper{
		Name: "alicloud_cen_route_service",
		F:    testSweepCenRouteService,
	})
}

func testSweepCenRouteService(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := cbn.CreateDescribeCensRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
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
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(v.Name), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping cen instance: %s ", v.Name)
			} else {
				cenIds = append(cenIds, v.CenId)
			}
		}
		if len(response.Cens.Cen) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	for _, cenId := range cenIds {
		request := cbn.CreateDescribeRouteServicesInCenRequest()
		request.CenId = cenId
		request.PageSize = requests.NewInteger(PageSizeLarge)
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
			if len(response.RouteServiceEntries.RouteServiceEntry) < PageSizeLarge {
				break
			}
			page, err := getNextpageNumber(request.PageNumber)
			if err != nil {
				return WrapError(err)
			}
			request.PageNumber = page
		}
	}
	return nil
}

func TestAccAlicloudCenRouteService_basic(t *testing.T) {
	var v cbn.RouteServiceEntry
	resourceId := "alicloud_cen_route_service.default"
	ra := resourceAttrInit(resourceId, CenRouteServiceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenRouteService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenRouteService%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CenRouteServiceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"access_region_id": defaultRegionToTest,
						"cen_id":           CHECKSET,
						"host":             "100.118.28.52/32",
						"host_region_id":   defaultRegionToTest,
						"host_vpc_id":      CHECKSET,
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

func TestAccAlicloudCenRouteService_basic1(t *testing.T) {
	var v cbn.RouteServiceEntry
	resourceId := "alicloud_cen_route_service.default"
	ra := resourceAttrInit(resourceId, CenRouteServiceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenRouteService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenRouteService%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CenRouteServiceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"access_region_id": defaultRegionToTest,
						"cen_id":           CHECKSET,
						"host":             "100.118.28.52/32",
						"host_region_id":   defaultRegionToTest,
						"host_vpc_id":      CHECKSET,
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
	"status": CHECKSET,
}

func CenRouteServiceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
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
`, name, defaultRegionToTest)
}

func deleteCenRouteService(id string, client *connectivity.AliyunClient) error {
	parts, err := ParseResourceId(id, 5)
	if err != nil {
		return WrapError(err)
	}
	request := cbn.CreateDeleteRouteServiceInCenRequest()
	request.AccessRegionId = parts[3]
	request.CenId = parts[0]
	request.Host = parts[2]
	request.HostRegionId = parts[1]
	request.HostVpcId = parts[4]
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DeleteRouteServiceInCen(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "InvalidOperation.CenInstanceStatus", "InvalidOperation.CloudRouteStatusNotAllow"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
