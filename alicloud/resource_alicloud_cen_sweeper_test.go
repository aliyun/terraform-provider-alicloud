package alicloud

// CEN (Cloud Enterprise Network / 云企业网) sweeper gap-fills.
//
// Product code: "Cbn" (Cloud Backbone Network). API version: "2017-09-12"
// (uniform across every CEN RpcPost call site in the provider — 127 occurrences).
//
// Two API styles coexist in the provider, both mirrored here:
//   - Older resources use the typed cbn SDK via client.WithCbnClient(...):
//     cen_bandwidth_package_attachment, cen_route_map, cen_private_zone,
//     cen_vbr_health_check.
//   - Newer resources use client.RpcPost("Cbn", "2017-09-12", action, nil/<query>,
//     request, true/false) with jsonpath extraction: the transit_router_* /
//     traffic_marking_policy / inter_region_traffic_qos_* /
//     transit_router_multicast_domain_* family, plus transit_route_table_aggregation
//     and child_instance_route_entry_to_attachment (the latter deferred — see below).
//
// Fleet-wipe gate (cenSweepNameOk): the NODELETING guard (protectedSweepVpcName) is
// ALWAYS active, even under ALICLOUD_SWEEP_ALL_RESOURCES; then the standard 6 test
// prefixes match (case-insensitive). Per the existing inline CEN sweepers
// (cen_route_service, cen_instance_attachment), parent-scoped child resources gate on
// the PARENT CEN instance Name when the child has no own Name; resources with their
// own Name field gate on that Name. The parent walk (cenSweepListCens) skips
// NODELETING-marked CENs so no child sweeper can reach their descendants.
//
// Pagination is inconsistent across CEN (verified per resource): ListTransitRouters
// and ListTransitRouterPrefixListAssociation use PageNumber/PageSize (int); most
// others use MaxResults/NextToken; ListTransitRouterCidr paginates not at all; and
// DescribeTransitRouteTableAggregation returns its list under the non-standard
// jsonpath "$.Data". Delete-action verbs are also inconsistent (Dissociate...,
// Disable..., Deregister..., Disassociate..., Revoke... rather than Delete...).
//
// Resources whose sweepers already exist inline in other *_test.go files (do NOT
// duplicate — verified via grep AddTestSweepers("alicloud_cen"):
//   alicloud_cen_instance, alicloud_cen_instance_attachment,
//   alicloud_cen_bandwidth_package, alicloud_cen_bandwidth_limit,
//   alicloud_cen_flowlog, alicloud_cen_route_service.
//
// Resources being added inline by the parallel cen-sub-sweep worktree (do NOT
// duplicate — would cause a duplicate AddTestSweepers init panic):
//   alicloud_cen_transit_router_vpc_attachment,
//   alicloud_cen_transit_router_peer_attachment.
//
// DEFERRED (no global list API / cross-service enumeration / brittle parent walk —
// recorded here, not coded, per the red-line that a sweeper never force-deletes or
// walks unrelated services):
//   alicloud_cen_route_entry — list (cbn.DescribePublishedRouteEntries) requires
//     CenId + ChildInstanceRouteTableId + a child instance derived via the VPC-side
//     CreateCenRouteEntryParas helper (QueryRouteTableById + VRouterId prefix check);
//     no global list API, sweeping would walk every CEN's every child-instance route
//     table across VPC/VBR. Left for an explicit cross-service sweep if ever needed.
//   alicloud_cen_instance_grant — uses the vpc SDK (WithVpcClient); the list API
//     DescribeGrantRulesToCen is per-child-instance (InstanceId + InstanceType, where
//     InstanceType is derived from the id prefix via GetCenChildInstanceType). No
//     global list API; sweeping would walk every VPC/VBR child instance in the region.
//   alicloud_cen_transit_router_grant_attachment — the provider has no list-all data
//     source and no ListTransitRouterGrantAttachments action; its per-id Describe uses
//     DescribeGrantRulesToResource (resource-side). The candidate list-all
//     (DescribeGrantRulesToCen) is a vpc-SDK API used by the sibling instance_grant
//     and would require iterating 5 ProductType values per CEN at 10 req/s. Too
//     ambiguous/complex for a safe sweeper; left for an explicit sweep if needed.
//   alicloud_cen_child_instance_route_entry_to_attachment — list
//     (ListCenChildInstanceRouteEntriesToAttachment) requires BOTH
//     ChildInstanceRouteTableId and TransitRouterAttachmentId; the provider has no
//     action to enumerate those pairs directly, and the route-table-associations list
//     item field set needed to derive them is not reliably documented. A 4-level
//     CEN→TR→RT→association walk with an uncertain field is too brittle; left for an
//     explicit sweep if the parent-pair source is confirmed.

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// cenSweepPrefixes is the standard 6-prefix test-resource marker set used across the
// single-file sweeper campaign. Matching is case-insensitive.
var cenSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf_example",
	"tf-example",
	"terraform-example",
	"tftestacc",
}

// cenSweepNameOk is the fleet-wipe gate for CEN resources.
// The NODELETING guard (protectedSweepVpcName) is ALWAYS active — even when
// ALICLOUD_SWEEP_ALL_RESOURCES=true — so a CEN/child marked do-not-delete is never
// touched. Under SWEEP_ALL every other resource passes; otherwise the resource must
// carry one of cenSweepPrefixes (case-insensitive) on the gated Name field.
func cenSweepNameOk(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	for _, prefix := range cenSweepPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// cenSweepListCens enumerates every CEN instance in the region via the cbn SDK
// (DescribeCens), returning the full cbn.Cen objects (CenId + Name + Status).
// NODELETING-marked CENs are skipped here so that no per-CEN child sweeper can reach
// their children, regardless of SWEEP_ALL. Each caller still applies cenSweepNameOk
// to the relevant Name (its own, or the parent CEN Name) to decide prefix vs
// fleet-wipe scope, matching the existing inline cen_route_service /
// cen_instance_attachment convention.
func cenSweepListCens(client *connectivity.AliyunClient) []cbn.Cen {
	var cens []cbn.Cen
	request := cbn.CreateDescribeCensRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCens(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve CEN instances for sweep: %s", err)
			return cens
		}
		response, _ := raw.(*cbn.DescribeCensResponse)
		for _, c := range response.Cens.Cen {
			if protectedSweepVpcName(c.Name) {
				log.Printf("[INFO] Skipping NODELETING CEN instance: %s (%s)", c.Name, c.CenId)
				continue
			}
			cens = append(cens, c)
		}
		if len(response.Cens.Cen) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return cens
		}
		request.PageNumber = page
	}
	return cens
}

// cenRpcPost wraps client.RpcPost("Cbn", "2017-09-12", action, nil, body, autoRetry)
// with the provider's standard NeedRetry/incrementalWait retry loop and addDebug
// call. CEN RpcPost APIs put all params in the body (request) map and pass nil for
// the query argument; the query arg is only an empty map in some auto-generated
// resource deletes, which is functionally equivalent to nil. autoRetry is passed
// through (most actions use true; DisableTransitRouterRouteTablePropagation uses
// false to match the provider).
func cenRpcPost(client *connectivity.AliyunClient, action string, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, body, autoRetry)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, body)
	return response, err
}

