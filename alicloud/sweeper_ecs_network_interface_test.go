package alicloud

// The alicloud_ecs_network_interface sweeper stays in the root test binary:
// the make sweep target runs only ./alicloud, and the SDKv2 sweeper helpers
// (sharedClientForRegion, WrapErrorf, PageSizeLarge, VpcService, sweepAll,
// BuildStateConf, DefaultTimeout) live in the root alicloud test package.
// The resource implementation itself lives in alicloud/service/ecs; the
// exported EcsNetworkInterfaceStateRefreshFunc method on *ecs.Service is
// reached via ecssvc.NewService(client) so the unexported client field is
// not constructed through a struct literal.

import (
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	ecssvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_network_interface", &resource.Sweeper{
		Name: "alicloud_ecs_network_interface",
		F:    testAliCloudEcsNetworkInterface,
	})
}

func testAliCloudEcsNetworkInterface(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf-example",
		"tf_example",
		"terraform-example",
		// Clean up the eni which created by cloud service
		"fc-eni",
		"bastion",
		"rds-eni",
	}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}
	action := "DescribeNetworkInterfaces"

	var response map[string]interface{}
	for {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Describe NetworkInterface failed, error: %s. Return!", err.Error())
			return nil
		}

		resp, err := jsonpath.Get("$.NetworkInterfaceSets.NetworkInterfaceSet", response)
		if err != nil {
			log.Printf("[ERROR] jsonpath Get NetworkInterface failed, %#v", err)
			continue
		}

		result, _ := resp.([]interface{})
		service := VpcService{client}
		ecsService := ecssvc.NewService(client)
		for _, v := range result {
			item := v.(map[string]interface{})
			if _, ok := item["NetworkInterfaceName"]; !ok {
				continue
			}
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["NetworkInterfaceName"].(string)), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				// If a nat gateway name is not set successfully, it should be fetched by vpc name and deleted.
				if skip {
					if need, err := service.needSweepVpc(item["VpcId"].(string), ""); err == nil {
						skip = !need
					}
				}
				if skip {
					log.Printf("[INFO] Skipping NetworkInterface: %s", item["NetworkInterfaceName"].(string))
					continue
				}
			}
			if item["InstanceId"] != "" {
				requestDetach := map[string]interface{}{
					"InstanceId":         item["InstanceId"],
					"NetworkInterfaceId": item["NetworkInterfaceId"],
					"RegionId":           client.RegionId,
				}
				actionDetach := "DetachNetworkInterface"

				response, err = client.RpcPost("Ecs", "2014-05-26", actionDetach, nil, requestDetach, true)
				if err != nil {
					log.Printf("[ERROR] Detach NetworkInterface failed, %#v", err)
					continue
				}
				stateConf := BuildStateConf([]string{}, []string{"Available"}, DefaultTimeout, 5*time.Second, ecsService.EcsNetworkInterfaceStateRefreshFunc(item["NetworkInterfaceId"].(string), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					log.Printf("[ERROR] Detach NetworkInterface failed, %#v", err)
					continue
				}
			}
			action = "DeleteNetworkInterface"
			request := map[string]interface{}{
				"NetworkInterfaceId": item["NetworkInterfaceId"],
				"RegionId":           client.RegionId,
			}
			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete NetworkInterface (%s): %s", item["NetworkInterfaceName"].(string), err)
				continue
			}
			log.Printf("[INFO] Delete NetworkInterface success: %s ", item["NetworkInterfaceName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil
}
