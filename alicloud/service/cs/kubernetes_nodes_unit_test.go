package cs

import (
	"fmt"
	"testing"

	roaCS "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea/tea"
)

// fakeKubernetesNodesClient replays canned responses so the node pool
// resolution, pagination and ordering logic runs without a cluster.
type fakeKubernetesNodesClient struct {
	poolResponses []*roaCS.DescribeClusterNodePoolsResponse
	poolErr       error
	poolCalls     int

	nodePages    []*roaCS.DescribeClusterNodesResponse
	nodeErr      error
	nodeCalls    int
	nodeRequests []*roaCS.DescribeClusterNodesRequest
}

func (f *fakeKubernetesNodesClient) DescribeClusterNodePools(clusterId *string, request *roaCS.DescribeClusterNodePoolsRequest) (*roaCS.DescribeClusterNodePoolsResponse, error) {
	if f.poolErr != nil {
		return nil, f.poolErr
	}
	if len(f.poolResponses) == 0 {
		return &roaCS.DescribeClusterNodePoolsResponse{Body: &roaCS.DescribeClusterNodePoolsResponseBody{}}, nil
	}
	index := f.poolCalls
	if index >= len(f.poolResponses) {
		index = len(f.poolResponses) - 1
	}
	f.poolCalls++
	return f.poolResponses[index], nil
}

func (f *fakeKubernetesNodesClient) DescribeClusterNodes(clusterId *string, request *roaCS.DescribeClusterNodesRequest) (*roaCS.DescribeClusterNodesResponse, error) {
	f.nodeRequests = append(f.nodeRequests, request)
	if f.nodeErr != nil {
		return nil, f.nodeErr
	}
	if len(f.nodePages) == 0 {
		return &roaCS.DescribeClusterNodesResponse{Body: &roaCS.DescribeClusterNodesResponseBody{}}, nil
	}
	index := f.nodeCalls
	if index >= len(f.nodePages) {
		index = len(f.nodePages) - 1
	}
	f.nodeCalls++
	return f.nodePages[index], nil
}

func fakeNode(instanceId, name, nodepoolId string) *roaCS.DescribeClusterNodesResponseBodyNodes {
	return &roaCS.DescribeClusterNodesResponseBodyNodes{
		InstanceId:   tea.String(instanceId),
		NodeName:     tea.String(name),
		NodepoolId:   tea.String(nodepoolId),
		HostName:     tea.String(name),
		InstanceRole: tea.String("Worker"),
		State:        tea.String("running"),
	}
}

func fakeNodesPage(totalCount, pageSize int32, nodes ...*roaCS.DescribeClusterNodesResponseBodyNodes) *roaCS.DescribeClusterNodesResponse {
	return &roaCS.DescribeClusterNodesResponse{
		Body: &roaCS.DescribeClusterNodesResponseBody{
			Nodes: nodes,
			Page: &roaCS.DescribeClusterNodesResponseBodyPage{
				TotalCount: tea.Int32(totalCount),
				PageSize:   tea.Int32(pageSize),
			},
		},
	}
}

func fakePoolsPage(pools ...*roaCS.DescribeClusterNodePoolsResponseBodyNodepools) *roaCS.DescribeClusterNodePoolsResponse {
	return &roaCS.DescribeClusterNodePoolsResponse{
		Body: &roaCS.DescribeClusterNodePoolsResponseBody{Nodepools: pools},
	}
}

func fakePool(id, name string) *roaCS.DescribeClusterNodePoolsResponseBodyNodepools {
	return &roaCS.DescribeClusterNodePoolsResponseBodyNodepools{
		NodepoolInfo: &roaCS.DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo{
			NodepoolId: tea.String(id),
			Name:       tea.String(name),
		},
	}
}

