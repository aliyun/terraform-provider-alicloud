// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	crsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cr"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Cr ArtifactLifecycleRule. >>> Resource test cases, automatically generated.
// Case resource_ArtifactLifecycleRule_test 12942
func TestAccAliCloudCrArtifactLifecycleRule_basic12942(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_artifact_lifecycle_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCrArtifactLifecycleRuleMap12942)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrArtifactLifecycleRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCrArtifactLifecycleRuleBasicDependence12942)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto":                "true",
					"namespace_name":      "${alicloud_cr_ee_namespace.namespaceCase_20260611_ArtifactLifecycleRule_1.name}",
					"retention_tag_count": "30",
					"schedule_time":       "WEEK",
					"scope":               "REPO",
					"instance_id":         "${alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id}",
					"tag_regexp":          ".*",
					"repo_name":           "${alicloud_cr_ee_repo.repoCase_20260611_ArtifactLifecycleRule_1.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto":                "true",
						"namespace_name":      acctest.CHECKSET,
						"retention_tag_count": "30",
						"schedule_time":       "WEEK",
						"scope":               "REPO",
						"instance_id":         acctest.CHECKSET,
						"tag_regexp":          ".*",
						"repo_name":           acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_tag_count": "50",
					"schedule_time":       "MONTH",
					"scope":               "REPO",
					"tag_regexp":          "v.*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_tag_count": "50",
						"schedule_time":       "MONTH",
						"scope":               "REPO",
						"tag_regexp":          "v.*",
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

var AlicloudCrArtifactLifecycleRuleMap12942 = map[string]string{
	"create_time":                acctest.CHECKSET,
	"modified_time":              acctest.CHECKSET,
	"artifact_lifecycle_rule_id": acctest.CHECKSET,
}

func AlicloudCrArtifactLifecycleRuleBasicDependence12942(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "resourceCase_20260526_Mmd6on_1" {
  default_oss_bucket = "true"
  instance_name      = var.name
  renewal_status     = "ManualRenewal"
  image_scanner      = "DISABLE"
  period             = "1"
  payment_type       = "Subscription"
  instance_type      = "Economy"
}

resource "alicloud_cr_ee_namespace" "namespaceCase_20260611_ArtifactLifecycleRule_1" {
  instance_id        = alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_repo" "repoCase_20260611_ArtifactLifecycleRule_1" {
  instance_id = alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id
  namespace   = alicloud_cr_ee_namespace.namespaceCase_20260611_ArtifactLifecycleRule_1.name
  name        = var.name
  repo_type   = "PRIVATE"
  summary     = "test repository for lifecycle rule"
}


`, name)
}

// Case 保留策略生命周期 5221
func TestAccAliCloudCrArtifactLifecycleRule_basic5221(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_artifact_lifecycle_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCrArtifactLifecycleRuleMap5221)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrArtifactLifecycleRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCrArtifactLifecycleRuleBasicDependence5221)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto":                "false",
					"retention_tag_count": "30",
					"scope":               "INSTANCE",
					"instance_id":         "${alicloud_cr_ee_instance.defaultnKIyBE.id}",
					"tag_regexp":          " ",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto":                "false",
						"retention_tag_count": "30",
						"scope":               "INSTANCE",
						"instance_id":         acctest.CHECKSET,
						"tag_regexp":          " ",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_tag_count": "31",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_tag_count": "31",
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

var AlicloudCrArtifactLifecycleRuleMap5221 = map[string]string{
	"create_time":                acctest.CHECKSET,
	"modified_time":              acctest.CHECKSET,
	"artifact_lifecycle_rule_id": acctest.CHECKSET,
}

func AlicloudCrArtifactLifecycleRuleBasicDependence5221(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "defaultnKIyBE" {
  instance_name  = var.name
  renewal_status = "ManualRenewal"
  payment_type   = "Subscription"
  image_scanner  = "ACR"
  period         = "1"
  instance_type  = "Basic"
}


`, name)
}

