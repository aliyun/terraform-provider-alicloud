package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/denverdino/aliyungo/cs"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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
		"tf_test_",
		"tf-test-",
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
