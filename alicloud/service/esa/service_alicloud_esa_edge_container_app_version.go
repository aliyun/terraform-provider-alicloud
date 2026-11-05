package esa

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	edgeContainerAppVersionReadTimeout = 2 * time.Minute
	edgeContainerAppVersionNotFoundErr = "InvalidVersion.NotFound"
	edgeContainerAppAppNotFoundErr     = "InvalidApp.NotFound"
)

// EdgeContainerAppVersionId parses the canonical "app_id:version_id" identity
// stored in state. Both segments are required and must be non-empty.
func EdgeContainerAppVersionId(id string) (appId, versionId string, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return "", "", err
	}
	return parts[0], parts[1], nil
}

// DescribeEsaEdgeContainerAppVersion returns the raw GetEdgeContainerAppVersion
// response for the version identified by "app_id:version_id". A missing
// version or application is reported as a NotFound error so the Framework Read
// can remove the resource from state on refresh.
func (s *EsaServiceV2) DescribeEsaEdgeContainerAppVersion(id string) (object map[string]interface{}, err error) {
	return s.DescribeEsaEdgeContainerAppVersionContext(context.Background(), id)
}

func (s *EsaServiceV2) DescribeEsaEdgeContainerAppVersionContext(ctx context.Context, id string) (object map[string]interface{}, err error) {
	appId, versionId, err := EdgeContainerAppVersionId(id)
	if err != nil {
		return object, err
	}
	client := s.client
	query := map[string]interface{}{
		"AppId":     appId,
		"VersionId": versionId,
		"RegionId":  client.RegionId,
	}
	action := "GetEdgeContainerAppVersion"

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, edgeContainerAppVersionReadTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcGet("ESA", "2024-09-10", action, query, map[string]interface{}{})
		if callErr != nil {
			if errs.IsExpectedErrors(callErr, []string{edgeContainerAppVersionNotFoundErr, edgeContainerAppAppNotFoundErr}) {
				return retry.NonRetryableError(callErr)
			}
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
		if errs.IsExpectedErrors(err, []string{edgeContainerAppVersionNotFoundErr, edgeContainerAppAppNotFoundErr}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("EdgeContainerAppVersion", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	versionRaw, ok := response["Version"].(map[string]interface{})
	if !ok {
		return object, errs.WrapErrorf(fmt.Errorf("unexpected Version response type %T", response["Version"]), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return versionRaw, nil
}

// EsaEdgeContainerAppVersionStateRefreshFunc waits for the version to reach one
// of the target statuses, treating NotFound as an empty status and any fail
// state as a terminal error.
func (s *EsaServiceV2) EsaEdgeContainerAppVersionStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEsaEdgeContainerAppVersion(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		if err != nil {
			return nil, "", errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, field, object)
		}
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// PublishEsaEdgeContainerAppVersion publishes the version with the supplied
// publish parameters. The publish API is invoked on Create when publish_env is
// configured and on Update whenever a publish parameter changes.
func (s *EsaServiceV2) PublishEsaEdgeContainerAppVersion(id string, request map[string]interface{}) error {
	return s.PublishEsaEdgeContainerAppVersionContext(context.Background(), id, request)
}

func (s *EsaServiceV2) PublishEsaEdgeContainerAppVersionContext(ctx context.Context, id string, request map[string]interface{}) error {
	appId, versionId, err := EdgeContainerAppVersionId(id)
	if err != nil {
		return err
	}
	client := s.client
	body := map[string]interface{}{
		"AppId":      appId,
		"PublishEnv": request["PublishEnv"],
	}
	if v, ok := request["Percentage"]; ok {
		body["Percentage"] = v
	}
	if v, ok := request["StartTime"]; ok {
		body["StartTime"] = v
	}
	if v, ok := request["Remarks"]; ok {
		body["Remarks"] = v
	}
	query := map[string]interface{}{
		"VersionId": versionId,
		"RegionId":  client.RegionId,
	}
	if v, ok := request["PublishType"]; ok {
		query["PublishType"] = v
	}
	if v, ok := request["Regions"]; ok {
		query["Regions"] = v
	}
	if v, ok := request["FullRelease"]; ok {
		query["FullRelease"] = v
	}
	action := "PublishEdgeContainerAppVersion"

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, edgeContainerAppVersionTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("ESA", "2024-09-10", action, query, body, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, body)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// DeleteEsaEdgeContainerAppVersion deletes the version. A missing version or
// application is treated as already deleted.
func (s *EsaServiceV2) DeleteEsaEdgeContainerAppVersion(id string) error {
	return s.DeleteEsaEdgeContainerAppVersionContext(context.Background(), id)
}

func (s *EsaServiceV2) DeleteEsaEdgeContainerAppVersionContext(ctx context.Context, id string) error {
	appId, versionId, err := EdgeContainerAppVersionId(id)
	if err != nil {
		return err
	}
	client := s.client
	query := map[string]interface{}{
		"AppId":     appId,
		"VersionId": versionId,
		"RegionId":  client.RegionId,
	}
	action := "DeleteEdgeContainerAppVersion"

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, edgeContainerAppVersionTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("ESA", "2024-09-10", action, query, map[string]interface{}{}, true)
		if callErr != nil {
			if errs.IsExpectedErrors(callErr, []string{edgeContainerAppVersionNotFoundErr, edgeContainerAppAppNotFoundErr}) {
				return retry.NonRetryableError(callErr)
			}
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
		if errs.IsExpectedErrors(err, []string{edgeContainerAppVersionNotFoundErr, edgeContainerAppAppNotFoundErr}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// edgeContainerAppVersionJoinId returns the canonical state identity.
func edgeContainerAppVersionJoinId(appId, versionId string) string {
	return fmt.Sprintf("%s:%s", strings.TrimSpace(appId), strings.TrimSpace(versionId))
}
