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

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudECSDiskAttachmentBasic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsDiskAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsDiskAttachment%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: AlicloudEcsDiskAttachmentBasicDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAliCloudECSDiskAttachmentMulti(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_attachment.default.1"
	ra := resourceAttrInit(resourceId, AlicloudEcsDiskAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsDiskAttachment%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: AlicloudEcsDiskAttachmentBasicDependenceMulti(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAliCloudECSDiskAttachmentBasic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsDiskAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsDiskAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsDiskAttachmentBasicDependence1)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id":              "${alicloud_ecs_disk.default.id}",
					"instance_id":          "${alicloud_instance.default.id}",
					"password":             "YouPassword123",
					"key_pair_name":        "${var.name}",
					"bootable":             false,
					"delete_with_instance": false,
					"force":                false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bootable", "delete_with_instance", "force", "key_pair_name", "password"},
			},
		},
	})
}

var AlicloudEcsDiskAttachmentMap = map[string]string{
	"disk_id":     CHECKSET,
	"instance_id": CHECKSET,
}

// newEcsDiskAttachmentTestData builds a fresh ResourceData for the
// alicloud_ecs_disk_attachment schema, marks it as a new resource and applies
// the given attributes. It is credential-free and exercises the real resource
// schema, so it can run in every CI stage (local-ci-check --quick changed-unit,
// ci-check-quick and full CI) without skipping on missing cloud credentials.
func newEcsDiskAttachmentTestData(attributes map[string]interface{}) *schema.ResourceData {
	d, _ := schema.InternalMap(resourceAlicloudEcsDiskAttachment().Schema).Data(nil, nil)
	d.MarkNewResource()
	// Set each schema attribute with a string-literal key. tfproviderlint R001
	// forbids a variable key argument to ResourceData.Set(), so the input map
	// is projected onto explicit per-attribute calls rather than a range loop.
	// Attributes absent from the map keep their schema defaults; "device" is
	// Computed-only and is refreshed by setEcsDiskAttachmentReadState instead.
	if v, ok := attributes["bootable"]; ok {
		if err := d.Set("bootable", v); err != nil {
			panic(fmt.Sprintf("failed to set bootable: %v", err))
		}
	}
	if v, ok := attributes["delete_with_instance"]; ok {
		if err := d.Set("delete_with_instance", v); err != nil {
			panic(fmt.Sprintf("failed to set delete_with_instance: %v", err))
		}
	}
	if v, ok := attributes["disk_id"]; ok {
		if err := d.Set("disk_id", v); err != nil {
			panic(fmt.Sprintf("failed to set disk_id: %v", err))
		}
	}
	if v, ok := attributes["force"]; ok {
		if err := d.Set("force", v); err != nil {
			panic(fmt.Sprintf("failed to set force: %v", err))
		}
	}
	if v, ok := attributes["instance_id"]; ok {
		if err := d.Set("instance_id", v); err != nil {
			panic(fmt.Sprintf("failed to set instance_id: %v", err))
		}
	}
	if v, ok := attributes["key_pair_name"]; ok {
		if err := d.Set("key_pair_name", v); err != nil {
			panic(fmt.Sprintf("failed to set key_pair_name: %v", err))
		}
	}
	if v, ok := attributes["password"]; ok {
		if err := d.Set("password", v); err != nil {
			panic(fmt.Sprintf("failed to set password: %v", err))
		}
	}
	return d
}

