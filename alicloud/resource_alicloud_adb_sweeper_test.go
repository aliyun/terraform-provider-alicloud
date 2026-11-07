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

// adbSweepPrefixes are the test-resource name prefixes used by the ADB
// acceptance tests. When ALICLOUD_SWEEP_ALL_RESOURCES is not set, only
// sub-resources whose name matches one of these prefixes are swept.
var adbSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tfaccadb",
	"tftestacc",
	"tf-testacc",
}

// listAdbClustersForSweep enumerates every AnalyticDB (adb 2019-03-15)
// cluster in the region via DescribeDBClusters. Clusters whose
// DBClusterDescription carries a do-not-delete marker (protectedSweepVpcName)
// are excluded from the returned slice so that every sub-resource sweeper
// consistently skips them.
func listAdbClustersForSweep(client *connectivity.AliyunClient) []map[string]interface{} {
	action := "DescribeDBClusters"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	clusters := make([]map[string]interface{}, 0)
	for {
		response, err := client.RpcPost("adb", "2019-03-15", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			return clusters
		}
		resp, err := jsonpath.Get("$.Items.DBCluster", response)
		if err != nil {
			log.Printf("[ERROR] %s response parse error: %v", action, err)
			return clusters
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprint(item["DBClusterDescription"])
			if protectedSweepVpcName(name) {
				log.Printf("[INFO] Skipping ADB cluster (do-not-delete marker): %s (%s)", name, fmt.Sprint(item["DBClusterId"]))
				continue
			}
			clusters = append(clusters, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return clusters
}

// adbSubSweepNameMatch reports whether a sub-resource name should be swept.
// It first rejects do-not-delete markers, then (outside SWEEP_ALL mode)
// requires a test-prefix match.
func adbSubSweepNameMatch(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	for _, prefix := range adbSweepPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

func init() {
	resource.AddTestSweepers("alicloud_adb_account", &resource.Sweeper{
		Name: "alicloud_adb_account",
		F:    testSweepAdbAccounts,
	})
	resource.AddTestSweepers("alicloud_adb_resource_group", &resource.Sweeper{
		Name: "alicloud_adb_resource_group",
		F:    testSweepAdbResourceGroups,
	})
	resource.AddTestSweepers("alicloud_adb_lake_account", &resource.Sweeper{
		Name: "alicloud_adb_lake_account",
		F:    testSweepAdbLakeAccounts,
	})
	resource.AddTestSweepers("alicloud_adb_connection", &resource.Sweeper{
		Name: "alicloud_adb_connection",
		F:    testSweepAdbConnections,
	})
}

// testSweepAdbAccounts sweeps orphaned alicloud_adb_account resources. It
// enumerates every ADB cluster, lists accounts (DescribeAccounts, adb
// 2019-03-15) on each, and deletes those matching the test prefixes (or all
// non-protected when SWEEP_ALL is set) via DeleteAccount.
func testSweepAdbAccounts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters := listAdbClustersForSweep(client)
	for _, cluster := range clusters {
		clusterId := fmt.Sprint(cluster["DBClusterId"])
		clusterName := fmt.Sprint(cluster["DBClusterDescription"])

		action := "DescribeAccounts"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"DBClusterId": clusterId,
		}
		response, err := client.RpcPost("adb", "2019-03-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				log.Printf("[ERROR] %s on cluster %s got a retryable error: %v", action, clusterId, err)
			} else {
				log.Printf("[ERROR] %s on cluster %s got an error: %v", action, clusterId, err)
			}
			continue
		}
		resp, err := jsonpath.Get("$.AccountList.DBAccount", response)
		if err != nil {
			log.Printf("[ERROR] %s response parse error on cluster %s: %v", action, clusterId, err)
			continue
		}
		accounts, _ := resp.([]interface{})
		for _, v := range accounts {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			accountName := fmt.Sprint(item["AccountName"])
			if !adbSubSweepNameMatch(accountName) && !adbSubSweepNameMatch(clusterName) {
				log.Printf("[INFO] Skipping ADB Account: %s on cluster %s", accountName, clusterId)
				continue
			}
			log.Printf("[INFO] Deleting ADB Account: %s on cluster %s", accountName, clusterId)
			deleteAction := "DeleteAccount"
			deleteRequest := map[string]interface{}{
				"RegionId":   client.RegionId,
				"DBClusterId": clusterId,
				"AccountName": accountName,
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				_, err := client.RpcPost("adb", "2019-03-15", deleteAction, nil, deleteRequest, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidAccountName.NotFound"}) || NotFoundError(err) {
						return nil
					}
					return retry.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				log.Printf("[ERROR] Deleting ADB Account %s on cluster %s failed: %v", accountName, clusterId, err)
			}
		}
	}
	return nil
}

// testSweepAdbResourceGroups sweeps orphaned alicloud_adb_resource_group
// resources. It enumerates every ADB cluster, lists resource groups
// (DescribeDBResourceGroup, adb 2019-03-15) on each, and deletes those
// matching the test prefixes (or all non-protected when SWEEP_ALL is set) via
// DeleteDBResourceGroup. The built-in USER_DEFAULT group is always spared.
func testSweepAdbResourceGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters := listAdbClustersForSweep(client)
	for _, cluster := range clusters {
		clusterId := fmt.Sprint(cluster["DBClusterId"])
		clusterName := fmt.Sprint(cluster["DBClusterDescription"])

		action := "DescribeDBResourceGroup"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"DBClusterId": clusterId,
		}
		response, err := client.RpcPost("adb", "2019-03-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				log.Printf("[ERROR] %s on cluster %s got a retryable error: %v", action, clusterId, err)
			} else {
				log.Printf("[ERROR] %s on cluster %s got an error: %v", action, clusterId, err)
			}
			continue
		}
		resp, err := jsonpath.Get("$.GroupsInfo", response)
		if err != nil {
			log.Printf("[ERROR] %s response parse error on cluster %s: %v", action, clusterId, err)
			continue
		}
		groups, _ := resp.([]interface{})
		for _, v := range groups {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			groupName := fmt.Sprint(item["GroupName"])
			// The default resource group USER_DEFAULT cannot be removed; skip it
			// to avoid noisy errors and to preserve cluster usability.
			if strings.EqualFold(groupName, "USER_DEFAULT") {
				log.Printf("[INFO] Skipping default ADB Resource Group: %s on cluster %s", groupName, clusterId)
				continue
			}
			if !adbSubSweepNameMatch(groupName) && !adbSubSweepNameMatch(clusterName) {
				log.Printf("[INFO] Skipping ADB Resource Group: %s on cluster %s", groupName, clusterId)
				continue
			}
			log.Printf("[INFO] Deleting ADB Resource Group: %s on cluster %s", groupName, clusterId)
			deleteAction := "DeleteDBResourceGroup"
			deleteRequest := map[string]interface{}{
				"RegionId":   client.RegionId,
				"DBClusterId": clusterId,
				"GroupName":   groupName,
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				_, err := client.RpcPost("adb", "2019-03-15", deleteAction, nil, deleteRequest, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) || NotFoundError(err) {
						return nil
					}
					return retry.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				log.Printf("[ERROR] Deleting ADB Resource Group %s on cluster %s failed: %v", groupName, clusterId, err)
			}
		}
	}
	return nil
}

