// Package polardb: state-machine tests for alicloud_polardb_backup. The file uses the internal
// package because it drives the real Create/Read/Update/Delete entry points with shortened poll
// intervals, which are package variables; the external polardb_test package cannot reach them.
package polardb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	fwpath "github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// These unit tests drive the resource through its public entry points (Create/Read/Update/Delete)
// against an isolated HTTP endpoint, so they exercise the real request serialization, the real
// SDK decoding and the real state machine. They never touch the cloud and never install a
// client interface seam: the provider client is the production client with its endpoint pointed
// at a local server.

const (
	backupUnitRegion    = "cn-beijing"
	backupUnitClusterID = "pc-unit-cluster"
	backupUnitBackupID  = "bk-unit-backup"
	backupUnitStartTime = "2026-10-20T00:00:00Z"
)

type backupUnitServer struct {
	t       *testing.T
	mu      sync.Mutex
	actions []string
	handler func(action string, params map[string]string) string
}

func (s *backupUnitServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	params := map[string]string{}
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	// The RPC client may send the API parameters in the POST body as JSON or as a form. The body
	// is read above, so it is decoded here explicitly instead of relying on ParseForm.
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if len(raw) > 0 {
		if strings.Contains(contentType, "json") {
			decoded := map[string]interface{}{}
			if json.Unmarshal(raw, &decoded) == nil {
				backupUnitFlattenParams("", decoded, params)
			}
		} else if values, err := url.ParseQuery(string(raw)); err == nil {
			for key := range values {
				params[key] = values.Get(key)
			}
		}
	}
	action := params["Action"]
	// A nested list parameter arrives as one JSON-encoded form value; expand it so the
	// assertions can address its leaves.
	for key, value := range params {
		trimmed := strings.TrimSpace(value)
		if !strings.HasPrefix(trimmed, "[") && !strings.HasPrefix(trimmed, "{") {
			continue
		}
		decoded := interface{}(nil)
		if json.Unmarshal([]byte(trimmed), &decoded) == nil {
			backupUnitFlattenParams(key, decoded, params)
		}
	}
	s.mu.Lock()
	s.actions = append(s.actions, action)
	s.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	body := "{}"
	if s.handler != nil {
		body = s.handler(action, params)
	}
	_, _ = w.Write([]byte(body))
}

// backupUnitFlattenParams flattens a decoded JSON body into dotted keys, so captured parameters
// are addressable regardless of whether the client nests entries or flattens them.
func backupUnitFlattenParams(prefix string, value interface{}, out map[string]string) {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, item := range typed {
			name := key
			if prefix != "" {
				name = prefix + "." + key
			}
			backupUnitFlattenParams(name, item, out)
		}
	case []interface{}:
		for index, item := range typed {
			backupUnitFlattenParams(fmt.Sprintf("%s.%d", prefix, index), item, out)
		}
	case nil:
		out[prefix] = ""
	default:
		out[prefix] = fmt.Sprint(typed)
	}
}

func (s *backupUnitServer) seen(action string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, seen := range s.actions {
		if seen == action {
			return true
		}
	}
	return false
}

func newBackupUnitClient(t *testing.T, handler func(action string, params map[string]string) string) (*connectivity.AliyunClient, *backupUnitServer) {
	t.Helper()
	server := &backupUnitServer{t: t, handler: handler}
	httpServer := httptest.NewServer(server)
	t.Cleanup(httpServer.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("building the test credential: %v", err)
	}
	endpoint := strings.TrimPrefix(httpServer.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             backupUnitRegion,
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("building the test client: %v", err)
	}
	endpoints.Store("polardb", endpoint)
	return client, server
}

func newBackupUnitResource(t *testing.T, client *connectivity.AliyunClient) (fwresource.Resource, fwschema.Schema) {
	t.Helper()
	r := NewBackupResource()
	configured, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatal("the backup resource does not accept the provider client")
	}
	var configureResp fwresource.ConfigureResponse
	configured.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", configureResp.Diagnostics)
	}
	var schemaResp fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building the schema: %v", schemaResp.Diagnostics)
	}
	return r, schemaResp.Schema
}

func backupUnitObjectType(t *testing.T, sch fwschema.Schema) tftypes.Object {
	t.Helper()
	objectType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("the resource schema does not describe an object")
	}
	return objectType
}

