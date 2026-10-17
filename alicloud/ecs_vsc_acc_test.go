package alicloud

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

	"github.com/PaesslerAG/jsonpath"
	credential "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccEcsVscProtoV5ProviderFactories serves the real muxed provider — the
// SDK v2 half and the framework half behind one protocol v5 server — so the
// Framework-native alicloud_ecs_vsc resource behaves under test exactly as it
// does in a real terraform run, whichever half implements it.
//
// This acceptance test lives in the flat alicloud package (rather than
// alicloud/service/ecs) because the remote acceptance runner discovers
// TestAcc functions only in the flat alicloud package. Package alicloud
// cannot import the alicloud/acctest helper (acctest imports alicloud for
// Provider(), which would form an import cycle), so the factory and precheck
// are constructed inline from alicloud.Provider() and
// provider.ProtoV5ProviderServerFactory. The SDK v2 helper/resource harness is
// used — not terraform-plugin-testing/helper/resource — because the flat
// alicloud test binary already links the SDK v2 harness across the rest of the
// suite, and linking both harnesses registers -sweep twice and panics with
// "flag redefined: sweep" before any test runs. The muxed server is identical
// to the one acctest.ProtoV5ProviderFactories serves, so resource behaviour is
// unchanged.
//
// The SDK v2 provider is captured into the callers' providers slice so
// CheckDestroy can reach its configured AliyunClient meta through the mux — the
// hbr replication vault test's pattern. getEcsVsc is unexported in package ecs,
// so the destroy check inlines the DescribeVscs RpcPost the resource uses.
func testAccEcsVscProtoV5ProviderFactories(providers *[]*schema.Provider) map[string]func() (tfprotov5.ProviderServer, error) {
	primary := Provider()
	*providers = append(*providers, primary)
	return map[string]func() (tfprotov5.ProviderServer, error){
		"alicloud": func() (tfprotov5.ProviderServer, error) {
			serverFactory, err := provider.ProtoV5ProviderServerFactory(context.Background(), primary)
			if err != nil {
				return nil, err
			}
			return serverFactory(), nil
		},
	}
}

