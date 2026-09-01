package alicloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// collection-order: bypass instance_ids -- pre-existing deprecated TypeList left untouched by this change (deprecated since provider 1.6.0); its suppressor is stock debt and intentionally returns true for every change, including member additions and removals, because CreateScalingConfiguration, ModifyScalingConfiguration and DescribeScalingConfigurations do not carry instance ids at all. A reorder apply cannot surface a diff, and a bounded suppressor unit net cannot have the required no-suppression membership case without changing that legacy behaviour.

// kms限制
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_name":           "${data.alicloud_images.default3.images.0.name}",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_essd",
					"password_inherit":     "true",
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

// kms限制
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_name":           "${data.alicloud_images.default3.images.0.name}",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd",
					"force_delete":         "true",
					"password_inherit":     "true",
					"active":               "true",
					"is_outdated":          "false",
					"substitute":           "false",
					"enable":               "true",
					"instance_ids":         []string{},
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_name":        "${data.alicloud_images.default4.images.0.name}",
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependenceImageName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default3.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"system_disk_category": "cloud_essd",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"force_delete":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id": "${data.alicloud_images.default4.images.0.id}",
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependenceImageName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default3.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd_xc1",
					"force_delete":         "true",
					"password_inherit":     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":   "${data.alicloud_images.default3.images.0.id}",
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

// kms限制
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_name":        "${data.alicloud_images.default3.images.0.name}",
					"instance_type":     "${data.alicloud_instance_types.default.instance_types.0.id}",
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

// kms限制
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"image_name":                "${data.alicloud_images.default3.images.0.name}",
					"instance_type":             "${data.alicloud_instance_types.default.instance_types.0.id}",
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
					"system_disk_size": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "60",
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

