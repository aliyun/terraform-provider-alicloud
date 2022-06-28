package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCREEInstance_Basic(t *testing.T) {
	var v *cr_ee.GetInstanceResponse
	resourceId := "alicloud_cr_ee_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEEInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-basic-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEEInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "Subscription",
					"period":         "1",
					"renew_period":   "0",
					"renewal_status": "ManualRenewal",
					"instance_type":  "Basic",
					"instance_name":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":         CHECKSET,
						"created_time":   CHECKSET,
						"end_time":       CHECKSET,
						"renew_period":   "0",
						"renewal_status": "ManualRenewal",
						"instance_name":  name,
						"instance_type":  "Basic",
						"payment_type":   "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "YourPassword123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{"period", "custom_oss_bucket", "password", "kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAlicloudCREEInstance_Standard(t *testing.T) {
	var v *cr_ee.GetInstanceResponse
	resourceId := "alicloud_cr_ee_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEEInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-standard-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEEInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "Subscription",
					"period":         "1",
					"renew_period":   "0",
					"renewal_status": "ManualRenewal",
					"instance_type":  "Standard",
					"instance_name":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":         CHECKSET,
						"created_time":   CHECKSET,
						"end_time":       CHECKSET,
						"renew_period":   "0",
						"renewal_status": "ManualRenewal",
						"instance_name":  name,
						"instance_type":  "Standard",
						"payment_type":   "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "YourPassword123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{"period", "custom_oss_bucket", "password", "kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAlicloudCREEInstance_Advanced(t *testing.T) {
	var v *cr_ee.GetInstanceResponse
	resourceId := "alicloud_cr_ee_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEEInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-advanced-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEEInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "Subscription",
					"period":         "1",
					"renew_period":   "0",
					"renewal_status": "ManualRenewal",
					"instance_type":  "Advanced",
					"instance_name":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":         CHECKSET,
						"created_time":   CHECKSET,
						"end_time":       CHECKSET,
						"renew_period":   "0",
						"renewal_status": "ManualRenewal",
						"instance_name":  name,
						"instance_type":  "Advanced",
						"payment_type":   "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "YourPassword123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{"period", "custom_oss_bucket", "password", "kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func resourceCrEEInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	data "alicloud_kms_keys" "default" {
	  status = "Enabled"
	}
	resource "alicloud_kms_key" "default" {
	  count = length(data.alicloud_kms_keys.default.ids) > 0 ? 0 : 1
	  description = var.name
	  status = "Enabled"
	  pending_window_in_days = 7
	}
	
	resource "alicloud_kms_ciphertext" "default" {
	  key_id = length(data.alicloud_kms_keys.default.ids) > 0 ? data.alicloud_kms_keys.default.ids.0 : concat(alicloud_kms_key.default.*.id, [""])[0]
	  plaintext = "YourPassword1234"
	  encryption_context = {
		"name" = var.name
	  }
	}
	`, name)
}
