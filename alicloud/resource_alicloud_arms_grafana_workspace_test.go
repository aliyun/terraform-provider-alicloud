package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudArmsGrafanaWorkspace_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsGrafanaWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsgrafanaworkspace%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence0)
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
					"grafana_workspace_name":    name,
					"grafana_version":           "9.0.x",
					"grafana_workspace_edition": "standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_workspace_name":    name,
						"grafana_version":           "9.0.x",
						"grafana_workspace_edition": "standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version": "8.2.x",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version": "8.2.x",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "grafana-rg-create-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "grafana-rg-create-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "grafana-rg-update-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "grafana-rg-update-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_workspace_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_workspace_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version": "9.0.x",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version": "9.0.x",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-rg-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-rg-create-test",
						"resource_group_id":         CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudArmsGrafanaWorkspaceMap0 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudArmsGrafanaWorkspaceBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}

// Case 0  twin
func TestAccAliCloudArmsGrafanaWorkspace_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsGrafanaWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsgrafanaworkspace%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence0)
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
					"grafana_version":           "9.0.x",
					"description":               "grafana-rg-update-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "9.0.x",
						"description":               "grafana-rg-update-test",
						"resource_group_id":         CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
						"tags.%":                    "2",
						"tags.Created":              "TF",
						"tags.For":                  "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Arms GrafanaWorkspace. >>> Resource test cases, automatically generated.
// Case Grafana Terraform 20241211 9473
func TestAccAliCloudArmsGrafanaWorkspace_basic9473(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap9473)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsGrafanaWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence9473)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "9.0.x",
					"description":               "grafana-rg-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
					"pricing_cycle":             "Month",
					"auto_renew":                "false",
					"aliyun_lang":               "zh",
					"duration":                  "1",
					"account_number":            "10",
					"custom_account_number":     "0",
					"password":                  "Arms@123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "9.0.x",
						"description":               "grafana-rg-create-test",
						"resource_group_id":         CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
						"pricing_cycle":             "Month",
						"auto_renew":                "false",
						"aliyun_lang":               "zh",
						"duration":                  CHECKSET,
						"account_number":            CHECKSET,
						"custom_account_number":     CHECKSET,
						"password":                  "Arms@123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            "grafana-rg-update-test",
					"grafana_workspace_name": name + "_update",
					"aliyun_lang":            "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "grafana-rg-update-test",
						"grafana_workspace_name": name + "_update",
						"aliyun_lang":            "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version": "10.0.x",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version": "10.0.x",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_number", "aliyun_lang", "auto_renew", "custom_account_number", "duration", "password", "pricing_cycle"},
			},
		},
	})
}

var AlicloudArmsGrafanaWorkspaceMap9473 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudArmsGrafanaWorkspaceBasicDependence9473(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case Grafana Terraform接入 5667
func TestAccAliCloudArmsGrafanaWorkspace_basic5667(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap5667)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsGrafanaWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence5667)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-rg-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-rg-create-test",
						"resource_group_id":         CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            "grafana-rg-update-test",
					"grafana_workspace_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "grafana-rg-update-test",
						"grafana_workspace_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version": "9.0.x",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version": "9.0.x",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_number", "aliyun_lang", "auto_renew", "custom_account_number", "duration", "password", "pricing_cycle"},
			},
		},
	})
}

var AlicloudArmsGrafanaWorkspaceMap5667 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudArmsGrafanaWorkspaceBasicDependence5667(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case Grafana RMC接入_副本1703555887959 5591
func TestAccAliCloudArmsGrafanaWorkspace_basic5591(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap5591)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsGrafanaWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence5591)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-rg-create-test",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-rg-create-test",
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            "grafana-rg-update-test",
					"grafana_workspace_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "grafana-rg-update-test",
						"grafana_workspace_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version": "9.0.x",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version": "9.0.x",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_number", "aliyun_lang", "auto_renew", "custom_account_number", "duration", "password", "pricing_cycle"},
			},
		},
	})
}

var AlicloudArmsGrafanaWorkspaceMap5591 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudArmsGrafanaWorkspaceBasicDependence5591(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case Grafana RMC接入 5241
func TestAccAliCloudArmsGrafanaWorkspace_basic5241(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap5241)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsGrafanaWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence5241)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-rg-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-rg-create-test",
						"resource_group_id":         CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            "grafana-rg-update-test",
					"grafana_workspace_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "grafana-rg-update-test",
						"grafana_workspace_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version": "9.0.x",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version": "9.0.x",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_number", "aliyun_lang", "auto_renew", "custom_account_number", "duration", "password", "pricing_cycle"},
			},
		},
	})
}

var AlicloudArmsGrafanaWorkspaceMap5241 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudArmsGrafanaWorkspaceBasicDependence5241(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case Grafana资源组接入RG-新 4414
func TestAccAliCloudArmsGrafanaWorkspace_basic4414(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap4414)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsGrafanaWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence4414)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-rg-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-rg-create-test",
						"resource_group_id":         CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            "grafana-rg-update-test",
					"grafana_workspace_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "grafana-rg-update-test",
						"grafana_workspace_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version": "9.0.x",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version": "9.0.x",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_number", "aliyun_lang", "auto_renew", "custom_account_number", "duration", "password", "pricing_cycle"},
			},
		},
	})
}

var AlicloudArmsGrafanaWorkspaceMap4414 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudArmsGrafanaWorkspaceBasicDependence4414(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case Grafana标签接入-用例 4029
func TestAccAliCloudArmsGrafanaWorkspace_basic4029(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap4029)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsGrafanaWorkspace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence4029)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-tag-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-tag-create-test",
						"resource_group_id":         CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            "grafana-tag-update-test",
					"grafana_workspace_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "grafana-tag-update-test",
						"grafana_workspace_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version": "9.0.x",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version": "9.0.x",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_number", "aliyun_lang", "auto_renew", "custom_account_number", "duration", "password", "pricing_cycle"},
			},
		},
	})
}

var AlicloudArmsGrafanaWorkspaceMap4029 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudArmsGrafanaWorkspaceBasicDependence4029(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test Arms GrafanaWorkspace. <<< Resource test cases, automatically generated.
