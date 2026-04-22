package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_log_project", &resource.Sweeper{
		Name: "alicloud_log_project",
		F:    testSweepLogProjects,
	})
}

func testSweepLogProjects(region string) error {
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf-example",
		"tf_example",
		"terraform-example",
		"tf_test_",
		"tf-test-",
		"k8s-log-",
		"dbaudit",
		"kms-log-",
	}
	return testSweepLogProjectsWithPrefixAndSuffix(region, prefixes, []string{})
}

func testSweepLogProjectsWithPrefixAndSuffix(region string, prefixes, suffixes []string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return slsClient.ListProject()
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving Log Projects: %s", WrapError(err))
	}
	names, _ := raw.([]string)

	for _, v := range names {
		name := v
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				for _, suffix := range suffixes {
					if strings.HasSuffix(strings.ToLower(name), strings.ToLower(suffix)) {
						skip = false
						break
					}
				}
			}
			// Sweep the project which from the k8s cluster
			if skip && strings.HasPrefix(name, "k8s-log-") {
				k8sId := strings.TrimPrefix(name, "k8s-log-")
				raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
					return csClient.DescribeCluster(k8sId)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
						skip = false
					} else {
						log.Printf("[ERROR] DescribeCluster got an error: %#v", err)
					}
				} else {
					cluster, _ := raw.(cs.ClusterType)
					if strings.HasPrefix(strings.ToLower(cluster.Name), "tf-testacc") || strings.HasPrefix(strings.ToLower(cluster.Name), "tf_testacc") {
						skip = false
					}
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Log Project: %s", name)
				continue
			}
		}
		log.Printf("[INFO] Deleting Log Project: %s", name)
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteProject(name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Log Project (%s): %s", name, err)
		}
	}
	return nil
}
func TestAccAliCloudLogProject_basic(t *testing.T) {
	var v *sls.LogProject
	resourceId := "alicloud_log_project.default"
	ra := resourceAttrInit(resourceId, logProjectMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogproject-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogProjectConfigDependence)

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
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"description": "tf unit test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf unit test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf unit test update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf unit test update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": `{\"Version\":\"1\",\"Statement\":[{\"Resource\":\"acs:log:*:*:project/exampleproject/*\",\"Effect\":\"Deny\",\"Action\":[\"log:PostLogStoreLogs\"],\"Condition\":{\"StringNotLike\":{\"acs:SourceVpc\":[\"vpc-*\"]}}}]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": "{\"Version\":\"1\",\"Statement\":[{\"Resource\":\"acs:log:*:*:project/exampleproject/*\",\"Effect\":\"Deny\",\"Action\":[\"log:PostLogStoreLogs\"],\"Condition\":{\"StringNotLike\":{\"acs:SourceVpc\":[\"vpc-*\"]}}}]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudLogProject_tags(t *testing.T) {
	var v *sls.LogProject
	resourceId := "alicloud_log_project.default"
	ra := resourceAttrInit(resourceId, logProjectMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogproject-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogProjectConfigDependence)

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
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
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
		},
	})
}

func TestAccAliCloudLogProject_multi(t *testing.T) {
	var v *sls.LogProject
	resourceId := "alicloud_log_project.default.4"
	ra := resourceAttrInit(resourceId, logProjectMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogproject-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogProjectConfigDependence)

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
					"name":  name + "${count.index}",
					"count": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceLogProjectConfigDependence(name string) string {
	return ""
}

var logProjectMap = map[string]string{
	"name": CHECKSET,
}

// Test Sls Project. >>> Resource test cases, automatically generated.
// Case 4209
func TestAccAliCloudSlsProject_basic4209(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_log_project.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsProjectMap4209)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsproject%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsProjectBasicDependence4209)
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
					"project_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.project_desc}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": CHECKSET,
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
					"description": "${var.project_desc2}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": CHECKSET,
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
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.project_desc}",
					"project_name":      name + "update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       CHECKSET,
						"project_name":      name + "update",
						"resource_group_id": CHECKSET,
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
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudSlsProjectMap4209 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudSlsProjectBasicDependence4209(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "project_desc2" {
  default = "test2"
}

variable "project_desc" {
  default = "project for test"
}

variable "project_name" {
  default = "project-test-2023823"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case 4209  twin
func TestAccAliCloudSlsProject_basic4209_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_log_project.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsProjectMap4209)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsproject%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsProjectBasicDependence4209)
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
					"description":       "${var.project_desc2}",
					"project_name":      name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       CHECKSET,
						"project_name":      name,
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
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

// Test Sls Project. <<< Resource test cases, automatically generated.

// TestAccAliCloudSlsProject_signVersionV4 is the path-A representative for
// v4 signature regression testing. It exercises the SDK v2 client.Do("Sls",
// ...) path under a provider configured with sign_version.sls = "v4" (and
// sign_version.oss = "v4"), matching the real-world provider.tf that first
// surfaced the bug:
//
//	provider "alicloud" {
//	  region = "cn-hangzhou"
//	  endpoints { log = "cn-hangzhou-acdr-ut-3.log.aliyuncs.com" }
//	  sign_version { sls = "v4" oss = "v4" }
//	}
//
// Prior to the fix in applyOpenapiSignVersion, sign_version.sls was silently
// dropped, and the acdr-ut endpoint rejected every request with
// SignatureVersionNotSupported. This test pins the hardened wiring against
// a real acdr-ut endpoint end-to-end.
//
// The companion path-B representative is TestAccAlicloudLogMachineGroup_-
// signVersionV4 (alicloud/resource_alicloud_log_machine_group_test.go),
// which exercises the legacy WithLogClient path. See that test's doc
// comment for the full coverage map across all 22 alicloud_sls_* /
// alicloud_log_* resources.
//
// The endpoint is hard-coded to cn-hangzhou-acdr-ut-3.log.aliyuncs.com
// because only a v4-only endpoint is a meaningful environment for this
// test: public regional endpoints (<region>.log.aliyuncs.com) accept v1
// and have inconsistent v4 coverage across individual SLS APIs (we
// empirically saw GetProjectPolicy return "ParameterInvalid: The signing
// region in credential is invalid" under v4 against eu-central-1 even
// though the credential-scope region matched the endpoint), which would
// produce flakes unrelated to the provider wiring we are actually testing.
// Running this test requires ALICLOUD_ACCESS_KEY whose account has access
// to cn-hangzhou-acdr-ut-3.
func TestAccAliCloudSlsProject_signVersionV4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_log_project.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsProjectMap4209)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsprojv4%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsProjectSignVersionV4Dependence)
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
					"project_name": name,
					"description":  "${var.project_desc}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name": name,
						"description":  CHECKSET,
					}),
				),
			},
			// ImportState is intentionally omitted: Terraform SDK v1's
			// EvalImportStateVerify builds a fresh eval context for the
			// post-import refresh, and inline `provider "alicloud" {...}`
			// blocks do not reliably propagate custom endpoints /
			// sign_version into it — Read then hits the default endpoint
			// and returns ProjectNotExist. The signature-version wiring is
			// already validated end-to-end by the Create step above and
			// unit-tested in TestApplyOpenapiSignVersion /
			// TestApplyLogClientSignVersion.
		},
	})
}

// AlicloudSlsProjectSignVersionV4Dependence injects an inline provider block
// targeting a v4-only SLS endpoint (cn-hangzhou-acdr-ut-3) so the test
// exercises the end-to-end wiring schema -> config.SignVersion -> openapi
// client that applyOpenapiSignVersion is responsible for.
//
// The endpoint, region and sign_version are all hard-coded. They must stay
// consistent: v4 signing puts a credential-scope region into the signature
// that the server validates against the endpoint-derived region — if the
// two drift (e.g. provider.region = "eu-central-1" but endpoint points at
// cn-hangzhou-acdr-ut-3) the server returns "ParameterInvalid: The signing
// region in credential is invalid". The provider code already derives the
// credential-scope region from the endpoint; the provider.region here is
// the account-level region the acdr-ut endpoint belongs to.
func AlicloudSlsProjectSignVersionV4Dependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "project_desc" {
  default = "tf acc test for sls sign_version v4"
}

provider "alicloud" {
  skip_region_validation = true
  region = "cn-hangzhou-acdr-ut-3"
  endpoints {
    log = "cn-hangzhou-acdr-ut-3.log.aliyuncs.com"
  }
  sign_version {
    sls = "v4"
    oss = "v4"
  }
}
`, name)
}
