package sls_test

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_sls_collection_policy", &resource.Sweeper{
		Name: "alicloud_sls_collection_policy",
		F:    testSweepSlsCollectionPolicies,
	})
}

func testSweepSlsCollectionPolicies(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"xc-test-oss-",
	}

	action := "/collectionpolicy"
	query := make(map[string]*string)
	query["size"] = helper.StringPointer("100")
	query["offset"] = helper.StringPointer("0")
	hostMap := make(map[string]*string)

	response, err := client.Do("Sls", helper.RoaParam("GET", "2020-12-30", "ListCollectionPolicies", action), query, nil, nil, hostMap, true)
	if err != nil {
		return fmt.Errorf("error listing SLS collection policies: %s", err)
	}

	dataRaw, ok := response["data"]
	if !ok || dataRaw == nil {
		return nil
	}

	policies, ok := dataRaw.([]interface{})
	if !ok {
		return nil
	}

	for _, p := range policies {
		policy, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		name, ok := policy["policyName"].(string)
		if !ok {
			continue
		}
		if internal, ok := policy["internalPolicy"].(bool); ok && internal {
			continue
		}
		shouldDelete := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) {
				shouldDelete = true
				break
			}
		}
		if !shouldDelete {
			continue
		}
		productCode, _ := policy["productCode"].(string)
		dataCode, _ := policy["dataCode"].(string)
		log.Printf("[INFO] Deleting SLS collection policy: %s", name)
		deleteAction := fmt.Sprintf("/collectionpolicy/%s", name)
		deleteQuery := make(map[string]*string)
		if productCode != "" {
			deleteQuery["productCode"] = helper.StringPointer(productCode)
		}
		if dataCode != "" {
			deleteQuery["dataCode"] = helper.StringPointer(dataCode)
		}
		_, delErr := client.Do("Sls", helper.RoaParam("DELETE", "2020-12-30", "DeleteCollectionPolicy", deleteAction), deleteQuery, nil, nil, hostMap, false)
		if delErr != nil {
			log.Printf("[ERROR] Failed to delete SLS collection policy (%s): %s", name, delErr)
		}
	}
	return nil
}

// Test Sls CollectionPolicy. >>> Resource test cases, automatically generated.
// Case TF验收_资源目录 8160
func TestAccAliCloudSlsCollectionPolicy_basic8160(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8160)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-oss-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8160)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  policyName,
					"enabled":      "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":    "metering_log",
						"product_code": "oss",
						"policy_name":  policyName,
						"enabled":      "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_enabled": "true",
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou", "cn-shanghai", "cn-beijing"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.update_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore002.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.update_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore002.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
							"members":            []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "true",
					"product_code":       "oss",
					"policy_name":        policyName,
					"enabled":            "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "true",
						"product_code":       "oss",
						"policy_name":        policyName,
						"enabled":            "true",
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

var AlicloudSlsCollectionPolicyMap8160 = map[string]string{}

func AlicloudSlsCollectionPolicyBasicDependence8160(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "project_create_01" {
  description = "xxccc"
  name        = var.name
}

resource "alicloud_log_store" "logstore_create_01" {
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.project_create_01.name
  name             = format("%%s1", var.name)
}

resource "alicloud_log_project" "update_01" {
  description = "dfdfd"
  name        = format("%%s2", var.name)
}

resource "alicloud_log_store" "logstore002" {
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.update_01.name
  name             = format("%%s3", var.name)
}


`, name)
}

// Case TF验收 8140
func TestAccAliCloudSlsCollectionPolicy_basic8140(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8140)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-oss-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8140)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"instance_ids": []string{
								"xcd", "sdd"},
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  policyName,
					"enabled":      "false",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":    "metering_log",
						"product_code": "oss",
						"policy_name":  policyName,
						"enabled":      "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
						},
					},
					"centralize_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
								"\"xc-sx\"":   "xxc",
							},
							"regions": []string{
								"cn-hangzhou", "cn-shanghai", "cn-beijing"},
							"instance_ids": []string{
								"xcd", "sdd", "des"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project02.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore002.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
								"\"xc-sx\"":   "xxc",
							},
							"regions": []string{
								"cn-hangzhou", "cn-beijing"},
							"instance_ids": []string{
								"xcd", "des"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"regions":       []string{},
							"instance_ids":  []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"instance_ids": []string{
								"xcd", "sdd"},
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "false",
					"product_code":       "oss",
					"policy_name":        policyName,
					"enabled":            "false",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "false",
						"product_code":       "oss",
						"policy_name":        policyName,
						"enabled":            "false",
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

var AlicloudSlsCollectionPolicyMap8140 = map[string]string{}

func AlicloudSlsCollectionPolicyBasicDependence8140(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "all_project01" {
  description = "xccx"
  name        = var.name
}

resource "alicloud_log_store" "all_logstore_01" {
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.all_project01.name
  name             = format("%%s1", var.name)
}

resource "alicloud_log_project" "all_project02" {
  description = "cvcvc"
  name        = format("%%s2", var.name)
}

resource "alicloud_log_store" "all_logstore002" {
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.all_project02.name
  name             = format("%%s3", var.name)
}


`, name)
}

