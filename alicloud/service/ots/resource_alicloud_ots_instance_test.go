package ots_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	sdkots "github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ots_instance", &resource.Sweeper{
		Name: "alicloud_ots_instance",
		F:    testSweepOtsInstances,
	})
}

func testSweepOtsInstances(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
		// Documentation examples use this prefix because OTS instance names are
		// limited to 16 bytes, which rules out the full "terraform-example" prefix.
		"tf-example-ots",
	}

	var insts []sdkots.InstanceInfo
	req := sdkots.CreateListInstanceRequest()
	req.RegionId = client.RegionId
	req.Method = "GET"
	req.PageSize = requests.NewInteger(helper.PageSizeLarge)
	req.PageNum = requests.NewInteger(1)
	for {
		raw, err := client.WithOtsClient(func(otsClient *sdkots.Client) (interface{}, error) {
			return otsClient.ListInstance(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving OTS Instances: %s", err)
		}
		resp, _ := raw.(*sdkots.ListInstanceResponse)
		if resp == nil || len(resp.InstanceInfos.InstanceInfo) < 1 {
			break
		}
		insts = append(insts, resp.InstanceInfos.InstanceInfo...)

		if len(resp.InstanceInfos.InstanceInfo) < helper.PageSizeLarge {
			break
		}

		if page, err := helper.GetNextpageNumber(req.PageNum); err != nil {
			return err
		} else {
			req.PageNum = page
		}
	}
	for _, v := range insts {
		name := v.InstanceName
		skip := true
		if !sweep.SweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping OTS Instance: %s", name)
				continue
			}
		}
		log.Printf("[INFO] Deleting OTS Instance %s table stores.", name)
		raw, err := client.WithTableStoreClient(name, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			return tableStoreClient.ListTable()
		})
		if err != nil {
			if isOtsObjectNotExist(err) {
				continue
			}
			return fmt.Errorf("error listing tables of OTS instance %s: %s", name, err)
		}
		tables, _ := raw.(*tablestore.ListTableResponse)
		if tables == nil {
			return fmt.Errorf("empty ListTable response for OTS instance %s", name)
		}
		for _, t := range tables.TableNames {
			// A search index must be deleted before its table, otherwise DeleteTable is rejected.
			if err := sweepSearchIndexesInTable(client, name, t); err != nil {
				return err
			}
			if err := sweepTunnelsInTable(client, name, t); err != nil {
				return err
			}
			req := new(tablestore.DeleteTableRequest)
			req.TableName = t
			if _, err := client.WithTableStoreClient(name, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
				return tableStoreClient.DeleteTable(req)
			}); err != nil && !isOtsObjectNotExist(err) {
				return fmt.Errorf("error deleting table %s of OTS instance %s: %s", t, name, err)
			}
		}
		if len(tables.TableNames) > 0 {
			if err := waitForOtsTablesDeleted(client, name); err != nil {
				return err
			}
		}
		log.Printf("[INFO] Deleting OTS Instance: %s", name)
		req := sdkots.CreateDeleteInstanceRequest()
		req.InstanceName = name
		_, err = client.WithOtsClient(func(otsClient *sdkots.Client) (interface{}, error) {
			return otsClient.DeleteInstance(req)
		})
		if err != nil && !isOtsObjectNotExist(err) {
			return fmt.Errorf("error deleting OTS instance %s: %s", name, err)
		}
	}

	return nil
}

func isOtsObjectNotExist(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "OTSObjectNotExist")
}

