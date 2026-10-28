package fc

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// Status mirrors the root package status type used by the FC service's wait helpers.
type Status string

const (
	Deleted = Status("Deleted")
)

// HeaderEnableEBTrigger mirrors the root package header constant used by FC trigger calls.
const HeaderEnableEBTrigger = "x-fc-enable-eventbridge-trigger"

func (s *Service) DescribeFcService(id string) (*fc.GetServiceOutput, error) {
	response := &fc.GetServiceOutput{}
	request := &fc.GetServiceInput{ServiceName: &id}
	var requestInfo *fc.Client
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.GetService(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ServiceNotFound"}) {
			err = errs.WrapErrorf(err, errs.NotFoundMsg, errs.FcGoSdk)
		} else {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetService", errs.FcGoSdk)
		}
		return response, err
	}
	logs.AddDebug("GetService", raw, requestInfo, request)
	response, _ = raw.(*fc.GetServiceOutput)
	if *response.ServiceName != id {
		err = errs.WrapErrorf(errs.NotFoundErr("Service", id), errs.NotFoundMsg, errs.FcGoSdk)
	}
	return response, err
}

func (s *Service) WaitForFcService(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeFcService(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if *object.ServiceName == id && status != Deleted {
			return nil
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, *object.ServiceName, id, errs.ProviderERROR)
		}
	}
}

func (s *Service) DescribeFcFunction(id string) (*fc.GetFunctionOutput, error) {
	response := &fc.GetFunctionOutput{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return response, errs.WrapError(err)
	}
	service, name := parts[0], parts[1]
	request := &fc.GetFunctionInput{
		ServiceName:  &service,
		FunctionName: &name,
	}
	var requestInfo *fc.Client
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.GetFunction(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ServiceNotFound", "FunctionNotFound"}) {
			err = errs.WrapErrorf(err, errs.NotFoundMsg, errs.FcGoSdk)
		} else {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetFunction", errs.FcGoSdk)
		}
		return response, err
	}
	logs.AddDebug("GetFunction", raw, requestInfo, request)
	response, _ = raw.(*fc.GetFunctionOutput)
	if *response.FunctionName == "" {
		err = errs.WrapErrorf(errs.NotFoundErr("FcFunction", id), errs.NotFoundMsg, errs.FcGoSdk)
	}
	return response, err
}

func (s *Service) WaitForFcFunction(id string, status Status, timeout int) error {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeFcFunction(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if *object.FunctionName == parts[1] && status != Deleted {
			break
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, *object.FunctionName, parts[1], errs.ProviderERROR)
		}
	}
	return nil
}

