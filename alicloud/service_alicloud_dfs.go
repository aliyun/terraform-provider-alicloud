package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DfsService struct {
	client *connectivity.AliyunClient
}

func (s *DfsService) DescribeDfsAccessGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAlidfsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetAccessGroup"
	request := map[string]interface{}{
		"InputRegionId": s.client.RegionId,
		"AccessGroupId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameter.AccessGroupNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("DFS:AccessGroup", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AccessGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AccessGroup", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DfsService) GetAccessGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAlidfsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetAccessGroup"
	request := map[string]interface{}{
		"InputRegionId": s.client.RegionId,
		"AccessGroupId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameter.AccessGroupNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("DFS:AccessGroup", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AccessGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AccessGroup", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DfsService) DescribeDfsFileSystem(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAlidfsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetFileSystem"
	request := map[string]interface{}{
		"InputRegionId": s.client.RegionId,
		"FileSystemId":  id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameter.FileSystemNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("DFS:FileSystem", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.FileSystem", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.FileSystem", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
