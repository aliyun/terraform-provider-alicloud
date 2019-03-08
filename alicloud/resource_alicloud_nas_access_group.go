package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudNasAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasAccessGroupCreate,
		Read:   resourceAlicloudNasAccessGroupRead,
		Update: resourceAlicloudNasAccessGroupUpdate,
		Delete: resourceAlicloudNasAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(Vpc), string(Classic)}),
				ForceNew:     true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudNasAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateCreateAccessGroupRequest()
	request.RegionId = string(client.Region)
	request.AccessGroupName = d.Get("name").(string)
	request.AccessGroupType = d.Get("type").(string)
	request.Description = d.Get("description").(string)
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateAccessGroup(request)
	})
	response, _ := raw.(*nas.CreateAccessGroupResponse)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "nas_access_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId(response.AccessGroupName)
	return resourceAlicloudNasAccessGroupRead(d, meta)
}

func resourceAlicloudNasAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := nas.CreateModifyAccessGroupRequest()
	request.AccessGroupName = d.Id()

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.ModifyAccessGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudNasAccessGroupRead(d, meta)
}

func resourceAlicloudNasAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}

	resp, err := nasService.DescribeNasAccessGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", resp.AccessGroupName)
	d.Set("type", resp.AccessGroupType)
	d.Set("description", resp.Description)

	return nil
}

func resourceAlicloudNasAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	request := nas.CreateDeleteAccessGroupRequest()
	request.AccessGroupName = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DeleteAccessGroup(request)
		})
		if err != nil {
			if IsExceptedError(err, ForbiddenNasNotFound) {
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		if _, err := nasService.DescribeNasAccessGroup(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
	})
}
