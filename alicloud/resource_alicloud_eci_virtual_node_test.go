package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_eci_virtual_node",
		&resource.Sweeper{
			Name: "alicloud_eci_virtual_node",
			F:    testSweepECIVirtualNode,
		})
}

func testSweepECIVirtualNode(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeVirtualNodes"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Eci", "2018-08-08", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.VirtualNodes", response)

		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.VirtualNodes", action, err)
			return nil
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["VirtualNodeName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ECI VirtualNode: %s", item["VirtualNodeName"].(string))
				continue
			}

			sweeped = true
			action := "DeleteVirtualNode"
			request := map[string]interface{}{
				"VirtualNodeId": item["VirtualNodeId"],
				"RegionId":      client.RegionId,
			}
			request["ClientToken"] = buildClientToken("DeleteVirtualNode")
			_, err = client.RpcPost("Eci", "2018-08-08", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete ECI VirtualNode (%s): %s", item["VirtualNodeName"].(string), err)
			}
			if sweeped {
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete ECI VirtualNode success: %s ", item["VirtualNodeName"].(string))
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return nil
}

func TestAccAliCloudEciVirtualNode_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_virtual_node.default"
	ra := resourceAttrInit(resourceId, AlicloudECIVirtualNodeMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciVirtualNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secivirtualnode%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECIVirtualNodeBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "KUBE_CONFIG")
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.default.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.1}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"virtual_node_name": name,
					"eip_instance_id":   "${alicloud_eip_address.default.id}",
					"taints": []map[string]interface{}{
						{
							"effect": "NoSchedule",
							"key":    "Tf1",
							"value":  "Test1",
						},
					},
					"zone_id":               "${data.alicloud_eci_zones.default.zones.0.zone_ids.0}",
					"enable_public_network": true,
					"tls_bootstrap_enabled": true,
					"custom_resources":      "CreateVirtualNodeValue",
					"kube_config":           "${var.kube_config}",
					"tags": map[string]string{
						"Created": "Tf1",
						"For":     "Test1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
						"vswitch_id":        CHECKSET,
						"resource_group_id": CHECKSET,
						"virtual_node_name": name,
						"tags.%":            "2",
						"tags.Created":      "Tf1",
						"tags.For":          "Test1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.default.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.1}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"virtual_node_name": name,
					"eip_instance_id":   "${alicloud_eip_address.default.id}",
					"taints": []map[string]interface{}{
						{
							"effect": "NoExecute",
							"key":    "Tf2",
							"value":  "Test2",
						},
					},
					"zone_id":               "${data.alicloud_eci_zones.default.zones.0.zone_ids.0}",
					"enable_public_network": true,
					"tls_bootstrap_enabled": true,
					"custom_resources":      "UpdateVirtualNodeValue",
					"kube_config":           "${var.kube_config}",
					"cluster_domain":        "cluster.local",
					"cluster_dns":           "192.168.0.10",
					"tags": map[string]string{
						"Created": "Tf1",
						"For":     "Test1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
						"vswitch_id":        CHECKSET,
						"resource_group_id": CHECKSET,
						"virtual_node_name": name,
						"cluster_domain":    "cluster.local",
						"cluster_dns":       "192.168.0.10",
						"tags.%":            "2",
						"tags.Created":      "Tf1",
						"tags.For":          "Test1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kube_config", "taints", "cluster_domain", "cluster_dns"},
			},
		},
	})
}

func TestAccAliCloudEciVirtualNode_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eci_virtual_node.default"
	ra := resourceAttrInit(resourceId, AlicloudECIVirtualNodeMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EciService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEciVirtualNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secivirtualnode%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECIVirtualNodeBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "KUBE_CONFIG")
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.default.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"virtual_node_name": name,
					"kube_config":       "${var.kube_config}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
						"vswitch_id":        CHECKSET,
						"resource_group_id": CHECKSET,
						"virtual_node_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kube_config", "taints"},
			},
		},
	})
}

var AlicloudECIVirtualNodeMap0 = map[string]string{}

func AlicloudECIVirtualNodeBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "kube_config" {
  default = "%s"
}

data "alicloud_eci_zones" "default"{}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_eci_zones.default.zones.0.zone_ids.0
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  name   = var.name
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name, os.Getenv("KUBE_CONFIG"))
}

// lintignore: R001
func TestUnitAlicloudEciVirtualNode(t *testing.T) {
	p := Provider().ResourcesMap
	fields := p["alicloud_eci_virtual_node"].Schema

	// Governance attributes mapped from CreateVirtualNode/UpdateVirtualNode must
	// exist with the correct type and diff semantics: custom_resources,
	// cluster_domain and cluster_dns are updatable; tls_bootstrap_enabled,
	// enable_public_network and zone_id force replacement.
	cases := []struct {
		name      string
		typ       schema.ValueType
		optional  bool
		forceNew  bool
		sensitive bool
	}{
		{"custom_resources", schema.TypeString, true, false, true},
		{"cluster_domain", schema.TypeString, true, false, true},
		{"cluster_dns", schema.TypeString, true, false, true},
		{"tls_bootstrap_enabled", schema.TypeBool, true, true, true},
		{"enable_public_network", schema.TypeBool, true, true, false},
		{"zone_id", schema.TypeString, true, true, false},
	}
	for _, c := range cases {
		field, ok := fields[c.name]
		if !assert.True(t, ok, "field %q should be present in schema", c.name) {
			continue
		}
		assert.Equal(t, c.typ, field.Type, "field %q type", c.name)
		assert.Equal(t, c.optional, field.Optional, "field %q optional", c.name)
		assert.Equal(t, c.forceNew, field.ForceNew, "field %q force_new", c.name)
		assert.Equal(t, c.sensitive, field.Sensitive, "field %q sensitive", c.name)
	}

	// Create/Read/Update/Delete callbacks and the Create+Update timeouts must be wired up.
	r := p["alicloud_eci_virtual_node"]
	assert.NotNil(t, r.Create)
	assert.NotNil(t, r.Read)
	assert.NotNil(t, r.Update)
	assert.NotNil(t, r.Delete)
	assert.NotNil(t, r.Timeouts)
}
