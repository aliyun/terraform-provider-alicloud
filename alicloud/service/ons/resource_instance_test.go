package ons_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	ons "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ons_instance", &resource.Sweeper{
		Name: "alicloud_ons_instance",
		F:    testSweepOnsInstance,
		Dependencies: []string{
			"alicloud_ons_topic",
			"alicloud_ons_group",
		},
	})
}

func testSweepOnsInstance(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
		"GID-tf-testAcc",
		"GID_tf-testacc",
		"CID-tf-testAcc",
		"CID_tf-testacc",
	}

	action := "OnsInstanceInServiceList"
	request := make(map[string]interface{})
	var response map[string]interface{}
	response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve ONS instance in service list: %s", err)
	}

	resp, err := jsonpath.Get("$.Data.InstanceVO", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.InstanceVO", response)
	}

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := item["InstanceName"].(string)
		skip := true
		if !sweep.SweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ons instance: %s ", name)
				continue
			}
		}
		log.Printf("[INFO] delete ons instance: %s ", name)

		action := "OnsInstanceDelete"
		request := map[string]interface{}{
			"InstanceId": item["InstanceId"],
		}
		_, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete ons instance (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAliCloudOnsInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ons_instance.default"
	ra := resourceAttrInit(resourceId, OnsInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ons.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOnsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sOnsInstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, OnsInstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
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
					"tags": map[string]string{
						"Created": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "1",
						"tags.Created": "TF",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "Test-Remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "Test-Remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}",
					"remark":        "Test-Remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"remark":        "Test-Remark",
					}),
				),
			},
		},
	})
}

var OnsInstanceMap = map[string]string{
	"instance_type": CHECKSET,
	"status":        CHECKSET,
}

func OnsInstanceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
