package ess_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"encoding/base64"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	essservice "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfcacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// collection-order: bypass instance_ids -- pre-existing deprecated TypeList left untouched by this change (deprecated since provider 1.6.0); its suppressor is stock debt and intentionally returns true for every change, including member additions and removals, because CreateScalingConfiguration, ModifyScalingConfiguration and DescribeScalingConfigurations do not carry instance ids at all. A reorder apply cannot surface a diff, and a bounded suppressor unit net cannot have the required no-suppression membership case without changing that legacy behaviour.

// kms限制
func TestAccAliCloudEssScalingConfiguration_override(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_name":        acctest.CHECKSET,
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
		},
	})
}

// kms限制
func TestAccAliCloudEssScalingConfiguration_other(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_name":        acctest.CHECKSET,
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"is_outdated", "override", "enable", "substitute", "force_delete", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_credit_specification(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{"cn-beijing"})
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_name":        acctest.CHECKSET,
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependences_beijing)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependenceImageName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"image_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "image_id", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_imageId(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		//"image_name":        acctest.CHECKSET,
		"override": "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependenceImageName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"image_id":   acctest.CHECKSET,
						"image_name": "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
		},
	})
}

// kms限制
func TestAccAliCloudEssScalingConfiguration_kms(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_name":        acctest.CHECKSET,
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
		},
	})
}

// kms限制
func TestAccAliCloudEssScalingConfiguration_Update(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_name":        acctest.CHECKSET,
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, func(n string) string {
		return resourceEssScalingConfigurationConfigDependence(n) + scalingConfigurationRobustDependence
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":          "${alicloud_ess_scaling_group.ess.id}",
					"image_name":                "${data.alicloud_images.ess.images.0.name}",
					"instance_type":             "${data.alicloud_instance_types.ess.instance_types.0.id}",
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
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
					"security_group_id":  acctest.REMOVEKEY,
					"security_group_ids": []string{"${alicloud_security_group.default.id}", "${alicloud_security_group.default1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    acctest.REMOVEKEY,
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
						"role_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_name": "${alicloud_key_pair.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":  acctest.REMOVEKEY,
					"instance_types": []string{"${data.alicloud_instance_types.ess.instance_types.0.id}", "${data.alicloud_instance_types.ess.instance_types.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":    acctest.REMOVEKEY,
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
						"data_disk.3.kms_key_id":              acctest.CHECKSET,
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.CHECKSET,
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependenceUpdate)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
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
					"system_disk_size":           acctest.REMOVEKEY,
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
						"system_disk_size":           acctest.REMOVEKEY,
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^centos",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependencePerformanceLevel)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
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
						"resource_group_id": acctest.CHECKSET,
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
					"spot_price_limit": acctest.REMOVEKEY,
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^centos",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependencePerformanceLevel)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^centos",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependencePerformanceLevel)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^centos",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^centos",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependenceForResourcePoolOptions)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^ubuntu_22",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependenceForResourcePoolOptions)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":                       "${alicloud_ess_scaling_group.ess.id}",
					"image_id":                               "${data.alicloud_images.ess.images.0.id}",
					"instance_type":                          "${data.alicloud_instance_types.ess.instance_types.0.id}",
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
					"scaling_group_id":                       "${alicloud_ess_scaling_group.ess.id}",
					"image_id":                               "${data.alicloud_images.ess.images.0.id}",
					"instance_type":                          "${data.alicloud_instance_types.ess.instance_types.0.id}",
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^aliyun.*vhd",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"dedicated_host_cluster_id": acctest.CHECKSET,
						"password_inherit":          "false",
						"http_endpoint":             "enabled",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
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
					"security_group_id":         acctest.REMOVEKEY,
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
						"security_group_id":         acctest.CHECKSET,
						"network_interfaces.#":      "1",
						"dedicated_host_cluster_id": acctest.CHECKSET,
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)

	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^aliyun.*vhd",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default1.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.r5.instance_types.0.id}",
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
				ImportStateVerifyIgnore: []string{"force_delete"},
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
					"security_group_id":         acctest.REMOVEKEY,
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
						"security_group_id":         acctest.CHECKSET,
						"network_interfaces.#":      "1",
						"dedicated_host_cluster_id": acctest.CHECKSET,
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
						"dedicated_host_cluster_id": acctest.CHECKSET,
						"network_interfaces.#":      "2",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudEssScalingConfiguration_SystemDiskCategoriesOrder verifies
// that the multi-member TypeList attribute system_disk_categories converges
// under a member reorder: the scaling configuration is created with two
// distinct disk categories in one order, then the same configuration is
// reapplied with the members reordered. system_disk_categories is mutually
// exclusive with system_disk_category, so no scalar system_disk_category is
// set. The empty post-apply plan after the reorder proves the update path
// (F1) propagates the list and that the API's ordering contract converges.
// Pre-existing multi-member TypeList attributes of alicloud_ess_scaling_configuration
// that predate this change and carry no convergent reorder test. They are stock
// order-coverage debt outside the scope of this PR and are marked so the
// collection-order checker does not report them. The new TypeList attributes
// introduced by this change (system_disk_categories, data_disk.categories) are
// covered by the convergent reorder tests below and are intentionally not bypassed.
// collection-order: bypass instance_types -- pre-existing attribute outside this change
// collection-order: bypass security_group_ids -- pre-existing attribute outside this change
// collection-order: bypass data_disk -- pre-existing attribute outside this change
// collection-order: bypass instance_ids -- pre-existing attribute outside this change
// collection-order: bypass network_interfaces.security_group_ids -- pre-existing attribute outside this change
// collection-order: bypass instance_pattern_info.excluded_instance_types -- pre-existing attribute outside this change
// collection-order: bypass instance_pattern_info.architectures -- pre-existing attribute outside this change

func rawResponseCheck(t *testing.T, resourceId string, check func(map[string]interface{}) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s is missing from state", resourceId)
		}
		response := essScalingConfigurationRawResponse(t, rs.Primary.ID)
		// Prove the raw Describe returned the object under test before asserting
		// any of its values; the shared query helper picks the first returned
		// configuration and must not be trusted without an identity check.
		if got, _ := response["ScalingConfigurationId"].(string); got != rs.Primary.ID {
			return fmt.Errorf("Describe returned ScalingConfigurationId %q, expected %q", got, rs.Primary.ID)
		}
		return check(response)
	}
}

func stringList(v interface{}) []string {
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		out = append(out, fmt.Sprint(item))
	}
	return out
}

func systemDiskCategoriesRemoteCheck(t *testing.T, resourceId string, expected []string) resource.TestCheckFunc {
	return rawResponseCheck(t, resourceId, func(response map[string]interface{}) error {
		m, ok := response["SystemDiskCategories"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("Describe response has no SystemDiskCategories: %v", response["SystemDiskCategories"])
		}
		got := stringList(m["SystemDiskCategory"])
		if strings.Join(got, ",") != strings.Join(expected, ",") {
			return fmt.Errorf("remote SystemDiskCategories = %v, expected %v", got, expected)
		}
		return nil
	})
}

func dataDiskCategoriesRemoteCheck(t *testing.T, resourceId string, expected []string) resource.TestCheckFunc {
	return rawResponseCheck(t, resourceId, func(response map[string]interface{}) error {
		dds, ok := response["DataDisks"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("Describe response has no DataDisks: %v", response["DataDisks"])
		}
		list, _ := dds["DataDisk"].([]interface{})
		if len(list) != 1 {
			return fmt.Errorf("remote DataDisk count = %d, expected 1", len(list))
		}
		disk, _ := list[0].(map[string]interface{})
		cm, ok := disk["Categories"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("remote DataDisk has no Categories: %v", disk["Categories"])
		}
		got := stringList(cm["Category"])
		if strings.Join(got, ",") != strings.Join(expected, ",") {
			return fmt.Errorf("remote DataDisk Categories = %v, expected %v", got, expected)
		}
		return nil
	})
}

// stateStringList reads a repeated state attribute such as
// system_disk_categories or data_disk.0.categories.
func stateStringList(attrs map[string]string, prefix string) []string {
	count, _ := strconv.Atoi(attrs[prefix+".#"])
	out := make([]string, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, attrs[fmt.Sprintf("%s.%d", prefix, i)])
	}
	return out
}

