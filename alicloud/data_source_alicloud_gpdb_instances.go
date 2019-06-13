package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudGpdbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGpdbInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_group_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudGpdbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	// name regex
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		} else {
			return WrapError(err)
		}
	}
	// availability zone
	var availabilityZone string
	if v, ok := d.GetOk("availability_zone"); ok {
		availabilityZone = strings.ToLower(v.(string))
	}
	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	request := gpdb.CreateDescribeDBInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var dbInstances []gpdb.DBInstanceAttribute
	for {
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.DescribeDBInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_gpdb_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*gpdb.DescribeDBInstancesResponse)
		addDebug(request.GetActionName(), response)
		if len(response.Items.DBInstance) < 1 {
			break
		}

		for _, item := range response.Items.DBInstance {
			// filter by description regex
			if nameRegex != nil {
				if !nameRegex.MatchString(item.DBInstanceDescription) {
					continue
				}
			}
			// filter by instance id
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.DBInstanceId]; !ok {
					continue
				}
			}
			// filter by availability zone
			if len(availabilityZone) > 0 && availabilityZone != strings.ToLower(string(item.ZoneId)) {
				continue
			}

			// describe instance
			instanceAttribute, err := gpdbService.DescribeGpdbInstance(item.DBInstanceId)
			if err != nil {
				if NotFoundError(err) {
					d.SetId("")
					return nil
				}
				return WrapError(err)
			}
			dbInstances = append(dbInstances, instanceAttribute)
		}

		if len(response.Items.DBInstance) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return describeGpdbInstances(d, dbInstances)
}

func describeGpdbInstances(d *schema.ResourceData, dbInstances []gpdb.DBInstanceAttribute) error {
	var instanceIds []string
	var instances []map[string]interface{}
	for _, item := range dbInstances {
		mapping := map[string]interface{}{
			"id":                    item.DBInstanceId,
			"description":           item.DBInstanceDescription,
			"region_id":             item.RegionId,
			"availability_zone":     item.ZoneId,
			"creation_time":         item.CreationTime,
			"status":                item.DBInstanceStatus,
			"engine":                item.Engine,
			"engine_version":        item.EngineVersion,
			"charge_type":           item.PayType,
			"instance_class":        item.DBInstanceClass,
			"instance_group_count":  item.DBInstanceGroupCount,
			"instance_network_type": item.InstanceNetworkType,
		}
		instanceIds = append(instanceIds, item.DBInstanceId)
		instances = append(instances, mapping)
	}
	d.SetId(dataResourceIdHash(instanceIds))
	if err := d.Set("instances", instances); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", instanceIds); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), instances)
	}
	return nil
}
