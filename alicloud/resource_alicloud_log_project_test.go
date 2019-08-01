package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_log_project", &resource.Sweeper{
		Name: "alicloud_log_project",
		F:    testSweepLogProjects,
	})
}

func testSweepLogProjects(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
	}

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

func testAccCheckAlicloudLogProjectExists(name string, project *sls.LogProject) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", name))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No Log project ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		logService := LogService{client}

		p, err := logService.DescribeLogProject(rs.Primary.ID)
		if err != nil {
			return WrapError(err)
		}
		if p == nil || p.Name == "" {
			return WrapError(fmt.Errorf("Log project %s is not exist.", rs.Primary.ID))
		}
		project = p

		return nil
	}
}
