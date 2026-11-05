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

// ecsReadDataSource drives one real Read of any ECS Framework data source
// against the loopback client. The configuration only carries strings, which is
// all these tests set. The returned state stays null when the read fails, which
// is how a test proves that no state was written.
func ecsReadDataSource(t *testing.T, ctx context.Context, client *connectivity.AliyunClient, ds datasource.DataSource, config map[string]string) *datasource.ReadResponse {
	t.Helper()
	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	require.True(t, ok, "the data source must accept the provider meta")
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
	return response
}

// ecsReadReservedInstances drives one real Read of the reserved instance data
// source and decodes the state it wrote.
func ecsReadReservedInstances(t *testing.T, ctx context.Context, client *connectivity.AliyunClient, config map[string]string) (*datasource.ReadResponse, *reservedInstancesModel) {
	t.Helper()
	response := ecsReadDataSource(t, ctx, client, NewReservedInstancesDataSource(), config)
	if response.Diagnostics.HasError() {
		return response, nil
	}

	var state reservedInstancesModel
	require.False(t, response.State.Get(ctx, &state).HasError(), "state decode diagnostics")
	return response, &state
}

// ecsReadAutoProvisioningGroups drives one real Read of the auto provisioning
// group data source and decodes the state it wrote.
func ecsReadAutoProvisioningGroups(t *testing.T, ctx context.Context, client *connectivity.AliyunClient, config map[string]string) (*datasource.ReadResponse, *autoProvisioningGroupsModel) {
	t.Helper()
	response := ecsReadDataSource(t, ctx, client, NewAutoProvisioningGroupsDataSource(), config)
	if response.Diagnostics.HasError() {
		return response, nil
	}

	var state autoProvisioningGroupsModel
	require.False(t, response.State.Get(ctx, &state).HasError(), "state decode diagnostics")
	return response, &state
}

// ecsReadImagePipelineExecutions drives one real Read of the image pipeline
// execution data source and decodes the state it wrote.
func ecsReadImagePipelineExecutions(t *testing.T, ctx context.Context, client *connectivity.AliyunClient, config map[string]string) (*datasource.ReadResponse, *imagePipelineExecutionsModel) {
	t.Helper()
	response := ecsReadDataSource(t, ctx, client, NewImagePipelineExecutionsDataSource(), config)
	if response.Diagnostics.HasError() {
		return response, nil
	}

	var state imagePipelineExecutionsModel
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
	response, _ := ecsReadReservedInstances(t, ctx, client, map[string]string{
		"scope": "Zone",
	})
	assert.True(t, response.Diagnostics.HasError(), "scope = Zone without zone_id must be rejected")
	assert.True(t, response.State.Raw.IsNull(), "a rejected configuration must not write state")
	assert.Zero(t, calls, "the rejected configuration must not reach the API")
}

// A 2xx body that is a JSON object but not the success envelope of the query
// must fail the read: decoding it would write a zero result over the state even
// though the service never answered the question.
func TestUnitEcsReservedInstancesReadRejectsANonEnvelopeBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"unexpected": "document",
		}); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(server.Close)
	client, _ := newEcsLoopbackClient(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, _ := ecsReadReservedInstances(t, ctx, client, map[string]string{
		"zone_id": "cn-hangzhou-k",
	})
	assert.True(t, response.Diagnostics.HasError(), "a non-envelope 2xx document must not be written as an empty result")
	assert.True(t, response.State.Raw.IsNull(), "no state may be written when the response is not the query envelope")
	for _, err := range response.Diagnostics.Errors() {
		assert.Contains(t, fmt.Sprintf("%s %s", err.Summary(), err.Detail()), "success envelope")
	}
}

// The legitimate empty results keep working: an explicit empty list, and the
// envelope the service returns when it omits the container because nothing
// matched.
func TestUnitEcsReservedInstancesReadKeepsLegitimateEmptyResults(t *testing.T) {
	cases := map[string]map[string]interface{}{
		"explicit empty list": {
			"RequestId":         "mock",
			"TotalCount":        0,
			"ReservedInstances": map[string]interface{}{"ReservedInstance": []interface{}{}},
		},
		"container omitted": {
			"RequestId":  "mock",
			"TotalCount": 0,
		},
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(body); err != nil {
					t.Error(err)
				}
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
			assert.Empty(t, state.Ids.Elements(), "an empty result must keep an empty id list")
			assert.Equal(t, int64(0), state.TotalCount.ValueInt64())
		})
	}
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
	response, _ := ecsReadReservedInstances(t, ctx, client, map[string]string{
		"zone_id": "cn-hangzhou-k",
	})
	assert.True(t, response.Diagnostics.HasError(), "an empty 2xx body must not be written as an empty result")
	assert.True(t, response.State.Raw.IsNull(), "no state may be written for an empty body")
	for _, err := range response.Diagnostics.Errors() {
		assert.Contains(t, fmt.Sprintf("%s %s", err.Summary(), err.Detail()), "without a JSON object body")
	}
}

