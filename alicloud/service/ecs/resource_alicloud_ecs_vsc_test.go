// alicloud_ecs_vsc acceptance tests. The resource is Framework-native
// (terraform-plugin-framework), so the tests serve the real muxed provider —
// the SDK v2 half and the framework half behind one protocol v5 server — via
// alicloud/acctest.ProtoV5ProviderFactories, the shared factory every service
// subpackage uses. package ecs_test (external test package) is required because
// alicloud/acctest imports package alicloud for Provider(); an internal
// package ecs test file would form that cycle. The destroy check reads the
// configured AliyunClient from acctest.Provider.Meta() — the same provider
// instance ProtoV5ProviderFactories serves — mirroring the cas/hbr subpackage
// pattern, so CheckDestroy reaches the same connectivity client the production
// resource uses.
//
// The Vsc instance type ecs.ebmg9i-vsc.48xlarge is available in
// cn-wulanchabu, so the acceptance run targets that region. Cloud coverage
// (TestAcc*) is gated behind the standard ALICLOUD_* credential envs via
// acctest.PreCheck; the pure-function destroy-parse regressions (TestUnit*)
// run without the cloud and exercise the fail-closed shapes the prior loose
// response["Vscs"].([]interface{}) assertion let pass.
package ecs_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

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
// flags resource acceptance tests whose TestCase lacks a CheckDestroy — passes,
// and because the test genuinely should fail if a Vsc survives the destroy step.
func TestAccAliCloudEcsVsc(t *testing.T) {
	r := tfacctest.RandIntRange(1000, 9999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEcsVscDestroy,
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
	r := tfacctest.RandIntRange(1000, 9999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEcsVscDestroy,
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

// testAccCheckEcsVscDestroy asserts every alicloud_ecs_vsc left in state is
// gone after the test destroys its stack. The check reads the configured
// AliyunClient from acctest.Provider.Meta() — the same provider instance
// ProtoV5ProviderFactories serves, so its meta is the AliyunClient the
// production resource reaches through connectivity — then inlines the
// DescribeVscs RpcPost the resource uses: the Ecs 2014-05-26 contract returns
// an empty Vscs array for a Vsc that no longer exists (the NotFound shape
// parseEcsVscResponse treats as nil,nil), and DeleteVsc returns before the Vsc
// is fully gone, so the check polls DescribeVscs until the array empties or the
// timeout elapses. A NotFound-class error is also treated as gone, mirroring
// the tolerance DeleteVsc itself applies (resource_alicloud_ecs_vsc.go).
//
// Fail-closed: if state contains an alicloud_ecs_vsc but the provider has no
// configured client (nil meta), the check fails rather than silently passing —
// the fail-open shape the prior provider.Meta()==nil continue produced.
func testAccCheckEcsVscDestroy(s *terraform.State) error {
	hasVsc := false
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alicloud_ecs_vsc" {
			hasVsc = true
			break
		}
	}
	client, ok := acctest.Provider.Meta().(*connectivity.AliyunClient)
	if !ok || client == nil {
		if hasVsc {
			return fmt.Errorf("alicloud_ecs_vsc destroy check could not run: state contains alicloud_ecs_vsc resources but the provider has no configured client (nil meta)")
		}
		return nil
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_vsc" {
			continue
		}
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
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
					return resource.NonRetryableError(fmt.Errorf("DescribeVscs for ECS Vsc %s failed: %w", rs.Primary.ID, rerr))
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
				return resource.NonRetryableError(fmt.Errorf("DescribeVscs response for ECS Vsc %s is malformed: %w", rs.Primary.ID, perr))
			}
			if !gone {
				return resource.RetryableError(fmt.Errorf("ECS Vsc %s still exists", rs.Primary.ID))
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
// resource_alicloud_ecs_vsc.go) strict contract so the destroy check cannot
// diverge from Read: a legitimate empty Vscs array is the Ecs 2014-05-26
// NotFound shape and reports gone; a missing, null, or non-array $.Vscs root,
// a non-object entry, a non-empty result whose entries do not contain the
// queried VscId, or a matched entry missing Status are surfaced as errors
// rather than silently treated as gone, so a garbage or inconsistent response
// can never let a surviving Vsc pass the destroy gate. parseEcsVscResponse is
// unexported in package ecs and this check lives in package ecs_test, so the
// contract is mirrored here (minimal consistent check per the existing package
// boundary) rather than calling it; keep the two in sync. The DescribeVscs
// RpcPost error (authoritative NotFound vs a hard failure) is classified by
// ecsVscDestroyRpcError, not here.
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

// ecsVscDestroyRpcError classifies a DescribeVscs RpcPost error for the destroy
// check. An authoritative NotFound means the Vsc is gone; any other error is a
// hard failure the check surfaces rather than silently retrying away. Pure so
// the NotFound-vs-hard-failure line is exercisable without a live client.
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
