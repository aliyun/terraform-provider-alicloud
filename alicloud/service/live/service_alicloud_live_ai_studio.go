package live

import (
	"context"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const liveAIStudioReadTimeout = 1 * time.Minute

// DescribeLiveAIStudio returns the AI studio whose StudioId (returned by
// CreateLiveAIStudio and named TemplateId in the DescribeLiveAIStudio response)
// equals id. DescribeLiveAIStudio is a collection API: StudioId is an optional
// filter, so the target is selected by its exact id instead of the first item.
func (s *Service) DescribeLiveAIStudio(id string) (object map[string]interface{}, err error) {
	return s.DescribeLiveAIStudioContext(context.Background(), id)
}

func (s *Service) DescribeLiveAIStudioContext(ctx context.Context, id string) (object map[string]interface{}, err error) {
	client := s.client
	// The tea RPC client sends the 5th argument (request) as the request body;
	// the 4th (query) is kept empty to match the existing live service convention
	// (see DescribeLiveDomainDetail).
	query := map[string]interface{}{}
	request := map[string]interface{}{
		"StudioId":   id,
		"PageNumber": 1,
		"PageSize":   helper.PageSizeLarge,
		"RegionId":   client.RegionId,
	}
	action := "DescribeLiveAIStudio"

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, liveAIStudioReadTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
		// A transport-level 404 is definitive; a body-level 404 is classified
		// after the container is validated below, so a malformed body cannot be
		// turned into NotFound.
		if errs.IsTypedNotFound(err) && !errs.IsResponseCodeNotFound(err) {
			return object, errs.WrapErrorf(errs.NotFoundErr("LiveAIStudio", id), errs.NotFoundMsg, response)
		}
		if !errs.IsResponseCodeNotFound(err) {
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	// The response nests the items under StudioConfigs.SubtitleConfig. Validate
	// the container before any NotFound classification: a null or scalar shape
	// is a malformed response, not evidence that the studio is gone.
	configsRaw, gerr := jsonpath.Get("$.StudioConfigs", response)
	if gerr != nil {
		return object, errs.WrapErrorf(gerr, errs.FailedGetAttributeMsg, id, "$.StudioConfigs", response)
	}
	configs, ok := configsRaw.(map[string]interface{})
	if !ok {
		// StudioConfigs is absent when the filter matches nothing; treat that as
		// NotFound only when it is not a scalar drift (e.g. a stray string).
		if configsRaw == nil {
			return object, errs.WrapErrorf(errs.NotFoundErr("LiveAIStudio", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(fmt.Errorf("unexpected StudioConfigs response type %T", configsRaw), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	raw, gerr := jsonpath.Get("$.SubtitleConfig", configs)
	if gerr != nil {
		return object, errs.WrapErrorf(gerr, errs.FailedGetAttributeMsg, id, "$.StudioConfigs.SubtitleConfig", response)
	}
	list, ok := raw.([]interface{})
	if !ok {
		if raw == nil {
			return object, errs.WrapErrorf(errs.NotFoundErr("LiveAIStudio", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(fmt.Errorf("unexpected SubtitleConfig response type %T", raw), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if len(list) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("LiveAIStudio", id), errs.NotFoundMsg, response)
	}

	for _, item := range list {
		studio, ok := item.(map[string]interface{})
		if !ok {
			return object, errs.WrapErrorf(fmt.Errorf("unexpected SubtitleConfig element type %T", item), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		// CreateLiveAIStudio returns StudioId; DescribeLiveAIStudio echoes it as
		// TemplateId. The two name the same resource identity.
		if templateId, _ := studio["TemplateId"].(string); templateId == id {
			return studio, nil
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("LiveAIStudio", id), errs.NotFoundMsg, response)
}

// DeleteLiveAIStudio deletes the AI studio identified by StudioId. A typed 404
// means the studio is already gone.
func (s *Service) DeleteLiveAIStudio(id string) error {
	return s.DeleteLiveAIStudioContext(context.Background(), id)
}

func (s *Service) DeleteLiveAIStudioContext(ctx context.Context, id string) error {
	client := s.client
	action := "DeleteLiveAIStudio"
	request := map[string]interface{}{
		"StudioId": id,
		"RegionId": client.RegionId,
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, liveAIStudioReadTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("live", "2016-11-01", action, map[string]interface{}{}, request, true)
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
