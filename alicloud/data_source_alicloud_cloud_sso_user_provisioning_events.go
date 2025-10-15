// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCloudSsoUserProvisioningEvents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudSsoUserProvisioningEventRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_provisioning_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"events": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deletion_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duplication_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_count": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"error_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_sync_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_provisioning_id": {
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

func dataSourceAliCloudCloudSsoUserProvisioningEventRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "ListUserProvisioningEvents"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	if v, ok := d.GetOk("directory_id"); ok {
		request["DirectoryId"] = v
	}
	request["DirectoryId"] = d.Get("directory_id")
	if v, ok := d.GetOk("user_provisioning_id"); ok {
		request["UserProvisioningId"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["MaxResults"] = PageSizeLarge
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.UserProvisioningEvents[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DirectoryId"], ":", item["EventId"])]; !ok {
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

		mapping["id"] = fmt.Sprint(objectRaw["DirectoryId"], ":", objectRaw["EventId"])

		mapping["content"] = objectRaw["Content"]
		mapping["create_time"] = objectRaw["GmtCreate"]
		mapping["deletion_strategy"] = objectRaw["DeletionStrategy"]
		mapping["duplication_strategy"] = objectRaw["DuplicationStrategy"]
		mapping["error_count"] = objectRaw["ErrorCount"]
		mapping["error_info"] = objectRaw["ErrorInfo"]
		mapping["last_sync_time"] = objectRaw["GmtLastExec"]
		mapping["principal_id"] = objectRaw["PrincipalId"]
		mapping["principal_name"] = objectRaw["PrincipalName"]
		mapping["principal_type"] = objectRaw["PrincipalType"]
		mapping["source_type"] = objectRaw["SourceType"]
		mapping["target_id"] = objectRaw["TargetId"]
		mapping["target_name"] = objectRaw["TargetName"]
		mapping["target_path"] = objectRaw["TargetPath"]
		mapping["target_type"] = objectRaw["TargetType"]
		mapping["update_time"] = objectRaw["GmtModified"]
		mapping["user_provisioning_id"] = objectRaw["UserProvisioningId"]
		mapping["directory_id"] = objectRaw["DirectoryId"]
		mapping["event_id"] = objectRaw["EventId"]

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw[""])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("events", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