func TestUnitKubernetesNodesPagination(t *testing.T) {
	client := &fakeKubernetesNodesClient{
		nodePages: []*roaCS.DescribeClusterNodesResponse{
			fakeNodesPage(3, 2, fakeNode("i-1", "node-1", "np-1"), fakeNode("i-2", "node-2", "np-1")),
			fakeNodesPage(3, 2, fakeNode("i-3", "node-3", "np-1")),
		},
	}

	result, err := fetchKubernetesNodes(client, "c-1", "", "")
	if err != nil {
		t.Fatalf("fetchKubernetesNodes returned error: %v", err)
	}
	if len(result.Nodes) != 3 {
		t.Fatalf("expected 3 nodes across two pages, got %d", len(result.Nodes))
	}
	if client.nodeCalls != 2 {
		t.Fatalf("expected 2 page requests, got %d", client.nodeCalls)
	}
	if got := tea.StringValue(client.nodeRequests[0].PageNumber); got != "1" {
		t.Fatalf("expected the first request on page 1, got %s", got)
	}
	if got := tea.StringValue(client.nodeRequests[1].PageNumber); got != "2" {
		t.Fatalf("expected the second request on page 2, got %s", got)
	}
	// Without a node pool filter the request must not narrow the query.
	if client.nodeRequests[0].NodepoolId != nil {
		t.Fatalf("expected no NodepoolId filter, got %s", tea.StringValue(client.nodeRequests[0].NodepoolId))
	}
}

func TestUnitKubernetesNodesStableSort(t *testing.T) {
	client := &fakeKubernetesNodesClient{
		nodePages: []*roaCS.DescribeClusterNodesResponse{
			fakeNodesPage(3, 3, fakeNode("i-c", "node-c", "np-1"), fakeNode("i-a", "node-a", "np-1"), fakeNode("i-b", "node-b", "np-1")),
		},
	}

	result, err := fetchKubernetesNodes(client, "c-1", "np-1", "")
	if err != nil {
		t.Fatalf("fetchKubernetesNodes returned error: %v", err)
	}
	want := []string{"node-a", "node-b", "node-c"}
	if len(result.Nodes) != len(want) {
		t.Fatalf("expected %d nodes, got %d", len(want), len(result.Nodes))
	}
	for index, name := range want {
		if got := tea.StringValue(result.Nodes[index].NodeName); got != name {
			t.Fatalf("expected node %d to be %s, got %s", index, name, got)
		}
	}
}

func TestUnitKubernetesNodesErrorPropagation(t *testing.T) {
	client := &fakeKubernetesNodesClient{
		nodeErr: fmt.Errorf("DescribeClusterNodes API error"),
	}

	result, err := fetchKubernetesNodes(client, "c-1", "np-1", "")
	if err == nil {
		t.Fatalf("expected the API error to be propagated, got nil")
	}
	if len(result.Nodes) != 0 {
		t.Fatalf("expected no nodes when the API fails, got %d", len(result.Nodes))
	}
}

func TestUnitKubernetesNodesNodepoolNameChange(t *testing.T) {
	// Name A and name B live in the same cluster: a name change must query B,
	// not keep returning A's nodes.
	client := &fakeKubernetesNodesClient{
		poolResponses: []*roaCS.DescribeClusterNodePoolsResponse{
			fakePoolsPage(fakePool("np-a", "pool-a"), fakePool("np-b", "pool-b")),
		},
		nodePages: []*roaCS.DescribeClusterNodesResponse{
			fakeNodesPage(1, 1, fakeNode("i-a", "node-a", "np-a")),
			fakeNodesPage(1, 1, fakeNode("i-b", "node-b", "np-b")),
		},
	}

	first, err := fetchKubernetesNodes(client, "c-1", "", "pool-a")
	if err != nil {
		t.Fatalf("fetchKubernetesNodes(pool-a) returned error: %v", err)
	}
	if first.NodepoolId != "np-a" {
		t.Fatalf("expected pool-a to resolve to np-a, got %s", first.NodepoolId)
	}
	if got := tea.StringValue(client.nodeRequests[0].NodepoolId); got != "np-a" {
		t.Fatalf("expected the first query to filter np-a, got %s", got)
	}

	second, err := fetchKubernetesNodes(client, "c-1", "", "pool-b")
	if err != nil {
		t.Fatalf("fetchKubernetesNodes(pool-b) returned error: %v", err)
	}
	if second.NodepoolId != "np-b" {
		t.Fatalf("expected pool-b to resolve to np-b, got %s", second.NodepoolId)
	}
	if got := tea.StringValue(client.nodeRequests[1].NodepoolId); got != "np-b" {
		t.Fatalf("expected the second query to filter np-b, got %s", got)
	}
}

