package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// PolarDB (polardb / 云数据库PolarDB) gap sweepers.
//
// Product: polardb, API version 2017-08-01, RPC style.
//
// RpcPost list/delete calls reuse the package helper polarDBRpcPost (retry on
// NeedRetry/LockTimeout via polarDBRpcPostRetryable, addDebug with sensitive-
// value redaction for PolarFs actions). The cluster-walking sweepers reuse
// listPolarDBClustersForSweep (DescribeDBClusters, SDK-typed) defined in
// resource_alicloud_polardb_cluster_test.go. endpoint_address mirrors the
// SDK-typed style of the inline alicloud_polardb_endpoint sweeper.
//
// Resources covered by this file (no inline sweeper existed):
//   - alicloud_polardb_global_security_ip_group
//     DescribeGlobalSecurityIPGroup (RegionId) -> $.GlobalSecurityIPGroup[*]
//     DeleteGlobalSecurityIPGroup {GlobalSecurityGroupId, RegionId, GlobalIgName}
//     NotFound / SecAPICallingFailed => success.
//   - alicloud_polardb_application
//     DescribeApplications (RegionId, paginated) -> $.Items.Applications[*]
//     DeleteApplication {ApplicationId}; PayType==Prepaid skipped (not deletable).
//   - alicloud_polardb_application_endpoint
//     walk apps -> DescribeApplicationAttribute {ApplicationId} -> Endpoints[*]
//     DeleteApplicationEndpointAddress {ApplicationId, EndpointId, NetType}
//     Endpoint.NotFound => success.
//   - alicloud_polardb_polar_fs
//     DescribePolarFs (RegionId, PolarFsType=all, paginated) -> Items.PolarFsPaths[*]
//     DeletePolarFs {PolarFsInstanceId}; NotFound => success.
//   - alicloud_polardb_polar_fs_endpoint
//     walk polar_fs -> DescribeDBClusterEndpoints {DBClusterId, PolarFsInstanceId}
//     -> Items[*] (EndpointType=="Custom" only)
//     DeleteDBClusterEndpoint {DBClusterId, DBEndpointId, PolarFsInstanceId}
//     NotFound / OperationDenied.EndpointNotFound => success.
//   - alicloud_polardb_endpoint_address
//     walk clusters -> DescribeDBClusterEndpoints (SDK) -> Items[*].AddressItems[*]
//     (NetType=="Public" only)
//     DeleteDBEndpointAddress {RegionId, DBClusterId, DBEndpointId, NetType=Public}
//     InvalidDBClusterId.NotFound / InvalidCurrentConnectionString.NotFound /
//     AtLeastOneNetTypeExists => success.
//   - alicloud_polardb_zonal_account
//     walk clusters -> DescribeAccountsZonal {DBClusterId} -> Accounts[*]
//     DeleteAccountZonal {DBClusterId, AccountName}; DynamoDB AccountType skipped.
//   - alicloud_polardb_zonal_endpoint
//     walk clusters -> DescribeDBClusterEndpointsZonal {RegionId, DBClusterId}
//     -> Items[*] (EndpointType=="Custom" only)
//     DeleteDBClusterEndpointZonal {RegionId, DBClusterId, DBEndpointId}
//     InvalidDBClusterId.NotFound / InvalidCurrentConnectionString.NotFound /
//     AtLeastOneNetTypeExists => success.
//
// DEFERRED (no inline sweeper, but not sweepable here):
//   - alicloud_polardb_gateway, gateway_consumer, gateway_consumer_group,
//     gateway_cost_rule, gateway_model_api, gateway_model_service:
//     no DescribeGateways (plural) list API; the provider Read uses the singular
//     DescribeGatewayAttribute by id. Child resources are parent-scoped
//     (need GwClusterId) and cannot enumerate parents server-side.
//   - alicloud_polardb_aicluster: no DescribeAIDBClusters (plural) list; only
//     the singular DescribeAIDBClusterAttribute by id.
//   - alicloud_polardb_batch_task: Delete is a no-op (return nil) and there is
//     no list API.
//   - alicloud_polardb_zonal_db_cluster: default PayType is Prepaid which can
//     only be released via BssOpenApi RefundInstance (needs explicit user
//     approval); PostPaid zonal clusters are already swept by the inline
//     alicloud_polardb_cluster sweeper; no dedicated zonal list action.
//   - alicloud_polardb_account_privilege: grant/revoke sub-object; no dedicated
//     list; deletion is per-(account, DB, privilege) RevokeAccountPrivilege.
//   - alicloud_polardb_backup_policy: modify-only; Delete is a no-op.
//   - alicloud_polardb_cluster_endpoint, alicloud_polardb_primary_endpoint:
//     adopt/modify-only built-in system endpoints; Delete is a no-op.
//   - alicloud_polardb_dynamo_item: no list API (items are per-key GetItem; a
//     sweep would need Scan per table and items carry no name prefix).
//   - alicloud_polardb_dynamo_table: deletable via the AWS DynamoDB SDK
//     (ListTables/DeleteTable) but connection resolution is tied to
//     schema.ResourceData and crosses SDKs; tables are cascade-cleaned when the
//     parent cluster is swept by the inline alicloud_polardb_cluster sweeper.
//
// Account 1511928242963727 / cn-beijing only. protectedSweepVpcName (NODELETING
// and the 21 do-not-delete markers) is honored in every mode, SWEEP_ALL
// included. No AK/SK/Token appear in this file.

