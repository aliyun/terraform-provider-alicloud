package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/oos"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudOosTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosTemplateCreate,
		Read:   resourceAlicloudOosTemplateRead,
		Update: resourceAlicloudOosTemplateUpdate,
		Delete: resourceAlicloudOosTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_delete_executions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_trigger": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"template_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudOosTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := oos.CreateCreateTemplateRequest()
	request.Content = d.Get("content").(string)
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("tags"); ok {
		request.Tags = v.(map[string]interface{})
	}
	request.TemplateName = d.Get("template_name").(string)
	if v, ok := d.GetOk("version_name"); ok {
		request.VersionName = v.(string)
	}

	raw, err := client.WithOosClient(func(oosClient *oos.Client) (interface{}, error) {
		return oosClient.CreateTemplate(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_template", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*oos.CreateTemplateResponse)
	d.SetId(fmt.Sprintf("%v", response.Template.TemplateName))

	return resourceAlicloudOosTemplateRead(d, meta)
}
func resourceAlicloudOosTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("template_name", d.Id())
	d.Set("created_by", object.CreatedBy)
	d.Set("created_date", object.CreatedDate)
	d.Set("description", object.Description)
	d.Set("has_trigger", object.HasTrigger)
	d.Set("share_type", object.ShareType)
	d.Set("tags", object.Tags)
	d.Set("template_format", object.TemplateFormat)
	d.Set("template_id", object.TemplateId)
	d.Set("template_type", object.TemplateType)
	d.Set("template_version", object.TemplateVersion)
	d.Set("updated_by", object.UpdatedBy)
	d.Set("updated_date", object.UpdatedDate)
	return nil
}
func resourceAlicloudOosTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := oos.CreateUpdateTemplateRequest()
	request.TemplateName = d.Id()
	if d.HasChange("content") {
		update = true
	}
	request.Content = d.Get("content").(string)
	request.RegionId = client.RegionId
	if d.HasChange("tags") {
		update = true
		request.Tags = d.Get("tags").(map[string]interface{})
	}
	if d.HasChange("version_name") {
		update = true
		request.VersionName = d.Get("version_name").(string)
	}
	if update {
		raw, err := client.WithOosClient(func(oosClient *oos.Client) (interface{}, error) {
			return oosClient.UpdateTemplate(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudOosTemplateRead(d, meta)
}
func resourceAlicloudOosTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := oos.CreateDeleteTemplateRequest()
	request.TemplateName = d.Id()
	if v, ok := d.GetOkExists("auto_delete_executions"); ok {
		request.AutoDeleteExecutions = requests.NewBoolean(v.(bool))
	}
	request.RegionId = client.RegionId
	raw, err := client.WithOosClient(func(oosClient *oos.Client) (interface{}, error) {
		return oosClient.DeleteTemplate(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Template"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
