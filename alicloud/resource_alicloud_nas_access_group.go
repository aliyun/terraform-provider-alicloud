package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"time"
)

func resourceAlicloudNAS_AccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNAS_AccessGroupCreate,
		Read:   resourceAlicloudNAS_AccessGroupRead,
		Update: resourceAlicloudNAS_AccessGroupUpdate,
		Delete: resourceAlicloudNAS_AccessGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"accessgroup_name": &schema.Schema{
				Type:         schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accessgroup_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Vpc",
					"Classic",
				}),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,

			},
		},
	}
}

func resourceAlicloudNAS_AccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateCreateAccessGroupRequest()
	request.RegionId = string(client.Region)
	request.AccessGroupName = d.Get("accessgroup_name").(string)
	request.AccessGroupType = d.Get("accessgroup_type").(string)
	request.Description = d.Get("description").(string)
	raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
		return nasClient.CreateAccessGroup(request)
	})
	fs, _ := raw.(*nas.CreateAccessGroupResponse)
	if err != nil {
		return fmt.Errorf("Error Waitting for NAS available: %#v", err)
	}

	d.SetId(fs.AccessGroupName)
	return resourceAlicloudNAS_AccessGroupUpdate(d, meta)
}

func resourceAlicloudNAS_AccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	attributeUpdate := false
	request := nas.CreateModifyAccessGroupRequest()
	request.AccessGroupName = d.Id()
	if d.HasChange("description") {
		attributeUpdate = true
		d.SetPartial("description")
		request.Description = d.Get("description").(string)
	}
	if attributeUpdate {
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.ModifyAccessGroup(request)
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAlicloudNAS_AccessGroupRead(d, meta)
}

func resourceAlicloudNAS_AccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}

	resp, err := nasService.DescribeAccessGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("accessgroup_name", resp.AccessGroupName)
	d.Set("accessgroup_type", resp.AccessGroupType)
	d.Set("RuleCount", resp.RuleCount)
	d.Set("MountTargetCount", resp.MountTargetCount)
	if resp.Description != "" {
		d.Set("description", resp.Description)
	}

	return nil
}

func resourceAlicloudNAS_AccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
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
			return resource.RetryableError(fmt.Errorf("Delete NAS timeout and got an error: %#v.", err))
		}

		if _, err := nasService.DescribeAccessGroup(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

