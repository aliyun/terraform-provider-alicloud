package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers("alicloud_maxcompute_project", &resource.Sweeper{
		Name: "alicloud_maxcompute_project",
		F:    testSweepMaxComputeProjects,
	})
}

// testSweepMaxComputeProjects lists every MaxCompute project in the region and
// deletes those created by acceptance tests. The NODELETING / do-not-delete
// name guard (protectedSweepVpcName) runs in BOTH sweepAll and normal mode so
// protected projects are never touched even when ALICLOUD_SWEEP_ALL_RESOURCES
// is set. Projects that the API refuses to delete (e.g. because they still
// contain tables/data) are logged and left as leftovers — they are NOT
// force-deleted, per the no-bypass-of-managed-protection red line.
func testSweepMaxComputeProjects(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	// List projects: GET /api/v1/projects (paginated by "marker").
	// Mirrors data_source_alicloud_maxcompute_projects.go which uses the same
	// RoaGet call and $.data.projects[*] / $.data.marker paths.
	action := "/api/v1/projects"
	objects := make([]interface{}, 0)
	query := make(map[string]*string)
	query["maxItem"] = StringPointer(strconv.Itoa(PageSizeLarge))
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RoaGet("MaxCompute", "2022-01-04", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, nil)
		if err != nil {
			log.Printf("[ERROR] ListMaxComputeProjects got an error: %v", err)
			break
		}
		resp, jsonErr := jsonpath.Get("$.data.projects[*]", response)
		if jsonErr != nil || resp == nil {
			break
		}
		result, _ := resp.([]interface{})
		objects = append(objects, result...)
		marker, _ := jsonpath.Get("$.data.marker", response)
		if nextToken, ok := marker.(string); ok && nextToken != "" {
			query["marker"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	for _, v := range objects {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["name"])

		// NODELETING / do-not-delete guard — runs in BOTH sweepAll and normal
		// mode (see hbase-sub #10923 / rds-sub #10924 for the pattern).
		if protectedSweepVpcName(name) {
			log.Printf("[INFO] Skipping MaxCompute Project (protected name): %s", name)
			continue
		}
		if !sweepAll() {
			matches := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					matches = true
					break
				}
			}
			if !matches {
				log.Printf("[INFO] Skipping MaxCompute Project (name does not match test prefixes): %s", name)
				continue
			}
		}
		// OLDER_THAN age filter. The list API returns "createdTime", which
		// sweepRecentCheck does not look for by default, so map it into a
		// recognized key. When ALICLOUD_SWEEP_OLDER_THAN is unset this is a
		// no-op; when set, it protects projects newer than the threshold.
		itemCopy := make(map[string]interface{}, len(item))
		for k, val := range item {
			itemCopy[k] = val
		}
		if ct, ok := item["createdTime"]; ok && ct != nil {
			itemCopy["CreateTime"] = ct
		}
		if skip, reason := sweepRecentCheck(itemCopy); skip {
			log.Printf("[INFO] Skipping MaxCompute Project (%s): %s", reason, name)
			continue
		}

		log.Printf("[INFO] Deleting MaxCompute Project: %s", name)
		// Mirror resourceAliCloudMaxComputeProjectDelete: DELETE /api/v1/projects/{name}
		// with isLogical=false (physical delete — the resource default).
		deleteAction := fmt.Sprintf("/api/v1/projects/%s", name)
		deleteQuery := make(map[string]*string)
		deleteQuery["isLogical"] = StringPointer("false")
		deleteBody := make(map[string]interface{})
		deleteBody["projectName"] = name
		if delErr := retry.Retry(5*time.Minute, func() *retry.RetryError {
			_, err := client.RoaDelete("MaxCompute", "2022-01-04", deleteAction, deleteQuery, nil, deleteBody, true)
			if err != nil {
				if NeedRetry(err) {
					wait := incrementalWait(3*time.Second, 5*time.Second)
					wait()
					return retry.RetryableError(err)
				}
				if NotFoundError(err) {
					return nil
				}
				return retry.NonRetryableError(err)
			}
			return nil
		}); delErr != nil {
			// Non-fatal: log and continue. A project that the API refuses to
			// delete (e.g. contains data/tables) is a documented leftover.
			log.Printf("[ERROR] Failed to delete MaxCompute Project %s: %v", name, delErr)
		}
	}

	return nil
}

