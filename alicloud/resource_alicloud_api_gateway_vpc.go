package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunApigatewayVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayVpcCreate,
		Read:   resourceAliyunApigatewayVpcRead,
		Delete: resourceAliyunApigatewayVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunApigatewayVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := cloudapi.CreateSetVpcAccessRequest()
	args.Name = d.Get("name").(string)
	args.VpcId = d.Get("vpc_id").(string)
	args.InstanceId = d.Get("instance_id").(string)
	args.Port = requests.NewInteger(d.Get("port").(int))
	_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.SetVpcAccess(args)
	})
	if err != nil {
		return fmt.Errorf("Creating apigatway Vpc error: %#v", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s%s%s", args.Name, COLON_SEPARATED, args.VpcId, COLON_SEPARATED, args.InstanceId, COLON_SEPARATED, args.Port))
	return resourceAliyunApigatewayVpcRead(d, meta)
}

func resourceAliyunApigatewayVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	vpc, err := cloudApiService.DescribeVpcAccess(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", vpc.Name)
	d.Set("vpc_id", vpc.VpcId)
	d.Set("instance_id", vpc.InstanceId)
	d.Set("port", vpc.Port)

	return nil
}

func resourceAliyunApigatewayVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}
	req := cloudapi.CreateRemoveVpcAccessRequest()
	req.VpcId = d.Get("vpc_id").(string)
	req.InstanceId = d.Get("instance_id").(string)
	req.Port = requests.NewInteger(d.Get("port").(int))

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.RemoveVpcAccess(req)
		})
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error deleting Vpc failed: %#v", err))
		}

		if _, err := cloudApiService.DescribeVpcAccess(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error describing Vpc failed when deleting Vpc: %#v", err))
		}
		return resource.RetryableError(fmt.Errorf("Delete Vpc %s timeout.", d.Id()))
	})
}
