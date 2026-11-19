package cs_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	roacs "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cs"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func autoscalingConfigTestBoolKeys() []string {
	return []string{
		"skip_nodes_with_system_pods",
		"skip_nodes_with_local_storage",
		"daemonset_eviction_for_nodes",
		"recycle_node_deletion_enabled",
		"scale_up_from_zero",
	}
}

func TestUnitCSAutoscalingConfigRequestBooleans(t *testing.T) {
	cases := []struct {
		name   string
		values map[string]interface{}
		raw    map[string]cty.Value
		want   map[string]bool
		update bool
	}{
		{
			name: "explicit false is sent",
			values: map[string]interface{}{
				"skip_nodes_with_system_pods":   false,
				"skip_nodes_with_local_storage": false,
				"daemonset_eviction_for_nodes":  false,
				"recycle_node_deletion_enabled": false,
				"scale_up_from_zero":            false,
			},
			raw: map[string]cty.Value{
				"skip_nodes_with_system_pods":   cty.False,
				"skip_nodes_with_local_storage": cty.False,
				"daemonset_eviction_for_nodes":  cty.False,
				"recycle_node_deletion_enabled": cty.False,
				"scale_up_from_zero":            cty.False,
			},
			want: map[string]bool{
				"skip_nodes_with_system_pods":   false,
				"skip_nodes_with_local_storage": false,
				"daemonset_eviction_for_nodes":  false,
				"recycle_node_deletion_enabled": false,
				"scale_up_from_zero":            false,
			},
		},
		{
			name: "explicit true is sent",
			values: map[string]interface{}{
				"skip_nodes_with_system_pods":   true,
				"skip_nodes_with_local_storage": true,
				"daemonset_eviction_for_nodes":  true,
				"recycle_node_deletion_enabled": true,
				"scale_up_from_zero":            true,
			},
			raw: map[string]cty.Value{
				"skip_nodes_with_system_pods":   cty.True,
				"skip_nodes_with_local_storage": cty.True,
				"daemonset_eviction_for_nodes":  cty.True,
				"recycle_node_deletion_enabled": cty.True,
				"scale_up_from_zero":            cty.True,
			},
			want: map[string]bool{
				"skip_nodes_with_system_pods":   true,
				"skip_nodes_with_local_storage": true,
				"daemonset_eviction_for_nodes":  true,
				"recycle_node_deletion_enabled": true,
				"scale_up_from_zero":            true,
			},
		},
		{
			name:   "omitted is not sent",
			values: map[string]interface{}{},
			raw: map[string]cty.Value{
				"skip_nodes_with_system_pods":   cty.NullVal(cty.Bool),
				"skip_nodes_with_local_storage": cty.NullVal(cty.Bool),
				"daemonset_eviction_for_nodes":  cty.NullVal(cty.Bool),
				"recycle_node_deletion_enabled": cty.NullVal(cty.Bool),
				"scale_up_from_zero":            cty.NullVal(cty.Bool),
			},
			want: map[string]bool{},
		},
		{
			name:   "update removal after false is not sent",
			update: true,
			values: map[string]interface{}{
				"scale_down_enabled": false,
			},
			raw: map[string]cty.Value{
				"scale_down_enabled":            cty.False,
				"skip_nodes_with_system_pods":   cty.NullVal(cty.Bool),
				"skip_nodes_with_local_storage": cty.NullVal(cty.Bool),
				"daemonset_eviction_for_nodes":  cty.NullVal(cty.Bool),
				"recycle_node_deletion_enabled": cty.NullVal(cty.Bool),
				"scale_up_from_zero":            cty.NullVal(cty.Bool),
			},
			want: map[string]bool{},
		},
		{
			name: "mixed explicit and omitted",
			values: map[string]interface{}{
				"skip_nodes_with_system_pods": false,
				"scale_up_from_zero":          true,
			},
			raw: map[string]cty.Value{
				"skip_nodes_with_system_pods":   cty.False,
				"skip_nodes_with_local_storage": cty.NullVal(cty.Bool),
				"daemonset_eviction_for_nodes":  cty.NullVal(cty.Bool),
				"recycle_node_deletion_enabled": cty.NullVal(cty.Bool),
				"scale_up_from_zero":            cty.True,
			},
			want: map[string]bool{
				"skip_nodes_with_system_pods": false,
				"scale_up_from_zero":          true,
			},
		},
		{
			// lifecycle.ignore_changes keeps the applied value while the raw
			// configuration still holds the ignored one; the planned value must
			// be sent so the ignored setting is not overwritten.
			name: "ignored raw value sends planned value",
			values: map[string]interface{}{
				"skip_nodes_with_system_pods": true,
			},
			raw: map[string]cty.Value{
				"skip_nodes_with_system_pods": cty.False,
			},
			want: map[string]bool{
				"skip_nodes_with_system_pods": true,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var body map[string]interface{}
			requestCount := 0
			client := cs.CsKubeconfigTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				switch {
				case req.Method == http.MethodPost && req.URL.Path == "/cluster/c1/autoscale/config/":
					defer req.Body.Close()
					requestCount++
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
						t.Errorf("decode CreateAutoscalingConfig body: %v", err)
					}
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"RequestId":"test"}`))
				case req.Method == http.MethodGet && req.URL.Path == "/clusters/c1/addon_instances/cluster-autoscaler":
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"name":"cluster-autoscaler","state":"active","version":"1.0","config":"{}"}`))
				default:
					http.NotFound(w, req)
				}
			})

			config := map[string]interface{}{"cluster_id": "c1"}
			for key, value := range tc.values {
				config[key] = value
			}
			resourceSchema := cs.ResourceCSAutoscalingConfig()
			var priorState *sdktf.InstanceState
			if tc.update {
				priorValues := map[string]interface{}{"cluster_id": "c1"}
				for _, key := range autoscalingConfigTestBoolKeys() {
					priorValues[key] = false
				}
				priorData := schema.TestResourceDataRaw(t, resourceSchema.Schema, priorValues)
				priorData.SetId("c1")
				priorState = priorData.State()
			}
			diff, err := resourceSchema.Diff(context.Background(), priorState, sdktf.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatalf("diff: %v", err)
			}
			if diff == nil {
				t.Fatal("expected a resource diff")
			}
			if tc.update && diff.RequiresNew() {
				t.Fatal("removal must use Update, not replacement")
			}
			raw := map[string]cty.Value{"cluster_id": cty.StringVal("c1")}
			for key, value := range tc.raw {
				raw[key] = value
			}
			diff.RawConfig = cty.ObjectVal(raw)
			applied, diags := resourceSchema.Apply(context.Background(), priorState, diff, client)
			if diags.HasError() {
				t.Fatalf("apply: %v", diags)
			}
			if applied == nil || applied.ID != "c1" {
				t.Fatalf("applied state ID = %#v, want c1", applied)
			}
			if requestCount != 1 {
				t.Fatalf("CreateAutoscalingConfig request count = %d, want 1", requestCount)
			}
			if body == nil {
				t.Fatal("CreateAutoscalingConfig was not called")
			}
			for _, key := range autoscalingConfigTestBoolKeys() {
				value, present := body[key]
				want, wantPresent := tc.want[key]
				if !wantPresent {
					if present {
						t.Fatalf("%s = %v in the request body, want it omitted", key, value)
					}
					continue
				}
				if !present {
					t.Fatalf("%s is absent from the request body, want %t", key, want)
				}
				got, ok := value.(bool)
				if !ok || got != want {
					t.Fatalf("%s = %v in the request body, want %t", key, value, want)
				}
			}
		})
	}
}

