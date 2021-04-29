package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_fnf_flow", &resource.Sweeper{
		Name:         "alicloud_fnf_flow",
		F:            testSweepFnfFlow,
		Dependencies: []string{"alicloud_fnf_schedule"},
	})
}

func testSweepFnfFlow(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	support := false
	for _, v := range connectivity.FnfSupportRegions {
		if v == connectivity.Region(region) {
			support = true
			break
		}
	}
	if !support {
		return nil
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListFlows"
	request := make(map[string]interface{})
	var response map[string]interface{}
	conn, err := client.NewFnfClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-03-15"), StringPointer("AK"), request, nil, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_fnf_flows", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	resp, err := jsonpath.Get("$.Flows", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Flows", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := item["Name"].(string)
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Fnf Flow: %s ", name)
			continue
		}
		log.Printf("[Info] Delete Fnf Flow: %s", name)

		action := "DeleteFlow"
		conn, err := client.NewFnfClient()
		if err != nil {
			return WrapError(err)
		}
		request := map[string]interface{}{
			"Name": name,
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-03-15"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Fnf Flow (%s): %s", name, err)
		}
	}
	return nil
}

func TestAccAlicloudFnfFlow_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fnf_flow.default"
	ra := resourceAttrInit(resourceId, AlicloudFnfFlowMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &FnfService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFnfFlow")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudFnfFlow%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFnfFlowBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.FnfSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"definition":  `version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld`,
					"description": "tf-testaccFnFFlow983041",
					"name":        "${var.name}",
					"type":        "FDL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"definition":  strings.Replace(`version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld`, `\n`, "\n", -1),
						"description": "tf-testaccFnFFlow983041",
						"name":        name,
						"type":        "FDL",
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
					"definition": `version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworldchange`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"definition": strings.Replace(`version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworldchange`, `\n`, "\n", -1),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testaccFnFFlow813242",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccFnFFlow813242",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_arn": `${alicloud_ram_role.default.arn}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "FDL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "FDL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"definition":  `version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld`,
					"description": "tf-testaccFnFFlow983041",
					"role_arn":    `${alicloud_ram_role.default.arn}`,
					"type":        "FDL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"definition":  strings.Replace(`version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld`, `\n`, "\n", -1),
						"description": "tf-testaccFnFFlow983041",
						"role_arn":    CHECKSET,
						"type":        "FDL",
					}),
				),
			},
		},
	})
}

var AlicloudFnfFlowMap0 = map[string]string{
	"flow_id":            CHECKSET,
	"last_modified_time": CHECKSET,
}

func AlicloudFnfFlowBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_ram_role" "default" {
  name = var.name
  document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "fnf.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
}
`, name)
}
