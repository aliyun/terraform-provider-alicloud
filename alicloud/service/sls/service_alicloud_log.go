package sls

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var SlsClientTimeoutCatcher = helper.Catcher{Reason: errs.LogClientTimeout, RetryCount: 15, RetryWaitSeconds: 5}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

type Status string

const (
	Pending     = Status("Pending")
	Creating    = Status("Creating")
	Running     = Status("Running")
	Available   = Status("Available")
	Unavailable = Status("Unavailable")
	Modifying   = Status("Modifying")
	Deleting    = Status("Deleting")
	Starting    = Status("Starting")
	Stopping    = Status("Stopping")
	Stopped     = Status("Stopped")
	Normal      = Status("Normal")
	Changing    = Status("Changing")
	Online      = Status("online")
	Configuring = Status("configuring")

	Associating   = Status("Associating")
	Unassociating = Status("Unassociating")
	InUse         = Status("InUse")
	DiskInUse     = Status("In_use")

	Active   = Status("Active")
	Inactive = Status("Inactive")
	Idle     = Status("Idle")

	SoldOut = Status("SoldOut")

	InService      = Status("InService")
	Removing       = Status("Removing")
	DisabledStatus = Status("Disabled")

	Init            = Status("Init")
	Provisioning    = Status("Provisioning")
	Updating        = Status("Updating")
	FinancialLocked = Status("FinancialLocked")

	PUBLISHED   = Status("Published")
	NOPUBLISHED = Status("NonPublished")

	Deleted = Status("Deleted")
	Null    = Status("Null")

	Enable = Status("Enable")
	BINDED = Status("BINDED")
)

const (
	COLON_SEPARATED = ":"
	DefaultTimeout  = 120
)

type Service struct {
	client *connectivity.AliyunClient
}

func (s *Service) DescribeLogProject(id string) (*sls.LogProject, error) {
	project := &sls.LogProject{}
	var requestInfo *sls.Client
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetProject(id)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout, "ProjectForbidden"}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetProject", raw, requestInfo, map[string]string{"name": id})
		}
		project, _ = raw.(*sls.LogProject)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist"}) {
			return project, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return project, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetProject", errs.AliyunLogGoSdkERROR)
	}
	if project == nil || project.Name == "" {
		return project, errs.WrapErrorf(errs.NotFoundErr("LogProject", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return project, nil
}

func (s *Service) WaitForLogProject(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeLogProject(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Name == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Name, id, errs.ProviderERROR)
		}
	}
}

func (s *Service) DescribeLogStore(id string) (*sls.LogStore, error) {
	store := &sls.LogStore{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return store, errs.WrapError(err)
	}
	projectName, name := parts[0], parts[1]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetLogStore(projectName, name)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetLogStore", raw, requestInfo, map[string]string{
				"project":  projectName,
				"logstore": name,
			})
		}
		store, _ = raw.(*sls.LogStore)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist", "LogStoreNotExist"}) {
			return store, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return store, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetLogStore", errs.AliyunLogGoSdkERROR)
	}
	if store == nil || store.Name == "" {
		return store, errs.WrapErrorf(errs.NotFoundErr("LogStore", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return store, nil
}

func (s *Service) WaitForLogStore(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogStore(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Name == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Name, name, errs.ProviderERROR)
		}
	}
}

