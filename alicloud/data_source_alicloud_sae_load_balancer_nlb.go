package alicloud

import (
	"fmt"
	"log"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudSaeLoadBalancerNlb() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSaeLoadBalancerNlbRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by_sae": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_ids": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSaeLoadBalancerNlbRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	appId := d.Get("app_id").(string)
	object, err := saeService.DescribeApplicationNlb(appId)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Data source alicloud_sae_load_balancer_nlb DescribeApplicationNlb Not Found: %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	instancesRaw, err := jsonpath.Get("$.Instances", object)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, appId, "$.Instances", object)
	}
	instances := extractInstanceList(instancesRaw)
	if len(instances) == 0 {
		return WrapError(Error("The specified app_id [%s] has no NLB load balancer binding.", appId))
	}
	instance := instances[0]
	d.SetId(fmt.Sprintf("%s", appId))
	d.Set("app_id", appId)
	d.Set("dns_name", instance["DnsName"])
	d.Set("created_by_sae", instance["CreatedBySae"])
	listenersArray := make([]interface{}, 0)
	if listenersRaw, ok := instance["Listeners"]; ok && listenersRaw != nil {
		listenersList := extractInstanceList(listenersRaw)
		for _, listener := range listenersList {
			listenersArray = append(listenersArray, map[string]interface{}{
				"port":        listener["Port"],
				"target_port": listener["TargetPort"],
				"protocol":    listener["Protocol"],
				"cert_ids":    listener["CertIds"],
			})
		}
	}
	if err := d.Set("listeners", listenersArray); err != nil {
		return WrapError(err)
	}
	return nil
}