// TestUnitEcsDiskAttachmentAttachDiskRequest_Force is a credential-free
// request-level regression for the force argument of AttachDisk. It runs the
// real buildAttachDiskRequest code path and asserts that an explicit force=true
// sends Force=true, an explicit force=false sends Force=false, and an omitted
// force (default false) does not fabricate a value that contradicts the OpenAPI
// contract (AttachDisk Force: Boolean/Optional, default false). It replaces the
// previous static single-line wiring assertion and is what proves the force=true
// dynamic request construction is actually exercised.
func TestUnitEcsDiskAttachmentAttachDiskRequest_Force(t *testing.T) {
	base := map[string]interface{}{
		"disk_id":     "d-test-disk",
		"instance_id": "i-test-instance",
	}

	// explicit force=true must construct a real AttachDisk request carrying Force=true
	dTrue := newEcsDiskAttachmentTestData(mergeForce(base, true))
	reqTrue := buildAttachDiskRequest(dTrue)
	assert.Equal(t, "d-test-disk", reqTrue["DiskId"], "DiskId must be mapped from disk_id")
	assert.Equal(t, "i-test-instance", reqTrue["InstanceId"], "InstanceId must be mapped from instance_id")
	forceVal, forceOk := reqTrue["Force"]
	assert.True(t, forceOk, "explicit force=true must add Force to the AttachDisk request")
	assert.Equal(t, true, forceVal, "explicit force=true must send Force=true")

	// explicit force=false must still be a valid AttachDisk request; Force=false is
	// semantically equal to the API default and is sent because the field is set
	dFalse := newEcsDiskAttachmentTestData(mergeForce(base, false))
	reqFalse := buildAttachDiskRequest(dFalse)
	assert.Equal(t, "d-test-disk", reqFalse["DiskId"])
	assert.Equal(t, "i-test-instance", reqFalse["InstanceId"])
	if v, ok := reqFalse["Force"]; ok {
		assert.Equal(t, false, v, "explicit force=false may send Force=false (== API default)")
	}

	// omitted force must not fabricate a Force that contradicts the default-false
	// contract; whether the checker omits it or sends false, it must never be true
	dOmitted := newEcsDiskAttachmentTestData(base)
	reqOmitted := buildAttachDiskRequest(dOmitted)
	assert.Equal(t, "d-test-disk", reqOmitted["DiskId"])
	assert.Equal(t, "i-test-instance", reqOmitted["InstanceId"])
	if v, ok := reqOmitted["Force"]; ok {
		assert.Equal(t, false, v, "omitted force must never resolve to Force=true")
	}
}

// TestUnitEcsDiskAttachmentReadState_NoForce is a credential-free regression
// proving Read does not read force back into state. force is a create-only
// AttachDisk argument (only effective for ESSD cloud_regional_disk_auto disks)
// and is never returned by DescribeDisks, so a refresh must leave the configured
// value untouched. This keeps an import (ImportStatePassthrough -> Read) and a
// no-op refresh stable instead of fabricating a force drift.
func TestUnitEcsDiskAttachmentReadState_NoForce(t *testing.T) {
	d := newEcsDiskAttachmentTestData(map[string]interface{}{
		"disk_id":     "d-preconfigured",
		"instance_id": "i-preconfigured",
		"force":       true,
	})
	assert.Equal(t, true, d.Get("force"), "precondition: force is configured as true")

	disk := map[string]interface{}{
		"InstanceId": "i-from-describe",
		"DiskId":     "d-from-describe",
		"Device":     "/dev/xvdb",
	}
	setEcsDiskAttachmentReadState(d, disk)

	assert.Equal(t, "i-from-describe", d.Get("instance_id"), "Read must refresh instance_id")
	assert.Equal(t, "d-from-describe", d.Get("disk_id"), "Read must refresh disk_id")
	assert.Equal(t, "/dev/xvdb", d.Get("device"), "Read must refresh device")
	assert.Equal(t, true, d.Get("force"), "Read must NOT read force back: configured value is preserved (create-only)")
}

func mergeForce(base map[string]interface{}, force bool) map[string]interface{} {
	out := make(map[string]interface{}, len(base)+1)
	for k, v := range base {
		out[k] = v
	}
	out["force"] = force
	return out
}

