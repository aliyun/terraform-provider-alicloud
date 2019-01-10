package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunApigatewayApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayAppCreate,
		Read:   resourceAliyunApigatewayAppRead,
		Update: resourceAliyunApigatewayAppUpdate,
		Delete: resourceAliyunApigatewayAppDelete,
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
				Optional: true,
			},
		},
	}
}

func resourceAliyunApigatewayAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := cloudapi.CreateCreateAppRequest()
	args.AppName = d.Get("name").(string)
	if v, exist := d.GetOk("description"); exist {
		args.Description = v.(string)
	}
	args.Description = d.Get("description").(string)

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateApp(args)
		})
		if err != nil {
			if IsExceptedError(err, RepeatedCommit) {
				return resource.RetryableError(fmt.Errorf("Create app got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Create app got an error: %#v.", err))
		}
		resp, _ := raw.(*cloudapi.CreateAppResponse)
		d.SetId(strconv.Itoa(resp.AppId))
		return nil
	}); err != nil {
		return fmt.Errorf("Creating apigatway group error: %#v", err)
	}

	return resourceAliyunApigatewayAppRead(d, meta)
}

func resourceAliyunApigatewayAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := cloudapi.CreateDescribeAppRequest()
	args.AppId = requests.Integer(d.Id())

	if err := resource.Retry(3*time.Second, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeApp(args)
		})
		if err != nil {
			if IsExceptedError(err, NotFoundApp) {
				return resource.RetryableError(fmt.Errorf("Create app got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Create app got an error: %#v.", err))
		}
		resp := raw.(*cloudapi.DescribeAppResponse)

		d.Set("name", resp.AppName)
		d.Set("description", resp.Description)
		return nil
	}); err != nil {
		return fmt.Errorf("Describe apigatway App error: %#v", err)
	}
	return nil
}

func resourceAliyunApigatewayAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("name") || d.HasChange("description") {
		req := cloudapi.CreateModifyAppRequest()
		req.AppId = requests.Integer(d.Id())
		req.AppName = d.Get("name").(string)
		if v, exist := d.GetOk("description"); exist {
			req.Description = v.(string)
		}

		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApp(req)
		})
		if err != nil {
			return fmt.Errorf("Modify App got an error: %#v", err)
		}
	}
	time.Sleep(3 * time.Second)
	return resourceAliyunApigatewayAppRead(d, meta)
}

func resourceAliyunApigatewayAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	req := cloudapi.CreateDeleteAppRequest()
	req.AppId = requests.Integer(d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DeleteApp(req)
		})
		if err != nil {
			if IsExceptedError(err, NotFoundApp) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting App failed: %#v", err))
		}

		if _, err := cloudApiService.DescribeApp(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error describing App failed when deleting App: %#v", err))
		}
		return nil
	})
}
