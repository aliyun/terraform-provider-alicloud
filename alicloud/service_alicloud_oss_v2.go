package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type OssServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeOssBucketReferer <<< Encapsulated get interface for Oss BucketReferer.

func (s *OssServiceV2) DescribeOssBucketReferer(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var hostMap map[string]*string
	action := fmt.Sprintf("/?referer")
	conn, err := client.NewOssClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap = make(map[string]*string)
	hostMap["bucket"] = StringPointer(id)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.Execute(genXmlParam("GetBucketReferer", "GET", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: nil, HostMap: hostMap}, &util.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.body", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RefererConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *OssServiceV2) OssBucketRefererStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOssBucketReferer(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOssBucketReferer >>> Encapsulated.
