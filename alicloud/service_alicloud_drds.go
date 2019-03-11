package alicloud

import (
	"fmt"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
)

type DrdsService struct {
	client *connectivity.AliyunClient
}

// crate Drdsinstance
func (s *DrdsService) CreateDrdsInstance(req *drds.CreateDrdsInstanceRequest) (response *drds.CreateDrdsInstanceResponse, err error) {

	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.CreateDrdsInstance(req)
	})
	resp, _ := raw.(*drds.CreateDrdsInstanceResponse)

	if err != nil {
		return resp, fmt.Errorf("createDrdsInstance got an error: %#v", err)
	}

	return resp, nil
}

func (s *DrdsService) DescribeDrdsInstance(drdsInstanceId string) (response *drds.DescribeDrdsInstanceResponse, err error) {
	req := drds.CreateDescribeDrdsInstanceRequest()
	req.DrdsInstanceId = drdsInstanceId
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstance(req)
	})
	if err != nil {
		return nil, fmt.Errorf("describe drds instance error: %#v", err)
	}
	resp, _ := raw.(*drds.DescribeDrdsInstanceResponse)

	if resp == nil {
		return resp, GetNotFoundErrorFromString(GetNotFoundMessage("Instance", drdsInstanceId))

	}
	return resp, nil
}

func (s *DrdsService) DescribeDrdsInstances(regionId string) (response *drds.DescribeDrdsInstancesResponse, err error) {
	req := drds.CreateDescribeDrdsInstancesRequest()
	req.Type = string(Private)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstances(req)
	})
	resp, _ := raw.(*drds.DescribeDrdsInstancesResponse)

	return resp, err

}

func (s *DrdsService) ModifyDrdsInstanceDescription(request *drds.ModifyDrdsInstanceDescriptionRequest) (response *drds.ModifyDrdsInstanceDescriptionResponse, err error) {
	req := drds.CreateModifyDrdsInstanceDescriptionRequest()
	req.DrdsInstanceId = request.DrdsInstanceId
	req.Description = request.Description
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.ModifyDrdsInstanceDescription(req)
	})
	resp, _ := raw.(*drds.ModifyDrdsInstanceDescriptionResponse)

	return resp, err

}

func (s *DrdsService) RemoveDrdsInstance(drdsInstanceId string) (response *drds.RemoveDrdsInstanceResponse, err error) {
	req := drds.CreateRemoveDrdsInstanceRequest()
	req.DrdsInstanceId = drdsInstanceId
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.RemoveDrdsInstance(req)
	})
	resp, _ := raw.(*drds.RemoveDrdsInstanceResponse)

	return resp, err
}

func (s *DrdsService) WaitForDrdsInstance(instanceId string, status string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := s.DescribeDrdsInstance(instanceId)

		if err != nil && strings.Contains(err.Error(), InvalidDRDSInstanceIdNotFound) {
			time.Sleep(DefaultIntervalMedium * time.Second)
		} else {
			return err
		}

		if instance != nil && instance.Data.Status == status {
			break
		}

		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("DRDS Instance", instanceId))
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}