var polardbSweepPrefixes = []string{"tf-testAcc", "tf_testAcc"}

// polardbSweepPrefixMatch reports whether name begins (case-insensitively)
// with one of the standard test prefixes.
func polardbSweepPrefixMatch(name string) bool {
	for _, prefix := range polardbSweepPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// polardbSweepNameOk reports whether a named resource may be swept: never for
// NODELETING/do-not-delete markers (every mode, SWEEP_ALL included); otherwise
// under SWEEP_ALL everything is swept, and otherwise only test-prefixed names.
func polardbSweepNameOk(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	return polardbSweepPrefixMatch(name)
}

// listPolarFsForSweep pages DescribePolarFs (PolarFsType=all) and returns each
// polar_fs item as a map (PolarFsInstanceId, PolarFsInstanceDescription,
// DBClusterId, RelativeDbClusterId, ...). Mirrors the DescribePolarFs call in
// data_source_alicloud_polardb_polar_fs_instances.go and reuses the package
// helper polarDBPolarFsListItems to unwrap Items.PolarFsPaths.
func listPolarFsForSweep(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	var insts []map[string]interface{}
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"PageSize":    100,
		"PageNumber":  1,
		"PolarFsType": "all",
	}
	for {
		response, err := polarDBRpcPost(client, "DescribePolarFs", request, 5*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving PolarDB PolarFs instances: %s", err)
		}
		items := polarDBPolarFsListItems(response["Items"])
		for _, v := range items {
			if item, ok := v.(map[string]interface{}); ok {
				insts = append(insts, item)
			}
		}
		if len(items) < 100 {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return insts, nil
}

// listPolarDBApplicationsForSweep pages DescribeApplications and returns each
// application item as a map (ApplicationId, Description, PayType, Status, ...).
// Items live at the nested path $.Items.Applications (verified via the OpenAPI
// definition: the response Items object wraps the Applications array).
func listPolarDBApplicationsForSweep(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	var insts []map[string]interface{}
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		response, err := polarDBRpcPost(client, "DescribeApplications", request, 5*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving PolarDB Applications: %s", err)
		}
		raw, err := jsonpath.Get("$.Items.Applications", response)
		if err != nil {
			return nil, fmt.Errorf("Error parsing PolarDB Applications response: %s", err)
		}
		items, _ := raw.([]interface{})
		for _, v := range items {
			if item, ok := v.(map[string]interface{}); ok {
				insts = append(insts, item)
			}
		}
		if len(items) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return insts, nil
}

func init() {
	resource.AddTestSweepers("alicloud_polardb_global_security_ip_group", &resource.Sweeper{
		Name: "alicloud_polardb_global_security_ip_group",
		F:    testSweepPolarDBGlobalSecurityIPGroup,
	})
	resource.AddTestSweepers("alicloud_polardb_application", &resource.Sweeper{
		Name: "alicloud_polardb_application",
		F:    testSweepPolarDBApplication,
	})
	resource.AddTestSweepers("alicloud_polardb_application_endpoint", &resource.Sweeper{
		Name: "alicloud_polardb_application_endpoint",
		F:    testSweepPolarDBApplicationEndpoint,
	})
	resource.AddTestSweepers("alicloud_polardb_polar_fs", &resource.Sweeper{
		Name: "alicloud_polardb_polar_fs",
		F:    testSweepPolarDBPolarFs,
	})
	resource.AddTestSweepers("alicloud_polardb_polar_fs_endpoint", &resource.Sweeper{
		Name: "alicloud_polardb_polar_fs_endpoint",
		F:    testSweepPolarDBPolarFsEndpoint,
	})
	resource.AddTestSweepers("alicloud_polardb_endpoint_address", &resource.Sweeper{
		Name: "alicloud_polardb_endpoint_address",
		F:    testSweepPolarDBEndpointAddress,
	})
	resource.AddTestSweepers("alicloud_polardb_zonal_account", &resource.Sweeper{
		Name: "alicloud_polardb_zonal_account",
		F:    testSweepPolarDBZonalAccount,
	})
	resource.AddTestSweepers("alicloud_polardb_zonal_endpoint", &resource.Sweeper{
		Name: "alicloud_polardb_zonal_endpoint",
		F:    testSweepPolarDBZonalEndpoint,
	})
}

// testSweepPolarDBGlobalSecurityIPGroup deletes every PolarDB
// global-security-IP-group whose GlobalIgName passes the sweep gate. Only the
// one list call is region-scoped; deletion needs both the id and the name.
func testSweepPolarDBGlobalSecurityIPGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	request := map[string]interface{}{"RegionId": client.RegionId}
	response, err := polarDBRpcPost(client, "DescribeGlobalSecurityIPGroup", request, 5*time.Minute)
	if err != nil {
		log.Printf("[ERROR] Failed to list PolarDB GlobalSecurityIPGroup: %s", err)
		return nil
	}
	raw, err := jsonpath.Get("$.GlobalSecurityIPGroup", response)
	if err != nil {
		log.Printf("[ERROR] PolarDB GlobalSecurityIPGroup: no list in response (treat as empty): %s", err)
		return nil
	}
	items, _ := raw.([]interface{})

	for _, v := range items {
		item, _ := v.(map[string]interface{})
		name := fmt.Sprint(item["GlobalIgName"])
		if !polardbSweepNameOk(name) {
			log.Printf("[INFO] Skipping PolarDB GlobalSecurityIPGroup (NODELETING/no-prefix): %s", name)
			continue
		}
		id := fmt.Sprint(item["GlobalSecurityGroupId"])
		log.Printf("[INFO] Deleting PolarDB GlobalSecurityIPGroup: %s (%s)", name, id)
		_, err := polarDBRpcPost(client, "DeleteGlobalSecurityIPGroup", map[string]interface{}{
			"GlobalSecurityGroupId": id,
			"RegionId":              client.RegionId,
			"GlobalIgName":          name,
		}, 3*time.Minute)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"SecAPICallingFailed"}) {
				continue
			}
			log.Printf("[ERROR] Failed to delete PolarDB GlobalSecurityIPGroup %s: %s", id, err)
		}
	}
	return nil
}

