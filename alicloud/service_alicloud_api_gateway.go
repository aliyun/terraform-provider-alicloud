package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
)

func (client *AliyunClient) DescribeApiGroup(groupId string) (apiGroup *cloudapi.DescribeApiGroupResponse, err error) {
	req := cloudapi.CreateDescribeApiGroupRequest()
	req.GroupId = groupId

	apiGroup, err = client.cloudapiconn.DescribeApiGroup(req)
	if err != nil {
		if IsExceptedError(err, ApiGroupNotFound) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("ApiGroup", groupId))
		}
		return
	}
	if apiGroup == nil || apiGroup.GroupId == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("ApiGroup", groupId))
	}
	return
}
