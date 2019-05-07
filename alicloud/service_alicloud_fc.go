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

func (s *FcService) DescribeFcTrigger(id string) (response *fc.GetTriggerOutput, err error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	service, function, name := parts[0], parts[1], parts[2]
	raw, err := s.client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.GetTrigger(fc.NewGetTriggerInput(service, function, name))
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound, TriggerNotFound}) {
			err = WrapErrorf(err, NotFoundMsg, FcGoSdk)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "FcTrigger", FcGoSdk)
		}
		return
	}
	response, _ = raw.(*fc.GetTriggerOutput)
	if *response.TriggerName == "" {
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
			}
			return WrapError(err)
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