// testSweepPolarDBApplication deletes every PolarDB Application whose
// Description passes the sweep gate. Prepaid applications are skipped because
// the provider's Delete refuses them (DeleteApplication on a Prepaid app
// returns an error); they would need a BSS refund, which is out of scope here.
func testSweepPolarDBApplication(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	apps, err := listPolarDBApplicationsForSweep(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list PolarDB Applications: %s", err)
		return nil
	}

	for _, app := range apps {
		name := fmt.Sprint(app["Description"])
		id := fmt.Sprint(app["ApplicationId"])
		if !polardbSweepNameOk(name) {
			log.Printf("[INFO] Skipping PolarDB Application (NODELETING/no-prefix): %s (%s)", name, id)
			continue
		}
		if strings.EqualFold(fmt.Sprint(app["PayType"]), "Prepaid") {
			log.Printf("[INFO] Skipping PolarDB Application (Prepaid, not deletable): %s (%s)", name, id)
			continue
		}
		if fmt.Sprint(app["Status"]) == "Deleting" {
			continue
		}
		log.Printf("[INFO] Deleting PolarDB Application: %s (%s)", name, id)
		_, err := polarDBRpcPost(client, "DeleteApplication", map[string]interface{}{
			"ApplicationId": id,
		}, 3*time.Minute)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			log.Printf("[ERROR] Failed to delete PolarDB Application %s: %s", id, err)
		}
	}
	return nil
}

