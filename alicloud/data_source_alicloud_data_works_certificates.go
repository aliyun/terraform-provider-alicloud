// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudDataWorksCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudDataWorksCertificateRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_size_in_bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
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

func dataSourceAliCloudDataWorksCertificateRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "ListCertificates"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["RegionId"] = client.RegionId
	if v, ok := d.GetOkExists("project_id"); ok {
		query["ProjectId"] = strconv.Itoa(v.(int))
	}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcGet("dataworks-public", "2024-05-18", action, query, request)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.PagingInfo.Certificates[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				// The resource id is composite "<project_id>:<certificate_id>"; accept both
				// the composite form and the raw certificate id for filtering.
				rawId := fmt.Sprint(item["Id"])
				compositeId := fmt.Sprintf("%v:%v", d.Get("project_id"), item["Id"])
				if _, ok := idsMap[rawId]; !ok {
					if _, ok := idsMap[compositeId]; !ok {
						continue
					}
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["Id"]

		mapping["create_time"] = dataWorksCertificateCreateTimeFormat(objectRaw["CreateTime"])
		mapping["create_user"] = objectRaw["CreateUser"]
		mapping["description"] = objectRaw["Description"]
		mapping["file_size_in_bytes"] = objectRaw["FileSizeInBytes"]
		mapping["name"] = objectRaw["Name"]
		mapping["id"] = objectRaw["Id"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["Id"])
		mapping, err = dataSourceAliCloudDataWorksCertificateReadDescription(d, id, mapping, meta)
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

	if err := d.Set("certificates", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudDataWorksCertificateReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	dataWorksServiceV2 := DataWorksServiceV2{client}
	// DescribeDataWorksCertificate expects a composite id "<project_id>:<certificate_id>".
	compositeId := fmt.Sprintf("%v:%v", d.Get("project_id"), id)
	getResp, err := dataWorksServiceV2.DescribeDataWorksCertificate(compositeId)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["create_time"] = dataWorksCertificateCreateTimeFormat(objectRaw["CreateTime"])
	mapping["create_user"] = objectRaw["CreateUser"]
	mapping["description"] = objectRaw["Description"]
	mapping["file_size_in_bytes"] = objectRaw["FileSizeInBytes"]
	mapping["name"] = objectRaw["Name"]
	mapping["project_id"] = objectRaw["ProjectId"]
	mapping["id"] = objectRaw["Id"]

	return mapping, nil
}