// backupUnitObject builds a fully typed object value: every attribute not listed is null.
func backupUnitObject(t *testing.T, objectType tftypes.Object, attributes map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	full := map[string]tftypes.Value{}
	for name, attributeType := range objectType.AttributeTypes {
		if value, ok := attributes[name]; ok {
			full[name] = value
			continue
		}
		full[name] = tftypes.NewValue(attributeType, nil)
	}
	return tftypes.NewValue(objectType, full)
}

func backupUnitString(value string) tftypes.Value {
	return tftypes.NewValue(tftypes.String, value)
}

func backupUnitBool(value bool) tftypes.Value {
	return tftypes.NewValue(tftypes.Bool, value)
}

// backupUnitPolicyList builds the advanced_data_policies block value from per-entry attributes.
func backupUnitPolicyList(t *testing.T, objectType tftypes.Object, entries []map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	listType, ok := objectType.AttributeTypes["advanced_data_policies"].(tftypes.List)
	if !ok {
		t.Fatal("advanced_data_policies is not a list block")
	}
	elementType, ok := listType.ElementType.(tftypes.Object)
	if !ok {
		t.Fatal("advanced_data_policies entries are not objects")
	}
	values := make([]tftypes.Value, 0, len(entries))
	for _, entry := range entries {
		full := map[string]tftypes.Value{}
		for name, attributeType := range elementType.AttributeTypes {
			if value, ok := entry[name]; ok {
				full[name] = value
				continue
			}
			full[name] = tftypes.NewValue(attributeType, nil)
		}
		values = append(values, tftypes.NewValue(elementType, full))
	}
	return tftypes.NewValue(listType, values)
}

func backupUnitNotFound(code string) string {
	return fmt.Sprintf(`{"Code":%q,"Message":"unit test not found"}`, code)
}

func backupUnitBackupsJSON(entries ...string) string {
	if len(entries) == 0 {
		return `{"Items":{"Backup":[]},"TotalRecordCount":0}`
	}
	return fmt.Sprintf(`{"Items":{"Backup":[%s]},"TotalRecordCount":%d}`, strings.Join(entries, ","), len(entries))
}

// backupUnitBackupEntry builds a DescribeBackups entry with only the fields that action
// documents; StoreStatus belongs to DescribeDetachedBackups and is added by that fixture.
func backupUnitBackupEntry(backupID, clusterID, startTime string) string {
	return fmt.Sprintf(`{"BackupId":%q,"DBClusterId":%q,"BackupStatus":"Success","BackupStartTime":%q,"BackupMode":"Manual"}`,
		backupID, clusterID, startTime)
}

func backupUnitStateAttribute(t *testing.T, state tfsdk.State, name string) (string, bool) {
	t.Helper()
	var value types.String
	if diagnostics := state.GetAttribute(context.Background(), fwpath.Root(name), &value); diagnostics.HasError() {
		return "", false
	}
	if value.IsNull() || value.IsUnknown() {
		return "", false
	}
	return value.ValueString(), true
}

func backupUnitDiagnosticsText(diagnostics interface{ HasError() bool }) string {
	return fmt.Sprint(diagnostics)
}

func withBackupUnitTimings(t *testing.T, poll time.Duration, grace time.Duration, resolve time.Duration) {
	t.Helper()
	previousPoll, previousGrace, previousResolve := backupPollInterval, backupGoneGrace, backupCreateResolveTimeout
	backupPollInterval = poll
	backupGoneGrace = grace
	backupCreateResolveTimeout = resolve
	t.Cleanup(func() {
		backupPollInterval = previousPoll
		backupGoneGrace = previousGrace
		backupCreateResolveTimeout = previousResolve
	})
}

