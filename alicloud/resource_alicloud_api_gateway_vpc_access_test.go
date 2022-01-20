package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_vpc_access", &resource.Sweeper{
		Name: "alicloud_api_gateway_vpc_access",
		F:    testSweepApiGatewayVpcAccess,
	})
}

func testSweepApiGatewayVpcAccess(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.ApiGatewayNoSupportedRegions) {
		log.Printf("[INFO] Skipping API Gateway unsupported region: %s", region)
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

	req := cloudapi.CreateDescribeVpcAccessesRequest()
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeVpcAccesses(req)
	})
	if err != nil {
		return fmt.Errorf("Error Describe Api Gateway Vpc: %s", err)
	}

	allVpcs, _ := raw.(*cloudapi.DescribeVpcAccessesResponse)

	swept := false

	for _, v := range allVpcs.VpcAccessAttributes.VpcAccessAttribute {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Api Gateway Vpc: %s", name)
			continue
		}
		swept = true

		req := cloudapi.CreateRemoveVpcAccessRequest()
		req.VpcId = v.VpcId
		req.InstanceId = v.InstanceId
		req.Port = requests.NewInteger(v.Port)
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.RemoveVpcAccess(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Api Gaiteway Vpc (%s): %s", name, err)
		}
	}
	if swept {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudApigatewayVpcAccess_basic(t *testing.T) {
	var v *cloudapi.VpcAccessAttribute
	resourceId := "alicloud_api_gateway_vpc_access.default"
	ra := resourceAttrInit(resourceId, apiGatewayVpcAccessMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sApiGatewayVpcAccess-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayVpcAccessConfigDependence)

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
					"name":        "${var.name}",
					"vpc_id":      "${data.alicloud_vpcs.default.ids.0}",
					"instance_id": "${alicloud_instance.default.id}",
					"port":        "8080",
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
		},
	})
}

func resourceApigatewayVpcAccessConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	%s
	`, name, ApigatewayVpcAccessConfigDependence)
}

var apiGatewayVpcAccessMap = map[string]string{
	"name":        CHECKSET,
	"vpc_id":      CHECKSET,
	"instance_id": CHECKSET,
	"port":        "8080",
}

const ApigatewayVpcAccessConfigDependence = `
	data "alicloud_zones" "default" {
	  available_disk_category = "cloud_efficiency"
	  available_resource_creation= "VSwitch"
	}

	data "alicloud_instance_types" "default" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

	data "alicloud_images" "default" {
	  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
	  most_recent = true
	  owners = "system"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
		vpc_id = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.0.id
	}
	
	resource "alicloud_vswitch" "vswitch" {
	  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	  vpc_id            = data.alicloud_vpcs.default.ids.0
	  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
	  zone_id           = data.alicloud_zones.default.zones.0.id
	  vswitch_name      = var.name
	}
	
	locals {
	  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
	}


	resource "alicloud_security_group" "default" {
	  name = "${var.name}"
	  description = "foo"
	  vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_instance" "default" {
	  vswitch_id = local.vswitch_id
	  image_id = "${data.alicloud_images.default.images.0.id}"

	  # series III
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  system_disk_category = "cloud_efficiency"

	  internet_charge_type = "PayByTraffic"
	  internet_max_bandwidth_out = 5
	  security_groups = ["${alicloud_security_group.default.id}"]
	  instance_name = "${var.name}"
	}`
