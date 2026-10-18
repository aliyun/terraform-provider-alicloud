package alicloud

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func init() {
	resource.AddTestSweepers("alicloud_mse_nacos_config", &resource.Sweeper{
		Name: "alicloud_mse_nacos_config",
		F:    testSweepMseNacosConfig,
	})
}

// testSweepMseNacosConfig sweeps leftover Nacos configs across every MSE
// cluster in the region. A Nacos config is identified by (InstanceId,
// NamespaceId, DataId, Group). The sweeper walks every cluster, enumerates its
// engine namespaces via ListEngineNamespaces (plus the default namespace "") so
// every NamespaceId is known, lists configs per namespace via ListNacosConfigs
// (the data source API), and deletes each via DeleteNacosConfig mirroring
// resourceAlicloudMseNacosConfigDelete. Errors are logged and skipped; the
// sweeper always returns nil.
func testSweepMseNacosConfig(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	clusters, err := listMseSweepClusters(client)
	if err != nil {
		log.Printf("[ERROR] ListClusters for mse_nacos_config sweep failed: %s", err)
		return nil
	}

	for _, cluster := range clusters {
		clusterName := fmt.Sprint(cluster["ClusterAliasName"])
		instanceId := fmt.Sprint(cluster["InstanceId"])
		if !mseSweepNameAllowed(clusterName) {
			log.Printf("[INFO] Skipping nacos configs of cluster %s (%s): name not allowed", clusterName, instanceId)
			continue
		}
		// Enumerate namespaces: the default namespace ("") plus every namespace
		// returned by ListEngineNamespaces, so configs in any namespace can be
		// listed and deleted with the correct NamespaceId.
		namespaceIds := []string{""}
		if ns, nerr := listMseSweepEngineNamespaces(client, instanceId); nerr == nil {
			for _, n := range ns {
				namespaceIds = append(namespaceIds, fmt.Sprint(n["Namespace"]))
			}
		} else {
			log.Printf("[ERROR] ListEngineNamespaces for cluster %s failed: %s", instanceId, nerr)
		}
		seen := map[string]bool{}
		for _, nsId := range namespaceIds {
			configs, cerr := listMseSweepNacosConfigs(client, instanceId, nsId)
			if cerr != nil {
				log.Printf("[ERROR] ListNacosConfigs for cluster %s namespace %s failed: %s", instanceId, nsId, cerr)
				continue
			}
			for _, cfg := range configs {
				dataId := fmt.Sprint(cfg["DataId"])
				group := fmt.Sprint(cfg["Group"])
				if !mseSweepNameAllowed(dataId) {
					log.Printf("[INFO] Skipping mse nacos config: %s", dataId)
					continue
				}
				id := fmt.Sprintf("%s:%s:%s:%s", instanceId, nsId, dataId, group)
				if seen[id] {
					continue
				}
				seen[id] = true
				deleteAction := "DeleteNacosConfig"
				deleteRequest := map[string]interface{}{
					"InstanceId":  instanceId,
					"NamespaceId": nsId,
					"DataId":      dataId,
					"Group":       group,
				}
				if _, derr := client.RpcPost("mse", "2019-05-31", deleteAction, nil, deleteRequest, false); derr != nil {
					log.Printf("[ERROR] Failed to delete mse nacos config %s: %s", id, derr)
					continue
				}
				log.Printf("[INFO] Delete mse nacos config success: %s", id)
			}
		}
	}
	return nil
}

// listMseSweepEngineNamespaces returns every engine namespace of a cluster. It
// is shared by the engine_namespace and nacos_config sweepers.
func listMseSweepEngineNamespaces(client *connectivity.AliyunClient, instanceId string) ([]map[string]interface{}, error) {
	action := "ListEngineNamespaces"
	request := map[string]interface{}{
		"InstanceId": instanceId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	var err error
	if err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	addDebug(action, response, request)
	if fmt.Sprint(response["Success"]) == "false" {
		return nil, nil
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return nil, err
	}
	result, _ := resp.([]interface{})
	out := make([]map[string]interface{}, 0, len(result))
	for _, v := range result {
		out = append(out, v.(map[string]interface{}))
	}
	return out, nil
}

// listMseSweepNacosConfigs paginates ListNacosConfigs for one cluster+namespace
// and returns the raw config items. Pagination follows the data source.
func listMseSweepNacosConfigs(client *connectivity.AliyunClient, instanceId, namespaceId string) ([]map[string]interface{}, error) {
	action := "ListNacosConfigs"
	request := map[string]interface{}{
		"InstanceId":  instanceId,
		"NamespaceId": namespaceId,
		"PageSize":    PageSizeLarge,
		"PageNum":     1,
	}
	var out []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		if err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return nil, err
		}
		addDebug(action, response, request)
		resp, err := jsonpath.Get("$.Configurations", response)
		if err != nil {
			return nil, err
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			out = append(out, v.(map[string]interface{}))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	return out, nil
}

func TestAccAliCloudMseNacosConfig_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mse_nacos_config.default"
	checkoutSupportedRegions(t, true, connectivity.MSESupportRegions)

	ra := resourceAttrInit(resourceId, AlicloudMseNacosConfigMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMseNacosConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-msenacosconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMseNacosConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":     "${alicloud_mse_cluster.example.id}",
					"data_id":         "${var.name}:dataId",
					"group":           "${var.name}:group",
					"namespace_id":    "${alicloud_mse_engine_namespace.example.namespace_id}",
					"content":         "test",
					"app_name":        "test",
					"desc":            "test",
					"accept_language": "zh",
					"type":            "text",
					"tags":            "test",
					"beta_ips":        "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_id":  name + ":dataId",
						"group":    name + ":group",
						"content":  "test",
						"app_name": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content": "test_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name": "test_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desc": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desc": "test_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags": "test_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "xml",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "xml",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"beta_ips": "test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"beta_ips": "test_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accept_language", "beta_ips"},
			},
		},
	})
}

var AlicloudMseNacosConfigMap0 = map[string]string{}

func AlicloudMseNacosConfigBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_mse_cluster" "example" {
  connection_type       = "slb"
  net_type              = "privatenet"
  vswitch_id            = alicloud_vswitch.example.id
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_version       = "NACOS_2_0_0"
  instance_count        = "3"
  pub_network_flow      = "1"
  cluster_alias_name    = var.name
  mse_version           = "mse_pro"
  cluster_type          = "Nacos-Ans"
}

resource "alicloud_mse_engine_namespace" "example" {
  instance_id          = alicloud_mse_cluster.example.id
  namespace_show_name  = var.name
  namespace_id         = var.name
}

`, name)
}
