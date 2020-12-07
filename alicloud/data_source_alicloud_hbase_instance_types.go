package alicloud

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudHBaseInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHBaseInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHBaseInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := hbase.CreateDescribeInstanceTypeRequest()
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = strings.TrimSpace(instanceType.(string))
	}

	raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
		return hbaseClient.DescribeInstanceType(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbase_instance_types", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*hbase.DescribeInstanceTypeResponse)
	var ids []string
	var types []map[string]interface{}
	for _, spec := range response.InstanceTypeSpecList.InstanceTypeSpec {
		e := map[string]interface{}{
			"value":    spec.InstanceType,
			"cpu_size": spec.CpuSize,
			"mem_size": spec.MemSize,
		}
		ids = append(ids, spec.InstanceType)
		types = append(types, e)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("types", types); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		_ = writeToFile(output.(string), types)
	}
	return nil
}
