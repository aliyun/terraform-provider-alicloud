package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

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
	var commonBandwidthPackage vpc.DescribeCommonBandwidthPackagesResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "alicloud_common_bandwidth_package.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCommonBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackagePayByTraffic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "name", "tf_testAccCommonBandwidthPackagePayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_CommonBandwidthPackagePayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "ratio", "100"),
				),
			},
			{
				Config: testAccCommonBandwidthPackagePayByTrafficName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "name", "tf_testAccCommonBandwidthPackagePayByTrafficUpdate"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_CommonBandwidthPackagePayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "ratio", "100"),
				),
			},
			{
				Config: testAccCommonBandwidthPackagePayByTrafficDescription,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "name", "tf_testAccCommonBandwidthPackagePayByTrafficUpdate"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_CommonBandwidthPackagePayByTraffic_Update"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "ratio", "100"),
				),
			},
			{
				Config: testAccCommonBandwidthPackagePayByTrafficBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "bandwidth", "20"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "name", "tf_testAccCommonBandwidthPackagePayByTrafficUpdate"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_CommonBandwidthPackagePayByTraffic_Update"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "ratio", "100"),
				),
			},
			{
				Config: testAccCommonBandwidthPackagePayByTraffic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "name", "tf_testAccCommonBandwidthPackagePayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_CommonBandwidthPackagePayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "ratio", "100"),
				),
			},
		},
	})
}

func TestAccAlicloudCommonBandwidthPackage_PayByBandwidth(t *testing.T) {
	var commonBandwidthPackage vpc.DescribeCommonBandwidthPackagesResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		// module name
		IDRefreshName: "alicloud_common_bandwidth_package.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCommonBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackagePayByBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "name", "tf_testAccCommonBandwidthPackagePayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_CommonBandwidthPackagePayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "ratio", "100"),
				),
			},
			{
				Config: testAccCommonBandwidthPackagePayByBandwidthName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "name", "tf_testAccCommonBandwidthPackagePayByBandwidthUpdate"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_CommonBandwidthPackagePayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "ratio", "100"),
				),
			},
			{
				Config: testAccCommonBandwidthPackagePayByBandwidthDescription,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "name", "tf_testAccCommonBandwidthPackagePayByBandwidthUpdate"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_CommonBandwidthPackagePayByBandwidth_Update"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "ratio", "100"),
				),
			},
			{
				Config: testAccCommonBandwidthPackagePayByBandwidthBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "bandwidth", "20"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "name", "tf_testAccCommonBandwidthPackagePayByBandwidthUpdate"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_CommonBandwidthPackagePayByBandwidth_Update"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "ratio", "100"),
				),
			},
			{
				Config: testAccCommonBandwidthPackagePayByBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "name", "tf_testAccCommonBandwidthPackagePayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_CommonBandwidthPackagePayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo", "ratio", "100"),
				),
			},
		},
	})
}

func TestAccAlicloudCommonBandwidthPackage_Multi(t *testing.T) {
	var commonBandwidthPackage vpc.DescribeCommonBandwidthPackagesResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		// module name
		IDRefreshName: "alicloud_common_bandwidth_package.foo.9",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCommonBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo.9", &commonBandwidthPackage),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo.9", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo.9", "name", "tf_testAcc_CommonBandwidthPackageMulti_9"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo.9", "description", "tf_testAcc_common_bandwidth_package"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo.9", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr("alicloud_common_bandwidth_package.foo.9", "ratio", "100"),
				),
			},
		},
	})
}

func testAccCheckCommonBandwidthPackageExists(n string, commonBandwidthPackage *vpc.DescribeCommonBandwidthPackagesResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Common Bandwidth Package ID is set")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		commonBandwidthPackageService := CommonBandwidthPackageService{client}
		_, err := commonBandwidthPackageService.DescribeCommonBandwidthPackage(rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckCommonBandwidthPackageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	commonBandwidthPackageService := CommonBandwidthPackageService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_common_bandwidth_package" {
			continue
		}
		_, err := commonBandwidthPackageService.DescribeCommonBandwidthPackage(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe Common Bandwidth Package error %#v", err)
		}
	}
	return nil
}

const testAccCommonBandwidthPackagePayByTraffic = `
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "10"
  name = "tf_testAccCommonBandwidthPackagePayByTraffic"
  description = "tf_testAcc_CommonBandwidthPackagePayByTraffic"
}
`
const testAccCommonBandwidthPackagePayByTrafficName = `
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "10"
  name = "tf_testAccCommonBandwidthPackagePayByTrafficUpdate"
  description = "tf_testAcc_CommonBandwidthPackagePayByTraffic"
}
`
const testAccCommonBandwidthPackagePayByTrafficDescription = `
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "10"
  name = "tf_testAccCommonBandwidthPackagePayByTrafficUpdate"
  description = "tf_testAcc_CommonBandwidthPackagePayByTraffic_Update"
}
`
const testAccCommonBandwidthPackagePayByTrafficBandwidth = `
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "20"
  name = "tf_testAccCommonBandwidthPackagePayByTrafficUpdate"
  description = "tf_testAcc_CommonBandwidthPackagePayByTraffic_Update"
}
`

const testAccCommonBandwidthPackagePayByBandwidth = `
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "10"
  internet_charge_type = "PayByBandwidth"
  name = "tf_testAccCommonBandwidthPackagePayByBandwidth"
  description = "tf_testAcc_CommonBandwidthPackagePayByBandwidth"
}
`
const testAccCommonBandwidthPackagePayByBandwidthName = `
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "10"
  internet_charge_type = "PayByBandwidth"
  name = "tf_testAccCommonBandwidthPackagePayByBandwidthUpdate"
  description = "tf_testAcc_CommonBandwidthPackagePayByBandwidth"
}
`
const testAccCommonBandwidthPackagePayByBandwidthDescription = `
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "10"
  internet_charge_type = "PayByBandwidth"
  name = "tf_testAccCommonBandwidthPackagePayByBandwidthUpdate"
  description = "tf_testAcc_CommonBandwidthPackagePayByBandwidth_Update"
}
`
const testAccCommonBandwidthPackagePayByBandwidthBandwidth = `
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "20"
  internet_charge_type = "PayByBandwidth"
  name = "tf_testAccCommonBandwidthPackagePayByBandwidthUpdate"
  description = "tf_testAcc_CommonBandwidthPackagePayByBandwidth_Update"
}
`

const testAccCommonBandwidthPackageMulti = `
resource "alicloud_common_bandwidth_package" "foo" {
  count = 10
  bandwidth = "10"
  internet_charge_type = "PayByBandwidth"
  name = "tf_testAcc_CommonBandwidthPackageMulti_${count.index}"
  description = "tf_testAcc_common_bandwidth_package"
}

`
