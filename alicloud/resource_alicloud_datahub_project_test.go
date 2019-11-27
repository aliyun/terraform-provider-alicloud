package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

var datahubProjectSuffixMin = 100000
var datahubProjectSuffixMax = 999999

func init() {
	resource.AddTestSweepers("alicloud_datahub_project", &resource.Sweeper{
		Name: "alicloud_datahub_project",
		F:    testSweepDatahubProject,
	})
}

func testSweepDatahubProject(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DatahubSupportedRegions) {
		log.Printf("[INFO] Skipping Datahub unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// List projects
	raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
		return dataHubClient.ListProjects()
	})
	if err != nil {
		// Now, only some region support Datahub
		log.Printf("[ERROR] Failed to list Datahub projects: %s", err)
	}
	projects, _ := raw.(*datahub.Projects)

	for _, projectName := range projects.Names {
		// a testing project?
		if !isTerraformTestingDatahubObject(projectName) {
			log.Printf("[INFO] Skipping Datahub project: %s", projectName)
			continue
		}
		log.Printf("[INFO] Deleting project: %s", projectName)

		// List topics
		raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			return dataHubClient.ListTopics(projectName)
		})
		if err != nil {
			return fmt.Errorf("error listing Datahub topics: %s", err)
		}
		topics, _ := raw.(*datahub.Topics)

		for _, topicName := range topics.Names {
			log.Printf("[INFO] Deleting topic: %s/%s", projectName, topicName)

			// List subscriptions
			raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
				return dataHubClient.ListSubscriptions(projectName, topicName)
			})

			if err != nil {
				return fmt.Errorf("error listing Datahub subscriptions: %s", err)
			}
			subscriptions, _ := raw.(*datahub.Subscriptions)

			for _, subscription := range subscriptions.Subscriptions {
				log.Printf("[INFO] Deleting subscription: %s/%s/%s", projectName, topicName, subscription.SubId)

				// Delete subscription
				_, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
					return nil, dataHubClient.DeleteSubscription(projectName, topicName, subscription.SubId)
				})
				if err != nil {
					log.Printf("[ERROR] Failed to delete Datahub subscriptions: %s/%s/%s", projectName, topicName, subscription.SubId)
					return fmt.Errorf("error deleting  Datahub subscriptions: %s/%s/%s", projectName, topicName, subscription.SubId)
				}
			}

			// Delete topic
			_, err = client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
				return nil, dataHubClient.DeleteTopic(projectName, topicName)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Datahub topic: %s/%s", projectName, topicName)
				return fmt.Errorf("[ERROR] Failed to delete Datahub topic: %s/%s", projectName, topicName)
			}
		}

		// Delete project
		_, err = client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			return nil, dataHubClient.DeleteProject(projectName)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Datahub project: %s", projectName)
			return fmt.Errorf("[ERROR] Failed to delete Datahub project: %s", projectName)
		}
	}

	return nil
}

func TestAccAlicloudDatahubProject_basic(t *testing.T) {
	var v *datahub.Project

	resourceId := "alicloud_datahub_project.default"
	ra := resourceAttrInit(resourceId, datahubProjectBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testaccdatahubproject%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubProjectConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
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
					"comment": "project for basic.",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "project for basic.",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "project added by terraform",
					}),
				),
			},
		},
	})
}
func TestAccAlicloudDatahubProject_multi(t *testing.T) {
	var v *datahub.Project

	resourceId := "alicloud_datahub_project.default.4"
	ra := resourceAttrInit(resourceId, datahubProjectBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testaccdatahubproject%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubProjectConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
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
func resourceDatahubProjectConfigDependence(name string) string {
	return ""
}

var datahubProjectBasicMap = map[string]string{
	"name":    CHECKSET,
	"comment": "project added by terraform",
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		_, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			return dataHubClient.GetProject(rs.Primary.ID)
		})

		if err != nil {
			return err
		}
		return nil
	}
}
