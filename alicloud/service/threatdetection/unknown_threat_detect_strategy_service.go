package threatdetection

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	unknownThreatDetectStrategyTypeName = "alicloud_threat_detection_unknown_threat_detect_strategy"
	// unknownThreatDetectStrategyReadTimeout bounds the read-back and delete RPC
	// retries. Framework has no user-configurable resource timeouts.
	unknownThreatDetectStrategyReadTimeout  = 2 * time.Minute
	unknownThreatDetectStrategyWriteTimeout = 5 * time.Minute
)

// maxExactFloatInteger is the largest integer float64 represents without loss;
// a strategy Id above it cannot be recovered from the number form.
const maxExactFloatInteger = 1 << 53

// UnknownThreatDetectStrategyId converts a strategy Id from
// CreateUnknownThreatDetectStrategy or ListUnknownThreatDetectStrategy into the
// canonical decimal identity used in state, list filters and delete requests.
// The published response contract declares Id as type integer/format int64 while
// JSON decoding yields float64; some accounts return a numeric string, so both
// are accepted. Containers, fractional numbers, non-finite values and numbers
// outside the lossless integer range are contract drift and are rejected instead
// of being stringified into an id.
func UnknownThreatDetectStrategyId(raw interface{}) (string, error) {
	switch v := raw.(type) {
	case string:
		return unknownThreatDetectStrategyDecimalInteger(v)
	case json.Number:
		return unknownThreatDetectStrategyDecimalInteger(v.String())
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return "", fmt.Errorf("Id %v is not finite", v)
		}
		if v != math.Trunc(v) {
			return "", fmt.Errorf("Id %v is not an integer", v)
		}
		if v <= 0 {
			return "", fmt.Errorf("Id %v is not positive", v)
		}
		if math.Abs(v) > maxExactFloatInteger {
			return "", fmt.Errorf("Id %v exceeds the lossless integer range", v)
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return unknownThreatDetectStrategyPositiveInt(int64(v))
	case int64:
		return unknownThreatDetectStrategyPositiveInt(v)
	case uint64:
		if v == 0 || v > uint64(1<<63-1) {
			return "", fmt.Errorf("Id %d is not a positive int64", v)
		}
		return strconv.FormatUint(v, 10), nil
	case nil:
		return "", fmt.Errorf("Id is missing")
	default:
		return "", fmt.Errorf("Id has unexpected type %T", raw)
	}
}

func unknownThreatDetectStrategyPositiveInt(v int64) (string, error) {
	if v <= 0 {
		return "", fmt.Errorf("Id %d is not positive", v)
	}
	return strconv.FormatInt(v, 10), nil
}

func unknownThreatDetectStrategyDecimalInteger(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("Id is empty")
	}
	if unknownThreatDetectStrategyHasNonIntegerRune(s) {
		return "", fmt.Errorf("Id %q is not an integer", s)
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return "", fmt.Errorf("Id %q is not an int64: %w", s, err)
	}
	return unknownThreatDetectStrategyPositiveInt(v)
}

func unknownThreatDetectStrategyHasNonIntegerRune(s string) bool {
	for _, c := range s {
		if c == '.' || c == 'e' || c == 'E' || c == '+' {
			return true
		}
	}
	return false
}

// DescribeUnknownThreatDetectStrategy returns the strategy whose Id equals id.
// ListUnknownThreatDetectStrategy is a collection API, so the target is selected
// by its exact Id instead of the first returned item.
func (s *Service) DescribeUnknownThreatDetectStrategy(id string) (object map[string]interface{}, err error) {
	return s.DescribeUnknownThreatDetectStrategyContext(context.Background(), id)
}

func (s *Service) DescribeUnknownThreatDetectStrategyContext(ctx context.Context, id string) (object map[string]interface{}, err error) {
	client := s.client
	query := map[string]interface{}{
		"Id":          id,
		"CurrentPage": 1,
		"PageSize":    helper.PageSizeXLarge,
	}
	action := "ListUnknownThreatDetectStrategy"

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, unknownThreatDetectStrategyReadTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("Sas", "2018-12-03", action, query, map[string]interface{}{}, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, query)
	if err != nil {
		// A transport-level 404 is definitive. A successful RPC carrying Code 404
		// only in its body is classified after the collection container is
		// validated below, so a malformed body cannot be turned into NotFound.
		if errs.IsTypedNotFound(err) && !errs.IsResponseCodeNotFound(err) {
			return object, errs.WrapErrorf(errs.NotFoundErr("UnknownThreatDetectStrategy", id), errs.NotFoundMsg, response)
		}
		if !errs.IsResponseCodeNotFound(err) {
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	code, _ := jsonpath.Get("$.Code", response)

	// The documented contract returns Data as an array. Validate the container
	// before any NotFound classification: a null, scalar or object shape is a
	// malformed response, not evidence that the strategy no longer exists. Only
	// a confirmed array that is empty or lacks the requested Id means the
	// strategy is gone.
	raw, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	list, ok := raw.([]interface{})
	if !ok {
		return object, errs.WrapErrorf(fmt.Errorf("unexpected Data response type %T", raw), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(code) == "404" || len(list) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("UnknownThreatDetectStrategy", id), errs.NotFoundMsg, response)
	}

	for _, item := range list {
		strategy, ok := item.(map[string]interface{})
		if !ok {
			return object, errs.WrapErrorf(fmt.Errorf("unexpected Data element type %T", item), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		itemId, idErr := UnknownThreatDetectStrategyId(strategy["Id"])
		if idErr != nil {
			return object, errs.WrapErrorf(idErr, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		if itemId == id {
			return strategy, nil
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("UnknownThreatDetectStrategy", id), errs.NotFoundMsg, response)
}

// DeleteUnknownThreatDetectStrategy deletes the strategy. DeleteUnknownThreatDetectStrategy
// takes IdList as a repeatList; a single resource is removed by passing its Id as the
// only element. A typed 404 means the strategy is already gone.
func (s *Service) DeleteUnknownThreatDetectStrategy(id string) error {
	return s.DeleteUnknownThreatDetectStrategyContext(context.Background(), id)
}

func (s *Service) DeleteUnknownThreatDetectStrategyContext(ctx context.Context, id string) error {
	client := s.client
	action := "DeleteUnknownThreatDetectStrategy"
	request := map[string]interface{}{
		"IdList.1": id,
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, unknownThreatDetectStrategyWriteTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsTypedNotFound(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
