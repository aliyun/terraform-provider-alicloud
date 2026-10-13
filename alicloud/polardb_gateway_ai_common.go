package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	polarDBGatewayAIProduct  = "polardb"
	polarDBGatewayAIVersion  = "2017-08-01"
	polarDBGatewayAIPageSize = 100
)

func polarDBGatewayAIChildID(gatewayID, childID string) string {
	return fmt.Sprintf("%s:%s", gatewayID, childID)
}

func parsePolarDBGatewayAIChildID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid PolarDB AI gateway child resource ID %q, expected <gateway_id>:<child_id>", id)
	}
	return parts[0], parts[1], nil
}

func normalizePolarDBGatewayAIDecimal(value interface{}) string {
	normalized := strings.TrimSpace(fmt.Sprint(value))
	if !strings.Contains(normalized, ".") {
		return normalized
	}
	normalized = strings.TrimRight(strings.TrimRight(normalized, "0"), ".")
	if normalized == "" || normalized == "-" {
		return "0"
	}
	return normalized
}

func normalizePolarDBGatewayAIJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	value, err := normalizeJsonString(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return value
}

func validatePolarDBGatewayAIRouteRules(v interface{}, k string) ([]string, []error) {
	value := v.(string)
	var rules []interface{}
	if err := json.Unmarshal([]byte(value), &rules); err != nil {
		return nil, []error{fmt.Errorf("%q must be a valid JSON array: %s", k, err)}
	}
	if len(rules) == 0 {
		return nil, []error{fmt.Errorf("%q must contain at least one routing rule", k)}
	}
	return nil, nil
}

func redactPolarDBGatewayAISecrets(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(typed))
		for key, item := range typed {
			if strings.EqualFold(key, "ApiKey") {
				result[key] = "***"
				continue
			}
			result[key] = redactPolarDBGatewayAISecrets(item)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(typed))
		for i, item := range typed {
			result[i] = redactPolarDBGatewayAISecrets(item)
		}
		return result
	default:
		return value
	}
}

func debugPolarDBGatewayAI(action string, response, request map[string]interface{}) {
	addDebug(action, redactPolarDBGatewayAISecrets(response), redactPolarDBGatewayAISecrets(request))
}

func callPolarDBGatewayAI(client *connectivity.AliyunClient, action string, request map[string]interface{}, timeout time.Duration, retryCodes ...string) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		response, err = client.RpcPost(polarDBGatewayAIProduct, polarDBGatewayAIVersion, action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, retryCodes) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	debugPolarDBGatewayAI(action, response, request)
	return response, err
}

func suppressPolarDBGatewayAIUnsetOptional(_, old, new string, _ *schema.ResourceData) bool {
	// ModifyModelApi and ModifyModelService retain the current value when these
	// optional fields are omitted or empty, so the API cannot clear them.
	return old != "" && new == ""
}

func listPolarDBGatewayAI(client *connectivity.AliyunClient, action string, request map[string]interface{}) ([]map[string]interface{}, error) {
	objects := make([]map[string]interface{}, 0)
	for pageNumber := 1; ; pageNumber++ {
		pageRequest := make(map[string]interface{}, len(request)+3)
		for key, value := range request {
			pageRequest[key] = value
		}
		pageRequest["RegionId"] = client.RegionId
		pageRequest["PageNumber"] = pageNumber
		pageRequest["PageSize"] = polarDBGatewayAIPageSize

		response, err := callPolarDBGatewayAI(client, action, pageRequest, 5*time.Minute)
		if err != nil {
			return nil, err
		}
		items := polarDBGatewayAIItems(response["Items"])
		objects = append(objects, items...)
		if len(items) < polarDBGatewayAIPageSize {
			break
		}
	}
	return objects, nil
}

func polarDBGatewayAIItems(raw interface{}) []map[string]interface{} {
	if wrapper, ok := raw.(map[string]interface{}); ok {
		for _, key := range []string{"Item", "ModelService", "ModelApi"} {
			if value, exists := wrapper[key]; exists {
				raw = value
				break
			}
		}
	}
	values, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(values))
	for _, value := range values {
		if item, ok := value.(map[string]interface{}); ok {
			result = append(result, item)
		}
	}
	return result
}

func findPolarDBGatewayAIItem(client *connectivity.AliyunClient, action, gatewayID, requestIDKey, responseIDKey, childID string) (map[string]interface{}, error) {
	request := map[string]interface{}{"GwClusterId": gatewayID}
	if requestIDKey != "" {
		request[requestIDKey] = childID
	}
	items, err := listPolarDBGatewayAI(client, action, request)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if fmt.Sprint(item[responseIDKey]) == childID {
			return item, nil
		}
	}
	return nil, nil
}

func copyPolarDBGatewayAIFilters(d *schema.ResourceData, request map[string]interface{}, fields map[string]string) {
	for key, requestKey := range fields {
		if value, ok := d.GetOk(key); ok {
			request[requestKey] = value
		}
	}
}

func matchesPolarDBGatewayAIFilters(d *schema.ResourceData, object map[string]interface{}, fields map[string][]string) bool {
	for key, responseKeys := range fields {
		expected, ok := d.GetOk(key)
		if !ok {
			continue
		}
		var actual interface{}
		for _, responseKey := range responseKeys {
			if value, exists := object[responseKey]; exists && value != nil {
				actual = value
				break
			}
		}
		if fmt.Sprint(actual) != fmt.Sprint(expected) {
			return false
		}
	}
	return true
}

func stringListFromSchema(raw interface{}) []string {
	values, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, fmt.Sprint(value))
	}
	return result
}
