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

// KMS (密钥管理服务 / Key Management Service, product "Kms" / "2016-01-20",
// RPC style) sweeper gap fill. The three region-scoped "core" resources
// (alicloud_kms_key, alicloud_kms_secret, alicloud_kms_alias) already have
// inline sweepers in their own _test.go files; this file fills the remaining
// gaps:
//
//   - alicloud_kms_network_rule          (ListNetworkRules  -> DeleteNetworkRule, by Name)
//   - alicloud_kms_policy                (ListPolicies      -> DeletePolicy,       by Name)
//   - alicloud_kms_application_access_point (ListApplicationAccessPoints -> DeleteApplicationAccessPoint, by Name)
//   - alicloud_kms_client_key            (parent-scoped: walks matching AAPs, ListClientKeys -> DeleteClientKey, by ClientKeyId)
//   - alicloud_kms_instance             (ListKmsInstances + GetKmsInstance -> ReleaseKmsInstance; only PayAsYouGo/POSTPAY, Subscription/PREPAY deferred as BSS RefundInstance 退订)
//
// Intentionally DEFERRED (no sweep written):
//   - alicloud_kms_key_version: resource Delete is a no-op (KMS key versions are
//     non-deletable; there is no DeleteKeyVersion API), and ListKeyVersions is
//     parent-scoped to a KeyId whose own deletion is scheduled (7-day window).
//     Keys are already swept by alicloud_kms_key; versions vanish with the key.
//   - alicloud_kms_ciphertext: resource Delete is a no-op (an Encrypt blob, not a
//     standalone deletable resource) and there is no list/enumerate API.
//   - alicloud_kms_value_added_service: resource Delete is BssOpenApi
//     RefundInstance (退订) only — no clean KMS release. Deferred per policy.
//
// Fleet-wipe pattern: SWEEP_ALL (ALICLOUD_SWEEP_ALL_RESOURCES=true) bypasses the
// tf-prefix filter so every non-NODELETING resource is swept; NORMAL mode keeps
// the prefix filter. The protectedSweepVpcName NODELETING guard (21 markers) is
// always-active in both modes. Client keys carry no tf-prefix (their id is an
// auto-generated KAAP.* token), so they are swept by walking their parent AAP
// and deleting every client key beneath a name-matching AAP. All sweepers
// return nil — errors are logged and never propagated, so one failing item
// does not abort the cascade.

func init() {
	resource.AddTestSweepers("alicloud_kms_network_rule", &resource.Sweeper{
		Name: "alicloud_kms_network_rule",
		F:    testSweepKmsNetworkRules,
	})
	resource.AddTestSweepers("alicloud_kms_policy", &resource.Sweeper{
		Name: "alicloud_kms_policy",
		F:    testSweepKmsPolicies,
	})
	resource.AddTestSweepers("alicloud_kms_application_access_point", &resource.Sweeper{
		Name: "alicloud_kms_application_access_point",
		F:    testSweepKmsApplicationAccessPoints,
		// Client keys are children of an AAP; sweep them before the AAP itself
		// is removed so DeleteApplicationAccessPoint is not blocked.
		Dependencies: []string{
			"alicloud_kms_client_key",
		},
	})
	resource.AddTestSweepers("alicloud_kms_client_key", &resource.Sweeper{
		Name: "alicloud_kms_client_key",
		F:    testSweepKmsClientKeys,
	})
	resource.AddTestSweepers("alicloud_kms_instance", &resource.Sweeper{
		Name: "alicloud_kms_instance",
		F:    testSweepKmsInstances,
		// Keys / secrets / aliases are swept (keys are scheduled for deletion)
		// before the keystore instance is released.
		Dependencies: []string{
			"alicloud_kms_key",
			"alicloud_kms_secret",
			"alicloud_kms_alias",
		},
	})
}

