package alicloud

import (
	"strings"

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
			err = WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, name, "GetService", AliyunLogGoSdkERROR)
		}
		return
	}
	service, _ = raw.(*fc.GetServiceOutput)
	if service == nil || *service.ServiceName == "" {
		err = WrapErrorf(Error(GetNotFoundMessage("fc_service", name)), NotFoundMsg, AliyunLogGoSdkERROR)
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
			err = WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, name, "GetFunction", AliyunLogGoSdkERROR)
		}
		return
	}
	function, _ = raw.(*fc.GetFunctionOutput)
	if function == nil || *function.FunctionName == "" {
		err = WrapErrorf(Error(GetNotFoundMessage("fc_function", name)), NotFoundMsg, AliyunLogGoSdkERROR)
	}
	return
}

func (s *FcService) DescribeFcTrigger(service, function, name string) (trigger *fc.GetTriggerOutput, err error) {
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.GetTrigger(fc.NewGetTriggerInput(service, function, name))
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound, TriggerNotFound}) {
			err = WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, name, "GetTrigger", AliyunLogGoSdkERROR)
		}
		return
	}
	trigger, _ = raw.(*fc.GetTriggerOutput)
	if trigger == nil || *trigger.TriggerName == "" {
		err = WrapErrorf(Error(GetNotFoundMessage("fc_trigger", name)), NotFoundMsg, AliyunLogGoSdkERROR)
	}
	return
}

func removeSpaceAndEnter(s string) string {
	if Trim(s) == "" {
		return Trim(s)
	}
	return strings.Replace(strings.Replace(strings.Replace(s, " ", "", -1), "\n", "", -1), "\t", "", -1)
}
