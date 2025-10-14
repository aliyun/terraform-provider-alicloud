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

func TestAccAliCloudCommonBandwidthPackageAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.default.id}",
					"instance_id":          "${alicloud_eip_address.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": CHECKSET,
						"instance_id":          CHECKSET,
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
					"bandwidth_package_bandwidth": "Cancelled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_bandwidth": "Cancelled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_bandwidth": "2",
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
				ImportStateVerifyIgnore: []string{"cancel_common_bandwidth_package_ip_bandwidth", "ip_type"},
			},
		},
	})
}

func TestAccAliCloudCommonBandwidthPackageAttachment_basic0_twin0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
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

func TestAccAliCloudCommonBandwidthPackageAttachment_basic0_twin1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id":        "${alicloud_common_bandwidth_package.default.id}",
					"instance_id":                 "${alicloud_eip_address.default.id}",
					"bandwidth_package_bandwidth": "Cancelled",
					"ip_type":                     "EIP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id":        CHECKSET,
						"instance_id":                 CHECKSET,
						"bandwidth_package_bandwidth": "Cancelled",
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

func TestAccAliCloudCommonBandwidthPackageAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.default.id}",
					"instance_id":          "${alicloud_eip_address.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": CHECKSET,
						"instance_id":          CHECKSET,
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
					"bandwidth_package_bandwidth":                  REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_bandwidth": "Cancelled",
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
				ImportStateVerifyIgnore: []string{"cancel_common_bandwidth_package_ip_bandwidth", "ip_type"},
			},
		},
	})
}

func TestAccAliCloudCommonBandwidthPackageAttachment_basic1_twin0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.default.id}",
					"instance_id":          "${alicloud_eip_address.default.id}",
					"ip_type":              "EIP",
					"cancel_common_bandwidth_package_ip_bandwidth": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id":        CHECKSET,
						"instance_id":                 CHECKSET,
						"bandwidth_package_bandwidth": "Cancelled",
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

func TestAccAliCloudCommonBandwidthPackageAttachment_basic1_twin1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id":        "${alicloud_common_bandwidth_package.default.id}",
					"instance_id":                 "${alicloud_eip_address.default.id}",
					"bandwidth_package_bandwidth": "Cancelled",
					"ip_type":                     "EIP",
					"cancel_common_bandwidth_package_ip_bandwidth": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id":        CHECKSET,
						"instance_id":                 CHECKSET,
						"bandwidth_package_bandwidth": "Cancelled",
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

func TestAccAliCloudCommonBandwidthPackageAttachment_Multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default.2"
	ra := resourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbwpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":                       "3",
					"bandwidth_package_id":        "${element(alicloud_common_bandwidth_package.default.*.id,count.index)}",
					"instance_id":                 "${element(alicloud_eip_address.default.*.id,count.index)}",
					"bandwidth_package_bandwidth": "2",
					"ip_type":                     "EIP",
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

var AliCloudCommonBandwidthPackageAttachmentMap0 = map[string]string{
	"status": CHECKSET,
}

func AliCloudCommonBandwidthPackageAttachmentBasicDependence0(name string) string {
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

func AliCloudCommonBandwidthPackageAttachmentBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_common_bandwidth_package" "default" {
  		count                = 3
  		bandwidth            = 3
  		internet_charge_type = "PayByBandwidth"
	}

	resource "alicloud_eip_address" "default" {
  		count                = 3
  		bandwidth            = "3"
  		internet_charge_type = "PayByTraffic"
	}
`, name)
}
