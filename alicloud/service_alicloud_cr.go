package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
		} `json:"namespace"`
	} `json:"data"`
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
		if NotFoundError(err) || IsExceptedError(err, ErrorNamespaceNotExist) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, namespaceName, req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return resp, nil
}