func (s *Service) DescribeFcTrigger(id string) (*fc.GetTriggerOutput, error) {
	response := &fc.GetTriggerOutput{}
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return response, errs.WrapError(err)
	}
	service, function, name := parts[0], parts[1], parts[2]
	request := fc.NewGetTriggerInput(service, function, name)
	request.WithHeader(HeaderEnableEBTrigger, "enable")
	var requestInfo *fc.Client
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.GetTrigger(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ServiceNotFound", "FunctionNotFound", "TriggerNotFound"}) {
			err = errs.WrapErrorf(err, errs.NotFoundMsg, errs.FcGoSdk)
		} else {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "FcTrigger", errs.FcGoSdk)
		}
		return response, err
	}
	logs.AddDebug("GetTrigger", raw, requestInfo, request)
	response, _ = raw.(*fc.GetTriggerOutput)
	if *response.TriggerName != name {
		err = errs.WrapErrorf(errs.NotFoundErr("FcTrigger", name), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, err
}

func (s *Service) DescribeFcAlias(id string) (*fc.GetAliasOutput, error) {
	response := &fc.GetAliasOutput{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	service, name := parts[0], parts[1]
	request := &fc.GetAliasInput{
		ServiceName: &service,
		AliasName:   &name,
	}
	var requestInfo *fc.Client
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.GetAlias(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ServiceNotFound", "AliasNotFound"}) {
			err = errs.WrapErrorf(err, errs.NotFoundMsg, errs.FcGoSdk)
		} else {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetAlias", errs.FcGoSdk)
		}
		return response, err
	}
	logs.AddDebug("GetAlias", raw, requestInfo, request)
	response, _ = raw.(*fc.GetAliasOutput)
	if *response.AliasName == "" {
		err = errs.WrapErrorf(errs.NotFoundErr("FcAlias", id), errs.NotFoundMsg, errs.FcGoSdk)
	}
	return response, err
}

func removeSpaceAndEnter(s string) string {
	if helper.Trim(s) == "" {
		return helper.Trim(s)
	}
	return strings.Replace(strings.Replace(strings.Replace(s, " ", "", -1), "\n", "", -1), "\t", "", -1)
}

func delEmptyPayloadIfExist(v string, k string) (string, error) {
	if v == "" {
		return v, nil
	}
	in := []byte(v)
	var raw map[string]interface{}
	if err := json.Unmarshal(in, &raw); err != nil {
		log.Printf("[ERROR] %q contains an invalid JSON: %s", k, err)
		return v, err
	}

	if v, ok := raw["payload"]; ok {
		if vStr, ok := v.(string); ok && vStr == "" {
			delete(raw, "payload")
		}
	}

	out, err := json.Marshal(raw)
	if err != nil {
		log.Printf("[ERROR] %q contains an invalid JSON: %s", k, err)
	}
	return string(out), err
}

func ValidateFcTriggerConfig(v interface{}, k string) (ws []string, errors []error) {
	if v == nil {
		return
	}
	_, errors = validation.StringIsJSON(v, k)
	if errors != nil && len(errors) > 0 {
		return
	}
	in := []byte(removeSpaceAndEnter(fmt.Sprint(v)))
	var raw map[string]interface{}
	if err := json.Unmarshal(in, &raw); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
		return
	}

	if _, err := json.Marshal(raw); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}

	return
}

func resolveFcTriggerConfig(v string, k string) (string, error) {
	if v == "" {
		return v, nil
	}
	in := []byte(v)
	var raw map[string]interface{}
	if err := json.Unmarshal(in, &raw); err != nil {
		log.Printf("[ERROR] %q contains an invalid JSON: %s", k, err)
		return v, err
	}

	out, err := json.Marshal(raw)
	if err != nil {
		log.Printf("[ERROR] %q contains an invalid JSON: %s", k, err)
	}
	return string(out), err
}

func delNilEventSourceParams(v string, k string) (string, error) {
	if v == "" {
		return v, nil
	}
	in := []byte(v)
	var raw map[string]interface{}
	if err := json.Unmarshal(in, &raw); err != nil {
		log.Printf("[ERROR] %q contains an invalid JSON: %s", k, err)
		return v, err
	}
	if v, ok := raw["eventSourceConfig"]; ok {
		if eventSourceConfig, ok := v.(map[string]interface{}); ok {
			if v1, ok := eventSourceConfig["eventSourceParameters"]; ok {
				if eventSourceParams, ok := v1.(map[string]interface{}); ok {
					if vMNS, ok := eventSourceParams["sourceMNSParameters"]; ok {
						if _, ok := vMNS.(map[string]interface{}); ok {

						} else if vMNS == nil {
							// sourceMNSParameters is nil
							delete(eventSourceParams, "sourceMNSParameters")
						}

					}
					if vRocketMQ, ok := eventSourceParams["sourceRocketMQParameters"]; ok {
						if _, ok := vRocketMQ.(map[string]interface{}); ok {

						} else if vRocketMQ == nil {
							// sourceRocketMQParameters is nil
							delete(eventSourceParams, "sourceRocketMQParameters")
						}
					}
					if vRabbitMQ, ok := eventSourceParams["sourceRabbitMQParameters"]; ok {
						if _, ok := vRabbitMQ.(map[string]interface{}); ok {

						} else if vRabbitMQ == nil {
							// sourceRabbitMQParameters is nil
							delete(eventSourceParams, "sourceRabbitMQParameters")
						}
					}
				} else if v1 == nil {
					// eventSourceParams is nil
					delete(eventSourceConfig, "eventSourceParameters")
				}
			}
		}
	}
	out, err := json.Marshal(raw)
	if err != nil {
		log.Printf("[ERROR] %q contains an invalid JSON: %s", k, err)
	}
	return string(out), err
}

