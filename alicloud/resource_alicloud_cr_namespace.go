package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"time"
)

func resourceAlicloudCRNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCRNamespaceCreate,
		Read:   resourceAlicloudCRNamespaceRead,
		Update: resourceAlicloudCRNamespaceUpdate,
		Delete: resourceAlicloudCRNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateContainerRegistryNamespaceName,
			},
			"auto_create": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"default_visibility": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"PUBLIC", "PRIVATE"}),
			},
		},
	}
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

func resourceAlicloudCRNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	namespaceName := d.Get("name").(string)

	payload := &crCreateNamespaceRequestPayload{}
	payload.Namespace.Namespace = namespaceName
	serialized, err := json.Marshal(payload)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "new", "json marshal", ProviderERROR)
	}

	req := cr.CreateCreateNamespaceRequest()
	req.SetContent(serialized)
	// FIXME
	// Temporary hack, see https://github.com/aliyun/alibaba-cloud-sdk-go/issues/208
	req.SetDomain(fmt.Sprintf("cr.%s.aliyuncs.com", client.RegionId))

	if err := invoker.Run(func() error {
		_, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.CreateNamespace(req)
		})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "new", req.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(namespaceName)

	return resourceAlicloudCRNamespaceUpdate(d, meta)
}

type crUpdateNamespaceRequestPayload struct {
	Namespace struct {
		AutoCreate        bool   `json:"AutoCreate"`
		DefaultVisibility string `json:"DefaultVisibility"`
	} `json:"Namespace"`
}

func resourceAlicloudCRNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	if d.HasChange("auto_create") || d.HasChange("default_visibility") {
		payload := &crUpdateNamespaceRequestPayload{}
		payload.Namespace.DefaultVisibility = d.Get("default_visibility").(string)
		payload.Namespace.AutoCreate = d.Get("auto_create").(bool)

		serialized, err := json.Marshal(payload)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "new", "json marshal", ProviderERROR)
		}
		req := cr.CreateUpdateNamespaceRequest()
		// FIXME
		// Temporary hack, see https://github.com/aliyun/alibaba-cloud-sdk-go/issues/208
		req.SetDomain(fmt.Sprintf("cr.%s.aliyuncs.com", client.RegionId))
		req.SetContent(serialized)
		req.Namespace = d.Get("name").(string)

		if err := invoker.Run(func() error {
			_, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
				return crClient.UpdateNamespace(req)
			})
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCRNamespaceRead(d, meta)
}

func resourceAlicloudCRNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	resp, err := crService.DescribeNamespace(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", resp.Data.Namespace.Namespace)
	d.Set("auto_create", resp.Data.Namespace.AutoCreate)
	d.Set("default_visibility", resp.Data.Namespace.DefaultVisibility)

	return nil
}

func resourceAlicloudCRNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		req := cr.CreateDeleteNamespaceRequest()
		// FIXME
		// Temporary hack, see https://github.com/aliyun/alibaba-cloud-sdk-go/issues/208
		req.SetDomain(fmt.Sprintf("cr.%s.aliyuncs.com", client.RegionId))
		req.Namespace = d.Get("name").(string)

		var resp crDefaultResponse

		if err := invoker.Run(func() error {
			raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
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
			if NotFoundError(err) || IsExceptedError(err, ErrorNamespaceNotExist) {
				return nil
			}
			return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		return nil
	})
}
