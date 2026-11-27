// Unit coverage for the OTS service-layer V4 glue that sign_version.default
// = v4 adds on top of the client layer: the {"body": ...} response contract of
// the REST retry wrapper, the dara.Convert mapping of listVpcInfoByInstance,
// and the BindInstance2Vpc / UnbindInstance2Vpc parameters of the instance
// attachment. Everything runs against a local HTTP server that recomputes the
// V4 wire signature, so the requests are proven signed with the physical
// region rather than merely reaching the endpoint.
package ots_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	pop "github.com/alibabacloud-go/alibabacloud-gateway-pop/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	v4TestRegion     = "cn-hangzhou"
	v4TestInstance   = "test-instance"
	v4TestVpcName    = "vpc-name"
	v4TestVpcId      = "vpc-123"
	v4TestVSwitchId  = "vsw-123"
	v4TestAccessKey  = "test-access-key"
	v4TestSecretKey  = "test-secret-key"
	v4TestInstanceFq = "test-instance.cn-hangzhou.vpc.ots.aliyuncs.com"
)

// v4SigningFixture serves the endpoints the OTS V4 service glue calls: the
// REST data plane (getinstance / updateinstance), the OTS control plane
// (ListVpcInfoByInstance, BindInstance2Vpc, UnbindInstance2Vpc) and the VPC
// vswitch lookup the attachment create path resolves its VpcId through. The
// attachment is bound after BindInstance2Vpc and released after
// UnbindInstance2Vpc, so the resource read / wait paths observe real state.
type v4SigningFixture struct {
	t             *testing.T
	server        *httptest.Server
	client        *connectivity.AliyunClient
	mutex         sync.Mutex
	bound         bool
	bindQuery     url.Values
	unbindQuery   url.Values
	restPostBody  string
	bindLegacyRPC bool
	vswitchV4     bool
}

// newV4SigningFixture serves the OTS endpoints with sign_version default =
// v4. A non-empty otsOverride additionally stores the product-level "ots"
// override, the combination formal review F-02 reported (default = v4 with
// ots = v1 keeps VPC requests on the V4 path while the OTS requests go
// legacy).
func newV4SigningFixture(t *testing.T, otsOverride string) *v4SigningFixture {
	t.Helper()
	return newV4SigningFixtureWithHandler(t, otsOverride, nil)
}

func newV4SigningFixtureWithHandler(t *testing.T, otsOverride string, handler http.HandlerFunc) *v4SigningFixture {
	t.Helper()
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")
	for _, name := range []string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy"} {
		t.Setenv(name, "")
	}

	fixture := &v4SigningFixture{t: t}
	if handler == nil {
		handler = fixture.handle
	}
	fixture.server = httptest.NewServer(handler)
	t.Cleanup(fixture.server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(v4TestAccessKey).
		SetAccessKeySecret(v4TestSecretKey))
	if err != nil {
		t.Fatalf("failed to build the test credential: %v", err)
	}
	host := strings.TrimPrefix(fixture.server.URL, "http://")
	endpoints := &sync.Map{}
	endpoints.Store("ots", host)
	endpoints.Store("tablestore", host)
	endpoints.Store("vpc", host)
	signVersion := &sync.Map{}
	signVersion.Store("default", "v4")
	if otsOverride != "" {
		signVersion.Store("ots", otsOverride)
	}
	client, err := connectivity.NewClient(&connectivity.Config{
		AccessKey:            v4TestAccessKey,
		SecretKey:            v4TestSecretKey,
		Region:               connectivity.Region(v4TestRegion),
		RegionId:             v4TestRegion,
		AccountType:          "Domestic",
		Protocol:             "HTTP",
		ClientReadTimeout:    5000,
		ClientConnectTimeout: 5000,
		SkipRegionValidation: true,
		OtsEndpoint:          fixture.server.URL,
		Endpoints:            endpoints,
		SignVersion:          signVersion,
		Credential:           credential,
	})
	if err != nil {
		t.Fatalf("failed to build the test client: %v", err)
	}
	fixture.client = client
	return fixture
}

