package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type EhpcService struct {
	client *connectivity.AliyunClient
}

func (s *EhpcService) DescribeEhpcJobTemplate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewEhpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListJobTemplates"
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-04-12"), StringPointer("AK"), request, nil, &runtime)
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
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Templates.JobTemplates", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Templates.JobTemplates", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Ehpc", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Id"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Ehpc", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *EhpcService) DescribeEhpcCluster(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewEhsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeCluster"
	request := map[string]interface{}{
		"ClusterId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-04-12"), StringPointer("AK"), request, nil, &runtime)
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
		if IsExpectedErrors(err, []string{"ClusterNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Ehpc:Cluster", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ClusterInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ClusterInfo", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *EhpcService) EhpcClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEhpcCluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}