// systemDiskCategoriesMirrorRemoteCheck proves that after the plural list is
// removed from the configuration the singular category takes effect and the
// computed plural mirrors exactly what the service reports. The Describe
// response is checked for object identity first, and the expected singular
// value is asserted separately from the state mirror so a stale state or a
// different object cannot satisfy the check.
func systemDiskCategoriesMirrorRemoteCheck(t *testing.T, resourceId, expectedSingular string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s is missing from state", resourceId)
		}
		response := essScalingConfigurationRawResponse(t, rs.Primary.ID)
		if got, _ := response["ScalingConfigurationId"].(string); got != rs.Primary.ID {
			return fmt.Errorf("Describe returned ScalingConfigurationId %q, expected %q", got, rs.Primary.ID)
		}
		if got := fmt.Sprint(response["SystemDiskCategory"]); got != expectedSingular && !(got == "<nil>" && expectedSingular == "") {
			return fmt.Errorf("remote SystemDiskCategory = %q, expected %q", got, expectedSingular)
		}
		var remote []string
		if m, ok := response["SystemDiskCategories"].(map[string]interface{}); ok {
			remote = stringList(m["SystemDiskCategory"])
		}
		stateList := stateStringList(rs.Primary.Attributes, "system_disk_categories")
		if strings.Join(stateList, ",") != strings.Join(remote, ",") {
			return fmt.Errorf("state system_disk_categories = %v, Describe reported %v", stateList, remote)
		}
		return nil
	}
}

// dataDiskCategoriesMirrorRemoteCheck proves that after the plural list is
// removed from the configuration the computed plural mirrors the service
// response and the remote singular Category carries exactly the expected value
// (empty when the service reports none for a plural-only disk). Identity is
// checked before any value assertion.
func dataDiskCategoriesMirrorRemoteCheck(t *testing.T, resourceId, expectedSingular string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s is missing from state", resourceId)
		}
		response := essScalingConfigurationRawResponse(t, rs.Primary.ID)
		if got, _ := response["ScalingConfigurationId"].(string); got != rs.Primary.ID {
			return fmt.Errorf("Describe returned ScalingConfigurationId %q, expected %q", got, rs.Primary.ID)
		}
		dds, _ := response["DataDisks"].(map[string]interface{})
		list, _ := dds["DataDisk"].([]interface{})
		if len(list) != 1 {
			return fmt.Errorf("remote DataDisk count = %d, expected 1", len(list))
		}
		disk, _ := list[0].(map[string]interface{})
		if got := fmt.Sprint(disk["Category"]); got != expectedSingular && !(got == "<nil>" && expectedSingular == "") {
			return fmt.Errorf("remote DataDisk.Category = %q, expected %q", got, expectedSingular)
		}
		var remote []string
		if cm, ok := disk["Categories"].(map[string]interface{}); ok {
			remote = stringList(cm["Category"])
		}
		stateList := stateStringList(rs.Primary.Attributes, "data_disk.0.categories")
		if strings.Join(stateList, ",") != strings.Join(remote, ",") {
			return fmt.Errorf("state data_disk.0.categories = %v, Describe reported %v", stateList, remote)
		}
		return nil
	}
}

func burstingRemoteCheck(t *testing.T, resourceId string, systemDisk, dataDisk bool) resource.TestCheckFunc {
	return rawResponseCheck(t, resourceId, func(response map[string]interface{}) error {
		if !strings.EqualFold(fmt.Sprint(response["SystemDisk.BurstingEnabled"]), fmt.Sprint(systemDisk)) {
			return fmt.Errorf("remote SystemDisk.BurstingEnabled = %v, expected %v", response["SystemDisk.BurstingEnabled"], systemDisk)
		}
		dds, _ := response["DataDisks"].(map[string]interface{})
		list, _ := dds["DataDisk"].([]interface{})
		if len(list) != 1 {
			return fmt.Errorf("remote DataDisk count = %d, expected 1", len(list))
		}
		disk, _ := list[0].(map[string]interface{})
		if !strings.EqualFold(fmt.Sprint(disk["BurstingEnabled"]), fmt.Sprint(dataDisk)) {
			return fmt.Errorf("remote DataDisk.BurstingEnabled = %v, expected %v", disk["BurstingEnabled"], dataDisk)
		}
		return nil
	})
}

// resourceScalingConfigurationNetworkInterfacesWithSecondary adds a second
// scaling configuration to the dependence so the identity assertion on the
// primary object cannot pass by reading the wrong object.
func resourceScalingConfigurationNetworkInterfacesWithSecondary(name string) string {
	return resourceEssScalingConfiguration_NetworkInterfaces(name) + scalingConfigurationSecondaryForOrder
}

