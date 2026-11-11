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

// This file 补齐 (fills) sweeper gaps for bastionhost sub-resources. The
// alicloud_bastionhost_instance resource itself has NO sweeper because its
// Delete is intentionally a no-op: bastionhost instances are Subscription
// products purchased via BssOpenApi that cannot be released through the
// Yundun-bastionhost Delete API (the user must unsubscribe in the console or
// wait for the subscription to expire). The sub-resources below (Host,
// HostAccount, HostGroup, HostShareKey, User, UserGroup) all live under a
// bastionhost instance, have real Delete APIs, and previously had no sweepers,
// which left orphaned test resources after acceptance runs.

func init() {
	resource.AddTestSweepers("alicloud_bastionhost_host", &resource.Sweeper{
		Name: "alicloud_bastionhost_host",
		F:    testSweepBastionhostHosts,
	})
	resource.AddTestSweepers("alicloud_bastionhost_host_account", &resource.Sweeper{
		Name: "alicloud_bastionhost_host_account",
		F:    testSweepBastionhostHostAccounts,
	})
	resource.AddTestSweepers("alicloud_bastionhost_host_group", &resource.Sweeper{
		Name: "alicloud_bastionhost_host_group",
		F:    testSweepBastionhostHostGroups,
	})
	resource.AddTestSweepers("alicloud_bastionhost_host_share_key", &resource.Sweeper{
		Name: "alicloud_bastionhost_host_share_key",
		F:    testSweepBastionhostHostShareKeys,
	})
	resource.AddTestSweepers("alicloud_bastionhost_user", &resource.Sweeper{
		Name: "alicloud_bastionhost_user",
		F:    testSweepBastionhostUsers,
	})
	resource.AddTestSweepers("alicloud_bastionhost_user_group", &resource.Sweeper{
		Name: "alicloud_bastionhost_user_group",
		F:    testSweepBastionhostUserGroups,
	})
}

var bastionhostSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf_test_",
	"tf-test-",
	"testAcc",
	"tf-example",
	"tf-test",
}

// bastionhostSubResource describes a single bastionhost sub-resource so a
// shared helper can sweep all six sub-resource types without duplicated logic.
type bastionhostSubResource struct {
	sweeperName  string // e.g. "alicloud_bastionhost_host"
	listAction   string // e.g. "ListHosts"
	listPath     string // e.g. "$.Hosts"
	nameField    string // e.g. "HostName"
	idField      string // e.g. "HostId"
	deleteAction string // e.g. "DeleteHost"
	deleteParam  string // e.g. "HostId"
}

func testSweepBastionhostHosts(region string) error {
	return sweepBastionhostSubResources(region, bastionhostSubResource{
		sweeperName:  "alicloud_bastionhost_host",
		listAction:   "ListHosts",
		listPath:     "$.Hosts",
		nameField:    "HostName",
		idField:      "HostId",
		deleteAction: "DeleteHost",
		deleteParam:  "HostId",
	})
}

func testSweepBastionhostHostAccounts(region string) error {
	return sweepBastionhostSubResources(region, bastionhostSubResource{
		sweeperName:  "alicloud_bastionhost_host_account",
		listAction:   "ListHostAccounts",
		listPath:     "$.HostAccounts",
		nameField:    "HostAccountName",
		idField:      "HostAccountId",
		deleteAction: "DeleteHostAccount",
		deleteParam:  "HostAccountId",
	})
}

func testSweepBastionhostHostGroups(region string) error {
	return sweepBastionhostSubResources(region, bastionhostSubResource{
		sweeperName:  "alicloud_bastionhost_host_group",
		listAction:   "ListHostGroups",
		listPath:     "$.HostGroups",
		nameField:    "HostGroupName",
		idField:      "HostGroupId",
		deleteAction: "DeleteHostGroup",
		deleteParam:  "HostGroupId",
	})
}

func testSweepBastionhostHostShareKeys(region string) error {
	return sweepBastionhostSubResources(region, bastionhostSubResource{
		sweeperName:  "alicloud_bastionhost_host_share_key",
		listAction:   "ListHostShareKeys",
		listPath:     "$.HostShareKeys",
		nameField:    "HostShareKeyName",
		idField:      "HostShareKeyId",
		deleteAction: "DeleteHostShareKey",
		deleteParam:  "HostShareKeyId",
	})
}

func testSweepBastionhostUsers(region string) error {
	return sweepBastionhostSubResources(region, bastionhostSubResource{
		sweeperName:  "alicloud_bastionhost_user",
		listAction:   "ListUsers",
		listPath:     "$.Users",
		nameField:    "UserName",
		idField:      "UserId",
		deleteAction: "DeleteUser",
		deleteParam:  "UserId",
	})
}

