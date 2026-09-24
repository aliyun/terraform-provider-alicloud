package alicloud

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	alicloudCredentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		"tftest",
		"cri",
		"tf-example",
		"tf_example",
		"terraform-example",
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
	for _, v := range resp.Buckets {
		name := v.Name
		if !strings.HasSuffix(v.Location, client.RegionId) {
			continue
		}
		skip := true
		if !sweepAll() {
			if strings.HasPrefix(strings.ToLower(name), "terraform-alicloud-provider") {
				continue
			}
			if strings.HasPrefix(strings.ToLower(name), "terraform-remote-backend") {
				continue
			}
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
		}
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.Bucket(name)
		})
		if err != nil {
			return fmt.Errorf("Error getting bucket (%s): %#v", name, err)
		}
		bucket, _ := raw.(*oss.Bucket)
		var initialOptions []oss.Option
		nextMarker := ""
		for {
			var options []oss.Option
			options = append(options, initialOptions...)
			if nextMarker != "" {
				options = append(options, oss.Marker(nextMarker))
			}
			raw, err := client.WithOssBucketByName(bucket.BucketName, func(bucket *oss.Bucket) (interface{}, error) {
				return bucket.ListObjects(options...)
			})
			if err != nil {
				continue
			}
			response, _ := raw.(oss.ListObjectsResult)

			if response.Objects == nil || len(response.Objects) < 1 {
				break
			}

			if len(response.Objects) < 1 {
				break
			}

			if len(response.Objects) > 0 {
				for _, o := range response.Objects {
					if err := bucket.DeleteObject(o.Key); err != nil {
						log.Printf("[ERROR] Failed to delete object (%s): %s.", o.Key, err)
					}
				}
			}

			nextMarker = response.NextMarker
			if nextMarker == "" {
				break
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
	return nil
}

func TestAccAliCloudOssBucketBasic(t *testing.T) {
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
		"days":                     3,
		"created_before_date":      "",
		"storage_class":            "IA",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
	}))
	hashcode4 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     30,
		"created_before_date":      "",
		"storage_class":            "Archive",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
	}))
	hashcode5 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     0,
		"created_before_date":      "2023-11-11",
		"storage_class":            "IA",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
	}))
	hashcode6 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     0,
		"created_before_date":      "2023-11-10",
		"storage_class":            "Archive",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
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
						"bucket":                  name,
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			// Since public ACL creation is disabled by default, `public-*` ACL and policy test cases are disabled.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"acl": "public-read",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"acl": "public-read",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"acl": "public-read-write",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"acl": "public-read-write",
			//		}),
			//	),
			//},
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
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"policy": `{\"Statement\":[{\"Action\":[\"oss:*\"],\"Effect\":\"Allow\",\"Resource\":[\"acs:oss:*:*:*\"]}],\"Version\":\"1\"}`,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(nil),
			//	),
			//},
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
					"acl":            "private",
					"cors_rule":      REMOVEKEY,
					"tags":           REMOVEKEY,
					"website":        REMOVEKEY,
					"logging":        REMOVEKEY,
					"referer_config": REMOVEKEY,
					"lifecycle_rule": REMOVEKEY,
					//"policy":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl":                           "private",
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

func TestAccAliCloudOssBucketVersioning(t *testing.T) {
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
		"days":                     3,
		"storage_class":            "IA",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
	}))
	hashcode4 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     5,
		"storage_class":            "Archive",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
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
						"bucket":                  name,
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
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