// Case 1
func TestAccAliCloudMaxComputeProject_basic1968(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, AliCloudMaxComputeProjectMap1968)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)
	name := fmt.Sprintf("tf_testaccmp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMaxComputeProjectBasicDependence1968)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"default_quota":    "默认后付费Quota",
					"project_name":     "${var.name}",
					"comment":          "${var.name}",
					"product_type":     "PayAsYouGo",
					"is_logical":       "false",
					"three_tier_model": "true",
					"ip_white_list": []map[string]interface{}{
						{
							"ip_list":     "1.1.1.1,2.2.2.2",
							"vpc_ip_list": "10.10.10.10,11.11.11.11",
						},
					},
					"properties": []map[string]interface{}{
						{
							"allow_full_scan":  "false",
							"enable_dr":        "true",
							"enable_decimal2":  "true",
							"retention_days":   "1",
							"sql_metering_max": "0",
							"timezone":         "Asia/Shanghai",
							"type_system":      "2",
							"encryption": []map[string]interface{}{
								{
									"enable":    "true",
									"algorithm": "AESCTR",
									"key":       "f58d854d-7bc0-4a6e-9205-160e10ffedec",
								},
							},
							"table_lifecycle": []map[string]interface{}{
								{
									"type":  "optional",
									"value": "37231",
								},
							},
						},
					},
					"security_properties": []map[string]interface{}{
						{
							"enable_download_privilege":            "false",
							"label_security":                       "true",
							"object_creator_has_access_permission": "true",
							"object_creator_has_grant_permission":  "true",
							"project_protection": []map[string]interface{}{
								{
									"protected": "false",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_quota":         CHECKSET,
						"project_name":          CHECKSET,
						"comment":               CHECKSET,
						"product_type":          CHECKSET,
						"three_tier_model":      "true",
						"ip_white_list.#":       "1",
						"properties.#":          "1",
						"security_properties.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":      "${var.name}_u",
					"project_name": "${var.name}_u",
					"is_logical":   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":      CHECKSET,
						"project_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_logical": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "is_logical"},
			},
		},
	})
}

func TestAccAliCloudMaxComputeProject_basic1968_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, AliCloudMaxComputeProjectMap1968)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)
	name := fmt.Sprintf("tf_testaccmp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMaxComputeProjectBasicDependence1968)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"default_quota": "默认后付费Quota",
					"project_name":  "${var.name}",
					"comment":       "${var.name}",
					"product_type":  "PayAsYouGo",
					"is_logical":    "false",
					"ip_white_list": []map[string]interface{}{
						{
							"ip_list":     "1.1.1.1,2.2.2.2",
							"vpc_ip_list": "10.10.10.10,11.11.11.11",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_quota":   CHECKSET,
						"project_name":    CHECKSET,
						"comment":         CHECKSET,
						"product_type":    CHECKSET,
						"ip_white_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"three_tier_model": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"three_tier_model": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_logical": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "is_logical"},
			},
		},
	})
}

var AliCloudMaxComputeProjectMap1968 = map[string]string{}

