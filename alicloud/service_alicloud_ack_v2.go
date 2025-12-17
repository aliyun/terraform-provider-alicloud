package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type AckServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeAckNodepool <<< Encapsulated get interface for Ack Nodepool.

func (s *AckServiceV2) DescribeAckNodepool(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	ClusterId := parts[0]
	NodepoolId := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)

	action := fmt.Sprintf("/clusters/%s/nodepools/%s", ClusterId, NodepoolId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("CS", "2015-12-15", action, query, header, nil)

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
		if IsExpectedErrors(err, []string{"ErrorNodePoolNotFound", "ErrorClusterNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Nodepool", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *AckServiceV2) AckNodepoolStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.AckNodepoolStateRefreshFuncWithApi(id, field, failStates, s.DescribeAckNodepool)
}

func (s *AckServiceV2) AckNodepoolStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

func (s *AckServiceV2) DescribeAsyncAckNodepoolStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeTaskInfo(d, res)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
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
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAckNodepool >>> Encapsulated.

// DescribeAsyncDescribeTaskInfo <<< Encapsulated for Ack.
func (s *AckServiceV2) DescribeAsyncDescribeTaskInfo(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	task_id, err := jsonpath.Get("$.task_id", res)

	action := fmt.Sprintf("/tasks/%s", task_id)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("CS", "2015-12-15", action, query, header, nil)

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
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeAsyncDescribeTaskInfo >>> Encapsulated.

// DescribeAckPolicyInstance <<< Encapsulated get interface for Ack PolicyInstance.
func (s *AckServiceV2) DescribeAckPolicyInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	// Using the official SDK method DescribePolicyInstances
	csClient, clientErr := client.NewRoaCsClient()
	if clientErr != nil {
		return object, WrapError(clientErr)
	}

	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
		return object, err
	}

	cluster_id := parts[0]
	describeReq := &cs.DescribePolicyInstancesRequest{
		PolicyName:   tea.String(parts[1]),
		InstanceName: tea.String(parts[2]),
	}
	action := "DescribePolicyInstances"
	var response map[string]interface{}
	var rawResponse interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		// Call the official DescribePolicyInstances method
		var resp *cs.DescribePolicyInstancesResponse
		resp, err = csClient.DescribePolicyInstances(tea.String(cluster_id), describeReq)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		// Convert response to map[string]interface{} for compatibility
		respMap := make(map[string]interface{})
		if resp.Body != nil {
			byt, _ := json.Marshal(resp.Body)
			d := json.NewDecoder(strings.NewReader(string(byt)))
			d.UseNumber()
			d.Decode(&respMap)
		}
		response = respMap
		rawResponse = resp.Body

		return nil
	})
	addDebug(action, response, describeReq)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidResource.NotFound", "404"}) {
			return object, WrapErrorf(NotFoundErr("PolicyInstance", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	// Handle the case where response itself is an array (not wrapped in body)
	// This happens when the API returns a direct array instead of an object with a body field
	if rawResponse != nil {
		// Try to marshal and unmarshal to handle different response formats
		b, _ := json.Marshal(rawResponse)
		var directArray []interface{}
		if json.Unmarshal(b, &directArray) == nil && len(directArray) > 0 {
			if item, ok := directArray[0].(map[string]interface{}); ok {
				return item, nil
			}
		}
	}

	// Handle the case where response is an array wrapped in body field
	if bodyArray, ok := response["body"].([]interface{}); ok && len(bodyArray) > 0 {
		// Return first element if it's a map
		if firstElement, ok := bodyArray[0].(map[string]interface{}); ok {
			return firstElement, nil
		}
		return object, WrapErrorf(fmt.Errorf("unexpected response format: first element is not a map"), FailedGetAttributeMsg, id, "$[*]", response)
	}

	// Handle the case where response is a map with items in body field
	if bodyMap, ok := response["body"].(map[string]interface{}); ok {
		if itemsArray, ok := bodyMap["items"].([]interface{}); ok && len(itemsArray) > 0 {
			if firstElement, ok := itemsArray[0].(map[string]interface{}); ok {
				return firstElement, nil
			}
		}
	}

	// Fallback to original logic with jsonpath for backward compatibility
	v, err := jsonpath.Get("$[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$[*]", response)
	}

	// Check if the result is a slice and has at least one element
	resultSlice, ok := v.([]interface{})
	if !ok || len(resultSlice) == 0 {
		return object, WrapErrorf(NotFoundErr("PolicyInstance", id), NotFoundMsg, response)
	}

	// Check if the first element is a map
	firstElement, ok := resultSlice[0].(map[string]interface{})
	if !ok {
		return object, WrapErrorf(fmt.Errorf("unexpected response format: first element is not a map"), FailedGetAttributeMsg, id, "$[*]", response)
	}

	return firstElement, nil
}

// DescribeAckPolicyInstance >>> Encapsulated.