func sweepSearchIndexesInTable(client *connectivity.AliyunClient, instanceName string, tableName string) error {
	deadline := time.Now().Add(2 * time.Minute)
	for {
		raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			return tableStoreClient.ListSearchIndex(&tablestore.ListSearchIndexRequest{
				TableName: tableName,
			})
		})
		if err != nil {
			if isOtsObjectNotExist(err) {
				return nil
			}
			return fmt.Errorf("error listing search indexes of table %s in OTS instance %s: %s", tableName, instanceName, err)
		}
		indexes, _ := raw.(*tablestore.ListSearchIndexResponse)
		if indexes == nil {
			return fmt.Errorf("empty ListSearchIndex response for table %s in OTS instance %s", tableName, instanceName)
		}
		if len(indexes.IndexInfo) == 0 {
			return nil
		}
		for _, index := range indexes.IndexInfo {
			if index == nil || index.IndexName == "" {
				continue
			}
			if _, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
				return tableStoreClient.DeleteSearchIndex(&tablestore.DeleteSearchIndexRequest{
					TableName: tableName,
					IndexName: index.IndexName,
				})
			}); err != nil && !isOtsObjectNotExist(err) {
				return fmt.Errorf("error deleting search index %s of table %s in OTS instance %s: %s", index.IndexName, tableName, instanceName, err)
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for search indexes of table %s in OTS instance %s to be deleted", tableName, instanceName)
		}
		time.Sleep(5 * time.Second)
	}
}

func waitForOtsTablesDeleted(client *connectivity.AliyunClient, instanceName string) error {
	deadline := time.Now().Add(2 * time.Minute)
	for {
		raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			return tableStoreClient.ListTable()
		})
		if err != nil {
			if isOtsObjectNotExist(err) {
				return nil
			}
			return fmt.Errorf("error listing tables of OTS instance %s while waiting for deletion: %s", instanceName, err)
		}
		tables, _ := raw.(*tablestore.ListTableResponse)
		if tables == nil || len(tables.TableNames) == 0 {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for tables %v of OTS instance %s to be deleted", tables.TableNames, instanceName)
		}
		time.Sleep(5 * time.Second)
	}
}

func sweepTunnelsInTable(client *connectivity.AliyunClient, instanceName string, tableName string) error {
	log.Printf("[INFO] Deleting OTS Tunnels, instance: %s, table: %s", instanceName, tableName)
	raw, err := client.WithTableStoreTunnelClient(instanceName, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
		return tunnelClient.ListTunnel(&otsTunnel.ListTunnelRequest{
			TableName: tableName,
		})
	})
	if err != nil {
		return fmt.Errorf("error listing tunnels of table %s in OTS instance %s: %s", tableName, instanceName, err)
	}

	tunnels, _ := raw.(*otsTunnel.ListTunnelResponse)
	if tunnels != nil && len(tunnels.Tunnels) > 0 {
		for _, t := range tunnels.Tunnels {
			if _, err := client.WithTableStoreTunnelClient(instanceName, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
				return tunnelClient.DeleteTunnel(&otsTunnel.DeleteTunnelRequest{
					TableName:  tableName,
					TunnelName: t.TunnelName,
				})
			}); err != nil && !isOtsObjectNotExist(err) {
				return fmt.Errorf("error deleting tunnel %s of table %s in OTS instance %s: %s", t.TunnelName, tableName, instanceName, err)
			}
		}
	}
	return nil
}

