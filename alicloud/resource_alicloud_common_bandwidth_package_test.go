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

func TestAccAlicloudCommonBandwidthPackage_basic(t *testing.T) {
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
			resource.TestStep{
				Config: testAccCommonBandwidthPackageConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageExists("alicloud_common_bandwidth_package.foo", &commonBandwidthPackage),
					resource.TestCheckResourceAttr(
						"alicloud_common_bandwidth_package.foo", "bandwidth", "100"),
					resource.TestCheckResourceAttr(
						"alicloud_common_bandwidth_package.foo", "name", "tf_testAcc_common_bandwidth_package"),
					resource.TestCheckResourceAttr(
						"alicloud_common_bandwidth_package.foo", "description", "tf_testAcc_common_bandwidth_package"),
					resource.TestCheckResourceAttr(
						"alicloud_common_bandwidth_package.foo", "internet_charge_type", "PayByTraffic"),
					resource.TestCheckResourceAttr(
						"alicloud_common_bandwidth_package.foo2", "internet_charge_type", "PayByBandwidth"),
					resource.TestCheckResourceAttr(
						"alicloud_common_bandwidth_package.foo3", "internet_charge_type", "PayByBandwidth"),
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

const testAccCommonBandwidthPackageConfig = `

resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "100"
  name = "tf_testAcc_common_bandwidth_package"
  description = "tf_testAcc_common_bandwidth_package"
}

resource "alicloud_common_bandwidth_package" "foo2" {
  bandwidth = "200"
  internet_charge_type = "PayByBandwidth"
  name = "tf_testAcc_common_bandwidth_package"
  description = "tf_testAcc_common_bandwidth_package"
}

resource "alicloud_common_bandwidth_package" "foo3" {
  bandwidth = "2"
  internet_charge_type = "PayByBandwidth"
  name = "tf_testAcc_common_bandwidth_package"
  description = "tf_testAcc_common_bandwidth_package"
}

`
