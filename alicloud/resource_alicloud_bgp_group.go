package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceAliyunBgpGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunBgpGroupCreate,
		Read:   resourceAliyunBgpGroupRead,
		Update: resourceAliyunBgpGroupUpdate,
		Delete: resourceAliyunBgpGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_asn": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"auth_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_fake_asn": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliyunBgpGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	BgpGroupService := BgpService{client}

	request := vpc.CreateCreateBgpGroupRequest()
	request.RegionId = client.RegionId
	request.PeerAsn = requests.NewInteger(d.Get("peer_asn").(int))
	request.RouterId = d.Get("router_id").(string)
	request.ClientToken = buildClientToken(request.GetActionName())

	if v, ok := d.GetOk("is_fake_asn"); ok {
		request.IsFakeAsn = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("auth_key"); ok && v.(string) != "" {
		request.AuthKey = v.(string)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateBgpGroup(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bgp_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*vpc.CreateBgpGroupResponse)
	d.SetId(response.BgpGroupId)

	if err := BgpGroupService.WaitForBgpGroup(d.Id(), Available, 2*DefaultLongTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunBgpGroupRead(d, meta)
}

func resourceAliyunBgpGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	BgpGroupService := BgpService{client}

	object, err := BgpGroupService.DescribeBgpGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	peerAsn, err := strconv.Atoi(object.PeerAsn)
	if err == nil {
		d.Set("peer_asn", peerAsn)
	}
	d.Set("router_id", object.RouterId)
	d.Set("auth_key", object.AuthKey)
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	if object.IsFake == "true" {
		d.Set("is_fake_asn", true)
	} else {
		d.Set("is_fake_asn", false)
	}

	return nil
}

func resourceAliyunBgpGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false

	request := vpc.CreateModifyBgpGroupAttributeRequest()
	request.RegionId = client.RegionId
	request.BgpGroupId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
	}
	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if d.HasChange("is_fake_asn") {
		request.IsFakeAsn = requests.NewBoolean(d.Get("is_fake_asn").(bool))
		update = true
	}
	if d.HasChange("auth_key") {
		request.AuthKey = d.Get("auth_key").(string)
		update = true
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyBgpGroupAttribute(request)
		})
		if err != nil {
			WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAliyunBgpGroupRead(d, meta)
}

func resourceAliyunBgpGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	BgpGroupService := BgpService{client}

	request := vpc.CreateDeleteBgpGroupRequest()
	request.RegionId = client.RegionId
	request.BgpGroupId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())

	wait := incrementalWait(5*time.Second, 5*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteBgpGroup(request)
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
	return WrapError(BgpGroupService.WaitForBgpGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