func TestAccAliCloudOtsInstance_basic(t *testing.T) {
	var v sdkots.InstanceInfo

	resourceId := "alicloud_ots_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, otsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":          name,
					"description":   name,
					"instance_type": "Capacity",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"description":          name,
						"instance_type":        "Capacity",
						"network_type_acl.#":   "3",
						"network_source_acl.#": "1",
						"resource_group_id":    acctest.CHECKSET,
						"accessed_by":          "Any",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessed_by": "Vpc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessed_by":          "Vpc",
						"network_type_acl.#":   acctest.CHECKSET,
						"network_source_acl.#": acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "acceptance test",
						"network_type_acl.#":   acctest.CHECKSET,
						"network_source_acl.#": acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":               "3",
						"tags.Created":         "TF",
						"tags.For":             "acceptance test",
						"tags.Updated":         "TF",
						"network_type_acl.#":   acctest.CHECKSET,
						"network_source_acl.#": acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessed_by": acctest.REMOVEKEY,
					"tags":        acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessed_by":          "Vpc",
						"tags.%":               "0",
						"tags.Created":         acctest.REMOVEKEY,
						"tags.For":             acctest.REMOVEKEY,
						"tags.Updated":         acctest.REMOVEKEY,
						"network_type_acl.#":   acctest.CHECKSET,
						"network_source_acl.#": acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOtsInstance_acl(t *testing.T) {
	var v sdkots.InstanceInfo

	resourceId := "alicloud_ots_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, otsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"description":        name,
					"instance_type":      "Capacity",
					"network_type_acl":   []string{"INTERNET", "VPC"},
					"network_source_acl": []string{"TRUST_PROXY"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"description":          name,
						"instance_type":        "Capacity",
						"network_type_acl.#":   "2",
						"network_source_acl.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name,
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"description":          name,
						"network_type_acl.#":   "2",
						"network_source_acl.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"description":        name,
					"network_type_acl":   []string{"VPC"},
					"network_source_acl": []string{"TRUST_PROXY"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"description":          name,
						"network_type_acl.#":   "1",
						"network_source_acl.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"description":        name,
					"network_type_acl":   []string{"VPC"},
					"network_source_acl": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"description":          name,
						"network_type_acl.#":   "1",
						"network_source_acl.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"description":        name,
					"network_type_acl":   []string{"CLASSIC"},
					"network_source_acl": []string{"TRUST_PROXY"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"description":          name,
						"network_type_acl.#":   "1",
						"network_source_acl.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"description":        name,
					"network_type_acl":   []string{"INTERNET"},
					"network_source_acl": []string{"TRUST_PROXY"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"description":          name,
						"network_type_acl.#":   "1",
						"network_source_acl.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOtsInstanceHighPerformance(t *testing.T) {
	var v sdkots.InstanceInfo

	resourceId := "alicloud_ots_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, otsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":          name,
					"description":   name,
					"instance_type": "HighPerformance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"description":          name,
						"instance_type":        "HighPerformance",
						"network_type_acl.#":   acctest.CHECKSET,
						"network_source_acl.#": acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessed_by": "Vpc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessed_by":          "Vpc",
						"network_type_acl.#":   acctest.CHECKSET,
						"network_source_acl.#": acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "acceptance test",
						"network_type_acl.#":   acctest.CHECKSET,
						"network_source_acl.#": acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":               "3",
						"tags.Created":         "TF",
						"tags.For":             "acceptance test",
						"tags.Updated":         "TF",
						"network_type_acl.#":   acctest.CHECKSET,
						"network_source_acl.#": acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessed_by": acctest.REMOVEKEY,
					"tags":        acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessed_by":          "Vpc",
						"tags.%":               "0",
						"tags.Created":         acctest.REMOVEKEY,
						"tags.For":             acctest.REMOVEKEY,
						"tags.Updated":         acctest.REMOVEKEY,
						"network_type_acl.#":   acctest.CHECKSET,
						"network_source_acl.#": acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOtsInstance_multi(t *testing.T) {
	var v sdkots.InstanceInfo

	resourceId := "alicloud_ots_instance.default.4"
	ra := acctest.ResourceAttrInit(resourceId, otsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":          name + "${count.index}",
					"description":   name + "${count.index}",
					"instance_type": "Capacity",
					"count":         "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceOtsInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
			default = "%s"
		}
		
		data "alicloud_resource_manager_resource_groups" "default" {}`, name)
}

var otsInstanceBasicMap = map[string]string{
	"name":          acctest.CHECKSET,
	"accessed_by":   "Any",
	"instance_type": acctest.CHECKSET,
	"description":   acctest.CHECKSET,
}

func testAccCheckOtsInstanceExist(n string, instance *ots.RestOtsInstanceInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found OTS table: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OTS table ID is set")
		}

		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		otsService := ots.NewService(client)

		response, err := otsService.DescribeOtsInstance(rs.Primary.ID)

		if err != nil {
			return err
		}
		instance = &response
		return nil
	}
}