func TestAccAliCloudEssScalingConfiguration_UpdateSystemDiskSize(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  CHECKSET,
		"instance_type":     CHECKSET,
		"security_group_id": CHECKSET,
		"image_id":          CHECKSET,
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependenceUpdate)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"image_id":                  "${data.alicloud_images.default3.images.0.id}",
					"instance_type":             "${data.alicloud_instance_types.default.instance_types.0.id}",
					"security_group_id":         "${alicloud_security_group.default.id}",
					"force_delete":              "true",
					"internet_max_bandwidth_in": "1",
					"system_disk_size":          "60",
					"system_disk_category":      "cloud_ssd",
					//"password":          "123-abcABC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit":          "false",
						"internet_max_bandwidth_in": "1",
						"system_disk_size":          "60",
						"system_disk_category":      "cloud_ssd",
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
					"system_disk_size":           REMOVEKEY,
					"scaling_group_id":           "${alicloud_ess_scaling_group.default.id}",
					"image_id":                   "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":              "${data.alicloud_instance_types.n4.instance_types.1.id}",
					"security_group_id":          "${alicloud_security_group.default.id}",
					"system_disk_category":       "cloud_essd_xc1",
					"internet_max_bandwidth_out": "512",
					"system_disk_encrypted":      "false",
					"instance_name":              "test",
					"scaling_configuration_name": fmt.Sprintf("tf-testAccEssScalingConfigurati-%d", rand),
					"override":                   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size":           REMOVEKEY,
						"system_disk_category":       "cloud_essd_xc1",
						"scaling_configuration_name": fmt.Sprintf("tf-testAccEssScalingConfigurati-%d", rand),
						"internet_max_bandwidth_out": "512",
						"system_disk_encrypted":      "false",
						"instance_name":              "test",
						"override":                   "true",
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependencePerformanceLevel)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default3.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_essd",
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
					"system_disk_category": "cloud_essd_xc1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_essd_xc1",
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
						"instance_type": "${data.alicloud_instance_types.n4.instance_types.0.id}",
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
							"category":             "cloud_essd",
							"delete_with_instance": "false",
							"encrypted":            "true",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":     "20",
							"category": "cloud_essd",
							"name":     "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_essd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                 "20",
						"data_disk.1.category":             "cloud_essd",
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependencePerformanceLevel)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default3.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd",
					"force_delete":         "true",
					"custom_priorities": []map[string]string{{
						"instance_type": "${data.alicloud_instance_types.n4.instance_types.0.id}",
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
					"system_disk_category": "cloud_essd_xc1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_essd_xc1",
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependencePerformanceLevel)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default3.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd",
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
					"system_disk_category": "cloud_essd_xc1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category": "cloud_essd_xc1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_priorities": []map[string]string{{
						"instance_type": "${data.alicloud_instance_types.n4.instance_types.0.id}",
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

// kms限制
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                "${alicloud_ess_scaling_group.default.id}",
					"image_id":                        "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":                   "${data.alicloud_instance_types.n4.instance_types.0.id}",
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
					"system_disk_category":            "cloud_essd",
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"provisioned_iops":     "20",
							"category":             "cloud_essd",
							"delete_with_instance": "false",
							"encrypted":            "true",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":             "20",
							"provisioned_iops": "10",
							"category":         "cloud_essd",
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
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.provisioned_iops":     "20",
						"data_disk.0.category":             "cloud_essd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                 "20",
						"data_disk.1.category":             "cloud_essd",
						"data_disk.1.provisioned_iops":     "10",
						"data_disk.1.delete_with_instance": "true",
						"data_disk.1.encrypted":            "false",
						"data_disk.1.name":                 name,
						"data_disk.1.description":          "",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_ResourcePoolOptions_PoolTags(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependenceForResourcePoolOptions)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                "${alicloud_ess_scaling_group.default.id}",
					"image_id":                        "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":                   "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":               "${alicloud_security_group.default.id}",
					"security_enhancement_strategy":   "Deactive",
					"force_delete":                    "true",
					"instance_description":            "test",
					"spot_duration":                   "0",
					"spot_strategy":                   "SpotWithPriceLimit",
					"image_options_login_as_non_root": "true",
					"system_disk_encrypt_algorithm":   "AES-256",
					"system_disk_category":            "cloud_essd",
					"system_disk_provisioned_iops":    "10",
					"system_disk_encrypted":           "true",
					"resource_pool_options_strategy":  "PrivatePoolFirst",
					"resource_pool_options_private_pool_tags": []map[string]string{
						{
							"key":   "TagKey",
							"value": "TagValue",
						},
					},
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"provisioned_iops":     "20",
							"category":             "cloud_essd",
							"delete_with_instance": "false",
							"encrypted":            "true",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":             "20",
							"provisioned_iops": "10",
							"category":         "cloud_essd",
							"name":             "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit":                 "false",
						"security_enhancement_strategy":    "Deactive",
						"instance_description":             "test",
						"system_disk_category":             "cloud_essd",
						"spot_duration":                    "0",
						"image_options_login_as_non_root":  "true",
						"system_disk_encrypt_algorithm":    "AES-256",
						"system_disk_provisioned_iops":     "10",
						"system_disk_encrypted":            "true",
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.provisioned_iops":     "20",
						"data_disk.0.category":             "cloud_essd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                          "20",
						"data_disk.1.category":                      "cloud_essd",
						"data_disk.1.provisioned_iops":              "10",
						"data_disk.1.delete_with_instance":          "true",
						"data_disk.1.encrypted":                     "false",
						"data_disk.1.name":                          name,
						"data_disk.1.description":                   "",
						"resource_pool_options_strategy":            "PrivatePoolFirst",
						"resource_pool_options_private_pool_tags.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                "${alicloud_ess_scaling_group.default.id}",
					"image_id":                        "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":                   "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":               "${alicloud_security_group.default.id}",
					"security_enhancement_strategy":   "Deactive",
					"force_delete":                    "true",
					"instance_description":            "test",
					"spot_duration":                   "0",
					"spot_strategy":                   "SpotWithPriceLimit",
					"image_options_login_as_non_root": "true",
					"system_disk_encrypt_algorithm":   "AES-256",
					"system_disk_category":            "cloud_essd",
					"system_disk_provisioned_iops":    "10",
					"system_disk_encrypted":           "true",
					"resource_pool_options_strategy":  "PrivatePoolOnly",
					"resource_pool_options_private_pool_tags": []map[string]string{
						{
							"key":   "TagKey",
							"value": "TagValue",
						},
						{
							"key":   "TagKey1",
							"value": "TagValue2",
						},
					},
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"provisioned_iops":     "20",
							"category":             "cloud_essd",
							"delete_with_instance": "false",
							"encrypted":            "true",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":             "20",
							"provisioned_iops": "10",
							"category":         "cloud_essd",
							"name":             "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit":                 "false",
						"security_enhancement_strategy":    "Deactive",
						"instance_description":             "test",
						"system_disk_category":             "cloud_essd",
						"spot_duration":                    "0",
						"image_options_login_as_non_root":  "true",
						"system_disk_encrypt_algorithm":    "AES-256",
						"system_disk_provisioned_iops":     "10",
						"system_disk_encrypted":            "true",
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.provisioned_iops":     "20",
						"data_disk.0.category":             "cloud_essd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                          "20",
						"data_disk.1.category":                      "cloud_essd",
						"data_disk.1.provisioned_iops":              "10",
						"data_disk.1.delete_with_instance":          "true",
						"data_disk.1.encrypted":                     "false",
						"data_disk.1.name":                          name,
						"data_disk.1.description":                   "",
						"resource_pool_options_strategy":            "PrivatePoolOnly",
						"resource_pool_options_private_pool_tags.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_ResourcePoolOptions_PoolIds(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependenceForResourcePoolOptions)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                       "${alicloud_ess_scaling_group.default.id}",
					"image_id":                               "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":                          "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":                      "${alicloud_security_group.default.id}",
					"security_enhancement_strategy":          "Deactive",
					"force_delete":                           "true",
					"instance_description":                   "test",
					"spot_duration":                          "0",
					"spot_strategy":                          "SpotWithPriceLimit",
					"image_options_login_as_non_root":        "true",
					"system_disk_encrypt_algorithm":          "AES-256",
					"system_disk_category":                   "cloud_essd",
					"system_disk_provisioned_iops":           "10",
					"system_disk_encrypted":                  "true",
					"resource_pool_options_strategy":         "PrivatePoolFirst",
					"resource_pool_options_private_pool_ids": []string{"${alicloud_ecs_capacity_reservation.default.0.id}"},
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"provisioned_iops":     "20",
							"category":             "cloud_essd",
							"delete_with_instance": "false",
							"encrypted":            "true",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":             "20",
							"provisioned_iops": "10",
							"category":         "cloud_essd",
							"name":             "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit":                 "false",
						"security_enhancement_strategy":    "Deactive",
						"instance_description":             "test",
						"system_disk_category":             "cloud_essd",
						"spot_duration":                    "0",
						"image_options_login_as_non_root":  "true",
						"system_disk_encrypt_algorithm":    "AES-256",
						"system_disk_provisioned_iops":     "10",
						"system_disk_encrypted":            "true",
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.provisioned_iops":     "20",
						"data_disk.0.category":             "cloud_essd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                         "20",
						"data_disk.1.category":                     "cloud_essd",
						"data_disk.1.provisioned_iops":             "10",
						"data_disk.1.delete_with_instance":         "true",
						"data_disk.1.encrypted":                    "false",
						"data_disk.1.kms_key_id":                   "",
						"data_disk.1.name":                         name,
						"data_disk.1.description":                  "",
						"resource_pool_options_strategy":           "PrivatePoolFirst",
						"resource_pool_options_private_pool_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                       "${alicloud_ess_scaling_group.default.id}",
					"image_id":                               "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":                          "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":                      "${alicloud_security_group.default.id}",
					"security_enhancement_strategy":          "Deactive",
					"force_delete":                           "true",
					"instance_description":                   "test",
					"spot_duration":                          "0",
					"spot_strategy":                          "SpotWithPriceLimit",
					"image_options_login_as_non_root":        "true",
					"system_disk_encrypt_algorithm":          "AES-256",
					"system_disk_category":                   "cloud_essd",
					"system_disk_provisioned_iops":           "10",
					"system_disk_encrypted":                  "true",
					"resource_pool_options_strategy":         "PrivatePoolOnly",
					"resource_pool_options_private_pool_ids": []string{"${alicloud_ecs_capacity_reservation.default.0.id}", "${alicloud_ecs_capacity_reservation.default.1.id}"},
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"provisioned_iops":     "20",
							"category":             "cloud_essd",
							"delete_with_instance": "false",
							"encrypted":            "true",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":             "20",
							"provisioned_iops": "10",
							"category":         "cloud_essd",
							"name":             "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_inherit":                 "false",
						"security_enhancement_strategy":    "Deactive",
						"instance_description":             "test",
						"system_disk_category":             "cloud_essd",
						"spot_duration":                    "0",
						"image_options_login_as_non_root":  "true",
						"system_disk_encrypt_algorithm":    "AES-256",
						"system_disk_provisioned_iops":     "10",
						"system_disk_encrypted":            "true",
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.provisioned_iops":     "20",
						"data_disk.0.category":             "cloud_essd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                         "20",
						"data_disk.1.category":                     "cloud_essd",
						"data_disk.1.provisioned_iops":             "10",
						"data_disk.1.delete_with_instance":         "true",
						"data_disk.1.encrypted":                    "false",
						"data_disk.1.kms_key_id":                   "",
						"data_disk.1.name":                         name,
						"data_disk.1.description":                  "",
						"resource_pool_options_strategy":           "PrivatePoolOnly",
						"resource_pool_options_private_pool_ids.#": "2",
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.default.id}",
					"image_id":                  "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"dedicated_host_cluster_id": "${alicloud_ecs_dedicated_host_cluster.default.id}",
					"instance_type":             "ecs.c9i.large",
					"security_group_id":         "${alicloud_security_group.default1.id}",
					"force_delete":              "true",
					"system_disk_category":      "cloud_essd",
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":        "ecs.c9i.large",
					"security_group_id":    "${alicloud_security_group.default1.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_essd",
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

// kms限制
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd_xc1",
					"force_delete":         "true",
					"spot_duration":        "0",
					"spot_strategy":        "SpotWithPriceLimit",
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
					"instance_description": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_description": "cloud_essd",
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
						"instance_type": "${data.alicloud_instance_types.n4.instance_types.0.id}",
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
							"category":             "cloud_essd",
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
							"category":         "cloud_essd",
							"name":             "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.provisioned_iops":     "20",
						"data_disk.0.category":             "cloud_essd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "true",
						"data_disk.0.kms_key_id":           CHECKSET,
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                 "20",
						"data_disk.1.category":             "cloud_essd",
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":                "10",
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default3.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd",
					"force_delete":         "true",
					"password":             "123-abcABC",
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
    data "alicloud_instance_types" "default12" {
      instance_type_family = "ecs.n4"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_kms_keys" "default" {
		  status = "Enabled"
		}
	resource "alicloud_kms_key" "default" {
	  description = var.name
	  status = "Enabled"
	  pending_window_in_days = 7
	}
    data "alicloud_instance_types" "n4" {
      instance_type_family = "ecs.n4"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
     data "alicloud_images" "default4" {
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_images" "default2" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
  data "alicloud_images" "default5" {
        name_regex  = "^anolisos"
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

func resourceEssScalingConfigurationConfigDependencesCpuAndRange(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	  name_regex = "default"
	}
    data "alicloud_instance_types" "default12" {
      instance_type_family = "ecs.n4"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_kms_keys" "default" {
		  status = "Enabled"
		}
    data "alicloud_instance_types" "n4" {
      instance_type_family = "ecs.n4"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
     data "alicloud_images" "default4" {
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_images" "default2" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
  data "alicloud_images" "default5" {
        name_regex  = "^anolisos"
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
    data "alicloud_images" "default3" {
		name_regex  = "^win"
  		most_recent = true
  		owners      = "system"
	}
   data "alicloud_images" "default4" {
		name_regex  = "^debian"
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
    data "alicloud_instance_types" "n4" {
      instance_type_family = "ecs.g7h"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_images" "default3" {
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

func resourceEssScalingConfigurationConfigDependencePerformanceLevel(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	  name_regex = "default"
	}
    data "alicloud_instance_types" "n4" {
      instance_type_family = "ecs.g7h"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_kms_keys" "default" {
		  status = "Enabled"
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
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_images" "default3" {
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
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingConfigurationConfigDependenceImageName(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}
    data "alicloud_instance_types" "n4" {
      instance_type_family = "ecs.g7h"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_images" "default3" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_images" "default4" {
			name_regex  = "^anolisos"
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
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingConfigurationConfigDependenceTest(name string) string {
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
    data "alicloud_instance_types" "n4" {
      instance_type_family = "ecs.g7h"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_kms_keys" "default" {
		  status = "Enabled"
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
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_images" "default3" {
		name_regex  = "^aliyun"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_kms_ciphertext" "default" {
        key_id = "${data.alicloud_kms_keys.default.keys.0.id}"
        plaintext = "YourPassword1234"
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
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingConfigurationConfigDependenceUpdate(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	  name_regex = "default"
	}
    data "alicloud_instance_types" "n4" {
      instance_type_family = "ecs.g7h"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_kms_keys" "default" {
		  status = "Enabled"
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
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_images" "default3" {
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
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingConfigurationConfigDependenceForResourcePoolOptions(name string) string {
	return fmt.Sprintf(`
	%s
	
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	  name_regex = "default"
	}

	data "alicloud_kms_keys" "default" {
		  status = "Enabled"
		}

	data "alicloud_instance_types" "cloud_auto" {
	  availability_zone    = data.alicloud_zones.default.zones.0.id
	  kubernetes_node_role = "Worker"
	  system_disk_category = "cloud_auto"
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
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_instance_types" "t5" {
      instance_type_family = "ecs.t5"
	}
    data "alicloud_instance_types" "c6" {
      
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
     data "alicloud_instance_types" "n4" {
      instance_type_family = "ecs.g7h"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	 data "alicloud_instance_types" "default12" {
      instance_type_family = "ecs.n4"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

    resource "alicloud_ecs_capacity_reservation" "default" {
	  count = 2
	  description               = "terraform-example"
	  platform                  = "linux"
	  capacity_reservation_name = "terraform-example"
	  end_time_type             = "Unlimited"
	  instance_amount           = 1
	  instance_type             = data.alicloud_instance_types.n4.instance_types.0.id
	  match_criteria            = "Target"
	  zone_ids = ["${data.alicloud_zones.default.zones.0.id}"]
    }
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

func resourceEssScalingConfigurationConfigDependenceOnlyInstancePatternInfo(name string) string {
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
		name_regex  = "^centos.*_64"
  		most_recent = true
  		owners      = "system"
	}
	data "alicloud_images" "default2" {
		name_regex  = "^debian"
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
        instance_type_family = "ecs.c9i.lar"
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
    
   data "alicloud_instance_types" "n4" {
      instance_type_family = "ecs.g7h"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	resource "alicloud_security_group" "default1" {
	  name   = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	}
   
	data "alicloud_images" "default1" {
		name_regex  = "^debian"
  		most_recent = true
  		owners      = "system"
	}
    data "alicloud_images" "default3" {
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

// kms限制
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default3.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd_xc1",
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
					"data_disk": []map[string]string{
						{
							"size":                 "20",
							"category":             "cloud_essd",
							"delete_with_instance": "false",
							"encrypted":            "false",
							"kms_key_id":           "${alicloud_kms_key.key.id}",
							"name":                 "kms",
							"description":          "kms",
							"performance_level":    "PL1",
						},
						{
							"size":     "20",
							"category": "cloud_essd",
							"name":     "${var.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_disk.#":                      "2",
						"data_disk.0.size":                 "20",
						"data_disk.0.category":             "cloud_essd",
						"data_disk.0.delete_with_instance": "false",
						"data_disk.0.encrypted":            "false",
						"data_disk.0.kms_key_id":           CHECKSET,
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                 "20",
						"data_disk.1.category":             "cloud_essd",
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

func TestAccAliCloudEssScalingConfiguration_OnlyInstancePatternInfo(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":        CHECKSET,
		"security_group_id":       CHECKSET,
		"image_id":                REGEXMATCH + "^debian",
		"override":                "false",
		"instance_pattern_info.#": "1",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependenceOnlyInstancePatternInfo)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_id":          "${data.alicloud_images.default2.images.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":   "EntryLevel",
							"cores":                   "4",
							"memory":                  "8.00",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86"},
							"excluded_instance_types": []string{"ecs.c6.large"},
						},
					},
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
				ImportStateVerifyIgnore: []string{"force_delete", "security_group_id"},
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
							"instance_family_level":   "EntryLevel",
							"cores":                   "4",
							"memory":                  "8.00",
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
							"memory":                  "8.0",
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
							"memory":                                  "8.0",
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
							"memory":                                  "8.0",
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

// kms限制
func TestAccAliCloudEssScalingConfiguration_InstancePatternInfo_InstanceTypeFamiliesAndInstanceCategories(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand) //resourceEssScalingConfigurationConfigDependences
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default3.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.n4.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd",
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
		"image_id":          REGEXMATCH + "^anolisos",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependencesCpuAndRange)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_id":          "${data.alicloud_images.default5.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.n4.instance_types.0.id}",
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
						"data_disk.0.name":                 "kms",
						"data_disk.0.description":          "kms",
						"data_disk.0.performance_level":    "PL1",

						"data_disk.1.size":                 "20",
						"data_disk.1.category":             "cloud_ssd",
						"data_disk.1.delete_with_instance": "true",
						"data_disk.1.encrypted":            "false",
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

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "aliyun_3_9_x64_20G_alibase_20231219.vhd",
					"instance_type":        "ecs.c9i.large",
					"security_group_id":    "${alicloud_security_group.default1.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_essd",
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
		"image_id":          REGEXMATCH + "^anolisos",
		"override":          "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependencesCpuAndRange)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_id":          "${data.alicloud_images.default5.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.n4.instance_types.0.id}",
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
						"system_disk_category":     "cloud_ssd",
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

// TestAccAliCloudEssScalingConfiguration_userDataTagsConvergence covers the
// four regression scenarios of the perpetual-diff fix for issue #2079:
//  1. user_data/tags written out of band (for example by the ACK cluster
//     autoscaler) are preserved on refresh and the following plans stay empty.
//  2. an explicit configuration update of user_data/tags takes effect.
//  3. user_data/tags that are not configured are not cleared on the cloud side;
//     the remote values survive and no clear request is sent (contract R1).
//  4. removing user_data/tags from the configuration keeps the remote values
//     and converges instead of producing a perpetual diff.
func TestAccAliCloudEssScalingConfiguration_userDataTagsConvergence(t *testing.T) {
	var v ess.ScalingConfiguration
	var scalingConfigurationId string
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
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", acctest.RandIntRange(1000, 999999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationUserDataTagsDependence)

	externalUserData := base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\necho external-autoscaler"))
	externalUserDataPlain := "#!/bin/bash\necho external-autoscaler"
	updatedUserData := base64.StdEncoding.EncodeToString([]byte("#!/bin/bash\necho terraform-updated"))
	externalTags := map[string]string{"external": "ack", "managed-by": "autoscaler"}

	testAccConfig(map[string]interface{}{
		"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
		"image_name":           "${data.alicloud_images.default.images.0.name}",
		"instance_type":        "${data.alicloud_instance_types.compatible.instance_types.0.id}",
		"security_group_id":    "${alicloud_security_group.default.id}",
		"system_disk_category": "cloud_essd",
		"force_delete":         "true",
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"force_delete": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force_delete": "true",
					}),
					func(state *terraform.State) error {
						rs, ok := state.RootModule().Resources[resourceId]
						if !ok {
							return fmt.Errorf("resource %s not found in state", resourceId)
						}
						scalingConfigurationId = rs.Primary.ID
						// A configuration without instance tags either omits the
						// Tags key or returns an empty tag list. Both mean "no
						// remote tags", which is what Read has to turn into an
						// empty map instead of leaving the attribute unknown.
						response := essScalingConfigurationRawResponse(t, scalingConfigurationId)
						t.Logf("describe without tags returned Tags=%#v", response["Tags"])
						if tags, ok := response["Tags"].(map[string]interface{}); ok {
							if tagList, ok := tags["Tag"].([]interface{}); ok && len(tagList) > 0 {
								return fmt.Errorf("expected no remote tags, got %#v", tagList)
							}
						}
						if got := rs.Primary.Attributes["tags.%"]; got != "" && got != "0" {
							return fmt.Errorf("expected an empty tags map in state, got %q", got)
						}
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					modifyEssScalingConfigurationOutOfBand(t, scalingConfigurationId, externalUserData, externalTags)
				},
				Config: testAccConfig(map[string]interface{}{
					"force_delete": "true",
				}),
				Check: testAccCheck(map[string]string{
					"tags.external":   "ack",
					"tags.managed-by": "autoscaler",
					"user_data":       externalUserDataPlain,
				}),
			},
			{
				// Delete remote tags without going through Terraform: a
				// non-empty tag document replaces the whole set, so keeping
				// only "external" drops "managed-by" from the service. The
				// configuration still declares no tags, so the refresh has to
				// shrink the state to the observed tag and the plan after it
				// has to be empty instead of resurrecting the deleted tag.
				PreConfig: func() {
					modifyEssScalingConfigurationTagsOutOfBand(t, scalingConfigurationId, map[string]string{"external": "ack"})
				},
				Config: testAccConfig(map[string]interface{}{
					"force_delete": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.external":   "ack",
						"tags.managed-by": REMOVEKEY,
						"user_data":       externalUserDataPlain,
					}),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.managed-by"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.created"),
					func(state *terraform.State) error {
						response := essScalingConfigurationRawResponse(t, scalingConfigurationId)
						observed := tagsToMap(essTagList(response))
						t.Logf("after the out-of-band tag deletion: remote Tags=%#v observed=%#v", response["Tags"], observed)
						if len(observed) != 1 || observed["external"] != "ack" {
							return fmt.Errorf("expected the remote tags to shrink to {external: ack}, got %#v", observed)
						}
						return nil
					},
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]interface{}{
						"external":   "ack",
						"managed-by": "autoscaler",
						"created":    "tf",
					},
					"user_data": updatedUserData,
				}),
				Check: testAccCheck(map[string]string{
					"tags.external":   "ack",
					"tags.managed-by": "autoscaler",
					"tags.created":    "tf",
					"user_data":       updatedUserData,
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags":      REMOVEKEY,
					"user_data": REMOVEKEY,
				}),
				Check: testAccCheck(map[string]string{
					"tags.external":   "ack",
					"tags.managed-by": "autoscaler",
					"tags.created":    "tf",
					"user_data":       updatedUserData,
				}),
			},
			{
				// Send the clearing payloads to the service directly, without
				// going through the provider, so a surviving value proves that
				// ESS ignored the clear rather than the provider skipping it.
				PreConfig: func() {
					probeEssScalingConfigurationClearSemantics(t, scalingConfigurationId)
				},
				Config: testAccConfig(map[string]interface{}{
					"force_delete": "true",
				}),
				Check: testAccCheck(map[string]string{
					"tags.external":   "ack",
					"tags.managed-by": "autoscaler",
					"tags.created":    "tf",
					"user_data":       updatedUserData,
				}),
			},
			{
				// The two empty values do not leave the provider the same way.
				// An explicit empty tags map is still sent, because Update
				// publishes the whole tag document on any tags change; the
				// service ignores that empty document and keeps the previous
				// remote tags. An explicit empty user_data string is not sent at
				// all: GetOk reports an unset Optional+Computed string, so no
				// UserData parameter is built. In both cases the remote values
				// survive, Read restores them into state and the plan keeps
				// reporting the difference, so the residual diff is the
				// documented limitation rather than a convergent apply.
				Config: testAccConfig(map[string]interface{}{
					"tags":      CLEARMAP,
					"user_data": "",
				}),
				ExpectNonEmptyPlan: true,
				Check: testAccCheck(map[string]string{
					"tags.external":   "ack",
					"tags.managed-by": "autoscaler",
					"tags.created":    "tf",
					"user_data":       updatedUserData,
				}),
			},
		},
	})
}

// essScalingConfigurationRawResponse returns the DescribeScalingConfigurations
// body exactly as the common API decoded it, so a test can assert the service
// response instead of the values the provider wrote into state.
func essScalingConfigurationRawResponse(t *testing.T, scalingConfigurationId string) map[string]interface{} {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	response, err := (&EssService{client}).DescribeEssScalingConfigurationByCommonApi(scalingConfigurationId)
	if err != nil {
		t.Fatalf("failed to describe scaling configuration %s: %s", scalingConfigurationId, err)
	}
	return response
}

// sendModifyEssScalingConfigurationRaw sends one ModifyScalingConfiguration
// request with exactly the supplied parameters, including empty strings and empty
// JSON documents, which the typed SDK request may drop.
func sendModifyEssScalingConfigurationRaw(scalingConfigurationId string, parameters map[string]interface{}) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId":               client.RegionId,
		"ScalingConfigurationId": scalingConfigurationId,
	}
	for key, value := range parameters {
		request[key] = value
	}
	_, err := client.RpcPost("Ess", "2014-08-28", "ModifyScalingConfiguration", nil, request, true)
	return err
}

// probeEssScalingConfigurationClearSemantics calls ModifyScalingConfiguration
// with the payloads that would clear instance tags and user data and records what
// the service does with each of them.
func probeEssScalingConfigurationClearSemantics(t *testing.T, scalingConfigurationId string) {
	before := essScalingConfigurationRawResponse(t, scalingConfigurationId)
	t.Logf("before clearing attempts: Tags=%#v UserData=%#v", before["Tags"], before["UserData"])
	for _, parameters := range []map[string]interface{}{
		{"Tags": "{}"},
		{"UserData": ""},
	} {
		if err := sendModifyEssScalingConfigurationRaw(scalingConfigurationId, parameters); err != nil {
			t.Fatalf("failed to send %v to scaling configuration %s: %s", parameters, scalingConfigurationId, err)
		}
		after := essScalingConfigurationRawResponse(t, scalingConfigurationId)
		t.Logf("after ModifyScalingConfiguration %#v: Tags=%#v UserData=%#v", parameters, after["Tags"], after["UserData"])
	}
	// A JSON array is not an accepted tag document: the service rejects it with
	// InvalidTags.Format instead of reading it as a clear request.
	if err := sendModifyEssScalingConfigurationRaw(scalingConfigurationId, map[string]interface{}{"Tags": "[]"}); err == nil {
		t.Fatalf("expected the service to reject a JSON array tag document")
	} else {
		t.Logf("ModifyScalingConfiguration with Tags=[] was rejected: %s", err)
	}
	// The tag APIs are scoped to scaling groups, so a scaling configuration has
	// no untag endpoint either; the probe records the rejection instead of
	// claiming that an untag call was silently ignored.
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	untagRequest := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ResourceType": "scalingconfiguration",
		"ResourceId.1": scalingConfigurationId,
		"All":          true,
	}
	if _, err := client.RpcPost("Ess", "2014-08-28", "UntagResources", nil, untagRequest, true); err == nil {
		t.Logf("UntagResources with ResourceType=scalingconfiguration was accepted, revisit the tag removal semantics")
	} else {
		t.Logf("UntagResources with ResourceType=scalingconfiguration was rejected: %s", err)
	}
	after := essScalingConfigurationRawResponse(t, scalingConfigurationId)
	if beforeTags, afterTags := tagsToMap(essTagList(before)), tagsToMap(essTagList(after)); len(beforeTags) != len(afterTags) {
		t.Fatalf("the service changed the remote tags from %v to %v while clearing them", beforeTags, afterTags)
	} else {
		for key, value := range beforeTags {
			if afterTags[key] != value {
				t.Fatalf("the service changed the remote tag %q from %v to %v while clearing it", key, value, afterTags[key])
			}
		}
	}
	if before["UserData"] != after["UserData"] {
		t.Fatalf("the service changed the remote user data from %v to %v while clearing it", before["UserData"], after["UserData"])
	}
}

// essTagList returns the raw tag list of a DescribeScalingConfigurations
// response so tag comparisons reuse the same decoding the provider uses.
func essTagList(response map[string]interface{}) interface{} {
	tags, ok := response["Tags"].(map[string]interface{})
	if !ok {
		return nil
	}
	return tags["Tag"]
}

// modifyEssScalingConfigurationOutOfBand simulates a managed service such as
// the ACK cluster autoscaler writing user_data and tags directly through the
// ESS API, bypassing Terraform and its state entirely.
func modifyEssScalingConfigurationOutOfBand(t *testing.T, scalingConfigurationId, userData string, tags map[string]string) {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	tagPairs := make([]string, 0, len(tags))
	for key, value := range tags {
		tagPairs = append(tagPairs, fmt.Sprintf("%q:%q", key, value))
	}
	request := ess.CreateModifyScalingConfigurationRequest()
	request.RegionId = client.RegionId
	request.ScalingConfigurationId = scalingConfigurationId
	request.UserData = userData
	request.Tags = "{" + strings.Join(tagPairs, ",") + "}"
	if _, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingConfiguration(request)
	}); err != nil {
		t.Fatalf("failed to modify scaling configuration %s out of band: %s", scalingConfigurationId, err)
	}

	deadline := time.Now().Add(2 * time.Minute)
	for {
		config, err := (&EssService{client}).DescribeEssScalingConfiguration(scalingConfigurationId)
		if err == nil && essUserDataMatches(config.UserData, userData) {
			remoteTags := make(map[string]string)
			for _, tag := range config.Tags.Tag {
				remoteTags[tag.Key] = tag.Value
			}
			matched := true
			for key, value := range tags {
				if remoteTags[key] != value {
					matched = false
					break
				}
			}
			if matched {
				return
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for the out-of-band modification of scaling configuration %s, last error: %v", scalingConfigurationId, err)
		}
		time.Sleep(5 * time.Second)
	}
}

func essUserDataMatches(remote, submitted string) bool {
	if remote == submitted {
		return true
	}
	// The API may echo user data base64-decoded instead of verbatim.
	decoded, err := base64.StdEncoding.DecodeString(submitted)
	return err == nil && remote == string(decoded)
}

// modifyEssScalingConfigurationTagsOutOfBand replaces the instance tags of a
// scaling configuration with exactly the supplied set through the ESS API. A
// non-empty tag document replaces the previous set, so the tags left out of it
// are deleted remotely; an empty document is ignored by the service, which is
// why the last tag cannot be removed this way.
func modifyEssScalingConfigurationTagsOutOfBand(t *testing.T, scalingConfigurationId string, tags map[string]string) {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	tagPairs := make([]string, 0, len(tags))
	for key, value := range tags {
		tagPairs = append(tagPairs, fmt.Sprintf("%q:%q", key, value))
	}
	request := ess.CreateModifyScalingConfigurationRequest()
	request.RegionId = client.RegionId
	request.ScalingConfigurationId = scalingConfigurationId
	request.Tags = "{" + strings.Join(tagPairs, ",") + "}"
	if _, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingConfiguration(request)
	}); err != nil {
		t.Fatalf("failed to replace the tags of scaling configuration %s out of band: %s", scalingConfigurationId, err)
	}

	deadline := time.Now().Add(2 * time.Minute)
	for {
		config, err := (&EssService{client}).DescribeEssScalingConfiguration(scalingConfigurationId)
		if err == nil {
			remoteTags := make(map[string]string)
			for _, tag := range config.Tags.Tag {
				remoteTags[tag.Key] = tag.Value
			}
			matched := len(remoteTags) == len(tags)
			for key, value := range tags {
				if remoteTags[key] != value {
					matched = false
					break
				}
			}
			if matched {
				t.Logf("the remote tags of scaling configuration %s are now %v", scalingConfigurationId, remoteTags)
				return
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for the out-of-band tag replacement of scaling configuration %s, last error: %v", scalingConfigurationId, err)
		}
		time.Sleep(5 * time.Second)
	}
}

func resourceEssScalingConfigurationUserDataTagsDependence(name string) string {
	return fmt.Sprintf(`
	%s

	variable "name" {
		default = "%s"
	}
	data "alicloud_instance_types" "compatible" {
		availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
		image_id             = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_essd"
	}
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

// TestAccAliCloudEssScalingConfiguration_TypeListOrderConvergence proves that
// reordering the multi-member TypeList attributes of the scaling configuration
// converges. Every consecutive step permutes exactly one list and relies on the
// SDK post-apply plan check, so a canonicalizing API that does not echo the
// submitted order fails here instead of drifting silently in user plans.
func TestAccAliCloudEssScalingConfiguration_TypeListOrderConvergence(t *testing.T) {
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":     CHECKSET,
		"instance_types.#":     "2",
		"security_group_ids.#": "2",
		"override":             "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", acctest.RandIntRange(1000, 999999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationTypeListOrderDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default.images.0.id}",
					"instance_types":       []string{"${data.alicloud_instance_types.compatible.instance_types.0.id}", "${data.alicloud_instance_types.compatible.instance_types.1.id}"},
					"security_group_ids":   []string{"${alicloud_security_group.default.id}", "${alicloud_security_group.default1.id}"},
					"system_disk_category": "cloud_essd",
					"force_delete":         "true",
					"data_disk": []map[string]string{
						{
							"size":     "40",
							"category": "cloud_essd",
							"name":     "${var.name}-1",
						},
						{
							"size":     "50",
							"category": "cloud_essd",
							"name":     "${var.name}-2",
						},
					},
				}),
				Check: testAccCheck(map[string]string{
					"instance_types.#":     "2",
					"security_group_ids.#": "2",
					"data_disk.#":          "2",
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_types": []string{"${data.alicloud_instance_types.compatible.instance_types.1.id}", "${data.alicloud_instance_types.compatible.instance_types.0.id}"},
				}),
				Check: testAccCheck(map[string]string{
					"instance_types.#": "2",
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{"${alicloud_security_group.default1.id}", "${alicloud_security_group.default.id}"},
				}),
				Check: testAccCheck(map[string]string{
					"security_group_ids.#": "2",
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]string{
						{
							"size":     "50",
							"category": "cloud_essd",
							"name":     "${var.name}-2",
						},
						{
							"size":     "40",
							"category": "cloud_essd",
							"name":     "${var.name}-1",
						},
					},
				}),
				Check: testAccCheck(map[string]string{
					"data_disk.#":      "2",
					"data_disk.0.size": "50",
					"data_disk.1.size": "40",
				}),
			},
		},
	})
}

// TestAccAliCloudEssScalingConfiguration_InstancePatternInfoOrder covers the
// member order of the TypeLists nested in instance_pattern_info: a single outer
// set member keeps the container identity stable while architectures and
// excluded_instance_types are permuted in separate applies.
func TestAccAliCloudEssScalingConfiguration_InstancePatternInfoOrder(t *testing.T) {
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":        CHECKSET,
		"security_group_id":       CHECKSET,
		"image_id":                REGEXMATCH + "^debian",
		"override":                "false",
		"instance_pattern_info.#": "1",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", acctest.RandIntRange(1000, 999999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationConfigDependenceOnlyInstancePatternInfo)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_id":          "${data.alicloud_images.default2.images.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"force_delete":      "true",
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":   "EntryLevel",
							"cores":                   "4",
							"memory":                  "8.00",
							"burstable_performance":   "Include",
							"architectures":           []string{"X86", "Arm"},
							"excluded_instance_types": []string{"ecs.c6.large", "ecs.c6.xlarge"},
						},
					},
				}),
				Check: testAccCheck(map[string]string{
					"instance_pattern_info.#": "1",
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":   "EntryLevel",
							"cores":                   "4",
							"memory":                  "8.00",
							"burstable_performance":   "Include",
							"architectures":           []string{"Arm", "X86"},
							"excluded_instance_types": []string{"ecs.c6.large", "ecs.c6.xlarge"},
						},
					},
				}),
				Check: testAccCheck(map[string]string{
					"instance_pattern_info.#": "1",
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_pattern_info": []map[string]interface{}{
						{
							"instance_family_level":   "EntryLevel",
							"cores":                   "4",
							"memory":                  "8.00",
							"burstable_performance":   "Include",
							"architectures":           []string{"Arm", "X86"},
							"excluded_instance_types": []string{"ecs.c6.xlarge", "ecs.c6.large"},
						},
					},
				}),
				Check: testAccCheck(map[string]string{
					"instance_pattern_info.#": "1",
				}),
			},
		},
	})
}

