package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_common_bandwidth_package", &resource.Sweeper{
		Name: "alicloud_common_bandwidth_package",
		F:    testSweepCommonBandwidthPackage,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_common_bandwidth_package_attachment",
		},
	})
}

func testSweepCommonBandwidthPackage(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var commonBandwidthPackages []vpc.CommonBandwidthPackage
	req := vpc.CreateDescribeCommonBandwidthPackagesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCommonBandwidthPackages(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving CommonBandwidthPackages: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeCommonBandwidthPackagesResponse)
		if resp == nil || len(resp.CommonBandwidthPackages.CommonBandwidthPackage) < 1 {
			break
		}
		commonBandwidthPackages = append(commonBandwidthPackages, resp.CommonBandwidthPackages.CommonBandwidthPackage...)

		if len(resp.CommonBandwidthPackages.CommonBandwidthPackage) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, cbwp := range commonBandwidthPackages {
		name := cbwp.Name
		id := cbwp.BandwidthPackageId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Common Bandwidth Package: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Common Bandwidth Package: %s (%s)", name, id)
		req := vpc.CreateDeleteCommonBandwidthPackageRequest()
		req.BandwidthPackageId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteCommonBandwidthPackage(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Common Bandwidth Package (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudCommonBandwidthPackage_PayByTraffic(t *testing.T) {

	var v vpc.CommonBandwidthPackage
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"bandwidth":            "10",
		"name":                 fmt.Sprintf("tf-testAccCommonBandwidthPackage%d", rand),
		"description":          "",
		"internet_charge_type": "PayByTraffic",
		"ratio":                "100",
	})
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCommonBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageBasic(rand, "PayByTraffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccCommonBandwidthPackageName(rand, "PayByTraffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccCommonBandwidthPackage%d_change", rand),
					}),
				),
			},
			{
				Config: testAccCommonBandwidthPackageDescription(rand, "PayByTraffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccCommonBandwidthPackage%d_description", rand),
					}),
				),
			},
			{
				Config: testAccCommonBandwidthPackageBandwidth(rand, "PayByTraffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "20",
					}),
				),
			},
			{
				Config: testAccCommonBandwidthPackageAll(rand, "PayByTraffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAccCommonBandwidthPackage%d_all", rand),
						"description": fmt.Sprintf("tf-testAccCommonBandwidthPackage%d_all", rand),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCommonBandwidthPackage_PayByBandwidth(t *testing.T) {

	var v vpc.CommonBandwidthPackage
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_common_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"bandwidth":            "10",
		"name":                 fmt.Sprintf("tf-testAccCommonBandwidthPackage%d", rand),
		"description":          "",
		"internet_charge_type": "PayByBandwidth",
		"ratio":                "100",
	})
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCommonBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageBasic(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccCommonBandwidthPackageName(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccCommonBandwidthPackage%d_change", rand),
					}),
				),
			},
			{
				Config: testAccCommonBandwidthPackageDescription(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccCommonBandwidthPackage%d_description", rand),
					}),
				),
			},
			{
				Config: testAccCommonBandwidthPackageBandwidth(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "20",
					}),
				),
			},
			{
				Config: testAccCommonBandwidthPackageAll(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAccCommonBandwidthPackage%d_all", rand),
						"description": fmt.Sprintf("tf-testAccCommonBandwidthPackage%d_all", rand),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCommonBandwidthPackage_Multi(t *testing.T) {
	var v vpc.CommonBandwidthPackage
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_common_bandwidth_package.default.9"
	ra := resourceAttrInit(resourceId, map[string]string{
		"bandwidth":            "10",
		"name":                 fmt.Sprintf("tf-testAccCommonBandwidthPackage%d", rand),
		"description":          "",
		"internet_charge_type": "PayByBandwidth",
		"ratio":                "100",
	})
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCommonBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageMulti(rand, "PayByBandwidth"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckCommonBandwidthPackageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	VpcService := VpcService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_common_bandwidth_package" {
			continue
		}
		_, err := VpcService.DescribeCommonBandwidthPackage(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe Common Bandwidth Package error %#v", err)
		}
	}
	return nil
}

func testAccCommonBandwidthPackageBasic(rand int, internetChargeType string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccCommonBandwidthPackage%d"
}

resource "alicloud_common_bandwidth_package" "default" {
  internet_charge_type = "%s"
  bandwidth = "10"
  name = "${var.name}"
}
`, rand, internetChargeType)
}
func testAccCommonBandwidthPackageName(rand int, internetChargeType string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccCommonBandwidthPackage%d"
}

resource "alicloud_common_bandwidth_package" "default" {
  internet_charge_type = "%s"
  bandwidth = "10"
  name = "${var.name}_change"
}
`, rand, internetChargeType)
}
func testAccCommonBandwidthPackageDescription(rand int, internetChargeType string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccCommonBandwidthPackage%d"
}

resource "alicloud_common_bandwidth_package" "default" {
  internet_charge_type = "%s"
  bandwidth = "10"
  name = "${var.name}_change"
  description = "${var.name}_description"
}
`, rand, internetChargeType)
}

func testAccCommonBandwidthPackageBandwidth(rand int, internetChargeType string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccCommonBandwidthPackage%d"
}

resource "alicloud_common_bandwidth_package" "default" {
  internet_charge_type = "%s"
  bandwidth = "20"
  name = "${var.name}_change"
  description = "${var.name}_description"
}
`, rand, internetChargeType)
}

func testAccCommonBandwidthPackageAll(rand int, internetChargeType string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccCommonBandwidthPackage%d"
}

resource "alicloud_common_bandwidth_package" "default" {
  internet_charge_type = "%s"
  bandwidth = "20"
  name = "${var.name}_all"
  description = "${var.name}_all"
}
`, rand, internetChargeType)
}

func testAccCommonBandwidthPackageMulti(rand int, internetChargeType string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccCommonBandwidthPackage%d"
}

resource "alicloud_common_bandwidth_package" "default" {
  count = 10
  internet_charge_type = "%s"
  bandwidth = "10"
  name = "${var.name}"
}
`, rand, internetChargeType)
}
