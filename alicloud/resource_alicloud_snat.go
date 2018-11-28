package alicloud

import (
	"fmt"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSnatEntryCreate,
		Read:   resourceAliyunSnatEntryRead,
		Update: resourceAliyunSnatEntryUpdate,
		Delete: resourceAliyunSnatEntryDelete,

		Schema: map[string]*schema.Schema{
			"snat_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snat_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliyunSnatEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateCreateSnatEntryRequest()
	request.RegionId = string(client.Region)
	request.SnatTableId = d.Get("snat_table_id").(string)
	request.SourceVSwitchId = d.Get("source_vswitch_id").(string)
	request.SnatIp = d.Get("snat_ip").(string)

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ar := request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateSnatEntry(ar)
		})
		if err != nil {
			if IsExceptedError(err, EIP_NOT_IN_GATEWAY) {
				return resource.RetryableError(fmt.Errorf("CreateSnatEntry timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("CreateSnatEntry got error: %#v.", err))
		}
		resp, _ := raw.(*vpc.CreateSnatEntryResponse)
		d.SetId(resp.SnatEntryId)
		return nil
	}); err != nil {
		return err
	}

	return resourceAliyunSnatEntryRead(d, meta)
}

func resourceAliyunSnatEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	snatEntry, err := vpcService.DescribeSnatEntry(d.Get("snat_table_id").(string), d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("snat_table_id", snatEntry.SnatTableId)
	d.Set("source_vswitch_id", snatEntry.SourceVSwitchId)
	d.Set("snat_ip", snatEntry.SnatIp)

	return nil
}

func resourceAliyunSnatEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	snatEntry, err := vpcService.DescribeSnatEntry(d.Get("snat_table_id").(string), d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)
	attributeUpdate := false
	request := vpc.CreateModifySnatEntryRequest()
	request.RegionId = string(client.Region)
	request.SnatTableId = snatEntry.SnatTableId
	request.SnatEntryId = snatEntry.SnatEntryId

	if d.HasChange("snat_ip") {
		d.SetPartial("snat_ip")
		var snat_ip string
		if v, ok := d.GetOk("snat_ip"); ok {
			snat_ip = v.(string)
		} else {
			return fmt.Errorf("cann't change snap_ip to empty string")
		}
		request.SnatIp = snat_ip

		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifySnatEntry(request)
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAliyunSnatEntryRead(d, meta)
}

func resourceAliyunSnatEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateDeleteSnatEntryRequest()
	request.RegionId = string(client.Region)
	request.SnatTableId = d.Get("snat_table_id").(string)
	request.SnatEntryId = d.Id()

	_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteSnatEntry(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidSnatTableIdNotFound) {
			return nil
		}
		return err
	}

	return nil
}
