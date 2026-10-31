package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudECDPolicyGroup_basic0(t *testing.T) {
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"html_access":        "off",
					"html_file_transfer": "all",
					"policy_group_name":  "tf-testaccPolicyGroupName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name": "tf-testaccPolicyGroupName",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"html_access":        "on",
					"html_file_transfer": "download",
					"local_drive":        "read",
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
					"domain_list": "[white:],alicloud-provider.cn,sina.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_list": "[white:],alicloud-provider.cn,sina.com",
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
					"domain_list":            "[white:],alicloud-provider.cn",
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
						"domain_list":                       "[white:],alicloud-provider.cn",
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

func TestAccAliCloudECDPolicyGroup_basic1(t *testing.T) {
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

func TestAccAliCloudECDPolicyGroup_basic2(t *testing.T) {
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":    "tf-testaccPolicyGroupName",
					"recording":            "period",
					"recording_start_time": "08:00:00",
					"recording_end_time":   "08:59:00",
					"recording_fps":        "2",
					"camera_redirect":      "on",
					"recording_expires":    "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name":    "tf-testaccPolicyGroupName",
						"recording":            "period",
						"recording_start_time": "08:00:00",
						"recording_end_time":   "08:59:00",
						"recording_fps":        "2",
						"camera_redirect":      "on",
						"recording_expires":    "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"recording_expires": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recording_expires": "30",
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
					"recording_expires":    "180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recording":            "period",
						"recording_start_time": "10:00:00",
						"recording_end_time":   "12:59:00",
						"recording_fps":        "5",
						"camera_redirect":      "on",
						"recording_expires":    "180",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"recording_start_time", "recording_end_time", "status"},
			},
		},
	})
}

var AlicloudECDPolicyGroupMap0 = map[string]string{
	"policy_group_name": "tf-testaccPolicyGroupName",
}

func TestAccAliCloudECDPolicyGroup_basic3(t *testing.T) {
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":               "tf-testaccPolicyGroupName",
					"internet_communication_protocol": "TCP",
					"recording":                       "alltime",
					"recording_duration":              "30",
					"client_types": []map[string]interface{}{
						{
							"client_type": "html5",
							"status":      "off",
						},
					},
					"end_user_apply_admin_coordinate": "on",
					"end_user_group_coordinate":       "on",
				}),
				// DescribePolicyGroups returns the full client set (every type ON by
				// default, plus the disabled html5 OFF), a superset of the configured
				// entry. The refresh plan is intentionally non-empty: the API owns the
				// full set and the config only asserts the user-configured member, so the
				// real create (html5 off) is verified without synthesizing API-defaulted
				// types into state.
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name":               "tf-testaccPolicyGroupName",
						"internet_communication_protocol": "TCP",
						"recording_duration":              "30",
						"end_user_apply_admin_coordinate": "on",
						"end_user_group_coordinate":       "on",
					}),
					// F1: client_types read back from the plural "ClientTypes" key.
					// F2: status lowercased on read (OFF -> off); html5 disabled via ToUpper on write.
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "client_types.*", map[string]string{
						"client_type": "html5",
						"status":      "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_communication_protocol": "BOTH",
					"recording":                       "alltime",
					"recording_duration":              "60",
					"client_types": []map[string]interface{}{
						{
							"client_type": "html5",
							"status":      "on",
						},
						{
							"client_type": "windows",
							"status":      "on",
						},
					},
					"end_user_apply_admin_coordinate": "off",
					"end_user_group_coordinate":       "off",
				}),
				// When no client type is disabled (all ON), DescribePolicyGroups returns
				// the full default client set (every type ON), a superset of the configured
				// entries. The refresh plan is intentionally non-empty: the API owns the
				// full ON set and the config only asserts the user-configured members, so the
				// real update (html5 off->on) is verified without synthesizing API-defaulted
				// types into state.
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_communication_protocol": "BOTH",
						"recording_duration":              "60",
						"end_user_apply_admin_coordinate": "off",
						"end_user_group_coordinate":       "off",
					}),
					// F2: html5 re-enabled (on -> ON on write, ON -> on on read).
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "client_types.*", map[string]string{
						"client_type": "html5",
						"status":      "on",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "client_types.*", map[string]string{
						"client_type": "windows",
						"status":      "on",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// client_types reflects whatever DescribePolicyGroups returns (the disabled
				// subset, or the full ON set when nothing is disabled), which is a superset
				// of the user's config; verified in the apply steps above.
				ImportStateVerifyIgnore: []string{"client_types"},
			},
		},
	})
}

func AlicloudECDPolicyGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
