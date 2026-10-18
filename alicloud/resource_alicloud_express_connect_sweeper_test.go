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

// Express Connect sweepers cover the gap resources that the original inline
// virtual_border_router sweeper did not register. Physical connections are
// leased-circuit infrastructure, express_connect_router_interface dispatches on
// payment_type (a Subscription path uses BSS RefundInstance), and
// traffic_qos_association has no dedicated delete action (it mutates the parent
// QoS); those four are intentionally deferred and noted in the campaign report.
//
// Two Alibaba Cloud products are involved:
//   - Vpc / 2016-04-28 for the VBR / physical-connection / grant-rule /
//     traffic-qos family.
//   - ExpressConnectRouter / 2023-09-01 for the Express Connect Router family
//     (the ECR parent plus its grant / tr / vbr-child / vpc associations).
//
// Parent/child ordering is expressed through Dependencies so that children are
// always detached before their parent is deleted:
//   - ECR-children (grant/tr/vbr-child/vpc associations) run before the ECR
//     sweeper.
//   - traffic_qos_rule runs before traffic_qos_queue, which runs before
//     traffic_qos.
//   - The two VBR-children (grant_rule_to_cen, vbr_pconn_association) are wired
//     as dependencies of the existing virtual_border_router sweeper below, so
//     they run before VBRs are deleted.
//   - virtual_physical_connection depends on vbr_pconn_association so a virtual
//     physical connection is unassociated from its VBR before it is deleted.

func init() {
	// VBR-scoped children: revoke CEN grants and unassociate physical
	// connections before the existing virtual_border_router sweeper deletes
	// the VBRs.
	resource.AddTestSweepers("alicloud_express_connect_grant_rule_to_cen", &resource.Sweeper{
		Name: "alicloud_express_connect_grant_rule_to_cen",
		F:    testSweepExpressConnectGrantRuleToCens,
	})
	resource.AddTestSweepers("alicloud_express_connect_vbr_pconn_association", &resource.Sweeper{
		Name: "alicloud_express_connect_vbr_pconn_association",
		F:    testSweepExpressConnectVbrPconnAssociations,
	})

	// Express Connect Router family: children before the ECR parent.
	resource.AddTestSweepers("alicloud_express_connect_router_grant_association", &resource.Sweeper{
		Name: "alicloud_express_connect_router_grant_association",
		F:    testSweepExpressConnectRouterGrantAssociations,
	})
	resource.AddTestSweepers("alicloud_express_connect_router_tr_association", &resource.Sweeper{
		Name: "alicloud_express_connect_router_tr_association",
		F:    testSweepExpressConnectRouterTrAssociations,
	})
	resource.AddTestSweepers("alicloud_express_connect_router_vbr_child_instance", &resource.Sweeper{
		Name: "alicloud_express_connect_router_vbr_child_instance",
		F:    testSweepExpressConnectRouterVbrChildInstances,
	})
	resource.AddTestSweepers("alicloud_express_connect_router_vpc_association", &resource.Sweeper{
		Name: "alicloud_express_connect_router_vpc_association",
		F:    testSweepExpressConnectRouterVpcAssociations,
	})
	resource.AddTestSweepers("alicloud_express_connect_router_express_connect_router", &resource.Sweeper{
		Name: "alicloud_express_connect_router_express_connect_router",
		F:    testSweepExpressConnectRouterExpressConnectRouters,
		Dependencies: []string{
			"alicloud_express_connect_router_grant_association",
			"alicloud_express_connect_router_tr_association",
			"alicloud_express_connect_router_vbr_child_instance",
			"alicloud_express_connect_router_vpc_association",
		},
	})

	// Traffic QoS family: rules -> queues -> QoS.
	resource.AddTestSweepers("alicloud_express_connect_traffic_qos_rule", &resource.Sweeper{
		Name: "alicloud_express_connect_traffic_qos_rule",
		F:    testSweepExpressConnectTrafficQosRules,
	})
	resource.AddTestSweepers("alicloud_express_connect_traffic_qos_queue", &resource.Sweeper{
		Name: "alicloud_express_connect_traffic_qos_queue",
		F:    testSweepExpressConnectTrafficQosQueues,
		Dependencies: []string{
			"alicloud_express_connect_traffic_qos_rule",
		},
	})
	resource.AddTestSweepers("alicloud_express_connect_traffic_qos", &resource.Sweeper{
		Name: "alicloud_express_connect_traffic_qos",
		F:    testSweepExpressConnectTrafficQoses,
		Dependencies: []string{
			"alicloud_express_connect_traffic_qos_queue",
			"alicloud_express_connect_traffic_qos_rule",
		},
	})

	// Virtual physical connection: unassociate from VBRs before deleting.
	resource.AddTestSweepers("alicloud_express_connect_virtual_physical_connection", &resource.Sweeper{
		Name: "alicloud_express_connect_virtual_physical_connection",
		F:    testSweepExpressConnectVirtualPhysicalConnections,
		Dependencies: []string{
			"alicloud_express_connect_vbr_pconn_association",
		},
	})
}