// Case TF验收_ResourceMode 7850
func TestAccAliCloudSlsCollectionPolicy_basic7850(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap7850)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-project-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence7850)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
						},
					},
					"data_code":    "operation_log",
					"product_code": "project",
					"policy_name":  policyName,
					"enabled":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":    "operation_log",
						"product_code": "project",
						"policy_name":  policyName,
						"enabled":      "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"centralize_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "attributeMode",
							"resource_tags": map[string]interface{}{
								"\"xc-pro\"": "xc-sd",
							},
							"regions": []string{
								"cn-hangzhou", "cn-beijing"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "instanceMode",
							"regions":       []string{},
							"instance_ids": []string{
								"xxx"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
						},
					},
					"data_code":          "operation_log",
					"centralize_enabled": "false",
					"product_code":       "project",
					"policy_name":        policyName,
					"enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "operation_log",
						"centralize_enabled": "false",
						"product_code":       "project",
						"policy_name":        policyName,
						"enabled":            "false",
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

var AlicloudSlsCollectionPolicyMap7850 = map[string]string{}

func AlicloudSlsCollectionPolicyBasicDependence7850(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case TF验收_资源目录 8160  twin
func TestAccAliCloudSlsCollectionPolicy_basic8160_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8160)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-oss-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8160)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "true",
					"product_code":       "oss",
					"policy_name":        policyName,
					"enabled":            "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "true",
						"product_code":       "oss",
						"policy_name":        policyName,
						"enabled":            "true",
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

// Case TF验收 8140  twin
func TestAccAliCloudSlsCollectionPolicy_basic8140_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8140)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-oss-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8140)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"instance_ids": []string{
								"xcd", "sdd"},
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "false",
					"product_code":       "oss",
					"policy_name":        policyName,
					"enabled":            "false",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "false",
						"product_code":       "oss",
						"policy_name":        policyName,
						"enabled":            "false",
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

// Case TF验收_ResourceMode 7850  twin
func TestAccAliCloudSlsCollectionPolicy_basic7850_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap7850)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-project-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence7850)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
						},
					},
					"data_code":          "operation_log",
					"centralize_enabled": "false",
					"product_code":       "project",
					"policy_name":        policyName,
					"enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "operation_log",
						"centralize_enabled": "false",
						"product_code":       "project",
						"policy_name":        policyName,
						"enabled":            "false",
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

// Case TF验收_资源目录 8160  raw
func TestAccAliCloudSlsCollectionPolicy_basic8160_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8160)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-oss-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8160)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "true",
					"product_code":       "oss",
					"policy_name":        policyName,
					"enabled":            "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "true",
						"product_code":       "oss",
						"policy_name":        policyName,
						"enabled":            "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou", "cn-shanghai", "cn-beijing"},
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.update_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore002.name}",
							"dest_ttl":      "7",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
							"dest_ttl":      "7",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.update_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore002.name}",
							"dest_ttl":      "7",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
							"members":            []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

