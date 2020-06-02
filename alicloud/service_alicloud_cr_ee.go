package alicloud

import (
	"errors"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
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
		return response, WrapErrorf(err, DataDefaultErrorMsg, "ListInstance", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListInstanceResponse)
	if !response.ListInstanceIsSuccess {
		return response, WrapErrorf(errors.New(response.Code), DataDefaultErrorMsg, "ListInstance", action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DescribeCrEEInstance(instanceId string) (*cr_ee.GetInstanceResponse, error) {
	response := &cr_ee.GetInstanceResponse{}
	request := cr_ee.CreateGetInstanceRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	resource := instanceId
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.GetInstance(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.GetInstanceResponse)
	if !response.GetInstanceIsSuccess {
		return response, c.wrapCrServiceError(resource, action, response.Code)
	}
	return response, nil
}

func (c *CrService) GetCrEEInstanceUsage(instanceId string) (*cr_ee.GetInstanceUsageResponse, error) {
	response := &cr_ee.GetInstanceUsageResponse{}
	request := cr_ee.CreateGetInstanceUsageRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	resource := instanceId
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.GetInstanceUsage(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.GetInstanceUsageResponse)
	if !response.GetInstanceUsageIsSuccess {
		return response, c.wrapCrServiceError(resource, action, response.Code)
	}
	return response, nil
}

func (c *CrService) ListCrEEInstanceEndpoint(instanceId string) (*cr_ee.ListInstanceEndpointResponse, error) {
	response := &cr_ee.ListInstanceEndpointResponse{}
	request := cr_ee.CreateListInstanceEndpointRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	resource := instanceId
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListInstanceEndpoint(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListInstanceEndpointResponse)
	if !response.ListInstanceEndpointIsSuccess {
		return response, c.wrapCrServiceError(resource, action, response.Code)
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
	resource := c.GenResourceId(instanceId)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListNamespace(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DataDefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListNamespaceResponse)
	if !response.ListNamespaceIsSuccess {
		return response, WrapErrorf(errors.New(response.Code), DataDefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DescribeCrEENamespace(id string) (*cr_ee.GetNamespaceResponse, error) {
	strRet := c.ParseResourceId(id)
	instanceId := strRet[0]
	namespaceName := strRet[1]
	response := &cr_ee.GetNamespaceResponse{}
	request := cr_ee.CreateGetNamespaceRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	request.NamespaceName = namespaceName
	resource := c.GenResourceId(instanceId, namespaceName)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.GetNamespace(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NAMESPACE_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.GetNamespaceResponse)
	if !response.GetNamespaceIsSuccess {
		return response, c.wrapCrServiceError(resource, action, response.Code)
	}
	return response, nil
}

func (c *CrService) DeleteCrEENamespace(instanceId string, namespaceName string) (*cr_ee.DeleteNamespaceResponse, error) {
	response := &cr_ee.DeleteNamespaceResponse{}
	request := cr_ee.CreateDeleteNamespaceRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	request.NamespaceName = namespaceName
	resource := c.GenResourceId(instanceId, namespaceName)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.DeleteNamespace(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NAMESPACE_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.DeleteNamespaceResponse)
	if !response.DeleteNamespaceIsSuccess {
		return response, c.wrapCrServiceError(resource, action, response.Code)
	}
	return response, nil
}

func (c *CrService) WaitForCrEENamespace(instanceId string, namespaceName string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	resource := c.GenResourceId(instanceId, namespaceName)

	for {
		resp, err := c.DescribeCrEENamespace(c.GenResourceId(instanceId, namespaceName))
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if resp.NamespaceName == namespaceName && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, resource, GetFunc(1), timeout, resp.NamespaceName, namespaceName, ProviderERROR)
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
	resource := c.GenResourceId(instanceId, namespace)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListRepository(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DataDefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListRepositoryResponse)
	if !response.ListRepositoryIsSuccess {
		return response, WrapErrorf(errors.New(response.Code), DataDefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) DescribeCrEERepo(id string) (*cr_ee.GetRepositoryResponse, error) {
	strRet := c.ParseResourceId(id)
	instanceId := strRet[0]
	namespace := strRet[1]
	repo := strRet[2]
	response := &cr_ee.GetRepositoryResponse{}
	request := cr_ee.CreateGetRepositoryRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	request.RepoNamespaceName = namespace
	request.RepoName = repo
	resource := c.GenResourceId(instanceId, namespace, repo)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.GetRepository(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.GetRepositoryResponse)
	if !response.GetRepositoryIsSuccess {
		return response, c.wrapCrServiceError(resource, action, response.Code)
	}

	return response, nil
}

func (c *CrService) DeleteCrEERepo(instanceId, namespace, repo, repoId string) (*cr_ee.DeleteRepositoryResponse, error) {
	response := &cr_ee.DeleteRepositoryResponse{}
	request := cr_ee.CreateDeleteRepositoryRequest()
	request.RegionId = c.client.RegionId
	request.InstanceId = instanceId
	request.RepoId = repoId
	resource := c.GenResourceId(instanceId, namespace, repo)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.DeleteRepository(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.DeleteRepositoryResponse)
	if !response.DeleteRepositoryIsSuccess {
		return response, c.wrapCrServiceError(resource, action, response.Code)
	}
	return response, nil
}

func (c *CrService) WaitForCrEERepo(instanceId string, namespace string, repo string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	resource := c.GenResourceId(instanceId, namespace, repo)

	for {
		resp, err := c.DescribeCrEERepo(c.GenResourceId(instanceId, namespace, repo))
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if resp.RepoName == repo && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, resource, GetFunc(1), timeout, resp.RepoName, repo, ProviderERROR)
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
	resource := c.GenResourceId(instanceId, repoId)
	action := request.GetActionName()

	raw, err := c.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
		return creeClient.ListRepoTag(request)
	})
	if err != nil {
		return response, WrapErrorf(err, DataDefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, raw, request.RpcRequest, request)

	response, _ = raw.(*cr_ee.ListRepoTagResponse)
	if !response.ListRepoTagIsSuccess {
		return response, WrapErrorf(errors.New(response.Code), DataDefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (c *CrService) wrapCrServiceError(resource string, action string, code string) error {
	switch code {
	case "INSTANCE_NOT_EXIST", "NAMESPACE_NOT_EXIST", "REPO_NOT_EXIST":
		return WrapErrorf(errors.New(code), NotFoundMsg, AlibabaCloudSdkGoERROR)
	default:
		return WrapErrorf(errors.New(code), DefaultErrorMsg, resource, action, AlibabaCloudSdkGoERROR)
	}
}

func (c *CrService) GenResourceId(args ...string) string {
	return strings.Join(args, COLON_SEPARATED)
}

func (c *CrService) ParseResourceId(id string) []string {
	return strings.Split(id, COLON_SEPARATED)
}
