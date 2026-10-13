// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Cr ArtifactLifecycleRule. >>> Resource test cases, automatically generated.
// Case resource_ArtifactLifecycleRule_test 12942
func TestAccAliCloudCrArtifactLifecycleRule_basic12942(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_artifact_lifecycle_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCrArtifactLifecycleRuleMap12942)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrArtifactLifecycleRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrArtifactLifecycleRuleBasicDependence12942)
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
						"namespace_name":      CHECKSET,
						"retention_tag_count": "30",
						"schedule_time":       "WEEK",
						"scope":               "REPO",
						"instance_id":         CHECKSET,
						"tag_regexp":          ".*",
						"repo_name":           CHECKSET,
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
	"create_time":                CHECKSET,
	"modified_time":              CHECKSET,
	"artifact_lifecycle_rule_id": CHECKSET,
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
	ra := resourceAttrInit(resourceId, AlicloudCrArtifactLifecycleRuleMap5221)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrArtifactLifecycleRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrArtifactLifecycleRuleBasicDependence5221)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"instance_id":         CHECKSET,
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
	"create_time":                CHECKSET,
	"modified_time":              CHECKSET,
	"artifact_lifecycle_rule_id": CHECKSET,
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
	ra := resourceAttrInit(resourceId, AlicloudCrArtifactLifecycleRuleMap6046)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrArtifactLifecycleRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrArtifactLifecycleRuleBasicDependence6046)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"instance_id":         CHECKSET,
						"tag_regexp":          " .*",
						"namespace_name":      CHECKSET,
						"repo_name":           CHECKSET,
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
					"namespace_name":      REMOVEKEY,
					"repo_name":           REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_tag_count": "31",
						"scope":               "INSTANCE",
						"tag_regexp":          "release-v.*",
						"schedule_time":       "MONTH",
						"namespace_name":      REMOVEKEY,
						"repo_name":           REMOVEKEY,
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
						"namespace_name":      CHECKSET,
						"repo_name":           CHECKSET,
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
					"schedule_time":       REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto":                "false",
						"retention_tag_count": "30",
						"tag_regexp":          " .*",
						"namespace_name":      CHECKSET,
						"repo_name":           CHECKSET,
						"schedule_time":       REMOVEKEY,
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
	"create_time":                CHECKSET,
	"modified_time":              CHECKSET,
	"artifact_lifecycle_rule_id": CHECKSET,
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
	ra := resourceAttrInit(resourceId, AlicloudCrArtifactLifecycleRuleMap5605)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrArtifactLifecycleRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrArtifactLifecycleRuleBasicDependence5605)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"instance_id":         CHECKSET,
						"tag_regexp":          " .*",
						"namespace_name":      CHECKSET,
						"repo_name":           CHECKSET,
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
					"namespace_name":      REMOVEKEY,
					"repo_name":           REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_tag_count": "31",
						"scope":               "INSTANCE",
						"tag_regexp":          "release-v.*",
						"schedule_time":       "MONTH",
						"namespace_name":      REMOVEKEY,
						"repo_name":           REMOVEKEY,
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
						"namespace_name":      CHECKSET,
						"repo_name":           CHECKSET,
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
					"schedule_time":       REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto":                "false",
						"retention_tag_count": "30",
						"tag_regexp":          " .*",
						"namespace_name":      CHECKSET,
						"repo_name":           CHECKSET,
						"schedule_time":       REMOVEKEY,
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
	"create_time":                CHECKSET,
	"modified_time":              CHECKSET,
	"artifact_lifecycle_rule_id": CHECKSET,
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