// ecsDiskAttachmentTestClient stands up a credential-free AliyunClient whose ECS
// transport is pointed at an in-process httptest server, mirroring the
// ecs-snapshot / rds-account unit-test construction. The ECS endpoint is
// injected via Config.Endpoints ("ecs"); both the RPC path
// (client.RpcPost -> rpcRequest -> loadApiEndpoint) and the typed SDK path
// (WithEcsClient -> loadApiEndpoint) consult that map, so a single handler
// receives the real resource Create/Read/Delete SDK traffic (AttachDisk,
// DescribeDisks, DetachDisk) instead of any helper function. This uses httptest
// endpoint injection rather than binary patching, so it runs under a plain
// `go test` with no cloud credentials.
func ecsDiskAttachmentTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
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

// cannedAttachedDisk is the DescribeDisks response used across the Create/Read
// flows. Status is "In_use" so Create's BuildStateConf(WaitForState target
// "In_use") resolves on the first poll and completes deterministically without
// cloud resources. DeleteAutoSnapshot is stable so the post-AttachDisk re-read
// skips the ModifyDiskAttribute branch. DescribeDisks never carries a Force
// field (force is a create-only AttachDisk argument), which is what the
// Read-does-not-fabricate-force regression relies on.
func cannedAttachedDisk(diskID, instanceID string) map[string]interface{} {
	return map[string]interface{}{
		"RequestId": "00000000-0000-0000-0000-unit-test",
		"Disks": map[string]interface{}{"Disk": []map[string]interface{}{{
			"DiskId": diskID, "InstanceId": instanceID,
			"Status": "In_use", "DeleteAutoSnapshot": false, "Device": "/dev/xvdb",
		}}},
	}
}

func boolPtr(b bool) *bool { return &b }

// TestUnitEcsDiskAttachmentCreateSendsForce drives the real Create path
// (resourceAlicloudEcsDiskAttachmentCreate -> DescribeEcsDisk pre-read,
// AttachDisk via client.RpcPost, In_use state polling, post-read, then Read)
// against the httptest server and asserts the Force field actually carried by
// the AttachDisk HTTP request for force=true / force=false / omitted.
//
// buildAttachDiskRequest uses GetOkExists("force"): an explicit true or false is
// "set" and sends Force=true / Force=false. An omitted force (schema default
// false) must be key-absent or exactly Force=false -- both satisfy the OpenAPI
// contract (AttachDisk Force: Boolean/Optional, default false) -- and Force=true
// must never occur for an omitted field. This is the request-level proof that
// complements the helper-only buildAttachDiskRequest unit test.
func TestUnitEcsDiskAttachmentCreateSendsForce(t *testing.T) {
	for _, tc := range []struct {
		name      string
		setForce  *bool  // nil = omitted (schema default false, never explicitly set)
		wantForce string // expected AttachDisk Force form value; "" = omitted: key-absent (!Has) or exactly "false"
	}{
		{name: "force_true", setForce: boolPtr(true), wantForce: "true"},
		{name: "force_false", setForce: boolPtr(false), wantForce: "false"},
		{name: "force_omitted", setForce: nil, wantForce: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				attachForce      string
				attachHasForce   bool
				attachDiskID     string
				attachInstanceID string
				attachCount      int
				describeCount    int
			)
			client := ecsDiskAttachmentTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Form.Get("Action") {
				case "AttachDisk":
					attachCount++
					attachHasForce = r.Form.Has("Force")
					attachForce = r.Form.Get("Force")
					attachDiskID = r.Form.Get("DiskId")
					attachInstanceID = r.Form.Get("InstanceId")
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"RequestId": "00000000-attach-unit-test",
					})
				case "DescribeDisks":
					describeCount++
					_ = json.NewEncoder(w).Encode(cannedAttachedDisk("disk-test", "inst-test"))
				default:
					t.Errorf("unexpected ECS action in Create: %s", r.Form.Get("Action"))
					w.WriteHeader(http.StatusBadRequest)
					_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "UnexpectedAction"})
				}
			})

			r := resourceAlicloudEcsDiskAttachment()
			r.Timeouts = &schema.ResourceTimeout{Create: schema.DefaultTimeout(10 * time.Second)}
			d := r.Data(nil)
			d.MarkNewResource()
			if err := d.Set("disk_id", "disk-test"); err != nil {
				t.Fatalf("set disk_id: %v", err)
			}
			if err := d.Set("instance_id", "inst-test"); err != nil {
				t.Fatalf("set instance_id: %v", err)
			}
			if tc.setForce != nil {
				if err := d.Set("force", tc.setForce); err != nil {
					t.Fatalf("set force: %v", err)
				}
			}

			if err := resourceAlicloudEcsDiskAttachmentCreate(d, client); err != nil {
				t.Fatalf("Create failed for %s: %v", tc.name, err)
			}

			assert.Equal(t, "disk-test:inst-test", d.Id(), "id must be DiskId:InstanceId")

			assert.GreaterOrEqual(t, attachCount, 1, "Create must issue a real AttachDisk HTTP request")
			assert.Equal(t, "disk-test", attachDiskID, "AttachDisk must carry DiskId")
			assert.Equal(t, "inst-test", attachInstanceID, "AttachDisk must carry InstanceId")
			if tc.wantForce == "" {
				// Omitted force (schema default false) must be truly key-absent
				// (!attachHasForce) or exactly Force=false. r.Form.Get yields "" for
				// both a missing Force key and a present-but-empty one, so Has (not the
				// empty value) is the correct check on this HTTP layer; this rejects
				// Force=true, an empty-value Force key, and any illegal string. (The
				// Read-path InstanceState.Attributes check is a separate layer where an
				// absent map key looks up as "" -- not conflated here.)
				assert.True(t, !attachHasForce || attachForce == "false",
					"omitted force must be key-absent or exactly Force=false (not an empty-value key), got hasForce=%v force=%q", attachHasForce, attachForce)
			} else {
				assert.True(t, attachHasForce, "explicit %s must send a Force key", tc.name)
				assert.Equal(t, tc.wantForce, attachForce, "AttachDisk Force for %s", tc.name)
			}
			assert.GreaterOrEqual(t, describeCount, 1,
				"Create must issue real DescribeDisks (pre-read/state/read)")
		})
	}
}

