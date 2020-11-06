package alicloud

import (
	"time"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
			if IsExpectedErrors(err, []string{LogClientTimeout, "ProjectForbidden"}) {
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
		if IsExpectedErrors(err, []string{"ProjectNotExist"}) {
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
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
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
		if IsExpectedErrors(err, []string{"ProjectNotExist", "LogStoreNotExist"}) {
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
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
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
		if IsExpectedErrors(err, []string{"ProjectNotExist", "LogStoreNotExist", "IndexConfigNotExist"}) {
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
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
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
		if IsExpectedErrors(err, []string{"ProjectNotExist", "GroupNotExist", "MachineGroupNotExist"}) {
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
			if IsExpectedErrors(err, []string{"InternalServerError"}) {
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
		if IsExpectedErrors(err, []string{"ProjectNotExist", "LogStoreNotExist", "ConfigNotExist"}) {
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
			if IsExpectedErrors(err, []string{"InternalServerError"}) {
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
		if IsExpectedErrors(err, []string{"ProjectNotExist", "ConfigNotExist", "MachineGroupNotExist"}) {
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

func (s *LogService) DescribeLogAlert(id string) (*sls.Alert, error) {
	alert := &sls.Alert{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return alert, WrapError(err)
	}
	projectName, alertName := parts[0], parts[1]
	var requestInfo *sls.Client
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetAlert(projectName, alertName)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetLogstoreAlert", raw, requestInfo, map[string]string{
				"project":    projectName,
				"alert_name": alertName,
			})
		}
		alert, _ = raw.(*sls.Alert)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ProjectNotExist", "JobNotExist"}) {
			return alert, WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		}
		return alert, WrapErrorf(err, DefaultErrorMsg, id, "GetLogstoreAlert", AliyunLogGoSdkERROR)
	}

	if alert == nil || alert.Name == "" {
		return alert, WrapErrorf(Error(GetNotFoundMessage("LogstoreAlert", id)), NotFoundMsg, ProviderERROR)
	}
	return alert, nil
}

func (s *LogService) WaitForLogstoreAlert(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogAlert(id)
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

func (s *LogService) CreateLogDashboard(project, name string) error {
	dashboard := sls.Dashboard{
		DashboardName: name,
		ChartList:     []sls.Chart{},
	}
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.CreateDashboard(project, dashboard)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
				return resource.RetryableError(err)
			}
			if err.(*sls.Error).Message == "specified dashboard already exists" {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("CreateLogDashboard", raw, map[string]string{
				"project":        project,
				"dashboard_name": name,
			})
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateDashboard(project, name string, client *sls.Client) error {
	dashboard := sls.Dashboard{
		DashboardName: name,
		ChartList:     []sls.Chart{},
	}
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		err := client.CreateDashboard(project, dashboard)
		if err != nil {
			if err.(*sls.Error).Message == "specified dashboard already exists" {
				return nil
			}
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
				return resource.RetryableError(err)
			}

		}
		if debugOn() {
			addDebug("CreateLogDashboard", dashboard, map[string]string{
				"project":        project,
				"dashboard_name": name,
			})
		}
		return nil
	})
	return err
}

func (s *LogService) DescribeLogAudit(id string) (*slsPop.DescribeAppResponse, error) {
	request := slsPop.CreateDescribeAppRequest()
	response := &slsPop.DescribeAppResponse{}
	request.AppName = "audit"
	raw, err := s.client.WithLogPopClient(func(client *slsPop.Client) (interface{}, error) {
		return client.DescribeApp(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"AppNotExist"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*slsPop.DescribeAppResponse)
	return response, nil
}

func GetCharTitile(project, dashboard, char string, client *sls.Client) string {
	board, err := client.GetDashboard(project, dashboard)
	// If the query fails to ignore the error, return the original value.
	if err != nil {
		return char
	}
	for _, v := range board.ChartList {
		if v.Display.DisplayName == char {
			return v.Title
		} else {
			return char
		}

	}
	return char
}

func (s *LogService) DescribeLogDashboard(id string) (*sls.Dashboard, error) {
	dashboard := &sls.Dashboard{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return dashboard, WrapError(err)
	}
	projectName, dashboardName := parts[0], parts[1]
	var requestInfo *sls.Client
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetDashboard(projectName, dashboardName)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("GetLogstoreDashboard", raw, requestInfo, map[string]string{
				"project":        projectName,
				"dashboard_name": dashboardName,
			})
		}
		dashboard, _ = raw.(*sls.Dashboard)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ProjectNotExist", "DashboardNotExist"}) {
			return dashboard, WrapErrorf(err, NotFoundMsg, AliyunLogGoSdkERROR)
		}
		return dashboard, WrapErrorf(err, DefaultErrorMsg, id, "GetLogstoreDashboard", AliyunLogGoSdkERROR)
	}

	if dashboard == nil || dashboard.DashboardName == "" {
		return dashboard, WrapErrorf(Error(GetNotFoundMessage("LogstoreDashboard", id)), NotFoundMsg, ProviderERROR)
	}
	return dashboard, nil
}

func (s *LogService) WaitForLogDashboard(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogDashboard(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.DashboardName == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DashboardName, name, ProviderERROR)
		}
	}
}
