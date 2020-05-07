package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudRenewInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRenewInstanceRead,

		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Normal",
				ValidateFunc: validation.StringInSlice([]string{"Normal", "AutoRenewal", "NotRenewal"}, false),
			},
			// Computed values
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pricing_cycle": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRenewInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateDescribeInstanceAutoRenewAttributeRequest()
	request.RegionId = string(client.Region)
	request.InstanceType = d.Get("instance_type").(string)
	request.RenewalStatus = d.Get("renewal_status").(string)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = v.(string)
	}
	invoker := NewInvoker()
	var err error
	var ids []string
	var raw interface{}
	if err := invoker.Run(func() error {
		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeInstanceAutoRenewAttribute(request)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return err
	}); err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "DescribeInstanceAutoRenewAttribute", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response := raw.(*vpc.DescribeInstanceAutoRenewAttributeResponse)
	var s []map[string]interface{}
	if len(response.InstanceRenewAttributes.InstanceRenewAttribute) > 0 {
		for _, val := range response.InstanceRenewAttributes.InstanceRenewAttribute {
			mapping := map[string]interface{}{
				"duration":      val.Duration,
				"status":        val.RenewalStatus,
				"pricing_cycle": val.PricingCycle,
				"instance_id":   val.InstanceId,
			}
			s = append(s, mapping)
			ids = append(ids, val.InstanceId)
		}
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	return nil
}