// TestUnitAliCloudPolarDBBackupCreateRefusesForeignBackup covers PRR-01: the created backup set
// is bound only through the task id returned by CreateBackup, and a set that another actor
// created at a different time is never accepted as this task's backup set.
func TestUnitAliCloudPolarDBBackupCreateRefusesForeignBackup(t *testing.T) {
	withBackupUnitTimings(t, 5*time.Millisecond, time.Second, 120*time.Millisecond)
	foreign := backupUnitBackupEntry("bk-someone-else", backupUnitClusterID, "2026-10-19T23:58:00Z")
	client, server := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "CreateBackup":
			return fmt.Sprintf(`{"BackupJobId":%q}`, backupUnitBackupID)
		case "DescribeBackupTasks":
			return fmt.Sprintf(`{"Items":{"BackupJob":[{"BackupJobId":%q,"BackupProgressStatus":"Finished","StartTime":%q}]}}`, backupUnitBackupID, backupUnitStartTime)
		case "DescribeBackups":
			// Our task's backup set is never listed; another actor's set is.
			return backupUnitBackupsJSON(foreign)
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)

	plan := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"db_cluster_id": backupUnitString(backupUnitClusterID),
		"comment":       backupUnitString("unit-foreign"),
	})
	var resp fwresource.CreateResponse
	resp.State = tfsdk.State{Schema: sch}
	r.Create(context.Background(), fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: plan}}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("Create accepted a create whose own backup set never appeared: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if !server.seen("DescribeBackups") {
		t.Fatal("Create did not verify the created backup set against DescribeBackups")
	}
	if detail := backupUnitDiagnosticsText(resp.Diagnostics); !strings.Contains(detail, backupUnitBackupID) {
		t.Errorf("Create error does not name the backup task id it refused to bind: %s", detail)
	}
}

// TestUnitAliCloudPolarDBBackupCreateRefusesAmbiguousBackupSet covers the other side of the
// task-start-time binding: when two manual backup sets started at the task start time, the
// create must refuse instead of picking one.
func TestUnitAliCloudPolarDBBackupCreateRefusesAmbiguousBackupSet(t *testing.T) {
	withBackupUnitTimings(t, 5*time.Millisecond, time.Second, 120*time.Millisecond)
	first := backupUnitBackupEntry("bk-first", backupUnitClusterID, backupUnitStartTime)
	second := backupUnitBackupEntry("bk-second", backupUnitClusterID, backupUnitStartTime)
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "CreateBackup":
			return fmt.Sprintf(`{"BackupJobId":%q}`, backupUnitBackupID)
		case "DescribeBackupTasks":
			return fmt.Sprintf(`{"Items":{"BackupJob":[{"BackupJobId":%q,"BackupProgressStatus":"Finished","StartTime":%q}]}}`, backupUnitBackupID, backupUnitStartTime)
		case "DescribeBackups":
			return backupUnitBackupsJSON(first, second)
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)

	plan := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"db_cluster_id": backupUnitString(backupUnitClusterID),
		"comment":       backupUnitString("unit-ambiguous"),
	})
	var resp fwresource.CreateResponse
	resp.State = tfsdk.State{Schema: sch}
	r.Create(context.Background(), fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: plan}}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("Create bound one of two backup sets that started at the task start time: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if detail := backupUnitDiagnosticsText(resp.Diagnostics); !strings.Contains(detail, "backup sets starting at that second") {
		t.Errorf("Create error does not explain the ambiguous backup sets: %s", detail)
	}
}

// TestUnitAliCloudPolarDBBackupCreateBindsTaskBackupSet covers the positive branch: the task id
// from CreateBackup differs from the backup set id, and the set is bound because it is the only
// manual backup starting at the task start time.
func TestUnitAliCloudPolarDBBackupCreateBindsTaskBackupSet(t *testing.T) {
	withBackupUnitTimings(t, 5*time.Millisecond, time.Second, 2*time.Second)
	created := backupUnitBackupEntry("3218271464", backupUnitClusterID, backupUnitStartTime)
	automated := backupUnitBackupEntry("3217277072", backupUnitClusterID, backupUnitStartTime)
	automated = strings.Replace(automated, `"BackupMode":"Manual"`, `"BackupMode":"Automated"`, 1)
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "CreateBackup":
			return fmt.Sprintf(`{"BackupJobId":%q}`, backupUnitBackupID)
		case "DescribeBackupTasks":
			return fmt.Sprintf(`{"Items":{"BackupJob":[{"BackupJobId":%q,"BackupProgressStatus":"Finished","StartTime":%q}]}}`, backupUnitBackupID, backupUnitStartTime)
		case "DescribeBackups":
			return backupUnitBackupsJSON(automated, created)
		case "DescribeBackupPolicy":
			return `{"DBClusterId":"` + backupUnitClusterID + `","BackupPolicyLevel":"Normal"}`
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)

	plan := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"db_cluster_id": backupUnitString(backupUnitClusterID),
		"comment":       backupUnitString("unit-bind"),
	})
	var resp fwresource.CreateResponse
	resp.State = tfsdk.State{Schema: sch}
	r.Create(context.Background(), fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: plan}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed to bind the set produced by its task: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if id, ok := backupUnitStateAttribute(t, resp.State, "id"); !ok || id != backupUnitClusterID+":3218271464" {
		t.Fatalf("Create did not record the resolved backup set id: %q present=%v", id, ok)
	}
	if jobId, ok := backupUnitStateAttribute(t, resp.State, "backup_job_id"); !ok || jobId != backupUnitBackupID {
		t.Fatalf("Create did not record the authoritative task id: %q present=%v", jobId, ok)
	}
}