func (f *v4SigningFixture) handle(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/v2/openapi/getinstance":
		f.verifyWireV4(w, r, "tablestore")
		if got := r.URL.Query().Get("InstanceName"); got != v4TestInstance {
			f.fail(w, "getinstance InstanceName = %q, want %q", got, v4TestInstance)
			return
		}
		if got := r.URL.Query().Get("RegionId"); got != v4TestRegion {
			f.fail(w, "getinstance RegionId = %q, want %q", got, v4TestRegion)
			return
		}
		f.writeJSON(w, map[string]interface{}{
			"InstanceName":          v4TestInstance,
			"InstanceStatus":        "Running",
			"InstanceSpecification": "0-0-0-0",
		})
	case "/v2/openapi/updateinstance":
		body := f.verifyWireV4(w, r, "tablestore")
		if body == nil {
			return
		}
		f.mutex.Lock()
		f.restPostBody = string(body)
		f.mutex.Unlock()
		f.writeJSON(w, map[string]interface{}{"RequestId": "request-id"})
	default:
		f.handleControlPlane(w, r)
	}
}

func (f *v4SigningFixture) handleControlPlane(w http.ResponseWriter, r *http.Request) {
	// The common V4 RPC path routes the action through the x-acs-action
	// header (gateway signing), not the query parameter the legacy RPC
	// client used. The legacy POP SDK also sends an x-acs-action header on
	// RPC requests, so the signing style is told apart by the V4
	// authorization header only the gateway client sets.
	legacy := !isV4WireRequest(r)
	action := r.Header.Get("x-acs-action")
	if action == "" {
		action = r.URL.Query().Get("Action")
	}
	if action != "DescribeVSwitchAttributes" && r.URL.Query().Get("RegionId") != v4TestRegion {
		f.t.Errorf("%s: missing physical RegionId in the request query", action)
	}
	switch action {
	case "DescribeVSwitchAttributes":
		if f.verifyWire(w, r, "vpc") == nil {
			return
		}
		if r.Method != "POST" {
			f.fail(w, "DescribeVSwitchAttributes method = %s, want POST", r.Method)
			return
		}
		if got := r.URL.Query().Get("VSwitchId"); got != v4TestVSwitchId {
			f.fail(w, "DescribeVSwitchAttributes VSwitchId = %q, want %q", got, v4TestVSwitchId)
			return
		}
		if !legacy {
			f.mutex.Lock()
			f.vswitchV4 = true
			f.mutex.Unlock()
		}
		f.writeJSON(w, map[string]interface{}{
			"VpcId":     v4TestVpcId,
			"VSwitchId": v4TestVSwitchId,
			"RequestId": "request-id",
		})
	case "BindInstance2Vpc":
		if f.verifyWire(w, r, "ots") == nil {
			return
		}
		f.mutex.Lock()
		f.bindQuery = cloneQuery(r)
		f.bindLegacyRPC = legacy
		f.bound = true
		f.mutex.Unlock()
		f.writeJSON(w, map[string]interface{}{"RequestId": "request-id"})
	case "UnbindInstance2Vpc":
		if f.verifyWire(w, r, "ots") == nil {
			return
		}
		f.mutex.Lock()
		f.unbindQuery = cloneQuery(r)
		f.bound = false
		f.mutex.Unlock()
		f.writeJSON(w, map[string]interface{}{"RequestId": "request-id"})
	case "ListVpcInfoByInstance":
		if f.verifyWire(w, r, "ots") == nil {
			return
		}
		if got := r.URL.Query().Get("InstanceName"); got != v4TestInstance {
			f.fail(w, "ListVpcInfoByInstance InstanceName = %q, want %q", got, v4TestInstance)
			return
		}
		f.mutex.Lock()
		bound := f.bound
		f.mutex.Unlock()
		if bound {
			f.writeJSON(w, map[string]interface{}{
				"RequestId":    "request-id",
				"InstanceName": v4TestInstance,
				"TotalCount":   1,
				"PageNum":      1,
				"PageSize":     1,
				"VpcInfos": map[string]interface{}{
					"VpcInfo": []map[string]interface{}{{
						"Domain":          v4TestInstanceFq,
						"Endpoint":        v4TestInstanceFq,
						"RegionNo":        v4TestRegion,
						"InstanceName":    v4TestInstance,
						"InstanceVpcName": v4TestVpcName,
						"VpcId":           v4TestVpcId,
					}},
				},
			})
		} else {
			f.writeJSON(w, map[string]interface{}{
				"RequestId":  "request-id",
				"TotalCount": 0,
				"VpcInfos":   map[string]interface{}{"VpcInfo": []map[string]interface{}{}},
			})
		}
	default:
		f.fail(w, "unexpected request %s %s", r.Method, r.URL.String())
	}
}

