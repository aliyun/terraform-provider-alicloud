package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksConnectionsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"connection_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"env_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connection_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"env_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"content": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connect_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"binding_calc_engine_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gmt_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sequence": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"shared": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"default_engine": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ListConnections"
	request := map[string]interface{}{
		"ProjectId":  d.Get("project_id"),
		"PageNumber": 1,
		"PageSize":   100,
	}
	if v, ok := d.GetOk("connection_type"); ok && v.(string) != "" {
		request["ConnectionType"] = v
	}
	if v, ok := d.GetOk("sub_type"); ok && v.(string) != "" {
		request["SubType"] = v
	}
	if v, ok := d.GetOk("status"); ok && v.(string) != "" {
		request["Status"] = v
	}
	if v, ok := d.GetOk("env_type"); ok {
		request["EnvType"] = v
	}
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request["Name"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_connections", action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Data.Connections", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_data_works_connections", "$.Data.Connections", response)
	}
	connections, ok := v.([]interface{})
	if !ok {
		return WrapError(fmt.Errorf("Connections is not an array"))
	}
	var rawConnections []map[string]interface{}
	var ids []string
	for _, conn := range connections {
		c, ok := conn.(map[string]interface{})
		if !ok {
			continue
		}
		connectionId := fmt.Sprintf("%v", c["Id"])
		mapping := map[string]interface{}{
			"id":                     fmt.Sprintf("%d:%s", d.Get("project_id").(int), connectionId),
			"connection_id":          c["Id"],
			"project_id":             c["ProjectId"],
			"connection_name":        c["Name"],
			"connection_type":        c["ConnectionType"],
			"sub_type":               c["SubType"],
			"env_type":               c["EnvType"],
			"content":                c["Content"],
			"description":            c["Description"],
			"operator":               c["Operator"],
			"connect_status":         c["ConnectStatus"],
			"binding_calc_engine_id": c["BindingCalcEngineId"],
			"gmt_modified":           c["GmtModified"],
			"sequence":               c["Sequence"],
			"shared":                 c["Shared"],
			"default_engine":         c["DefaultEngine"],
			"create_time":            c["GmtCreate"],
			"tenant_id":              c["TenantId"],
			"region_id":              client.RegionId,
		}
		if c["Status"] != nil {
			mapping["status"] = fmt.Sprintf("%v", c["Status"])
		}
		rawConnections = append(rawConnections, mapping)
		ids = append(ids, mapping["id"].(string))
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("connections", rawConnections); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), rawConnections)
	}
	return nil
}
