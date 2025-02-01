package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRdsClassDetails() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsRdsClassDetailsRead,
		Schema: map[string]*schema.Schema{
			"commodity_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"bards", "rds", "rords", "rds_rordspre_public_cn", "bards_intl", "rds_intl", "rords_intl", "rds_rordspre_public_intl"}, false),
			},
			"class_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: StringInSlice([]string{
					string(MySQL),
					string(SQLServer),
					string(PostgreSQL),
					string(MariaDB),
				}, false),
			},
			"max_iombps": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_connections": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"class_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instruction_set_arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_iops": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reference_price": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudRdsRdsClassDetailsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeClassDetails"
	request := map[string]interface{}{
		"SourceIp":      client.SourceIp,
		"RegionId":      client.RegionId,
		"CommodityCode": d.Get("commodity_code"),
		"ClassCode":     d.Get("class_code"),
		"EngineVersion": d.Get("engine_version"),
		"Engine":        d.Get("engine"),
	}
	id := ""
	if v, ok := d.GetOk("class_code"); ok {
		id = v.(string)
	}
	var response map[string]interface{}
	var err error
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_class_details", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "", response)
	}
	object := resp.(map[string]interface{})

	d.SetId(id)
	d.Set("max_iombps", object["MaxIOMBPS"])
	d.Set("max_connections", object["MaxConnections"])
	d.Set("class_group", object["ClassGroup"])
	d.Set("cpu", object["Cpu"])
	d.Set("instruction_set_arch", object["InstructionSetArch"])
	d.Set("memory_class", object["MemoryClass"])
	d.Set("max_iops", object["MaxIOPS"])
	d.Set("reference_price", object["ReferencePrice"])
	d.Set("category", object["Category"])
	d.Set("db_instance_storage_type", object["DBInstanceStorageType"])

	return nil
}