// isV4WireRequest reports whether the request carries a gateway/V4
// authorization. The legacy POP SDK signs RPC requests through the query
// string but also sends an x-acs-action header, so the presence of that
// header alone cannot tell the two signing styles apart.
func isV4WireRequest(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Authorization"), "ACS4-HMAC-SHA256")
}

// verifyWire dispatches per signing style: a gateway/V4 request carries an
// ACS4-HMAC-SHA256 authorization and is verified by recomputing the V4
// signature, while a request from the legacy POP SDK keeps its RPC signature
// in the query and only needs the signature markers asserted.
func (f *v4SigningFixture) verifyWire(w http.ResponseWriter, r *http.Request, product string) []byte {
	if !isV4WireRequest(r) {
		query := r.URL.Query()
		if query.Get("Action") == "" || query.Get("SignatureMethod") != "HMAC-SHA1" || query.Get("Signature") == "" {
			f.fail(w, "legacy %s request is not RPC-signed: %s", product, r.URL.String())
			return nil
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			f.fail(w, "reading the legacy request body: %v", err)
			return nil
		}
		return body
	}
	return f.verifyWireV4(w, r, product)
}

// verifyWireV4 recomputes the ACS4-HMAC-SHA256 authorization from what the
// HTTP server actually received, including the payload digest and the Host
// header, so the assertion covers the physical signing region and the real
// request rather than an algorithm name. It returns the request body it
// consumed so handlers can assert the posted content.
func (f *v4SigningFixture) verifyWireV4(w http.ResponseWriter, r *http.Request, product string) []byte {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		f.fail(w, "reading the request body: %v", err)
		return nil
	}
	payload := fmt.Sprintf("%x", sha256.Sum256(body))
	date := r.Header.Get("X-Acs-Date")
	if len(date) < 10 || r.Header.Get("X-Acs-Content-Sha256") != payload {
		f.fail(w, "missing X-Acs-Date or wire payload digest mismatch")
		return nil
	}
	headers := map[string]*string{"host": &r.Host}
	for name := range r.Header {
		headers[strings.ToLower(name)] = tea.String(r.Header.Get(name))
	}
	query := make(map[string]*string)
	for name := range r.URL.Query() {
		query[name] = tea.String(r.URL.Query().Get(name))
	}
	gateway, err := pop.NewClient()
	if err != nil {
		f.fail(w, "building the pop gateway client: %v", err)
		return nil
	}
	algorithm := "ACS4-HMAC-SHA256"
	day := strings.ReplaceAll(date[:10], "-", "")
	key := gateway.GetSigningkey(&algorithm, tea.String(v4TestSecretKey), &product, tea.String(v4TestRegion), &day)
	expected, err := gateway.GetAuthorization(tea.String(r.URL.EscapedPath()), &r.Method, query, headers,
		&algorithm, &payload, tea.String(v4TestAccessKey), key, &product, tea.String(v4TestRegion), &day)
	if err != nil || tea.StringValue(expected) != r.Header.Get("Authorization") {
		f.fail(w, "wire signature does not match the recomputed %s V4 signature", product)
		return nil
	}
	return body
}

