package wafopenapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Wafv3ServiceV2 struct {
	client *connectivity.AliyunClient
}

func NewWafv3ServiceV2(client *connectivity.AliyunClient) *Wafv3ServiceV2 {
	return &Wafv3ServiceV2{client: client}
}

// DescribeWafv3DefenseTemplate <<< Encapsulated get interface for Wafv3 DefenseTemplate.

func (s *Wafv3ServiceV2) DescribeWafv3DefenseTemplate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TemplateId"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	action := "DescribeDefenseTemplate"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

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

	v, err := jsonpath.Get("$.Template", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("DefenseTemplate", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["TemplateName"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("DefenseTemplate", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}
func (s *Wafv3ServiceV2) DescribeDefenseTemplateDescribeTemplateResources(id string, resourceType string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TemplateId"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ResourceType"] = resourceType
	action := "DescribeTemplateResources"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"Defense.Control.DefenseTemplateNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DefenseTemplate", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *Wafv3ServiceV2) Wafv3DefenseTemplateStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.Wafv3DefenseTemplateStateRefreshFuncWithApi(id, field, failStates, s.DescribeWafv3DefenseTemplate)
}

func (s *Wafv3ServiceV2) Wafv3DefenseTemplateStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeWafv3DefenseTemplate >>> Encapsulated.

// DescribeWafv3Domain <<< Encapsulated get interface for Wafv3 Domain.

func (s *Wafv3ServiceV2) DescribeWafv3Domain(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	action := "DescribeDomainDetail"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"Waf.Instance.ValidFaild"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Domain", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	currentStatus := response["Status"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Domain", id), errs.NotFoundMsg, response)
	}

	return response, nil
}
func (s *Wafv3ServiceV2) DescribeDomainListTagResources(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId.1"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ResourceType"] = "ALIYUN::WAF::DEFENSERESOURCE"
	action := "ListTagResources"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

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

	return response, nil
}

func (s *Wafv3ServiceV2) Wafv3DomainStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.Wafv3DomainStateRefreshFuncWithApi(id, field, failStates, s.DescribeWafv3Domain)
}

func (s *Wafv3ServiceV2) Wafv3DomainStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeWafv3Domain >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Wafv3.
func (s *Wafv3ServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
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
			parts := strings.Split(d.Id(), ":")
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = parts[1]
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			parts := strings.Split(d.Id(), ":")
			action = "TagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = parts[1]
			request["RegionId"] = client.RegionId
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeWafv3DefenseRule <<< Encapsulated get interface for Wafv3 DefenseRule.

func (s *Wafv3ServiceV2) DescribeWafv3DefenseRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DefenseType"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RuleId"] = parts[2]
	request["RegionId"] = client.RegionId
	action := "DescribeDefenseRule"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"Defense.Control.DefenseRuleNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DefenseRule", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Rule", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Rule", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Wafv3ServiceV2) Wafv3DefenseRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.Wafv3DefenseRuleStateRefreshFuncWithApi(id, field, failStates, s.DescribeWafv3DefenseRule)
}

func (s *Wafv3ServiceV2) Wafv3DefenseRuleStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeWafv3DefenseRule >>> Encapsulated.

// DescribeWafv3DefenseResourceGroup <<< Encapsulated get interface for Wafv3 DefenseResourceGroup.

func (s *Wafv3ServiceV2) DescribeWafv3DefenseResourceGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["GroupName"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	action := "DescribeDefenseResourceGroup"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

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

	v, err := jsonpath.Get("$.Group", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("DefenseResourceGroup", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["GroupName"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("DefenseResourceGroup", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Wafv3ServiceV2) Wafv3DefenseResourceGroupStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.Wafv3DefenseResourceGroupStateRefreshFuncWithApi(id, field, failStates, s.DescribeWafv3DefenseResourceGroup)
}

func (s *Wafv3ServiceV2) Wafv3DefenseResourceGroupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		if field == "$.ResourceList" {
			e := jsonata.MustCompile("$split($.ResourceList, \",\")")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}

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

// DescribeWafv3DefenseResourceGroup >>> Encapsulated.

// DescribeWafv3AddressBook <<< Encapsulated get interface for Wafv3 AddressBook.

func (s *Wafv3ServiceV2) DescribeWafv3AddressBook(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RuleId"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["DefenseType"] = "global"
	action := "DescribeDefenseRule"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"Defense.Control.DefenseRuleNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("AddressBook", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Rule.Config", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Rule.Config", response)
	}

	if configStr, ok := v.(string); ok {
		object = make(map[string]interface{})
		if err := json.Unmarshal([]byte(configStr), &object); err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Rule.Config", response)
		}
		return object, nil
	}
	return v.(map[string]interface{}), nil
}
func (s *Wafv3ServiceV2) DescribeWafv3AddressBookAddresses(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}

	action := "DescribeAddresses"
	aggregatedAddresses := make([]interface{}, 0)
	nextToken := ""

	for {
		request := map[string]interface{}{
			"RuleId":     parts[1],
			"InstanceId": parts[0],
			"RegionId":   client.RegionId,
			"MaxResults": 100,
		}
		if nextToken != "" {
			request["NextToken"] = nextToken
		}
		query := make(map[string]interface{})

		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
			if errs.IsExpectedErrors(err, []string{"Defense.Control.DefenseRuleNotExist"}) {
				return object, errs.WrapErrorf(errs.NotFoundErr("AddressBook", id), errs.NotFoundMsg, response)
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}

		if pageList, ok := response["AddressList"].([]interface{}); ok {
			aggregatedAddresses = append(aggregatedAddresses, pageList...)
		}

		nt, _ := response["NextToken"].(string)
		if nt == "" || nt == nextToken {
			break
		}
		nextToken = nt
	}

	return map[string]interface{}{"AddressList": aggregatedAddresses}, nil
}

// DescribeWafv3AddressBook >>> Encapsulated.
