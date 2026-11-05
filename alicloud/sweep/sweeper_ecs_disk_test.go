package sweep_test

// Legacy product sweepers run in this dedicated test binary.

import (
	"fmt"
	"log"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	tfsweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_disk", &resource.Sweeper{
		Name: "alicloud_ecs_disk",
		F:    testAliCloudEcsDisk,
	})
}

func testAliCloudEcsDisk(region string) error {
	rawClient, err := tfsweep.SharedClientForRegion(region)
	if err != nil {
		return tferrs.WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"PageSize":   tfhelper.PageSizeLarge,
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
		tflogs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.Disks.Disk", response)
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, action, "$.Disks.Disk", response)
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
		if len(result) < tfhelper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}
