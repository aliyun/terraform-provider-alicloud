package alicloud

import (
	"fmt"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunForwardEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunForwardEntryCreate,
		Read:   resourceAliyunForwardEntryRead,
		Update: resourceAliyunForwardEntryUpdate,
		Delete: resourceAliyunForwardEntryDelete,

		Schema: map[string]*schema.Schema{
			"forward_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_port": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateForwardPort,
			},
			"ip_protocol": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"tcp", "udp", "any"}),
			},
			"internal_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_port": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateForwardPort,
			},
		},
	}
}

func resourceAliyunForwardEntryCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := vpc.CreateCreateForwardEntryRequest()
	args.RegionId = string(getRegion(d, meta))
	args.ForwardTableId = d.Get("forward_table_id").(string)
	args.ExternalIp = d.Get("external_ip").(string)
	args.ExternalPort = d.Get("external_port").(string)
	args.IpProtocol = d.Get("ip_protocol").(string)
	args.InternalIp = d.Get("internal_ip").(string)
	args.InternalPort = d.Get("internal_port").(string)

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		ar := args
		resp, err := conn.CreateForwardEntry(ar)
		if err != nil {
			if IsExceptedError(err, InvalidIpNotInNatgw) {
				return resource.RetryableError(fmt.Errorf("CreateForwardEntry timeout and got error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("CreateNatGateway got error: %#v", err))
		}
		d.SetId(resp.ForwardEntryId)
		d.Set("forward_table_id", d.Get("forward_table_id").(string))
		return nil
	}); err != nil {
		return err
	}

	return resourceAliyunForwardEntryRead(d, meta)
}

func resourceAliyunForwardEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	forwardEntry, err := client.DescribeForwardEntry(d.Get("forward_table_id").(string), d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
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
	client := meta.(*AliyunClient)

	forwardEntry, err := client.DescribeForwardEntry(d.Get("forward_table_id").(string), d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)
	attributeUpdate := false
	args := vpc.CreateModifyForwardEntryRequest()
	args.RegionId = string(getRegion(d, meta))
	args.ForwardTableId = forwardEntry.ForwardTableId
	args.ForwardEntryId = forwardEntry.ForwardEntryId
	args.ExternalIp = forwardEntry.ExternalIp
	args.IpProtocol = forwardEntry.IpProtocol
	args.ExternalPort = forwardEntry.ExternalPort
	args.InternalIp = forwardEntry.InternalIp
	args.InternalPort = forwardEntry.InternalPort

	if d.HasChange("external_port") {
		d.SetPartial("external_port")
		args.ExternalPort = d.Get("external_port").(string)
		attributeUpdate = true
	}

	if d.HasChange("ip_protocol") {
		d.SetPartial("ip_protocol")
		args.IpProtocol = d.Get("ip_protocol").(string)
		attributeUpdate = true
	}

	if d.HasChange("internal_port") {
		d.SetPartial("internal_port")
		args.InternalPort = d.Get("internal_port").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		if _, err := client.vpcconn.ModifyForwardEntry(args); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAliyunForwardEntryRead(d, meta)
}

func resourceAliyunForwardEntryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	args := vpc.CreateDeleteForwardEntryRequest()
	args.RegionId = string(getRegion(d, meta))
	args.ForwardTableId = d.Get("forward_table_id").(string)
	args.ForwardEntryId = d.Id()

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		if _, err := client.vpcconn.DeleteForwardEntry(args); err != nil {
			if IsExceptedError(err, InvalidForwardEntryIdNotFound) ||
				IsExceptedError(err, InvalidForwardTableIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		forwardEntry, err := client.DescribeForwardEntry(d.Get("forward_table_id").(string), d.Id())

		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if forwardEntry.ForwardEntryId == d.Id() {
			return resource.RetryableError(fmt.Errorf("Delete Forward Entry timeout and got an error:%#v.", err))
		}

		return nil
	})
}
