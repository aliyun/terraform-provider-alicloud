package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/denverdino/aliyungo/common"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type KvstoreService struct {
	client *connectivity.AliyunClient
}

func (s *KvstoreService) DescribeRKVInstanceById(id string) (instance *r_kvstore.DBInstanceAttribute, err error) {
	request := r_kvstore.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = id
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeInstanceAttribute(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore instance", id))
		}
		return nil, err
	}
	resp, _ := raw.(*r_kvstore.DescribeInstanceAttributeResponse)
	if resp == nil || len(resp.Instances.DBInstanceAttribute) <= 0 {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore instance", id))
	}

	return &resp.Instances.DBInstanceAttribute[0], nil
}

func (s *KvstoreService) DescribeRKVInstancebackupPolicy(id string) (policy *r_kvstore.DescribeBackupPolicyResponse, err error) {
	request := r_kvstore.CreateDescribeBackupPolicyRequest()
	request.InstanceId = id
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeBackupPolicy(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore Instance Policy", id))
		}
		return nil, err
	}
	policy, _ = raw.(*r_kvstore.DescribeBackupPolicyResponse)

	if policy == nil {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("KVStore Instance Policy", id))
	}

	return
}

func (s *KvstoreService) WaitForRKVInstance(instanceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := s.DescribeRKVInstanceById(instanceId)
		if err != nil && !NotFoundError(err) {
			return err
		}

		if instance != nil && instance.InstanceStatus == string(status) {
			break
		}

		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *KvstoreService) WaitForRKVInstanceVpcAuthMode(instanceId string, status string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := s.DescribeRKVInstanceById(instanceId)
		if err != nil && !NotFoundError(err) {
			return err
		}

		if instance != nil && instance.VpcAuthMode == string(status) {
			break
		}

		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *KvstoreService) DescribeParameters(instanceId string) (ds *r_kvstore.DescribeParametersResponse, err error) {
	request := r_kvstore.CreateDescribeParametersRequest()
	request.DBInstanceId = instanceId

	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeParameters(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBInstanceIdNotFound, InvalidDBInstanceNameNotFound}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("DB Instance", instanceId))
		}
		return nil, err
	}
	resp, _ := raw.(*r_kvstore.DescribeParametersResponse)
	if resp == nil {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("KVStore Instance Parameter", instanceId))
	}
	return resp, err
}

func (s *KvstoreService) ModifyInstanceConfig(instanceId string, config string) error {
	request := r_kvstore.CreateModifyInstanceConfigRequest()
	request.InstanceId = instanceId
	request.Config = config

	if err := s.WaitForRKVInstance(instanceId, Normal, DefaultLongTimeout); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}
	_, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.ModifyInstanceConfig(request)
	})
	if err != nil {
		return fmt.Errorf("ModifyInstanceConfig got an error: %#v", err)
	}

	return nil
}
