package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
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
		// convert error code
		if IsExceptedErrors(err, []string{InvalidGpdbInstanceIdNotFound, InvalidGpdbNameNotFound}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
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

	return nil
}

func (s *GpdbService) DescribeGpdbConnection(id string) (*gpdb.DBInstanceNetInfo, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	// Describe DB Instance Net Info
	request := gpdb.CreateDescribeDBInstanceNetInfoRequest()
	request.DBInstanceId = parts[0]
	raw, err := s.client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
		return gpdbClient.DescribeDBInstanceNetInfo(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidGpdbInstanceIdNotFound, InvalidGpdbNameNotFound}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*gpdb.DescribeDBInstanceNetInfoResponse)
	if response.DBInstanceNetInfos.DBInstanceNetInfo != nil {
		for _, o := range response.DBInstanceNetInfos.DBInstanceNetInfo {
			if strings.HasPrefix(o.ConnectionString, parts[1]) {
				return &o, nil
			}
		}
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("GpdbConnection", id)), NotFoundMsg, ProviderERROR)
}

func (s *GpdbService) GpdbInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGpdbInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBInstanceStatus == failState {
				return object, object.DBInstanceStatus, WrapError(Error(FailedToReachTargetStatus, object.DBInstanceStatus))
			}
		}
		return object, object.DBInstanceStatus, nil
	}
}

func (s *GpdbService) WaitForGpdbConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeGpdbConnection(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ConnectionString != "" && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ConnectionString, id, ProviderERROR)
		}
	}
}
