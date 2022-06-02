package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECDPolicyGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_policy_group.default"
	ra := resourceAttrInit(resourceId, AlicloudECDPolicyGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdPolicyGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "tf-testaccPolicyGroupName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name": "tf-testaccPolicyGroupName",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_drive": "read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_drive": "read",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_drive": "readwrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_drive": "readwrite",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_drive": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_drive": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"usb_redirect": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"usb_redirect": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"usb_redirect": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"usb_redirect": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "tf-testAccNameUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name": "tf-testAccNameUpdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"clipboard": "read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"clipboard": "read",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"clipboard": "readwrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"clipboard": "readwrite",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"clipboard": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"clipboard": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"watermark": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"watermark": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"watermark": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"watermark": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"watermark_transparency": "LIGHT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"watermark_transparency": "LIGHT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"watermark_transparency": "MIDDLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"watermark_transparency": "MIDDLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"watermark_transparency": "DARK",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"watermark_transparency": "DARK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"watermark": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"watermark": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"watermark_type": "HostName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"watermark_type": "HostName",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"watermark_type": "EndUserId",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"watermark_type": "EndUserId",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"authorize_security_policy_rules": []map[string]interface{}{
						{
							"type":        "inflow",
							"policy":      "accept",
							"description": "Terraform-Description",
							"port_range":  "43/43",
							"ip_protocol": "TCP",
							"priority":    "1",
							"cidr_ip":     "0.0.0.0/0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authorize_security_policy_rules.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"authorize_access_policy_rules": []map[string]interface{}{
						{
							"description": "Terraform-Description1",
							"cidr_ip":     "1.2.4.1/24",
						},
						{
							"description": "Terraform-Description2",
							"cidr_ip":     "1.2.4.2/24",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authorize_access_policy_rules.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"visual_quality": "low",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"visual_quality": "low",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"visual_quality": "lossless",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"visual_quality": "lossless",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"visual_quality": "high",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"visual_quality": "high",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_list": "[white:],baidu.com,sina.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_list": "[white:],baidu.com,sina.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_drive":       "off",
					"usb_redirect":      "off",
					"policy_group_name": "PolicyGroupNameAll",
					"clipboard":         "off",
					"watermark":         "on",
					"authorize_security_policy_rules": []map[string]interface{}{
						{
							"type":        "inflow",
							"policy":      "accept",
							"description": "Terraform-Description",
							"port_range":  "43/43",
							"ip_protocol": "TCP",
							"priority":    "1",
							"cidr_ip":     "0.0.0.0/3",
						},
						{
							"type":        "inflow",
							"policy":      "accept",
							"description": "Terraform-Description",
							"port_range":  "43/43",
							"ip_protocol": "TCP",
							"priority":    "1",
							"cidr_ip":     "0.0.0.0/4",
						},
					},
					"authorize_access_policy_rules": []map[string]interface{}{
						{
							"description": "Terraform-Description1",
							"cidr_ip":     "1.2.4.1/24",
						},
						{
							"description": "Terraform-Description2",
							"cidr_ip":     "1.2.4.2/24",
						},
					},
					"watermark_type":         "EndUserId",
					"domain_list":            "[white:],baidu.com",
					"watermark_transparency": "LIGHT",
					"visual_quality":         "medium",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_drive":                       "off",
						"usb_redirect":                      "off",
						"policy_group_name":                 "PolicyGroupNameAll",
						"clipboard":                         "off",
						"watermark":                         "on",
						"authorize_security_policy_rules.#": "2",
						"authorize_access_policy_rules.#":   "2",
						"watermark_type":                    "EndUserId",
						"domain_list":                       "[white:],baidu.com",
						"watermark_transparency":            "LIGHT",
						"visual_quality":                    "medium",
					}),
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

func TestAccAlicloudECDPolicyGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_policy_group.default"
	ra := resourceAttrInit(resourceId, AlicloudECDPolicyGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdPolicyGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":    "tf-testaccPolicyGroupName",
					"recording":            "period",
					"recording_start_time": "08:00:00",
					"recording_end_time":   "08:59:00",
					"recording_fps":        "2",
					"camera_redirect":      "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name":    "tf-testaccPolicyGroupName",
						"recording":            "period",
						"recording_start_time": "08:00:00",
						"recording_end_time":   "08:59:00",
						"recording_fps":        "2",
						"camera_redirect":      "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"recording": "alltime",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recording": "alltime",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"recording_fps": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recording_fps": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"camera_redirect": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"camera_redirect": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"recording":            "period",
					"recording_start_time": "10:00:00",
					"recording_end_time":   "12:59:00",
					"recording_fps":        "5",
					"camera_redirect":      "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recording":            "period",
						"recording_start_time": "10:00:00",
						"recording_end_time":   "12:59:00",
						"recording_fps":        "5",
						"camera_redirect":      "on",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"recording_start_time", "recording_end_time"},
			},
		},
	})
}

var AlicloudECDPolicyGroupMap0 = map[string]string{
	"policy_group_name": "tf-testaccPolicyGroupName",
}

func AlicloudECDPolicyGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