// cenSweepItems extracts a jsonpath array from a response, returning []map.
// Returns nil (no error) when the path is absent or the value is not a slice — the
// caller treats absence as an empty list.
func cenSweepItems(response map[string]interface{}, path string) []map[string]interface{} {
	v, err := jsonpath.Get(path, response)
	if err != nil {
		return nil
	}
	arr, _ := v.([]interface{})
	out := make([]map[string]interface{}, 0, len(arr))
	for _, raw := range arr {
		if m, ok := raw.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

// cenSweepNextToken returns the NextToken pagination cursor from a response, if any.
func cenSweepNextToken(response map[string]interface{}) string {
	if t, ok := response["NextToken"].(string); ok {
		return t
	}
	return ""
}

// cenSweepListTransitRouters enumerates every Enterprise Edition Transit Router in the
// given CEN via ListTransitRouters (PageNumber/PageSize). TransitRouterId is the child
// id and TransitRouterName is the human Name used for prefix gating. Used as the
// parent walk for TransitRouter-scoped child resources.
func cenSweepListTransitRouters(client *connectivity.AliyunClient, cenId string) []map[string]interface{} {
	action := "ListTransitRouters"
	request := map[string]interface{}{
		"CenId":      cenId,
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	var trs []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve Transit Routers for CEN %s: %s", cenId, err)
			return trs
		}
		items := cenSweepItems(response, "$.TransitRouters")
		trs = append(trs, items...)
		if len(items) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return trs
}

// cenSweepListTransitRouterRouteTables enumerates route tables under a Transit Router
// via ListTransitRouterRouteTables (MaxResults/NextToken). Each item has
// TransitRouterRouteTableId and TransitRouterRouteTableName.
func cenSweepListTransitRouterRouteTables(client *connectivity.AliyunClient, transitRouterId string) []map[string]interface{} {
	action := "ListTransitRouterRouteTables"
	request := map[string]interface{}{
		"TransitRouterId": transitRouterId,
		"MaxResults":      PageSizeLarge,
	}
	var rts []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve route tables for TR %s: %s", transitRouterId, err)
			return rts
		}
		items := cenSweepItems(response, "$.TransitRouterRouteTables")
		rts = append(rts, items...)
		if token := cenSweepNextToken(response); token != "" {
			request["NextToken"] = token
		} else {
			break
		}
	}
	return rts
}

// cenSweepListTransitRouterMulticastDomains enumerates multicast domains under a
// Transit Router via ListTransitRouterMulticastDomains (MaxResults/NextToken). Each
// item has TransitRouterMulticastDomainId and TransitRouterMulticastDomainName.
func cenSweepListTransitRouterMulticastDomains(client *connectivity.AliyunClient, transitRouterId string) []map[string]interface{} {
	action := "ListTransitRouterMulticastDomains"
	request := map[string]interface{}{
		"TransitRouterId": transitRouterId,
		"MaxResults":      PageSizeLarge,
	}
	var mds []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve multicast domains for TR %s: %s", transitRouterId, err)
			return mds
		}
		items := cenSweepItems(response, "$.TransitRouterMulticastDomains")
		mds = append(mds, items...)
		if token := cenSweepNextToken(response); token != "" {
			request["NextToken"] = token
		} else {
			break
		}
	}
	return mds
}

// cenSweepListMulticastGroups enumerates multicast group entries (sources/members/
// peer members share this single list action) under a multicast domain via
// ListTransitRouterMulticastGroups (MaxResults/NextToken). Each item carries
// TransitRouterMulticastDomainId, GroupIpAddress, NetworkInterfaceId (source/member)
// and PeerTransitRouterMulticastDomainId (peer member).
func cenSweepListMulticastGroups(client *connectivity.AliyunClient, domainId string) []map[string]interface{} {
	action := "ListTransitRouterMulticastGroups"
	request := map[string]interface{}{
		"TransitRouterMulticastDomainId": domainId,
		"MaxResults":                     PageSizeLarge,
	}
	var groups []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve multicast groups for domain %s: %s", domainId, err)
			return groups
		}
		items := cenSweepItems(response, "$.TransitRouterMulticastGroups")
		groups = append(groups, items...)
		if token := cenSweepNextToken(response); token != "" {
			request["NextToken"] = token
		} else {
			break
		}
	}
	return groups
}

// cenSweepListAttachmentsByCen enumerates attachments of one type under a CEN via the
// given list action (ListTransitRouterVpcAttachments / VbrAttachments /
// VpnAttachments / EcrAttachments / PeerAttachments). All return their items under
// $.TransitRouterAttachments, each carrying TransitRouterAttachmentId,
// TransitRouterAttachmentName and TransitRouterId. MaxResults/NextToken pagination.
func cenSweepListAttachmentsByCen(client *connectivity.AliyunClient, cenId, action string) []map[string]interface{} {
	request := map[string]interface{}{
		"CenId":      cenId,
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeLarge,
	}
	var atts []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve attachments (%s) for CEN %s: %s", action, cenId, err)
			return atts
		}
		items := cenSweepItems(response, "$.TransitRouterAttachments")
		atts = append(atts, items...)
		if token := cenSweepNextToken(response); token != "" {
			request["NextToken"] = token
		} else {
			break
		}
	}
	return atts
}

// cenSweepListQosPolicies enumerates inter-region QoS policies for one
// (TransitRouterId, TransitRouterAttachmentId) parent pair via
// ListCenInterRegionTrafficQosPolicies (MaxResults/NextToken, $.TrafficQosPolicies).
func cenSweepListQosPolicies(client *connectivity.AliyunClient, transitRouterId, attachmentId string) []map[string]interface{} {
	action := "ListCenInterRegionTrafficQosPolicies"
	request := map[string]interface{}{
		"TransitRouterId":           transitRouterId,
		"TransitRouterAttachmentId": attachmentId,
		"MaxResults":                PageSizeLarge,
	}
	var policies []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve QoS policies (TR=%s, att=%s): %s", transitRouterId, attachmentId, err)
			return policies
		}
		items := cenSweepItems(response, "$.TrafficQosPolicies")
		policies = append(policies, items...)
		if token := cenSweepNextToken(response); token != "" {
			request["NextToken"] = token
		} else {
			break
		}
	}
	return policies
}

// cenSweepListAllQosPolicies walks CENs → VPC+Peer attachments → QoS policies and
// returns every policy item (TrafficQosPolicyId + TrafficQosPolicyName). The CEN is
// gated by Name (inter-region QoS lives on inter-region/peer attachments; both
// attachment list actions are tried, since the provider's peer-attachment list shares
// the $.TransitRouterAttachments shape). Used by the qos_policy and qos_queue
// sweepers.
func cenSweepListAllQosPolicies(client *connectivity.AliyunClient) []map[string]interface{} {
	var policies []map[string]interface{}
	for _, cen := range cenSweepListCens(client) {
		if !cenSweepNameOk(cen.Name) {
			continue
		}
		for _, action := range []string{"ListTransitRouterPeerAttachments", "ListTransitRouterVpcAttachments"} {
			for _, att := range cenSweepListAttachmentsByCen(client, cen.CenId, action) {
				trId := fmt.Sprint(att["TransitRouterId"])
				attId := fmt.Sprint(att["TransitRouterAttachmentId"])
				if trId == "" || attId == "" {
					continue
				}
				policies = append(policies, cenSweepListQosPolicies(client, trId, attId)...)
			}
		}
	}
	return policies
}

