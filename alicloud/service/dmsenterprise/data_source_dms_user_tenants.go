package dmsenterprise

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudDmsUserTenants() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDmsUserTenantsRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "IN_ACTIVE"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenants": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDmsUserTenantsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListUserTenants"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var objects []map[string]interface{}
	status := ""
	if v, ok := d.GetOk("status"); ok {
		status = fmt.Sprint(v)
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
	var response map[string]interface{}
	var err error
	response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DataDefaultErrorMsg, "alicloud_dms_user_tenants", action, tferrs.AlibabaCloudSdkGoERROR)
	}
	tflogs.AddDebug(action, response, request)

	resp, err := jsonpath.Get("$.TenantList", response)
	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, action, "$.TenantList", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["Tid"])]; !ok {
				continue
			}
		}
		if status != "" && fmt.Sprint(item["Status"]) != status {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"tid":         object["Tid"],
			"tenant_name": object["TenantName"],
			"status":      object["Status"],
			"id":          object["Tid"],
		}
		ids = append(ids, fmt.Sprint(object["Tid"]))
		names = append(names, object["TenantName"])
		s = append(s, mapping)
	}

	d.SetId(tfhelper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return tferrs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return tferrs.WrapError(err)
	}

	if err := d.Set("tenants", s); err != nil {
		return tferrs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		tfhelper.WriteToFile(output.(string), s)
	}

	return nil
}
