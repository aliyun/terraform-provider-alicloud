package expressconnectrouter

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeExpressConnectRouterExpressConnectRouter <<< Encapsulated get interface for ExpressConnectRouter ExpressConnectRouter.

func (s *Service) DescribeExpressConnectRouterExpressConnectRouter(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeExpressConnectRouter"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.EcrList[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouter", id), errs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouter", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["AlibabaSideAsn"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouter", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *Service) DescribeDescribeInstanceGrantedToExpressConnectRouter(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstanceGrantedToExpressConnectRouter"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouter", id), errs.NotFoundMsg, response)
	}

	return response, nil
}
func (s *Service) DescribeDescribeExpressConnectRouterRouteEntries(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeExpressConnectRouterRouteEntries"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouter", id), errs.NotFoundMsg, response)
	}

	return response, nil
}
func (s *Service) DescribeDescribeDisabledExpressConnectRouterRouteEntries(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeDisabledExpressConnectRouterRouteEntries"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouter", id), errs.NotFoundMsg, response)
	}

	return response, nil
}
func (s *Service) DescribeDescribeExpressConnectRouterInterRegionTransitMode(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeExpressConnectRouterInterRegionTransitMode"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouter", id), errs.NotFoundMsg, response)
	}

	return response, nil
}
func (s *Service) DescribeDescribeExpressConnectRouterRegion(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeExpressConnectRouterRegion"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouter", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *Service) ExpressConnectRouterExpressConnectRouterStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterExpressConnectRouter(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func (s *Service) DescribeAsyncDescribeInstanceGrantedToExpressConnectRouter(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstanceGrantedToExpressConnectRouter"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["EcrId"] = id
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.EcrId"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouter", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *Service) DescribeAsyncExpressConnectRouterExpressConnectRouterStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeInstanceGrantedToExpressConnectRouter(d, res)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}
func (s *Service) DescribeAsyncGetExpressConnectRouter(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetExpressConnectRouter"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["EcrId"] = id
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeExpressConnectRouterExpressConnectRouter >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for ExpressConnectRouter.
func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["ClientToken"] = helper.BuildClientToken(action)
			request["ResourceType"] = resourceType
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = "TagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["ClientToken"] = helper.BuildClientToken(action)
			request["ResourceType"] = resourceType
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance <<< Encapsulated get interface for ExpressConnectRouter ExpressConnectRouterVbrChildInstance.

func (s *Service) DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterChildInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ChildInstanceId"] = parts[1]
	request["ChildInstanceType"] = parts[2]
	request["EcrId"] = parts[0]

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

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
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.EcrId"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterVbrChildInstance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ChildInstanceList[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ChildInstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterVbrChildInstance", id), errs.NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["ChildInstanceId"]) != parts[1] {
			continue
		}
		if fmt.Sprint(item["ChildInstanceType"]) != parts[2] {
			continue
		}
		if fmt.Sprint(item["ChildInstanceType"]) != "VBR" {
			continue
		}
		if fmt.Sprint(item["EcrId"]) != parts[0] {
			continue
		}
		return item, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterVbrChildInstance", id), errs.NotFoundMsg, response)
}

func (s *Service) ExpressConnectRouterExpressConnectRouterVbrChildInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func (s *Service) DescribeAsyncExpressConnectRouterVbrChildInstanceDescribeExpressConnectRouterChildInstance(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterChildInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["ChildInstanceId"] = parts[1]

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

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
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *Service) DescribeAsyncExpressConnectRouterExpressConnectRouterVbrChildInstanceStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncExpressConnectRouterVbrChildInstanceDescribeExpressConnectRouterChildInstance(d, res)
		if err != nil {
			if errs.NotFoundError(err) {
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
					return _err, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance >>> Encapsulated.
// DescribeExpressConnectRouterExpressConnectRouterTrAssociation <<< Encapsulated get interface for ExpressConnectRouter ExpressConnectRouterTrAssociation.

func (s *Service) DescribeExpressConnectRouterExpressConnectRouterTrAssociation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterAssociation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AssociationId"] = parts[1]
	request["EcrId"] = parts[0]
	request["TransitRouterId"] = parts[2]

	request["ClientToken"] = helper.BuildClientToken(action)

	request["AssociationNodeType"] = "TR"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterTrAssociation", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.AssociationList[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AssociationList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterTrAssociation", id), errs.NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if item["AssociationId"] != parts[1] {
			continue
		}
		if item["AssociationNodeType"] != "TR" {
			continue
		}
		if item["EcrId"] != parts[0] {
			continue
		}
		if item["TransitRouterId"] != parts[2] {
			continue
		}
		return item, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterTrAssociation", id), errs.NotFoundMsg, response)
}