func init() {
	// cbn-SDK (v1) resources.
	resource.AddTestSweepers("alicloud_cen_bandwidth_package_attachment", &resource.Sweeper{
		Name: "alicloud_cen_bandwidth_package_attachment",
		F:    testSweepCenBandwidthPackageAttachment,
	})
	resource.AddTestSweepers("alicloud_cen_route_map", &resource.Sweeper{
		Name: "alicloud_cen_route_map",
		F:    testSweepCenRouteMap,
	})
	resource.AddTestSweepers("alicloud_cen_private_zone", &resource.Sweeper{
		Name: "alicloud_cen_private_zone",
		F:    testSweepCenPrivateZone,
	})
	resource.AddTestSweepers("alicloud_cen_vbr_health_check", &resource.Sweeper{
		Name: "alicloud_cen_vbr_health_check",
		F:    testSweepCenVbrHealthCheck,
	})

	// RpcPost (v2) resources. Dependencies list the resource's children so they run
	// first (children enumerate parents themselves, but a parent sweeper deleting the
	// parent first would orphan them; this ordering keeps a single sweep pass clean).
	resource.AddTestSweepers("alicloud_cen_transit_router_route_entry", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_route_entry",
		F:    testSweepCenTransitRouterRouteEntry,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_route_table_association", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_route_table_association",
		F:    testSweepCenTransitRouterRouteTableAssociation,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_route_table_propagation", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_route_table_propagation",
		F:    testSweepCenTransitRouterRouteTablePropagation,
	})
	resource.AddTestSweepers("alicloud_cen_transit_route_table_aggregation", &resource.Sweeper{
		Name: "alicloud_cen_transit_route_table_aggregation",
		F:    testSweepCenTransitRouteTableAggregation,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_prefix_list_association", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_prefix_list_association",
		F:    testSweepCenTransitRouterPrefixListAssociation,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_route_table", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_route_table",
		F:    testSweepCenTransitRouterRouteTable,
		Dependencies: []string{
			"alicloud_cen_transit_router_route_entry",
			"alicloud_cen_transit_router_route_table_association",
			"alicloud_cen_transit_router_route_table_propagation",
			"alicloud_cen_transit_route_table_aggregation",
			"alicloud_cen_transit_router_prefix_list_association",
		},
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_cidr", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_cidr",
		F:    testSweepCenTransitRouterCidr,
	})
	resource.AddTestSweepers("alicloud_cen_traffic_marking_policy", &resource.Sweeper{
		Name: "alicloud_cen_traffic_marking_policy",
		F:    testSweepCenTrafficMarkingPolicy,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_multicast_domain_source", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_multicast_domain_source",
		F:    testSweepCenTransitRouterMulticastDomainSource,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_multicast_domain_member", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_multicast_domain_member",
		F:    testSweepCenTransitRouterMulticastDomainMember,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_multicast_domain_peer_member", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_multicast_domain_peer_member",
		F:    testSweepCenTransitRouterMulticastDomainPeerMember,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_multicast_domain_association", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_multicast_domain_association",
		F:    testSweepCenTransitRouterMulticastDomainAssociation,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_multicast_domain", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_multicast_domain",
		F:    testSweepCenTransitRouterMulticastDomain,
		Dependencies: []string{
			"alicloud_cen_transit_router_multicast_domain_source",
			"alicloud_cen_transit_router_multicast_domain_member",
			"alicloud_cen_transit_router_multicast_domain_peer_member",
			"alicloud_cen_transit_router_multicast_domain_association",
		},
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_ecr_attachment", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_ecr_attachment",
		F:    testSweepCenTransitRouterEcrAttachment,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_vbr_attachment", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_vbr_attachment",
		F:    testSweepCenTransitRouterVbrAttachment,
	})
	resource.AddTestSweepers("alicloud_cen_transit_router_vpn_attachment", &resource.Sweeper{
		Name: "alicloud_cen_transit_router_vpn_attachment",
		F:    testSweepCenTransitRouterVpnAttachment,
	})
	resource.AddTestSweepers("alicloud_cen_inter_region_traffic_qos_queue", &resource.Sweeper{
		Name: "alicloud_cen_inter_region_traffic_qos_queue",
		F:    testSweepCenInterRegionTrafficQosQueue,
	})
	resource.AddTestSweepers("alicloud_cen_inter_region_traffic_qos_policy", &resource.Sweeper{
		Name: "alicloud_cen_inter_region_traffic_qos_policy",
		F:    testSweepCenInterRegionTrafficQosPolicy,
		Dependencies: []string{
			"alicloud_cen_inter_region_traffic_qos_queue",
		},
	})
	resource.AddTestSweepers("alicloud_cen_transit_router", &resource.Sweeper{
		Name: "alicloud_cen_transit_router",
		F:    testSweepCenTransitRouter,
		Dependencies: []string{
			"alicloud_cen_transit_router_route_table",
			"alicloud_cen_transit_router_cidr",
			"alicloud_cen_traffic_marking_policy",
			"alicloud_cen_transit_router_multicast_domain",
			"alicloud_cen_transit_router_ecr_attachment",
			"alicloud_cen_transit_router_vbr_attachment",
			"alicloud_cen_transit_router_vpn_attachment",
			"alicloud_cen_inter_region_traffic_qos_policy",
		},
	})
}

// ===========================================================================
// cbn-SDK (v1) sweepers
// ===========================================================================

// alicloud_cen_bandwidth_package_attachment
//
// There is no dedicated list API for the attachment (association).
// DescribeCenBandwidthPackages lists packages globally; a package "is" an attachment
// when its CenIds.CenId is populated. The sweeper lists all packages and, for each
// attached one whose package Name passes the gate, calls UnassociateCenBandwidthPackage.
//
// IMPORTANT: this calls Unassociate (NOT DeleteCenBandwidthPackage) — deleting the
// package itself is the job of the existing alicloud_cen_bandwidth_package sweeper.
//
// LIST: cbn.DescribeCenBandwidthPackages — .CenBandwidthPackages.CenBandwidthPackage
//
//	id = CenBandwidthPackageId ; name = Name ; attached CEN = CenIds.CenId[0]
//
// DELETE: cbn.UnassociateCenBandwidthPackage(CenId, CenBandwidthPackageId)
func testSweepCenBandwidthPackageAttachment(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := cbn.CreateDescribeCenBandwidthPackagesRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenBandwidthPackages(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve CEN bandwidth packages: %s", err)
			break
		}
		response, _ := raw.(*cbn.DescribeCenBandwidthPackagesResponse)
		items := response.CenBandwidthPackages.CenBandwidthPackage
		for _, item := range items {
			if len(item.CenIds.CenId) == 0 {
				continue // package is not attached to any CEN
			}
			if !cenSweepNameOk(item.Name) {
				log.Printf("[INFO] Skipping CEN bandwidth package attachment: %s (%s)", item.Name, item.CenBandwidthPackageId)
				continue
			}
			cenId := item.CenIds.CenId[0]
			bwpId := item.CenBandwidthPackageId
			log.Printf("[INFO] Deleting CEN bandwidth package attachment: %s (cen=%s, bwp=%s)", item.Name, cenId, bwpId)
			req := cbn.CreateUnassociateCenBandwidthPackageRequest()
			req.CenId = cenId
			req.CenBandwidthPackageId = bwpId
			if _, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.UnassociateCenBandwidthPackage(req)
			}); err != nil {
				if NotFoundError(err) {
					continue
				}
				log.Printf("[ERROR] Failed to unassociate CEN bandwidth package %s from CEN %s: %s", bwpId, cenId, err)
			}
		}
		if len(items) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	return nil
}

// alicloud_cen_route_map
//
// LIST: cbn.DescribeCenRouteMaps(CenId) — .RouteMaps.RouteMap
//
//	id = RouteMapId ; no Name (Description is the human-set field); CenRegionId per item
//
// DELETE: cbn.DeleteCenRouteMap(CenId, RouteMapId, CenRegionId)
// Parent-scoped by CenId; gate on the parent CEN Name.
func testSweepCenRouteMap(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		if !cenSweepNameOk(cen.Name) {
			log.Printf("[INFO] Skipping CEN for route map sweep: %s (%s)", cen.Name, cen.CenId)
			continue
		}
		cenId := cen.CenId
		request := cbn.CreateDescribeCenRouteMapsRequest()
		request.CenId = cenId
		request.PageSize = requests.NewInteger(PageSizeLarge)
		request.PageNumber = requests.NewInteger(1)
		for {
			raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.DescribeCenRouteMaps(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to retrieve CEN route maps (cen=%s): %s", cenId, err)
				break
			}
			response, _ := raw.(*cbn.DescribeCenRouteMapsResponse)
			items := response.RouteMaps.RouteMap
			for _, item := range items {
				log.Printf("[INFO] Deleting CEN route map: %s (cen=%s)", item.RouteMapId, cenId)
				req := cbn.CreateDeleteCenRouteMapRequest()
				req.CenId = cenId
				req.RouteMapId = item.RouteMapId
				req.CenRegionId = item.CenRegionId
				if _, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
					return cbnClient.DeleteCenRouteMap(req)
				}); err != nil {
					if NotFoundError(err) {
						continue
					}
					log.Printf("[ERROR] Failed to delete CEN route map %s (cen=%s): %s", item.RouteMapId, cenId, err)
				}
			}
			if len(items) < PageSizeLarge {
				break
			}
			page, err := getNextpageNumber(request.PageNumber)
			if err != nil {
				return WrapError(err)
			}
			request.PageNumber = page
		}
	}
	return nil
}

