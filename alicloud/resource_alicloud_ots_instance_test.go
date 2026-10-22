package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ots_instance", &resource.Sweeper{
		Name: "alicloud_ots_instance",
		F:    testSweepOtsInstances,
	})
}

func testSweepOtsInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

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
			return fmt.Errorf("Error retrieving OTS Instances: %s", err)
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
			return err
		} else {
			req.PageNum = page
		}
	}
	for _, v := range insts {
		name := v.InstanceName
		skip := true
		if !sweepAll() {
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
		raw, err := otsService.client.WithTableStoreClient(name, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			return tableStoreClient.ListTable()
		})
		if err != nil {
			log.Printf("[ERROR] List OTS Instance %s table stores got an error: %#v.", name, err)
		}
		tables, _ := raw.(*tablestore.ListTableResponse)
		if tables != nil && len(tables.TableNames) > 0 {
			for _, t := range tables.TableNames {
				sweepTunnelsInTable(otsService.client, name, t)
				req := new(tablestore.DeleteTableRequest)
				req.TableName = t
				if _, err := otsService.client.WithTableStoreClient(name, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
					return tableStoreClient.DeleteTable(req)
				}); err != nil {
					log.Printf("[ERROR] Delete OTS Instance %s table store %s got an error: %#v.", name, t, err)
				}
			}
			time.Sleep(30 * time.Second)
		}
		log.Printf("[INFO] Deleting OTS Instance: %s", name)
		req := ots.CreateDeleteInstanceRequest()
		req.InstanceName = name
		_, err = client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete OTS Instance (%s): %s", name, err)
		}
	}

	return nil
}

func sweepTunnelsInTable(client *connectivity.AliyunClient, instanceName string, tableName string) {
	log.Printf("[INFO] Deleting OTS Tunnels, instance: %s, table: %s", instanceName, tableName)
	raw, err := client.WithTableStoreTunnelClient(instanceName, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
		return tunnelClient.ListTunnel(&otsTunnel.ListTunnelRequest{
			TableName: tableName,
		})
	})
	if err != nil {
		log.Printf("[ERROR] List OTS Tunnel failed, instance: %s, table: %s, error: %#v", instanceName, tableName, err)
		return
	}

	tunnels, _ := raw.(*otsTunnel.ListTunnelResponse)
	if tunnels != nil && len(tunnels.Tunnels) > 0 {
		for _, t := range tunnels.Tunnels {
			if _, err := client.WithTableStoreTunnelClient(instanceName, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
				return tunnelClient.DeleteTunnel(&otsTunnel.DeleteTunnelRequest{
					TableName:  tableName,
					TunnelName: t.TunnelName,
				})
			}); err != nil {
				log.Printf("[ERROR] Delete OTS Tunnel failed, instance: %s, table: %s, tunnel: %s, error: %#v", instanceName, tableName, t.TunnelName, err)
			}
		}
	}
}

