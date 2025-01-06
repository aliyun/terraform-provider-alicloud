package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMongoDBPublicNetworkAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongoDBPublicNetworkAddressCreate,
		Read:   resourceAlicloudMongoDBPublicNetworkAddressRead,
		Update: resourceAlicloudMongoDBPublicNetworkAddressUpdate,
		Delete: resourceAlicloudMongoDBPublicNetworkAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connection_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "3717",
			},
			"connection_strings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudMongoDBPublicNetworkAddressCreate(d *schema.ResourceData, meta interface{}) error {
	// only one public network address per instance.
	instanceId := d.Get("db_instance_id").(string)
	d.SetId(instanceId)

	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	if err := ddsService.AllocatePublicNetworkAddress(instanceId); err != nil {
		return err
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudMongoDBPublicNetworkAddressUpdate(d, meta)
}

func resourceAlicloudMongoDBPublicNetworkAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	object, err := ddsService.DescribeReplicaSetRole(d.Id())
	if err != nil {
		return err
	}

	if replicaSetsMap, ok := object["ReplicaSets"].(map[string]interface{}); ok && replicaSetsMap != nil {
		if replicaSetsList, ok := replicaSetsMap["ReplicaSet"]; ok && replicaSetsList != nil {
			replicaSetsMaps := make([]map[string]interface{}, 0)
			for _, replicaSets := range replicaSetsList.([]interface{}) {
				replicaSetsArg := replicaSets.(map[string]interface{})
				replicaSetsItemMap := make(map[string]interface{})

				if networkType, ok := replicaSetsArg["NetworkType"]; !ok || networkType.(string) != "Public" {
					continue
				}

				if connectionType, ok := replicaSetsArg["ConnectionType"]; ok && connectionType.(string) == "SRV" {
					continue
				}

				if connectionDomain, ok := replicaSetsArg["ConnectionDomain"]; ok {
					replicaSetsItemMap["connection_domain"] = connectionDomain
				}

				if connectionPort, ok := replicaSetsArg["ConnectionPort"]; ok {
					replicaSetsItemMap["connection_port"] = connectionPort
				}

				if role, ok := replicaSetsArg["ReplicaSetRole"]; ok {
					replicaSetsItemMap["role"] = role
				}

				if roleID, ok := replicaSetsArg["RoleId"]; ok {
					replicaSetsItemMap["role_id"] = roleID
				}

				replicaSetsMaps = append(replicaSetsMaps, replicaSetsItemMap)
			}

			d.Set("connection_strings", replicaSetsMaps)
		}
	}

	return nil
}

func resourceAlicloudMongoDBPublicNetworkAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	needUpdate := false
	connectionPrefixChanged := false
	if d.HasChange("connection_prefix") {
		needUpdate = true
		connectionPrefixChanged = true
	}
	if d.HasChange("port") {
		needUpdate = true
	}

	if needUpdate {
		connectionPrefix := ""
		port := 3717
		if v, ok := d.GetOk("connection_prefix"); ok {
			connectionPrefix = v.(string)
		}
		if v, ok := d.GetOk("port"); ok {
			switch v.(type) {
			case int:
				port = v.(int)
			case string:
				var err error
				port, err = strconv.Atoi(v.(string))
				if err != nil {
					return WrapError(fmt.Errorf("invalid port, expected int, got: %v", v))
				}
			}
		}

		if err := ddsService.ModifyAllPublicNetworkAddress(d, connectionPrefix, port, connectionPrefixChanged); err != nil {
			return err
		}
	}

	return resourceAlicloudMongoDBPublicNetworkAddressRead(d, meta)
}

func resourceAlicloudMongoDBPublicNetworkAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	if err := ddsService.ReleasePublicNetworkAddress(d.Id()); err != nil {
		return err
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return nil
}
