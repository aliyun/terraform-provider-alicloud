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

// ENS (边缘计算节点服务 / Edge Node Service, product "Ens" / "2017-11-10",
// RPC style) sweeper gap fill. None of the alicloud_ens_* resources had a
// sweeper registered (verified via grep for AddTestSweepers / testSweepEns
// across alicloud/*_ens_*_test.go); this file fills all 18 gaps.
//
// All 18 ENS resources have both a list (Describe*) API and a delete API, and
// none is a BssOpenApi RefundInstance (退订) class — the instance release path
// uses the ENS-native ReleasePrePaidInstance / ReleasePostPaidInstance actions
// (mirroring the resource Delete), so instance is swept, not deferred.
//
// Sweep strategies:
//   - region-scoped (walk a single Describe* list, delete each): disk, eip,
//     image, instance, key_pair, load_balancer, nat_gateway, network,
//     network_route_table, security_group, snapshot, vswitch.
//   - parent-scoped (walk a parent list, enumerate children, delete each):
//     disk_instance_attachment (walk disks), eip_instance_attachment (walk
//     eips), instance_security_group_attachment (walk instances → SG list),
//     load_balancer_http_listener / load_balancer_udp_listener (walk LBs →
//     DescribeLoadBalancerListeners), bucket_lifecycle (walk buckets →
//     GetBucketLifecycle). These resources carry no name of their own, so the
//     prefix/NODELETING gate is applied to the PARENT's name (mirrors the kms
//     client_key pattern that gates on its parent AAP).
//
// Notable conventions:
//   - DescribeLoadBalancers / DescribeVSwitches / DescribeNatGateways are
//     GET-only (or GET-preferred) list actions the provider does not call; they
//     are reached via RpcGet. DescribeNatGateways returns a direct array
//     ($.NatGateways[*]); DescribeLoadBalancers / DescribeVSwitches nest one
//     level ($.LoadBalancers.LoadBalancer[*] / $.VSwitches.VSwitch[*]).
//   - DescribeSecurityGroups paginates with MaxResults + NextToken; every other
//     list uses PageNumber / PageSize.
//   - DeleteEip's API parameter is named "InstanceId" but holds the EIP
//     AllocationId (the EIP's own id) — the resource Delete passes
//     request["InstanceId"] = d.Id() (= AllocationId); the sweeper mirrors this.
//   - DeleteKeyPairs accepts KeyPairName (and KeyPairId) only; the "Version"
//     field in the terraform schema is a synthetic id segment (the list API
//     neither takes nor returns Version), so the sweeper deletes by KeyPairName.
//   - Instance DeletionProtection is a cloud-product managed protection: per
//     policy it is NEVER force-cleared. Instances with DeletionProtection on
//     are ALWAYS-ACTIVE skipped (not bypassed by SWEEP_ALL).
//
// Fleet-wipe pattern: SWEEP_ALL (ALICLOUD_SWEEP_ALL_RESOURCES=true) bypasses the
// tf-prefix filter so every non-NODELETING, non-protected resource is swept;
// NORMAL mode keeps the prefix filter. The protectedSweepVpcName NODELETING
// guard is always-active in both modes. All sweepers return nil — errors are
// logged at [ERROR] and never propagated, so one failing item does not abort
// the cascade.

