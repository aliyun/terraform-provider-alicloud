package ecd_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecd"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudECDPolicyGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_policy_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudECDPolicyGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdPolicyGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
				// Import records the legacy form: the server's DomainList is adopted
				// when the first import Read has no configuration and no recorded
				// representation, so the import verify compares domain_list directly.
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudECDPolicyGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdPolicyGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudECDPolicyGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdPolicyGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudECDPolicyGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdPolicyGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
				// DescribePolicyGroups returns every client type, but state mirrors the
				// configured subset, so the step converges without ExpectNonEmptyPlan.
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name":               "tf-testaccPolicyGroupName",
						"internet_communication_protocol": "TCP",
						"recording_duration":              "30",
						"end_user_apply_admin_coordinate": "on",
						"end_user_group_coordinate":       "on",
					}),
					resource.TestCheckResourceAttr(resourceId, "client_types.#", "1"),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_communication_protocol": "BOTH",
						"recording_duration":              "60",
						"end_user_apply_admin_coordinate": "off",
						"end_user_group_coordinate":       "off",
					}),
					resource.TestCheckResourceAttr(resourceId, "client_types.#", "2"),
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
				// An entry without a status must keep an identical element hash across
				// refresh: Read must not backfill the API status, otherwise the set
				// element identity never holds and every plan is non-empty.
				Config: testAccConfig(map[string]interface{}{
					"client_types": []map[string]interface{}{
						{
							"client_type": "windows",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_communication_protocol": "BOTH",
						"recording_duration":              "60",
					}),
					resource.TestCheckResourceAttr(resourceId, "client_types.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "client_types.*", map[string]string{
						"client_type": "windows",
					}),
				),
			},
			{
				// Removal transition: dropping every client_types block sends
				// ModifyPolicyGroup without ClientType parameters (the contract then
				// defaults to allowing every client type) and converges to an empty set.
				Config: testAccConfig(map[string]interface{}{
					"client_types": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_communication_protocol": "BOTH",
						"recording_duration":              "60",
					}),
					resource.TestCheckResourceAttr(resourceId, "client_types.#", "0"),
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

func TestUnitAliCloudEcdPolicyGroupClientTypesValidation(t *testing.T) {
	cases := []struct {
		name        string
		clientTypes []interface{}
		wantErr     bool
	}{
		{
			name:        "valid subset",
			clientTypes: []interface{}{map[string]interface{}{"client_type": "html5", "status": "off"}},
			wantErr:     false,
		},
		{
			name:        "valid entry without status",
			clientTypes: []interface{}{map[string]interface{}{"client_type": "windows"}},
			wantErr:     false,
		},
		{
			name:        "empty list is allowed",
			clientTypes: []interface{}{},
			wantErr:     false,
		},
		{
			name:        "empty client_type is rejected",
			clientTypes: []interface{}{map[string]interface{}{"client_type": "", "status": "on"}},
			wantErr:     true,
		},
		{
			name:        "nil element from an empty block is rejected",
			clientTypes: []interface{}{nil},
			wantErr:     true,
		},
		{
			name: "duplicate client_type is rejected",
			clientTypes: []interface{}{
				map[string]interface{}{"client_type": "html5", "status": "off"},
				map[string]interface{}{"client_type": "html5", "status": "on"},
			},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ecd.ValidatePolicyGroupClientTypes(c.clientTypes)
			switch {
			case c.wantErr && err == nil:
				t.Fatalf("expected an error, got nil")
			case !c.wantErr && err != nil:
				t.Fatalf("expected no error, got: %v", err)
			}
		})
	}
}

