package alicloud_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// This test exercises the framework-native resource end to end through the
// muxed (SDK v2 + framework) provider, using the SDK v2 helper/resource
// TestCase.ProtoV5ProviderFactories field (terraform-plugin-sdk/v2@v2.40.1
// helper/resource/plugin.go:263 explicitly supports supplying a muxed provider
// to that field; testcase_validate.go checks it). The SDK v2 harness is used on
// purpose rather than terraform-plugin-testing/helper/resource: the flat
// provider package transitively imports terraform-plugin-sdk/v2/helper/resource
// (its init registers a -sweep flag), so linking the testing variant alongside
// it would redefine -sweep and panic the test binary at startup. Using the SDK
// v2 harness links that single -sweep flag once and still drives the framework
// resource through the same acctest.ProtoV5ProviderFactories() mux, so the
// business resource stays Framework-native and registered; the harness choice
// does not migrate it back to SDK v2.

// TestAccAliCloudEcsSecurityGroupSnapshotPolicy_basic exercises the framework-native
// resource end to end through the muxed (SDK v2 + framework) provider: create,
// two partial updates that each change only one Optional+Computed field (so the
// Update guard skips the unchanged one), import, and a real CheckDestroy that
// confirms Describe no longer returns the policy.
//
// The partial updates are the cloud analogue of the F3 unknown-field regression:
// with only status changing, plan name equals the prior state name, so Update's
// IsUnknown+Equal guard omits name from Modify; with only name changing, status
// is omitted. A correct Modify therefore never clears an unchanged field.
func TestAccAliCloudEcsSecurityGroupSnapshotPolicy_basic(t *testing.T) {
	r := sdkacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc-sgsp-%d", r)
	renamed := fmt.Sprintf("tf-testacc-sgsp-%d-r", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckSecurityGroupSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig(name, "enable"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_security_group_snapshot_policy.default", "id"),
					// security_group_snapshot_policy_id is the API-generated id, not the
					// random name: assert it is set and equals the resource id rather than
					// fabricating an expected value from the name (no such contract).
					resource.TestCheckResourceAttrSet("alicloud_security_group_snapshot_policy.default", "security_group_snapshot_policy_id"),
					resource.TestCheckResourceAttrPair("alicloud_security_group_snapshot_policy.default", "id", "alicloud_security_group_snapshot_policy.default", "security_group_snapshot_policy_id"),
					resource.TestCheckResourceAttr("alicloud_security_group_snapshot_policy.default", "security_group_snapshot_policy_name", name),
					resource.TestCheckResourceAttr("alicloud_security_group_snapshot_policy.default", "status", "enable"),
					resource.TestCheckResourceAttrSet("alicloud_security_group_snapshot_policy.default", "region_id"),
				),
			},
			// Partial update 1: name unchanged, status enable→disable. Update must
			// send Status only; the equal-state name is skipped (F3 guard).
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig(name, "disable"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_security_group_snapshot_policy.default", "security_group_snapshot_policy_name", name),
					resource.TestCheckResourceAttr("alicloud_security_group_snapshot_policy.default", "status", "disable"),
				),
			},
			// Partial update 2: status unchanged, name→renamed. Update must send
			// SecurityGroupSnapshotPolicyName only; the equal-state status is skipped.
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig(renamed, "disable"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_security_group_snapshot_policy.default", "security_group_snapshot_policy_name", renamed),
					resource.TestCheckResourceAttr("alicloud_security_group_snapshot_policy.default", "status", "disable"),
				),
			},
			{
				ResourceName:      "alicloud_security_group_snapshot_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityGroupSnapshotPolicyConfig(name, status string) string {
	return fmt.Sprintf(`
resource "alicloud_security_group_snapshot_policy" "default" {
  security_group_snapshot_policy_name = "%[1]s"
  status                               = "%[2]s"
}
`, name, status)
}