func TestAccAliCloudOtsInstance_basic(t *testing.T) {
	var v ots.InstanceInfo

	resourceId := "alicloud_ots_instance.default"
	ra := resourceAttrInit(resourceId, otsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
						"resource_group_id":    CHECKSET,
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
						"network_type_acl.#":   CHECKSET,
						"network_source_acl.#": CHECKSET,
						"resource_group_id":    CHECKSET,
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
						"network_type_acl.#":   CHECKSET,
						"network_source_acl.#": CHECKSET,
						"resource_group_id":    CHECKSET,
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
						"network_type_acl.#":   CHECKSET,
						"network_source_acl.#": CHECKSET,
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessed_by": REMOVEKEY,
					"tags":        REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessed_by":          "Vpc",
						"tags.%":               "0",
						"tags.Created":         REMOVEKEY,
						"tags.For":             REMOVEKEY,
						"tags.Updated":         REMOVEKEY,
						"network_type_acl.#":   CHECKSET,
						"network_source_acl.#": CHECKSET,
						"resource_group_id":    CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOtsInstance_acl(t *testing.T) {
	var v ots.InstanceInfo

	resourceId := "alicloud_ots_instance.default"
	ra := resourceAttrInit(resourceId, otsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	var v ots.InstanceInfo

	resourceId := "alicloud_ots_instance.default"
	ra := resourceAttrInit(resourceId, otsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},
		// module name
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
						"network_type_acl.#":   CHECKSET,
						"network_source_acl.#": CHECKSET,
						"resource_group_id":    CHECKSET,
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
						"network_type_acl.#":   CHECKSET,
						"network_source_acl.#": CHECKSET,
						"resource_group_id":    CHECKSET,
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
						"network_type_acl.#":   CHECKSET,
						"network_source_acl.#": CHECKSET,
						"resource_group_id":    CHECKSET,
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
						"network_type_acl.#":   CHECKSET,
						"network_source_acl.#": CHECKSET,
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessed_by": REMOVEKEY,
					"tags":        REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessed_by":          "Vpc",
						"tags.%":               "0",
						"tags.Created":         REMOVEKEY,
						"tags.For":             REMOVEKEY,
						"tags.Updated":         REMOVEKEY,
						"network_type_acl.#":   CHECKSET,
						"network_source_acl.#": CHECKSET,
						"resource_group_id":    CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudOtsInstance_multi(t *testing.T) {
	var v ots.InstanceInfo

	resourceId := "alicloud_ots_instance.default.4"
	ra := resourceAttrInit(resourceId, otsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
	"name":          CHECKSET,
	"accessed_by":   "Any",
	"instance_type": CHECKSET,
	"description":   CHECKSET,
}

func testAccCheckOtsInstanceExist(n string, instance *RestOtsInstanceInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found OTS table: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OTS table ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		otsService := OtsService{client}

		response, err := otsService.DescribeOtsInstance(rs.Primary.ID)

		if err != nil {
			return err
		}
		instance = &response
		return nil
	}
}

// Reuse the native attachment resource's APIs, never DeleteInstance: an OTS
// instance can be shared by several VPCs and contains persistent table data.
type vpcOtsAttachmentAPI struct {
	listVpc      func(string, int, int) (*ots.ListVpcInfoByVpcResponse, error)
	listInstance func(string, int, int) (*ots.ListVpcInfoByInstanceResponse, error)
	unbind       func(string, string) error
	pause        func(time.Duration)
	attempts     int
}

func nativeVpcOtsAttachmentSweeper(region, vpcID string, client *connectivity.AliyunClient) *resource.Sweeper {
	api := vpcOtsAttachmentAPI{
		listVpc: func(id string, page, size int) (*ots.ListVpcInfoByVpcResponse, error) {
			request := ots.CreateListVpcInfoByVpcRequest()
			request.RegionId = region
			request.Method = "GET"
			request.VpcId = id
			request.PageNum = requests.NewInteger(page)
			request.PageSize = requests.NewInteger(size)
			raw, err := client.WithOtsClient(func(c *ots.Client) (interface{}, error) { return c.ListVpcInfoByVpc(request) })
			if err != nil {
				return nil, err
			}
			response, ok := raw.(*ots.ListVpcInfoByVpcResponse)
			if !ok {
				return nil, fmt.Errorf("invalid OTS VPC response")
			}
			return response, nil
		},
		listInstance: func(name string, page, size int) (*ots.ListVpcInfoByInstanceResponse, error) {
			request := ots.CreateListVpcInfoByInstanceRequest()
			request.RegionId = region
			request.Method = "GET"
			request.InstanceName = name
			request.PageNum = requests.NewInteger(page)
			request.PageSize = requests.NewInteger(size)
			raw, err := client.WithOtsClient(func(c *ots.Client) (interface{}, error) { return c.ListVpcInfoByInstance(request) })
			if err != nil {
				return nil, err
			}
			response, ok := raw.(*ots.ListVpcInfoByInstanceResponse)
			if !ok {
				return nil, fmt.Errorf("invalid OTS instance attachment response")
			}
			return response, nil
		},
		unbind: func(name, binding string) error {
			request := ots.CreateUnbindInstance2VpcRequest()
			request.RegionId = region
			request.RegionNo = region
			request.InstanceName = name
			request.InstanceVpcName = binding
			_, err := client.WithOtsClient(func(c *ots.Client) (interface{}, error) { return c.UnbindInstance2Vpc(request) })
			return err
		}, pause: time.Sleep, attempts: 120,
	}
	return vpcOtsAttachmentSweeper(region, vpcID, api)
}

func vpcOtsAttachmentSweeper(region, vpcID string, api vpcOtsAttachmentAPI) *resource.Sweeper {
	return &resource.Sweeper{Name: "alicloud_ots_instance_attachment", F: func(actual string) error {
		if region == "" || vpcID == "" || actual != region || api.listVpc == nil || api.listInstance == nil || api.unbind == nil || api.pause == nil || api.attempts < 1 {
			return fmt.Errorf("OTS attachment requires verified VPC/region and bounded native API")
		}
		candidates, err := listVpcOtsAttachments(region, vpcID, "", api)
		if err != nil {
			return err
		}
		var failures []string
		for _, candidate := range candidates {
			rows, err := listVpcOtsAttachments(region, vpcID, candidate.InstanceName, api)
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
			present := false
			for _, row := range rows {
				if row.InstanceVpcName == candidate.InstanceVpcName {
					if row.VpcId != vpcID {
						err = fmt.Errorf("OTS attachment ownership changed; refusing unbind")
						break
					}
					present = true
				}
			}
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
			if !present {
				continue
			}
			if err = api.unbind(candidate.InstanceName, candidate.InstanceVpcName); err != nil {
				failures = append(failures, fmt.Sprintf("OTS attachment unbind: %v", err))
				continue
			}
			gone := false
			for attempt := 0; attempt < api.attempts; attempt++ {
				rows, err = listVpcOtsAttachments(region, vpcID, candidate.InstanceName, api)
				if err != nil {
					break
				}
				present = false
				for _, row := range rows {
					if row.InstanceVpcName == candidate.InstanceVpcName {
						if row.VpcId != vpcID {
							err = fmt.Errorf("OTS attachment changed VPC during wait")
							break
						}
						present = true
					}
				}
				if err != nil {
					break
				}
				if !present {
					gone = true
					break
				}
				if attempt+1 < api.attempts {
					api.pause(5 * time.Second)
				}
			}
			if err != nil {
				failures = append(failures, err.Error())
			} else if !gone {
				failures = append(failures, "OTS attachment still present after bounded unbind wait")
			}
		}
		if len(failures) > 0 {
			return fmt.Errorf("OTS VPC attachment cleanup incomplete: %s", strings.Join(failures, "; "))
		}
		return nil
	}}
}

// By-VPC responses provide VpcId in the response context, whereas by-instance
// responses provide InstanceName there. Validate context before filling only
// that contextual field, preserving every other attachment's independent VPC.
func listVpcOtsAttachments(region, vpcID, instance string, api vpcOtsAttachmentAPI) ([]ots.VpcInfo, error) {
	var all []ots.VpcInfo
	seen := map[string]bool{}
	const size = 100
	total := -1
	for page := 1; page <= 10000; page++ {
		var rows []ots.VpcInfo
		var count, returnedPage, returnedSize int
		if instance == "" {
			response, err := api.listVpc(vpcID, page, size)
			if err != nil {
				return nil, err
			}
			if response == nil || response.VpcId != vpcID {
				return nil, fmt.Errorf("OTS VPC inventory ownership unknown or changed")
			}
			rows = response.VpcInfos.VpcInfo
			count = response.TotalCount
			returnedPage = response.PageNum
			returnedSize = response.PageSize
		} else {
			response, err := api.listInstance(instance, page, size)
			if err != nil {
				return nil, err
			}
			if response == nil || response.InstanceName != instance {
				return nil, fmt.Errorf("OTS instance inventory ownership unknown or changed")
			}
			rows = response.VpcInfos.VpcInfo
			count = response.TotalCount
			returnedPage = response.PageNum
			returnedSize = response.PageSize
		}
		if returnedPage != page || returnedSize != size || count < 0 || (total >= 0 && count != total) {
			return nil, fmt.Errorf("OTS attachment pagination metadata inconsistent")
		}
		total = count
		for _, row := range rows {
			if instance == "" {
				if row.VpcId != "" && row.VpcId != vpcID {
					return nil, fmt.Errorf("OTS VPC inventory crossed scope")
				}
				row.VpcId = vpcID
			} else {
				if row.InstanceName != "" && row.InstanceName != instance {
					return nil, fmt.Errorf("OTS instance inventory crossed scope")
				}
				row.InstanceName = instance
			}
			if row.RegionNo != region || row.VpcId == "" || row.InstanceName == "" || row.InstanceVpcName == "" {
				return nil, fmt.Errorf("OTS attachment identity unknown; refusing unbind")
			}
			key := row.InstanceName + "/" + row.InstanceVpcName
			if seen[key] {
				return nil, fmt.Errorf("OTS attachment pagination repeated identity")
			}
			seen[key] = true
			all = append(all, row)
		}
		if len(all) > total {
			return nil, fmt.Errorf("OTS attachment count exceeded inventory")
		}
		if len(all) == total {
			return all, nil
		}
		if len(rows) != size {
			return nil, fmt.Errorf("OTS attachment inventory truncated")
		}
	}
	return nil, fmt.Errorf("OTS attachment pagination exceeded bound")
}

func vpcOtsTestAPI(t *testing.T, scenario string) (vpcOtsAttachmentAPI, *int, *int) {
	t.Helper()
	deletes, reads := 0, 0
	target := ots.VpcInfo{InstanceName: "instance", InstanceVpcName: "target-link", RegionNo: "cn-shanghai"}
	protected := ots.VpcInfo{InstanceVpcName: "protected-link", RegionNo: "cn-shanghai", VpcId: "vpc-protected"}
	api := vpcOtsAttachmentAPI{pause: func(time.Duration) {}, attempts: 2}
	api.listVpc = func(vpc string, page, size int) (*ots.ListVpcInfoByVpcResponse, error) {
		if vpc != "vpc-target" || page != 1 || size != 100 {
			t.Fatalf("unexpected inventory scope")
		}
		r := &ots.ListVpcInfoByVpcResponse{VpcId: vpc, PageNum: page, PageSize: size, TotalCount: 1}
		r.VpcInfos.VpcInfo = []ots.VpcInfo{target}
		return r, nil
	}
	api.listInstance = func(instance string, page, size int) (*ots.ListVpcInfoByInstanceResponse, error) {
		reads++
		if instance != "instance" {
			t.Fatal("wrong shared instance")
		}
		if scenario == "read-error" || (scenario == "wait-read-error" && deletes > 0) {
			return nil, fmt.Errorf("read denied")
		}
		row := target
		row.InstanceName = ""
		row.VpcId = "vpc-target"
		switch scenario {
		case "wrong-vpc":
			row.VpcId = "vpc-protected"
		case "unknown-region":
			row.RegionNo = ""
		case "unknown-vpc":
			row.VpcId = ""
		case "unknown-binding":
			row.InstanceVpcName = ""
		}
		if scenario == "wait-scope-change" && deletes > 0 {
			row.VpcId = "vpc-protected"
		}
		r := &ots.ListVpcInfoByInstanceResponse{InstanceName: instance, PageNum: page, PageSize: size}
		r.VpcInfos.VpcInfo = []ots.VpcInfo{protected}
		if deletes == 0 || scenario == "timeout" || scenario == "wait-scope-change" || (scenario == "async" && reads == 2) {
			r.VpcInfos.VpcInfo = append(r.VpcInfos.VpcInfo, row)
		}
		r.TotalCount = len(r.VpcInfos.VpcInfo)
		if scenario == "wrong-context" {
			r.InstanceName = "other"
		}
		return r, nil
	}
	api.unbind = func(instance, binding string) error {
		if instance != "instance" || binding != "target-link" {
			t.Fatal("shared/protected attachment touched")
		}
		deletes++
		if scenario == "write-error" {
			return fmt.Errorf("unbind denied")
		}
		return nil
	}
	return api, &deletes, &reads
}

func TestUnitVpcSweepOtsAttachmentSharedInstanceExactUnbind(t *testing.T) {
	api, deletes, reads := vpcOtsTestAPI(t, "async")
	pauses := 0
	api.pause = func(time.Duration) { pauses++ }
	if err := vpcOtsAttachmentSweeper("cn-shanghai", "vpc-target", api).F("cn-shanghai"); err != nil {
		t.Fatal(err)
	}
	if *deletes != 1 || *reads != 3 || pauses != 1 {
		t.Fatalf("unbind=%d fresh reads=%d waits=%d", *deletes, *reads, pauses)
	}
}

func TestUnitVpcSweepOtsAttachmentUnknownOwnershipNeverUnbinds(t *testing.T) {
	for _, scenario := range []string{"wrong-vpc", "unknown-region", "unknown-vpc", "unknown-binding", "wrong-context", "read-error"} {
		t.Run(scenario, func(t *testing.T) {
			api, deletes, _ := vpcOtsTestAPI(t, scenario)
			if err := vpcOtsAttachmentSweeper("cn-shanghai", "vpc-target", api).F("cn-shanghai"); err == nil || *deletes != 0 {
				t.Fatalf("unsafe scope err=%v unbind=%d", err, *deletes)
			}
		})
	}
}

func TestUnitVpcSweepOtsAttachmentWaitAndWriteFailures(t *testing.T) {
	for _, scenario := range []string{"timeout", "write-error", "wait-read-error", "wait-scope-change"} {
		t.Run(scenario, func(t *testing.T) {
			api, deletes, _ := vpcOtsTestAPI(t, scenario)
			if err := vpcOtsAttachmentSweeper("cn-shanghai", "vpc-target", api).F("cn-shanghai"); err == nil || *deletes != 1 {
				t.Fatalf("release failure hidden err=%v unbind=%d", err, *deletes)
			}
		})
	}
}

func TestUnitVpcSweepOtsAttachmentSDKContextMapping(t *testing.T) {
	var byVpc ots.ListVpcInfoByVpcResponse
	if err := json.Unmarshal([]byte(`{"VpcId":"vpc-target","TotalCount":1,"PageNum":1,"PageSize":100,"VpcInfos":{"VpcInfo":[{"InstanceName":"instance","InstanceVpcName":"link","RegionNo":"cn-shanghai"}]}}`), &byVpc); err != nil {
		t.Fatal(err)
	}
	var byInstance ots.ListVpcInfoByInstanceResponse
	if err := json.Unmarshal([]byte(`{"InstanceName":"instance","TotalCount":1,"PageNum":1,"PageSize":100,"VpcInfos":{"VpcInfo":[{"VpcId":"vpc-target","InstanceVpcName":"link","RegionNo":"cn-shanghai"}]}}`), &byInstance); err != nil {
		t.Fatal(err)
	}
	api := vpcOtsAttachmentAPI{listVpc: func(string, int, int) (*ots.ListVpcInfoByVpcResponse, error) { return &byVpc, nil }, listInstance: func(string, int, int) (*ots.ListVpcInfoByInstanceResponse, error) { return &byInstance, nil }}
	for _, instance := range []string{"", "instance"} {
		rows, err := listVpcOtsAttachments("cn-shanghai", "vpc-target", instance, api)
		if err != nil || len(rows) != 1 || rows[0].InstanceName != "instance" || rows[0].VpcId != "vpc-target" {
			t.Fatalf("native contextual identity lost rows=%v err=%v", rows, err)
		}
	}
}

func TestUnitVpcSweepOtsAttachmentPaginationRejectsIncomplete(t *testing.T) {
	for _, scenario := range []string{"truncated", "wrong-page", "wrong-context", "cross-vpc", "repeated"} {
		t.Run(scenario, func(t *testing.T) {
			api := vpcOtsAttachmentAPI{listVpc: func(vpc string, page, size int) (*ots.ListVpcInfoByVpcResponse, error) {
				r := &ots.ListVpcInfoByVpcResponse{VpcId: vpc, PageNum: page, PageSize: size, TotalCount: 2}
				row := ots.VpcInfo{InstanceName: "instance", InstanceVpcName: "link", RegionNo: "cn-shanghai"}
				r.VpcInfos.VpcInfo = []ots.VpcInfo{row}
				switch scenario {
				case "wrong-page":
					r.PageNum = 0
				case "wrong-context":
					r.VpcId = "vpc-protected"
				case "cross-vpc":
					r.VpcInfos.VpcInfo[0].VpcId = "vpc-protected"
				case "repeated":
					r.VpcInfos.VpcInfo = append(r.VpcInfos.VpcInfo, row)
				}
				return r, nil
			}}
			if _, err := listVpcOtsAttachments("cn-shanghai", "vpc-target", "", api); err == nil {
				t.Fatal("unverified pagination accepted")
			}
		})
	}
}

func TestUnitVpcSweepOtsAttachmentCompletePagination(t *testing.T) {
	calls := 0
	api := vpcOtsAttachmentAPI{listVpc: func(vpc string, page, size int) (*ots.ListVpcInfoByVpcResponse, error) {
		calls++
		r := &ots.ListVpcInfoByVpcResponse{VpcId: vpc, PageNum: page, PageSize: size, TotalCount: 101}
		count := 100
		if page == 2 {
			count = 1
		}
		for i := 0; i < count; i++ {
			r.VpcInfos.VpcInfo = append(r.VpcInfos.VpcInfo, ots.VpcInfo{InstanceName: "instance", InstanceVpcName: fmt.Sprintf("page%d-item%d", page, i), RegionNo: "cn-shanghai"})
		}
		return r, nil
	}}
	rows, err := listVpcOtsAttachments("cn-shanghai", "vpc-target", "", api)
	if err != nil || len(rows) != 101 || calls != 2 {
		t.Fatalf("inventory incomplete rows=%d calls=%d err=%v", len(rows), calls, err)
	}
}
