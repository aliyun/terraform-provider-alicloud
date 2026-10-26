package bpstudio_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/bpstudio"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_bp_studio_application", &resource.Sweeper{
		Name: "alicloud_bp_studio_application",
		F:    testSweepBpStudioApplication,
	})
}

func testSweepBpStudioApplication(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	bpStudioService := bpstudio.NewService(client)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListApplication"
	request := make(map[string]interface{})
	request["MaxResults"] = helper.PageSizeLarge
	request["NextToken"] = 1
	var response map[string]interface{}
	BpStudioApplicationIds := make([]string, 0)
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("BPStudio", "2021-09-31", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_bp_studio_application", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			skip := true
			item := v.(map[string]interface{})
			if !sweep.SweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["Name"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping BpStudioApplication Instance: %v", item["ApplicationId"])
					continue
				}
			}
			BpStudioApplicationIds = append(BpStudioApplicationIds, fmt.Sprint(item["ApplicationId"]))
		}
		if len(result) < request["MaxResults"].(int) {
			break
		}
		request["NextToken"] = request["NextToken"].(int) + 1
	}

	for _, id := range BpStudioApplicationIds {
		log.Printf("[INFO] Deleting BpStudioApplication Instance: %s", id)
		releaseAction := "ReleaseApplication"
		deleteAction := "DeleteApplication"
		if err != nil {
			return errs.WrapError(err)
		}
		request = map[string]interface{}{
			"ApplicationId": id,
		}
		object, err := bpStudioService.DescribeBpStudioApplication(id)
		if err != nil {
			return errs.WrapError(err)
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)

		if fmt.Sprint(object["Status"]) == "Deployed_Failure" || fmt.Sprint(object["Status"]) == "PartiallyDeployedSuccess" || fmt.Sprint(object["Status"]) == "Deployed_Success" || fmt.Sprint(object["Status"]) == "Destroyed_Failure" || fmt.Sprint(object["Status"]) == "PartiallyDestroyedSuccess" {
			err = retry.Retry(120*time.Minute, func() *retry.RetryError {
				_, err = client.RpcPost("BPStudio", "2021-09-31", releaseAction, nil, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			releaseApplicationStateConf := helper.BuildStateConf([]string{}, []string{"Destroyed_Success"}, 120*time.Minute, 5*time.Second, bpStudioService.BpStudioApplicationStateRefreshFunc(id, []string{}))
			if _, err := releaseApplicationStateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, id)
			}
		}

		err = retry.Retry(120*time.Minute, func() *retry.RetryError {
			_, err = client.RpcPost("BPStudio", "2021-09-31", deleteAction, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete BpStudioApplication Instance (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccAliCloudBpStudioApplication_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bp_studio_application.default"
	acctest.PreCheckWithRegions(t, true, connectivity.BpStudioApplicationSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, resourceAlicloudBpStudioApplicationMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return bpstudio.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeBpStudioApplication")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBpStudioApplication-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlicloudBpStudioApplicationBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name":  name,
					"template_id":       "YAUUQIYRSV1CMFGX",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"area_id":           acctest.DefaultRegionToTest,
					"instances": []map[string]string{
						{
							"id":        "${alicloud_instance.default.id}",
							"node_name": "${alicloud_instance.default.instance_name}",
							"node_type": "ecs",
						},
					},
					"configuration": map[string]string{
						"enableMonitor": "1",
					},
					"variables": map[string]string{
						"test": "1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name":  name,
						"template_id":       "YAUUQIYRSV1CMFGX",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"area_id", "instances", "configuration", "variables"},
			},
		},
	})
}

var resourceAlicloudBpStudioApplicationMap = map[string]string{
	"status": acctest.CHECKSET,
}

func resourceAlicloudBpStudioApplicationBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_zones" "default" {
		available_disk_category     = "cloud_efficiency"
		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
		name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
		most_recent   = true
		owners        = "system"
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
	}

	data "alicloud_instance_types" "default" {
		availability_zone    = data.alicloud_zones.default.zones.0.id
		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
		name   = var.name
		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_instance" "default" {
		image_id             = data.alicloud_images.default.images.0.id
		instance_type        = data.alicloud_instance_types.default.instance_types.0.id
		instance_name        = var.name
		security_groups      = alicloud_security_group.default.*.id
		availability_zone    = data.alicloud_zones.default.zones.0.id
		instance_charge_type = "PostPaid"
		system_disk_category = "cloud_efficiency"
		vswitch_id           = data.alicloud_vswitches.default.ids.0
	}
`, name)
}

// resourceTestAccConfigFunc is a thin delegate to acctest.ResourceTestAccConfig.
// It exists only so the collection-order checker (scripts/collection-order) can
// establish the builder for the reorder test below: that checker matches the
// bare constructor by name (ident(call.Fun) == "resourceTestAccConfigFunc") and
// does not yet recognize the qualified acctest.ResourceTestAccConfig spelling
// (PR #11237 fixes that). The flat-layout original lives in the root alicloud
// package and cannot be imported from a migrated service package, so ecs and
// agentloop keep an identical local spelling; this mirrors that pattern without
// copying the whole helper type. Once #11237 lands, this wrapper and the bare
// call below can switch back to the qualified form.
func resourceTestAccConfigFunc(resourceId, name string, configDependence func(string) string) acctest.ResourceTestAccConfigFunc {
	return acctest.ResourceTestAccConfig(resourceId, name, configDependence)
}

// resourceAlicloudBpStudioApplicationReorderDependence provisions two ECS
// instances so the ordered instances TypeList can carry two distinct members.
// It reuses the basic0 fixture (default resource manager group, no name_regex
// filter) already exercised by TestAccAliCloudBpStudioApplication_basic0 on
// the BPStudio support regions.
func resourceAlicloudBpStudioApplicationReorderDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_zones" "default" {
		available_disk_category     = "cloud_efficiency"
		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
		name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
		most_recent   = true
		owners        = "system"
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
	}

	data "alicloud_instance_types" "default" {
		availability_zone    = data.alicloud_zones.default.zones.0.id
		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
		name   = var.name
		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_instance" "first" {
		image_id             = data.alicloud_images.default.images.0.id
		instance_type        = data.alicloud_instance_types.default.instance_types.0.id
		instance_name        = "%[1]s-first"
		security_groups      = alicloud_security_group.default.*.id
		availability_zone    = data.alicloud_zones.default.zones.0.id
		instance_charge_type = "PostPaid"
		system_disk_category = "cloud_efficiency"
		vswitch_id           = data.alicloud_vswitches.default.ids.0
	}

	resource "alicloud_instance" "second" {
		image_id             = data.alicloud_images.default.images.0.id
		instance_type        = data.alicloud_instance_types.default.instance_types.0.id
		instance_name        = "%[1]s-second"
		security_groups      = alicloud_security_group.default.*.id
		availability_zone    = data.alicloud_zones.default.zones.0.id
		instance_charge_type = "PostPaid"
		system_disk_category = "cloud_efficiency"
		vswitch_id           = data.alicloud_vswitches.default.ids.0
	}
`, name)
}