// testAccCheckSecurityGroupSnapshotPolicyDestroy is a real, fail-closed
// CheckDestroy: after the test deletes the policy, DescribeSecurityGroupSnapshot
// Policies must no longer return it. A NotFound error or a well-formed empty list
// (or a list of other valid policies none matching) means it is gone; the id
// still present means the resource leaked and the test fails. A missing root, a
// wrong-type root, an inner object without its array, or a non-map element
// cannot be ruled out as the deleted policy, so each is reported as an error
// rather than swallowed as "absent" (a false PASS). It reaches the same
// Ecs/2016-03-14 Inner contract and the same strict list/id parsing the resource
// uses, through the client the muxed provider configured (acctest.Provider.Meta());
// it does not rely on the resource's unexported describe helper.
func testAccCheckSecurityGroupSnapshotPolicyDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_security_group_snapshot_policy" {
			continue
		}
		id := rs.Primary.ID
		request := map[string]interface{}{
			"RegionId":                       client.RegionId,
			"SecurityGroupSnapshotPolicyIds": fmt.Sprintf(`["%s"]`, id),
		}
		response, err := client.RpcPost("Ecs", "2016-03-14", "DescribeSecurityGroupSnapshotPolicies", nil, request, true)
		if err != nil {
			// A NotFound after delete is the expected terminal state.
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("CheckDestroy DescribeSecurityGroupSnapshotPolicies(%s): %w", id, err)
		}
		items, err := sgSnapshotPolicyDestroyItems(response)
		if err != nil {
			return fmt.Errorf("CheckDestroy %s: %w", id, err)
		}
		for _, raw := range items {
			obj, ok := raw.(map[string]interface{})
			if !ok {
				// A non-map element cannot be ruled out as the deleted policy:
				// confirm-destroyed fails closed rather than skipping it.
				return fmt.Errorf("CheckDestroy %s: response carried a non-object policy element %v; cannot confirm destroy", id, raw)
			}
			pid, parseable := destroyPolicyIdOf(obj)
			if !parseable {
				// An element whose id is missing or the wrong type cannot be
				// ruled out as the deleted policy: confirm-destroyed fails
				// closed rather than continuing to an empty pass (which would
				// report a false destroy). This is the caller-level guard the
				// strict destroyPolicyIdOf alone could not enforce.
				return fmt.Errorf("CheckDestroy %s: response carried a policy element with an unparseable id %v; cannot confirm destroy", id, obj)
			}
			if pid == id {
				return fmt.Errorf("security group snapshot policy %s still exists", id)
			}
		}
	}
	return nil
}

// sgSnapshotPolicyDestroyItems extracts the policy array from a Describe response
// in either POP list shape (top-level array, or inner object's array), mirroring
// the resource's strict securityGroupSnapshotPolicyItems. A missing root, a
// wrong-type root, or an inner object without its array is an error — not nil —
// so the caller cannot treat an untrustworthy response as a confirmed destroy.
func sgSnapshotPolicyDestroyItems(response map[string]interface{}) ([]interface{}, error) {
	if sets, ok := response["SecurityGroupSnapshotPolicySets"].([]interface{}); ok {
		return sets, nil
	}
	if inner, ok := response["SecurityGroupSnapshotPolicySets"].(map[string]interface{}); ok {
		if arr, ok := inner["SecurityGroupSnapshotPolicy"].([]interface{}); ok {
			return arr, nil
		}
		return nil, fmt.Errorf("SecurityGroupSnapshotPolicySets present but no SecurityGroupSnapshotPolicy array")
	}
	return nil, fmt.Errorf("SecurityGroupSnapshotPolicySets missing from response")
}

// destroyPolicyIdOf reads a policy's id with the same strict contract as the
// resource's securityGroupSnapshotPolicyIdOf: a non-empty string or a JSON
// number is parseable; a missing, nil or wrong-type id returns ("", false) so
// the caller treats the element as unconfirmed-destroyed rather than coercing
// it. RpcPost delivers JSON numbers as json.Number (tea-utils ReadAsJSON uses
// json.Decoder.UseNumber()), so json.Number is normalized via String() the same
// way policyIdOfOk normalizes it; the float64 branch covers response maps built
// without RpcPost. Without the json.Number branch a numeric id returned by
// DescribeSecurityGroupSnapshotPolicies would fall through to default and the
// CheckDestroy caller could not match it, misreporting a still-existing
// resource as destroyed.
func destroyPolicyIdOf(obj map[string]interface{}) (string, bool) {
	raw, ok := obj["SecurityGroupSnapshotPolicyId"]
	if !ok || raw == nil {
		return "", false
	}
	switch v := raw.(type) {
	case string:
		if v == "" {
			return "", false
		}
		return v, true
	case json.Number:
		s := v.String()
		if s == "" {
			return "", false
		}
		return s, true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	default:
		return "", false
	}
}

