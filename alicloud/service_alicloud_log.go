package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
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

func (s *LogService) DescribeLogStore(projectName, name string) (store *sls.LogStore, err error) {
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetLogStore(projectName, name)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist}) {
				return resource.NonRetryableError(GetNotFoundErrorFromString(GetNotFoundMessage("Log Store", name)))
			}
			if IsExceptedErrors(err, []string{InternalServerError, LogClientTimeout}) {
				return resource.RetryableError(fmt.Errorf("GetLogStore %s got an error: %#v.", name, err))
			}
			return resource.NonRetryableError(fmt.Errorf("GetLogStore %s got an error: %s.", name, err))
		}
		store, _ = raw.(*sls.LogStore)
		return nil
	})

	if err != nil {
		return
	}

	if store == nil || store.Name == "" {
		return store, GetNotFoundErrorFromString(GetNotFoundMessage("Log Store", name))
	}
	return
}

func (s *LogService) DescribeLogStoreIndex(projectName, name string) (index *sls.Index, err error) {
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetIndex(projectName, name)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist, IndexConfigNotExist}) {
				return resource.NonRetryableError(GetNotFoundErrorFromString(GetNotFoundMessage("Log Store", name)))
			}
			if IsExceptedErrors(err, []string{InternalServerError, LogClientTimeout}) {
				return resource.RetryableError(fmt.Errorf("GetLogStore %s got an error: %#v.", name, err))
			}
			return resource.NonRetryableError(fmt.Errorf("GetLogStore %s got an error: %#v.", name, err))
		}
		index, _ = raw.(*sls.Index)
		return nil
	})

	if err != nil {
		return
	}

	if index == nil || (index.Line == nil && index.Keys == nil) {
		return index, GetNotFoundErrorFromString(GetNotFoundMessage("Log Store Index", name))
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

func (s *LogService) DescribeLogtailToMachineGroup(projectName, configName string) (groupNames []string, err error) {
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		group_names, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetAppliedMachineGroups(projectName, configName)
		})
		if len(group_names.([]string)) == 0 {
			return resource.NonRetryableError(WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR))
		}
		if err != nil {
			if IsExceptedErrors(err, []string{ProjectNotExist, LogConfigNotExist, MachineGroupNotExist}) {
				return resource.NonRetryableError(WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR))
			}
			if IsExceptedErrors(err, []string{InternalServerError}) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, configName, "GetAppliedMachineGroups", AliyunLogGoSdkERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, configName, "GetAppliedMachineGroups", AliyunLogGoSdkERROR))
		}
		groupNames, _ = group_names.([]string)
		return nil
	})
	if err != nil {
		return
	}
	if groupNames == nil {
		return nil, WrapError(fmt.Errorf("Configuration not found for application on machine set"))
	}
	return groupNames, nil
}
