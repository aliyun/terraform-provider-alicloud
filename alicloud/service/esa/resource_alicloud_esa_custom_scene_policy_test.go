package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA CustomScenePolicy. >>> Resource test cases, automatically generated.
// Case 0
func TestAccAliCloudEsaCustomScenePolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_scene_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaCustomScenePolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaCustomScenePolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scsp%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaCustomScenePolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"site_ids":                 acctest.CHECKSET,
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
						"site_ids": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaCustomScenePolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaCustomScenePolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scsp%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaCustomScenePolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"site_ids":                 acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaCustomScenePolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaCustomScenePolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scsp%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaCustomScenePolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"site_ids":                 acctest.CHECKSET,
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
						"site_ids": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaCustomScenePolicyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaCustomScenePolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scsp%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaCustomScenePolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"site_ids":                 acctest.CHECKSET,
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
	"status": acctest.CHECKSET,
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