func (f *v4SigningFixture) writeJSON(w http.ResponseWriter, value map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func (f *v4SigningFixture) fail(w http.ResponseWriter, format string, args ...interface{}) {
	f.t.Errorf("v4 signing fixture: "+format, args...)
	http.Error(w, fmt.Sprintf(format, args...), http.StatusInternalServerError)
}

func cloneQuery(r *http.Request) url.Values {
	clone := make(url.Values, len(r.URL.Query()))
	for name, values := range r.URL.Query() {
		clone[name] = append([]string(nil), values...)
	}
	return clone
}

func (f *v4SigningFixture) setBound(bound bool) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.bound = bound
}

func (f *v4SigningFixture) isBound() bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.bound
}

func (f *v4SigningFixture) capturedBindQuery() url.Values {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.bindQuery
}

func (f *v4SigningFixture) capturedUnbindQuery() url.Values {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.unbindQuery
}

func (f *v4SigningFixture) capturedRestPostBody() string {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.restPostBody
}

func (f *v4SigningFixture) bindUsedLegacyRPC() bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.bindLegacyRPC
}

func (f *v4SigningFixture) sawV4VSwitchLookup() bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.vswitchV4
}

// TestUnitOtsV4RestResponseContract pins the legacy response contract of the
// REST retry wrapper on the V4 branch: OtsRestV4 already unwraps the Tea
// response, and invokeOtsRestApiWithRetry re-wraps it as {"body": ...} so the
// describe callers keep reading resp["body"] exactly as with the legacy ROA
// client (the contract fixed for AI review finding 260806450).
func TestUnitOtsV4RestResponseContract(t *testing.T) {
	fixture := newV4SigningFixture(t, "")
	service := ots.NewService(fixture.client)

	instance, err := service.DescribeOtsInstance(v4TestInstance)
	if err != nil {
		t.Fatalf("DescribeOtsInstance through the V4 REST path: %v", err)
	}
	if instance.InstanceName != v4TestInstance || instance.InstanceStatus != "Running" {
		t.Errorf("DescribeOtsInstance = %+v, want instance name and status from resp[\"body\"]", instance)
	}

	response, err := ots.OtsRestApiPostWithRetry(fixture.client, "tablestore", "2020-12-09", "/v2/openapi/updateinstance", map[string]interface{}{
		"InstanceName":   v4TestInstance,
		"NetworkTypeACL": []interface{}{"NORMAL"},
	})
	if err != nil {
		t.Fatalf("OtsRestApiPostWithRetry through the V4 REST path: %v", err)
	}
	body, ok := response["body"].(map[string]interface{})
	if !ok {
		t.Fatalf("resp[\"body\"] contract broken on the V4 branch: %v", response)
	}
	if body["RequestId"] != "request-id" {
		t.Errorf("resp[\"body\"][\"RequestId\"] = %v, want request-id", body["RequestId"])
	}
	posted := fixture.capturedRestPostBody()
	if !strings.Contains(posted, `"InstanceName":"`+v4TestInstance+`"`) || !strings.Contains(posted, `"NetworkTypeACL":["NORMAL"]`) {
		t.Errorf("updateinstance posted body = %q, want the request body passed through to the endpoint", posted)
	}
}

// TestUnitOtsV4ListVpcInfoByInstanceMapping covers the V4 branch of
// listVpcInfoByInstance: the raw OtsPopV4 map must convert into the legacy
// ListVpcInfoByInstanceResponse struct with the same fields the v1 SDK
// response carried, and an unbound instance must surface the shared
// not-found diagnostic.
func TestUnitOtsV4ListVpcInfoByInstanceMapping(t *testing.T) {
	fixture := newV4SigningFixture(t, "")
	service := ots.NewService(fixture.client)

	fixture.setBound(true)
	info, err := service.DescribeOtsInstanceAttachment(v4TestInstance)
	if err != nil {
		t.Fatalf("DescribeOtsInstanceAttachment through the V4 control plane: %v", err)
	}
	if info.InstanceVpcName != v4TestVpcName || info.VpcId != v4TestVpcId ||
		info.InstanceName != v4TestInstance || info.RegionNo != v4TestRegion {
		t.Errorf("dara.Convert mapping of ListVpcInfoByInstance = %+v, want the vpc binding fields", info)
	}

	fixture.setBound(false)
	_, err = service.DescribeOtsInstanceAttachment(v4TestInstance)
	if err == nil || !errs.NotFoundError(err) {
		t.Fatalf("expected the not-found diagnostic for an unbound instance, got %v", err)
	}
}