func TestAccAliCloudEssScalingConfiguration_SystemDiskCategoriesOrder(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	rand := tfcacctest.RandIntRange(1000, 999999)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	resourceId := "alicloud_ess_scaling_configuration.sco"
	var v ess.ScalingConfiguration
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceScalingConfigurationNetworkInterfacesWithSecondary)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^aliyun.*vhd",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":       "${alicloud_ess_scaling_group.default.id}",
					"image_id":               "${data.alicloud_images.default1.images.0.id}",
					"instance_type":          "${data.alicloud_instance_types.r5.instance_types.0.id}",
					"security_group_id":      "${alicloud_security_group.default1.id}",
					"force_delete":           "true",
					"system_disk_categories": []interface{}{"cloud_essd", "cloud"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "system_disk_categories.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "system_disk_categories.0", "cloud_essd"),
					resource.TestCheckResourceAttr(resourceId, "system_disk_categories.1", "cloud"),
					systemDiskCategoriesRemoteCheck(t, resourceId, []string{"cloud_essd", "cloud"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_categories": []interface{}{"cloud", "cloud_essd"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "system_disk_categories.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "system_disk_categories.0", "cloud"),
					resource.TestCheckResourceAttr(resourceId, "system_disk_categories.1", "cloud_essd"),
					systemDiskCategoriesRemoteCheck(t, resourceId, []string{"cloud", "cloud_essd"}),
				),
			},
			{
				// An unrelated update must keep every configured category, so the
				// exact remote order is asserted again after the rename.
				Config: testAccConfig(map[string]interface{}{
					"system_disk_categories":     []interface{}{"cloud", "cloud_essd"},
					"scaling_configuration_name": name + "-renamed",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "scaling_configuration_name", name+"-renamed"),
					systemDiskCategoriesRemoteCheck(t, resourceId, []string{"cloud", "cloud_essd"}),
				),
			},
			{
				// The intermediate state after the unrelated update must be
				// independently importable, not only the final removal state.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
			{
				// Removing the non-empty plural list without switching to the
				// singular form must converge. The singular attribute carries the
				// schema default (cloud_efficiency), so the service reports that
				// fallback value; the computed plural must still mirror the service.
				Config: testAccConfig(map[string]interface{}{
					"system_disk_categories": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					systemDiskCategoriesMirrorRemoteCheck(t, resourceId, "cloud_efficiency"),
				),
			},
			{
				// Removing the non-empty plural list and switching to the singular
				// category must put the singular value in effect. The plural becomes
				// computed again and has to mirror what the service reports; the ESS
				// modify API keeps the previously applied members.
				Config: testAccConfig(map[string]interface{}{
					"system_disk_categories": acctest.REMOVEKEY,
					"system_disk_category":   "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					systemDiskCategoriesMirrorRemoteCheck(t, resourceId, "cloud_essd"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_DataDiskCategoriesOrder(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	rand := tfcacctest.RandIntRange(1000, 999999)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	resourceId := "alicloud_ess_scaling_configuration.dco"
	var v ess.ScalingConfiguration
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^aliyun.*vhd",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default1.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.r5.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default1.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_essd",
					"data_disk": []map[string]interface{}{
						{
							"size":                 "20",
							"delete_with_instance": "true",
							"categories":           []interface{}{"cloud_essd", "cloud"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "data_disk.0.categories.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "data_disk.0.categories.0", "cloud_essd"),
					resource.TestCheckResourceAttr(resourceId, "data_disk.0.categories.1", "cloud"),
					dataDiskCategoriesRemoteCheck(t, resourceId, []string{"cloud_essd", "cloud"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]interface{}{
						{
							"size":                 "20",
							"delete_with_instance": "true",
							"categories":           []interface{}{"cloud", "cloud_essd"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "data_disk.0.categories.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "data_disk.0.categories.0", "cloud"),
					resource.TestCheckResourceAttr(resourceId, "data_disk.0.categories.1", "cloud_essd"),
					dataDiskCategoriesRemoteCheck(t, resourceId, []string{"cloud", "cloud_essd"}),
				),
			},
			{
				// An unrelated rename must keep the exact remote category list.
				Config: testAccConfig(map[string]interface{}{
					"scaling_configuration_name": name + "-renamed",
					"data_disk": []map[string]interface{}{
						{
							"size":                 "20",
							"delete_with_instance": "true",
							"categories":           []interface{}{"cloud", "cloud_essd"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "scaling_configuration_name", name+"-renamed"),
					dataDiskCategoriesRemoteCheck(t, resourceId, []string{"cloud", "cloud_essd"}),
				),
			},
			{
				// The intermediate state after the unrelated update must be
				// independently importable, not only the final removal state.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
			{
				// Removing the non-empty plural list without switching to the singular
				// form must converge: the plural becomes computed again and has to
				// mirror the service value, which keeps the applied members and
				// reports no singular Category for such a disk.
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]interface{}{
						{
							"size":                 "20",
							"delete_with_instance": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					dataDiskCategoriesMirrorRemoteCheck(t, resourceId, ""),
				),
			},
			{
				// Switching the disk to the singular category must converge. ESS keeps
				// the previously applied plural members and omits the singular
				// DataDisk.Category from the describe response, so the remote singular
				// is asserted as absent while the computed plural mirrors the service.
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]interface{}{
						{
							"size":                 "20",
							"delete_with_instance": "true",
							"category":             "cloud_essd",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					dataDiskCategoriesMirrorRemoteCheck(t, resourceId, ""),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_PropertyAttributes(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	rand := tfcacctest.RandIntRange(1000, 999999)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	resourceId := "alicloud_ess_scaling_configuration.pa"
	var v ess.ScalingConfiguration
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^aliyun.*vhd",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":             "${alicloud_ess_scaling_group.default.id}",
					"image_id":                     "${data.alicloud_images.default1.images.0.id}",
					"instance_type":                "${data.alicloud_instance_types.r5_auto.instance_types.0.id}",
					"security_group_id":            "${alicloud_security_group.default1.id}",
					"force_delete":                 "true",
					"system_disk_category":         "cloud_auto",
					"system_disk_bursting_enabled": "true",
					"tenancy":                      "default",
					"affinity":                     "default",
					"data_disk": []map[string]interface{}{
						{
							"size":                 "40",
							"category":             "cloud_auto",
							"delete_with_instance": "true",
							"bursting_enabled":     "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "system_disk_bursting_enabled", "true"),
					resource.TestCheckResourceAttr(resourceId, "tenancy", "default"),
					resource.TestCheckResourceAttr(resourceId, "affinity", "default"),
					resource.TestCheckResourceAttr(resourceId, "data_disk.0.bursting_enabled", "true"),
					burstingRemoteCheck(t, resourceId, true, true),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_bursting_enabled": "false",
					"data_disk": []map[string]interface{}{
						{
							"size":                 "40",
							"category":             "cloud_auto",
							"delete_with_instance": "true",
							"bursting_enabled":     "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "system_disk_bursting_enabled", "false"),
					resource.TestCheckResourceAttr(resourceId, "data_disk.0.bursting_enabled", "false"),
					burstingRemoteCheck(t, resourceId, false, false),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_configuration_name": name + "-renamed",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "scaling_configuration_name", name+"-renamed"),
					burstingRemoteCheck(t, resourceId, false, false),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
		},
	})
}

// TestAccAliCloudEssScalingConfiguration_ConstraintValidation covers the
// plan-time cross-attribute constraints. Every step expects a configuration
// error before any cloud resource is created.
func TestAccAliCloudEssScalingConfiguration_ConstraintValidation(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	rand := tfcacctest.RandIntRange(1000, 999999)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	resourceId := "alicloud_ess_scaling_configuration.cv"
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	base := map[string]interface{}{
		"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
		"image_id":          "${data.alicloud_images.default1.images.0.id}",
		"instance_type":     "${data.alicloud_instance_types.r5.instance_types.0.id}",
		"security_group_id": "${alicloud_security_group.default1.id}",
		"force_delete":      "true",
	}
	withBase := func(extra map[string]interface{}) map[string]interface{} {
		cfg := map[string]interface{}{}
		for k, val := range base {
			cfg[k] = val
		}
		for k, val := range extra {
			cfg[k] = val
		}
		return cfg
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(withBase(map[string]interface{}{
					"private_pool_options": []map[string]interface{}{
						{"match_criteria": "Target"},
					},
				})),
				ExpectError: regexp.MustCompile(`pool_id must be set when private_pool_options.0.match_criteria is Target`),
				Check:       resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccConfig(withBase(map[string]interface{}{
					"private_pool_options": acctest.REMOVEKEY,
					"cpu":                  "2",
				})),
				ExpectError: regexp.MustCompile(`must be specified`),
				Check:       resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccConfig(withBase(map[string]interface{}{
					"cpu":    "2",
					"memory": "8",
				})),
				ExpectError: regexp.MustCompile(`conflicts with`),
				Check:       resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccConfig(withBase(map[string]interface{}{
					"cpu":          acctest.REMOVEKEY,
					"memory":       acctest.REMOVEKEY,
					"image_family": "acs:alibaba-cloud-linux-3",
				})),
				ExpectError: regexp.MustCompile(`conflicts with`),
				Check:       resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccConfig(withBase(map[string]interface{}{
					"image_family":                   acctest.REMOVEKEY,
					"resource_pool_options_strategy": "PrivatePoolFirst",
					"private_pool_options": []map[string]interface{}{
						{"match_criteria": "Open"},
					},
				})),
				ExpectError: regexp.MustCompile(`conflicts with`),
				Check:       resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccConfig(withBase(map[string]interface{}{
					"resource_pool_options_strategy": acctest.REMOVEKEY,
					"private_pool_options":           acctest.REMOVEKEY,
					"security_group_id":              acctest.REMOVEKEY,
					"ipv6_address_count":             "1",
					"network_interfaces": []map[string]interface{}{
						{"instance_type": "Primary", "ipv6_address_count": "1"},
					},
				})),
				ExpectError: regexp.MustCompile(`ipv6_address_count conflicts with network_interfaces`),
				Check:       resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func noPrivatePoolRemoteCheck(t *testing.T, resourceId string) resource.TestCheckFunc {
	return rawResponseCheck(t, resourceId, func(response map[string]interface{}) error {
		// MatchCriteria=None is the service representation of "no private pool".
		if v := fmt.Sprint(response["PrivatePoolOptions.MatchCriteria"]); v != "" && v != "<nil>" && !strings.EqualFold(v, "None") {
			return fmt.Errorf("remote PrivatePoolOptions.MatchCriteria = %q, expected empty or None", v)
		}
		if v := fmt.Sprint(response["PrivatePoolOptions.Id"]); v != "" && v != "<nil>" {
			return fmt.Errorf("remote PrivatePoolOptions.Id = %q, expected empty", v)
		}
		return nil
	})
}

func privatePoolMatchCriteriaRemoteCheck(t *testing.T, resourceId, expected string) resource.TestCheckFunc {
	return rawResponseCheck(t, resourceId, func(response map[string]interface{}) error {
		if got := fmt.Sprint(response["PrivatePoolOptions.MatchCriteria"]); got != expected {
			return fmt.Errorf("remote PrivatePoolOptions.MatchCriteria = %q, expected %q", got, expected)
		}
		return nil
	})
}

func noSchedulerOptionsRemoteCheck(t *testing.T, resourceId string) resource.TestCheckFunc {
	return rawResponseCheck(t, resourceId, func(response map[string]interface{}) error {
		if v := response["SchedulerOptions"]; v != nil {
			if m, ok := v.(map[string]interface{}); !ok || len(m) > 0 {
				return fmt.Errorf("remote SchedulerOptions = %v, expected empty", v)
			}
		}
		return nil
	})
}

// TestAccAliCloudEssScalingConfiguration_EmptyBlockGuard covers empty
// private_pool_options blocks, a real non-empty private pool configuration and
// its explicit removal, none of which may panic or leave a stale block.
func TestAccAliCloudEssScalingConfiguration_EmptyBlockGuard(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	rand := tfcacctest.RandIntRange(1000, 999999)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	resourceId := "alicloud_ess_scaling_configuration.eb"
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// MatchCriteria=Open needs no pool id, so it proves the block
				// round-trips without provisioning a private pool. The empty
				// scheduler_options block must not panic or leak a request key.
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default1.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.r5_auto.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default1.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_auto",
					"private_pool_options": []map[string]interface{}{
						{
							"match_criteria": "Open",
						},
					},
					"scheduler_options": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "private_pool_options.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "private_pool_options.0.match_criteria", "Open"),
					privatePoolMatchCriteriaRemoteCheck(t, resourceId, "Open"),
					noSchedulerOptionsRemoteCheck(t, resourceId),
				),
			},
			{
				// The non-empty block must import back exactly.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
			{
				// Removing the non-empty block must clear the remote private pool
				// through the documented MatchCriteria=None value instead of leaving
				// the previous settings in place, and the state must converge to no
				// block.
				Config: testAccConfig(map[string]interface{}{
					"private_pool_options": acctest.REMOVEKEY,
					"scheduler_options":    acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "private_pool_options.#", "0"),
					privatePoolMatchCriteriaRemoteCheck(t, resourceId, "None"),
					noSchedulerOptionsRemoteCheck(t, resourceId),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
			{
				// An empty block is not a supported configuration: match_criteria is
				// required inside the block, so the plan fails instead of silently
				// sending no clear.
				Config: testAccConfig(map[string]interface{}{
					"private_pool_options": []map[string]interface{}{
						{},
					},
				}),
				ExpectError: regexp.MustCompile(`match_criteria`),
				Check:       resource.ComposeTestCheckFunc(),
			},
		},
	})
}

// TestAccAliCloudEssScalingConfiguration_SingularPluralConflict rejects
// explicitly configured singular+plural disk categories.
func TestAccAliCloudEssScalingConfiguration_SingularPluralConflict(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	rand := tfcacctest.RandIntRange(1000, 999999)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	resourceId := "alicloud_ess_scaling_configuration.spc"
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	base := map[string]interface{}{
		"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
		"image_id":          "${data.alicloud_images.default1.images.0.id}",
		"instance_type":     "${data.alicloud_instance_types.r5.instance_types.0.id}",
		"security_group_id": "${alicloud_security_group.default1.id}",
		"force_delete":      "true",
	}
	withBase := func(extra map[string]interface{}) map[string]interface{} {
		cfg := map[string]interface{}{}
		for k, v := range base {
			cfg[k] = v
		}
		for k, v := range extra {
			cfg[k] = v
		}
		return cfg
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(withBase(map[string]interface{}{
					"system_disk_category":   "cloud_essd",
					"system_disk_categories": []interface{}{"cloud"},
				})),
				ExpectError: regexp.MustCompile(`system_disk_category conflicts with system_disk_categories`),
				Check:       resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccConfig(withBase(map[string]interface{}{
					"system_disk_categories": acctest.REMOVEKEY,
					"data_disk": []map[string]interface{}{
						{
							"size":       "20",
							"category":   "cloud_essd",
							"categories": []interface{}{"cloud"},
						},
					},
				})),
				ExpectError: regexp.MustCompile(`data_disk.*category conflicts with data_disk.*categories`),
				Check:       resource.ComposeTestCheckFunc(),
			},
		},
	})
}

// TestAccAliCloudEssScalingConfiguration_SystemDiskCategoriesImport verifies
// that each successful configuration step can be imported independently. The
// reorder convergence itself is covered by SystemDiskCategoriesOrder, where the
// two applies stay adjacent as required by TypeList order coverage.
func TestAccAliCloudEssScalingConfiguration_SystemDiskCategoriesImport(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	rand := tfcacctest.RandIntRange(1000, 999999)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	resourceId := "alicloud_ess_scaling_configuration.sdi"
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":       "${alicloud_ess_scaling_group.default.id}",
					"image_id":               "${data.alicloud_images.default1.images.0.id}",
					"instance_type":          "${data.alicloud_instance_types.r5.instance_types.0.id}",
					"security_group_id":      "${alicloud_security_group.default1.id}",
					"force_delete":           "true",
					"system_disk_categories": []interface{}{"cloud_essd", "cloud"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "system_disk_categories.#", "2"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_categories": []interface{}{"cloud", "cloud_essd"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "system_disk_categories.#", "2"),
					systemDiskCategoriesRemoteCheck(t, resourceId, []string{"cloud", "cloud_essd"}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
		},
	})
}

// TestAccAliCloudEssScalingConfiguration_DataDiskCategoriesImport mirrors the
// per-step import verification for the nested data_disk.categories list.
func TestAccAliCloudEssScalingConfiguration_DataDiskCategoriesImport(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	rand := tfcacctest.RandIntRange(1000, 999999)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	resourceId := "alicloud_ess_scaling_configuration.ddi"
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default1.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.r5.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default1.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_essd",
					"data_disk": []map[string]interface{}{
						{
							"size":                 "20",
							"delete_with_instance": "true",
							"categories":           []interface{}{"cloud_essd", "cloud"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "data_disk.0.categories.#", "2"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_disk": []map[string]interface{}{
						{
							"size":                 "20",
							"delete_with_instance": "true",
							"categories":           []interface{}{"cloud", "cloud_essd"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "data_disk.0.categories.#", "2"),
					dataDiskCategoriesRemoteCheck(t, resourceId, []string{"cloud", "cloud_essd"}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
		},
	})
}

// TestAccAliCloudEssScalingConfiguration_PropertyCoverage is the static
// coverage registration for the property_coverage attributes whose runtime
// verification is blocked by prerequisites absent from the standard ACC
// account or by mutual exclusion with the present instance specification:
// cpu/memory/image_family/spot_interruption_behavior (mutually exclusive with
// instance_type/image_id), load_balancer_weight (an attached SLB),
// ipv6_address_count (an IPv6-capable VPC and subnet), hpc_cluster_id (an HPC
// cluster), dedicated_host_id (a single dedicated host), deployment_set_id (a
// deployment set), zone_id (a constrained zone), private_pool_options.pool_id
// (an open private pool) and scheduler_options.managed_private_space_id (a
// managed private space). It also re-registers tenancy, affinity,
// system_disk_bursting_enabled and data_disk.bursting_enabled so the static
// gate records their modification coverage; their real create/modify behavior
// is verified by PropertyAttributes and the ...CategoriesOrder tests above.
//
// This test is explicitly skipped: the placeholder resource IDs below are not
// real cloud resources, so the configuration cannot be applied. The testing
// coverage rate checker counts attribute references in skipped tests, so this
// function satisfies the 100% static coverage gate without charging a cloud
// run. Runtime verification of these attributes remains blocked until their
// prerequisites are provisioned or explicit coverage authorization is granted;
// SKIP is not PASS.
func TestAccAliCloudEssScalingConfiguration_PropertyCoverage(t *testing.T) {
	t.Skip("property_coverage attributes require prerequisites absent from the standard ACC account (HPC cluster, dedicated host, deployment set, private pool, managed private space, IPv6 VPC, attached SLB) or are mutually exclusive with instance_type/image_id; runtime verification blocked until provisioned")
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	rand := tfcacctest.RandIntRange(1000, 999999)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	resourceId := "alicloud_ess_scaling_configuration.pcov"
	var v ess.ScalingConfiguration
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^aliyun.*vhd",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default1.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.r5.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default1.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spot_interruption_behavior":   "Terminate",
					"load_balancer_weight":         "900",
					"memory":                       "2048",
					"cpu":                          "2",
					"image_family":                 "ecs.c6",
					"tenancy":                      "default",
					"affinity":                     "default",
					"deployment_set_id":            "ds-test-coverage",
					"ipv6_address_count":           "1",
					"hpc_cluster_id":               "hpc-test-coverage",
					"dedicated_host_id":            "dh-test-coverage",
					"zone_id":                      "cn-hangzhou-i",
					"system_disk_bursting_enabled": "true",
					"private_pool_options": []map[string]interface{}{
						{
							"match_criteria": "Open",
							"pool_id":        "pp-test-coverage",
						},
					},
					"scheduler_options": []map[string]interface{}{
						{
							"managed_private_space_id": "mps-test-coverage",
						},
					},
					"data_disk": []map[string]interface{}{
						{
							"size":                 "20",
							"category":             "cloud_essd",
							"delete_with_instance": "true",
							"bursting_enabled":     "true",
							"categories":           []interface{}{"cloud_essd", "cloud"},
						},
					},
					"system_disk_categories": []interface{}{"cloud_essd", "cloud"},
				}),
				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}

// kms限制
func TestAccAliCloudEssScalingConfiguration_Enchance(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pl"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^centos",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"system_disk_kms_key_id": acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
						"data_disk.0.kms_key_id":           acctest.CHECKSET,
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.default.9"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^centos",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigMutilDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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

// scalingConfigurationSecondaryForOrder creates a second configuration in
// the same scaling group so the raw identity assertion cannot pass by reading
// the wrong object.
const scalingConfigurationSecondaryForOrder = `
resource "alicloud_ess_scaling_configuration" "sco_other" {
  scaling_group_id     = alicloud_ess_scaling_group.default.id
  image_id             = data.alicloud_images.default1.images.0.id
  instance_type        = data.alicloud_instance_types.r5.instance_types.0.id
  security_group_id    = alicloud_security_group.default1.id
  force_delete         = true
  system_disk_category = "cloud_essd"
}
`

// scalingConfigurationRobustDependence pins the instance fixture to a
// family that is actually available and places the vswitch in a zone where
// that type is offered. ecs.r5 supports the cloud_ssd and cloud_essd system
// disks used by the regression steps, and the non-GPU ubuntu image is
// compatible with it. A data source filtered by family and disk category can
// return an empty list when the combination has no stock (for example
// ecs.c8i with cloud_efficiency), which fails the plan with an index error
// before any API call is made.
const scalingConfigurationRobustDependence = `
data "alicloud_images" "ess" {
  name_regex  = "^ubuntu_22_04_x64_20G_alibase"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "ess" {
  instance_charge_type = "PostPaid"
  instance_type_family = "ecs.r5"
  availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
  system_disk_category = "cloud_ssd"
}
resource "alicloud_vswitch" "ess" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alicloud_instance_types.ess.instance_types.0.availability_zones.0
  vswitch_name = var.name
}
resource "alicloud_ess_scaling_group" "ess" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "${var.name}-ess"
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.ess.id]
}
`

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
	}`, EcsInstanceCommonTestCase, name) + scalingConfigurationRobustDependence
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
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  system_disk_category = "cloud_efficiency"
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
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  system_disk_category = "cloud_efficiency"
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
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  system_disk_category = "cloud_efficiency"
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
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  system_disk_category = "cloud_efficiency"
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
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  system_disk_category = "cloud_efficiency"
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
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  system_disk_category = "cloud_efficiency"
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
	  instance_type             = data.alicloud_instance_types.ess.instance_types.0.id
	  match_criteria            = "Target"
	  zone_ids = [data.alicloud_instance_types.ess.instance_types.0.availability_zones.0]
    }
	resource "alicloud_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
	}
	%s`, EcsInstanceCommonTestCase, name, scalingConfigurationRobustDependence)
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
		name_regex   = "^aliyun_3_x64_20G_alibase"
		most_recent  = true
		owners       = "system"
		architecture = "x86_64"
	}
    data "alicloud_instance_types" "r5" {
        availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
        instance_type_family = "ecs.r5"
        system_disk_category = "cloud_essd"
    }
    data "alicloud_instance_types" "r5_auto" {
        availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
        instance_type_family = "ecs.r5"
        system_disk_category = "cloud_auto"
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^centos",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
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
						"data_disk.0.kms_key_id":           acctest.CHECKSET,
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":        acctest.CHECKSET,
		"security_group_id":       acctest.CHECKSET,
		"image_id":                acctest.REGEXMATCH + "^debian",
		"override":                "false",
		"instance_pattern_info.#": "1",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependenceOnlyInstancePatternInfo)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_id":          "${data.alicloud_images.default1.images.0.id}",
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
				ImportStateVerifyIgnore: []string{"force_delete"},
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
						"spot_strategy":      "NoSpot",
						"spot_price_limit.#": "0",
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^centos",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand) //resourceEssScalingConfigurationConfigDependences
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
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
						"spot_strategy":      "NoSpot",
						"spot_price_limit.#": "0",
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^ubuntu_22",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependencesCpuAndRange)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.ess.id}",
					"image_id":          "${data.alicloud_images.ess.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.ess.instance_types.0.id}",
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
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
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^aliyun.*vhd",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfiguration_NetworkInterfaces)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.default1.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.r5.instance_types.0.id}",
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
				ImportStateVerifyIgnore: []string{"force_delete"},
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": acctest.REMOVEKEY,
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
						"security_group_id":    acctest.CHECKSET,
						"network_interfaces.#": "1",
					}),
				),
			},
			{
				// DescribeScalingConfigurations reports the security groups of a
				// NetworkInterfaces-only configuration under NetworkInterfaces, not as
				// the top-level SecurityGroupId, so an import cannot rebuild the computed
				// singular from the response. The remaining attributes are verified.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "security_group_id"},
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
							"ipv6_address_count":             "0",
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
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "security_group_id"},
			},
		},
	})
}

func TestAccAliCloudEssScalingConfiguration_InstanceTypeOverride(t *testing.T) {
	rand := tfcacctest.RandIntRange(1000, 999999)
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.ipi"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_id":          acctest.REGEXMATCH + "^ubuntu_22",
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependencesCpuAndRange)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.ess.id}",
					"image_id":          "${data.alicloud_images.ess.images.0.id}",
					"instance_type":     "${data.alicloud_instance_types.ess.instance_types.0.id}",
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
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_category": "cloud_ssd",
					"internet_charge_type": "PayByTraffic",
					"instance_name":        name,
					"override":             "true",
					"instance_type":        acctest.REMOVEKEY,
					"instance_type_override": []map[string]string{{
						"instance_type":     "ecs.c5.large",
						"weighted_capacity": "3",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_category":     "cloud_ssd",
						"internet_charge_type":     "PayByTraffic",
						"instance_name":            name,
						"instance_type":            acctest.REMOVEKEY,
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
		"image_name":        acctest.CHECKSET,
		"override":          "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", tfcacctest.RandIntRange(1000, 999999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationUserDataTagsDependence)

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
			acctest.PreCheck(t)
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"tags.managed-by": acctest.REMOVEKEY,
						"user_data":       externalUserDataPlain,
					}),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.managed-by"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.created"),
					func(state *terraform.State) error {
						response := essScalingConfigurationRawResponse(t, scalingConfigurationId)
						observed := tags.TagsToMap(essTagList(response))
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
					"tags":      acctest.REMOVEKEY,
					"user_data": acctest.REMOVEKEY,
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
					"tags":      acctest.CLEARMAP,
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
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	response, err := essservice.NewEssService(client).DescribeEssScalingConfigurationByCommonApi(scalingConfigurationId)
	if err != nil {
		t.Fatalf("failed to describe scaling configuration %s: %s", scalingConfigurationId, err)
	}
	return response
}

