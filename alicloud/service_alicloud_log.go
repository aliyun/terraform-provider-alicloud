package alicloud

import (
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

var SlsClientTimeoutCatcher = Catcher{LogClientTimeout, 15, 5}

type LogService struct {
	client *connectivity.AliyunClient
}

func (s *LogService) DescribeLogProject(id string) (*sls.LogProject, error) {
	project := &sls.LogProject{}
	var requestInfo *sls.Client
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetProject(id)
		})
		if err != nil {
			if IsExceptedError(err, LogClientTimeout) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetProject", raw, requestInfo, map[string]string{"name": id})
		}
		project, _ = raw.(*sls.LogProject)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, ProjectNotExist) {
			return project, WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		}
		return project, WrapErrorf(err, DefaultErrorMsg, id, "GetProject", AliyunLogGoSdkERROR)
	}
	if project == nil || project.Name == "" {
		return project, WrapErrorf(Error(GetNotFoundMessage("LogProject", id)), NotFoundMsg, ProviderERROR)
	}
	return project, nil
}

func (s *LogService) WaitForLogProject(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeLogProject(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Name == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Name, id, ProviderERROR)
		}
	}
}

func (s *LogService) DescribeLogStore(id string) (*sls.LogStore, error) {
	store := &sls.LogStore{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return store, WrapError(err)
	}
	projectName, name := parts[0], parts[1]
	var requestInfo *sls.Client
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetLogStore(projectName, name)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InternalServerError, LogClientTimeout}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetLogStore", raw, requestInfo, map[string]string{
				"project":  projectName,
				"logstore": name,
			})
		}
		store, _ = raw.(*sls.LogStore)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist}) {
			return store, WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		}
		return store, WrapErrorf(err, DefaultErrorMsg, id, "GetLogStore", AliyunLogGoSdkERROR)
	}
	if store == nil || store.Name == "" {
		return store, WrapErrorf(Error(GetNotFoundMessage("LogStore", id)), NotFoundMsg, ProviderERROR)
	}
	return store, nil
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

func (s *LogService) DescribeLogStoreIndex(id string) (*sls.Index, error) {
	index := &sls.Index{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return index, WrapError(err)
	}
	projectName, name := parts[0], parts[1]
	var requestInfo *sls.Client
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetIndex(projectName, name)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InternalServerError, LogClientTimeout}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetIndex", raw, requestInfo, map[string]string{
				"project":  projectName,
				"logstore": name,
			})
		}
		index, _ = raw.(*sls.Index)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist, IndexConfigNotExist}) {
			return index, WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		}
		return index, WrapErrorf(err, DefaultErrorMsg, id, "GetIndex", AliyunLogGoSdkERROR)
	}

	if index == nil || (index.Line == nil && index.Keys == nil) {
		return index, WrapErrorf(Error(GetNotFoundMessage("LogStoreIndex", id)), NotFoundMsg, ProviderERROR)
	}
	return index, nil
}

func (s *LogService) DescribeLogMachineGroup(id string) (*sls.MachineGroup, error) {
	group := &sls.MachineGroup{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return group, WrapError(err)
	}
	projectName, groupName := parts[0], parts[1]
	var requestInfo *sls.Client
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetMachineGroup(projectName, groupName)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InternalServerError, LogClientTimeout}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetMachineGroup", raw, requestInfo, map[string]string{
				"project":      projectName,
				"machineGroup": groupName,
			})
		}
		group, _ = raw.(*sls.MachineGroup)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, []string{ProjectNotExist, GroupNotExist, MachineGroupNotExist}) {
			return group, WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		}
		return group, WrapErrorf(err, DefaultErrorMsg, id, "GetMachineGroup", AliyunLogGoSdkERROR)
	}

	if group == nil || group.Name == "" {
		return group, WrapErrorf(Error(GetNotFoundMessage("LogMachineGroup", id)), NotFoundMsg, ProviderERROR)
	}
	return group, nil
}

func (s *LogService) WaitForLogMachineGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogMachineGroup(id)
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

func (s *LogService) DescribeLogtailConfig(id string) (*sls.LogConfig, error) {
	response := &sls.LogConfig{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return response, WrapError(err)
	}
	projectName, configName := parts[0], parts[2]
	var requestInfo *sls.Client
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetConfig(projectName, configName)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InternalServerError}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetConfig", raw, requestInfo, map[string]string{
				"project": projectName,
				"config":  configName,
			})
		}
		response, _ = raw.(*sls.LogConfig)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist, LogConfigNotExist}) {
			return response, WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, "GetConfig", AliyunLogGoSdkERROR)
	}
	if response == nil || response.Name == "" {
		return response, WrapErrorf(Error(GetNotFoundMessage("LogTailConfig", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func (s *LogService) WaitForLogtailConfig(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	name := parts[2]
	for {
		object, err := s.DescribeLogtailConfig(id)
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

func (s *LogService) DescribeLogtailAttachment(id string) (groupName string, err error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return groupName, WrapError(err)
	}
	projectName, configName, name := parts[0], parts[1], parts[2]
	var groupNames []string
	var requestInfo *sls.Client
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetAppliedMachineGroups(projectName, configName)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InternalServerError}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetAppliedMachineGroups", raw, requestInfo, map[string]string{
				"project":  projectName,
				"confName": configName,
			})
		}
		groupNames, _ = raw.([]string)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ProjectNotExist, LogConfigNotExist, MachineGroupNotExist}) {
			return groupName, WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		}
		return groupName, WrapErrorf(err, DefaultErrorMsg, id, "GetAppliedMachineGroups", AliyunLogGoSdkERROR)
	}
	for _, group_name := range groupNames {
		if group_name == name {
			groupName = group_name
		}
	}
	if groupName == "" {
		return groupName, WrapErrorf(Error(GetNotFoundMessage("LogtailAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	return groupName, nil
}

func (s *LogService) WaitForLogtailAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	name := parts[2]
	for {
		object, err := s.DescribeLogtailAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object, name, ProviderERROR)
		}
	}
}