// newEcsStaticServer answers every request with the same status and body, so a
// real Read can be pinned to one response shape.
func newEcsStaticServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != "" {
			_, _ = w.Write([]byte(body))
		}
	}))
	t.Cleanup(server.Close)
	return server
}

// ecsDataSourceRead drives one real Read of one data source and reports the
// number of identities and the total count it wrote.
type ecsDataSourceRead struct {
	name string
	read func(t *testing.T, ctx context.Context, client *connectivity.AliyunClient) (*datasource.ReadResponse, int, int64)
}

func ecsDataSourceReads() []ecsDataSourceRead {
	return []ecsDataSourceRead{
		{
			name: "alicloud_ecs_auto_provisioning_groups",
			read: func(t *testing.T, ctx context.Context, client *connectivity.AliyunClient) (*datasource.ReadResponse, int, int64) {
				response, state := ecsReadAutoProvisioningGroups(t, ctx, client, map[string]string{"auto_provisioning_group_name": "target-group"})
				if state == nil {
					return response, 0, 0
				}
				return response, len(state.Ids.Elements()), state.TotalCount.ValueInt64()
			},
		},
		{
			name: "alicloud_ecs_image_pipeline_executions",
			read: func(t *testing.T, ctx context.Context, client *connectivity.AliyunClient) (*datasource.ReadResponse, int, int64) {
				response, state := ecsReadImagePipelineExecutions(t, ctx, client, map[string]string{"image_pipeline_id": "pipeline-1"})
				if state == nil {
					return response, 0, 0
				}
				return response, len(state.Ids.Elements()), state.TotalCount.ValueInt64()
			},
		},
		{
			name: "alicloud_reserved_instances",
			read: func(t *testing.T, ctx context.Context, client *connectivity.AliyunClient) (*datasource.ReadResponse, int, int64) {
				response, state := ecsReadReservedInstances(t, ctx, client, map[string]string{"zone_id": "cn-hangzhou-k"})
				if state == nil {
					return response, 0, 0
				}
				return response, len(state.Ids.Elements()), state.TotalCount.ValueInt64()
			},
		},
	}
}

func ecsDiagnosticsText(response *datasource.ReadResponse) string {
	var builder strings.Builder
	for _, err := range response.Diagnostics.Errors() {
		builder.WriteString(err.Summary())
		builder.WriteString(" ")
		builder.WriteString(err.Detail())
		builder.WriteString(" ")
	}
	return builder.String()
}

// A 2xx document that only carries a generic RequestId is not the empty success
// envelope of any of the three queries: the read must fail and must not write a
// zero result over the state.
func TestUnitEcsDataSourceReadsRejectAGenericRequestIdBody(t *testing.T) {
	for _, readCase := range ecsDataSourceReads() {
		t.Run(readCase.name, func(t *testing.T) {
			server := newEcsStaticServer(t, http.StatusOK, `{"RequestId":"fixture","error":"upstream"}`)
			client, _ := newEcsLoopbackClient(t, server.URL)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			response, _, _ := readCase.read(t, ctx, client)
			require.True(t, response.Diagnostics.HasError(), "a generic RequestId body must not be written as an empty result")
			assert.True(t, response.State.Raw.IsNull(), "no state may be written when the body is not the query envelope")
			assert.Contains(t, ecsDiagnosticsText(response), "success envelope")
		})
	}
}

// The documented empty envelope, which omits the container but keeps the
// RequestId and the total count, is a successful empty read for all three
// sources.
func TestUnitEcsDataSourceReadsKeepTheEmptyEnvelope(t *testing.T) {
	for _, readCase := range ecsDataSourceReads() {
		t.Run(readCase.name, func(t *testing.T) {
			server := newEcsStaticServer(t, http.StatusOK, `{"RequestId":"mock","TotalCount":0}`)
			client, _ := newEcsLoopbackClient(t, server.URL)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			response, ids, total := readCase.read(t, ctx, client)
			require.False(t, response.Diagnostics.HasError(), "diagnostics: %+v", response.Diagnostics)
			assert.False(t, response.State.Raw.IsNull(), "a successful read writes its empty result")
			assert.Zero(t, ids, "the empty envelope must read as an empty result")
			assert.Zero(t, total)
		})
	}
}