// alicloud_cen_private_zone
//
// LIST: cbn.DescribeCenPrivateZoneRoutes(CenId) — .PrivateZoneInfos.PrivateZoneInfo
//
//	id (composite) = CenId:AccessRegionId ; no Name ; AccessRegionId is a region
//
// DELETE: cbn.UnroutePrivateZoneInCenToVpc(CenId, AccessRegionId)
// Parent-scoped by CenId with no child Name; gate on the parent CEN Name.
func testSweepCenPrivateZone(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		if !cenSweepNameOk(cen.Name) {
			log.Printf("[INFO] Skipping CEN for private zone sweep: %s (%s)", cen.Name, cen.CenId)
			continue
		}
		cenId := cen.CenId
		request := cbn.CreateDescribeCenPrivateZoneRoutesRequest()
		request.CenId = cenId
		request.PageSize = requests.NewInteger(PageSizeLarge)
		request.PageNumber = requests.NewInteger(1)
		for {
			raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.DescribeCenPrivateZoneRoutes(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to retrieve CEN private zones (cen=%s): %s", cenId, err)
				break
			}
			response, _ := raw.(*cbn.DescribeCenPrivateZoneRoutesResponse)
			items := response.PrivateZoneInfos.PrivateZoneInfo
			for _, item := range items {
				log.Printf("[INFO] Deleting CEN private zone: cen=%s region=%s", cenId, item.AccessRegionId)
				req := cbn.CreateUnroutePrivateZoneInCenToVpcRequest()
				req.CenId = cenId
				req.AccessRegionId = item.AccessRegionId
				if _, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
					return cbnClient.UnroutePrivateZoneInCenToVpc(req)
				}); err != nil {
					if NotFoundError(err) {
						continue
					}
					log.Printf("[ERROR] Failed to delete CEN private zone (cen=%s, region=%s): %s", cenId, item.AccessRegionId, err)
				}
			}
			if len(items) < PageSizeLarge {
				break
			}
			page, err := getNextpageNumber(request.PageNumber)
			if err != nil {
				return WrapError(err)
			}
			request.PageNumber = page
		}
	}
	return nil
}

// alicloud_cen_vbr_health_check
//
// LIST: cbn.DescribeCenVbrHealthCheck(VbrInstanceRegionId=client.RegionId) —
//
//	.VbrHealthChecks.VbrHealthCheck ; each item carries the CenId required by Delete
//	(which is NOT in the resource id VbrInstanceId:VbrInstanceRegionId).
//
// DELETE: cbn.DisableCenVbrHealthCheck(VbrInstanceId, VbrInstanceRegionId, CenId)
// The list is per-region; gate on the parent CEN Name (built into a sweepable-cenId set).
func testSweepCenVbrHealthCheck(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	sweepableCenIds := make(map[string]bool)
	for _, cen := range cenSweepListCens(client) {
		if cenSweepNameOk(cen.Name) {
			sweepableCenIds[cen.CenId] = true
		}
	}

	request := cbn.CreateDescribeCenVbrHealthCheckRequest()
	request.VbrInstanceRegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenVbrHealthCheck(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve CEN VBR health checks: %s", err)
			break
		}
		response, _ := raw.(*cbn.DescribeCenVbrHealthCheckResponse)
		items := response.VbrHealthChecks.VbrHealthCheck
		for _, item := range items {
			if !sweepableCenIds[item.CenId] {
				continue
			}
			log.Printf("[INFO] Deleting CEN VBR health check: %s (cen=%s)", item.VbrInstanceId, item.CenId)
			req := cbn.CreateDisableCenVbrHealthCheckRequest()
			req.VbrInstanceId = item.VbrInstanceId
			req.VbrInstanceRegionId = item.VbrInstanceRegionId
			req.CenId = item.CenId
			if _, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.DisableCenVbrHealthCheck(req)
			}); err != nil {
				if NotFoundError(err) {
					continue
				}
				log.Printf("[ERROR] Failed to delete CEN VBR health check %s (cen=%s): %s", item.VbrInstanceId, item.CenId, err)
			}
		}
		if len(items) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	return nil
}

// ===========================================================================
// RpcPost (v2) sweepers
// ===========================================================================

// alicloud_cen_transit_router
//
// LIST: ListTransitRouters(CenId) per CEN — $.TransitRouters
//
//	id = TransitRouterId ; name = TransitRouterName
//
// DELETE: DeleteTransitRouter(TransitRouterId)
// Parent-scoped by CenId; gate on TransitRouterName.
func testSweepCenTransitRouter(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			name := fmt.Sprint(tr["TransitRouterName"])
			trId := fmt.Sprint(tr["TransitRouterId"])
			if !cenSweepNameOk(name) {
				log.Printf("[INFO] Skipping CEN Transit Router: %s (%s)", name, trId)
				continue
			}
			if err := cenSweepDeleteTransitRouter(client, trId); err != nil {
				log.Printf("[ERROR] %s", err)
			}
		}
	}
	return nil
}

func cenSweepDeleteTransitRouter(client *connectivity.AliyunClient, trId string) error {
	action := "DeleteTransitRouter"
	body := map[string]interface{}{
		"TransitRouterId": trId,
		"ClientToken":     buildClientToken(action),
	}
	log.Printf("[INFO] Deleting CEN Transit Router: %s", trId)
	if _, err := cenRpcPost(client, action, body, true); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return fmt.Errorf("failed to delete Transit Router %s: %s", trId, err)
	}
	return nil
}

// alicloud_cen_transit_router_route_table
//
// LIST: ListTransitRouterRouteTables(TransitRouterId) per TR — $.TransitRouterRouteTables
//
//	id = TransitRouterRouteTableId ; name = TransitRouterRouteTableName
//
// DELETE: DeleteTransitRouterRouteTable(TransitRouterRouteTableId)
// Parent-scoped by TransitRouterId (walk CEN→TR); gate on TransitRouterRouteTableName.
func testSweepCenTransitRouterRouteTable(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			for _, rt := range cenSweepListTransitRouterRouteTables(client, trId) {
				name := fmt.Sprint(rt["TransitRouterRouteTableName"])
				rtId := fmt.Sprint(rt["TransitRouterRouteTableId"])
				if !cenSweepNameOk(name) {
					log.Printf("[INFO] Skipping CEN TR route table: %s (%s)", name, rtId)
					continue
				}
				action := "DeleteTransitRouterRouteTable"
				body := map[string]interface{}{
					"TransitRouterRouteTableId": rtId,
					"ClientToken":               buildClientToken(action),
				}
				log.Printf("[INFO] Deleting CEN TR route table: %s (tr=%s)", rtId, trId)
				if _, err := cenRpcPost(client, action, body, true); err != nil {
					if NotFoundError(err) {
						continue
					}
					log.Printf("[ERROR] Failed to delete TR route table %s: %s", rtId, err)
				}
			}
		}
	}
	return nil
}

