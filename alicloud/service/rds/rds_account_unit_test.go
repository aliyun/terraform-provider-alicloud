package rds

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Provider is the package-internal test provider built from the RDS
// ServicePackage registration. It is used only by the package-rds unit tests
// (service_test.go and the TestUnitRdsAccount* cases below) that stand up
// credential-free clients and exercise unexported CRUD functions directly. It
// cannot import the root alicloud package's Provider (that would form a cycle:
// rds -> acctest -> alicloud -> rds), so it is reconstructed here from the same
// SDKResources factories that registration wires up.
func Provider() *schema.Provider {
	sp := ServicePackage()
	p := &schema.Provider{ResourcesMap: map[string]*schema.Resource{}}
	for _, r := range sp.SDKResources {
		p.ResourcesMap[r.TypeName] = r.Factory()
	}
	return p
}

// rdsAccountTestClient stands up a credential-free AliyunClient whose RDS
// transport is pointed at an in-process httptest server, mirroring the
// ecs-snapshot unit-test construction. This avoids gomonkey binary patching
// (which the macOS arm64 kernel refuses) and runs under a plain `go test`.
func rdsAccountTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
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
	endpoints.Store("rds", endpoint)
	return client
}

func writeGone403(t *testing.T, w http.ResponseWriter) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"Code":    "OperationDenied.DBInstanceStatus",
		"Message": "the request was not permitted in the current state",
	}); err != nil {
		t.Error(err)
	}
}

func rdsAccountResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	p := Provider().ResourcesMap
	d, err := schema.InternalMap(p["alicloud_rds_account"].Schema).Data(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	d.SetId("rm-test:test-account")
	return d
}

func TestUnitRdsAccountReadParentGone(t *testing.T) {
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
			return
		case "DescribeAccounts":
			writeGone403(t, w)
			return
		}
		t.Errorf("unexpected RDS action in Read: %s", r.Form.Get("Action"))
		w.WriteHeader(http.StatusBadRequest)
	})
	d := rdsAccountResourceData(t)

	err := resourceAliCloudRdsAccountRead(d, client)
	assert.Nil(t, err, "Read must not hard-fail when the parent instance is gone")
	assert.Equal(t, "", d.Id(), "account state must be cleared when the parent instance is gone")
}

func TestUnitRdsAccountDeleteParentGone(t *testing.T) {
	var deleteAccountCalls int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
			return
		case "DeleteAccount":
			atomic.AddInt32(&deleteAccountCalls, 1)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
			return
		}
		t.Errorf("unexpected RDS action in Delete: %s", r.Form.Get("Action"))
		w.WriteHeader(http.StatusBadRequest)
	})
	d := rdsAccountResourceData(t)

	err := resourceAliCloudRdsAccountDelete(d, client)
	assert.Nil(t, err, "Delete must be idempotent when the parent instance is gone")
	assert.Equal(t, int32(0), atomic.LoadInt32(&deleteAccountCalls),
		"DeleteAccount must not be called when the parent instance is already gone")
}

func TestUnitRdsAccountDeleteDeleteAccountTerminal403(t *testing.T) {
	var describeDBInstanceCalls int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			n := atomic.AddInt32(&describeDBInstanceCalls, 1)
			if n == 1 {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Items": map[string]interface{}{
						"DBInstanceAttribute": []map[string]interface{}{
							{"DBInstanceStatus": "Running", "DBInstanceId": "rm-test"},
						},
					},
				})
				return
			}
			writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
			return
		case "DeleteAccount":
			writeGone403(t, w)
			return
		case "DescribeAccounts":
			writeGone403(t, w)
			return
		}
		t.Errorf("unexpected RDS action in Delete: %s", r.Form.Get("Action"))
		w.WriteHeader(http.StatusBadRequest)
	})
	d := rdsAccountResourceData(t)

	err := resourceAliCloudRdsAccountDelete(d, client)
	assert.Nil(t, err, "Delete must finish idempotently when DeleteAccount returns the terminal 403")
	assert.Equal(t, int32(2), atomic.LoadInt32(&describeDBInstanceCalls),
		"the follow-up DescribeDBInstance must run exactly once after the DeleteAccount 403 (pre-check + follow-up)")
}
