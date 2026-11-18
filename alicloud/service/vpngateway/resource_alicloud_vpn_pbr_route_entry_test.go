package vpngateway_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpngateway"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccAliCloudVPNPbrRouteEntry_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_pbr_route_entry.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpnPbrRouteEntryMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vpngateway.NewVpcService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpnPbrRouteEntry", "weight", "priority")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%spbr%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpnPbrRouteEntryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_id": "${local.vpn_gateway_id}",
					"route_dest":     "10.0.0.0/24",
					"next_hop":       "${alicloud_vpn_connection.default.id}",
					"weight":         "100",
					"route_source":   "192.168.1.0/24",
					"publish_vpc":    "false",
					"description":    "tf-testacc-desc",
					"priority":       "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_source":   "192.168.1.0/24",
						"route_dest":     "10.0.0.0/24",
						"weight":         "100",
						"publish_vpc":    "false",
						"next_hop":       acctest.CHECKSET,
						"vpn_gateway_id": acctest.CHECKSET,
						"description":    "tf-testacc-desc",
						"priority":       "20",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"publish_vpc": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"publish_vpc": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight":      "100",
					"publish_vpc": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight":      "100",
						"publish_vpc": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
			{
				// Update priority via ModifyVpnPbrRouteEntryPriority (old/new pattern).
				Config: testAccConfig(map[string]interface{}{
					"priority": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "50",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
			{
				// Update weight and priority in the same apply: the weight request
				// carries the old priority and the priority request carries the new
				// weight, so both asynchronous modifications must converge.
				Config: testAccConfig(map[string]interface{}{
					"weight":   "0",
					"priority": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight":   "0",
						"priority": "30",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
			{
				// description is ForceNew: changing it must replace the route entry.
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testacc-desc-replaced",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testacc-desc-replaced",
						"weight":      "0",
						"priority":    "30",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "30",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: testAccConfig(map[string]interface{}{
					"description": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": acctest.REMOVEKEY,
						// Replacement Create omits priority, so the new entry gets the server default.
						"priority": "10",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
		},
	})
}

// TestAccAliCloudVPNPbrRouteEntry_defaultPriority covers the server-side default
// priority: the entry is created without priority, then its weight is updated and
// the entry deleted while the priority value is read back from the API.
func TestAccAliCloudVPNPbrRouteEntry_defaultPriority(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_pbr_route_entry.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpnPbrRouteEntryMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vpngateway.NewVpcService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpnPbrRouteEntry", "weight", "priority")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%spbrdefault%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpnPbrRouteEntryDefaultPriorityDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_id": "${local.vpn_gateway_id}",
					"route_dest":     "10.2.0.0/24",
					"next_hop":       "${alicloud_vpn_connection.default.id}",
					"weight":         "100",
					"route_source":   "192.168.2.0/24",
					"publish_vpc":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_source":   "192.168.2.0/24",
						"route_dest":     "10.2.0.0/24",
						"weight":         "100",
						"publish_vpc":    "false",
						"next_hop":       acctest.CHECKSET,
						"vpn_gateway_id": acctest.CHECKSET,
						"priority":       "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Weight update for an entry whose priority comes from the server.
				Config: testAccConfig(map[string]interface{}{
					"weight": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight":   "0",
						"priority": "10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudVpnPbrRouteEntryMap0 = map[string]string{}

func AlicloudVpnPbrRouteEntryBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
%s
resource "alicloud_vpn_customer_gateway" "defaultCustomerGateway" {
  description           = "defaultCustomerGateway"
  ip_address            = "2.2.2.25"
  asn                   = "2224"
  customer_gateway_name = var.name
}

resource "alicloud_vpn_connection" "default" {
  vpn_gateway_id      = local.vpn_gateway_id
  vpn_connection_name = var.name
  local_subnet = [
    "3.0.0.0/24"
  ]
  remote_subnet = [
    "10.0.0.0/24",
    "10.0.1.0/24"
  ]
  tags = {
    Created = "TF"
    For     = "example"
  }
  enable_tunnels_bgp = "true"
  tunnel_options_specification {
    tunnel_ipsec_config {
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes256"
      ipsec_lifetime = "16400"
      ipsec_pfs      = "group5"
    }

    customer_gateway_id = alicloud_vpn_customer_gateway.defaultCustomerGateway.id
    role                = "master"
    tunnel_bgp_config {
      local_asn    = "1219002"
      tunnel_cidr  = "169.254.30.0/30"
      local_bgp_ip = "169.254.30.1"
    }

    tunnel_ike_config {
      ike_mode     = "aggressive"
      ike_version  = "ikev2"
      local_id     = "localid_tunnel2"
      psk          = "12345678"
      remote_id    = "remote2"
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes256"
      ike_lifetime = "3600"
      ike_pfs      = "group14"
    }

  }
  tunnel_options_specification {
    tunnel_ike_config {
      remote_id    = "remote24"
      ike_enc_alg  = "aes256"
      ike_lifetime = "27000"
      ike_mode     = "aggressive"
      ike_pfs      = "group5"
      ike_auth_alg = "md5"
      ike_version  = "ikev2"
      local_id     = "localid_tunnel2"
      psk          = "12345678"
    }

    tunnel_ipsec_config {
      ipsec_lifetime = "2700"
      ipsec_pfs      = "group14"
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes256"
    }

    customer_gateway_id = alicloud_vpn_customer_gateway.defaultCustomerGateway.id
    role                = "slave"
    tunnel_bgp_config {
      local_asn    = "1219002"
      local_bgp_ip = "169.254.40.1"
      tunnel_cidr  = "169.254.40.0/30"
    }
  }
}
`, name, AlicloudVpnPbrRouteEntryGatewaySelection)
}

func AlicloudVpnPbrRouteEntryDefaultPriorityDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
%s
resource "alicloud_vpn_customer_gateway" "defaultCustomerGateway" {
  description           = "defaultCustomerGateway"
  ip_address            = "2.2.2.26"
  asn                   = "2226"
  customer_gateway_name = var.name
}

resource "alicloud_vpn_connection" "default" {
  vpn_gateway_id      = local.vpn_gateway_id
  vpn_connection_name = var.name
  local_subnet = [
    "4.0.0.0/24"
  ]
  remote_subnet = [
    "10.2.0.0/24",
    "10.2.1.0/24"
  ]
  tags = {
    Created = "TF"
    For     = "example"
  }
  enable_tunnels_bgp = "true"
  tunnel_options_specification {
    tunnel_ipsec_config {
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes256"
      ipsec_lifetime = "16400"
      ipsec_pfs      = "group5"
    }

    customer_gateway_id = alicloud_vpn_customer_gateway.defaultCustomerGateway.id
    role                = "master"
    tunnel_bgp_config {
      local_asn    = "1219002"
      tunnel_cidr  = "169.254.50.0/30"
      local_bgp_ip = "169.254.50.1"
    }

    tunnel_ike_config {
      ike_mode     = "aggressive"
      ike_version  = "ikev2"
      local_id     = "localid_tunnel3"
      psk          = "12345678"
      remote_id    = "remote3"
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes256"
      ike_lifetime = "3600"
      ike_pfs      = "group14"
    }

  }
  tunnel_options_specification {
    tunnel_ike_config {
      remote_id    = "remote34"
      ike_enc_alg  = "aes256"
      ike_lifetime = "27000"
      ike_mode     = "aggressive"
      ike_pfs      = "group5"
      ike_auth_alg = "md5"
      ike_version  = "ikev2"
      local_id     = "localid_tunnel3"
      psk          = "12345678"
    }

    tunnel_ipsec_config {
      ipsec_lifetime = "2700"
      ipsec_pfs      = "group14"
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes256"
    }

    customer_gateway_id = alicloud_vpn_customer_gateway.defaultCustomerGateway.id
    role                = "slave"
    tunnel_bgp_config {
      local_asn    = "1219002"
      local_bgp_ip = "169.254.60.1"
      tunnel_cidr  = "169.254.60.0/30"
    }
  }
}
`, name, AlicloudVpnPbrRouteEntryGatewaySelection)
}

// AlicloudVpnPbrRouteEntryGatewaySelection resolves the task-owned fixture by
// its reusable test tag. one() rejects zero or multiple matches instead of
// selecting an arbitrary gateway.
const AlicloudVpnPbrRouteEntryGatewaySelection = `
data "alicloud_vpn_gateways" "default" {
}

locals {
  vpn_gateway_ids = [
    for gateway in data.alicloud_vpn_gateways.default.gateways : gateway.id
    if lookup(gateway.tags, "TerraformTestFixture", "") == "vpn-pbr-route-entry"
  ]
  vpn_gateway_id = one(local.vpn_gateway_ids)
}
`

// lintignore: R001

func TestUnitVpnPbrRouteEntryDescriptionValidationAndCreateRequest(t *testing.T) {
	for _, tc := range []struct {
		name            string
		description     *string
		wantErr         string
		wantDescription *string
	}{
		{
			name: "explicit empty is rejected before create",
			description: func() *string {
				value := ""
				return &value
			}(),
			wantErr: "description",
		},
		{
			name: "omitted description is absent from create request",
		},
		{
			name: "valid description is sent in create request",
			description: func() *string {
				value := "unit-test-description"
				return &value
			}(),
			wantDescription: func() *string {
				value := "unit-test-description"
				return &value
			}(),
		},
		{
			name: "one hundred unicode characters are sent in create request",
			description: func() *string {
				value := strings.Repeat("界", 100)
				return &value
			}(),
			wantDescription: func() *string {
				value := strings.Repeat("界", 100)
				return &value
			}(),
		},
		{
			name: "description over one hundred characters is rejected before create",
			description: func() *string {
				value := strings.Repeat("x", 101)
				return &value
			}(),
			wantErr: "description",
		},
		{
			name: "one hundred one unicode characters are rejected before create",
			description: func() *string {
				value := strings.Repeat("界", 101)
				return &value
			}(),
			wantErr: "description",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			createCalls := 0
			var createForm map[string][]string
			client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse RPC request: %v", err)
					return
				}
				switch r.Form.Get("Action") {
				case "CreateVpnPbrRouteEntry":
					createCalls++
					createForm = r.Form
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock"})
				case "DescribeVpnPbrRouteEntries":
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(
						vpnPbrRouteEntryUnitEntry(0, 20),
						vpnPbrRouteEntryUnitEntry(100, 10),
					))
				default:
					t.Errorf("unexpected RPC action: %s", r.Form.Get("Action"))
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{"Code": "UnexpectedAction"})
				}
			})
			r := vpnPbrRouteEntryUnitResource(t)
			config := vpnPbrRouteEntryUnitConfig()
			delete(config, "priority")
			delete(config, "description")
			if tc.description != nil {
				config["description"] = *tc.description
			}

			validationDiags := r.Validate(sdktf.NewResourceConfigRaw(config))
			if tc.wantErr != "" {
				if !validationDiags.HasError() {
					t.Fatal("Validate succeeded, want a description length validation error")
				}
				diagnosticText := validationDiags[0].Summary + " " + validationDiags[0].Detail
				if !strings.Contains(diagnosticText, tc.wantErr) || !strings.Contains(diagnosticText, "length") {
					t.Fatalf("Validate diagnostics = %s, want a description length validation error", diagnosticText)
				}
				if createCalls != 0 {
					t.Fatalf("Create calls = %d, want no API request", createCalls)
				}
				return
			}
			if validationDiags.HasError() {
				t.Fatalf("Validate diagnostics = %v", validationDiags)
			}
			diff, err := r.Diff(context.Background(), nil, sdktf.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			diff.RawConfig = vpnPbrRouteEntryUnitRawConfigValue(r, config)
			if diff.RawConfig == cty.NilVal || diff.RawConfig.IsNull() {
				t.Fatal("Diff RawConfig is null; the Create identity path is not exercised")
			}
			if !diff.RawConfig.GetAttr("priority").IsNull() {
				t.Fatal("Diff RawConfig Priority must be null when omitted")
			}

			state, diags := r.Apply(context.Background(), nil, diff, client)
			if diags.HasError() {
				t.Fatalf("Apply diagnostics = %v", diags)
			}
			if createCalls != 1 {
				t.Fatalf("Create calls = %d, want 1", createCalls)
			}
			if state.Attributes["weight"] != "100" || state.Attributes["priority"] != "10" {
				t.Fatalf("Create Read selected weight/priority = %q/%q, want 100/10", state.Attributes["weight"], state.Attributes["priority"])
			}
			values, present := createForm["Description"]
			if tc.wantDescription == nil {
				if present {
					t.Fatalf("Create Description = %q, want the parameter omitted", values)
				}
			} else if !present || len(values) != 1 || values[0] != *tc.wantDescription {
				t.Fatalf("Create Description = %q (present=%v), want %q", values, present, *tc.wantDescription)
			}
		})
	}
}

func TestUnitVpnPbrRouteEntryConfiguredToOmittedTransitions(t *testing.T) {
	t.Run("description replacement omits the create parameter", func(t *testing.T) {
		createCalls := 0
		deleteCalls := 0
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Errorf("parse RPC request: %v", err)
				return
			}
			switch r.Form.Get("Action") {
			case "CreateVpnPbrRouteEntry":
				createCalls++
				if values, present := r.Form["Description"]; present {
					t.Errorf("Replacement Create Description = %q, want the parameter omitted", values)
				}
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock"})
			case "DeleteVpnPbrRouteEntry":
				deleteCalls++
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock"})
			case "DescribeVpnPbrRouteEntries":
				if deleteCalls > 0 && createCalls == 0 {
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse())
					return
				}
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(vpnPbrRouteEntryUnitEntry(100, 20)))
			default:
				t.Errorf("unexpected RPC action: %s", r.Form.Get("Action"))
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{"Code": "UnexpectedAction"})
			}
		})
		r := vpnPbrRouteEntryUnitResource(t)
		oldConfig := vpnPbrRouteEntryUnitConfig()
		old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
		old.SetId("vpn-unit:vco-unit:192.168.1.0/24:10.0.0.0/24")
		oldState := old.State()
		oldState.RawConfig = vpnPbrRouteEntryUnitRawConfigValue(r, oldConfig)
		oldState.RawState = vpnPbrRouteEntryUnitRawStateValue(r, oldConfig, oldState.ID)
		newConfig := vpnPbrRouteEntryUnitConfig()
		delete(newConfig, "description")
		diff, err := r.Diff(context.Background(), oldState, sdktf.NewResourceConfigRaw(newConfig), nil)
		if err != nil {
			t.Fatal(err)
		}
		diff.RawConfig = vpnPbrRouteEntryUnitRawConfigValue(r, newConfig)
		if !diff.RawConfig.GetAttr("description").IsNull() {
			t.Fatal("Diff RawConfig Description is not null after removal")
		}
		if !diff.RequiresNew() {
			t.Fatal("Removing configured description did not require replacement")
		}
		state, diags := r.Apply(context.Background(), oldState, diff, client)
		if diags.HasError() {
			t.Fatalf("Apply diagnostics = %v", diags)
		}
		if deleteCalls != 1 || createCalls != 1 {
			t.Fatalf("Delete/Create calls = %d/%d, want 1/1", deleteCalls, createCalls)
		}
		if got := state.Attributes["description"]; got != "" {
			t.Fatalf("Replacement state description = %q, want empty after omission", got)
		}
	})

	t.Run("priority removal preserves the remote value without modify", func(t *testing.T) {
		describeCalls := 0
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Errorf("parse RPC request: %v", err)
				return
			}
			if r.Form.Get("Action") != "DescribeVpnPbrRouteEntries" {
				t.Errorf("unexpected RPC action: %s", r.Form.Get("Action"))
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{"Code": "UnexpectedAction"})
				return
			}
			describeCalls++
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(vpnPbrRouteEntryUnitEntry(100, 20)))
		})
		r := vpnPbrRouteEntryUnitResource(t)
		oldConfig := vpnPbrRouteEntryUnitConfig()
		old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
		old.SetId("vpn-unit:vco-unit:192.168.1.0/24:10.0.0.0/24")
		oldState := old.State()
		oldState.RawConfig = vpnPbrRouteEntryUnitRawConfigValue(r, oldConfig)
		oldState.RawState = vpnPbrRouteEntryUnitRawStateValue(r, oldConfig, oldState.ID)
		newConfig := vpnPbrRouteEntryUnitConfig()
		delete(newConfig, "priority")
		diff, err := r.Diff(context.Background(), oldState, sdktf.NewResourceConfigRaw(newConfig), nil)
		if err != nil {
			t.Fatal(err)
		}
		if diff != nil {
			t.Fatalf("Priority removal produced a diff = %#v, want no update", diff)
		}
		newRawConfig := vpnPbrRouteEntryUnitRawConfigValue(r, newConfig)
		if !newRawConfig.GetAttr("priority").IsNull() {
			t.Fatal("RawConfig Priority is not null after removal")
		}
		provider := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{"alicloud_vpn_pbr_route_entry": r},
		}
		provider.SetMeta(client)
		server := schema.NewGRPCProviderServer(provider)
		refreshState := vpnPbrRouteEntryUnitRawStateValue(r, newConfig, oldState.ID)
		stateBytes, err := msgpack.Marshal(refreshState, r.CoreConfigSchema().ImpliedType())
		if err != nil {
			t.Fatal(err)
		}
		read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
			TypeName:     "alicloud_vpn_pbr_route_entry",
			CurrentState: &tfprotov5.DynamicValue{MsgPack: stateBytes},
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(read.Diagnostics) != 0 {
			t.Fatalf("Read diagnostics = %v", read.Diagnostics)
		}
		state, err := r.ShimInstanceStateFromValue(vpnPbrRouteEntryUnitDynamicValue(t, r, read.NewState))
		if err != nil {
			t.Fatal(err)
		}
		if describeCalls == 0 {
			t.Fatal("Priority removal did not refresh the remote value")
		}
		if state.Attributes["priority"] != "20" {
			t.Fatalf("Priority after removal = %q, want retained remote value 20", state.Attributes["priority"])
		}
	})
}

