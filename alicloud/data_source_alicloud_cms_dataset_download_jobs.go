// Package alicloud. This file is hand-written from the Cms 2024-03-30 CloudSpec DatasetDownloadJob definition.
package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsDatasetDownloadJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsDatasetDownloadJobsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"job_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataset_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"download_jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compression_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataset_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudCmsDatasetDownloadJobsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	var jobNameRegex *regexp.Regexp
	if v, ok := d.GetOk("job_name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		jobNameRegex = r
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

	workspace := d.Get("workspace")
	datasetName := d.Get("dataset_name")
	action := fmt.Sprintf("/workspace/%s/dataset/%s/downloadjob", workspace, datasetName)
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	query["maxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, query)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.downloadJobs[*]", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			itemId := fmt.Sprintf("%v:%v:%v", item["workspace"], item["datasetName"], item["jobName"])
			if jobNameRegex != nil {
				if !jobNameRegex.MatchString(fmt.Sprint(item["jobName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[itemId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0, len(objects))
	s := make([]map[string]interface{}, 0, len(objects))
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}
		mapping["compression_type"] = objectRaw["compressionType"]
		mapping["content_type"] = objectRaw["contentType"]
		mapping["create_time"] = objectRaw["createTime"]
		mapping["dataset_name"] = objectRaw["datasetName"]
		mapping["job_name"] = objectRaw["jobName"]
		mapping["query"] = objectRaw["query"]
		mapping["update_time"] = objectRaw["updateTime"]
		mapping["workspace"] = objectRaw["workspace"]
		mapping["id"] = fmt.Sprintf("%v:%v:%v", objectRaw["workspace"], objectRaw["datasetName"], objectRaw["jobName"])

		ids = append(ids, mapping["id"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("download_jobs", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}

	return nil
}
