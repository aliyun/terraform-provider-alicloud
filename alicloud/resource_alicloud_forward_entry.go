package alicloud

import (
	"strings"
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
				ForceNew: true,
				Required: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Required: true,
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forward_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunForwardEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateForwardEntryRequest()
	request.RegionId = string(client.Region)
	request.ForwardTableId = d.Get("forward_table_id").(string)
	request.ExternalIp = d.Get("external_ip").(string)
	request.ExternalPort = d.Get("external_port").(string)
	request.IpProtocol = d.Get("ip_protocol").(string)
	request.InternalIp = d.Get("internal_ip").(string)
	request.InternalPort = d.Get("internal_port").(string)
	if name, ok := d.GetOk("name"); ok {
		request.ForwardEntryName = name.(string)
	}
	var raw interface{}
	var err error
	if err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		ar := request
		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateForwardEntry(ar)
		})
		if err != nil {
			if IsExceptedError(err, InvalidIpNotInNatgw) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_forward_entry", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateForwardEntryResponse)
	d.SetId(request.ForwardTableId + COLON_SEPARATED + response.ForwardEntryId)
	if err := vpcService.WaitForForwardEntry(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return resourceAliyunForwardEntryRead(d, meta)
}

func resourceAliyunForwardEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}
	forwardEntry, err := vpcService.DescribeForwardEntry(d.Id())
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
	d.Set("name", forwardEntry.ForwardEntryName)
	d.Set("forward_entry_id", forwardEntry.ForwardEntryId)

	return nil
}

func resourceAliyunForwardEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := vpc.CreateModifyForwardEntryRequest()
	request.RegionId = string(client.Region)
	request.ForwardEntryId = parts[1]
	request.ForwardTableId = parts[0]

	if d.HasChange("external_ip") {
		request.ExternalIp = d.Get("external_ip").(string)
	}

	if d.HasChange("external_port") {
		request.ExternalPort = d.Get("external_port").(string)
	}

	if d.HasChange("ip_protocol") {
		request.IpProtocol = d.Get("ip_protocol").(string)
	}

	if d.HasChange("internal_ip") {
		request.InternalIp = d.Get("internal_ip").(string)
	}

	if d.HasChange("internal_port") {
		request.InternalPort = d.Get("internal_port").(string)
	}

	if d.HasChange("name") {
		request.ForwardEntryName = d.Get("name").(string)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyForwardEntry(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err := vpcService.WaitForForwardEntry(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunForwardEntryRead(d, meta)
}

func resourceAliyunForwardEntryDelete(d *schema.ResourceData, meta interface{}) error {
	if !strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(d.Get("forward_table_id").(string) + COLON_SEPARATED + d.Id())
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteForwardEntryRequest()
	request.RegionId = string(client.Region)
	request.ForwardTableId = parts[0]
	request.ForwardEntryId = parts[1]

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteForwardEntry(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{UnknownError}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidForwardEntryIdNotFound, InvalidForwardTableIdNotFound}) {
			return nil
		}
		WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForForwardEntry(d.Id(), Deleted, DefaultTimeout))
}