func init() {
	resource.AddTestSweepers("alicloud_ens_bucket_lifecycle", &resource.Sweeper{
		Name: "alicloud_ens_bucket_lifecycle",
		F:    testSweepEnsBucketLifecycles,
	})
	resource.AddTestSweepers("alicloud_ens_disk_instance_attachment", &resource.Sweeper{
		Name: "alicloud_ens_disk_instance_attachment",
		F:    testSweepEnsDiskInstanceAttachments,
	})
	resource.AddTestSweepers("alicloud_ens_eip_instance_attachment", &resource.Sweeper{
		Name: "alicloud_ens_eip_instance_attachment",
		F:    testSweepEnsEipInstanceAttachments,
	})
	resource.AddTestSweepers("alicloud_ens_instance_security_group_attachment", &resource.Sweeper{
		Name: "alicloud_ens_instance_security_group_attachment",
		F:    testSweepEnsInstanceSecurityGroupAttachments,
	})
	resource.AddTestSweepers("alicloud_ens_load_balancer_http_listener", &resource.Sweeper{
		Name: "alicloud_ens_load_balancer_http_listener",
		F:    testSweepEnsLoadBalancerHttpListeners,
	})
	resource.AddTestSweepers("alicloud_ens_load_balancer_udp_listener", &resource.Sweeper{
		Name: "alicloud_ens_load_balancer_udp_listener",
		F:    testSweepEnsLoadBalancerUdpListeners,
	})
	resource.AddTestSweepers("alicloud_ens_snapshot", &resource.Sweeper{
		Name: "alicloud_ens_snapshot",
		F:    testSweepEnsSnapshots,
	})
	resource.AddTestSweepers("alicloud_ens_nat_gateway", &resource.Sweeper{
		Name: "alicloud_ens_nat_gateway",
		F:    testSweepEnsNatGateways,
	})
	resource.AddTestSweepers("alicloud_ens_network_route_table", &resource.Sweeper{
		Name: "alicloud_ens_network_route_table",
		F:    testSweepEnsNetworkRouteTables,
	})
	resource.AddTestSweepers("alicloud_ens_key_pair", &resource.Sweeper{
		Name: "alicloud_ens_key_pair",
		F:    testSweepEnsKeyPairs,
		// Key pairs bound to a running instance cannot be deleted
		// (PubKeyAttachedWithInstance); sweeping instances first releases
		// them and frees the key binding.
		Dependencies: []string{"alicloud_ens_instance"},
	})
	resource.AddTestSweepers("alicloud_ens_image", &resource.Sweeper{
		Name: "alicloud_ens_image",
		F:    testSweepEnsImages,
		// Images referenced by running instances are typically not deletable;
		// sweep instances first.
		Dependencies: []string{"alicloud_ens_instance"},
	})
	resource.AddTestSweepers("alicloud_ens_disk", &resource.Sweeper{
		Name: "alicloud_ens_disk",
		F:    testSweepEnsDisks,
		// An In-use disk cannot be deleted; detach it (and release its
		// instance) first.
		Dependencies: []string{
			"alicloud_ens_disk_instance_attachment",
			"alicloud_ens_instance",
		},
	})
	resource.AddTestSweepers("alicloud_ens_eip", &resource.Sweeper{
		Name: "alicloud_ens_eip",
		F:    testSweepEnsEips,
		// An associated EIP cannot be deleted; unassociate it first.
		Dependencies: []string{
			"alicloud_ens_eip_instance_attachment",
			"alicloud_ens_instance",
		},
	})
	resource.AddTestSweepers("alicloud_ens_security_group", &resource.Sweeper{
		Name: "alicloud_ens_security_group",
		F:    testSweepEnsSecurityGroups,
		// A security group with instances attached cannot be deleted; detach
		// (and release instances) first.
		Dependencies: []string{
			"alicloud_ens_instance_security_group_attachment",
			"alicloud_ens_instance",
		},
	})
	resource.AddTestSweepers("alicloud_ens_load_balancer", &resource.Sweeper{
		Name: "alicloud_ens_load_balancer",
		F:    testSweepEnsLoadBalancers,
		// A load balancer with active listeners cannot be deleted; sweep the
		// listeners first.
		Dependencies: []string{
			"alicloud_ens_load_balancer_http_listener",
			"alicloud_ens_load_balancer_udp_listener",
		},
	})
	resource.AddTestSweepers("alicloud_ens_vswitch", &resource.Sweeper{
		Name: "alicloud_ens_vswitch",
		F:    testSweepEnsVswitches,
		// A vswitch with child resources (LB / instance / nat gateway) cannot
		// be deleted; sweep them first.
		Dependencies: []string{
			"alicloud_ens_load_balancer",
			"alicloud_ens_instance",
			"alicloud_ens_nat_gateway",
		},
	})
	resource.AddTestSweepers("alicloud_ens_network", &resource.Sweeper{
		Name: "alicloud_ens_network",
		F:    testSweepEnsNetworks,
		// A network with child resources cannot be deleted; sweep the children
		// first.
		Dependencies: []string{
			"alicloud_ens_vswitch",
			"alicloud_ens_network_route_table",
			"alicloud_ens_nat_gateway",
			"alicloud_ens_load_balancer",
			"alicloud_ens_security_group",
			"alicloud_ens_instance",
		},
	})
	resource.AddTestSweepers("alicloud_ens_instance", &resource.Sweeper{
		Name: "alicloud_ens_instance",
		F:    testSweepEnsInstances,
	})
}

