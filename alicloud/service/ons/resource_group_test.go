package ons_test

import (
	"fmt"
	"log"
	"regexp"
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
	resource.AddTestSweepers("alicloud_ons_group", &resource.Sweeper{
		Name: "alicloud_ons_group",
		F:    testSweepOnsGroup,
	})
}

func testSweepOnsGroup(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
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
		log.Printf("[ERROR] Failed to retrieve ons instance in service list: %s", err)
	}
	resp, err := jsonpath.Get("$.Data.InstanceVO", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.InstanceVO", response)
	}

	var instanceIds []string
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		instanceIds = append(instanceIds, item["InstanceId"].(string))
	}

	for _, instanceId := range instanceIds {

		action := "OnsGroupList"
		request := make(map[string]interface{})
		var response map[string]interface{}
		response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ons_groups", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.SubscribeInfoDo", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.SubscribeInfoDo", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := item["GroupId"].(string)
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ons group: %s ", name)
				continue
			}
			log.Printf("[INFO] delete ons group: %s ", name)

			action := "OnsGroupDelete"
			request := map[string]interface{}{
				"GroupId":    name,
				"InstanceId": instanceId,
			}
			response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete ons group (%s): %s", name, err)

			}
		}
	}

	return nil
}

func TestAccAliCloudOnsGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ons_group.default"
	ra := resourceAttrInit(resourceId, onsGroupBasicMap)
	serviceFunc := func() interface{} {
		return ons.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOnsGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_ons_instance.default.id}",
					"group_id":    "${var.group_id}",
					"remark":      "alicloud_ons_group_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand),
						"remark":   "alicloud_ons_group_remark",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"read_enable"},
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TFM",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "1",
						"tags.Created": "TFM",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"group_id": "${var.group_id}_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": fmt.Sprintf("GID-tf-testacconsgroupbasic%v_change", rand)}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"group_id": "${var.group_id}",
					"remark":   "alicloud_ons_group_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand),
						"remark":   "alicloud_ons_group_remark",
					}),
				),
			},

			// read_enable can only be set after the group is used by a client
			// (see docs/r/ons_group.html.markdown). OnsGroupConsumerUpdate returns
			// -2075 ("The Group ID has never been online") for a freshly created
			// group and is not in helper.NeedRetry, so the error is non-retryable;
			// ExpectError asserts that failure path so the attribute is exercised
			// in coverage rather than exempted or left uncovered. Placed after the
			// group_id/remark update steps: the apply fails (ExpectError) and the
			// state is not updated, so the subsequent ForceNew recreate is unaffected.
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_ons_instance.default.id}",
					"group_id":    "${var.group_id}",
					"remark":      "alicloud_ons_group_remark",
					"read_enable": true,
				}),
				ExpectError: regexp.MustCompile(`(?i).*OnsGroupConsumerUpdate.*`),
			},

			// Switch from the deprecated group_id to group_name (ConflictWith
			// group_id) and exercise group_type=http; both are ForceNew so the
			// resource is recreated, covering group_name and group_type.
			// configUpdate merges each step's map over the persistent
			// attributeMap, so group_id (from the create step) and read_enable
			// (from the ExpectError step above) are still present here; remove
			// them via the REMOVEKEY sentinel so group_name does not conflict
			// with group_id and the recreate does not re-trigger read_enable.
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_ons_instance.default.id}",
					"group_id":    REMOVEKEY,
					"read_enable": REMOVEKEY,
					"group_name":  "${var.group_id}",
					"group_type":  "http",
					"remark":      "alicloud_ons_group_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand),
						"group_type": "http",
						"remark":     "alicloud_ons_group_remark",
					}),
				),
			},
		},
	})

}

func TestAccAliCloudOnsGroup_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ons_group.default"
	ra := resourceAttrInit(resourceId, onsGroupBasicMap)
	serviceFunc := func() interface{} {
		return ons.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOnsGroupMultiConfigMulti)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_ons_instance.default.id}",
					"group_id":    "${var.group_id}",
					"remark":      "alicloud_ons_group_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand),
						"remark":   "alicloud_ons_group_remark",
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

func resourceOnsGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ons_instance" "default" {
  name = "%s"
}

variable "group_id" {
 default = "%s"
}
`, name, name)
}

func resourceOnsGroupMultiConfigMulti(name string) string {
	return fmt.Sprintf(`
variable "group_id" {
  default = "%s"
}
variable "groups" {
  type = list(map(string))
  default = [
    {
      "group_name" : "GID_TF_TEST_1_2",
      "remark" : "TF_TEST_1_2"
    },
    {
      "group_name" : "GID_TF_TEST_1_1",
      "remark" : "TF_TEST_1_1"
    },
      {
      "group_name" : "GID_TF_TEST_1",
      "remark" : "TF_TEST_1"
    },
      {
      "group_name" : "GID_TF_TEST_2_1",
      "remark" : "TF_TEST_2_1"
    },
      {
      "group_name" : "GID_TF_TEST_2",
      "remark" : "TF_TEST_2"
    }
  ]
  description = "The specification of ons groups."
}

resource "alicloud_ons_instance" "default" {
  instance_name = "%s"
}

resource "alicloud_ons_group" "default1" {
  count       = length(var.groups)
  instance_id = alicloud_ons_instance.default.id
  group_name  = var.groups[count.index].group_name
  remark      = var.groups[count.index].remark
}
`, name, name)
}

var onsGroupBasicMap = map[string]string{
	"group_id": "${var.group_id}",
	"remark":   "alicloud_ons_group_remark",
}
