package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_vswitch", &resource.Sweeper{
		Name: "alicloud_vswitch",
		F:    testSweepVSwitches,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_instance",
			"alicloud_db_instance",
			"alicloud_slb",
			"alicloud_ess_scalinggroup",
			"alicloud_fc_service",
			"alicloud_cs_cluster",
			"alicloud_kvstore_instance",
			"alicloud_route_table_attachment",
			//"alicloud_havip",
			"alicloud_network_interface",
			"alicloud_drds_instance",
			"alicloud_elasticsearch_instance",
			"alicloud_vpn_gateway",
			"alicloud_mongodb_instance",
			"alicloud_mongodb_sharding_instance",
			"alicloud_gpdb_instance",
		},
	})
}

func testSweepVSwitches(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var vswitches []vpc.VSwitch
	req := vpc.CreateDescribeVSwitchesRequest()
	req.RegionId = client.RegionId
	// API DescribeVSwitches has some limitations
	// If there is no vpc_id, setting PageSizeSmall can avoid ServiceUnavailable Error
	req.PageSize = requests.NewInteger(PageSizeSmall)
	req.PageNumber = requests.NewInteger(1)
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			rsp, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVSwitches(req)
			})
			raw = rsp
			return err
		}); err != nil {
			log.Printf("[ERROR] Error retrieving VSwitches: %s", WrapError(err))
		}
		resp, _ := raw.(*vpc.DescribeVSwitchesResponse)
		if resp == nil || len(resp.VSwitches.VSwitch) < 1 {
			break
		}
		vswitches = append(vswitches, resp.VSwitches.VSwitch...)

		if len(resp.VSwitches.VSwitch) < PageSizeSmall {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
		req.PageNumber = page
	}
	sweeped := false
	service := VpcService{client}
	for _, vsw := range vswitches {
		name := vsw.VSwitchName
		id := vsw.VSwitchId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a vswitch name is set by other service, it should be fetched by vpc name and deleted.
		if skip {
			if need, err := service.needSweepVpc(vsw.VpcId, ""); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping VSwitch: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting VSwitch: %s (%s)", name, id)
		if err := service.sweepVSwitch(id); err != nil {
			log.Printf("[ERROR] Failed to delete VSwitch (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func testAccCheckVSwitchExists(n string, vsw *vpc.DescribeVSwitchAttributesResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vswitch ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}
		instance, err := vpcService.DescribeVSwitch(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vsw = instance
		return nil
	}
}

func testAccCheckVSwitchDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vswitch" {
			continue
		}

		// Try to find the Vswitch
		if _, err := vpcService.DescribeVSwitch(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(Error("Vswitch still exist"))
	}

	return nil
}

func TestAccAlicloudVSwitchBasic(t *testing.T) {
	var v vpc.DescribeVSwitchAttributesResponse
	resourceId := "alicloud_vswitch.default"
	ra := resourceAttrInit(resourceId, testAccCheckVSwitchCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeVSwitch")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVSwitchConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccVswitchConfig%d", rand),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVSwitchConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccVswitchConfig%d_change", rand),
					}),
				),
			},
			{
				Config: testAccVSwitchConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccVswitchConfig%d_description", rand),
					}),
				),
			},
			{
				Config: testAccVSwitchConfig_tags(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccVSwitchConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         fmt.Sprintf("tf-testAccVswitchConfig%d_all", rand),
						"description":  fmt.Sprintf("tf-testAccVswitchConfig%d_description_all", rand),
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudVSwitchMulti(t *testing.T) {
	var v vpc.DescribeVSwitchAttributesResponse
	resourceId := "alicloud_vswitch.default.2"
	ra := resourceAttrInit(resourceId, testAccCheckVSwitchCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeVSwitch")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVSwitchConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block": "172.16.2.0/24",
						"name":       fmt.Sprintf("tf-testAccVswitchConfig%d", rand),
					}),
				),
			},
		},
	})
}

func testAccVSwitchConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
variable "name" {
  default = "tf-testAccVswitchConfig%d"
}
resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}
`, rand)
}

func testAccVSwitchConfig_name(rand int) string {
	return fmt.Sprintf(
		`
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
variable "name" {
  default = "tf-testAccVswitchConfig%d"
}
resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}_change"
}
`, rand)
}

func testAccVSwitchConfig_description(rand int) string {
	return fmt.Sprintf(
		`
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
variable "name" {
  default = "tf-testAccVswitchConfig%d"
}
resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}_change"
  description = "${var.name}_description"
}
`, rand)
}

func testAccVSwitchConfig_tags(rand int) string {
	return fmt.Sprintf(
		`
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
variable "name" {
  default = "tf-testAccVswitchConfig%d"
}
resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}_change"
  description = "${var.name}_description"
  tags 		= {
		Created = "TF"
		For 	= "acceptance test"
  }
}
`, rand)
}

func testAccVSwitchConfig_all(rand int) string {
	return fmt.Sprintf(
		`
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
variable "name" {
  default = "tf-testAccVswitchConfig%d"
}
resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}_all"
  description = "${var.name}_description_all"
}
`, rand)
}

func testAccVSwitchConfigMulti(rand int) string {
	return fmt.Sprintf(
		`
variable "number" {
	default = "3"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
variable "name" {
  default = "tf-testAccVswitchConfig%d"
}
resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  count = "${var.number}"
  vpc_id = "${ alicloud_vpc.default.id }"
  cidr_block = "172.16.${count.index}.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}
`, rand)
}

var testAccCheckVSwitchCheckMap = map[string]string{
	"vpc_id":            CHECKSET,
	"cidr_block":        "172.16.0.0/24",
	"availability_zone": CHECKSET,
	"description":       "",
}