// expressConnectSweepPrefixes are the test-resource name prefixes the sweepers
// target when ALICLOUD_SWEEP_ALL_RESOURCES is not set.
var expressConnectSweepPrefixes = []string{
	"tf-testAcc", "tf_testAcc", "tf-example", "tf_example", "terraform-example", "tftestacc",
}

// expressConnectSweepNameOk is the fleet-wipe gate shared by all express_connect
// sweepers. It always spares resources whose name carries a NODELETING marker
// (protectedSweepVpcName, an always-active guard that SWEEP_ALL cannot bypass),
// and otherwise honours the tf-prefix filter unless sweepAll() is set.
func expressConnectSweepNameOk(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	lower := strings.ToLower(name)
	for _, p := range expressConnectSweepPrefixes {
		if strings.HasPrefix(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// expressConnectGetField returns the first non-nil value on item whose key
// case-insensitively matches any of names. Express Connect list APIs are
// inconsistent about casing (NodeId vs InstanceId, QosName vs Name), so this
// lets a sweeper read whichever spelling the response actually uses.
func expressConnectGetField(item map[string]interface{}, names ...string) interface{} {
	for _, name := range names {
		target := strings.ToLower(name)
		for k, v := range item {
			if v == nil {
				continue
			}
			if strings.ToLower(k) == target {
				return v
			}
		}
	}
	return nil
}

// expressConnectRpcPost wraps client.RpcPost against the Vpc 2016-04-28 product
// with the standard retry pattern (NeedRetry + incrementalWait). retryCodes are
// additional expected error codes that should be retried.
func expressConnectRpcPost(client *connectivity.AliyunClient, action string, query, body map[string]interface{}, retryCodes []string) (map[string]interface{}, error) {
	return expressConnectRpcPostProduct(client, "Vpc", "2016-04-28", action, query, body, retryCodes)
}

// expressConnectEcrRpcPost is the ExpressConnectRouter 2023-09-01 counterpart.
func expressConnectEcrRpcPost(client *connectivity.AliyunClient, action string, query, body map[string]interface{}, retryCodes []string) (map[string]interface{}, error) {
	return expressConnectRpcPostProduct(client, "ExpressConnectRouter", "2023-09-01", action, query, body, retryCodes)
}

func expressConnectRpcPostProduct(client *connectivity.AliyunClient, product, version, action string, query, body map[string]interface{}, retryCodes []string) (map[string]interface{}, error) {
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	err := retry.Retry(1*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost(product, version, action, query, body, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, retryCodes) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		return nil
	})
	return response, err
}

// expressConnectListEcrIds lists every Express Connect Router whose Name passes
// the fleet-wipe gate. Child sweepers (grant/tr/vpc/vbr associations) walk the
// children under each returned EcrId.
func expressConnectListEcrIds(client *connectivity.AliyunClient) []string {
	var ids []string
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeLarge,
	}
	for {
		request["ClientToken"] = buildClientToken("DescribeExpressConnectRouter")
		response, err := expressConnectEcrRpcPost(client, "DescribeExpressConnectRouter", nil, request, nil)
		if err != nil {
			log.Printf("[ERROR] DescribeExpressConnectRouter got an error: %v", err)
			return ids
		}
		raw, err := jsonpath.Get("$.EcrList", response)
		if err != nil {
			log.Printf("[ERROR] parsing DescribeExpressConnectRouter response got an error: %v", err)
			return ids
		}
		items, _ := raw.([]interface{})
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := fmt.Sprint(expressConnectGetField(item, "Name"))
			ecrId := fmt.Sprint(item["EcrId"])
			if !expressConnectSweepNameOk(name) {
				log.Printf("[INFO] Skipping ExpressConnectRouter: %s (%s)", name, ecrId)
				continue
			}
			ids = append(ids, ecrId)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return ids
}

// expressConnectListQosIds lists every Express Connect Traffic QoS whose QosName
// passes the fleet-wipe gate. Child sweepers (queues and rules) walk the children
// under each returned QosId.
func expressConnectListQosIds(client *connectivity.AliyunClient) []string {
	var ids []string
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeLarge,
	}
	for {
		request["ClientToken"] = buildClientToken("DescribeExpressConnectTrafficQos")
		response, err := expressConnectRpcPost(client, "DescribeExpressConnectTrafficQos", nil, request, nil)
		if err != nil {
			log.Printf("[ERROR] DescribeExpressConnectTrafficQos got an error: %v", err)
			return ids
		}
		raw, err := jsonpath.Get("$.QosList", response)
		if err != nil {
			log.Printf("[ERROR] parsing DescribeExpressConnectTrafficQos response got an error: %v", err)
			return ids
		}
		items, _ := raw.([]interface{})
		for _, v := range items {
			item, _ := v.(map[string]interface{})
			name := fmt.Sprint(expressConnectGetField(item, "QosName", "Name"))
			qosId := fmt.Sprint(item["QosId"])
			if !expressConnectSweepNameOk(name) {
				log.Printf("[INFO] Skipping ExpressConnectTrafficQos: %s (%s)", name, qosId)
				continue
			}
			ids = append(ids, qosId)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return ids
}

// expressConnectListVbrItems lists every Virtual Border Router in the region.
// Callers apply the fleet-wipe gate to each item's Name; the full item map is
// returned because vbr_pconn_association reads the nested
// AssociatedPhysicalConnections from each VBR item.
func expressConnectListVbrItems(client *connectivity.AliyunClient) []map[string]interface{} {
	var items []map[string]interface{}
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := expressConnectRpcPost(client, "DescribeVirtualBorderRouters", nil, request, nil)
		if err != nil {
			log.Printf("[ERROR] DescribeVirtualBorderRouters got an error: %v", err)
			return items
		}
		raw, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		if err != nil {
			log.Printf("[ERROR] parsing DescribeVirtualBorderRouters response got an error: %v", err)
			return items
		}
		result, _ := raw.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			items = append(items, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return items
}

// testSweepExpressConnectGrantRuleToCens revokes CEN grants whose granted
// instance is a tf-prefixed (or SWEEP_ALL) Virtual Border Router. The list API
// DescribeGrantRulesToCen requires an InstanceId (the VBR), so the sweeper walks
// every tf-prefixed VBR and revokes each grant found under it.
func testSweepExpressConnectGrantRuleToCens(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, vbr := range expressConnectListVbrItems(client) {
		name := fmt.Sprint(expressConnectGetField(vbr, "Name"))
		if !expressConnectSweepNameOk(name) {
			continue
		}
		vbrId := fmt.Sprint(vbr["VbrId"])
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"InstanceId":   vbrId,
			"InstanceType": "VBR",
			"PageSize":     PageSizeLarge,
			"PageNumber":   1,
		}
		for {
			request["ClientToken"] = buildClientToken("DescribeGrantRulesToCen")
			response, err := expressConnectRpcPost(client, "DescribeGrantRulesToCen", nil, request, nil)
			if err != nil {
				log.Printf("[ERROR] DescribeGrantRulesToCen for VBR %s got an error: %v", vbrId, err)
				break
			}
			raw, err := jsonpath.Get("$.CenGrantRules.CbnGrantRule", response)
			if err != nil {
				log.Printf("[ERROR] parsing DescribeGrantRulesToCen response got an error: %v", err)
				break
			}
			result, _ := raw.([]interface{})
			for _, v := range result {
				item, _ := v.(map[string]interface{})
				cenId := fmt.Sprint(expressConnectGetField(item, "CenInstanceId", "CenId"))
				cenOwnerId := fmt.Sprint(expressConnectGetField(item, "CenOwnerId"))
				delReq := map[string]interface{}{
					"CenId":        cenId,
					"CenOwnerId":   cenOwnerId,
					"InstanceId":   vbrId,
					"InstanceType": "VBR",
					"RegionId":     client.RegionId,
					"ClientToken":  buildClientToken("RevokeInstanceFromCen"),
				}
				_, err := expressConnectRpcPost(client, "RevokeInstanceFromCen", nil, delReq, nil)
				if err != nil {
					if NotFoundError(err) {
						log.Printf("[INFO] express_connect grant_rule_to_cen %s:%s:%s already gone", cenId, cenOwnerId, vbrId)
					} else {
						log.Printf("[ERROR] RevokeInstanceFromCen %s:%s:%s got an error: %v", cenId, cenOwnerId, vbrId, err)
					}
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

// testSweepExpressConnectVbrPconnAssociations unassociates physical connections
// from tf-prefixed (or SWEEP_ALL) Virtual Border Routers. The association is
// nested under each VBR item as AssociatedPhysicalConnections.
func testSweepExpressConnectVbrPconnAssociations(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, vbr := range expressConnectListVbrItems(client) {
		name := fmt.Sprint(expressConnectGetField(vbr, "Name"))
		if !expressConnectSweepNameOk(name) {
			continue
		}
		vbrId := fmt.Sprint(vbr["VbrId"])
		conns := expressConnectVbrAssociatedPconns(vbr)
		for _, c := range conns {
			conn, _ := c.(map[string]interface{})
			pconnId := fmt.Sprint(expressConnectGetField(conn, "PhysicalConnectionId"))
			if pconnId == "" {
				continue
			}
			delReq := map[string]interface{}{
				"VbrId":                vbrId,
				"PhysicalConnectionId": pconnId,
				"RegionId":             client.RegionId,
				"ClientToken":          buildClientToken("UnassociatePhysicalConnectionFromVirtualBorderRouter"),
			}
			_, err := expressConnectRpcPost(client, "UnassociatePhysicalConnectionFromVirtualBorderRouter", nil, delReq, nil)
			if err != nil {
				if NotFoundError(err) {
					log.Printf("[INFO] express_connect vbr_pconn_association %s:%s already gone", vbrId, pconnId)
				} else {
					log.Printf("[ERROR] UnassociatePhysicalConnectionFromVirtualBorderRouter %s:%s got an error: %v", vbrId, pconnId, err)
				}
			}
		}
	}
	return nil
}

// expressConnectVbrAssociatedPconns returns the AssociatedPhysicalConnection
// array nested under a DescribeVirtualBorderRouters item, or nil when the item
// does not carry any associations.
func expressConnectVbrAssociatedPconns(vbr map[string]interface{}) []interface{} {
	raw, ok := vbr["AssociatedPhysicalConnections"]
	if !ok || raw == nil {
		return nil
	}
	group, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	list, ok := group["AssociatedPhysicalConnection"].([]interface{})
	if !ok {
		return nil
	}
	return list
}

// testSweepExpressConnectRouterExpressConnectRouters deletes Express Connect
// Routers whose Name passes the fleet-wipe gate. Child associations are detached
// first via the dependency-declared child sweepers.
func testSweepExpressConnectRouterExpressConnectRouters(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, ecrId := range expressConnectListEcrIds(client) {
		delReq := map[string]interface{}{
			"EcrId":       ecrId,
			"ClientToken": buildClientToken("DeleteExpressConnectRouter"),
		}
		_, err := expressConnectEcrRpcPost(client, "DeleteExpressConnectRouter", nil, delReq, []string{"Conflict.Lock"})
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"ResourceNotFound.EcrId"}) {
				log.Printf("[INFO] express_connect_router_express_connect_router %s already gone", ecrId)
			} else {
				log.Printf("[ERROR] DeleteExpressConnectRouter %s got an error: %v", ecrId, err)
			}
		}
	}
	return nil
}

// testSweepExpressConnectRouterGrantAssociations revokes instance grants from
// tf-prefixed (or SWEEP_ALL) Express Connect Routers. The list API
// DescribeInstanceGrantedToExpressConnectRouter requires an EcrId, so the
// sweeper walks every tf-prefixed ECR and revokes each grant found under it.
// The list response spells the granted instance fields NodeId/NodeRegionId/
// NodeType while the revoke API expects InstanceId/InstanceRegionId/
// InstanceType; expressConnectGetField smooths that over.
func testSweepExpressConnectRouterGrantAssociations(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, ecrId := range expressConnectListEcrIds(client) {
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"EcrId":      ecrId,
			"MaxResults": PageSizeLarge,
		}
		for {
			request["ClientToken"] = buildClientToken("DescribeInstanceGrantedToExpressConnectRouter")
			response, err := expressConnectEcrRpcPost(client, "DescribeInstanceGrantedToExpressConnectRouter", nil, request, nil)
			if err != nil {
				log.Printf("[ERROR] DescribeInstanceGrantedToExpressConnectRouter for ECR %s got an error: %v", ecrId, err)
				break
			}
			raw, err := jsonpath.Get("$.EcrGrantedInstanceList", response)
			if err != nil {
				log.Printf("[ERROR] parsing DescribeInstanceGrantedToExpressConnectRouter response got an error: %v", err)
				break
			}
			result, _ := raw.([]interface{})
			for _, v := range result {
				item, _ := v.(map[string]interface{})
				delReq := map[string]interface{}{
					"EcrId":            ecrId,
					"InstanceId":       expressConnectGetField(item, "InstanceId", "NodeId"),
					"InstanceRegionId": expressConnectGetField(item, "InstanceRegionId", "NodeRegionId"),
					"EcrOwnerAliUid":   expressConnectGetField(item, "EcrOwnerAliUid"),
					"InstanceType":     expressConnectGetField(item, "InstanceType", "NodeType"),
					"ClientToken":      buildClientToken("RevokeInstanceFromExpressConnectRouter"),
				}
				_, err := expressConnectEcrRpcPost(client, "RevokeInstanceFromExpressConnectRouter", nil, delReq, []string{"Conflict.Lock"})
				if err != nil {
					if NotFoundError(err) || IsExpectedErrors(err, []string{"ResourceNotFound.EcrId"}) {
						log.Printf("[INFO] express_connect_router_grant_association under ECR %s already gone", ecrId)
					} else {
						log.Printf("[ERROR] RevokeInstanceFromExpressConnectRouter under ECR %s got an error: %v", ecrId, err)
					}
				}
			}
			if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
				request["NextToken"] = nextToken
			} else {
				break
			}
		}
	}
	return nil
}

// testSweepExpressConnectRouterTrAssociations deletes Transit Router
// associations from tf-prefixed (or SWEEP_ALL) Express Connect Routers. The list
// API DescribeExpressConnectRouterAssociation requires an EcrId and is scoped to
// TR via AssociationNodeType.
func testSweepExpressConnectRouterTrAssociations(region string) error {
	return expressConnectSweepRouterAssociations(region, "TR")
}

// testSweepExpressConnectRouterVpcAssociations deletes VPC associations from
// tf-prefixed (or SWEEP_ALL) Express Connect Routers. Same list API as TR
// associations but scoped to VPC via AssociationNodeType.
func testSweepExpressConnectRouterVpcAssociations(region string) error {
	return expressConnectSweepRouterAssociations(region, "VPC")
}

// expressConnectSweepRouterAssociations is the shared body for the TR and VPC
// association sweepers. nodeType is "TR" or "VPC".
func expressConnectSweepRouterAssociations(region, nodeType string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, ecrId := range expressConnectListEcrIds(client) {
		request := map[string]interface{}{
			"RegionId":            client.RegionId,
			"EcrId":               ecrId,
			"AssociationNodeType": nodeType,
			"MaxResults":          PageSizeLarge,
		}
		for {
			request["ClientToken"] = buildClientToken("DescribeExpressConnectRouterAssociation")
			response, err := expressConnectEcrRpcPost(client, "DescribeExpressConnectRouterAssociation", nil, request, nil)
			if err != nil {
				log.Printf("[ERROR] DescribeExpressConnectRouterAssociation(%s) for ECR %s got an error: %v", nodeType, ecrId, err)
				break
			}
			raw, err := jsonpath.Get("$.AssociationList", response)
			if err != nil {
				log.Printf("[ERROR] parsing DescribeExpressConnectRouterAssociation(%s) response got an error: %v", nodeType, err)
				break
			}
			result, _ := raw.([]interface{})
			for _, v := range result {
				item, _ := v.(map[string]interface{})
				associationId := fmt.Sprint(item["AssociationId"])
				if associationId == "" {
					continue
				}
				delReq := map[string]interface{}{
					"EcrId":         ecrId,
					"AssociationId": associationId,
					"ClientToken":   buildClientToken("DeleteExpressConnectRouterAssociation"),
				}
				_, err := expressConnectEcrRpcPost(client, "DeleteExpressConnectRouterAssociation", nil, delReq, []string{"Conflict.Lock"})
				if err != nil {
					if NotFoundError(err) || IsExpectedErrors(err, []string{"ResourceNotFound.AssociationId", "ResourceNotFound.EcrId"}) {
						log.Printf("[INFO] express_connect_router_%s_association %s under ECR %s already gone", strings.ToLower(nodeType), associationId, ecrId)
					} else {
						log.Printf("[ERROR] DeleteExpressConnectRouterAssociation(%s) %s under ECR %s got an error: %v", nodeType, associationId, ecrId, err)
					}
				}
			}
			if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
				request["NextToken"] = nextToken
			} else {
				break
			}
		}
	}
	return nil
}

// testSweepExpressConnectRouterVbrChildInstances detaches child VBR instances
// from tf-prefixed (or SWEEP_ALL) Express Connect Routers. The list API
// DescribeExpressConnectRouterChildInstance requires an EcrId.
func testSweepExpressConnectRouterVbrChildInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, ecrId := range expressConnectListEcrIds(client) {
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"EcrId":      ecrId,
			"MaxResults": PageSizeLarge,
		}
		for {
			request["ClientToken"] = buildClientToken("DescribeExpressConnectRouterChildInstance")
			response, err := expressConnectEcrRpcPost(client, "DescribeExpressConnectRouterChildInstance", nil, request, nil)
			if err != nil {
				log.Printf("[ERROR] DescribeExpressConnectRouterChildInstance for ECR %s got an error: %v", ecrId, err)
				break
			}
			raw, err := jsonpath.Get("$.ChildInstanceList", response)
			if err != nil {
				log.Printf("[ERROR] parsing DescribeExpressConnectRouterChildInstance response got an error: %v", err)
				break
			}
			result, _ := raw.([]interface{})
			for _, v := range result {
				item, _ := v.(map[string]interface{})
				childInstanceId := fmt.Sprint(expressConnectGetField(item, "ChildInstanceId"))
				childInstanceType := fmt.Sprint(expressConnectGetField(item, "ChildInstanceType"))
				if childInstanceId == "" || childInstanceType == "" {
					continue
				}
				delReq := map[string]interface{}{
					"EcrId":             ecrId,
					"ChildInstanceId":   childInstanceId,
					"ChildInstanceType": childInstanceType,
					"ClientToken":       buildClientToken("DetachExpressConnectRouterChildInstance"),
				}
				_, err := expressConnectEcrRpcPost(client, "DetachExpressConnectRouterChildInstance", nil, delReq, []string{"Conflict.Lock", "OperationDenied.VbrHasRouteEntryInMiddleStatus"})
				if err != nil {
					if NotFoundError(err) || IsExpectedErrors(err, []string{"ResourceNotFound.AssociationId", "ResourceNotFound.ChildInstanceId", "ResourceNotFound.EcrId"}) {
						log.Printf("[INFO] express_connect_router_vbr_child_instance %s:%s under ECR %s already gone", childInstanceId, childInstanceType, ecrId)
					} else {
						log.Printf("[ERROR] DetachExpressConnectRouterChildInstance %s:%s under ECR %s got an error: %v", childInstanceId, childInstanceType, ecrId, err)
					}
				}
			}
			if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
				request["NextToken"] = nextToken
			} else {
				break
			}
		}
	}
	return nil
}

