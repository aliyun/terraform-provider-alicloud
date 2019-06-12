package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type GpdbService struct {
	client *connectivity.AliyunClient
}

func (s *GpdbService) DescribeGpdbInstance(id string) (instanceAttribute gpdb.DBInstanceAttribute, err error) {
	request := gpdb.CreateDescribeDBInstanceAttributeRequest()
	request.DBInstanceId = id
	raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
		return client.DescribeDBInstanceAttribute(request)
	})

	response, _ := raw.(*gpdb.DescribeDBInstanceAttributeResponse)
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidGpdbInstanceIdNotFound, InvalidGpdbNameNotFound}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, "GetFunction", AlibabaCloudSdkGoERROR)
		}
		return
	}

	addDebug(request.GetActionName(), response)
	if len(response.Items.DBInstanceAttribute) == 0 {
		return instanceAttribute, WrapErrorf(Error(GetNotFoundMessage("Gpdb Instance", id)), NotFoundMsg, ProviderERROR)
	}

	return response.Items.DBInstanceAttribute[0], nil
}

func (s *GpdbService) GetSecurityIps(id string) ([]string, error) {
	arr, err := s.DescribeGpdbSecurityIps(id)
	if err != nil {
		return nil, WrapError(err)
	}

	var ips, separator string
	ipsMap := make(map[string]string)
	for _, ip := range arr {
		ips += separator + ip.SecurityIPList
		separator = COMMA_SEPARATED
	}
	for _, ip := range strings.Split(ips, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}
	return finalIps, nil
}

func (s *GpdbService) DescribeGpdbSecurityIps(id string) (ips []gpdb.DBInstanceIPArray, err error) {
	request := gpdb.CreateDescribeDBInstanceIPArrayListRequest()
	request.DBInstanceId = id

	raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
		return client.DescribeDBInstanceIPArrayList(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidGpdbInstanceIdNotFound, InvalidGpdbNameNotFound}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		return
	}
	response, _ := raw.(*gpdb.DescribeDBInstanceIPArrayListResponse)
	addDebug(request.GetActionName(), response)
	return response.Items.DBInstanceIPArray, nil
}

func (s *GpdbService) ModifyGpdbSecurityIps(id, ips string) error {
	request := gpdb.CreateModifySecurityIpsRequest()
	request.DBInstanceId = id
	request.SecurityIPList = ips
	raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
		return client.ModifySecurityIps(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response := raw.(*gpdb.ModifySecurityIpsResponse)
	addDebug(request.GetActionName(), response)
	if err := s.WaitForGpdbInstance(id, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *GpdbService) NotFoundGpdbInstance(err error) bool {
	if NotFoundError(err) || IsExceptedErrors(err, []string{InvalidGpdbInstanceIdNotFound, InvalidGpdbNameNotFound}) {
		return true
	}
	return false
}

// Wait for instance to given status
func (s *GpdbService) WaitForGpdbInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		instance, err := s.DescribeGpdbInstance(id)
		if err != nil {
			if s.NotFoundGpdbInstance(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		// if status equal to given status, then success
		if instance.DBInstanceStatus == string(status) {
			return nil
		}

		// timeout
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, instance.DBInstanceStatus, string(status), ProviderERROR)
		}

		time.Sleep(DefaultIntervalShort * time.Second)
	}
}