// TestUnitAliCloudPolarDBBackupReadDetachedRetained covers G1/G2: after the cluster is released,
// Read resolves the retained backup set from DescribeDetachedBackups, keeps the resource and
// reports the StoreStatus the released-cluster list is the only source for.
func TestUnitAliCloudPolarDBBackupReadDetachedRetained(t *testing.T) {
	withBackupUnitTimings(t, 5*time.Millisecond, time.Second, time.Second)
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DescribeBackupPolicy":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeBackups":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeDetachedBackups":
			// Only the fields DescribeDetachedBackups documents: the comment is not one of them,
			// so it must be preserved from state rather than taken from the response.
			return backupUnitBackupsJSON(`{"BackupId":"` + backupUnitBackupID + `","DBClusterId":"` + backupUnitClusterID + `","BackupStatus":"Success","StoreStatus":"Success","BackupStartTime":"` + backupUnitStartTime + `","IsAvail":"1"}`)
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":     backupUnitString(backupUnitClusterID),
		"backup_id":         backupUnitString(backupUnitBackupID),
		"backup_start_time": backupUnitString(backupUnitStartTime),
		"backup_retention_policy_on_cluster_deletion": backupUnitString("ALL"),
		"preferred_backup_time":                       backupUnitString("10:00Z-11:00Z"),
		"comment":                                     backupUnitString("unit-retained"),
	})

	var resp fwresource.ReadResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed on the retained detached branch: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if id, ok := backupUnitStateAttribute(t, resp.State, "id"); !ok || id != backupUnitClusterID+":"+backupUnitBackupID {
		t.Fatalf("Read dropped the retained backup set from state: id=%q present=%v", id, ok)
	}
	if storeStatus, ok := backupUnitStateAttribute(t, resp.State, "store_status"); !ok || storeStatus != "Success" {
		t.Fatalf("Read did not report store_status from DescribeDetachedBackups: %q present=%v", storeStatus, ok)
	}
	if retention, ok := backupUnitStateAttribute(t, resp.State, "backup_retention_policy_on_cluster_deletion"); !ok || retention != "ALL" {
		t.Fatalf("Read lost the cluster policy fields of a released cluster: %q present=%v", retention, ok)
	}
	if comment, ok := backupUnitStateAttribute(t, resp.State, "comment"); !ok || comment != "unit-retained" {
		t.Fatalf("Read did not map the detached backup-set fields: comment=%q present=%v", comment, ok)
	}
}

// TestUnitAliCloudPolarDBBackupReadKeepsRetainedAbsence covers N2: retention ALL keeps the backup
// set by contract, so an empty released-cluster listing must keep the state instead of removing a
// possibly billed backup.
func TestUnitAliCloudPolarDBBackupReadKeepsRetainedAbsence(t *testing.T) {
	withBackupUnitTimings(t, 5*time.Millisecond, 50*time.Millisecond, time.Second)
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DescribeBackupPolicy", "DescribeBackups":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeDetachedBackups":
			return backupUnitBackupsJSON()
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":     backupUnitString(backupUnitClusterID),
		"backup_id":         backupUnitString(backupUnitBackupID),
		"backup_start_time": backupUnitString(backupUnitStartTime),
		"backup_retention_policy_on_cluster_deletion": backupUnitString("ALL"),
	})

	var resp fwresource.ReadResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed although the retention contract keeps the backup set: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if len(resp.Diagnostics.Warnings()) == 0 {
		t.Fatal("Read did not warn that the retained backup set is not listed yet")
	}
	if id, ok := backupUnitStateAttribute(t, resp.State, "id"); !ok || id == "" {
		t.Fatal("Read removed a backup set that retention ALL keeps")
	}
}