func TestAccAliCloudECDPolicyGroup_domainResolveRule(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_policy_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudECDPolicyGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdPolicyGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Structured domain resolve rules: enable and configure a blocklist.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_resolve_rule_type": "on",
					"domain_resolve_rule": []map[string]interface{}{
						{
							"domain":      "*.trycloudflare.com",
							"policy":      "block",
							"description": "blocked domain",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_rule_type": "on",
						"domain_resolve_rule.#":    "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "domain_resolve_rule.*", map[string]string{
						"domain":      "*.trycloudflare.com",
						"policy":      "block",
						"description": "blocked domain",
					}),
				),
			},
			{
				// Import of the structured form: the first import Read has no
				// configuration and no recorded representation, so it adopts the
				// structured form from the server return — the rules carry
				// descriptions only the structured form can represent. The import
				// verify no longer ignores the domain fields: state keeps the
				// rules that are actually in effect.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Update the rule set: replace the blocklist entry and add an allow entry.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_resolve_rule_type": "on",
					"domain_resolve_rule": []map[string]interface{}{
						{
							"domain":      "blocked.example.com",
							"policy":      "block",
							"description": "updated blocked domain",
						},
						{
							"domain":      "allowed.example.com",
							"policy":      "allow",
							"description": "allowed domain",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_rule_type": "on",
						"domain_resolve_rule.#":    "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Shrink the rule set from two entries to one. The submitted
				// DomainResolveRule list must define the full server-side list, so
				// the dropped entry stops being enforced after the update instead
				// of lingering server-side while the apply reports success.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_resolve_rule_type": "on",
					"domain_resolve_rule": []map[string]interface{}{
						{
							"domain":      "allowed.example.com",
							"policy":      "allow",
							"description": "allowed domain",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_rule_type": "on",
						"domain_resolve_rule.#":    "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "domain_resolve_rule.*", map[string]string{
						"domain": "allowed.example.com",
						"policy": "allow",
					}),
					testAccCheckEcdPolicyGroupDomainResolveRules(resourceId, []string{"allowed.example.com"}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Remove the switch and every rule block in the same apply: the
				// removal path must send DomainResolveRuleType=OFF and state must
				// converge to an empty configuration.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_resolve_rule_type": acctest.REMOVEKEY,
					"domain_resolve_rule":      acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_rule_type": "off",
						"domain_resolve_rule.#":    "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Re-enable with an explicit switch, then disable with an explicit
				// type=off to cover the configured-off path separately from removal.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_resolve_rule_type": "on",
					"domain_resolve_rule": []map[string]interface{}{
						{
							"domain": "reenabled.example.com",
							"policy": "block",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_rule_type": "on",
						"domain_resolve_rule.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_resolve_rule_type": "off",
					"domain_resolve_rule":      acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_rule_type": "off",
						"domain_resolve_rule.#":    "0",
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

func TestAccAliCloudECDPolicyGroup_domainListRemoval(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_policy_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudECDPolicyGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdPolicyGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Legacy domain_list: enable domain access control.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "tf-testaccPolicyGroupName",
					"domain_list":       "[black:],blocked.example.com,www.example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_list":              "[black:],blocked.example.com,www.example.com",
						"domain_resolve_rule_type": "on",
					}),
				),
			},
			{
				// Import records the legacy form: the first import Read adopts the
				// server's DomainList (its derived rules carry no descriptions, so
				// the lossless legacy form is chosen), so the import verify can
				// compare domain_list directly.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Adopt the imported legacy form: the configuration takes over the
				// same domain_list the import recorded, converging with no domain
				// diff instead of re-applying.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "tf-testaccPolicyGroupName",
					"domain_list":       "[black:],blocked.example.com,www.example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_list":              "[black:],blocked.example.com,www.example.com",
						"domain_resolve_rule_type": "on",
					}),
				),
			},
			{
				// The documented single-apply cleanup: drop the managed
				// domain_list and set domain_resolve_rule_type=off in the same
				// apply. The plan accepts it and ModifyPolicyGroup sends
				// DomainResolveRuleType=OFF, so state converges to empty.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_list":              acctest.REMOVEKEY,
					"domain_resolve_rule_type": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_rule_type": "off",
						"domain_list":              "",
						"domain_resolve_rule.#":    "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Re-enable the legacy form from the off state: the apply sends
				// DomainList without DomainResolveRuleType, the server turns itself
				// back on, and Read must store the echoed list and on switch instead
				// of clearing state per the prior off.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_list":              "[black:],reenabled.example.com",
					"domain_resolve_rule_type": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_list":              "[black:],reenabled.example.com",
						"domain_resolve_rule_type": "on",
						"domain_resolve_rule.#":    "0",
					}),
				),
			},
			{
				// The same configuration must produce an empty plan: the
				// re-enabled state matches the configuration, so the diff that
				// never converged while Read cleared the list is gone.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_list":              "[black:],reenabled.example.com",
					"domain_resolve_rule_type": acctest.REMOVEKEY,
				}),
			},
			{
				// Import records the re-enabled legacy form again, so the verify
				// compares domain_list directly.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Remove domain_list: the provider sends DomainResolveRuleType=OFF so
				// the configuration is actually cleared on the server instead of being
				// silently kept by ModifyPolicyGroup, and state converges to empty.
				// REMOVEKEY clears domain_list from the accumulated test config map so
				// the generated HCL genuinely omits it, exercising the removal path.
				// If the server did not converge to off, Read would keep mirroring
				// the rules it still reports and the removal would stay a visible
				// plan diff for the next apply.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "tf-testaccPolicyGroupName",
					"domain_list":       acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_rule_type": "off",
						"domain_list":              "",
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

// TestAccAliCloudECDPolicyGroup_domainUnmanagedUpdate covers domain access control
// the configuration does not manage, through the import that records it: adopting
// the recorded form keeps the server-side configuration untouched, while dropping
// it is a visible removal diff whose apply deliberately disables the control.
func TestAccAliCloudECDPolicyGroup_domainUnmanagedUpdate(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_policy_group.default"
	// No baseline attributes: the imported policy group has a generated name.
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdPolicyGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	// Renaming the group guarantees a plan diff against the imported name, so
	// the unrelated update always reaches ModifyPolicyGroup instead of being a
	// no-op apply.
	renamed := name + "-renamed"
	unmanagedId := ""
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
			unmanagedId = createPolicyGroupWithDomainList(t, name)
			t.Cleanup(func() {
				deletePolicyGroupIfExists(t, unmanagedId)
			})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Import a policy group that Terraform never created, so the domain
				// access control configuration is not managed by this configuration.
				// The first import Read adopts the server's representation: the
				// derived rules carry no descriptions, so the legacy DomainList is
				// recorded. ImportStatePersist keeps the imported state for the
				// following step; importing over an existing state is rejected.
				// An import step still needs the config to render the working
				// directory, and the import does not apply it.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": renamed,
					"local_drive":       "read",
				}),
				ResourceName:       resourceId,
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return unmanagedId, nil
				},
			},
			{
				// Adopt the recorded domain_list together with the rename: the
				// configuration manages the same value the import recorded, so the
				// plan has no domain removal diff and the update must not send
				// DomainResolveRuleType=OFF — the server keeps enforcing the
				// imported configuration. The applied rename proves the update
				// happened, so the domain assertions cannot pass on a no-op apply.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": renamed,
					"local_drive":       "read",
					"domain_list":       "[black:],blocked.example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name":        renamed,
						"local_drive":              "read",
						"domain_list":              "[black:],blocked.example.com",
						"domain_resolve_rule_type": "on",
					}),
					testAccCheckEcdPolicyGroupDomainControl(resourceId, "[black:],blocked.example.com"),
				),
			},
			{
				// Drop the adopted domain_list: the removal is a visible plan diff
				// — the import recorded the rules, so the plan shows exactly what
				// would be removed — and applying it deliberately disables the
				// domain access control instead of leaving it enforced but
				// invisible in state. REMOVEKEY clears domain_list from the
				// accumulated test config map so the generated HCL genuinely omits
				// it: without it the step would keep re-sending the adopted list
				// and the removal would never reach the plan.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": renamed,
					"local_drive":       "read",
					"domain_list":       acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name":        renamed,
						"local_drive":              "read",
						"domain_list":              "",
						"domain_resolve_rule_type": "off",
					}),
					testAccCheckEcdPolicyGroupDomainControlDisabled(resourceId),
				),
			},
		},
	})
}