func TestUnitVpnPbrRouteEntryCreateReadUsesConfiguredIdentity(t *testing.T) {
	for _, weight := range []int{0, 100} {
		t.Run(fmt.Sprintf("weight %d with another same-priority candidate", weight), func(t *testing.T) {
			createCalls := 0
			client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse RPC request: %v", err)
					return
				}
				switch r.Form.Get("Action") {
				case "CreateVpnPbrRouteEntry":
					createCalls++
					if r.Form.Get("Weight") != fmt.Sprint(weight) || r.Form.Get("Priority") != "20" {
						t.Errorf("Create Weight/Priority = %q/%q, want %d/20", r.Form.Get("Weight"), r.Form.Get("Priority"), weight)
					}
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock"})
				case "DescribeVpnPbrRouteEntries":
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(
						vpnPbrRouteEntryUnitEntry(100-weight, 20),
						vpnPbrRouteEntryUnitEntry(weight, 20),
					))
				default:
					t.Errorf("unexpected RPC action: %s", r.Form.Get("Action"))
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{"Code": "UnexpectedAction"})
				}
			})
			r := vpnPbrRouteEntryUnitResource(t)
			config := vpnPbrRouteEntryUnitConfig()
			config["weight"] = weight
			config["priority"] = 20
			diff, err := r.Diff(context.Background(), nil, sdktf.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			diff.RawConfig = vpnPbrRouteEntryUnitRawConfigValue(r, config)
			if diff.RawConfig == cty.NilVal || diff.RawConfig.IsNull() {
				t.Fatal("Diff RawConfig is null; the Create identity path is not exercised")
			}
			configuredWeight, _ := diff.RawConfig.GetAttr("weight").AsBigFloat().Int64()
			if configuredWeight != int64(weight) {
				t.Fatalf("Diff RawConfig Weight = %d, want %d", configuredWeight, weight)
			}
			state, diags := r.Apply(context.Background(), nil, diff, client)
			if diags.HasError() {
				t.Fatalf("Apply diagnostics = %v", diags)
			}
			if createCalls != 1 {
				t.Fatalf("Create calls = %d, want 1", createCalls)
			}
			if state.Attributes["weight"] != fmt.Sprint(weight) || state.Attributes["priority"] != "20" {
				t.Fatalf("Create Read selected weight/priority = %q/%q, want %d/20", state.Attributes["weight"], state.Attributes["priority"], weight)
			}
		})
	}
}

