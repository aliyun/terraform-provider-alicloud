package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEssScalingConfiguration_override(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_name":        CHECKSET,
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
					"image_name":        "${data.alicloud_images.default2.images.0.name}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"password_inherit":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
					},
					"override": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "credit_specification", "security_group_id", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_other(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_name":        CHECKSET,
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
					"image_name":        "${data.alicloud_images.default2.images.0.name}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"password_inherit":  "true",
					"active":            "true",
					"is_outdated":       "false",
					"substitute":        "false",
					"enable":            "true",
					"instance_ids":      []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
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
				ImportStateVerifyIgnore: []string{"is_outdated", "override", "enable", "substitute", "force_delete", "instance_type", "credit_specification", "security_group_id", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_credit_specification(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, []connectivity.Region{"cn-beijing"})
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_name":        CHECKSET,
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependences_beijing)
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
					"image_name":        "${data.alicloud_images.default1.images.0.name}",
					"instance_type":     "${data.alicloud_instance_types.t6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"password_inherit":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
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
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "credit_specification", "security_group_id", "kms_encryption_context"},
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
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_imageName(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          CHECKSET,
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
					"image_id":          "${data.alicloud_images.default2.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id": "${data.alicloud_images.default2.images.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "image_id", "instance_type", "credit_specification", "security_group_id", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_imageId(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		//"image_name":        CHECKSET,
		"override": "false",
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
					"image_id":          "${data.alicloud_images.default2.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"password_inherit":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":   "${data.alicloud_images.default2.images.0.id}",
					"image_name": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id":   CHECKSET,
						"image_name": "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "credit_specification", "security_group_id", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_kms(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_name":        CHECKSET,
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
					"image_name":        "${data.alicloud_images.default2.images.0.name}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"password_inherit":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
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
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "credit_specification", "security_group_id", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_Update(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_name":        CHECKSET,
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
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"image_name":                "${data.alicloud_images.default2.images.0.name}",
					"instance_type":             "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id":         "${alicloud_security_group.default.id}",
					"force_delete":              "true",
					"internet_max_bandwidth_in": "1",
					//"password":          "123-abcABC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit":          "false",
						"internet_max_bandwidth_in": "1",
					}),
				),
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
					"system_disk_encrypted": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_encrypted": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_in": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_in": "2",
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
					"internet_max_bandwidth_out": "1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_out": "1024",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "credit_specification", "security_group_id", "kms_encryption_context"},
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
					"password_inherit": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "false",
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
					"system_disk_auto_snapshot_policy_id": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_auto_snapshot_policy_id": "",
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
					"instance_type":  REMOVEKEY,
					"instance_types": []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.1.id}"},
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
					"data_disk": []map[string]string{
						{
							"size":     "20",
							"category": "cloud_ssd",
							"name":     "${var.name}",
						},
						{
							"size":        "30",
							"category":    "cloud_ssd",
							"name":        "${var.name}",
							"description": "${var.name}",
						},
						{
							"size":                 "40",
							"category":             "cloud_efficiency",
							"delete_with_instance": "false",
							"name":                 "${var.name}",
							"description":          "${var.name}",
						},
						{
							"size":                    "50",
							"category":                "cloud_ssd",
							"delete_with_instance":    "false",
							"encrypted":               "true",
							"kms_key_id":              "${alicloud_kms_key.key.id}",
							"snapshot_id":             "",
							"auto_snapshot_policy_id": "",
							"device":                  "/dev/xvdb",
							"name":                    "${var.name}",
							"description":             "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "4",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_ssd",
						"data_disk.0.delete_with_instance": "true",
						"data_disk.0.encrypted":            "false",
						"data_disk.0.kms_key_id":           "",
						"data_disk.0.name":                 name,
						"data_disk.0.description":          "",

						"data_disk.1.size":                 "30",
						"data_disk.1.category":             "cloud_ssd",
						"data_disk.1.delete_with_instance": "true",
						"data_disk.1.encrypted":            "false",
						"data_disk.1.kms_key_id":           "",
						"data_disk.1.name":                 name,
						"data_disk.1.description":          name,

						"data_disk.2.size":                 "40",
						"data_disk.2.category":             "cloud_efficiency",
						"data_disk.2.delete_with_instance": "false",
						"data_disk.2.encrypted":            "false",
						"data_disk.2.kms_key_id":           "",
						"data_disk.2.name":                 name,
						"data_disk.2.description":          name,

						"data_disk.3.size":                    "50",
						"data_disk.3.category":                "cloud_ssd",
						"data_disk.3.delete_with_instance":    "false",
						"data_disk.3.encrypted":               "true",
						"data_disk.3.kms_key_id":              CHECKSET,
						"data_disk.3.snapshot_id":             "",
						"data_disk.3.auto_snapshot_policy_id": "",
						"data_disk.3.device":                  "/dev/xvdb",
						"data_disk.3.name":                    name,
						"data_disk.3.description":             name,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_PerformanceLevel(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun",
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
					"image_id":          "${data.alicloud_images.default2.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
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
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "kms_encryption_context"},
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
					"internet_max_bandwidth_in": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_in": "3",
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
					"spot_strategy": "SpotWithPriceLimit",
					"spot_price_limit": []map[string]string{{
						"instance_type": "${data.alicloud_instance_types.c6.instance_types.0.id}",
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
					"spot_strategy":    "NoSpot",
					"spot_price_limit": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_strategy":      "NoSpot",
						"spot_price_limit.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"category":             "cloud_ssd",
							"delete_with_instance": "false",
							"encrypted":            "true",
							"kms_key_id":           "${alicloud_kms_key.key.id}",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":     "20",
							"category": "cloud_ssd",
							"name":     "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_ssd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.kms_key_id":           CHECKSET,
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                 "20",
						"data_disk.1.category":             "cloud_ssd",
						"data_disk.1.delete_with_instance": "true",
						"data_disk.1.encrypted":            "false",
						"data_disk.1.kms_key_id":           "",
						"data_disk.1.name":                 name,
						"data_disk.1.description":          "",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_CustomPrioritiesCreate(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun",
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
					"image_id":          "${data.alicloud_images.default2.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"custom_priorities": []map[string]string{{
						"instance_type": "${data.alicloud_instance_types.c6.instance_types.0.id}",
						"vswitch_id":    "${alicloud_vswitch.default.id}",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit":    "false",
						"custom_priorities.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "kms_encryption_context"},
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
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_CustomPrioritiesModify(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun",
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
					"image_id":          "${data.alicloud_images.default2.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
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
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "kms_encryption_context"},
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
					"custom_priorities": []map[string]string{{
						"instance_type": "${data.alicloud_instance_types.c6.instance_types.0.id}",
						"vswitch_id":    "${alicloud_vswitch.default.id}",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_priorities.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_EnchanceCreate(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun",
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
					"scaling_group_id":                "${alicloud_ess_scaling_group.default.id}",
					"image_id":                        "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":                   "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id":               "${alicloud_security_group.default.id}",
					"security_enhancement_strategy":   "Deactive",
					"force_delete":                    "true",
					"instance_description":            "test",
					"spot_duration":                   "0",
					"spot_strategy":                   "SpotWithPriceLimit",
					"image_options_login_as_non_root": "true",
					"system_disk_encrypt_algorithm":   "AES-256",
					"system_disk_provisioned_iops":    "10",
					"system_disk_encrypted":           "true",
					"system_disk_kms_key_id":          "${alicloud_kms_key.key.id}",
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"provisioned_iops":     "20",
							"category":             "cloud_ssd",
							"delete_with_instance": "false",
							"encrypted":            "true",
							"kms_key_id":           "${alicloud_kms_key.key.id}",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":             "20",
							"provisioned_iops": "10",
							"category":         "cloud_ssd",
							"name":             "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit":                 "false",
						"security_enhancement_strategy":    "Deactive",
						"instance_description":             "test",
						"spot_duration":                    "0",
						"image_options_login_as_non_root":  "true",
						"system_disk_encrypt_algorithm":    "AES-256",
						"system_disk_provisioned_iops":     "10",
						"system_disk_encrypted":            "true",
						"system_disk_kms_key_id":           CHECKSET,
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.provisioned_iops":     "20",
						"data_disk.0.category":             "cloud_ssd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.kms_key_id":           CHECKSET,
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                 "20",
						"data_disk.1.category":             "cloud_ssd",
						"data_disk.1.provisioned_iops":     "10",
						"data_disk.1.delete_with_instance": "true",
						"data_disk.1.encrypted":            "false",
						"data_disk.1.kms_key_id":           "",
						"data_disk.1.name":                 name,
						"data_disk.1.description":          "",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_DedicatedHostClusterId_HttpEndpoint_Create(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun.*vhd",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
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
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"image_id":                  "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"dedicated_host_cluster_id": "${alicloud_ecs_dedicated_host_cluster.default.id}",
					"instance_type":             "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id":         "${alicloud_security_group.default1.id}",
					"force_delete":              "true",
					"system_disk_category":      "cloud_ssd",
					"http_endpoint":             "enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_cluster_id": CHECKSET,
						"password_inherit":          "false",
						"http_endpoint":             "enabled",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id"},
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
					"security_group_id":         REMOVEKEY,
					"dedicated_host_cluster_id": "${alicloud_ecs_dedicated_host_cluster.default1.id}",
					"http_endpoint":             "disabled",
					"network_interfaces": []map[string]interface{}{
						{
							"instance_type":                  "Primary",
							"ipv6_address_count":             "1",
							"network_interface_traffic_mode": "Standard",
							"security_group_ids":             []string{"${alicloud_security_group.default1.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":         "",
						"network_interfaces.#":      "1",
						"dedicated_host_cluster_id": CHECKSET,
						"http_endpoint":             "disabled",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": []map[string]interface{}{
						{
							"instance_type":                  "Primary",
							"ipv6_address_count":             "1",
							"network_interface_traffic_mode": "Standard",
							"security_group_ids":             []string{"${alicloud_security_group.default1.id}"},
						},
						{
							"instance_type":                  "Secondary",
							"network_interface_traffic_mode": "HighPerformance",
							"security_group_ids":             []string{"${alicloud_security_group.default1.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interfaces.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_DedicatedHostClusterId_HttpEndpoint(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun.*vhd",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
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
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":        "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default1.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_ssd",
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
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id"},
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
					"security_group_id":         REMOVEKEY,
					"dedicated_host_cluster_id": "${alicloud_ecs_dedicated_host_cluster.default1.id}",
					"http_endpoint":             "disabled",
					"network_interfaces": []map[string]interface{}{
						{
							"instance_type":                  "Primary",
							"ipv6_address_count":             "1",
							"network_interface_traffic_mode": "Standard",
							"security_group_ids":             []string{"${alicloud_security_group.default1.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":         "",
						"network_interfaces.#":      "1",
						"dedicated_host_cluster_id": CHECKSET,
						"http_endpoint":             "disabled",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"http_endpoint":             "enabled",
					"dedicated_host_cluster_id": "${alicloud_ecs_dedicated_host_cluster.default.id}",
					"network_interfaces": []map[string]interface{}{
						{
							"instance_type":                  "Primary",
							"ipv6_address_count":             "1",
							"network_interface_traffic_mode": "Standard",
							"security_group_ids":             []string{"${alicloud_security_group.default1.id}"},
						},
						{
							"instance_type":                  "Secondary",
							"network_interface_traffic_mode": "HighPerformance",
							"security_group_ids":             []string{"${alicloud_security_group.default1.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_endpoint":             "enabled",
						"dedicated_host_cluster_id": CHECKSET,
						"network_interfaces.#":      "2",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_Enchance(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun",
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
					"image_id":          "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"spot_duration":     "0",
					"spot_strategy":     "SpotWithPriceLimit",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "false",
						"spot_duration":    "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_encrypted":  "true",
					"system_disk_kms_key_id": "${alicloud_kms_key.key.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_encrypted":  "true",
						"system_disk_kms_key_id": CHECKSET,
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
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_duration": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_duration": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_provisioned_iops":    "20",
					"system_disk_encrypt_algorithm":   "SM4-128",
					"image_options_login_as_non_root": "false",
					"deletion_protection":             "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_provisioned_iops":    "20",
						"system_disk_encrypt_algorithm":   "SM4-128",
						"image_options_login_as_non_root": "false",
						"deletion_protection":             "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_description": "cloud_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_description": "cloud_ssd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_description": "testd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_description": "testd",
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
					"spot_price_limit": []map[string]string{{
						"instance_type": "${data.alicloud_instance_types.c6.instance_types.0.id}",
						"price_limit":   "2.1",
					},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spot_price_limit.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"provisioned_iops":     "20",
							"category":             "cloud_ssd",
							"delete_with_instance": "false",
							"encrypted":            "true",
							"kms_key_id":           "${alicloud_kms_key.key.id}",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":             "20",
							"provisioned_iops": "10",
							"category":         "cloud_ssd",
							"name":             "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.provisioned_iops":     "20",
						"data_disk.0.category":             "cloud_ssd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.kms_key_id":           CHECKSET,
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                 "20",
						"data_disk.1.category":             "cloud_ssd",
						"data_disk.1.provisioned_iops":     "10",
						"data_disk.1.delete_with_instance": "true",
						"data_disk.1.encrypted":            "false",
						"data_disk.1.kms_key_id":           "",
						"data_disk.1.name":                 name,
						"data_disk.1.description":          "",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_Multi(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default.9"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigMutilDependence)
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
					"image_id":          "${data.alicloud_images.default1.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
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

func resourceEssScalingConfigurationConfigDependences(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}
	resource "alicloud_key_pair" "default" {
  		key_pair_name = "${var.name}"
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	  name_regex = "default"
	}

	data "alicloud_kms_keys" "default" {
		  status = "Enabled"
		}
	resource "alicloud_kms_key" "default" {
	  description = var.name
	  status = "Enabled"
	  pending_window_in_days = 7
	}
	
	resource "alicloud_kms_ciphertext" "default" {
	  key_id = "${alicloud_kms_key.default.id}"
	  plaintext = "YourPassword1234"
	  encryption_context = {
		"name" = var.name
	  }
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
		name_regex  = "^ubu"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_images" "default2" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "t5" {
      instance_type_family = "ecs.t5"
	}
    data "alicloud_instance_types" "t6" {
      instance_type_family = "ecs.t6"
	}
    data "alicloud_instance_types" "c6" {
      
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	 data "alicloud_instance_types" "default12" {
      instance_type_family = "ecs.n4"
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

func resourceEssScalingConfigurationConfigDependences_beijing(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}
	resource "alicloud_key_pair" "default" {
  		key_pair_name = "${var.name}"
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	  name_regex = "default"
	}

	data "alicloud_kms_keys" "default" {
		  status = "Enabled"
		}
	resource "alicloud_kms_key" "default" {
	  description = var.name
	  status = "Enabled"
	  pending_window_in_days = 7
	}
	
	resource "alicloud_kms_ciphertext" "default" {
	  key_id = "${alicloud_kms_key.default.id}"
	  plaintext = "YourPassword1234"
	  encryption_context = {
		"name" = var.name
	  }
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
		name_regex  = "^ubu"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_images" "default2" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
    resource "alicloud_vswitch" "default1" {
  		vpc_id            = "${alicloud_vpc.default.id}"
  		cidr_block        = "172.16.1.0/24"
  		zone_id = "cn-beijing-f"
        vswitch_name              = "${var.name}"
 	}
	data "alicloud_instance_types" "t5" {
      instance_type_family = "ecs.t5"
	}
    data "alicloud_instance_types" "t6" {
      instance_type_family = "ecs.t5"
	}
    data "alicloud_instance_types" "c6" {
      
	  availability_zone = "cn-beijing-f"
	}
	 data "alicloud_instance_types" "default12" {
      instance_type_family = "ecs.n4"
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
		vswitch_ids = ["${alicloud_vswitch.default1.id}"]
	}`, EcsInstanceCommonTestCase, name)
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
	data "alicloud_resource_manager_resource_groups" "default" {
	  name_regex = "default"
	}

	data "alicloud_kms_keys" "default" {
		  status = "Enabled"
		}
	resource "alicloud_kms_key" "default" {
	  description = var.name
	  status = "Enabled"
	  pending_window_in_days = 7
	}
	
	resource "alicloud_kms_ciphertext" "default" {
	  key_id = "${alicloud_kms_key.default.id}"
	  plaintext = "YourPassword1234"
	  encryption_context = {
		"name" = var.name
	  }
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
	data "alicloud_images" "default2" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "t5" {
      instance_type_family = "ecs.t5"
	}
    data "alicloud_instance_types" "c6" {
      
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	 data "alicloud_instance_types" "default12" {
      instance_type_family = "ecs.n4"
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

func resourceEssScalingConfiguration_NetworkInterfaces(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}
    resource "alicloud_vpc" "vpc" {
      cidr_block = "192.168.0.0/16"
      vpc_name   = var.name
      ipv6_isp    = "BGP"
      enable_ipv6  = "true"
	}
	resource "alicloud_ecs_dedicated_host_cluster" "default" {
  		dedicated_host_cluster_name = var.name
  		description                 = var.name
  		zone_id                     = data.alicloud_zones.default.zones.0.id
  		tags                        = {
    		Create = "TF"
    		For    = "DDH_Cluster_Test",
  		}
	}
	resource "alicloud_ecs_dedicated_host_cluster" "default1" {
  		dedicated_host_cluster_name = var.name
  		description                 = var.name
  		zone_id                     = data.alicloud_zones.default.zones.0.id
  		tags                        = {
    		Create = "TF1"
    		For    = "DDH_Cluster_Test1",
  		}
	}
	resource "alicloud_vswitch" "vswtich" {
      vpc_id       = "${alicloud_vpc.vpc.id}"
      zone_id      = "${data.alicloud_zones.default.zones.0.id}"
      vswitch_name = var.name
      cidr_block   = "192.168.10.0/24"
	  ipv6_cidr_block_mask = "8"
    } 
	resource "alicloud_security_group" "default1" {
	  security_group_name   = "${var.name}"
	  vpc_id = "${alicloud_vpc.vpc.id}"
	}
	data "alicloud_images" "default1" {
		name_regex  = "^aliyun.*vhd"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_instance_types" "c6" {
        instance_type_family = "ecs.sn1ne"
    }
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.vswtich.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingConfigurationConfigMutilDependence(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}

	resource "alicloud_security_group" "default1" {
	  name   = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	}

	data "alicloud_images" "default1" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_instance_types" "c6" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

func TestAccAliCloudEssScalingConfiguration_InstancePatternInfo(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun",
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
					"image_id":          "${data.alicloud_images.default2.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
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
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "kms_encryption_context"},
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
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"category":             "cloud_ssd",
							"delete_with_instance": "false",
							"encrypted":            "false",
							"kms_key_id":           "${alicloud_kms_key.key.id}",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":     "20",
							"category": "cloud_ssd",
							"name":     "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_ssd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "false",
						"data_disk.0.kms_key_id":           CHECKSET,
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                 "20",
						"data_disk.1.category":             "cloud_ssd",
						"data_disk.1.delete_with_instance": "true",
						"data_disk.1.encrypted":            "false",
						"data_disk.1.kms_key_id":           "",
						"data_disk.1.name":                 name,
						"data_disk.1.description":          "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":   "EntryLevel",
							"cores":                   "4",
							"memory":                  "4.00",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":   "EntryLevel",
							"cores":                   "4",
							"memory":                  "4.0",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
						},
						{
							"instance_family_level":   "EntryLevel",
							"cores":                   "2",
							"memory":                  "8.0",
							"max_price":               "2.1",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":                   "EntryLevel",
							"minimum_cpu_core_count":                  "4",
							"maximum_cpu_core_count":                  "4",
							"memory":                                  "4.0",
							"burstable_performance":                   "Include",
							"architectures":                           []string{"X86"},
							"excluded_instance_types":                 []string{"ecs.c6.large"},
							"minimum_eni_quantity":                    "1",
							"minimum_eni_private_ip_address_quantity": "1",
							"minimum_eni_ipv6_address_quantity":       "1",
							"minimum_baseline_credit":                 "1",
							"minimum_gpu_amount":                      "1",
							"maximum_gpu_amount":                      "1",
						},
						{
							"instance_family_level":   "EntryLevel",
							"minimum_cpu_core_count":  "2",
							"maximum_cpu_core_count":  "2",
							"memory":                  "8.0",
							"max_price":               "2.1",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":                   "EntryLevel",
							"minimum_cpu_core_count":                  "2",
							"maximum_cpu_core_count":                  "4",
							"memory":                                  "4.0",
							"burstable_performance":                   "Include",
							"architectures":                           []string{"X86"},
							"excluded_instance_types":                 []string{"ecs.c6.large"},
							"minimum_eni_quantity":                    "2",
							"minimum_eni_private_ip_address_quantity": "2",
							"minimum_eni_ipv6_address_quantity":       "2",
							"minimum_baseline_credit":                 "2",
							"minimum_gpu_amount":                      "2",
							"maximum_gpu_amount":                      "2",
						},
						{
							"instance_family_level":   "EntryLevel",
							"minimum_cpu_core_count":  "4",
							"maximum_cpu_core_count":  "8",
							"memory":                  "8.0",
							"max_price":               "2.1",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_InstancePatternInfo_InstanceTypeFamiliesAndInstanceCategories(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^ubu",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependences)
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
					"image_id":          "${data.alicloud_images.default1.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
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
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "kms_encryption_context"},
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
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":     "EntryLevel",
							"cores":                     "4",
							"memory":                    "4.00",
							"burstable_performance":     "Include",
							"instance_type_families":    []string{"ecs.c6"},
							"instance_categories":       []string{"General-purpose"},
							"physical_processor_models": []string{"Intel Xeon(Ice Lake) Platinum 8369B"},
							"cpu_architectures":         []string{"X86"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":     "EntryLevel",
							"cores":                     "4",
							"memory":                    "4.0",
							"burstable_performance":     "Include",
							"instance_type_families":    []string{"ecs.g6"},
							"instance_categories":       []string{"General-purpose", "Enhanced"},
							"physical_processor_models": []string{"Intel Xeon (Skylake) Platinum 8163"},
							"cpu_architectures":         []string{"X86", "ARM"},
						},
						{
							"instance_family_level":     "EntryLevel",
							"cores":                     "2",
							"memory":                    "8.0",
							"max_price":                 "2.1",
							"burstable_performance":     "Include",
							"instance_type_families":    []string{"ecs.c6", "ecs.g6"},
							"instance_categories":       []string{"General-purpose", "Compute-optimized"},
							"physical_processor_models": []string{"Intel Xeon (Skylake) Platinum 8163", "Intel Xeon(Ice Lake) Platinum 8369B"},
							"cpu_architectures":         []string{"ARM"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":                   "EntryLevel",
							"minimum_cpu_core_count":                  "4",
							"maximum_cpu_core_count":                  "4",
							"memory":                                  "4.0",
							"architectures":                           []string{"X86"},
							"excluded_instance_types":                 []string{"ecs.c6.large"},
							"minimum_eni_quantity":                    "1",
							"minimum_eni_private_ip_address_quantity": "1",
							"minimum_eni_ipv6_address_quantity":       "1",
							"minimum_baseline_credit":                 "1",
							"minimum_initial_credit":                  "1",
							"minimum_gpu_amount":                      "1",
							"maximum_gpu_amount":                      "1",
						},
						{
							"instance_family_level":   "EntryLevel",
							"minimum_cpu_core_count":  "2",
							"maximum_cpu_core_count":  "2",
							"memory":                  "8.0",
							"max_price":               "2.1",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":                   "EntryLevel",
							"minimum_cpu_core_count":                  "2",
							"maximum_cpu_core_count":                  "4",
							"memory":                                  "4.0",
							"architectures":                           []string{"X86"},
							"excluded_instance_types":                 []string{"ecs.c6.large"},
							"minimum_eni_quantity":                    "2",
							"minimum_eni_private_ip_address_quantity": "2",
							"minimum_eni_ipv6_address_quantity":       "2",
							"minimum_baseline_credit":                 "2",
							"minimum_initial_credit":                  "2",
							"minimum_gpu_amount":                      "2",
							"maximum_gpu_amount":                      "2",
						},
						{
							"instance_family_level":   "EntryLevel",
							"minimum_cpu_core_count":  "4",
							"maximum_cpu_core_count":  "8",
							"memory":                  "8.0",
							"max_price":               "2.1",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_InstancePatternInfo_CpuAndRange(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^ubu",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependences)
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
					"image_id":          "${data.alicloud_images.default1.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
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
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "kms_encryption_context"},
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
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"category":             "cloud_ssd",
							"delete_with_instance": "false",
							"encrypted":            "false",
							"kms_key_id":           "${alicloud_kms_key.key.id}",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":     "20",
							"category": "cloud_ssd",
							"name":     "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_ssd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "false",
						"data_disk.0.kms_key_id":           CHECKSET,
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                 "20",
						"data_disk.1.category":             "cloud_ssd",
						"data_disk.1.delete_with_instance": "true",
						"data_disk.1.encrypted":            "false",
						"data_disk.1.kms_key_id":           "",
						"data_disk.1.name":                 name,
						"data_disk.1.description":          "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":   "EntryLevel",
							"minimum_cpu_core_count":  "2",
							"maximum_cpu_core_count":  "4",
							"minimum_memory_size":     "4",
							"maximum_memory_size":     "8",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
							"gpu_specs":               []string{"NVIDIA T4", "NVIDIA V100"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":   "EntryLevel",
							"minimum_cpu_core_count":  "4",
							"maximum_cpu_core_count":  "4",
							"minimum_memory_size":     "8",
							"maximum_memory_size":     "8",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
							"gpu_specs":               []string{"NVIDIA V100"},
						},
						{
							"instance_family_level":   "EntryLevel",
							"minimum_cpu_core_count":  "4",
							"maximum_cpu_core_count":  "4",
							"minimum_memory_size":     "8",
							"maximum_memory_size":     "8",
							"max_price":               "2.1",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
							"gpu_specs":               []string{"NVIDIA T4"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":                   "EntryLevel",
							"minimum_cpu_core_count":                  "4",
							"maximum_cpu_core_count":                  "4",
							"memory":                                  "4.0",
							"burstable_performance":                   "Include",
							"architectures":                           []string{"X86"},
							"excluded_instance_types":                 []string{"ecs.c6.large"},
							"minimum_eni_quantity":                    "1",
							"minimum_eni_private_ip_address_quantity": "1",
							"minimum_eni_ipv6_address_quantity":       "1",
							"minimum_baseline_credit":                 "1",
							"minimum_gpu_amount":                      "1",
							"maximum_gpu_amount":                      "1",
						},
						{
							"instance_family_level":   "EntryLevel",
							"minimum_cpu_core_count":  "2",
							"maximum_cpu_core_count":  "2",
							"memory":                  "8.0",
							"max_price":               "2.1",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":                   "EntryLevel",
							"minimum_cpu_core_count":                  "2",
							"maximum_cpu_core_count":                  "4",
							"memory":                                  "4.0",
							"burstable_performance":                   "Include",
							"architectures":                           []string{"X86"},
							"excluded_instance_types":                 []string{"ecs.c6.large"},
							"minimum_eni_quantity":                    "2",
							"minimum_eni_private_ip_address_quantity": "2",
							"minimum_eni_ipv6_address_quantity":       "2",
							"minimum_baseline_credit":                 "2",
							"minimum_gpu_amount":                      "2",
							"maximum_gpu_amount":                      "2",
						},
						{
							"instance_family_level":   "EntryLevel",
							"minimum_cpu_core_count":  "4",
							"maximum_cpu_core_count":  "8",
							"memory":                  "8.0",
							"max_price":               "2.1",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_pattern_info.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_NetworkInterfaces(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun.*vhd",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
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
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":        "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default1.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_ssd",
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
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id"},
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
					"security_group_id": REMOVEKEY,
					"network_interfaces": []map[string]interface{}{
						{
							"instance_type":                  "Primary",
							"ipv6_address_count":             "1",
							"network_interface_traffic_mode": "Standard",
							"security_group_ids":             []string{"${alicloud_security_group.default1.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    "",
						"network_interfaces.#": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": []map[string]interface{}{
						{
							"instance_type":                  "Primary",
							"ipv6_address_count":             "1",
							"network_interface_traffic_mode": "Standard",
							"security_group_ids":             []string{"${alicloud_security_group.default1.id}"},
						},
						{
							"instance_type":                  "Secondary",
							"network_interface_traffic_mode": "HighPerformance",
							"security_group_ids":             []string{"${alicloud_security_group.default1.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interfaces.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_InstanceTypeOverride(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          REGEXMATCH + "^aliyun",
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
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"system_disk_category": "cloud_efficiency",
					"image_id":             "${data.alicloud_images.default2.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.c6.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"force_delete":         "true",
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
				ImportStateVerifyIgnore: []string{"force_delete", "instance_type", "security_group_id", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_essd",
					"internet_charge_type": "PayByTraffic",
					"instance_name":        name,
					"override":             "true",
					"instance_type":        REMOVEKEY,
					"instance_type_override": []map[string]string{{
						"instance_type":     "ecs.c6.large",
						"weighted_capacity": "3",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category":     "cloud_essd",
						"internet_charge_type":     "PayByTraffic",
						"instance_name":            name,
						"instance_type":            REMOVEKEY,
						"instance_type_override.#": "1",
						"override":                 "true",
					}),
				),
			},
		},
	})
}