// TestAccAliCloudECDPolicyGroup_domainUnmanagedRefresh covers domain access
// control that is enabled out of band on a group this configuration created
// but never managed: unlike the first import Read, a refresh whose state
// already records the switch must not adopt the server's rules into state,
// so the following unrelated update plans no domain removal and never sends
// DomainResolveRuleType=OFF — the out-of-band control stays enforced.
func TestAccAliCloudECDPolicyGroup_domainUnmanagedRefresh(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_policy_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudECDPolicyGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdPolicyGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	driftedId := ""
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create the group without managing any domain field: state records
				// only the switch the server reports for a fresh group, and no
				// domain form. The switch value itself is not asserted — only the
				// follow-up behavior that depends on the switch being recorded.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "tf-testaccPolicyGroupName",
				}),
				Check: resource.ComposeTestCheckFunc(
					captureEcdPolicyGroupId(resourceId, &driftedId),
				),
			},
			{
				// The control is re-enabled out of band. The refresh records the
				// switch the server now reports but adopts no rules, so the
				// unrelated rename shows no domain diff and the apply sends no
				// OFF: the server-side control must still be on afterwards. The
				// applied rename proves the update really happened, so the domain
				// assertions cannot pass on a no-op apply.
				PreConfig: func() {
					reEnablePolicyGroupDomainList(t, driftedId)
				},
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "tf-testaccPolicyGroupNameRenamed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_name":        "tf-testaccPolicyGroupNameRenamed",
						"domain_resolve_rule_type": "on",
						"domain_list":              "",
						"domain_resolve_rule.#":    "0",
					}),
					testAccCheckEcdPolicyGroupDomainControl(resourceId, "drifted.example.com"),
				),
			},
		},
	})
}

