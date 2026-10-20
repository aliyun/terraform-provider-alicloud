package sas

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// policyFieldString reads a field from a ListCheckPolicies element or a
// CreateCheckPolicy response as a clean string. Numeric ids (PolicyId,
// DependentPolicyId) may decode as json.Number or float64 depending on the rpc
// client; fmt.Sprint of a float64 policy id lands in scientific notation and
// would corrupt the composite id, so every numeric field flows through here.
func policyFieldString(object map[string]interface{}, key string) string {
	v, ok := object[key]
	if !ok || v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case json.Number:
		return val.String()
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	default:
		return fmt.Sprintf("%v", val)
	}
}

const (
	sasProduct = "Sas"
	sasVersion = "2018-12-03"
)

// parsePolicyID splits a composite resource id of the form
// <policy_id>:<policy_type> used by the custom check standard policy.
func parsePolicyID(id string) (policyId, policyType string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid resource id %s, expected <policy_id>:<policy_type>", id)
	}
	return parts[0], parts[1], nil
}

// listCustomCheckStandardPolicies calls ListCheckPolicies. When policyId is
// set it is forwarded so the API can narrow the result; callers that need a
// specific policy still match on PolicyId to stay correct when the API
// ignores the filter.
func listCustomCheckStandardPolicies(client *connectivity.AliyunClient, policyType, policyId string) ([]map[string]interface{}, error) {
	const pageSize = 50
	out := make([]map[string]interface{}, 0)
	for page := 1; ; page++ {
		request := map[string]interface{}{
			"PageSize":    pageSize,
			"CurrentPage": page,
			"PolicyType":  policyType,
		}
		if policyId != "" {
			request["PolicyId"] = policyId
		}

		response, err := client.RpcPost(sasProduct, sasVersion, "ListCheckPolicies", nil, request, true)
		if err != nil {
			return nil, err
		}

		raw, _ := response["Policies"].([]interface{})
		for _, item := range raw {
			if m, ok := item.(map[string]interface{}); ok {
				out = append(out, m)
			}
		}
		// Stop once a page is short of the requested size; the describe path
		// sends PolicyId, so it always returns a single item on page one.
		if len(raw) < pageSize {
			break
		}
	}
	return out, nil
}

// describeCustomCheckStandardPolicy returns the policy object for the given
// composite id, or nil when no matching policy is found (so the caller can
// remove the resource from state on Read).
func describeCustomCheckStandardPolicy(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	policyId, policyType, err := parsePolicyID(id)
	if err != nil {
		return nil, err
	}

	policies, err := listCustomCheckStandardPolicies(client, policyType, policyId)
	if err != nil {
		return nil, err
	}

	for _, p := range policies {
		if policyFieldString(p, "PolicyId") == policyId {
			return p, nil
		}
	}
	return nil, nil
}