// Case 保留策略生命周期_换账号可用_副本1709104286506 6046
func TestAccAliCloudCrArtifactLifecycleRule_basic6046(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_artifact_lifecycle_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCrArtifactLifecycleRuleMap6046)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrArtifactLifecycleRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCrArtifactLifecycleRuleBasicDependence6046)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto":                "true",
					"retention_tag_count": "30",
					"scope":               "REPO",
					"instance_id":         "${alicloud_cr_ee_instance.default2Rk8gT.id}",
					"tag_regexp":          " .*",
					"namespace_name":      "${alicloud_cr_ee_namespace.defaultevafKF.name}",
					"repo_name":           "${alicloud_cr_ee_repo.defaultKunw72.name}",
					"schedule_time":       "WEEK",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto":                "true",
						"retention_tag_count": "30",
						"scope":               "REPO",
						"instance_id":         acctest.CHECKSET,
						"tag_regexp":          " .*",
						"namespace_name":      acctest.CHECKSET,
						"repo_name":           acctest.CHECKSET,
						"schedule_time":       "WEEK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_tag_count": "31",
					"scope":               "INSTANCE",
					"tag_regexp":          "release-v.*",
					"schedule_time":       "MONTH",
					"namespace_name":      acctest.REMOVEKEY,
					"repo_name":           acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_tag_count": "31",
						"scope":               "INSTANCE",
						"tag_regexp":          "release-v.*",
						"schedule_time":       "MONTH",
						"namespace_name":      acctest.REMOVEKEY,
						"repo_name":           acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_tag_count": "28",
					"scope":               "REPO",
					"namespace_name":      "${alicloud_cr_ee_namespace.defaultGPiaHQ.name}",
					"repo_name":           "${alicloud_cr_ee_repo.defaultkCdOJ6.name}",
					"schedule_time":       "WEEK",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_tag_count": "28",
						"scope":               "REPO",
						"namespace_name":      acctest.CHECKSET,
						"repo_name":           acctest.CHECKSET,
						"schedule_time":       "WEEK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto":                "false",
					"retention_tag_count": "30",
					"tag_regexp":          " .*",
					"namespace_name":      "${alicloud_cr_ee_namespace.defaultevafKF.name}",
					"repo_name":           "${alicloud_cr_ee_repo.defaultKunw72.name}",
					"schedule_time":       acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto":                "false",
						"retention_tag_count": "30",
						"tag_regexp":          " .*",
						"namespace_name":      acctest.CHECKSET,
						"repo_name":           acctest.CHECKSET,
						"schedule_time":       acctest.REMOVEKEY,
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

var AlicloudCrArtifactLifecycleRuleMap6046 = map[string]string{
	"create_time":                acctest.CHECKSET,
	"modified_time":              acctest.CHECKSET,
	"artifact_lifecycle_rule_id": acctest.CHECKSET,
}