// testSweepAdbLakeAccounts sweeps orphaned alicloud_adb_lake_account
// resources. Lake accounts live on ADB lake-version clusters and are managed
// through the adb 2021-12-01 API family. The sweeper enumerates every ADB
// cluster, lists accounts (DescribeAccounts, adb 2021-12-01) on each, and
// deletes those matching the test prefixes (or all non-protected when
// SWEEP_ALL is set) via DeleteAccount.
func testSweepAdbLakeAccounts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters := listAdbClustersForSweep(client)
	for _, cluster := range clusters {
		clusterId := fmt.Sprint(cluster["DBClusterId"])
		clusterName := fmt.Sprint(cluster["DBClusterDescription"])

		action := "DescribeAccounts"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"DBClusterId": clusterId,
		}
		response, err := client.RpcPost("adb", "2021-12-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				log.Printf("[ERROR] %s on cluster %s got a retryable error: %v", action, clusterId, err)
			} else {
				log.Printf("[ERROR] %s on cluster %s got an error: %v", action, clusterId, err)
			}
			continue
		}
		resp, err := jsonpath.Get("$.AccountList.DBAccount", response)
		if err != nil {
			log.Printf("[ERROR] %s response parse error on cluster %s: %v", action, clusterId, err)
			continue
		}
		accounts, _ := resp.([]interface{})
		for _, v := range accounts {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			accountName := fmt.Sprint(item["AccountName"])
			if !adbSubSweepNameMatch(accountName) && !adbSubSweepNameMatch(clusterName) {
				log.Printf("[INFO] Skipping ADB Lake Account: %s on cluster %s", accountName, clusterId)
				continue
			}
			log.Printf("[INFO] Deleting ADB Lake Account: %s on cluster %s", accountName, clusterId)
			deleteAction := "DeleteAccount"
			deleteRequest := map[string]interface{}{
				"RegionId":   client.RegionId,
				"DBClusterId": clusterId,
				"AccountName": accountName,
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				_, err := client.RpcPost("adb", "2021-12-01", deleteAction, nil, deleteRequest, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					if NotFoundError(err) {
						return nil
					}
					return retry.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				log.Printf("[ERROR] Deleting ADB Lake Account %s on cluster %s failed: %v", accountName, clusterId, err)
			}
		}
	}
	return nil
}