func (s *Service) DescribeLogStoreIndex(id string) (*sls.Index, error) {
	index := &sls.Index{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return index, errs.WrapError(err)
	}
	projectName, name := parts[0], parts[1]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetIndex(projectName, name)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetIndex", raw, requestInfo, map[string]string{
				"project":  projectName,
				"logstore": name,
			})
		}
		index, _ = raw.(*sls.Index)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist", "LogStoreNotExist", "IndexConfigNotExist"}) {
			return index, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return index, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetIndex", errs.AliyunLogGoSdkERROR)
	}

	if index == nil || (index.Line == nil && index.Keys == nil) {
		return index, errs.WrapErrorf(errs.NotFoundErr("LogStoreIndex", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return index, nil
}

func (s *Service) DescribeLogMachineGroup(id string) (*sls.MachineGroup, error) {
	group := &sls.MachineGroup{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return group, errs.WrapError(err)
	}
	projectName, groupName := parts[0], parts[1]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetMachineGroup(projectName, groupName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetMachineGroup", raw, requestInfo, map[string]string{
				"project":      projectName,
				"machineGroup": groupName,
			})
		}
		group, _ = raw.(*sls.MachineGroup)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist", "GroupNotExist", "MachineGroupNotExist"}) {
			return group, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return group, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetMachineGroup", errs.AliyunLogGoSdkERROR)
	}

	if group == nil || group.Name == "" {
		return group, errs.WrapErrorf(errs.NotFoundErr("LogMachineGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return group, nil
}

func (s *Service) WaitForLogMachineGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogMachineGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Name == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Name, name, errs.ProviderERROR)
		}
	}
}

func (s *Service) DescribeLogtailConfig(id string) (*sls.LogConfig, error) {
	response := &sls.LogConfig{}
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return response, errs.WrapError(err)
	}
	projectName, configName := parts[0], parts[2]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			slsClient.RetryTimeOut = 30 * time.Second
			return slsClient.GetConfig(projectName, configName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"unconvertable", "Unconvertable"}) {
				return retry.NonRetryableError(err)
			}
			if errs.IsExpectedErrors(err, []string{"InternalServerError"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetConfig", raw, requestInfo, map[string]string{
				"project": projectName,
				"config":  configName,
			})
		}
		response, _ = raw.(*sls.LogConfig)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist", "LogStoreNotExist", "ConfigNotExist"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetConfig", errs.AliyunLogGoSdkERROR)
	}
	if response == nil || response.Name == "" {
		return response, errs.WrapErrorf(errs.NotFoundErr("LogTailConfig", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, nil
}

func (s *Service) WaitForLogtailConfig(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	name := parts[2]
	for {
		object, err := s.DescribeLogtailConfig(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Name == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Name, name, errs.ProviderERROR)
		}
	}
}

func (s *Service) DescribeLogtailAttachment(id string) (groupName string, err error) {
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return groupName, errs.WrapError(err)
	}
	projectName, configName, name := parts[0], parts[1], parts[2]
	var groupNames []string
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {

		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetAppliedMachineGroups(projectName, configName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetAppliedMachineGroups", raw, requestInfo, map[string]string{
				"project":  projectName,
				"confName": configName,
			})
		}
		groupNames, _ = raw.([]string)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist", "ConfigNotExist", "MachineGroupNotExist"}) {
			return groupName, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return groupName, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetAppliedMachineGroups", errs.AliyunLogGoSdkERROR)
	}
	for _, group_name := range groupNames {
		if group_name == name {
			groupName = group_name
		}
	}
	if groupName == "" {
		return groupName, errs.WrapErrorf(errs.NotFoundErr("LogtailAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return groupName, nil
}

func (s *Service) WaitForLogtailAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	name := parts[2]
	for {
		object, err := s.DescribeLogtailAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object, name, errs.ProviderERROR)
		}
	}
}

