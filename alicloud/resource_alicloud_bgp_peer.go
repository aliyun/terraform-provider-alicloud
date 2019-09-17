package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunBgpPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunBgpPeerCreate,
		Read:   resourceAliyunBgpPeerRead,
		Delete: resourceAliyunBgpPeerDelete,
		Update: resourceAliyunBgpPeerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bgp_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_bfd": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"peer_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bgp_peer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliyunBgpPeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bgpPeerService := BgpService{client}
	request := vpc.CreateCreateBgpPeerRequest()
	request.RegionId = client.RegionId
	request.BgpGroupId = d.Get("bgp_group_id").(string)
	request.EnableBfd = requests.NewBoolean(d.Get("enable_bfd").(bool))
	if v, ok := d.GetOk("client_token"); ok && v.(string) != "" {
		request.ClientToken = v.(string)
	} else {
		request.ClientToken = buildClientToken(request.GetActionName())
	}

	if v, ok := d.GetOk("peer_ip_address"); ok && v.(string) != "" {
		request.PeerIpAddress = v.(string)
	}
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateBgpPeer(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bgp_peer", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*vpc.CreateBgpPeerResponse)
	d.SetId(response.BgpPeerId)

	if err := bgpPeerService.WaitForBgpPeer(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return resourceAliyunBgpPeerRead(d, meta)
}

func resourceAliyunBgpPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bgpPeerService := BgpService{client}

	object, err := bgpPeerService.DescribeBgpPeer(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bgp_group_id", object.BgpGroupId)
	d.Set("peer_ip_address", object.PeerIpAddress)
	d.Set("enable_bfd", object.EnableBfd)
	d.Set("bgp_peer_id", object.BgpPeerId)

	return nil
}

func resourceAliyunBgpPeerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false

	request := vpc.CreateModifyBgpPeerAttributeRequest()
	request.RegionId = client.RegionId
	request.BgpPeerId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())

	if d.HasChange("bgp_group_id") {
		request.BgpGroupId = d.Get("bgp_group_id").(string)
		d.SetPartial("bgp_group_id")
		update = true
	}

	if d.HasChange("peer_ip_address") {
		request.PeerIpAddress = d.Get("peer_ip_address").(string)
		d.SetPartial("peer_ip_address")
		update = true
	}

	if d.HasChange("enable_bfd") {
		request.EnableBfd = requests.NewBoolean(d.Get("enable_bfd").(bool))
		d.SetPartial("enable_bfd")
		update = true
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyBgpPeerAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAliyunBgpPeerRead(d, meta)
}

func resourceAliyunBgpPeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bgpPeerService := BgpService{client}
	if err := bgpPeerService.WaitForBgpPeer(d.Id(), Available, 2*DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	request := vpc.CreateDeleteBgpPeerRequest()
	request.RegionId = client.RegionId
	request.BgpPeerId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())

	wait := incrementalWait(5*time.Second, 5*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteBgpPeer(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{BgpInvalidStatus, BgpPeerDependencyViolation}) {
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
	return WrapError(bgpPeerService.WaitForBgpPeer(d.Id(), Deleted, DefaultTimeoutMedium))
}
