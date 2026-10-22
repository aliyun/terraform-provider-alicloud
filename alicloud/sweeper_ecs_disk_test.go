package alicloud

// The alicloud_ecs_disk sweeper stays in the root test binary: the make
// sweep target runs only ./alicloud, and the SDKv2 and plugin-testing
// helper/resource packages cannot both register -sweep in one binary. The
// resource implementation itself lives in alicloud/service/ecs.

import (
	"fmt"
	"log"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_disk", &resource.Sweeper{
		Name: "alicloud_ecs_disk",
		F:    testAliCloudEcsDisk,
	})
}

func testAliCloudEcsDisk(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}

	var response map[string]interface{}
	for {
		action := "DescribeDisks"
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %s", action, err)
			return nil
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Disks.Disk", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Disks.Disk", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["DiskName"])
			if !ecs.IsDiskSweepCandidateName(name) {
				log.Printf("[INFO] Skipping Disk: %s", name)
				continue
			}
			id := fmt.Sprint(item["DiskId"])
			log.Printf("[INFO] Deleting Disk: %s (%s)", name, id)
			// SweepDisk runs the forced release sequence, so a PrePaid disk left
			// behind by the force_delete acceptance cases is converted and
			// detached before DeleteDisk; the dependencies of those cases
			// (instance, vpc, vswitch, security group) carry the same tf-testAcc
			// name and are covered by their own sweepers.
			if err := ecs.NewService(client).SweepDisk(id); err != nil {
				log.Printf("[ERROR] Failed to delete Disk (%s): %s", name, err)
				continue
			}
			log.Printf("[INFO] Delete Disk success: %s ", name)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}
