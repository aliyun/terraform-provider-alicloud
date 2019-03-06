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

type crCreateNamespaceRequestPayload struct {
	Namespace struct {
		Namespace string `json:"Namespace"`
	} `json:"Namespace"`
}

type crCreateNamespaceResponse struct {
	RequestId string `json:"requestId"`
	Data      struct {
		NamespaceId int64 `json:"namespaceId"`
	} `json:"data"`
}

func (c *CrService) CreateNamespace(namespaceName string) (*crCreateNamespaceResponse, error) {
	invoker := NewInvoker()

	payload := &crCreateNamespaceRequestPayload{}
	payload.Namespace.Namespace = namespaceName
	serialized, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req := cr.CreateCreateNamespaceRequest()
	req.SetContent(serialized)
	// FIXME
	// Temporary hack, see https://github.com/aliyun/alibaba-cloud-sdk-go/issues/208
	req.SetDomain(fmt.Sprintf("cr.%s.aliyuncs.com", c.client.RegionId))

	var resp crCreateNamespaceResponse

	if err := invoker.Run(func() error {
		raw, err := c.client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.CreateNamespace(req)
		})
		if err != nil {
			return err
		}
		err = json.Unmarshal(raw.(*cr.CreateNamespaceResponse).GetHttpContentBytes(), &resp)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}

type crUpdateNamespaceRequestPayload struct {
	Namespace struct {
		AutoCreate        bool   `json:"AutoCreate"`
		DefaultVisibility string `json:"DefaultVisibility"`
	} `json:"Namespace"`
}

func (c *CrService) UpdateNamespace(namespaceName string, autoCreate bool, defaultVisibility string) (*crDefaultResponse, error) {
	invoker := NewInvoker()

	payload := &crUpdateNamespaceRequestPayload{}
	payload.Namespace.DefaultVisibility = defaultVisibility
	payload.Namespace.AutoCreate = autoCreate
	serialized, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req := cr.CreateUpdateNamespaceRequest()
	// FIXME
	// Temporary hack, see https://github.com/aliyun/alibaba-cloud-sdk-go/issues/208
	req.SetDomain(fmt.Sprintf("cr.%s.aliyuncs.com", c.client.RegionId))
	req.SetContent(serialized)
	req.Namespace = namespaceName

	var resp crDefaultResponse

	if err := invoker.Run(func() error {
		raw, err := c.client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.UpdateNamespace(req)
		})
		if err != nil {
			return err
		}
		err = json.Unmarshal(raw.(*cr.UpdateNamespaceResponse).GetHttpContentBytes(), &resp)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}

type crGetNamespaceResponse struct {
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

func (c *CrService) GetNamespace(namespaceName string) (*crGetNamespaceResponse, error) {
	invoker := NewInvoker()

	req := cr.CreateGetNamespaceRequest()
	// FIXME
	// Temporary hack, see https://github.com/aliyun/alibaba-cloud-sdk-go/issues/208
	req.SetDomain(fmt.Sprintf("cr.%s.aliyuncs.com", c.client.RegionId))
	req.Namespace = namespaceName

	var resp crGetNamespaceResponse

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
		return nil, err
	}
	return &resp, nil
}

func (c *CrService) DeleteNamespace(namespaceName string) (*crDefaultResponse, error) {
	invoker := NewInvoker()

	req := cr.CreateDeleteNamespaceRequest()
	// FIXME
	// Temporary hack, see https://github.com/aliyun/alibaba-cloud-sdk-go/issues/208
	req.SetDomain(fmt.Sprintf("cr.%s.aliyuncs.com", c.client.RegionId))
	req.Namespace = namespaceName

	var resp crDefaultResponse

	if err := invoker.Run(func() error {
		raw, err := c.client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.DeleteNamespace(req)
		})
		if err != nil {
			return err
		}
		err = json.Unmarshal(raw.(*cr.DeleteNamespaceResponse).GetHttpContentBytes(), &resp)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("deleting namespace got an error: %#v", err)
	}
	return &resp, nil
}