func TestUnitVpnPbrRouteEntryUpdateReadUsesConfiguredIdentity(t *testing.T) {
	weightCalls := 0
	priorityCalls := 0
	client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse RPC request: %v", err)
			return
		}
		switch r.Form.Get("Action") {
		case "ModifyVpnPbrRouteEntryWeight":
			weightCalls++
			if r.Form.Get("Weight") != "100" || r.Form.Get("NewWeight") != "0" || r.Form.Get("Priority") != "20" {
				t.Errorf("Weight update Weight/NewWeight/Priority = %q/%q/%q, want 100/0/20", r.Form.Get("Weight"), r.Form.Get("NewWeight"), r.Form.Get("Priority"))
			}
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock"})
		case "ModifyVpnPbrRouteEntryPriority":
			priorityCalls++
			if r.Form.Get("Weight") != "0" || r.Form.Get("Priority") != "20" || r.Form.Get("NewPriority") != "30" {
				t.Errorf("Priority update Weight/Priority/NewPriority = %q/%q/%q, want 0/20/30", r.Form.Get("Weight"), r.Form.Get("Priority"), r.Form.Get("NewPriority"))
			}
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock"})
		case "DescribeVpnGateway":
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock", "Status": "active"})
		case "DescribeVpnPbrRouteEntries":
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(
				vpnPbrRouteEntryUnitEntry(100, 30),
				vpnPbrRouteEntryUnitEntry(0, 30),
			))
		default:
			t.Errorf("unexpected RPC action: %s", r.Form.Get("Action"))
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{"Code": "UnexpectedAction"})
		}
	})
	r := vpnPbrRouteEntryUnitResource(t)
	oldConfig := vpnPbrRouteEntryUnitConfig()
	oldConfig["weight"] = 100
	oldConfig["priority"] = 20
	old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
	old.SetId("vpn-unit:vco-unit:192.168.1.0/24:10.0.0.0/24")
	oldState := old.State()
	oldState.RawConfig = vpnPbrRouteEntryUnitRawConfigValue(r, oldConfig)
	oldState.RawState = vpnPbrRouteEntryUnitRawStateValue(r, oldConfig, oldState.ID)
	newConfig := vpnPbrRouteEntryUnitConfig()
	newConfig["weight"] = 0
	newConfig["priority"] = 30
	diff, err := r.Diff(context.Background(), oldState, sdktf.NewResourceConfigRaw(newConfig), nil)
	if err != nil {
		t.Fatal(err)
	}
	diff.RawConfig = vpnPbrRouteEntryUnitRawConfigValue(r, newConfig)
	if diff.RawConfig == cty.NilVal || diff.RawConfig.IsNull() {
		t.Fatal("Diff RawConfig is null; the Update identity path is not exercised")
	}
	if diff.RawState == cty.NilVal || diff.RawState.IsNull() {
		t.Fatal("Diff RawState is null; the pre-Update identity path is not exercised")
	}
	previousWeight, _ := diff.RawState.GetAttr("weight").AsBigFloat().Int64()
	if previousWeight != 100 {
		t.Fatalf("Diff RawState Weight = %d, want 100", previousWeight)
	}
	currentWeight, _ := diff.RawConfig.GetAttr("weight").AsBigFloat().Int64()
	currentPriority, _ := diff.RawConfig.GetAttr("priority").AsBigFloat().Int64()
	if currentWeight != 0 || currentPriority != 30 {
		t.Fatalf("Diff RawConfig Weight/Priority = %d/%d, want 0/30", currentWeight, currentPriority)
	}
	state, diags := r.Apply(context.Background(), oldState, diff, client)
	if diags.HasError() {
		t.Fatalf("Apply diagnostics = %v", diags)
	}
	if weightCalls != 1 || priorityCalls != 1 {
		t.Fatalf("Weight/Priority update calls = %d/%d, want 1/1", weightCalls, priorityCalls)
	}
	if state.Attributes["weight"] != "0" || state.Attributes["priority"] != "30" {
		t.Fatalf("Update Read selected weight/priority = %q/%q, want 0/30", state.Attributes["weight"], state.Attributes["priority"])
	}
}