func (s *Service) DescribeLogAlertResource(id string) (map[string]string, error) {
	var result = map[string]string{}
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return result, errs.WrapError(err)
	}
	resourceType := parts[1]
	if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		_, err := s.client.WithLogPopClient(func(slsPopClient *slsPop.Client) (interface{}, error) {
			switch resourceType {
			case "user":
				_, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
					record, err := slsClient.GetResourceRecord("sls.alert.global_config", "default_config")
					if err != nil {
						return nil, err
					}
					var alertGlobalConfig AlertGlobalConfig
					err = json.Unmarshal([]byte(record.Value), &alertGlobalConfig)
					if err != nil {
						return nil, err
					}
					region := alertGlobalConfig.ConfigDetail.AlertCenterLog.Region
					accountId, err := s.client.AccountId()
					if err != nil {
						return nil, err
					}
					projectName := fmt.Sprintf("sls-alert-%s-%s", accountId, region)
					endpoint := slsClient.Endpoint
					slsClient.Endpoint = strings.Replace(endpoint, s.client.RegionId, region, 1)
					_, err = slsClient.GetProject(projectName)
					if err != nil {
						slsClient.Endpoint = endpoint
						return nil, err
					}
					_, err = slsClient.GetLogStore(projectName, "internal-alert-center-log")
					slsClient.Endpoint = endpoint
					if err != nil {
						return nil, err
					}
					return nil, nil
				})
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"ProjectNotExist"}) || errs.IsExpectedErrors(err, []string{"LogStoreNotExist"}) {
						return result, nil
					}
					return nil, err
				}
				lang := parts[2]
				result["type"] = resourceType
				result["project"] = ""
				result["lang"] = lang
				return result, nil
			case "project":
				project := parts[2]
				_, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
					return slsClient.GetLogStore(project, "internal-alert-history")
				})
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"LogStoreNotExist"}) {
						return nil, nil
					}
					return nil, err
				}
				result["type"] = resourceType
				result["project"] = project
				result["lang"] = ""
				return result, nil
			default:
				return result, errs.WrapErrorf(errors.New("type error"), errs.DefaultErrorMsg, "alicloud_log_alert_resource", "ReadAlertResource", errs.AliyunLogGoSdkERROR)
			}
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return result, errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_alert_resource", "ReadAlertResource", errs.AliyunLogGoSdkERROR)
	}
	return result, nil
}

func (s *Service) DescribeLogAlert(id string) (*sls.Alert, error) {
	alert := &sls.Alert{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return alert, errs.WrapError(err)
	}
	projectName, alertName := parts[0], parts[1]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetAlert(projectName, alertName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetLogstoreAlert", raw, requestInfo, map[string]string{
				"project":    projectName,
				"alert_name": alertName,
			})
		}
		alert, _ = raw.(*sls.Alert)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist", "JobNotExist"}) {
			return alert, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return alert, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetLogstoreAlert", errs.AliyunLogGoSdkERROR)
	}

	if alert == nil || alert.Name == "" {
		return alert, errs.WrapErrorf(errs.NotFoundErr("LogstoreAlert", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return alert, nil
}

func (s *Service) WaitForLogstoreAlert(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	name := parts[1]
	for {
		object, err := s.DescribeLogAlert(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Name == name && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Name, name, errs.ProviderERROR)
		}
	}
}

func (s *Service) CreateLogDashboard(project, name string) error {
	dashboard := sls.Dashboard{
		DashboardName: name,
		ChartList:     []sls.Chart{},
	}
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.CreateDashboard(project, dashboard)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			if err.(*sls.Error).Message == "specified dashboard already exists" {
				return nil
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("CreateLogDashboard", raw, map[string]string{
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
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		err := client.CreateDashboard(project, dashboard)
		if err != nil {
			if err.(*sls.Error).Message == "specified dashboard already exists" {
				return nil
			}
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}

		}
		if logs.DebugOn() {
			logs.AddDebug("CreateLogDashboard", dashboard, map[string]string{
				"project":        project,
				"dashboard_name": name,
			})
		}
		return nil
	})
	return err
}

func (s *Service) DescribeLogAudit(id string) (*slsPop.DescribeAppResponse, error) {
	request := slsPop.CreateDescribeAppRequest()
	response := &slsPop.DescribeAppResponse{}
	request.AppName = "audit"
	raw, err := s.client.WithLogPopClient(func(client *slsPop.Client) (interface{}, error) {
		return client.DescribeApp(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AppNotExist"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
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

func (s *Service) DescribeLogDashboard(id string) (string, error) {
	dashboard := ""
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return dashboard, errs.WrapError(err)
	}
	projectName, dashboardName := parts[0], parts[1]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetDashboardString(projectName, dashboardName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetLogstoreDashboard", raw, requestInfo, map[string]string{
				"project":        projectName,
				"dashboard_name": dashboardName,
			})
		}
		dashboard, _ = raw.(string)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist", "DashboardNotExist"}) {
			return dashboard, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return dashboard, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetLogstoreDashboard", errs.AliyunLogGoSdkERROR)
	}

	if dashboard == "" {
		return dashboard, errs.WrapErrorf(errs.NotFoundErr("LogstoreDashboard", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return dashboard, nil
}

func (s *Service) WaitForLogDashboard(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	name := parts[1]
	for {
		objectString, err := s.DescribeLogDashboard(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if objectString != "" && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, objectString, name, errs.ProviderERROR)
		}
	}
}

func (s *Service) DescribeLogProjectPolicy(id string) (string, error) {
	policy := ""
	projectName := id
	var requestInfo *sls.Client
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetProjectPolicy(projectName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetLogProjectPolicy", raw, requestInfo, map[string]string{
				"project": projectName,
			})
		}
		policy, _ = raw.(string)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist"}) {
			return policy, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return policy, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetProjectPolicy", errs.AliyunLogGoSdkERROR)
	}

	return policy, nil
}