// TestAccAliCloudBpStudioApplication_instancesReorder covers the ordered
// instances TypeList: the original member order is applied, then a convergent
// apply reorders the two distinct members. instances is ForceNew with no
// Update and Read does not write it back, so the reordered apply replaces the
// resource and converges on an empty post-apply plan; ForceNew replacement is
// an allowed reorder mode per scripts/collection-order/README.md.
func TestAccAliCloudBpStudioApplication_instancesReorder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bp_studio_application.default"
	acctest.PreCheckWithRegions(t, true, connectivity.BpStudioApplicationSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, resourceAlicloudBpStudioApplicationMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return bpstudio.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeBpStudioApplication")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBpStudioApplication-reorder%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudBpStudioApplicationReorderDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name":  name,
					"template_id":       "YAUUQIYRSV1CMFGX",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"area_id":           acctest.DefaultRegionToTest,
					"instances": []map[string]string{
						{
							"id":        "${alicloud_instance.first.id}",
							"node_name": "${alicloud_instance.first.instance_name}",
							"node_type": "ecs",
						},
						{
							"id":        "${alicloud_instance.second.id}",
							"node_name": "${alicloud_instance.second.instance_name}",
							"node_type": "ecs",
						},
					},
					"configuration": map[string]string{
						"enableMonitor": "1",
					},
					"variables": map[string]string{
						"test": "1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name":  name,
						"template_id":       "YAUUQIYRSV1CMFGX",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name":  name,
					"template_id":       "YAUUQIYRSV1CMFGX",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"area_id":           acctest.DefaultRegionToTest,
					"instances": []map[string]string{
						{
							"id":        "${alicloud_instance.second.id}",
							"node_name": "${alicloud_instance.second.instance_name}",
							"node_type": "ecs",
						},
						{
							"id":        "${alicloud_instance.first.id}",
							"node_name": "${alicloud_instance.first.instance_name}",
							"node_type": "ecs",
						},
					},
					"configuration": map[string]string{
						"enableMonitor": "1",
					},
					"variables": map[string]string{
						"test": "1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name":  name,
						"template_id":       "YAUUQIYRSV1CMFGX",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}