// testSweepPolarDBApplicationEndpoint deletes the endpoints of every PolarDB
// Application whose Description passes the sweep gate. Endpoints are listed via
// DescribeApplicationAttribute (singular) which returns the Endpoints array at
// the top level; each endpoint carries EndpointId and NetType, both required
// by DeleteApplicationEndpointAddress. Application PayType does not block
// endpoint deletion, so endpoints on Prepaid apps are swept too.
func testSweepPolarDBApplicationEndpoint(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	apps, err := listPolarDBApplicationsForSweep(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list PolarDB Applications: %s", err)
		return nil
	}

	for _, app := range apps {
		appName := fmt.Sprint(app["Description"])
		appId := fmt.Sprint(app["ApplicationId"])
		if !polardbSweepNameOk(appName) {
			continue
		}
		resp, err := polarDBRpcPost(client, "DescribeApplicationAttribute", map[string]interface{}{
			"ApplicationId": appId,
		}, 5*time.Minute)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			log.Printf("[ERROR] Failed to describe PolarDB Application %s endpoints: %s", appId, err)
			continue
		}
		rawEps, _ := resp["Endpoints"].([]interface{})
		for _, e := range rawEps {
			ep, _ := e.(map[string]interface{})
			endpointId := fmt.Sprint(ep["EndpointId"])
			netType := fmt.Sprint(ep["NetType"])
			log.Printf("[INFO] Deleting PolarDB Application Endpoint: %s (%s) on %s", endpointId, netType, appId)
			_, err := polarDBRpcPost(client, "DeleteApplicationEndpointAddress", map[string]interface{}{
				"ApplicationId": appId,
				"EndpointId":    endpointId,
				"NetType":       netType,
			}, 3*time.Minute)
			if err != nil {
				if NotFoundError(err) || IsExpectedErrors(err, []string{"Endpoint.NotFound"}) {
					continue
				}
				log.Printf("[ERROR] Failed to delete PolarDB Application Endpoint %s on %s: %s", endpointId, appId, err)
			}
		}
	}
	return nil
}

// testSweepPolarDBPolarFs deletes every PolarDB PolarFs whose
// PolarFsInstanceDescription passes the sweep gate. The list unwraps
// Items.PolarFsPaths (handled by polarDBPolarFsListItems); each item carries
// PolarFsInstanceId. DeletePolarFs only needs the id.
func testSweepPolarDBPolarFs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	insts, err := listPolarFsForSweep(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list PolarDB PolarFs: %s", err)
		return nil
	}

	for _, inst := range insts {
		name := fmt.Sprint(inst["PolarFsInstanceDescription"])
		id := fmt.Sprint(inst["PolarFsInstanceId"])
		if !polardbSweepNameOk(name) {
			log.Printf("[INFO] Skipping PolarDB PolarFs (NODELETING/no-prefix): %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting PolarDB PolarFs: %s (%s)", name, id)
		_, err := polarDBRpcPost(client, "DeletePolarFs", map[string]interface{}{
			"PolarFsInstanceId": id,
		}, 3*time.Minute)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			log.Printf("[ERROR] Failed to delete PolarDB PolarFs %s: %s", id, err)
		}
	}
	return nil
}

// testSweepPolarDBPolarFsEndpoint deletes the Custom endpoints of every PolarDB
// PolarFs whose PolarFsInstanceDescription passes the sweep gate. Endpoints
// live under the parent polar_fs, so the sweep lists polar_fs (to get both
// PolarFsInstanceId and DBClusterId) then calls DescribeDBClusterEndpoints with
// both ids; only Custom endpoints are deletable. DBClusterId is read from the
// polar_fs item (falling back to RelativeDbClusterId) because the describe
// helper requires it.
func testSweepPolarDBPolarFsEndpoint(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	insts, err := listPolarFsForSweep(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list PolarDB PolarFs: %s", err)
		return nil
	}

	for _, inst := range insts {
		pfsName := fmt.Sprint(inst["PolarFsInstanceDescription"])
		polarFsInstanceId := fmt.Sprint(inst["PolarFsInstanceId"])
		if !polardbSweepNameOk(pfsName) {
			log.Printf("[INFO] Skipping PolarDB PolarFs endpoints (NODELETING/no-prefix polar_fs): %s (%s)", pfsName, polarFsInstanceId)
			continue
		}
		dbClusterId := fmt.Sprint(inst["DBClusterId"])
		if dbClusterId == "" || dbClusterId == "<nil>" {
			dbClusterId = fmt.Sprint(inst["RelativeDbClusterId"])
		}
		if dbClusterId == "" || dbClusterId == "<nil>" {
			continue
		}
		resp, err := polarDBRpcPost(client, "DescribeDBClusterEndpoints", map[string]interface{}{
			"DBClusterId":       dbClusterId,
			"PolarFsInstanceId": polarFsInstanceId,
		}, 5*time.Minute)
		if err != nil {
			log.Printf("[ERROR] Failed to list PolarDB PolarFs endpoints on %s/%s: %s", dbClusterId, polarFsInstanceId, err)
			continue
		}
		rawEps, _ := resp["Items"].([]interface{})
		for _, e := range rawEps {
			ep, _ := e.(map[string]interface{})
			if fmt.Sprint(ep["EndpointType"]) != "Custom" {
				continue
			}
			desc := fmt.Sprint(ep["DBEndpointDescription"])
			if !polardbSweepNameOk(desc) {
				continue
			}
			endpointId := fmt.Sprint(ep["DBEndpointId"])
			log.Printf("[INFO] Deleting PolarDB PolarFs Endpoint: %s on %s", endpointId, polarFsInstanceId)
			_, err := polarDBRpcPost(client, "DeleteDBClusterEndpoint", map[string]interface{}{
				"DBClusterId":       dbClusterId,
				"DBEndpointId":      endpointId,
				"PolarFsInstanceId": polarFsInstanceId,
			}, 3*time.Minute)
			if err != nil {
				if NotFoundError(err) || IsExpectedErrors(err, []string{"OperationDenied.EndpointNotFound"}) {
					continue
				}
				log.Printf("[ERROR] Failed to delete PolarDB PolarFs Endpoint %s on %s: %s", endpointId, polarFsInstanceId, err)
			}
		}
	}
	return nil
}

