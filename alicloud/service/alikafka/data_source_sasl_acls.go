package alikafka

import (
	"fmt"

	alikafkasdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceSaslAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSaslAclsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"acl_resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Group", "Topic"}, false),
			},
			"acl_resource_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_resource_pattern_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_operation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSaslAclsRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := NewAlikafkaService(client)

	request := alikafkasdk.CreateDescribeAclsRequest()
	request.InstanceId = d.Get("instance_id").(string)
	request.RegionId = client.RegionId
	request.Username = d.Get("username").(string)
	request.AclResourceType = d.Get("acl_resource_type").(string)
	request.AclResourceName = d.Get("acl_resource_name").(string)

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
		return alikafkaClient.DescribeAcls(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"BIZ_SUBSCRIPTION_NOT_FOUND", "BIZ_TOPIC_NOT_FOUND"}) {
			var emptyValue []alikafkasdk.KafkaAclVO
			return alikafkaSaslAclsDecriptionAttributes(d, emptyValue, meta)
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_alikafka_sasl_acls", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alikafkasdk.DescribeAclsResponse)

	return alikafkaSaslAclsDecriptionAttributes(d, response.KafkaAclList.KafkaAclVO, meta)
}

func alikafkaSaslAclsDecriptionAttributes(d *schema.ResourceData, kafkaAclsInfo []alikafkasdk.KafkaAclVO, meta interface{}) error {

	var names []string
	var s []map[string]interface{}

	instanceId := d.Get("instance_id").(string)

	for _, item := range kafkaAclsInfo {
		mapping := map[string]interface{}{
			"id":                        fmt.Sprintf("%s:%s:%s:%s:%s:%s", instanceId, item.Username, item.AclResourceType, item.AclResourceName, item.AclResourcePatternType, item.AclOperationType),
			"username":                  item.Username,
			"acl_resource_type":         item.AclResourceType,
			"acl_resource_name":         item.AclResourceName,
			"acl_resource_pattern_type": item.AclResourcePatternType,
			"host":                      item.Host,
			"acl_operation_type":        item.AclOperationType,
		}

		name := fmt.Sprintf("%s:%s:%s:%s:%s", item.Username, item.AclResourceType, item.AclResourceName, item.AclResourcePatternType, item.AclOperationType)
		names = append(names, name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))

	if err := d.Set("acls", s); err != nil {
		return errs.WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
