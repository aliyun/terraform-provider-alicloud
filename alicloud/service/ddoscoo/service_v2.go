package ddoscoo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

// DescribeDdosCooPort <<< Encapsulated get interface for DdosCoo Port.

func (s *Service) DescribeDdosCooPort(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribePort"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FrontendPort"] = parts[1]
	query["FrontendProtocol"] = parts[2]
	query["InstanceId"] = parts[0]

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"anycast_controller3006"}) || helper.NeedRetry(err) {
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

	v, err := jsonpath.Get("$.NetworkRules[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.NetworkRules[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Port", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *Service) DescribeDescribeNetworkRuleAttributes(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeNetworkRuleAttributes"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	jsonString := "{}"
	frontendPort, _ := strconv.ParseInt(parts[1], 10, 64)
	jsonString, _ = sjson.Set(jsonString, "NetworkRules.0.FrontendPort", frontendPort)
	jsonString, _ = sjson.Set(jsonString, "NetworkRules.0.Protocol", parts[2])
	jsonString, _ = sjson.Set(jsonString, "NetworkRules.0.InstanceId", parts[0])
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return object, errs.WrapError(err)
	}
	request["NetworkRules"] = convert.ConvertObjectToJsonString(request["NetworkRules"])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"anycast_controller3006"}) || helper.NeedRetry(err) {
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

	v, err := jsonpath.Get("$.NetworkRuleAttributes[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.NetworkRuleAttributes[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Port", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) DdosCooPortStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDdosCooPort(id)
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

// DescribeDdosCooPort >>> Encapsulated.

// DescribeDdosCooDomainResource <<< Encapsulated get interface for DdosCoo DomainResource.

func (s *Service) DescribeDdosCooDomainResource(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = id
	action := "DescribeDomainResource"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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

	v, err := jsonpath.Get("$.WebRules[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.WebRules[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DomainResource", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *Service) DescribeDomainResourceDescribeWebRules(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = id
	action := "DescribeWebRules"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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

	v, err := jsonpath.Get("$.WebRules[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.WebRules[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DomainResource", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *Service) DescribeDomainResourceDescribeDomainCcProtectSwitch(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domains.1"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeDomainCcProtectSwitch"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"DomainNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DomainResource", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ProtectSwitchList[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ProtectSwitchList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DomainResource", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *Service) DescribeDomainResourceDescribeHeaders(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Domain"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeHeaders"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"DomainNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DomainResource", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.CustomHeader", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CustomHeader", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) DdosCooDomainResourceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.DdosCooDomainResourceStateRefreshFuncWithApi(id, field, failStates, s.DescribeDdosCooDomainResource)
}

func (s *Service) DdosCooDomainResourceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		object["OcspEnabled"] = convertDdosCooDomainResourceWebRulesOcspEnabledResponse(object["OcspEnabled"])
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

// DescribeDdosCooDomainResource >>> Encapsulated.

// DescribeDdosCooInstance <<< Encapsulated get interface for DdosCoo Instance.

func (s *Service) DescribeDdosCooInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIds.1"] = id
	request["PageSize"] = helper.PageSizeSmall
	request["PageNumber"] = 1
	action := "DescribeInstances"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Instances[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Instances[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) DescribeInstanceDescribeInstanceSpecs(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIds.1"] = id
	action := "DescribeInstanceSpecs"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidInstanceId"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InstanceNotFound", "ddos_coop3301"}) || errs.NotFoundError(err) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.InstanceSpecs[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.InstanceSpecs[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) DescribeInstanceDescribeInstanceExt(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id

	action := "DescribeInstanceExt"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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

	v, err := jsonpath.Get("$.InstanceExtSpecs[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.InstanceExtSpecs[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) DescribeInstanceDescribeTagResources(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceIds.1"] = id
	request["RegionId"] = client.RegionId
	request["ResourceType"] = "INSTANCE"
	action := "DescribeTagResources"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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

func (s *Service) ModifyDdosCooInstance(schemaName string, specName string, d *schema.ResourceData) (err error) {
	client := s.client
	var response map[string]interface{}
	action := "ModifyInstance"
	request := map[string]interface{}{
		"InstanceId":  d.Id(),
		"ProductType": "ddoscoo",
		"ModifyType":  "UPGRADE",
	}

	if v, ok := d.GetOk("product_type"); ok {
		request["ProductType"] = v
	}

	o, n := d.GetChange(schemaName)
	oi, _ := strconv.Atoi(o.(string))
	ni, _ := strconv.Atoi(n.(string))
	if ni < oi {
		request["ModifyType"] = "DOWNGRADE"
	}

	if v, ok := d.GetOk(schemaName); ok {
		request[specName] = v
	}

	if v, ok := d.GetOk("modify_type"); ok {
		request["ModifyType"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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

	return nil
}

func (s *Service) DdosCooInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDdosCooInstance(id)
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

// DescribeDdosCooInstance >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for DdosCoo.
func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
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
			action = "DeleteTagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceIds.1"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ResourceType"] = resourceType
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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
			action = "CreateTagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceIds.1"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ResourceType"] = resourceType
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tags.%d.Key", count)] = key
				request[fmt.Sprintf("Tags.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)
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

// DescribeDdosCooDomainPreciseAccessRule <<< Encapsulated get interface for DdosCoo DomainPreciseAccessRule.

func (s *Service) DescribeDdosCooDomainPreciseAccessRule(id string) (object map[string]interface{}, err error) {
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
	request["Domains.1"] = parts[0]
	request["RegionId"] = client.RegionId
	request["Owner"] = "manual"
	action := "DescribeWebPreciseAccessRule"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"DomainNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DomainPreciseAccessRule", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.PreciseAccessConfigList[0].RuleList", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("DomainPreciseAccessRule", id), errs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DomainPreciseAccessRule", id), errs.NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Name"]) != parts[1] {
			continue
		}
		return item, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("DomainPreciseAccessRule", id), errs.NotFoundMsg, response)
}

func (s *Service) DdosCooDomainPreciseAccessRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.DdosCooDomainPreciseAccessRuleStateRefreshFuncWithApi(id, field, failStates, s.DescribeDdosCooDomainPreciseAccessRule)
}

func (s *Service) DdosCooDomainPreciseAccessRuleStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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
		if field == "$.RuleList[*].ConditionList[*]" {
			e := jsonata.MustCompile("$.RuleList.ConditionList")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.RuleList[*].ConditionList[*].Content" {
			e := jsonata.MustCompile("$.RuleList.ConditionList.($exists(Content) ? Content : null)")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.RuleList[*].ConditionList[*].Field" {
			e := jsonata.MustCompile("$.RuleList.ConditionList.($exists(Field) ? Field : null)")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.RuleList[*].ConditionList[*].HeaderName" {
			e := jsonata.MustCompile("$.RuleList.ConditionList.($exists(HeaderName) ? HeaderName : null)")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.RuleList[*].ConditionList[*].MatchMethod" {
			e := jsonata.MustCompile("$.RuleList.ConditionList.($exists(MatchMethod) ? MatchMethod : null)")
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

// DescribeDdosCooDomainPreciseAccessRule >>> Encapsulated.

// DescribeDdosCooWebCcRule <<< Encapsulated get interface for DdosCoo WebCcRule.

func (s *Service) DescribeDdosCooWebCcRule(id string) (object map[string]interface{}, err error) {
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
	request["Domain"] = parts[0]

	request["Owner"] = "manual"
	action := "DescribeWebCCRulesV2"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"DomainNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("WebCcRule", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.WebCCRules[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("WebCcRule", id), errs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("WebCcRule", id), errs.NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Expires"]) != "0" {
			continue
		}
		if fmt.Sprint(item["Name"]) != parts[1] {
			continue
		}
		return item, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("WebCcRule", id), errs.NotFoundMsg, response)
}

func (s *Service) DdosCooWebCcRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.DdosCooWebCcRuleStateRefreshFuncWithApi(id, field, failStates, s.DescribeDdosCooWebCcRule)
}

func (s *Service) DdosCooWebCcRuleStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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
		if field == "$.RuleDetail.Condition[*]" {
			e := jsonata.MustCompile("$.RuleDetail.Condition")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.RuleDetail.Condition[*].Content" {
			e := jsonata.MustCompile("$.RuleDetail.Condition.($exists(Content) ? Content : null)")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.RuleDetail.Condition[*].Field" {
			e := jsonata.MustCompile("$.RuleDetail.Condition.($exists(Field) ? Field : null)")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.RuleDetail.Condition[*].HeaderName" {
			e := jsonata.MustCompile("$.RuleDetail.Condition.($exists(HeaderName) ? HeaderName : null)")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.RuleDetail.Condition[*].MatchMethod" {
			e := jsonata.MustCompile("$.RuleDetail.Condition.($exists(MatchMethod) ? MatchMethod : null)")
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

// DescribeDdosCooWebCcRule >>> Encapsulated.
