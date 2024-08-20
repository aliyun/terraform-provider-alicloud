package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudLogStore_basic(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-log-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependence)

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
					"name":                  name,
					"project":               "${alicloud_log_project.foo.name}",
					"shard_count":           "1",
					"auto_split":            "true",
					"max_split_shard_count": "1",
					"mode":                  "standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"project":               name,
						"shard_count":           "1",
						"auto_split":            "true",
						"max_split_shard_count": "1",
						"mode":                  "standard",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_split_shard_count": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_split_shard_count": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "3000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "3000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "35",
					"hot_ttl":          "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "35",
						"hot_ttl":          "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "30",
					"hot_ttl":          "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "30",
						"hot_ttl":          "0",
					}),
				),
			},
			// TODO: because auto_split and max_split_shard_count affect each other, when auto_split = false, max_split_shard_count will be set to 0, and when updating auto_split = true, max_split_shard_count must be set to be greater than 0, so in the test, auto_split = true in step 0, omitting this step
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"auto_split": "true",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"auto_split": "true",
			// 		}),
			// 	),
			// },
			{
				Config: testAccConfig(map[string]interface{}{
					"append_meta": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"append_meta": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_web_tracking": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_web_tracking": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period":      REMOVEKEY,
					"auto_split":            REMOVEKEY,
					"max_split_shard_count": REMOVEKEY,
					"append_meta":           REMOVEKEY,
					"enable_web_tracking":   REMOVEKEY,
					"hot_ttl":               REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period":      "30",
						"auto_split":            "false",
						"max_split_shard_count": "0",
						"append_meta":           "true",
						"enable_web_tracking":   "false",
						"hot_ttl":               "0",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudLogStore_mode(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-log-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SlsTestRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"project":               "${alicloud_log_project.foo.name}",
					"shard_count":           "1",
					"auto_split":            "true",
					"max_split_shard_count": "1",
					"mode":                  "query",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"project":               name,
						"shard_count":           "1",
						"auto_split":            "true",
						"max_split_shard_count": "1",
						"mode":                  "query",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mode": "standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mode": "standard",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudLogStore_lite(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-log-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SlsTestRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"project":               "${alicloud_log_project.foo.name}",
					"shard_count":           "1",
					"auto_split":            "true",
					"max_split_shard_count": "1",
					"mode":                  "lite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"project":               name,
						"shard_count":           "1",
						"auto_split":            "true",
						"max_split_shard_count": "1",
						"mode":                  "lite",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mode": "standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mode": "standard",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudLogStore_metric(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-metric-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependenceWithEncrypt)

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
					"name":           name,
					"project":        "${alicloud_log_project.foo.name}",
					"shard_count":    "1",
					"telemetry_type": "Metrics",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"project":        name,
						"shard_count":    "1",
						"telemetry_type": "Metrics",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_count": "2",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_count":    "2",
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "m4",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_ttl": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_ttl": "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudLogStore_metric_fix_bug_using_logstoreName(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-metric-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependenceWithEncrypt)

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
					"logstore_name":  name,
					"project":        "${alicloud_log_project.foo.name}",
					"shard_count":    "1",
					"telemetry_type": "Metrics",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name":  name,
						"project":        name,
						"shard_count":    "1",
						"telemetry_type": "Metrics",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_count": "2",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_count":    "2",
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "m4",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_ttl": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_ttl": "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudLogStore_metric_fix_bug_using_projectName(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-metric-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependenceWithEncrypt)

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
					"name":           name,
					"project_name":   "${alicloud_log_project.foo.name}",
					"shard_count":    "1",
					"telemetry_type": "Metrics",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"project_name":   name,
						"shard_count":    "1",
						"telemetry_type": "Metrics",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_count": "2",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_count":    "2",
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "m4",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_ttl": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_ttl": "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudLogStore_metric_fix_bug_using_projectName_and_logstoreName(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-metric-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependenceWithEncrypt)

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
					"logstore_name":  name,
					"project_name":   "${alicloud_log_project.foo.name}",
					"shard_count":    "1",
					"telemetry_type": "Metrics",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name":  name,
						"project_name":   name,
						"shard_count":    "1",
						"telemetry_type": "Metrics",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_count": "2",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_count":    "2",
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "m4",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_ttl": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_ttl": "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudLogStore_create_with_encrypt(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-log-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependenceWithEncrypt)

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
					"name":        name,
					"project":     "${alicloud_log_project.foo.name}",
					"shard_count": "1",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"project":        name,
						"shard_count":    "1",
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "m4",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "false",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudLogStore_create_with_encrypt_bugfix(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-log-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependenceWithEncrypt)

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
					"name":        name,
					"project":     "${alicloud_log_project.foo.name}",
					"shard_count": "1",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"project":        name,
						"shard_count":    "1",
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "m4",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "false",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudLogStore_create_with_encrypt_enable(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-log-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependenceWithEncrypt)

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
					"name":        name,
					"project":     "${alicloud_log_project.foo.name}",
					"shard_count": "1",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]string{
								{
									"cmk_key_id": "${alicloud_kms_key.key.id}",
									"arn":        "acs:ram::${data.alicloud_account.default.id}:role/aliyunlogdefaultrole",
									"region_id":  defaultRegionToTest,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"project":        name,
						"shard_count":    "1",
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_conf": []map[string]interface{}{
						{
							"enable": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_conf.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudLogStore_multi(t *testing.T) {
	var v *sls.LogStore
	resourceId := "alicloud_log_store.default.4"
	ra := resourceAttrInit(resourceId, logStoreMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-log-store-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogStoreConfigDependence)

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
					"name":    name + "${count.index}",
					"project": "${alicloud_log_project.foo.name}",
					"count":   "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceLogStoreConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_project" "foo" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	`, name)
}

func resourceLogStoreConfigDependenceWithEncrypt(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	data "alicloud_account" "default"{
	}

	data "alicloud_vpcs" "default" {
	  name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
	  vpc_id  = data.alicloud_vpcs.default.ids.0
	}
	
	resource "alicloud_kms_instance" "default" {
	  product_version = "3"
	  vpc_id          = data.alicloud_vpcs.default.ids.0
	  zone_ids = [
		data.alicloud_vswitches.default.vswitches.0.zone_id,
        data.alicloud_vswitches.default.vswitches.1.zone_id
	  ]
	  vswitch_ids = [
		data.alicloud_vswitches.default.ids.0
	  ]
	  vpc_num    = "1"
	  key_num    = "1000"
	  secret_num = "0"
	  spec       = "1000"
      force_delete_without_backup = true
      payment_type = "PayAsYouGo"
	}

	resource "alicloud_kms_key" "key" {
  		description             = "${var.name}"
  		pending_window_in_days  = "7"
  		status                  = "Enabled"
        dkms_instance_id        = alicloud_kms_instance.default.id
	}
	resource "alicloud_log_project" "foo" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	`, name)
}

var logStoreMap = map[string]string{
	"name":    CHECKSET,
	"project": CHECKSET,
}

func testAccCheckSLSLogStoreDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckSLSLogStoreDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckSLSLogStoreDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_log_store" {
			continue
		}

		_, err := slsServiceV2.DescribeSlsLogStore(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func TestAccAliCloudSlsLogStore_basic5614_old(t *testing.T) {
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsLogStoreMap5614_old)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslslogstore%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsLogStoreBasicDependence5614_old)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SLSLiteSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSLSLogStoreDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"logstore_name":         name,
					"project_name":          "${alicloud_log_project.defaultbRFbyS.name}",
					"shard_count":           "2",
					"retention_period":      "30",
					"metering_mode":         "ChargeByFunction",
					"telemetry_type":        "None",
					"mode":                  "query",
					"auto_split":            "true",
					"max_split_shard_count": "6",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]interface{}{
								{
									"cmk_key_id": "${alicloud_kms_key.default.id}",
									"region_id":  "cn-hangzhou",
									"arn":        "acs:ram::1511928242963727:role/${alicloud_ram_role.defaultRole.name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name":         name,
						"project_name":          CHECKSET,
						"shard_count":           "2",
						"retention_period":      "30",
						"metering_mode":         "ChargeByFunction",
						"telemetry_type":        "None",
						"mode":                  "query",
						"auto_split":            "true",
						"max_split_shard_count": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_ttl": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_ttl": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_split": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_split": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mode": "query",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mode": "query",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"append_meta": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"append_meta": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_split_shard_count": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_split_shard_count": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mode": "standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mode": "standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"metering_mode": "ChargeByDataIngest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metering_mode": "ChargeByDataIngest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_ttl": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_ttl": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_split_shard_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_split_shard_count": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "31",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "31",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_ttl": "11",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_ttl": "11",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_split": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_split": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"append_meta": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"append_meta": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_split_shard_count": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_split_shard_count": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_web_tracking": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_web_tracking": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_ttl":        "7",
					"logstore_name":  name + "_update",
					"project_name":   "${alicloud_log_project.defaultbRFbyS.name}",
					"shard_count":    "2",
					"auto_split":     "true",
					"mode":           "query",
					"telemetry_type": "None",
					"append_meta":    "true",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "true",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]interface{}{
								{
									"cmk_key_id": "${alicloud_kms_key.default.id}",
									"region_id":  "cn-hangzhou",
									"arn":        "acs:ram::1511928242963727:role/${alicloud_ram_role.defaultRole.name}",
								},
							},
						},
					},
					"max_split_shard_count": "6",
					"retention_period":      "30",
					"metering_mode":         "ChargeByFunction",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_ttl":               "7",
						"logstore_name":         name + "_update",
						"project_name":          CHECKSET,
						"shard_count":           "2",
						"auto_split":            "true",
						"mode":                  "query",
						"telemetry_type":        "None",
						"append_meta":           "true",
						"max_split_shard_count": "6",
						"retention_period":      "30",
						"metering_mode":         "ChargeByFunction",
					}),
				),
			},
		},
	})
}

