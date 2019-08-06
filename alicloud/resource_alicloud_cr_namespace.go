package alicloud

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
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
				ForceNew:     true,
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

	namespaceName := d.Get("name").(string)

	payload := &crCreateNamespaceRequestPayload{}
	payload.Namespace.Namespace = namespaceName
	serialized, err := json.Marshal(payload)
	if err != nil {
		return WrapError(err)
	}

	request := cr.CreateCreateNamespaceRequest()
	request.SetContent(serialized)

	raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
		return crClient.CreateNamespace(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_namespace", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	d.SetId(namespaceName)

	return resourceAlicloudCRNamespaceUpdate(d, meta)
}

func resourceAlicloudCRNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("auto_create") || d.HasChange("default_visibility") {
		payload := &crUpdateNamespaceRequestPayload{}
		payload.Namespace.DefaultVisibility = d.Get("default_visibility").(string)
		payload.Namespace.AutoCreate = d.Get("auto_create").(bool)

		serialized, err := json.Marshal(payload)
		if err != nil {
			return WrapError(err)
		}
		request := cr.CreateUpdateNamespaceRequest()
		request.SetContent(serialized)
		request.Namespace = d.Get("name").(string)

		raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			return crClient.UpdateNamespace(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}
	return resourceAlicloudCRNamespaceRead(d, meta)
}

func resourceAlicloudCRNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	object, err := crService.DescribeCrNamespace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	var response crDescribeNamespaceResponse
	err = json.Unmarshal(object.GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}

	d.Set("name", response.Data.Namespace.Namespace)
	d.Set("auto_create", response.Data.Namespace.AutoCreate)
	d.Set("default_visibility", response.Data.Namespace.DefaultVisibility)

	return nil
}

func resourceAlicloudCRNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	request := cr.CreateDeleteNamespaceRequest()
	request.Namespace = d.Id()

	raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
		return crClient.DeleteNamespace(request)
	})
	if err != nil {
		if IsExceptedError(err, ErrorNamespaceNotExist) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return WrapError(crService.WaitForCRNamespace(d.Id(), Deleted, DefaultTimeout))
}
