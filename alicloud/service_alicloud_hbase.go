package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

const (
	Hb_LAUNCHING                  = "LAUNCHING"
	Hb_CREATING                   = "CREATING"
	Hb_ACTIVATION                 = "ACTIVATION"
	Hb_DELETING                   = "DELETING"
	Hb_RESTARTING                 = "RESTARTING"
	Hb_CLASS_CHANGING             = "CLASS_CHANGING"
	Hb_TRANSING                   = "TRANSING"
	Hb_VERSION_TRANSING           = "VERSION_TRANSING"
	Hb_MINOR_VERSION_TRANSING     = "MINOR_VERSION_TRANSING"
	Hb_DATABASE_TRANSING          = "DATABASE_TRANSING"
	Hb_GUARD_CREATING             = "GUARD_CREATING"
	Hb_BACKUP_OPENNING            = "BACKUP_OPENNING"
	Hb_BACKUP_RECOVERING          = "BACKUP_RECOVERING"
	Hb_SUBINS_CREATING            = "SUBINS_CREATING"
	Hb_DATA_IMPORTING             = "DATA_IMPORTING"
	Hb_DATABASE_IMPORTING         = "DATABASE_IMPORTING"
	Hb_NODE_RESIZING              = "NODE_RESIZING"
	Hb_DISK_RESIZING              = "DISK_RESIZING"
	Hb_STATE_HBASE_EXPANDING      = "HBASE_EXPANDING"
	Hb_STATE_HBASE_SCALE_OUT      = "HBASE_SCALE_OUT"
	Hb_STATE_MULTMOD_SCALE_OUT    = "HBASE_SCALE_OUT"
	Hb_STATE_HBASE_COLD_EXPANDING = "HBASE_COLD_EXPANDING"
	Hb_NET_SWITCHING              = "NET_SWITCHING"
	Hb_NET_CREATING               = "NET_CREATING"
	Hb_NET_DELETING               = "NET_DELETING"
	Hb_ADD_COMPONENT              = "COMP_ADDING"
	Hb_COMP_REMOVING              = "COMP_REMOVING"
	Hb_STATE_NET_MODIFYING        = "NET_MODIFYING"
	Hb_INSTANCE_LEVEL_MODIFY      = "INSTANCE_LEVEL_MODIFY"
	Hb_GUARD_SWITCHING            = "GUARD_SWITCHING"
	Hb_LINK_SWITCHING             = "LINK_SWITCHING"
	Hb_CREATE_FAILED              = "CREATE_FAILED"
	Hb_DELETED                    = "DELETED"
	Hb_LOCKED                     = "LOCKED"
	Hb_EXCEPTION                  = "EXCEPTION"
)

type HBaseService struct {
	client *connectivity.AliyunClient
}

func (s *HBaseService) NotFoundHBaseInstance(err error) bool {
	if NotFoundError(err) || IsExceptedErrors(err, []string{InvalidHBaseInstanceIdNotFound, InvalidHBaseNameNotFound}) {
		return true
	}
	return false
}

func (s *HBaseService) DescribeHBaseInstance(id string) (instance hbase.DescribeInstanceResponse, err error) {
	request := hbase.CreateDescribeInstanceRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeInstance(request)
	})
	response, _ := raw.(*hbase.DescribeInstanceResponse)
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidHBaseInstanceIdNotFound}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response == nil || response.InstanceId == "" {
		return instance, WrapErrorf(Error(GetNotFoundMessage("HBase Instance", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	// hbase 响应值参数填充
	instance = *response
	return instance, nil
}

/**
pop has limit, support next.
*/
func (s *HBaseService) DescribeIpWhitelist(id string) (instance hbase.DescribeIpWhitelistResponse, err error) {
	request := hbase.CreateDescribeIpWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeIpWhitelist(request)
	})
	response, _ := raw.(*hbase.DescribeIpWhitelistResponse)
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidHBaseInstanceIdNotFound}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *response, nil
}

func (s *HBaseService) HBaseClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHBaseInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}
