package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type NasService struct {
	client *connectivity.AliyunClient
}

func (s *NasService) DescribeNasFileSystem(id string) (fs nas.DescribeFileSystemsFileSystem1, err error) {

	request := nas.CreateDescribeFileSystemsRequest()
	request.RegionId = s.client.RegionId
	request.FileSystemId = id
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeFileSystems(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ServiceUnavailable, Throttling}) {
				return resource.RetryableError(err)
			}
			if IsExceptedErrors(err, []string{InvalidFileSystemIDNotFound, ForbiddenNasNotFound}) {
				return resource.NonRetryableError(WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*nas.DescribeFileSystemsResponse)
		if response.TotalCount <= 0 {
			return resource.NonRetryableError(WrapErrorf(Error(GetNotFoundMessage("NasFileSystem", id)), NotFoundMsg, ProviderERROR))
		}
		fs = response.FileSystems.FileSystem[0]
		return nil
	})
	return fs, WrapError(err)
}

func (s *NasService) DescribeNasMountTarget(id string) (fs nas.DescribeMountTargetsMountTarget1, err error) {

	request := nas.CreateDescribeMountTargetsRequest()
	request.RegionId = string(s.client.Region)
	split := strings.Split(id, "-")
	request.FileSystemId = split[0]
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeMountTargets(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ServiceUnavailable, Throttling}) {
				return resource.RetryableError(err)
			}
			if IsExceptedErrors(err, NasNotFound) {
				return resource.NonRetryableError(WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*nas.DescribeMountTargetsResponse)
		for _, mountTarget := range response.MountTargets.MountTarget {
			if id == mountTarget.MountTargetDomain {
				fs = mountTarget
				return nil
			}
		}
		return resource.NonRetryableError(WrapErrorf(Error(GetNotFoundMessage("NasMountTarget", id)), NotFoundMsg, ProviderERROR))
	})
	return fs, WrapError(err)
}

func (s *NasService) DescribeNasAccessGroup(id string) (ag nas.DescribeAccessGroupsAccessGroup1, err error) {

	request := nas.CreateDescribeAccessGroupsRequest()
	request.RegionId = string(s.client.Region)
	request.AccessGroupName = id

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessGroups(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ServiceUnavailable, Throttling}) {
				return resource.RetryableError(err)
			}
			if IsExceptedErrors(err, []string{InvalidAccessGroupNotFound, ForbiddenNasNotFound}) {
				return resource.NonRetryableError(WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*nas.DescribeAccessGroupsResponse)
		if len(response.AccessGroups.AccessGroup) <= 0 {
			return resource.NonRetryableError(WrapErrorf(Error(GetNotFoundMessage("NasAccessGroup", id)), NotFoundMsg, ProviderERROR))
		}
		ag = response.AccessGroups.AccessGroup[0]
		return nil
	})
	return ag, WrapError(err)
}

func (s *NasService) DescribeNasAccessRule(id string) (fs nas.DescribeAccessRulesAccessRule1, err error) {

	request := nas.CreateDescribeAccessRulesRequest()
	request.RegionId = string(s.client.Region)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request.AccessGroupName = parts[0]

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessRules(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ServiceUnavailable, Throttling}) {
				return resource.RetryableError(err)
			}
			if IsExceptedErrors(err, []string{InvalidAccessGroupNotFound, ForbiddenNasNotFound}) {
				return resource.NonRetryableError(WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*nas.DescribeAccessRulesResponse)
		for _, accessRule := range response.AccessRules.AccessRule {
			if parts[1] == accessRule.AccessRuleId {
				fs = accessRule
				return nil
			}
		}
		return resource.NonRetryableError(WrapErrorf(Error(GetNotFoundMessage("NasAccessRule", id)), NotFoundMsg, ProviderERROR))
	})
	return fs, WrapError(err)
}

func (s *NasService) WaitForNasMountTarget(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNasMountTarget(id)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		} else if strings.ToLower(object.Status) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *NasService) WaitForNasFileSystem(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNasFileSystem(id)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strings.ToLower(object.FileSystemId) == strings.ToLower(id) && status != Deleted {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.FileSystemId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *NasService) WaitForNasAccessRule(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNasAccessRule(id)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		} else if strings.ToLower(object.AccessRuleId) == strings.ToLower(id) && status != Deleted {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AccessRuleId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *NasService) WaitForNasAccessGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNasAccessGroup(id)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strings.ToLower(object.AccessGroupName) == strings.ToLower(id) && status != Deleted {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AccessGroupName, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}