// TestAccAliCloudEssScalingConfiguration_NetworkInterfaceSecurityGroupOrder
// covers the TypeList nested in the network_interfaces set: the outer member
// keeps its identity while its security_group_ids members are permuted.
func TestAccAliCloudEssScalingConfiguration_NetworkInterfaceSecurityGroupOrder(t *testing.T) {
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"override":         "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", acctest.RandIntRange(1000, 999999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationTypeListOrderDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.compatible.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd",
					"force_delete":         "true",
				}),
				Check: testAccCheck(map[string]string{
					"security_group_id": CHECKSET,
				}),
			},
			{
				// ESS rejects NetworkInterfaces.N for the primary ENI together with the
				// top-level SecurityGroupId/SecurityGroupIds, so the primary ENI moves
				// into NetworkInterfaces by removing security_group_id.
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": REMOVEKEY,
					"network_interfaces": []map[string]interface{}{
						{
							"instance_type":                  "Primary",
							"network_interface_traffic_mode": "Standard",
							"security_group_ids":             []string{"${alicloud_security_group.default.id}", "${alicloud_security_group.default1.id}"},
						},
					},
				}),
				Check: testAccCheck(map[string]string{
					"security_group_id":    "",
					"network_interfaces.#": "1",
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": []map[string]interface{}{
						{
							"instance_type":                  "Primary",
							"network_interface_traffic_mode": "Standard",
							"security_group_ids":             []string{"${alicloud_security_group.default1.id}", "${alicloud_security_group.default.id}"},
						},
					},
				}),
				Check: testAccCheck(map[string]string{
					"network_interfaces.#": "1",
				}),
			},
		},
	})
}

