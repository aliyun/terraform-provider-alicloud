package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_pvtz_zone", &resource.Sweeper{
		Name: "alicloud_pvtz_zone",
		F:    testSweepPvtzZones,
	})
}

func testSweepPvtzZones(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	var zones []pvtz.Zone
	req := pvtz.CreateDescribeZonesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DescribeZones(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Private Zones: %s", err)
		}
		resp, _ := raw.(*pvtz.DescribeZonesResponse)
		if resp == nil || len(resp.Zones.Zone) < 1 {
			break
		}
		zones = append(zones, resp.Zones.Zone...)

		if len(resp.Zones.Zone) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}
	sweeped := false

	for _, v := range zones {
		name := v.ZoneName
		id := v.ZoneId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Private Zone: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Unbinding VPC from Private Zone: %s (%s)", name, id)
		request := pvtz.CreateBindZoneVpcRequest()
		request.ZoneId = id
		vpcs := make([]pvtz.BindZoneVpcVpcs, 0)
		request.Vpcs = &vpcs

		_, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.BindZoneVpc(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to unbind VPC from Private Zone (%s (%s)): %s ", name, id, err)
		}

		log.Printf("[INFO] Deleting Private Zone: %s (%s)", name, id)
		req := pvtz.CreateDeleteZoneRequest()
		req.ZoneId = id
		_, err = client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DeleteZone(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Private Zone (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudPvtzZone_basic(t *testing.T) {
	var v pvtz.DescribeZoneInfoResponse

	resourceId := "alicloud_pvtz_zone.default"
	ra := resourceAttrInit(resourceId, pvtzZoneBasicMap)

	serviceFunc := func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePvtzZoneConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"remark": "remark-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "remark-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "remark-test-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "remark-test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlicloudPvtzZone_multi(t *testing.T) {
	var v pvtz.DescribeZoneInfoResponse

	resourceId := "alicloud_pvtz_zone.default.4"
	ra := resourceAttrInit(resourceId, pvtzZoneBasicMap)

	serviceFunc := func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePvtzZoneConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":  fmt.Sprintf("tf-testacc%d${count.index}.test.com", rand),
					"count": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
func resourcePvtzZoneConfigDependence(name string) string {
	return ""
}

var pvtzZoneBasicMap = map[string]string{
	"name": CHECKSET,
}
