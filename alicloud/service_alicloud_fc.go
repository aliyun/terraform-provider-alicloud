package alicloud

import (
	"fmt"

	"github.com/aliyun/fc-go-sdk"
)

func (client *AliyunClient) DescribeFcService(name string) (service *fc.GetServiceOutput, err error) {
	service, err = client.fcconn.GetService(&fc.GetServiceInput{
		ServiceName: &name,
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Service", name))
		} else {
			err = fmt.Errorf("GetService %s got an error: %#v.", name, err)
		}
		return
	}
	if service == nil || *service.ServiceName == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Service", name))
	}
	return
}

func (client *AliyunClient) DescribeFcFunction(service, name string) (function *fc.GetFunctionOutput, err error) {
	function, err = client.fcconn.GetFunction(&fc.GetFunctionInput{
		ServiceName:  &service,
		FunctionName: &name,
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Function", name))
		} else {
			err = fmt.Errorf("GetFunction %s got an error: %#v.", name, err)
		}
		return
	}
	if function == nil || *function.FunctionName == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Function", name))
	}
	return
}

func (client *AliyunClient) DescribeFcTrigger(service, function, name string) (trigger *fc.GetTriggerOutput, err error) {
	trigger, err = client.fcconn.GetTrigger(&fc.GetTriggerInput{
		ServiceName:  &service,
		FunctionName: &function,
		TriggerName:  &name,
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound, TriggerNotFound}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Trigger", name))
		} else {
			err = fmt.Errorf("GetTrigger %s got an error: %#v.", name, err)
		}
		return
	}
	if trigger == nil || *trigger.TriggerName == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("FC Trigger", name))
	}
	return
}
