package cbwp_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cbwp"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_common_bandwidth_package_attachment", &resource.Sweeper{
		Name: "alicloud_common_bandwidth_package_attachment",
		F:    testSweepCommonBandwidthPackageAttachment,
	})
}

func testSweepCommonBandwidthPackageAttachment(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeCommonBandwidthPackages"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	type attachmentTarget struct {
		name string
		id   string
		eip  string
	}
	targets := make([]attachmentTarget, 0)
	for {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			return fmt.Errorf("Error retrieving CommonBandwidthPackage: %s", err)
		}
		resp, err := jsonpath.Get("$.CommonBandwidthPackages.CommonBandwidthPackage", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.CommonBandwidthPackages.CommonBandwidthPackage", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["Name"])
			id := fmt.Sprint(item["BandwidthPackageId"])
			publicIpAddresseList := make([]interface{}, 0)
			if publicIpAddresseRaw, ok := item["PublicIpAddresses"].(map[string]interface{}); ok {
				if list, ok := publicIpAddresseRaw["PublicIpAddresse"].([]interface{}); ok {
					publicIpAddresseList = list
				}
			}
			for _, e := range publicIpAddresseList {
				if m1, ok := e.(map[string]interface{}); ok {
					skip := true
					if !sweep.SweepAll() {
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
					targets = append(targets, attachmentTarget{name: name, id: id, eip: fmt.Sprint(m1["AllocationId"])})
				}
			}
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, target := range targets {
		log.Printf("[INFO] Unassociating Common Bandwidth Package: %s (%s)", target.name, target.id)
		req := map[string]interface{}{
			"BandwidthPackageId": target.id,
			"IpInstanceId":       target.eip,
		}
		action := "RemoveCommonBandwidthPackageIp"
		req["RegionId"] = client.RegionId
		_, err := client.RpcPost("Vpc", "2016-04-28", action, nil, req, false)
		if err != nil {
			log.Printf("[ERROR] Failed to unassociate Common Bandwidth Package (%s (%s)): %s", target.name, target.id, err)
		}
	}
	return nil
}

func TestAccAliCloudCommonBandwidthPackageAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_common_bandwidth_package_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cbwp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.default.id}",
					"instance_id":          "${alicloud_eip_address.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": acctest.CHECKSET,
						"instance_id":          acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cbwp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"bandwidth_package_id":        acctest.CHECKSET,
						"instance_id":                 acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cbwp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"bandwidth_package_id":        acctest.CHECKSET,
						"instance_id":                 acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cbwp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.default.id}",
					"instance_id":          "${alicloud_eip_address.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": acctest.CHECKSET,
						"instance_id":          acctest.CHECKSET,
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
					"bandwidth_package_bandwidth":                  acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cbwp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"bandwidth_package_id":        acctest.CHECKSET,
						"instance_id":                 acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cbwp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"bandwidth_package_id":        acctest.CHECKSET,
						"instance_id":                 acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCommonBandwidthPackageAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cbwp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageAttachment-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCommonBandwidthPackageAttachmentBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"bandwidth_package_id":        acctest.CHECKSET,
						"instance_id":                 acctest.CHECKSET,
						"bandwidth_package_bandwidth": "2",
					}),
				),
			},
		},
	})
}

var AliCloudCommonBandwidthPackageAttachmentMap0 = map[string]string{
	"status": acctest.CHECKSET,
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
