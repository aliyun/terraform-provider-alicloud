// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEnsSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsSecurityGroupCreate,
		Read:   resourceAliCloudEnsSecurityGroupRead,
		Update: resourceAliCloudEnsSecurityGroupUpdate,
		Delete: resourceAliCloudEnsSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      ensSecurityGroupPermissionHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Accept", "Drop"}, false),
						},
						"port_range": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_port_range": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 100),
						},
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"TCP", "UDP", "ICMP", "ALL"}, false),
						},
						"dest_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipv6_source_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipv6_dest_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"direction": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"ingress", "egress"}, false),
						},
					},
				},
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEnsSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSecurityGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("permissions"); ok {
		permissionsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := buildEnsPermissionMap(dataLoopTmp)
			if desc := dataLoopTmp["description"].(string); desc != "" {
				dataLoopMap["Description"] = desc
			}
			permissionsMapsArray = append(permissionsMapsArray, dataLoopMap)
		}
		permissionsMapsJson, err := json.Marshal(permissionsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Permissions"] = string(permissionsMapsJson)
	}

	if v, ok := d.GetOk("security_group_name"); ok {
		request["SecurityGroupName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_security_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SecurityGroupId"]))

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(response["SecurityGroupId"])}, d.Timeout(schema.TimeoutCreate), 30*time.Second, ensServiceV2.EnsSecurityGroupStateRefreshFunc(d.Id(), "SecurityGroupId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEnsSecurityGroupRead(d, meta)
}

func resourceAliCloudEnsSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsSecurityGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_security_group DescribeEnsSecurityGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("security_group_name", objectRaw["SecurityGroupName"])

	permissionRaw, _ := jsonpath.Get("$.Permissions.Permission", objectRaw)
	permissionsMaps := make([]map[string]interface{}, 0)
	if permissionRaw != nil {
		for _, permissionChildRaw := range convertToInterfaceArray(permissionRaw) {
			permissionsMap := make(map[string]interface{})
			permissionChildRaw := permissionChildRaw.(map[string]interface{})
			permissionsMap["creation_time"] = permissionChildRaw["CreationTime"]
			permissionsMap["description"] = permissionChildRaw["Description"]
			permissionsMap["dest_cidr_ip"] = permissionChildRaw["DestCidrIp"]
			permissionsMap["direction"] = permissionChildRaw["Direction"]
			permissionsMap["ip_protocol"] = permissionChildRaw["IpProtocol"]
			permissionsMap["ipv6_dest_cidr_ip"] = permissionChildRaw["Ipv6DestCidrIp"]
			permissionsMap["ipv6_source_cidr_ip"] = permissionChildRaw["Ipv6SourceCidrIp"]
			permissionsMap["policy"] = permissionChildRaw["Policy"]
			permissionsMap["port_range"] = permissionChildRaw["PortRange"]
			permissionsMap["priority"] = permissionChildRaw["Priority"]
			permissionsMap["source_cidr_ip"] = permissionChildRaw["SourceCidrIp"]
			permissionsMap["source_port_range"] = permissionChildRaw["SourcePortRange"]

			permissionsMaps = append(permissionsMaps, permissionsMap)
		}
	}
	if err := d.Set("permissions", permissionsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudEnsSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifySecurityGroupAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SecurityGroupId"] = d.Id()

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["Description"] = v
	}
	if d.HasChange("security_group_name") {
		update = true
	}
	if v, ok := d.GetOk("security_group_name"); ok || d.HasChange("security_group_name") {
		request["SecurityGroupName"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("permissions") {
		oldEntry, newEntry := d.GetChange("permissions")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action := "DeleteSecurityGroupPermissions"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["SecurityGroupId"] = d.Id()
			permissionsMapsArray := make([]interface{}, 0)
			for _, dataLoop := range removed.List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := buildEnsPermissionMap(dataLoopTmp)
				permissionsMapsArray = append(permissionsMapsArray, dataLoopMap)
			}
			permissionsMapsJson, err := json.Marshal(permissionsMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["Permissions"] = string(permissionsMapsJson)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		if added.Len() > 0 {
			action := "CreateSecurityGroupPermissions"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["SecurityGroupId"] = d.Id()
			permissionsMapsArray := make([]interface{}, 0)
			for _, dataLoop := range added.List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := buildEnsPermissionMap(dataLoopTmp)
				if desc := dataLoopTmp["description"].(string); desc != "" {
					dataLoopMap["Description"] = desc
				}
				permissionsMapsArray = append(permissionsMapsArray, dataLoopMap)
			}
			permissionsMapsJson, err := json.Marshal(permissionsMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["Permissions"] = string(permissionsMapsJson)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
	}

	return resourceAliCloudEnsSecurityGroupRead(d, meta)
}

func resourceAliCloudEnsSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSecurityGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SecurityGroupId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, ensServiceV2.EnsSecurityGroupStateRefreshFunc(d.Id(), "SecurityGroupId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

// buildEnsPermissionMap builds the permission payload for ENS SecurityGroup
// permission actions, skipping optional string fields that are empty. The
// ENS API rejects requests with InvalidParameter when an unset permission
// attribute is sent as an empty string (for example a v4-only rule leaves
// Ipv6SourceCidrIp/Ipv6DestCidrIp empty), so only set fields are forwarded.
// Description is intentionally not included here because the
// DeleteSecurityGroupPermissions action does not accept it; callers that
// need Description add it after invoking this helper.
func buildEnsPermissionMap(dataLoopTmp map[string]interface{}) map[string]interface{} {
	dataLoopMap := make(map[string]interface{})
	if priority, ok := dataLoopTmp["priority"].(int); ok && priority > 0 {
		dataLoopMap["Priority"] = priority
	}
	if v := dataLoopTmp["dest_cidr_ip"].(string); v != "" {
		dataLoopMap["DestCidrIp"] = v
	}
	if v := dataLoopTmp["direction"].(string); v != "" {
		dataLoopMap["Direction"] = v
	}
	if v := dataLoopTmp["ip_protocol"].(string); v != "" {
		dataLoopMap["IpProtocol"] = v
	}
	if v := dataLoopTmp["ipv6_dest_cidr_ip"].(string); v != "" {
		dataLoopMap["Ipv6DestCidrIp"] = v
	}
	if v := dataLoopTmp["ipv6_source_cidr_ip"].(string); v != "" {
		dataLoopMap["Ipv6SourceCidrIp"] = v
	}
	if v := dataLoopTmp["policy"].(string); v != "" {
		dataLoopMap["Policy"] = v
	}
	if v := dataLoopTmp["port_range"].(string); v != "" {
		dataLoopMap["PortRange"] = v
	}
	if v := dataLoopTmp["source_cidr_ip"].(string); v != "" {
		dataLoopMap["SourceCidrIp"] = v
	}
	if v := dataLoopTmp["source_port_range"].(string); v != "" {
		dataLoopMap["SourcePortRange"] = v
	}
	return dataLoopMap
}

// ensSecurityGroupPermissionHash computes the set identity hash for an ENS
// security group permission rule. It intentionally excludes creation_time
// (a server-populated Computed field) and source_port_range (Optional+Computed,
// defaulted by the API when the user leaves it unset) so that a rule written
// in config without those fields still hashes identically to the same rule
// read back from the API. Without this the TypeSet would show a perpetual
// diff, because the default schema.HashResource hashes Computed fields that
// are empty in config but populated in state. The remaining identity fields
// are sufficient to uniquely distinguish a permission within one security
// group.
func ensSecurityGroupPermissionHash(i interface{}) int {
	m, ok := i.(map[string]interface{})
	if !ok {
		return 0
	}
	return schema.HashString(fmt.Sprintf("%s|%s|%s|%s|%v|%s|%s|%s|%s|%s",
		m["direction"], m["ip_protocol"], m["port_range"], m["policy"], m["priority"],
		m["source_cidr_ip"], m["dest_cidr_ip"], m["ipv6_source_cidr_ip"], m["ipv6_dest_cidr_ip"], m["description"]))
}
