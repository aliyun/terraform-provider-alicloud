package alicloud

import (
	"fmt"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

var SlsClientTimeoutCatcher = Catcher{LogClientTimeout, 15, 5}

type LogService struct {
	client *connectivity.AliyunClient
}

func (s *LogService) DescribeLogProject(name string) (project *sls.LogProject, err error) {
	invoker := NewInvoker()
	invoker.AddCatcher(SlsClientTimeoutCatcher)
	err = invoker.Run(func() error {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetProject(name)
		})
		if err != nil {
			return err
		}
		project, _ = raw.(*sls.LogProject)
		if project == nil || project.Name == "" {
			return GetNotFoundErrorFromString(GetNotFoundMessage("Log Project", name))
		}
		return nil
	})
	return
}

func (s *LogService) DescribeLogStore(id string) (store *sls.LogStore, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	projectName, name := parts[0], parts[1]
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetLogStore(projectName, name)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InternalServerError, LogClientTimeout}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("GetLogStore", raw)
		store, _ = raw.(*sls.LogStore)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist}) {
			return store, WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, "GetLogStore", AliyunLogGoSdkERROR)
	}
	if store == nil || store.Name == "" {
		return store, WrapErrorf(Error(GetNotFoundMessage("LogStore", id)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *LogService) WaitForLogStore(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogStore(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Name == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Name, name, ProviderERROR)
		}
	}
}

func (s *LogService) DescribeLogStoreIndex(id string) (index *sls.Index, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	projectName, name := parts[0], parts[1]
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetIndex(projectName, name)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InternalServerError, LogClientTimeout}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("GetIndex", raw)
		index, _ = raw.(*sls.Index)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist, IndexConfigNotExist}) {
			return index, WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, "GetIndex", AliyunLogGoSdkERROR)
	}

	if index == nil || (index.Line == nil && index.Keys == nil) {
		return index, WrapErrorf(Error(GetNotFoundMessage("LogStoreIndex", id)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *LogService) DescribeLogMachineGroup(projectName, groupName string) (group *sls.MachineGroup, err error) {

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetMachineGroup(projectName, groupName)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ProjectNotExist, GroupNotExist, MachineGroupNotExist}) {
				return resource.NonRetryableError(GetNotFoundErrorFromString(GetNotFoundMessage("Log Machine Group", groupName)))
			}
			if IsExceptedErrors(err, []string{InternalServerError, LogClientTimeout}) {
				return resource.RetryableError(fmt.Errorf("GetLogMachineGroup %s got an error: %#v.", groupName, err))
			}
			return resource.NonRetryableError(fmt.Errorf("GetLogMachineGroup %s got an error: %#v.", groupName, err))
		}
		group, _ = raw.(*sls.MachineGroup)
		return nil
	})

	if err != nil {
		return
	}

	if group == nil || group.Name == "" {
		return group, GetNotFoundErrorFromString(GetNotFoundMessage("Log Machine Group", groupName))
	}
	return
}

func (s *LogService) DescribeLogLogtailConfig(projectName, configName string) (logconfig *sls.LogConfig, err error) {
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetConfig(projectName, configName)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist, LogConfigNotExist}) {
				return resource.NonRetryableError(WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR))
			}
			if IsExceptedErrors(err, []string{InternalServerError}) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, configName, "GetConfig", AliyunLogGoSdkERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, configName, "GetConfig", AliyunLogGoSdkERROR))
		}
		logconfig, _ = raw.(*sls.LogConfig)
		return nil
	})
	if err != nil {
		return
	}
	if logconfig == nil || logconfig.Name == "" {
		return logconfig, WrapErrorf(Error(GetNotFoundMessage("Log LogTail Config", configName)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *LogService) DescribeLogtailAttachment(id string) (groupNames []string, err error) {
	split := strings.Split(id, COLON_SEPARATED)
	projectName, configName := split[0], split[1]
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		group_names, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetAppliedMachineGroups(projectName, configName)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ProjectNotExist, LogConfigNotExist, MachineGroupNotExist}) {
				return resource.NonRetryableError(WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR))
			}
			if IsExceptedErrors(err, []string{InternalServerError}) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, id, "GetAppliedMachineGroups", AliyunLogGoSdkERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, id, "GetAppliedMachineGroups", AliyunLogGoSdkERROR))
		}

		groupNames, _ = group_names.([]string)
		return nil
	})
	if err != nil {
		return
	}
	if len(groupNames) == 0 {
		return groupNames, WrapErrorf(Error(GetNotFoundMessage("Logtail Attachment", id)), NotFoundMsg, ProviderERROR)
	}
	return groupNames, nil
}
