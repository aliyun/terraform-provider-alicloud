// Node group matching, pagination and the update path are not reachable through a
// public entry point: the resource exposes no lookup helper and DescribeNodeGroups
// needs a live instance, so these cases live in an internal unit test file with an
// HTTP fixture standing in for the StarRocks API.
package starrocks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func starRocksTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build credential: %s", err)
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
		t.Fatalf("build client: %s", err)
	}
	endpoints.Store("starrocks", endpoint)
	return client
}

func unitTestBackendNodeGroup(id, zone string) map[string]interface{} {
	return map[string]interface{}{
		"NodeGroupId":              id,
		"ComponentType":            "CN",
		"ZoneId":                   zone,
		"Cu":                       float64(8),
		"DiskNumber":               float64(1),
		"ResidentNodeNumber":       float64(3),
		"SpecType":                 "standard",
		"StoragePerformanceLevel":  "pl1",
		"LocalStorageInstanceType": "",
		"StorageSize":              float64(200),
		"Status":                   "RUNNING",
	}
}

func TestUnitMatchBackendNodeGroup(t *testing.T) {
	config := map[string]interface{}{
		"cu":                          8,
		"disk_number":                 1,
		"resident_node_number":        3,
		"zone_id":                     "cn-hangzhou-i",
		"spec_type":                   "standard",
		"storage_performance_level":   "pl1",
		"local_storage_instance_type": "null",
		"storage_size":                300,
	}

	t.Run("resolves the single matching group", func(t *testing.T) {
		groups := []map[string]interface{}{
			unitTestBackendNodeGroup("ng-unrelated-zone", "cn-hangzhou-h"),
			unitTestBackendNodeGroup("ng-target", "cn-hangzhou-i"),
		}
		id, err := matchBackendNodeGroup(groups, config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != "ng-target" {
			t.Fatalf("expected ng-target, got %q", id)
		}
	})

	t.Run("matches by the remaining attributes when the zone is not reported", func(t *testing.T) {
		groups := []map[string]interface{}{
			unitTestBackendNodeGroup("ng-target", ""),
		}
		id, err := matchBackendNodeGroup(groups, config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != "ng-target" {
			t.Fatalf("expected ng-target, got %q", id)
		}
	})

	t.Run("fails when no group matches", func(t *testing.T) {
		groups := []map[string]interface{}{
			unitTestBackendNodeGroup("ng-other", ""),
		}
		groups[0]["Cu"] = float64(16)
		_, err := matchBackendNodeGroup(groups, config)
		if err == nil || !strings.Contains(err.Error(), "no backend node group matches") {
			t.Fatalf("expected a no-match error, got %v", err)
		}
	})

	t.Run("fails instead of guessing between identical groups", func(t *testing.T) {
		groups := []map[string]interface{}{
			unitTestBackendNodeGroup("ng-1", "cn-hangzhou-i"),
			unitTestBackendNodeGroup("ng-2", "cn-hangzhou-i"),
		}
		id, err := matchBackendNodeGroup(groups, config)
		if err == nil {
			t.Fatalf("expected an ambiguity error, got id %q", id)
		}
		if !strings.Contains(err.Error(), "refusing to guess") {
			t.Fatalf("expected an ambiguity error, got %v", err)
		}
	})

	t.Run("fails when two groups only differ by an unreported zone", func(t *testing.T) {
		groups := []map[string]interface{}{
			unitTestBackendNodeGroup("ng-1", ""),
			unitTestBackendNodeGroup("ng-2", ""),
		}
		id, err := matchBackendNodeGroup(groups, config)
		if err == nil {
			t.Fatalf("expected an ambiguity error, got id %q", id)
		}
		if !strings.Contains(err.Error(), "refusing to guess") {
			t.Fatalf("expected an ambiguity error, got %v", err)
		}
	})
}

func TestUnitBackendNodeGroupAttributesMatch(t *testing.T) {
	config := map[string]interface{}{
		"cu":                          8,
		"disk_number":                 1,
		"resident_node_number":        3,
		"zone_id":                     "cn-hangzhou-i",
		"spec_type":                   "standard",
		"storage_performance_level":   "pl1",
		"local_storage_instance_type": "null",
		"storage_size":                300,
	}

	group := unitTestBackendNodeGroup("ng-target", "cn-hangzhou-i")
	if !backendNodeGroupAttributesMatch(group, config, true) {
		t.Fatal("expected the group to match the configured attributes")
	}
	if !backendNodeGroupAttributesMatch(group, config, false) {
		t.Fatal("expected the group to match when the zone is not checked")
	}

	group["Cu"] = float64(16)
	if backendNodeGroupAttributesMatch(group, config, true) || backendNodeGroupAttributesMatch(group, config, false) {
		t.Fatal("expected a different cu to make the group not match")
	}
	group["Cu"] = float64(8)

	group["StorageSize"] = float64(400)
	if !backendNodeGroupAttributesMatch(group, config, true) {
		t.Fatal("expected storage_size to be ignored when matching")
	}
}

func TestUnitBackendNodeGroupStorageSizeChangeForcesReplacement(t *testing.T) {
	cases := []struct {
		name  string
		old   int
		oldOk bool
		new   int
		newOk bool
		want  bool
	}{
		{name: "growth is in place", old: 200, oldOk: true, new: 300, newOk: true, want: false},
		{name: "same size is in place", old: 200, oldOk: true, new: 200, newOk: true, want: false},
		{name: "shrink replaces", old: 300, oldOk: true, new: 200, newOk: true, want: true},
		{name: "removing the value replaces", old: 300, oldOk: true, new: 0, newOk: false, want: true},
		{name: "setting a previously unset value replaces", old: 0, oldOk: false, new: 300, newOk: true, want: true},
		{name: "unset on both sides never replaces", old: 0, oldOk: false, new: 0, newOk: false, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := backendNodeGroupStorageSizeChangeForcesReplacement(tc.old, tc.oldOk, tc.new, tc.newOk)
			if got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

// starRocksNodeGroupsHandler serves the node group and instance calls of one update
// and records every action it sees.
type starRocksNodeGroupsHandler struct {
	t      *testing.T
	groups []map[string]interface{}
	total  int
	// modifyResponses are consumed in order; an empty entry means success.
	modifyResponses []int
	modifyCalls     int
	actions         []string
	modifyQueries   []map[string]string
	instanceStatus  string
}

func (h *starRocksNodeGroupsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	action := r.URL.Path
	h.actions = append(h.actions, action)
	switch action {
	case "/webapi/nodegroup/describeNodeGroups":
		body := map[string]interface{}{}
		if raw, err := io.ReadAll(r.Body); err == nil && len(raw) > 0 {
			_ = json.Unmarshal(raw, &body)
		}
		if _, byId := body["nodeGroupIds"]; byId {
			item := h.groupById(fmt.Sprint(body["nodeGroupIds"].([]interface{})[0]))
			item["Status"] = "RUNNING"
			if h.instanceStatus != "" {
				item["Status"] = h.instanceStatus
			}
			writeJSON(w, map[string]interface{}{"Data": []interface{}{item}})
			return
		}
		pageNumber := 1
		pageSize := 100
		fmt.Sscan(query.Get("PageNumber"), &pageNumber)
		fmt.Sscan(query.Get("PageSize"), &pageSize)
		componentType := fmt.Sprint(body["componentType"])
		filtered := make([]map[string]interface{}, 0, len(h.groups))
		for _, group := range h.groups {
			if componentType != "" && fmt.Sprint(group["ComponentType"]) != componentType {
				continue
			}
			filtered = append(filtered, group)
		}
		start := (pageNumber - 1) * pageSize
		end := start + pageSize
		if start > len(filtered) {
			start = len(filtered)
		}
		if end > len(filtered) {
			end = len(filtered)
		}
		page := make([]interface{}, 0, end-start)
		for _, group := range filtered[start:end] {
			page = append(page, group)
		}
		total := h.total
		if total > len(filtered) {
			total = len(filtered)
		}
		writeJSON(w, map[string]interface{}{"Data": page, "Total": total})
	case "/webapi/resourceChange/modifyDiskSize":
		h.modifyCalls++
		h.modifyQueries = append(h.modifyQueries, map[string]string{
			"InstanceId":  query.Get("InstanceId"),
			"NodeGroupId": query.Get("NodeGroupId"),
			"Target":      query.Get("Target"),
		})
		// Apply the resize before reporting an injected failure: the lost response case
		// is exactly what the retry has to detect by re-reading the group.
		target := 0
		fmt.Sscan(query.Get("Target"), &target)
		for _, group := range h.groups {
			if fmt.Sprint(group["NodeGroupId"]) == query.Get("NodeGroupId") {
				group["StorageSize"] = float64(target)
			}
		}
		status := http.StatusOK
		if len(h.modifyResponses) > 0 {
			status = h.modifyResponses[0]
			h.modifyResponses = h.modifyResponses[1:]
		}
		if status != http.StatusOK {
			w.WriteHeader(status)
			writeJSON(w, map[string]interface{}{"Code": "InternalError", "Message": "code: 500, injected failure", "RequestId": "rid"})
			return
		}
		writeJSON(w, map[string]interface{}{"Success": true})
	case "/webapi/starrocks/describeInstances":
		writeJSON(w, map[string]interface{}{"Data": []interface{}{map[string]interface{}{
			"InstanceId":     "c-fixture",
			"InstanceStatus": "running",
			"InstanceName":   "tf-testacc-starrocks",
			"Tags":           []interface{}{},
			"VSwitches":      []interface{}{},
		}}})
	default:
		writeJSON(w, map[string]interface{}{"Success": true})
	}
}

func (h *starRocksNodeGroupsHandler) groupById(id string) map[string]interface{} {
	for _, group := range h.groups {
		if fmt.Sprint(group["NodeGroupId"]) == id {
			return group
		}
	}
	return unitTestBackendNodeGroup(id, "")
}

func (h *starRocksNodeGroupsHandler) countAction(action string) int {
	count := 0
	for _, item := range h.actions {
		if item == action {
			count++
		}
	}
	return count
}

func writeJSON(w http.ResponseWriter, body map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

func backendGroupStateAttributes(size string) map[string]string {
	return map[string]string{
		"backend_node_groups.#":                             "1",
		"backend_node_groups.0.cu":                          "8",
		"backend_node_groups.0.zone_id":                     "cn-hangzhou-i",
		"backend_node_groups.0.spec_type":                   "standard",
		"backend_node_groups.0.disk_number":                 "1",
		"backend_node_groups.0.storage_performance_level":   "pl1",
		"backend_node_groups.0.resident_node_number":        "3",
		"backend_node_groups.0.storage_size":                size,
		"backend_node_groups.0.local_storage_instance_type": "null",
	}
}

func starRocksInstanceData(t *testing.T, attributes map[string]string, changes map[string]*terraform.ResourceAttrDiff) *schema.ResourceData {
	t.Helper()
	drain := map[string]string{"instance_name": "tf-testacc-starrocks"}
	for key, value := range attributes {
		drain[key] = value
	}
	state := &terraform.InstanceState{ID: "c-fixture", Attributes: drain}
	diff := terraform.NewInstanceDiff()
	if diff.Attributes == nil {
		diff.Attributes = make(map[string]*terraform.ResourceAttrDiff)
	}
	for key, change := range changes {
		diff.Attributes[key] = change
	}
	data, err := schema.InternalMap(ResourceInstance().Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("build resource data: %s", err)
	}
	return data
}

func TestUnitDescribeNodeGroupsPagination(t *testing.T) {
	groups := make([]map[string]interface{}, 0, 150)
	for i := 0; i < 149; i++ {
		groups = append(groups, unitTestBackendNodeGroup(fmt.Sprintf("ng-filler-%03d", i), "cn-hangzhou-h"))
		groups[i]["Cu"] = float64(4)
	}
	groups = append(groups, unitTestBackendNodeGroup("ng-last-page", "cn-hangzhou-i"))
	handler := &starRocksNodeGroupsHandler{t: t, groups: groups, total: len(groups)}
	client := starRocksTestClient(t, handler.ServeHTTP)

	service := NewService(client)
	listed, err := service.DescribeNodeGroupsByComponentType("c-fixture", "CN", 10*time.Second)
	if err != nil {
		t.Fatalf("list node groups: %s", err)
	}
	if len(listed) != len(groups) {
		t.Fatalf("expected %d groups across all pages, got %d", len(groups), len(listed))
	}
	if handler.countAction("/webapi/nodegroup/describeNodeGroups") < 2 {
		t.Fatal("expected the listing to read a second page")
	}
	id, err := matchBackendNodeGroup(listed, map[string]interface{}{
		"cu": 8, "disk_number": 1, "resident_node_number": 3,
		"zone_id": "cn-hangzhou-i", "spec_type": "standard", "storage_performance_level": "pl1",
	})
	if err != nil {
		t.Fatalf("match group on the later page: %s", err)
	}
	if id != "ng-last-page" {
		t.Fatalf("expected ng-last-page, got %q", id)
	}
}

func TestUnitBackendNodeGroupResizeRequest(t *testing.T) {
	handler := &starRocksNodeGroupsHandler{
		t:              t,
		groups:         []map[string]interface{}{unitTestBackendNodeGroup("ng-target", "cn-hangzhou-i")},
		total:          1,
		instanceStatus: "RUNNING",
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	data := starRocksInstanceData(t, backendGroupStateAttributes("200"), map[string]*terraform.ResourceAttrDiff{
		"backend_node_groups.0.storage_size": {Old: "200", New: "300"},
	})

	if err := resourceInstanceUpdate(data, client); err != nil {
		t.Fatalf("update: %s", err)
	}
	if handler.modifyCalls != 1 {
		t.Fatalf("expected exactly one ModifyDiskSize call, got %d", handler.modifyCalls)
	}
	sent := handler.modifyQueries[0]
	if sent["InstanceId"] != "c-fixture" || sent["NodeGroupId"] != "ng-target" || sent["Target"] != "300" {
		t.Fatalf("unexpected ModifyDiskSize request: %#v", sent)
	}
	if size := fmt.Sprint(handler.groups[0]["StorageSize"]); size != "300" {
		t.Fatalf("expected the remote group to reach 300 GiB, got %s", size)
	}
	if handler.countAction("/webapi/cluster/update_name") != 0 {
		t.Fatal("the name must not be updated when only storage_size changes")
	}
}

func TestUnitBackendNodeGroupAmbiguityWritesNothing(t *testing.T) {
	handler := &starRocksNodeGroupsHandler{
		t: t,
		groups: []map[string]interface{}{
			unitTestBackendNodeGroup("ng-1", "cn-hangzhou-i"),
			unitTestBackendNodeGroup("ng-2", "cn-hangzhou-i"),
		},
		total:          2,
		instanceStatus: "RUNNING",
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	data := starRocksInstanceData(t, backendGroupStateAttributes("200"), map[string]*terraform.ResourceAttrDiff{
		"backend_node_groups.0.storage_size": {Old: "200", New: "300"},
		"instance_name":                      {Old: "tf-testacc-starrocks", New: "tf-testacc-starrocks-renamed"},
	})

	err := resourceInstanceUpdate(data, client)
	if err == nil || !strings.Contains(err.Error(), "refusing to guess") {
		t.Fatalf("expected an ambiguity error, got %v", err)
	}
	if handler.modifyCalls != 0 {
		t.Fatalf("no ModifyDiskSize call may be sent on ambiguity, got %d", handler.modifyCalls)
	}
	if handler.countAction("/webapi/cluster/update_name") != 0 {
		t.Fatal("the instance name must not change when the node group match is ambiguous")
	}
	if handler.countAction("/webapi/resourceGroup/change") != 0 {
		t.Fatal("the resource group must not change when the node group match is ambiguous")
	}
}

func TestUnitBackendNodeGroupRetryRechecksRemoteState(t *testing.T) {
	handler := &starRocksNodeGroupsHandler{
		t:               t,
		groups:          []map[string]interface{}{unitTestBackendNodeGroup("ng-target", "cn-hangzhou-i")},
		total:           1,
		modifyResponses: []int{http.StatusInternalServerError},
		instanceStatus:  "RUNNING",
	}
	client := starRocksTestClient(t, handler.ServeHTTP)
	data := starRocksInstanceData(t, backendGroupStateAttributes("200"), map[string]*terraform.ResourceAttrDiff{
		"backend_node_groups.0.storage_size": {Old: "200", New: "300"},
	})

	// The fixture applies the resize before returning the injected 500, which is the
	// lost response case: the retry must re-read the group and skip the second write.
	if err := resourceInstanceUpdate(data, client); err != nil {
		t.Fatalf("update: %s", err)
	}
	if handler.modifyCalls != 1 {
		t.Fatalf("expected the retry to be skipped after the state reached the target, got %d writes", handler.modifyCalls)
	}
}

func backendGroupConfig(storageSize interface{}, overrides map[string]interface{}) map[string]interface{} {
	group := map[string]interface{}{
		"cu":                          8,
		"zone_id":                     "cn-hangzhou-i",
		"spec_type":                   "standard",
		"disk_number":                 1,
		"storage_performance_level":   "pl1",
		"resident_node_number":        3,
		"local_storage_instance_type": "null",
	}
	if storageSize != nil {
		group["storage_size"] = storageSize
	}
	for key, value := range overrides {
		group[key] = value
	}
	return map[string]interface{}{
		"admin_password":      "1qaz@QAZ",
		"cluster_zone_id":     "cn-hangzhou-i",
		"instance_name":       "tf-testacc-starrocks",
		"backend_node_groups": []interface{}{group},
	}
}

func TestUnitResourceInstanceBackendNodeGroupPlan(t *testing.T) {
	resource := ResourceInstance()
	state := func(cfg map[string]interface{}) *terraform.InstanceState {
		data := schema.TestResourceDataRaw(t, resource.Schema, cfg)
		data.SetId("c-fixture")
		return data.State()
	}
	diffRequiresNew := func(t *testing.T, oldCfg, newCfg map[string]interface{}) bool {
		t.Helper()
		diff, err := resource.Diff(context.Background(), state(oldCfg), terraform.NewResourceConfigRaw(newCfg), nil)
		if err != nil {
			t.Fatalf("diff: %s", err)
		}
		if diff == nil {
			return false
		}
		return diff.RequiresNew()
	}

	grow := backendGroupConfig(300, nil)
	shrink := backendGroupConfig(200, nil)

	t.Run("growing is in place", func(t *testing.T) {
		if diffRequiresNew(t, backendGroupConfig(200, nil), grow) {
			t.Fatal("growing storage_size must stay an in-place update")
		}
	})
	t.Run("shrinking replaces", func(t *testing.T) {
		if !diffRequiresNew(t, backendGroupConfig(300, nil), shrink) {
			t.Fatal("shrinking storage_size must replace the instance")
		}
	})
	t.Run("removing replaces", func(t *testing.T) {
		if !diffRequiresNew(t, backendGroupConfig(300, nil), backendGroupConfig(nil, nil)) {
			t.Fatal("removing storage_size must replace the instance")
		}
	})
	t.Run("setting a previously unset value replaces", func(t *testing.T) {
		if !diffRequiresNew(t, backendGroupConfig(nil, nil), backendGroupConfig(300, nil)) {
			t.Fatal("setting storage_size without a previous value must replace the instance")
		}
	})
	for _, field := range []struct {
		name  string
		value interface{}
	}{
		{"cu", 16},
		{"spec_type", "localSSD"},
		{"disk_number", 2},
		{"resident_node_number", 6},
		{"storage_performance_level", "pl2"},
		{"zone_id", "cn-hangzhou-h"},
		{"local_storage_instance_type", "local_ssd"},
	} {
		field := field
		t.Run("changing "+field.name+" replaces", func(t *testing.T) {
			changed := backendGroupConfig(200, map[string]interface{}{field.name: field.value})
			if !diffRequiresNew(t, backendGroupConfig(200, nil), changed) {
				t.Fatalf("changing %s must replace the instance", field.name)
			}
		})
	}
	t.Run("adding a backend group replaces", func(t *testing.T) {
		second := backendGroupConfig(200, map[string]interface{}{"zone_id": "cn-hangzhou-h"})["backend_node_groups"].([]interface{})[0]
		cfg := backendGroupConfig(200, nil)
		cfg["backend_node_groups"] = []interface{}{cfg["backend_node_groups"].([]interface{})[0], second}
		if !diffRequiresNew(t, backendGroupConfig(200, nil), cfg) {
			t.Fatal("adding a backend group must replace the instance")
		}
	})
	t.Run("removing a backend group replaces", func(t *testing.T) {
		second := backendGroupConfig(200, map[string]interface{}{"zone_id": "cn-hangzhou-h"})["backend_node_groups"].([]interface{})[0]
		oldCfg := backendGroupConfig(200, nil)
		oldCfg["backend_node_groups"] = []interface{}{oldCfg["backend_node_groups"].([]interface{})[0], second}
		if !diffRequiresNew(t, oldCfg, backendGroupConfig(200, nil)) {
			t.Fatal("removing a backend group must replace the instance")
		}
	})
	t.Run("changing frontend groups replaces", func(t *testing.T) {
		oldCfg := backendGroupConfig(200, nil)
		oldCfg["frontend_node_groups"] = []interface{}{map[string]interface{}{
			"cu": 8, "zone_id": "cn-hangzhou-i", "spec_type": "standard", "disk_number": 1,
			"storage_performance_level": "pl1", "resident_node_number": 3, "storage_size": 100,
		}}
		newCfg := backendGroupConfig(200, nil)
		newCfg["frontend_node_groups"] = []interface{}{map[string]interface{}{
			"cu": 16, "zone_id": "cn-hangzhou-i", "spec_type": "standard", "disk_number": 1,
			"storage_performance_level": "pl1", "resident_node_number": 3, "storage_size": 100,
		}}
		if !diffRequiresNew(t, oldCfg, newCfg) {
			t.Fatal("changing frontend_node_groups must replace the instance")
		}
	})
}