var ensSweeperPrefixes = []string{
	"tf-testacc",
	"tf_testacc",
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// ensGetField returns the first non-empty string value of item[name] across the
// given name variants. ENS response casing is mostly PascalCase but the variant
// list tolerates camelCase fallbacks.
func ensGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// ensSweepNameOk applies the always-active NODELETING guard and the
// sweepAll()-aware tf-prefix filter. Returns true when the resource should be
// swept. SWEEP_ALL bypasses the prefix filter (fleet-wipe); NORMAL mode keeps
// it.
func ensSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping ens resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range ensSweeperPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping ens resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// ensRpcPost is the shared ENS RpcPost list/read wrapper: RpcPost with
// autoRetry=false wrapped in retry.Retry + NeedRetry + incrementalWait so
// throttling and 5xx are retried at the sweeper level (the SDK does not
// auto-retry here). List params are carried in the request body (query is nil),
// matching the ENS datasource convention.
func ensRpcPost(client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, nil, request, false)
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

// ensRpcGet is the RpcGet list/read wrapper for the GET-only / GET-preferred
// ENS list actions the provider does not otherwise call (DescribeLoadBalancers,
// DescribeVSwitches, DescribeNatGateways). List params are carried in the query
// string, matching the provider's DescribeEnsNatGateway convention.
func ensRpcGet(client *connectivity.AliyunClient, action string, query map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Ens", "2017-11-10", action, query, nil)
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

// ensListPagedPost paginates an RpcPost list action with PageNumber/PageSize
// (PageSizeLarge=50, within every ENS list cap — DescribeDisks caps at 50, the
// rest at 100) and returns every item from jsonPath across all pages.
func ensListPagedPost(client *connectivity.AliyunClient, action, jsonPath string) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := ensRpcPost(client, action, request)
		if err != nil {
			return nil, err
		}
		raw, err := jsonpath.Get(jsonPath, response)
		if err != nil {
			return nil, fmt.Errorf("%s jsonpath %s: %s", action, jsonPath, err)
		}
		result, _ := raw.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return items, nil
}

// ensListPagedGet is the RpcGet counterpart of ensListPagedPost (params in query).
func ensListPagedGet(client *connectivity.AliyunClient, action, jsonPath string) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	query := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := ensRpcGet(client, action, query)
		if err != nil {
			return nil, err
		}
		raw, err := jsonpath.Get(jsonPath, response)
		if err != nil {
			return nil, fmt.Errorf("%s jsonpath %s: %s", action, jsonPath, err)
		}
		result, _ := raw.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		query["PageNumber"] = query["PageNumber"].(int) + 1
	}
	return items, nil
}

// ensListTokenPost paginates an RpcPost list action with MaxResults + NextToken
// (DescribeSecurityGroups) and returns every item from jsonPath.
func ensListTokenPost(client *connectivity.AliyunClient, action, jsonPath string) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	request := map[string]interface{}{"MaxResults": PageSizeLarge}
	for {
		response, err := ensRpcPost(client, action, request)
		if err != nil {
			return nil, err
		}
		raw, err := jsonpath.Get(jsonPath, response)
		if err != nil {
			return nil, fmt.Errorf("%s jsonpath %s: %s", action, jsonPath, err)
		}
		result, _ := raw.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return items, nil
}

