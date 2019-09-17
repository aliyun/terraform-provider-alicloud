package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunBgpNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunBgpNetworkCreate,
		Read:   resourceAliyunBgpNetworkRead,
		Delete: resourceAliyunBgpNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"dst_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunBgpNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bgpNetworkService := BgpService{client}
	request := vpc.CreateAddBgpNetworkRequest()
	request.RegionId = client.RegionId
	request.DstCidrBlock = d.Get("dst_cidr_block").(string)
	request.RouterId = d.Get("router_id").(string)
	request.ClientToken = buildClientToken(request.GetActionName())

	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request.VpcId = v.(string)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.AddBgpNetwork(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bgp_network", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	id := d.Get("router_id").(string) + ":" + d.Get("dst_cidr_block").(string)
	d.SetId(id)

	if err := bgpNetworkService.WaitForBgpNetwork(d.Id(), Active, 2*DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return resourceAliyunBgpNetworkRead(d, meta)
}

func resourceAliyunBgpNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bgpNetworkService := BgpService{client}

	object, err := bgpNetworkService.DescribeBgpNetwork(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("dst_cidr_block", object.DstCidrBlock)
	d.Set("router_id", object.RouterId)
	d.Set("vpc_id", object.VpcId)
	d.Set("status", object.Status)

	return nil
}

func resourceAliyunBgpNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bgpNetworkService := BgpService{client}

	request := vpc.CreateDeleteBgpNetworkRequest()
	request.RegionId = client.RegionId
	request.DstCidrBlock = d.Get("dst_cidr_block").(string)
	request.RouterId = d.Get("router_id").(string)
	request.ClientToken = buildClientToken(request.GetActionName())

	wait := incrementalWait(5*time.Second, 5*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteBgpNetwork(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{BgpInvalidStatus}) {
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
	return WrapError(bgpNetworkService.WaitForBgpNetwork(d.Id(), Deleted, DefaultTimeoutMedium))
}
