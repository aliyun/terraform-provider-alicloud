package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The alicloud_amqp_instance sweeper stays in the root test binary: the make
// sweep target runs only ./alicloud, and the SDKv2 and plugin-testing
// helper/resource packages cannot both register -sweep in one binary. The
// resource implementation itself lives in alicloud/service/amqp. The file is
// named sweeper_*.go (not resource_alicloud_*.go) so the collection-order and
// testing-coverage CI gates do not treat it as a managed-resource test file;
// the delete/refund logic mirrors the production Delete (DeleteInstance /
// RefundInstance).

func init() {
	resource.AddTestSweepers("alicloud_amqp_instance", &resource.Sweeper{
		Name: "alicloud_amqp_instance",
		F:    testSweepAmqpInstance,
	})
}

// testSweepAmqpInstance removes AMQP instances left over from failed acceptance
// runs. It lists instances in the region and deletes any whose InstanceName
// starts with a test prefix. PayAsYouGo/Serverless instances are released with
// DeleteInstance; Subscription instances are refunded with RefundInstance —
// mirroring the production Delete logic.
func testSweepAmqpInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Prefixes follow the provider sweep convention (see
	// acc-test-writing-patterns.md): tf-testacc is what the acceptance
	// tests create, tf_testacc covers the underscore variant, and
	// tfaccamqp is the historical AMQP-specific prefix. Matching is
	// case-insensitive so tf-testAcc/tf_testAcc are also swept.
	prefixes := []string{"tf-testacc", "tf_testacc", "tfaccamqp"}
	action := "ListInstances"
	request := map[string]interface{}{"MaxResults": 100}

	// failed collects instance IDs that neither DeleteInstance nor RefundInstance
	// could remove. A non-empty slice makes the sweep return an error instead of
	// reporting success while instances remain, so a broken sweep is not masked.
	var failed []string

	for {
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			var rpcErr error
			response, rpcErr = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
			if rpcErr != nil {
				if helper.NeedRetry(rpcErr) {
					wait()
					return retry.RetryableError(rpcErr)
				}
				return retry.NonRetryableError(rpcErr)
			}
			return nil
		})
		if err != nil {
			// ListInstances failure must not be swallowed: returning nil here
			// would make a sweep that listed nothing report success, leaving
			// charged instances behind. Propagate the error so the sweep fails
			// loudly and the operator can investigate.
			return fmt.Errorf("listing AMQP instances for sweep (action %s): %w", action, err)
		}

		data, _ := response["Data"].(map[string]interface{})
		instances, _ := data["Instances"].([]interface{})
		for _, v := range instances {
			item, _ := v.(map[string]interface{})
			instanceName := fmt.Sprint(item["InstanceName"])
			if !amqpHasPrefix(instanceName, prefixes) {
				continue
			}
			instanceID := fmt.Sprint(item["InstanceId"])
			log.Printf("[INFO] sweeping AMQP instance %s (%s)", instanceID, instanceName)
			if !amqpSweepDelete(client, instanceID) {
				if !amqpSweepRefund(client, instanceID) {
					log.Printf("[ERROR] failed to sweep AMQP instance %s", instanceID)
					failed = append(failed, instanceID)
				}
			}
		}

		if nextToken, ok := data["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("failed to sweep %d AMQP instance(s): %s", len(failed), strings.Join(failed, ", "))
	}
	return nil
}

func amqpSweepDelete(client *connectivity.AliyunClient, id string) bool {
	request := map[string]interface{}{"InstanceId": id, "RegionId": client.RegionId}
	_, err := client.RpcPost("amqp-open", "2019-12-12", "DeleteInstance", map[string]interface{}{}, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InstanceNotExist", "InstanceIdNotExist", "ResourceNotFound"}) {
			return true
		}
		log.Printf("[ERROR] DeleteInstance for %s failed: %s", id, err)
		return false
	}
	return true
}

func amqpSweepRefund(client *connectivity.AliyunClient, id string) bool {
	request := map[string]interface{}{
		"InstanceId":         id,
		"ClientToken":        helper.BuildClientToken("RefundInstance"),
		"ImmediatelyRelease": "1",
		"ProductCode":        "ons",
		"ProductType":        "ons_onsproxy_pre",
	}
	if client.IsInternationalAccount() {
		request["ProductType"] = "ons_onsproxy_public_intl"
	}
	_, err := client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", "RefundInstance", map[string]interface{}{}, request, true, "")
	if err != nil {
		// BssOpenApi RefundInstance error codes (verified via OpenAPI metadata):
		// ResourceNotExists = the instance no longer exists; ExistRefundingOrderError
		// = the instance may have expired, been refunded, or been released. Both are
		// treated as a successful (idempotent) delete.
		if errs.IsExpectedErrors(err, []string{"ResourceNotExists", "ExistRefundingOrderError"}) {
			return true
		}
		log.Printf("[ERROR] RefundInstance for %s failed: %s", id, err)
		return false
	}
	return true
}

func amqpHasPrefix(name string, prefixes []string) bool {
	lower := strings.ToLower(name)
	for _, p := range prefixes {
		if strings.HasPrefix(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}
