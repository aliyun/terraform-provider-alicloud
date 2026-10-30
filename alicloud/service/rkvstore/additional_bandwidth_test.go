// Package rkvstore internal tests for the additional_bandwidth helpers that
// moved from the root alicloud package into package rkvstore. The unit tests
// exercise unexported prod helpers (buildAdditionalBandwidthState,
// collectAdditionalBandwidthNodes, refreshAdditionalBandwidthState,
// buildAdditionalBandwidthRequest) and the kvstoreAdditionalBandwidth type, so
// they live in an internal test file (package rkvstore), mirroring
// validators_local_test.go. Behavior preserved verbatim from PR #11409.
package rkvstore

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitKvstoreAdditionalBandwidthState(t *testing.T) {
	node := func(insName, role string, defaultBandWidth, currentBandWidth int64) map[string]interface{} {
		return map[string]interface{}{
			"CustinsId":        insName + "-shard",
			"InsName":          insName,
			"Role":             role,
			"DefaultBandWidth": defaultBandWidth,
			"CurrentBandWidth": currentBandWidth,
		}
	}
	twoShards := []map[string]interface{}{
		node("r-x-db-0", "MASTER", 96, 116),
		node("r-x-db-1", "MASTER", 96, 126),
	}

	state, ok := buildAdditionalBandwidthState(nil, nil)
	assert.False(t, ok)
	assert.Nil(t, state)

	// An import (no managed shards) reports every shard with purchased bandwidth, keyed by name.
	state, ok = buildAdditionalBandwidthState(twoShards, nil)
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{
		{"node_id": "r-x-db-0", "bandwidth": 20},
		{"node_id": "r-x-db-1", "bandwidth": 30},
	}, state)

	// A shard without purchased bandwidth is not reported by an import.
	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{
		node("r-x-db-0", "MASTER", 96, 116),
		node("r-x-db-1", "MASTER", 96, 96),
	}, nil)
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{{"node_id": "r-x-db-0", "bandwidth": 20}}, state)

	// Managed shards are reported with their own value and keep the configured node_id, so a
	// shard that the configuration does not manage does not appear.
	state, ok = buildAdditionalBandwidthState(twoShards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-1", Bandwidth: 0},
		{NodeId: "r-x-db-0#100", Bandwidth: 0},
	})
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{
		{"node_id": "db-1", "bandwidth": 30},
		{"node_id": "r-x-db-0#100", "bandwidth": 20},
	}, state)

	// A shard reports a master and a slave row; the master row wins, and read-only nodes
	// (InsType 3) are not bandwidth targets.
	readOnly := node("r-x-db-1", "MASTER", 96, 200)
	readOnly["InsType"] = 3
	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{
		node("r-x-db-0", "SLAVE", 96, 140),
		node("r-x-db-0", "MASTER", 96, 116),
		readOnly,
	}, nil)
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{{"node_id": "r-x-db-0", "bandwidth": 20}}, state)

	// Bandwidth below the specification default is reported as no additional bandwidth.
	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{
		node("r-x-db-0", "MASTER", 96, 90),
	}, nil)
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{}, state)

	// A managed shard without a bandwidth row means the listing is incomplete: the caller keeps
	// the previous state instead of reporting a smaller value.
	state, ok = buildAdditionalBandwidthState(twoShards[:1], []kvstoreAdditionalBandwidth{
		{NodeId: "db-1", Bandwidth: 30},
	})
	assert.False(t, ok)
	assert.Nil(t, state)

	// A managed node_id that no longer maps to exactly one shard cannot be refreshed.
	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{
		node("r-x-db-0", "MASTER", 96, 116),
		node("r-y-db-0", "MASTER", 96, 116),
	}, []kvstoreAdditionalBandwidth{{NodeId: "db-0", Bandwidth: 20}})
	assert.False(t, ok)
	assert.Nil(t, state)

	// An unambiguous suffix is resolved to its shard.
	state, ok = buildAdditionalBandwidthState(twoShards, []kvstoreAdditionalBandwidth{{NodeId: "db-0", Bandwidth: 20}})
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{{"node_id": "db-0", "bandwidth": 20}}, state)

	// A standard instance configures the bandwidth without a node_id; its single data shard is
	// the target, including for a refreshed value.
	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{node("r-x", "MASTER", 96, 116)}, []kvstoreAdditionalBandwidth{{Bandwidth: 20}})
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{{"node_id": "", "bandwidth": 20}}, state)

	state, ok = buildAdditionalBandwidthState([]map[string]interface{}{node("r-x", "MASTER", 96, 136)}, []kvstoreAdditionalBandwidth{{Bandwidth: 20}})
	assert.True(t, ok)
	assert.Equal(t, []map[string]interface{}{{"node_id": "", "bandwidth": 40}}, state)

	// A node_id-less block cannot be refreshed when the listing has several data shards.
	state, ok = buildAdditionalBandwidthState(twoShards, []kvstoreAdditionalBandwidth{{Bandwidth: 20}})
	assert.False(t, ok)
	assert.Nil(t, state)
}