func (s *Service) ExpressConnectRouterExpressConnectRouterTrAssociationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterExpressConnectRouterTrAssociation(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectRouterExpressConnectRouterTrAssociation >>> Encapsulated.
// DescribeExpressConnectRouterExpressConnectRouterVpcAssociation <<< Encapsulated get interface for ExpressConnectRouter ExpressConnectRouterVpcAssociation.

func (s *Service) DescribeExpressConnectRouterExpressConnectRouterVpcAssociation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterAssociation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AssociationId"] = parts[1]
	request["EcrId"] = parts[0]
	request["VpcId"] = parts[2]

	request["ClientToken"] = helper.BuildClientToken(action)

	request["AssociationNodeType"] = "VPC"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterVpcAssociation", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.AssociationList[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AssociationList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterVpcAssociation", id), errs.NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if item["AssociationId"] != parts[1] {
			continue
		}
		if item["AssociationNodeType"] != "VPC" {
			continue
		}
		if item["EcrId"] != parts[0] {
			continue
		}
		if item["VpcId"] != parts[2] {
			continue
		}
		return item, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterVpcAssociation", id), errs.NotFoundMsg, response)
}
func (s *Service) DescribeDescribeExpressConnectRouterAllowedPrefixHistory(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterAllowedPrefixHistory"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["InstanceId"] = parts[2]

	request["ClientToken"] = helper.BuildClientToken(action)

	request["InstanceType"] = "VPC"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.EcrId", "ResourceNotFound.AssociationId", "ResourceNotFound.InstanceId"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnectRouterVpcAssociation", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.AllowedPrefixHistoryList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AllowedPrefixHistoryList", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterExpressConnectRouterVpcAssociation(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectRouterExpressConnectRouterVpcAssociation >>> Encapsulated.

// DescribeExpressConnectRouterGrantAssociation <<< Encapsulated get interface for ExpressConnectRouter GrantAssociation.

func (s *Service) DescribeExpressConnectRouterGrantAssociation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 5 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 5, len(parts)))
	}
	action := "DescribeInstanceGrantedToExpressConnectRouter"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["InstanceId"] = parts[1]
	request["InstanceRegionId"] = parts[2]
	request["InstanceType"] = parts[4]

	request["ClientToken"] = helper.BuildClientToken(action)

	request["CallerType"] = "OTHER"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

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

	v, err := jsonpath.Get("$.EcrGrantedInstanceList[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.EcrGrantedInstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("GrantAssociation", id), errs.NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["EcrId"]) != parts[0] {
			continue
		}
		if fmt.Sprint(item["EcrOwnerAliUid"]) != parts[3] {
			continue
		}
		if fmt.Sprint(item["NodeId"]) != parts[1] {
			continue
		}
		if fmt.Sprint(item["NodeRegionId"]) != parts[2] {
			continue
		}
		if fmt.Sprint(item["NodeType"]) != parts[4] {
			continue
		}
		return item, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("GrantAssociation", id), errs.NotFoundMsg, response)
}

func (s *Service) ExpressConnectRouterGrantAssociationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterGrantAssociation(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectRouterGrantAssociation >>> Encapsulated.