// ensDeletePost issues a single best-effort ENS RpcPost delete (autoRetry=false,
// no wrapper retry): not-found / expected errors are swallowed, everything else
// is logged and swallowed so a single failure never aborts the cascade. params
// is placed in whichever of query/request the resource's own Delete uses (the
// caller passes the populated map and nil for the other); nil is normalized to
// an empty map so the underlying SDK never sees a nil body.
func ensDeletePost(client *connectivity.AliyunClient, action, label string, query, request map[string]interface{}, expectedCodes ...string) {
	if query == nil {
		query = make(map[string]interface{})
	}
	if request == nil {
		request = make(map[string]interface{})
	}
	if _, err := client.RpcPost("Ens", "2017-11-10", action, query, request, false); err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, expectedCodes) {
			return
		}
		log.Printf("[ERROR] Failed to delete ens %s (%s): %s", label, action, err)
	}
}

// ensDeletionProtectionEnabled reports whether an instance's DeletionProtection
// field is on. The DescribeInstances response surfaces it as a bool, but the
// JSON-decoded form may also be a string; both are handled. Per policy this is
// NEVER force-cleared — protected instances are always-active skipped (not
// bypassed by SWEEP_ALL), mirroring the fcv3 domain-bearing always-active skip.
func ensDeletionProtectionEnabled(v interface{}) bool {
	switch dp := v.(type) {
	case bool:
		return dp
	case string:
		return strings.EqualFold(dp, "true") || strings.EqualFold(dp, "Enabled") || dp == "1"
	}
	return false
}

// ensStopListener best-effort stops a LoadBalancer listener so a subsequent
// DeleteLoadBalancerListener / DeleteLoadBalancerHTTPListener can succeed (a
// running listener cannot be deleted). IncorrectListenerStatus (already
// stopped) and any other error is swallowed.
func ensStopListener(client *connectivity.AliyunClient, loadBalancerId, listenerPort string) {
	request := map[string]interface{}{
		"LoadBalancerId": loadBalancerId,
		"ListenerPort":   listenerPort,
	}
	if _, err := client.RpcPost("Ens", "2017-11-10", "StopLoadBalancerListener", nil, request, false); err != nil {
		log.Printf("[DEBUG] ens stop listener (%s:%s): %s (ignored)", loadBalancerId, listenerPort, err)
	}
}

func testSweepEnsDisks(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedPost(client, "DescribeDisks", "$.Disks.Disks[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens disks: %s", err)
		return nil
	}
	for _, item := range items {
		diskId := ensGetField(item, "DiskId")
		if diskId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "DiskName")) {
			continue
		}
		log.Printf("[INFO] Deleting ens disk: %s", diskId)
		ensDeletePost(client, "DeleteDisk", "disk", nil, map[string]interface{}{"DiskId": diskId})
	}
	return nil
}

func testSweepEnsEips(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedPost(client, "DescribeEnsEipAddresses", "$.EipAddresses.EipAddress[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens eips: %s", err)
		return nil
	}
	for _, item := range items {
		allocationId := ensGetField(item, "AllocationId")
		if allocationId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "Name")) {
			continue
		}
		log.Printf("[INFO] Deleting ens eip: %s", allocationId)
		// DeleteEip's "InstanceId" param carries the EIP AllocationId.
		ensDeletePost(client, "DeleteEip", "eip", nil, map[string]interface{}{"InstanceId": allocationId})
	}
	return nil
}

func testSweepEnsImages(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedPost(client, "DescribeSelfImages", "$.Images.Image[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens images: %s", err)
		return nil
	}
	for _, item := range items {
		imageId := ensGetField(item, "ImageId")
		if imageId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "ImageName")) {
			continue
		}
		log.Printf("[INFO] Deleting ens image: %s", imageId)
		ensDeletePost(client, "DeleteImage", "image", nil, map[string]interface{}{"ImageId": imageId})
	}
	return nil
}

func testSweepEnsInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedPost(client, "DescribeInstances", "$.Instances.Instance[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens instances: %s", err)
		return nil
	}
	for _, item := range items {
		instanceId := ensGetField(item, "InstanceId")
		if instanceId == "" {
			continue
		}
		instanceName := ensGetField(item, "InstanceName")
		if !ensSweepNameOk(instanceName) {
			continue
		}
		// DeletionProtection is a cloud-product managed protection: never
		// force-cleared. Always-active skip (not bypassed by SWEEP_ALL).
		if ensDeletionProtectionEnabled(item["DeletionProtection"]) {
			log.Printf("[INFO] Skipping ens instance %s (%s): DeletionProtection enabled (always-active skip)", instanceId, instanceName)
			continue
		}
		// Charge type (InstanceResourceType): EnsInstance = Subscription/PrePaid
		// → ReleasePrePaidInstance; everything else (EnsService /
		// EnsPostPaidInstance / EnsPostInstance / BuildMachine) is
		// PayAsYouGo/PostPaid → ReleasePostPaidInstance. Both are ENS-native
		// release actions (NOT BssOpenApi RefundInstance / 退订).
		resourceType := ensGetField(item, "InstanceResourceType")
		var action string
		if resourceType == "EnsInstance" {
			action = "ReleasePrePaidInstance"
		} else {
			action = "ReleasePostPaidInstance"
		}
		log.Printf("[INFO] Releasing ens instance: %s (%s, %s)", instanceId, instanceName, action)
		ensDeletePost(client, action, "instance", map[string]interface{}{"InstanceId": instanceId}, nil)
	}
	return nil
}

func testSweepEnsKeyPairs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedPost(client, "DescribeKeyPairs", "$.KeyPairs.KeyPair")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens key_pairs: %s", err)
		return nil
	}
	for _, item := range items {
		keyPairName := ensGetField(item, "KeyPairName")
		if keyPairName == "" {
			continue
		}
		if !ensSweepNameOk(keyPairName) {
			continue
		}
		log.Printf("[INFO] Deleting ens key_pair: %s", keyPairName)
		// DeleteKeyPairs takes KeyPairName; Version is a terraform-only id
		// segment the API neither requires nor returns, so it is omitted.
		ensDeletePost(client, "DeleteKeyPairs", "key_pair", nil, map[string]interface{}{"KeyPairName": keyPairName},
			"PubKeyNameNotExist", "PubKeyAttachedWithInstance")
	}
	return nil
}

func testSweepEnsLoadBalancers(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedGet(client, "DescribeLoadBalancers", "$.LoadBalancers.LoadBalancer[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens load_balancers: %s", err)
		return nil
	}
	for _, item := range items {
		loadBalancerId := ensGetField(item, "LoadBalancerId")
		if loadBalancerId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "LoadBalancerName")) {
			continue
		}
		log.Printf("[INFO] Deleting ens load_balancer: %s", loadBalancerId)
		ensDeletePost(client, "DeleteLoadBalancer", "load_balancer", nil, map[string]interface{}{"LoadBalancerId": loadBalancerId})
	}
	return nil
}

func testSweepEnsNatGateways(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedGet(client, "DescribeNatGateways", "$.NatGateways[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens nat_gateways: %s", err)
		return nil
	}
	for _, item := range items {
		natGatewayId := ensGetField(item, "NatGatewayId")
		if natGatewayId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "Name")) {
			continue
		}
		log.Printf("[INFO] Deleting ens nat_gateway: %s", natGatewayId)
		ensDeletePost(client, "DeleteNatGateway", "nat_gateway", map[string]interface{}{"NatGatewayId": natGatewayId}, nil)
	}
	return nil
}

