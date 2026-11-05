package ecs

// The tests in this file drive the data source read path against a real HTTP
// server. They pin the transport behaviour (connection reuse across pages and
// cancellation while a page is in flight) and the request the real Read entry
// builds, which a request-map fixture cannot show.

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newEcsLoopbackClient points the ECS product at the given test server, so the
// whole SDK stack (signing, transport, decoding) runs against real HTTP.
func newEcsLoopbackClient(t *testing.T, serverURL string) (*connectivity.AliyunClient, string) {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	require.NoError(t, err)

	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(serverURL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	client, err := (&connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		ClientReadTimeout:    60000,
		ClientConnectTimeout: 60000,
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}).Client()
	require.NoError(t, err)
	endpoints.Store("ecs", endpoint)
	return client, endpoint
}

// ecsReadReservedInstances drives one real Read of the reserved instance data
// source. The configuration only carries strings, which is all these tests set.
func ecsReadReservedInstances(t *testing.T, ctx context.Context, client *connectivity.AliyunClient, config map[string]string) (*datasource.ReadResponse, *reservedInstancesModel) {
	t.Helper()
	ds := NewReservedInstancesDataSource()

	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	require.True(t, ok, "the reserved instance data source must accept the provider meta")
	var configureResponse datasource.ConfigureResponse
	configurable.Configure(ctx, datasource.ConfigureRequest{ProviderData: client}, &configureResponse)
	require.False(t, configureResponse.Diagnostics.HasError(), "configure diagnostics: %+v", configureResponse.Diagnostics)

	var schemaResponse datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResponse)
	require.False(t, schemaResponse.Diagnostics.HasError(), "schema diagnostics: %+v", schemaResponse.Diagnostics)

	objectType, ok := schemaResponse.Schema.Type().TerraformType(ctx).(tftypes.Object)
	require.True(t, ok, "unexpected data source schema type %T", schemaResponse.Schema.Type().TerraformType(ctx))

	raw := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attrType := range objectType.AttributeTypes {
		raw[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range config {
		attrType, ok := objectType.AttributeTypes[name]
		require.True(t, ok, "config attribute %q is not in the schema", name)
		raw[name] = tftypes.NewValue(attrType, value)
	}

	response := &datasource.ReadResponse{
		State: tfsdk.State{Schema: schemaResponse.Schema, Raw: tftypes.NewValue(objectType, nil)},
	}
	ds.Read(ctx, datasource.ReadRequest{
		Config: tfsdk.Config{Schema: schemaResponse.Schema, Raw: tftypes.NewValue(objectType, raw)},
	}, response)
	if response.Diagnostics.HasError() {
		return response, nil
	}

	var state reservedInstancesModel
	require.False(t, response.State.Get(ctx, &state).HasError(), "state decode diagnostics")
	return response, &state
}

func writeReservedInstancesPage(t *testing.T, w http.ResponseWriter, instances []map[string]interface{}, total int) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"RequestId":  "mock",
		"TotalCount": total,
		"ReservedInstances": map[string]interface{}{
			"ReservedInstance": instances,
		},
	}); err != nil {
		t.Error(err)
	}
}

// The pages of one read must share the process-level connection pool. The SDK
// keys its client and transport on the runtime timeouts, so a per-page budget
// would open a new connection per page; the read deadline travels with the
// context instead, and the pages are read over the same keep-alive connection.
func TestUnitFwCollectNumberedPagesReusesOneConnectionAcrossPages(t *testing.T) {
	var mu sync.Mutex
	remoteAddrs := make([]string, 0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
		}
		mu.Lock()
		remoteAddrs = append(remoteAddrs, r.RemoteAddr)
		mu.Unlock()

		switch r.Form.Get("PageNumber") {
		case "1":
			writeReservedInstancesPage(t, w, []map[string]interface{}{
				{"ReservedInstanceId": "ri-1"},
				{"ReservedInstanceId": "ri-2"},
			}, 4)
		case "2":
			writeReservedInstancesPage(t, w, []map[string]interface{}{
				{"ReservedInstanceId": "ri-3"},
				{"ReservedInstanceId": "ri-4"},
			}, 4)
		default:
			t.Errorf("unexpected page %q", r.Form.Get("PageNumber"))
		}
	}))
	t.Cleanup(server.Close)
	client, _ := newEcsLoopbackClient(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	objects, totalCount, err := fwCollectNumberedPages(ctx, 10, 1, 2, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
		request := map[string]interface{}{
			"RegionId":   "cn-hangzhou",
			"PageNumber": int(pageNumber),
			"PageSize":   int(2),
		}
		response, err := fwRpcPost(ctx, client, reservedInstancesTypeName, reservedInstancesAction, request)
		if err != nil {
			return fwNumberedPage{}, err
		}
		items, err := fwDecodeRecordList(response, "ReservedInstances", "ReservedInstance", reservedInstancesAction)
		if err != nil {
			return fwNumberedPage{}, err
		}
		total, totalKnown, err := fwDecodeCount(response, "TotalCount", reservedInstancesAction)
		if err != nil {
			return fwNumberedPage{}, err
		}
		return fwNumberedPage{items: items, total: total, totalKnown: totalKnown}, nil
	})

	require.NoError(t, err)
	require.Len(t, objects, 4)
	assert.Equal(t, int64(4), totalCount)
	mu.Lock()
	defer mu.Unlock()
	require.Len(t, remoteAddrs, 2, "both pages must be read")
	assert.Equal(t, remoteAddrs[0], remoteAddrs[1], "the pages must reuse the same keep-alive connection")
}

