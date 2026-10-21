package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers("alicloud_mongodb_global_security_ip_group", &resource.Sweeper{
		Name: "alicloud_mongodb_global_security_ip_group",
		F:    testSweepMongoDBGlobalSecurityIPGroups,
	})
}

// testSweepMongoDBGlobalSecurityIPGroups sweeps alicloud_mongodb_global_security_ip_group
// resources. It lists all global IP whitelists in the region via
// DescribeGlobalSecurityIPGroup (GlobalSecurityGroupId is optional) and deletes
// those matching the test prefixes (or all when ALICLOUD_SWEEP_ALL_RESOURCES=true),
// skipping any whose name carries a NODELETING/do-not-delete marker.
func testSweepMongoDBGlobalSecurityIPGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tfacc",
	}

	action := "DescribeGlobalSecurityIPGroup"
	request := map[string]interface{}{}
	query := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Dds", "2015-12-01", action, query, request)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.GlobalSecurityIPGroup", response)
	if err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.GlobalSecurityIPGroup", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["GlobalIgName"])
		if protectedSweepVpcName(name) {
			log.Printf("[INFO] Skipping Mongodb GlobalSecurityIPGroup (NODELETING): %s", name)
			continue
		}

		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Mongodb GlobalSecurityIPGroup: %s", name)
				continue
			}
		}
		delAction := "DeleteGlobalSecurityIPGroup"
		delRequest := map[string]interface{}{
			"GlobalSecurityGroupId": item["GlobalSecurityGroupId"],
			"GlobalIgName":          name,
			"RegionId":              client.RegionId,
		}
		_, err = client.RpcPost("Dds", "2015-12-01", delAction, nil, delRequest, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Mongodb GlobalSecurityIPGroup (%s): %s", name, err)
			continue
		}
		log.Printf("[INFO] Delete Mongodb GlobalSecurityIPGroup success: %s ", name)
	}
	return nil
}