func resourceEssScalingConfigurationTypeListOrderDependence(name string) string {
	return fmt.Sprintf(`
	%s

	variable "name" {
		default = "%s"
	}
	data "alicloud_instance_types" "compatible" {
		availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
		image_id             = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_essd"
	}
	resource "alicloud_security_group" "default1" {
		name   = "${var.name}-1"
		vpc_id = "${alicloud_vpc.default.id}"
	}
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}`, EcsInstanceCommonTestCase, name)
}

// TestAccAliCloudEssScalingConfiguration_NetworkInterfaceIPv6Transition covers a
// primary ENI that stops requesting IPv6 addresses: the configuration starts with
// one IPv6 address, then asks for zero, then omits the attribute. Each step pins
// what the service keeps, so the request filter that drops a zero count has
// evidence instead of an assumption.
func TestAccAliCloudEssScalingConfiguration_NetworkInterfaceIPv6Transition(t *testing.T) {
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipv6"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"instance_type":    CHECKSET,
		"override":         "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", acctest.RandIntRange(1000, 999999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationIPv6Dependence)

	primaryNetworkInterface := func(ipv6AddressCount string) []map[string]interface{} {
		item := map[string]interface{}{
			"instance_type":                  "Primary",
			"network_interface_traffic_mode": "Standard",
			"security_group_ids":             []string{"${alicloud_security_group.default.id}"},
		}
		if ipv6AddressCount != "" {
			item["ipv6_address_count"] = ipv6AddressCount
		}
		return []map[string]interface{}{item}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_name":           "${data.alicloud_images.default.images.0.name}",
					"instance_type":        "${data.alicloud_instance_types.compatible.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd",
					"force_delete":         "true",
				}),
				Check: testAccCheck(map[string]string{
					"security_group_id": CHECKSET,
				}),
			},
			{
				// The primary ENI takes over the security group from the
				// top-level attribute and requests one IPv6 address.
				Config: testAccConfig(map[string]interface{}{
					"security_group_id":  REMOVEKEY,
					"network_interfaces": primaryNetworkInterface("1"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    "",
						"network_interfaces.#": "1",
					}),
					essScalingConfigurationNetworkInterfaceIPv6Check(t, resourceId, "1", "1"),
				),
			},
			{
				// Removing the count from a non-zero configuration has to clear
				// it as well, without passing through an explicit zero. The
				// attribute leaves the state together with the configuration, so
				// an absent state value and an absent remote count are the
				// converged form here.
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": primaryNetworkInterface(""),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interfaces.#": "1",
					}),
					essScalingConfigurationNetworkInterfaceIPv6Check(t, resourceId, "", ""),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": primaryNetworkInterface("1"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interfaces.#": "1",
					}),
					essScalingConfigurationNetworkInterfaceIPv6Check(t, resourceId, "1", "1"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": primaryNetworkInterface("0"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interfaces.#": "1",
					}),
					essScalingConfigurationNetworkInterfaceIPv6Check(t, resourceId, "0", ""),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": primaryNetworkInterface(""),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_interfaces.#": "1",
					}),
					essScalingConfigurationNetworkInterfaceIPv6Check(t, resourceId, "0", ""),
				),
			},
		},
	})
}

