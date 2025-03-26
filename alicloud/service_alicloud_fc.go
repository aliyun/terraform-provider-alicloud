package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type FcService struct {
	client *connectivity.AliyunClient
}

func (s *FcService) DescribeFcService(id string) (*fc.GetServiceOutput, error) {
	response := &fc.GetServiceOutput{}
	request := &fc.GetServiceInput{ServiceName: &id}
	var requestInfo *fc.Client
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.GetService(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ServiceNotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, FcGoSdk)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "GetService", FcGoSdk)
		}
		return response, err
	}
	addDebug("GetService", raw, requestInfo, request)
	response, _ = raw.(*fc.GetServiceOutput)
	if *response.ServiceName != id {
		err = WrapErrorf(NotFoundErr("FcService", id), NotFoundMsg, FcGoSdk)
	}
	return response, err
}

func (s *FcService) WaitForFcService(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeFcService(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if *object.ServiceName == id && status != Deleted {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, *object.ServiceName, id, ProviderERROR)
		}
	}
}

func (s *FcService) DescribeFcFunction(id string) (*fc.GetFunctionOutput, error) {
	response := &fc.GetFunctionOutput{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return response, WrapError(err)
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
		if IsExpectedErrors(err, []string{"ServiceNotFound", "FunctionNotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, FcGoSdk)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "GetFunction", FcGoSdk)
		}
		return response, err
	}
	addDebug("GetFunction", raw, requestInfo, request)
	response, _ = raw.(*fc.GetFunctionOutput)
	if *response.FunctionName == "" {
		err = WrapErrorf(NotFoundErr("FcFunction", id), NotFoundMsg, FcGoSdk)
	}
	return response, err
}

func (s *FcService) WaitForFcFunction(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeFcFunction(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if *object.FunctionName == parts[1] && status != Deleted {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, *object.FunctionName, parts[1], ProviderERROR)
		}
	}
	return nil
}

func (s *FcService) DescribeFcTrigger(id string) (*fc.GetTriggerOutput, error) {
	response := &fc.GetTriggerOutput{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return response, WrapError(err)
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
		if IsExpectedErrors(err, []string{"ServiceNotFound", "FunctionNotFound", "TriggerNotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, FcGoSdk)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "FcTrigger", FcGoSdk)
		}
		return response, err
	}
	addDebug("GetTrigger", raw, requestInfo, request)
	response, _ = raw.(*fc.GetTriggerOutput)
	if *response.TriggerName != name {
		err = WrapErrorf(NotFoundErr("FcTrigger", name), NotFoundMsg, ProviderERROR)
	}
	return response, err
}

func (s *FcService) DescribeFcAlias(id string) (*fc.GetAliasOutput, error) {
	response := &fc.GetAliasOutput{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
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
		if IsExpectedErrors(err, []string{"ServiceNotFound", "AliasNotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, FcGoSdk)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "GetAlias", FcGoSdk)
		}
		return response, err
	}
	addDebug("GetAlias", raw, requestInfo, request)
	response, _ = raw.(*fc.GetAliasOutput)
	if *response.AliasName == "" {
		err = WrapErrorf(NotFoundErr("FcAlias", id), NotFoundMsg, FcGoSdk)
	}
	return response, err
}

func removeSpaceAndEnter(s string) string {
	if Trim(s) == "" {
		return Trim(s)
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
	_, errors = validation.ValidateJsonString(v, k)
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

func (s *FcService) WaitForFcTrigger(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeFcTrigger(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if *object.TriggerName == parts[2] {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, *object.TriggerName, parts[2], ProviderERROR)
		}
	}
	return nil
}

func (s *FcService) DescribeFcCustomDomain(id string) (*fc.GetCustomDomainOutput, error) {
	request := &fc.GetCustomDomainInput{DomainName: &id}
	response := &fc.GetCustomDomainOutput{}

	var requestInfo *fc.Client
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.GetCustomDomain(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainNameNotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, FcGoSdk)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "FcCustomDomain", FcGoSdk)
		}
		return response, err
	}
	addDebug("GetCustomDomain", raw, requestInfo, request)
	response, _ = raw.(*fc.GetCustomDomainOutput)
	if *response.DomainName != id {
		err = WrapErrorf(NotFoundErr("FcCustomDomain", id), NotFoundMsg, ProviderERROR)
	}
	return response, err
}

func (s *FcService) WaitForFcCustomDomain(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeFcCustomDomain(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if *object.DomainName == id && status != Deleted {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, *object.DomainName, id, ProviderERROR)
		}
	}
}

func (s *FcService) DescribeFcFunctionAsyncInvokeConfig(id string) (*fc.GetFunctionAsyncInvokeConfigOutput, error) {
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
		if IsExpectedErrors(err, []string{"ServiceNotFound", "FunctionNotFound", "AsyncConfigNotExists"}) {
			err = WrapErrorf(err, NotFoundMsg, FcGoSdk)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "FcFunctionAsyncInvokeConfig", FcGoSdk)
		}
		return response, err
	}
	addDebug("GetFunctionAsyncInvokeConfig", raw, requestInfo, request)
	response, _ = raw.(*fc.GetFunctionAsyncInvokeConfigOutput)
	return response, err
}

func (s *FcService) SetResourceTags(d *schema.ResourceData, resourceArn *string) error {
	if d.HasChange("tags") {
		added, removed := parsingTags(d)

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "FcUnTagResource", FcGoSdk)
			}
			addDebug("FcUnTagResource", raw, requestInfo, request)
		}
		if len(addedTags) > 0 {
			request := &fc.TagResourceInput{Tags: addedTags, ResourceArn: resourceArn}
			var requestInfo *fc.Client
			raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
				requestInfo = fcClient
				return fcClient.TagResource(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "FcTagResource", FcGoSdk)
			}
			addDebug("FcTagResource", raw, requestInfo, request)
		}
	}
	return nil
}

func (s *FcService) WaitForFcFunctionAsyncInvokeConfig(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		_, err := s.DescribeFcFunctionAsyncInvokeConfig(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, id, ProviderERROR)
		}
	}
}