var AlicloudSlsLogStoreMap5614_old = map[string]string{
	"create_time":    CHECKSET,
	"encrypt_conf.#": CHECKSET,
}

func AlicloudSlsLogStoreBasicDependence5614_old(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "logstore_name" {
  default = "logstore"
}

variable "project_name" {
  default = "terraform-logstore-test"
}

provider "alicloud" {
	alias  = "hz"
	region = "cn-hangzhou"
}

resource "alicloud_log_project" "defaultbRFbyS" {
  description = "terraform-logstore-test"
  name        = var.name

}

data "alicloud_vpcs" "default" {
  provider   = alicloud.hz
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  provider   = alicloud.hz
  vpc_id  = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_kms_instance" "default" {
  provider   = alicloud.hz
  product_version = "3"
  vpc_id          = data.alicloud_vpcs.default.ids.0
  zone_ids = [
	data.alicloud_vswitches.default.vswitches.0.zone_id,
	data.alicloud_vswitches.default.vswitches.1.zone_id
  ]
  vswitch_ids = [
	data.alicloud_vswitches.default.ids.0
  ]
  vpc_num    = "1"
  key_num    = "1000"
  secret_num = "0"
  spec       = "1000"
  force_delete_without_backup = true
  payment_type = "PayAsYouGo"
}

resource "alicloud_kms_key" "default" {
  provider          = alicloud.hz
  description = "Default"
  status = "Enabled"
  pending_window_in_days = 7
  dkms_instance_id = alicloud_kms_instance.default.id
}

resource "alicloud_ram_role" "defaultRole" {
  name = var.name

  description                 = "tf-test-role-two"
  document = <<EOF
{
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
                "Service": [
                    "log.aliyuncs.com"
                ]
            }
        }
    ],
    "Version": "1"
}
  EOF
}

