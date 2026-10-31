package rds

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

const (
	rcClusterRdsProduct = "Rds"
	rcClusterRdsVersion = "2014-08-15"

	rcClusterPollInterval = 10 * time.Second
	rcClusterCreateWindow = 30 * time.Minute
	rcClusterDeleteWindow = 10 * time.Minute

	// rcClusterStatusUninitialized is the status a cluster reports before it is
	// imported into RDS Custom, and again once the import has been removed.
	rcClusterStatusUninitialized = "UNINITIALIZED"
	// rcClusterStatusAvailable is the status of a cluster whose import has
	// finished and that can carry RDS Custom instances and node pools.
	rcClusterStatusAvailable = "avail"
)

// describeNotFoundRCVClusterCodes are the API codes that mean the cluster is not
// imported into RDS Custom. OperationDenied.ClusterNotSupport is returned for
// clusters that do not carry the RDS Custom extension.
var describeNotFoundRCVClusterCodes = []string{
	"InvalidParam.ClusterIdNotExist",
	"OperationDenied.ClusterNotSupport",
}

// deleteNotFoundRCVClusterCodes are the API codes that mean the import is
// already gone. OperationDenied.ClusterNotSupport is deliberately absent: the
// delete API rejects a cluster whose import is still converging with that code,
// so it must be retried instead of being treated as a finished deletion.
var deleteNotFoundRCVClusterCodes = []string{
	"InvalidParam.ClusterIdNotExist",
}

// unwrapRCVClusterResponse returns the payload of a DescribeRCVCluster
// response. The operation declares a data wrapper while the gateway currently
// returns the fields flat, so both shapes are accepted.
func unwrapRCVClusterResponse(response map[string]interface{}) map[string]interface{} {
	if data, ok := response["data"].(map[string]interface{}); ok && len(data) > 0 {
		return data
	}
	return response
}

// describeRCVCluster returns the raw DescribeRCVCluster payload, or nil when the
// cluster is not imported into RDS Custom.
func describeRCVCluster(client *connectivity.AliyunClient, clusterId string) (map[string]interface{}, error) {
	response, err := client.RpcPost(rcClusterRdsProduct, rcClusterRdsVersion, "DescribeRCVCluster", nil, map[string]interface{}{
		"RegionId":  client.RegionId,
		"ClusterId": clusterId,
	}, true)
	if err != nil {
		if isRCVClusterNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	info := unwrapRCVClusterResponse(response)
	status := rcVClusterStatus(info)
	log.Printf("[DEBUG] DescribeRCVCluster(%s) status=%q", clusterId, status)
	if strings.EqualFold(status, rcClusterStatusUninitialized) {
		return nil, nil
	}
	return info, nil
}

// waitRCVClusterAvailable polls until the import finished and the cluster can
// carry RDS Custom node pools and instances.
func waitRCVClusterAvailable(ctx context.Context, client *connectivity.AliyunClient, clusterId string, window time.Duration) (map[string]interface{}, error) {
	var latest map[string]interface{}
	_, err := pollRCVCluster(ctx, window, func() (bool, error) {
		info, err := describeRCVCluster(client, clusterId)
		if err != nil {
			return false, err
		}
		if info == nil {
			return false, nil
		}
		latest = info
		return strings.EqualFold(rcVClusterStatus(info), rcClusterStatusAvailable), nil
	})
	if err != nil {
		return nil, err
	}
	return latest, nil
}

// deleteRCVCluster removes the import. The delete API rejects a cluster whose
// import is still converging with OperationDenied.ClusterNotSupport, so the call
// is retried until that condition clears, then the removal is awaited.
func deleteRCVCluster(ctx context.Context, client *connectivity.AliyunClient, clusterId string, window time.Duration) (string, error) {
	deadline := time.Now().Add(window)
	lastStatus := ""
	for {
		response, err := client.RpcPost(rcClusterRdsProduct, rcClusterRdsVersion, "DeleteRCVCluster", nil, map[string]interface{}{
			"RegionId":  client.RegionId,
			"ClusterId": clusterId,
		}, true)
		log.Printf("[DEBUG] DeleteRCVCluster(%s) response=%v err=%v", clusterId, response, err)
		if err == nil {
			break
		}
		if isRCVClusterDeleteNotFound(err) {
			return lastStatus, nil
		}
		if !isRCVClusterNotSupported(err) {
			return lastStatus, err
		}

		info, describeErr := describeRCVCluster(client, clusterId)
		if describeErr != nil {
			return lastStatus, describeErr
		}
		if info == nil {
			return lastStatus, nil
		}
		lastStatus = rcVClusterStatus(info)
		if time.Now().Add(rcClusterPollInterval).After(deadline) {
			return lastStatus, fmt.Errorf("timed out after %s waiting for the RDS Custom container cluster import to converge before deletion", window)
		}
		select {
		case <-ctx.Done():
			return lastStatus, ctx.Err()
		case <-time.After(rcClusterPollInterval):
		}
	}

	return waitRCVClusterDeleted(ctx, client, clusterId, window)
}

// waitRCVClusterDeleted polls until the cluster is no longer imported and
// returns the last observed VClusterStatus for diagnostics.
func waitRCVClusterDeleted(ctx context.Context, client *connectivity.AliyunClient, clusterId string, window time.Duration) (string, error) {
	var lastStatus string
	_, err := pollRCVCluster(ctx, window, func() (bool, error) {
		info, err := describeRCVCluster(client, clusterId)
		if err != nil {
			return false, err
		}
		if info == nil {
			return true, nil
		}
		lastStatus = rcVClusterStatus(info)
		return false, nil
	})
	return lastStatus, err
}

// pollRCVCluster re-runs probe every poll interval until it reports done, the
// window expires or the context is cancelled. A probe error is fatal.
func pollRCVCluster(ctx context.Context, window time.Duration, probe func() (bool, error)) (map[string]interface{}, error) {
	deadline := time.Now().Add(window)
	for {
		done, err := probe()
		if err != nil {
			return nil, err
		}
		if done {
			return nil, nil
		}
		if time.Now().Add(rcClusterPollInterval).After(deadline) {
			return nil, fmt.Errorf("timed out after %s waiting for the RDS Custom container cluster to converge", window)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(rcClusterPollInterval):
		}
	}
}

// isRCVClusterNotFound reports whether err means the cluster is not imported
// into RDS Custom.
func isRCVClusterNotFound(err error) bool {
	return containsAnyCode(err, describeNotFoundRCVClusterCodes)
}

// isRCVClusterDeleteNotFound reports whether the import is already gone.
func isRCVClusterDeleteNotFound(err error) bool {
	return containsAnyCode(err, deleteNotFoundRCVClusterCodes)
}

// isRCVClusterNotSupported reports whether the API rejected the call because the
// cluster import is not ready yet.
func isRCVClusterNotSupported(err error) bool {
	return containsAnyCode(err, []string{"OperationDenied.ClusterNotSupport"})
}

func containsAnyCode(err error, codes []string) bool {
	if err == nil {
		return false
	}
	message := err.Error()
	for _, code := range codes {
		if strings.Contains(message, code) {
			return true
		}
	}
	return false
}

// rcVClusterStatus returns the VClusterStatus of a DescribeRCVCluster payload.
func rcVClusterStatus(object map[string]interface{}) string {
	status, _ := object["VClusterStatus"].(string)
	return status
}
