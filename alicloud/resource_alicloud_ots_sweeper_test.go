package alicloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// OTS sub-resource sweepers.
//
// The instance sweeper (resource_alicloud_ots_instance_test.go) already deletes
// tables and tunnels transitively while walking each instance. These dedicated
// sweepers cover the gaps it does NOT close:
//   - search indexes and secondary indexes are never deleted by the instance sweeper;
//   - instance attachments (VPC bindings) are never unbound by the instance sweeper;
//   - tables/tunnels left orphaned on instances that fail ListTable (e.g. VPC-bound
//     OTSAuthFailed while the instance is stuck-deleting) are not reaped.
//
// Each sub-resource sweeper scopes by the parent INSTANCE name (tf-testAcc* test
// prefix) rather than by the sub-resource name, because OTS test tables/tunnels/
// indexes are named testAcc* (no tf- prefix) and OTS tables are data-bearing —
// production sub-resources must never be swept. The parent-instance prefix check
// therefore applies EVEN under SWEEP_ALL for the data-bearing sub-resources, as a
// guard against silent production data loss (cf. the real存-preserve red-line
// principle). NODELETING markers are honored on both the instance and sub-resource
// names via protectedSweepVpcName.

func init() {
	resource.AddTestSweepers("alicloud_ots_table", &resource.Sweeper{
		Name: "alicloud_ots_table",
		F:    testSweepOtsTables,
	})
	resource.AddTestSweepers("alicloud_ots_tunnel", &resource.Sweeper{
		Name: "alicloud_ots_tunnel",
		F:    testSweepOtsTunnels,
	})
	resource.AddTestSweepers("alicloud_ots_search_index", &resource.Sweeper{
		Name: "alicloud_ots_search_index",
		F:    testSweepOtsSearchIndexes,
	})
	resource.AddTestSweepers("alicloud_ots_secondary_index", &resource.Sweeper{
		Name: "alicloud_ots_secondary_index",
		F:    testSweepOtsSecondaryIndexes,
	})
	resource.AddTestSweepers("alicloud_ots_instance_attachment", &resource.Sweeper{
		Name: "alicloud_ots_instance_attachment",
		F:    testSweepOtsInstanceAttachments,
	})
}

var otsSweepInstancePrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf_test_",
	"tf-test-",
	"tftest",
}

// otsInstanceShouldSweep decides whether a sub-resource's parent instance should
// be swept. It honors NODELETING markers and, when not in SWEEP_ALL mode, requires
// a test prefix on the instance name.
func otsInstanceShouldSweep(instanceName string) bool {
	if protectedSweepVpcName(instanceName) {
		log.Printf("[INFO] Skipping OTS sub-resources of NODELETING-protected instance: %s", instanceName)
		return false
	}
	if !sweepAll() {
		for _, prefix := range otsSweepInstancePrefixes {
			if strings.HasPrefix(strings.ToLower(instanceName), strings.ToLower(prefix)) {
				return true
			}
		}
		log.Printf("[INFO] Skipping OTS sub-resources (instance not a test prefix): %s", instanceName)
		return false
	}
	return true
}

// otsSubResourceShouldSweep decides whether an individual data-bearing sub-resource
// (table/tunnel/index) should be swept. It ALWAYS requires the parent instance to
// carry a test prefix (even under SWEEP_ALL) and honors NODELETING markers on the
// sub-resource name.
func otsSubResourceShouldSweep(instanceName, subResourceName string) bool {
	if protectedSweepVpcName(instanceName) {
		log.Printf("[INFO] Skipping OTS sub-resource (NODELETING-protected instance %s): %s", instanceName, subResourceName)
		return false
	}
	for _, prefix := range otsSweepInstancePrefixes {
		if strings.HasPrefix(strings.ToLower(instanceName), strings.ToLower(prefix)) {
			if protectedSweepVpcName(subResourceName) {
				log.Printf("[INFO] Skipping OTS sub-resource (NODELETING-protected name): %s/%s", instanceName, subResourceName)
				return false
			}
			return true
		}
	}
	log.Printf("[INFO] Skipping OTS sub-resource (instance %s not a test prefix): %s/%s", instanceName, instanceName, subResourceName)
	return false
}

