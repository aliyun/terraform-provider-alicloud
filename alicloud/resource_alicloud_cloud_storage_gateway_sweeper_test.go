package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// csgSweepPrefixes is the set of test-resource name prefixes the CSG sweepers
// match when ALICLOUD_SWEEP_ALL_RESOURCES is not set.
var csgSweepPrefixes = []string{"tf-testAcc", "tf_testAcc"}

// csgShouldSweepName decides whether a single CSG resource (by its human-readable
// name) should be swept. It enforces the red-line NODELETING name-based
// protection first (reusing the generic VPC name matcher), then honours
// sweepAll() to bypass the prefix filter.
func csgShouldSweepName(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	ln := strings.ToLower(name)
	for _, p := range csgSweepPrefixes {
		if strings.HasPrefix(ln, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// csgListStorageBundles enumerates every Cloud Storage Gateway storage bundle
// visible to the caller in the region (DescribeStorageBundles, paginated).
// BackendBucketRegionId is mandatory for this action, so it is set to the
// client's region (the sweep region); test storage bundles created in the
// region carry the same backend bucket region.
func csgListStorageBundles(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	action := "DescribeStorageBundles"
	out := make([]map[string]interface{}, 0)
	request := map[string]interface{}{
		"BackendBucketRegionId": client.RegionId,
		"PageSize":              PageSizeLarge,
		"PageNumber":            1,
	}
	for {
		var response map[string]interface{}
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
			log.Printf("[ERROR] %s get an error: %s", action, err)
			return out, err
		}
		if fmt.Sprint(response["Success"]) == "false" {
			log.Printf("[ERROR] %s failed, response: %v", action, response)
			return out, nil
		}
		raw, err := jsonpath.Get("$.StorageBundles.StorageBundle", response)
		if err != nil {
			if formatInt(response["TotalCount"]) == 0 {
				return out, nil
			}
			log.Printf("[ERROR] Getting %s attribute %s failed: %s", action, "$.StorageBundles.StorageBundle", err)
			return out, nil
		}
		result, _ := raw.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				out = append(out, item)
			}
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return out, nil
}

// csgListGateways enumerates every Cloud Storage Gateway gateway in the region.
// DescribeGateways is scoped to a storage bundle, so this iterates all storage
// bundles first and aggregates their gateways. Each returned item carries the
// parent StorageBundleId, the gateway Name, GatewayId and IsPostPaid flag.
func csgListGateways(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	out := make([]map[string]interface{}, 0)
	bundles, err := csgListStorageBundles(client)
	if err != nil {
		return out, err
	}
	action := "DescribeGateways"
	for _, b := range bundles {
		storageBundleId := fmt.Sprint(b["StorageBundleId"])
		request := map[string]interface{}{
			"StorageBundleId": storageBundleId,
			"PageSize":        PageSizeLarge,
			"PageNumber":      1,
		}
		for {
			var response map[string]interface{}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(5*time.Minute, func() *retry.RetryError {
				response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
				log.Printf("[ERROR] %s (storage bundle %s) get an error: %s", action, storageBundleId, err)
				break
			}
			if fmt.Sprint(response["Success"]) == "false" {
				log.Printf("[ERROR] %s (storage bundle %s) failed, response: %v", action, storageBundleId, response)
				break
			}
			raw, jerr := jsonpath.Get("$.Gateways.Gateway", response)
			if jerr != nil {
				if formatInt(response["TotalCount"]) == 0 {
					break
				}
				log.Printf("[ERROR] Getting %s attribute %s failed: %s", action, "$.Gateways.Gateway", jerr)
				break
			}
			result, _ := raw.([]interface{})
			for _, v := range result {
				if item, ok := v.(map[string]interface{}); ok {
					out = append(out, item)
				}
			}
			if len(result) < request["PageSize"].(int) {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return out, nil
}

// csgListGatewayChildren lists child resources of one gateway via the supplied
// sgw action/path. The gateway id is injected as "GatewayId". Pagination is
// honoured when the action supports it.
func csgListGatewayChildren(client *connectivity.AliyunClient, gatewayId, action, path string, paginate bool) ([]map[string]interface{}, error) {
	out := make([]map[string]interface{}, 0)
	request := map[string]interface{}{
		"GatewayId": gatewayId,
	}
	if paginate {
		request["PageSize"] = PageSizeLarge
		request["PageNumber"] = 1
	}
	for {
		var response map[string]interface{}
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
			if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
				return out, nil
			}
			log.Printf("[ERROR] %s (gateway %s) get an error: %s", action, gatewayId, err)
			return out, err
		}
		if fmt.Sprint(response["Success"]) == "false" {
			log.Printf("[ERROR] %s (gateway %s) failed, response: %v", action, gatewayId, response)
			return out, nil
		}
		raw, jerr := jsonpath.Get(path, response)
		if jerr != nil {
			return out, nil
		}
		result, _ := raw.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				out = append(out, item)
			}
		}
		if !paginate || len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return out, nil
}

func init() {
	// Children must be swept before their parents so that gateways and storage
	// bundles can be deleted cleanly. Dependencies lists sweepers that must run
	// before the current one.
	resource.AddTestSweepers("alicloud_cloud_storage_gateway_storage_bundle",
		&resource.Sweeper{
			Name:         "alicloud_cloud_storage_gateway_storage_bundle",
			Dependencies: []string{"alicloud_cloud_storage_gateway_gateway"},
			F:            testSweepCloudStorageGatewayStorageBundle,
		})

	resource.AddTestSweepers("alicloud_cloud_storage_gateway_gateway",
		&resource.Sweeper{
			Name: "alicloud_cloud_storage_gateway_gateway",
			Dependencies: []string{
				"alicloud_cloud_storage_gateway_gateway_file_share",
				"alicloud_cloud_storage_gateway_gateway_cache_disk",
				"alicloud_cloud_storage_gateway_gateway_block_volume",
				"alicloud_cloud_storage_gateway_gateway_smb_user",
				"alicloud_cloud_storage_gateway_gateway_logging",
			},
			F: testSweepCloudStorageGatewayGateway,
		})

	resource.AddTestSweepers("alicloud_cloud_storage_gateway_gateway_smb_user",
		&resource.Sweeper{
			Name: "alicloud_cloud_storage_gateway_gateway_smb_user",
			F:    testSweepCloudStorageGatewayGatewaySmbUser,
		})

	resource.AddTestSweepers("alicloud_cloud_storage_gateway_gateway_cache_disk",
		&resource.Sweeper{
			Name: "alicloud_cloud_storage_gateway_gateway_cache_disk",
			F:    testSweepCloudStorageGatewayGatewayCacheDisk,
		})

	resource.AddTestSweepers("alicloud_cloud_storage_gateway_gateway_file_share",
		&resource.Sweeper{
			Name: "alicloud_cloud_storage_gateway_gateway_file_share",
			F:    testSweepCloudStorageGatewayGatewayFileShare,
		})

	resource.AddTestSweepers("alicloud_cloud_storage_gateway_gateway_block_volume",
		&resource.Sweeper{
			Name: "alicloud_cloud_storage_gateway_gateway_block_volume",
			F:    testSweepCloudStorageGatewayGatewayBlockVolume,
		})

	resource.AddTestSweepers("alicloud_cloud_storage_gateway_gateway_logging",
		&resource.Sweeper{
			Name: "alicloud_cloud_storage_gateway_gateway_logging",
			F:    testSweepCloudStorageGatewayGatewayLogging,
		})

	resource.AddTestSweepers("alicloud_cloud_storage_gateway_express_sync_share_attachment",
		&resource.Sweeper{
			Name:         "alicloud_cloud_storage_gateway_express_sync_share_attachment",
			Dependencies: []string{"alicloud_cloud_storage_gateway_express_sync"},
			F:            testSweepCloudStorageGatewayExpressSyncShareAttachment,
		})
}

func testSweepCloudStorageGatewayStorageBundle(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	bundles, _ := csgListStorageBundles(client)
	for _, b := range bundles {
		name := fmt.Sprint(b["Name"])
		id := fmt.Sprint(b["StorageBundleId"])
		if !csgShouldSweepName(name) {
			log.Printf("[INFO] Skipping CloudStorageGateway StorageBundle: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting CloudStorageGateway StorageBundle: %s (%s)", name, id)
		action := "DeleteStorageBundle"
		request := map[string]interface{}{"StorageBundleId": id}
		if _, err := client.RpcPost("sgw", "2018-05-11", action, nil, request, false); err != nil {
			if IsExpectedErrors(err, []string{"StorageBundleNotExist"}) {
				continue
			}
			log.Printf("[ERROR] Failed to delete CloudStorageGateway StorageBundle (%s): %s", id, err)
		}
	}
	return nil
}

func testSweepCloudStorageGatewayGateway(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	gateways, _ := csgListGateways(client)
	for _, g := range gateways {
		name := fmt.Sprint(g["Name"])
		id := fmt.Sprint(g["GatewayId"])
		if !csgShouldSweepName(name) {
			log.Printf("[INFO] Skipping CloudStorageGateway Gateway: %s (%s)", name, id)
			continue
		}
		// Subscription (prepaid) gateways cannot be released by DeleteGateway;
		// the provider Delete refuses them, so skip them here too.
		if v, ok := g["IsPostPaid"].(bool); ok && !v {
			log.Printf("[INFO] Skipping CloudStorageGateway Gateway (Subscription): %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting CloudStorageGateway Gateway: %s (%s)", name, id)
		action := "DeleteGateway"
		request := map[string]interface{}{"GatewayId": id}
		if _, err := client.RpcPost("sgw", "2018-05-11", action, nil, request, true); err != nil {
			if IsExpectedErrors(err, []string{"GatewayNotExist"}) || NotFoundError(err) {
				continue
			}
			log.Printf("[ERROR] Failed to delete CloudStorageGateway Gateway (%s): %s", id, err)
		}
	}
	return nil
}

// csgForEachSweepableGateway iterates every gateway whose own name passes the
// sweep filter, calling fn(gatewayId, gatewayName) for each. Children of
// non-sweepable gateways (real存 or non-test) are never touched.
func csgForEachSweepableGateway(client *connectivity.AliyunClient, fn func(gatewayId, gatewayName string)) {
	gateways, _ := csgListGateways(client)
	for _, g := range gateways {
		name := fmt.Sprint(g["Name"])
		id := fmt.Sprint(g["GatewayId"])
		if !csgShouldSweepName(name) {
			log.Printf("[INFO] Skipping CloudStorageGateway Gateway children: %s (%s)", name, id)
			continue
		}
		fn(id, name)
	}
}

func testSweepCloudStorageGatewayGatewaySmbUser(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	csgForEachSweepableGateway(client, func(gatewayId, gatewayName string) {
		users, _ := csgListGatewayChildren(client, gatewayId, "DescribeGatewaySMBUsers", "$.Users.User", true)
		for _, u := range users {
			username := fmt.Sprint(u["Username"])
			if !csgShouldSweepName(username) {
				log.Printf("[INFO] Skipping CloudStorageGateway GatewaySmbUser: %s (%s:%s)", username, gatewayId, username)
				continue
			}
			log.Printf("[INFO] Deleting CloudStorageGateway GatewaySmbUser: %s (%s:%s)", username, gatewayId, username)
			action := "DeleteGatewaySMBUser"
			request := map[string]interface{}{
				"GatewayId": gatewayId,
				"Username":  username,
			}
			if _, err := client.RpcPost("sgw", "2018-05-11", action, nil, request, false); err != nil {
				if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
					return
				}
				log.Printf("[ERROR] Failed to delete CloudStorageGateway GatewaySmbUser (%s:%s): %s", gatewayId, username, err)
			}
		}
	})
	return nil
}

func testSweepCloudStorageGatewayGatewayCacheDisk(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	csgForEachSweepableGateway(client, func(gatewayId, gatewayName string) {
		caches, _ := csgListGatewayChildren(client, gatewayId, "DescribeGatewayCaches", "$.Caches.Cache", false)
		for _, c := range caches {
			cacheId := fmt.Sprint(c["CacheId"])
			localFilePath := fmt.Sprint(c["LocalFilePath"])
			// Cache disks have no Name field; the local file path is the only
			// user-set identifier, so use it for NODELETING/prefix filtering.
			if !csgShouldSweepName(localFilePath) {
				log.Printf("[INFO] Skipping CloudStorageGateway GatewayCacheDisk: %s (%s:%s:%s)", localFilePath, gatewayId, cacheId, localFilePath)
				continue
			}
			log.Printf("[INFO] Deleting CloudStorageGateway GatewayCacheDisk: %s (%s:%s:%s)", localFilePath, gatewayId, cacheId, localFilePath)
			action := "DeleteGatewayCacheDisk"
			request := map[string]interface{}{
				"GatewayId":     gatewayId,
				"CacheId":       cacheId,
				"LocalFilePath": localFilePath,
			}
			if _, err := client.RpcPost("sgw", "2018-05-11", action, nil, request, true); err != nil {
				if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
					return
				}
				log.Printf("[ERROR] Failed to delete CloudStorageGateway GatewayCacheDisk (%s:%s): %s", gatewayId, cacheId, err)
			}
		}
	})
	return nil
}

func testSweepCloudStorageGatewayGatewayFileShare(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	csgForEachSweepableGateway(client, func(gatewayId, gatewayName string) {
		shares, _ := csgListGatewayChildren(client, gatewayId, "DescribeGatewayFileShares", "$.FileShares.FileShare", false)
		for _, s := range shares {
			name := fmt.Sprint(s["Name"])
			indexId := fmt.Sprint(s["IndexId"])
			if !csgShouldSweepName(name) {
				log.Printf("[INFO] Skipping CloudStorageGateway GatewayFileShare: %s (%s:%s)", name, gatewayId, indexId)
				continue
			}
			log.Printf("[INFO] Deleting CloudStorageGateway GatewayFileShare: %s (%s:%s)", name, gatewayId, indexId)
			action := "DeleteGatewayFileShares"
			request := map[string]interface{}{
				"GatewayId": gatewayId,
				"IndexId":   indexId,
				"Force":     true,
			}
			if _, err := client.RpcPost("sgw", "2018-05-11", action, nil, request, false); err != nil {
				if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
					return
				}
				log.Printf("[ERROR] Failed to delete CloudStorageGateway GatewayFileShare (%s:%s): %s", gatewayId, indexId, err)
			}
		}
	})
	return nil
}

func testSweepCloudStorageGatewayGatewayBlockVolume(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	csgForEachSweepableGateway(client, func(gatewayId, gatewayName string) {
		volumes, _ := csgListGatewayChildren(client, gatewayId, "DescribeGatewayBlockVolumes", "$.BlockVolumes.BlockVolume", false)
		for _, bv := range volumes {
			name := fmt.Sprint(bv["Name"])
			indexId := fmt.Sprint(bv["IndexId"])
			if !csgShouldSweepName(name) {
				log.Printf("[INFO] Skipping CloudStorageGateway GatewayBlockVolume: %s (%s:%s)", name, gatewayId, indexId)
				continue
			}
			log.Printf("[INFO] Deleting CloudStorageGateway GatewayBlockVolume: %s (%s:%s)", name, gatewayId, indexId)
			action := "DeleteGatewayBlockVolumes"
			request := map[string]interface{}{
				"GatewayId":       gatewayId,
				"IndexId":         indexId,
				"IsSourceDeletion": true,
			}
			if _, err := client.RpcPost("sgw", "2018-05-11", action, nil, request, false); err != nil {
				if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
					return
				}
				log.Printf("[ERROR] Failed to delete CloudStorageGateway GatewayBlockVolume (%s:%s): %s", gatewayId, indexId, err)
			}
		}
	})
	return nil
}

func testSweepCloudStorageGatewayGatewayLogging(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	csgForEachSweepableGateway(client, func(gatewayId, gatewayName string) {
		action := "DescribeGatewayLogging"
		request := map[string]interface{}{"GatewayId": gatewayId}
		response, err := client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
				return
			}
			log.Printf("[ERROR] %s (gateway %s) get an error: %s", action, gatewayId, err)
			return
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return
		}
		status := fmt.Sprint(response["GatewayLoggingStatus"])
		if status == "" || status == "None" {
			return
		}
		log.Printf("[INFO] Deleting CloudStorageGateway GatewayLogging: %s (%s)", gatewayName, gatewayId)
		delAction := "DeleteGatewayLogging"
		delRequest := map[string]interface{}{"GatewayId": gatewayId}
		if _, err := client.RpcPost("sgw", "2018-05-11", delAction, nil, delRequest, false); err != nil {
			if IsExpectedErrors(err, []string{"GatewayNotExist"}) {
				return
			}
			log.Printf("[ERROR] Failed to delete CloudStorageGateway GatewayLogging (%s): %s", gatewayId, err)
		}
	})
	return nil
}

