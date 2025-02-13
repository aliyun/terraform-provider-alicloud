package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudMongodbShardingNetworkPrivateAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongodbShardingNetworkPrivateAddressCreate,
		Read:   resourceAliCloudMongodbShardingNetworkPrivateAddressRead,
		Update: resourceAliCloudMongodbShardingNetworkPrivateAddressUpdate,
		Delete: resourceAliCloudMongodbShardingNetworkPrivateAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^[\w!@#$%^&*()_+=]{6,32}$`), "The account password must be 6 to 32 characters in length, and can contain letters, digits, and special characters（!@#$%^&*()_+-=)."),
			},
			"network_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudMongodbShardingNetworkPrivateAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	var response map[string]interface{}
	action := "AllocateNodePrivateNetworkAddress"
	request := make(map[string]interface{})
	var err error

	request["DBInstanceId"] = d.Get("db_instance_id")
	request["NodeId"] = d.Get("node_id")
	request["ZoneId"] = d.Get("zone_id")

	if v, ok := d.GetOk("account_name"); ok {
		request["AccountName"] = v
	}

	if v, ok := d.GetOk("account_password"); ok {
		request["AccountPassword"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_sharding_network_private_address", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], request["NodeId"]))

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ddsService.RdsMongodbDBInstanceStateRefreshFunc(fmt.Sprint(request["DBInstanceId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAliCloudMongodbShardingNetworkPrivateAddressRead(d, meta)
}

func resourceAliCloudMongodbShardingNetworkPrivateAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	object, err := ddsService.DescribeMongodbShardingNetworkPrivateAddress(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_sharding_network_private_address MongoDBService.DescribeMongodbShardingNetworkPrivateAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("db_instance_id", parts[0])
	d.Set("node_id", object["NodeId"])

	networkAddressMaps := make([]map[string]interface{}, 0)
	for _, networkAddressArg := range object["NetworkAddress"].([]map[string]interface{}) {
		networkAddressMap := map[string]interface{}{}

		if nodeId, ok := networkAddressArg["NodeId"]; ok {
			networkAddressMap["node_id"] = nodeId
		}

		if nodeType, ok := networkAddressArg["NodeType"]; ok {
			networkAddressMap["node_type"] = nodeType
		}

		if role, ok := networkAddressArg["Role"]; ok {
			networkAddressMap["role"] = role
		}

		if vpcId, ok := networkAddressArg["VPCId"]; ok {
			networkAddressMap["vpc_id"] = vpcId
		}

		if vswitchId, ok := networkAddressArg["VswitchId"]; ok {
			networkAddressMap["vswitch_id"] = vswitchId
		}

		if networkType, ok := networkAddressArg["NetworkType"]; ok {
			networkAddressMap["network_type"] = networkType
		}

		if networkAddress, ok := networkAddressArg["NetworkAddress"]; ok {
			networkAddressMap["network_address"] = networkAddress
		}

		if ipAddress, ok := networkAddressArg["IPAddress"]; ok {
			networkAddressMap["ip_address"] = ipAddress
		}

		if port, ok := networkAddressArg["Port"]; ok {
			networkAddressMap["port"] = port
		}

		if expiredTime, ok := networkAddressArg["ExpiredTime"]; ok {
			networkAddressMap["expired_time"] = expiredTime
		}

		networkAddressMaps = append(networkAddressMaps, networkAddressMap)
	}

	d.Set("network_address", networkAddressMaps)

	shardingInstanceObject, err := ddsService.DescribeMongoDBShardingInstance(parts[0])
	if err != nil {
		return WrapError(err)
	}

	d.Set("zone_id", shardingInstanceObject["ZoneId"])

	return nil
}

func resourceAliCloudMongodbShardingNetworkPrivateAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAliCloudMongodbShardingNetworkPrivateAddressRead(d, meta)
}

func resourceAliCloudMongodbShardingNetworkPrivateAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}
	action := "ReleaseNodePrivateNetworkAddress"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"NodeId":       parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidStatus.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ddsService.RdsMongodbDBInstanceStateRefreshFunc(fmt.Sprint(parts[0]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return nil
}