// TestAccAliCloudEssScalingConfiguration_NetworkInterfaceIPv6Negative proves that
// a negative ipv6_address_count is rejected while the configuration is planned,
// for both creation and update, instead of being silently dropped by the request
// builder.
func TestAccAliCloudEssScalingConfiguration_NetworkInterfaceIPv6Negative(t *testing.T) {
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipv6_negative"
	checkoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"instance_type":    CHECKSET,
		"override":         "false",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", acctest.RandIntRange(1000, 999999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssScalingConfigurationIPv6Dependence)

	primaryNetworkInterface := func(ipv6AddressCount string) []map[string]interface{} {
		item := map[string]interface{}{
			"instance_type":                  "Primary",
			"network_interface_traffic_mode": "Standard",
			"security_group_ids":             []string{"${alicloud_security_group.default.id}"},
		}
		if ipv6AddressCount != "" {
			item["ipv6_address_count"] = ipv6AddressCount
		}
		return []map[string]interface{}{item}
	}
	negativeCountError := regexp.MustCompile(`ipv6_address_count.*at least \(0\)`)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// A negative count is rejected while the configuration is
				// planned, so creation never reaches the request builder.
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_name":           "${data.alicloud_images.default.images.0.name}",
					"instance_type":        "${data.alicloud_instance_types.compatible.instance_types.0.id}",
					"system_disk_category": "cloud_essd",
					"force_delete":         "true",
					"network_interfaces":   primaryNetworkInterface("-1"),
				}),
				ExpectError: negativeCountError,
			},
			{
				// A valid creation the negative update can start from.
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_name":           "${data.alicloud_images.default.images.0.name}",
					"instance_type":        "${data.alicloud_instance_types.compatible.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd",
					"force_delete":         "true",
					"network_interfaces":   REMOVEKEY,
				}),
				Check: testAccCheck(map[string]string{
					"security_group_id": CHECKSET,
				}),
			},
			{
				// The same rejection has to happen when an existing
				// configuration is updated, not only on creation.
				Config: testAccConfig(map[string]interface{}{
					"security_group_id":  REMOVEKEY,
					"network_interfaces": primaryNetworkInterface("-1"),
				}),
				ExpectError: negativeCountError,
			},
			{
				// The configuration is unchanged by the rejected update and a
				// valid count still applies normally.
				Config: testAccConfig(map[string]interface{}{
					"network_interfaces": primaryNetworkInterface("1"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    "",
						"network_interfaces.#": "1",
					}),
					essScalingConfigurationNetworkInterfaceIPv6Check(t, resourceId, "1", "1"),
				),
			},
		},
	})
}