func AliCloudMaxComputeProjectBasicDependence1968(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

var AliCloudMaxComputeProjectMap7168 = map[string]string{
	"status":       CHECKSET,
	"owner":        CHECKSET,
	"project_name": CHECKSET,
	"create_time":  CHECKSET,
	"type":         CHECKSET,
}

func AliCloudMaxComputeProjectBasicDependence7168(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case Terraform发布#57547636 7168  raw
func TestAccAliCloudMaxComputeProject_basic7168_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, AliCloudMaxComputeProjectMap7168)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccmp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMaxComputeProjectBasicDependence7168)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"default_quota": "默认后付费Quota",
					"project_name":  name,
					"comment":       "terraform测试项目",
					"properties": []map[string]interface{}{
						{
							"type_system":      "2",
							"sql_metering_max": "10240",
							"encryption": []map[string]interface{}{
								{
									"enable":    "true",
									"algorithm": "AESCTR",
									"key":       "f58d854d-7bc0-4a6e-9205-160e10ffedec",
								},
							},
						},
					},
					"status": "AVAILABLE",
					"ip_white_list": []map[string]interface{}{
						{
							"ip_list":     "10.0.0.0/8",
							"vpc_ip_list": "10.0.0.0/8",
						},
					},
					"security_properties": []map[string]interface{}{
						{
							"using_acl":                            "true",
							"using_policy":                         "true",
							"object_creator_has_access_permission": "true",
							"object_creator_has_grant_permission":  "true",
							"label_security":                       "true",
							"enable_download_privilege":            "true",
							"project_protection": []map[string]interface{}{
								{
									"protected":        "true",
									"exception_policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"odps:*\\\"],\\\"Resource\\\":[\\\"acs:odps:*:projects/ludong/tables/*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":[\\\"ALIYUN$ludong@aliyun.com\\\"]}]}",
								},
							},
						},
					},
					"tags": map[string]string{
						"Created": "TF-CI",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_quota": "默认后付费Quota",
						"project_name":  name,
						"comment":       "terraform测试项目",
						"status":        "AVAILABLE",
						"tags.%":        "2",
						"tags.Created":  "TF-CI",
						"tags.For":      "Test",
					}),
				),
			},
			{
				// service denies SetProperty once the project is READONLY/FROZEN, so
				// meta changes are applied while AVAILABLE and the status flips alone
				Config: testAccConfig(map[string]interface{}{
					"ip_white_list": []map[string]interface{}{
						{
							"ip_list":     "0.0.0.0/0",
							"vpc_ip_list": "0.0.0.0/0",
						},
					},
					"security_properties": []map[string]interface{}{
						{
							"using_acl":                            "true",
							"using_policy":                         "true",
							"object_creator_has_access_permission": "true",
							"object_creator_has_grant_permission":  "true",
							"label_security":                       "true",
							"enable_download_privilege":            "true",
							"project_protection": []map[string]interface{}{
								{
									"protected":        "true",
									"exception_policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"odps:*\\\"],\\\"Resource\\\":[\\\"acs:odps:*:projects/ludong/tables/*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":[\\\"ALIYUN$ludong@aliyun.com\\\",\\\"ALIYUN$liuhao@aliyun.com\\\"]}]}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "AVAILABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "READONLY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "READONLY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "FROZEN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// the service stores user-initiated freeze as FROZEN_BY_USER while the
						// schema only accepts FROZEN as input, so the post-apply plan is not empty
						"status": "FROZEN_BY_USER",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "AVAILABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "AVAILABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_properties": []map[string]interface{}{
						{
							"project_protection": []map[string]interface{}{
								{
									"protected": "false",
								},
							},
						},
					},
					"properties": []map[string]interface{}{
						{
							"enable_decimal2":  "false",
							"type_system":      "1",
							"sql_metering_max": "1024",
							"retention_days":   "2",
							"allow_full_scan":  "false",
							"enable_dr":        "true",
							"timezone":         "Asia/Shanghai",
							"table_lifecycle": []map[string]interface{}{
								{
									"type":  "optional",
									"value": "37231",
								},
							},
							"encryption": []map[string]interface{}{
								{
									"enable":    "true",
									"algorithm": "AESCTR",
									"key":       "f58d854d-7bc0-4a6e-9205-160e10ffedec",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":       "AVAILABLE",
						"properties.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type"},
			},
		},
	})
}
