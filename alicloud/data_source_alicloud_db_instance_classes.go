package alicloud

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

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
			"db_instance_class": {
				Type:     schema.TypeString,
				Optional: true,
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
			"multi_zone": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
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
						"zone_ids": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sub_zone_ids": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Computed: true,
									},
								},
							},
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
	multiZone := d.Get("multi_zone").(bool)
	if instanceChargeType == string(PostPaid) {
		instanceChargeType = string(Postpaid)
	} else {
		instanceChargeType = string(Prepaid)
	}
	request.InstanceChargeType = instanceChargeType
	var response = &rds.DescribeAvailableResourceResponse{}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeAvailableResource(request)
		})
		if err != nil {
			if IsExceptedError(err, Throttling) {
				time.Sleep(time.Duration(5) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response = raw.(*rds.DescribeAvailableResourceResponse)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_classes", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	type ClassInfosItem struct {
		Index        int
		StorageRange map[string]string
		ZoneIds      []map[string]interface{}
	}

	classInfos := make(map[string]ClassInfosItem)
	indexMap := make(map[string]int)
	ids := []string{}

	engine, engineGot := d.GetOk("engine")
	engineVersion, engineVersionGot := d.GetOk("engine_version")
	dbInstanceClass, dbInstanceClassGot := d.GetOk("db_instance_class")
	storageType, storageTypeGot := d.GetOk("storage_type")
	category, categoryGot := d.GetOk("category")

	for _, AvailableZone := range response.AvailableZones.AvailableZone {
		id_item := []string{}
		if multiZone {
			if !strings.Contains(AvailableZone.ZoneId, "MAZ") {
				continue
			}
			for _, v := range splitMultiZoneId(AvailableZone.ZoneId) {
				id_item = append(id_item, v)
			}
		} else {
			if strings.Contains(AvailableZone.ZoneId, "MAZ") {
				continue
			}
			id_item = []string{AvailableZone.ZoneId}
		}

		zoneId := map[string]interface{}{
			"id":           AvailableZone.ZoneId,
			"sub_zone_ids": id_item,
		}

		ids = append(ids, AvailableZone.ZoneId)
		for _, SupportedEngine := range AvailableZone.SupportedEngines.SupportedEngine {
			if engineGot && engine.(string) != SupportedEngine.Engine {
				continue
			}
			ids = append(ids, SupportedEngine.Engine)
			for _, SupportedEngineVersion := range SupportedEngine.SupportedEngineVersions.SupportedEngineVersion {
				if engineVersionGot && engineVersion != SupportedEngineVersion.Version {
					continue
				}
				ids = append(ids, SupportedEngineVersion.Version)
				for _, SupportedCategory := range SupportedEngineVersion.SupportedCategorys.SupportedCategory {
					if categoryGot && category.(string) != SupportedCategory.Category {
						continue
					}
					for _, SupportedStorageType := range SupportedCategory.SupportedStorageTypes.SupportedStorageType {
						if storageTypeGot && storageType.(string) != SupportedStorageType.StorageType {
							continue
						}
						for _, AvailableResource := range SupportedStorageType.AvailableResources.AvailableResource {
							if dbInstanceClassGot && dbInstanceClass.(string) != AvailableResource.DBInstanceClass {
								continue
							}
							zoneIds := []map[string]interface{}{}
							if _, ok := classInfos[AvailableResource.DBInstanceClass]; ok {
								zoneIds = append(classInfos[AvailableResource.DBInstanceClass].ZoneIds, zoneId)
							} else {
								zoneIds = []map[string]interface{}{zoneId}
								indexMap[AvailableResource.DBInstanceClass] = len(classInfos)
							}
							classInfos[AvailableResource.DBInstanceClass] = ClassInfosItem{
								Index: indexMap[AvailableResource.DBInstanceClass],
								StorageRange: map[string]string{
									"min":  strconv.Itoa(AvailableResource.DBInstanceStorageRange.Min),
									"max":  strconv.Itoa(AvailableResource.DBInstanceStorageRange.Max),
									"step": strconv.Itoa(AvailableResource.DBInstanceStorageRange.Step),
								},
								ZoneIds: zoneIds,
							}
						}
					}
				}
			}
		}
	}

	infos := make([]map[string]interface{}, len(classInfos))
	for k, v := range classInfos {
		infos[v.Index] = map[string]interface{}{
			"zone_ids":       v.ZoneIds,
			"storage_range":  v.StorageRange,
			"instance_class": k,
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