// testSweepExpressConnectTrafficQoses deletes Express Connect Traffic QoS
// objects whose QosName passes the fleet-wipe gate. Child queues and rules are
// deleted first via the dependency-declared child sweepers.
func testSweepExpressConnectTrafficQoses(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, qosId := range expressConnectListQosIds(client) {
		delReq := map[string]interface{}{
			"QosId":       qosId,
			"RegionId":    client.RegionId,
			"ClientToken": buildClientToken("DeleteExpressConnectTrafficQos"),
		}
		_, err := expressConnectRpcPost(client, "DeleteExpressConnectTrafficQos", nil, delReq, nil)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"IllegalParam.QosId"}) {
				log.Printf("[INFO] express_connect_traffic_qos %s already gone", qosId)
			} else {
				log.Printf("[ERROR] DeleteExpressConnectTrafficQos %s got an error: %v", qosId, err)
			}
		}
	}
	return nil
}

// testSweepExpressConnectTrafficQosQueues deletes queues from tf-prefixed (or
// SWEEP_ALL) Express Connect Traffic QoS objects. DescribeExpressConnectTrafficQosQueue
// takes a QosId and returns the full queue list (no pagination).
func testSweepExpressConnectTrafficQosQueues(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, qosId := range expressConnectListQosIds(client) {
		query := map[string]interface{}{
			"QosId": qosId,
		}
		body := map[string]interface{}{
			"RegionId":    client.RegionId,
			"ClientToken": buildClientToken("DescribeExpressConnectTrafficQosQueue"),
		}
		response, err := expressConnectRpcPost(client, "DescribeExpressConnectTrafficQosQueue", query, body, nil)
		if err != nil {
			log.Printf("[ERROR] DescribeExpressConnectTrafficQosQueue for QoS %s got an error: %v", qosId, err)
			continue
		}
		raw, err := jsonpath.Get("$.QueueList", response)
		if err != nil {
			log.Printf("[ERROR] parsing DescribeExpressConnectTrafficQosQueue response got an error: %v", err)
			continue
		}
		result, _ := raw.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			queueId := fmt.Sprint(expressConnectGetField(item, "QueueId"))
			if queueId == "" {
				continue
			}
			delReq := map[string]interface{}{
				"QosId":       qosId,
				"QueueId":     queueId,
				"RegionId":    client.RegionId,
				"ClientToken": buildClientToken("DeleteExpressConnectTrafficQosQueue"),
			}
			_, err := expressConnectRpcPost(client, "DeleteExpressConnectTrafficQosQueue", nil, delReq, nil)
			if err != nil {
				if NotFoundError(err) {
					log.Printf("[INFO] express_connect_traffic_qos_queue %s:%s already gone", qosId, queueId)
				} else {
					log.Printf("[ERROR] DeleteExpressConnectTrafficQosQueue %s:%s got an error: %v", qosId, queueId, err)
				}
			}
		}
	}
	return nil
}