func TestAccAliCloudOssBucketCheckSseRule(t *testing.T) {
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
						"bucket":                  name,
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
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
							"sse_algorithm": "SM4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.0.sse_algorithm":     "SM4",
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
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm":       "KMS",
							"kms_master_key_id":   "kms-id",
							"kms_data_encryption": "SM4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.0.sse_algorithm":       "KMS",
						"server_side_encryption_rule.0.kms_master_key_id":   "kms-id",
						"server_side_encryption_rule.0.kms_data_encryption": "SM4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.#":                     "0",
						"server_side_encryption_rule.0.sse_algorithm":       REMOVEKEY,
						"server_side_encryption_rule.0.kms_master_key_id":   REMOVEKEY,
						"server_side_encryption_rule.0.kms_data_encryption": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOssBucketCreateWithSseRuleKms(t *testing.T) {
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
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm": "KMS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                                            name,
						"access_monitor.#":                                  "1",
						"access_monitor.0.status":                           "Disabled",
						"resource_group_id":                                 CHECKSET,
						"server_side_encryption_rule.#":                     "1",
						"server_side_encryption_rule.0.sse_algorithm":       "KMS",
						"server_side_encryption_rule.0.kms_master_key_id":   "",
						"server_side_encryption_rule.0.kms_data_encryption": "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
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
						"server_side_encryption_rule.#":                     "1",
						"server_side_encryption_rule.0.sse_algorithm":       "AES256",
						"server_side_encryption_rule.0.kms_master_key_id":   "",
						"server_side_encryption_rule.0.kms_data_encryption": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.#":                     "0",
						"server_side_encryption_rule.0.sse_algorithm":       REMOVEKEY,
						"server_side_encryption_rule.0.kms_master_key_id":   REMOVEKEY,
						"server_side_encryption_rule.0.kms_data_encryption": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOssBucketCreateWithSseRuleKmsId(t *testing.T) {
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
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm":     "KMS",
							"kms_master_key_id": "kms-id",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                                          name,
						"access_monitor.#":                                "1",
						"access_monitor.0.status":                         "Disabled",
						"resource_group_id":                               CHECKSET,
						"server_side_encryption_rule.#":                   "1",
						"server_side_encryption_rule.0.sse_algorithm":     "KMS",
						"server_side_encryption_rule.0.kms_master_key_id": "kms-id",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm":     "SM4",
							"kms_master_key_id": "",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.#":                   "1",
						"server_side_encryption_rule.0.sse_algorithm":     "SM4",
						"server_side_encryption_rule.0.kms_master_key_id": "",
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

func TestAccAliCloudOssBucketCreateWithSseRuleKmsSm4(t *testing.T) {
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
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm":       "KMS",
							"kms_master_key_id":   "",
							"kms_data_encryption": "SM4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                                            name,
						"access_monitor.#":                                  "1",
						"access_monitor.0.status":                           "Disabled",
						"resource_group_id":                                 CHECKSET,
						"server_side_encryption_rule.#":                     "1",
						"server_side_encryption_rule.0.sse_algorithm":       "KMS",
						"server_side_encryption_rule.0.kms_master_key_id":   "",
						"server_side_encryption_rule.0.kms_data_encryption": "SM4",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm":       "KMS",
							"kms_master_key_id":   "kms-id",
							"kms_data_encryption": "",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.#":                     "1",
						"server_side_encryption_rule.0.sse_algorithm":       "KMS",
						"server_side_encryption_rule.0.kms_master_key_id":   "kms-id",
						"server_side_encryption_rule.0.kms_data_encryption": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.#":                     "0",
						"server_side_encryption_rule.0.sse_algorithm":       REMOVEKEY,
						"server_side_encryption_rule.0.kms_master_key_id":   REMOVEKEY,
						"server_side_encryption_rule.0.kms_data_encryption": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOssBucketCreateWithSseRuleKmsIdSm4(t *testing.T) {
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
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm":       "KMS",
							"kms_master_key_id":   "kms-id",
							"kms_data_encryption": "SM4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                                            name,
						"access_monitor.#":                                  "1",
						"access_monitor.0.status":                           "Disabled",
						"resource_group_id":                                 CHECKSET,
						"server_side_encryption_rule.#":                     "1",
						"server_side_encryption_rule.0.sse_algorithm":       "KMS",
						"server_side_encryption_rule.0.kms_master_key_id":   "kms-id",
						"server_side_encryption_rule.0.kms_data_encryption": "SM4",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": []map[string]interface{}{
						{
							"sse_algorithm":       "SM4",
							"kms_master_key_id":   "",
							"kms_data_encryption": "",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.#":                     "1",
						"server_side_encryption_rule.0.sse_algorithm":       "SM4",
						"server_side_encryption_rule.0.kms_master_key_id":   "",
						"server_side_encryption_rule.0.kms_data_encryption": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption_rule": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_side_encryption_rule.#":                     "0",
						"server_side_encryption_rule.0.sse_algorithm":       REMOVEKEY,
						"server_side_encryption_rule.0.kms_master_key_id":   REMOVEKEY,
						"server_side_encryption_rule.0.kms_data_encryption": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOssBucketCheckTransferAcc(t *testing.T) {
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
						"bucket":                  name,
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
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

func TestAccAliCloudOssBucketBasic1(t *testing.T) {
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
					"acl":    "private",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                  name,
						"acl":                     "private",
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
		},
	})
}

func TestAccAliCloudOssBucketBasic_no_set_name(t *testing.T) {
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                  CHECKSET,
						"acl":                     "private",
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
		},
	})
}
func TestAccAliCloudOssBucketColdArchive(t *testing.T) {
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
	hashcode3 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     3,
		"created_before_date":      "",
		"storage_class":            "IA",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
	}))
	hashcode4 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     30,
		"created_before_date":      "",
		"storage_class":            "ColdArchive",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
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
					"bucket":        name,
					"acl":           "public-read",
					"storage_class": "ColdArchive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                  name,
						"acl":                     "public-read",
						"storage_class":           "ColdArchive",
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule": []map[string]interface{}{
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
									"storage_class": "ColdArchive",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "1",
						"lifecycle_rule.0.id":                                          "rule3",
						"lifecycle_rule.0.prefix":                                      "path3/",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode3 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode3 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.transitions." + hashcode4 + ".days":          "30",
						"lifecycle_rule.0.transitions." + hashcode4 + ".storage_class": string(oss.StorageColdArchive),
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOssBucketLifeCycleRuleOverlap(t *testing.T) {
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
	hashcode1 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     3,
		"created_before_date":      "",
		"storage_class":            "IA",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
	}))
	hashcode2 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     30,
		"created_before_date":      "",
		"storage_class":            "IA",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
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
					"acl":    "public-read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                  name,
						"acl":                     "public-read",
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule_allow_same_action_overlap": true,
					"lifecycle_rule": []map[string]interface{}{
						{
							"id":      "rule1",
							"prefix":  "path3/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
						},
						{
							"id":      "rule2",
							"prefix":  "path3/abc",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":          "30",
									"storage_class": "IA",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "2",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.prefix":                                      "path3/",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.1.id":                                          "rule2",
						"lifecycle_rule.1.prefix":                                      "path3/abc",
						"lifecycle_rule.1.enabled":                                     "true",
						"lifecycle_rule.1.transitions." + hashcode2 + ".days":          "30",
						"lifecycle_rule.1.transitions." + hashcode2 + ".storage_class": string(oss.StorageIA),
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOssBucketAccessMonitor(t *testing.T) {
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
	hashcode1 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     3,
		"created_before_date":      "",
		"storage_class":            "IA",
		"is_access_time":           true,
		"return_to_std_when_visit": false,
	}))
	hashcode2 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     30,
		"created_before_date":      "",
		"storage_class":            "IA",
		"is_access_time":           true,
		"return_to_std_when_visit": true,
	}))
	hashcode3 := strconv.Itoa(expirationHash(map[string]interface{}{
		"days":                         0,
		"date":                         "",
		"created_before_date":          "",
		"expired_object_delete_marker": true,
	}))
	hashcode4 := strconv.Itoa(expirationHash(map[string]interface{}{
		"days": 10,
	}))
	hashcode5 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     3,
		"storage_class":            "IA",
		"is_access_time":           true,
		"return_to_std_when_visit": true,
	}))
	hashcode6 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     5,
		"storage_class":            "Archive",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
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
					"acl":    "public-read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                  name,
						"acl":                     "public-read",
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			// enable accesss monitor
			{
				Config: testAccConfig(map[string]interface{}{
					"access_monitor": []map[string]interface{}{
						{
							"status": "Enabled",
						},
					},
					"lifecycle_rule": []map[string]interface{}{
						{
							"id":      "rule1",
							"prefix":  "path1/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":                     "3",
									"storage_class":            "IA",
									"is_access_time":           "true",
									"return_to_std_when_visit": "false",
								},
							},
						},
						{
							"id":      "rule2",
							"prefix":  "path2/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":                     "30",
									"storage_class":            "IA",
									"is_access_time":           "true",
									"return_to_std_when_visit": "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                                        "2",
						"lifecycle_rule.0.id":                                                     "rule1",
						"lifecycle_rule.0.prefix":                                                 "path1/",
						"lifecycle_rule.0.enabled":                                                "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":                     "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class":            string(oss.StorageIA),
						"lifecycle_rule.0.transitions." + hashcode1 + ".is_access_time":           "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".return_to_std_when_visit": "false",
						"lifecycle_rule.1.id":                                                     "rule2",
						"lifecycle_rule.1.prefix":                                                 "path2/",
						"lifecycle_rule.1.enabled":                                                "true",
						"lifecycle_rule.1.transitions." + hashcode2 + ".days":                     "30",
						"lifecycle_rule.1.transitions." + hashcode2 + ".storage_class":            string(oss.StorageIA),
						"lifecycle_rule.1.transitions." + hashcode2 + ".is_access_time":           "true",
						"lifecycle_rule.1.transitions." + hashcode2 + ".return_to_std_when_visit": "true",
						"access_monitor.#":                                                        "1",
						"access_monitor.0.status":                                                 "Enabled",
					}),
				),
			},
			// disable accesss monitor
			{
				Config: testAccConfig(map[string]interface{}{
					"access_monitor": []map[string]interface{}{
						{
							"status": "Disabled",
						},
					},
					"lifecycle_rule": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_monitor.#":                                                        "1",
						"access_monitor.0.status":                                                 "Disabled",
						"lifecycle_rule.#":                                                        "0",
						"lifecycle_rule.0.id":                                                     REMOVEKEY,
						"lifecycle_rule.0.prefix":                                                 REMOVEKEY,
						"lifecycle_rule.0.enabled":                                                REMOVEKEY,
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":                     REMOVEKEY,
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class":            REMOVEKEY,
						"lifecycle_rule.0.transitions." + hashcode1 + ".is_access_time":           REMOVEKEY,
						"lifecycle_rule.0.transitions." + hashcode1 + ".return_to_std_when_visit": REMOVEKEY,
						"lifecycle_rule.1.id":                                                     REMOVEKEY,
						"lifecycle_rule.1.prefix":                                                 REMOVEKEY,
						"lifecycle_rule.1.enabled":                                                REMOVEKEY,
						"lifecycle_rule.1.transitions." + hashcode2 + ".days":                     REMOVEKEY,
						"lifecycle_rule.1.transitions." + hashcode2 + ".storage_class":            REMOVEKEY,
						"lifecycle_rule.1.transitions." + hashcode2 + ".is_access_time":           REMOVEKEY,
						"lifecycle_rule.1.transitions." + hashcode2 + ".return_to_std_when_visit": REMOVEKEY,
					}),
				),
			},
			// enable versioning and accesss monitor
			{
				Config: testAccConfig(map[string]interface{}{
					"access_monitor": []map[string]interface{}{
						{
							"status": "Enabled",
						},
					},
					"versioning": []map[string]interface{}{
						{
							"status": "Enabled",
						},
					},
					"lifecycle_rule": []map[string]interface{}{
						{
							"id":      "rule1",
							"prefix":  "path1/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":                     "3",
									"storage_class":            "IA",
									"is_access_time":           "true",
									"return_to_std_when_visit": "false",
								},
							},
						},
						{
							"id":      "rule2",
							"prefix":  "path2/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":                     "30",
									"storage_class":            "IA",
									"is_access_time":           "true",
									"return_to_std_when_visit": "true",
								},
							},
						},
						{
							"id":      "rule3",
							"prefix":  "path3/",
							"enabled": "true",
							"expiration": []map[string]string{
								{
									"expired_object_delete_marker": "true",
								},
							},
						},
						{
							"id":      "rule4",
							"prefix":  "path4/",
							"enabled": "true",
							"noncurrent_version_expiration": []map[string]string{
								{
									"days": "10",
								},
							},
							"noncurrent_version_transition": []map[string]string{
								{
									"days":                     "3",
									"storage_class":            "IA",
									"is_access_time":           "true",
									"return_to_std_when_visit": "true",
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
						"access_monitor.#":         "1",
						"access_monitor.0.status":  "Enabled",
						"versioning.#":             "1",
						"versioning.0.status":      "Enabled",
						"lifecycle_rule.#":         "4",
						"lifecycle_rule.0.id":      "rule1",
						"lifecycle_rule.0.prefix":  "path1/",
						"lifecycle_rule.0.enabled": "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":                     "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class":            string(oss.StorageIA),
						"lifecycle_rule.0.transitions." + hashcode1 + ".is_access_time":           "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".return_to_std_when_visit": "false",
						"lifecycle_rule.1.id":                                                     "rule2",
						"lifecycle_rule.1.prefix":                                                 "path2/",
						"lifecycle_rule.1.enabled":                                                "true",
						"lifecycle_rule.1.transitions." + hashcode2 + ".days":                     "30",
						"lifecycle_rule.1.transitions." + hashcode2 + ".storage_class":            string(oss.StorageIA),
						"lifecycle_rule.1.transitions." + hashcode2 + ".is_access_time":           "true",
						"lifecycle_rule.1.transitions." + hashcode2 + ".return_to_std_when_visit": "true",

						"lifecycle_rule.2.id":      "rule3",
						"lifecycle_rule.2.prefix":  "path3/",
						"lifecycle_rule.2.enabled": "true",
						"lifecycle_rule.2.expiration." + hashcode3 + ".expired_object_delete_marker": "true",

						"lifecycle_rule.3.id":      "rule4",
						"lifecycle_rule.3.prefix":  "path4/",
						"lifecycle_rule.3.enabled": "true",
						"lifecycle_rule.3.noncurrent_version_expiration." + hashcode4 + ".days":                     "10",
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode5 + ".days":                     "3",
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode5 + ".storage_class":            string(oss.StorageIA),
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode5 + ".is_access_time":           "true",
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode5 + ".return_to_std_when_visit": "true",
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode6 + ".days":                     "5",
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode6 + ".storage_class":            string(oss.StorageArchive),
					}),
				),
			},
			// disable accesss monitor and status
			{
				Config: testAccConfig(map[string]interface{}{
					"access_monitor": []map[string]interface{}{
						{
							"status": "Disabled",
						},
					},
					"versioning": []map[string]interface{}{
						{
							"status": "Suspended",
						},
					},
					"lifecycle_rule": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_monitor.#":         "1",
						"access_monitor.0.status":  "Disabled",
						"versioning.#":             "1",
						"versioning.0.status":      "Suspended",
						"lifecycle_rule.#":         "0",
						"lifecycle_rule.0.id":      REMOVEKEY,
						"lifecycle_rule.0.prefix":  REMOVEKEY,
						"lifecycle_rule.0.enabled": REMOVEKEY,
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":                     REMOVEKEY,
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class":            REMOVEKEY,
						"lifecycle_rule.0.transitions." + hashcode1 + ".is_access_time":           REMOVEKEY,
						"lifecycle_rule.0.transitions." + hashcode1 + ".return_to_std_when_visit": REMOVEKEY,
						"lifecycle_rule.1.id":                                                     REMOVEKEY,
						"lifecycle_rule.1.prefix":                                                 REMOVEKEY,
						"lifecycle_rule.1.enabled":                                                REMOVEKEY,
						"lifecycle_rule.1.transitions." + hashcode2 + ".days":                     REMOVEKEY,
						"lifecycle_rule.1.transitions." + hashcode2 + ".storage_class":            REMOVEKEY,
						"lifecycle_rule.1.transitions." + hashcode2 + ".is_access_time":           REMOVEKEY,
						"lifecycle_rule.1.transitions." + hashcode2 + ".return_to_std_when_visit": REMOVEKEY,

						"lifecycle_rule.2.id":      REMOVEKEY,
						"lifecycle_rule.2.prefix":  REMOVEKEY,
						"lifecycle_rule.2.enabled": REMOVEKEY,
						"lifecycle_rule.2.expiration." + hashcode3 + ".expired_object_delete_marker": REMOVEKEY,

						"lifecycle_rule.3.id":      REMOVEKEY,
						"lifecycle_rule.3.prefix":  REMOVEKEY,
						"lifecycle_rule.3.enabled": REMOVEKEY,
						"lifecycle_rule.3.noncurrent_version_expiration." + hashcode4 + ".days":                     REMOVEKEY,
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode5 + ".days":                     REMOVEKEY,
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode5 + ".storage_class":            REMOVEKEY,
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode5 + ".is_access_time":           REMOVEKEY,
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode5 + ".return_to_std_when_visit": REMOVEKEY,
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode6 + ".days":                     REMOVEKEY,
						"lifecycle_rule.3.noncurrent_version_transition." + hashcode6 + ".storage_class":            REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOssBucketDeepColdArchive(t *testing.T) {
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
	hashcode3 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     3,
		"created_before_date":      "",
		"storage_class":            "IA",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
	}))
	hashcode4 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     30,
		"created_before_date":      "",
		"storage_class":            "DeepColdArchive",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
	}))
	hashcode6 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     5,
		"storage_class":            "DeepColdArchive",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
	}))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.OssDeepColdArchiveSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":        name,
					"acl":           "public-read",
					"storage_class": "DeepColdArchive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                  name,
						"acl":                     "public-read",
						"storage_class":           "DeepColdArchive",
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule": []map[string]interface{}{
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
									"storage_class": "DeepColdArchive",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "1",
						"lifecycle_rule.0.id":                                          "rule3",
						"lifecycle_rule.0.prefix":                                      "path3/",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode3 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode3 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.transitions." + hashcode4 + ".days":          "30",
						"lifecycle_rule.0.transitions." + hashcode4 + ".storage_class": string(oss.StorageDeepColdArchive),
					}),
				),
			},
			// enable versioning
			{
				Config: testAccConfig(map[string]interface{}{
					"versioning": []map[string]interface{}{
						{
							"status": "Enabled",
						},
					},
					"lifecycle_rule": []map[string]interface{}{
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
									"storage_class": "DeepColdArchive",
								},
							},
						},
						{
							"id":      "rule4",
							"prefix":  "path4/",
							"enabled": "true",
							"noncurrent_version_transition": []map[string]string{
								{
									"days":          "5",
									"storage_class": "DeepColdArchive",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"versioning.#":             "1",
						"versioning.0.status":      "Enabled",
						"lifecycle_rule.#":         "2",
						"lifecycle_rule.0.id":      "rule3",
						"lifecycle_rule.0.prefix":  "path3/",
						"lifecycle_rule.0.enabled": "true",
						"lifecycle_rule.0.transitions." + hashcode3 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode3 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.transitions." + hashcode4 + ".days":          "30",
						"lifecycle_rule.0.transitions." + hashcode4 + ".storage_class": string(oss.StorageDeepColdArchive),

						"lifecycle_rule.1.id":      "rule4",
						"lifecycle_rule.1.prefix":  "path4/",
						"lifecycle_rule.1.enabled": "true",
						"lifecycle_rule.1.noncurrent_version_transition." + hashcode6 + ".days":          "5",
						"lifecycle_rule.1.noncurrent_version_transition." + hashcode6 + ".storage_class": string(oss.StorageDeepColdArchive),
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOssBucketLifeCycleTags(t *testing.T) {
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
	hashcode1 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     3,
		"created_before_date":      "",
		"storage_class":            "IA",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
	}))
	hashcode2 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     30,
		"created_before_date":      "",
		"storage_class":            "IA",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
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
					"acl":    "public-read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                  name,
						"acl":                     "public-read",
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule": []map[string]interface{}{
						{
							"id":      "rule1",
							"prefix":  "path1/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
							"tags": map[string]string{
								"key1": "value1",
								"key2": "value2",
							},
						},
						{
							"id":      "rule2",
							"prefix":  "path2/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":          "30",
									"storage_class": "IA",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "2",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.prefix":                                      "path1/",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.tags.%":                                      "2",
						"lifecycle_rule.0.tags.key1":                                   "value1",
						"lifecycle_rule.0.tags.key2":                                   "value2",

						"lifecycle_rule.1.id":                                          "rule2",
						"lifecycle_rule.1.prefix":                                      "path2/",
						"lifecycle_rule.1.enabled":                                     "true",
						"lifecycle_rule.1.transitions." + hashcode2 + ".days":          "30",
						"lifecycle_rule.1.transitions." + hashcode2 + ".storage_class": string(oss.StorageIA),
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
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
							"tags": map[string]string{
								"key1": "value1-1",
							},
						},
						{
							"id":      "rule2",
							"prefix":  "path2/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":          "30",
									"storage_class": "IA",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "2",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.prefix":                                      "path1/",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.tags.%":                                      "1",
						"lifecycle_rule.0.tags.key1":                                   "value1-1",
						"lifecycle_rule.0.tags.key2":                                   REMOVEKEY,

						"lifecycle_rule.1.id":                                          "rule2",
						"lifecycle_rule.1.prefix":                                      "path2/",
						"lifecycle_rule.1.enabled":                                     "true",
						"lifecycle_rule.1.transitions." + hashcode2 + ".days":          "30",
						"lifecycle_rule.1.transitions." + hashcode2 + ".storage_class": string(oss.StorageIA),
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
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
						},
						{
							"id":      "rule2",
							"prefix":  "path2/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":          "30",
									"storage_class": "IA",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "2",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.prefix":                                      "path1/",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.tags.%":                                      "0",
						"lifecycle_rule.0.tags.key1":                                   REMOVEKEY,

						"lifecycle_rule.1.id":                                          "rule2",
						"lifecycle_rule.1.prefix":                                      "path2/",
						"lifecycle_rule.1.enabled":                                     "true",
						"lifecycle_rule.1.transitions." + hashcode2 + ".days":          "30",
						"lifecycle_rule.1.transitions." + hashcode2 + ".storage_class": string(oss.StorageIA),
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOssBucketLifeCycleFilter(t *testing.T) {
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
	hashcode1 := strconv.Itoa(transitionsHash(map[string]interface{}{
		"days":                     3,
		"created_before_date":      "",
		"storage_class":            "IA",
		"is_access_time":           false,
		"return_to_std_when_visit": false,
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
					"acl":    "public-read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                  name,
						"acl":                     "public-read",
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule": []map[string]interface{}{
						{
							"id":      "rule1",
							"prefix":  "path1/",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
							"filter": []map[string]interface{}{
								{
									"not": []map[string]interface{}{
										{
											"prefix": "path1/sub",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "1",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.prefix":                                      "path1/",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.filter.#":                                    "1",
						"lifecycle_rule.0.filter.0.not.#":                              "1",
						"lifecycle_rule.0.filter.0.not.0.prefix":                       "path1/sub",
						"lifecycle_rule.0.filter.0.not.0.tag.#":                        "0",
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
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
							"filter": []map[string]interface{}{
								{
									"not": []map[string]interface{}{
										{
											"prefix": "path1/sub1",
											"tag": []map[string]interface{}{
												{
													"key":   "key1",
													"value": "value1",
												},
											},
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "1",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.prefix":                                      "path1/",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.filter.#":                                    "1",
						"lifecycle_rule.0.filter.0.not.#":                              "1",
						"lifecycle_rule.0.filter.0.not.0.prefix":                       "path1/sub1",
						"lifecycle_rule.0.filter.0.not.0.tag.#":                        "1",
						"lifecycle_rule.0.filter.0.not.0.tag.0.key":                    "key1",
						"lifecycle_rule.0.filter.0.not.0.tag.0.value":                  "value1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule": []map[string]interface{}{
						{
							"id":      "rule1",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
							"filter": []map[string]interface{}{
								{
									"object_size_greater_than": "1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "1",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.prefix":                                      REMOVEKEY,
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.filter.#":                                    "1",
						"lifecycle_rule.0.filter.0.not.#":                              "0",
						"lifecycle_rule.0.filter.0.not.0.prefix":                       REMOVEKEY,
						"lifecycle_rule.0.filter.0.not.0.tag.#":                        "0",
						"lifecycle_rule.0.filter.0.not.0.tag.0.key":                    REMOVEKEY,
						"lifecycle_rule.0.filter.0.not.0.tag.0.value":                  REMOVEKEY,
						"lifecycle_rule.0.filter.0.object_size_greater_than":           "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule": []map[string]interface{}{
						{
							"id":      "rule1",
							"enabled": "true",
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
							"filter": []map[string]interface{}{
								{
									"object_size_less_than": "1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "1",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.filter.#":                                    "1",
						"lifecycle_rule.0.filter.0.not.#":                              "0",
						"lifecycle_rule.0.filter.0.object_size_greater_than":           REMOVEKEY,
						"lifecycle_rule.0.filter.0.object_size_less_than":              "1",
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
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
							"filter": []map[string]interface{}{
								{
									"not": []map[string]interface{}{
										{
											"prefix": "path1/sub1",
											"tag": []map[string]interface{}{
												{
													"key":   "key2",
													"value": "value2",
												},
											},
										},
									},
									"object_size_greater_than": "2",
									"object_size_less_than":    "4",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "1",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.prefix":                                      "path1/",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.filter.#":                                    "1",
						"lifecycle_rule.0.filter.0.not.#":                              "1",
						"lifecycle_rule.0.filter.0.not.0.prefix":                       "path1/sub1",
						"lifecycle_rule.0.filter.0.not.0.tag.#":                        "1",
						"lifecycle_rule.0.filter.0.not.0.tag.0.key":                    "key2",
						"lifecycle_rule.0.filter.0.not.0.tag.0.value":                  "value2",
						"lifecycle_rule.0.filter.0.object_size_greater_than":           "2",
						"lifecycle_rule.0.filter.0.object_size_less_than":              "4",
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
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "1",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.prefix":                                      "path1/",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.filter.#":                                    "0",
						"lifecycle_rule.0.filter.0.not.#":                              "0",
						"lifecycle_rule.0.filter.0.not.0.prefix":                       REMOVEKEY,
						"lifecycle_rule.0.filter.0.not.0.tag.#":                        "0",
						"lifecycle_rule.0.filter.0.not.0.tag.0.key":                    REMOVEKEY,
						"lifecycle_rule.0.filter.0.not.0.tag.0.value":                  REMOVEKEY,
						"lifecycle_rule.0.filter.0.object_size_greater_than":           REMOVEKEY,
						"lifecycle_rule.0.filter.0.object_size_less_than":              REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule": []map[string]interface{}{
						{
							"id":      "rule1",
							"enabled": "true",
							"prefix":  "",
							"transitions": []map[string]interface{}{
								{
									"days":          "3",
									"storage_class": "IA",
								},
							},
							"filter": []map[string]interface{}{
								{
									"not": []map[string]interface{}{
										{
											"prefix": "",
											"tag": []map[string]interface{}{
												{
													"key":   "key2",
													"value": "value2",
												},
											},
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule.#":                                             "1",
						"lifecycle_rule.0.id":                                          "rule1",
						"lifecycle_rule.0.prefix":                                      "",
						"lifecycle_rule.0.enabled":                                     "true",
						"lifecycle_rule.0.transitions." + hashcode1 + ".days":          "3",
						"lifecycle_rule.0.transitions." + hashcode1 + ".storage_class": string(oss.StorageIA),
						"lifecycle_rule.0.filter.#":                                    "1",
						"lifecycle_rule.0.filter.0.not.#":                              "1",
						"lifecycle_rule.0.filter.0.not.0.prefix":                       "",
						"lifecycle_rule.0.filter.0.not.0.tag.#":                        "1",
						"lifecycle_rule.0.filter.0.not.0.tag.0.key":                    "key2",
						"lifecycle_rule.0.filter.0.not.0.tag.0.value":                  "value2",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOssBucketResourceGroup(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketResourceGroupDependence)
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
					"bucket":            name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"redundancy_type":   "LRS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                  name,
						"access_monitor.#":        "1",
						"access_monitor.0.status": "Disabled",
						"resource_group_id":       CHECKSET,
						"redundancy_type":         "LRS",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "lifecycle_rule_allow_same_action_overlap"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"force_destroy":     true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceOssBucketConfigBasic(name string) string {
	return fmt.Sprintf("")
}

func resourceOssBucketConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "target"{
  bucket = "%s-t"
}
`, name)
}

func resourceOssBucketResourceGroupDependence(name string) string {
	return fmt.Sprintf(`
data "alicloud_resource_manager_resource_groups" "default" {}
locals {
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  resource_group_id1 = data.alicloud_resource_manager_resource_groups.default.groups.1.id
}`)
}

var ossBucketBasicMap = map[string]string{
	"creation_date":    CHECKSET,
	"lifecycle_rule.#": "0",
}

// The tests in this file lock the create/delete pre-check contract of
// alicloud_oss_bucket against an in-process HTTP mock of the OSS XML API.
// Every scenario also asserts that no service-level ListBuckets request is
// issued on any create or delete path.

type ossMockExchange struct {
	Method   string
	Path     string
	RawQuery string
}

type ossMockServer struct {
	*httptest.Server
	mu        sync.Mutex
	exchanges []ossMockExchange
}

func newOssMockServer(handler http.HandlerFunc) *ossMockServer {
	m := &ossMockServer{}
	m.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.mu.Lock()
		m.exchanges = append(m.exchanges, ossMockExchange{Method: r.Method, Path: r.URL.Path, RawQuery: r.URL.RawQuery})
		m.mu.Unlock()
		handler(w, r)
	}))
	return m
}

func (m *ossMockServer) count(method, path, queryKey string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := 0
	for _, e := range m.exchanges {
		if e.Method != method || strings.TrimRight(e.Path, "/") != strings.TrimRight(path, "/") {
			continue
		}
		if queryKey == "" || ossQueryKeyPresent(e.RawQuery, queryKey) {
			n++
		}
	}
	return n
}

func (m *ossMockServer) assertNoListBuckets(t *testing.T) {
	t.Helper()
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.exchanges {
		if e.Method == http.MethodGet && e.Path == "/" {
			t.Fatalf("unexpected service-level ListBuckets request: GET /?%s", e.RawQuery)
		}
	}
}

func ossQueryKeyPresent(rawQuery, key string) bool {
	for _, part := range strings.Split(rawQuery, "&") {
		if part == key || strings.HasPrefix(part, key+"=") {
			return true
		}
	}
	return false
}

func hasOssSubresource(q url.Values, key string) bool {
	_, ok := q[key]
	return ok
}

func writeOssXML(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	fmt.Fprint(w, body)
}

func writeOssError(w http.ResponseWriter, status int, code string) {
	writeOssXML(w, status, fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><Error><Code>%s</Code><Message>mock %s</Message><RequestId>mock-request-id</RequestId><HostId>mock-host</HostId></Error>`, code, code))
}

func ossBucketInfoXML(name string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><BucketInfo><Bucket><Name>%s</Name><Location>oss-cn-hangzhou</Location><CreationDate>2026-09-21T00:00:00.000Z</CreationDate><ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint><IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint><AccessControlList><Grant>private</Grant></AccessControlList><DataRedundancyType>LRS</DataRedundancyType><Owner><ID>1234567890</ID><DisplayName>1234567890</DisplayName></Owner><StorageClass>Standard</StorageClass></Bucket></BucketInfo>`, name)
}

// serveEmptyOssBucketSubresources answers every subresource read that
// resourceAlicloudOssBucketRead issues for a freshly created bucket.
func serveEmptyOssBucketSubresources(w http.ResponseWriter, r *http.Request) bool {
	q := r.URL.Query()
	switch {
	case hasOssSubresource(q, "cors"):
		writeOssError(w, http.StatusNotFound, "NoSuchCORSConfiguration")
	case hasOssSubresource(q, "website"):
		writeOssError(w, http.StatusNotFound, "NoSuchWebsiteConfiguration")
	case hasOssSubresource(q, "logging"):
		writeOssXML(w, http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?><BucketLoggingStatus></BucketLoggingStatus>`)
	case hasOssSubresource(q, "referer"):
		writeOssXML(w, http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?><RefererConfiguration><AllowEmptyReferer>true</AllowEmptyReferer><RefererList></RefererList></RefererConfiguration>`)
	case hasOssSubresource(q, "lifecycle"):
		writeOssError(w, http.StatusNotFound, "NoSuchLifecycle")
	case hasOssSubresource(q, "policy"):
		writeOssError(w, http.StatusNotFound, "NoSuchPolicy")
	case hasOssSubresource(q, "tagging"):
		writeOssXML(w, http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?><Tagging><TagSet></TagSet></Tagging>`)
	case hasOssSubresource(q, "transferAcceleration"):
		writeOssError(w, http.StatusNotFound, "NoSuchTransferAccelerationConfiguration")
	case hasOssSubresource(q, "resourceGroup"):
		writeOssXML(w, http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?><BucketResourceGroupConfiguration><ResourceGroupId></ResourceGroupId></BucketResourceGroupConfiguration>`)
	case r.Method == http.MethodPut && hasOssSubresource(q, "acl"):
		writeOssXML(w, http.StatusOK, "")
	default:
		return false
	}
	return true
}

func newOssMockClient(t *testing.T, endpoint string) *connectivity.AliyunClient {
	t.Helper()
	t.Setenv("HTTP_PROXY", "")
	t.Setenv("http_proxy", "")
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")
	staticCredential, err := alicloudCredentials.NewCredential(new(alicloudCredentials.Config).
		SetType("access_key").
		SetAccessKeyId("mock-access-key").
		SetAccessKeySecret("mock-secret-key"))
	if err != nil {
		t.Fatalf("building mock credential failed: %v", err)
	}
	endpoints := &sync.Map{}
	endpoints.Store("oss", endpoint)
	config := connectivity.Config{
		AccessKey:            "mock-access-key",
		SecretKey:            "mock-secret-key",
		Credential:           staticCredential,
		Region:               "cn-hangzhou",
		RegionId:             "cn-hangzhou",
		AccountType:          "Domestic",
		Protocol:             "HTTP",
		SkipRegionValidation: true,
		Endpoints:            endpoints,
		SignVersion:          &sync.Map{},
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("building mock OSS client failed: %v", err)
	}
	return client
}

func ossBucketResourceData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceAlicloudOssBucket().Schema, raw)
}

func TestUnitAliCloudOssBucketCreateFailsWhenBucketAlreadyExists(t *testing.T) {
	bucketName := "mock-existing-bucket"
	mock := newOssMockServer(func(w http.ResponseWriter, r *http.Request) {
		if hasOssSubresource(r.URL.Query(), "bucketInfo") {
			writeOssXML(w, http.StatusOK, ossBucketInfoXML(bucketName))
			return
		}
		writeOssError(w, http.StatusInternalServerError, "UnexpectedMockRequest")
	})
	defer mock.Close()
	client := newOssMockClient(t, mock.URL)

	d := ossBucketResourceData(t, map[string]interface{}{"bucket": bucketName})
	err := resourceAlicloudOssBucketCreate(d, client)
	if err == nil {
		t.Fatal("expected an error when the bucket already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected an import/rename hint, got: %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("expected empty resource ID, got %q", d.Id())
	}
	if n := mock.count(http.MethodPut, "/"+bucketName, ""); n != 0 {
		t.Fatalf("expected no CreateBucket request, got %d", n)
	}
	mock.assertNoListBuckets(t)
}

func TestUnitAliCloudOssBucketCreateFailsWhenPrecheckDenied(t *testing.T) {
	bucketName := "mock-denied-bucket"
	mock := newOssMockServer(func(w http.ResponseWriter, r *http.Request) {
		if hasOssSubresource(r.URL.Query(), "bucketInfo") {
			writeOssError(w, http.StatusForbidden, "AccessDenied")
			return
		}
		writeOssError(w, http.StatusInternalServerError, "UnexpectedMockRequest")
	})
	defer mock.Close()
	client := newOssMockClient(t, mock.URL)

	d := ossBucketResourceData(t, map[string]interface{}{"bucket": bucketName})
	err := resourceAlicloudOssBucketCreate(d, client)
	if err == nil {
		t.Fatal("expected an error when the pre-creation check is denied")
	}
	if d.Id() != "" {
		t.Fatalf("expected empty resource ID, got %q", d.Id())
	}
	if n := mock.count(http.MethodPut, "/"+bucketName, ""); n != 0 {
		t.Fatalf("expected no CreateBucket request, got %d", n)
	}
	mock.assertNoListBuckets(t)
}

func TestUnitAliCloudOssBucketCreateFailsOnNameConflict(t *testing.T) {
	bucketName := "mock-conflict-bucket"
	mock := newOssMockServer(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch {
		case hasOssSubresource(q, "bucketInfo"):
			writeOssError(w, http.StatusNotFound, "NoSuchBucket")
		case r.Method == http.MethodPut:
			writeOssError(w, http.StatusConflict, "BucketAlreadyExists")
		default:
			writeOssError(w, http.StatusInternalServerError, "UnexpectedMockRequest")
		}
	})
	defer mock.Close()
	client := newOssMockClient(t, mock.URL)

	d := ossBucketResourceData(t, map[string]interface{}{"bucket": bucketName})
	err := resourceAlicloudOssBucketCreate(d, client)
	if err == nil {
		t.Fatal("expected an error when creation hits a name conflict")
	}
	if !strings.Contains(err.Error(), "not available") {
		t.Fatalf("expected the name-conflict message, got: %v", err)
	}
	if d.Id() != "" {
		t.Fatalf("expected empty resource ID, got %q", d.Id())
	}
	if n := mock.count(http.MethodPut, "/"+bucketName, ""); n != 1 {
		t.Fatalf("expected exactly one CreateBucket request, got %d", n)
	}
	mock.assertNoListBuckets(t)
}

func TestUnitAliCloudOssBucketCreateKeepsIdWhenPostCheckFails(t *testing.T) {
	bucketName := "mock-postcheck-denied-bucket"
	var bucketInfoCalls int32
	mock := newOssMockServer(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch {
		case hasOssSubresource(q, "bucketInfo"):
			if atomic.AddInt32(&bucketInfoCalls, 1) == 1 {
				writeOssError(w, http.StatusNotFound, "NoSuchBucket")
				return
			}
			writeOssError(w, http.StatusForbidden, "AccessDenied")
		case r.Method == http.MethodPut && len(q) == 0:
			writeOssXML(w, http.StatusOK, "")
		default:
			writeOssError(w, http.StatusInternalServerError, "UnexpectedMockRequest")
		}
	})
	defer mock.Close()
	client := newOssMockClient(t, mock.URL)

	d := ossBucketResourceData(t, map[string]interface{}{"bucket": bucketName})
	err := resourceAlicloudOssBucketCreate(d, client)
	if err == nil {
		t.Fatal("expected an error when the post-creation check is denied")
	}
	if d.Id() != bucketName {
		t.Fatalf("expected the created bucket %q to stay tracked, got %q", bucketName, d.Id())
	}
	if n := mock.count(http.MethodGet, "/"+bucketName, "bucketInfo"); n != 2 {
		t.Fatalf("expected no retry on a deterministic post-check failure, got %d GetBucketInfo calls", n)
	}
	mock.assertNoListBuckets(t)
}

func TestUnitAliCloudOssBucketCreateRetriesUntilBucketVisible(t *testing.T) {
	bucketName := "mock-retry-bucket"
	var bucketInfoCalls int32
	mock := newOssMockServer(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch {
		case hasOssSubresource(q, "bucketInfo"):
			if atomic.AddInt32(&bucketInfoCalls, 1) <= 2 {
				writeOssError(w, http.StatusNotFound, "NoSuchBucket")
				return
			}
			writeOssXML(w, http.StatusOK, ossBucketInfoXML(bucketName))
		case r.Method == http.MethodPut && len(q) == 0:
			writeOssXML(w, http.StatusOK, "")
		case serveEmptyOssBucketSubresources(w, r):
		default:
			writeOssError(w, http.StatusInternalServerError, "UnexpectedMockRequest")
		}
	})
	defer mock.Close()
	client := newOssMockClient(t, mock.URL)

	d := ossBucketResourceData(t, map[string]interface{}{"bucket": bucketName, "acl": "private"})
	d.MarkNewResource()
	if err := resourceAlicloudOssBucketCreate(d, client); err != nil {
		t.Fatalf("expected creation to succeed after the bucket becomes visible, got: %v", err)
	}
	if d.Id() != bucketName {
		t.Fatalf("expected resource ID %q, got %q", bucketName, d.Id())
	}
	if n := mock.count(http.MethodGet, "/"+bucketName, "bucketInfo"); n < 3 {
		t.Fatalf("expected the post-creation wait to retry, got %d GetBucketInfo calls", n)
	}
	if n := mock.count(http.MethodPut, "/"+bucketName, ""); n != 1 {
		t.Fatalf("expected a single CreateBucket request on first creation, got %d PUT requests", n)
	}
	if n := mock.count(http.MethodPut, "/"+bucketName, "acl"); n != 0 {
		t.Fatalf("expected no acl update request on first creation, got %d", n)
	}
	mock.assertNoListBuckets(t)
}

func TestUnitAliCloudOssBucketCreateKeepsIdWhenWaitTimesOut(t *testing.T) {
	oldTimeout := ossBucketCreateWaitTimeout
	ossBucketCreateWaitTimeout = 2 * time.Second
	defer func() { ossBucketCreateWaitTimeout = oldTimeout }()

	bucketName := "mock-timeout-bucket"
	mock := newOssMockServer(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch {
		case hasOssSubresource(q, "bucketInfo"):
			writeOssError(w, http.StatusNotFound, "NoSuchBucket")
		case r.Method == http.MethodPut && len(q) == 0:
			writeOssXML(w, http.StatusOK, "")
		default:
			writeOssError(w, http.StatusInternalServerError, "UnexpectedMockRequest")
		}
	})
	defer mock.Close()
	client := newOssMockClient(t, mock.URL)

	d := ossBucketResourceData(t, map[string]interface{}{"bucket": bucketName})
	err := resourceAlicloudOssBucketCreate(d, client)
	if err == nil {
		t.Fatal("expected an error when the post-creation wait times out")
	}
	if d.Id() != bucketName {
		t.Fatalf("expected the created bucket %q to stay tracked, got %q", bucketName, d.Id())
	}
	if n := mock.count(http.MethodGet, "/"+bucketName, "bucketInfo"); n < 3 {
		t.Fatalf("expected the wait to poll until the timeout, got %d GetBucketInfo calls", n)
	}
	mock.assertNoListBuckets(t)
}

func TestUnitAliCloudOssBucketDeleteSucceedsWhenAlreadyAbsent(t *testing.T) {
	bucketName := "mock-absent-bucket"
	mock := newOssMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			writeOssError(w, http.StatusNotFound, "NoSuchBucket")
			return
		}
		writeOssError(w, http.StatusInternalServerError, "UnexpectedMockRequest")
	})
	defer mock.Close()
	client := newOssMockClient(t, mock.URL)

	d := ossBucketResourceData(t, map[string]interface{}{"bucket": bucketName})
	d.SetId(bucketName)
	if err := resourceAlicloudOssBucketDelete(d, client); err != nil {
		t.Fatalf("expected delete of an absent bucket to succeed, got: %v", err)
	}
	if n := mock.count(http.MethodDelete, "/"+bucketName, ""); n != 1 {
		t.Fatalf("expected exactly one DeleteBucket request, got %d", n)
	}
	if n := mock.count(http.MethodGet, "/"+bucketName, "bucketInfo"); n != 0 {
		t.Fatalf("expected the post-delete wait to be skipped, got %d GetBucketInfo calls", n)
	}
	mock.assertNoListBuckets(t)
}

func TestUnitAliCloudOssBucketDeleteFailsWhenAccessDenied(t *testing.T) {
	bucketName := "mock-delete-denied-bucket"
	mock := newOssMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			writeOssError(w, http.StatusForbidden, "AccessDenied")
			return
		}
		writeOssError(w, http.StatusInternalServerError, "UnexpectedMockRequest")
	})
	defer mock.Close()
	client := newOssMockClient(t, mock.URL)

	d := ossBucketResourceData(t, map[string]interface{}{"bucket": bucketName})
	d.SetId(bucketName)
	if err := resourceAlicloudOssBucketDelete(d, client); err == nil {
		t.Fatal("expected an error when delete is denied")
	}
	if n := mock.count(http.MethodDelete, "/"+bucketName, ""); n != 1 {
		t.Fatalf("expected no retry on a deterministic delete failure, got %d DeleteBucket requests", n)
	}
	mock.assertNoListBuckets(t)
}

func TestUnitAliCloudOssBucketDeleteFailsWhenNotEmptyWithoutForceDestroy(t *testing.T) {
	bucketName := "mock-notempty-bucket"
	mock := newOssMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			writeOssError(w, http.StatusConflict, "BucketNotEmpty")
			return
		}
		writeOssError(w, http.StatusInternalServerError, "UnexpectedMockRequest")
	})
	defer mock.Close()
	client := newOssMockClient(t, mock.URL)

	d := ossBucketResourceData(t, map[string]interface{}{"bucket": bucketName, "force_destroy": false})
	d.SetId(bucketName)
	if err := resourceAlicloudOssBucketDelete(d, client); err == nil {
		t.Fatal("expected an error when the bucket is not empty and force_destroy is false")
	}
	if n := mock.count(http.MethodGet, "/"+bucketName, "versions"); n != 0 {
		t.Fatalf("expected no object cleanup without force_destroy, got %d ListObjectVersions requests", n)
	}
	mock.assertNoListBuckets(t)
}

func TestUnitAliCloudOssBucketDeleteForceDestroyVersions(t *testing.T) {
	bucketName := "mock-force-destroy-bucket"
	var deleteCalls int32
	mock := newOssMockServer(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch {
		case r.Method == http.MethodDelete:
			if atomic.AddInt32(&deleteCalls, 1) == 1 {
				writeOssError(w, http.StatusConflict, "BucketNotEmpty")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		case hasOssSubresource(q, "versions"):
			writeOssXML(w, http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?><ListVersionsResult><Name>`+bucketName+`</Name><Version><Key>obj1</Key><VersionId>v1</VersionId><IsLatest>true</IsLatest></Version><DeleteMarker><Key>obj2</Key><VersionId>v2</VersionId><IsLatest>false</IsLatest></DeleteMarker></ListVersionsResult>`)
		case hasOssSubresource(q, "delete"):
			writeOssXML(w, http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?><DeleteResult><Deleted><Key>obj1</Key><VersionId>v1</VersionId></Deleted><Deleted><Key>obj2</Key><VersionId>v2</VersionId></Deleted></DeleteResult>`)
		case hasOssSubresource(q, "bucketInfo"):
			writeOssError(w, http.StatusNotFound, "NoSuchBucket")
		default:
			writeOssError(w, http.StatusInternalServerError, "UnexpectedMockRequest")
		}
	})
	defer mock.Close()
	client := newOssMockClient(t, mock.URL)

	d := ossBucketResourceData(t, map[string]interface{}{"bucket": bucketName, "force_destroy": true})
	d.SetId(bucketName)
	if err := resourceAlicloudOssBucketDelete(d, client); err != nil {
		t.Fatalf("expected force_destroy delete to succeed, got: %v", err)
	}
	if n := mock.count(http.MethodDelete, "/"+bucketName, ""); n != 2 {
		t.Fatalf("expected DeleteBucket to be retried after cleanup, got %d requests", n)
	}
	if n := mock.count(http.MethodGet, "/"+bucketName, "versions"); n != 1 {
		t.Fatalf("expected one ListObjectVersions request, got %d", n)
	}
	if n := mock.count(http.MethodPost, "/"+bucketName, "delete"); n != 1 {
		t.Fatalf("expected one DeleteObjectVersions request, got %d", n)
	}
	mock.assertNoListBuckets(t)
}