// lintignore: AT001
func TestAccAliCloudCSAutoscalingConfig_basic(t *testing.T) {
	var v *roacs.CreateAutoscalingConfigRequest
	resourceId := "alicloud_cs_autoscaling_config.default"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, baseMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsAutoscalingConfig"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSAutoscalingConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAuoscalingConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Create with explicit false values: they must reach the API instead
				// of falling back to the server defaults (skip_nodes_with_system_pods
				// and scale_up_from_zero default to true). The elastic node pool is
				// created first because priorities references its scaling group.
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":                    "${alicloud_cs_managed_kubernetes.default.id}",
					"cool_down_duration":            "10m",
					"unneeded_duration":             "10m",
					"utilization_threshold":         "0.5",
					"gpu_utilization_threshold":     "0.5",
					"scan_interval":                 "30s",
					"scale_down_enabled":            "true",
					"expander":                      "priority",
					"scaler_type":                   "cluster-autoscaler",
					"skip_nodes_with_system_pods":   "false",
					"skip_nodes_with_local_storage": "false",
					"daemonset_eviction_for_nodes":  "false",
					"recycle_node_deletion_enabled": "false",
					"scale_up_from_zero":            "false",
					"priorities": map[string]interface{}{
						"10": "${alicloud_cs_kubernetes_node_pool.default.scaling_group_id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// The autoscaler is still starting right after the cluster
						// and node pool are created; wait before the next update so
						// CreateAutoscalingConfig does not hit InvalidAutoScalerStatus.
						acctest.ForceSleep:              "180",
						"cluster_id":                    acctest.CHECKSET,
						"cool_down_duration":            acctest.CHECKSET,
						"unneeded_duration":             acctest.CHECKSET,
						"utilization_threshold":         acctest.CHECKSET,
						"gpu_utilization_threshold":     acctest.CHECKSET,
						"scan_interval":                 acctest.CHECKSET,
						"scale_down_enabled":            acctest.CHECKSET,
						"expander":                      acctest.CHECKSET,
						"scaler_type":                   acctest.CHECKSET,
						"skip_nodes_with_system_pods":   "false",
						"skip_nodes_with_local_storage": "false",
						"daemonset_eviction_for_nodes":  "false",
						"recycle_node_deletion_enabled": "false",
						"scale_up_from_zero":            "false",
						"priorities.%":                  "1",
						"priorities.10":                 acctest.CHECKSET,
					}),
				),
			},
			{
				// false -> true for the same batch of boolean arguments.
				Config: testAccConfig(map[string]interface{}{
					"skip_nodes_with_system_pods":   "true",
					"skip_nodes_with_local_storage": "true",
					"daemonset_eviction_for_nodes":  "true",
					"max_graceful_termination_sec":  "14400",
					"min_replica_count":             "0",
					"recycle_node_deletion_enabled": "true",
					"scale_up_from_zero":            "true",
				}),
				// min_replica_count keeps using GetOk, so an explicit 0 is dropped
				// (pre-existing); its zero state is not evidence of delivery, so it
				// is configured but not asserted here.
				// The cluster holds a single autoscaler operation lock: a config
				// change issued while the previous one is still in progress is
				// rejected with InvalidAutoScalerStatus. Every update step waits
				// before the next one, as the create step does.
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						acctest.ForceSleep:              "180",
						"skip_nodes_with_system_pods":   "true",
						"skip_nodes_with_local_storage": "true",
						"daemonset_eviction_for_nodes":  "true",
						"max_graceful_termination_sec":  "14400",
						"recycle_node_deletion_enabled": "true",
						"scale_up_from_zero":            "true",
					}),
				),
			},
			{
				// true -> false again. The state assertions below check the
				// planned values converge; the request body itself is covered by
				// TestUnitCSAutoscalingConfigRequest*.
				Config: testAccConfig(map[string]interface{}{
					"skip_nodes_with_system_pods":   "false",
					"skip_nodes_with_local_storage": "false",
					"daemonset_eviction_for_nodes":  "false",
					"recycle_node_deletion_enabled": "false",
					"scale_up_from_zero":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						acctest.ForceSleep:              "180",
						"skip_nodes_with_system_pods":   "false",
						"skip_nodes_with_local_storage": "false",
						"daemonset_eviction_for_nodes":  "false",
						"recycle_node_deletion_enabled": "false",
						"scale_up_from_zero":            "false",
					}),
				),
			},
			{
				// Remove the boolean arguments. They are Optional without a
				// default or Computed, so the legacy SDK plans the zero value and
				// the request omits them: the provider stops sending the argument.
				// Removing a value that is true in the cluster is not supported;
				// write false explicitly instead.
				Config: testAccConfig(map[string]interface{}{
					"skip_nodes_with_system_pods":   acctest.REMOVEKEY,
					"skip_nodes_with_local_storage": acctest.REMOVEKEY,
					"daemonset_eviction_for_nodes":  acctest.REMOVEKEY,
					"recycle_node_deletion_enabled": acctest.REMOVEKEY,
					"scale_up_from_zero":            acctest.REMOVEKEY,
					"scale_down_enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						acctest.ForceSleep:              "180",
						"scale_down_enabled":            "false",
						"skip_nodes_with_system_pods":   "false",
						"skip_nodes_with_local_storage": "false",
						"daemonset_eviction_for_nodes":  "false",
						"recycle_node_deletion_enabled": "false",
						"scale_up_from_zero":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cool_down_duration":        "5m",
					"unneeded_duration":         "5m",
					"utilization_threshold":     "0.6",
					"gpu_utilization_threshold": "0.6",
					"scan_interval":             "40s",
					"expander":                  "least-waste",
					"priorities":                acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						acctest.ForceSleep:          "180",
						"cluster_id":                acctest.CHECKSET,
						"cool_down_duration":        acctest.CHECKSET,
						"unneeded_duration":         acctest.CHECKSET,
						"utilization_threshold":     acctest.CHECKSET,
						"gpu_utilization_threshold": acctest.CHECKSET,
						"scan_interval":             acctest.CHECKSET,
						"expander":                  acctest.CHECKSET,
						// The check map accumulates across steps: drop the
						// priorities assertions of the earlier steps and expect
						// the removed map instead.
						"priorities.%":  "0",
						"priorities.10": acctest.REMOVEKEY,
					}),
				),
			},
			{
				// Modify priorities so both setting and updating the map are
				// exercised (the read path does not populate priorities).
				Config: testAccConfig(map[string]interface{}{
					"expander": "priority",
					"priorities": map[string]interface{}{
						"20": "${alicloud_cs_kubernetes_node_pool.default.scaling_group_id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						acctest.ForceSleep: "180",
						"expander":         acctest.CHECKSET,
						"priorities.%":     "1",
						"priorities.20":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expander":   "least-waste",
					"priorities": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						acctest.ForceSleep: "180",
						"expander":         acctest.CHECKSET,
						"priorities.%":     "0",
						"priorities.20":    acctest.REMOVEKEY,
					}),
				),
			},
			{
				// The cluster in the test account ships ack-goatscaler, which
				// conflicts with cluster-autoscaler, so the read path cannot
				// resolve an addon instance here; this step only checks that the
				// scaler_type update is accepted and planned.
				Config: testAccConfig(map[string]interface{}{
					"scaler_type": "goatscaler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaler_type": acctest.CHECKSET,
					}),
				),
			},
			{
				// Import smoke cover: the importer copies the import ID (the
				// cluster ID) into cluster_id so the read path can locate the
				// cluster addon. Attribute-level comparison is not available in
				// this environment because the read path cannot resolve the
				// cluster-autoscaler addon instance, so this step only proves the
				// import succeeds and keeps the import identity.
				ResourceName: resourceId,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					cluster, ok := s.RootModule().Resources["alicloud_cs_managed_kubernetes.default"]
					if !ok || cluster.Primary == nil || cluster.Primary.ID == "" {
						return "", fmt.Errorf("cannot find the test cluster state")
					}
					return cluster.Primary.ID, nil
				},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected one imported state, got %d", len(states))
					}
					importId := states[0].ID
					if importId == "" {
						return fmt.Errorf("imported state has an empty ID")
					}
					if got := states[0].Attributes["cluster_id"]; got != importId {
						return fmt.Errorf("imported cluster_id = %q, want the import ID %q", got, importId)
					}
					return nil
				},
			},
		},
	})
}