// testAccCheckEcdPolicyGroupDomainControl reads the policy group back from the
// service and asserts that domain access control is still enabled with the
// domains the group was created with: an unrelated update of an imported,
// unmanaged policy group must not clear the server-side configuration.
func testAccCheckEcdPolicyGroupDomainControl(resourceId, wantDomainList string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no ID in state", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := ecd.NewService(client).DescribeEcdPolicyGroup(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("reading the policy group %s from the service: %s", rs.Primary.ID, err)
		}
		ruleType, _ := object["DomainResolveRuleType"].(string)
		if strings.ToLower(ruleType) != "on" {
			return fmt.Errorf("server-side DomainResolveRuleType of %s is %q, want on", rs.Primary.ID, ruleType)
		}
		domainList, _ := object["DomainList"].(string)
		if !strings.Contains(domainList, wantDomainList) {
			return fmt.Errorf("server-side DomainList of %s is %q, want it to contain %q", rs.Primary.ID, domainList, wantDomainList)
		}
		return nil
	}
}

// testAccCheckEcdPolicyGroupDomainControlDisabled reads the policy group back from
// the service and asserts that domain access control is fully disabled: a removal
// that applied must actually have cleared the server-side configuration.
func testAccCheckEcdPolicyGroupDomainControlDisabled(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no ID in state", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := ecd.NewService(client).DescribeEcdPolicyGroup(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("reading the policy group %s from the service: %s", rs.Primary.ID, err)
		}
		ruleType, _ := object["DomainResolveRuleType"].(string)
		if strings.ToLower(ruleType) != "off" {
			return fmt.Errorf("server-side DomainResolveRuleType of %s is %q, want off", rs.Primary.ID, ruleType)
		}
		domainList, _ := object["DomainList"].(string)
		if domainList != "" {
			return fmt.Errorf("server-side DomainList of %s is %q, want empty", rs.Primary.ID, domainList)
		}
		return nil
	}
}