// testSweepAdbConnections sweeps orphaned alicloud_adb_connection resources.
// A public connection is a cluster-level attribute allocated via
// AllocateClusterPublicConnection and released via
// ReleaseClusterPublicConnection. The sweeper enumerates every ADB cluster,
// inspects its network addresses (DescribeDBClusterNetInfo) for a Public
// endpoint, and releases the public connection on clusters matching the test
// prefixes (or all non-protected when SWEEP_ALL is set).
func testSweepAdbConnections(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters := listAdbClustersForSweep(client)
	for _, cluster := range clusters {
		clusterId := fmt.Sprint(cluster["DBClusterId"])
		clusterName := fmt.Sprint(cluster["DBClusterDescription"])

		if !adbSubSweepNameMatch(clusterName) {
			log.Printf("[INFO] Skipping ADB Connection on cluster %s (%s): name does not match sweep filters", clusterId, clusterName)
			continue
		}

		action := "DescribeDBClusterNetInfo"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"DBClusterId": clusterId,
		}
		response, err := client.RpcPost("adb", "2019-03-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				log.Printf("[ERROR] %s on cluster %s got a retryable error: %v", action, clusterId, err)
			} else {
				log.Printf("[ERROR] %s on cluster %s got an error: %v", action, clusterId, err)
			}
			continue
		}
		resp, err := jsonpath.Get("$.Items.Address", response)
		if err != nil {
			log.Printf("[ERROR] %s response parse error on cluster %s: %v", action, clusterId, err)
			continue
		}
		addresses, _ := resp.([]interface{})
		hasPublic := false
		for _, v := range addresses {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if fmt.Sprint(item["NetType"]) == "Public" {
				hasPublic = true
				break
			}
		}
		if !hasPublic {
			continue
		}
		log.Printf("[INFO] Releasing ADB public connection on cluster %s (%s)", clusterId, clusterName)
		deleteAction := "ReleaseClusterPublicConnection"
		deleteRequest := map[string]interface{}{
			"RegionId":   client.RegionId,
			"DBClusterId": clusterId,
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
			_, err := client.RpcPost("adb", "2019-03-15", deleteAction, nil, deleteRequest, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) || NotFoundError(err) {
					return nil
				}
				return retry.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			log.Printf("[ERROR] Releasing ADB public connection on cluster %s failed: %v", clusterId, err)
		}
	}
	return nil
}
