package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		if !sweepAll() {
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

func TestAccAliCloudDRDSInstance_Vpc(t *testing.T) {
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
						"description":   name,
						"mysql_version": "5",
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

func TestAccAliCloudDRDSInstance_Multi(t *testing.T) {
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
					"mysql_version":        "5",
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

func TestAccAliCloudDRDSInstance_VpcId(t *testing.T) {
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
					"vpc_id":               "${data.alicloud_vpcs.default.ids.0}",
					"mysql_version":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
						"vpc_id":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccAliCloudDRDSInstance_MySQLVersion(t *testing.T) {
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
					"mysql_version":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   name,
						"mysql_version": "5",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
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
	"connection_string":    CHECKSET,
	"port":                 CHECKSET,
}

// TestUnitAliCloudDRDSInstanceFlattenVips guards the DRDS instance Read against a
// panic when the DescribeDrdsInstance response returns an empty VIP list. A valid
// instance can transiently report no VIPs, and indexing Vip[0] directly used to
// crash the provider. The empty case must return zero values; the populated case
// must surface the first VIP's VpcId plus the intranet connection string/port.
func TestUnitAliCloudDRDSInstanceFlattenVips(t *testing.T) {
	// Empty VIP list: must not panic and must yield zero values.
	vpcId, connectionString, port := flattenDrdsInstanceVips([]drds.Vip{})
	if vpcId != "" || connectionString != "" || port != "" {
		t.Fatalf("empty VIP list should yield zero values, got vpcId=%q connectionString=%q port=%q", vpcId, connectionString, port)
	}

	// Nil VIP list: same guarantee.
	vpcId, connectionString, port = flattenDrdsInstanceVips(nil)
	if vpcId != "" || connectionString != "" || port != "" {
		t.Fatalf("nil VIP list should yield zero values, got vpcId=%q connectionString=%q port=%q", vpcId, connectionString, port)
	}

	// Populated VIP list: vpc_id from the first VIP, connection_string/port from the intranet VIP.
	vips := []drds.Vip{
		{Type: "internet", VpcId: "vpc-external", Dns: "public.example.com", Port: "3306"},
		{Type: "intranet", VpcId: "vpc-internal", Dns: "intranet.example.com", Port: "3307"},
	}
	vpcId, connectionString, port = flattenDrdsInstanceVips(vips)
	if vpcId != "vpc-external" {
		t.Fatalf("vpcId should come from the first VIP, got %q", vpcId)
	}
	if connectionString != "intranet.example.com" {
		t.Fatalf("connectionString should come from the intranet VIP, got %q", connectionString)
	}
	if port != "3307" {
		t.Fatalf("port should come from the intranet VIP, got %q", port)
	}
}
