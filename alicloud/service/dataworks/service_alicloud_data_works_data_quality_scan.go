package dataworks

import (
	"context"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// dataQualityScanNotFoundCodes are the response codes the dataworks-public
// 2024-05-18 product returns for a missing entity (seen on Project and
// DwResourceGroup reads). GetDataQualityScan shares the product error
// taxonomy, so these codes are treated as NotFound.
var dataQualityScanNotFoundCodes = []string{"9990040003"}

// DescribeDataQualityScanContext calls GetDataQualityScan and returns the
// DataQualityScan object. A confirmed missing scan (a not-found response code
// or an empty/absent DataQualityScan object) is reported as a NotFound error
// so the Framework Read path can remove the resource from state.
func (s *Service) DescribeDataQualityScanContext(ctx context.Context, id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "GetDataQualityScan"
	request := map[string]interface{}{
		"Id":       id,
		"RegionId": client.RegionId,
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, dataQualityScanReadTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
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
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, dataQualityScanNotFoundCodes) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DataQualityScan", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DataQualityScan", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("DataQualityScan", id), errs.NotFoundMsg, response)
	}
	obj, ok := v.(map[string]interface{})
	if !ok || len(obj) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DataQualityScan", id), errs.NotFoundMsg, response)
	}
	return obj, nil
}

// DeleteDataQualityScanContext deletes the scan by Id. A not-found response is
// treated as success so the delete is idempotent.
func (s *Service) DeleteDataQualityScanContext(ctx context.Context, id string) error {
	client := s.client
	action := "DeleteDataQualityScan"
	request := map[string]interface{}{
		"Id":       id,
		"RegionId": client.RegionId,
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, dataQualityScanTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = client.RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
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
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, dataQualityScanNotFoundCodes) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
