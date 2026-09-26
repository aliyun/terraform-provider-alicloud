package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudThreatDetectionDataConnectors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionDataConnectorsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_connectors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_connector_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_connector_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_connector_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_connector_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_connector_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_data_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_data_source_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_store_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_config_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_config_vendor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_config_product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_ingestion_job_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_ingestion_job_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionDataConnectorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("region_id"); ok {
		request["RegionId"] = v
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

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	request["PageNumber"] = 1
	request["PageSize"] = PageSizeMedium

	var objects []interface{}
	var response map[string]interface{}
	for {
		action := "ListDataConnectors"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_data_connectors", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DataConnector", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DataConnector", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DataConnectorId"])]; !ok {
					continue
				}
			}
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["DataConnectorName"])) {
				continue
			}
			objects = append(objects, item)
		}
		if isPagingRequest(d) || len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	cloudSiemService := CloudSiemService{client}
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                      fmt.Sprint(object["DataConnectorId"]),
			"data_connector_id":       fmt.Sprint(object["DataConnectorId"]),
			"data_connector_name":     object["DataConnectorName"],
			"data_connector_type":     object["DataConnectorType"],
			"data_connector_status":   object["DataConnectorStatus"],
			"data_connector_config":   object["DataConnectorConfig"],
			"src_data_type":           object["SrcDataType"],
			"dest_data_source_id":     object["DestDataSourceId"],
			"log_project_name":        object["LogProjectName"],
			"log_store_name":          object["LogStoreName"],
			"log_region_id":           object["LogRegionId"],
			"auth_config_id":          object["AuthConfigId"],
			"auth_config_vendor":      object["AuthConfigVendor"],
			"auth_config_product":     object["AuthConfigProduct"],
			"sls_ingestion_job_name":  object["SlsIngestionJobName"],
			"sls_ingestion_job_state": object["SlsIngestionJobState"],
			"creation_time":           object["CreationTime"],
			"update_time":             object["UpdateTime"],
		}
		ids = append(ids, fmt.Sprint(object["DataConnectorId"]))

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DataConnectorId"])
		obj, err := cloudSiemService.DescribeThreatDetectionDataConnector(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["data_connector_config"] = obj["DataConnectorConfig"]
		mapping["data_connector_status"] = obj["DataConnectorStatus"]
		mapping["auth_config_id"] = obj["AuthConfigId"]
		mapping["auth_config_vendor"] = obj["AuthConfigVendor"]
		mapping["auth_config_product"] = obj["AuthConfigProduct"]
		mapping["sls_ingestion_job_name"] = obj["SlsIngestionJobName"]
		mapping["sls_ingestion_job_state"] = obj["SlsIngestionJobState"]
		mapping["creation_time"] = obj["CreationTime"]
		mapping["update_time"] = obj["UpdateTime"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("data_connectors", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
