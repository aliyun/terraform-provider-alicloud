package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlicloudMongodbShardingAuditFilters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMongodbShardingAuditFiltersRead,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"audit_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"hot_storage_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMongodbShardingAuditFiltersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := MongodbServiceV2{client}
	dbInstanceId := d.Get("db_instance_id").(string)

	// DescribeMongoDBLogConfig — provides ServiceType, TtlForStandard, HotTtlForV2Standard.
	logConfig, err := mongodbServiceV2.DescribeMongodbAuditPolicy(dbInstanceId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("MongodbShardingAuditFilter")
			return nil
		}
		return WrapError(err)
	}

	// DescribeAuditPolicy — provides LogAuditStatus (enable/disabled).
	auditPolicy, err := mongodbServiceV2.DescribeAuditPolicyDescribeAuditPolicy(dbInstanceId)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	// DescribeAuditLogFilter — provides Filter / RoleType.
	auditLogFilter, err := mongodbServiceV2.DescribeAuditPolicyDescribeAuditLogFilter(dbInstanceId)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	// DescribeMongoDBLogConfig reports the retention in edition-specific fields: the Standard edition
	// keeps it in TtlForStandard, while the V2_Standard edition reports the cold storage duration in
	// TtlForV2Standard together with the hot storage duration in HotTtlForV2Standard. Reading
	// TtlForStandard for a V2_Standard audit log returns the untouched Standard-edition default, so the
	// data source selects the pair that belongs to the reported edition. The data source has no
	// declared edition of its own, so an unrecognised reported edition (the Trial value seen while the
	// audit log is disabled) falls back to the Standard fields.
	storagePeriod, hotStoragePeriod := mongodbShardingAuditFilterStorageValues(logConfig, fmt.Sprint(logConfig["ServiceType"]), "")

	mapping := map[string]interface{}{
		"id":             dbInstanceId,
		"db_instance_id": dbInstanceId,
		"service_type":   logConfig["ServiceType"],
		"region_id":      client.RegionId,
		"filter":         parseMongodbShardingAuditFilter(fmt.Sprint(auditLogFilter["Filter"]), ""),
		"role_type":      auditLogFilter["RoleType"],
	}
	// When the API omits a retention field, the edition selection yields no value and the computed
	// attribute reads back as zero; setting it is unconditional only for the fields the edition owns.
	if storagePeriod != nil {
		mapping["storage_period"] = storagePeriod
	}
	if hotStoragePeriod != nil {
		mapping["hot_storage_period"] = hotStoragePeriod
	}
	if v, ok := auditPolicy["LogAuditStatus"].(string); ok {
		mapping["audit_status"] = convertMongodbAuditPolicyResponse(v)
	}

	s := []map[string]interface{}{mapping}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("filters", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
