// Package rocketmq hosts the terraform-plugin-framework native resources and
// data sources for the RocketMQ (2022-08-01) service. The legacy Plugin SDK
// resources live in package alicloud; new types are added here on Framework.
package rocketmq

import (
	"fmt"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

const (
	rocketmqProduct  = "RocketMQ"
	rocketmqVersion  = "2022-08-01"
	planNotFoundCode = "DisasterRecoveryPlan.NotFound"

	// planStatusRunning is the live steady state; CloudSpec's statusDefine only
	// lists CREATED, so Create must also accept RUNNING or it times out on a
	// plan that has already settled.
	planStatusCreated = "CREATED"
	planStatusRunning = "RUNNING"
)

// describeDisasterRecoveryPlan issues GetDisasterRecoveryPlan and returns the
// plan object (the response "data" payload), or nil when the plan no longer
// exists. A nil object with a nil error signals removal to the caller.
func describeDisasterRecoveryPlan(client *connectivity.AliyunClient, planId string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/disaster_recovery/%s", planId)
	response, err := client.RoaGet(rocketmqProduct, rocketmqVersion, action, map[string]*string{}, nil, nil)
	if err != nil {
		if isPlanNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	data, _ := response["data"].(map[string]interface{})
	return data, nil
}

// listDisasterRecoveryPlans issues ListDisasterRecoveryPlans with pagination
// and returns the raw plan objects. instanceId, when non-empty, is sent as a
// query filter so the server narrows the result set before client filtering.
func listDisasterRecoveryPlans(client *connectivity.AliyunClient, instanceId string) ([]map[string]interface{}, error) {
	pageNumber := 1
	const pageSize = 100
	objects := make([]map[string]interface{}, 0)
	for {
		query := map[string]*string{
			"pageNumber": tea.String(fmt.Sprint(pageNumber)),
			"pageSize":   tea.String(fmt.Sprint(pageSize)),
		}
		if instanceId != "" {
			query["instanceId"] = tea.String(instanceId)
		}
		response, err := client.RoaGet(rocketmqProduct, rocketmqVersion, "/disaster_recovery", query, nil, nil)
		if err != nil {
			return nil, err
		}
		data, _ := response["data"].(map[string]interface{})
		rawList, _ := data["list"].([]interface{})
		if len(rawList) == 0 {
			break
		}
		for _, raw := range rawList {
			if obj, ok := raw.(map[string]interface{}); ok {
				objects = append(objects, obj)
			}
		}
		totalCount := toInt(data["totalCount"])
		if totalCount <= 0 || len(objects) >= totalCount {
			break
		}
		pageNumber++
	}
	return objects, nil
}

// waitForPlanStatus polls the plan until status reaches one of target, the
// timeout elapses, or the plan disappears. Used by Create to settle a freshly
// created plan into CREATED or RUNNING.
func waitForPlanStatus(client *connectivity.AliyunClient, planId string, target []string, timeout, interval time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		object, err := describeDisasterRecoveryPlan(client, planId)
		if err != nil {
			return err
		}
		if object == nil {
			return fmt.Errorf("disaster recovery plan %s disappeared before reaching %s", planId, strings.Join(target, "/"))
		}
		status := fmt.Sprint(object["planStatus"])
		for _, t := range target {
			if status == t {
				return nil
			}
		}
		time.Sleep(interval)
	}
	return fmt.Errorf("waiting for disaster recovery plan %s to reach %s timed out after %s", planId, strings.Join(target, "/"), timeout)
}

// waitForPlanGone polls until Describe returns no plan, used by Delete to
// confirm removal. A missing plan is success; transient errors propagate.
func waitForPlanGone(client *connectivity.AliyunClient, planId string, timeout, interval time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		object, err := describeDisasterRecoveryPlan(client, planId)
		if err != nil {
			return err
		}
		if object == nil {
			return nil
		}
		time.Sleep(interval)
	}
	return fmt.Errorf("waiting for disaster recovery plan %s to disappear timed out after %s", planId, timeout)
}

// isPlanNotFound reports whether err is the DisasterRecoveryPlan.NotFound
// signal or a generic 404 from the ROA tea client. The Code arrives on a
// *tea.SDKError; the substring fallback covers wrapped ComplexError forms.
func isPlanNotFound(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok {
		if tea.IntValue(e.StatusCode) == 404 {
			return true
		}
		if e.Code != nil && (*e.Code == planNotFoundCode || strings.Contains(*e.Code, "NotFound")) {
			return true
		}
	}
	return strings.Contains(err.Error(), planNotFoundCode)
}

// toInt coerces the JSON-decoded numeric forms (float64, json.Number, string)
// the tea ROA client hands back into an int. The default branch parses the
// value's string form so totalCount survives whatever decoder the ROA client
// uses; otherwise ListDisasterRecoveryPlans stops after the first page.
func toInt(v interface{}) int {
	switch n := v.(type) {
	case nil:
		return 0
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case string:
		var i int
		_, _ = fmt.Sscanf(n, "%d", &i)
		return i
	default:
		var i int
		_, _ = fmt.Sscanf(fmt.Sprint(v), "%d", &i)
		return i
	}
}
