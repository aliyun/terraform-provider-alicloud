// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudSlsEtls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSlsEtlRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"etls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"script": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"to_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"parameters": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"sink": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"datasets": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"project": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"endpoint": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"logstore": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"role_arn": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"logstore": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"lang": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"from_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"role_arn": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"schedule_id": {
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
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudSlsEtlRead(d *schema.ResourceData, meta interface{}) error {
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
	var query map[string]*string
	action := fmt.Sprintf("/etls")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(d.Get("project").(string))
	if v, ok := d.GetOkExists("logstore"); ok {
		query["logstore"] = StringPointer(fmt.Sprint(v))
	}

	if v, ok := d.GetOkExists("offset"); ok {
		query["offset"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOkExists("max_results"); ok {
		query["size"] = StringPointer(strconv.Itoa(v.(int)))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["MaxResults"] = StringPointer("50")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("GET", "2020-12-30", "ListETLs", action), query, nil, nil, hostMap, false)

			if err != nil {
				if IsExpectedErrors(err, []string{"403"}) || NeedRetry(err) {
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

		resp, _ := jsonpath.Get("$.results[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(*hostMap["project"], ":", item["name"])]; !ok {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(*hostMap["project"], ":", objectRaw["name"])

		mapping["create_time"] = objectRaw["createTime"]
		mapping["description"] = objectRaw["displayName"]
		mapping["display_name"] = objectRaw["description"]
		mapping["last_modified_time"] = objectRaw["lastModifiedTime"]
		mapping["schedule_id"] = objectRaw["scheduleId"]
		mapping["status"] = objectRaw["status"]
		mapping["job_name"] = objectRaw["name"]

		configurationMaps := make([]map[string]interface{}, 0)
		configurationMap := make(map[string]interface{})
		configurationRaw := make(map[string]interface{})
		if objectRaw["configuration"] != nil {
			configurationRaw = objectRaw["configuration"].(map[string]interface{})
		}
		if len(configurationRaw) > 0 {
			configurationMap["from_time"] = configurationRaw["fromTime"]
			configurationMap["lang"] = configurationRaw["lang"]
			configurationMap["logstore"] = configurationRaw["logstore"]
			configurationMap["parameters"] = configurationRaw["parameters"]
			configurationMap["role_arn"] = configurationRaw["roleArn"]
			configurationMap["script"] = configurationRaw["script"]
			configurationMap["to_time"] = configurationRaw["toTime"]

			sinksRaw := configurationRaw["sinks"]
			sinkMaps := make([]map[string]interface{}, 0)
			if sinksRaw != nil {
				for _, sinksChildRaw := range sinksRaw.([]interface{}) {
					sinkMap := make(map[string]interface{})
					sinksChildRaw := sinksChildRaw.(map[string]interface{})
					sinkMap["endpoint"] = sinksChildRaw["endpoint"]
					sinkMap["logstore"] = sinksChildRaw["logstore"]
					sinkMap["name"] = sinksChildRaw["name"]
					sinkMap["project"] = sinksChildRaw["project"]
					sinkMap["role_arn"] = sinksChildRaw["roleArn"]

					datasetsRaw := make([]interface{}, 0)
					if sinksChildRaw["datasets"] != nil {
						datasetsRaw = sinksChildRaw["datasets"].([]interface{})
					}

					sinkMap["datasets"] = datasetsRaw
					sinkMaps = append(sinkMaps, sinkMap)
				}
			}
			configurationMap["sink"] = sinkMaps
			configurationMaps = append(configurationMaps, configurationMap)
		}
		mapping["configuration"] = configurationMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw[""])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("etls", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
