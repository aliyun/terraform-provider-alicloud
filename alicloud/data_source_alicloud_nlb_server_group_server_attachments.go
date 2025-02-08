package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudNlbServerGroupServerAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNlbServerGroupServerAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"server_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"server_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"server_ips": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"attachments": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"port": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"server_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"server_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"server_ip": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"server_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"weight": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"zone_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudNlbServerGroupServerAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("server_group_id"); ok {
		request["ServerGroupId"] = v
	}
	if v, ok := d.GetOk("server_ids"); ok {
		request["ServerIds"] = v.([]interface{})
	}
	if v, ok := d.GetOk("server_ips"); ok {
		request["ServerIps"] = v.([]interface{})
	}
	request["MaxResults"] = PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListServerGroupServers"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nlb_server_group_server_attachments", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Servers", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Servers", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ServerGroupId"], ":", item["ServerId"], ":", item["ServerType"], ":", item["Port"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":              fmt.Sprint(object["ServerGroupId"], ":", object["ServerId"], ":", object["ServerType"], ":", object["Port"]),
			"description":     object["Description"],
			"port":            formatInt(object["Port"]),
			"server_group_id": object["ServerGroupId"],
			"server_id":       object["ServerId"],
			"server_ip":       object["ServerIp"],
			"server_type":     object["ServerType"],
			"status":          object["Status"],
			"weight":          formatInt(object["Weight"]),
			"zone_id":         object["ZoneId"],
		}

		ids = append(ids, fmt.Sprint(object["ServerGroupId"], ":", object["ServerId"], ":", object["ServerType"], ":", object["Port"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("attachments", s); err != nil {
		return WrapError(err)
	}
	return nil
}
