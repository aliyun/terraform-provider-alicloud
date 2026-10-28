package cr

import (
	"errors"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func (c *CrService) ListCrEEInstances(pageNo int, pageSize int) (*cr_ee.ListInstanceResponse, error) {
	response := &cr_ee.ListInstanceResponse{}
	request := cr_ee.CreateListInstanceRequest()
	request.RegionId = c.client.RegionId
	request.PageNo = requests.NewInteger(pageNo)
	request.PageSize = requests.NewInteger(pageSize)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListInstance(request)
	})
	if err != nil {
		return response, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "ListInstance", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListInstanceResponse)
	if !response.ListInstanceIsSuccess {
		return response, errs.WrapErrorf(errors.New(response.Code), errs.DataDefaultErrorMsg, "ListInstance", action, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DescribeCrEEInstance(instanceId string) (*cr_ee.GetInstanceResponse, error) {
	response := &cr_ee.GetInstanceResponse{}
	request := cr_ee.CreateGetInstanceRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	action := request.GetActionName()

	err := retry.Retry(6*time.Second, func() *retry.RetryError {
		raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.GetInstance(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
				time.Sleep(time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, raw, request.RpcRequest, request)
		response, _ = raw.(*cr_ee.GetInstanceResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, instanceId, action, errs.AlibabaCloudSdkGoERROR)
	}
	if !response.GetInstanceIsSuccess {
		if response.Code == "INSTANCE_NOT_EXIST" {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(fmt.Errorf("%v", response), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) GetCrEEInstanceUsage(instanceId string) (*cr_ee.GetInstanceUsageResponse, error) {
	response := &cr_ee.GetInstanceUsageResponse{}
	request := cr_ee.CreateGetInstanceUsageRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	action := request.GetActionName()

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.GetInstanceUsage(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, raw, request.RpcRequest, request)
		response, _ = raw.(*cr_ee.GetInstanceUsageResponse)
		if !response.GetInstanceUsageIsSuccess && response.Code == "INSTANCE_STATUS_NOT_SUPPORT" {
			wait()
			return retry.RetryableError(errs.WrapErrorf(fmt.Errorf("%v", response), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR))
		}
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, instanceId, action, errs.AlibabaCloudSdkGoERROR)
	}

	if !response.GetInstanceUsageIsSuccess {
		if response.Code == "INSTANCE_NOT_EXIST" {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(fmt.Errorf("%v", response), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) ListCrEEInstanceEndpoint(instanceId string) (*cr_ee.ListInstanceEndpointResponse, error) {
	response := &cr_ee.ListInstanceEndpointResponse{}
	request := cr_ee.CreateListInstanceEndpointRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	action := request.GetActionName()

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.ListInstanceEndpoint(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, raw, request.RpcRequest, request)
		response, _ = raw.(*cr_ee.ListInstanceEndpointResponse)
		if !response.ListInstanceEndpointIsSuccess && response.Code == "INSTANCE_STATUS_NOT_SUPPORT" {
			// INSTANCE_STATUS_NOT_SUPPORT is not a transient condition — the instance does
			// not support managed public/VPC endpoints (e.g. retired instances). Stop
			// retrying immediately so data.alicloud_cr_ee_instances can skip this instance
			// instead of exhausting the 5-minute retry budget and failing the whole read.
			return retry.NonRetryableError(fmt.Errorf("INSTANCE_STATUS_NOT_SUPPORT"))
		}
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, instanceId, action, errs.AlibabaCloudSdkGoERROR)
	}
	if !response.ListInstanceEndpointIsSuccess {
		if response.Code == "INSTANCE_NOT_EXIST" {
			return response, errs.WrapErrorf(fmt.Errorf("%v", response), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(fmt.Errorf("%v", response), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) ListCrEENamespaces(instanceId string, pageNo int, pageSize int) (*cr_ee.ListNamespaceResponse, error) {
	response := &cr_ee.ListNamespaceResponse{}
	request := cr_ee.CreateListNamespaceRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	request.PageNo = requests.NewInteger(pageNo)
	request.PageSize = requests.NewInteger(pageSize)
	action := request.GetActionName()

	var raw interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.ListNamespace(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, raw, request.RpcRequest, request)

	if err != nil {
		return response, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, instanceId, action, errs.AlibabaCloudSdkGoERROR)
	}

	response, _ = raw.(*cr_ee.ListNamespaceResponse)

	if !response.ListNamespaceIsSuccess {
		return response, errs.WrapErrorf(errors.New(response.Code), errs.DataDefaultErrorMsg, instanceId, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (c *CrService) DescribeCrEENamespace(id string) (*cr_ee.GetNamespaceResponse, error) {
	response := &cr_ee.GetNamespaceResponse{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return response, errs.WrapError(err)
	}

	request := cr_ee.CreateGetNamespaceRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = parts[0]
	request.NamespaceName = parts[1]

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.GetNamespace(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.GetNamespaceResponse)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NAMESPACE_NOT_EXIST"}) {
			return response, errs.WrapErrorf(errs.NotFoundErr("CrEE:Namespace", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response.RequestId))
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if !response.GetNamespaceIsSuccess {
		if response.Code == "NAMESPACE_NOT_EXIST" {
			return response, errs.WrapErrorf(errs.NotFoundErr("CrEE:Namespace", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response.RequestId))
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (c *CrService) DeleteCrEENamespace(id string) (*cr_ee.DeleteNamespaceResponse, error) {
	response := &cr_ee.DeleteNamespaceResponse{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return response, errs.WrapError(err)
	}
	request := cr_ee.CreateDeleteNamespaceRequest()
	request.RegionId = c.client.RegionId

	request.InstanceId = parts[0]
	request.NamespaceName = parts[1]
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.DeleteNamespace(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NAMESPACE_NOT_EXIST"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.DeleteNamespaceResponse)
	if !response.DeleteNamespaceIsSuccess {
		if response.Code == "NAMESPACE_NOT_EXIST" {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(fmt.Errorf("%v", response), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) WaitForCrEENamespace(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		resp, err := c.DescribeCrEENamespace(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		if resp.NamespaceName == parts[1] && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, resp.NamespaceName, parts[1], errs.ProviderERROR)
		}
		time.Sleep(3 * time.Second)
	}
}

func (c *CrService) ListCrEERepos(instanceId string, namespace string, pageNo int, pageSize int) (*cr_ee.ListRepositoryResponse, error) {
	response := &cr_ee.ListRepositoryResponse{}
	request := cr_ee.CreateListRepositoryRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	request.RepoNamespaceName = namespace
	request.RepoStatus = "ALL"
	request.PageNo = requests.NewInteger(pageNo)
	request.PageSize = requests.NewInteger(pageSize)
	action := request.GetActionName()

	var raw interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.ListRepository(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, raw, request.RpcRequest, request)

	if err != nil {
		return response, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, fmt.Sprint(instanceId, ":", namespace), action, errs.AlibabaCloudSdkGoERROR)
	}

	response, _ = raw.(*cr_ee.ListRepositoryResponse)

	if !response.ListRepositoryIsSuccess {
		return response, errs.WrapErrorf(errors.New(response.Code), errs.DataDefaultErrorMsg, fmt.Sprint(instanceId, ":", namespace), action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (c *CrService) DescribeCrEERepo(id string) (*cr_ee.GetRepositoryResponse, error) {
	response := &cr_ee.GetRepositoryResponse{}
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return response, errs.WrapError(err)
	}

	request := cr_ee.CreateGetRepositoryRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = parts[0]
	request.RepoNamespaceName = parts[1]
	request.RepoName = parts[2]

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.GetRepository(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.GetRepositoryResponse)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
			return response, errs.WrapErrorf(errs.NotFoundErr("CrEE:Repo", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response.RequestId))
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if !response.GetRepositoryIsSuccess {
		if response.Code == "REPO_NOT_EXIST" {
			return response, errs.WrapErrorf(errs.NotFoundErr("CrEE:Repo", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response.RequestId))
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (c *CrService) DeleteCrEERepo(id, repoId string) (*cr_ee.DeleteRepositoryResponse, error) {
	response := &cr_ee.DeleteRepositoryResponse{}
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return response, errs.WrapError(err)
	}
	request := cr_ee.CreateDeleteRepositoryRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = parts[0]
	request.RepoId = repoId
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.DeleteRepository(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.DeleteRepositoryResponse)
	if !response.DeleteRepositoryIsSuccess {
		if response.Code == "REPO_NOT_EXIST" {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(fmt.Errorf("%v", response), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) WaitForCrEERepo(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		resp, err := c.DescribeCrEERepo(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if resp.RepoName == parts[2] && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, resp.RepoName, parts[2], errs.ProviderERROR)
		}
		time.Sleep(3 * time.Second)
	}
}

func (c *CrService) ListCrEERepoTags(instanceId string, repoId string, pageNo int, pageSize int) (*cr_ee.ListRepoTagResponse, error) {
	response := &cr_ee.ListRepoTagResponse{}
	request := cr_ee.CreateListRepoTagRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	request.RepoId = repoId
	request.PageNo = requests.NewInteger(pageNo)
	request.PageSize = requests.NewInteger(pageSize)
	action := request.GetActionName()

	var raw interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.ListRepoTag(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, raw, request.RpcRequest, request)

	if err != nil {
		return response, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, fmt.Sprint(instanceId, ":", repoId), action, errs.AlibabaCloudSdkGoERROR)
	}

	response, _ = raw.(*cr_ee.ListRepoTagResponse)

	if !response.ListRepoTagIsSuccess {
		return response, errs.WrapErrorf(errors.New(response.Code), errs.DataDefaultErrorMsg, fmt.Sprint(instanceId, ":", repoId), action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (c *CrService) DescribeCrEESyncRule(id string) (*cr_ee.SyncRulesItem, error) {
	response := &cr_ee.ListRepoSyncRuleResponse{}
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := cr_ee.CreateListRepoSyncRuleRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = parts[0]
	request.NamespaceName = parts[1]
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNo = requests.NewInteger(1)

	idExist := false
	for {
		var raw interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			raw, err = c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
				return creeClient.ListRepoSyncRule(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if err != nil {
			return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		response, _ = raw.(*cr_ee.ListRepoSyncRuleResponse)

		if !response.ListRepoSyncRuleIsSuccess {
			return nil, errs.WrapErrorf(fmt.Errorf("%v", response), errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		if len(response.SyncRules) < 1 {
			return nil, errs.WrapErrorf(errs.NotFoundErr("CrEE:SyncRule", id), errs.NotFoundWithResponse, response)
		}

		for _, rule := range response.SyncRules {
			if rule.LocalInstanceId == parts[0] && rule.LocalNamespaceName == parts[1] && rule.SyncRuleId == parts[2] {
				idExist = true
				return &rule, nil
			}
		}

		if len(response.SyncRules) < helper.PageSizeLarge {
			break
		}

		pageNum, err := helper.GetNextpageNumber(request.PageNo)
		if err != nil {
			return nil, errs.WrapError(err)
		}

		request.PageNo = pageNum
	}

	if !idExist {
		return nil, errs.WrapErrorf(errs.NotFoundErr("CrEE:SyncRule", id), errs.NotFoundWithResponse, response)
	}

	return nil, nil
}
