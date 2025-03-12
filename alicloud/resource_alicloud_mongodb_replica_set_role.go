package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/samber/lo"
)

func resourceAlicloudMongoDBReplicaSetRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongoDBReplicaSetRoleCreate,
		Read:   resourceAlicloudMongoDBReplicaSetRoleRead,
		Update: resourceAlicloudMongoDBReplicaSetRoleUpdate,
		Delete: resourceAlicloudMongoDBReplicaSetRoleDelete,

		// we can't import this resource.

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connection_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replica_set_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIfChange("role_id", func(old, new, meta interface{}) bool {
				return old.(string) != new.(string)
			}),
			customdiff.ForceNewIfChange("connection_prefix", func(old, new, meta interface{}) bool {
				return old.(string) != new.(string)
			}),
			customdiff.ForceNewIfChange("network_type", func(old, new, meta interface{}) bool {
				return old.(string) != new.(string)
			}),
		),
	}
}

func makeKeyOfMongoReplica(instanceId, roleId string, networkType string) string {
	return fmt.Sprintf("%s:%s:%s", instanceId, roleId, networkType)
}

func parseKeyOfMongoReplica(key string) (instanceId, roleId, networkType string) {
	parts := strings.Split(key, ":")
	return parts[0], parts[1], parts[2]
}

func makeSecondaryIndexOfMongoReplica(roleId, networkType string) string {
	return fmt.Sprintf("%s:%s", roleId, networkType)
}

func resourceAlicloudMongoDBReplicaSetRoleCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlicloudMongoDBReplicaSetRoleUpdate(d, meta)
}

func resourceAlicloudMongoDBReplicaSetRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceId := d.Get("db_instance_id").(string)

	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	replicaSetsObjects, err := ddsService.DescribeReplicaSetRole(instanceId)
	if err != nil {
		return WrapError(err)
	}
	replicaSets := transferToMongoReplicaSets(replicaSetsObjects, false)

	replicasGroupByRoleIdAndNetworkType := lo.FilterSliceToMap(replicaSets, func(replica map[string]interface{}) (string, map[string]interface{}, bool) {
		roleID, roleIdOk := replica["role_id"]
		networkType, networkTypeOk := replica["network_type"]
		if roleIdOk && networkTypeOk {
			return makeSecondaryIndexOfMongoReplica(roleID.(string), networkType.(string)), replica, true
		}
		return "", replica, false
	})

	roleId := d.Get("role_id").(string)
	networkType := d.Get("network_type").(string)

	replica, ok := replicasGroupByRoleIdAndNetworkType[makeSecondaryIndexOfMongoReplica(roleId, networkType)]
	if !ok {
		return WrapError(fmt.Errorf("connection address not found, roleId: %s, network type: %s", roleId, networkType))
	}

	d.SetId(makeKeyOfMongoReplica(instanceId, roleId, networkType))

	currentConnectionDomain, ok := replica["connection_domain"]
	if !ok {
		// though should not happen.
		return WrapError(fmt.Errorf("current connection domain not found, roleId: %s, network type: %s", roleId, networkType))
	}

	currentConnectionPort, ok := replica["connection_port"]
	if !ok {
		// though should not happen.
		return WrapError(fmt.Errorf("current connection port not found, roleId: %s, network type: %s", roleId, networkType))
	}

	connectionChanged := false
	newConnectionPrefix := getPrefixOfConnectionDomain(currentConnectionDomain.(string))
	newPort := currentConnectionPort.(string)

	if d.HasChange("connection_prefix") {
		if connectionPrefix, ok := d.GetOk("connection_prefix"); ok && connectionPrefix.(string) != newConnectionPrefix {
			newConnectionPrefix = connectionPrefix.(string)
			connectionChanged = true
		}
	}

	if d.HasChange("connection_port") {
		if connectionPort, ok := d.GetOk("connection_port"); ok && connectionPort.(string) != newPort {
			newPort = connectionPort.(string)
			connectionChanged = true
		}
	}

	if connectionChanged {
		err = ddsService.ModifyDBInstanceConnectionString(d, instanceId, currentConnectionDomain.(string), newConnectionPrefix, newPort)
		if err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudMongoDBReplicaSetRoleRead(d, meta)
}

func resourceAlicloudMongoDBReplicaSetRoleRead(d *schema.ResourceData, meta interface{}) error {
	instanceId, roleId, networkType := parseKeyOfMongoReplica(d.Id())
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	replicaSetsObjects, err := ddsService.DescribeReplicaSetRole(instanceId)
	if err != nil {
		return WrapError(err)
	}
	replicaSets := transferToMongoReplicaSets(replicaSetsObjects, false)

	replicasGroupByRoleIdAndNetworkType := lo.FilterSliceToMap(replicaSets, func(replica map[string]interface{}) (string, map[string]interface{}, bool) {
		roleID, roleIdOk := replica["role_id"]
		networkType, networkTypeOk := replica["network_type"]
		if roleIdOk && networkTypeOk {
			return makeSecondaryIndexOfMongoReplica(roleID.(string), networkType.(string)), replica, true
		}
		return "", replica, false
	})

	replica, ok := replicasGroupByRoleIdAndNetworkType[makeSecondaryIndexOfMongoReplica(roleId, networkType)]
	if !ok {
		return WrapError(fmt.Errorf("connection address not found, roleId: %s, network type: %s", roleId, networkType))
	}

	if connectionDomain, ok := replica["connection_domain"]; ok {
		if err := d.Set("connection_domain", connectionDomain); err != nil {
			return WrapError(err)
		}
		if err := d.Set("connection_prefix", getPrefixOfConnectionDomain(connectionDomain.(string))); err != nil {
			return WrapError(err)
		}
	}

	if connectionPort, ok := replica["connection_port"]; ok {
		if err := d.Set("connection_port", connectionPort); err != nil {
			return WrapError(err)
		}
	}

	if replicaSetRole, ok := replica["replica_set_role"]; ok {
		if err := d.Set("replica_set_role", replicaSetRole); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func resourceAlicloudMongoDBReplicaSetRoleDelete(d *schema.ResourceData, meta interface{}) error {
	// dont do anything.
	return nil
}
