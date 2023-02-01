package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
	"time"
)

func dataSourceAlicloudHbrHanaRestoreJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrHanaRestoreJobsRead,
		Schema: map[string]*schema.Schema{
			"backup_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"cluster_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"database_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"restore_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"restore_status": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"token": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"vault_id": {
				Optional: true,
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
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  10,
			},
			"jobs": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"backup_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"backup_prefix": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"check_access": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"clear_log": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"cluster_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"current_phase": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"current_progress": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"database_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"database_restore_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"end_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"log_position": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"max_phase": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"max_progress": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"message": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"mode": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"phase": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"reached_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"recovery_point_in_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"restore_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source_cluster_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"start_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"state": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"system_copy": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"use_catalog": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"use_delta": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"vault_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"volume_id": {
							Computed: true,
							Type:     schema.TypeInt,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHbrHanaRestoreJobsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("backup_id"); ok {
		request["BackupId"] = v
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}
	if v, ok := d.GetOk("database_name"); ok {
		request["DatabaseName"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("restore_id"); ok {
		request["RestoreId"] = v
	}
	if v, ok := d.GetOk("restore_status"); ok {
		request["RestoreStatus"] = v
	}
	if v, ok := d.GetOk("token"); ok {
		request["Token"] = v
	}
	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
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

	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeHanaRestores"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_hana_restore_jobs", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.HanaRestore.HanaRestores", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.HanaRestore.HanaRestores", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RestoreId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                     fmt.Sprint(object["RestoreId"]),
			"backup_id":              object["BackupID"],
			"backup_prefix":          object["BackupPrefix"],
			"check_access":           object["CheckAccess"],
			"clear_log":              object["ClearLog"],
			"cluster_id":             object["ClusterId"],
			"current_phase":          object["CurrentPhase"],
			"current_progress":       object["CurrentProgress"],
			"database_name":          object["DatabaseName"],
			"database_restore_id":    object["DatabaseRestoreId"],
			"end_time":               object["EndTime"],
			"log_position":           object["LogPosition"],
			"max_phase":              object["MaxPhase"],
			"max_progress":           object["MaxProgress"],
			"message":                object["Message"],
			"mode":                   object["Mode"],
			"phase":                  object["Phase"],
			"reached_time":           object["ReachedTime"],
			"recovery_point_in_time": object["RecoveryPointInTime"],
			"restore_id":             object["RestoreId"],
			"source":                 object["Source"],
			"source_cluster_id":      object["SourceClusterId"],
			"start_time":             object["StartTime"],
			"state":                  object["State"],
			"status":                 object["Status"],
			"system_copy":            object["SystemCopy"],
			"use_catalog":            object["UseCatalog"],
			"use_delta":              object["UseDelta"],
			"vault_id":               object["VaultId"],
			"volume_id":              object["VolumeId"],
		}

		ids = append(ids, fmt.Sprint(object["RestoreId"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("jobs", s); err != nil {
		return WrapError(err)
	}
	return nil
}
