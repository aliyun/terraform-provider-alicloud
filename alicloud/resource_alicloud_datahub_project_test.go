package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_datahub_project", &resource.Sweeper{
		Name: "alicloud_datahub_project",
		F:    testSweepDatahubProject,
	})
}

func testSweepDatahubProject(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	dh := client.(*AliyunClient).dhconn

	// List projects
	projects, err := dh.ListProjects()
	if err != nil {
		return fmt.Errorf("error listing Datahub projects: %s", err)
	}

	for _, projectName := range projects.Names {
		// a testing project?
		if !isTerraformTestingDatahubObject(projectName) {
			log.Printf("[INFO] Skipping Datahub project: %s", projectName)
			continue
		}
		log.Printf("[INFO] Deleting project: %s", projectName)

		// List topics
		topics, err := dh.ListTopics(projectName)
		if err != nil {
			return fmt.Errorf("error listing Datahub topics: %s", err)
		}

		for _, topicName := range topics.Names {
			log.Printf("[INFO] Deleting topic: %s/%s", projectName, topicName)

			// List subscriptions
			subscriptions, err := dh.ListSubscriptions(projectName, topicName)

			if err != nil {
				return fmt.Errorf("error listing Datahub subscriptions: %s", err)
			}

			for _, subscription := range subscriptions.Subscriptions {
				log.Printf("[INFO] Deleting subscription: %s/%s/%s", projectName, topicName, subscription.SubId)

				// Delete subscription
				err := dh.DeleteSubscription(projectName, topicName, subscription.SubId)
				if err != nil {
					log.Printf("[ERROR] Failed to delete Datahub subscriptions: %s/%s/%s", projectName, topicName, subscription.SubId)
					return fmt.Errorf("error deleting  Datahub subscriptions: %s/%s/%s", projectName, topicName, subscription.SubId)
				}
			}

			// Delete topic
			err = dh.DeleteTopic(projectName, topicName)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Datahub topic: %s/%s", projectName, topicName)
				return fmt.Errorf("[ERROR] Failed to delete Datahub topic: %s/%s", projectName, topicName)
			}
		}

		// Delete project
		err = dh.DeleteProject(projectName)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Datahub project: %s", projectName)
			return fmt.Errorf("[ERROR] Failed to delete Datahub project: %s", projectName)
		}
	}

	return nil
}

func TestAccAlicloudDatahubProject_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_datahub_project.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDatahubProjectDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubProject,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_project.basic",
						"name", "tf_testacc_datahub_project"),
				),
			},
		},
	})
}

func TestAccAlicloudDatahubProject_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_datahub_project.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDatahubProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubProject,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_project.basic",
						"comment", "project for basic."),
				),
			},

			resource.TestStep{
				Config: testAccDatahubProjectUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_project.basic",
						"comment", "project for update."),
				),
			},
		},
	})
}

func testAccCheckDatahubProjectExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found Datahub project: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Datahub project ID is set")
		}

		dh := testAccProvider.Meta().(*AliyunClient).dhconn
		_, err := dh.GetProject(rs.Primary.ID)

		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckDatahubProjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_datahub_project" {
			continue
		}

		dh := testAccProvider.Meta().(*AliyunClient).dhconn
		_, err := dh.GetProject(rs.Primary.ID)

		if err != nil && isDatahubNotExistError(err) {
			continue
		}

		return fmt.Errorf("Datahub project %s may still exist", rs.Primary.ID)
	}

	return nil
}

const testAccDatahubProject = `
variable "name" {
  default = "tf_testacc_datahub_project"
}
resource "alicloud_datahub_project" "basic" {
  name = "${var.name}"
  comment = "project for basic."
}
`

const testAccDatahubProjectUpdate = `
variable "name" {
  default = "tf_testacc_datahub_project"
}
resource "alicloud_datahub_project" "basic" {
  name = "${var.name}"
  comment = "project for update."
}
`