resource "alicloud_ram_role_policy_attachment" "RolePolicyAttachment" {
  policy_type = "System"
  role_name   = alicloud_ram_role.defaultRole.name
  policy_name = "AliyunKMSReadOnlyAccess"
}

resource "alicloud_ram_role_policy_attachment" "default83dHsl" {
  policy_type = "System"
  role_name   = alicloud_ram_role.defaultRole.name
  policy_name = "AliyunKMSCryptoUserAccess"
}

resource "alicloud_ram_role_policy_attachment" "default83dHes" {
  policy_type = "System"
  role_name   = alicloud_ram_role.defaultRole.name
  policy_name = "AliyunLogRolePolicy"
}

resource "alicloud_ram_policy" "defaultLPolicy" {
  description = "tf-test-policy-two"
  policy_name = var.name

  document = <<EOF
{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ram:PassRole"
      ],
      "Resource": [
        "acs::ram::*"
      ],
      "Condition": {}
    }
  ]
}
  EOF
}

`, name)
}

// Case 5614_old  twin
func TestAccAliCloudSlsLogStore_basic5614_old_twin(t *testing.T) {
	resourceId := "alicloud_log_store.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsLogStoreMap5614_old)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslslogstore%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsLogStoreBasicDependence5614_old)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SLSLiteSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSLSLogStoreDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_ttl":        "7",
					"logstore_name":  name,
					"project_name":   "${alicloud_log_project.defaultbRFbyS.name}",
					"shard_count":    "2",
					"auto_split":     "true",
					"mode":           "query",
					"telemetry_type": "None",
					"append_meta":    "true",
					"encrypt_conf": []map[string]interface{}{
						{
							"enable":       "false",
							"encrypt_type": "default",
							"user_cmk_info": []map[string]interface{}{
								{
									"cmk_key_id": "${alicloud_kms_key.default.id}",
									"region_id":  "cn-hangzhou",
									"arn":        "acs:ram::1511928242963727:role/${alicloud_ram_role.defaultRole.name}",
								},
							},
						},
					},
					"max_split_shard_count": "6",
					"retention_period":      "200",
					"enable_web_tracking":   "false",
					"metering_mode":         "ChargeByFunction",
					"infrequent_access_ttl": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_ttl":               "7",
						"logstore_name":         name,
						"project_name":          CHECKSET,
						"shard_count":           "2",
						"auto_split":            "true",
						"mode":                  "query",
						"telemetry_type":        "None",
						"append_meta":           "true",
						"max_split_shard_count": "6",
						"retention_period":      "200",
						"enable_web_tracking":   "false",
						"metering_mode":         "ChargeByFunction",
						"infrequent_access_ttl": "40",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_ttl":               "11",
					"auto_split":            "true",
					"mode":                  "standard",
					"append_meta":           "true",
					"max_split_shard_count": "20",
					"retention_period":      "202",
					"enable_web_tracking":   "true",
					"metering_mode":         "ChargeByDataIngest",
					"infrequent_access_ttl": "43",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_ttl":               "11",
						"auto_split":            "true",
						"mode":                  "standard",
						"append_meta":           "true",
						"max_split_shard_count": "20",
						"retention_period":      "202",
						"enable_web_tracking":   "true",
						"metering_mode":         "ChargeByDataIngest",
						"infrequent_access_ttl": "43",
					}),
				),
			},
		},
	})
}