// listOtsSweepInstances paginates ListInstance and returns InstanceInfo for the region.
func listOtsSweepInstances(client *connectivity.AliyunClient) ([]ots.InstanceInfo, error) {
	var insts []ots.InstanceInfo
	req := ots.CreateListInstanceRequest()
	req.RegionId = client.RegionId
	req.Method = "GET"
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNum = requests.NewInteger(1)
	for {
		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.ListInstance(req)
		})
		if err != nil {
			return nil, fmt.Errorf("Error retrieving OTS Instances: %s", err)
		}
		resp, _ := raw.(*ots.ListInstanceResponse)
		if resp == nil || len(resp.InstanceInfos.InstanceInfo) < 1 {
			break
		}
		insts = append(insts, resp.InstanceInfos.InstanceInfo...)
		if len(resp.InstanceInfos.InstanceInfo) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(req.PageNum); err != nil {
			return nil, err
		} else {
			req.PageNum = page
		}
	}
	log.Printf("[INFO] listOtsSweepInstances: ListInstance returned %d OTS instance(s) for region %s", len(insts), client.RegionId)
	return insts, nil
}

// listOtsSweepTables lists table names in an instance (best-effort; returns nil on
// error so the caller can continue with other instances).
func listOtsSweepTables(client *connectivity.AliyunClient, instanceName string) []string {
	raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
		return tableStoreClient.ListTable()
	})
	if err != nil {
		log.Printf("[ERROR] List OTS Instance %s table stores got an error: %#v.", instanceName, err)
		return nil
	}
	tables, _ := raw.(*tablestore.ListTableResponse)
	if tables == nil {
		return nil
	}
	return tables.TableNames
}

func testSweepOtsTables(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	insts, err := listOtsSweepInstances(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list ots instances: %s", err)
		return nil
	}
	for _, inst := range insts {
		instanceName := inst.InstanceName
		if !otsInstanceShouldSweep(instanceName) {
			continue
		}
		for _, tableName := range listOtsSweepTables(client, instanceName) {
			if !otsSubResourceShouldSweep(instanceName, tableName) {
				continue
			}
			log.Printf("[INFO] Deleting OTS Table: instance=%s, table=%s", instanceName, tableName)
			req := new(tablestore.DeleteTableRequest)
			req.TableName = tableName
			if _, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
				return tableStoreClient.DeleteTable(req)
			}); err != nil {
				log.Printf("[ERROR] Failed to delete OTS Table (instance=%s, table=%s): %s", instanceName, tableName, err)
			}
		}
	}
	return nil
}

func testSweepOtsTunnels(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	insts, err := listOtsSweepInstances(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list ots instances: %s", err)
		return nil
	}
	for _, inst := range insts {
		instanceName := inst.InstanceName
		if !otsInstanceShouldSweep(instanceName) {
			continue
		}
		for _, tableName := range listOtsSweepTables(client, instanceName) {
			raw, err := client.WithTableStoreTunnelClient(instanceName, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
				return tunnelClient.ListTunnel(&otsTunnel.ListTunnelRequest{
					TableName: tableName,
				})
			})
			if err != nil {
				log.Printf("[ERROR] List OTS Tunnel failed, instance: %s, table: %s, error: %#v", instanceName, tableName, err)
				continue
			}
			tunnels, _ := raw.(*otsTunnel.ListTunnelResponse)
			if tunnels == nil {
				continue
			}
			for _, t := range tunnels.Tunnels {
				if !otsSubResourceShouldSweep(instanceName, t.TunnelName) {
					continue
				}
				log.Printf("[INFO] Deleting OTS Tunnel: instance=%s, table=%s, tunnel=%s", instanceName, tableName, t.TunnelName)
				if _, err := client.WithTableStoreTunnelClient(instanceName, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
					return tunnelClient.DeleteTunnel(&otsTunnel.DeleteTunnelRequest{
						TableName:  tableName,
						TunnelName: t.TunnelName,
					})
				}); err != nil {
					log.Printf("[ERROR] Failed to delete OTS Tunnel (instance=%s, table=%s, tunnel=%s): %s", instanceName, tableName, t.TunnelName, err)
				}
			}
		}
	}
	return nil
}

