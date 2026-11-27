package cs_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cs"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_apiserver(t *testing.T) {
	var v *cs.Component

	resourceId := "alicloud_cs_kubernetes_addon.apiserver"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsKubernetesAddon"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		//CheckDestroy:  rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// import
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "kube-apiserver",
					"config":     "{\\\"requestTimeout\\\":\\\"2m0s\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   acctest.CHECKSET,
						"name":         "kube-apiserver",
						"version":      acctest.CHECKSET,
						"next_version": acctest.CHECKSET,
						"can_upgrade":  acctest.CHECKSET,
						"required":     acctest.CHECKSET,
					}),
				),
			},
			{
				// config
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"requestTimeout\\\":\\\"3m10s\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}

// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_terway_eniip(t *testing.T) {
	var v *cs.Component

	resourceId := "alicloud_cs_kubernetes_addon.terway-eniip"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsKubernetesAddon"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// import
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "terway-eniip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   acctest.CHECKSET,
						"name":         "terway-eniip",
						"version":      acctest.CHECKSET,
						"next_version": acctest.CHECKSET,
						"can_upgrade":  acctest.CHECKSET,
						"required":     acctest.CHECKSET,
					}),
				),
			},
			{
				// config
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"MaxPoolSize\\\":3}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}

// create new addon, upgrade, update and delete
// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_logtail_ds(t *testing.T) {
	var v *cs.Component

	resourceId := "alicloud_cs_kubernetes_addon.logtail-ds"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsKubernetesAddon"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "logtail-ds",
					"version":    "v1.8.5.0-aliyun",
					"depends_on": []string{"alicloud_cs_kubernetes_node_pool.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   acctest.CHECKSET,
						"name":         "logtail-ds",
						"version":      "v1.8.5.0-aliyun",
						"next_version": acctest.CHECKSET,
						"can_upgrade":  acctest.CHECKSET,
						"required":     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "v1.8.6.0-aliyun",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "v1.8.6.0-aliyun",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"LogtailDSLimitCPU\\\":\\\"3\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}

// create new addon, upgrade, update and delete
// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_ack_node_problem_detector(t *testing.T) {
	var v *cs.Component

	resourceId := "alicloud_cs_kubernetes_addon.ack-node-problem-detector"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsKubernetesAddon"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "ack-node-problem-detector",
					"depends_on": []string{"alicloud_cs_kubernetes_node_pool.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":   acctest.CHECKSET,
						"name":         "ack-node-problem-detector",
						"next_version": acctest.CHECKSET,
						"can_upgrade":  acctest.CHECKSET,
						"required":     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"npd\\\":{\\\"resources\\\":{\\\"limits\\\":{\\\"cpu\\\":\\\"2.0\\\",\\\"memory\\\":\\\"800Mi\\\"}}}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": "{\\\"npd\\\":{\\\"resources\\\":{\\\"limits\\\":{\\\"cpu\\\":\\\"1.0\\\",\\\"memory\\\":\\\"800Mi\\\"}}}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
	})
}