// TestAccAliCloudEcsVsc exercises the Framework-native ECS Vsc resource end to
// end on the live account. The Vsc instance type ecs.ebmg9i-vsc.48xlarge is
// available in cn-wulanchabu, so the acceptance run must target that region.
//
// Step 1 creates the Vsc with two tags (CreateVsc). Steps 2-3 are in-place tag
// updates (add then remove) that drive the TagResources/UntagResources Update
// path. The final step imports and verifies.
//
// resource_group_id is RequiresReplace: CreateVsc accepts ResourceGroupId on
// create, but JoinResourceGroup rejects ResourceType="vsc" (404
// InvalidResourceType.NotFound; its enum is instance/disk/snapshot/image/
// securitygroup/ddh/ddhcluster/eni/keypair/launchtemplate/command/activation/
// managedinstance), so an in-place move is unsupported and a change forces
// replace. resource_group_id is therefore not moved in-place. This test does
// not set resource_group_id on create: a live run observed a server-side 500
// UnknownError from CreateVsc when it carried ResourceGroupId, even though the
// referenced resource group was already Active/OK and propagated (so the 500 is
// a server-side CreateVsc behaviour, not a provider bug). The field stays
// Optional+Computed+RequiresReplace and reads back from DescribeVscs; whether
// TagResources/UntagResources accept ResourceType="vsc" is decided here
// (steps 2-3 confirm they do), not assumed from the public enum, which is
// non-exhaustive for the global tagging APIs.
//
// CheckDestroy is set (not omitted) so tfproviderlint's AT001 rule — which
// flags resource acceptance tests whose TestCase lacks a CheckDestroy —
// passes, and because the test genuinely should fail if a Vsc survives the
// destroy step.
func TestAccAliCloudEcsVsc(t *testing.T) {
	r := acctest.RandIntRange(1000, 9999)
	var providers []*schema.Provider

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: testAccEcsVscProtoV5ProviderFactories(&providers),
		CheckDestroy:             testAccCheckEcsVscDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			// Create the Vsc with two tags (CreateVsc). resource_group_id is
			// not set on create (see field comment above).
			{
				Config: testAccEcsVscConfig(r, "Created = \"tfacc\"\n  Env     = \"dev\"\n", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ecs_vsc.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_ecs_vsc.default", "status"),
					resource.TestCheckResourceAttr("alicloud_ecs_vsc.default", "tags.%", "2"),
				),
			},
			// In-place update: keep Env, drop Created, add Team. UntagResources
			// (TagKey.1=Created) + TagResources (Tag.1.Key=Team) run.
			{
				Config: testAccEcsVscConfig(r, "Env  = \"dev\"\n  Team = \"qa\"\n", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ecs_vsc.default", "tags.%", "2"),
				),
			},
			// In-place update: reduce to a single tag, exercising UntagResources.
			{
				Config: testAccEcsVscConfig(r, "Env = \"dev\"\n", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ecs_vsc.default", "tags.%", "1"),
				),
			},
			{
				ResourceName:      "alicloud_ecs_vsc.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAliCloudEcsVscEmptyTags is the tags={} plan-state consistency
// regression the static review flagged: creating with an explicit empty tags
// map (tags={}) must not trip Framework's "inconsistent result after apply".
// Before the Create-path plan retention, ecsVscState normalized the known empty
// map to null when Describe carried no user tag, diverging from plan tags={}
// and failing apply. It also exercises the end-delete-to-zero path (tags ->
// Env -> {}): after UntagResources removes the last tag, the Update-path plan
// retention keeps the applied state an empty map (not null), so the next apply
// sees no spurious diff. The instance/region scaffolding, factory and
// CheckDestroy are shared with TestAccAliCloudEcsVsc; import coverage stays in
// that test: ImportStatePassthroughID seeds only id, so import's prior state
// carries no tags={} and ImportStateVerify on tags={} would report a
// null-vs-empty diff. A bare refresh of a state that already holds tags={} is
// now stable — applyReadResult carries the known empty map forward — but
// import cannot, so its coverage stays in the tags-bearing test above.
func TestAccAliCloudEcsVscEmptyTags(t *testing.T) {
	r := acctest.RandIntRange(1000, 9999)
	var providers []*schema.Provider

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: testAccEcsVscProtoV5ProviderFactories(&providers),
		CheckDestroy:             testAccCheckEcsVscDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			// Create with an explicit empty tags map (tags={}). Before the fix
			// this fails apply with "inconsistent result after apply": plan
			// tags={} (known empty map) vs applied state tags=null.
			{
				Config: testAccEcsVscConfig(r, "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ecs_vsc.default", "id"),
					resource.TestCheckResourceAttr("alicloud_ecs_vsc.default", "tags.%", "0"),
				),
			},
			// Add a tag (TagResources Update path).
			{
				Config: testAccEcsVscConfig(r, "Env = \"dev\"\n", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ecs_vsc.default", "tags.%", "1"),
				),
			},
			// End-delete to zero: tags={} again (UntagResources). The applied
			// state must stay an empty map, not null, so the next apply has no
			// spurious diff.
			{
				Config: testAccEcsVscConfig(r, "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ecs_vsc.default", "tags.%", "0"),
				),
			},
		},
	})
}

// testAccEcsVscConfig renders the Vsc with a caller-chosen tags block on top
// of a shared ECS instance dependency. The instance type/region are fixed to
// where the Vsc instance type is available. rgAttr is kept as a parameter for
// callers that need to place the Vsc in a resource group on create (CreateVsc
// takes ResourceGroupId); an empty rgAttr creates the Vsc in the default
// resource group. resource_group_id is RequiresReplace, so when set it only
// changes on replace, never in-place.
func testAccEcsVscConfig(r int, tagsBlock, rgAttr string) string {
	rgBlock := ""
	if rgAttr != "" {
		rgBlock = fmt.Sprintf(`
resource "alicloud_resource_manager_resource_group" "default" {
  resource_group_name = "tf-vsc-example-rg-%[1]d"
  display_name         = "tf-vsc-example-rg-%[1]d"
}
`, r)
	}
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-wulanchabu"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_instance_type      = "ecs.ebmg9i-vsc.48xlarge"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "tf-vsc-example-vpc-%[1]d"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "tf-vsc-example-vsw-%[1]d"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = "tf-vsc-example-sg-%[1]d"
  vpc_id = alicloud_vpc.default.id
}

// The bare-metal Vsc instance type only supports specific images; querying with
// instance_type set (DescribeImages InstanceType filter) returns only compatible
// system images, so images.0 is guaranteed purchasable with the type (RunInstances
// 403 InvalidParameter.NotMatch otherwise).
data "alicloud_images" "default" {
  most_recent   = true
  owners        = "system"
  instance_type = "ecs.ebmg9i-vsc.48xlarge"
}