// testAccCheckEcdPolicyGroupDomainResolveRules reads the policy group back from
// the service and asserts the server-side DomainResolveRule entries are exactly
// the expected domains: shrinking the configured rule set must drop the removed
// entries server-side instead of leaving them enforced while the apply reports
// success.
func testAccCheckEcdPolicyGroupDomainResolveRules(resourceId string, wantDomains []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no ID in state", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := ecd.NewService(client).DescribeEcdPolicyGroup(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("reading the policy group %s from the service: %s", rs.Primary.ID, err)
		}
		rules, _ := object["DomainResolveRule"].([]interface{})
		gotDomains := make([]string, 0, len(rules))
		for _, rule := range rules {
			ruleMap, ok := rule.(map[string]interface{})
			if !ok {
				continue
			}
			gotDomains = append(gotDomains, fmt.Sprint(ruleMap["Domain"]))
		}
		if len(gotDomains) != len(wantDomains) {
			return fmt.Errorf("server-side DomainResolveRule of %s is %v, want exactly %v", rs.Primary.ID, gotDomains, wantDomains)
		}
		for _, wantDomain := range wantDomains {
			found := false
			for _, gotDomain := range gotDomains {
				if gotDomain == wantDomain {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("server-side DomainResolveRule of %s is %v, want exactly %v", rs.Primary.ID, gotDomains, wantDomains)
			}
		}
		return nil
	}
}

// createPolicyGroupWithDomainList creates a policy group with domain access
// control enabled outside Terraform, so the regression test can import a policy
// group whose domain configuration this configuration does not manage.
func createPolicyGroupWithDomainList(t *testing.T, name string) string {
	client := ecdTestPolicyGroupClient(t)
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"Name":       name,
		"DomainList": "[black:],blocked.example.com",
	}
	response, err := client.RpcPost("ecd", "2020-09-30", "CreatePolicyGroup", nil, request, false)
	if err != nil {
		t.Fatalf("failed to create the unmanaged policy group: %s", err)
	}
	id := fmt.Sprint(response["PolicyGroupId"])
	if id == "" || id == "<nil>" {
		t.Fatalf("CreatePolicyGroup returned no PolicyGroupId: %v", response)
	}
	return id
}

func deletePolicyGroupIfExists(t *testing.T, id string) {
	if id == "" {
		return
	}
	client := ecdTestPolicyGroupClient(t)
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"PolicyGroupId": []string{id},
	}
	if _, err := client.RpcPost("ecd", "2020-09-30", "DeletePolicyGroups", nil, request, false); err != nil {
		t.Logf("cleanup of policy group %s finished: %s", id, err)
	}
}

// ecdTestPolicyGroupClient builds a client from the test environment; the
// provider metadata is not configured yet when PreCheck runs.
func ecdTestPolicyGroupClient(t *testing.T) *connectivity.AliyunClient {
	rawClient, err := ecdTestClientForRegion(acctest.DefaultRegionToTest)
	if err != nil {
		t.Fatalf("failed to build an ECD client: %s", err)
	}
	client, ok := rawClient.(*connectivity.AliyunClient)
	if !ok {
		t.Fatalf("unexpected ECD client type %T", rawClient)
	}
	return client
}

// TestAccAliCloudECDPolicyGroup_domainDriftReenable covers the refresh semantics
// of domain_resolve_rule_type: the configuration keeps type=off while the domain
// access control is re-enabled out of band. Refresh must record the switch the
// server actually reports instead of the configured off, so the drift shows up as
// a plan diff and the apply re-sends DomainResolveRuleType=OFF.
func TestAccAliCloudECDPolicyGroup_domainDriftReenable(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_policy_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudECDPolicyGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdPolicyGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	driftedId := ""
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Manage the legacy form so the state records the domain
				// configuration, and capture the ID for the out-of-band drift.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "tf-testaccPolicyGroupName",
					"domain_list":       "[black:],drift.example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_list":              "[black:],drift.example.com",
						"domain_resolve_rule_type": "on",
					}),
					captureEcdPolicyGroupId(resourceId, &driftedId),
				),
			},
			{
				// The documented single-apply cleanup: the configuration keeps
				// only domain_resolve_rule_type=off.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_list":              acctest.REMOVEKEY,
					"domain_resolve_rule_type": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_rule_type": "off",
						"domain_list":              "",
					}),
				),
			},
			{
				// Out-of-band drift: the domain access control is re-enabled
				// directly through ModifyPolicyGroup while the configuration
				// stays at type=off. Refresh must record the server's "on" so
				// the plan is not empty: the drift is visible, not masked by
				// the configured off.
				PreConfig: func() {
					reEnablePolicyGroupDomainList(t, driftedId)
				},
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_resolve_rule_type": "off",
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the visible drift re-sends DomainResolveRuleType=OFF;
				// state and the server converge back to disabled.
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_resolve_rule_type": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_resolve_rule_type": "off",
						"domain_list":              "",
					}),
					testAccCheckEcdPolicyGroupDomainControlDisabled(resourceId),
				),
			},
		},
	})
}

