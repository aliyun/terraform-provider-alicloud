package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudHBRRestoreJob_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	checkoutAccount(t, true)
	defer checkoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_restore_job.default"
	ra := resourceAttrInit(resourceId, AlicloudHBRRestoreJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrRestoreJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrrestorejob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRRestoreJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_hash":         "${data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_hash}",
					"vault_id":              "${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}",
					"source_type":           "NAS",
					"restore_type":          "NAS",
					"snapshot_id":           "${data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_id}",
					"target_file_system_id": "${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}",
					"target_create_time":    "${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}",
					"target_path":           "/",
					"options":               "{\\\"includes\\\":[],\\\"excludes\\\":[]}",
					"include":               "[\\\"/proc\\\"]",
					"exclude":               "[\\\"/home\\\", \\\"/var/\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":  "NAS",
						"restore_type": "NAS",
						"target_path":  "/",
						"options":      "{\"includes\":[],\"excludes\":[]}",
						"include":      "[\"/proc\"]",
						"exclude":      "[\"/home\", \"/var/\"]",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"include", "exclude", "udm_region_id"},
			},
		},
	})
}

func TestAccAlicloudHBRRestoreJob_basic1(t *testing.T) {
	checkoutAccount(t, true)
	defer checkoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_restore_job.default"
	ra := resourceAttrInit(resourceId, AlicloudHBRRestoreJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrRestoreJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrrestorejob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRRestoreJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_hash": "${data.alicloud_hbr_snapshots.oss_snapshots.snapshots.0.snapshot_hash}",
					"vault_id":      "${data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id}",
					"source_type":   "OSS",
					"restore_type":  "OSS",
					"snapshot_id":   "${data.alicloud_hbr_snapshots.oss_snapshots.snapshots.0.snapshot_id}",
					"target_bucket": "${data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket}",
					"target_prefix": "",
					"options":       "{\\\"includes\\\":[],\\\"excludes\\\":[]}",
					"include":       "[\\\"/proc\\\"]",
					"exclude":       "[\\\"/home\\\", \\\"/var/\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":   "OSS",
						"restore_type":  "OSS",
						"target_prefix": "",
						"options":       "{\"includes\":[],\"excludes\":[]}",
						"include":       "[\"/proc\"]",
						"exclude":       "[\"/home\", \"/var/\"]",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"include", "exclude", "udm_region_id"},
			},
		},
	})
}

func TestAccAlicloudHBRRestoreJob_basic2(t *testing.T) {
	checkoutAccount(t, true)
	defer checkoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_restore_job.default"
	ra := resourceAttrInit(resourceId, AlicloudHBRRestoreJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrRestoreJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrrestorejob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRRestoreJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_hash":      "${data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_hash}",
					"vault_id":           "${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}",
					"source_type":        "ECS_FILE",
					"restore_type":       "ECS_FILE",
					"snapshot_id":        "${data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_id}",
					"target_instance_id": "${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}",
					"target_path":        "/",
					"include":            "[\\\"/proc\\\"]",
					"exclude":            "[\\\"/home\\\", \\\"/var/\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":  "ECS_FILE",
						"restore_type": "ECS_FILE",
						"target_path":  "/",
						"include":      "[\"/proc\"]",
						"exclude":      "[\"/home\", \"/var/\"]",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"include", "exclude", "udm_region_id"},
			},
		},
	})
}

func TestAccAlicloudHBRRestoreJob_basic3(t *testing.T) {
	checkoutAccount(t, true)
	defer checkoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_restore_job.default"
	ra := resourceAttrInit(resourceId, AlicloudHBRRestoreJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrRestoreJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrrestorejob%d", defaultRegionToTest, rand)
	ecsId := fmt.Sprintf("tf-testacc%d", rand+3)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRRestoreJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"restore_job_id":     ecsId,
					"snapshot_hash":      "${data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_hash}",
					"vault_id":           "${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}",
					"source_type":        "ECS_FILE",
					"restore_type":       "ECS_FILE",
					"snapshot_id":        "${data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_id}",
					"target_instance_id": "${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}",
					"target_path":        "/",
					"include":            "[\\\"/proc\\\"]",
					"exclude":            "[\\\"/home\\\", \\\"/var/\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"restore_job_id": ecsId,
						"source_type":    "ECS_FILE",
						"restore_type":   "ECS_FILE",
						"target_path":    "/",
						"include":        "[\"/proc\"]",
						"exclude":        "[\"/home\", \"/var/\"]",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"include", "exclude", "udm_region_id"},
			},
		},
	})
}

var AlicloudHBRRestoreJobMap0 = map[string]string{
	"include":       NOSET,
	"status":        CHECKSET,
	"exclude":       NOSET,
	"udm_detail":    NOSET,
	"udm_region_id": NOSET,
}

func AlicloudHBRRestoreJobBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_hbr_ecs_backup_plans" "default" {
    name_regex = "plan-tf-used-dont-delete"
}
data "alicloud_hbr_oss_backup_plans" "default" {
	name_regex = "plan-tf-used-dont-delete"
}
data "alicloud_hbr_nas_backup_plans" "default" {
	name_regex = "plan-tf-used-dont-delete"
}

data "alicloud_hbr_snapshots" "ecs_snapshots" {
    source_type  = "ECS_FILE"
	vault_id     =  data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id
	instance_id  =  data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id
}

data "alicloud_hbr_snapshots" "oss_snapshots" {
    source_type  = "OSS"
	vault_id     =  data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id
	bucket       =  data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket
}

data "alicloud_hbr_snapshots" "nas_snapshots" {
    source_type     = "NAS"
	vault_id        =  data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
	file_system_id  =  data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
    create_time     =  data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time
}

`, name)
}
