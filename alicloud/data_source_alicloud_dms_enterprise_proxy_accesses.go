package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDmsEnterpriseProxyAccesses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDmsEnterpriseProxyAccessesRead,
		Schema: map[string]*schema.Schema{
			"proxy_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_details": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"accesses": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"access_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"access_secret": {
							Computed:  true,
							Sensitive: true,
							Type:      schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"indep_account": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"instance_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"origin_info": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"proxy_access_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"proxy_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"user_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"user_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"user_uid": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDmsEnterpriseProxyAccessesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("proxy_id"); ok {
		request["ProxyId"] = v
	}

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
	action := "ListProxyAccesses"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_proxy_accesses", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.ProxyAccessList", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ProxyAccessList", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["ProxyAccessId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	dmsEnterpriseService := DmsEnterpriseService{client}
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":              fmt.Sprint(object["ProxyAccessId"]),
			"access_id":       object["AccessId"],
			"create_time":     object["GmtCreate"],
			"indep_account":   object["IndepAccount"],
			"instance_id":     object["InstanceId"],
			"origin_info":     object["OriginInfo"],
			"proxy_access_id": object["ProxyAccessId"],
			"proxy_id":        object["ProxyId"],
			"user_id":         object["UserId"],
			"user_name":       object["UserName"],
			"user_uid":        object["UserUid"],
		}

		ids = append(ids, fmt.Sprint(object["ProxyAccessId"]))

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["ProxyAccessId"])
		object, err = dmsEnterpriseService.DescribeDmsEnterpriseProxyAccess(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["access_id"] = object["AccessId"]
		mapping["create_time"] = object["GmtCreate"]
		mapping["indep_account"] = object["IndepAccount"]
		mapping["instance_id"] = object["InstanceId"]
		mapping["origin_info"] = object["OriginInfo"]
		mapping["proxy_access_id"] = object["ProxyAccessId"]
		mapping["proxy_id"] = object["ProxyId"]
		mapping["user_id"] = object["UserId"]
		mapping["user_name"] = object["UserName"]
		mapping["user_uid"] = object["UserUid"]
		object, err = dmsEnterpriseService.InspectProxyAccessSecret(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["access_secret"] = object["AccessSecret"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("accesses", s); err != nil {
		return WrapError(err)
	}
	return nil
}