// TestUnitOtsV4InstanceAttachmentBindParams runs the public instance
// attachment create path under sign_version.default = v4 and asserts the
// BindInstance2Vpc request parameters on the wire, including the VpcId the
// vswitch lookup resolves through the V4 RPC path.
func TestUnitOtsV4InstanceAttachmentBindParams(t *testing.T) {
	fixture := newV4SigningFixture(t, "")
	resourceOtsAttachment := ots.ResourceOtsInstanceAttachment()
	diff, err := resourceOtsAttachment.Diff(context.Background(), nil, sdkterraform.NewResourceConfigRaw(map[string]interface{}{
		"instance_name": v4TestInstance,
		"vpc_name":      v4TestVpcName,
		"vswitch_id":    v4TestVSwitchId,
	}), fixture.client)
	if err != nil {
		t.Fatalf("planning the instance attachment: %v", err)
	}
	if diff == nil {
		t.Fatal("no diff produced for the attachment configuration")
	}
	state, diags := resourceOtsAttachment.Apply(context.Background(), nil, diff, fixture.client)
	if diags.HasError() {
		t.Fatalf("creating the instance attachment under V4: %v", diags)
	}
	if state.ID != v4TestInstance {
		t.Errorf("attachment id = %q, want the instance name", state.ID)
	}
	bindQuery := fixture.capturedBindQuery()
	if bindQuery == nil {
		t.Fatal("no BindInstance2Vpc request reached the endpoint")
	}
	for param, want := range map[string]string{
		"InstanceName":    v4TestInstance,
		"InstanceVpcName": v4TestVpcName,
		"VpcId":           v4TestVpcId,
		"VirtualSwitchId": v4TestVSwitchId,
	} {
		if got := bindQuery.Get(param); got != want {
			t.Errorf("BindInstance2Vpc %s = %q, want %q", param, got, want)
		}
	}
	if got := state.Attributes["vpc_id"]; got != v4TestVpcId {
		t.Errorf("state vpc_id = %v, want %q read back through the V4 control plane", got, v4TestVpcId)
	}
	if got := state.Attributes["vpc_name"]; got != v4TestVpcName {
		t.Errorf("state vpc_name = %v, want %q read back through the V4 control plane", got, v4TestVpcName)
	}
}

// TestUnitOtsV4InstanceAttachmentOtsV1Override pins the product-level
// override semantics formal review F-02 reported: with
// sign_version{default = v4, ots = v1} the attachment create must not be
// blocked by the provider's own V4 gate. The vswitch lookup is a VPC request
// and stays on the V4 RPC path, while the bind itself goes through the
// legacy POP client with its parameters on the query string.
func TestUnitOtsV4InstanceAttachmentOtsV1Override(t *testing.T) {
	fixture := newV4SigningFixture(t, "v1")
	resourceOtsAttachment := ots.ResourceOtsInstanceAttachment()
	diff, err := resourceOtsAttachment.Diff(context.Background(), nil, sdkterraform.NewResourceConfigRaw(map[string]interface{}{
		"instance_name": v4TestInstance,
		"vpc_name":      v4TestVpcName,
		"vswitch_id":    v4TestVSwitchId,
	}), fixture.client)
	if err != nil {
		t.Fatalf("planning the instance attachment: %v", err)
	}
	state, diags := resourceOtsAttachment.Apply(context.Background(), nil, diff, fixture.client)
	if diags.HasError() {
		t.Fatalf("creating the instance attachment with the ots = v1 override: %v", diags)
	}
	if state.ID != v4TestInstance {
		t.Errorf("attachment id = %q, want the instance name", state.ID)
	}
	if !fixture.sawV4VSwitchLookup() {
		t.Error("the vswitch lookup did not sign V4 as a VPC request under default = v4")
	}
	bindQuery := fixture.capturedBindQuery()
	if bindQuery == nil {
		t.Fatal("no BindInstance2Vpc request reached the endpoint")
	}
	if !fixture.bindUsedLegacyRPC() {
		t.Error("BindInstance2Vpc should go through the legacy POP client under the ots = v1 override")
	}
	for param, want := range map[string]string{
		"InstanceName":    v4TestInstance,
		"InstanceVpcName": v4TestVpcName,
		"VpcId":           v4TestVpcId,
		"VirtualSwitchId": v4TestVSwitchId,
	} {
		if got := bindQuery.Get(param); got != want {
			t.Errorf("BindInstance2Vpc %s = %q, want %q", param, got, want)
		}
	}
	if got := state.Attributes["vpc_id"]; got != v4TestVpcId {
		t.Errorf("state vpc_id = %v, want %q read back through the legacy control plane", got, v4TestVpcId)
	}
}

