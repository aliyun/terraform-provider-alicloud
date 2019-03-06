package alicloud

import (
	"encoding/json"
	"fmt"
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

func (c *CrService) DescribeNamespace(namespaceName string) (*crDescribeNamespaceResponse, error) {
	invoker := NewInvoker()

	req := cr.CreateGetNamespaceRequest()
	// FIXME
	// Temporary hack, see https://github.com/aliyun/alibaba-cloud-sdk-go/issues/208
	req.SetDomain(fmt.Sprintf("cr.%s.aliyuncs.com", c.client.RegionId))
	req.Namespace = namespaceName

	var resp crDescribeNamespaceResponse

	if err := invoker.Run(func() error {
		raw, err := c.client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.GetNamespace(req)
		})
		if err != nil {
			return err
		}
		err = json.Unmarshal(raw.(*cr.GetNamespaceResponse).GetHttpContentBytes(), &resp)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, namespaceName, req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return &resp, nil
}
