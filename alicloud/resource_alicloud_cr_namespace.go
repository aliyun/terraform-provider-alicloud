package alicloud

import (
	"encoding/json"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

func resourceAlicloudCRNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	namespaceName := d.Get("name").(string)

	payload := &crCreateNamespaceRequestPayload{}
	payload.Namespace.Namespace = namespaceName
	serialized, err := json.Marshal(payload)
	if err != nil {
		return WrapError(err)
	}

	req := cr.CreateCreateNamespaceRequest()
	req.SetContent(serialized)

	if err := invoker.Run(func() error {
		_, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.CreateNamespace(req)
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_namespace", req.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(namespaceName)

	return resourceAlicloudCRNamespaceUpdate(d, meta)
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
			return WrapError(err)
		}
		req := cr.CreateUpdateNamespaceRequest()
		req.SetContent(serialized)
		req.Namespace = d.Get("name").(string)

		if err := invoker.Run(func() error {
			_, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
				return crClient.UpdateNamespace(req)
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCRNamespaceRead(d, meta)
}

func resourceAlicloudCRNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	raw, err := crService.DescribeNamespace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	var resp crDescribeNamespaceResponse
	err = json.Unmarshal(raw.GetHttpContentBytes(), &resp)
	if err != nil {
		return WrapError(err)
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
		req.Namespace = d.Id()

		if err := invoker.Run(func() error {
			_, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
				return crClient.DeleteNamespace(req)
			})
			return err
		}); err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorNamespaceNotExist) {
				return nil
			}
			return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		return nil
	})
}
