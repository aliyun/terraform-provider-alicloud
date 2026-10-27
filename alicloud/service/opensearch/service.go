package opensearch

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the OpenSearch API calls used by the resources and data
// sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the flat
// layout's acceptance tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeAppGroup <<< Encapsulated get interface for OpenSearch AppGroup.
func (s *Service) DescribeAppGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/v4/openapi/app-groups/" + id
	body := map[string]interface{}{
		"appGroupIdentity": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("OpenSearch", "2017-12-25", action, nil, nil, body)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug("GET "+action, response, body)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"App.NotFound"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.result", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

// AppStateRefreshFunc <<< Encapsulated the WaitForState method for OpenSearch AppGroup.
func (s *Service) AppStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAppGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["status"].(string) == failState {
				return object, object["status"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["status"].(string)))
			}
		}
		return object, object["status"].(string), nil
	}
}

func convertAppGroupPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	case "Subscription":
		return "PREPAY"
	}
	return source
}

func convertAppGroupPaymentTypeResponse(source string) string {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}

func convertChargingWayRequest(source string) interface{} {
	switch source {
	case "compute_resource":
		return 1
	case "qps":
		return 2
	}
	return source
}

func convertChargingWayResponse(source string) interface{} {
	switch source {
	case "1":
		return "compute_resource"
	case "2":
		return "qps"
	}
	return source
}
