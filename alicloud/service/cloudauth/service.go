package cloudauth

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the Cloudauth API calls used by the resources and data
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

func (s *Service) DescribeFaceConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeFaceConfig"
	request := map[string]interface{}{
		"Lang": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cloudauth", "2019-03-07", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cloudauth", id), errs.NotFoundWithResponse, response)
	} else {
		for _, obj := range v.([]interface{}) {
			if fmt.Sprint(obj.(map[string]interface{})["BizType"]) == id {
				object = v.([]interface{})[0].(map[string]interface{})
				return object, nil
			}
		}
		return object, errs.WrapErrorf(errs.NotFoundErr("Cloudauth", id), errs.NotFoundWithResponse, response)
	}
}
