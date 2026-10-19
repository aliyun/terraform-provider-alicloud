package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// rocketmqTestClient stands up a credential-free AliyunClient whose RocketMQ
// ROA transport is pointed at an in-process httptest server, mirroring the
// rds-account unit-test construction. Requests are signed with throwaway
// AK/SK; the server matches on path and method and returns canned JSON.
func rocketmqTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(handler))
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
		RegionId: "cn-beijing", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("rocketmq", endpoint)
	return client
}

func writeJSON(t *testing.T, w http.ResponseWriter, v interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Error(err)
	}
}

// planObject is the GetDisasterRecoveryPlan response payload. It deliberately
// omits the write-exclusive fields (password, endpointUrl, networkType, vpcId,
// vSwitchId) the live API never returns, so carry-over is exercised; it keeps
// consumerGroupId, which is the property this governance item restores.
func planObject() map[string]interface{} {
	return map[string]interface{}{
		"planId":                "plan-900",
		"planName":              "tf-plan",
		"planType":              "FILE",
		"planStatus":            "RUNNING",
		"autoSyncCheckpoint":    true,
		"syncCheckpointEnabled": false,
		"createTime":            "2026-10-09T00:00:00Z",
		"updateTime":            "2026-10-09T00:00:00Z",
		"instances": []interface{}{map[string]interface{}{
			"instanceId":      "i-1",
			"consumerGroupId": "cg-1",
			"securityGroupId": "sg-1",
			"username":        "u-1",
			"instanceRole":    "ACTIVE",
			"instanceType":    "ALIYUN_ROCKETMQ",
			"regionId":        "cn-beijing",
			"authType":        "ACL_AUTH",
		}},
	}
}

// TestUnitRocketmqDisasterRecoveryPlanServiceMock drives the real ROA HTTP
// round-trip (POST create, GET describe, DELETE) against an httptest server
// and verifies consumerGroupId lands in the create body, the API omits the
// write-exclusive fields, Read carry-over restores them, and Delete settles
// via NotFound.
func TestUnitRocketmqDisasterRecoveryPlanServiceMock(t *testing.T) {
	var deleted atomic.Bool
	var createBodyRaw string
	var createBody map[string]interface{}

	client := rocketmqTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			raw, _ := io.ReadAll(r.Body)
			createBodyRaw = string(raw)
			_ = json.Unmarshal(raw, &createBody)
			writeJSON(t, w, map[string]interface{}{"data": "plan-900"})
		case http.MethodGet:
			if strings.Contains(r.URL.Path, "/disaster_recovery/") {
				if deleted.Load() {
					// The live ROA API returns a JSON error envelope on 404,
					// not an empty body: the tea client's AssertAsMap panics on
					// nil, so the mock must mirror that shape. isPlanNotFound
					// then maps the Code/404 to a gone signal.
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusNotFound)
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"Code":    "DisasterRecoveryPlan.NotFound",
						"Message": "The specified disaster recovery plan does not exist.",
						"Success": false,
					})
					return
				}
				writeJSON(t, w, map[string]interface{}{"data": planObject()})
				return
			}
			writeJSON(t, w, map[string]interface{}{"data": map[string]interface{}{
				"list":       []interface{}{planObject()},
				"totalCount": float64(1),
			}})
		case http.MethodPatch:
			writeJSON(t, w, map[string]interface{}{})
		case http.MethodDelete:
			deleted.Store(true)
			writeJSON(t, w, map[string]interface{}{})
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	ctx := context.Background()

	// expandInstances maps the snake_case consumer_group_id plan field into
	// the camelCase API body key, including every write-exclusive field.
	in := []instanceModel{{
		ConsumerGroupId: types.StringValue("cg-1"),
		InstanceId:      types.StringValue("i-1"),
		InstanceRole:    types.StringValue("ACTIVE"),
		InstanceType:    types.StringValue("ALIYUN_ROCKETMQ"),
		SecurityGroupId: types.StringValue("sg-1"),
		AuthType:        types.StringValue("ACL_AUTH"),
		NetworkType:     types.StringValue("TCP_VPC"),
		VpcId:           types.StringValue("vpc-1"),
		VswitchId:       types.StringValue("vsw-1"),
		RegionId:        types.StringValue("cn-beijing"),
		EndpointUrl:     types.StringValue("rmq.example.com"),
		Password:        types.StringValue("p@ss"),
		Username:        types.StringValue("u-1"),
		MessageProperty: types.ListNull(messagePropertyElementType()),
	}}
	listVal, diags := types.ListValueFrom(ctx, instancesElementType(), in)
	if diags.HasError() {
		t.Fatalf("build instances list: %v", diags)
	}
	expanded, err := expandInstances(ctx, listVal)
	assert.NoError(t, err)
	assert.Equal(t, "cg-1", expanded[0]["consumerGroupId"])
	assert.Equal(t, "vpc-1", expanded[0]["vpcId"])
	assert.Equal(t, "p@ss", expanded[0]["password"])

	// Real ROA POST round-trip: the signed create request reaches the mock
	// server, which records the body and confirms consumerGroupId is present.
	body := map[string]interface{}{
		"planName":  "tf-plan",
		"planType":  "FILE",
		"instances": expanded,
	}
	resp, err := client.RoaPost(rocketmqProduct, rocketmqVersion, "/disaster_recovery", map[string]*string{}, nil, body, true)
	assert.NoError(t, err)
	assert.Equal(t, "plan-900", fmt.Sprint(resp["data"]))
	assert.Contains(t, createBodyRaw, "consumerGroupId", "create body must carry the restored consumerGroupId parameter")
	if instances, ok := createBody["instances"].([]interface{}); ok && len(instances) > 0 {
		if item, ok := instances[0].(map[string]interface{}); ok {
			assert.Equal(t, "cg-1", item["consumerGroupId"])
		}
	}

	// Real ROA GET round-trip: describe returns the plan; the API omits the
	// write-exclusive fields that Read must carry over from prior state.
	obj, err := describeDisasterRecoveryPlan(client, "plan-900")
	assert.NoError(t, err)
	assert.Equal(t, "RUNNING", fmt.Sprint(obj["planStatus"]))
	rawInstances, _ := obj["instances"].([]interface{})
	rawInst, _ := rawInstances[0].(map[string]interface{})
	assert.Equal(t, "cg-1", rawInst["consumerGroupId"])
	assert.Empty(t, rawInst["password"], "password is write-exclusive and must not be returned")
	assert.Empty(t, rawInst["vpcId"], "vpcId is write-exclusive and must not be returned")

	// Carry-over: flattenInstances restores the write-exclusive fields from
	// the prior state element, so a refresh does not drift.
	prior := []instanceModel{{
		Password:    types.StringValue("p@ss"),
		VpcId:       types.StringValue("vpc-1"),
		EndpointUrl: types.StringValue("rmq.example.com"),
		NetworkType: types.StringValue("TCP_VPC"),
		VswitchId:   types.StringValue("vsw-1"),
	}}
	flattened := flattenInstances(obj["instances"], prior)
	assert.Len(t, flattened, 1)
	assert.Equal(t, "p@ss", flattened[0].Password.ValueString())
	assert.Equal(t, "vpc-1", flattened[0].VpcId.ValueString())
	assert.Equal(t, "rmq.example.com", flattened[0].EndpointUrl.ValueString())
	assert.Equal(t, "cg-1", flattened[0].ConsumerGroupId.ValueString())

	// Real ROA DELETE round-trip + NotFound settle: after delete, describe
	// 404s and waitForPlanGone returns without error.
	_, err = client.RoaDelete(rocketmqProduct, rocketmqVersion, "/disaster_recovery/plan-900", map[string]*string{}, nil, nil, true)
	assert.NoError(t, err)
	assert.NoError(t, waitForPlanGone(client, "plan-900", 2*time.Second, 100*time.Millisecond))
}

