package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
)

type CrService struct {
	client *connectivity.AliyunClient
}

type crDefaultResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
	} `json:"data"`
}

type crCreateNamespaceRequestPayload struct {
	Namespace struct {
		Namespace string `json:"Namespace"`
	} `json:"Namespace"`
}

type crUpdateNamespaceRequestPayload struct {
	Namespace struct {
		AutoCreate        bool   `json:"AutoCreate"`
		DefaultVisibility string `json:"DefaultVisibility"`
	} `json:"Namespace"`
}

type crDescribeNamespaceResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Namespace struct {
			Namespace         string `json:"namespace"`
			AuthorizeType     string `json:"authorizeType"`
			DefaultVisibility string `json:"defaultVisibility"`
			AutoCreate        bool   `json:"autoCreate"`
			NamespaceStatus   string `json:"namespaceStatus"`
		} `json:"namespace"`
	} `json:"data"`
}

type crDescribeNamespaceListResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Namespace []struct {
			Namespace       string `json:"namespace"`
			AuthorizeType   string `json:"authorizeType"`
			NamespaceStatus string `json:"namespaceStatus"`
		} `json:"namespaces"`
	} `json:"data"`
}

const (
	RepoTypePublic  = "PUBLIC"
	RepoTypePrivate = "PRIVATE"
)

type crCreateRepoRequestPayload struct {
	Repo struct {
		RepoNamespace string `json:"RepoNamespace"`
		RepoName      string `json:"RepoName"`
		Summary       string `json:"Summary"`
		Detail        string `json:"Detail"`
		RepoType      string `json:"RepoType"`
	} `json:"Repo"`
}

type crUpdateRepoRequestPayload struct {
	Repo struct {
		Summary  string `json:"Summary"`
		Detail   string `json:"Detail"`
		RepoType string `json:"RepoType"`
	} `json:"Repo"`
}

type crDescribeRepoResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Repo struct {
			Summary        string `json:"summary"`
			Detail         string `json:"detail"`
			RepoNamespace  string `json:"repoNamespace"`
			RepoName       string `json:"repoName"`
			RepoType       string `json:"repoType"`
			RepoDomainList struct {
				Public   string `json:"public"`
				Internal string `json:"internal"`
				Vpc      string `json:"vpc"`
			}
		} `json:"repo"`
	} `json:"data"`
}

type crDescribeReposResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Repos    []crRepo `json:"repos"`
		Total    int      `json:"total"`
		PageSize int      `json:"pageSize"`
		Page     int      `json:"page"`
	} `json:"data"`
}

type crRepo struct {
	Summary        string `json:"summary"`
	RepoNamespace  string `json:"repoNamespace"`
	RepoName       string `json:"repoName"`
	RepoType       string `json:"repoType"`
	RegionId       string `json:"regionId"`
	RepoDomainList struct {
		Public   string `json:"public"`
		Internal string `json:"internal"`
		Vpc      string `json:"vpc"`
	} `json:"repoDomainList"`
}

type crDescribeRepoTagsResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		Tags     []crTag `json:"tags"`
		Total    int     `json:"total"`
		PageSize int     `json:"pageSize"`
		Page     int     `json:"page"`
	} `json:"data"`
}

type crTag struct {
	ImageId     string `json:"imageId"`
	Digest      string `json:"digest"`
	Tag         string `json:"tag"`
	Status      string `json:"status"`
	ImageUpdate int    `json:"imageUpdate"`
	ImageCreate int    `json:"imageCreate"`
	ImageSize   int    `json:"imageSize"`
}

func (c *CrService) DescribeNamespace(namespaceName string) (*cr.GetNamespaceResponse, error) {
	invoker := NewInvoker()

	req := cr.CreateGetNamespaceRequest()
	req.Namespace = namespaceName

	var resp *cr.GetNamespaceResponse

	if err := invoker.Run(func() error {
		var err error
		raw, err := c.client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.GetNamespace(req)
		})
		resp, _ = raw.(*cr.GetNamespaceResponse)
		return err
	}); err != nil {
		if IsExceptedError(err, ErrorNamespaceNotExist) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, namespaceName, req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return resp, nil
}

func (c *CrService) DescribeRepo(repoPath string) (*cr.GetRepoResponse, error) {
	invoker := NewInvoker()

	sli := strings.Split(repoPath, SLASH_SEPARATED)
	repoNamespace := sli[0]
	repoName := sli[1]

	req := cr.CreateGetRepoRequest()
	req.RepoNamespace = repoNamespace
	req.RepoName = repoName

	var resp *cr.GetRepoResponse

	if err := invoker.Run(func() error {
		var err error
		raw, err := c.client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.GetRepo(req)
		})
		resp, _ = raw.(*cr.GetRepoResponse)
		return err
	}); err != nil {
		if IsExceptedError(err, ErrorRepoNotExist) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, fmt.Sprintf("%s%s%s", repoNamespace, COLON_SEPARATED, repoName), req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return resp, nil
}
