package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPolarDBNodeClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPolarDBInstanceClassesRead,

		Schema: map[string]*schema.Schema{
			"pay_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
			},
			"db_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "Basic", "ArchiveNormal", "NormalMultimaster", "SENormal"}, false),
			},
			// Computed values.
			"classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supported_engines": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"available_resources": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_node_class": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"engine": {
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

func dataSourceAlicloudPolarDBInstanceClassesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := polardb.CreateDescribeDBClusterAvailableResourcesRequest()
	payType := d.Get("pay_type").(string)
	if payType == string(PostPaid) {
		request.PayType = "Postpaid"
	} else if payType == string(PrePaid) {
		request.PayType = "Prepaid"
	}

	dbType, dbTypeGot := d.GetOk("db_type")
	dbVersion, dbVersionGot := d.GetOk("db_version")
	if dbTypeGot && dbVersionGot {
		request.DBType = dbType.(string)
		request.DBVersion = dbVersion.(string)
	}
	checkDBNodeClass := ""
	if dbNodeClass, ok := d.GetOk("db_node_class"); ok {
		request.DBNodeClass = dbNodeClass.(string)
		checkDBNodeClass = dbNodeClass.(string)
	}
	if regionId, ok := d.GetOk("region_id"); ok {
		request.RegionId = regionId.(string)
	}
	if zoneId, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = zoneId.(string)
	}

	var category string
	if s, ok := d.GetOk("category"); ok && s.(string) != "" {
		category = s.(string)
	}

	var response = &polardb.DescribeDBClusterAvailableResourcesResponse{}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeDBClusterAvailableResources(request)
		})
		if err != nil {
			if NeedRetry(err) {
				time.Sleep(time.Duration(5) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*polardb.DescribeDBClusterAvailableResourcesResponse)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_node_classes", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	ids := []string{}
	var availableClasses []interface{}
	for _, AvailableZone := range response.AvailableZones {
		zondId := AvailableZone.ZoneId
		ids = append(ids, zondId)

		supportedEngines := make([]interface{}, 0)
		for _, supportedEngine := range AvailableZone.SupportedEngines {
			if len(supportedEngine.AvailableResources) == 0 {
				continue
			}
			if dbTypeGot {
				engineType, _ := polarDBEngineInfo(supportedEngine.Engine)
				if normalizePolarDBEngineType(engineType) != normalizePolarDBEngineType(dbType.(string)) {
					continue
				}
			}
			if dbVersionGot {
				_, version := polarDBEngineInfo(supportedEngine.Engine)
				if version != strings.ToLower(dbVersion.(string)) {
					continue
				}
			}
			var dbNodeClasses []map[string]string
			for _, availableResource := range supportedEngine.AvailableResources {
				if "" != checkDBNodeClass && availableResource.DBNodeClass != checkDBNodeClass {
					continue
				}
				dbNodeClass := map[string]string{"db_node_class": availableResource.DBNodeClass}

				if "" != category {
					// 匹配过滤条件，返回符合条件的数据
					resultCategory := availableResource.Category
					if category == resultCategory {
						dbNodeClasses = append(dbNodeClasses, dbNodeClass)
					}
				} else {
					// category是空没有过滤条件返回所有数据
					dbNodeClasses = append(dbNodeClasses, dbNodeClass)
				}

			}
			// 过滤掉不支持的可用区数据
			if len(dbNodeClasses) == 0 {
				continue
			}
			availableResources := map[string]interface{}{
				"engine":              supportedEngine.Engine,
				"available_resources": dbNodeClasses,
			}
			supportedEngines = append(supportedEngines, availableResources)
			ids = append(ids, supportedEngine.Engine)
		}

		var availableClass map[string]interface{}
		if len(supportedEngines) > 0 {

			availableClass = map[string]interface{}{
				"zone_id":           zondId,
				"supported_engines": supportedEngines,
			}
		}
		if len(availableClass) > 0 {
			availableClasses = append(availableClasses, availableClass)
		}

	}
	d.SetId(dataResourceIdHash(ids))

	err = d.Set("classes", availableClasses)
	if err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), availableClasses)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

// polarDBEngineInfo splits a PolarDB engine identifier returned by the
// DescribeDBClusterAvailableResources API into its engine prefix and the
// remaining version suffix. The API returns the engine as the full DB type
// name followed by the version (e.g. "PostgreSQL11", "MySQL5.6", "Oracle");
// some responses may also use compact forms (e.g. "pg14", "mysql8"). Splitting
// at the first digit preserves version suffixes that contain dots.
func polarDBEngineInfo(engine string) (engineType string, version string) {
	e := strings.ToLower(engine)
	i := 0
	for i < len(e) && (e[i] < '0' || e[i] > '9') {
		i++
	}
	return e[:i], e[i:]
}

// normalizePolarDBEngineType maps a DB type to the canonical engine prefix,
// applied to both the API-returned engine prefix (from polarDBEngineInfo) and
// the user-specified db_type argument so the two are comparable. It handles
// full-name forms ("PostgreSQL", "MySQL", "Oracle") and compact forms ("pg")
// alike, mapping them to "pg", "mysql", "oracle". Unrecognized values are
// lower-cased as-is.
func normalizePolarDBEngineType(dbType string) string {
	switch strings.ToLower(dbType) {
	case "postgresql", "pg":
		return "pg"
	case "mysql":
		return "mysql"
	case "oracle":
		return "oracle"
	default:
		return strings.ToLower(dbType)
	}
}
