package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_oos_state_configuration",
		&resource.Sweeper{
			Name: "alicloud_oos_state_configuration",
			F:    testSweepOOSStateConfiguration,
		})
}

func testSweepOOSStateConfiguration(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.OOSSupportRegions) {
		log.Printf("[INFO] Skipping Oos Patch Baseline unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListStateConfigurations"
	request := map[string]interface{}{}
	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	conn, err := client.NewOosClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.StateConfigurations", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.StateConfigurations", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Description"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Oos Patch Baseline: %s", item["Description"].(string))
				continue
			}
			action := "DeleteStateConfigurations"
			request := map[string]interface{}{
				"StateConfigurationIds": convertListToJsonString([]interface{}{item["StateConfigurationId"]}),
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Oos Patch Baseline (%s): %s", item["Description"].(string), err)
			}
			log.Printf("[INFO] Delete Oos Patch Baseline success: %s ", item["Description"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudOOSStateConfiguration_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_state_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOOSStateConfigurationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosStateConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soosstateconfiguration%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOOSStateConfigurationBasicDependence0)
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
					"template_name":       "ACS-ECS-InventoryDataCollection",
					"configure_mode":      "ApplyOnly",
					"description":         name,
					"schedule_type":       "rate",
					"schedule_expression": "1 hour",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"targets":             `{\"Filters\": [{\"Type\": \"All\", \"Parameters\": {\"InstanceChargeType\": \"PrePaid\"}}], \"ResourceType\": \"ALIYUN::ECS::Instance\"}`,
					"parameters":          `{\"policy\": {\"ACS:Application\": {\"Collection\": \"Enabled\"}}}`,
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "test1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name":       "ACS-ECS-InventoryDataCollection",
						"configure_mode":      "ApplyOnly",
						"description":         name,
						"schedule_type":       "rate",
						"schedule_expression": "1 hour",
						"resource_group_id":   CHECKSET,
						"parameters":          CHECKSET,
						"targets":             CHECKSET,
						"tags.%":              "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_expression": "30 minutes",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_expression": "30 minutes",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": `{\"policy\":{\"ACS:InstanceDetailedInformation\":{\"Collection\":\"Enabled\"},\"ACS:WindowsUpdate\":{\"Collection\":\"Enabled\"}}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": `{\"Filters\": [{\"Type\": \"All\", \"Parameters\": {\"Status\":\"Running\"}}], \"ResourceType\": \"ALIYUN::ECS::Instance\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets": CHECKSET,
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
					"configure_mode": "ApplyAndMonitor",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configure_mode": "ApplyAndMonitor",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF2",
						"For":     "test2",
						"From":    "Accept",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configure_mode":      "ApplyOnly",
					"description":         name,
					"schedule_type":       "rate",
					"schedule_expression": "1 hour",
					"targets":             `{\"Filters\": [{\"Type\": \"All\", \"Parameters\": {\"InstanceChargeType\": \"PrePaid\"}}], \"ResourceType\": \"ALIYUN::ECS::Instance\"}`,
					"parameters":          `{\"policy\": {\"ACS:Application\": {\"Collection\": \"Enabled\"}}}`,
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "test1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configure_mode":      "ApplyOnly",
						"description":         name,
						"schedule_type":       "rate",
						"schedule_expression": "1 hour",
						"parameters":          CHECKSET,
						"targets":             CHECKSET,
						"tags.%":              "2",
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

func TestAccAlicloudOOSStateConfiguration_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_state_configuration.default"
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOOSStateConfigurationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosStateConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soosstateconfiguration%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOOSStateConfigurationBasicDependence0)
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
					"template_name":       "ACS-ECS-InventoryDataCollection",
					"description":         name,
					"schedule_type":       "rate",
					"schedule_expression": "1 hour",
					"targets":             `{\"Filters\": [{\"Type\": \"All\", \"Parameters\": {\"InstanceChargeType\": \"PrePaid\"}}], \"ResourceType\": \"ALIYUN::ECS::Instance\"}`,
					"parameters":          `{\"policy\": {\"ACS:Application\": {\"Collection\": \"Enabled\"}}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name":       "ACS-ECS-InventoryDataCollection",
						"description":         name,
						"schedule_type":       "rate",
						"schedule_expression": "1 hour",
						"resource_group_id":   CHECKSET,
						"parameters":          CHECKSET,
						"targets":             CHECKSET,
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

var AlicloudOOSStateConfigurationMap0 = map[string]string{}

func AlicloudOOSStateConfigurationBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}