var kmsSweeperPrefixes = []string{
	"tf-testacc",
	"tf_testacc",
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// kmsGetField returns the first non-empty string value of item[name] across the
// given name variants. KMS response field casing is consistent (PascalCase) but
// the variant list keeps the helper tolerant of camelCase fallbacks.
func kmsGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// kmsSweepNameOk applies the always-active NODELETING guard and the
// sweepAll()-aware tf-prefix filter. Returns true when the resource should be
// swept. SWEEP_ALL bypasses the prefix filter (fleet-wipe: clean all
// non-NODELETING); NORMAL mode keeps the prefix filter.
func kmsSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping kms resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range kmsSweeperPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping kms resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// kmsRpcPost is the shared KMS list/read wrapper: RpcPost with autoRetry=false
// wrapped in retry.Retry + NeedRetry + incrementalWait so throttling and 5xx
// are retried at the sweeper level (the SDK itself does not auto-retry here).
func kmsRpcPost(client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// kmsDeletionProtectionEnabled reports whether the GetKmsInstance field
// DeletionProtection is on. The schema is TypeBool, but the JSON-decoded
// response may also surface a string form; both are handled.
func kmsDeletionProtectionEnabled(v interface{}) bool {
	switch dp := v.(type) {
	case bool:
		return dp
	case string:
		return strings.EqualFold(dp, "true") || strings.EqualFold(dp, "Enabled") || dp == "1"
	}
	return false
}

// listKmsMatchingAapNames enumerates application access points whose Name
// passes kmsSweepNameOk. Shared by the AAP sweeper (to delete them) and the
// client_key sweeper (to walk each AAP and delete its client keys).
func listKmsMatchingAapNames(client *connectivity.AliyunClient) ([]string, error) {
	names := make([]string, 0)
	request := map[string]interface{}{
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
	}
	action := "ListApplicationAccessPoints"
	for {
		response, err := kmsRpcPost(client, action, request)
		if err != nil {
			return nil, fmt.Errorf("ListApplicationAccessPoints: %s", err)
		}
		resp, err := jsonpath.Get("$.ApplicationAccessPoints.ApplicationAccessPoint", response)
		if err != nil {
			return nil, fmt.Errorf("ListApplicationAccessPoints jsonpath $.ApplicationAccessPoints.ApplicationAccessPoint: %s", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := kmsGetField(item, "Name")
			if kmsSweepNameOk(name) {
				names = append(names, name)
			}
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return names, nil
}

func testSweepKmsNetworkRules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
	}
	action := "ListNetworkRules"
	for {
		response, err := kmsRpcPost(client, action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list kms network_rules: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.NetworkRules.NetworkRule", response)
		if err != nil {
			log.Printf("[ERROR] ListNetworkRules jsonpath $.NetworkRules.NetworkRule: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := kmsGetField(item, "Name")
			if !kmsSweepNameOk(name) {
				continue
			}
			log.Printf("[INFO] Deleting kms network_rule: %s", name)
			deleteRequest := map[string]interface{}{"Name": name}
			if _, e := client.RpcPost("Kms", "2016-01-20", "DeleteNetworkRule", nil, deleteRequest, false); e != nil {
				log.Printf("[ERROR] Failed to delete kms network_rule (%s): %s", name, e)
				continue
			}
			log.Printf("[INFO] Delete kms network_rule success: %s", name)
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func testSweepKmsPolicies(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
	}
	action := "ListPolicies"
	for {
		response, err := kmsRpcPost(client, action, request)
		if err != nil {
			log.Printf("[ERROR] Failed to list kms policies: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			log.Printf("[ERROR] ListPolicies jsonpath $.Policies.Policy: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := kmsGetField(item, "Name")
			if !kmsSweepNameOk(name) {
				continue
			}
			log.Printf("[INFO] Deleting kms policy: %s", name)
			deleteRequest := map[string]interface{}{"Name": name}
			if _, e := client.RpcPost("Kms", "2016-01-20", "DeletePolicy", nil, deleteRequest, false); e != nil {
				log.Printf("[ERROR] Failed to delete kms policy (%s): %s", name, e)
				continue
			}
			log.Printf("[INFO] Delete kms policy success: %s", name)
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func testSweepKmsApplicationAccessPoints(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	names, err := listKmsMatchingAapNames(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list kms application_access_points: %s", err)
		return nil
	}
	for _, name := range names {
		log.Printf("[INFO] Deleting kms application_access_point: %s", name)
		deleteRequest := map[string]interface{}{"Name": name}
		if _, e := client.RpcPost("Kms", "2016-01-20", "DeleteApplicationAccessPoint", nil, deleteRequest, false); e != nil {
			log.Printf("[ERROR] Failed to delete kms application_access_point (%s): %s", name, e)
			continue
		}
		log.Printf("[INFO] Delete kms application_access_point success: %s", name)
	}
	return nil
}

func testSweepKmsClientKeys(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Client keys have no tf-prefix (their id is an auto-generated KAAP.* token),
	// so sweep by walking each name-matching AAP and deleting every client key
	// beneath it. AapName is the documented (optional) filter on ListClientKeys.
	aapNames, err := listKmsMatchingAapNames(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list kms application_access_points for client_key sweep: %s", err)
		return nil
	}
	for _, aapName := range aapNames {
		request := map[string]interface{}{
			"AapName":    aapName,
			"PageSize":   PageSizeXLarge,
			"PageNumber": 1,
		}
		action := "ListClientKeys"
		for {
			response, err := kmsRpcPost(client, action, request)
			if err != nil {
				log.Printf("[ERROR] Failed to list kms client_keys (aap %s): %s", aapName, err)
				break
			}
			// $.ClientKeys is a direct array (no nested itemNode).
			resp, err := jsonpath.Get("$.ClientKeys", response)
			if err != nil {
				log.Printf("[ERROR] ListClientKeys jsonpath $.ClientKeys (aap %s): %s", aapName, err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, _ := v.(map[string]interface{})
				clientKeyId := kmsGetField(item, "ClientKeyId")
				if clientKeyId == "" {
					continue
				}
				log.Printf("[INFO] Deleting kms client_key: %s (aap %s)", clientKeyId, aapName)
				deleteRequest := map[string]interface{}{"ClientKeyId": clientKeyId}
				if _, e := client.RpcPost("Kms", "2016-01-20", "DeleteClientKey", nil, deleteRequest, false); e != nil {
					log.Printf("[ERROR] Failed to delete kms client_key (%s, aap %s): %s", clientKeyId, aapName, e)
					continue
				}
				log.Printf("[INFO] Delete kms client_key success: %s (aap %s)", clientKeyId, aapName)
			}
			if len(result) < PageSizeXLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return nil
}

func testSweepKmsInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// ListKmsInstances returns only {KmsInstanceId, KmsInstanceArn} — no name,
	// no charge type — so each instance is described with GetKmsInstance to read
	// InstanceName + ChargeType before the prefix/NODELETING gate is applied.
	listRequest := map[string]interface{}{
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
	}
	listAction := "ListKmsInstances"
	for {
		response, err := kmsRpcPost(client, listAction, listRequest)
		if err != nil {
			log.Printf("[ERROR] Failed to list kms instances: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.KmsInstances.KmsInstance", response)
		if err != nil {
			log.Printf("[ERROR] ListKmsInstances jsonpath $.KmsInstances.KmsInstance: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			instanceId := kmsGetField(item, "KmsInstanceId")
			if instanceId == "" {
				continue
			}
			// Describe to obtain InstanceName + ChargeType.
			describeRequest := map[string]interface{}{"KmsInstanceId": instanceId}
			describeResponse, dErr := kmsRpcPost(client, "GetKmsInstance", describeRequest)
			if dErr != nil {
				log.Printf("[ERROR] Failed to describe kms instance %s: %s", instanceId, dErr)
				continue
			}
			objRaw, jErr := jsonpath.Get("$.KmsInstance", describeResponse)
			if jErr != nil {
				log.Printf("[ERROR] GetKmsInstance jsonpath $.KmsInstance (%s): %s", instanceId, jErr)
				continue
			}
			obj, _ := objRaw.(map[string]interface{})
			if obj == nil {
				continue
			}
			instanceName := kmsGetField(obj, "InstanceName")
			if !kmsSweepNameOk(instanceName) {
				continue
			}
			// ChargeType: PREPAY = Subscription (BSS RefundInstance 退订) -> defer;
			// POSTPAY = PayAsYouGo -> ReleaseKmsInstance (clean KMS delete).
			chargeType := kmsGetField(obj, "ChargeType")
			if strings.EqualFold(chargeType, "PREPAY") {
				log.Printf("[INFO] Skipping kms instance %s (%s): Subscription/PREPAY — deferred (BSS RefundInstance 退订, requires explicit approval)", instanceId, instanceName)
				continue
			}
			// Disable deletion protection (best-effort) before release, mirroring
			// the resource Delete path which clears it via SetDeletionProtection.
			if kmsDeletionProtectionEnabled(obj["DeletionProtection"]) {
				dpRequest := map[string]interface{}{
					"KmsInstanceId":            instanceId,
					"EnableDeletionProtection": false,
				}
				if _, e := client.RpcPost("Kms", "2016-01-20", "SetDeletionProtection", nil, dpRequest, false); e != nil {
					log.Printf("[ERROR] Failed to clear kms instance deletion_protection (%s): %s", instanceId, e)
				}
			}
			log.Printf("[INFO] Deleting kms instance: %s (%s)", instanceId, instanceName)
			releaseRequest := map[string]interface{}{"KmsInstanceId": instanceId}
			if _, e := client.RpcPost("Kms", "2016-01-20", "ReleaseKmsInstance", nil, releaseRequest, false); e != nil {
				log.Printf("[ERROR] Failed to release kms instance (%s, %s): %s", instanceId, instanceName, e)
				continue
			}
			log.Printf("[INFO] Delete kms instance success: %s (%s)", instanceId, instanceName)
		}
		if len(result) < PageSizeXLarge {
			break
		}
		listRequest["PageNumber"] = listRequest["PageNumber"].(int) + 1
	}
	return nil
}