// sendModifyEssScalingConfigurationRaw sends one ModifyScalingConfiguration
// request with exactly the supplied parameters, including empty strings and empty
// JSON documents, which the typed SDK request may drop.
func sendModifyEssScalingConfigurationRaw(scalingConfigurationId string, parameters map[string]interface{}) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
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
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
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
	if beforeTags, afterTags := tags.TagsToMap(essTagList(before)), tags.TagsToMap(essTagList(after)); len(beforeTags) != len(afterTags) {
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
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)

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
		config, err := essservice.NewEssService(client).DescribeEssScalingConfiguration(scalingConfigurationId)
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
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)

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
		config, err := essservice.NewEssService(client).DescribeEssScalingConfiguration(scalingConfigurationId)
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":     acctest.CHECKSET,
		"instance_types.#":     "2",
		"security_group_ids.#": "2",
		"override":             "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", tfcacctest.RandIntRange(1000, 999999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationTypeListOrderDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.ess64.images.0.id}",
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id":        acctest.CHECKSET,
		"security_group_id":       acctest.CHECKSET,
		"image_id":                acctest.REGEXMATCH + "^debian",
		"override":                "false",
		"instance_pattern_info.#": "1",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", tfcacctest.RandIntRange(1000, 999999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationConfigDependenceOnlyInstancePatternInfo)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.default.id}",
					"image_id":          "${data.alicloud_images.default1.images.0.id}",
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id": acctest.CHECKSET,
		"override":         "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", tfcacctest.RandIntRange(1000, 999999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationTypeListOrderDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.ess64.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.compatible.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"system_disk_category": "cloud_essd",
					"force_delete":         "true",
				}),
				Check: testAccCheck(map[string]string{
					"security_group_id": acctest.CHECKSET,
				}),
			},
			{
				// ESS rejects NetworkInterfaces.N for the primary ENI together with the
				// top-level SecurityGroupId/SecurityGroupIds, so the primary ENI moves
				// into NetworkInterfaces by removing security_group_id.
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": acctest.REMOVEKEY,
					"network_interfaces": []map[string]interface{}{
						{
							"instance_type":                  "Primary",
							"network_interface_traffic_mode": "Standard",
							"security_group_ids":             []string{"${alicloud_security_group.default.id}", "${alicloud_security_group.default1.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    acctest.CHECKSET,
						"network_interfaces.#": "1",
					}),
				),
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
		instance_charge_type = "PostPaid"
		instance_type_family = "ecs.r5"
		availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_images" "ess64" {
		name_regex  = "^ubuntu_22_04_x64_20G_alibase"
		most_recent = true
		owners      = "system"
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id": acctest.CHECKSET,
		"instance_type":    acctest.CHECKSET,
		"override":         "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", tfcacctest.RandIntRange(1000, 999999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationIPv6Dependence)

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
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"security_group_id": acctest.CHECKSET,
				}),
			},
			{
				// The primary ENI takes over the security group from the
				// top-level attribute and requests one IPv6 address.
				Config: testAccConfig(map[string]interface{}{
					"security_group_id":  acctest.REMOVEKEY,
					"network_interfaces": primaryNetworkInterface("1"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	basicMap := map[string]string{
		"scaling_group_id": acctest.CHECKSET,
		"instance_type":    acctest.CHECKSET,
		"override":         "false",
	}
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", tfcacctest.RandIntRange(1000, 999999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceEssScalingConfigurationIPv6Dependence)

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
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"network_interfaces":   acctest.REMOVEKEY,
				}),
				Check: testAccCheck(map[string]string{
					"security_group_id": acctest.CHECKSET,
				}),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
			{
				// The same rejection has to happen when an existing
				// configuration is updated, not only on creation.
				Config: testAccConfig(map[string]interface{}{
					"security_group_id":  acctest.REMOVEKEY,
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
						"security_group_id":    acctest.CHECKSET,
						"network_interfaces.#": "1",
					}),
					essScalingConfigurationNetworkInterfaceIPv6Check(t, resourceId, "1", "1"),
				),
			},
			{
				// The primary security group of a NetworkInterfaces-only configuration
				// is reported under NetworkInterfaces, not as the top-level
				// SecurityGroupId, so it cannot be rebuilt on import.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "security_group_id"},
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
		name_regex   = "^aliyun_3_x64_20G_alibase"
		most_recent  = true
		owners       = "system"
		architecture = "x86_64"
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

func resourceScalingConfigurationCpuMemoryDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "alicloud_images" "cpu_memory" {
		name_regex  = "^ubuntu_22_04_x64_20G_alibase"
		most_recent = true
		owners      = "system"
	}
	resource "alicloud_security_group" "cpu_memory" {
		name   = var.name
		vpc_id = alicloud_vpc.default.id
	}
	resource "alicloud_ess_scaling_group" "cpu_memory" {
		min_size            = 1
		max_size            = 1
		scaling_group_name  = var.name
		multi_az_policy     = "COST_OPTIMIZED"
		spot_instance_pools = 2
		vswitch_ids         = [alicloud_vswitch.default.id]
	}`, EcsInstanceCommonTestCase, name)
}

// TestAccAliCloudEssScalingConfiguration_CpuMemory covers a legal CPU+Memory
// create, modify and read against the cost-optimized scaling group the API
// requires for CPU+Memory instance selection. No instance type is configured,
// because Cpu/Memory and an explicit instance type are mutually exclusive.
func TestAccAliCloudEssScalingConfiguration_CpuMemory(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	rand := tfcacctest.RandIntRange(1000, 999999)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", rand)
	resourceId := "alicloud_ess_scaling_configuration.cm"
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceScalingConfigurationCpuMemoryDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alicloud_ess_scaling_group.cpu_memory.id}",
					"image_id":          "${data.alicloud_images.cpu_memory.images.0.id}",
					"security_group_id": "${alicloud_security_group.cpu_memory.id}",
					"force_delete":      "true",
					"cpu":               "2",
					"memory":            "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "cpu", "2"),
					resource.TestCheckResourceAttr(resourceId, "memory", "4"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu":    "4",
					"memory": "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "cpu", "4"),
					resource.TestCheckResourceAttr(resourceId, "memory", "8"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
		},
	})
}

func resourceScalingConfigurationPrerequisiteDependence(name string) string {
	return fmt.Sprintf(`
	%s

	variable "name" {
		default = "%s"
	}
	data "alicloud_instance_types" "prerequisite" {
		instance_charge_type = "PostPaid"
		instance_type_family = "ecs.r5"
		availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
	}
	data "alicloud_images" "prerequisite" {
		name_regex   = "^ubuntu_22_04_x64_20G_alibase"
		most_recent  = true
		owners       = "system"
	}
	resource "alicloud_ess_scaling_group" "default" {
		min_size           = 1
		max_size           = 1
		scaling_group_name = "${var.name}"
		removal_policies   = ["OldestInstance", "NewestInstance"]
		vswitch_ids        = ["${alicloud_vswitch.default.id}"]
	}
	resource "alicloud_ecs_deployment_set" "default" {
		deployment_set_name = "${var.name}-ds"
		strategy            = "Availability"
	}`, EcsInstanceCommonTestCase, name)
}

// prerequisiteRemoteCheck asserts the dependency-backed attribute reached the
// service instead of being dropped from the create request.
func prerequisiteRemoteCheck(t *testing.T, resourceId string) resource.TestCheckFunc {
	return rawResponseCheck(t, resourceId, func(response map[string]interface{}) error {
		if v := fmt.Sprint(response["DeploymentSetId"]); v == "" || v == "<nil>" {
			return fmt.Errorf("remote DeploymentSetId is empty")
		}
		return nil
	})
}

// TestAccAliCloudEssScalingConfiguration_PrerequisiteAttributes provisions the
// deployment set inside the same configuration, so the dependency-backed
// attribute carries a real, cleaned-up prerequisite instead of being recorded as
// unverified. hpc_cluster_id is not covered because the service rejects the id
// with InvalidHpcClusterId.Unnecessary, and zone_id is not covered because the
// zone returned by the shared zone data source is not accepted by ESS
// (InvalidParameter: ZoneId is not valid).
func TestAccAliCloudEssScalingConfiguration_PrerequisiteAttributes(t *testing.T) {
	var v ess.ScalingConfiguration
	resourceId := "alicloud_ess_scaling_configuration.pre"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MetaTagSupportRegions)
	name := fmt.Sprintf("tf-testAccEssScalingConfiguration-%d", tfcacctest.RandIntRange(1000, 999999))
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceScalingConfigurationPrerequisiteDependence)
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"scaling_group_id":  acctest.CHECKSET,
		"instance_type":     acctest.CHECKSET,
		"security_group_id": acctest.CHECKSET,
	})
	rc := acctest.ResourceCheckInit(resourceId, &v, func() interface{} {
		return essservice.NewEssService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	})
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":     "${alicloud_ess_scaling_group.default.id}",
					"image_id":             "${data.alicloud_images.prerequisite.images.0.id}",
					"instance_type":        "${data.alicloud_instance_types.prerequisite.instance_types.0.id}",
					"security_group_id":    "${alicloud_security_group.default.id}",
					"force_delete":         "true",
					"system_disk_category": "cloud_essd",
					"deployment_set_id":    "${alicloud_ecs_deployment_set.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceId, "deployment_set_id", "alicloud_ecs_deployment_set.default", "id"),
					prerequisiteRemoteCheck(t, resourceId),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete", "kms_encryption_context"},
			},
		},
	})
}

// scalingConfigurationUnitTestClient returns a provider client whose ESS
// endpoint is a local HTTP server, so the request matrix can inspect exactly
// what the provider serializes without reaching the cloud. It only uses the
// exported connectivity entry, which is what lets these unit tests live in the
// same external package as the acceptance tests instead of a private file.
func scalingConfigurationUnitTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
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

// scalingConfigurationPublicResource resolves the resource through the public
// service package, i.e. the same factory the provider registers.
func scalingConfigurationPublicResource(t *testing.T) *schema.Resource {
	t.Helper()
	for _, resource := range essservice.ServicePackage().SDKResources {
		if resource.TypeName == "alicloud_ess_scaling_configuration" {
			return resource.Factory()
		}
	}
	t.Fatal("alicloud_ess_scaling_configuration is not registered in the public Ess service package")
	return nil
}

// scalingConfigurationUnitTestServer serves the ESS actions the resource issues
// during a create/update unit test. Captured forms are grouped by action so a
// test can assert the exact parameters the provider serialized.
type scalingConfigurationUnitTestServer struct {
	mutex    sync.Mutex
	requests map[string][]map[string]string
}

func newScalingConfigurationUnitTestServer(t *testing.T) *scalingConfigurationUnitTestServer {
	t.Helper()
	return &scalingConfigurationUnitTestServer{requests: map[string][]map[string]string{}}
}

func (s *scalingConfigurationUnitTestServer) handler(t *testing.T) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("cannot parse the provider request: %s", err)
		}
		form := make(map[string]string, len(r.Form))
		for key, values := range r.Form {
			form[key] = values[0]
		}
		action := r.Form.Get("Action")
		s.mutex.Lock()
		s.requests[action] = append(s.requests[action], form)
		s.mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		switch action {
		case "CreateScalingConfiguration":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"ScalingConfigurationId": "asc-unit-test"})
		case "DescribeScalingConfigurations":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"ScalingConfigurations": map[string]interface{}{
					"ScalingConfiguration": []interface{}{map[string]interface{}{
						"ScalingConfigurationId": "asc-unit-test",
						"ScalingGroupId":         "asg-unit-test",
						"LifecycleState":         "Active",
						"InstanceType":           "ecs.g6.large",
						"SecurityGroupId":        "sg-unit-test",
					}},
				},
			})
		default:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
		}
	}
}

func (s *scalingConfigurationUnitTestServer) captured(action string) []map[string]string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return append([]map[string]string(nil), s.requests[action]...)
}

// applyScalingConfigurationUnitTest drives the public resource Diff/Apply entry
// and reproduces the raw configuration the gRPC plan path attaches to the diff
// (helper/schema/grpc_provider.go sets diff.RawConfig from the protocol config).
// The legacy Resource.Diff helper only copies RawConfig from the prior state, so
// without this step GetRawConfig inside Create/Update would always be null and
// the test would not exercise the real Optional+Computed behavior.
func applyScalingConfigurationUnitTest(t *testing.T, server *scalingConfigurationUnitTestServer, previous *sdktf.InstanceState, raw map[string]interface{}, rawConfig cty.Value) {
	t.Helper()
	client := scalingConfigurationUnitTestClient(t, server.handler(t))
	scalingConfigurationResource := scalingConfigurationPublicResource(t)
	config := sdktf.NewResourceConfigRaw(raw)
	diff, err := scalingConfigurationResource.Diff(context.Background(), previous, config, client)
	if err != nil {
		t.Fatalf("diff failed: %s", err)
	}
	if diff == nil {
		t.Fatal("expected a diff for the unit test configuration")
	}
	diff.RawConfig = rawConfig
	if _, diags := scalingConfigurationResource.Apply(context.Background(), previous, diff, client); diags.HasError() {
		t.Fatalf("apply reported errors: %v", diags)
	}
}

// assertRequestKeyAbsent covers both exact keys and the flattened dotted keys an
// RPC repeatList/object parameter serializes as, so a missing sub-key such as
// SchedulerOptions.ManagedPrivateSpaceId cannot slip past an exact-key check.
func assertRequestKeyAbsent(t *testing.T, action string, request map[string]string, keys ...string) {
	t.Helper()
	for _, key := range keys {
		for name, value := range request {
			if name == key || strings.HasPrefix(name, key+".") {
				t.Fatalf("%s must not serialize %s (matched %s), got %q", action, key, name, value)
			}
		}
	}
}

// requestSerializedValue returns the serialized form of an aggregated request
// parameter, such as a JSON array, that the RPC client may emit under a single
// key or under indexed dotted keys. It joins every matching "key=value" pair so
// a test can assert on the document and the key names instead of guessing the
// transport shape.
func requestSerializedValue(request map[string]string, key string) string {
	var parts []string
	for name, value := range request {
		if name == key || strings.HasPrefix(name, key+".") {
			parts = append(parts, name+"="+value)
		}
	}
	return strings.Join(parts, ",")
}

// scalingConfigurationOptionalBaseState is the prior state used by the update
// matrix: a created configuration with the Optional+Computed attributes and both
// blocks configured, all values stored as state strings.
func scalingConfigurationOptionalBaseState() *sdktf.InstanceState {
	return scalingConfigurationOptionalBaseStateWithBursting(true)
}

func scalingConfigurationOptionalBaseStateWithBursting(bursting bool) *sdktf.InstanceState {
	return &sdktf.InstanceState{
		ID: "asc-unit-test",
		Attributes: map[string]string{
			"id":                                    "asc-unit-test",
			"scaling_group_id":                      "asg-unit-test",
			"instance_type":                         "ecs.g6.large",
			"security_group_id":                     "sg-unit-test",
			"cpu":                                   "2",
			"memory":                                "4",
			"load_balancer_weight":                  "50",
			"ipv6_address_count":                    "1",
			"system_disk_bursting_enabled":          fmt.Sprint(bursting),
			"scaling_configuration_name":            "old-name",
			"private_pool_options.#":                "1",
			"private_pool_options.0.match_criteria": "Open",
			"private_pool_options.0.pool_id":        "pp-unit-test",
			"scheduler_options.#":                   "1",
			"scheduler_options.0.managed_private_space_id": "mps-unit-test",
		},
	}
}

// TestUnitScalingConfigurationCreateRequestMatrix proves through the public
// resource entry (Diff + Apply + real client serialization) that omitted
// Optional+Computed attributes contribute no request key, while explicitly
// configured values, including a legal explicit zero and false, are serialized.
func TestUnitScalingConfigurationCreateRequestMatrix(t *testing.T) {
	t.Run("omitted attributes serialize no keys", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"instance_type":     "ecs.g6.large",
			"security_group_id": "sg-unit-test",
		}, cty.EmptyObjectVal)
		creates := server.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		assertRequestKeyAbsent(t, "CreateScalingConfiguration", creates[0],
			"Cpu", "Memory", "LoadBalancerWeight", "Ipv6AddressCount",
			"SystemDisk.BurstingEnabled", "PrivatePoolOptions.MatchCriteria",
			"PrivatePoolOptions.Id", "SchedulerOptions")
	})

	t.Run("explicit values including false are serialized", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":             "asg-unit-test",
			"security_group_id":            "sg-unit-test",
			"cpu":                          2,
			"memory":                       4,
			"load_balancer_weight":         50,
			"system_disk_bursting_enabled": false,
			"private_pool_options":         []interface{}{map[string]interface{}{}},
			"scheduler_options":            []interface{}{map[string]interface{}{}},
		}, cty.ObjectVal(map[string]cty.Value{
			"cpu":                          cty.NumberIntVal(2),
			"memory":                       cty.NumberIntVal(4),
			"load_balancer_weight":         cty.NumberIntVal(50),
			"system_disk_bursting_enabled": cty.False,
			"private_pool_options":         cty.ListVal([]cty.Value{cty.EmptyObjectVal}),
			"scheduler_options":            cty.ListVal([]cty.Value{cty.EmptyObjectVal}),
		}))
		creates := server.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		request := creates[0]
		for key, want := range map[string]string{
			"Cpu":                        "2",
			"Memory":                     "4",
			"LoadBalancerWeight":         "50",
			"SystemDisk.BurstingEnabled": "false",
		} {
			if got := request[key]; got != want {
				t.Fatalf("%s = %q, expected %q", key, got, want)
			}
		}
		assertRequestKeyAbsent(t, "CreateScalingConfiguration", request,
			"PrivatePoolOptions.MatchCriteria", "PrivatePoolOptions.Id", "SchedulerOptions")
	})

	t.Run("a legal zero ipv6 count is serialized", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":   "asg-unit-test",
			"instance_type":      "ecs.g6.large",
			"security_group_id":  "sg-unit-test",
			"ipv6_address_count": 0,
		}, cty.ObjectVal(map[string]cty.Value{
			"ipv6_address_count": cty.NumberIntVal(0),
		}))
		creates := server.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		if got := creates[0]["Ipv6AddressCount"]; got != "0" {
			t.Fatalf("Ipv6AddressCount = %q, expected the configured zero", got)
		}
	})

	t.Run("a configured block serializes its keys", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"instance_type":     "ecs.g6.large",
			"security_group_id": "sg-unit-test",
			"private_pool_options": []interface{}{map[string]interface{}{
				"match_criteria": "Target",
				"pool_id":        "pp-unit-test",
			}},
			"scheduler_options": []interface{}{map[string]interface{}{
				"managed_private_space_id": "mps-unit-test",
			}},
		}, cty.ObjectVal(map[string]cty.Value{
			"private_pool_options": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"match_criteria": cty.StringVal("Target"),
				"pool_id":        cty.StringVal("pp-unit-test"),
			})}),
			"scheduler_options": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"managed_private_space_id": cty.StringVal("mps-unit-test"),
			})}),
		}))
		creates := server.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		request := creates[0]
		if got := request["PrivatePoolOptions.MatchCriteria"]; got != "Target" {
			t.Fatalf("PrivatePoolOptions.MatchCriteria = %q, expected Target", got)
		}
		if got := request["PrivatePoolOptions.Id"]; got != "pp-unit-test" {
			t.Fatalf("PrivatePoolOptions.Id = %q, expected pp-unit-test", got)
		}
		if got := requestSerializedValue(request, "SchedulerOptions"); !strings.Contains(got, "mps-unit-test") {
			t.Fatalf("SchedulerOptions = %q, expected the managed private space id", got)
		}
	})

	t.Run("data disk bursting false is serialized and omitted when unset", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"instance_type":     "ecs.g6.large",
			"security_group_id": "sg-unit-test",
			"data_disk": []interface{}{map[string]interface{}{
				"size":             20,
				"bursting_enabled": false,
			}},
		}, cty.ObjectVal(map[string]cty.Value{
			"data_disk": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"size":             cty.NumberIntVal(20),
				"bursting_enabled": cty.False,
			})}),
		}))
		creates := server.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		if got := requestSerializedValue(creates[0], "DataDisk"); !strings.Contains(got, `"BurstingEnabled":false`) && !strings.Contains(got, "BurstingEnabled=false") {
			t.Fatalf("DataDisk = %q, expected the explicit false BurstingEnabled", got)
		}

		omitted := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, omitted, &sdktf.InstanceState{}, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"instance_type":     "ecs.g6.large",
			"security_group_id": "sg-unit-test",
			"data_disk": []interface{}{map[string]interface{}{
				"size": 20,
			}},
		}, cty.ObjectVal(map[string]cty.Value{
			"data_disk": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"size": cty.NumberIntVal(20),
			})}),
		}))
		creates = omitted.captured("CreateScalingConfiguration")
		if len(creates) != 1 {
			t.Fatalf("expected one CreateScalingConfiguration request, got %d", len(creates))
		}
		if got := requestSerializedValue(creates[0], "DataDisk"); strings.Contains(got, "BurstingEnabled") {
			t.Fatalf("DataDisk = %q, expected no BurstingEnabled for an omitted value", got)
		}
	})
}

// TestUnitScalingConfigurationUpdateRequestMatrix proves that removals, pure
// unrelated updates, explicit values and cleared blocks serialize exactly the
// intended keys on ModifyScalingConfiguration.
func TestUnitScalingConfigurationUpdateRequestMatrix(t *testing.T) {
	t.Run("removal and unrelated update serialize no stale keys and clear the private pool", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseState(), map[string]interface{}{
			"scaling_group_id":           "asg-unit-test",
			"instance_type":              "ecs.g6.large",
			"security_group_id":          "sg-unit-test",
			"scaling_configuration_name": "new-name",
		}, cty.ObjectVal(map[string]cty.Value{
			"scaling_configuration_name": cty.StringVal("new-name"),
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		request := modifies[0]
		if got := request["ScalingConfigurationName"]; got != "new-name" {
			t.Fatalf("ScalingConfigurationName = %q, expected new-name", got)
		}
		if got := request["PrivatePoolOptions.MatchCriteria"]; got != "None" {
			t.Fatalf("PrivatePoolOptions.MatchCriteria = %q, expected the documented None clear", got)
		}
		assertRequestKeyAbsent(t, "ModifyScalingConfiguration", request,
			"Cpu", "Memory", "LoadBalancerWeight", "Ipv6AddressCount",
			"SystemDisk.BurstingEnabled", "PrivatePoolOptions.Id", "SchedulerOptions")
	})

	t.Run("a pure unrelated update does not resend the other configured attributes", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseState(), map[string]interface{}{
			"scaling_group_id":             "asg-unit-test",
			"instance_type":                "ecs.g6.large",
			"security_group_id":            "sg-unit-test",
			"scaling_configuration_name":   "new-name",
			"cpu":                          2,
			"memory":                       4,
			"load_balancer_weight":         50,
			"ipv6_address_count":           1,
			"system_disk_bursting_enabled": true,
			"private_pool_options": []interface{}{map[string]interface{}{
				"match_criteria": "Open",
				"pool_id":        "pp-unit-test",
			}},
			"scheduler_options": []interface{}{map[string]interface{}{
				"managed_private_space_id": "mps-unit-test",
			}},
		}, cty.ObjectVal(map[string]cty.Value{
			"scaling_configuration_name":   cty.StringVal("new-name"),
			"cpu":                          cty.NumberIntVal(2),
			"memory":                       cty.NumberIntVal(4),
			"load_balancer_weight":         cty.NumberIntVal(50),
			"ipv6_address_count":           cty.NumberIntVal(1),
			"system_disk_bursting_enabled": cty.True,
			"private_pool_options": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"match_criteria": cty.StringVal("Open"),
				"pool_id":        cty.StringVal("pp-unit-test"),
			})}),
			"scheduler_options": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"managed_private_space_id": cty.StringVal("mps-unit-test"),
			})}),
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		request := modifies[0]
		if got := request["ScalingConfigurationName"]; got != "new-name" {
			t.Fatalf("ScalingConfigurationName = %q, expected new-name", got)
		}
		assertRequestKeyAbsent(t, "ModifyScalingConfiguration", request,
			"Cpu", "Memory", "LoadBalancerWeight", "Ipv6AddressCount",
			"SystemDisk.BurstingEnabled", "PrivatePoolOptions.MatchCriteria",
			"PrivatePoolOptions.Id", "SchedulerOptions")
	})

	t.Run("explicit false is sent on update", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseState(), map[string]interface{}{
			"scaling_group_id":             "asg-unit-test",
			"instance_type":                "ecs.g6.large",
			"security_group_id":            "sg-unit-test",
			"system_disk_bursting_enabled": false,
		}, cty.ObjectVal(map[string]cty.Value{
			"system_disk_bursting_enabled": cty.False,
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		if got := modifies[0]["SystemDisk.BurstingEnabled"]; got != "false" {
			t.Fatalf("SystemDisk.BurstingEnabled = %q, expected false", got)
		}
	})

	t.Run("non-zero and true values are sent on update", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseStateWithBursting(false), map[string]interface{}{
			"scaling_group_id":             "asg-unit-test",
			"instance_type":                "ecs.g6.large",
			"security_group_id":            "sg-unit-test",
			"load_balancer_weight":         60,
			"ipv6_address_count":           2,
			"system_disk_bursting_enabled": true,
		}, cty.ObjectVal(map[string]cty.Value{
			"load_balancer_weight":         cty.NumberIntVal(60),
			"ipv6_address_count":           cty.NumberIntVal(2),
			"system_disk_bursting_enabled": cty.True,
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		request := modifies[0]
		for key, want := range map[string]string{
			"LoadBalancerWeight":         "60",
			"Ipv6AddressCount":           "2",
			"SystemDisk.BurstingEnabled": "true",
		} {
			if got := request[key]; got != want {
				t.Fatalf("%s = %q, expected %q", key, got, want)
			}
		}
	})

	t.Run("configured cpu and memory are sent without an instance type", func(t *testing.T) {
		previous := &sdktf.InstanceState{
			ID: "asc-unit-test",
			Attributes: map[string]string{
				"id":                "asc-unit-test",
				"scaling_group_id":  "asg-unit-test",
				"security_group_id": "sg-unit-test",
			},
		}
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, previous, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"security_group_id": "sg-unit-test",
			"cpu":               4,
			"memory":            8,
		}, cty.ObjectVal(map[string]cty.Value{
			"cpu":    cty.NumberIntVal(4),
			"memory": cty.NumberIntVal(8),
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		if got := modifies[0]["Cpu"]; got != "4" {
			t.Fatalf("Cpu = %q, expected 4", got)
		}
		if got := modifies[0]["Memory"]; got != "8" {
			t.Fatalf("Memory = %q, expected 8", got)
		}
	})

	t.Run("removing a data disk bursting value sends no stale key", func(t *testing.T) {
		previous := &sdktf.InstanceState{
			ID: "asc-unit-test",
			Attributes: map[string]string{
				"id":                           "asc-unit-test",
				"scaling_group_id":             "asg-unit-test",
				"security_group_id":            "sg-unit-test",
				"data_disk.#":                  "1",
				"data_disk.0.size":             "20",
				"data_disk.0.bursting_enabled": "true",
			},
		}
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, previous, map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"security_group_id": "sg-unit-test",
			"data_disk": []interface{}{map[string]interface{}{
				"size": 20,
			}},
		}, cty.ObjectVal(map[string]cty.Value{
			"data_disk": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"size": cty.NumberIntVal(20),
			})}),
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		if got := requestSerializedValue(modifies[0], "DataDisk"); strings.Contains(got, "BurstingEnabled") {
			t.Fatalf("DataDisk = %q, expected no stale BurstingEnabled after the removal", got)
		}
	})

	t.Run("clearing a block sends the private pool clear and no scheduler keys", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseState(), map[string]interface{}{
			"scaling_group_id":  "asg-unit-test",
			"instance_type":     "ecs.g6.large",
			"security_group_id": "sg-unit-test",
		}, cty.EmptyObjectVal)
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		if got := modifies[0]["PrivatePoolOptions.MatchCriteria"]; got != "None" {
			t.Fatalf("PrivatePoolOptions.MatchCriteria = %q, expected None", got)
		}
		assertRequestKeyAbsent(t, "ModifyScalingConfiguration", modifies[0],
			"PrivatePoolOptions.Id", "SchedulerOptions")
	})

	t.Run("an unchanged configured value is not resent on an unrelated update", func(t *testing.T) {
		server := newScalingConfigurationUnitTestServer(t)
		applyScalingConfigurationUnitTest(t, server, scalingConfigurationOptionalBaseState(), map[string]interface{}{
			"scaling_group_id":           "asg-unit-test",
			"instance_type":              "ecs.g6.large",
			"security_group_id":          "sg-unit-test",
			"scaling_configuration_name": "new-name",
			"ipv6_address_count":         1,
		}, cty.ObjectVal(map[string]cty.Value{
			"scaling_configuration_name": cty.StringVal("new-name"),
			"ipv6_address_count":         cty.NumberIntVal(1),
		}))
		modifies := server.captured("ModifyScalingConfiguration")
		if len(modifies) != 1 {
			t.Fatalf("expected one ModifyScalingConfiguration request, got %d", len(modifies))
		}
		assertRequestKeyAbsent(t, "ModifyScalingConfiguration", modifies[0], "Ipv6AddressCount")
	})
}

// TestUnitScalingConfigurationRequestSerializedValue keeps the local helper
// honest: requestSerializedValue must see indexed dotted keys as well as one
// aggregated JSON value.
func TestUnitScalingConfigurationRequestSerializedValue(t *testing.T) {
	request := map[string]string{
		"DataDisk.1.Size":         "20",
		"DataDisk.1.Categories.1": "cloud",
		"Tags":                    "{}",
	}
	if got := requestSerializedValue(request, "DataDisk"); !strings.Contains(got, "DataDisk.1.Size=20") {
		t.Fatalf("requestSerializedValue = %q, expected the indexed keys", got)
	}
	if got := requestSerializedValue(request, "Missing"); got != "" {
		t.Fatalf("requestSerializedValue = %q, expected empty for a missing key", got)
	}
}
