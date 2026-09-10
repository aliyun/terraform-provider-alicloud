package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudDatahubTopic_basic(t *testing.T) {
	var v *datahub.GetTopicResult

	resourceId := "alicloud_datahub_topic.default"
	ra := resourceAttrInit(resourceId, datahubTopicBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testacc_datahub_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubTopicConfigDependence)

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
					"name":         name,
					"project_name": "${alicloud_datahub_project.default.name}",
					"record_schema": map[string]string{
						"bigint_field":    "BIGINT",
						"timestamp_field": "TIMESTAMP",
						"string_field":    "STRING",
						"double_field":    "DOUBLE",
						"boolean_field":   "BOOLEAN",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            name,
						"project_name":    name,
						"record_schema.%": "5",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// TODO There is a GetTopic bug that it will return diff comment value when invkoing twice.
			// After it is fixed, reopen this case.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"comment": "topic added by terraform update",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"comment": "topic added by terraform update",
			//		}),
			//	),
			//},
		},
	})
}

func TestAccAlicloudDatahubTopic_blob(t *testing.T) {
	var v *datahub.GetTopicResult

	resourceId := "alicloud_datahub_topic.default"
	ra := resourceAttrInit(resourceId, datahubTopicBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testacc_datahub_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubTopicConfigDependence)

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
					"name":         name,
					"project_name": "${alicloud_datahub_project.default.name}",
					"record_type":  "BLOB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"project_name": name,
						"record_type":  "BLOB",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// TODO There is a GetTopic bug that it will return diff comment value when invkoing twice.
			// After it is fixed, reopen this case.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"comment": "topic added by terraform update",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"comment": "topic added by terraform update",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"comment": REMOVEKEY,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"comment": "topic added by terraform",
			//		}),
			//	),
			//},
		},
	})
}

func TestAccAlicloudDatahubTopic_multi(t *testing.T) {
	var v *datahub.GetTopicResult

	resourceId := "alicloud_datahub_topic.default.4"
	ra := resourceAttrInit(resourceId, datahubTopicBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testacc_datahub_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubTopicConfigDependence)

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
					"name":         name + "${count.index}",
					"project_name": "${alicloud_datahub_project.default.name}",
					"record_schema": map[string]string{
						"bigint_field":    "BIGINT",
						"timestamp_field": "TIMESTAMP",
						"string_field":    "STRING",
						"double_field":    "DOUBLE",
						"boolean_field":   "BOOLEAN",
					},
					"count": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceDatahubTopicConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alicloud_datahub_project" "default" {
	  name = "${var.name}"
	  comment = "project for basic."
	}
	`, name)
}

var datahubTopicBasicMap = map[string]string{
	"name":             CHECKSET,
	"project_name":     CHECKSET,
	"shard_count":      "1",
	"life_cycle":       "3",
	"comment":          "topic added by terraform",
	"record_type":      "TUPLE",
	"create_time":      CHECKSET,
	"last_modify_time": CHECKSET,
}

func testAccCheckDatahubTopicExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found Datahub topic: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Datahub topic ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		projectName := split[0]
		topicName := split[1]
		_, err := client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
			return dataHubClient.GetTopic(projectName, topicName)
		})

		if err != nil {
			return err
		}
		return nil
	}
}
