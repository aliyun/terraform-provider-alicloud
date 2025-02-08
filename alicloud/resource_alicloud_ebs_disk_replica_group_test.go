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

func TestAccAlicloudEBSDiskReplicaGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_disk_replica_group.default"
	checkoutSupportedRegions(t, true, connectivity.EBSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudEbsDiskReplicaGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsDiskReplicaGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sebsdiskreplicagroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsDiskReplicaGroupBasicDependence0)
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
					"source_region_id":      "${var.region}",
					"source_zone_id":        "${data.alicloud_ebs_regions.default.regions[0].zones[0].zone_id}",
					"destination_region_id": "${var.region}",
					"destination_zone_id":   "${data.alicloud_ebs_regions.default.regions[0].zones[1].zone_id}",
					"group_name":            name,
					"description":           name,
					"rpo":                   "900",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_zone_id":        CHECKSET,
						"destination_region_id": CHECKSET,
						"destination_zone_id":   CHECKSET,
						"group_name":            name,
						"description":           name,
						"rpo":                   "900",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name + "update",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"group_name":  name,
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":  name,
						"description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
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

variable "region" {
  default = "%s"
}

data "alicloud_ebs_regions" "default"{
  region_id = var.region
}

`, name, defaultRegionToTest)
}
