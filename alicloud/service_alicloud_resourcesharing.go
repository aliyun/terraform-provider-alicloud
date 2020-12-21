package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ResourcesharingService struct {
	client *connectivity.AliyunClient
}

func (s *ResourcesharingService) DescribeResourceManagerResourceShare(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewRessharingClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListResourceShares"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"ResourceShareIds": []string{id},
		"ResourceOwner":    "Self",
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-10"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ResourceShares", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ResourceShares", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ResourceManager", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["ResourceShareId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("ResourceManager", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ResourcesharingService) ResourceManagerResourceShareStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerResourceShare(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ResourceShareStatus"].(string) == failState {
				return object, object["ResourceShareStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ResourceShareStatus"].(string)))
			}
		}
		return object, object["ResourceShareStatus"].(string), nil
	}
}
