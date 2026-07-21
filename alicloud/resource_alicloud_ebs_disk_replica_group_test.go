package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ebs_disk_replica_group",
		&resource.Sweeper{
			Name: "alicloud_ebs_disk_replica_group",
			F:    testSweepEbsDiskReplicaGroup,
		})
}

func testSweepEbsDiskReplicaGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDiskReplicaGroups"
	request := map[string]interface{}{}
	request["MaxResults"] = PageSizeXLarge
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.ReplicaGroups", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ReplicaGroups", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["GroupName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["GroupName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ebs DiskReplicaGroup: %s", item["GroupName"].(string))
				continue
			}
			action := "DeleteDiskReplicaGroup"
			request := map[string]interface{}{
				"ReplicaGroupId": item["ReplicaGroupId"],
				"RegionId":       client.RegionId,
			}
			request["ClientToken"] = buildClientToken("DeleteDiskReplicaGroup")
			_, err = client.RpcPost("ebs", "2021-07-30", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ebs DiskReplicaGroup (%s): %s", item["ReplicaGroupId"].(string), err)
			}
			log.Printf("[INFO] Delete Ebs DiskReplicaGroup success: %s ", item["ReplicaGroupId"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudEbsDiskReplicaGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_disk_replica_group.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsDiskReplicaGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsDiskReplicaGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccebs%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsDiskReplicaGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rpo":                     "900",
					"source_region_id":        "${var.disk-region}",
					"description":             "cctest",
					"destination_region_id":   "${var.disk-region}",
					"destination_zone_id":     "${var.dst-disk-zone}",
					"source_zone_id":          "${var.src-disk-zone}",
					"disk_replica_group_name": name,
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rpo":                     "900",
						"source_region_id":        CHECKSET,
						"description":             "cctest",
						"destination_region_id":   CHECKSET,
						"destination_zone_id":     CHECKSET,
						"source_zone_id":          CHECKSET,
						"disk_replica_group_name": name,
						"resource_group_id":       CHECKSET,
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
					"description": "description-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_replica_group_name": name + "_update",
					"pair_ids": []string{
						"${alicloud_ebs_disk_replica_pair.defaultUCZMS9.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_replica_group_name": name + "_update",
						"pair_ids.#":              "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":   "normal",
					"one_shot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":   "normal",
						"one_shot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "failovered",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "failovered",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":            "stopped",
					"reverse_replicate": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":            "stopped",
						"reverse_replicate": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pair_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pair_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"one_shot", "reverse_replicate"},
			},
		},
	})
}

var AlicloudEbsDiskReplicaGroupMap0 = map[string]string{}

func AlicloudEbsDiskReplicaGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "disk-region" {
  default = "cn-hangzhou"
}

variable "dst-disk-zone" {
  default = "cn-hangzhou-h"
}

variable "src-disk-zone" {
  default = "cn-hangzhou-i"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ecs_disk" "defaultOxDXka" {
  zone_id   = var.src-disk-zone
  size      = "20"
  disk_name = "fg-tf-e-case3"
  category  = "cloud_essd"
  lifecycle {
	ignore_changes = [tags]
  }
}

resource "alicloud_ecs_disk" "defaultQkzwDg" {
  zone_id   = var.dst-disk-zone
  size      = "20"
  disk_name = "fg-tf-b-case3"
  category  = "cloud_essd"
  lifecycle {
	ignore_changes = [tags]
  }
}

resource "alicloud_ebs_disk_replica_pair" "defaultUCZMS9" {
  destination_disk_id   = alicloud_ecs_disk.defaultQkzwDg.id
  destination_region_id = var.disk-region
  destination_zone_id   = var.dst-disk-zone
  payment_type          = "PayAsYouGo"
  source_zone_id        = var.src-disk-zone
  disk_id               = alicloud_ecs_disk.defaultOxDXka.id
}


`, name)
}