resource "alicloud_instance" "default" {
  instance_type              = "ecs.ebmg9i-vsc.48xlarge"
  image_id                   = data.alicloud_images.default.images.0.id
  security_groups            = [alicloud_security_group.default.id]
  instance_name              = "tf-vsc-example-instance-%[1]d"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_charge_type       = "PostPaid"
  internet_max_bandwidth_out = 0
  system_disk_category       = "cloud_essd"
  system_disk_size           = 40
}
%[2]s
resource "alicloud_ecs_vsc" "default" {
  instance_id = alicloud_instance.default.id
  vsc_name    = "tf-vsc-example-%[1]d"
  description = "tf-vsc-example-%[1]d"
  tags = {
%[3]s  }
%[4]s}
`, r, rgBlock, tagsBlock, rgAttr)
}

// testAccCheckEcsVscDestroyWithProviders asserts every alicloud_ecs_vsc left
// in state is gone after the test destroys its stack. getEcsVsc is unexported
// in package ecs, so the check inlines the DescribeVscs RpcPost the resource
// uses: the Ecs 2014-05-26 contract returns an empty Vscs array for a Vsc that
// no longer exists (the NotFound shape parseEcsVscResponse treats as nil,nil),
// and DeleteVsc returns before the Vsc is fully gone, so the check polls
// DescribeVscs until the array empties or the timeout elapses. A NotFound-class
// error is also treated as gone, mirroring the tolerance DeleteVsc itself
// applies (resource_alicloud_ecs_vsc.go).
func testAccCheckEcsVscDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Determine whether the test left any alicloud_ecs_vsc in state. A
		// destroy check that cannot run because no provider has a configured
		// client (nil meta) must fail when Vsc resources remain: silently
		// passing would let a surviving Vsc escape the destroy gate, the
		// fail-open shape the prior provider.Meta()==nil continue produced.
		hasVsc := false
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "alicloud_ecs_vsc" {
				hasVsc = true
				break
			}
		}
		checked := false
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			checked = true
			if err := testAccCheckEcsVscDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		if hasVsc && !checked {
			return fmt.Errorf("alicloud_ecs_vsc destroy check could not run: state contains alicloud_ecs_vsc resources but no provider had a configured client (nil meta)")
		}
		return nil
	}
}

func testAccCheckEcsVscDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_vsc" {
			continue
		}
		if err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("Ecs", "2014-05-26", "DescribeVscs", nil, map[string]interface{}{
				"VscIds.1": rs.Primary.ID,
				"RegionId": client.RegionId,
			}, true)
			if err != nil {
				// An authoritative NotFound means the Vsc is gone (DeleteVsc
				// tolerates the same); anything else is a hard failure the test
				// must surface, not silently retry away. ecsVscDestroyRpcError
				// is extracted so the NotFound-vs-hard-failure line is exercised
				// without a live client.
				if rerr := ecsVscDestroyRpcError(err); rerr != nil {
					return retry.NonRetryableError(fmt.Errorf("DescribeVscs for ECS Vsc %s failed: %w", rs.Primary.ID, rerr))
				}
				return nil
			}
			// Parse the response strictly instead of the prior loose
			// response["Vscs"].([]interface{}) assertion: that assertion treated
			// a missing, null, or non-array $.Vscs root as "gone" (ok=false ->
			// return nil), letting a garbage response pass a surviving Vsc.
			// ecsVscDestroyGone mirrors the production parseEcsVscResponse
			// (package ecs) contract so the destroy check cannot diverge from
			// Read: an empty array is the Ecs 2014-05-26 NotFound shape (gone);
			// a missing/null/non-array root, a non-object entry, a non-empty
			// result without the queried VscId, or a matched entry missing
			// Status is an error (fail) rather than a silent pass.
			gone, perr := ecsVscDestroyGone(response, rs.Primary.ID)
			if perr != nil {
				return retry.NonRetryableError(fmt.Errorf("DescribeVscs response for ECS Vsc %s is malformed: %w", rs.Primary.ID, perr))
			}
			if !gone {
				return retry.RetryableError(fmt.Errorf("ECS Vsc %s still exists", rs.Primary.ID))
			}
			return nil
		}); err != nil {
			return fmt.Errorf("ECS Vsc %s destroy check failed: %w", rs.Primary.ID, err)
		}
	}
	return nil
}

// ecsVscDestroyGone decides from a DescribeVscs response whether the Vsc is
// gone. It mirrors the production parseEcsVscResponse (package ecs,
// resource_alicloud_ecs_vsc_unit.go) strict contract so the destroy check
// cannot diverge from Read: a legitimate empty Vscs array is the Ecs
// 2014-05-26 NotFound shape and reports gone; a missing, null, or non-array
// $.Vscs root, a non-object entry, a non-empty result whose entries do not
// contain the queried VscId, or a matched entry missing Status are surfaced
// as errors rather than silently treated as gone, so a garbage or
// inconsistent response can never let a surviving Vsc pass the destroy gate.
// parseEcsVscResponse is unexported in package ecs and this check lives in
// package alicloud, so the contract is mirrored here (minimal consistent
// check per the existing package boundary) rather than calling it; keep the
// two in sync. The DescribeVscs RpcPost error (authoritative NotFound vs a
// hard failure) is classified by ecsVscDestroyRpcError, not here.
func ecsVscDestroyGone(response map[string]interface{}, id string) (bool, error) {
	v, err := jsonpath.Get("$.Vscs", response)
	if err != nil {
		return false, fmt.Errorf("DescribeVscs response missing the $.Vscs root: %w", err)
	}
	if v == nil {
		return false, fmt.Errorf("DescribeVscs response has a null $.Vscs root")
	}
	arr, ok := v.([]interface{})
	if !ok {
		return false, fmt.Errorf("DescribeVscs $.Vscs is not an array: %T", v)
	}
	if len(arr) == 0 {
		// Legitimate empty array: the Vsc no longer exists.
		return true, nil
	}
	for _, raw := range arr {
		obj, ok := raw.(map[string]interface{})
		if !ok {
			// A non-object entry is malformed; surface it rather than
			// silently deciding the response is unusable or the Vsc is gone.
			return false, fmt.Errorf("DescribeVscs $.Vscs contains a non-object entry: %T", raw)
		}
		if vid, _ := obj["VscId"].(string); vid != id {
			// A sibling entry; skip rather than misattribute it.
			continue
		}
		if obj["Status"] == nil {
			// A matched entry must carry a Status (see parseEcsVscResponse);
			// a present-but-Status-less entry is inconsistent/malformed, so
			// fail rather than treating it as gone.
			return false, fmt.Errorf("DescribeVscs entry for VscId %q is missing Status: %v", id, obj)
		}
		// A matched, Status-bearing entry: the Vsc still exists.
		return false, nil
	}
	// Non-empty result but no entry whose VscId matches the queried id: the
	// response is inconsistent (DescribeVscs was queried with VscIds.1=id),
	// so fail rather than a silent NotFound.
	return false, fmt.Errorf("DescribeVscs result does not contain the queried VscId %q", id)
}

// ecsVscDestroyRpcError classifies a DescribeVscs RpcPost error for the
// destroy check. An authoritative NotFound means the Vsc is gone; any other
// error is a hard failure the check surfaces rather than silently retrying
// away. Pure so the NotFound-vs-hard-failure line is exercisable without a
// live client.
func ecsVscDestroyRpcError(err error) error {
	if errs.NotFoundError(err) {
		return nil
	}
	return err
}

// TestUnitEcsVscDestroyGone is the cloud-free regression for the destroy check's
// strict response parsing (mirrors parseEcsVscResponse). It covers the
// fail-open shapes the prior loose response["Vscs"].([]interface{}) assertion
// let pass: a missing, null, or non-array $.Vscs root and a non-object entry
// must fail, not report gone; a sibling entry must not be misread as this Vsc;
// a matched entry missing Status must fail. The legitimate empty array
// (NotFound shape) reports gone, and a matched, Status-bearing entry reports
// still present.
func TestUnitEcsVscDestroyGone(t *testing.T) {
	const id = "vsc-1"
	cases := []struct {
		name     string
		resp     map[string]interface{}
		wantGone bool
		wantErr  bool
	}{
		{
			name:     "empty array is gone (NotFound shape)",
			resp:     map[string]interface{}{"Vscs": []interface{}{}},
			wantGone: true,
		},
		{
			name: "matched entry with Status is still present",
			resp: map[string]interface{}{"Vscs": []interface{}{
				map[string]interface{}{"VscId": id, "Status": "In_use"},
			}},
			wantGone: false,
		},
		{
			name: "matched entry found after skipping a sibling",
			resp: map[string]interface{}{"Vscs": []interface{}{
				map[string]interface{}{"VscId": "vsc-other", "Status": "In_use"},
				map[string]interface{}{"VscId": id, "Status": "Available"},
			}},
			wantGone: false,
		},
		{
			name:    "missing $.Vscs root is an error, not gone",
			resp:    map[string]interface{}{"RequestId": "r-1"},
			wantErr: true,
		},
		{
			name:    "null $.Vscs root is an error, not gone",
			resp:    map[string]interface{}{"Vscs": nil},
			wantErr: true,
		},
		{
			name:    "non-array $.Vscs is an error, not gone",
			resp:    map[string]interface{}{"Vscs": "not-an-array"},
			wantErr: true,
		},
		{
			name:    "non-object entry is an error, not gone",
			resp:    map[string]interface{}{"Vscs": []interface{}{"not-an-object"}},
			wantErr: true,
		},
		{
			name: "non-empty result without the queried id is an error (query-id isolation)",
			resp: map[string]interface{}{"Vscs": []interface{}{
				map[string]interface{}{"VscId": "vsc-other", "Status": "In_use"},
			}},
			wantErr: true,
		},
		{
			name: "matched entry without Status is an error, not gone",
			resp: map[string]interface{}{"Vscs": []interface{}{
				map[string]interface{}{"VscId": id},
			}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gone, err := ecsVscDestroyGone(c.resp, id)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got gone=%v err=nil", gone)
				}
				if gone {
					t.Fatalf("expected gone=false on error, got true")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gone != c.wantGone {
				t.Fatalf("gone = %v, want %v", gone, c.wantGone)
			}
		})
	}
}

// TestUnitEcsVscDestroyRpcError is the cloud-free regression for the destroy
// check's RpcPost error classification: an authoritative NotFound reports the
// Vsc gone, and any other error (a permission failure, a 500, etc.) is
// surfaced as a hard failure rather than silently treated as gone.
func TestUnitEcsVscDestroyRpcError(t *testing.T) {
	if rerr := ecsVscDestroyRpcError(errs.GetNotFoundErrorFromString("NotFound")); rerr != nil {
		t.Fatalf("NotFound should be gone (nil), got %v", rerr)
	}
	// A non-NotFound ProviderError and a plain error must both surface, not
	// be swallowed as gone.
	for _, e := range []error{
		errs.NewProviderError("InvalidUser.NotSupported", "not allowed"),
		fmt.Errorf("Forbidden.KeyPairNotActive"),
	} {
		if rerr := ecsVscDestroyRpcError(e); rerr == nil {
			t.Fatalf("non-NotFound error should be surfaced, not treated as gone: %v", e)
		}
	}
}

// TestUnitEcsVscDestroyNilClient is the cloud-free regression for the destroy
// check's nil-client fail-closed rule: a provider with no configured client
// (nil meta) must fail when the state still holds an alicloud_ecs_vsc, rather
// than silently continuing past it (the fail-open shape the prior
// provider.Meta()==nil continue produced). With no Vsc in state a nil-meta
// provider is harmless and the check passes.
func TestUnitEcsVscDestroyNilClient(t *testing.T) {
	vscState := func(hasVsc bool) *terraform.State {
		res := map[string]*terraform.ResourceState{}
		if hasVsc {
			res["alicloud_ecs_vsc.default"] = &terraform.ResourceState{
				Type:    "alicloud_ecs_vsc",
				Primary: &terraform.InstanceState{ID: "vsc-1"},
			}
		}
		return &terraform.State{
			Modules: []*terraform.ModuleState{{
				Path:      []string{"root"},
				Resources: res,
			}},
		}
	}

	// Positive: Vsc in state, no configured client -> must fail.
	nilMetaProviders := []*schema.Provider{&schema.Provider{}}
	if err := testAccCheckEcsVscDestroyWithProviders(&nilMetaProviders)(vscState(true)); err == nil {
		t.Fatalf("expected destroy check to fail when a Vsc is in state but no provider has a configured client (nil meta)")
	}

	// Negative: no Vsc in state, a nil-meta provider is harmless -> pass.
	if err := testAccCheckEcsVscDestroyWithProviders(&nilMetaProviders)(vscState(false)); err != nil {
		t.Fatalf("expected destroy check to pass when no Vsc is in state, got %v", err)
	}
}

// ecsVscDestroyHttpState builds a minimal *terraform.State carrying one
// alicloud_ecs_vsc with the given id — the shape CheckDestroy inspects after a
// test destroys its stack. It is the httptest regression's inline state;
// TestUnitEcsVscDestroyNilClient keeps its own closure, and the two are not
// shared so this round's change stays a pure rename of the existing destroy
// regressions plus this one new HTTP-path test.
func ecsVscDestroyHttpState(id string) *terraform.State {
	return &terraform.State{
		Modules: []*terraform.ModuleState{{
			Path: []string{"root"},
			Resources: map[string]*terraform.ResourceState{
				"alicloud_ecs_vsc.default": {
					Type:    "alicloud_ecs_vsc",
					Primary: &terraform.InstanceState{ID: id},
				},
			},
		}},
	}
}

// ecsVscDestroyTestClient builds a real *connectivity.AliyunClient whose Ecs
// DescribeVscs RpcPost calls route to srv, mirroring the production
// provider.Configure client-build path with test credentials and the same
// bare-host:port endpoint injection the service/ecs Read/Create cloud-free
// tests use (the tea-rpc client builds the URL as {protocol}://{endpoint}/).
// No production code is changed; the client is only exercised through the
// destroy check's existing RpcPost call.
func ecsVscDestroyTestClient(t *testing.T, srv *httptest.Server) *connectivity.AliyunClient {
	t.Helper()
	endpoints := &sync.Map{}
	endpoints.Store("ecs", strings.TrimPrefix(srv.URL, "http://"))
	credConfig := new(credential.Config).
		SetType("access_key").
		SetAccessKeyId("vsc-test-ak").
		SetAccessKeySecret("vsc-test-sk")
	cred, err := credential.NewCredential(credConfig)
	if err != nil {
		t.Fatalf("NewCredential: %v", err)
	}
	cfg := &connectivity.Config{
		AccessKey:            "vsc-test-ak",
		SecretKey:            "vsc-test-sk",
		RegionId:             "cn-hangzhou",
		Protocol:             "HTTP",
		AccountType:          "Domestic",
		SkipRegionValidation: true,
		Endpoints:            endpoints,
		Credential:           cred,
		SignVersion:          new(sync.Map),
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("cfg.Client: %v", err)
	}
	return client
}

// TestUnitEcsVscDestroyHttpTest is the cloud-free regression that drives the
// destroy check through the real CheckDestroy -> RpcPost -> HTTP path against a
// local httptest.Server. It proves the destroy check queries DescribeVscs with
// VscIds.1 set to the target VscId (never a sibling, empty, or absent), and
// that a bad DescribeVscs response — a 200 missing the $.Vscs root — is never
// judged destroy-success: the check fails rather than silently passing a
// surviving Vsc, while a legitimate empty-array NotFound reports gone. The
// pure TestUnitEcsVscDestroyGone/RpcError regressions cover the decision logic;
// this test closes the loop through the actual client.RpcPost call and the
// strict parser. No production code is changed; no live API is exercised.
func TestUnitEcsVscDestroyHttpTest(t *testing.T) {
	const id = "vsc-http-1"
	cases := []struct {
		name    string
		body    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "200 missing the $.Vscs root is not destroy-success",
			body:    map[string]interface{}{"RequestId": "r-bad"},
			wantErr: true,
		},
		{
			name:    "legitimate empty-array NotFound reports gone",
			body:    map[string]interface{}{"Vscs": []interface{}{}},
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				gotVscID string
				mu       sync.Mutex
			)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.FormValue("Action") != "DescribeVscs" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if v := r.FormValue("VscIds.1"); v != "" {
					mu.Lock()
					gotVscID = v
					mu.Unlock()
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(c.body)
			}))
			t.Cleanup(srv.Close)

			client := ecsVscDestroyTestClient(t, srv)
			p := &schema.Provider{}
			p.SetMeta(client)

			err := testAccCheckEcsVscDestroyWithProvider(ecsVscDestroyHttpState(id), p)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected the destroy check to fail for a bad DescribeVscs response, got nil (destroy judged success)")
				}
			} else {
				if err != nil {
					t.Fatalf("expected the destroy check to report gone for a legitimate empty-array NotFound, got %v", err)
				}
			}
			mu.Lock()
			got := gotVscID
			mu.Unlock()
			if got != id {
				t.Fatalf("DescribeVscs was queried with VscIds.1=%q, want %q (the target VscId)", got, id)
			}
		})
	}
}