func testSweepEnsNetworks(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedPost(client, "DescribeNetworks", "$.Networks.Network[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens networks: %s", err)
		return nil
	}
	for _, item := range items {
		networkId := ensGetField(item, "NetworkId")
		if networkId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "NetworkName")) {
			continue
		}
		log.Printf("[INFO] Deleting ens network: %s", networkId)
		ensDeletePost(client, "DeleteNetwork", "network", map[string]interface{}{"NetworkId": networkId}, nil)
	}
	return nil
}

func testSweepEnsNetworkRouteTables(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedPost(client, "DescribeEnsRouteTables", "$.RouteTables[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens network_route_tables: %s", err)
		return nil
	}
	for _, item := range items {
		routeTableId := ensGetField(item, "RouteTableId")
		if routeTableId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "RouteTableName")) {
			continue
		}
		log.Printf("[INFO] Deleting ens network_route_table: %s", routeTableId)
		ensDeletePost(client, "DeleteEnsRouteTable", "network_route_table", map[string]interface{}{"RouteTableId": routeTableId}, nil)
	}
	return nil
}

func testSweepEnsSecurityGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListTokenPost(client, "DescribeSecurityGroups", "$.SecurityGroups.SecurityGroup[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens security_groups: %s", err)
		return nil
	}
	for _, item := range items {
		securityGroupId := ensGetField(item, "SecurityGroupId")
		if securityGroupId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "SecurityGroupName")) {
			continue
		}
		log.Printf("[INFO] Deleting ens security_group: %s", securityGroupId)
		ensDeletePost(client, "DeleteSecurityGroup", "security_group", nil, map[string]interface{}{"SecurityGroupId": securityGroupId})
	}
	return nil
}

func testSweepEnsSnapshots(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedPost(client, "DescribeSnapshots", "$.Snapshots[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens snapshots: %s", err)
		return nil
	}
	for _, item := range items {
		snapshotId := ensGetField(item, "SnapshotId")
		if snapshotId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "SnapshotName")) {
			continue
		}
		log.Printf("[INFO] Deleting ens snapshot: %s", snapshotId)
		ensDeletePost(client, "DeleteSnapshot", "snapshot", map[string]interface{}{"SnapshotId": snapshotId}, nil)
	}
	return nil
}

func testSweepEnsVswitches(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	items, err := ensListPagedGet(client, "DescribeVSwitches", "$.VSwitches.VSwitch[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens vswitches: %s", err)
		return nil
	}
	for _, item := range items {
		vswitchId := ensGetField(item, "VSwitchId")
		if vswitchId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "VSwitchName")) {
			continue
		}
		log.Printf("[INFO] Deleting ens vswitch: %s", vswitchId)
		ensDeletePost(client, "DeleteVSwitch", "vswitch", nil, map[string]interface{}{"VSwitchId": vswitchId})
	}
	return nil
}

// --- parent-scoped sweepers ---

func testSweepEnsDiskInstanceAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Walk disks; each In-use disk item carries the InstanceId it is attached
	// to. The gate is on the disk (parent) name.
	items, err := ensListPagedPost(client, "DescribeDisks", "$.Disks.Disks[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens disks for disk_instance_attachment sweep: %s", err)
		return nil
	}
	for _, item := range items {
		diskId := ensGetField(item, "DiskId")
		if diskId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "DiskName")) {
			continue
		}
		instanceId := ensGetField(item, "InstanceId")
		if instanceId == "" {
			continue // not attached
		}
		log.Printf("[INFO] Detaching ens disk_instance_attachment: %s:%s", diskId, instanceId)
		ensDeletePost(client, "DetachDisk", "disk_instance_attachment",
			map[string]interface{}{"DiskId": diskId, "InstanceId": instanceId}, nil)
	}
	return nil
}

func testSweepEnsEipInstanceAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Walk eips; an associated eip item carries a non-empty InstanceId. The gate
	// is on the eip (parent) Name.
	items, err := ensListPagedPost(client, "DescribeEnsEipAddresses", "$.EipAddresses.EipAddress[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens eips for eip_instance_attachment sweep: %s", err)
		return nil
	}
	for _, item := range items {
		allocationId := ensGetField(item, "AllocationId")
		if allocationId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "Name")) {
			continue
		}
		instanceId := ensGetField(item, "InstanceId")
		if instanceId == "" {
			continue // not associated
		}
		log.Printf("[INFO] Unassociating ens eip_instance_attachment: %s", allocationId)
		ensDeletePost(client, "UnAssociateEnsEipAddress", "eip_instance_attachment",
			map[string]interface{}{"AllocationId": allocationId}, nil)
	}
	return nil
}

func testSweepEnsInstanceSecurityGroupAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Walk instances; each instance item carries SecurityGroupIds.SecurityGroupId
	// (array). The gate is on the instance (parent) InstanceName.
	items, err := ensListPagedPost(client, "DescribeInstances", "$.Instances.Instance[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens instances for instance_security_group_attachment sweep: %s", err)
		return nil
	}
	for _, item := range items {
		instanceId := ensGetField(item, "InstanceId")
		if instanceId == "" {
			continue
		}
		if !ensSweepNameOk(ensGetField(item, "InstanceName")) {
			continue
		}
		raw, jerr := jsonpath.Get("$.SecurityGroupIds.SecurityGroupId", item)
		if jerr != nil {
			continue
		}
		sgs, _ := raw.([]interface{})
		for _, sg := range sgs {
			securityGroupId := fmt.Sprint(sg)
			if securityGroupId == "" || securityGroupId == "<nil>" {
				continue
			}
			log.Printf("[INFO] Leaving ens instance_security_group_attachment: %s:%s", instanceId, securityGroupId)
			ensDeletePost(client, "LeaveSecurityGroup", "instance_security_group_attachment",
				map[string]interface{}{"SecurityGroupId": securityGroupId, "InstanceId": instanceId}, nil)
		}
	}
	return nil
}