// alicloud_cen_transit_router_route_entry
//
// LIST: ListTransitRouterRouteEntries(TransitRouterRouteTableId) — $.TransitRouterRouteEntries
//
//	id = TransitRouterRouteEntryId ; name = TransitRouterRouteEntryName
//	each item also carries TransitRouterRouteEntryNextHopType (always sent to Delete)
//	and TransitRouterRouteEntryNextHopId (sent when non-empty).
//
// DELETE: DeleteTransitRouterRouteEntry(TransitRouterRouteEntryId, NextHopType[, NextHopId])
// Parent-scoped by TransitRouterRouteTableId (walk CEN→TR→RT); gate on the entry Name.
func testSweepCenTransitRouterRouteEntry(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			for _, rt := range cenSweepListTransitRouterRouteTables(client, trId) {
				rtId := fmt.Sprint(rt["TransitRouterRouteTableId"])
				for _, re := range cenSweepListRouteEntries(client, rtId) {
					name := fmt.Sprint(re["TransitRouterRouteEntryName"])
					reId := fmt.Sprint(re["TransitRouterRouteEntryId"])
					if !cenSweepNameOk(name) {
						log.Printf("[INFO] Skipping CEN TR route entry: %s (%s)", name, reId)
						continue
					}
					nextHopType := fmt.Sprint(re["TransitRouterRouteEntryNextHopType"])
					nextHopId := fmt.Sprint(re["TransitRouterRouteEntryNextHopId"])
					action := "DeleteTransitRouterRouteEntry"
					body := map[string]interface{}{
						"TransitRouterRouteEntryId":          reId,
						"TransitRouterRouteEntryNextHopType": nextHopType,
						"ClientToken":                        buildClientToken(action),
					}
					if nextHopId != "" {
						body["TransitRouterRouteEntryNextHopId"] = nextHopId
					}
					log.Printf("[INFO] Deleting CEN TR route entry: %s (rt=%s)", reId, rtId)
					if _, err := cenRpcPost(client, action, body, true); err != nil {
						if NotFoundError(err) {
							continue
						}
						log.Printf("[ERROR] Failed to delete TR route entry %s: %s", reId, err)
					}
				}
			}
		}
	}
	return nil
}

func cenSweepListRouteEntries(client *connectivity.AliyunClient, rtId string) []map[string]interface{} {
	action := "ListTransitRouterRouteEntries"
	request := map[string]interface{}{
		"TransitRouterRouteTableId": rtId,
		"MaxResults":                PageSizeLarge,
	}
	var out []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve TR route entries (rt=%s): %s", rtId, err)
			return out
		}
		items := cenSweepItems(response, "$.TransitRouterRouteEntries")
		out = append(out, items...)
		if token := cenSweepNextToken(response); token != "" {
			request["NextToken"] = token
		} else {
			break
		}
	}
	return out
}

// alicloud_cen_transit_router_route_table_association
//
// LIST: ListTransitRouterRouteTableAssociations(TransitRouterRouteTableId) —
//
//	$.TransitRouterAssociations ; id = TransitRouterAttachmentId (no Name).
//
// DELETE: DissociateTransitRouterAttachmentFromRouteTable(TransitRouterAttachmentId,
//
//	TransitRouterRouteTableId)  — note the Dissociate verb (not Delete).
//
// Parent-scoped by TransitRouterRouteTableId (walk CEN→TR→RT); gate on the parent CEN Name.
func testSweepCenTransitRouterRouteTableAssociation(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		if !cenSweepNameOk(cen.Name) {
			continue
		}
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			for _, rt := range cenSweepListTransitRouterRouteTables(client, trId) {
				rtId := fmt.Sprint(rt["TransitRouterRouteTableId"])
				for _, assoc := range cenSweepListRouteTableAssociations(client, rtId) {
					attId := fmt.Sprint(assoc["TransitRouterAttachmentId"])
					if attId == "" {
						continue
					}
					action := "DissociateTransitRouterAttachmentFromRouteTable"
					body := map[string]interface{}{
						"TransitRouterAttachmentId": attId,
						"TransitRouterRouteTableId": rtId,
						"ClientToken":               buildClientToken(action),
					}
					log.Printf("[INFO] Deleting CEN TR route table association: att=%s rt=%s", attId, rtId)
					if _, err := cenRpcPost(client, action, body, true); err != nil {
						if NotFoundError(err) {
							continue
						}
						log.Printf("[ERROR] Failed to delete TR route table association %s: %s", attId, err)
					}
				}
			}
		}
	}
	return nil
}

func cenSweepListRouteTableAssociations(client *connectivity.AliyunClient, rtId string) []map[string]interface{} {
	action := "ListTransitRouterRouteTableAssociations"
	request := map[string]interface{}{
		"TransitRouterRouteTableId": rtId,
		"MaxResults":                PageSizeLarge,
	}
	var out []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve TR route table associations (rt=%s): %s", rtId, err)
			return out
		}
		items := cenSweepItems(response, "$.TransitRouterAssociations")
		out = append(out, items...)
		if token := cenSweepNextToken(response); token != "" {
			request["NextToken"] = token
		} else {
			break
		}
	}
	return out
}

// alicloud_cen_transit_router_route_table_propagation
//
// LIST: ListTransitRouterRouteTablePropagations(TransitRouterRouteTableId) —
//
//	$.TransitRouterPropagations ; id = TransitRouterAttachmentId (no Name).
//
// DELETE: DisableTransitRouterRouteTablePropagation(TransitRouterAttachmentId,
//
//	TransitRouterRouteTableId) — note the Disable verb (not Delete); autoRetry=false.
//
// Parent-scoped by TransitRouterRouteTableId (walk CEN→TR→RT); gate on the parent CEN Name.
func testSweepCenTransitRouterRouteTablePropagation(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		if !cenSweepNameOk(cen.Name) {
			continue
		}
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			for _, rt := range cenSweepListTransitRouterRouteTables(client, trId) {
				rtId := fmt.Sprint(rt["TransitRouterRouteTableId"])
				for _, prop := range cenSweepListRouteTablePropagations(client, rtId) {
					attId := fmt.Sprint(prop["TransitRouterAttachmentId"])
					if attId == "" {
						continue
					}
					action := "DisableTransitRouterRouteTablePropagation"
					body := map[string]interface{}{
						"TransitRouterAttachmentId": attId,
						"TransitRouterRouteTableId": rtId,
					}
					log.Printf("[INFO] Deleting CEN TR route table propagation: att=%s rt=%s", attId, rtId)
					if _, err := cenRpcPost(client, action, body, false); err != nil {
						if NotFoundError(err) {
							continue
						}
						log.Printf("[ERROR] Failed to delete TR route table propagation %s: %s", attId, err)
					}
				}
			}
		}
	}
	return nil
}

func cenSweepListRouteTablePropagations(client *connectivity.AliyunClient, rtId string) []map[string]interface{} {
	action := "ListTransitRouterRouteTablePropagations"
	request := map[string]interface{}{
		"TransitRouterRouteTableId": rtId,
		"MaxResults":                PageSizeLarge,
	}
	var out []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve TR route table propagations (rt=%s): %s", rtId, err)
			return out
		}
		items := cenSweepItems(response, "$.TransitRouterPropagations")
		out = append(out, items...)
		if token := cenSweepNextToken(response); token != "" {
			request["NextToken"] = token
		} else {
			break
		}
	}
	return out
}