func TestUnitVpnPbrRouteEntryReadRestoresWeightIdentity(t *testing.T) {
	for _, weight := range []int{0, 100} {
		t.Run(fmt.Sprintf("import weight %d", weight), func(t *testing.T) {
			client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(vpnPbrRouteEntryUnitEntry(weight, 20)))
			})
			r := vpnPbrRouteEntryUnitResource(t)
			provider := &schema.Provider{
				ResourcesMap: map[string]*schema.Resource{"alicloud_vpn_pbr_route_entry": r},
			}
			provider.SetMeta(client)
			server := schema.NewGRPCProviderServer(provider)
			imported, err := server.ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
				TypeName: "alicloud_vpn_pbr_route_entry",
				ID:       "vpn-unit:vco-unit:192.168.1.0/24:10.0.0.0/24",
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(imported.ImportedResources) != 1 {
				t.Fatalf("Imported resources = %d, want 1", len(imported.ImportedResources))
			}
			read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
				TypeName:     "alicloud_vpn_pbr_route_entry",
				CurrentState: imported.ImportedResources[0].State,
				Private:      imported.ImportedResources[0].Private,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(read.Diagnostics) != 0 {
				t.Fatalf("Read diagnostics = %v", read.Diagnostics)
			}
			state, err := r.ShimInstanceStateFromValue(vpnPbrRouteEntryUnitDynamicValue(t, r, read.NewState))
			if err != nil {
				t.Fatal(err)
			}
			if state.Attributes["weight"] != fmt.Sprint(weight) || state.Attributes["priority"] != "20" {
				t.Fatalf("Read restored weight/priority = %q/%q, want %d/20", state.Attributes["weight"], state.Attributes["priority"], weight)
			}
		})
	}

	t.Run("normal refresh keeps weight zero identity", func(t *testing.T) {
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(
				vpnPbrRouteEntryUnitEntry(100, 20),
				vpnPbrRouteEntryUnitEntry(0, 10),
			))
		})
		r := vpnPbrRouteEntryUnitResource(t)
		provider := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{"alicloud_vpn_pbr_route_entry": r},
		}
		provider.SetMeta(client)
		server := schema.NewGRPCProviderServer(provider)
		d := schema.TestResourceDataRaw(t, r.Schema, vpnPbrRouteEntryUnitConfig())
		d.SetId("vpn-unit:vco-unit:192.168.1.0/24:10.0.0.0/24")
		state := d.State()
		state.Attributes["weight"] = "0"
		state.Attributes["priority"] = "10"
		stateValue := cty.ObjectVal(map[string]cty.Value{
			"id":             cty.StringVal(state.ID),
			"vpn_gateway_id": cty.StringVal("vpn-unit"),
			"next_hop":       cty.StringVal("vco-unit"),
			"route_source":   cty.StringVal("192.168.1.0/24"),
			"route_dest":     cty.StringVal("10.0.0.0/24"),
			"publish_vpc":    cty.False,
			"weight":         cty.NumberIntVal(0),
			"priority":       cty.NumberIntVal(10),
			"description":    cty.StringVal("unit-test"),
			"status":         cty.StringVal("normal"),
			"timeouts":       cty.NullVal(r.CoreConfigSchema().ImpliedType().AttributeType("timeouts")),
		})
		stateBytes, err := msgpack.Marshal(stateValue, r.CoreConfigSchema().ImpliedType())
		if err != nil {
			t.Fatal(err)
		}
		read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
			TypeName:     "alicloud_vpn_pbr_route_entry",
			CurrentState: &tfprotov5.DynamicValue{MsgPack: stateBytes},
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(read.Diagnostics) != 0 {
			t.Fatalf("Read diagnostics = %v", read.Diagnostics)
		}
		refreshed, err := r.ShimInstanceStateFromValue(vpnPbrRouteEntryUnitDynamicValue(t, r, read.NewState))
		if err != nil {
			t.Fatal(err)
		}
		if refreshed.Attributes["weight"] != "0" || refreshed.Attributes["priority"] != "10" {
			t.Fatalf("Refresh restored weight/priority = %q/%q, want 0/10", refreshed.Attributes["weight"], refreshed.Attributes["priority"])
		}
	})

	for _, weight := range []int{0, 100} {
		t.Run(fmt.Sprintf("legacy state without priority selects weight %d through public read", weight), func(t *testing.T) {
			describeCalls := 0
			client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse RPC request: %v", err)
					return
				}
				if r.Form.Get("Action") != "DescribeVpnPbrRouteEntries" {
					t.Errorf("unexpected RPC action: %s", r.Form.Get("Action"))
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{"Code": "UnexpectedAction"})
					return
				}
				describeCalls++
				if r.Form.Get("VpnGatewayId") != "vpn-unit" || r.Form.Get("PageNumber") != "1" {
					t.Errorf("Describe VpnGatewayId/PageNumber = %q/%q, want vpn-unit/1", r.Form.Get("VpnGatewayId"), r.Form.Get("PageNumber"))
				}
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(
					vpnPbrRouteEntryUnitEntry(100-weight, 20),
					vpnPbrRouteEntryUnitEntry(weight, 10),
				))
			})
			r := vpnPbrRouteEntryUnitResource(t)
			provider := &schema.Provider{
				ResourcesMap: map[string]*schema.Resource{"alicloud_vpn_pbr_route_entry": r},
			}
			provider.SetMeta(client)
			server := schema.NewGRPCProviderServer(provider)
			legacyConfig := vpnPbrRouteEntryUnitConfig()
			legacyConfig["weight"] = weight
			delete(legacyConfig, "priority")
			legacyState := vpnPbrRouteEntryUnitRawStateValue(r, legacyConfig, "vpn-unit:vco-unit:192.168.1.0/24:10.0.0.0/24")
			if !legacyState.GetAttr("priority").IsNull() {
				t.Fatal("Legacy fixture Priority is not null; the RawState fallback is not exercised")
			}
			stateBytes, err := msgpack.Marshal(legacyState, r.CoreConfigSchema().ImpliedType())
			if err != nil {
				t.Fatal(err)
			}
			read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
				TypeName:     "alicloud_vpn_pbr_route_entry",
				CurrentState: &tfprotov5.DynamicValue{MsgPack: stateBytes},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(read.Diagnostics) != 0 {
				t.Fatalf("Read diagnostics = %v", read.Diagnostics)
			}
			state, err := r.ShimInstanceStateFromValue(vpnPbrRouteEntryUnitDynamicValue(t, r, read.NewState))
			if err != nil {
				t.Fatal(err)
			}
			if describeCalls == 0 {
				t.Fatal("Public Read did not issue DescribeVpnPbrRouteEntries")
			}
			if state.Attributes["weight"] != fmt.Sprint(weight) || state.Attributes["priority"] != "10" {
				t.Fatalf("Legacy Read selected weight/priority = %q/%q, want %d/10", state.Attributes["weight"], state.Attributes["priority"], weight)
			}
		})
	}

	t.Run("ambiguous import returns an error", func(t *testing.T) {
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(
				vpnPbrRouteEntryUnitEntry(0, 10),
				vpnPbrRouteEntryUnitEntry(100, 20),
			))
		})
		r := vpnPbrRouteEntryUnitResource(t)
		provider := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{"alicloud_vpn_pbr_route_entry": r},
		}
		provider.SetMeta(client)
		server := schema.NewGRPCProviderServer(provider)
		imported, err := server.ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
			TypeName: "alicloud_vpn_pbr_route_entry",
			ID:       "vpn-unit:vco-unit:192.168.1.0/24:10.0.0.0/24",
		})
		if err != nil {
			t.Fatal(err)
		}
		read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
			TypeName:     "alicloud_vpn_pbr_route_entry",
			CurrentState: imported.ImportedResources[0].State,
			Private:      imported.ImportedResources[0].Private,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(read.Diagnostics) == 0 {
			t.Fatal("Ambiguous Import returned no diagnostics, want a multiple-candidate error")
		}
		diagnosticText := read.Diagnostics[0].Summary + " " + read.Diagnostics[0].Detail
		if !strings.Contains(diagnosticText, "multiple") {
			t.Fatalf("Ambiguous Import diagnostics = %s, want a multiple-candidate error", diagnosticText)
		}
	})
}

