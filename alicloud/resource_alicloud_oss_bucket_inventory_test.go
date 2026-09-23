package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketInventory. >>> Resource test cases.

func TestAccAliCloudOssBucketInventory_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_inventory.default"
	ra := resourceAttrInit(resourceId, AliCloudOssBucketInventoryMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketInventory")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccossinv%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudOssBucketInventoryBasicDependence)
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
					"bucket":                   "${alicloud_oss_bucket.source.id}",
					"inventory_id":             name,
					"is_enabled":               "true",
					"included_object_versions": "All",
					"schedule": []map[string]interface{}{
						{
							"frequency": "Daily",
						},
					},
					"optional_fields": []map[string]interface{}{
						{
							"field": []string{"Size", "LastModifiedDate"},
						},
					},
					"destination": []map[string]interface{}{
						{
							"oss_bucket_destination": []map[string]interface{}{
								{
									"format":     "CSV",
									"account_id": "${data.alicloud_account.current.id}",
									"role_arn":   "acs:ram::${data.alicloud_account.current.id}:role/${alicloud_ram_role_policy_attachment.default.role_name}",
									"bucket":     "acs:oss:::${alicloud_oss_bucket.dest.id}",
									"prefix":     "inv/",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                                 CHECKSET,
						"inventory_id":                           name,
						"is_enabled":                             "true",
						"included_object_versions":               "All",
						"schedule.#":                             "1",
						"schedule.0.frequency":                   "Daily",
						"optional_fields.#":                      "1",
						"optional_fields.0.field.#":              "2",
						"destination.#":                          "1",
						"destination.0.oss_bucket_destination.#": "1",
						"destination.0.oss_bucket_destination.0.format":     "CSV",
						"destination.0.oss_bucket_destination.0.account_id": CHECKSET,
						"destination.0.oss_bucket_destination.0.role_arn":   CHECKSET,
						"destination.0.oss_bucket_destination.0.bucket":     CHECKSET,
						"destination.0.oss_bucket_destination.0.prefix":     "inv/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_enabled":               "false",
					"included_object_versions": "Current",
					"schedule": []map[string]interface{}{
						{
							"frequency": "Weekly",
						},
					},
					"optional_fields": []map[string]interface{}{
						{
							"field": []string{
								"Size", "LastModifiedDate", "ETag", "StorageClass",
								"IsMultipartUploaded", "EncryptionStatus",
							},
						},
					},
					"filter": []map[string]interface{}{
						{
							"prefix":                       "logs/",
							"storage_class":                "Standard",
							"lower_size_bound":             "1024",
							"upper_size_bound":             "1048576",
							"last_modify_begin_time_stamp": "1262275200",
							"last_modify_end_time_stamp":   "2000000000",
						},
					},
					"destination": []map[string]interface{}{
						{
							"oss_bucket_destination": []map[string]interface{}{
								{
									"format":     "CSV",
									"account_id": "${data.alicloud_account.current.id}",
									"role_arn":   "acs:ram::${data.alicloud_account.current.id}:role/${alicloud_ram_role.default.name}",
									"bucket":     "acs:oss:::${alicloud_oss_bucket.dest.id}",
									"prefix":     "inv-v2/",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_enabled":                                    "false",
						"included_object_versions":                      "Current",
						"schedule.0.frequency":                          "Weekly",
						"optional_fields.0.field.#":                     "6",
						"filter.#":                                      "1",
						"filter.0.prefix":                               "logs/",
						"filter.0.storage_class":                        "Standard",
						"filter.0.lower_size_bound":                     "1024",
						"filter.0.upper_size_bound":                     "1048576",
						"filter.0.last_modify_begin_time_stamp":         "1262275200",
						"filter.0.last_modify_end_time_stamp":           "2000000000",
						"destination.0.oss_bucket_destination.0.format": "CSV",
						"destination.0.oss_bucket_destination.0.prefix": "inv-v2/",
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

var AliCloudOssBucketInventoryMap = map[string]string{}

func AliCloudOssBucketInventoryBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "current" {
}

resource "alicloud_oss_bucket" "source" {
  bucket        = "${var.name}-src"
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "dest" {
  bucket        = "${var.name}-dst"
  storage_class = "Standard"
}

resource "alicloud_kms_key" "default" {
  description             = "${var.name}"
  pending_window_in_days  = 7
  key_state               = "Enabled"
}

resource "alicloud_ram_role" "default" {
  name        = "${var.name}"
  document    = <<EOF
{
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "oss.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
EOF
  description = "this is a test for bucket inventory."
  force       = true
}

resource "alicloud_ram_policy" "default" {
  name     = "${var.name}"
  document = <<EOF
{
  "Statement": [
    {
      "Action": ["oss:PutObject", "oss:AbortMultipartUpload"],
      "Effect": "Allow",
      "Resource": ["acs:oss:*:*:${alicloud_oss_bucket.dest.id}/*"]
    }
  ],
  "Version": "1"
}
EOF
  force    = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  policy_name = "${alicloud_ram_policy.default.name}"
  policy_type = "${alicloud_ram_policy.default.type}"
  role_name   = "${alicloud_ram_role.default.name}"
}

`, name)
}

// TestAccAliCloudOssBucketInventory_incremental covers the incremental_inventory
// attribute tree. It is skipped by default because OSS's Incremental Inventory
// feature is not enabled on every account; when enabled, this test creates a
// bucket inventory rule that also exports an incremental report every 600s.
func TestAccAliCloudOssBucketInventory_incremental(t *testing.T) {
	t.Skip("OSS Incremental Inventory is an account-level opt-in feature; skip by default.")
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_inventory.default"
	ra := resourceAttrInit(resourceId, AliCloudOssBucketInventoryMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketInventory")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccossinv%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudOssBucketInventoryBasicDependence)
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
					"bucket":                   "${alicloud_oss_bucket.source.id}",
					"inventory_id":             name,
					"is_enabled":               "true",
					"included_object_versions": "All",
					"schedule": []map[string]interface{}{
						{
							"frequency": "Daily",
						},
					},
					"optional_fields": []map[string]interface{}{
						{
							"field": []string{"Size", "LastModifiedDate"},
						},
					},
					"destination": []map[string]interface{}{
						{
							"oss_bucket_destination": []map[string]interface{}{
								{
									"format":     "CSV",
									"account_id": "${data.alicloud_account.current.id}",
									"role_arn":   "acs:ram::${data.alicloud_account.current.id}:role/${alicloud_ram_role_policy_attachment.default.role_name}",
									"bucket":     "acs:oss:::${alicloud_oss_bucket.dest.id}",
									"prefix":     "inv/",
									"encryption": []map[string]interface{}{
										{
											"sseoss": "",
											"ssekms": []map[string]interface{}{
												{
													"key_id": "${alicloud_kms_key.default.id}",
												},
											},
										},
									},
								},
							},
						},
					},
					"incremental_inventory": []map[string]interface{}{
						{
							"is_enabled": "true",
							"schedule": []map[string]interface{}{
								{
									"frequency": "600",
								},
							},
							"optional_fields": []map[string]interface{}{
								{
									"field": []string{"Size", "LastModifiedDate"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination.0.oss_bucket_destination.0.encryption.#":                 "1",
						"destination.0.oss_bucket_destination.0.encryption.0.ssekms.#":        "1",
						"destination.0.oss_bucket_destination.0.encryption.0.ssekms.0.key_id": CHECKSET,
						"destination.0.oss_bucket_destination.0.encryption.0.sseoss":          "",
						"incremental_inventory.#":                                             "1",
						"incremental_inventory.0.is_enabled":                                  "true",
						"incremental_inventory.0.schedule.#":                                  "1",
						"incremental_inventory.0.schedule.0.frequency":                        "600",
						"incremental_inventory.0.optional_fields.#":                           "1",
						"incremental_inventory.0.optional_fields.0.field.#":                   "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"incremental_inventory": []map[string]interface{}{
						{
							"is_enabled": "false",
							"schedule": []map[string]interface{}{
								{
									"frequency": "600",
								},
							},
							"optional_fields": []map[string]interface{}{
								{
									"field": []string{"Size", "LastModifiedDate", "ETag"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"incremental_inventory.0.is_enabled":                "false",
						"incremental_inventory.0.optional_fields.0.field.#": "3",
					}),
				),
			},
		},
	})
}

// Test Oss BucketInventory. <<< Resource test cases.
