package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test ESA CustomScenePolicy. >>> Resource test cases, automatically generated.
// Case 0
func TestAccAliCloudEsaCustomScenePolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_scene_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaCustomScenePolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCustomScenePolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scsp%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaCustomScenePolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_scene_policy_name": name,
					"end_time":                 "2025-09-08T18:00:00Z",
					"site_ids":                 "${data.alicloud_esa_sites.default.sites.0.id}",
					"start_time":               "2025-08-08T18:00:00Z",
					"template":                 "promotion",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_scene_policy_name": name,
						"end_time":                 "2025-09-08T18:00:00Z",
						"site_ids":                 CHECKSET,
						"start_time":               "2025-08-08T18:00:00Z",
						"template":                 "promotion",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_scene_policy_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_scene_policy_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time": "2025-10-08T18:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time": "2025-10-08T18:00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"site_ids": "${data.alicloud_esa_sites.default.sites.0.id},${data.alicloud_esa_sites.default.sites.1.id},${data.alicloud_esa_sites.default.sites.2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_ids": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"start_time": "2025-09-08T18:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"start_time": "2025-09-08T18:00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template": "promotion",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template": "promotion",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
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

func TestAccAliCloudEsaCustomScenePolicy_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_scene_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaCustomScenePolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCustomScenePolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scsp%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaCustomScenePolicyBasicDependence0)
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
					"custom_scene_policy_name": name,
					"end_time":                 "2025-09-08T18:00:00Z",
					"site_ids":                 "${data.alicloud_esa_sites.default.sites.0.id},${data.alicloud_esa_sites.default.sites.1.id},${data.alicloud_esa_sites.default.sites.2.id}",
					"start_time":               "2025-08-08T18:00:00Z",
					"template":                 "promotion",
					"status":                   "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_scene_policy_name": name,
						"end_time":                 "2025-09-08T18:00:00Z",
						"site_ids":                 CHECKSET,
						"start_time":               "2025-08-08T18:00:00Z",
						"template":                 "promotion",
						"status":                   "Disabled",
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

// Case 1
func TestAccAliCloudEsaCustomScenePolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_scene_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaCustomScenePolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCustomScenePolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scsp%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaCustomScenePolicyBasicDependence0)
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
					"custom_scene_policy_name": name,
					"end_time":                 "2025-09-08T18:00:00Z",
					"site_ids":                 "${data.alicloud_esa_sites.default.sites.0.id}",
					"create_time":              "2025-08-08T18:00:00Z",
					"template":                 "promotion",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_scene_policy_name": name,
						"end_time":                 "2025-09-08T18:00:00Z",
						"site_ids":                 CHECKSET,
						"create_time":              "2025-08-08T18:00:00Z",
						"template":                 "promotion",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_scene_policy_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_scene_policy_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time": "2025-10-08T18:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time": "2025-10-08T18:00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"site_ids": "${data.alicloud_esa_sites.default.sites.0.id},${data.alicloud_esa_sites.default.sites.1.id},${data.alicloud_esa_sites.default.sites.2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_ids": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"create_time": "2025-09-08T18:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"create_time": "2025-09-08T18:00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template": "promotion",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template": "promotion",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
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

func TestAccAliCloudEsaCustomScenePolicy_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_scene_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaCustomScenePolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCustomScenePolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scsp%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaCustomScenePolicyBasicDependence0)
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
					"custom_scene_policy_name": name,
					"end_time":                 "2025-09-08T18:00:00Z",
					"site_ids":                 "${data.alicloud_esa_sites.default.sites.0.id},${data.alicloud_esa_sites.default.sites.1.id},${data.alicloud_esa_sites.default.sites.2.id}",
					"create_time":              "2025-08-08T18:00:00Z",
					"template":                 "promotion",
					"status":                   "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_scene_policy_name": name,
						"end_time":                 "2025-09-08T18:00:00Z",
						"site_ids":                 CHECKSET,
						"create_time":              "2025-08-08T18:00:00Z",
						"template":                 "promotion",
						"status":                   "Disabled",
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

var AliCloudEsaCustomScenePolicyMap0 = map[string]string{
	"status": CHECKSET,
}

func AliCloudEsaCustomScenePolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}
`, name)
}

// Test ESA CustomScenePolicy. <<< Resource test cases, automatically generated.
