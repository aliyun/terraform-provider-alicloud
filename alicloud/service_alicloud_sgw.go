package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type SgwService struct {
	client *connectivity.AliyunClient
}

func (s *SgwService) DescribeCloudStorageGatewayStorageBundle(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewHcsSgwClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeStorageBundle"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"StorageBundleId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"StorageBundleNotExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("CloudStorageGatewayStorageBundle", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
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
