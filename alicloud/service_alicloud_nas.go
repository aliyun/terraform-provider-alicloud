package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type NasService struct {
	client *connectivity.AliyunClient
}

func (s *NasService) DescribeNasFileSystem(id string) (fs nas.FileSystem, err error) {

	request := nas.CreateDescribeFileSystemsRequest()
	request.RegionId = s.client.RegionId
	request.FileSystemId = id
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeFileSystems(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ServiceUnavailable, Throttling}) {
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"InvalidFileSystem.NotFound", "Forbidden.NasNotFound"}) {
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

func (s *NasService) DescribeNasMountTarget(id string) (object nas.MountTarget, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := nas.CreateDescribeMountTargetsRequest()
	request.RegionId = s.client.RegionId
	request.FileSystemId = parts[0]
	request.MountTargetDomain = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeMountTargets(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable", "Throttling"}) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"Forbidden.NasNotFound", "InvalidFileSystem.NotFound", "InvalidLBid.NotFound", "InvalidMountTarget.NotFound", "VolumeUnavailable"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("NasMountTarget", id)), NotFoundMsg, ProviderERROR)
				return resource.NonRetryableError(err)
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*nas.DescribeMountTargetsResponse)

		if len(response.MountTargets.MountTarget) < 1 {
			err = WrapErrorf(Error(GetNotFoundMessage("NasMountTarget", id)), NotFoundMsg, ProviderERROR, response.RequestId)
			return resource.NonRetryableError(err)
		}
		object = response.MountTargets.MountTarget[0]
		return nil
	})
	return object, WrapError(err)
}

func (s *NasService) DescribeNasAccessGroup(id string) (object nas.AccessGroup, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := nas.CreateDescribeAccessGroupsRequest()
	request.RegionId = s.client.RegionId
	request.AccessGroupName = parts[0]
	request.FileSystemType = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessGroups(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable", "Throttling"}) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"Forbidden.NasNotFound", "InvalidAccessGroup.NotFound", "Resource.NotFound"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("NasAccessGroup", id)), NotFoundMsg, ProviderERROR)
				return resource.NonRetryableError(err)
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*nas.DescribeAccessGroupsResponse)

		if len(response.AccessGroups.AccessGroup) < 1 {
			err = WrapErrorf(Error(GetNotFoundMessage("NasAccessGroup", id)), NotFoundMsg, ProviderERROR)
			return resource.NonRetryableError(err)
		}
		object = response.AccessGroups.AccessGroup[0]
		return nil
	})
	return object, WrapError(err)
}

func (s *NasService) DescribeNasAccessRule(id string) (fs nas.AccessRule, err error) {

	request := nas.CreateDescribeAccessRulesRequest()
	request.RegionId = string(s.client.Region)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request.AccessGroupName = parts[0]
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeXLarge)

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessRules(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ServiceUnavailable, Throttling}) {
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"InvalidAccessGroup.NotFound", "Forbidden.NasNotFound"}) {
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
		err = WrapErrorf(Error(GetNotFoundMessage("NasAccessRule", id)), NotFoundMsg, ProviderERROR)
		if len(response.AccessRules.AccessRule) < PageSizeXLarge {
			return resource.NonRetryableError(err)
		}

		page, e := getNextpageNumber(request.PageNumber)
		if e != nil {
			return resource.NonRetryableError(WrapError(e))
		}
		request.PageNumber = page
		return resource.RetryableError(err)
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

func (s *NasService) NasMountTargetStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasMountTarget(id)
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