func vpnPbrRouteEntryUnitDynamicValue(t *testing.T, r *schema.Resource, value *tfprotov5.DynamicValue) cty.Value {
	t.Helper()
	if value == nil || len(value.MsgPack) == 0 {
		t.Fatal("Read returned no MessagePack state")
	}
	state, err := msgpack.Unmarshal(value.MsgPack, r.CoreConfigSchema().ImpliedType())
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func TestUnitVpnPbrRouteEntryDeleteRetryClassification(t *testing.T) {
	for _, tc := range []struct {
		name              string
		failures          []string
		status            int
		wantErr           string
		wantDeleteCalls   int
		wantDescribeCalls int
	}{
		{
			name:              "configuring retries then succeeds",
			failures:          []string{"VpnRouteEntry.Configuring"},
			status:            http.StatusBadRequest,
			wantDeleteCalls:   2,
			wantDescribeCalls: 1,
		},
		{
			name:              "not found is idempotent success",
			failures:          []string{"InvalidPbrRouteEntry.NotFound"},
			status:            http.StatusNotFound,
			wantDeleteCalls:   1,
			wantDescribeCalls: 1,
		},
		{
			name:              "route entry not found is idempotent success",
			failures:          []string{"InvalidRouteEntry.NotFound"},
			status:            http.StatusBadRequest,
			wantDeleteCalls:   1,
			wantDescribeCalls: 1,
		},
		{
			name:              "other errors are terminal",
			failures:          []string{"InvalidParameter"},
			status:            http.StatusBadRequest,
			wantErr:           "InvalidParameter",
			wantDeleteCalls:   1,
			wantDescribeCalls: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			deleteCalls := 0
			describeCalls := 0
			client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse RPC request: %v", err)
					return
				}
				switch r.Form.Get("Action") {
				case "DeleteVpnPbrRouteEntry":
					deleteCalls++
					if deleteCalls <= len(tc.failures) {
						vpnPbrRouteEntryUnitWriteJSON(t, w, tc.status, map[string]interface{}{
							"RequestId": "mock",
							"Code":      tc.failures[deleteCalls-1],
							"Message":   "mock request failed",
						})
						return
					}
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock"})
				case "DescribeVpnPbrRouteEntries":
					describeCalls++
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{
						"RequestId": "mock",
						"VpnPbrRouteEntries": map[string]interface{}{
							"VpnPbrRouteEntry": []interface{}{},
						},
					})
				default:
					t.Errorf("unexpected RPC action: %s", r.Form.Get("Action"))
					vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{"Code": "UnexpectedAction"})
				}
			})
			r := vpnPbrRouteEntryUnitResource(t)
			r.Timeouts.Delete = schema.DefaultTimeout(12 * time.Second)
			d := schema.TestResourceDataRaw(t, r.Schema, vpnPbrRouteEntryUnitConfig())
			d.SetId("vpn-unit:vco-unit:192.168.1.0/24:10.0.0.0/24")

			err := r.Delete(d, client)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("Delete returned an unexpected error: %v", err)
			}
			if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Fatalf("Delete error = %v, want an error containing %q", err, tc.wantErr)
			}
			if deleteCalls != tc.wantDeleteCalls || describeCalls != tc.wantDescribeCalls {
				t.Fatalf("Delete/Describe calls = %d/%d, want %d/%d", deleteCalls, describeCalls, tc.wantDeleteCalls, tc.wantDescribeCalls)
			}
		})
	}
}