// lintignore: AT001
func TestAccAliCloudCSAutoscalingConfig_complete(t *testing.T) {
	var v *roacs.CreateAutoscalingConfigRequest
	resourceId := "alicloud_cs_autoscaling_config.default"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, baseMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsAutoscalingConfig"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSAutoscalingConfigComplete-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAuoscalingConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Create with explicit true values so the create path is exercised
				// with non-default booleans (the basic case creates with false).
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":                    "${alicloud_cs_managed_kubernetes.default.id}",
					"cool_down_duration":            "10m",
					"unneeded_duration":             "10m",
					"utilization_threshold":         "0.5",
					"gpu_utilization_threshold":     "0.5",
					"scan_interval":                 "30s",
					"scale_down_enabled":            "true",
					"expander":                      "priority",
					"scaler_type":                   "cluster-autoscaler",
					"skip_nodes_with_system_pods":   "true",
					"skip_nodes_with_local_storage": "true",
					"daemonset_eviction_for_nodes":  "true",
					"max_graceful_termination_sec":  "14400",
					"min_replica_count":             "1",
					"recycle_node_deletion_enabled": "true",
					"scale_up_from_zero":            "true",
					"priorities": map[string]interface{}{
						"10": "${alicloud_cs_kubernetes_node_pool.default.scaling_group_id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                    acctest.CHECKSET,
						"cool_down_duration":            "10m",
						"unneeded_duration":             "10m",
						"utilization_threshold":         "0.5",
						"gpu_utilization_threshold":     "0.5",
						"scan_interval":                 "30s",
						"scale_down_enabled":            "true",
						"expander":                      "priority",
						"scaler_type":                   "cluster-autoscaler",
						"skip_nodes_with_system_pods":   "true",
						"skip_nodes_with_local_storage": "true",
						"daemonset_eviction_for_nodes":  "true",
						"max_graceful_termination_sec":  "14400",
						"min_replica_count":             "1",
						"recycle_node_deletion_enabled": "true",
						"scale_up_from_zero":            "true",
						"priorities.%":                  "1",
						"priorities.10":                 acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

var baseMap = map[string]string{
	"cluster_id": acctest.CHECKSET,
}

func resourceCSAuoscalingConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_efficiency"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id           = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitch.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  vswitch_ids          = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 37)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  node_pool_name       = var.name
  vswitch_ids          = [local.vswitch_id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40

  scaling_config {
    enable   = true
    min_size = 0
    max_size = 10
  }
}

`, name)
}