// testSweepExpressConnectTrafficQosRules deletes rules from tf-prefixed (or
// SWEEP_ALL) Express Connect Traffic QoS queues. DescribeExpressConnectTrafficQosRule
// takes a QosId and QueueId (both in the query) and returns the full rule list.
// Rules live under queues, so the queue sweeper runs first via Dependencies.
func testSweepExpressConnectTrafficQosRules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	for _, qosId := range expressConnectListQosIds(client) {
		queues := expressConnectListQosQueues(client, qosId)
		for _, queueId := range queues {
			query := map[string]interface{}{
				"QosId":   qosId,
				"QueueId": queueId,
			}
			body := map[string]interface{}{
				"RegionId":    client.RegionId,
				"ClientToken": buildClientToken("DescribeExpressConnectTrafficQosRule"),
			}
			response, err := expressConnectRpcPost(client, "DescribeExpressConnectTrafficQosRule", query, body, nil)
			if err != nil {
				log.Printf("[ERROR] DescribeExpressConnectTrafficQosRule for QoS %s queue %s got an error: %v", qosId, queueId, err)
				continue
			}
			raw, err := jsonpath.Get("$.RuleList", response)
			if err != nil {
				log.Printf("[ERROR] parsing DescribeExpressConnectTrafficQosRule response got an error: %v", err)
				continue
			}
			result, _ := raw.([]interface{})
			for _, v := range result {
				item, _ := v.(map[string]interface{})
				ruleId := fmt.Sprint(expressConnectGetField(item, "RuleId"))
				if ruleId == "" {
					continue
				}
				delQuery := map[string]interface{}{
					"QosId":   qosId,
					"QueueId": queueId,
					"RuleId":  ruleId,
				}
				delBody := map[string]interface{}{
					"RegionId":    client.RegionId,
					"ClientToken": buildClientToken("DeleteExpressConnectTrafficQosRule"),
				}
				_, err := expressConnectRpcPost(client, "DeleteExpressConnectTrafficQosRule", delQuery, delBody, nil)
				if err != nil {
					if NotFoundError(err) {
						log.Printf("[INFO] express_connect_traffic_qos_rule %s:%s:%s already gone", qosId, queueId, ruleId)
					} else {
						log.Printf("[ERROR] DeleteExpressConnectTrafficQosRule %s:%s:%s got an error: %v", qosId, queueId, ruleId, err)
					}
				}
			}
		}
	}
	return nil
}