// Case TF验收 8140  raw
func TestAccAliCloudSlsCollectionPolicy_basic8140_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8140)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-oss-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8140)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"instance_ids": []string{
								"xcd", "sdd"},
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":          "metering_log",
					"centralize_enabled": "false",
					"product_code":       "oss",
					"policy_name":        policyName,
					"enabled":            "false",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "metering_log",
						"centralize_enabled": "false",
						"product_code":       "oss",
						"policy_name":        policyName,
						"enabled":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
								"\"xc-sx\"":   "xxc",
							},
							"regions": []string{
								"cn-hangzhou", "cn-shanghai", "cn-beijing"},
							"instance_ids": []string{
								"xcd", "sdd", "des"},
						},
					},
					"centralize_enabled": "true",
					"enabled":            "true",
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project02.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore002.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "true",
						"enabled":            "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
								"\"xc-sx\"":   "xxc",
							},
							"regions": []string{
								"cn-hangzhou", "cn-beijing"},
							"instance_ids": []string{
								"xcd", "des"},
						},
					},
					"centralize_enabled": "false",
					"enabled":            "false",
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"centralize_enabled": "false",
						"enabled":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"regions":       []string{},
							"instance_ids":  []string{},
						},
					},
					"centralize_config": []map[string]interface{}{
						{
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.all_project01.name}",
							"dest_logstore": "${alicloud_log_store.all_logstore_01.name}",
							"dest_ttl":      "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

// Case TF验收_ResourceMode 7850  raw
func TestAccAliCloudSlsCollectionPolicy_basic7850_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap7850)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-project-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence7850)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
						},
					},
					"data_code":          "operation_log",
					"centralize_enabled": "false",
					"product_code":       "project",
					"policy_name":        policyName,
					"enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":          "operation_log",
						"centralize_enabled": "false",
						"product_code":       "project",
						"policy_name":        policyName,
						"enabled":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "attributeMode",
							"resource_tags": map[string]interface{}{
								"\"xc-pro\"": "xc-sd",
							},
							"regions": []string{
								"cn-hangzhou", "cn-beijing"},
						},
					},
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "instanceMode",
							"regions":       []string{},
							"instance_ids": []string{
								"xxx"},
						},
					},
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
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

