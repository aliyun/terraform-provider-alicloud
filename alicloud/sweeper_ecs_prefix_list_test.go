package alicloud

// The alicloud_ecs_prefix_list sweeper stays in the root test binary: the
// make sweep target runs only ./alicloud, and the SDKv2 sweeper helpers
// (sharedClientForRegion, PageSizeXLarge) live in the root alicloud test
// package. The resource implementation itself lives in alicloud/service/ecs.

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_prefix_list", &resource.Sweeper{
		Name: "alicloud_ecs_prefix_list",
		F:    testAlicloudEcsPrefixList,
	})
}

func testAlicloudEcsPrefixList(region string) error {
	prefixes := []string{
		"tf-testacc",
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := make(map[string]interface{})
	var response map[string]interface{}
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeXLarge

	action := "DescribePrefixLists"
	response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.PrefixLists.PrefixList", response)
	if err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.PrefixLists.PrefixList", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		skip := true
		for _, prefix := range prefixes {
			PrefixListName := ""
			appName, exist := item["PrefixListName"]
			if exist {
				PrefixListName = appName.(string)
			}
			if strings.Contains(strings.ToLower(PrefixListName), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ecs SnapShot Policy: %s (%s)", item["Description"], item["Id"])
			continue
		}
		action = "DeletePrefixList"
		request := map[string]interface{}{
			"PrefixListId": item["PrefixListId"],
		}
		request["RegionId"] = client.RegionId
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ecs PrefixList (%s (%v)): %s", item["Description"].(string), item["Id"].(float64), err)
		}
		// Waiting 30 seconds to ensure these Ecs SnapShot Policy have been deleted.
		time.Sleep(30 * time.Second)

		log.Printf("[INFO] Delete Ecs PrefixList success: %v ", item["Id"].(float64))
	}
	return nil
}
