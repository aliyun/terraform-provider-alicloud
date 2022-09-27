package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEssScalingConfigurationUpdate(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.ESSSupportRegions)
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^ubuntu",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_id":          "${data.alicloud_images.default.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.t5.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"password":          "123-abcABC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "password", "kms_encrypted_password", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"active": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"active": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_configuration_name": fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_configuration_name": fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_charge_type": "PayByTraffic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_charge_type": "PayByTraffic",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_out": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_out": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"credit_specification": "Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"credit_specification": "Standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_ssd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_name": "kms",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_name": "kms",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_description": "kms",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_description": "kms",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]string{{
						"size":                 "20",
						"category":             "cloud_ssd",
						"delete_with_instance": "false",
						"encrypted":            "true",
						"kms_key_id":           "${alicloud_kms_key.key.id}",
						"name":                 "kms",
						"description":          "kms",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "1",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_ssd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.kms_key_id":           CHECKSET,
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_data": `#!/bin/bash\necho \"hello\"\n`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data": "#!/bin/bash\necho \"hello\"\n",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name": "${alicloud_ram_role.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_name": "${alicloud_key_pair.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"name": "tf-test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.name": "tf-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":  REMOVEKEY,
					"instance_types": []string{"${data.alicloud_instance_types.t5.instance_types.0.id}", "${data.alicloud_instance_types.t5.instance_types.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":    REMOVEKEY,
						"instance_types.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id":  REMOVEKEY,
					"security_group_ids": []string{"${alicloud_security_group.default.id}", "${alicloud_security_group.default1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    REMOVEKEY,
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id": "${data.alicloud_images.default1.images.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id": REGEXMATCH + "^centos",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password_inherit": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":              REMOVEKEY,
					"instance_types":             REMOVEKEY,
					"internet_charge_type":       REMOVEKEY,
					"security_group_ids":         REMOVEKEY,
					"internet_max_bandwidth_out": REMOVEKEY,
					"credit_specification":       REMOVEKEY,
					"system_disk_category":       REMOVEKEY,
					"key_name":                   REMOVEKEY,
					"role_name":                  REMOVEKEY,
					"data_disk":                  REMOVEKEY,
					"tags":                       REMOVEKEY,
					"user_data":                  REMOVEKEY,
					"scaling_configuration_name": REMOVEKEY,
					"scaling_group_id":           "${alicloud_ess_scaling_group.default.id}",
					"image_id":                   "${data.alicloud_images.default.images.0.id}",
					"instance_type":              "${data.alicloud_instance_types.default.instance_types.0.id}",
					"security_group_id":          "${alicloud_security_group.default.id}",
					"force_delete":               "true",
					"override":                   "true",
					"system_disk_size":           "100",
					"password_inherit":           "false",
					"system_disk_name":           REMOVEKEY,
					"system_disk_description":    REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":                    REMOVEKEY,
						"instance_types.#":                 REMOVEKEY,
						"internet_charge_type":             REMOVEKEY,
						"security_group_ids.#":             REMOVEKEY,
						"internet_max_bandwidth_out":       REMOVEKEY,
						"credit_specification":             REMOVEKEY,
						"system_disk_category":             REMOVEKEY,
						"key_name":                         REMOVEKEY,
						"role_name":                        REMOVEKEY,
						"data_disk.#":                      REMOVEKEY,
						"data_disk.0.delete_with_instance": REMOVEKEY,
						"data_disk.0.size":                 REMOVEKEY,
						"data_disk.0.category":             REMOVEKEY,
						"data_disk.0.name":                 REMOVEKEY,
						"data_disk.0.encrypted":            REMOVEKEY,
						"data_disk.0.description":          REMOVEKEY,
						"data_disk.0.kms_key_id":           REMOVEKEY,
						"tags.name":                        REMOVEKEY,
						"user_data":                        REMOVEKEY,
						"scaling_configuration_name":       REMOVEKEY,
						"system_disk_size":                 "100",
						"system_disk_name":                 REMOVEKEY,
						"system_disk_description":          REMOVEKEY,
						"scaling_group_id":                 CHECKSET,
						"instance_type":                    CHECKSET,
						"security_group_id":                CHECKSET,
						"image_id":                         REGEXMATCH + "^ubuntu",
						"override":                         "true",
						"password_inherit":                 "false",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfigurationPerformanceLevel(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^ubuntu",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_id":          "${data.alicloud_images.default.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.default.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"password":          "123-abcABC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "password", "kms_encrypted_password", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"active": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"active": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_performance_level": "PL1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_performance_level": "PL1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_strategy": "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy":    "NoSpot",
						"spot_price_limit": NOSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_strategy": "SpotWithPriceLimit",
					"spot_price_limit": []map[string]string{{
						"instance_type": "${data.alicloud_instance_types.default.instance_types.0.id}",
						"price_limit":   "2.1",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy":      "SpotWithPriceLimit",
						"spot_price_limit.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]string{{
						"size":                 "20",
						"category":             "cloud_essd",
						"delete_with_instance": "false",
						"encrypted":            "true",
						"kms_key_id":           "${alicloud_kms_key.key.id}",
						"name":                 "kms",
						"description":          "kms",
						"performance_level":    "PL1",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "1",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_essd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.kms_key_id":           CHECKSET,
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEssScalingConfigurationMulti(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default.9"
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^ubuntu",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":             "10",
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_id":          "${data.alicloud_images.default.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.default.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"password":          "123-abcABC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceEssScalingConfigurationConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}
	
	resource "alicloud_key_pair" "default" {
  		key_pair_name = "${var.name}"
	}

	resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = true
	}

	resource "alicloud_security_group" "default1" {
	  name   = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	}

	data "alicloud_images" "default1" {
		name_regex  = "^centos.*_64"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "t5" {
      instance_type_family = "ecs.t5"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	resource "alicloud_kms_key" "key" {
		description             = var.name
		pending_window_in_days  = "7"
		key_state               = "Enabled"
	}
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

func TestAccAlicloudEssScalingConfigurationInstancePatternInfo(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^ubuntu",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_id":          "${data.alicloud_images.default.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.default.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"password":          "123-abcABC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "password", "kms_encrypted_password", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"active": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"active": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_performance_level": "PL1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_performance_level": "PL1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_strategy": "NoSpot",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy":    "NoSpot",
						"spot_price_limit": NOSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_strategy": "SpotWithPriceLimit",
					"spot_price_limit": []map[string]string{{
						"instance_type": "${data.alicloud_instance_types.default.instance_types.0.id}",
						"price_limit":   "2.1",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy":      "SpotWithPriceLimit",
						"spot_price_limit.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]string{{
						"size":                 "20",
						"category":             "cloud_essd",
						"delete_with_instance": "false",
						"encrypted":            "true",
						"kms_key_id":           "${alicloud_kms_key.key.id}",
						"name":                 "kms",
						"description":          "kms",
						"performance_level":    "PL1",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "1",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_essd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.kms_key_id":           CHECKSET,
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]string{{
						"instance_family_level": "EntryLevel",
						"cores":                 "4",
						"memory":                "4.0",
						"max_price":             "2.1",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#":                       "1",
						"instance_pattern_info.0.instance_family_level": "EntryLevel",
						"instance_pattern_info.0.cores":                 "4",
						"instance_pattern_info.0.memory":                "4.0",
						"instance_pattern_info.0.max_price":             "2.1",
					}),
				),
			},
		},
	})
}