// Case TypeListOrderCoverage memberOrder
// Convergent member-order coverage for the nested TypeLists
// policy_config.instance_ids, policy_config.regions and
// resource_directory.members. Each nested list gets its own 2-step
// function (pattern B): step1 applies a full base config with the
// target list in order A (>=2 distinct members) plus all Required
// attributes; step2 re-applies the same full base with ONLY the target
// nested list reordered (everything else byte-identical). The empty
// post-apply plan proves the resource converged on the reordered
// members.
func TestAccAliCloudSlsCollectionPolicy_memberOrder_policyConfigInstanceIds(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8160)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-oss-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8160)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions":      []string{"cn-hangzhou", "cn-shanghai"},
							"instance_ids": []string{"xcd", "sdd"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  policyName,
					"enabled":      "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":    "metering_log",
						"product_code": "oss",
						"policy_name":  policyName,
						"enabled":      "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions":      []string{"cn-hangzhou", "cn-shanghai"},
							"instance_ids": []string{"sdd", "xcd"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  policyName,
					"enabled":      "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

func TestAccAliCloudSlsCollectionPolicy_memberOrder_policyConfigRegions(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8160)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-oss-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8160)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions":      []string{"cn-hangzhou", "cn-shanghai"},
							"instance_ids": []string{"xcd", "sdd"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  policyName,
					"enabled":      "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":    "metering_log",
						"product_code": "oss",
						"policy_name":  policyName,
						"enabled":      "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions":      []string{"cn-shanghai", "cn-hangzhou"},
							"instance_ids": []string{"xcd", "sdd"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  policyName,
					"enabled":      "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

func TestAccAliCloudSlsCollectionPolicy_memberOrder_resourceDirectoryMembers(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8160)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-oss-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8160)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{"cn-hangzhou"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  policyName,
					"enabled":      "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_enabled": "true",
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "custom",
							"members":            []string{"1458992460067836", "1184291612526208"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":    "metering_log",
						"product_code": "oss",
						"policy_name":  policyName,
						"enabled":      "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{"cn-hangzhou"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  policyName,
					"enabled":      "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
					"centralize_enabled": "true",
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "custom",
							"members":            []string{"1184291612526208", "1458992460067836"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

// Case resource_directory remove detection
// Reproduces the reported defect: removing the resource_directory
// block from configuration must be detected as a change and trigger an
// Update. Previously the Optional+Computed schema caused HasChange to
// return false so the removal was silently ignored. After removing
// Computed, HasChange returns true and the Update calls
// UpsertCollectionPolicy. The SLS API rejects removing resource_directory
// from an existing policy with NotAllowed ("you can't change from resource
// directory to non-resource-directory, you should delete the policy first"),
// so the removal step expects that error instead of silent success. After
// the rejected removal, refresh-only and plan-only steps verify that Read
// restores the remote resource_directory into state and that the residual
// drift stays visible in the plan.
func TestAccAliCloudSlsCollectionPolicy_resourceDirectoryRemove(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_collection_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsCollectionPolicyMap8160)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsCollectionPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslscollectionpolicy%d", acctest.DefaultRegionToTest, rand)
	policyName := fmt.Sprintf("xc-test-rd-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsCollectionPolicyBasicDependence8160)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Step 1: Create without resource_directory
				Config: testAccConfig(map[string]interface{}{
					"policy_config": []map[string]interface{}{
						{
							"resource_mode": "all",
							"resource_tags": map[string]interface{}{
								"\"xc-test\"": "xc-test-val",
							},
							"regions": []string{
								"cn-hangzhou"},
						},
					},
					"data_code":    "metering_log",
					"product_code": "oss",
					"policy_name":  policyName,
					"enabled":      "true",
					"data_config": []map[string]interface{}{
						{
							"data_region": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_code":    "metering_log",
						"product_code": "oss",
						"policy_name":  policyName,
						"enabled":      "true",
					}),
					resource.TestCheckResourceAttr(resourceId, "resource_directory.#", "0"),
				),
			},
			{
				// Step 2: Add centralize_config and resource_directory (enable multi-account)
				// dest_region must match the test region enforced by PreCheckWithRegions
				// (cn-shanghai), not acctest.DefaultRegionToTest which is captured
				// before PreCheck runs and may still hold the init default cn-beijing.
				// account_group_type "all" avoids depending on specific member UIDs
				// that may be removed from the resource directory between runs.
				Config: testAccConfig(map[string]interface{}{
					"centralize_enabled": "true",
					"centralize_config": []map[string]interface{}{
						{
							"dest_ttl":      "3",
							"dest_region":   "cn-shanghai",
							"dest_project":  "${alicloud_log_project.project_create_01.name}",
							"dest_logstore": "${alicloud_log_store.logstore_create_01.name}",
						},
					},
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
					resource.TestCheckResourceAttr(resourceId, "resource_directory.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "resource_directory.0.account_group_type", "all"),
				),
			},
			{
				// Step 3: Remove resource_directory block - HasChange must detect
				// the removal (the core fix) and trigger an Update. The SLS API
				// rejects removing resource_directory from an existing policy:
				// NotAllowed "you can't change from resource directory to
				// non-resource-directory, you should delete the policy first".
				// Expecting this error confirms the Update was actually invoked
				// (previously the Computed schema suppressed HasChange so the
				// removal was silently ignored with no API call at all).
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": acctest.REMOVEKEY,
				}),
				// Match the exact rejection reason, not just the error code,
				// so the assertion fails if the API starts rejecting removal
				// for a different cause.
				ExpectError: regexp.MustCompile(`you can't change from resource directory to non-resource-directory`),
			},
			{
				// Step 3.5: Refresh-only verification that Read restores the
				// remote resource_directory. Step 3's failed Apply left
				// state.resource_directory empty (the framework does not
				// refresh after a failed apply, and the diff's removal had
				// already been applied to state), while the remote policy
				// still holds the Step 2 value because the API rejected the
				// removal. Read must populate resource_directory from the API
				// response regardless of prior state; the pre-fix guard keyed
				// on d.GetOk("resource_directory") skipped the write here,
				// leaving state silently diverged from the remote policy.
				// The step keeps Step 3's configuration (block still removed),
				// so the post-refresh plan is expected to stay non-empty.
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "resource_directory.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "resource_directory.0.account_group_type", "all"),
				),
			},
			{
				// Step 3.6: Plan-only check with the block still removed from
				// configuration. The step framework refreshes before planning;
				// the refreshed state carries the remote resource_directory, so
				// the plan must remain non-empty (residual drift stays
				// visible). This fails closed: if Read stopped restoring the
				// remote value, the plan would be empty and the framework
				// errors with "Expected a non-empty plan, but got an empty
				// plan".
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": acctest.REMOVEKEY,
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Step 4: Re-add resource_directory (same as Step 2). Steps
				// 3.5/3.6 already refreshed state from the remote policy, so
				// this configuration matches state: the plan is empty, no
				// Update is issued, and the resource converges. The checks
				// confirm state keeps reflecting the remote policy after
				// re-declaring the block, ending the test with state matching
				// configuration.
				Config: testAccConfig(map[string]interface{}{
					"resource_directory": []map[string]interface{}{
						{
							"account_group_type": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
					resource.TestCheckResourceAttr(resourceId, "resource_directory.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "resource_directory.0.account_group_type", "all"),
				),
			},
		},
	})
}

// Test Sls CollectionPolicy. <<< Resource test cases, automatically generated.