func TestUnitKvstoreAdditionalBandwidthStandardRefresh(t *testing.T) {
	node := func(insName string, defaultBandWidth, currentBandWidth int64) map[string]interface{} {
		return map[string]interface{}{
			"InsName":          insName,
			"Role":             "MASTER",
			"DefaultBandWidth": defaultBandWidth,
			"CurrentBandWidth": currentBandWidth,
		}
	}
	stateOf := func(t *testing.T, d *schema.ResourceData) map[string]int {
		value, ok := d.Get("additional_bandwidth").(*schema.Set)
		assert.True(t, ok)
		blocks := make(map[string]int)
		for _, item := range value.List() {
			entry := item.(map[string]interface{})
			blocks[fmt.Sprint(entry["node_id"])] = entry["bandwidth"].(int)
		}
		return blocks
	}
	newResourceData := func(raw map[string]interface{}) *schema.ResourceData {
		d := schema.TestResourceDataRaw(t, ResourceAliCloudKvstoreInstance().Schema, raw)
		d.SetId("r-x")
		return d
	}

	// A standard instance configures the bandwidth without a node_id. The refresh keeps that block
	// instead of treating the instance as an import, so the state still matches the configuration.
	d := newResourceData(map[string]interface{}{
		"additional_bandwidth": []interface{}{map[string]interface{}{"bandwidth": 20}},
	})
	assert.True(t, refreshAdditionalBandwidthState(d, []map[string]interface{}{node("r-x", 96, 116)}))
	assert.Equal(t, map[string]int{"": 20}, stateOf(t, d))

	// A second refresh of the same listing produces the same state, so the next plan is empty.
	assert.True(t, refreshAdditionalBandwidthState(d, []map[string]interface{}{node("r-x", 96, 116)}))
	assert.Equal(t, map[string]int{"": 20}, stateOf(t, d))

	// A bandwidth changed outside Terraform reaches the state, so the drift shows up in the plan.
	assert.True(t, refreshAdditionalBandwidthState(d, []map[string]interface{}{node("r-x", 96, 136)}))
	assert.Equal(t, map[string]int{"": 40}, stateOf(t, d))

	// An explicit zero releases the purchased bandwidth but keeps the block, so the release stays
	// in the state and does not turn into a pending creation.
	released := newResourceData(map[string]interface{}{
		"additional_bandwidth": []interface{}{map[string]interface{}{"bandwidth": 0}},
	})
	assert.True(t, refreshAdditionalBandwidthState(released, []map[string]interface{}{node("r-x", 96, 96)}))
	assert.Equal(t, map[string]int{"": 0}, stateOf(t, released))

	// Only a state without a managed block (an import) rebuilds the blocks from the node names.
	imported := newResourceData(map[string]interface{}{})
	assert.True(t, refreshAdditionalBandwidthState(imported, []map[string]interface{}{node("r-x", 96, 116)}))
	assert.Equal(t, map[string]int{"r-x": 20}, stateOf(t, imported))

	// A node_id-less block cannot be matched to a listing with several data shards, so the refresh
	// keeps the previous state rather than writing another shard's value.
	ambiguous := newResourceData(map[string]interface{}{
		"additional_bandwidth": []interface{}{map[string]interface{}{"bandwidth": 20}},
	})
	assert.False(t, refreshAdditionalBandwidthState(ambiguous, []map[string]interface{}{
		node("r-x-db-0", 96, 116),
		node("r-x-db-1", 96, 116),
	}))
	assert.Equal(t, map[string]int{"": 20}, stateOf(t, ambiguous))
}