func testSweepBastionhostUserGroups(region string) error {
	return sweepBastionhostSubResources(region, bastionhostSubResource{
		sweeperName:  "alicloud_bastionhost_user_group",
		listAction:   "ListUserGroups",
		listPath:     "$.UserGroups",
		nameField:    "UserGroupName",
		idField:      "UserGroupId",
		deleteAction: "DeleteUserGroup",
		deleteParam:  "UserGroupId",
	})
}

// sweepBastionhostSubResources lists every bastionhost instance in the region
// and, for each one, lists and deletes the named sub-resource type. The
// NODELETING guard (protectedSweepVpcName) runs at BOTH the parent-instance
// level and the per-sub-resource level, and in BOTH sweepAll() and normal mode
// — this matches the canonical VPC sweeper pattern and ensures do-not-delete
// markers are honored even when SWEEP_ALL=1 bypasses the prefix filter.
func sweepBastionhostSubResources(region string, r bastionhostSubResource) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instances, err := listBastionhostInstancesForSweep(client)
	if err != nil {
		return err
	}

	for _, inst := range instances {
		instanceID := fmt.Sprint(inst["InstanceId"])
		instanceDesc := fmt.Sprint(inst["Description"])

		// NODELETING guard at the parent-instance level — runs in BOTH modes.
		if protectedSweepVpcName(instanceDesc) {
			log.Printf("[INFO] Skipping bastionhost instance %s (%s): NODELETING marker on parent instance; %s under it protected", instanceID, instanceDesc, r.sweeperName)
			continue
		}
		if !sweepAll() {
			if !bastionhostNameMatchesPrefix(instanceDesc) {
				log.Printf("[INFO] Skipping bastionhost instance %s (%s): parent description does not match test prefixes; %s under it skipped", instanceID, instanceDesc, r.sweeperName)
				continue
			}
		}

		if err := sweepBastionhostSubResourcesForInstance(client, instanceID, r); err != nil {
			// Non-fatal: keep sweeping other instances.
			log.Printf("[ERROR] %s sweep for instance %s ended with: %s", r.sweeperName, instanceID, err)
		}
	}
	return nil
}

func sweepBastionhostSubResourcesForInstance(client *connectivity.AliyunClient, instanceID string, r bastionhostSubResource) error {
	request := map[string]interface{}{
		"InstanceId": instanceID,
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := client.RpcPost("Yundun-bastionhost", "2019-12-09", r.listAction, nil, request, true)
		if err != nil {
			return fmt.Errorf("error listing %s for instance %s: %s", r.listAction, instanceID, err)
		}
		resp, err := jsonpath.Get(r.listPath, response)
		if err != nil {
			return fmt.Errorf("error parsing %s response for instance %s: %s", r.listAction, instanceID, err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprint(item[r.nameField])
			id := fmt.Sprint(item[r.idField])

			// NODELETING guard at the per-sub-resource level — runs in BOTH modes.
			if protectedSweepVpcName(name) {
				log.Printf("[INFO] Skipping %s %s (%s) under instance %s: NODELETING marker", r.sweeperName, id, name, instanceID)
				continue
			}
			if !sweepAll() {
				if !bastionhostNameMatchesPrefix(name) {
					log.Printf("[INFO] Skipping %s %s (%s) under instance %s: name does not match test prefixes", r.sweeperName, id, name, instanceID)
					continue
				}
			}

			log.Printf("[INFO] Deleting bastionhost %s: %s (%s) under instance %s", r.sweeperName, id, name, instanceID)
			deleteReq := map[string]interface{}{
				r.deleteParam: id,
				"InstanceId":  instanceID,
				"RegionId":    client.RegionId,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(5*time.Minute, func() *retry.RetryError {
				_, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", r.deleteAction, nil, deleteReq, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					// Already-deleted or transitional states are non-fatal.
					if IsExpectedErrors(err, []string{"OBJECT_NOT_FOUND", "Commodity.BizError.InvalidStatus"}) {
						return nil
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete bastionhost %s %s (%s) under instance %s: %s", r.sweeperName, id, name, instanceID, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func listBastionhostInstancesForSweep(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	var instances []map[string]interface{}
	action := "DescribeInstances"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
		if err != nil {
			return nil, fmt.Errorf("error listing bastionhost instances: %s", err)
		}
		resp, err := jsonpath.Get("$.Instances", response)
		if err != nil {
			return nil, fmt.Errorf("error parsing DescribeInstances response: %s", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				instances = append(instances, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return instances, nil
}

func bastionhostNameMatchesPrefix(name string) bool {
	for _, prefix := range bastionhostSweepPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}
