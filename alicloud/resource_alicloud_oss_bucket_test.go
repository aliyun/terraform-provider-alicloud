package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_oss_bucket", &resource.Sweeper{
		Name: "alicloud_oss_bucket",
		F:    testSweepOSSBuckets,
	})
}

func testSweepOSSBuckets(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf-test-",
		"test-bucket-",
		"tf-oss-test-",
		"tf-object-test-",
		"test-acc-alicloud-",
	}

	var options []oss.Option
	options = append(options, oss.MaxKeys(1000))

	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.ListBuckets(options...)
	})
	if err != nil {
		return fmt.Errorf("Error retrieving OSS buckets: %s", err)
	}
	resp, _ := raw.(oss.ListBucketsResult)
	sweeped := false

	for _, v := range resp.Buckets {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping OSS bucket: %s", name)
			continue
		}
		sweeped = true
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.Bucket(name)
		})
		if err != nil {
			return fmt.Errorf("Error getting bucket (%s): %#v", name, err)
		}
		bucket, _ := raw.(*oss.Bucket)
		if objects, err := bucket.ListObjects(options...); err != nil {
			log.Printf("[ERROR] Failed to list objects: %s", err)
		} else if len(objects.Objects) > 0 {
			for _, o := range objects.Objects {
				if err := bucket.DeleteObject(o.Key); err != nil {
					log.Printf("[ERROR] Failed to delete object (%s): %s.", o.Key, err)
				}
			}

		}

		log.Printf("[INFO] Deleting OSS bucket: %s", name)

		_, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return nil, ossClient.DeleteBucket(name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete OSS bucket (%s): %s", name, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudOssBucketBasic(t *testing.T) {
	var v oss.GetBucketInfoResult

	resourceId := "alicloud_oss_bucket.default"
	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketConfigDependence)
	hashcode1 := strconv.Itoa(expirationHash(map[string]interface{}{
		"days":                         365,
		"date":                         "",
		"created_before_date":          "",
		"expired_object_delete_marker": false,
	}))
	hashcode2 := strconv.Itoa(expirationHash(map[string]interface{}{
		"days":                         0,
		"date":                         "2018-01-12",
		"created_before_date":          "",
		"expired_object_delete_marker": false,
	}))
	hashcode3 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                3,
		"created_before_date": "",
		"storage_class":       "IA",
	}))
	hashcode4 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                30,
		"created_before_date": "",
		"storage_class":       "Archive",
	}))
	hashcode5 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                0,
		"created_before_date": "2023-11-11",
		"storage_class":       "IA",
	}))
	hashcode6 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                0,
		"created_before_date": "2023-11-10",
		"storage_class":       "Archive",
	}))
	hashcode7 := strconv.Itoa(expirationHash(map[string]interface{}{
		"days":                         0,
		"date":                         "",
		"created_before_date":          "2018-01-12",
		"expired_object_delete_marker": false,
	}))
	hashcode8 := strconv.Itoa(abortMultipartUploadHash(map[string]interface{}{
		"days":                0,
		"created_before_date": "2018-01-22",
	}))
	hashcode9 := strconv.Itoa(abortMultipartUploadHash(map[string]interface{}{
		"days":                10,
		"created_before_date": "",
	}))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "public-read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "public-read",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "public-read-write",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "public-read-write",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cors_rule": []map[string]interface{}{
						{
							"allowed_origins": []string{"*"},
							"allowed_methods": []string{"PUT", "GET"},
							"allowed_headers": []string{"authorization"},
						},
						{
							"allowed_origins": []string{"http://www.a.com", "http://www.b.com"},
							"allowed_methods": []string{"GET"},
							"allowed_headers": []string{"authorization"},
							"expose_headers":  []string{"x-oss-test", "x-oss-test1"},
							"max_age_seconds": "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cors_rule.#":                   "2",
						"cors_rule.0.allowed_headers.0": "authorization",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"website": []map[string]interface{}{
						{
							"index_document": "index.html",
							"error_document": "error.html",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"website.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{
						{
							"target_bucket": "${alicloud_oss_bucket.target.id}",
							"target_prefix": "log/",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logging.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"referer_config": []map[string]interface{}{
						{
							"allow_empty": "false",
							"referers": []string{
								"http://www.aliyun.com",
								"https://www.aliyun.com",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"referer_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule": []map[string]interface{}{
						{
							"id":      "rule1",
							"prefix":  "path1/",
							"enabled": "true",
							"expiration": []map[string]string{
								{
									"days": "365",
								},
							},
						},
						{
							"id":      "rule2",
							"prefix":  "path2/",
							"enabled": "true",
							"expiration": []map[string]string{
								{
									"date": "2018-01-12",
								},
							},
						},
						{
							"id":      "rule3",
							"prefix":  "path3/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
								{
									"days":          "30",
									"storage_class": "Archive",
								},
							},
						},
						{
							"id":      "rule4",
							"prefix":  "path4/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"created_before_date": "2023-11-11",
									"storage_class":       "IA",
								},
								{
									"created_before_date": "2023-11-10",
									"storage_class":       "Archive",
								},
							},
						},
						{
							"id":      "rule5",
							"prefix":  "path5/",
							"enabled": "true",
							"expiration": []map[string]string{
								{
									"created_before_date": "2018-01-12",
								},
							},
							"abort_multipart_upload": []map[string]string{
								{
									"created_before_date": "2018-01-22",
								},
							},
						},
						{
							"id":      "rule6",
							"prefix":  "path6/",
							"enabled": "true",
							"abort_multipart_upload": []map[string]string{
								{
									"days": "10",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                   "6",
						"lifecycle_rule.0.id":                                "rule1",
						"lifecycle_rule.0.prefix":                            "path1/",
						"lifecycle_rule.0.enabled":                           "true",
						"lifecycle_rule.0.expiration." + hashcode1 + ".days": "365",
						"lifecycle_rule.1.id":                                "rule2",
						"lifecycle_rule.1.prefix":                            "path2/",
						"lifecycle_rule.1.enabled":                           "true",
						"lifecycle_rule.1.expiration." + hashcode2 + ".date": "2018-01-12",

						"lifecycle_rule.2.id":                                          "rule3",
						"lifecycle_rule.2.prefix":                                      "path3/",
						"lifecycle_rule.2.enabled":                                     "true",
						"lifecycle_rule.2.transitions." + hashcode3 + ".days":          "3",
						"lifecycle_rule.2.transitions." + hashcode3 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.2.transitions." + hashcode4 + ".days":          "30",
						"lifecycle_rule.2.transitions." + hashcode4 + ".storage_class": string(oss.StorageArchive),

						"lifecycle_rule.3.id":      "rule4",
						"lifecycle_rule.3.prefix":  "path4/",
						"lifecycle_rule.3.enabled": "true",
						"lifecycle_rule.3.transitions." + hashcode5 + ".created_before_date": "2023-11-11",
						"lifecycle_rule.3.transitions." + hashcode5 + ".storage_class":       string(oss.StorageIA),
						"lifecycle_rule.3.transitions." + hashcode6 + ".created_before_date": "2023-11-10",
						"lifecycle_rule.3.transitions." + hashcode6 + ".storage_class":       string(oss.StorageArchive),

						"lifecycle_rule.4.id":      "rule5",
						"lifecycle_rule.4.prefix":  "path5/",
						"lifecycle_rule.4.enabled": "true",
						"lifecycle_rule.4.expiration." + hashcode7 + ".created_before_date":             "2018-01-12",
						"lifecycle_rule.4.abort_multipart_upload." + hashcode8 + ".created_before_date": "2018-01-22",

						"lifecycle_rule.5.id":      "rule6",
						"lifecycle_rule.5.prefix":  "path6/",
						"lifecycle_rule.5.enabled": "true",
						"lifecycle_rule.5.abort_multipart_upload." + hashcode9 + ".days": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": `{\"Statement\":[{\"Action\":[\"oss:*\"],\"Effect\":\"Allow\",\"Resource\":[\"acs:oss:*:*:*\"]}],\"Version\":\"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"key1": "value1",
						"Key2": "Value2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":    "2",
						"tags.key1": "value1",
						"tags.Key2": "Value2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"key1-update": "value1-update",
						"Key2-update": "Value2-update",
						"key3-new":    "value3-new",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":           "3",
						"tags.key1-update": "value1-update",
						"tags.Key2-update": "Value2-update",
						"tags.key3-new":    "value3-new",
						"tags.key1":        REMOVEKEY,
						"tags.Key2":        REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl":            "public-read",
					"cors_rule":      REMOVEKEY,
					"tags":           REMOVEKEY,
					"website":        REMOVEKEY,
					"logging":        REMOVEKEY,
					"referer_config": REMOVEKEY,
					"lifecycle_rule": REMOVEKEY,
					"policy":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl":                           "public-read",
						"cors_rule.#":                   "0",
						"cors_rule.0.allowed_headers.0": REMOVEKEY,
						"website.#":                     "0",
						"logging.#":                     "0",
						"referer_config.#":              "0",
						"lifecycle_rule.#":              "0",
						"lifecycle_rule.0.id":           REMOVEKEY,
						"lifecycle_rule.0.prefix":       REMOVEKEY,
						"lifecycle_rule.0.enabled":      REMOVEKEY,
						"lifecycle_rule.0.expiration." + hashcode1 + ".days": REMOVEKEY,
						"lifecycle_rule.1.id":                                REMOVEKEY,
						"lifecycle_rule.1.prefix":                            REMOVEKEY,
						"lifecycle_rule.1.enabled":                           REMOVEKEY,
						"lifecycle_rule.1.expiration." + hashcode2 + ".date": REMOVEKEY,

						"lifecycle_rule.2.id":                                          REMOVEKEY,
						"lifecycle_rule.2.prefix":                                      REMOVEKEY,
						"lifecycle_rule.2.enabled":                                     REMOVEKEY,
						"lifecycle_rule.2.transitions." + hashcode3 + ".days":          REMOVEKEY,
						"lifecycle_rule.2.transitions." + hashcode3 + ".storage_class": REMOVEKEY,
						"lifecycle_rule.2.transitions." + hashcode4 + ".days":          REMOVEKEY,
						"lifecycle_rule.2.transitions." + hashcode4 + ".storage_class": REMOVEKEY,

						"lifecycle_rule.3.id":      REMOVEKEY,
						"lifecycle_rule.3.prefix":  REMOVEKEY,
						"lifecycle_rule.3.enabled": REMOVEKEY,
						"lifecycle_rule.3.transitions." + hashcode5 + ".created_before_date": REMOVEKEY,
						"lifecycle_rule.3.transitions." + hashcode5 + ".storage_class":       REMOVEKEY,
						"lifecycle_rule.3.transitions." + hashcode6 + ".created_before_date": REMOVEKEY,
						"lifecycle_rule.3.transitions." + hashcode6 + ".storage_class":       REMOVEKEY,

						"lifecycle_rule.4.id":      REMOVEKEY,
						"lifecycle_rule.4.prefix":  REMOVEKEY,
						"lifecycle_rule.4.enabled": REMOVEKEY,
						"lifecycle_rule.4.expiration." + hashcode7 + ".created_before_date":             REMOVEKEY,
						"lifecycle_rule.4.abort_multipart_upload." + hashcode8 + ".created_before_date": REMOVEKEY,

						"lifecycle_rule.5.id":      REMOVEKEY,
						"lifecycle_rule.5.prefix":  REMOVEKEY,
						"lifecycle_rule.5.enabled": REMOVEKEY,
						"lifecycle_rule.5.abort_multipart_upload." + hashcode9 + ".days": REMOVEKEY,

						"tags.%":           "0",
						"tags.key1-update": REMOVEKEY,
						"tags.Key2-update": REMOVEKEY,
						"tags.key3-new":    REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudOssBucketVersioning(t *testing.T) {
	var v oss.GetBucketInfoResult

	resourceId := "alicloud_oss_bucket.default"
	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketConfigDependence)
	hashcode1 := strconv.Itoa(expirationHash(map[string]interface{}{
		"days":                         0,
		"date":                         "",
		"created_before_date":          "",
		"expired_object_delete_marker": true,
	}))
	hashcode2 := strconv.Itoa(expirationHash(map[string]interface{}{
		"days": 10,
	}))
	hashcode3 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":          3,
		"storage_class": "IA",
	}))
	hashcode4 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":          5,
		"storage_class": "Archive",
	}))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"versioning": []map[string]interface{}{
						{
							"status": "Enabled",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"versioning.0.status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"versioning": []map[string]interface{}{
						{
							"status": "Suspended",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"versioning.0.status": "Suspended",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule": []map[string]interface{}{
						{
							"id":      "rule1",
							"prefix":  "path1/",
							"enabled": "true",
							"expiration": []map[string]string{
								{
									"expired_object_delete_marker": "true",
								},
							},
						},
						{
							"id":      "rule2",
							"prefix":  "path2/",
							"enabled": "true",
							"noncurrent_version_expiration": []map[string]string{
								{
									"days": "10",
								},
							},
							"noncurrent_version_transition": []map[string]string{
								{
									"days":          "3",
									"storage_class": "IA",
								},
								{
									"days":          "5",
									"storage_class": "Archive",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":         "2",
						"lifecycle_rule.0.id":      "rule1",
						"lifecycle_rule.0.prefix":  "path1/",
						"lifecycle_rule.0.enabled": "true",
						"lifecycle_rule.0.expiration." + hashcode1 + ".expired_object_delete_marker": "true",

						"lifecycle_rule.1.id":      "rule2",
						"lifecycle_rule.1.prefix":  "path2/",
						"lifecycle_rule.1.enabled": "true",
						"lifecycle_rule.1.noncurrent_version_expiration." + hashcode2 + ".days":          "10",
						"lifecycle_rule.1.noncurrent_version_transition." + hashcode3 + ".days":          "3",
						"lifecycle_rule.1.noncurrent_version_transition." + hashcode3 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.1.noncurrent_version_transition." + hashcode4 + ".days":          "5",
						"lifecycle_rule.1.noncurrent_version_transition." + hashcode4 + ".storage_class": string(oss.StorageArchive),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudOssBucketCheckSseRule(t *testing.T) {
	var v oss.GetBucketInfoResult

	resourceId := "alicloud_oss_bucket.default"
	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm": "AES256",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.0.sse_algorithm": "AES256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm": "KMS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.0.sse_algorithm":     "KMS",
						"server_side_encryption_rule.0.kms_master_key_id": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm":     "KMS",
							"kms_master_key_id": "kms-id",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.0.sse_algorithm":     "KMS",
						"server_side_encryption_rule.0.kms_master_key_id": "kms-id",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.#":                   "0",
						"server_side_encryption_rule.0.sse_algorithm":     REMOVEKEY,
						"server_side_encryption_rule.0.kms_master_key_id": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudOssBucketCheckTransferAcc(t *testing.T) {
	var v oss.GetBucketInfoResult

	resourceId := "alicloud_oss_bucket.default"
	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transfer_acceleration": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transfer_acceleration.0.enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transfer_acceleration": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transfer_acceleration.0.enabled": "false",
					}),
				),
			},
		},
	})
}

func resourceOssBucketConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "target"{
	bucket = "%s-t"
}
`, name)
}

func TestAccAlicloudOssBucketBasic1(t *testing.T) {
	var v oss.GetBucketInfoResult

	resourceId := "alicloud_oss_bucket.default"
	ra := resourceAttrInit(resourceId, ossBucketBasicMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketConfigBasic)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": name,
					"acl":    "public-read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
						"acl":    "public-read",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func resourceOssBucketConfigBasic(name string) string {
	return fmt.Sprintf("")
}

var ossBucketBasicMap = map[string]string{
	"creation_date":    CHECKSET,
	"lifecycle_rule.#": "0",
}
