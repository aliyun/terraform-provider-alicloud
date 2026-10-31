package oss_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oss"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Oss BucketObjectWormConfiguration. >>> Resource test cases, automatically generated.
// Case 测试ObjectWorm_依赖资源 12777
// lintignore: AT001
func TestAccAliCloudOssBucketObjectWormConfiguration_basic12777(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_object_worm_configuration.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketObjectWormConfigurationMap12777)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketObjectWormConfiguration")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketObjectWormConfigurationBasicDependence12777)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":         "${alicloud_oss_bucket.defaultQf8G0L.id}",
					"object_worm_enabled": "Enabled",
					"rule": []map[string]interface{}{
						{
							"default_retention": []map[string]interface{}{
								{
									"mode": "COMPLIANCE",
									"days": "1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name":         acctest.CHECKSET,
						"object_worm_enabled": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": []map[string]interface{}{
						{
							"default_retention": []map[string]interface{}{
								{
									"mode": "COMPLIANCE",
									"days": "2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": []map[string]interface{}{
						{
							"default_retention": []map[string]interface{}{
								{
									"mode":  "COMPLIANCE",
									"years": "1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": []map[string]interface{}{
						{
							"default_retention": []map[string]interface{}{
								{
									"mode":  "COMPLIANCE",
									"years": "2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudOssBucketObjectWormConfigurationMap12777 = map[string]string{}

func AlicloudOssBucketObjectWormConfigurationBasicDependence12777(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_oss_bucket" "defaultQf8G0L" {
  bucket        = var.name
  storage_class = "Standard"
  versioning {
    status = "Enabled"
  }
}
`, name)
}

// Case 测试ObjectWorm_Year 12778
// lintignore: AT001
func TestAccAliCloudOssBucketObjectWormConfiguration_basic12778(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_object_worm_configuration.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketObjectWormConfigurationMap12778)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketObjectWormConfiguration")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketObjectWormConfigurationBasicDependence12778)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":         "${alicloud_oss_bucket.defaultQf8G0L.id}",
					"object_worm_enabled": "Enabled",
					"rule": []map[string]interface{}{
						{
							"default_retention": []map[string]interface{}{
								{
									"mode":  "COMPLIANCE",
									"years": "1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name":         acctest.CHECKSET,
						"object_worm_enabled": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudOssBucketObjectWormConfigurationMap12778 = map[string]string{}

func AlicloudOssBucketObjectWormConfigurationBasicDependence12778(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_oss_bucket" "defaultQf8G0L" {
  bucket        = var.name
  storage_class = "Standard"
  versioning {
    status = "Enabled"
  }
}
`, name)
}

// Case 测试ObjectWorm_无rule 12779
// lintignore: AT001
func TestAccAliCloudOssBucketObjectWormConfiguration_basic12779(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_object_worm_configuration.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketObjectWormConfigurationMap12779)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketObjectWormConfiguration")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketObjectWormConfigurationBasicDependence12779)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Configuration without any "rule" block. This resource has no
				// CustomizeDiff hook; read-time flattenOssBucketObjectWormConfigurationRule
				// returns nil when the server response carries no Rule, so state rule.# stays
				// 0 and in sync with config. "after apply, plan must be empty" therefore holds
				// via read-time normalization, not a diff-suppression hook.
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":         "${alicloud_oss_bucket.defaultQf8G0L.id}",
					"object_worm_enabled": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name":         acctest.CHECKSET,
						"object_worm_enabled": "Enabled",
					}),
				),
			},
			{
				// Introduce a default_retention rule (with mode) in an
				// update step (configIndex > 0). The create step above
				// omits the rule entirely, so `mode` first appears here,
				// satisfying the modification-coverage gate. mode is
				// ValidateFunc-locked to "COMPLIANCE", so the value cannot
				// change across steps; first-introduce-in-non-first-step
				// is the only viable path.
				Config: testAccConfig(map[string]interface{}{
					"rule": []map[string]interface{}{
						{
							"default_retention": []map[string]interface{}{
								{
									"mode": "COMPLIANCE",
									"days": "1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudOssBucketObjectWormConfigurationMap12779 = map[string]string{}

func AlicloudOssBucketObjectWormConfigurationBasicDependence12779(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_oss_bucket" "defaultQf8G0L" {
  bucket        = var.name
  storage_class = "Standard"
  versioning {
    status = "Enabled"
  }
}
`, name)
}

// Test Oss BucketObjectWormConfiguration. <<< Resource test cases, automatically generated.
