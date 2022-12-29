package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type EbsService struct {
	client *connectivity.AliyunClient
}

func (s *EbsService) DescribeEbsDiskReplicaGroup(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewEbsClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"GroupIds": id,
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeDiskReplicaGroups"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-07-30"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NoSuchResource"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ReplicaGroups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ReplicaGroups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("EBS", id)), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *EbsService) EbsDiskReplicaGroupStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEbsDiskReplicaGroup(d.Id())
		if err != nil {
			if NotFoundError(err) {
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
func (s *EbsService) DescribeEbsDedicatedBlockStorageCluster(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewEbsClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"DedicatedBlockStorageClusterId": []string{id},
		"RegionId":                       s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeDedicatedBlockStorageClusters"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-07-30"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DedicatedBlockStorageClusters", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DedicatedBlockStorageClusters", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("EBS", id)), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *EbsService) EbsDedicatedBlockStorageClusterStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEbsDedicatedBlockStorageCluster(d.Id())
		if err != nil {
			if NotFoundError(err) {
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
