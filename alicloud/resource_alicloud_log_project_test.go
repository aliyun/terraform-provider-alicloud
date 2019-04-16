package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/aliyun-log-go-sdk"
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
	var project sls.LogProject
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogProjectBasic(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.foo", &project),
					resource.TestCheckResourceAttr("alicloud_log_project.foo", "description", "tf unit test"),
				),
			},
			{
				Config: testAccLogProjectUpdate(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudLogProjectExists("alicloud_log_project.foo", &project),
					resource.TestCheckResourceAttr("alicloud_log_project.foo", "description", "tf unit test update"),
				),
			},
		},
	})
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

func testAccCheckAlicloudLogProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_log_project" {
			continue
		}

		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.CheckProjectExist(rs.Primary.ID)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "logtail_project", "CheckProjectExist", AliyunLogGoSdkERROR)
		}
		exist, _ := raw.(bool)
		if !exist {
			return nil
		}

		return WrapError(fmt.Errorf("Log project %s still exists.", rs.Primary.ID))
	}

	return nil
}

func testAccLogProjectBasic(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_log_project" "foo" {
	    name = "tf-testacclogproject-%d"
	    description = "tf unit test"
	}`, rand)
}

func testAccLogProjectUpdate(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_log_project" "foo"{
		name = "tf-testacclogproject-%d"
		description = "tf unit test update"
}
`, rand)
}
