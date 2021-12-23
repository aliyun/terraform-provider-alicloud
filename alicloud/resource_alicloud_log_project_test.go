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
func TestAccAlicloudLogProject_basic(t *testing.T) {
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
		},
	})
}

func TestAccAlicloudLogProject_tags(t *testing.T) {
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

func TestAccAlicloudLogProject_multi(t *testing.T) {
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