func TestUnitKubernetesNodesNodepoolNameReplaced(t *testing.T) {
	// The same name is served by a different node pool on a later read: the
	// resolved ID of the current read wins over any previously resolved value.
	client := &fakeKubernetesNodesClient{
		poolResponses: []*roaCS.DescribeClusterNodePoolsResponse{
			fakePoolsPage(fakePool("np-old", "pool-a")),
			fakePoolsPage(fakePool("np-new", "pool-a")),
		},
		nodePages: []*roaCS.DescribeClusterNodesResponse{
			fakeNodesPage(1, 1, fakeNode("i-old", "node-old", "np-old")),
			fakeNodesPage(1, 1, fakeNode("i-new", "node-new", "np-new")),
		},
	}

	first, err := fetchKubernetesNodes(client, "c-1", "", "pool-a")
	if err != nil {
		t.Fatalf("first fetchKubernetesNodes returned error: %v", err)
	}
	if first.NodepoolId != "np-old" {
		t.Fatalf("expected the first read to resolve np-old, got %s", first.NodepoolId)
	}

	second, err := fetchKubernetesNodes(client, "c-1", "", "pool-a")
	if err != nil {
		t.Fatalf("second fetchKubernetesNodes returned error: %v", err)
	}
	if second.NodepoolId != "np-new" {
		t.Fatalf("expected the second read to resolve np-new, got %s", second.NodepoolId)
	}
	if got := tea.StringValue(client.nodeRequests[1].NodepoolId); got != "np-new" {
		t.Fatalf("expected the second query to filter np-new, got %s", got)
	}
}

func TestUnitKubernetesNodesConfiguredIdPrecedence(t *testing.T) {
	client := &fakeKubernetesNodesClient{
		nodePages: []*roaCS.DescribeClusterNodesResponse{
			fakeNodesPage(1, 1, fakeNode("i-a", "node-a", "np-a")),
		},
	}

	result, err := fetchKubernetesNodes(client, "c-1", "np-a", "pool-b")
	if err != nil {
		t.Fatalf("fetchKubernetesNodes returned error: %v", err)
	}
	if result.NodepoolId != "np-a" {
		t.Fatalf("expected the configured ID np-a to win, got %s", result.NodepoolId)
	}
	if client.poolCalls != 0 {
		t.Fatalf("expected no node pool lookup when an ID is configured, got %d calls", client.poolCalls)
	}
	if got := tea.StringValue(client.nodeRequests[0].NodepoolId); got != "np-a" {
		t.Fatalf("expected the query to filter np-a, got %s", got)
	}
}

func TestUnitKubernetesNodesUnknownNameReturnsEmpty(t *testing.T) {
	client := &fakeKubernetesNodesClient{
		poolResponses: []*roaCS.DescribeClusterNodePoolsResponse{
			fakePoolsPage(fakePool("np-a", "pool-a")),
		},
		nodePages: []*roaCS.DescribeClusterNodesResponse{
			fakeNodesPage(1, 1, fakeNode("i-a", "node-a", "np-a")),
		},
	}

	result, err := fetchKubernetesNodes(client, "c-1", "", "pool-missing")
	if err != nil {
		t.Fatalf("fetchKubernetesNodes returned error: %v", err)
	}
	if result.NodepoolId != "" || len(result.Nodes) != 0 {
		t.Fatalf("expected an empty result, got pool %q with %d nodes", result.NodepoolId, len(result.Nodes))
	}
	// Falling back to the whole cluster would answer a different question.
	if client.nodeCalls != 0 {
		t.Fatalf("expected no node query for an unknown pool name, got %d calls", client.nodeCalls)
	}
}

func TestUnitKubernetesNodesNodepoolLookupError(t *testing.T) {
	client := &fakeKubernetesNodesClient{
		poolErr: fmt.Errorf("DescribeClusterNodePools API error"),
	}

	result, err := fetchKubernetesNodes(client, "c-1", "", "pool-a")
	if err == nil {
		t.Fatalf("expected the node pool lookup error to be propagated, got nil")
	}
	if len(result.Nodes) != 0 {
		t.Fatalf("expected no nodes when the lookup fails, got %d", len(result.Nodes))
	}
}