// mustDynamicValue msgpack-encodes a cty.Value against the resource's implied
// type and wraps it as a tfprotov5.DynamicValue. The SDK's own gRPC server drives
// github.com/hashicorp/go-cty/cty/msgpack straight onto DynamicValue.MsgPack
// (helper/schema/grpc_provider.go); it never uses tfprotov5.NewDynamicValue,
// which is tftypes-based and would force a cty->tftypes bridge. Mirroring the
// SDK avoids the JSON-only EOF pitfall in the plugin protocol.
func mustDynamicValue(t *testing.T, ty cty.Type, val cty.Value) *tfprotov5.DynamicValue {
	t.Helper()
	b, err := msgpack.Marshal(val, ty)
	if err != nil {
		t.Fatalf("msgpack.Marshal failed: %v\n  type: %v\n  value: %v", err, ty, val)
	}
	return &tfprotov5.DynamicValue{MsgPack: b}
}

// ctyFromDynamicValue msgpack-decodes a tfprotov5.DynamicValue back to a cty.Value
// against the resource's implied type, mirroring the SDK server's unmarshal path
// (grpc_provider.go uses msgpack.Unmarshal(dv.MsgPack, impliedType)).
func ctyFromDynamicValue(t *testing.T, dv *tfprotov5.DynamicValue, ty cty.Type) cty.Value {
	t.Helper()
	if dv == nil {
		t.Fatal("expected a non-nil DynamicValue")
	}
	val, err := msgpack.Unmarshal(dv.MsgPack, ty)
	if err != nil {
		t.Fatalf("msgpack.Unmarshal failed: %v", err)
	}
	return val
}