// create new addon, and delete with cleanup_cloud_resources
// lintignore: AT001
func TestAccAliCloudCSKubernetesAddon_vk(t *testing.T) {
	var v *cs.Component

	resourceId := "alicloud_cs_kubernetes_addon.ack-virtual-node"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, csdKubernetesAddonBasicMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsKubernetesAddon"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAddon-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAddonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${local.cluster_id}",
					"name":       "ack-virtual-node",
					//"version":    "v2.10.5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id": acctest.CHECKSET,
						"name":       "ack-virtual-node",
						//"version":      "v2.10.5",
						"next_version": acctest.CHECKSET,
						"can_upgrade":  acctest.CHECKSET,
						"required":     acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cleanup_cloud_resources": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

// TestUnitCSKubernetesAddonTimeoutsPlanConsistency reproduces, at the gRPC
// protocol level, the plan serialization asymmetry behind the Terraform apply
// error "Provider produced inconsistent final plan ... .timeouts: was absent,
// but now present" reported for alicloud_cs_kubernetes_addon.
//
// Terraform Core keeps the planned state produced during planning and, while
// applying, re-expands the plan for each instance once configuration
// expressions that were unknown at plan time have been resolved. helper/schema
// serializes the declared timeouts block into the planned state only when the
// proposed value is wholly known (copyTimeoutValues drops it otherwise), so a
// timeouts block that is wholly unknown at plan time - for example a dynamic
// block whose for_each depends on another resource created in the same apply -
// leaves the saved plan with .timeouts null. When the instance is re-planned
// during apply with the resolved timeouts values, the block is copied and the
// re-expanded plan has .timeouts present. Terraform Core compares the two and
// rejects the difference with ".timeouts: was absent, but now present".
//
// The "original plan" case asserts the wholly-unknown timeouts block is
// silently dropped, the "apply-time re-expansion" case asserts the resolved
// block is present: the two planned states differ only in the timeouts
// known-ness of the configuration, which is the inconsistency Core reports.
// The last case records the related variant where a single timeouts attribute
// is unknown at plan time: planning fails outright with a duration parsing
// error, so that variant cannot reach the apply-time error.
//
// These expectations document current SDK behavior while it reproduces the
// defect; an SDK fix that preserves the timeouts block would flip them.
func TestUnitCSKubernetesAddonTimeoutsPlanConsistency(t *testing.T) {
	const typeName = "alicloud_cs_kubernetes_addon"

	server := (&schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			typeName: cs.ResourceCSKubernetesAddon(),
		},
	}).GRPCProvider()

	res := cs.ResourceCSKubernetesAddon()
	objType := res.CoreConfigSchema().ImpliedType()
	timeoutsType := objType.AttributeTypes()["timeouts"]

	// A configuration for creating the node-local-dns addon: required
	// arguments set, every other attribute null.
	base := make(map[string]cty.Value, len(objType.AttributeTypes()))
	for name, attrType := range objType.AttributeTypes() {
		base[name] = cty.NullVal(attrType)
	}
	base["cluster_id"] = cty.StringVal("cluster-unit-test")
	base["name"] = cty.StringVal("node-local-dns")
	delete(base, "timeouts")

	withTimeouts := func(timeouts cty.Value) cty.Value {
		attrs := make(map[string]cty.Value, len(base)+1)
		for name, value := range base {
			attrs[name] = value
		}
		attrs["timeouts"] = timeouts
		return cty.ObjectVal(attrs)
	}

	// The resolved timeouts block, as decoded from `timeouts { create = "30m" }`.
	resolved := make(map[string]cty.Value, len(timeoutsType.AttributeTypes()))
	for name, attrType := range timeoutsType.AttributeTypes() {
		resolved[name] = cty.NullVal(attrType)
	}
	resolved["create"] = cty.StringVal("30m")

	marshal := func(v cty.Value) *tfprotov5.DynamicValue {
		t.Helper()
		raw, err := msgpack.Marshal(v, objType)
		if err != nil {
			t.Fatalf("marshaling value: %s", err)
		}
		return &tfprotov5.DynamicValue{MsgPack: raw}
	}

	nullPriorState := marshal(cty.NullVal(objType))

	plan := func(config cty.Value) (*tfprotov5.PlanResourceChangeResponse, cty.Value) {
		t.Helper()
		resp, err := server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
			TypeName:         typeName,
			PriorState:       nullPriorState,
			ProposedNewState: marshal(config),
			Config:           marshal(config),
		})
		if err != nil {
			t.Fatalf("PlanResourceChange call failed: %s", err)
		}
		planned := cty.NullVal(objType)
		if resp.PlannedState != nil && len(resp.PlannedState.MsgPack) > 0 {
			planned, err = msgpack.Unmarshal(resp.PlannedState.MsgPack, objType)
			if err != nil {
				t.Fatalf("unmarshaling planned state: %s", err)
			}
		}
		return resp, planned
	}

	diagnosticsText := func(resp *tfprotov5.PlanResourceChangeResponse) string {
		var parts []string
		for _, d := range resp.Diagnostics {
			if d.Summary != "" {
				parts = append(parts, d.Summary)
			}
			if d.Detail != "" {
				parts = append(parts, d.Detail)
			}
		}
		return strings.Join(parts, "; ")
	}

	t.Run("original plan drops wholly unknown timeouts", func(t *testing.T) {
		resp, planned := plan(withTimeouts(cty.UnknownVal(timeoutsType)))

		if len(resp.Diagnostics) > 0 {
			t.Fatalf("planning with wholly unknown timeouts returned diagnostics: %s", diagnosticsText(resp))
		}
		if !planned.GetAttr("timeouts").IsNull() {
			t.Fatalf("planned timeouts = %#v, want null: the block is dropped while it is not wholly known", planned.GetAttr("timeouts"))
		}
	})

	t.Run("apply-time re-expansion keeps resolved timeouts", func(t *testing.T) {
		resp, planned := plan(withTimeouts(cty.ObjectVal(resolved)))

		if len(resp.Diagnostics) > 0 {
			t.Fatalf("planning with resolved timeouts returned diagnostics: %s", diagnosticsText(resp))
		}
		timeouts := planned.GetAttr("timeouts")
		if timeouts.IsNull() {
			t.Fatal("planned timeouts is null, want the configured block to be copied into the planned state")
		}
		if got := timeouts.GetAttr("create"); !got.RawEquals(cty.StringVal("30m")) {
			t.Fatalf("planned timeouts.create = %#v, want \"30m\"", got)
		}
	})

	t.Run("partially unknown timeouts fail planning with a duration error", func(t *testing.T) {
		unknown := make(map[string]cty.Value, len(resolved))
		for name, value := range resolved {
			unknown[name] = value
		}
		unknown["create"] = cty.UnknownVal(cty.String)

		resp, _ := plan(withTimeouts(cty.ObjectVal(unknown)))

		if len(resp.Diagnostics) == 0 {
			t.Fatal("planning with a partially unknown timeouts block returned no diagnostics, want a duration parsing error")
		}
		if text := diagnosticsText(resp); !strings.Contains(text, "Error parsing") {
			t.Fatalf("unexpected diagnostics for a partially unknown timeouts block: %s", text)
		}
	})
}

