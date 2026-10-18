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

// testSweepMongoDBGlobalSecurityIPGroups sweeps global IP-whitelist templates
// (a standalone mongodb resource not tied to a single instance — instance
// deletion does NOT cascade-delete these). DescribeGlobalSecurityIPGroup with
// no GlobalSecurityGroupId lists every template in the region. Each is deleted
// via DeleteGlobalSecurityIPGroup. NODELETING/do-not-delete markers on the
// template name (GlobalIgName) are honored in ALL modes (including SWEEP_ALL)
// via protectedSweepVpcName.
func testSweepMongoDBGlobalSecurityIPGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeGlobalSecurityIPGroup"
	query := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	request := map[string]interface{}{}

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
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

	resp, err := jsonpath.Get("$.GlobalSecurityIPGroup[*]", response)
	if err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.GlobalSecurityIPGroup[*]", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		name := fmt.Sprint(item["GlobalIgName"])

		// NODELETING / do-not-delete marker guard — applies in ALL modes
		// (including SWEEP_ALL) so named-protected templates are never swept.
		if protectedSweepVpcName(name) {
			log.Printf("[INFO] Skipping Mongodb GlobalSecurityIPGroup (do-not-delete marker): %s", name)
			continue
		}
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Mongodb GlobalSecurityIPGroup: %s", name)
				continue
			}
		}
		delAction := "DeleteGlobalSecurityIPGroup"
		delQuery := map[string]interface{}{}
		delRequest := map[string]interface{}{
			"GlobalSecurityGroupId": item["GlobalSecurityGroupId"],
			"RegionId":               client.RegionId,
			"GlobalIgName":          item["GlobalIgName"],
		}
		_, err = client.RpcPost("Dds", "2015-12-01", delAction, delQuery, delRequest, true)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Mongodb GlobalSecurityIPGroup (%s): %s", name, err)
			continue
		}
		log.Printf("[INFO] Delete Mongodb GlobalSecurityIPGroup success: %s ", name)
	}
	return nil
}