// listEnsLoadBalancerIds returns the LoadBalancerId of every load balancer whose
// LoadBalancerName passes ensSweepNameOk. Shared by the http/udp listener
// sweepers (each walks the parent LB list, enumerates its listeners, and deletes
// its own protocol's listeners).
func listEnsLoadBalancerIds(client *connectivity.AliyunClient) ([]string, error) {
	items, err := ensListPagedGet(client, "DescribeLoadBalancers", "$.LoadBalancers.LoadBalancer[*]")
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		if !ensSweepNameOk(ensGetField(item, "LoadBalancerName")) {
			continue
		}
		if id := ensGetField(item, "LoadBalancerId"); id != "" {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// listEnsListeners enumerates all listeners beneath a load balancer (across
// protocols) via DescribeLoadBalancerListeners, returning the LoadBalancerId,
// ListenerPort and Protocol of each.
func listEnsListeners(client *connectivity.AliyunClient, loadBalancerId string) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	request := map[string]interface{}{
		"LoadBalancerId": loadBalancerId,
		"PageSize":       PageSizeLarge,
		"PageNumber":     1,
	}
	for {
		response, err := ensRpcPost(client, "DescribeLoadBalancerListeners", request)
		if err != nil {
			return nil, err
		}
		raw, err := jsonpath.Get("$.Listeners.Listener[*]", response)
		if err != nil {
			return nil, fmt.Errorf("DescribeLoadBalancerListeners jsonpath $.Listeners.Listener[*]: %s", err)
		}
		result, _ := raw.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return items, nil
}

func testSweepEnsLoadBalancerHttpListeners(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	lbIds, err := listEnsLoadBalancerIds(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list ens load_balancers for http_listener sweep: %s", err)
		return nil
	}
	for _, lbId := range lbIds {
		listeners, err := listEnsListeners(client, lbId)
		if err != nil {
			log.Printf("[ERROR] Failed to list ens listeners for load_balancer %s: %s", lbId, err)
			continue
		}
		for _, listener := range listeners {
			protocol := ensGetField(listener, "Protocol")
			if !strings.EqualFold(protocol, "http") {
				continue
			}
			listenerPort := fmt.Sprint(listener["ListenerPort"])
			if listenerPort == "" || listenerPort == "<nil>" {
				continue
			}
			log.Printf("[INFO] Deleting ens load_balancer_http_listener: %s:%s", lbId, listenerPort)
			ensStopListener(client, lbId, listenerPort)
			ensDeletePost(client, "DeleteLoadBalancerHTTPListener", "load_balancer_http_listener",
				nil, map[string]interface{}{"LoadBalancerId": lbId, "ListenerPort": listenerPort},
				"IncorrectListenerStatus", "IncorrectInstanceStatus")
		}
	}
	return nil
}

func testSweepEnsLoadBalancerUdpListeners(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	lbIds, err := listEnsLoadBalancerIds(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list ens load_balancers for udp_listener sweep: %s", err)
		return nil
	}
	for _, lbId := range lbIds {
		listeners, err := listEnsListeners(client, lbId)
		if err != nil {
			log.Printf("[ERROR] Failed to list ens listeners for load_balancer %s: %s", lbId, err)
			continue
		}
		for _, listener := range listeners {
			protocol := ensGetField(listener, "Protocol")
			if !strings.EqualFold(protocol, "udp") {
				continue
			}
			listenerPort := fmt.Sprint(listener["ListenerPort"])
			if listenerPort == "" || listenerPort == "<nil>" {
				continue
			}
			log.Printf("[INFO] Deleting ens load_balancer_udp_listener: %s:%s", lbId, listenerPort)
			ensStopListener(client, lbId, listenerPort)
			// DeleteLoadBalancerListener is the generic listener delete (by
			// LoadBalancerId + ListenerPort); only invoked for UDP listeners
			// here, so HTTP listeners are left to the http_listener sweeper.
			ensDeletePost(client, "DeleteLoadBalancerListener", "load_balancer_udp_listener",
				nil, map[string]interface{}{"LoadBalancerId": lbId, "ListenerPort": listenerPort},
				"IncorrectListenerStatus", "IncorrectInstanceStatus")
		}
	}
	return nil
}

func testSweepEnsBucketLifecycles(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Walk buckets; the lifecycle rule carries no name of its own so the gate
	// is on the bucket (parent) BucketName. ListBuckets returns a direct
	// array ($.BucketInfos[*]).
	buckets, err := ensListPagedGet(client, "ListBuckets", "$.BucketInfos[*]")
	if err != nil {
		log.Printf("[ERROR] Failed to list ens buckets for bucket_lifecycle sweep: %s", err)
		return nil
	}
	for _, bucket := range buckets {
		bucketName := ensGetField(bucket, "BucketName")
		if bucketName == "" {
			continue
		}
		if !ensSweepNameOk(bucketName) {
			continue
		}
		// GetBucketLifecycle with only BucketName returns every lifecycle rule
		// on the bucket ($.Rule[*]); each rule's id is the "ID" field.
		request := map[string]interface{}{"BucketName": bucketName}
		response, err := ensRpcPost(client, "GetBucketLifecycle", request)
		if err != nil {
			log.Printf("[ERROR] Failed to list ens bucket_lifecycles for bucket %s: %s", bucketName, err)
			continue
		}
		raw, jerr := jsonpath.Get("$.Rule[*]", response)
		if jerr != nil {
			continue
		}
		rules, _ := raw.([]interface{})
		for _, r := range rules {
			rule, ok := r.(map[string]interface{})
			if !ok {
				continue
			}
			ruleId := fmt.Sprint(rule["ID"])
			if ruleId == "" || ruleId == "<nil>" {
				continue
			}
			log.Printf("[INFO] Deleting ens bucket_lifecycle: %s:%s", bucketName, ruleId)
			ensDeletePost(client, "DeleteBucketLifecycle", "bucket_lifecycle",
				map[string]interface{}{"BucketName": bucketName, "RuleId": ruleId}, nil,
				"NoSuchBucket", "NoSuchLifecycle")
		}
	}
	return nil
}
