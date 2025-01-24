package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCREEInstance_Basic(t *testing.T) {
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
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"period":       "1",
					//"renew_period":   "0",
					"renewal_status": "ManualRenewal",
					"instance_type":  "Basic",
					"instance_name":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":       CHECKSET,
						"created_time": CHECKSET,
						"end_time":     CHECKSET,
						//"renew_period":   "0",
						"renewal_status": "ManualRenewal",
						"instance_name":  name,
						"instance_type":  "Basic",
						"payment_type":   "Subscription",
					}),
				),
			},
			// Currently, the API does not support sts
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"password": "YourPassword123",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
			//		"kms_encryption_context": map[string]string{
			//			"name": name,
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{}),
			//	),
			//},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "custom_oss_bucket", "password", "kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudCREEInstance_Standard(t *testing.T) {
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
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
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
			// Currently, the API does not support sts
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"password": "YourPassword123",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
			//		"kms_encryption_context": map[string]string{
			//			"name": name,
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{}),
			//	),
			//},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "custom_oss_bucket", "password", "kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudCREEInstance_Advanced(t *testing.T) {
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
			// Currently, the API does not support sts
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"password": "YourPassword123",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
			//		"kms_encryption_context": map[string]string{
			//			"name": name,
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{}),
			//	),
			//},
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

// Case 实例生命周期测试_2_new 7970_modified
func TestAccAliCloudCrInstance_basic7970_modified(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCrInstanceMap7970_modified)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-basic-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrInstanceBasicDependence7970_modified)
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
					"instance_name":      name,
					"period":             "1",
					"renewal_status":     "AutoRenewal",
					"instance_type":      "Standard",
					"payment_type":       "Subscription",
					"renew_period":       "1",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"default_oss_bucket": "false",
					"custom_oss_bucket":  "${alicloud_oss_bucket.defaultkcvHCP.bucket}",
					"image_scanner":      "ACR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":      name,
						"period":             "1",
						"renewal_status":     "AutoRenewal",
						"instance_type":      "Standard",
						"payment_type":       "Subscription",
						"renew_period":       "1",
						"resource_group_id":  CHECKSET,
						"default_oss_bucket": "false",
						"custom_oss_bucket":  CHECKSET,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_oss_bucket", "default_oss_bucket", "instance_type", "period", "image_scanner"},
			},
		},
	})
}

var AlicloudCrInstanceMap7970_modified = map[string]string{
	"status":               CHECKSET,
	"end_time":             CHECKSET,
	"create_time":          CHECKSET,
	"instance_endpoints.#": CHECKSET,
}

func AlicloudCrInstanceBasicDependence7970_modified(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ram_role" "defaultRole" {
  name = "AliyunContainerRegistryCustomizedOSSBucketRole"

  description = var.name
  document = <<EOF
{
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
                "Service": [
                    "cr.aliyuncs.com"
                ]
            }
        }
    ],
    "Version": "1"
}
  EOF
}

resource "alicloud_ram_policy" "defaultLPolicy" {
  policy_name = "AliyunContainerRegistryCustomizedOSSBucketRolePolicy"
  description = var.name

  document = <<EOF

{
    "Version": "1",
    "Statement": [
        {
            "Action": [
                "oss:GetObject",
                "oss:PutObject",
                "oss:DeleteObject",
                "oss:ListParts",
                "oss:AbortMultipartUpload",
                "oss:InitiateMultipartUpload",
                "oss:CompleteMultipartUpload",
                "oss:DeleteMultipleObjects",
                "oss:ListMultipartUploads",
                "oss:ListObjects",
                "oss:DeleteObjectVersion",
                "oss:GetObjectVersion",
                "oss:ListObjectVersions",
                "oss:PutObjectTagging",
                "oss:GetObjectTagging",
                "oss:DeleteObjectTagging"
            ],
            "Resource": [
                "acs:oss:*:*:cri-*",
                "acs:oss:*:*:cri-*/*",
                "acs:oss:*:*:${var.name}",
                "acs:oss:*:*:${var.name}/*"
            ],
            "Effect": "Allow",
            "Condition": {

            }
        },
        {
            "Action": [
                "oss:PutBucket",
                "oss:GetBucket",
                "oss:GetBucketLocation",
                "oss:PutBucketEncryption",
                "oss:GetBucketEncryption",
                "oss:PutBucketAcl",
                "oss:GetBucketAcl",
                "oss:PutBucketLogging",
                "oss:GetBucketReferer",
                "oss:PutBucketReferer",
                "oss:GetBucketLogging",
                "oss:PutBucketVersioning",
                "oss:GetBucketVersioning",
                "oss:GetBucketLifecycle",
                "oss:PutBucketLifecycle",
                "oss:DeleteBucketLifecycle",
                "oss:GetBucketTransferAcceleration"
            ],
            "Resource": [
                "acs:oss:*:*:cri-*",
                "acs:oss:*:*:cri-*/*",
                "acs:oss:*:*:${var.name}",
                "acs:oss:*:*:${var.name}/*"
            ],
            "Effect": "Allow",
            "Condition": {

            }
        },
        {
            "Effect": "Allow",
            "Action": "oss:ListBuckets",
            "Resource": [
                "acs:oss:*:*:*",
                "acs:oss:*:*:*/*"
            ],
            "Condition": {

            }
        },
        {
            "Action": [
                "vpc:DescribeVpcs"
            ],
            "Resource": "acs:vpc:*:*:vpc/*",
            "Effect": "Allow",
            "Condition": {

            }
        },
        {
            "Action": [
                "cms:QueryMetricLast",
                "cms:QueryMetricList"
            ],
            "Resource": "*",
            "Effect": "Allow"
        }
    ]
}
  EOF
}

resource "alicloud_ram_role_policy_attachment" "RolePolicyAttachment" {
  policy_type = "Custom"
  role_name   = alicloud_ram_role.defaultRole.name
  policy_name = alicloud_ram_policy.defaultLPolicy.policy_name
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oss_bucket" "defaultkcvHCP" {
	depends_on = [
			alicloud_ram_role_policy_attachment.RolePolicyAttachment]
  storage_class = "Standard"
  bucket = var.name
}


`, name)
}
