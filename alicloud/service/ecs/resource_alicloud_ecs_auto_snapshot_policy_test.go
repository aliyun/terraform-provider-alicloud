package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Ecs AutoSnapshotPolicy. >>> Resource test cases, automatically generated.
// Case 自动快照策略 7893
func TestAccAliCloudECSAutoSnapshotPolicy_basic7893(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"1", "2"},
					"retention_days":  "1",
					"time_points":     []string{"1", "2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "2",
						"retention_days":    "1",
						"time_points.#":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_snapshot_policy_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_snapshot_policy_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"copied_snapshots_retention_days": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"copied_snapshots_retention_days": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_cross_region_copy": "true",
					"target_copy_regions":      []string{"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_cross_region_copy": "true",
						"target_copy_regions.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"copy_encryption_configuration": []map[string]interface{}{
						{
							"encrypted":  "true",
							"kms_key_id": "${alicloud_kms_key.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"copy_encryption_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"3", "4", "5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "3",
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
					"retention_days": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_days": "2",
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
				Config: testAccConfig(map[string]interface{}{
					"time_points": []string{"7", "6", "5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_points.#": "3",
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

func TestAccAliCloudECSAutoSnapshotPolicy_basic7893_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays":                 []string{"1", "2"},
					"retention_days":                  "1",
					"time_points":                     []string{"1", "2"},
					"auto_snapshot_policy_name":       name,
					"copied_snapshots_retention_days": "2",
					"copy_encryption_configuration": []map[string]interface{}{
						{
							"encrypted":  "true",
							"kms_key_id": "${alicloud_kms_key.default.id}",
						},
					},
					"enable_cross_region_copy": "true",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"target_copy_regions":      []string{"cn-hangzhou"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#":               "2",
						"retention_days":                  "1",
						"time_points.#":                   "2",
						"auto_snapshot_policy_name":       name,
						"copied_snapshots_retention_days": "2",
						"copy_encryption_configuration.#": "1",
						"enable_cross_region_copy":        "true",
						"resource_group_id":               CHECKSET,
						"target_copy_regions.#":           "1",
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
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

var AliCloudEcsAutoSnapshotPolicyMap7893 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudEcsAutoSnapshotPolicyBasicDependence7893(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = 7
  		key_state              = "Enabled"
	}
`, name)
}

// Case 适配废弃字段name 8035
func TestAccAliCloudECSAutoSnapshotPolicy_basic8035(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"1", "2"},
					"retention_days":  "1",
					"time_points":     []string{"1", "2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "2",
						"retention_days":    "1",
						"time_points.#":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"copied_snapshots_retention_days": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"copied_snapshots_retention_days": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_cross_region_copy": "true",
					"target_copy_regions":      []string{"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_cross_region_copy": "true",
						"target_copy_regions.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"copy_encryption_configuration": []map[string]interface{}{
						{
							"encrypted":  "true",
							"kms_key_id": "${alicloud_kms_key.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"copy_encryption_configuration.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays": []string{"3", "4", "5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#": "3",
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
					"retention_days": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_days": "2",
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
				Config: testAccConfig(map[string]interface{}{
					"time_points": []string{"7", "6", "5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_points.#": "3",
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

func TestAccAliCloudECSAutoSnapshotPolicy_basic8035_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_weekdays":                 []string{"1", "2"},
					"retention_days":                  "1",
					"time_points":                     []string{"1", "2"},
					"name":                            name,
					"copied_snapshots_retention_days": "2",
					"copy_encryption_configuration": []map[string]interface{}{
						{
							"encrypted":  "true",
							"kms_key_id": "${alicloud_kms_key.default.id}",
						},
					},
					"enable_cross_region_copy": "true",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"target_copy_regions":      []string{"cn-hangzhou"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_weekdays.#":               "2",
						"retention_days":                  "1",
						"time_points.#":                   "2",
						"name":                            name,
						"copied_snapshots_retention_days": "2",
						"copy_encryption_configuration.#": "1",
						"enable_cross_region_copy":        "true",
						"resource_group_id":               CHECKSET,
						"target_copy_regions.#":           "1",
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
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

// Case 关联实例标签自动快照策略：association_type=AssociatedWithInstanceTag + target_tags 创建/回读 no-drift
func TestAccAliCloudECSAutoSnapshotPolicy_associationInstanceTag(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_snapshot_policy_name": name,
					"repeat_weekdays":           []string{"1", "2"},
					"retention_days":            "1",
					"time_points":               []string{"1", "2"},
					"association_type":          "AssociatedWithInstanceTag",
					"target_tags": []map[string]interface{}{
						{
							"tag_key":   fmt.Sprintf("tf-testacc-tag-key-%d", rand),
							"tag_value": fmt.Sprintf("tf-testacc-tag-value-%d", rand),
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_type":        "AssociatedWithInstanceTag",
						"target_tags.#":           "1",
						"target_tags.0.tag_key":   fmt.Sprintf("tf-testacc-tag-key-%d", rand),
						"target_tags.0.tag_value": fmt.Sprintf("tf-testacc-tag-value-%d", rand),
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

// Case association_type ForceNew：改 association_type 触发 destroy+recreate（验证 ForceNew 生效、不走 Modify 永续 diff）
func TestAccAliCloudECSAutoSnapshotPolicy_associationForceNew(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_snapshot_policy_name": name,
					"repeat_weekdays":           []string{"1", "2"},
					"retention_days":            "1",
					"time_points":               []string{"1", "2"},
					"association_type":          "AssociatedWithDisk",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_type": "AssociatedWithDisk",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"association_type": "AssociatedWithInstanceTag",
					"target_tags": []map[string]interface{}{
						{
							"tag_key":   fmt.Sprintf("tf-testacc-fn-key-%d", rand),
							"tag_value": fmt.Sprintf("tf-testacc-fn-value-%d", rand),
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_type":        "AssociatedWithInstanceTag",
						"target_tags.#":           "1",
						"target_tags.0.tag_key":   fmt.Sprintf("tf-testacc-fn-key-%d", rand),
						"target_tags.0.tag_value": fmt.Sprintf("tf-testacc-fn-value-%d", rand),
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

// Case target_tags 可更新：改 target_tags 不重建、走 Modify 更新成功 + Read 回读新值
func TestAccAliCloudECSAutoSnapshotPolicy_targetTagsUpdate(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsAutoSnapshotPolicyMap7893)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsAutoSnapshotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsAutoSnapshotPolicyBasicDependence7893)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_snapshot_policy_name": name,
					"repeat_weekdays":           []string{"1", "2"},
					"retention_days":            "1",
					"time_points":               []string{"1", "2"},
					"association_type":          "AssociatedWithInstanceTag",
					"target_tags": []map[string]interface{}{
						{
							"tag_key":   fmt.Sprintf("tf-testacc-tu-key-%d", rand),
							"tag_value": "v1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_type":        "AssociatedWithInstanceTag",
						"target_tags.#":           "1",
						"target_tags.0.tag_key":   fmt.Sprintf("tf-testacc-tu-key-%d", rand),
						"target_tags.0.tag_value": "v1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_tags": []map[string]interface{}{
						{
							"tag_key":   fmt.Sprintf("tf-testacc-tu-key2-%d", rand),
							"tag_value": "v2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_type":        "AssociatedWithInstanceTag",
						"target_tags.#":           "1",
						"target_tags.0.tag_key":   fmt.Sprintf("tf-testacc-tu-key2-%d", rand),
						"target_tags.0.tag_value": "v2",
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

// Test Ecs AutoSnapshotPolicy. <<< Resource test cases, automatically generated.
