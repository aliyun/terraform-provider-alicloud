package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliyunApigatewayGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := cloudapi.CreateCreateApiGroupRequest()
	args.GroupName = d.Get("name").(string)
	args.Description = d.Get("description").(string)

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateApiGroup(args)
		})
		if err != nil {
			if IsExceptedError(err, RepeatedCommit) {
				return resource.RetryableError(fmt.Errorf("Create api group got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Create api group got an error: %#v.", err))
		}
		resp, _ := raw.(*cloudapi.CreateApiGroupResponse)
		d.SetId(resp.GroupId)
		return nil
	}); err != nil {
		return fmt.Errorf("Creating apigatway group error: %#v", err)
	}

	return resourceAliyunApigatewayGroupRead(d, meta)
}

func resourceAliyunApigatewayGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	apiGroup, err := cloudApiService.DescribeApiGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", apiGroup.GroupName)
	d.Set("description", apiGroup.Description)

	return nil
}

func resourceAliyunApigatewayGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("name") || d.HasChange("description") {
		req := cloudapi.CreateModifyApiGroupRequest()
		req.Description = d.Get("description").(string)
		req.GroupName = d.Get("name").(string)
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApiGroup(req)
		})
		if err != nil {
			return fmt.Errorf("ModifyApiGroup got an error: %#v", err)
		}
	}
	return resourceAliyunApigatewayGroupRead(d, meta)
}

func resourceAliyunApigatewayGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	req := cloudapi.CreateDeleteApiGroupRequest()
	req.GroupId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DeleteApiGroup(req)
		})
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error deleting ApiGroup failed: %#v", err))
		}

		if _, err := cloudApiService.DescribeApiGroup(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error describing apiGroup failed when deleting apiGroup: %#v", err))
		}
		return resource.RetryableError(fmt.Errorf("Delete ApiGroup %s timeout.", d.Id()))
	})
}
