package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunApigatewayGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayGroupCreate,
		Read:   resourceAliyunApigatewayGroupRead,
		Update: resourceAliyunApigatewayGroupUpdate,
		Delete: resourceAliyunApigatewayGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunApigatewayGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	args := cloudapi.CreateCreateApiGroupRequest()

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		args.GroupName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.Description = v.(string)
	}

	resp, err := client.cloudapiconn.CreateApiGroup(args)

	if err != nil {
		return fmt.Errorf("Creating apigatway group error: %#v", err)
	}

	d.SetId(resp.GroupId)

	return resourceAliyunApigatewayGroupRead(d, meta)
}

func resourceAliyunApigatewayGroupRead(d *schema.ResourceData, meta interface{}) error {
	apiGroup, err := meta.(*AliyunClient).DescribeApiGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("group_id", apiGroup.GroupId)
	d.Set("name", apiGroup.GroupName)
	d.Set("description", apiGroup.Description)

	/*
		d.Set("groupId", apiGroup.GroupId)
		d.Set("groupName", apiGroup.GroupName)
		d.Set("SubDomain", apiGroup.SubDomain)
		d.Set("Description", apiGroup.Description)
		d.Set("RegionId", apiGroup.RegionId)
		d.Set("status", apiGroup.Status)
		d.Set("billingStatus", apiGroup.BillingStatus)
		d.Set("illegalStatus", apiGroup.IllegalStatus)
		d.Set("trafficLimit", apiGroup.TrafficLimit)
		d.Set("customDomains", apiGroup.CustomDomains)
		d.Set("stageItems", apiGroup.StageItems)
	*/

	return nil
}

func resourceAliyunApigatewayGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	cloudapiconn := meta.(*AliyunClient).cloudapiconn

	if d.HasChange("name") || d.HasChange("description") {
		req := cloudapi.CreateModifyApiGroupRequest()
		req.Description = d.Get("description").(string)
		req.GroupName = d.Get("name").(string)
		if _, err := cloudapiconn.ModifyApiGroup(req); err != nil {
			return fmt.Errorf("ModifyApiGroup got an error: %#v", err)
		}
	}
	return resourceAliyunApigatewayGroupRead(d, meta)
}

func resourceAliyunApigatewayGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	req := cloudapi.CreateDeleteApiGroupRequest()
	req.GroupId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.cloudapiconn.DeleteApiGroup(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error deleting ApiGroup failed: %#v", err))
		}

		if _, err := client.DescribeApiGroup(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error describing apiGroup failed when deleting apiGroup: %#v", err))
		}
		return resource.RetryableError(fmt.Errorf("Delete ApiGroup %s timeout.", d.Id()))
	})
}
