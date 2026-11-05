package vms

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
)

// Dyvmsapi (Voice Messaging Service) RPC product/version constants shared by every
// call in this package; ServicePackage.Name in register.go mirrors the product code.
const (
	vmsProduct = "Dyvmsapi"
	vmsVersion = "2017-05-25"
)

const (
	skillCreateTimeout = 5 * time.Minute
	skillReadTimeout   = 5 * time.Minute
	skillUpdateTimeout = 5 * time.Minute
	skillDeleteTimeout = 5 * time.Minute
	skillListPageSize  = 100
)

// Service wraps the Voice Messaging Service (Dyvmsapi) RPC calls used by the
// resources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService constructs a Voice Messaging Service from the provider client.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// retryRpcPost wraps RpcPost in the helper retry loop with the shared transient
// classifier. ctx carries the lifecycle deadline, so a cancelled operation stops
// retrying.
func retryRpcPost(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	return helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPost(vmsProduct, vmsVersion, action, nil, request, true)
	}, timeout, retryableVmsError)
}

// retryableVmsError retries only the shared transient classifier's errors. The
// Dyvmsapi Skill APIs return terminal 4xx for validation conflicts (duplicate
// skill name, not found), which the callers resolve instead of a blind retry
// that would only delay the same failure.
func retryableVmsError(err error) bool {
	return helper.NeedRetry(err)
}

// parseSkillId splits the composite resource id "<enterprise_id>:<skill_id>"
// into its numeric parts. ok is false when the id is not the expected shape.
func parseSkillId(id string) (enterpriseId, skillId int64, ok bool) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	eid, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return 0, 0, false
	}
	sid, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		return 0, 0, false
	}
	return eid, sid, true
}

// toInt64 coerces an RPC response value (float64, json.Number, int64, int, numeric
// string) to int64. The tea RPC client decodes large integers as float64 or
// json.Number depending on the SDK path; this keeps the skill Id exact.
func toInt64(v interface{}) (int64, error) {
	switch n := v.(type) {
	case float64:
		return int64(n), nil
	case float32:
		return int64(n), nil
	case int64:
		return n, nil
	case int32:
		return int64(n), nil
	case int:
		return int64(n), nil
	case json.Number:
		return n.Int64()
	case string:
		if strings.TrimSpace(n) == "" {
			return 0, fmt.Errorf("empty integer value")
		}
		return strconv.ParseInt(n, 10, 64)
	default:
		return 0, fmt.Errorf("unsupported integer type %T", v)
	}
}
