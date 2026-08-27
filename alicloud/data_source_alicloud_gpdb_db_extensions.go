// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudGpdbDbExtensions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudGpdbDbExtensionRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"extensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extension_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extension_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_install_need_restart": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_latest_version": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"latest_version": {
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudGpdbDbExtensionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListDatabaseExtensions"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	if v, ok := d.GetOk("database_name"); ok {
		request["DatabaseName"] = v
	}
	request["DBInstanceId"] = d.Get("db_instance_id")
	request["DatabaseName"] = d.Get("database_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	resp, _ := jsonpath.Get("$.Extensions[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(request["DBInstanceId"], ":", request["DatabaseName"], ":", item["ExtensionName"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(request["DBInstanceId"], ":", request["DatabaseName"], ":", objectRaw["ExtensionName"])

		mapping["description"] = objectRaw["Description"]
		mapping["status"] = objectRaw["Status"]
		mapping["extension_name"] = objectRaw["ExtensionName"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(request["DBInstanceId"], ":", request["DatabaseName"], ":", objectRaw["ExtensionName"])
		mapping, err = dataSourceAliCloudGpdbDbExtensionReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("extensions", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudGpdbDbExtensionReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	gpdbServiceV2 := GpdbServiceV2{client}
	getResp, err := gpdbServiceV2.DescribeGpdbDbExtension(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["current_version"] = objectRaw["CurrentVersion"]
	mapping["description"] = objectRaw["Description"]
	mapping["extension_id"] = objectRaw["ExtensionId"]
	mapping["is_install_need_restart"] = objectRaw["IsInstallNeedRestart"]
	mapping["is_latest_version"] = objectRaw["IsLatestVersion"]
	mapping["latest_version"] = objectRaw["LatestVersion"]
	mapping["status"] = objectRaw["Status"]
	mapping["extension_name"] = objectRaw["ExtensionName"]

	return mapping, nil
}
