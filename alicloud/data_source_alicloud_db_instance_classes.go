package alicloud

import (
	"encoding/json"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudDBInstanceClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDBInstanceClassesRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validateAllowedStringValue([]string{string(PostPaid), string(PrePaid)}),
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"cloud_ssd", "local_ssd"}),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values.
			"instance_classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_range": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"max": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"step": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDBInstanceClassesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := rds.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	request.ZoneId = d.Get("zone_id").(string)
	instanceChargeType := d.Get("instance_charge_type").(string)
	if instanceChargeType == string(PostPaid) {
		instanceChargeType = string(Postpaid)
	} else {
		instanceChargeType = string(Prepaid)
	}
	request.InstanceChargeType = instanceChargeType

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeAvailableResource(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_classes", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	response := raw.(*rds.DescribeAvailableResourceResponse)

	infos := []map[string]interface{}{}
	ids := []string{}

	engine, engineGot := d.GetOk("engine")
	engineVersion, engineVersionGot := d.GetOk("engine_version")
	storageType, storageTypeGot := d.GetOk("storage_type")
	category, categoryGot := d.GetOk("category")

	for _, AvailableZone := range response.AvailableZones.AvailableZone {
		info := make(map[string]interface{})
		info["zone_id"] = AvailableZone.ZoneId
		ids = append(ids, AvailableZone.ZoneId)
		for _, SupportedEngine := range AvailableZone.SupportedEngine.SupportedEngines {
			if engineGot && engine.(string) != SupportedEngine.Engine {
				continue
			}
			ids = append(ids, SupportedEngine.Engine)
			for _, SupportedEngineVersion := range SupportedEngine.SupportedEngineVersion.SupportedEngineVersions {
				if engineVersionGot && engineVersion != SupportedEngineVersion.Version {
					continue
				}
				ids = append(ids, SupportedEngineVersion.Version)
				for _, SupportedCategory := range SupportedEngineVersion.SupportedCategory.SupportedCategorys {
					if categoryGot && category.(string) != SupportedCategory.Category {
						continue
					}
					for _, SupportedStorageType := range SupportedCategory.SupportedStorageType.SupportedStorageTypes {
						if storageTypeGot && storageType.(string) != SupportedStorageType.StorageType {
							continue
						}
						for _, AvailableResource := range SupportedStorageType.AvailableResource.AvailableResources {
							storageRange := make(map[string]interface{})
							err = json.Unmarshal([]byte(AvailableResource.StorageRange), &storageRange)
							if err != nil {
								return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_classes", request.GetActionName(), AlibabaCloudSdkGoERROR)
							}
							//json number to string
							for _, key := range []string{"min", "max", "step"} {
								if v, ok := storageRange[key]; ok {
									storageRange[key] = strconv.Itoa(int(v.(float64)))
								} else {
									return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_classes", request.GetActionName(), AlibabaCloudSdkGoERROR)
								}
							}
							info["storage_range"] = storageRange
							info["instance_class"] = AvailableResource.DBInstanceClass
							temp := make(map[string]interface{}, len(info))
							for key, value := range info {
								temp[key] = value
							}
							infos = append(infos, temp)
						}
					}
				}
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))
	err = d.Set("instance_classes", infos)
	if err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), infos)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}