// TestUnitAliCloudPolarDBBackupReadRemovesUnretainedAbsence covers the other half of N2: when the
// retention policy does not keep the backup set, repeated complete empty answers remove it.
func TestUnitAliCloudPolarDBBackupReadRemovesUnretainedAbsence(t *testing.T) {
	withBackupUnitTimings(t, 5*time.Millisecond, time.Second, time.Second)
	client, server := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DescribeBackupPolicy", "DescribeBackups":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeDetachedBackups":
			return backupUnitBackupsJSON()
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":     backupUnitString(backupUnitClusterID),
		"backup_id":         backupUnitString(backupUnitBackupID),
		"backup_start_time": backupUnitString(backupUnitStartTime),
		"backup_retention_policy_on_cluster_deletion": backupUnitString("NONE"),
	})

	var resp fwresource.ReadResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed while removing an unretained backup set: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if !server.seen("DescribeDetachedBackups") {
		t.Fatal("Read removed the backup set without consulting the released-cluster list")
	}
	if _, ok := backupUnitStateAttribute(t, resp.State, "id"); ok {
		t.Fatal("Read kept a backup set that is gone and that retention does not keep")
	}
}

// TestUnitAliCloudPolarDBBackupDeletePropagatesNonNotFound covers QAP-01: a DeleteBackup failure
// that is not a NotFound must propagate; it must not be reinterpreted as "retained" success.
func TestUnitAliCloudPolarDBBackupDeletePropagatesNonNotFound(t *testing.T) {
	withBackupUnitTimings(t, 5*time.Millisecond, time.Second, time.Second)
	client, server := newBackupUnitClient(t, func(action string, params map[string]string) string {
		if action == "DeleteBackup" {
			return `{"Code":"InvalidParameter","Message":"the backup set can not be deleted"}`
		}
		if action == "DescribeDetachedBackups" {
			return backupUnitBackupsJSON(backupUnitBackupEntry(backupUnitBackupID, backupUnitClusterID, backupUnitStartTime))
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":            backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id": backupUnitString(backupUnitClusterID),
		"backup_id":     backupUnitString(backupUnitBackupID),
	})

	var resp fwresource.DeleteResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("Delete swallowed a non-NotFound DeleteBackup failure: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if server.seen("DescribeDetachedBackups") {
		t.Fatal("Delete consulted the retention list before classifying the original error")
	}
}

// TestUnitAliCloudPolarDBBackupDeleteRetainedWarning covers the G1 delete branch: DeleteBackup
// reports the released cluster as absent, the detached list still reports the exact set, so the
// destroy stops tracking it with a warning instead of failing or claiming a deletion.
func TestUnitAliCloudPolarDBBackupDeleteRetainedWarning(t *testing.T) {
	withBackupUnitTimings(t, 5*time.Millisecond, time.Second, time.Second)
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DeleteBackup":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeDetachedBackups":
			return backupUnitBackupsJSON(backupUnitBackupEntry(backupUnitBackupID, backupUnitClusterID, backupUnitStartTime))
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":     backupUnitString(backupUnitClusterID),
		"backup_id":         backupUnitString(backupUnitBackupID),
		"backup_start_time": backupUnitString(backupUnitStartTime),
	})

	var resp fwresource.DeleteResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete failed although the backup set is retained by a released cluster: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if len(resp.Diagnostics.Warnings()) == 0 {
		t.Fatal("Delete did not warn that the retained backup set stays in the cloud")
	}
}

