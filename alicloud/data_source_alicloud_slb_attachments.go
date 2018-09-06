package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/schema"
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
	conn := meta.(*AliyunClient).slbconn

	args := slb.CreateDescribeLoadBalancerAttributeRequest()
	args.LoadBalancerId = d.Get("load_balancer_id").(string)

	instanceIdsMap := make(map[string]string)
	if v, ok := d.GetOk("instance_ids"); ok {
		for _, vv := range v.([]interface{}) {
			instanceIdsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	resp, err := conn.DescribeLoadBalancerAttribute(args)
	if err != nil {
		return fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", err)
	}
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

	if len(filteredBackendServersTemp) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_slb_attachments - Slb attachments found: %#v", filteredBackendServersTemp)

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

		log.Printf("[DEBUG] alicloud_slb_attachments - adding slb_attachment mapping: %v", mapping)
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