func TestUnitVpnPbrRouteEntryUpdateUsesOneTimeoutBudget(t *testing.T) {
	weightCalls := 0
	var weightClientTokens []string
	priorityCalls := 0
	gatewayDescribeCalls := 0
	routeDescribeCalls := 0
	client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse RPC request: %v", err)
			return
		}
		switch r.Form.Get("Action") {
		case "ModifyVpnPbrRouteEntryWeight":
			weightCalls++
			weightClientTokens = append(weightClientTokens, r.Form.Get("ClientToken"))
			if weightCalls == 1 {
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{
					"RequestId": "mock",
					"Code":      "VpnGateway.Configuring",
					"Message":   "mock request failed",
				})
				return
			}
			time.Sleep(2 * time.Second)
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock"})
		case "ModifyVpnPbrRouteEntryPriority":
			priorityCalls++
			time.Sleep(2 * time.Second)
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock"})
		case "DescribeVpnGateway":
			gatewayDescribeCalls++
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{"RequestId": "mock", "Status": "active"})
		case "DescribeVpnPbrRouteEntries":
			routeDescribeCalls++
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, map[string]interface{}{
				"RequestId": "mock",
				"VpnPbrRouteEntries": map[string]interface{}{
					"VpnPbrRouteEntry": []interface{}{map[string]interface{}{
						"VpnInstanceId": "vpn-unit",
						"NextHop":       "vco-unit",
						"RouteSource":   "192.168.1.0/24",
						"RouteDest":     "10.0.0.0/24",
						"Weight":        0,
						"Priority":      30,
						"State":         "normal",
					}},
				},
			})
		default:
			t.Errorf("unexpected RPC action: %s", r.Form.Get("Action"))
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusBadRequest, map[string]interface{}{"Code": "UnexpectedAction"})
		}
	})
	r := vpnPbrRouteEntryUnitResource(t)
	// Across the full 0-3s retry jitter, the first poll and priority call
	// complete within 13s, while the second poll starts no earlier than 14s.
	r.Timeouts.Update = schema.DefaultTimeout(13 * time.Second)
	oldConfig := vpnPbrRouteEntryUnitConfig()
	old := schema.TestResourceDataRaw(t, r.Schema, oldConfig)
	old.SetId("vpn-unit:vco-unit:192.168.1.0/24:10.0.0.0/24")
	oldState := old.State()
	oldState.RawConfig = vpnPbrRouteEntryUnitRawConfigValue(r, oldConfig)
	oldState.RawState = vpnPbrRouteEntryUnitRawStateValue(r, oldConfig, oldState.ID)
	newConfig := vpnPbrRouteEntryUnitConfig()
	newConfig["weight"] = 0
	newConfig["priority"] = 30
	diff, err := r.Diff(context.Background(), oldState, sdktf.NewResourceConfigRaw(newConfig), nil)
	if err != nil {
		t.Fatal(err)
	}
	diff.RawConfig = vpnPbrRouteEntryUnitRawConfigValue(r, newConfig)
	if diff.RawConfig == cty.NilVal || diff.RawConfig.IsNull() {
		t.Fatal("Diff RawConfig is null; the Update identity path is not exercised")
	}
	if diff.RawState == cty.NilVal || diff.RawState.IsNull() {
		t.Fatal("Diff RawState is null; the pre-Update identity path is not exercised")
	}
	previousWeight, _ := diff.RawState.GetAttr("weight").AsBigFloat().Int64()
	previousPriority, _ := diff.RawState.GetAttr("priority").AsBigFloat().Int64()
	currentWeight, _ := diff.RawConfig.GetAttr("weight").AsBigFloat().Int64()
	currentPriority, _ := diff.RawConfig.GetAttr("priority").AsBigFloat().Int64()
	if previousWeight != 100 || previousPriority != 20 || currentWeight != 0 || currentPriority != 30 {
		t.Fatalf("RawState Weight/Priority = %d/%d, RawConfig Weight/Priority = %d/%d; want 100/20 and 0/30", previousWeight, previousPriority, currentWeight, currentPriority)
	}

	_, diags := r.Apply(context.Background(), oldState, diff, client)
	if !diags.HasError() || !strings.Contains(strings.ToLower(fmt.Sprintf("%v", diags)), "timeout") {
		t.Fatalf("Update diagnostics = %v, want the shared timeout to expire", diags)
	}
	// The second gateway poll would require a fresh budget. A single budget
	// expires before that poll and prevents the final Read.
	if weightCalls != 2 || priorityCalls != 1 || gatewayDescribeCalls != 1 || routeDescribeCalls != 0 {
		t.Fatalf("weight/priority/gateway Describe/route Describe calls = %d/%d/%d/%d, want 2/1/1/0", weightCalls, priorityCalls, gatewayDescribeCalls, routeDescribeCalls)
	}
	if len(weightClientTokens) != 2 || weightClientTokens[0] == "" || weightClientTokens[0] != weightClientTokens[1] {
		t.Fatalf("weight ClientTokens = %q, want the same non-empty token across retries", weightClientTokens)
	}
}