func (s *Service) WaitForFcTrigger(id string, status Status, timeout int) error {
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeFcTrigger(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if *object.TriggerName == parts[2] {
			break
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, *object.TriggerName, parts[2], errs.ProviderERROR)
		}
	}
	return nil
}

func (s *Service) DescribeFcCustomDomain(id string) (*fc.GetCustomDomainOutput, error) {
	request := &fc.GetCustomDomainInput{DomainName: &id}
	response := &fc.GetCustomDomainOutput{}

	var requestInfo *fc.Client
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.GetCustomDomain(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"DomainNameNotFound"}) {
			err = errs.WrapErrorf(err, errs.NotFoundMsg, errs.FcGoSdk)
		} else {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "FcCustomDomain", errs.FcGoSdk)
		}
		return response, err
	}
	logs.AddDebug("GetCustomDomain", raw, requestInfo, request)
	response, _ = raw.(*fc.GetCustomDomainOutput)
	if *response.DomainName != id {
		err = errs.WrapErrorf(errs.NotFoundErr("FcCustomDomain", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, err
}

func (s *Service) WaitForFcCustomDomain(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeFcCustomDomain(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if *object.DomainName == id && status != Deleted {
			return nil
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, *object.DomainName, id, errs.ProviderERROR)
		}
	}
}

func (s *Service) DescribeFcFunctionAsyncInvokeConfig(id string) (*fc.GetFunctionAsyncInvokeConfigOutput, error) {
	serviceName, functionName, qualifier, err := parseFCDestinationConfigId(id)
	if err != nil {
		return nil, err
	}
	request := &fc.GetFunctionAsyncInvokeConfigInput{
		ServiceName:  &serviceName,
		FunctionName: &functionName,
	}
	if qualifier != "" {
		request.Qualifier = &qualifier
	}
	response := &fc.GetFunctionAsyncInvokeConfigOutput{}

	var requestInfo *fc.Client
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.GetFunctionAsyncInvokeConfig(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ServiceNotFound", "FunctionNotFound", "AsyncConfigNotExists"}) {
			err = errs.WrapErrorf(err, errs.NotFoundMsg, errs.FcGoSdk)
		} else {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "FcFunctionAsyncInvokeConfig", errs.FcGoSdk)
		}
		return response, err
	}
	logs.AddDebug("GetFunctionAsyncInvokeConfig", raw, requestInfo, request)
	response, _ = raw.(*fc.GetFunctionAsyncInvokeConfigOutput)
	return response, err
}

func (s *Service) SetResourceTags(d *schema.ResourceData, resourceArn *string) error {
	if d.HasChange("tags") {
		added, removed := tags.parsingTags(d)

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}

		addedTags := make(map[string]string)
		for k, v := range added {
			addedTags[k] = v.(string)
		}

		if len(removedTagKeys) > 0 {
			request := &fc.UnTagResourceInput{TagKeys: removedTagKeys, ResourceArn: resourceArn}
			var requestInfo *fc.Client
			raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
				requestInfo = fcClient
				return fcClient.UnTagResource(request)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "FcUnTagResource", errs.FcGoSdk)
			}
			logs.AddDebug("FcUnTagResource", raw, requestInfo, request)
		}
		if len(addedTags) > 0 {
			request := &fc.TagResourceInput{Tags: addedTags, ResourceArn: resourceArn}
			var requestInfo *fc.Client
			raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
				requestInfo = fcClient
				return fcClient.TagResource(request)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "FcTagResource", errs.FcGoSdk)
			}
			logs.AddDebug("FcTagResource", raw, requestInfo, request)
		}
	}
	return nil
}

func (s *Service) WaitForFcFunctionAsyncInvokeConfig(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		_, err := s.DescribeFcFunctionAsyncInvokeConfig(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, id, errs.ProviderERROR)
		}
	}
}
