package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudMongodbAuditPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMongodbAuditPoliciesRead,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
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
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMongodbAuditPoliciesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := MongodbServiceV2{client}
	dbInstanceId := d.Get("db_instance_id").(string)

	// DescribeMongoDBLogConfig — provides ServiceType, TtlForStandard, HotTtlForV2Standard.
	logConfig, err := mongodbServiceV2.DescribeMongodbAuditPolicy(dbInstanceId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("MongodbAuditPolicy")
			return nil
		}
		return WrapError(err)
	}

	// DescribeAuditPolicy — provides LogAuditStatus (enable/disabled).
	auditPolicy, err := mongodbServiceV2.DescribeAuditPolicyDescribeAuditPolicy(dbInstanceId)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	mapping := map[string]interface{}{
		"id":                 dbInstanceId,
		"db_instance_id":     dbInstanceId,
		"service_type":       logConfig["ServiceType"],
		"storage_period":     logConfig["TtlForStandard"],
		"hot_storage_period": logConfig["HotTtlForV2Standard"],
	}
	if v, ok := auditPolicy["LogAuditStatus"].(string); ok {
		mapping["audit_status"] = convertMongodbAuditPolicyResponse(v)
	}

	s := []map[string]interface{}{mapping}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
