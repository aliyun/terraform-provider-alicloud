package cr

import (
	"encoding/json"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	crsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		DeprecationMessage: "This resource has been deprecated since v1.276.0 and will be removed in the future. Please use 'alicloud_cr_ee_namespace' instead.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 30),
			},
			"auto_create": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"default_visibility": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLIC", "PRIVATE"}, false),
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
		return errs.WrapError(err)
	}

	request := crsdk.CreateCreateNamespaceRequest()
	request.SetContent(serialized)

	raw, err := client.WithCrClient(func(crClient *crsdk.Client) (interface{}, error) {
		return crClient.CreateNamespace(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cr_namespace", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RoaRequest, request)

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
			return errs.WrapError(err)
		}
		request := crsdk.CreateUpdateNamespaceRequest()
		request.RegionId = client.RegionId
		request.SetContent(serialized)
		request.Namespace = d.Get("name").(string)

		raw, err := client.WithCrClient(func(crClient *crsdk.Client) (interface{}, error) {
			return crClient.UpdateNamespace(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RoaRequest, request)
	}
	return resourceAlicloudCRNamespaceRead(d, meta)
}

func resourceAlicloudCRNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	object, err := crService.DescribeCrNamespace(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	var response crDescribeNamespaceResponse
	err = json.Unmarshal(object.GetHttpContentBytes(), &response)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("name", response.Data.Namespace.Namespace)
	d.Set("auto_create", response.Data.Namespace.AutoCreate)
	d.Set("default_visibility", response.Data.Namespace.DefaultVisibility)

	return nil
}

func resourceAlicloudCRNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	request := crsdk.CreateDeleteNamespaceRequest()
	request.Namespace = d.Id()

	raw, err := client.WithCrClient(func(crClient *crsdk.Client) (interface{}, error) {
		return crClient.DeleteNamespace(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NAMESPACE_NOT_EXIST"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RoaRequest, request)
	return errs.WrapError(crService.WaitForCRNamespace(d.Id(), Deleted, DefaultTimeout))
}