// assertNoErrorDiagnostics fails the test if any diagnostic carries Error
// severity; Warning diagnostics are tolerated, matching the SDK's own test
// convention (helper/schema/grpc_provider_test.go).
func assertNoErrorDiagnostics(t *testing.T, diags []*tfprotov5.Diagnostic) {
	t.Helper()
	for _, d := range diags {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("unexpected error diagnostic: %s: %s", d.Summary, d.Detail)
		}
	}
}

// TestUnitEcsDiskAttachmentReadPreservesForce drives the REAL Terraform plugin
// protocol path for alicloud_ecs_disk_attachment -- schema.NewGRPCProviderServer
// -> ImportResourceState -> ReadResource -> PlanResourceChange -- against the
// same credential-free httptest server, instead of the low-level r.Diff helper.
//
// Why the real protocol path, not r.Diff: force/bootable are Optional+Default
// false ForceNew. A real ImportStatePassthrough produces an id-only state whose
// force is null. The low-level schemaMap.Diff (r.Diff) compares that null prior
// against the config-defaulted "false" and sees Old("") != New("false") ->
// RequiresNew, which is a harness distortion: the real plugin protocol path
// computes PlanResourceChange through grpc_provider.go's ValuesSDKEquivalent /
// forceNoChanges, which treats null and false as SDK-equivalent. So an
// omitted/null or explicit-false force against a null prior yields
// RequiresReplace=[] and a PlannedState that preserves the prior (a genuine
// no-op), while a force=true against the null prior is a real ForceNew change
// that must require replacement -- the true negative proving the harness does
// not always report no-op.
//
// force is a create-only AttachDisk argument (effective only for ESSD
// cloud_regional_disk_auto disks) and never appears in DescribeDisks output, so
// the Read request chain must not fabricate it. Two scenarios:
//   - refresh_force_configured: a resource created with force=true is driven
//     through ReadResource (DescribeDisks HTTP); force stays true (preserved,
//     not read from the server) and the matching force=true plan is a no-op.
//   - import_empty_state: ImportResourceState("disk:instance") runs the real
//     ImportStatePassthrough importer (id-only output, force null), ReadResource
//     refreshes without fabricating force, then PlanResourceChange is asserted a
//     no-op for force omitted/null and force=false, and a real replacement for
//     force=true.
func TestUnitEcsDiskAttachmentReadPreservesForce(t *testing.T) {
	const diskID, instanceID = "disk-read", "inst-read"
	const typeName = "alicloud_ecs_disk_attachment"
	ctx := context.Background()

	// The real resource (with its ForceNew/Default/Timeouts schema). CoreConfigSchema
	// returns an internal *configschema.Block; its exported ImpliedType()/CoerceValue()
	// are reachable via type inference without importing that internal package (this is
	// the same pattern alicloud production code uses, e.g. resource_alicloud_*.go).
	schemaBlock := resourceAlicloudEcsDiskAttachment().CoreConfigSchema()
	impliedType := schemaBlock.ImpliedType()

	// newProtocolServer stands up a per-subtest httptest ECS endpoint and returns a
	// schema.NewGRPCProviderServer whose Meta is the credential-free AliyunClient
	// (set via the public SetMeta setter), so ReadResource/PlanResourceChange hit the
	// real Read HTTP path with no cloud credentials and no Configure call.
	newProtocolServer := func(t *testing.T, describeCount *int) tfprotov5.ProviderServer {
		client := ecsDiskAttachmentTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "DescribeDisks":
				*describeCount++
				_ = json.NewEncoder(w).Encode(cannedAttachedDisk(diskID, instanceID))
			default:
				t.Errorf("unexpected ECS action in Read: %s", r.Form.Get("Action"))
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "UnexpectedAction"})
			}
		})
		p := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				typeName: resourceAlicloudEcsDiskAttachment(),
			},
		}
		p.SetMeta(client)
		return schema.NewGRPCProviderServer(p)
	}

	// planChange runs PlanResourceChange(prior, proposedNewState, config). ProposedNewState
	// is built per real Core rules: preserve the prior's computed id/device, take
	// disk_id/instance_id/force from the config, and let CoerceValue null-fill the
	// rest (bootable/delete_with_instance/key_pair_name/password and the implicit
	// timeouts block). Config mirrors what the user wrote in HCL.
	planChange := func(t *testing.T, server tfprotov5.ProviderServer, prior cty.Value, force cty.Value) *tfprotov5.PlanResourceChangeResponse {
		cfg, err := schemaBlock.CoerceValue(cty.ObjectVal(map[string]cty.Value{
			"disk_id":     cty.StringVal(diskID),
			"instance_id": cty.StringVal(instanceID),
			"force":       force,
		}))
		if err != nil {
			t.Fatalf("config CoerceValue failed: %v", err)
		}
		proposed, err := schemaBlock.CoerceValue(cty.ObjectVal(map[string]cty.Value{
			"id":          prior.GetAttr("id"),
			"device":      prior.GetAttr("device"),
			"disk_id":     cty.StringVal(diskID),
			"instance_id": cty.StringVal(instanceID),
			"force":       force,
		}))
		if err != nil {
			t.Fatalf("proposed CoerceValue failed: %v", err)
		}
		resp, err := server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
			TypeName:         typeName,
			PriorState:       mustDynamicValue(t, impliedType, prior),
			ProposedNewState: mustDynamicValue(t, impliedType, proposed),
			Config:           mustDynamicValue(t, impliedType, cfg),
		})
		if err != nil {
			t.Fatalf("PlanResourceChange failed: %v", err)
		}
		assertNoErrorDiagnostics(t, resp.Diagnostics)
		return resp
	}

	t.Run("refresh_force_configured", func(t *testing.T) {
		var describeCount int
		server := newProtocolServer(t, &describeCount)

		// A resource created with force=true: prior carries id, disk_id, instance_id,
		// force=true and the computed device; CoerceValue null-fills the rest.
		prior, err := schemaBlock.CoerceValue(cty.ObjectVal(map[string]cty.Value{
			"id":          cty.StringVal(diskID + ":" + instanceID),
			"disk_id":     cty.StringVal(diskID),
			"instance_id": cty.StringVal(instanceID),
			"force":       cty.True,
			"device":      cty.StringVal("/dev/xvdb"),
		}))
		if err != nil {
			t.Fatalf("prior CoerceValue failed: %v", err)
		}

		// Real Read HTTP via ReadResource (DescribeEcsDiskAttachment -> DescribeDisks).
		readResp, err := server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
			TypeName:     typeName,
			CurrentState: mustDynamicValue(t, impliedType, prior),
		})
		if err != nil {
			t.Fatalf("ReadResource failed: %v", err)
		}
		assertNoErrorDiagnostics(t, readResp.Diagnostics)
		refreshed := ctyFromDynamicValue(t, readResp.NewState, impliedType)
		assert.GreaterOrEqual(t, describeCount, 1, "Read must issue a real DescribeDisks HTTP request")
		// force is create-only and absent from DescribeDisks: Read must preserve the
		// configured true, not fabricate it from (or lose it to) the server response.
		assert.True(t, refreshed.GetAttr("force").RawEquals(cty.True),
			"refresh must preserve configured force=true (create-only, not read from DescribeDisks)")
		assert.True(t, refreshed.GetAttr("device").RawEquals(cty.StringVal("/dev/xvdb")), "Read must refresh device")
		assert.True(t, refreshed.GetAttr("disk_id").RawEquals(cty.StringVal(diskID)), "Read must refresh disk_id")
		assert.True(t, refreshed.GetAttr("instance_id").RawEquals(cty.StringVal(instanceID)), "Read must refresh instance_id")

		// Plan against the matching force=true config must be a no-op (no recreation).
		resp := planChange(t, server, refreshed, cty.True)
		assert.Empty(t, resp.RequiresReplace,
			"force=true preserved + matching config must be a no-op plan (RequiresReplace empty)")
		planned := ctyFromDynamicValue(t, resp.PlannedState, impliedType)
		assert.True(t, refreshed.RawEquals(planned),
			"no-op plan must preserve the refreshed prior in PlannedState (real no-op, not a field subset)")
	})

	t.Run("import_empty_state", func(t *testing.T) {
		var describeCount int
		server := newProtocolServer(t, &describeCount)

		// Real ImportResourceState runs the ImportStatePassthrough importer on the
		// import string "disk:instance" -> id-only state (force null, not fabricated).
		importResp, err := server.ImportResourceState(ctx, &tfprotov5.ImportResourceStateRequest{
			TypeName: typeName,
			ID:       diskID + ":" + instanceID,
		})
		if err != nil {
			t.Fatalf("ImportResourceState failed: %v", err)
		}
		assertNoErrorDiagnostics(t, importResp.Diagnostics)
		if len(importResp.ImportedResources) != 1 {
			t.Fatalf("ImportStatePassthrough must return one imported resource, got %d", len(importResp.ImportedResources))
		}
		importedRes := importResp.ImportedResources[0]
		imported := ctyFromDynamicValue(t, importedRes.State, impliedType)
		assert.True(t, imported.GetAttr("id").RawEquals(cty.StringVal(diskID+":"+instanceID)),
			"importer must set the id from the import string")
		// Pre-Read imported state carries only the id: force is null (not auto-filled to
		// the schema default false), and disk_id/instance_id are null (they come from
		// Read, not import).
		assert.True(t, imported.GetAttr("force").IsNull(),
			"pre-Read imported force must be null (id-only import, not fabricated to false/true)")
		assert.True(t, imported.GetAttr("disk_id").IsNull(),
			"pre-Read imported disk_id must be null (comes from Read)")
		assert.True(t, imported.GetAttr("instance_id").IsNull(),
			"pre-Read imported instance_id must be null (comes from Read)")

		// Real Read HTTP via ReadResource on the importer's ACTUAL returned state
		// (importedRes.State passed straight through, not re-marshalled or discarded).
		readResp, err := server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
			TypeName:     typeName,
			CurrentState: importedRes.State,
			Private:      importedRes.Private,
		})
		if err != nil {
			t.Fatalf("ReadResource failed: %v", err)
		}
		assertNoErrorDiagnostics(t, readResp.Diagnostics)
		refreshed := ctyFromDynamicValue(t, readResp.NewState, impliedType)
		assert.GreaterOrEqual(t, describeCount, 1, "Read must issue a real DescribeDisks HTTP request")
		// DescribeDisks carries no Force field; Read must not fabricate force=true.
		assert.True(t, !refreshed.GetAttr("force").RawEquals(cty.True),
			"import Read must not fabricate force=true from DescribeDisks")
		assert.True(t, refreshed.GetAttr("device").RawEquals(cty.StringVal("/dev/xvdb")), "Read must refresh device")

		// PlanResourceChange: omitted/null and explicit false vs the null prior are
		// SDK-equivalent (ValuesSDKEquivalent) -> no-op; force=true vs the null prior
		// is a real ForceNew change -> RequiresReplace (true negative).
		for _, tc := range []struct {
			name     string
			force    cty.Value
			wantNoop bool
		}{
			{name: "force_omitted", force: cty.NullVal(cty.Bool), wantNoop: true},
			{name: "force_false", force: cty.False, wantNoop: true},
			{name: "force_true_negative", force: cty.True, wantNoop: false},
		} {
			t.Run(tc.name, func(t *testing.T) {
				resp := planChange(t, server, refreshed, tc.force)
				if tc.wantNoop {
					assert.Empty(t, resp.RequiresReplace,
						"force %s vs null prior is SDK-equivalent (ValuesSDKEquivalent), no replacement", tc.name)
					planned := ctyFromDynamicValue(t, resp.PlannedState, impliedType)
					assert.True(t, refreshed.RawEquals(planned),
						"no-op plan must preserve the refreshed prior in PlannedState (real no-op, not a field subset)")
					// force stays null (not fabricated to true) through the plan.
					assert.True(t, !planned.GetAttr("force").RawEquals(cty.True),
						"no-op plan must not fabricate force=true")
				} else {
					assert.NotEmpty(t, resp.RequiresReplace,
						"true negative: force=true vs null prior is a real ForceNew change, must require replacement")
					seenForce := false
					for _, p := range resp.RequiresReplace {
						if strings.Contains(p.String(), "force") {
							seenForce = true
						}
					}
					assert.True(t, seenForce,
						"RequiresReplace must name force as the replaced ForceNew attribute, got %v", resp.RequiresReplace)
				}
			})
		}
	})
}

func AlicloudEcsDiskAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name = "tf-test"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_images" "default" {
   name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
 	most_recent = true
	owners = "system"
}
resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_name   = "${var.name}"
  host_name       = "tf-testAcc"
  image_id        = data.alicloud_images.default.images.0.id
  instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  security_groups = [alicloud_security_group.default.id]
  vswitch_id      = data.alicloud_vswitches.default.ids.0
}

data "alicloud_zones" "disk" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_ecs_disk" "default" {
	zone_id = "${data.alicloud_zones.disk.zones.0.id}"
	category = "cloud_efficiency"
	delete_auto_snapshot = "true"
	description = "Test For Terraform"
	disk_name = var.name
	enable_auto_snapshot = "true"
	encrypted = "true"
	size = "500"
  	tags = {
    	Created     = "TF"
    	Environment = "Acceptance-test"
  	}
}

resource "alicloud_ecs_disk_attachment" "default" {
  disk_id = alicloud_ecs_disk.default.id
  instance_id = alicloud_instance.default.id
}
`, name)
}
func AlicloudEcsDiskAttachmentBasicDependenceMulti(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

variable "number" {
	default = "2"
}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name = "tf-test"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}


data "alicloud_images" "default" {
    name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  	most_recent = true
	owners = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_name   = "${var.name}"
  host_name       = "tf-testAcc"
  image_id        = data.alicloud_images.default.images.0.id
  instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  security_groups = [alicloud_security_group.default.id]
  vswitch_id      = data.alicloud_vswitches.default.ids.0
}

data "alicloud_zones" "disk" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_ecs_disk" "default" {
	count = "${var.number}"
	zone_id = "${data.alicloud_zones.disk.zones.0.id}"
	category = "cloud_efficiency"
	delete_auto_snapshot = "true"
	description = "Test For Terraform"
	disk_name = var.name
	enable_auto_snapshot = "true"
	encrypted = "true"
	size = "500"
  	tags = {
    	Created     = "TF"
    	Environment = "Acceptance-test"
  	}
}

resource "alicloud_ecs_disk_attachment" "default" {
  count = "${var.number}"
  disk_id = "${element(alicloud_ecs_disk.default.*.id, count.index)}"
  instance_id = alicloud_instance.default.id
}
`, name)
}
func AlicloudEcsDiskAttachmentBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g6"
}

data "alicloud_images" "default" {
  name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent   = true
  owners        = "system"
  instance_type = data.alicloud_instance_types.default.instance_types.0.id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  instance_name     = var.name
  host_name         = "tf-testAcc"
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = alicloud_vswitch.default.id
}

resource "alicloud_ecs_disk" "default" {
  zone_id              = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  category             = "cloud_efficiency"
  delete_auto_snapshot = "true"
  description          = "Test For Terraform"
  disk_name            = var.name
  enable_auto_snapshot = "true"
  encrypted            = "true"
  size                 = "500"
  tags = {
    Created     = "TF"
    Environment = "Acceptance-test"
  }
}
`, name)
}
