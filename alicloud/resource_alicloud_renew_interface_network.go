package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"time"
)

func resourceAlicloudRenewInterfaceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRenewInterfaceNetworkCreate,
		Read:   resourceAlicloudRenewInterfaceNetworkRead,
		Update: resourceAlicloudRenewInterfaceNetworkUpdate,
		Delete: resourceAlicloudRenewInterfaceNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"pricing_cycle": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_owner_account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudRenewInterfaceNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateRenewInstanceRequest()
	request.RegionId = client.RegionId
	request.Scheme = "https"
	request.InstanceId = d.Get("instance_id").(string)
	request.PricingCycle = d.Get("pricing_cycle").(string)
	request.Duration = requests.NewInteger(d.Get("duration").(int))
	request.InstanceType = d.Get("instance_type").(string)

	if v, ok := d.GetOk("owner_id"); ok {
		request.OwnerId = requests.NewInteger(v.(int))
	}
	if v := d.Get("resource_owner_account").(string); v != "" {
		request.ResourceOwnerAccount = v
	}
	if v := d.Get("resource_owner_id"); v != "" {
		request.ResourceOwnerId = requests.NewInteger(v.(int))
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.RenewInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "RenewNetworkInstance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*vpc.RenewInstanceResponse)
	if response != nil {
		d.SetId(d.Get("instance_id").(string))
	} else {
		return WrapError(Error("The response is nil"))
	}
	time.Sleep(10 * time.Second)
	return resourceAlicloudRenewInterfaceNetworkRead(d, meta)
}

func resourceAlicloudRenewInterfaceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlicloudRenewInterfaceNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	wait := incrementalWait(5*time.Second, 5*time.Second)

	if d.HasChange("instance_id") || d.HasChange("pricing_cycle") || d.HasChange("duration") || d.HasChange("instance_type") ||
		d.HasChange("owner_id") || d.HasChange("resource_owner_account") || d.HasChange("resource_owner_id") {
		request := vpc.CreateModifyInstanceAutoRenewalAttributeRequest()
		request.RegionId = client.RegionId
		request.Scheme = "https"
		request.InstanceId = d.Get("instance_id").(string)
		request.PricingCycle = d.Get("pricing_cycle").(string)
		request.Duration = requests.NewInteger(d.Get("duration").(int))
		request.InstanceType = d.Get("instance_type").(string)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.ModifyInstanceAutoRenewalAttribute(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{RenewIncorrectStatus}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("instance_id")
		d.SetPartial("pricing_cycle")
		d.SetPartial("duration")
		d.SetPartial("instance_type")
		d.SetPartial("owner_id")
		d.SetPartial("resource_owner_account")
		d.SetPartial("resource_owner_id")
	}
	d.Partial(false)
	return resourceAlicloudRenewInterfaceNetworkRead(d, meta)
}

func resourceAlicloudRenewInterfaceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
