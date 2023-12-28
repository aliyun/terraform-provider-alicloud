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
					"vpc_id":      "${alicloud_vpc.default.id}",
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
	}

	resource "alicloud_vpc" "default" {
	  vpc_name    = var.name
	  enable_ipv6 = "true"
	  cidr_block = "172.16.0.0/12"
	}
	
	resource "alicloud_vswitch" "vsw" {
	  vpc_id = "${alicloud_vpc.default.id}"
	  cidr_block = "172.16.0.0/21"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  name = var.name
	  ipv6_cidr_block_mask = "22"
	}
	
	resource "alicloud_security_group" "group" {
	  name        = var.name
	  description = "foo"
	  vpc_id      = alicloud_vpc.default.id
	}
	
	data "alicloud_instance_types" "default" {
	  availability_zone = data.alicloud_zones.default.zones.0.id
	  system_disk_category = "cloud_efficiency"
	  cpu_core_count = 4
	  minimum_eni_ipv6_address_quantity = 1
	}
	
	data "alicloud_images" "default" {
	  name_regex  = "^ubuntu_18.*64"
	  most_recent = true
	  owners      = "system"
	}
	
	resource "alicloud_instance" "default" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  ipv6_address_count = 1
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  system_disk_category = "cloud_efficiency"
	  image_id = "${data.alicloud_images.default.images.0.id}"
	  instance_name = var.name
	  vswitch_id = "${alicloud_vswitch.vsw.id}"
	  internet_max_bandwidth_out = 10
	  security_groups = "${alicloud_security_group.group.*.id}"
	}`
