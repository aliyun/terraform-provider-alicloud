package alicloud

import (
	"fmt"

	"github.com/aliyun/fc-go-sdk"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type FcService struct {
	client *connectivity.AliyunClient
}

func (s *FcService) DescribeFcService(name string) (service *fc.GetServiceOutput, err error) {
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.GetService(&fc.GetServiceInput{ServiceName: &name})
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Service", name))
		} else {
			err = fmt.Errorf("GetService %s got an error: %#v.", name, err)
		}
		return
	}
	service, _ = raw.(*fc.GetServiceOutput)
	if service == nil || *service.ServiceName == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Service", name))
	}
	return
}

func (s *FcService) DescribeFcFunction(service, name string) (function *fc.GetFunctionOutput, err error) {
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.GetFunction(&fc.GetFunctionInput{
			ServiceName:  &service,
			FunctionName: &name,
		})
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Function", name))
		} else {
			err = fmt.Errorf("GetFunction %s got an error: %#v.", name, err)
		}
		return
	}
	function, _ = raw.(*fc.GetFunctionOutput)
	if function == nil || *function.FunctionName == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Function", name))
	}
	return
}

func (s *FcService) DescribeFcTrigger(service, function, name string) (trigger *fc.GetTriggerOutput, err error) {
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.GetTrigger(&fc.GetTriggerInput{
			ServiceName:  &service,
			FunctionName: &function,
			TriggerName:  &name,
		})
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound, TriggerNotFound}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Trigger", name))
		} else {
			err = fmt.Errorf("GetTrigger %s got an error: %#v.", name, err)
		}
		return
	}
	trigger, _ = raw.(*fc.GetTriggerOutput)
	if trigger == nil || *trigger.TriggerName == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Trigger", name))
	}
	return
}
