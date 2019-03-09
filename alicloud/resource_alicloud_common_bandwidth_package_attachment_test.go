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
	resource.AddTestSweepers("alicloud_common_bandwidth_package_attachment", &resource.Sweeper{
		Name: "alicloud_common_bandwidth_package_attachment",
		F:    testSweepCommonBandwidthPackageAttachment,
	})
}

func testSweepCommonBandwidthPackageAttachment(region string) error {
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
			return fmt.Errorf("Error retrieving CommonBandwidthPackage: %s", err)
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
		for _, eip := range cbwp.PublicIpAddresses.PublicIpAddresse {
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
			log.Printf("[INFO] Unassociating Common Bandwidth Package: %s (%s)", name, id)
			req := vpc.CreateRemoveCommonBandwidthPackageIpRequest()
			req.BandwidthPackageId = id
			req.IpInstanceId = eip.AllocationId
			_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.RemoveCommonBandwidthPackageIp(req)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to unassociate Common Bandwidth Package (%s (%s)): %s", name, id, err)
			}
		}
	}
	return nil
}

func TestAccAlicloudCommonBandwidthPackageAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "alicloud_common_bandwidth_package_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCommonBandwidthPackageAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCommonBandwidthPackageAttachmentExists("alicloud_common_bandwidth_package_attachment.foo"),
					resource.TestCheckResourceAttrSet(
						"alicloud_common_bandwidth_package_attachment.foo", "bandwidth_package_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_common_bandwidth_package_attachment.foo", "instance_id"),
				),
			},
		},
	})
}

func testAccCheckCommonBandwidthPackageAttachmentExists(n string) resource.TestCheckFunc {
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
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if len(parts) != 2 {
			return fmt.Errorf("invalid resource id")
		}
		err := commonBandwidthPackageService.DescribeCommonBandwidthPackageAttachment(parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("Describe Common Bandwidth Package attachment error %#v", err)
		}
		return nil
	}
}

func testAccCheckCommonBandwidthPackageAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	commonBandwidthPackageService := CommonBandwidthPackageService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_common_bandwidth_package_attachment" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if len(parts) != 2 {
			return fmt.Errorf("invalid resource id")
		}
		err := commonBandwidthPackageService.DescribeCommonBandwidthPackageAttachment(parts[0], parts[1])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe Common Bandwidth Package attachment error %#v", err)
		}
	}
	return nil
}

const testAccCommonBandwidthPackageAttachmentConfig = `

resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "2"
  name = "tf_testAcc_cbwp"
  description = "tf_testAcc_cbwp"
}

resource "alicloud_eip" "foo" {
  bandwidth            = "2"
  internet_charge_type = "PayByTraffic"
}

resource "alicloud_common_bandwidth_package_attachment" "foo" {
  bandwidth_package_id = "${alicloud_common_bandwidth_package.foo.id}"
  instance_id = "${alicloud_eip.foo.id}"
}
`
