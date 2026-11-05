package dataworks

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

const dataWorksFileReadTimeout = 5 * time.Minute

// maxExactFileInteger is the largest integer float64 represents without loss;
// a FileId above it cannot be recovered from the number form.
const maxExactFileInteger = 1 << 53

// DataWorksFileId converts a FileId from the CreateFile response or GetFile
// response into the canonical decimal identity used in state. The published
// contract declares FileId as a numeric value while some accounts return a
// numeric string, so both are accepted. Containers, fractional numbers,
// non-finite values and numbers outside the lossless integer range are
// contract drift and are rejected instead of being stringified into an id.
func DataWorksFileId(raw interface{}) (string, error) {
	switch v := raw.(type) {
	case string:
		if v == "" {
			return "", fmt.Errorf("FileId is empty")
		}
		return dataWorksFileDecimalInteger(v)
	case json.Number:
		return dataWorksFileDecimalInteger(v.String())
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return "", fmt.Errorf("FileId %v is not finite", v)
		}
		if v != math.Trunc(v) {
			return "", fmt.Errorf("FileId %v is not an integer", v)
		}
		if v <= 0 {
			return "", fmt.Errorf("FileId %v is not positive", v)
		}
		if math.Abs(v) > maxExactFileInteger {
			return "", fmt.Errorf("FileId %v exceeds the lossless integer range", v)
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return dataWorksFilePositiveInt(int64(v))
	case int64:
		return dataWorksFilePositiveInt(v)
	case nil:
		return "", fmt.Errorf("FileId is missing")
	default:
		return "", fmt.Errorf("FileId has unexpected type %T", raw)
	}
}

func dataWorksFilePositiveInt(v int64) (string, error) {
	if v <= 0 {
		return "", fmt.Errorf("FileId %d is not positive", v)
	}
	return strconv.FormatInt(v, 10), nil
}

func dataWorksFileDecimalInteger(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("FileId is empty")
	}
	for _, c := range s {
		if c == '.' || c == 'e' || c == 'E' || c == '+' {
			return "", fmt.Errorf("FileId %q is not an integer", s)
		}
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return "", fmt.Errorf("FileId %q is not an int64: %w", s, err)
	}
	return dataWorksFilePositiveInt(v)
}

// DescribeDataWorksFile returns the file object for the composite id
// "fileId:projectId". GetFile is used as the primary read; an absent or empty
// Data object means the file no longer exists.
func (s *Service) DescribeDataWorksFile(id string) (object map[string]interface{}, err error) {
	return s.DescribeDataWorksFileContext(context.Background(), id, "", "")
}

// DescribeDataWorksFileContext reads a file. The composite id is
// "fileId:<workspace>", where <workspace> is the configured project_id (or the
// project_identifier when only that was set). GetFile accepts either ProjectId
// or ProjectIdentifier, so the explicit identity from state is sent under the
// matching key; the id segment is the final fallback for imports that supply
// "fileId:projectId".
func (s *Service) DescribeDataWorksFileContext(ctx context.Context, id, projectId, projectIdentifier string) (object map[string]interface{}, err error) {
	client := s.client
	action := "GetFile"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}
	request := map[string]interface{}{"FileId": parts[0]}
	dataWorksFileSetWorkspace(request, projectId, projectIdentifier, parts[1])
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.RetryContext(ctx, dataWorksFileReadTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"File.NotFound", "InvalidFile.NotFound", "1000030004", "1101080008"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("dataworks_file", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	// GetFile nests the file object under $.Data.File (not $.Data). An absent or
	// empty File means the file no longer exists and is treated as NotFound so
	// Read removes the resource and Delete stays idempotent.
	v, gpErr := jsonpath.Get("$.Data.File", response)
	if gpErr != nil || v == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("dataworks_file", id), errs.NotFoundWithResponse, response)
	}
	obj, ok := v.(map[string]interface{})
	if !ok || len(obj) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("dataworks_file", id), errs.NotFoundWithResponse, response)
	}
	return obj, nil
}

// DeleteDataWorksFile deletes the file for the composite id "fileId:projectId".
// A typed NotFound means the file is already gone.
func (s *Service) DeleteDataWorksFile(id string) error {
	return s.DeleteDataWorksFileContext(context.Background(), id, "", "")
}

func (s *Service) DeleteDataWorksFileContext(ctx context.Context, id, projectId, projectIdentifier string) error {
	client := s.client
	action := "DeleteFile"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{"FileId": parts[0]}
	dataWorksFileSetWorkspace(request, projectId, projectIdentifier, parts[1])
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.RetryContext(ctx, dataWorksFileReadTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
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
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"File.NotFound", "InvalidFile.NotFound", "1000030004", "1101080008"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// dataWorksFileSetWorkspace sets the workspace identity on a File API request.
// The DataWorks File APIs (Get/Create/Update/Delete/Deploy/Submit) accept
// either ProjectId or ProjectIdentifier and require one. The caller passes the
// configured identity from state/plan: project_id is preferred, then the
// identifier, then the id segment as a project_id fallback for imports.
func dataWorksFileSetWorkspace(request map[string]interface{}, projectId, projectIdentifier, idWorkspace string) {
	switch {
	case projectId != "":
		request["ProjectId"] = projectId
	case projectIdentifier != "":
		request["ProjectIdentifier"] = projectIdentifier
	default:
		request["ProjectId"] = idWorkspace
	}
}
