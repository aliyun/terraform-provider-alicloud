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
		"alicloud_eais_instance",
		&resource.Sweeper{
			Name: "alicloud_eais_instance",
			F:    testSweepEaisInstance,
		})
}

func testSweepEaisInstance(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeEais"
	request := map[string]interface{}{}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId

	var response map[string]interface{}
	conn, err := client.NewEaisClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-24"), StringPointer("AK"), nil, request, &runtime)
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
		v, err := jsonpath.Get("$.Instances.Instance", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Instances.Instance", action, err)
			return nil
		}
		result, _ := v.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["InstanceName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["InstanceName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Eais Instance: %s", item["InstanceName"].(string))
				continue
			}
			action := "DeleteEai"
			request := map[string]interface{}{
				"ElasticAcceleratedInstanceId": item["ElasticAcceleratedInstanceId"],
				"Force":                        true,
			}
			request["ClientToken"] = buildClientToken("DeleteEai")
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-24"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Eais Instance (%s): %s", item["ElasticAcceleratedInstanceId"].(string), err)
			}
			log.Printf("[INFO] Delete Eais Instance success: %s ", item["ElasticAcceleratedInstanceId"].(string))

		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudEaisInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eais_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEaisInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EaisService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEaisInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seaisinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEaisInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EAISSystemSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name,
					"instance_type":     "eais.ei-a6.medium",
					"security_group_id": "${alicloud_security_group.default.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"instance_type": "eais.ei-a6.medium",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vswitch_id", "force", "security_group_id"},
			},
		},
	})
}

var AlicloudEaisInstanceMap0 = map[string]string{}

func AlicloudEaisInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
		variable "name" {
		  default = "%s"
		}
		data "alicloud_zones" "default" {
			available_resource_creation = "VSwitch"
		}
		data "alicloud_vpcs" "default"{
			name_regex = "default-NODELETING"
		}
		data "alicloud_vswitches" "default" {
		  vpc_id  = data.alicloud_vpcs.default.ids.0
          zone_id = data.alicloud_zones.default.ids.0
		}
		
		resource "alicloud_security_group" "default" {
		  name        = var.name
		  description = "tf test"
		  vpc_id      = data.alicloud_vpcs.default.ids.0
		}
`, name)
}
