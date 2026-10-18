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

// gaSweepPrefixes is the canonical Terraform test-name prefix list used by
// every alicloud_ga_* sweeper when SWEEP_ALL is not set. When SWEEP_ALL=1 the
// prefix filter is bypassed (sweepAll()==true), but the NODELETING name guard
// via protectedSweepVpcName still applies, so manual/production Ga resources
// are never touched even in SWEEP_ALL.
var gaSweepPrefixes = []string{"tf-testAcc", "tf_testAcc"}

// gaSweepSkipName decides whether a Ga resource should be skipped based on
// its name. It returns true (skip) when the name carries a NODELETING /
// do-not-delete marker (protectedSweepVpcName), or when SWEEP_ALL is off and
// the name does not match the canonical Terraform test prefixes.
func gaSweepSkipName(name string) bool {
	if protectedSweepVpcName(name) {
		return true
	}
	if sweepAll() {
		return false
	}
	for _, prefix := range gaSweepPrefixes {
		if strings.HasPrefix(name, prefix) {
			return false
		}
	}
	return true
}

// listGaAcceleratorIds enumerates every Global Accelerator instance in the
// caller's region. Several Ga child resources (listener, endpoint group, ip
// set, forwarding rule, custom-routing endpoint group) expose no global list
// API; they require an AcceleratorId, so the sweepers iterate accelerators
// first.
func listGaAcceleratorIds(client *connectivity.AliyunClient) ([]string, error) {
	var ids []string
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"PageSize":    PageSizeLarge,
		"PageNumber":  1,
	}
	for {
		action := "ListAccelerators"
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", action, err)
		}
		addDebug(action, response, request)
		resp, err := jsonpath.Get("$.Accelerators", response)
		if err != nil {
			return nil, err
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if acc, ok := v.(map[string]interface{}); ok {
				if id := fmt.Sprint(acc["AcceleratorId"]); id != "" && id != "<nil>" {
					ids = append(ids, id)
				}
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return ids, nil
}

// listGaListenerIds enumerates every listener attached to the given
// accelerator. Forwarding-rule sweepers iterate listeners per accelerator
// because ListForwardingRules requires both AcceleratorId and ListenerId.
func listGaListenerIds(client *connectivity.AliyunClient, acceleratorId string) ([]map[string]interface{}, error) {
	var listeners []map[string]interface{}
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"AcceleratorId": acceleratorId,
		"PageSize":      PageSizeLarge,
		"PageNumber":    1,
	}
	for {
		action := "ListListeners"
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", action, err)
		}
		addDebug(action, response, request)
		resp, err := jsonpath.Get("$.Listeners", response)
		if err != nil {
			return nil, err
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if l, ok := v.(map[string]interface{}); ok {
				listeners = append(listeners, l)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return listeners, nil
}

func init() {
	resource.AddTestSweepers("alicloud_ga_listener", &resource.Sweeper{
		Name: "alicloud_ga_listener",
		F:    testSweepGaListener,
	})
	resource.AddTestSweepers("alicloud_ga_endpoint_group", &resource.Sweeper{
		Name: "alicloud_ga_endpoint_group",
		F:    testSweepGaEndpointGroup,
	})
	resource.AddTestSweepers("alicloud_ga_ip_set", &resource.Sweeper{
		Name: "alicloud_ga_ip_set",
		F:    testSweepGaIpSet,
	})
	resource.AddTestSweepers("alicloud_ga_forwarding_rule", &resource.Sweeper{
		Name: "alicloud_ga_forwarding_rule",
		F:    testSweepGaForwardingRule,
	})
	resource.AddTestSweepers("alicloud_ga_domain", &resource.Sweeper{
		Name: "alicloud_ga_domain",
		F:    testSweepGaDomain,
	})
	resource.AddTestSweepers("alicloud_ga_basic_accelerator", &resource.Sweeper{
		Name: "alicloud_ga_basic_accelerator",
		F:    testSweepGaBasicAccelerator,
	})
	resource.AddTestSweepers("alicloud_ga_custom_routing_endpoint_group", &resource.Sweeper{
		Name: "alicloud_ga_custom_routing_endpoint_group",
		F:    testSweepGaCustomRoutingEndpointGroup,
	})
}

// testSweepGaListener deletes tf-testAcc* Ga listeners. ListListeners requires
// an AcceleratorId, so the sweeper enumerates accelerators first and then
// walks listeners per accelerator.
func testSweepGaListener(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	accIds, err := listGaAcceleratorIds(client)
	if err != nil {
		log.Printf("[ERROR] testSweepGaListener list accelerators failed: %v", err)
		return nil
	}

	for _, accId := range accIds {
		listeners, err := listGaListenerIds(client, accId)
		if err != nil {
			log.Printf("[ERROR] ListListeners for accelerator %s failed: %v", accId, err)
			continue
		}
		for _, listener := range listeners {
			listenerId := fmt.Sprint(listener["ListenerId"])
			listenerName := fmt.Sprint(listener["Name"])
			if gaSweepSkipName(listenerName) {
				log.Printf("[INFO] Skipping Ga listener: %s(%s)", listenerId, listenerName)
				continue
			}
			log.Printf("[Info] Delete Ga listener: %s(%s)", listenerId, listenerName)
			req := map[string]interface{}{
				"ClientToken":   buildClientToken("DeleteListener"),
				"AcceleratorId": accId,
				"ListenerId":    listenerId,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(1*time.Minute, func() *retry.RetryError {
				_, e := client.RpcPost("Ga", "2019-11-20", "DeleteListener", nil, req, true)
				if e != nil {
					if IsExpectedErrors(e, []string{"StateError.Accelerator", "NotActive.Listener", "Exist.ForwardingRule", "Exist.EndpointGroup"}) || NeedRetry(e) {
						wait()
						return retry.RetryableError(e)
					}
					return retry.NonRetryableError(e)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Deleting Ga listener %s got an error: %s", listenerId, err)
			}
		}
	}
	return nil
}

// testSweepGaEndpointGroup deletes tf-testAcc* Ga endpoint groups. The
// ListEndpointGroups API requires an AcceleratorId, so the sweeper iterates
// accelerators first.
func testSweepGaEndpointGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	accIds, err := listGaAcceleratorIds(client)
	if err != nil {
		log.Printf("[ERROR] testSweepGaEndpointGroup list accelerators failed: %v", err)
		return nil
	}

	for _, accId := range accIds {
		request := map[string]interface{}{
			"RegionId":      client.RegionId,
			"AcceleratorId": accId,
			"PageSize":      PageSizeLarge,
			"PageNumber":    1,
		}
		for {
			action := "ListEndpointGroups"
			response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] %s got an error: %v", action, err)
				break
			}
			addDebug(action, response, request)
			resp, err := jsonpath.Get("$.EndpointGroups", response)
			if err != nil {
				log.Println(err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				eg := v.(map[string]interface{})
				egId := fmt.Sprint(eg["EndpointGroupId"])
				egName := fmt.Sprint(eg["Name"])
				if gaSweepSkipName(egName) {
					log.Printf("[INFO] Skipping Ga endpoint group: %s(%s)", egId, egName)
					continue
				}
				log.Printf("[Info] Delete Ga endpoint group: %s(%s)", egId, egName)
				req := map[string]interface{}{
					"ClientToken":     buildClientToken("DeleteEndpointGroup"),
					"AcceleratorId":   accId,
					"EndpointGroupId": egId,
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(1*time.Minute, func() *retry.RetryError {
					_, e := client.RpcPost("Ga", "2019-11-20", "DeleteEndpointGroup", nil, req, true)
					if e != nil {
						if IsExpectedErrors(e, []string{"StateError.Accelerator", "StateError.EndPointGroup", "NotActive.Listener"}) || NeedRetry(e) {
							wait()
							return retry.RetryableError(e)
						}
						return retry.NonRetryableError(e)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] Deleting Ga endpoint group %s got an error: %s", egId, err)
				}
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return nil
}

// testSweepGaIpSet deletes tf-testAcc* Ga ip sets. ListIpSets requires an
// AcceleratorId, so the sweeper iterates accelerators first.
func testSweepGaIpSet(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	accIds, err := listGaAcceleratorIds(client)
	if err != nil {
		log.Printf("[ERROR] testSweepGaIpSet list accelerators failed: %v", err)
		return nil
	}

	for _, accId := range accIds {
		request := map[string]interface{}{
			"RegionId":      client.RegionId,
			"AcceleratorId": accId,
			"PageSize":      PageSizeLarge,
			"PageNumber":    1,
		}
		for {
			action := "ListIpSets"
			response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] %s got an error: %v", action, err)
				break
			}
			addDebug(action, response, request)
			resp, err := jsonpath.Get("$.IpSets", response)
			if err != nil {
				log.Println(err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				ipSet := v.(map[string]interface{})
				ipSetId := fmt.Sprint(ipSet["IpSetId"])
				ipSetName := fmt.Sprint(ipSet["Name"])
				if gaSweepSkipName(ipSetName) {
					log.Printf("[INFO] Skipping Ga ip set: %s(%s)", ipSetId, ipSetName)
					continue
				}
				log.Printf("[Info] Delete Ga ip set: %s(%s)", ipSetId, ipSetName)
				req := map[string]interface{}{
					"IpSetIds.1": ipSetId,
					"RegionId":   client.RegionId,
					"ClientToken": buildClientToken("DeleteIpSet"),
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(1*time.Minute, func() *retry.RetryError {
					_, e := client.RpcPost("Ga", "2019-11-20", "DeleteIpSets", nil, req, true)
					if e != nil {
						if IsExpectedErrors(e, []string{"StateError.Accelerator", "StateError.IpSet", "NotActive.Listener"}) || NeedRetry(e) {
							wait()
							return retry.RetryableError(e)
						}
						return retry.NonRetryableError(e)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] Deleting Ga ip set %s got an error: %s", ipSetId, err)
				}
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return nil
}

// testSweepGaForwardingRule deletes tf-testAcc* Ga forwarding rules.
// ListForwardingRules requires both AcceleratorId and ListenerId, so the
// sweeper walks accelerators -> listeners -> forwarding rules.
func testSweepGaForwardingRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	accIds, err := listGaAcceleratorIds(client)
	if err != nil {
		log.Printf("[ERROR] testSweepGaForwardingRule list accelerators failed: %v", err)
		return nil
	}

	for _, accId := range accIds {
		listeners, err := listGaListenerIds(client, accId)
		if err != nil {
			log.Printf("[ERROR] ListListeners for accelerator %s failed: %v", accId, err)
			continue
		}
		for _, listener := range listeners {
			listenerId := fmt.Sprint(listener["ListenerId"])
			request := map[string]interface{}{
				"RegionId":      client.RegionId,
				"AcceleratorId": accId,
				"ListenerId":    listenerId,
				"PageSize":      PageSizeLarge,
				"PageNumber":    1,
			}
			for {
				action := "ListForwardingRules"
				response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
				if err != nil {
					log.Printf("[ERROR] %s got an error: %v", action, err)
					break
				}
				addDebug(action, response, request)
				resp, err := jsonpath.Get("$.ForwardingRules", response)
				if err != nil {
					log.Println(err)
					break
				}
				result, _ := resp.([]interface{})
				for _, v := range result {
					fr := v.(map[string]interface{})
					frId := fmt.Sprint(fr["ForwardingRuleId"])
					frName := fmt.Sprint(fr["Name"])
					if gaSweepSkipName(frName) {
						log.Printf("[INFO] Skipping Ga forwarding rule: %s(%s)", frId, frName)
						continue
					}
					log.Printf("[Info] Delete Ga forwarding rule: %s(%s)", frId, frName)
					req := map[string]interface{}{
						"RegionId":          client.RegionId,
						"ClientToken":       buildClientToken("DeleteForwardingRules"),
						"AcceleratorId":     accId,
						"ListenerId":        listenerId,
						"ForwardingRuleIds": []string{frId},
					}
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = retry.Retry(1*time.Minute, func() *retry.RetryError {
						_, e := client.RpcPost("Ga", "2019-11-20", "DeleteForwardingRules", nil, req, true)
						if e != nil {
							if IsExpectedErrors(e, []string{"StateError.Accelerator", "StateError.ForwardingRule"}) || NeedRetry(e) {
								wait()
								return retry.RetryableError(e)
							}
							return retry.NonRetryableError(e)
						}
						return nil
					})
					if err != nil {
						log.Printf("[ERROR] Deleting Ga forwarding rule %s got an error: %s", frId, err)
					}
				}
				if len(result) < PageSizeLarge {
					break
				}
				request["PageNumber"] = request["PageNumber"].(int) + 1
			}
		}
	}
	return nil
}

// testSweepGaDomain deletes tf-testAcc* Ga domain-accelerator relations. The
// Ga domain resource models the relation between an accelerator and a domain
// name; its Delete API is DeleteDomainAcceleratorRelation. The list API
// (ListDomains) takes no required accelerator arg, so the sweeper walks it
// directly. The Ga domain resource does not carry an ICP record; the real存
// red-line guard is the NODELETING name check plus the tf-testAcc prefix
// filter, which together protect any production domain.
func testSweepGaDomain(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		action := "ListDomains"
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			break
		}
		addDebug(action, response, request)
		resp, err := jsonpath.Get("$.Domains", response)
		if err != nil {
			log.Println(err)
			break
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			domain := v.(map[string]interface{})
			domainName := fmt.Sprint(domain["Domain"])
			acceleratorId := fmt.Sprint(domain["AcceleratorId"])
			// Ga domain names are real存/ICP-backed in production. Per the
			// red-line rule, ALWAYS require the tf-testAcc prefix EVEN under
			// SWEEP_ALL — never bypass the prefix filter for domains. The
			// NODELETING name guard applies on top.
			if protectedSweepVpcName(domainName) {
				log.Printf("[INFO] Skipping Ga domain (NODELETING): %s", domainName)
				continue
			}
			domainSkip := true
			for _, prefix := range gaSweepPrefixes {
				if strings.HasPrefix(domainName, prefix) {
					domainSkip = false
					break
				}
			}
			if domainSkip {
				log.Printf("[INFO] Skipping Ga domain (no tf-testAcc prefix; real存 guard): %s", domainName)
				continue
			}
			log.Printf("[Info] Delete Ga domain relation: %s -> %s", domainName, acceleratorId)
			req := map[string]interface{}{
				"AcceleratorIds.1": acceleratorId,
				"Domain":           domainName,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(1*time.Minute, func() *retry.RetryError {
				_, e := client.RpcPost("Ga", "2019-11-20", "DeleteDomainAcceleratorRelation", nil, req, false)
				if e != nil {
					if NeedRetry(e) {
						wait()
						return retry.RetryableError(e)
					}
					return retry.NonRetryableError(e)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Deleting Ga domain relation %s got an error: %s", domainName, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

// testSweepGaBasicAccelerator deletes tf-testAcc* Ga basic accelerators. The
// basic accelerator list API is global per region. PREPAY (Subscription)
// basic accelerators cannot be destroyed via API; the sweeper records them as
// leftovers and continues.
func testSweepGaBasicAccelerator(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		action := "ListBasicAccelerators"
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			break
		}
		addDebug(action, response, request)
		resp, err := jsonpath.Get("$.Accelerators", response)
		if err != nil {
			log.Println(err)
			break
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			acc := v.(map[string]interface{})
			accId := fmt.Sprint(acc["AcceleratorId"])
			accName := fmt.Sprint(acc["Name"])
			if gaSweepSkipName(accName) {
				log.Printf("[INFO] Skipping Ga basic accelerator: %s(%s)", accId, accName)
				continue
			}
			if fmt.Sprint(acc["InstanceChargeType"]) == "PREPAY" {
				log.Printf("[INFO] Skipping PREPAY Ga basic accelerator (cannot destroy; wait for expiry): %s(%s)", accId, accName)
				continue
			}
			log.Printf("[Info] Delete Ga basic accelerator: %s(%s)", accId, accName)
			req := map[string]interface{}{
				"RegionId":      client.RegionId,
				"AcceleratorId": accId,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(1*time.Minute, func() *retry.RetryError {
				_, e := client.RpcPost("Ga", "2019-11-20", "DeleteBasicAccelerator", nil, req, true)
				if e != nil {
					if NeedRetry(e) {
						wait()
						return retry.RetryableError(e)
					}
					return retry.NonRetryableError(e)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Deleting Ga basic accelerator %s got an error: %s", accId, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

// testSweepGaCustomRoutingEndpointGroup deletes tf-testAcc* Ga custom routing
// endpoint groups. The list API requires either AcceleratorId or ListenerId;
// the sweeper iterates accelerators first.
func testSweepGaCustomRoutingEndpointGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	accIds, err := listGaAcceleratorIds(client)
	if err != nil {
		log.Printf("[ERROR] testSweepGaCustomRoutingEndpointGroup list accelerators failed: %v", err)
		return nil
	}

	for _, accId := range accIds {
		request := map[string]interface{}{
			"RegionId":      client.RegionId,
			"AcceleratorId": accId,
			"PageSize":      PageSizeLarge,
			"PageNumber":    1,
		}
		for {
			action := "ListCustomRoutingEndpointGroups"
			response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] %s got an error: %v", action, err)
				break
			}
			addDebug(action, response, request)
			resp, err := jsonpath.Get("$.EndpointGroups", response)
			if err != nil {
				log.Println(err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				eg := v.(map[string]interface{})
				egId := fmt.Sprint(eg["EndpointGroupId"])
				egName := fmt.Sprint(eg["Name"])
				if gaSweepSkipName(egName) {
					log.Printf("[INFO] Skipping Ga custom routing endpoint group: %s(%s)", egId, egName)
					continue
				}
				log.Printf("[Info] Delete Ga custom routing endpoint group: %s(%s)", egId, egName)
				req := map[string]interface{}{
					"RegionId":         client.RegionId,
					"ClientToken":      buildClientToken("DeleteCustomRoutingEndpointGroups"),
					"EndpointGroupIds": []string{egId},
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(1*time.Minute, func() *retry.RetryError {
					_, e := client.RpcPost("Ga", "2019-11-20", "DeleteCustomRoutingEndpointGroups", nil, req, true)
					if e != nil {
						if IsExpectedErrors(e, []string{"NotActive.Listener", "StateError.Accelerator"}) || NeedRetry(e) {
							wait()
							return retry.RetryableError(e)
						}
						return retry.NonRetryableError(e)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] Deleting Ga custom routing endpoint group %s got an error: %s", egId, err)
				}
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return nil
}