// testSweepPolarDBEndpointAddress deletes the Public address allocation of
// every PolarDB cluster endpoint whose cluster description passes the sweep
// gate. It mirrors the SDK-typed style of the inline alicloud_polardb_endpoint
// sweeper: DescribeDBClusterEndpoints (SDK) -> Items[].AddressItems[]; only
// addresses with NetType=="Public" are deletable (Private/Inner addresses are
// managed by the cluster). The child gate uses the connection prefix (the
// label portion of ConnectionString before the first dot); an address with no
// ConnectionString is swept only under SWEEP_ALL.
func testSweepPolarDBEndpointAddress(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters, err := listPolarDBClustersForSweep(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list PolarDB clusters: %s", err)
		return nil
	}

	for _, c := range clusters {
		clusterName := c.DBClusterDescription
		clusterId := c.DBClusterId
		if protectedSweepVpcName(clusterName) {
			log.Printf("[INFO] Skipping PolarDB EndpointAddress on NODELETING cluster: %s (%s)", clusterName, clusterId)
			continue
		}
		if !sweepAll() && !polardbSweepPrefixMatch(clusterName) {
			continue
		}

		ereq := polardb.CreateDescribeDBClusterEndpointsRequest()
		ereq.RegionId = client.RegionId
		ereq.DBClusterId = clusterId
		raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeDBClusterEndpoints(ereq)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list PolarDB endpoints on cluster %s: %s", clusterId, err)
			continue
		}
		eresp, _ := raw.(*polardb.DescribeDBClusterEndpointsResponse)
		if eresp == nil {
			continue
		}
		for _, ep := range eresp.Items {
			for _, addr := range ep.AddressItems {
				if addr.NetType != "Public" {
					continue
				}
				connPrefix := ""
				if cs := addr.ConnectionString; cs != "" {
					if parts := strings.Split(cs, "."); len(parts) > 0 {
						connPrefix = parts[0]
					}
				}
				if protectedSweepVpcName(connPrefix) {
					continue
				}
				if !sweepAll() && !polardbSweepPrefixMatch(connPrefix) {
					continue
				}
				log.Printf("[INFO] Deleting PolarDB EndpointAddress (Public) %s on %s", ep.DBEndpointId, clusterId)
				dreq := polardb.CreateDeleteDBEndpointAddressRequest()
				dreq.RegionId = client.RegionId
				dreq.DBClusterId = clusterId
				dreq.DBEndpointId = ep.DBEndpointId
				dreq.NetType = "Public"
				if _, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
					return polardbClient.DeleteDBEndpointAddress(dreq)
				}); err != nil {
					if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) || NotFoundError(err) {
						continue
					}
					log.Printf("[ERROR] Failed to delete PolarDB EndpointAddress %s on %s: %s", ep.DBEndpointId, clusterId, err)
				}
			}
		}
	}
	return nil
}

