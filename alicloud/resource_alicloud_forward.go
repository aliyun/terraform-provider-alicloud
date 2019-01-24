package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunForwardEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunForwardEntryCreate,
		Read:   resourceAliyunForwardEntryRead,
		Update: resourceAliyunForwardEntryUpdate,
		Delete: resourceAliyunForwardEntryDelete,

		Schema: map[string]*schema.Schema{
			"forward_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateForwardPort,
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"tcp", "udp", "any"}),
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateForwardPort,
			},
		},
	}
}

func resourceAliyunForwardEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	args := vpc.CreateCreateForwardEntryRequest()
	args.RegionId = string(client.Region)
	args.ForwardTableId = d.Get("forward_table_id").(string)
	args.ExternalIp = d.Get("external_ip").(string)
	args.ExternalPort = d.Get("external_port").(string)
	args.IpProtocol = d.Get("ip_protocol").(string)
	args.InternalIp = d.Get("internal_ip").(string)
	args.InternalPort = d.Get("internal_port").(string)

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		ar := args
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateForwardEntry(ar)
		})
		if err != nil {
			if IsExceptedError(err, InvalidIpNotInNatgw) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, "forward_entry", args.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, "forward_entry", args.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		resp, _ := raw.(*vpc.CreateForwardEntryResponse)
		d.SetId(resp.ForwardEntryId)
		d.Set("forward_table_id", d.Get("forward_table_id").(string))
		return nil
	}); err != nil {
		return WrapError(err)
	}

	if err := vpcService.WaitForForwardEntry(args.ForwardTableId, d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunForwardEntryRead(d, meta)
}

func resourceAliyunForwardEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	forwardEntry, err := vpcService.DescribeForwardEntry(d.Get("forward_table_id").(string), d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("forward_table_id", forwardEntry.ForwardTableId)
	d.Set("external_ip", forwardEntry.ExternalIp)
	d.Set("external_port", forwardEntry.ExternalPort)
	d.Set("ip_protocol", forwardEntry.IpProtocol)
	d.Set("internal_ip", forwardEntry.InternalIp)
	d.Set("internal_port", forwardEntry.InternalPort)

	return nil
}

func resourceAliyunForwardEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	forwardEntry, err := vpcService.DescribeForwardEntry(d.Get("forward_table_id").(string), d.Id())
	if err != nil {
		return WrapError(err)
	}

	attributeUpdate := false
	args := vpc.CreateModifyForwardEntryRequest()
	args.RegionId = string(client.Region)
	args.ForwardTableId = forwardEntry.ForwardTableId
	args.ForwardEntryId = forwardEntry.ForwardEntryId
	args.ExternalIp = forwardEntry.ExternalIp
	args.IpProtocol = forwardEntry.IpProtocol
	args.ExternalPort = forwardEntry.ExternalPort
	args.InternalIp = forwardEntry.InternalIp
	args.InternalPort = forwardEntry.InternalPort

	if d.HasChange("external_port") {
		args.ExternalPort = d.Get("external_port").(string)
		attributeUpdate = true
	}

	if d.HasChange("ip_protocol") {
		args.IpProtocol = d.Get("ip_protocol").(string)
		attributeUpdate = true
	}

	if d.HasChange("internal_port") {
		args.InternalPort = d.Get("internal_port").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyForwardEntry(args)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if err := vpcService.WaitForForwardEntry(args.ForwardTableId, d.Id(), Available, DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliyunForwardEntryRead(d, meta)
}

func resourceAliyunForwardEntryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	args := vpc.CreateDeleteForwardEntryRequest()
	args.RegionId = string(client.Region)
	args.ForwardTableId = d.Get("forward_table_id").(string)
	args.ForwardEntryId = d.Id()

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteForwardEntry(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidForwardEntryIdNotFound, InvalidForwardTableIdNotFound}) {
				return nil
			}
			if IsExceptedErrors(err, []string{UnknownError}) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), args.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), args.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		if _, err := vpcService.DescribeForwardEntry(d.Get("forward_table_id").(string), d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), args.GetActionName(), ProviderERROR))

		return nil
	})
}
