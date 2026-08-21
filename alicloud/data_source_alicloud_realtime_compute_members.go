package alicloud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRealtimeComputeMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRealtimeComputeMembersRead,
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"total_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAlicloudRealtimeComputeMembersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	realtimeComputeServiceV2 := RealtimeComputeServiceV2{client}
	resourceId := d.Get("resource_id").(string)
	namespace := d.Get("namespace").(string)

	objects, resp, err := realtimeComputeServiceV2.ListRealtimeComputeFlinkMembers(resourceId, namespace)
	if err != nil {
		return WrapError(err)
	}

	members := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	for _, object := range objects {
		memberMap := object.(map[string]interface{})
		memberName, _ := memberMap["member"].(string)
		role, _ := memberMap["role"].(string)
		members = append(members, map[string]interface{}{
			"member": memberName,
			"role":   role,
		})
		ids = append(ids, memberName)
	}

	d.SetId(fmt.Sprintf("%s:%s", resourceId, namespace))
	d.Set("resource_id", resourceId)
	d.Set("namespace", namespace)
	d.Set("ids", ids)
	if resp != nil {
		if totalSize, ok := resp["totalSize"]; ok && totalSize != nil {
			switch t := totalSize.(type) {
			case json.Number:
				if n, e := t.Int64(); e == nil {
					d.Set("total_size", int(n))
				}
			case float64:
				d.Set("total_size", int(t))
			}
		}
	}
	if err := d.Set("members", members); err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), members); err != nil {
			log.Printf("[WARN] write output file failed: %s", err)
		}
	}

	return nil
}
