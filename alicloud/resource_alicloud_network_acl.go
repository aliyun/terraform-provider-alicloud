// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcNetworkAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcNetworkAclCreate,
		Read:   resourceAliCloudVpcNetworkAclRead,
		Update: resourceAliCloudVpcNetworkAclUpdate,
		Delete: resourceAliCloudVpcNetworkAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"egress_acl_entries": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"accept", "drop"}, false),
						},
						"destination_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringLenBetween(1, 256),
						},
						"entry_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"icmp", "gre", "tcp", "udp", "all"}, false),
						},
						"network_acl_entry_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringLenBetween(1, 128),
						},
						"network_acl_entry_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ingress_acl_entries": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"accept", "drop"}, false),
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringLenBetween(1, 256),
						},
						"entry_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"icmp", "gre", "tcp", "udp", "all"}, false),
						},
						"network_acl_entry_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringLenBetween(1, 128),
						},
						"network_acl_entry_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"network_acl_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
				ValidateFunc:  StringLenBetween(2, 128),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"source_network_acl_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'name' has been deprecated since provider version 1.122.0. New field 'network_acl_name' instead.",
				ValidateFunc: StringLenBetween(2, 128),
			},
		},
	}
}

func resourceAliCloudVpcNetworkAclCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateNetworkAcl"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("name"); ok || d.HasChange("name") {
		request["NetworkAclName"] = v
	}

	if v, ok := d.GetOk("network_acl_name"); ok && len(v.(string)) > 0 {
		request["NetworkAclName"] = v
	}
	if v, ok := d.GetOk("description"); ok && len(v.(string)) > 0 {
		request["Description"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_network_acl", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.NetworkAclAttribute.NetworkAclId", response)
	d.SetId(fmt.Sprint(id))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcNetworkAclStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcNetworkAclUpdate(d, meta)
}

func resourceAliCloudVpcNetworkAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcNetworkAcl(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_network_acl DescribeVpcNetworkAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("network_acl_name", objectRaw["NetworkAclName"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_id", objectRaw["VpcId"])

	egressAclEntryRaw, _ := jsonpath.Get("$.EgressAclEntries.EgressAclEntry", objectRaw)
	egressAclEntriesMaps := make([]map[string]interface{}, 0)
	if egressAclEntryRaw != nil {
		for _, egressAclEntryChildRaw := range egressAclEntryRaw.([]interface{}) {
			egressAclEntriesMap := make(map[string]interface{})
			egressAclEntryChildRaw := egressAclEntryChildRaw.(map[string]interface{})
			if egressAclEntryChildRaw["EntryType"] == "service" || egressAclEntryChildRaw["EntryType"] == "system" {
				continue
			}
			egressAclEntriesMap["description"] = egressAclEntryChildRaw["Description"]
			egressAclEntriesMap["destination_cidr_ip"] = egressAclEntryChildRaw["DestinationCidrIp"]
			egressAclEntriesMap["entry_type"] = egressAclEntryChildRaw["EntryType"]
			egressAclEntriesMap["ip_version"] = egressAclEntryChildRaw["IpVersion"]
			egressAclEntriesMap["network_acl_entry_id"] = egressAclEntryChildRaw["NetworkAclEntryId"]
			egressAclEntriesMap["network_acl_entry_name"] = egressAclEntryChildRaw["NetworkAclEntryName"]
			egressAclEntriesMap["policy"] = egressAclEntryChildRaw["Policy"]
			egressAclEntriesMap["port"] = egressAclEntryChildRaw["Port"]
			egressAclEntriesMap["protocol"] = egressAclEntryChildRaw["Protocol"]

			egressAclEntriesMaps = append(egressAclEntriesMaps, egressAclEntriesMap)
		}
	}
	if err := d.Set("egress_acl_entries", egressAclEntriesMaps); err != nil {
		return err
	}
	ingressAclEntryRaw, _ := jsonpath.Get("$.IngressAclEntries.IngressAclEntry", objectRaw)
	ingressAclEntriesMaps := make([]map[string]interface{}, 0)
	if ingressAclEntryRaw != nil {
		for _, ingressAclEntryChildRaw := range ingressAclEntryRaw.([]interface{}) {
			ingressAclEntriesMap := make(map[string]interface{})
			ingressAclEntryChildRaw := ingressAclEntryChildRaw.(map[string]interface{})
			if ingressAclEntryChildRaw["EntryType"] == "service" || ingressAclEntryChildRaw["EntryType"] == "system" {
				continue
			}
			ingressAclEntriesMap["description"] = ingressAclEntryChildRaw["Description"]
			ingressAclEntriesMap["entry_type"] = ingressAclEntryChildRaw["EntryType"]
			ingressAclEntriesMap["ip_version"] = ingressAclEntryChildRaw["IpVersion"]
			ingressAclEntriesMap["network_acl_entry_id"] = ingressAclEntryChildRaw["NetworkAclEntryId"]
			ingressAclEntriesMap["network_acl_entry_name"] = ingressAclEntryChildRaw["NetworkAclEntryName"]
			ingressAclEntriesMap["policy"] = ingressAclEntryChildRaw["Policy"]
			ingressAclEntriesMap["port"] = ingressAclEntryChildRaw["Port"]
			ingressAclEntriesMap["protocol"] = ingressAclEntryChildRaw["Protocol"]
			ingressAclEntriesMap["source_cidr_ip"] = ingressAclEntryChildRaw["SourceCidrIp"]

			ingressAclEntriesMaps = append(ingressAclEntriesMaps, ingressAclEntriesMap)
		}
	}
	if err := d.Set("ingress_acl_entries", ingressAclEntriesMaps); err != nil {
		return err
	}
	resourceRaw, _ := jsonpath.Get("$.Resources.Resource", objectRaw)
	resourcesMaps := make([]map[string]interface{}, 0)
	if resourceRaw != nil {
		for _, resourceChildRaw := range resourceRaw.([]interface{}) {
			resourcesMap := make(map[string]interface{})
			resourceChildRaw := resourceChildRaw.(map[string]interface{})
			resourcesMap["resource_id"] = resourceChildRaw["ResourceId"]
			resourcesMap["resource_type"] = resourceChildRaw["ResourceType"]
			resourcesMap["status"] = resourceChildRaw["Status"]

			resourcesMaps = append(resourcesMaps, resourcesMap)
		}
	}
	if err := d.Set("resources", resourcesMaps); err != nil {
		return err
	}
	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("name", d.Get("network_acl_name"))
	return nil
}

func resourceAliCloudVpcNetworkAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyNetworkAclAttributes"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["NetworkAclId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["NetworkAclName"] = d.Get("name")
	}

	if !d.IsNewResource() && d.HasChange("network_acl_name") {
		update = true
		request["NetworkAclName"] = d.Get("network_acl_name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
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
		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcNetworkAclStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateNetworkAclEntries"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["NetworkAclId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("ingress_acl_entries") {
		request["UpdateIngressAclEntries"] = "true"
		update = true
		if v, ok := d.GetOk("ingress_acl_entries"); ok || d.HasChange("ingress_acl_entries") {
			ingressAclEntriesMapsArray := make([]interface{}, 0)
			for _, dataLoop := range v.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Policy"] = dataLoopTmp["policy"]
				if len(dataLoopTmp["network_acl_entry_name"].(string)) > 0 {
					dataLoopMap["NetworkAclEntryName"] = dataLoopTmp["network_acl_entry_name"]
				}
				dataLoopMap["SourceCidrIp"] = dataLoopTmp["source_cidr_ip"]
				dataLoopMap["Protocol"] = dataLoopTmp["protocol"]
				dataLoopMap["Port"] = dataLoopTmp["port"]
				if len(dataLoopTmp["description"].(string)) > 0 {
					dataLoopMap["Description"] = dataLoopTmp["description"]
				}
				dataLoopMap["IpVersion"] = dataLoopTmp["ip_version"]
				dataLoopMap["EntryType"] = dataLoopTmp["entry_type"]
				if len(dataLoopTmp["network_acl_entry_id"].(string)) > 0 {
					dataLoopMap["NetworkAclEntryId"] = dataLoopTmp["network_acl_entry_id"]
				}
				ingressAclEntriesMapsArray = append(ingressAclEntriesMapsArray, dataLoopMap)
			}
			request["IngressAclEntries"] = ingressAclEntriesMapsArray
		}
	}

	if d.HasChange("egress_acl_entries") {
		request["UpdateEgressAclEntries"] = "true"
		update = true
		if v, ok := d.GetOk("egress_acl_entries"); ok || d.HasChange("egress_acl_entries") {
			egressAclEntriesMapsArray := make([]interface{}, 0)
			for _, dataLoop1 := range v.([]interface{}) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["Policy"] = dataLoop1Tmp["policy"]
				if len(dataLoop1Tmp["network_acl_entry_name"].(string)) > 0 {
					dataLoop1Map["NetworkAclEntryName"] = dataLoop1Tmp["network_acl_entry_name"]
				}
				if len(dataLoop1Tmp["description"].(string)) > 0 {
					dataLoop1Map["Description"] = dataLoop1Tmp["description"]
				}
				dataLoop1Map["Protocol"] = dataLoop1Tmp["protocol"]
				dataLoop1Map["DestinationCidrIp"] = dataLoop1Tmp["destination_cidr_ip"]
				dataLoop1Map["Port"] = dataLoop1Tmp["port"]
				dataLoop1Map["EntryType"] = dataLoop1Tmp["entry_type"]
				dataLoop1Map["IpVersion"] = dataLoop1Tmp["ip_version"]
				if len(dataLoop1Tmp["network_acl_entry_id"].(string)) > 0 {
					dataLoop1Map["NetworkAclEntryId"] = dataLoop1Tmp["network_acl_entry_id"]
				}
				egressAclEntriesMapsArray = append(egressAclEntriesMapsArray, dataLoop1Map)
			}
			request["EgressAclEntries"] = egressAclEntriesMapsArray
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
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
		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcNetworkAclStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "CopyNetworkAclEntries"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["NetworkAclId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	request["SourceNetworkAclId"] = d.Get("source_network_acl_id")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
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
		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, vpcServiceV2.VpcNetworkAclStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("resources") {
		oldEntry, newEntry := d.GetChange("resources")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action := "UnassociateNetworkAcl"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["NetworkAclId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := removed.List()
			resourceMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["ResourceType"] = dataLoopTmp["resource_type"]
				dataLoopMap["ResourceId"] = dataLoopTmp["resource_id"]
				resourceMapsArray = append(resourceMapsArray, dataLoopMap)
			}
			request["Resource"] = resourceMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy", "ResourceStatus.Error", "NetworkAclExistBinding"}) || NeedRetry(err) {
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
			vpcServiceV2 := VpcServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcNetworkAclStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if added.Len() > 0 {
			action := "AssociateNetworkAcl"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["NetworkAclId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := added.List()
			resourceMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["ResourceType"] = dataLoopTmp["resource_type"]
				dataLoopMap["ResourceId"] = dataLoopTmp["resource_id"]
				resourceMapsArray = append(resourceMapsArray, dataLoopMap)
			}
			request["Resource"] = resourceMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
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
			vpcServiceV2 := VpcServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcNetworkAclStateRefreshFunc(d.Id(), "Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

	}
	if d.HasChange("tags") {
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "NETWORKACL"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudVpcNetworkAclRead(d, meta)
}

func resourceAliCloudVpcNetworkAclDelete(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*connectivity.AliyunClient)}
	_, errMsg := vpcService.DeleteAclResources(d.Id())
	if errMsg != nil {
		return WrapError(errMsg)
	}

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteNetworkAcl"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["NetworkAclId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"NetworkAclExistBinding", "IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
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

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcNetworkAclStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
