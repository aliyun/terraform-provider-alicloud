package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudSlbAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"slb_attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := slb.CreateDescribeLoadBalancerAttributeRequest()
	args.LoadBalancerId = d.Get("load_balancer_id").(string)

	instanceIdsMap := make(map[string]string)
	if v, ok := d.GetOk("instance_ids"); ok {
		for _, vv := range v.([]interface{}) {
			instanceIdsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(args)
	})
	if err != nil {
		return fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", err)
	}
	resp, _ := raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if resp == nil {
		return fmt.Errorf("there is no SLB with the ID %s. Please change your search criteria and try again", args.LoadBalancerId)
	}

	var filteredBackendServersTemp []slb.BackendServer
	if len(instanceIdsMap) > 0 {
		for _, backendServer := range resp.BackendServers.BackendServer {
			if len(instanceIdsMap) > 0 {
				if _, ok := instanceIdsMap[backendServer.ServerId]; !ok {
					continue
				}
			}

			filteredBackendServersTemp = append(filteredBackendServersTemp, backendServer)
		}
	} else {
		filteredBackendServersTemp = resp.BackendServers.BackendServer
	}

	return slbAttachmentsDescriptionAttributes(d, filteredBackendServersTemp)
}

func slbAttachmentsDescriptionAttributes(d *schema.ResourceData, backendServers []slb.BackendServer) error {
	var ids []string
	var s []map[string]interface{}

	for _, backendServer := range backendServers {
		mapping := map[string]interface{}{
			"instance_id": backendServer.ServerId,
			"weight":      backendServer.Weight,
		}

		ids = append(ids, backendServer.ServerId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slb_attachments", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
