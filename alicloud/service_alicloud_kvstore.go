package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/denverdino/aliyungo/common"
)

func (client *AliyunClient) DescribeRKVInstanceById(id string) (instance *r_kvstore.DBInstanceAttribute, err error) {
	request := r_kvstore.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = id
	resp, err := client.rkvconn.DescribeInstanceAttribute(request)
	if err != nil {
		if IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore instance", id))
		}
		return nil, err
	}

	if resp == nil || len(resp.Instances.DBInstanceAttribute) <= 0 {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore instance", id))
	}

	return &resp.Instances.DBInstanceAttribute[0], nil
}

func (client *AliyunClient) DescribeRKVInstancebackupPolicy(id string) (policy *r_kvstore.DescribeBackupPolicyResponse, err error) {
	request := r_kvstore.CreateDescribeBackupPolicyRequest()
	request.InstanceId = id
	policy, err = client.rkvconn.DescribeBackupPolicy(request)
	if err != nil {
		if IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore Instance Policy", id))
		}
		return nil, err
	}

	if policy == nil {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("KVStore Instance Policy", id))
	}

	return
}

func (client *AliyunClient) WaitForRKVInstance(instanceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := client.DescribeRKVInstanceById(instanceId)
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
