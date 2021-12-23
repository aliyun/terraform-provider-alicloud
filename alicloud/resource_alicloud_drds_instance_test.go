package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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
		"tf-testAcc",
		"tf_testAcc",
	}

	request := drds.CreateDescribeDrdsInstancesRequest()
	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstances(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving DRDS Instances: %s", WrapError(err))
	}
	response, _ := raw.(*drds.DescribeDrdsInstancesResponse)

	vpcService := VpcService{client}
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
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			instanceDetailRequest := drds.CreateDescribeDrdsInstanceRequest()
			instanceDetailRequest.DrdsInstanceId = id
			raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
				return drdsClient.DescribeDrdsInstance(instanceDetailRequest)
			})
			if err != nil {
				log.Printf("[ERROR] Error retrieving DRDS Instance: %s. %s", id, WrapError(err))
			}
			instanceDetailResponse, _ := raw.(*drds.DescribeDrdsInstanceResponse)
			for _, vip := range instanceDetailResponse.Data.Vips.Vip {
				if need, err := vpcService.needSweepVpc(vip.VpcId, ""); err == nil {
					skip = !need
					break
				}
			}

		}
		if skip {
			log.Printf("[INFO] Skipping DRDS Instance: %s (%s)", name, id)
			continue
		}

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
	return nil
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
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"vswitch_id":           "${data.alicloud_vswitches.default.vswitches.0.id}",
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
					"description": "${var.name}",
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

	resourceId := "alicloud_drds_instance.default.2"
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
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"zone_id":              "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"instance_series":      "${var.instance_series}",
					"instance_charge_type": "PostPaid",
					"specification":        "drds.sn1.4c8g.8C16G",
					"vswitch_id":           "${data.alicloud_vswitches.default.vswitches.0.id}",
					"count":                "3",
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
	
	data "alicloud_vpcs" "default"	{
        name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
	  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
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