// Cancelling the read while a page is in flight must abort that page and stop
// the loop: the HTTP exchange is bound to the read context, not only to a
// numeric timeout.
func TestUnitFwCollectNumberedPagesAbortsAnInFlightPageOnCancel(t *testing.T) {
	pageTwoArrived := make(chan struct{})
	var arrivedOnce sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("PageNumber") == "1" {
			writeReservedInstancesPage(t, w, []map[string]interface{}{
				{"ReservedInstanceId": "ri-1"},
				{"ReservedInstanceId": "ri-2"},
			}, 4)
			return
		}
		arrivedOnce.Do(func() { close(pageTwoArrived) })
		<-r.Context().Done()
	}))
	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})
	client, _ := newEcsLoopbackClient(t, server.URL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type outcome struct {
		objects []map[string]interface{}
		err     error
	}
	done := make(chan outcome, 1)
	go func() {
		objects, _, err := fwCollectNumberedPages(ctx, 10, 1, 2, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
			request := map[string]interface{}{
				"RegionId":   "cn-hangzhou",
				"PageNumber": int(pageNumber),
				"PageSize":   int(2),
			}
			response, err := fwRpcPost(ctx, client, reservedInstancesTypeName, reservedInstancesAction, request)
			if err != nil {
				return fwNumberedPage{}, err
			}
			items, err := fwDecodeRecordList(response, "ReservedInstances", "ReservedInstance", reservedInstancesAction)
			if err != nil {
				return fwNumberedPage{}, err
			}
			return fwNumberedPage{items: items, total: 4, totalKnown: true}, nil
		})
		done <- outcome{objects: objects, err: err}
	}()

	select {
	case <-pageTwoArrived:
	case <-time.After(10 * time.Second):
		t.Fatal("the second page was never requested")
	}
	start := time.Now()
	cancel()

	select {
	case result := <-done:
		assert.Error(t, result.err, "a cancelled read must fail instead of returning a partial result")
		assert.Nil(t, result.objects)
		assert.Less(t, time.Since(start), 4*time.Second, "the in-flight page must be aborted by the read context")
	case <-time.After(15 * time.Second):
		t.Fatal("the read did not stop after the context was cancelled")
	}
}

// The reserved instance data source derives Scope = Zone from zone_id in the
// real Read entry; the request that reaches the server is what proves it.
func TestUnitEcsReservedInstancesReadInfersTheZoneScope(t *testing.T) {
	var mu sync.Mutex
	scope := ""
	zoneId := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
		}
		mu.Lock()
		scope = r.Form.Get("Scope")
		zoneId = r.Form.Get("ZoneId")
		mu.Unlock()
		writeReservedInstancesPage(t, w, []map[string]interface{}{
			{"ReservedInstanceId": "ri-zone", "Status": "Active", "ZoneId": "cn-hangzhou-k", "Scope": "Zone"},
		}, 1)
	}))
	t.Cleanup(server.Close)
	client, _ := newEcsLoopbackClient(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := ecsReadReservedInstances(t, ctx, client, map[string]string{
		"zone_id": "cn-hangzhou-k",
	})
	require.False(t, response.Diagnostics.HasError(), "diagnostics: %+v", response.Diagnostics)
	require.NotNil(t, state)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, "Zone", scope, "zone_id without scope must request Scope = Zone")
	assert.Equal(t, "cn-hangzhou-k", zoneId)
	require.Len(t, state.Ids.Elements(), 1)
}

// The same real Read entry must reject a configuration that asks for a zonal
// query without a zone, before it sends anything.
func TestUnitEcsReservedInstancesReadRejectsZoneScopeWithoutZone(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		writeReservedInstancesPage(t, w, nil, 0)
	}))
	t.Cleanup(server.Close)
	client, _ := newEcsLoopbackClient(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := ecsReadReservedInstances(t, ctx, client, map[string]string{
		"scope": "Zone",
	})
	assert.True(t, response.Diagnostics.HasError(), "scope = Zone without zone_id must be rejected")
	assert.Nil(t, state)
	assert.Zero(t, calls, "the rejected configuration must not reach the API")
}

// A 2xx response without a JSON object body must fail the read instead of being
// written as an empty result.
func TestUnitEcsReservedInstancesReadRejectsAnEmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	client, _ := newEcsLoopbackClient(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, state := ecsReadReservedInstances(t, ctx, client, map[string]string{
		"zone_id": "cn-hangzhou-k",
	})
	assert.True(t, response.Diagnostics.HasError(), "an empty 2xx body must not be written as an empty result")
	assert.Nil(t, state)
	for _, err := range response.Diagnostics.Errors() {
		assert.Contains(t, fmt.Sprintf("%s %s", err.Summary(), err.Detail()), "without a JSON object body")
	}
}
