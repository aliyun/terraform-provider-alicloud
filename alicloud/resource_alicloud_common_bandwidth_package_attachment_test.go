package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
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

func TestAccAlicloudCommonBandwidthPackageAttachmentBasic(t *testing.T) {
	var v vpc.CommonBandwidthPackage
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"bandwidth_package_id": CHECKSET,
		"instance_id":          CHECKSET,
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
		CheckDestroy:  testAccCheckCommonBandwidthPackageAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageAttachmentConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func TestAccAlicloudCommonBandwidthPackageAttachmentMulti(t *testing.T) {
	var v vpc.CommonBandwidthPackage
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_common_bandwidth_package_attachment.default.1"
	ra := resourceAttrInit(resourceId, map[string]string{
		"bandwidth_package_id": CHECKSET,
		"instance_id":          CHECKSET,
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
		CheckDestroy:  testAccCheckCommonBandwidthPackageAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageAttachmentConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckCommonBandwidthPackageAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	VpcService := VpcService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_common_bandwidth_package_attachment" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if len(parts) != 2 {
			return WrapError(Error("invalid resource id"))
		}
		_, err = VpcService.DescribeCommonBandwidthPackageAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func testAccCommonBandwidthPackageAttachmentConfigBasic(rand int) string {
	return fmt.Sprintf(`
    variable "name"{
    	default = "tf-testAccBandwidtchPackage%d"
    }

	resource "alicloud_common_bandwidth_package" "default" {
		bandwidth = 2
		internet_charge_type = "PayByBandwidth"
		name = "${var.name}"
		description = "${var.name}_description"
	}

	resource "alicloud_eip_address" "default" {
		address_name = "${var.name}"
		bandwidth            = "2"
		internet_charge_type = "PayByTraffic"
	}

	resource "alicloud_common_bandwidth_package_attachment" "default" {
		bandwidth_package_id = "${alicloud_common_bandwidth_package.default.id}"
		instance_id = "${alicloud_eip_address.default.id}"
	}
	`, rand)
}

func testAccCommonBandwidthPackageAttachmentConfigMulti(rand int) string {
	return fmt.Sprintf(`
    variable "name"{
    	default = "tf-testAccBandwidtchPackage%d"
    }

	variable "number" {
    	default = "2"
    }

	resource "alicloud_common_bandwidth_package" "default" {
		count = "${var.number}"
		bandwidth = 2
		internet_charge_type = "PayByBandwidth"
		name = "${var.name}"
		description = "${var.name}_description"
	}

	resource "alicloud_eip_address" "default" {
		count = "${var.number}"
		address_name = "${var.name}"
		bandwidth            = "2"
		internet_charge_type = "PayByTraffic"
	}

	resource "alicloud_common_bandwidth_package_attachment" "default" {
		count = "${var.number}"
		bandwidth_package_id = "${element(alicloud_common_bandwidth_package.default.*.id,count.index)}"
		instance_id = "${element(alicloud_eip_address.default.*.id,count.index)}"
	}
	`, rand)
}