func TestUnitKvstoreAdditionalBandwidthRequest(t *testing.T) {
	shards := []string{"r-x-db-0", "r-x-db-1", "r-x-db-10"}

	nodeIds, bandwidths, err := buildAdditionalBandwidthRequest("cluster", 3, shards, nil)
	assert.NoError(t, err)
	assert.Equal(t, "", nodeIds)
	assert.Equal(t, "", bandwidths)

	// Every bandwidth value is emitted next to its own shard, and the vectors are ordered by
	// shard name so a plan produces a stable request.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "r-x-db-1", Bandwidth: 30},
		{NodeId: "r-x-db-0", Bandwidth: 20},
		{NodeId: "r-x-db-10", Bandwidth: 40},
	})
	assert.NoError(t, err)
	assert.Equal(t, "r-x-db-0,r-x-db-1,r-x-db-10", nodeIds)
	assert.Equal(t, "20,30,40", bandwidths)

	// A suffix addresses one shard without repeating the instance ID.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-10", Bandwidth: 50},
	})
	assert.NoError(t, err)
	assert.Equal(t, "r-x-db-10", nodeIds)
	assert.Equal(t, "50", bandwidths)

	// The topology form is accepted, but only the shard node name is sent.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "r-x-db-0#123456", Bandwidth: 20},
	})
	assert.NoError(t, err)
	assert.Equal(t, "r-x-db-0", nodeIds)
	assert.Equal(t, "20", bandwidths)

	// Releasing the purchased bandwidth is an explicit zero.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-0", Bandwidth: 0},
	})
	assert.NoError(t, err)
	assert.Equal(t, "r-x-db-0", nodeIds)
	assert.Equal(t, "0", bandwidths)

	// A standard instance takes Bandwidth alone; the API does not accept NodeId for it.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("standard", 1, nil, []kvstoreAdditionalBandwidth{
		{Bandwidth: 20},
	})
	assert.NoError(t, err)
	assert.Equal(t, "", nodeIds)
	assert.Equal(t, "20", bandwidths)

	// A read/write splitting architecture identified only by its shard count still requires the
	// shard name, and a standard instance rejects a block that names one.
	nodeIds, bandwidths, err = buildAdditionalBandwidthRequest("", 2, []string{"r-x-db-0", "r-x-db-1"}, []kvstoreAdditionalBandwidth{
		{NodeId: "db-1", Bandwidth: 20},
	})
	assert.NoError(t, err)
	assert.Equal(t, "r-x-db-1", nodeIds)
	assert.Equal(t, "20", bandwidths)

	_, _, err = buildAdditionalBandwidthRequest("standard", 1, nil, []kvstoreAdditionalBandwidth{
		{NodeId: "r-x", Bandwidth: 20},
	})
	assert.Error(t, err)

	// A cluster instance rejects a block without a shard name, and a standard instance accepts a
	// single block only.
	_, _, err = buildAdditionalBandwidthRequest("cluster", 2, shards, []kvstoreAdditionalBandwidth{
		{Bandwidth: 20},
	})
	assert.Error(t, err)

	_, _, err = buildAdditionalBandwidthRequest("standard", 1, nil, []kvstoreAdditionalBandwidth{
		{Bandwidth: 20},
		{Bandwidth: 30},
	})
	assert.Error(t, err)

	// The same shard spelled two ways is rejected instead of being sent twice.
	_, _, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-0", Bandwidth: 20},
		{NodeId: "r-x-db-0#123456", Bandwidth: 30},
	})
	assert.Error(t, err)

	// A node_id that matches no shard is rejected before the API call.
	_, _, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-9", Bandwidth: 20},
	})
	assert.Error(t, err)

	// A suffix shared by two shards is ambiguous and rejected.
	_, _, err = buildAdditionalBandwidthRequest("cluster", 2, []string{"r-x-db-0", "r-y-db-0"}, []kvstoreAdditionalBandwidth{
		{NodeId: "db-0", Bandwidth: 20},
	})
	assert.Error(t, err)

	// A negative bandwidth value is rejected.
	_, _, err = buildAdditionalBandwidthRequest("cluster", 3, shards, []kvstoreAdditionalBandwidth{
		{NodeId: "db-0", Bandwidth: -1},
	})
	assert.Error(t, err)
}

func TestUnitKvstoreAdditionalBandwidthPagination(t *testing.T) {
	pages := map[int][]map[string]interface{}{
		1: {{"InsName": "r-x-db-0"}, {"InsName": "r-x-db-1"}},
		2: {{"InsName": "r-x-db-2"}},
	}
	var requested []int
	nodes, err := collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		requested = append(requested, pageNumber)
		return pages[pageNumber], 3, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, requested)
	assert.Len(t, nodes, 3)

	// A late page is still collected when the earlier pages are full.
	requested = nil
	nodes, err = collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		requested = append(requested, pageNumber)
		if pageNumber == 1 {
			return []map[string]interface{}{{"InsName": "r-x-db-0"}}, 2, nil
		}
		if pageNumber == 2 {
			return []map[string]interface{}{{"InsName": "r-x-db-1"}}, 2, nil
		}
		return nil, 2, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, requested)
	assert.Len(t, nodes, 2)

	// An empty page before TotalCount is reached is an incomplete listing.
	_, err = collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		if pageNumber == 1 {
			return []map[string]interface{}{{"InsName": "r-x-db-0"}}, 3, nil
		}
		return nil, 3, nil
	})
	assert.ErrorIs(t, err, errIncompleteAdditionalBandwidthNodes)

	// Without TotalCount a short page ends the walk.
	requested = nil
	nodes, err = collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		requested = append(requested, pageNumber)
		if pageNumber == 1 {
			return []map[string]interface{}{{"InsName": "r-x-db-0"}}, 0, nil
		}
		return nil, 0, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{1}, requested)
	assert.Len(t, nodes, 1)

	// A page error is surfaced with the nodes collected so far.
	_, err = collectAdditionalBandwidthNodes(func(pageNumber, pageSize int) ([]map[string]interface{}, int, error) {
		if pageNumber == 2 {
			return nil, 0, fmt.Errorf("boom")
		}
		return []map[string]interface{}{{"InsName": "r-x-db-0"}}, 2, nil
	})
	assert.Error(t, err)
}