// expressConnectListQosQueues lists the QueueId values under a single QoS. It is
// shared by the traffic_qos_rule sweeper (which needs to enumerate each queue's
// rules).
func expressConnectListQosQueues(client *connectivity.AliyunClient, qosId string) []string {
	var ids []string
	query := map[string]interface{}{
		"QosId": qosId,
	}
	body := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("DescribeExpressConnectTrafficQosQueue"),
	}
	response, err := expressConnectRpcPost(client, "DescribeExpressConnectTrafficQosQueue", query, body, nil)
	if err != nil {
		log.Printf("[ERROR] DescribeExpressConnectTrafficQosQueue for QoS %s got an error: %v", qosId, err)
		return ids
	}
	raw, err := jsonpath.Get("$.QueueList", response)
	if err != nil {
		log.Printf("[ERROR] parsing DescribeExpressConnectTrafficQosQueue response got an error: %v", err)
		return ids
	}
	result, _ := raw.([]interface{})
	for _, v := range result {
		item, _ := v.(map[string]interface{})
		ids = append(ids, fmt.Sprint(expressConnectGetField(item, "QueueId")))
	}
	return ids
}

// testSweepExpressConnectVirtualPhysicalConnections deletes virtual physical
// connections whose Name passes the fleet-wipe gate. vbr_pconn_association runs
// first (via Dependencies) so a virtual physical connection is unassociated
// from its VBR before it is deleted.
func testSweepExpressConnectVirtualPhysicalConnections(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeLarge,
	}
	for {
		request["ClientToken"] = buildClientToken("ListVirtualPhysicalConnections")
		response, err := expressConnectRpcPost(client, "ListVirtualPhysicalConnections", nil, request, nil)
		if err != nil {
			log.Printf("[ERROR] ListVirtualPhysicalConnections got an error: %v", err)
			break
		}
		raw, err := jsonpath.Get("$.VirtualPhysicalConnections", response)
		if err != nil {
			log.Printf("[ERROR] parsing ListVirtualPhysicalConnections response got an error: %v", err)
			break
		}
		result, _ := raw.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			name := fmt.Sprint(expressConnectGetField(item, "Name"))
			pconnId := fmt.Sprint(expressConnectGetField(item, "PhysicalConnectionId"))
			if !expressConnectSweepNameOk(name) {
				log.Printf("[INFO] Skipping VirtualPhysicalConnection: %s (%s)", name, pconnId)
				continue
			}
			delReq := map[string]interface{}{
				"PhysicalConnectionId": pconnId,
				"RegionId":             client.RegionId,
				"ClientToken":          buildClientToken("DeletePhysicalConnection"),
			}
			_, err := expressConnectRpcPost(client, "DeletePhysicalConnection", nil, delReq, nil)
			if err != nil {
				if NotFoundError(err) {
					log.Printf("[INFO] express_connect_virtual_physical_connection %s already gone", pconnId)
				} else {
					log.Printf("[ERROR] DeletePhysicalConnection %s got an error: %v", pconnId, err)
				}
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}