// alicloud_cen_transit_route_table_aggregation
//
// LIST: DescribeTransitRouteTableAggregation(TransitRouteTableId) — $.Data (odd path)
//
//	item fields: TrRouteTableId, TransitRouteTableAggregationCidr, Name.
//	NOTE: the LIST scope field is TransitRouteTableId (NOT TransitRouterRouteTableId —
//	no "Router"), but its VALUE equals a route table's TransitRouterRouteTableId.
//
// DELETE: DeleteTransitRouteTableAggregation(TransitRouteTableId,
//
//	TransitRouteTableAggregationCidr). The resource id uses '#' separator but the
//	sweeper passes the two fields directly from the list item.
//
// Parent-scoped by TransitRouteTableId (walk CEN→TR→RT); gate on the aggregation Name.
func testSweepCenTransitRouteTableAggregation(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			for _, rt := range cenSweepListTransitRouterRouteTables(client, trId) {
				rtId := fmt.Sprint(rt["TransitRouterRouteTableId"])
				for _, agg := range cenSweepListTransitRouteTableAggregations(client, rtId) {
					name := fmt.Sprint(agg["Name"])
					cidr := fmt.Sprint(agg["TransitRouteTableAggregationCidr"])
					// TrRouteTableId is the route-table id echoed back; prefer it, fall back to rtId.
					aggRtId := fmt.Sprint(agg["TrRouteTableId"])
					if aggRtId == "" {
						aggRtId = rtId
					}
					if !cenSweepNameOk(name) {
						log.Printf("[INFO] Skipping CEN transit route table aggregation: %s (%s)", name, cidr)
						continue
					}
					action := "DeleteTransitRouteTableAggregation"
					body := map[string]interface{}{
						"TransitRouteTableId":              aggRtId,
						"TransitRouteTableAggregationCidr": cidr,
						"ClientToken":                      buildClientToken(action),
					}
					log.Printf("[INFO] Deleting CEN transit route table aggregation: %s (rt=%s)", cidr, aggRtId)
					if _, err := cenRpcPost(client, action, body, true); err != nil {
						if NotFoundError(err) {
							continue
						}
						log.Printf("[ERROR] Failed to delete transit route table aggregation %s: %s", cidr, err)
					}
				}
			}
		}
	}
	return nil
}

func cenSweepListTransitRouteTableAggregations(client *connectivity.AliyunClient, rtId string) []map[string]interface{} {
	action := "DescribeTransitRouteTableAggregation"
	request := map[string]interface{}{
		"TransitRouteTableId": rtId,
		"MaxResults":          PageSizeLarge,
		"ClientToken":         buildClientToken(action),
	}
	var out []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve transit route table aggregations (rt=%s): %s", rtId, err)
			return out
		}
		items := cenSweepItems(response, "$.Data")
		out = append(out, items...)
		if token := cenSweepNextToken(response); token != "" {
			request["NextToken"] = token
		} else {
			break
		}
	}
	return out
}

// alicloud_cen_transit_router_prefix_list_association
//
// LIST: ListTransitRouterPrefixListAssociation (singular) — $.PrefixLists
//
//	scope: TransitRouterId + TransitRouterTableId (no "Route"); PageNumber/PageSize
//	pagination (the only one of this family not using MaxResults/NextToken).
//	item id (composite) = PrefixListId:TransitRouterId:TransitRouterTableId:NextHop
//
// DELETE: DeleteTransitRouterPrefixListAssociation(PrefixListId, TransitRouterId,
//
//	TransitRouterTableId, NextHop). No Name on items; gate on the parent CEN Name.
func testSweepCenTransitRouterPrefixListAssociation(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		if !cenSweepNameOk(cen.Name) {
			continue
		}
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			for _, rt := range cenSweepListTransitRouterRouteTables(client, trId) {
				rtId := fmt.Sprint(rt["TransitRouterRouteTableId"])
				for _, pla := range cenSweepListPrefixListAssociations(client, trId, rtId) {
					prefixListId := fmt.Sprint(pla["PrefixListId"])
					nextHop := fmt.Sprint(pla["NextHop"])
					plaTrId := fmt.Sprint(pla["TransitRouterId"])
					if plaTrId == "" {
						plaTrId = trId
					}
					plaRtId := fmt.Sprint(pla["TransitRouterTableId"])
					if plaRtId == "" {
						plaRtId = rtId
					}
					if prefixListId == "" {
						continue
					}
					action := "DeleteTransitRouterPrefixListAssociation"
					body := map[string]interface{}{
						"RegionId":             client.RegionId,
						"PrefixListId":         prefixListId,
						"TransitRouterId":      plaTrId,
						"TransitRouterTableId": plaRtId,
						"NextHop":              nextHop,
						"ClientToken":          buildClientToken(action),
					}
					log.Printf("[INFO] Deleting CEN TR prefix list association: %s (tr=%s rt=%s)", prefixListId, plaTrId, plaRtId)
					if _, err := cenRpcPost(client, action, body, true); err != nil {
						if NotFoundError(err) {
							continue
						}
						log.Printf("[ERROR] Failed to delete TR prefix list association %s: %s", prefixListId, err)
					}
				}
			}
		}
	}
	return nil
}

func cenSweepListPrefixListAssociations(client *connectivity.AliyunClient, trId, rtId string) []map[string]interface{} {
	action := "ListTransitRouterPrefixListAssociation"
	request := map[string]interface{}{
		"RegionId":             client.RegionId,
		"TransitRouterId":      trId,
		"TransitRouterTableId": rtId,
		"PageSize":             PageSizeLarge,
		"PageNumber":           1,
	}
	var out []map[string]interface{}
	for {
		response, err := cenRpcPost(client, action, request, true)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve TR prefix list associations (tr=%s rt=%s): %s", trId, rtId, err)
			return out
		}
		items := cenSweepItems(response, "$.PrefixLists")
		out = append(out, items...)
		if len(items) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return out
}

// alicloud_cen_transit_router_cidr
//
// LIST: ListTransitRouterCidr(TransitRouterId) — $.CidrLists (NO pagination loop;
//
//	a single call returns all cidrs for the TR). item: TransitRouterCidrId, Name.
//
// DELETE: DeleteTransitRouterCidr(TransitRouterCidrId, TransitRouterId)
// Parent-scoped by TransitRouterId (walk CEN→TR); gate on the cidr Name.
func testSweepCenTransitRouterCidr(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			action := "ListTransitRouterCidr"
			request := map[string]interface{}{
				"TransitRouterId": trId,
				"RegionId":        client.RegionId,
				"ClientToken":     buildClientToken(action),
			}
			response, err := cenRpcPost(client, action, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to retrieve TR cidrs (tr=%s): %s", trId, err)
				continue
			}
			for _, cidr := range cenSweepItems(response, "$.CidrLists") {
				name := fmt.Sprint(cidr["Name"])
				cidrId := fmt.Sprint(cidr["TransitRouterCidrId"])
				if !cenSweepNameOk(name) {
					log.Printf("[INFO] Skipping CEN TR cidr: %s (%s)", name, cidrId)
					continue
				}
				delAction := "DeleteTransitRouterCidr"
				body := map[string]interface{}{
					"TransitRouterCidrId": cidrId,
					"TransitRouterId":     trId,
					"RegionId":            client.RegionId,
					"ClientToken":         buildClientToken(delAction),
				}
				log.Printf("[INFO] Deleting CEN TR cidr: %s (tr=%s)", cidrId, trId)
				if _, err := cenRpcPost(client, delAction, body, true); err != nil {
					if NotFoundError(err) {
						continue
					}
					log.Printf("[ERROR] Failed to delete TR cidr %s: %s", cidrId, err)
				}
			}
		}
	}
	return nil
}

// alicloud_cen_traffic_marking_policy
//
// LIST: ListTrafficMarkingPolicies(TransitRouterId) — $.TrafficMarkingPolicies
//
//	id = TrafficMarkingPolicyId ; name = TrafficMarkingPolicyName
//
// DELETE: DeleteTrafficMarkingPolicy(TrafficMarkingPolicyId)  (TransitRouterId NOT sent)
// Parent-scoped by TransitRouterId (walk CEN→TR); gate on the policy Name.
func testSweepCenTrafficMarkingPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			action := "ListTrafficMarkingPolicies"
			request := map[string]interface{}{
				"TransitRouterId": trId,
				"MaxResults":      PageSizeLarge,
			}
			for {
				response, err := cenRpcPost(client, action, request, true)
				if err != nil {
					log.Printf("[ERROR] Failed to retrieve traffic marking policies (tr=%s): %s", trId, err)
					break
				}
				for _, pol := range cenSweepItems(response, "$.TrafficMarkingPolicies") {
					name := fmt.Sprint(pol["TrafficMarkingPolicyName"])
					polId := fmt.Sprint(pol["TrafficMarkingPolicyId"])
					if !cenSweepNameOk(name) {
						log.Printf("[INFO] Skipping CEN traffic marking policy: %s (%s)", name, polId)
						continue
					}
					delAction := "DeleteTrafficMarkingPolicy"
					body := map[string]interface{}{
						"TrafficMarkingPolicyId": polId,
						"ClientToken":            buildClientToken(delAction),
					}
					log.Printf("[INFO] Deleting CEN traffic marking policy: %s (tr=%s)", polId, trId)
					if _, err := cenRpcPost(client, delAction, body, true); err != nil {
						if NotFoundError(err) {
							continue
						}
						log.Printf("[ERROR] Failed to delete traffic marking policy %s: %s", polId, err)
					}
				}
				if token := cenSweepNextToken(response); token != "" {
					request["NextToken"] = token
				} else {
					break
				}
			}
		}
	}
	return nil
}

