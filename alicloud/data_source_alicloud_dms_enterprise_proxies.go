package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDmsEnterpriseProxies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDmsEnterpriseProxiesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"tid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"https_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"private_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"public_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDmsEnterpriseProxiesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListProxies"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("tid"); ok {
		request["Tid"] = v
	}
	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_proxies", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.ProxyList", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ProxyList", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["ProxyId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"creator_id":     fmt.Sprint(object["CreatorId"]),
			"creator_name":   object["CreatorName"],
			"https_port":     formatInt(object["HttpsPort"]),
			"instance_id":    fmt.Sprint(object["InstanceId"]),
			"private_enable": object["PrivateEnable"],
			"private_host":   object["PrivateHost"],
			"protocol_port":  formatInt(object["ProtocolPort"]),
			"protocol_type":  object["ProtocolType"],
			"id":             fmt.Sprint(object["ProxyId"]),
			"proxy_id":       fmt.Sprint(object["ProxyId"]),
			"public_enable":  object["PublicEnable"],
			"public_host":    object["PublicHost"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("proxies", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
