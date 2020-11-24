package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssAttachLb() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssAttachLbCreate,
		Read:   resourceAliyunEssAttachLbRead,
		Update: resourceAliyunEssAttachLbUpdate,
		Delete: resourceAliyunEssAttachLbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"loadbalancer_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 5,
				MinItems: 1,
			},
		},
	}
}

func resourceAliyunEssAttachLbCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("scaling_group_id").(string))
	return resourceAliyunEssAttachLbUpdate(d, meta)
}

func resourceAliyunEssAttachLbUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)
	if d.HasChange("loadbalancer_ids") {
		oldLoadbalancers, newLoadbalancers := d.GetChange("loadbalancer_ids")
		err := attachOrDetachLoadbalancers(d, client, oldLoadbalancers.(*schema.Set), newLoadbalancers.(*schema.Set))
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("loadbalancer_ids")
	}

	d.Partial(false)

	return resourceAliyunEssAttachLbRead(d, meta)
}

func resourceAliyunEssAttachLbRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssScalingGroup(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	var slbIds []string
	if len(object.LoadBalancerIds.LoadBalancerId) > 0 {
		for _, v := range object.LoadBalancerIds.LoadBalancerId {
			slbIds = append(slbIds, v)
		}
	}
	d.Set("loadbalancer_ids", slbIds)
	d.Set("scaling_group_id", object.ScalingGroupId)

	return nil
}

func resourceAliyunEssAttachLbDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	removed := convertArrayInterfaceToArrayString(d.Get("loadbalancer_ids").(*schema.Set).List())

	if len(removed) < 1 {
		return nil
	}
	object, err := essService.DescribeEssScalingGroup(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if err := essService.WaitForEssScalingGroup(object.ScalingGroupId, Active, DefaultTimeout); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	var subLists = partition(d.Get("loadbalancer_ids").(*schema.Set), int(AttachDetachLoadbalancersBatchsize))
	for _, subList := range subLists {
		detachLoadbalancersRequest := ess.CreateDetachLoadBalancersRequest()
		detachLoadbalancersRequest.RegionId = client.RegionId
		detachLoadbalancersRequest.ScalingGroupId = d.Id()
		detachLoadbalancersRequest.ForceDetach = requests.NewBoolean(true)
		detachLoadbalancersRequest.LoadBalancer = &subList
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DetachLoadBalancers(detachLoadbalancersRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), detachLoadbalancersRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(detachLoadbalancersRequest.GetActionName(), raw, detachLoadbalancersRequest.RpcRequest, detachLoadbalancersRequest)
	}
	return nil
}

