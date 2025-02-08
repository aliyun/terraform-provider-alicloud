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

func dataSourceAlicloudDbsBackupPlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDbsBackupPlansRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"backup_plan_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"check_pass", "init", "locked", "pause", "running", "stop", "wait"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"plans": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_objects": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_plan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_plan_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_retention_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backup_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_aliyun_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duplication_archive_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"duplication_infrequent_access_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enable_backup_log": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_database_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_sid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func dataSourceAlicloudDbsBackupPlansRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeBackupPlanList"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("backup_plan_name"); ok {
		request["BackupPlanName"] = v
	}
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNum"] = 0
	var objects []map[string]interface{}
	var backupPlanNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		backupPlanNameRegex = r
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Dbs", "2019-03-06", action, nil, request, true)
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
			if IsExpectedErrors(err, []string{"InvalidParameterValid"}) {
				d.SetId("DescribeBackupPlanList")
				d.Set("ids", []string{})
				return nil
			}
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dbs_backup_plans", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Items.BackupPlanDetail", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.BackupPlanDetail", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if backupPlanNameRegex != nil && !backupPlanNameRegex.MatchString(fmt.Sprint(item["BackupPlanName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["BackupPlanId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["BackupPlanStatus"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"backup_method":                        object["BackupMethod"],
			"backup_objects":                       object["BackupObjects"],
			"backup_period":                        convertDbsBackupPeriodResponse(object["BackupPeriod"]),
			"id":                                   fmt.Sprint(object["BackupPlanId"]),
			"backup_plan_id":                       fmt.Sprint(object["BackupPlanId"]),
			"backup_plan_name":                     object["BackupPlanName"],
			"backup_retention_period":              formatInt(object["BackupRetentionPeriod"]),
			"backup_start_time":                    object["BackupStartTime"],
			"backup_storage_type":                  object["BackupStorageType"],
			"cross_aliyun_id":                      object["CrossAliyunId"],
			"cross_role_name":                      object["CrossRoleName"],
			"database_type":                        object["DatabaseType"],
			"duplication_archive_period":           formatInt(object["DuplicationArchivePeriod"]),
			"duplication_infrequent_access_period": formatInt(object["DuplicationInfrequentAccessPeriod"]),
			"enable_backup_log":                    object["EnableBackupLog"],
			"instance_class":                       object["InstanceClass"],
			"oss_bucket_name":                      object["OSSBucketName"],
			"source_endpoint_database_name":        object["SourceEndpointDatabaseName"],
			"source_endpoint_instance_id":          object["SourceEndpointInstanceID"],
			"source_endpoint_instance_type":        convertDbsSourceEndpointInstanceTypeResponse(object["SourceEndpointInstanceType"]),
			"source_endpoint_region":               object["SourceEndpointRegion"],
			"source_endpoint_sid":                  object["SourceEndpointOracleSID"],
			"source_endpoint_user_name":            object["SourceEndpointUserName"],
			"status":                               object["BackupPlanStatus"],
			"resource_group_id":                    object["ResourceGroupId"],
		}
		if v, ok := object["BackupGatewayId"]; ok && fmt.Sprint(object["BackupGatewayId"]) != "0" {
			mapping["backup_gateway_id"] = fmt.Sprint(v)
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["BackupPlanName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["BackupPlanId"])
		dbsService := DbsService{client}
		getResp, err := dbsService.DescribeBackupPlanBilling(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["payment_type"] = convertDbsBackupPlanPaymentTypeResponse(getResp["BuyChargeType"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("plans", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