func TestUnitVpnPbrRouteEntryDescribeMatchesIdentity(t *testing.T) {
	other := vpnPbrRouteEntryUnitEntry(0, 10)
	target := vpnPbrRouteEntryUnitEntry(100, 20)
	id := "vpn-unit:vco-unit:192.168.1.0/24:10.0.0.0/24"

	t.Run("same key selects weight and priority", func(t *testing.T) {
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(other, target))
		})
		object, err := vpngateway.NewVpcService(client).DiscoverVpnPbrRouteEntry(id, 100, 20)
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprint(object["Weight"]) != "100" || fmt.Sprint(object["Priority"]) != "20" {
			t.Fatalf("Describe returned weight/priority %v/%v, want 100/20", object["Weight"], object["Priority"])
		}
	})

	t.Run("later page selects target", func(t *testing.T) {
		firstPage := make([]interface{}, helper.PageSizeMedium)
		for i := range firstPage {
			firstPage[i] = other
		}
		requestedSecondPage := false
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Errorf("parse RPC request: %v", err)
				return
			}
			if r.Form.Get("Action") != "DescribeVpnPbrRouteEntries" {
				t.Errorf("unexpected RPC action: %s", r.Form.Get("Action"))
				return
			}
			switch r.Form.Get("PageNumber") {
			case "1":
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(firstPage...))
			case "2":
				requestedSecondPage = true
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(target))
			default:
				vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse())
			}
		})
		object, err := vpngateway.NewVpcService(client).DiscoverVpnPbrRouteEntry(id, 100, 20)
		if err != nil {
			t.Fatal(err)
		}
		if !requestedSecondPage || fmt.Sprint(object["Priority"]) != "20" {
			t.Fatalf("Describe second-page result = %v, requestedSecondPage=%v, want target priority 20", object, requestedSecondPage)
		}
	})

	t.Run("import without identity recovers weight 100", func(t *testing.T) {
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(target))
		})
		object, err := vpngateway.NewVpcService(client).DiscoverVpnPbrRouteEntry(id)
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprint(object["Weight"]) != "100" {
			t.Fatalf("Import recovered weight %v, want 100", object["Weight"])
		}
	})

	t.Run("import without identity recovers weight 0", func(t *testing.T) {
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(other))
		})
		object, err := vpngateway.NewVpcService(client).DiscoverVpnPbrRouteEntry(id)
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprint(object["Weight"]) != "0" {
			t.Fatalf("Import recovered weight %v, want 0", object["Weight"])
		}
	})

	t.Run("unique drift returns current object", func(t *testing.T) {
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(other))
		})
		object, err := vpngateway.NewVpcService(client).DiscoverVpnPbrRouteEntry(id, 100, 20)
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprint(object["Weight"]) != "0" || fmt.Sprint(object["Priority"]) != "10" {
			t.Fatalf("Drift returned weight/priority %v/%v, want current 0/10", object["Weight"], object["Priority"])
		}
	})

	t.Run("ambiguous import returns error", func(t *testing.T) {
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(other, target))
		})
		object, err := vpngateway.NewVpcService(client).DiscoverVpnPbrRouteEntry(id)
		if err == nil || errs.NotFoundError(err) {
			t.Fatalf("Ambiguous Import error = %v, want an ambiguity error rather than NotFound", err)
		}
		if object != nil {
			t.Fatalf("Ambiguous Import returned %#v, want no object", object)
		}
	})

	t.Run("remaining entry does not hide not found", func(t *testing.T) {
		client := vpnPbrRouteEntryUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			vpnPbrRouteEntryUnitWriteJSON(t, w, http.StatusOK, vpnPbrRouteEntryUnitEntriesResponse(other))
		})
		object, err := vpngateway.NewVpcService(client).DescribeVpnPbrRouteEntry(id, 100, 20)
		if !errs.NotFoundError(err) {
			t.Fatalf("Describe error = %v, want NotFound when only another entry remains", err)
		}
		if object != nil {
			t.Fatalf("Describe returned %#v, want no object for a missing target", object)
		}
	})
}

