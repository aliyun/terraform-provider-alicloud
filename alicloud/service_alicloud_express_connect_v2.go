package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

type ExpressConnectServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeExpressConnectEcFailoverTestJob <<< Encapsulated get interface for ExpressConnect EcFailoverTestJob.

func (s *ExpressConnectServiceV2) DescribeExpressConnectEcFailoverTestJob(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeFailoverTestJob"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["JobId"] = id
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"IllegalParam.JobId"}) {
			return object, WrapErrorf(NotFoundErr("EcFailoverTestJob", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.FailoverTestJobModel", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.FailoverTestJobModel", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ExpressConnectServiceV2) ExpressConnectEcFailoverTestJobStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectEcFailoverTestJob(id)
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

// DescribeExpressConnectEcFailoverTestJob >>> Encapsulated.
// DescribeExpressConnectBgpPeer <<< Encapsulated get interface for ExpressConnect BgpPeer.

func (s *ExpressConnectServiceV2) DescribeExpressConnectBgpPeer(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeBgpPeers"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["BgpPeerId"] = id
	query["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

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

	v, err := jsonpath.Get("$.BgpPeers.BgpPeer[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BgpPeers.BgpPeer[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("BgpPeer", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *ExpressConnectServiceV2) ExpressConnectBgpPeerStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectBgpPeer(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
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

// DescribeExpressConnectBgpPeer >>> Encapsulated.

// DescribeExpressConnectTrafficQos <<< Encapsulated get interface for ExpressConnect TrafficQos.

func (s *ExpressConnectServiceV2) DescribeExpressConnectTrafficQos(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["QosIdList.1"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeExpressConnectTrafficQos"
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	v, err := jsonpath.Get("$.QosList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.QosList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("TrafficQos", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *ExpressConnectServiceV2) ExpressConnectTrafficQosStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectTrafficQos(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectTrafficQos >>> Encapsulated.

// DescribeExpressConnectTrafficQosAssociation <<< Encapsulated get interface for ExpressConnect TrafficQosAssociation.

func (s *ExpressConnectServiceV2) DescribeExpressConnectTrafficQosAssociation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectTrafficQos"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["QosIdList.1"] = parts[0]
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	v, err := jsonpath.Get("$.QosList[*].AssociatedInstanceList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.QosList[*].AssociatedInstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("TrafficQosAssociation", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if item["InstanceId"] != parts[1] {
			continue
		}
		if item["InstanceType"] != parts[2] {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(NotFoundErr("TrafficQosAssociation", id), NotFoundMsg, response)
}

func (s *ExpressConnectServiceV2) ExpressConnectTrafficQosAssociationStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectTrafficQosAssociation(id)
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

func (s *ExpressConnectServiceV2) DescribeAsyncTrafficQosAssociationDescribeExpressConnectTrafficQos(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectTrafficQos"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	return response, nil
}

func (s *ExpressConnectServiceV2) DescribeAsyncExpressConnectTrafficQosAssociationStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncTrafficQosAssociationDescribeExpressConnectTrafficQos(d, res)
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
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectTrafficQosAssociation >>> Encapsulated.

// DescribeExpressConnectTrafficQosQueue <<< Encapsulated get interface for ExpressConnect TrafficQosQueue.

func (s *ExpressConnectServiceV2) DescribeExpressConnectTrafficQosQueue(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeExpressConnectTrafficQosQueue"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["QosId"] = parts[0]
	request["QueueIdList.1"] = parts[1]
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	v, err := jsonpath.Get("$.QueueList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.QueueList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("TrafficQosQueue", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *ExpressConnectServiceV2) ExpressConnectTrafficQosQueueStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectTrafficQosQueue(id)
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

func (s *ExpressConnectServiceV2) DescribeAsyncTrafficQosQueueDescribeExpressConnectTrafficQos(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeExpressConnectTrafficQos"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	return response, nil
}

func (s *ExpressConnectServiceV2) DescribeAsyncExpressConnectTrafficQosQueueStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncTrafficQosQueueDescribeExpressConnectTrafficQos(d, res)
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
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectTrafficQosQueue >>> Encapsulated.

// DescribeExpressConnectTrafficQosRule <<< Encapsulated get interface for ExpressConnect TrafficQosRule.

func (s *ExpressConnectServiceV2) DescribeExpressConnectTrafficQosRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectTrafficQosRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["QosId"] = parts[0]
	query["QueueId"] = parts[1]
	request["RuleIdList.1"] = parts[2]
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	v, err := jsonpath.Get("$.RuleList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RuleList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("TrafficQosRule", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *ExpressConnectServiceV2) ExpressConnectTrafficQosRuleStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectTrafficQosRule(id)
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

func (s *ExpressConnectServiceV2) DescribeAsyncTrafficQosRuleDescribeExpressConnectTrafficQos(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectTrafficQos"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	return response, nil
}

func (s *ExpressConnectServiceV2) DescribeAsyncExpressConnectTrafficQosRuleStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncTrafficQosRuleDescribeExpressConnectTrafficQos(d, res)
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
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectTrafficQosRule >>> Encapsulated.
// DescribeExpressConnectVbrPconnAssociation <<< Encapsulated get interface for ExpressConnect VbrPconnAssociation.

func (s *ExpressConnectServiceV2) DescribeExpressConnectVbrPconnAssociation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "Filter.0.Value.0", parts[0])
	jsonString, _ = sjson.Set(jsonString, "Filter.0.Key", "VbrId")
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return object, WrapError(err)
	}
	request["RegionId"] = client.RegionId
	request["IncludeCrossAccountVbr"] = true

	action := "DescribeVirtualBorderRouters"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

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

	v, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType[*].AssociatedPhysicalConnections.AssociatedPhysicalConnection[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VirtualBorderRouterSet.VirtualBorderRouterType[*].AssociatedPhysicalConnections.AssociatedPhysicalConnection[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("VbrPconnAssociation", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["PhysicalConnectionId"]) != parts[1] {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(NotFoundErr("VbrPconnAssociation", id), NotFoundMsg, response)
}

func (s *ExpressConnectServiceV2) ExpressConnectVbrPconnAssociationStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ExpressConnectVbrPconnAssociationStateRefreshFuncWithApi(id, field, failStates, s.DescribeExpressConnectVbrPconnAssociation)
}

func (s *ExpressConnectServiceV2) ExpressConnectVbrPconnAssociationStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectVbrPconnAssociation >>> Encapsulated.

// DescribeExpressConnectBgpNetwork <<< Encapsulated get interface for ExpressConnect BgpNetwork.

func (s *ExpressConnectServiceV2) DescribeExpressConnectBgpNetwork(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RouterId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["PageNumber"] = 1
	request["PageSize"] = PageSizeLarge

	action := "DescribeBgpNetworks"
	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

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
			if IsExpectedErrors(err, []string{"InvalidResource.NotFound"}) {
				return object, WrapErrorf(NotFoundErr("BgpNetwork", id), NotFoundMsg, response)
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.BgpNetworks.BgpNetwork", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BgpNetworks.BgpNetwork", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("BgpNetwork", id), NotFoundMsg, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["DstCidrBlock"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("BgpNetwork", id), NotFoundMsg, response)
	}

	return object, nil
}

func (s *ExpressConnectServiceV2) ExpressConnectBgpNetworkStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectBgpNetwork(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectBgpNetwork >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for ExpressConnect.
func (s *ExpressConnectServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResourcesForExpressConnect"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = "TagResourcesForExpressConnect"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeExpressConnectBgpGroup <<< Encapsulated get interface for ExpressConnect BgpGroup.

func (s *ExpressConnectServiceV2) DescribeExpressConnectBgpGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BgpGroupId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeBgpGroups"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

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

	v, err := jsonpath.Get("$.BgpGroups.BgpGroup[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BgpGroups.BgpGroup[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("BgpGroup", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *ExpressConnectServiceV2) ExpressConnectBgpGroupStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ExpressConnectBgpGroupStateRefreshFuncWithApi(id, field, failStates, s.DescribeExpressConnectBgpGroup)
}

func (s *ExpressConnectServiceV2) ExpressConnectBgpGroupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectBgpGroup >>> Encapsulated.
// DescribeExpressConnectVbrHa <<< Encapsulated get interface for ExpressConnect VbrHa.

func (s *ExpressConnectServiceV2) DescribeExpressConnectVbrHa(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VbrHaId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeVbrHa"
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if vbrHaId, ok := object["VbrHaId"]; !ok || fmt.Sprint(vbrHaId) == "" {
		return object, WrapErrorf(NotFoundErr("VbrHa", id), NotFoundWithResponse, response)
	}

	return response, nil
}

func (s *ExpressConnectServiceV2) ExpressConnectVbrHaStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ExpressConnectVbrHaStateRefreshFuncWithApi(id, field, failStates, s.DescribeExpressConnectVbrHa)
}

func (s *ExpressConnectServiceV2) ExpressConnectVbrHaStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectVbrHa >>> Encapsulated.

// DescribeExpressConnectRouterInterface <<< Encapsulated get interface for ExpressConnect RouterInterface.

func (s *ExpressConnectServiceV2) DescribeExpressConnectRouterInterface(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeRouterInterfaceAttribute"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

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

	currentStatus := response["RouterInterfaceId"]
	if fmt.Sprint(currentStatus) == "" {
		return object, WrapErrorf(NotFoundErr("RouterInterface", id), NotFoundMsg, response)
	}

	return response, nil
}

func (s *ExpressConnectServiceV2) ExpressConnectRouterInterfaceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ExpressConnectRouterInterfaceStateRefreshFuncWithApi(id, field, failStates, s.DescribeExpressConnectRouterInterface)
}

func (s *ExpressConnectServiceV2) ExpressConnectRouterInterfaceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		object["ChargeType"] = convertExpressConnectRouterInterfaceChargeTypeResponse(object["ChargeType"])
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectRouterInterface >>> Encapsulated.

// DescribeExpressConnectGrantRuleToCen <<< Encapsulated get interface for ExpressConnect GrantRuleToCen.

func (s *ExpressConnectServiceV2) DescribeExpressConnectGrantRuleToCen(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[2]
	request["RegionId"] = client.RegionId
	request["InstanceType"] = "VBR"
	action := "DescribeGrantRulesToCen"
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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

	v, err := jsonpath.Get("$.CenGrantRules.CbnGrantRule[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CenGrantRules.CbnGrantRule[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("GrantRuleToCen", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["CenInstanceId"]) != parts[0] {
			continue
		}
		if fmt.Sprint(item["CenOwnerId"]) != parts[1] {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(NotFoundErr("GrantRuleToCen", id), NotFoundMsg, response)
}

func (s *ExpressConnectServiceV2) ExpressConnectGrantRuleToCenStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ExpressConnectGrantRuleToCenStateRefreshFuncWithApi(id, field, failStates, s.DescribeExpressConnectGrantRuleToCen)
}

func (s *ExpressConnectServiceV2) ExpressConnectGrantRuleToCenStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectGrantRuleToCen >>> Encapsulated.

// DescribeExpressConnectVirtualBorderRouter <<< Encapsulated get interface for ExpressConnect VirtualBorderRouter.

func (s *ExpressConnectServiceV2) DescribeExpressConnectVirtualBorderRouter(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["IncludeCrossAccountVbr"] = "true"
	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "Filter.0.Key", "VbrId")
	jsonString, _ = sjson.Set(jsonString, "Filter.0.Value.0", id)
	_ = json.Unmarshal([]byte(jsonString), &request)

	action := "DescribeVirtualBorderRouters"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

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

	v, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VirtualBorderRouterSet.VirtualBorderRouterType[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("VirtualBorderRouter", id), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["VbrId"]
	if fmt.Sprint(currentStatus) == "" {
		return object, WrapErrorf(NotFoundErr("VirtualBorderRouter", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *ExpressConnectServiceV2) ExpressConnectVirtualBorderRouterStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ExpressConnectVirtualBorderRouterStateRefreshFuncWithApi(id, field, failStates, s.DescribeExpressConnectVirtualBorderRouter)
}

func (s *ExpressConnectServiceV2) ExpressConnectVirtualBorderRouterStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectVirtualBorderRouter >>> Encapsulated.