// A 2xx envelope that reports a positive TotalCount but carries neither the
// records nor a continuation token does not describe an empty result: the read
// must fail instead of writing zero rows over the state.
func TestUnitEcsDataSourceReadsRejectAPositiveTotalWithoutRecords(t *testing.T) {
	for _, readCase := range ecsDataSourceReads() {
		t.Run(readCase.name, func(t *testing.T) {
			server := newEcsStaticServer(t, http.StatusOK, `{"RequestId":"mock","TotalCount":3}`)
			client, _ := newEcsLoopbackClient(t, server.URL)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			response, _, _ := readCase.read(t, ctx, client)
			require.True(t, response.Diagnostics.HasError(), "a positive TotalCount without records is not an empty result")
			assert.True(t, response.State.Raw.IsNull(), "no state may be written for an inconsistent page")
			assert.Contains(t, ecsDiagnosticsText(response), "TotalCount")
		})
	}
}

// A page that carries fewer records than its TotalCount and then ends the
// listing is incomplete, and the read must not write the partial result.
func TestUnitEcsImagePipelineExecutionsReadRejectsAnIncompleteTotal(t *testing.T) {
	server := newEcsStaticServer(t, http.StatusOK, `{"RequestId":"mock","TotalCount":2,"ImagePipelineExecution":{"ImagePipelineExecutionSet":[{"ExecutionId":"exec-1","Status":"SUCCESS"}]}}`)
	client, _ := newEcsLoopbackClient(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, _ := ecsReadImagePipelineExecutions(t, ctx, client, map[string]string{"image_pipeline_id": "pipeline-1"})
	require.True(t, response.Diagnostics.HasError(), "a partial page must not be written as the complete result")
	assert.True(t, response.State.Raw.IsNull(), "no state may be written for an incomplete result")
	assert.Contains(t, ecsDiagnosticsText(response), "TotalCount")
}

// A failed call keeps its SDK code and message through all three read paths;
// the failure is never treated as an empty envelope.
func TestUnitEcsDataSourceReadsKeepTheSDKError(t *testing.T) {
	for _, readCase := range ecsDataSourceReads() {
		t.Run(readCase.name, func(t *testing.T) {
			server := newEcsStaticServer(t, http.StatusForbidden, `{"RequestId":"mock","Code":"Forbidden.RAM","Message":"denied"}`)
			client, _ := newEcsLoopbackClient(t, server.URL)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			response, _, _ := readCase.read(t, ctx, client)
			require.True(t, response.Diagnostics.HasError(), "a failed call must stay failed")
			assert.True(t, response.State.Raw.IsNull(), "a failed call must not write state")
			text := ecsDiagnosticsText(response)
			assert.Contains(t, text, "Forbidden.RAM", "the SDK error code must be preserved")
			assert.Contains(t, text, "denied", "the SDK error message must be preserved")
		})
	}
}

// A 2xx error-shaped body may carry a numeric StatusCode in any shape the
// transports produce: reading it must never panic and must keep the code and
// message the service reported. This drives the real Read instead of calling
// the formatter directly.
func TestUnitEcsReadKeepsANumericErrorStatus(t *testing.T) {
	bodies := map[string]string{
		"json number":      `{"RequestId":"mock","Code":"InternalError","Message":"upstream","StatusCode":503}`,
		"malformed status": `{"RequestId":"mock","Code":"InternalError","Message":"upstream","StatusCode":{"unexpected":1}}`,
		"text status":      `{"RequestId":"mock","Code":"InternalError","Message":"upstream","StatusCode":"not-a-number"}`,
	}
	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			server := newEcsStaticServer(t, http.StatusOK, body)
			client, _ := newEcsLoopbackClient(t, server.URL)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			response, _ := ecsReadReservedInstances(t, ctx, client, map[string]string{"zone_id": "cn-hangzhou-k"})
			require.True(t, response.Diagnostics.HasError(), "a 2xx error body must fail the read instead of panicking")
			// The helper returns a nil model whenever the read reported a
			// diagnostic, so a nil model cannot prove the read left the state
			// untouched. The raw state is checked instead: a read that wrote a
			// result and then failed would still be caught here.
			assert.True(t, response.State.Raw.IsNull(), "no state may be written for an error body")
			text := ecsDiagnosticsText(response)
			assert.Contains(t, text, "InternalError", "the error code must be preserved")
			assert.Contains(t, text, "upstream", "the error message must be preserved")
		})
	}
}