func vpnPbrRouteEntryUnitEntry(weight, priority int) map[string]interface{} {
	return map[string]interface{}{
		"VpnInstanceId": "vpn-unit",
		"NextHop":       "vco-unit",
		"RouteSource":   "192.168.1.0/24",
		"RouteDest":     "10.0.0.0/24",
		"Weight":        weight,
		"Priority":      priority,
		"State":         "normal",
	}
}

func vpnPbrRouteEntryUnitEntriesResponse(entries ...interface{}) map[string]interface{} {
	if entries == nil {
		entries = []interface{}{}
	}
	return map[string]interface{}{
		"RequestId": "mock",
		"VpnPbrRouteEntries": map[string]interface{}{
			"VpnPbrRouteEntry": entries,
		},
	}
}

func vpnPbrRouteEntryUnitResource(t *testing.T) *schema.Resource {
	t.Helper()
	for _, registered := range vpngateway.ServicePackage().SDKResources {
		if registered.TypeName == "alicloud_vpn_pbr_route_entry" {
			return registered.Factory()
		}
	}
	t.Fatal("alicloud_vpn_pbr_route_entry is not registered")
	return nil
}

func vpnPbrRouteEntryUnitConfig() map[string]interface{} {
	return map[string]interface{}{
		"vpn_gateway_id": "vpn-unit",
		"next_hop":       "vco-unit",
		"route_source":   "192.168.1.0/24",
		"route_dest":     "10.0.0.0/24",
		"publish_vpc":    false,
		"weight":         100,
		"description":    "unit-test",
		"priority":       20,
	}
}

func vpnPbrRouteEntryUnitConfigValues(config map[string]interface{}) map[string]cty.Value {
	description := cty.NullVal(cty.String)
	if value, ok := config["description"].(string); ok {
		description = cty.StringVal(value)
	}
	priority := cty.NullVal(cty.Number)
	if value, ok := config["priority"].(int); ok {
		priority = cty.NumberIntVal(int64(value))
	}
	return map[string]cty.Value{
		"vpn_gateway_id": cty.StringVal(config["vpn_gateway_id"].(string)),
		"next_hop":       cty.StringVal(config["next_hop"].(string)),
		"route_source":   cty.StringVal(config["route_source"].(string)),
		"route_dest":     cty.StringVal(config["route_dest"].(string)),
		"publish_vpc":    cty.BoolVal(config["publish_vpc"].(bool)),
		"weight":         cty.NumberIntVal(int64(config["weight"].(int))),
		"description":    description,
		"priority":       priority,
		"status":         cty.NullVal(cty.String),
	}
}

func vpnPbrRouteEntryUnitRawConfigValue(r *schema.Resource, config map[string]interface{}) cty.Value {
	values := vpnPbrRouteEntryUnitConfigValues(config)
	values["timeouts"] = cty.NullVal(r.CoreConfigSchema().ImpliedType().AttributeType("timeouts"))
	return cty.ObjectVal(values)
}

func vpnPbrRouteEntryUnitRawStateValue(r *schema.Resource, config map[string]interface{}, id string) cty.Value {
	values := vpnPbrRouteEntryUnitConfigValues(config)
	values["id"] = cty.StringVal(id)
	values["status"] = cty.StringVal("normal")
	values["timeouts"] = cty.NullVal(r.CoreConfigSchema().ImpliedType().AttributeType("timeouts"))
	return cty.ObjectVal(values)
}

func vpnPbrRouteEntryUnitClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	staticCredential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("unit-test-access-key").SetAccessKeySecret("unit-test-secret-key"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := &sync.Map{}
	endpoints.Store("vpc", strings.TrimPrefix(server.URL, "http://"))
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")
	config := &connectivity.Config{
		AccessKey: "unit-test-access-key", SecretKey: "unit-test-secret-key", Credential: staticCredential,
		RegionId: "cn-hangzhou", AccountType: "ResourceOwner", Protocol: "HTTP", SkipRegionValidation: true,
		Endpoints: endpoints, SignVersion: &sync.Map{}, ClientReadTimeout: 5000, ClientConnectTimeout: 5000,
	}
	client, err := connectivity.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func vpnPbrRouteEntryUnitWriteJSON(t *testing.T, w http.ResponseWriter, status int, response map[string]interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		t.Errorf("write RPC response: %v", err)
	}
}