// TestUnitOtsV4InstanceAttachmentUnbindParams runs the public instance
// attachment delete path under sign_version.default = v4 and asserts the
// UnbindInstance2Vpc request parameters on the wire plus the wait loop
// observing the released binding.
func TestUnitOtsV4InstanceAttachmentUnbindParams(t *testing.T) {
	fixture := newV4SigningFixture(t, "")
	fixture.setBound(true)
	resourceOtsAttachment := ots.ResourceOtsInstanceAttachment()
	d := resourceOtsAttachment.Data(&sdkterraform.InstanceState{
		ID: v4TestInstance,
		Attributes: map[string]string{
			"instance_name": v4TestInstance,
			"vpc_name":      v4TestVpcName,
			"vswitch_id":    v4TestVSwitchId,
		},
	})
	if err := resourceOtsAttachment.Delete(d, fixture.client); err != nil {
		t.Fatalf("deleting the instance attachment under V4: %v", err)
	}
	unbindQuery := fixture.capturedUnbindQuery()
	if unbindQuery == nil {
		t.Fatal("no UnbindInstance2Vpc request reached the endpoint")
	}
	for param, want := range map[string]string{
		"InstanceName":    v4TestInstance,
		"InstanceVpcName": v4TestVpcName,
	} {
		if got := unbindQuery.Get(param); got != want {
			t.Errorf("UnbindInstance2Vpc %s = %q, want %q", param, got, want)
		}
	}
	if fixture.isBound() {
		t.Error("the attachment is still bound after the delete wait completed")
	}
}

func TestUnitOtsV4AttachmentAbsent(t *testing.T) {
	for _, tc := range []struct {
		name, code, message string
		status              int
		absent              bool
	}{
		{"instance code", "Instance.Notfound", "gone", http.StatusBadRequest, true},
		{"legacy message", "InvalidParameter", "instance is not found", http.StatusBadRequest, true},
		{"permission denied", "Forbidden", "instance is not found", http.StatusForbidden, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newV4SigningFixtureWithHandler(t, "", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_ = json.NewEncoder(w).Encode(map[string]string{"Code": tc.code, "Message": tc.message, "RequestId": "test"})
			})
			resource := ots.ResourceOtsInstanceAttachment()
			data := resource.Data(&sdkterraform.InstanceState{ID: v4TestInstance})
			err := resource.Read(data, fixture.client)
			if (err == nil) != tc.absent || (data.Id() == "") != tc.absent {
				t.Fatalf("read: id=%q, err=%v", data.Id(), err)
			}
			data.SetId(v4TestInstance)
			if err := resource.Delete(data, fixture.client); (err == nil) != tc.absent {
				t.Fatalf("delete: %v", err)
			}
			_, err = ots.NewService(fixture.client).DescribeOtsInstance(v4TestInstance)
			if err == nil || errs.NotFoundError(err) != tc.absent {
				t.Fatalf("instance read: %v", err)
			}
		})
	}
}
