package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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
				ValidateFunc: validation.StringInSlice([]string{"SplitRW", "cluster", "standard"}, false),
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
				ValidateFunc: validation.StringInSlice([]string{"Memcache", "Redis"}, false),
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
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_renew_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeMap,
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
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_type": {
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
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"qps": {
							Type:     schema.TypeInt,
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
						"ssl_enable": {
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
						"secondary_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_domain": {
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
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
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
				"Value": value,
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
		request["ZoneId"] = v.(string)
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
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"architecture_type":      object["ArchitectureType"],
			"bandwidth":              object["Bandwidth"],
			"capacity":               object["Capacity"],
			"connection_mode":        object["ConnectionMode"],
			"id":                     object["InstanceId"],
			"db_instance_id":         object["InstanceId"],
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
			"max_connections":        object["Connections"],
			"connections":            object["Connections"],
			"network_type":           object["NetworkType"],
			"node_type":              object["NodeType"],
			"package_type":           object["PackageType"],
			"payment_type":           object["ChargeType"],
			"charge_type":            object["ChargeType"],
			"port":                   object["Port"],
			"private_ip":             object["PrivateIp"],
			"qps":                    object["QPS"],
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
			"connection_domain":      object["ConnectionDomain"],
		}

		configs, _ := convertJsonStringToMap(fmt.Sprint(object["Config"]))
		config := make(map[string]string, len(configs))
		for k, v := range configs {
			config[k] = fmt.Sprint(v)
		}
		mapping["config"] = config

		tags, _ := jsonpath.Get("$.Tags.Tag", object)
		if tags != nil {
			mapping["tags"] = tagsToMap(tags)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, fmt.Sprint(mapping["db_instance_name"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		rKvstoreService := RKvstoreService{client}
		resp1, err := rKvstoreService.DescribeInstanceAttribute(fmt.Sprint(object["InstanceId"]))
		if err != nil {
			return WrapError(err)
		}

		mapping["instance_release_protection"] = resp1["InstanceReleaseProtection"]
		mapping["maintain_end_time"] = resp1["MaintainEndTime"]
		mapping["maintain_start_time"] = resp1["MaintainStartTime"]
		mapping["vpc_auth_mode"] = resp1["VpcAuthMode"]
		mapping["secondary_zone_id"] = resp1["SecondaryZoneId"]

		resp2, err := rKvstoreService.DescribeInstanceAutoRenewalAttribute(fmt.Sprint(object["InstanceId"]))
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		if v, ok := resp2["AutoRenew"]; ok {
			mapping["auto_renew"] = convertStringToBool(fmt.Sprint(v))
		}

		if v, ok := resp2["Duration"]; ok {
			mapping["auto_renew_period"] = formatInt(v)
		}

		resp3, err := rKvstoreService.DescribeInstanceSSL(fmt.Sprint(object["InstanceId"]))
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		mapping["ssl_enable"] = resp3["SSLEnabled"]

		resp4, err := rKvstoreService.DescribeSecurityGroupConfiguration(fmt.Sprint(object["InstanceId"]))
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		mapping["security_group_id"] = resp4["SecurityGroupId"]

		resp5, err := rKvstoreService.DescribeSecurityIps(fmt.Sprint(object["InstanceId"]))
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		mapping["security_ip_group_attribute"] = resp5["SecurityIpGroupAttribute"]
		mapping["security_ip_group_name"] = resp5["SecurityIpGroupName"]
		mapping["security_ips"] = strings.Split(fmt.Sprint(resp5["SecurityIpList"]), ",")

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
