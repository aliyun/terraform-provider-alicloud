package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type EipanycastService struct {
	client *connectivity.AliyunClient
}

func (s *EipanycastService) DescribeEipanycastAnycastEipAddress(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewEipanycastClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeAnycastEipAddress"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"AnycastId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-09"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *EipanycastService) EipanycastAnycastEipAddressStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEipanycastAnycastEipAddress(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}
