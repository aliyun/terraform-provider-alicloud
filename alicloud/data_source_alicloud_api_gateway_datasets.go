package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudApiGatewayDatasets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApiGatewayDatasetsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datasets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataset_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataset_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataset_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudApiGatewayDatasetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeDatasetList"
	request := make(map[string]interface{})
	query := make(map[string]interface{})

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var objects []map[string]interface{}
	pageNumber := 1
	pageSize := PageSizeLarge
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	for {
		query["PageNumber"] = pageNumber
		query["PageSize"] = pageSize
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, "", action, AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$.DatasetInfoList[*]", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, "", "$.DatasetInfoList[*]", response)
		}
		if v == nil {
			break
		}
		for _, datasetRaw := range v.([]interface{}) {
			dataset := datasetRaw.(map[string]interface{})
			if len(idsMap) > 0 {
				id := fmt.Sprint(dataset["DatasetId"])
				if _, ok := idsMap[id]; !ok {
					continue
				}
			}
			objects = append(objects, dataset)
		}

		totalCount, _ := jsonpath.Get("$.TotalCount", response)
		count := 0
		switch v := totalCount.(type) {
		case float64:
			count = int(v)
		case json.Number:
			n, _ := v.Int64()
			count = int(n)
		}
		if pageNumber*pageSize >= count {
			break
		}
		if len(idsMap) > 0 && len(objects) >= len(idsMap) {
			break
		}
		pageNumber++
	}

	return describeApiGatewayDatasetsFilters(d, objects)
}

func describeApiGatewayDatasetsFilters(d *schema.ResourceData, objects []map[string]interface{}) error {
	var ids []string
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}
	var s []map[string]interface{}
	for _, object := range objects {
		name := fmt.Sprint(object["DatasetName"])
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		datasetId := fmt.Sprint(object["DatasetId"])
		mapping := map[string]interface{}{
			"id":            datasetId,
			"dataset_id":    datasetId,
			"dataset_name":  object["DatasetName"],
			"dataset_type":  object["DatasetType"],
			"description":   object["Description"],
			"create_time":   object["CreatedTime"],
			"modified_time": object["ModifiedTime"],
		}
		ids = append(ids, datasetId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	d.Set("ids", ids)
	if err := d.Set("datasets", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
