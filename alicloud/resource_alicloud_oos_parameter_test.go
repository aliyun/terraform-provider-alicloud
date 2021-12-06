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
		"alicloud_oos_parameter",
		&resource.Sweeper{
			Name: "alicloud_oos_parameter",
			F:    testSweepOosParameter,
		})
}

func testSweepOosParameter(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.OOSSupportRegions) {
		log.Printf("[INFO] Skipping Oos Parameter unsupported region: %s", region)
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
	action := "ListParameters"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
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

		resp, err := jsonpath.Get("$.Parameters", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Parameters", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Oos Parameter: %s", item["Name"].(string))
				continue
			}
			action := "DeleteParameter"
			request := map[string]interface{}{
				"Name": item["Name"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Oos Parameter (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Oos Parameter success: %s ", item["Name"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudOOSParameter_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_parameter.default"
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOOSParameterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosParameter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soosparameter%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOOSParameterBasicDependence0)
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
					"parameter_name": "${var.name}",
					"type":           "String",
					"value":          "tf-testacc-oos_parameter",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_name": name,
						"type":           "String",
						"value":          "tf-testacc-oos_parameter",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"value": "tf-testacc-oos_parameter_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": "tf-testacc-oos_parameter_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc",
						"For":     "OosParameter",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc",
						"tags.For":     "OosParameter",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
					"value":       "tf-testacc-oos_parameter",
					"tags": map[string]string{
						"Created": "tfTestAcc1",
						"For":     "OosParameter1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  name + "_update",
						"value":        "tf-testacc-oos_parameter",
						"tags.%":       "2",
						"tags.Created": "tfTestAcc1",
						"tags.For":     "OosParameter1",
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

func TestAccAlicloudOOSParameter_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_parameter.default"
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOOSParameterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosParameter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soosparameter%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOOSParameterBasicDependence0)
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
					"parameter_name": "${var.name}",
					"type":           "String",
					"value":          "tf-testacc-oos_parameter",
					"description":    "${var.name}",
					"tags": map[string]string{
						"Created": "tfTestAcc",
						"For":     "OosParameter",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"constraints":       `{\"AllowedValues\":[\"tf-testacc-oos_parameter\"], \"AllowedPattern\": \"tf-testacc-oos_parameter\", \"MinLength\": 1, \"MaxLength\": 100}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_name":    name,
						"type":              "String",
						"value":             "tf-testacc-oos_parameter",
						"description":       name,
						"tags.%":            "2",
						"tags.Created":      "tfTestAcc",
						"tags.For":          "OosParameter",
						"resource_group_id": CHECKSET,
						"constraints":       CHECKSET,
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

var AlicloudOOSParameterMap0 = map[string]string{
	"value":             CHECKSET,
	"resource_group_id": CHECKSET,
	"parameter_name":    CHECKSET,
}

func AlicloudOOSParameterBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}
