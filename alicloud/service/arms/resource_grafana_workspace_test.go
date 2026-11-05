package arms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudArmsGrafanaWorkspace_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsGrafanaWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsgrafanaworkspace%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_workspace_name":    name,
					"grafana_version":           "9.0.x",
					"grafana_workspace_edition": "standard",
					"password":                  "Arms@123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_workspace_name":    name,
						"grafana_version":           "9.0.x",
						"grafana_workspace_edition": "standard",
						"password":                  "Arms@123",
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
						"resource_group_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
						"resource_group_id":         acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AlicloudArmsGrafanaWorkspaceMap0 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsGrafanaWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsgrafanaworkspace%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "9.0.x",
					"description":               "grafana-rg-update-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
					"password":                  "Arms@123",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "9.0.x",
						"description":               "grafana-rg-update-test",
						"resource_group_id":         acctest.CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
						"password":                  "Arms@123",
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
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

// Test Arms GrafanaWorkspace. >>> Resource test cases, automatically generated.
// Case Grafana Terraform 20241211 9473
func TestAccAliCloudArmsGrafanaWorkspace_basic9473(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_grafana_workspace.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap9473)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsGrafanaWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence9473)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"resource_group_id":         acctest.CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
						"pricing_cycle":             "Month",
						"auto_renew":                "false",
						"aliyun_lang":               "zh",
						"duration":                  acctest.CHECKSET,
						"account_number":            acctest.CHECKSET,
						"custom_account_number":     acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap5667)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsGrafanaWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence5667)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-rg-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
					"password":                  "Arms@123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-rg-create-test",
						"resource_group_id":         acctest.CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
						"password":                  "Arms@123",
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap5591)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsGrafanaWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence5591)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-rg-create-test",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
					"password":                  "Arms@123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-rg-create-test",
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
						"password":                  "Arms@123",
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap5241)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsGrafanaWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence5241)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-rg-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
					"password":                  "Arms@123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-rg-create-test",
						"resource_group_id":         acctest.CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
						"password":                  "Arms@123",
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap4414)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsGrafanaWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence4414)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-rg-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
					"password":                  "Arms@123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-rg-create-test",
						"resource_group_id":         acctest.CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
						"password":                  "Arms@123",
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsGrafanaWorkspaceMap4029)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsGrafanaWorkspace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccarms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsGrafanaWorkspaceBasicDependence4029)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"grafana_version":           "8.2.x",
					"description":               "grafana-tag-create-test",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"grafana_workspace_edition": "standard",
					"grafana_workspace_name":    name,
					"password":                  "Arms@123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"grafana_version":           "8.2.x",
						"description":               "grafana-tag-create-test",
						"resource_group_id":         acctest.CHECKSET,
						"grafana_workspace_edition": "standard",
						"grafana_workspace_name":    name,
						"password":                  "Arms@123",
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
						"resource_group_id": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
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