// TestUnitAliCloudPolarDBBackupUpdateRejectsSameCountReplacement covers QAP-02: replacing one
// advanced policy with a different one keeps the list length but removes an entry, so Update must
// reject it instead of leaving the old cloud policy unmanaged.
func TestUnitAliCloudPolarDBBackupUpdateRejectsSameCountReplacement(t *testing.T) {
	withBackupUnitTimings(t, 5*time.Millisecond, time.Second, time.Second)
	client, server := newBackupUnitClient(t, func(action string, params map[string]string) string {
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)

	policy := func(srcType, destType, retention string) map[string]tftypes.Value {
		return map[string]tftypes.Value{
			"action_type":     backupUnitString("UPDATE"),
			"bak_type":        backupUnitString("F"),
			"src_type":        backupUnitString(srcType),
			"dest_type":       backupUnitString(destType),
			"src_region":      backupUnitString(backupUnitRegion),
			"dest_region":     backupUnitString(backupUnitRegion),
			"retention_type":  backupUnitString("delay"),
			"retention_value": backupUnitString(retention),
		}
	}
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                     backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":          backupUnitString(backupUnitClusterID),
		"backup_id":              backupUnitString(backupUnitBackupID),
		"backup_policy_level":    backupUnitString("Advanced"),
		"advanced_data_policies": backupUnitPolicyList(t, objectType, []map[string]tftypes.Value{policy("db", "level1", "7")}),
	})
	plan := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                     backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":          backupUnitString(backupUnitClusterID),
		"backup_id":              backupUnitString(backupUnitBackupID),
		"backup_policy_level":    backupUnitString("Advanced"),
		"advanced_data_policies": backupUnitPolicyList(t, objectType, []map[string]tftypes.Value{policy("db", "level2", "7")}),
	})

	var resp fwresource.UpdateResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Update(context.Background(),
		fwresource.UpdateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: plan}, State: tfsdk.State{Schema: sch, Raw: state}},
		&resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("Update accepted a same-sized advanced-data-policy replacement: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if server.seen("ModifyBackupPolicy") {
		t.Fatal("Update sent ModifyBackupPolicy although an advanced entry was implicitly removed")
	}
}

