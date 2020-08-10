package alicloud

import (
	"fmt"

	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunSnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSnatEntryCreate,
		Read:   resourceAliyunSnatEntryRead,
		Update: resourceAliyunSnatEntryUpdate,
		Delete: resourceAliyunSnatEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"snat_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_vswitch_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: strings.Fields("source_cidr"),
			},
			"snat_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_cidr": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: strings.Fields("source_vswitch_id"),
			},
			"snat_entry_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snat_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunSnatEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateCreateSnatEntryRequest()
	request.RegionId = string(client.Region)
	request.SnatTableId = d.Get("snat_table_id").(string)
	request.SourceVSwitchId = d.Get("source_vswitch_id").(string)
	request.SnatIp = d.Get("snat_ip").(string)
	if v, ok := d.GetOk("source_cidr"); ok {
		request.SourceCIDR = v.(string)
	}

	if v, ok := d.GetOk("snat_entry_name"); ok {
		request.SnatEntryName = v.(string)
	}
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ar := request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateSnatEntry(ar)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"EIP_NOT_IN_GATEWAY", "OperationUnsupported.EipNatBWPCheck", "OperationUnsupported.EipInBinding"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.CreateSnatEntryResponse)
		d.SetId(fmt.Sprintf("%s%s%s", request.SnatTableId, COLON_SEPARATED, response.SnatEntryId))
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_snat_entry", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if err := vpcService.WaitForSnatEntry(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunSnatEntryRead(d, meta)
}

func resourceAliyunSnatEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	// compatible with previous id which in under 1.37.0
	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}

	object, err := vpcService.DescribeSnatEntry(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("snat_table_id", object.SnatTableId)
	if _, ok := d.GetOk("source_cidr"); ok {
		d.Set("source_cidr", object.SourceCIDR)
	} else {
		d.Set("source_vswitch_id", object.SourceVSwitchId)
	}
	d.Set("snat_ip", object.SnatIp)
	d.Set("snat_entry_id", object.SnatEntryId)
	d.Set("snat_entry_name", object.SnatEntryName)

	return nil
}

func resourceAliyunSnatEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	// compatible with previous id which in under 1.37.0
	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := vpc.CreateModifySnatEntryRequest()
	request.RegionId = string(client.Region)
	request.SnatTableId = parts[0]
	request.SnatEntryId = parts[1]
	update := false

	if d.HasChange("snat_entry_name") {
		update = true
		request.SnatEntryName = d.Get("snat_entry_name").(string)
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifySnatEntry(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err := vpcService.WaitForSnatEntry(d.Id(), Available, DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliyunSnatEntryRead(d, meta)
}

func resourceAliyunSnatEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	// compatible with previous id which in under 1.37.0
	if strings.HasPrefix(d.Id(), "snat-") {
		d.SetId(fmt.Sprintf("%s%s%s", d.Get("snat_table_id").(string), COLON_SEPARATED, d.Id()))
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := vpc.CreateDeleteSnatEntryRequest()
	request.RegionId = string(client.Region)
	request.SnatTableId = parts[0]
	request.SnatEntryId = parts[1]
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteSnatEntry(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorretSnatEntryStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSnatTableId.NotFound", "InvalidSnatEntryId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForSnatEntry(d.Id(), Deleted, DefaultTimeout))
}