// sgSnapshotPolicyTestClient builds an AliyunClient whose Ecs RPC calls are
// redirected at an in-process httptest server with test credentials. No cloud
// account or network egress is involved: the server owns every response,
// including error and malformed envelopes. This mirrors the established
// ecsSnapshotTestClient pattern (alicloud/resource_alicloud_ecs_snapshot_test.go)
// so the fail-closed CheckDestroy caller regression drives the real provider
// meta path (acctest.Provider.SetMeta) rather than a copied helper.
func sgSnapshotPolicyTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("ecs", endpoint)
	return client
}

// CheckDestroy helper unit cases (no cloud, no PreCheck): they prove the
// fail-closed contract — a malformed response cannot masquerade as a confirmed
// destroy — so a future loose reimplementation cannot silently PASS.

func TestUnitSgSnapshotPolicyDestroyItems_MalformedRootErrors(t *testing.T) {
	cases := []struct {
		name     string
		response map[string]interface{}
		wantErr  bool
		wantItem int
	}{
		{"top-level array well-formed", map[string]interface{}{"SecurityGroupSnapshotPolicySets": []interface{}{map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-1"}}}, false, 1},
		{"top-level empty array is well-formed (gone)", map[string]interface{}{"SecurityGroupSnapshotPolicySets": []interface{}{}}, false, 0},
		{"inner array well-formed", map[string]interface{}{"SecurityGroupSnapshotPolicySets": map[string]interface{}{"SecurityGroupSnapshotPolicy": []interface{}{map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-2"}}}}, false, 1},
		{"root missing is malformed", map[string]interface{}{"Other": 1}, true, 0},
		{"root wrong scalar type is malformed", map[string]interface{}{"SecurityGroupSnapshotPolicySets": "unexpected"}, true, 0},
		{"inner object without array is malformed", map[string]interface{}{"SecurityGroupSnapshotPolicySets": map[string]interface{}{"Other": 1}}, true, 0},
		{"bare response is malformed", map[string]interface{}{}, true, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			items, err := sgSnapshotPolicyDestroyItems(c.response)
			if (err != nil) != c.wantErr {
				t.Fatalf("sgSnapshotPolicyDestroyItems(%v) err = %v; want error=%v", c.response, err, c.wantErr)
			}
			if err == nil && len(items) != c.wantItem {
				t.Fatalf("sgSnapshotPolicyDestroyItems(%v) items = %d; want %d", c.response, len(items), c.wantItem)
			}
		})
	}
}

// A non-map element in the array must make CheckDestroy fail (cannot confirm
// destroy), not silently pass as absent.
func TestUnitSgSnapshotPolicyDestroyItems_NonMapElementSurfaced(t *testing.T) {
	items, err := sgSnapshotPolicyDestroyItems(map[string]interface{}{
		"SecurityGroupSnapshotPolicySets": []interface{}{"not-an-object"},
	})
	if err != nil {
		t.Fatalf("well-formed root must not error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected the non-map element surfaced to the caller, got %v", items)
	}
	// The caller (testAccCheckSecurityGroupSnapshotPolicyDestroy) returns an error
	// for a non-map element; verify destroyPolicyIdOf rejects the wrong-type id
	// rather than coercing it.
	if _, parseable := destroyPolicyIdOf(map[string]interface{}{"SecurityGroupSnapshotPolicyId": map[string]interface{}{"x": "sp-1"}}); parseable {
		t.Fatalf("a wrong-type id must not be parseable")
	}
}

// destroyPolicyIdOf mirrors policyIdOfOk's strict type contract.
func TestUnitDestroyPolicyIdOf_StrictTypeContract(t *testing.T) {
	cases := []struct {
		obj    map[string]interface{}
		wantId string
		wantOk bool
	}{
		{map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-1"}, "sp-1", true},
		{map[string]interface{}{"SecurityGroupSnapshotPolicyId": float64(9)}, "9", true},
		{map[string]interface{}{"SecurityGroupSnapshotPolicyId": json.Number("9")}, "9", true},
		{map[string]interface{}{"Status": "enable"}, "", false},
		{map[string]interface{}{"SecurityGroupSnapshotPolicyId": nil}, "", false},
		{map[string]interface{}{"SecurityGroupSnapshotPolicyId": ""}, "", false},
		{map[string]interface{}{"SecurityGroupSnapshotPolicyId": map[string]interface{}{"x": "sp-1"}}, "", false},
	}
	for _, c := range cases {
		gotId, gotOk := destroyPolicyIdOf(c.obj)
		if gotId != c.wantId || gotOk != c.wantOk {
			t.Fatalf("destroyPolicyIdOf(%v) = (%q, %v); want (%q, %v)", c.obj, gotId, gotOk, c.wantId, c.wantOk)
		}
	}
}

// F4/#11/#13: the fail-closed contract must hold at the CheckDestroy CALLER,
// not only at the helper. testAccCheckSecurityGroupSnapshotPolicyDestroy drives
// the real Describe through acctest.Provider's injected client (the SDK exposes
// SetMeta, restored on Cleanup), so a response that cannot be confirmed
// destroyed — a map missing its id, a map with a wrong-type id, a non-map
// element, or a response missing the list root — fails instead of continuing to
// an empty pass. A well-formed empty list or a list of other legitimate ids is a
// confirmed destroy (nil); the target id present in the response fails. This is
// not a helper assertion: it calls the real caller through the configured
// provider meta, the same path CheckDestroy takes in an ACC run.
func TestUnitAccCheckSecurityGroupSnapshotPolicyDestroy_FailClosedCaller(t *testing.T) {
	// Preserve and restore the real provider meta so a fake client never leaks
	// past this test.
	original := acctest.Provider.Meta()
	t.Cleanup(func() { acctest.Provider.SetMeta(original) })

	state := &terraform.State{
		Modules: []*terraform.ModuleState{
			{
				Path: []string{"root"},
				Resources: map[string]*terraform.ResourceState{
					"alicloud_security_group_snapshot_policy.default": {
						Type:    "alicloud_security_group_snapshot_policy",
						Primary: &terraform.InstanceState{ID: "sp-1"},
					},
				},
			},
		},
	}

	cases := []struct {
		name     string
		response map[string]interface{}
		wantErr  bool
	}{
		{"map missing id", map[string]interface{}{"SecurityGroupSnapshotPolicySets": []interface{}{map[string]interface{}{"Status": "enable"}}}, true},
		{"map wrong-type id", map[string]interface{}{"SecurityGroupSnapshotPolicySets": []interface{}{map[string]interface{}{"SecurityGroupSnapshotPolicyId": map[string]interface{}{"x": "sp-1"}}}}, true},
		{"non-map element", map[string]interface{}{"SecurityGroupSnapshotPolicySets": []interface{}{"not-an-object"}}, true},
		{"missing root", map[string]interface{}{"RequestId": "r-1"}, true},
		{"empty list is confirmed gone", map[string]interface{}{"SecurityGroupSnapshotPolicySets": []interface{}{}}, false},
		{"other legitimate id is gone", map[string]interface{}{"SecurityGroupSnapshotPolicySets": []interface{}{map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-other"}}}, false},
		{"target id still exists", map[string]interface{}{"SecurityGroupSnapshotPolicySets": []interface{}{map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-1"}}}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := sgSnapshotPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(c.response)
			})
			acctest.Provider.SetMeta(client)
			err := testAccCheckSecurityGroupSnapshotPolicyDestroy(state)
			if (err != nil) != c.wantErr {
				t.Fatalf("testAccCheckSecurityGroupSnapshotPolicyDestroy(%s) err = %v; want error=%v", c.name, err, c.wantErr)
			}
		})
	}
}