// alicloud_cen_transit_router_multicast_domain
//
// LIST: ListTransitRouterMulticastDomains(TransitRouterId) — $.TransitRouterMulticastDomains
//
//	id = TransitRouterMulticastDomainId ; name = TransitRouterMulticastDomainName
//
// DELETE: DeleteTransitRouterMulticastDomain(TransitRouterMulticastDomainId)
// Parent-scoped by TransitRouterId (walk CEN→TR); gate on the domain Name.
func testSweepCenTransitRouterMulticastDomain(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			for _, md := range cenSweepListTransitRouterMulticastDomains(client, trId) {
				name := fmt.Sprint(md["TransitRouterMulticastDomainName"])
				mdId := fmt.Sprint(md["TransitRouterMulticastDomainId"])
				if !cenSweepNameOk(name) {
					log.Printf("[INFO] Skipping CEN TR multicast domain: %s (%s)", name, mdId)
					continue
				}
				action := "DeleteTransitRouterMulticastDomain"
				body := map[string]interface{}{
					"TransitRouterMulticastDomainId": mdId,
					"ClientToken":                    buildClientToken(action),
				}
				log.Printf("[INFO] Deleting CEN TR multicast domain: %s (tr=%s)", mdId, trId)
				if _, err := cenRpcPost(client, action, body, true); err != nil {
					if NotFoundError(err) {
						continue
					}
					log.Printf("[ERROR] Failed to delete TR multicast domain %s: %s", mdId, err)
				}
			}
		}
	}
	return nil
}

// alicloud_cen_transit_router_multicast_domain_source
//
// LIST: ListTransitRouterMulticastGroups(TransitRouterMulticastDomainId) —
//
//	$.TransitRouterMulticastGroups (shared with member & peer_member).
//	id (composite) = TRMDId:GroupIpAddress:NetworkInterfaceId
//
// DELETE: DeregisterTransitRouterMulticastGroupSources(TransitRouterMulticastDomainId,
//
//	GroupIpAddress, NetworkInterfaceIds.1)  — note the dotted-key string form.
//
// Parent-scoped by TransitRouterMulticastDomainId (walk CEN→TR→MD); gate on the parent CEN Name.
func testSweepCenTransitRouterMulticastDomainSource(region string) error {
	return cenSweepMulticastGroupMembers(region, "source")
}

// alicloud_cen_transit_router_multicast_domain_member
//
// Same LIST as source. DELETE: DeregisterTransitRouterMulticastGroupMembers(
//
//	TransitRouterMulticastDomainId, GroupIpAddress, NetworkInterfaceIds []string)  —
//	note the []string array form (vs the dotted-key form used by source).
func testSweepCenTransitRouterMulticastDomainMember(region string) error {
	return cenSweepMulticastGroupMembers(region, "member")
}

// alicloud_cen_transit_router_multicast_domain_peer_member
//
// Same LIST (ListTransitRouterMulticastGroups). The peer member's 3rd id segment is
// PeerTransitRouterMulticastDomainId (NOT NetworkInterfaceId). DELETE reuses the member
// action (DeregisterTransitRouterMulticastGroupMembers) but with the dotted-key
// PeerTransitRouterMulticastDomains.1 instead of NetworkInterfaceIds.
func testSweepCenTransitRouterMulticastDomainPeerMember(region string) error {
	return cenSweepMulticastGroupMembers(region, "peer_member")
}

func cenSweepMulticastGroupMembers(region, kind string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		if !cenSweepNameOk(cen.Name) {
			continue
		}
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			for _, md := range cenSweepListTransitRouterMulticastDomains(client, trId) {
				mdId := fmt.Sprint(md["TransitRouterMulticastDomainId"])
				for _, g := range cenSweepListMulticastGroups(client, mdId) {
					groupIp := fmt.Sprint(g["GroupIpAddress"])
					domainId := fmt.Sprint(g["TransitRouterMulticastDomainId"])
					if domainId == "" {
						domainId = mdId
					}
					if groupIp == "" {
						continue
					}
					if err := cenSweepDeregisterMulticastMember(client, kind, domainId, groupIp, g); err != nil {
						log.Printf("[ERROR] %s", err)
					}
				}
			}
		}
	}
	return nil
}

func cenSweepDeregisterMulticastMember(client *connectivity.AliyunClient, kind, domainId, groupIp string, g map[string]interface{}) error {
	var action string
	body := map[string]interface{}{
		"TransitRouterMulticastDomainId": domainId,
		"GroupIpAddress":                 groupIp,
	}
	switch kind {
	case "source":
		niId := fmt.Sprint(g["NetworkInterfaceId"])
		action = "DeregisterTransitRouterMulticastGroupSources"
		body["NetworkInterfaceIds.1"] = niId
		body["ClientToken"] = buildClientToken(action)
		if niId == "" {
			return nil
		}
	case "member":
		niId := fmt.Sprint(g["NetworkInterfaceId"])
		action = "DeregisterTransitRouterMulticastGroupMembers"
		body["NetworkInterfaceIds"] = []string{niId}
		body["ClientToken"] = buildClientToken(action)
		if niId == "" {
			return nil
		}
	case "peer_member":
		peerMdId := fmt.Sprint(g["PeerTransitRouterMulticastDomainId"])
		action = "DeregisterTransitRouterMulticastGroupMembers"
		body["PeerTransitRouterMulticastDomains.1"] = peerMdId
		body["ClientToken"] = buildClientToken(action)
		if peerMdId == "" {
			return nil
		}
	default:
		return nil
	}
	log.Printf("[INFO] Deleting CEN multicast %s: domain=%s group=%s", kind, domainId, groupIp)
	if _, err := cenRpcPost(client, action, body, true); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return fmt.Errorf("failed to delete multicast %s (domain=%s group=%s): %s", kind, domainId, groupIp, err)
	}
	return nil
}