var csdKubernetesAddonBasicMap = map[string]string{
	"cluster_id": acctest.CHECKSET,
}

func resourceCSAddonConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_instance_types" "default" {
  cpu_core_count    = 4
  memory_size       = 8
  eni_amount        = 4
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name = var.name
}


resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix                  = var.name
  cluster_spec                 = "ack.pro.small"
  vswitch_ids                  = [local.vswitch_id]
  pod_vswitch_ids              = [local.vswitch_id]
  new_nat_gateway              = false
  service_cidr                 = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled         = false
  is_enterprise_security_group = true
  addons {
    name = "terway-eniip"
  }
  delete_options {
    delete_mode   = "delete"
    resource_type = "ALB"
  }

  delete_options {
    delete_mode   = "delete"
    resource_type = "SLB"
  }

  delete_options {
    delete_mode   = "delete"
    resource_type = "SLS_Data"
  }

  delete_options {
    delete_mode   = "delete"
    resource_type = "SLS_ControlPlane"
  }

  delete_options {
    delete_mode   = "delete"
    resource_type = "PrivateZone"
  }
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [local.vswitch_id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_essd"
  system_disk_size     = 40
  desired_size         = 2
}

locals {
  vswitch_id = alicloud_vswitch.vswitch.id
  cluster_id = alicloud_cs_managed_kubernetes.default.id
}
`, name)
}
