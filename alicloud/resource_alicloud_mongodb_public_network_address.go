package alicloud

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
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
			"connection_strings": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"connection_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_port": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					if id, ok := m["role_id"]; ok {
						buf.WriteString(fmt.Sprintf("%v", id))
					}
					return hashcode.String(buf.String())
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

func transferToReplicaSetMaps(replicaSetsList interface{}) map[string]map[string]interface{} {
	replicaSetsMaps := make(map[string]map[string]interface{})
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
			prefix, err := getPrefixOfConnectionDomain(connectionDomain.(string))
			if err == nil {
				replicaSetsItemMap["connection_prefix"] = prefix
			}
		}

		if connectionPort, ok := replicaSetsArg["ConnectionPort"]; ok {
			replicaSetsItemMap["connection_port"] = connectionPort
		}

		var roleID string
		if roleIDGot, ok := replicaSetsArg["RoleId"]; ok {
			replicaSetsItemMap["role_id"] = roleIDGot
			switch t := roleIDGot.(type) {
			case string:
				roleID = t
			default:
				roleID = fmt.Sprintf("%v", t)
			}
		}

		if role, ok := replicaSetsArg["ReplicaSetRole"]; ok {
			replicaSetsItemMap["role"] = role
		}

		replicaSetsMaps[roleID] = replicaSetsItemMap
	}
	return replicaSetsMaps
}

// Values creates an array of the map values.
func Values[K comparable, V any](in ...map[K]V) []V {
	size := 0
	for i := range in {
		size += len(in[i])
	}
	result := make([]V, 0, size)

	for i := range in {
		for k := range in[i] {
			result = append(result, in[i][k])
		}
	}

	return result
}

// MapToSlice transforms a map into a slice based on specific iteratee
func MapToSlice[K comparable, V any, R any](in map[K]V, iteratee func(key K, value V) R) []R {
	result := make([]R, 0, len(in))

	for k := range in {
		result = append(result, iteratee(k, in[k]))
	}

	return result
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
			replicaSetsMaps := transferToReplicaSetMaps(replicaSetsList)
			d.Set("connection_strings", Values(replicaSetsMaps))
		}
	}

	return nil
}

func resourceAlicloudMongoDBPublicNetworkAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	if !d.HasChange("connection_strings") {
		return resourceAlicloudMongoDBPublicNetworkAddressRead(d, meta)
	}

	_, current := d.GetChange("connection_strings")

	object, err := ddsService.DescribeReplicaSetRole(d.Id())
	if err != nil {
		return err
	}

	// todo: better to make this a util function.
	replicaSetsMaps := make(map[string]map[string]interface{})
	if replicaSetsMap, ok := object["ReplicaSets"].(map[string]interface{}); ok && replicaSetsMap != nil {
		if replicaSetsList, ok := replicaSetsMap["ReplicaSet"]; ok && replicaSetsList != nil {
			replicaSetsMaps = transferToReplicaSetMaps(replicaSetsList)
		}
	}

	changes := current.(*schema.Set).List()
	// changes := current.([]interface{})
	log.Printf("[WARN] changes: %v", changes)
	// changes := current.([]interface{})
	for _, change := range changes {
		changeMap, ok := change.(map[string]interface{})
		if !ok {
			continue
		}

		roleIDGot, ok := changeMap["role_id"]
		if !ok {
			return WrapError(fmt.Errorf("no role_id specified"))
			// or we could skip this modification?
			// log.Printf("[WARN] no role_id specified, skip this modification, change: %v", change)
			// continue
		}

		roleID := roleIDGot.(string)

		connectionInfo, ok := replicaSetsMaps[roleID]
		if !ok {
			return WrapError(fmt.Errorf("cannot find role_id in replica set, role_id: %v", roleID))
			// or we could skip this modification?
			// log.Printf("[WARN] no role_id specified, skip this modification, change: %v", change)
			// continue
		}

		connectionDomain, ok := connectionInfo["connection_domain"]
		if !ok {
			return WrapError(fmt.Errorf("no connection domain found in replica set, role_id: %v", roleID))
		}

		currentPrefix, err := getPrefixOfConnectionDomain(connectionDomain.(string))
		if err != nil {
			return WrapError(err)
		}

		var port int
		portNow, ok := connectionInfo["connection_port"]
		if ok {
			port, err = strconv.Atoi(portNow.(string))
			if err != nil {
				return WrapError(fmt.Errorf("port of connection info is invalid, port: %v", portNow))
			}
		}

		prefixToBeChanged := currentPrefix
		prefixChange, ok := changeMap["connection_prefix"]
		if ok {
			prefixToBeChanged = prefixChange.(string)
		}

		portToBeChanged := port
		portChange, ok := changeMap["connection_port"]
		if ok {
			portToBeChanged, err = strconv.Atoi(portChange.(string))
			if err != nil {
				return WrapError(fmt.Errorf("port of connection info is invalid, should be integer, got: %v", portChange))
			}
		}

		if prefixToBeChanged == currentPrefix && portToBeChanged == port {
			log.Printf("[INFO] no need to modify connection string, current prefix: %s, prefix to be changed: %s, port to be changed: %d",
				currentPrefix, prefixToBeChanged, portToBeChanged)
			continue
		}

		log.Printf("[INFO] ready to modify connection string, current prefix: %s, prefix to be changed: %s, port to be changed: %d",
			currentPrefix, prefixToBeChanged, portToBeChanged)
		if err := ddsService.ModifyDBInstanceConnectionString(d, connectionDomain.(string), prefixToBeChanged, portToBeChanged); err != nil {
			return WrapError(err)
		}
		log.Printf("[INFO] done to modify connection string, current prefix: %s, prefix to be changed: %s, port to be changed: %d",
			currentPrefix, prefixToBeChanged, portToBeChanged)
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
