package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type FcService struct {
	client *connectivity.AliyunClient
}

func (s *FcService) DescribeFcService(id string) (response *fc.GetServiceOutput, err error) {
	request := &fc.GetServiceInput{ServiceName: &id}
	var requestInfo *fc.Client
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.GetService(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound}) {
			err = WrapErrorf(err, NotFoundMsg, FcGoSdk)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "GetService", FcGoSdk)
		}
		return
	}
	addDebug("GetService", raw, requestInfo, request)
	response, _ = raw.(*fc.GetServiceOutput)
	if *response.ServiceName != id {
		err = WrapErrorf(Error(GetNotFoundMessage("FcService", id)), NotFoundMsg, FcGoSdk)
	}
	return
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

func (s *FcService) DescribeFcFunction(id string) (response *fc.GetFunctionOutput, err error) {

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
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
		if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound}) {
			err = WrapErrorf(err, NotFoundMsg, FcGoSdk)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "GetFunction", FcGoSdk)
		}
		return
	}
	addDebug("GetFunction", raw, requestInfo, request)
	response, _ = raw.(*fc.GetFunctionOutput)
	if *response.FunctionName == "" {
		err = WrapErrorf(Error(GetNotFoundMessage("FcFunction", id)), NotFoundMsg, FcGoSdk)
	}
	return
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

func (s *FcService) DescribeFcTrigger(id string) (response *fc.GetTriggerOutput, err error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	service, function, name := parts[0], parts[1], parts[2]
	request := fc.NewGetTriggerInput(service, function, name)
	var requestInfo *fc.Client
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.GetTrigger(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound, TriggerNotFound}) {
			err = WrapErrorf(err, NotFoundMsg, FcGoSdk)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "FcTrigger", FcGoSdk)
		}
		return
	}
	addDebug("GetTrigger", raw, requestInfo, request)
	response, _ = raw.(*fc.GetTriggerOutput)
	if *response.TriggerName != name {
		err = WrapErrorf(Error(GetNotFoundMessage("FcTrigger", name)), NotFoundMsg, ProviderERROR)
	}
	return
}

func removeSpaceAndEnter(s string) string {
	if Trim(s) == "" {
		return Trim(s)
	}
	return strings.Replace(strings.Replace(strings.Replace(s, " ", "", -1), "\n", "", -1), "\t", "", -1)
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
