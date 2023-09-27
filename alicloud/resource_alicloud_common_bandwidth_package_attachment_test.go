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
			if !sweepAll() {
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

func TestAccAliCloudCommonBandwidthPackageAttachment_Multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default.1"
	ra := resourceAttrInit(resourceId, resourceAlicloudCommonBandwidthPackageAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCommonBandwidthPackageAttachmentMultiDependence)

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
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id":        "${element(alicloud_common_bandwidth_package.default.*.id,count.index)}",
					"instance_id":                 "${element(alicloud_eip_address.default.*.id,count.index)}",
					"bandwidth_package_bandwidth": "2",
					"count":                       "${var.number}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id":        CHECKSET,
						"instance_id":                 CHECKSET,
						"bandwidth_package_bandwidth": "2",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCommonBandwidthPackageAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCommonBandwidthPackageAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCommonBandwidthPackageAttachmentBasicDependence)
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
					"bandwidth_package_id":        "${alicloud_common_bandwidth_package.default.id}",
					"instance_id":                 "${alicloud_eip_address.default.id}",
					"bandwidth_package_bandwidth": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id":        CHECKSET,
						"instance_id":                 CHECKSET,
						"bandwidth_package_bandwidth": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_bandwidth": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_bandwidth": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cancel_common_bandwidth_package_ip_bandwidth": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_bandwidth": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cancel_common_bandwidth_package_ip_bandwidth": "false",
					"bandwidth_package_bandwidth":                  "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_bandwidth": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cancel_common_bandwidth_package_ip_bandwidth"},
			},
		},
	})
}

func TestAccAliCloudCommonBandwidthPackageAttachment_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCommonBandwidthPackageAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCommonBandwidthPackageAttachmentBasicDependence)
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
					"bandwidth_package_id":        "${alicloud_common_bandwidth_package.default.id}",
					"instance_id":                 "${alicloud_eip_address.default.id}",
					"bandwidth_package_bandwidth": "2",
					"ip_type":                     "EIP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id":        CHECKSET,
						"instance_id":                 CHECKSET,
						"bandwidth_package_bandwidth": "2",
						"ip_type":                     "EIP",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cancel_common_bandwidth_package_ip_bandwidth", "ip_type"},
			},
		},
	})
}

var resourceAlicloudCommonBandwidthPackageAttachmentMap = map[string]string{}

func resourceAlicloudCommonBandwidthPackageAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_common_bandwidth_package" "default" {
  		bandwidth            = 3
  		internet_charge_type = "PayByBandwidth"
	}

	resource "alicloud_eip_address" "default" {
  		bandwidth            = "3"
  		internet_charge_type = "PayByTraffic"
	}
`, name)
}

func testAccCheckCommonBandwidthPackageAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cbwpServiceV2 := CbwpServiceV2{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_common_bandwidth_package_attachment" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if len(parts) != 2 {
			return WrapError(Error("invalid resource id"))
		}
		_, err = cbwpServiceV2.DescribeCommonBandwidthPackageAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}
func resourceAlicloudCommonBandwidthPackageAttachmentMultiDependence(name string) string {
	return fmt.Sprintf(`
    variable "name"{
    	default = "%s"
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
	`, name)
}