func testSweepCloudStorageGatewayExpressSyncShareAttachment(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeExpressSyncs"
	request := map[string]interface{}{}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
		log.Printf("[ERROR] %s get an error: %s", action, err)
		return nil
	}
	if fmt.Sprint(response["Success"]) == "false" {
		log.Printf("[ERROR] %s failed, response: %v", action, response)
		return nil
	}
	raw, jerr := jsonpath.Get("$.ExpressSyncs.ExpressSync", response)
	if jerr != nil {
		if formatInt(response["TotalCount"]) != 0 {
			log.Printf("[ERROR] Getting %s attribute %s failed: %s", action, "$.ExpressSyncs.ExpressSync", jerr)
		}
		return nil
	}
	expressSyncs, _ := raw.([]interface{})
	for _, v := range expressSyncs {
		es, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		expressSyncId := fmt.Sprint(es["ExpressSyncId"])
		expressSyncName := fmt.Sprint(es["Name"])
		// Gate share deletion by the express sync's own sweepability too: a
		// NODELETING or non-test express sync's attachments are left alone.
		if !csgShouldSweepName(expressSyncName) {
			log.Printf("[INFO] Skipping CloudStorageGateway ExpressSyncShareAttachment (express sync): %s (%s)", expressSyncName, expressSyncId)
			continue
		}
		shares, _ := csgListExpressSyncShares(client, expressSyncId)
		for _, sh := range shares {
			shareName := fmt.Sprint(sh["ShareName"])
			gatewayId := fmt.Sprint(sh["GatewayId"])
			if !csgShouldSweepName(shareName) {
				log.Printf("[INFO] Skipping CloudStorageGateway ExpressSyncShareAttachment: %s (%s:%s:%s)", shareName, expressSyncId, gatewayId, shareName)
				continue
			}
			log.Printf("[INFO] Deleting CloudStorageGateway ExpressSyncShareAttachment: %s (%s:%s:%s)", shareName, expressSyncId, gatewayId, shareName)
			sharesMaps := map[string]interface{}{gatewayId: shareName}
			jsonStr, cerr := convertMaptoJsonString(sharesMaps)
			if cerr != nil {
				log.Printf("[ERROR] Failed to marshal GatewayShares for ExpressSyncShareAttachment (%s:%s:%s): %s", expressSyncId, gatewayId, shareName, cerr)
				continue
			}
			delAction := "RemoveSharesFromExpressSync"
			delRequest := map[string]interface{}{
				"ExpressSyncId": expressSyncId,
				"GatewayShares": jsonStr,
			}
			if _, err := client.RpcPost("sgw", "2018-05-11", delAction, nil, delRequest, false); err != nil {
				if IsExpectedErrors(err, []string{"ExpressSyncNotExist", "GatewayNotExist"}) {
					continue
				}
				log.Printf("[ERROR] Failed to delete CloudStorageGateway ExpressSyncShareAttachment (%s:%s:%s): %s", expressSyncId, gatewayId, shareName, err)
			}
		}
	}
	return nil
}

// csgListExpressSyncShares lists the shares attached to one express sync.
func csgListExpressSyncShares(client *connectivity.AliyunClient, expressSyncId string) ([]map[string]interface{}, error) {
	out := make([]map[string]interface{}, 0)
	action := "DescribeExpressSyncShares"
	request := map[string]interface{}{"ExpressSyncIds": expressSyncId}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		var err error
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"GatewayNotExist", "ExpressSyncNotExist"}) {
			return out, nil
		}
		log.Printf("[ERROR] %s (express sync %s) get an error: %s", action, expressSyncId, err)
		return out, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return out, nil
	}
	raw, jerr := jsonpath.Get("$.Shares.Share", response)
	if jerr != nil {
		return out, nil
	}
	result, _ := raw.([]interface{})
	for _, v := range result {
		if item, ok := v.(map[string]interface{}); ok {
			out = append(out, item)
		}
	}
	return out, nil
}
