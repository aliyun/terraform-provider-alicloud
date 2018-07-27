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