// essScalingConfigurationNetworkInterfaceIPv6Check compares the IPv6 address
// count in state with the expected value and logs what the service reports.
func essScalingConfigurationNetworkInterfaceIPv6Check(t *testing.T, resourceId, expectedState, expectedRemote string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		count := ""
		for key, value := range rs.Primary.Attributes {
			if strings.HasPrefix(key, "network_interfaces.") && strings.HasSuffix(key, ".ipv6_address_count") {
				count = value
			}
		}
		response := essScalingConfigurationRawResponse(t, rs.Primary.ID)
		// The service omits Ipv6AddressCount once the count is cleared, so an
		// absent key is the remote form of zero.
		remote := ""
		if interfaces, ok := response["NetworkInterfaces"].(map[string]interface{}); ok {
			if items, ok := interfaces["NetworkInterface"].([]interface{}); ok {
				for _, item := range items {
					pack, ok := item.(map[string]interface{})
					if !ok || pack["InstanceType"] != "Primary" {
						continue
					}
					if value, ok := pack["Ipv6AddressCount"]; ok {
						remote = fmt.Sprintf("%v", value)
					}
				}
			}
		}
		t.Logf("state ipv6_address_count=%q, remote Ipv6AddressCount=%q, remote NetworkInterfaces=%#v", count, remote, response["NetworkInterfaces"])
		if count != expectedState {
			return fmt.Errorf("expected network_interfaces ipv6_address_count %q, got %q", expectedState, count)
		}
		if remote != expectedRemote {
			return fmt.Errorf("expected remote Ipv6AddressCount %q, got %q", expectedRemote, remote)
		}
		return nil
	}
}