func testSweepOtsSearchIndexes(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	insts, err := listOtsSweepInstances(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list ots instances: %s", err)
		return nil
	}
	for _, inst := range insts {
		instanceName := inst.InstanceName
		if !otsInstanceShouldSweep(instanceName) {
			continue
		}
		for _, tableName := range listOtsSweepTables(client, instanceName) {
			indexes, err := otsService.ListOtsSearchIndex(instanceName, tableName)
			if err != nil {
				log.Printf("[ERROR] List OTS SearchIndex failed, instance: %s, table: %s, error: %#v", instanceName, tableName, err)
				continue
			}
			for _, idx := range indexes {
				if !otsSubResourceShouldSweep(instanceName, idx.IndexName) {
					continue
				}
				log.Printf("[INFO] Deleting OTS SearchIndex: instance=%s, table=%s, index=%s", instanceName, tableName, idx.IndexName)
				if err := otsService.DeleteSearchIndex(instanceName, tableName, idx.IndexName); err != nil {
					log.Printf("[ERROR] Failed to delete OTS SearchIndex (instance=%s, table=%s, index=%s): %s", instanceName, tableName, idx.IndexName, err)
				}
			}
		}
	}
	return nil
}

func testSweepOtsSecondaryIndexes(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	insts, err := listOtsSweepInstances(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list ots instances: %s", err)
		return nil
	}
	for _, inst := range insts {
		instanceName := inst.InstanceName
		if !otsInstanceShouldSweep(instanceName) {
			continue
		}
		for _, tableName := range listOtsSweepTables(client, instanceName) {
			indexes, err := otsService.ListOtsSecondaryIndex(instanceName, tableName)
			if err != nil {
				log.Printf("[ERROR] List OTS SecondaryIndex failed, instance: %s, table: %s, error: %#v", instanceName, tableName, err)
				continue
			}
			for _, idx := range indexes {
				if !otsSubResourceShouldSweep(instanceName, idx.IndexName) {
					continue
				}
				log.Printf("[INFO] Deleting OTS SecondaryIndex: instance=%s, table=%s, index=%s", instanceName, tableName, idx.IndexName)
				req := &tablestore.DeleteIndexRequest{
					MainTableName: tableName,
					IndexName:     idx.IndexName,
				}
				if _, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
					return tableStoreClient.DeleteIndex(req)
				}); err != nil {
					log.Printf("[ERROR] Failed to delete OTS SecondaryIndex (instance=%s, table=%s, index=%s): %s", instanceName, tableName, idx.IndexName, err)
				}
			}
		}
	}
	return nil
}

func testSweepOtsInstanceAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	insts, err := listOtsSweepInstances(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list ots instances: %s", err)
		return nil
	}
	for _, inst := range insts {
		instanceName := inst.InstanceName
		// Capture each instance's human-readable status (incl. stuck-Deleting
		// state) for the sweep report — OTS has no aliyun CLI plugin, so the
		// sweeper log is the only source of record for instance state.
		if desc, derr := otsService.DescribeOtsInstance(instanceName); derr != nil {
			log.Printf("[INFO] OTS Instance %s status: unknown (describe failed: %s)", instanceName, derr)
		} else {
			log.Printf("[INFO] OTS Instance %s status: %s", instanceName, desc.InstanceStatus)
		}
		if protectedSweepVpcName(instanceName) {
			log.Printf("[INFO] Skipping OTS InstanceAttachment of NODELETING-protected instance: %s", instanceName)
			continue
		}
		if !sweepAll() {
			matches := false
			for _, prefix := range otsSweepInstancePrefixes {
				if strings.HasPrefix(strings.ToLower(instanceName), strings.ToLower(prefix)) {
					matches = true
					break
				}
			}
			if !matches {
				log.Printf("[INFO] Skipping OTS InstanceAttachment (instance not a test prefix): %s", instanceName)
				continue
			}
		}
		// List VPC bindings bound to this instance and unbind each.
		vpcs, verr := otsService.ListOtsInstanceVpc(instanceName)
		if verr != nil {
			// NotFound (no bindings) is the common case — not an error here.
			log.Printf("[INFO] OTS Instance %s has no VPC bindings to sweep: %s", instanceName, verr)
			continue
		}
		for _, vpc := range vpcs {
			log.Printf("[INFO] Deleting OTS InstanceAttachment: instance=%s, vpc_name=%s", instanceName, vpc.InstanceVpcName)
			request := ots.CreateUnbindInstance2VpcRequest()
			request.RegionId = client.RegionId
			request.InstanceName = instanceName
			request.InstanceVpcName = vpc.InstanceVpcName
			if _, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.UnbindInstance2Vpc(request)
			}); err != nil {
				log.Printf("[ERROR] Failed to delete OTS InstanceAttachment (instance=%s, vpc_name=%s): %s", instanceName, vpc.InstanceVpcName, err)
			}
		}
	}
	return nil
}