// TestUnitAliCloudPolarDBBackupUpdateAdvancedFieldsSerialized covers G4: every advanced-policy
// field the configuration can carry must reach ModifyBackupPolicy with its own value, and every
// readable field must come back into state. Fields that a cloud run cannot safely construct
// (auto_created, policy_id are server-assigned) are covered here with fixed values.
func TestUnitAliCloudPolarDBBackupUpdateAdvancedFieldsSerialized(t *testing.T) {
	withBackupUnitTimings(t, 5*time.Millisecond, time.Second, time.Second)
	region := backupUnitRegion
	entry := map[string]tftypes.Value{
		"action_type":                 backupUnitString("UPDATE"),
		"only_preserve_one_each_hour": backupUnitBool(true),
		"storage_class":               backupUnitString("standard"),
		"retention_value":             backupUnitString("14"),
		"bak_type":                    backupUnitString("F"),
		"filter_key":                  backupUnitString("dayOfWeek"),
		"filter_value":                backupUnitString("1,2,3,4,5,6,7"),
		"dest_type":                   backupUnitString("level1"),
		"filter_type":                 backupUnitString("crontab"),
		"src_region":                  backupUnitString(region),
		"auto_created":                backupUnitBool(false),
		"only_preserve_one_each_day":  backupUnitBool(true),
		"src_type":                    backupUnitString("db"),
		"dump_action":                 backupUnitString("copy"),
		"dest_region":                 backupUnitString(region),
		"policy_id":                   backupUnitString("pid-unit"),
		"retention_type":              backupUnitString("delay"),
	}
	var modifyMu sync.Mutex
	var modifyRequest map[string]string
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "ModifyBackupPolicy":
			modifyMu.Lock()
			modifyRequest = params
			modifyMu.Unlock()
			return "{}"
		case "DescribeBackupPolicy":
			return fmt.Sprintf(`{"DBClusterId":%q,"BackupPolicyLevel":"Advanced","BackupRetentionPolicyOnClusterDeletion":"ALL",`+
				`"AdvancedDataPolicies":{"AdvancedDataPolicy":[{"AutoCreated":false,"BakType":"F","DestRegion":%q,"DestType":"level1",`+
				`"DumpAction":"copy","FilterKey":"dayOfWeek","FilterType":"crontab","FilterValue":"1,2,3,4,5,6,7",`+
				`"OnlyPreserveOneEachDay":true,"OnlyPreserveOneEachHour":true,"PolicyId":"pid-unit","RetentionType":"delay",`+
				`"RetentionValue":"14","SrcRegion":%q,"SrcType":"db","StorageClass":"standard"}]}}`,
				backupUnitClusterID, region, region)
		case "DescribeBackups":
			return backupUnitBackupsJSON(backupUnitBackupEntry(backupUnitBackupID, backupUnitClusterID, backupUnitStartTime))
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)

	prior := map[string]tftypes.Value{}
	for key, value := range entry {
		prior[key] = value
	}
	prior["retention_value"] = backupUnitString("7")
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                     backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":          backupUnitString(backupUnitClusterID),
		"backup_id":              backupUnitString(backupUnitBackupID),
		"backup_start_time":      backupUnitString(backupUnitStartTime),
		"backup_policy_level":    backupUnitString("Advanced"),
		"advanced_data_policies": backupUnitPolicyList(t, objectType, []map[string]tftypes.Value{prior}),
	})
	plan := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                     backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":          backupUnitString(backupUnitClusterID),
		"backup_id":              backupUnitString(backupUnitBackupID),
		"backup_start_time":      backupUnitString(backupUnitStartTime),
		"backup_policy_level":    backupUnitString("Advanced"),
		"advanced_data_policies": backupUnitPolicyList(t, objectType, []map[string]tftypes.Value{entry}),
	})

	var resp fwresource.UpdateResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Update(context.Background(),
		fwresource.UpdateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: plan}, State: tfsdk.State{Schema: sch, Raw: state}},
		&resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed on the full advanced policy: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	modifyMu.Lock()
	captured := modifyRequest
	modifyMu.Unlock()
	if captured == nil {
		t.Fatal("Update did not call ModifyBackupPolicy for the changed advanced policy")
	}
	// RPC requests carry nested entries either as AdvancedDataPolicies.N.<Field> form keys or
	// inside one JSON body; match on the leaf name so the assertion does not depend on which.
	leaf := func(name, want string) {
		for key, value := range captured {
			if strings.EqualFold(key, name) || strings.HasSuffix(strings.ToLower(key), "."+strings.ToLower(name)) {
				if strings.EqualFold(value, want) {
					return
				}
			}
		}
		t.Errorf("ModifyBackupPolicy request does not carry %s=%q (params: %v)", name, want, captured)
	}
	leaf("DBClusterId", backupUnitClusterID)
	leaf("ActionType", "UPDATE")
	leaf("BakType", "F")
	leaf("DestType", "level1")
	leaf("SrcType", "db")
	leaf("SrcRegion", region)
	leaf("DestRegion", region)
	leaf("FilterType", "crontab")
	leaf("FilterKey", "dayOfWeek")
	leaf("FilterValue", "1,2,3,4,5,6,7")
	leaf("OnlyPreserveOneEachDay", "true")
	leaf("OnlyPreserveOneEachHour", "true")
	leaf("RetentionType", "delay")
	leaf("RetentionValue", "14")
	leaf("StorageClass", "standard")
	leaf("DumpAction", "copy")
	leaf("PolicyId", "pid-unit")
	leaf("AutoCreated", "false")
	payloadPresent := false
	for key := range captured {
		if strings.Contains(key, "AdvancedDataPolicies.") {
			payloadPresent = true
			break
		}
	}
	if !payloadPresent {
		t.Errorf("ModifyBackupPolicy request carries no AdvancedDataPolicies payload (params: %v)", captured)
	}
	// G4 read-back: the API-readable advanced policy fields are mapped back into state.
	if value, ok := backupUnitPolicyAttribute(t, resp.State, "retention_value"); !ok || value != "14" {
		t.Errorf("advanced_data_policies.retention_value in state = %q present=%v, want 14", value, ok)
	}
	if value, ok := backupUnitPolicyAttribute(t, resp.State, "policy_id"); !ok || value != "pid-unit" {
		t.Errorf("advanced_data_policies.policy_id in state = %q present=%v, want pid-unit", value, ok)
	}
}

// backupUnitPolicyAttribute returns the first known value of a string attribute inside the
// configured advanced_data_policies block of the read-back state.
func backupUnitPolicyAttribute(t *testing.T, state tfsdk.State, name string) (string, bool) {
	t.Helper()
	var list types.List
	if diagnostics := state.GetAttribute(context.Background(), fwpath.Root("advanced_data_policies"), &list); diagnostics.HasError() {
		return "", false
	}
	if list.IsNull() || list.IsUnknown() {
		return "", false
	}
	for _, element := range list.Elements() {
		policy, ok := element.(types.Object)
		if !ok {
			continue
		}
		value, ok := policy.Attributes()[name].(types.String)
		if ok && !value.IsNull() && !value.IsUnknown() {
			return value.ValueString(), true
		}
	}
	return "", false
}