func resourceEssScalingConfigurationIPv6Dependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	data "alicloud_zones" "default" {
		available_disk_category     = "cloud_essd"
		available_resource_creation = "VSwitch"
	}
	data "alicloud_images" "default" {
		name_regex  = "^aliyun.*vhd"
		most_recent = true
		owners      = "system"
	}
	data "alicloud_instance_types" "compatible" {
		availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
		image_id             = "${data.alicloud_images.default.images.0.id}"
		system_disk_category = "cloud_essd"
	}
	resource "alicloud_vpc" "default" {
		vpc_name    = "${var.name}"
		cidr_block  = "172.16.0.0/16"
		enable_ipv6 = "true"
		ipv6_isp    = "BGP"
	}
	resource "alicloud_vswitch" "default" {
		vpc_id               = "${alicloud_vpc.default.id}"
		cidr_block           = "172.16.1.0/24"
		zone_id              = "${data.alicloud_zones.default.zones.0.id}"
		vswitch_name         = "${var.name}"
		ipv6_cidr_block_mask = "8"
	}
	resource "alicloud_security_group" "default" {
		name   = "${var.name}"
		vpc_id = "${alicloud_vpc.default.id}"
	}
	resource "alicloud_ess_scaling_group" "default" {
		min_size           = 1
		max_size           = 1
		scaling_group_name = "${var.name}"
		removal_policies   = ["OldestInstance", "NewestInstance"]
		vswitch_ids        = ["${alicloud_vswitch.default.id}"]
	}`, name)
}

// essScalingConfigurationUnitTestClient returns a provider client whose ESS
// endpoint is a local HTTP server, so a unit test can inspect the requests the
// provider itself serializes without reaching the cloud.
func essScalingConfigurationUnitTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("ess", endpoint)
	return client
}

// TestUnitEssScalingConfigurationProviderClearRequest captures the request the
// provider serializes when tags become an explicit empty map and user_data
// becomes an empty string: the empty tag document is sent, while no UserData
// parameter is built at all.
func TestUnitEssScalingConfigurationProviderClearRequest(t *testing.T) {
	var mutex sync.Mutex
	modifyRequests := make([]map[string]string, 0)
	described := 0
	client := essScalingConfigurationUnitTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("cannot parse the provider request: %s", err)
		}
		form := make(map[string]string, len(r.Form))
		for key, values := range r.Form {
			form[key] = strings.Join(values, ",")
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "ModifyScalingConfiguration":
			mutex.Lock()
			modifyRequests = append(modifyRequests, form)
			mutex.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"ScalingConfigurationId": "asc-unit-test"})
		case "DescribeScalingConfigurations":
			mutex.Lock()
			described++
			mutex.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"ScalingConfigurations": map[string]interface{}{
					"ScalingConfiguration": []interface{}{map[string]interface{}{
						"ScalingConfigurationId": "asc-unit-test",
						"ScalingGroupId":         "asg-unit-test",
						"LifecycleState":         "Active",
					}},
				},
			})
		default:
			t.Errorf("unexpected Action %q", r.Form.Get("Action"))
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	scalingConfigurationResource := resourceAlicloudEssScalingConfiguration()
	previous := &sdktf.InstanceState{
		ID: "asc-unit-test",
		Attributes: map[string]string{
			"id":               "asc-unit-test",
			"scaling_group_id": "asg-unit-test",
			"force_delete":     "true",
			"tags.%":           "3",
			"tags.external":    "ack",
			"tags.managed-by":  "autoscaler",
			"tags.created":     "tf",
			"user_data":        "IyEvYmluL2Jhc2gKZWNobyB0ZXJyYWZvcm0tdXBkYXRlZA==",
		},
	}
	config := sdktf.NewResourceConfigRaw(map[string]interface{}{
		"scaling_group_id": "asg-unit-test",
		"force_delete":     true,
		"tags":             map[string]interface{}{},
		"user_data":        "",
	})
	diff, err := scalingConfigurationResource.Diff(context.Background(), previous, config, client)
	if err != nil {
		t.Fatalf("diff failed: %s", err)
	}
	if diff == nil || diff.Empty() {
		t.Fatal("expected the explicit empty values to produce a diff")
	}
	newState, diags := scalingConfigurationResource.Apply(context.Background(), previous, diff, client)
	if diags.HasError() {
		t.Fatalf("apply reported errors: %v", diags)
	}
	if newState == nil {
		t.Fatal("apply returned no state")
	}

	mutex.Lock()
	defer mutex.Unlock()
	if len(modifyRequests) != 1 {
		t.Fatalf("expected exactly one ModifyScalingConfiguration request, got %d: %v", len(modifyRequests), modifyRequests)
	}
	request := modifyRequests[0]
	if request["ScalingConfigurationId"] != "asc-unit-test" {
		t.Fatalf("unexpected ScalingConfigurationId %q", request["ScalingConfigurationId"])
	}
	if request["Tags"] != "{}" {
		t.Fatalf("expected the provider to send Tags=\"{}\", got %q", request["Tags"])
	}
	if value, ok := request["UserData"]; ok {
		t.Fatalf("expected no UserData parameter, got %q", value)
	}
	if described == 0 {
		t.Fatal("the update did not run the resource Read")
	}
	t.Logf("captured ModifyScalingConfiguration parameters: Tags=%q, UserData=<absent>, Read described %d times", request["Tags"], described)
}