// alicloud_cen_transit_router_multicast_domain_association
//
// LIST: ListTransitRouterMulticastDomainAssociations(TransitRouterMulticastDomainId) —
//
//	$.TransitRouterMulticastAssociations
//	id (composite) = TRMDId:TransitRouterAttachmentId:VSwitchId
//
// DELETE: DisassociateTransitRouterMulticastDomain(TransitRouterMulticastDomainId,
//
//	TransitRouterAttachmentId, VSwitchIds []string)  — note the []string array form.
//
// Parent-scoped by TransitRouterMulticastDomainId (walk CEN→TR→MD); gate on parent CEN Name.
func testSweepCenTransitRouterMulticastDomainAssociation(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		if !cenSweepNameOk(cen.Name) {
			continue
		}
		for _, tr := range cenSweepListTransitRouters(client, cen.CenId) {
			trId := fmt.Sprint(tr["TransitRouterId"])
			for _, md := range cenSweepListTransitRouterMulticastDomains(client, trId) {
				mdId := fmt.Sprint(md["TransitRouterMulticastDomainId"])
				action := "ListTransitRouterMulticastDomainAssociations"
				request := map[string]interface{}{
					"TransitRouterMulticastDomainId": mdId,
					"MaxResults":                     PageSizeLarge,
				}
				for {
					response, err := cenRpcPost(client, action, request, true)
					if err != nil {
						log.Printf("[ERROR] Failed to retrieve multicast domain associations (md=%s): %s", mdId, err)
						break
					}
					for _, assoc := range cenSweepItems(response, "$.TransitRouterMulticastAssociations") {
						attId := fmt.Sprint(assoc["TransitRouterAttachmentId"])
						vSwitchId := fmt.Sprint(assoc["VSwitchId"])
						domainId := fmt.Sprint(assoc["TransitRouterMulticastDomainId"])
						if domainId == "" {
							domainId = mdId
						}
						if attId == "" || vSwitchId == "" {
							continue
						}
						delAction := "DisassociateTransitRouterMulticastDomain"
						body := map[string]interface{}{
							"TransitRouterMulticastDomainId": domainId,
							"TransitRouterAttachmentId":      attId,
							"VSwitchIds":                     []string{vSwitchId},
							"ClientToken":                    buildClientToken(delAction),
						}
						log.Printf("[INFO] Deleting CEN multicast domain association: att=%s vsw=%s", attId, vSwitchId)
						if _, err := cenRpcPost(client, delAction, body, true); err != nil {
							if NotFoundError(err) {
								continue
							}
							log.Printf("[ERROR] Failed to delete multicast domain association %s: %s", attId, err)
						}
					}
					if token := cenSweepNextToken(response); token != "" {
						request["NextToken"] = token
					} else {
						break
					}
				}
			}
		}
	}
	return nil
}

// alicloud_cen_transit_router_ecr_attachment
//
// LIST: ListTransitRouterEcrAttachments(CenId) per CEN — $.TransitRouterAttachments
//
//	id = TransitRouterAttachmentId (plain) ; name = TransitRouterAttachmentName
//
// DELETE: DeleteTransitRouterEcrAttachment(TransitRouterAttachmentId)
// Parent-scoped by CenId; gate on the attachment Name.
func testSweepCenTransitRouterEcrAttachment(region string) error {
	return cenSweepAttachmentByCen(region, "ListTransitRouterEcrAttachments", "DeleteTransitRouterEcrAttachment", "ecr")
}

// alicloud_cen_transit_router_vbr_attachment
//
// LIST: ListTransitRouterVbrAttachments(CenId) — same shape. DELETE:
//
//	DeleteTransitRouterVbrAttachment(TransitRouterAttachmentId).
func testSweepCenTransitRouterVbrAttachment(region string) error {
	return cenSweepAttachmentByCen(region, "ListTransitRouterVbrAttachments", "DeleteTransitRouterVbrAttachment", "vbr")
}

// alicloud_cen_transit_router_vpn_attachment
//
// LIST: ListTransitRouterVpnAttachments(CenId) — same shape. DELETE:
//
//	DeleteTransitRouterVpnAttachment(TransitRouterAttachmentId).
func testSweepCenTransitRouterVpnAttachment(region string) error {
	return cenSweepAttachmentByCen(region, "ListTransitRouterVpnAttachments", "DeleteTransitRouterVpnAttachment", "vpn")
}

func cenSweepAttachmentByCen(region, listAction, deleteAction, kind string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, cen := range cenSweepListCens(client) {
		for _, att := range cenSweepListAttachmentsByCen(client, cen.CenId, listAction) {
			name := fmt.Sprint(att["TransitRouterAttachmentName"])
			attId := fmt.Sprint(att["TransitRouterAttachmentId"])
			if !cenSweepNameOk(name) {
				log.Printf("[INFO] Skipping CEN TR %s attachment: %s (%s)", kind, name, attId)
				continue
			}
			body := map[string]interface{}{
				"TransitRouterAttachmentId": attId,
				"ClientToken":               buildClientToken(deleteAction),
			}
			log.Printf("[INFO] Deleting CEN TR %s attachment: %s", kind, attId)
			if _, err := cenRpcPost(client, deleteAction, body, true); err != nil {
				if NotFoundError(err) {
					continue
				}
				log.Printf("[ERROR] Failed to delete TR %s attachment %s: %s", kind, attId, err)
			}
		}
	}
	return nil
}

// alicloud_cen_inter_region_traffic_qos_policy
//
// LIST: ListCenInterRegionTrafficQosPolicies(TransitRouterId, TransitRouterAttachmentId)
//
//	— $.TrafficQosPolicies ; id = TrafficQosPolicyId ; name = TrafficQosPolicyName.
//	Both scope fields are required, so the sweeper walks CENs → VPC+Peer attachments
//	(cenSweepListAllQosPolicies) to obtain every (TR, attachment) parent pair.
//
// DELETE: DeleteCenInterRegionTrafficQosPolicy(TrafficQosPolicyId)
// Gate on the policy Name.
func testSweepCenInterRegionTrafficQosPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, pol := range cenSweepListAllQosPolicies(client) {
		name := fmt.Sprint(pol["TrafficQosPolicyName"])
		polId := fmt.Sprint(pol["TrafficQosPolicyId"])
		if !cenSweepNameOk(name) {
			log.Printf("[INFO] Skipping CEN inter-region QoS policy: %s (%s)", name, polId)
			continue
		}
		action := "DeleteCenInterRegionTrafficQosPolicy"
		body := map[string]interface{}{
			"TrafficQosPolicyId": polId,
			"ClientToken":        buildClientToken(action),
		}
		log.Printf("[INFO] Deleting CEN inter-region QoS policy: %s", polId)
		if _, err := cenRpcPost(client, action, body, true); err != nil {
			if NotFoundError(err) {
				continue
			}
			log.Printf("[ERROR] Failed to delete inter-region QoS policy %s: %s", polId, err)
		}
	}
	return nil
}

// alicloud_cen_inter_region_traffic_qos_queue
//
// LIST: ListCenInterRegionTrafficQosQueues(TrafficQosPolicyId) — $.TrafficQosQueues
//
//	item id key = TrafficQosQueueId ; name = TrafficQosQueueName.
//
// DELETE: DeleteCenInterRegionTrafficQosQueue(QosQueueId) — NOTE the delete param key
//
//	is QosQueueId (no "Traffic" prefix), mapping from the list item's
//	TrafficQosQueueId. Same id value, different key.
//
// Parent-scoped by TrafficQosPolicyId; enumerate policies via cenSweepListAllQosPolicies.
// Gate on the queue Name.
func testSweepCenInterRegionTrafficQosQueue(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, pol := range cenSweepListAllQosPolicies(client) {
		polId := fmt.Sprint(pol["TrafficQosPolicyId"])
		if polId == "" {
			continue
		}
		action := "ListCenInterRegionTrafficQosQueues"
		request := map[string]interface{}{
			"TrafficQosPolicyId": polId,
			"MaxResults":         PageSizeLarge,
		}
		for {
			response, err := cenRpcPost(client, action, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to retrieve QoS queues (policy=%s): %s", polId, err)
				break
			}
			for _, q := range cenSweepItems(response, "$.TrafficQosQueues") {
				name := fmt.Sprint(q["TrafficQosQueueName"])
				queueId := fmt.Sprint(q["TrafficQosQueueId"]) // list key
				if !cenSweepNameOk(name) {
					log.Printf("[INFO] Skipping CEN inter-region QoS queue: %s (%s)", name, queueId)
					continue
				}
				delAction := "DeleteCenInterRegionTrafficQosQueue"
				body := map[string]interface{}{
					"QosQueueId":  queueId, // delete key (no "Traffic" prefix)
					"ClientToken": buildClientToken(delAction),
				}
				log.Printf("[INFO] Deleting CEN inter-region QoS queue: %s (policy=%s)", queueId, polId)
				if _, err := cenRpcPost(client, delAction, body, true); err != nil {
					if NotFoundError(err) {
						continue
					}
					log.Printf("[ERROR] Failed to delete inter-region QoS queue %s: %s", queueId, err)
				}
			}
			if token := cenSweepNextToken(response); token != "" {
				request["NextToken"] = token
			} else {
				break
			}
		}
	}
	return nil
}
