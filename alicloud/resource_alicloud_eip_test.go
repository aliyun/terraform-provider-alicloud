package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_eip", &resource.Sweeper{
		Name: "alicloud_eip",
		F:    testSweepEips,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_instance",
			"alicloud_slb",
			"alicloud_nat_gateway",
		},
	})
}

func testSweepEips(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var eips []vpc.EipAddress
	req := vpc.CreateDescribeEipAddressesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeEipAddresses(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving EIPs: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeEipAddressesResponse)
		if resp == nil || len(resp.EipAddresses.EipAddress) < 1 {
			break
		}
		eips = append(eips, resp.EipAddresses.EipAddress...)

		if len(resp.EipAddresses.EipAddress) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range eips {
		name := v.Name
		id := v.AllocationId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping EIP: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting EIP: %s (%s)", name, id)
		req := vpc.CreateReleaseEipAddressRequest()
		req.AllocationId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ReleaseEipAddress(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete EIP (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckEIPExists(n string, eip *vpc.EipAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No EIP ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}
		d, err := vpcService.DescribeEip(rs.Primary.ID)

		log.Printf("[WARN] eip id %#v", rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		*eip = d
		return nil
	}
}

func testAccCheckEIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_eip" {
			continue
		}

		_, err := vpcService.DescribeEip(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func TestAccAlicloudEipBasic_PayByBandwidth(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "alicloud_eip.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEipConfigBasic(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccCheckEipConfig_bandwidth(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_name(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAcceEipName%d", rand),
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_description(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAcceEipName%d_description", rand),
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_all(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcceEipName%d_all", rand),
						"description": fmt.Sprintf("tf-testAcceEipName%d_description_all", rand),
					}),
				),
			},
		},
	})

}

func TestAccAlicloudEipBasic_PayByTraffic(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "alicloud_eip.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEipConfigBasic(rand, "PayByTraffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type": "PayByTraffic",
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_bandwidth(rand, "PayByTraffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_name(rand, "PayByTraffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAcceEipName%d", rand),
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_description(rand, "PayByTraffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAcceEipName%d_description", rand),
					}),
				),
			},
			{
				Config: testAccCheckEipConfig_all(rand, "PayByTraffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcceEipName%d_all", rand),
						"description": fmt.Sprintf("tf-testAcceEipName%d_description_all", rand),
					}),
				),
			},
		},
	})

}

func TestAccAlicloudEipMulti(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "alicloud_eip.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckEipCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEipConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckEipConfigBasic(rand int, internet_charge_type string) string {
	return fmt.Sprintf(`
resource "alicloud_eip" "default" {
	instance_charge_type = "PostPaid"
	internet_charge_type = "%s"
	bandwidth = "5"
	period = "1"
}
`, internet_charge_type)
}

func testAccCheckEipConfig_bandwidth(rand int, internet_charge_type string) string {
	return fmt.Sprintf(`
resource "alicloud_eip" "default" {
	instance_charge_type = "PostPaid"
	internet_charge_type = "%s"
	bandwidth = "10"
	period = "1"
}
`, internet_charge_type)
}

func testAccCheckEipConfig_name(rand int, internet_charge_type string) string {
	return fmt.Sprintf(`
variable "name"{
	default = "tf-testAcceEipName%d"
}

resource "alicloud_eip" "default" {
	instance_charge_type = "PostPaid"
	internet_charge_type = "%s"
	bandwidth = "10"
	period = "1"
	name = "${var.name}"
}
`, rand, internet_charge_type)
}

func testAccCheckEipConfig_description(rand int, internet_charge_type string) string {
	return fmt.Sprintf(`
variable "name"{
	default = "tf-testAcceEipName%d"
}

resource "alicloud_eip" "default" {
	instance_charge_type = "PostPaid"
	internet_charge_type = "%s"
	bandwidth = "10"
	period = "1"
	name = "${var.name}"
    description = "${var.name}_description"
}
`, rand, internet_charge_type)
}

func testAccCheckEipConfig_all(rand int, internet_charge_type string) string {
	return fmt.Sprintf(`
variable "name"{
	default = "tf-testAcceEipName%d"
}

resource "alicloud_eip" "default" {
	instance_charge_type = "PostPaid"
	internet_charge_type = "%s"
	bandwidth = "10"
	period = "1"
	name = "${var.name}_all"
    description = "${var.name}_description_all"
}
`, rand, internet_charge_type)
}

func testAccCheckEipConfig_multi(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_eip" "default" {
    count = 10
	instance_charge_type = "PostPaid"
	internet_charge_type = "PayByBandwidth"
	bandwidth = "5"
	period = "1"
}
`)
}

var testAccCheckEipCheckMap = map[string]string{
	"name":                 "",
	"description":          "",
	"bandwidth":            "5",
	"instance_charge_type": "PostPaid",
	"internet_charge_type": "PayByBandwidth",
	// read method does't return a value for the period attribute, so it is not tested
	// "period" : "1",
	"ip_address": CHECKSET,
	"status":     CHECKSET,
}