func (s *Service) DescribeLogProjectTags(project_name string) ([]*sls.ResourceTagResponse, error) {
	var requestInfo *sls.Client
	var respTags []*sls.ResourceTagResponse

	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			raw, _, err := slsClient.ListTagResources(project_name, "project", []string{project_name}, []sls.ResourceFilterTag{}, "")
			return raw, err
		})

		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetProjectTags", raw, requestInfo, map[string]string{"project_name": project_name})
		}
		respTags = raw.([]*sls.ResourceTagResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist"}) {
			return respTags, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return respTags, errs.WrapErrorf(err, errs.DefaultErrorMsg, project_name, "GetProejctTags", errs.AliyunLogGoSdkERROR)
	}
	return respTags, nil
}

func (s *Service) DescribeLogEtl(id string) (*sls.ETL, error) {
	etl := &sls.ETL{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return etl, errs.WrapError(err)
	}
	projectName, etlName := parts[0], parts[1]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetETL(projectName, etlName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetLogETL", raw, requestInfo, map[string]string{
				"project":  projectName,
				"etl_name": etlName,
			})
		}
		etl, _ = raw.(*sls.ETL)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist", "JobNotExist"}) {
			return etl, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return etl, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetETL", errs.AliyunLogGoSdkERROR)
	}
	return etl, nil
}

func (s *Service) WaitForLogETL(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeLogEtl(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Name == parts[1] && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Name, id, errs.ProviderERROR)
		}
	}
}

func (s *Service) LogETLStateRefreshFunc(id string, failStates []string, slsClient *sls.Client) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		parts, err := helper.ParseResourceId(id, 2)
		if err != nil {
			return nil, "", errs.WrapError(err)
		}
		object, err := slsClient.GetETL(parts[0], parts[1])
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}

		return object, object.Status, nil
	}
}

func (s *Service) LogOssShipperStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeLogEtl(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}

		return object, object.Status, nil
	}
}

func (s *Service) DescribeLogOssShipper(id string) (*sls.Shipper, error) {
	var shipper *sls.Shipper
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return shipper, errs.WrapError(err)
	}
	projectName, logstoreName, shipperName := parts[0], parts[1], parts[2]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			project, _ := sls.NewLogProject(projectName, slsClient.Endpoint, slsClient.AccessKeyID, slsClient.AccessKeySecret)
			project, _ = project.WithToken(slsClient.SecurityToken)
			logstore, _ := sls.NewLogStore(logstoreName, project)
			return logstore.GetShipper(shipperName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetLogOssShipper", raw, requestInfo, map[string]string{
				"project_name":  projectName,
				"logstore_name": logstoreName,
				"shipper_name":  shipperName,
			})
		}
		shipper, _ = raw.(*sls.Shipper)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist"}) {
			return shipper, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		// SLS server problem, temporarily by returning nil value to solve.
		if d, ok := err.(*sls.Error); ok {
			if d.Message == fmt.Sprintf("shipperName %s does not exist", parts[2]) {
				return shipper, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
			}

		}
		return shipper, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetLogOssShipper", errs.AliyunLogGoSdkERROR)
	}
	return shipper, nil
}