func AlicloudCrArtifactLifecycleRuleBasicDependence6046(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "default2Rk8gT" {
  instance_name  = var.name
  renewal_status = "ManualRenewal"
  payment_type   = "Subscription"
  image_scanner  = "ACR"
  period         = "1"
  instance_type  = "Basic"
}

resource "alicloud_cr_ee_namespace" "defaultevafKF" {
  instance_id        = alicloud_cr_ee_instance.default2Rk8gT.id
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "defaultKunw72" {
  instance_id = alicloud_cr_ee_instance.default2Rk8gT.id
  namespace   = alicloud_cr_ee_namespace.defaultevafKF.name
  name        = var.name
  repo_type   = "PUBLIC"
  summary     = "dd"
}

resource "alicloud_cr_ee_namespace" "defaultGPiaHQ" {
  instance_id        = alicloud_cr_ee_instance.default2Rk8gT.id
  name               = "${var.name}-2"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "defaultkCdOJ6" {
  instance_id = alicloud_cr_ee_instance.default2Rk8gT.id
  namespace   = alicloud_cr_ee_namespace.defaultGPiaHQ.name
  name        = "${var.name}-2"
  repo_type   = "PUBLIC"
  summary     = "dddd"
}


`, name)
}

// Case 保留策略生命周期_使用固定Instance-SUCC 5605
func TestAccAliCloudCrArtifactLifecycleRule_basic5605(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_artifact_lifecycle_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCrArtifactLifecycleRuleMap5605)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrArtifactLifecycleRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCrArtifactLifecycleRuleBasicDependence5605)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto":                "true",
					"retention_tag_count": "30",
					"scope":               "REPO",
					"instance_id":         "${alicloud_cr_ee_instance.default2Rk8gT.id}",
					"tag_regexp":          " .*",
					"namespace_name":      "${alicloud_cr_ee_namespace.defaultevafKF.name}",
					"repo_name":           "${alicloud_cr_ee_repo.defaultKunw72.name}",
					"schedule_time":       "WEEK",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto":                "true",
						"retention_tag_count": "30",
						"scope":               "REPO",
						"instance_id":         acctest.CHECKSET,
						"tag_regexp":          " .*",
						"namespace_name":      acctest.CHECKSET,
						"repo_name":           acctest.CHECKSET,
						"schedule_time":       "WEEK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_tag_count": "31",
					"scope":               "INSTANCE",
					"tag_regexp":          "release-v.*",
					"schedule_time":       "MONTH",
					"namespace_name":      acctest.REMOVEKEY,
					"repo_name":           acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_tag_count": "31",
						"scope":               "INSTANCE",
						"tag_regexp":          "release-v.*",
						"schedule_time":       "MONTH",
						"namespace_name":      acctest.REMOVEKEY,
						"repo_name":           acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_tag_count": "28",
					"scope":               "REPO",
					"namespace_name":      "${alicloud_cr_ee_namespace.defaultGPiaHQ.name}",
					"repo_name":           "${alicloud_cr_ee_repo.defaultkCdOJ6.name}",
					"schedule_time":       "WEEK",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_tag_count": "28",
						"scope":               "REPO",
						"namespace_name":      acctest.CHECKSET,
						"repo_name":           acctest.CHECKSET,
						"schedule_time":       "WEEK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto":                "false",
					"retention_tag_count": "30",
					"tag_regexp":          " .*",
					"namespace_name":      "${alicloud_cr_ee_namespace.defaultevafKF.name}",
					"repo_name":           "${alicloud_cr_ee_repo.defaultKunw72.name}",
					"schedule_time":       acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto":                "false",
						"retention_tag_count": "30",
						"tag_regexp":          " .*",
						"namespace_name":      acctest.CHECKSET,
						"repo_name":           acctest.CHECKSET,
						"schedule_time":       acctest.REMOVEKEY,
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

var AlicloudCrArtifactLifecycleRuleMap5605 = map[string]string{
	"create_time":                acctest.CHECKSET,
	"modified_time":              acctest.CHECKSET,
	"artifact_lifecycle_rule_id": acctest.CHECKSET,
}

func AlicloudCrArtifactLifecycleRuleBasicDependence5605(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "default2Rk8gT" {
  instance_name  = var.name
  renewal_status = "ManualRenewal"
  payment_type   = "Subscription"
  image_scanner  = "ACR"
  period         = "1"
  instance_type  = "Basic"
}

resource "alicloud_cr_ee_namespace" "defaultevafKF" {
  instance_id        = alicloud_cr_ee_instance.default2Rk8gT.id
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "defaultKunw72" {
  instance_id = alicloud_cr_ee_instance.default2Rk8gT.id
  namespace   = alicloud_cr_ee_namespace.defaultevafKF.name
  name        = var.name
  repo_type   = "PUBLIC"
  summary     = "dd"
}

resource "alicloud_cr_ee_namespace" "defaultGPiaHQ" {
  instance_id        = alicloud_cr_ee_instance.default2Rk8gT.id
  name               = "${var.name}-2"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "defaultkCdOJ6" {
  instance_id = alicloud_cr_ee_instance.default2Rk8gT.id
  namespace   = alicloud_cr_ee_namespace.defaultGPiaHQ.name
  name        = "${var.name}-2"
  repo_type   = "PUBLIC"
  summary     = "dddd"
}


`, name)
}

// Test Cr ArtifactLifecycleRule. <<< Resource test cases, automatically generated.