func listPlan(planId, name string) map[string]interface{} {
	return map[string]interface{}{
		"planId":                planId,
		"planName":              name,
		"planType":              "FILE",
		"planStatus":            "RUNNING",
		"autoSyncCheckpoint":    true,
		"syncCheckpointEnabled": false,
		"createTime":            "2026-10-09T00:00:00Z",
		"updateTime":            "2026-10-09T00:00:00Z",
		"instances":             []interface{}{},
	}
}

// TestUnitRocketmqDisasterRecoveryPlansDataSourceMock drives the real ROA GET
// list round-trip with pagination, then exercises the client-side name_regex
// and ids filters and the data source id hash.
func TestUnitRocketmqDisasterRecoveryPlansDataSourceMock(t *testing.T) {
	client := rocketmqTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("pageNumber")
		if page == "" || page == "1" {
			writeJSON(t, w, map[string]interface{}{"data": map[string]interface{}{
				"list":       []interface{}{listPlan("plan-a", "alpha-1"), listPlan("plan-b", "beta-2")},
				"totalCount": float64(3),
			}})
			return
		}
		writeJSON(t, w, map[string]interface{}{"data": map[string]interface{}{
			"list":       []interface{}{listPlan("plan-c", "alpha-3")},
			"totalCount": float64(3),
		}})
	})
	ctx := context.Background()

	objects, err := listDisasterRecoveryPlans(client, "")
	assert.NoError(t, err)
	assert.Len(t, objects, 3, "pagination must concatenate both pages")

	var resp datasource.ReadResponse

	// name_regex filters client-side: alpha-* matches plan-a and plan-c.
	re := regexp.MustCompile(`^alpha`)
	items, _ := filterPlans(ctx, &resp, "cn-beijing", objects, re, map[string]struct{}{})
	assert.Len(t, items, 2)

	// ids filter narrows to the requested plan id.
	items2, ids := filterPlans(ctx, &resp, "cn-beijing", objects, nil, map[string]struct{}{"plan-b": {}})
	assert.Len(t, items2, 1)
	assert.Equal(t, []string{"plan-b"}, ids)

	// data source id hash is deterministic.
	h := dataResourceIdHash(ids)
	assert.NotEmpty(t, h)
	assert.Equal(t, h, dataResourceIdHash(ids))
	assert.NotEqual(t, h, dataResourceIdHash([]string{"plan-c"}))
}