// captureEcdPolicyGroupId records the resource ID of a step's state so a later
// PreConfig can act on the same instance out of band.
func captureEcdPolicyGroupId(resourceId string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no ID in state", resourceId)
		}
		*id = rs.Primary.ID
		return nil
	}
}

// reEnablePolicyGroupDomainList re-enables the domain access control of a policy
// group out of band, without Terraform, to simulate server-side drift. It waits
// for the service to report the re-enabled switch so the following refresh is
// guaranteed to observe the drift.
func reEnablePolicyGroupDomainList(t *testing.T, id string) {
	client := ecdTestPolicyGroupClient(t)
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"PolicyGroupId": id,
		"DomainList":    "[black:],drifted.example.com",
	}
	if _, err := client.RpcPost("ecd", "2020-09-30", "ModifyPolicyGroup", nil, request, false); err != nil {
		t.Fatalf("failed to re-enable domain access control out of band for %s: %s", id, err)
	}
	deadline := time.Now().Add(30 * time.Second)
	var lastErr error
	for {
		object, err := ecd.NewService(client).DescribeEcdPolicyGroup(id)
		if err == nil {
			if ruleType, _ := object["DomainResolveRuleType"].(string); strings.ToLower(ruleType) == "on" {
				return
			}
		} else {
			lastErr = err
		}
		if time.Now().After(deadline) {
			t.Fatalf("domain access control of %s did not re-enable out of band: %v", id, lastErr)
		}
		time.Sleep(3 * time.Second)
	}
}

// TestAccAliCloudECDPolicyGroup_domainResolveRuleGuard covers the plan-time
// cross-field checks: an enabled switch needs at least one rule, and every rule
// block needs a non-empty domain.
func TestAccAliCloudECDPolicyGroup_domainResolveRuleGuard(t *testing.T) {
	resourceId := "alicloud_ecd_policy_group.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_resolve_rule_type": "on",
				}),
				ExpectError: regexp.MustCompile("domain_resolve_rule_type=on requires at least one domain_resolve_rule block"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_resolve_rule_type": "on",
					"domain_resolve_rule": []map[string]interface{}{
						{
							"domain": "",
							"policy": "block",
						},
					},
				}),
				ExpectError: regexp.MustCompile("each domain_resolve_rule block requires a non-empty domain"),
			},
		},
	})
}

// TestAccAliCloudECDPolicyGroup_policyGroupNameComputed covers an omitted
// policy_group_name: CreatePolicyGroup generates the name and the state records
// the value read back from the service.
func TestAccAliCloudECDPolicyGroup_policyGroupNameComputed(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_policy_group.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdPolicyGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"html_access": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"html_access": "off",
					}),
					resource.TestMatchResourceAttr(resourceId, "policy_group_name", regexp.MustCompile(".+")),
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

// TestAccAliCloudECDPolicyGroup_domainFormSwitchGuard covers the plan-time
// rejection of switching between domain_list and domain_resolve_rule in one
// apply. ModifyPolicyGroup keeps omitted parameters, so a switch would leave the
// other representation on the server while Read hides it.
func TestAccAliCloudECDPolicyGroup_domainFormSwitchGuard(t *testing.T) {
	resourceId := "alicloud_ecd_policy_group.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdpolicygroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudECDPolicyGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name": "tf-testaccPolicyGroupName",
					"domain_list":       "[black:],switch.example.com",
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_name":        "tf-testaccPolicyGroupName",
					"domain_list":              acctest.REMOVEKEY,
					"domain_resolve_rule_type": "on",
					"domain_resolve_rule": []map[string]interface{}{
						{
							"domain": "switch.example.com",
							"policy": "block",
						},
					},
				}),
				ExpectError: regexp.MustCompile("cannot switch between domain_list and domain_resolve_rule in a single apply"),
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