func (s *Service) WaitForLogOssShipper(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeLogOssShipper(id)
		if err != nil {
			if object == nil {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.ShipperName == parts[2] && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.ShipperName, id, errs.ProviderERROR)
		}
	}
}

func (s *Service) LogProjectStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeLogProject(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *Service) DescribeLogIngestion(id string) (*sls.Ingestion, error) {
	var ingestion *sls.Ingestion
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return ingestion, errs.WrapError(err)
	}
	projectName, logstoreName, ingestionName := parts[0], parts[1], parts[2]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetIngestion(projectName, ingestionName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetIngestion", raw, requestInfo, map[string]string{
				"project":        projectName,
				"logstore":       logstoreName,
				"ingestion_name": ingestionName,
			})
		}
		ingestion, _ = raw.(*sls.Ingestion)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist", "JobNotExist"}) {
			return ingestion, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return ingestion, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetIngestion", errs.AliyunLogGoSdkERROR)
	}
	return ingestion, nil
}

func (s *Service) WaitForLogIngestion(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeLogIngestion(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Name == parts[2] && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Name, id, errs.ProviderERROR)
		}
	}
}
func (s *Service) DescribeLogOssExport(id string) (*sls.Export, error) {
	return s.describeLogExport(id)
}

func (s *Service) describeLogExport(id string) (*sls.Export, error) {
	var export *sls.Export
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return export, errs.WrapError(err)
	}
	projectName, logstoreName, exportName := parts[0], parts[1], parts[2]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetExport(projectName, exportName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetExport", raw, requestInfo, map[string]string{
				"project":     projectName,
				"logstore":    logstoreName,
				"export_name": exportName,
			})
		}
		export, _ = raw.(*sls.Export)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ProjectNotExist", "JobNotExist"}) {
			return export, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return export, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetExport", errs.AliyunLogGoSdkERROR)
	}
	return export, nil
}

func (s *Service) WaitForLogOssExport(id string, status Status, timeout int) error {
	return s.waitForLogExport(id, status, timeout)
}

func (s *Service) waitForLogExport(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeLogOssExport(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Name == parts[2] && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Name, id, errs.ProviderERROR)
		}
	}
}

func (s *Service) DescribeLogResource(id string) (*sls.Resource, error) {
	res := &sls.Resource{}
	resourceName := id
	var requestInfo *sls.Client
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetResource(resourceName)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetResource", raw, requestInfo, map[string]string{
				"resource_name": resourceName,
			})
		}
		res, _ = raw.(*sls.Resource)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotExist"}) {
			return res, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return res, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetLogResource", errs.AliyunLogGoSdkERROR)
	}

	if res == nil {
		return res, errs.WrapErrorf(errs.NotFoundErr("LogResource", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return res, nil
}

func (s *Service) WaitForLogResource(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	resourceName := id
	for {
		object, err := s.DescribeLogResource(resourceName)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Name == resourceName && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Name, resourceName, errs.ProviderERROR)
		}
	}
}

func (s *Service) DescribeLogResourceRecord(id string) (*sls.ResourceRecord, error) {
	res := &sls.ResourceRecord{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return res, errs.WrapError(err)
	}
	resourceName, recordId := parts[0], parts[1]
	var requestInfo *sls.Client
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.GetResourceRecord(resourceName, recordId)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("GetResourceRecord", raw, requestInfo, map[string]string{
				"resource_name": resourceName,
				"record_id":     recordId,
			})
		}
		res, _ = raw.(*sls.ResourceRecord)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotExist", "ResourceRecordNotExist"}) {
			return res, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunLogGoSdkERROR)
		}
		return res, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetLogResourceRecord", errs.AliyunLogGoSdkERROR)
	}

	if res == nil {
		return res, errs.WrapErrorf(errs.NotFoundErr("LogResourceRecord", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return res, nil
}

func (s *Service) WaitForLogResourceRecord(id string, status Status, timeout int) error {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	recordId := parts[1]
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeLogResourceRecord(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Id == recordId && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Id, recordId, errs.ProviderERROR)
		}
	}
}
