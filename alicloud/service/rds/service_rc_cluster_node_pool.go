package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

const (
	rcClusterNodePoolRdsProduct = "Rds"
	rcClusterNodePoolRdsVersion = "2014-08-15"

	rcClusterNodePoolPollInterval = 10 * time.Second
	rcClusterNodePoolCreateWindow = 30 * time.Minute
	rcClusterNodePoolUpdateWindow = 30 * time.Minute
	rcClusterNodePoolDeleteWindow = 30 * time.Minute
)

// nodePoolTransientStates lists the states reported while the node pool is
// still converging. Anything else (including an omitted state) counts as done.
var nodePoolTransientStates = map[string]struct{}{
	"creating":     {},
	"initializing": {},
	"pending":      {},
	"scaling":      {},
	"removing":     {},
	"deleting":     {},
}

// notFoundNodePoolCodes are the API codes that mean the node pool cannot be
// read any more: the pool itself, its cluster, or a cluster that does not carry
// the RDS Custom extension.
var notFoundNodePoolCodes = []string{
	"InvalidParam.ClusterIdNotExist",
	"InvalidParam.NodePoolIdNotExist",
	"OperationDenied.ClusterNotSupport",
	// DescribeRCClusterNodePoolDetail reports InvalidDBInstanceId.NotFound once
	// the node pool has been removed.
	"InvalidDBInstanceId.NotFound",
}

// describeRCClusterNodePool returns the raw DescribeRCClusterNodePoolDetail
// response, or nil when the node pool does not exist.
func describeRCClusterNodePool(client *connectivity.AliyunClient, clusterId, nodePoolId string) (map[string]interface{}, error) {
	response, err := client.RpcPost(rcClusterNodePoolRdsProduct, rcClusterNodePoolRdsVersion, "DescribeRCClusterNodePoolDetail", nil, map[string]interface{}{
		"RegionId":   client.RegionId,
		"ClusterId":  clusterId,
		"NodePoolId": nodePoolId,
	}, true)
	if err != nil {
		if isRCClusterNodePoolNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	info, _ := response["NodePoolInfo"].(map[string]interface{})
	if id, _ := info["NodePoolId"].(string); id == "" {
		return nil, nil
	}
	return response, nil
}

// waitRCClusterNodePoolVisible polls until the node pool can be read and is no
// longer converging, then returns its detail response.
func waitRCClusterNodePoolVisible(ctx context.Context, client *connectivity.AliyunClient, clusterId, nodePoolId string, window time.Duration) (map[string]interface{}, error) {
	latest := map[string]interface{}(nil)
	_, err := pollRCClusterNodePool(ctx, window, func() (bool, error) {
		object, err := describeRCClusterNodePool(client, clusterId, nodePoolId)
		if err != nil {
			return false, err
		}
		if object == nil {
			return false, nil
		}
		latest = object
		if _, transient := nodePoolTransientStates[strings.ToLower(rcClusterNodePoolDetailState(object))]; transient {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return latest, nil
}

// waitRCClusterNodePoolDeleted polls until the node pool can no longer be read.
func waitRCClusterNodePoolDeleted(ctx context.Context, client *connectivity.AliyunClient, clusterId, nodePoolId string, window time.Duration) error {
	_, err := pollRCClusterNodePool(ctx, window, func() (bool, error) {
		object, err := describeRCClusterNodePool(client, clusterId, nodePoolId)
		if err != nil {
			return false, err
		}
		return object == nil, nil
	})
	return err
}

// pollRCClusterNodePool re-runs probe every poll interval until it reports done,
// the window expires or the context is cancelled. An error from probe is fatal.
func pollRCClusterNodePool(ctx context.Context, window time.Duration, probe func() (bool, error)) (map[string]interface{}, error) {
	deadline := time.Now().Add(window)
	for {
		done, err := probe()
		if err != nil {
			return nil, err
		}
		if done {
			return nil, nil
		}
		if time.Now().Add(rcClusterNodePoolPollInterval).After(deadline) {
			return nil, fmt.Errorf("timed out after %s waiting for the node pool to converge", window)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(rcClusterNodePoolPollInterval):
		}
	}
}

// isRCClusterNodePoolNotFound reports whether err is one of the API codes that
// mean the node pool cannot be read any more.
func isRCClusterNodePoolNotFound(err error) bool {
	if err == nil {
		return false
	}
	message := err.Error()
	for _, code := range notFoundNodePoolCodes {
		if strings.Contains(message, code) {
			return true
		}
	}
	return false
}

// rcClusterNodePoolDetailState returns the status state of a detail response.
func rcClusterNodePoolDetailState(object map[string]interface{}) string {
	status, _ := object["Status"].(map[string]interface{})
	state, _ := status["State"].(string)
	return state
}

// expandRCClusterNodePoolNodePoolInfo renders the NodePoolInfo request parameter.
func expandRCClusterNodePoolNodePoolInfo(name string) (string, error) {
	if name == "" {
		return "", nil
	}
	payload, err := json.Marshal(map[string]interface{}{"Name": name})
	if err != nil {
		return "", err
	}
	return string(payload), nil
}
