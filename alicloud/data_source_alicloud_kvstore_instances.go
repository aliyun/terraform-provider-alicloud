package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudKvstoreInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKvstoreInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"architecture_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "cluster", "SplitRW"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"edition_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enterprise", "Community"}, false),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"2.8", "4.0", "5.0", "6.0"}, false),
			},
			"expired": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"global_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Memcache", "Redis", "tair_essd", "tair_scm"}, false),
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CLASSIC", "VPC"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"search_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Changing", "CleaningUpExpiredData", "Creating", "Flushing", "HASwitching", "Inactive", "MajorVersionUpgrading", "Migrating", "NetworkModifying", "Normal", "Rebooting", "SSLModifying", "Transforming", "ZoneMigrating"}, false),
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_renew": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_renew_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_period": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"backup_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capacity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destroy_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_renew_change_order": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_release_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_rds": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"maintain_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connections": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"qps": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replacate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"search_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ip_group_attribute": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ip_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ssl_enable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_auth_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_cloud_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
							Removed:  "Attribute 'user_name' has been removed and using '' instead.",
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudKvstoreInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeInstances"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("architecture_type"); ok {
		request["ArchitectureType"] = v
	}
	if v, ok := d.GetOk("edition_type"); ok {
		request["EditionType"] = v
	}
	if v, ok := d.GetOk("engine_version"); ok {
		request["EngineVersion"] = v
	}
	if v, ok := d.GetOk("expired"); ok {
		request["Expired"] = v
	}
	if v, ok := d.GetOkExists("global_instance"); ok {
		request["GlobalInstance"] = v
	}
	if v, ok := d.GetOk("instance_class"); ok {
		request["InstanceClass"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("network_type"); ok {
		request["NetworkType"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("search_key"); ok {
		request["SearchKey"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["InstanceStatus"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	} else if v, ok := d.GetOk("availability_zone"); ok {
		request["ZoneId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var dBInstanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		dBInstanceNameRegex = r
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
	var response map[string]interface{}
	conn, err := client.NewRedisaClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kvstore_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Instances.KVStoreInstance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Instances.KVStoreInstance", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if dBInstanceNameRegex != nil {
				if !dBInstanceNameRegex.MatchString(fmt.Sprint(item["InstanceName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		m := make(map[string]string)
		err := json.Unmarshal([]byte(object["Config"].(string)), &m)
		if err != nil {
			return WrapError(err)
		}
		mapping := map[string]interface{}{
			"architecture_type":      object["ArchitectureType"],
			"bandwidth":              fmt.Sprint(object["Bandwidth"]),
			"capacity":               fmt.Sprint(object["Capacity"]),
			"connection_domain":      object["ConnectionDomain"],
			"connection_mode":        object["ConnectionMode"],
			"id":                     fmt.Sprint(object["InstanceId"]),
			"db_instance_id":         fmt.Sprint(object["InstanceId"]),
			"db_instance_name":       object["InstanceName"],
			"name":                   object["InstanceName"],
			"destroy_time":           object["DestroyTime"],
			"end_time":               object["EndTime"],
			"expire_time":            object["EndTime"],
			"engine_version":         object["EngineVersion"],
			"has_renew_change_order": object["HasRenewChangeOrder"],
			"instance_class":         object["InstanceClass"],
			"instance_type":          object["InstanceType"],
			"is_rds":                 object["IsRds"],
			"max_connections":        fmt.Sprint(object["Connections"]),
			"connections":            fmt.Sprint(object["Connections"]),
			"network_type":           object["NetworkType"],
			"package_type":           object["PackageType"],
			"payment_type":           object["ChargeType"],
			"charge_type":            object["ChargeType"],
			"port":                   fmt.Sprint(object["Port"]),
			"private_ip":             object["PrivateIp"],
			"qps":                    fmt.Sprint(object["QPS"]),
			"replacate_id":           object["ReplacateId"],
			"resource_group_id":      object["ResourceGroupId"],
			"search_key":             object["SearchKey"],
			"status":                 object["InstanceStatus"],
			"vswitch_id":             object["VSwitchId"],
			"vpc_cloud_instance_id":  object["VpcCloudInstanceId"],
			"vpc_id":                 object["VpcId"],
			"zone_id":                object["ZoneId"],
			"availability_zone":      object["ZoneId"],
			"region_id":              object["RegionId"],
			"create_time":            object["CreateTime"],
			"user_name":              object["UserName"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["Key"].(string)
				value := t.(map[string]interface{})["Value"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["InstanceId"]))
			names = append(names, object["InstanceName"])
			s = append(s, mapping)
			continue
		}

		r_kvstoreService := R_kvstoreService{client}
		id := fmt.Sprint(object["InstanceId"])
		getResp, err := r_kvstoreService.DescribeBackupPolicy(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["backup_period"] = strings.Split(getResp["PreferredBackupPeriod"].(string), ",")
		mapping["backup_time"] = getResp["PreferredBackupTime"]
		getResp1, err := r_kvstoreService.DescribeKvstoreInstance(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["config"] = getResp1["Config"]
		mapping["instance_release_protection"] = getResp1["InstanceReleaseProtection"]
		mapping["maintain_end_time"] = getResp1["MaintainEndTime"]
		mapping["maintain_start_time"] = getResp1["MaintainStartTime"]
		mapping["vpc_auth_mode"] = getResp1["VpcAuthMode"]

		getResp2, err := r_kvstoreService.DescribeInstanceAutoRenewalAttribute(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["auto_renew"] = getResp2["AutoRenew"]
		mapping["auto_renew_period"] = getResp2["Duration"]

		getResp3, err := r_kvstoreService.DescribeInstanceSSL(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["ssl_enable"] = getResp3["SSLEnabled"]

		getResp4, err := r_kvstoreService.DescribeSecurityGroupConfiguration(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["security_group_id"] = getResp4["SecurityGroupId"]

		getResp5, err := r_kvstoreService.DescribeSecurityIps(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["security_ip_group_attribute"] = getResp5["SecurityIpGroupAttribute"]
		mapping["security_ip_group_name"] = getResp5["SecurityIpGroupName"]
		mapping["security_ips"] = strings.Split(getResp5["SecurityIpList"].(string), ",")

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["InstanceName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
