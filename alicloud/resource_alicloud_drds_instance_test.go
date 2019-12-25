package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_drds_instance", &resource.Sweeper{
		Name: "alicloud_drds_instance",
		F:    testSweepDRDSInstances,
	})
}

func testSweepDRDSInstances(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DrdsSupportedRegions) {
		log.Printf("[INFO] Skipping DRDS Instance unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		fmt.Sprintf("tf-testAcc%s", region),
		fmt.Sprintf("tf_testAcc%s", region),
	}

	request := drds.CreateDescribeDrdsInstancesRequest()
	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstances(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving DRDS Instances: %s", WrapError(err))
	}
	response, _ := raw.(*drds.DescribeDrdsInstancesResponse)

	sweeped := false
	for _, v := range response.Instances.Instance {
		name := v.Description
		id := v.DrdsInstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping DRDS Instance: %s (%s)", name, id)
			continue
		}

		sweeped = true
		log.Printf("[INFO] Deleting DRDS Instance: %s (%s)", name, id)
		req := drds.CreateRemoveDrdsInstanceRequest()
		req.DrdsInstanceId = id
		_, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
			return drdsClient.RemoveDrdsInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete DRDS Instance (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudDRDSInstance_Basic(t *testing.T) {
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sDrdsdatabase-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			testAccPreCheckWithRegions(t, false, connectivity.DrdsClassicNoSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          var.name,
					"zone_id":              data.alicloud_zones.default.zones.0.id,
					"instance_series":      var.instance_series,
					"instance_charge_type": "PostPaid",
					"specification":        "drds.sn1.4c8g.8C16G",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": var.name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudDRDSInstance_Vpc(t *testing.T) {
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          var.name,
					"zone_id":              data.alicloud_zones.default.zones.0.id,
					"instance_series":      var.instance_series,
					"instance_charge_type": "PostPaid",
					"vswitch_id":           alicloud_vswitch.foo.id,
					"specification":        "drds.sn1.4c8g.8C16G",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": var.name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudDRDSInstance_Multi(t *testing.T) {
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alicloud_drds_instance.default.4"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sDrdsdatabase-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDRDSInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			testAccPreCheckWithRegions(t, false, connectivity.DrdsClassicNoSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}-${count.index}",
					"zone_id":              data.alicloud_zones.default.zones.0.id,
					"instance_series":      var.instance_series,
					"instance_charge_type": "PostPaid",
					"specification":        "drds.sn1.4c8g.8C16G",
					"count":                "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-4",
					}),
				),
			},
		},
	})
}

func resourceDRDSInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

variable "instance_series" {
	default = "drds.sn1.4c8g"
}

resource "alicloud_vpc" "foo" {
	name = var.name
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = alicloud_vpc.foo.id
 	cidr_block = "172.16.0.0/21"
 	availability_zone = data.alicloud_zones.default.zones.0.id
 	name = var.name
}

`, name)
}

var drdsInstancebasicMap = map[string]string{
	"description":          CHECKSET,
	"zone_id":              CHECKSET,
	"instance_series":      "drds.sn1.4c8g",
	"instance_charge_type": "PostPaid",
	"specification":        "drds.sn1.4c8g.8C16G",
}