// testSweepPolarDBZonalAccount deletes the zonal accounts of every PolarDB
// cluster whose description passes the sweep gate. Zonal accounts are listed
// via DescribeAccountsZonal {DBClusterId} (omitting AccountName lists all) and
// deleted via DeleteAccountZonal {DBClusterId, AccountName}. DynamoDB-type
// accounts are skipped (DeleteAccountZonal does not handle them, matching the
// inline alicloud_polardb_account sweeper). A non-zonal cluster errors on the
// list call; that is logged and the sweep continues.
func testSweepPolarDBZonalAccount(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters, err := listPolarDBClustersForSweep(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list PolarDB clusters: %s", err)
		return nil
	}

	for _, c := range clusters {
		clusterName := c.DBClusterDescription
		clusterId := c.DBClusterId
		if protectedSweepVpcName(clusterName) {
			log.Printf("[INFO] Skipping PolarDB zonal accounts on NODELETING cluster: %s (%s)", clusterName, clusterId)
			continue
		}
		if !sweepAll() && !polardbSweepPrefixMatch(clusterName) {
			continue
		}
		resp, err := polarDBRpcPost(client, "DescribeAccountsZonal", map[string]interface{}{
			"DBClusterId": clusterId,
		}, 5*time.Minute)
		if err != nil {
			log.Printf("[ERROR] Failed to list PolarDB zonal accounts on %s: %s", clusterId, err)
			continue
		}
		rawAccs, _ := resp["Accounts"].([]interface{})
		for _, a := range rawAccs {
			acc, _ := a.(map[string]interface{})
			if fmt.Sprint(acc["AccountType"]) == "DynamoDB" {
				continue
			}
			name := fmt.Sprint(acc["AccountName"])
			if !polardbSweepNameOk(name) {
				continue
			}
			log.Printf("[INFO] Deleting PolarDB zonal account: %s on %s", name, clusterId)
			_, err := polarDBRpcPost(client, "DeleteAccountZonal", map[string]interface{}{
				"DBClusterId": clusterId,
				"AccountName": name,
			}, 3*time.Minute)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				log.Printf("[ERROR] Failed to delete PolarDB zonal account %s on %s: %s", name, clusterId, err)
			}
		}
	}
	return nil
}

// testSweepPolarDBZonalEndpoint deletes the Custom zonal endpoints of every
// PolarDB cluster whose description passes the sweep gate. Zonal endpoints are
// listed via DescribeDBClusterEndpointsZonal {RegionId, DBClusterId} and
// deleted via DeleteDBClusterEndpointZonal {RegionId, DBClusterId, DBEndpointId};
// only Custom endpoints are deletable (Primary/Cluster are managed by the
// cluster). A non-zonal cluster errors on the list call; that is logged and the
// sweep continues.
func testSweepPolarDBZonalEndpoint(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters, err := listPolarDBClustersForSweep(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list PolarDB clusters: %s", err)
		return nil
	}

	for _, c := range clusters {
		clusterName := c.DBClusterDescription
		clusterId := c.DBClusterId
		if protectedSweepVpcName(clusterName) {
			log.Printf("[INFO] Skipping PolarDB zonal endpoints on NODELETING cluster: %s (%s)", clusterName, clusterId)
			continue
		}
		if !sweepAll() && !polardbSweepPrefixMatch(clusterName) {
			continue
		}
		resp, err := polarDBRpcPost(client, "DescribeDBClusterEndpointsZonal", map[string]interface{}{
			"RegionId":    client.RegionId,
			"DBClusterId": clusterId,
		}, 5*time.Minute)
		if err != nil {
			log.Printf("[ERROR] Failed to list PolarDB zonal endpoints on %s: %s", clusterId, err)
			continue
		}
		rawEps, _ := resp["Items"].([]interface{})
		for _, e := range rawEps {
			ep, _ := e.(map[string]interface{})
			if fmt.Sprint(ep["EndpointType"]) != "Custom" {
				continue
			}
			desc := fmt.Sprint(ep["DBEndpointDescription"])
			if !polardbSweepNameOk(desc) {
				continue
			}
			endpointId := fmt.Sprint(ep["DBEndpointId"])
			log.Printf("[INFO] Deleting PolarDB zonal endpoint: %s on %s", endpointId, clusterId)
			_, err := polarDBRpcPost(client, "DeleteDBClusterEndpointZonal", map[string]interface{}{
				"RegionId":     client.RegionId,
				"DBClusterId":  clusterId,
				"DBEndpointId": endpointId,
			}, 3*time.Minute)
			if err != nil {
				if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
					continue
				}
				log.Printf("[ERROR] Failed to delete PolarDB zonal endpoint %s on %s: %s", endpointId, clusterId, err)
			}
		}
	}
	return nil
}
