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
	resource.AddTestSweepers(
		"alicloud_brain_industrial_pid_project",
		&resource.Sweeper{
			Name: "alicloud_brain_industrial_pid_project",
			F:    testSweepBrainIndustrialPidProject,
		})
}

func testSweepBrainIndustrialPidProject(region string) error {
	if !testSweepPreCheckWithRegions(region, false, connectivity.BrainIndustrialRegions) {
		log.Printf("[INFO] Skipping Brain Industrial unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"CurrentPage": 1,
		"PageSize":    20,
	}
	var response map[string]interface{}
	action := "ListPidProjects"
	conn, err := client.NewAistudioClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, _ = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &runtime)
		if fmt.Sprintf(`%v`, response["Code"]) != "200" {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_brain_industrial_pid_project", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.PidProjectList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PidProjectList", response)
		}
		for _, v := range resp.([]interface{}) {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["PidProjectName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Brain Industrial: %s", item["PidProjectName"].(string))
				continue
			}
			actionDelete := "DeletePidProject"
			requestDelete := map[string]interface{}{
				"PidProjectId": item["PidProjectId"],
			}
			response, err = conn.DoRequest(StringPointer(actionDelete), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, requestDelete, &util.RuntimeOptions{})
			if fmt.Sprintf(`%v`, response["Code"]) == "200" {
				log.Printf("[INFO] Delete Brain Industrial Project success: %s ", item["PidProjectName"].(string))
			} else if fmt.Sprintf(`%v`, response["Code"]) == "-100" && strings.Contains(response["Message"].(string), "存在回路") {
				log.Printf("[INFO] Firstly, Delete Loop belongs to Project")
				actionLoopList := "ListPidLoops"
				requestLoopList := map[string]interface{}{
					"PidProjectId": item["PidProjectId"],
					"PageSize":     20,
					"CurrentPage":  1,
				}
				for {
					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)

					responseLoop, _ := conn.DoRequest(StringPointer(actionLoopList), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, requestLoopList, &util.RuntimeOptions{})
					respLoop, _ := jsonpath.Get("$.PidLoopList", responseLoop)

					for _, v := range respLoop.([]interface{}) {
						itemLoop := v.(map[string]interface{})
						actionLoopDelete := "DeletePidLoop"
						requestLoopDelete := map[string]interface{}{
							"PidLoopId": itemLoop["PidLoopId"],
						}
						responseLoopDelete, _ := conn.DoRequest(StringPointer(actionLoopDelete), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, requestLoopDelete, &util.RuntimeOptions{})
						if fmt.Sprintf(`%v`, responseLoopDelete["Code"]) != "200" {
							log.Printf("[ERROR] Failed to delete Brain Industrial Loop (%s): %s", itemLoop["PidLoopId"].(string), responseLoopDelete["Message"].(string))
						} else {
							log.Printf("[INFO] Delete Brain Industrial Loop success (%s): %s", itemLoop["PidLoopId"].(string), responseLoopDelete["Message"].(string))
						}
					}
					if len(respLoop.([]interface{})) < request["PageSize"].(int) {
						break
					}
					request["CurrentPage"] = request["CurrentPage"].(int) + 1
				}
				log.Printf("[INFO] Delete Loop Done, Then delete Brain Industrial Project again")
				responseAgain, _ := conn.DoRequest(StringPointer(actionDelete), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, requestDelete, &util.RuntimeOptions{})
				if fmt.Sprintf(`%v`, responseAgain["Code"]) != "200" {
					log.Printf("[ERROR] Failed to again delete Brain Industrial Project  (%s): %s", item["PidProjectName"].(string), responseAgain["Message"].(string))
				} else {
					log.Printf("[INFO] Delete Brain Industrial Project again Success(%s): %s", item["PidProjectName"].(string), responseAgain["Message"].(string))
				}
			} else if fmt.Sprintf(`%v`, response["Code"]) != "200" {
				log.Printf("[ERROR] Failed to delete Brain Industrial Project (%s): %s", item["PidProjectName"].(string), response["Message"].(string))
			}
		}
		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}
	return nil
}

func TestAccAlicloudBrainIndustrialPidProject_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_brain_industrial_pid_project.default"
	ra := resourceAttrInit(resourceId, AlicloudBrainIndustrialPidProjectMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Brain_industrialService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBrainIndustrialPidProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBrainIndustrialPidProjectBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_organisation_id": "${alicloud_brain_industrial_pid_organization.default.id}",
					"pid_project_desc":    "tf test",
					"pid_project_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organisation_id": CHECKSET,
						"pid_project_desc":    "tf test",
						"pid_project_name":    name,
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
					"pid_organisation_id": "${alicloud_brain_industrial_pid_organization.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organisation_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_project_desc": "tf test update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_project_desc": "tf test update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_project_name": "tf-testAccUp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_project_name": "tf-testAccUp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_organisation_id": "${alicloud_brain_industrial_pid_organization.default.id}",
					"pid_project_desc":    "tf test",
					"pid_project_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organisation_id": CHECKSET,
						"pid_project_desc":    "tf test",
						"pid_project_name":    name,
					}),
				),
			},
		},
	})
}

var AlicloudBrainIndustrialPidProjectMap = map[string]string{}

func AlicloudBrainIndustrialPidProjectBasicDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_brain_industrial_pid_organization" "default" {
		pid_organization_name = "%s"
	}
	resource "alicloud_brain_industrial_pid_organization" "update" {
		pid_organization_name = "tf-testAccUp"
	}`, name)
}
