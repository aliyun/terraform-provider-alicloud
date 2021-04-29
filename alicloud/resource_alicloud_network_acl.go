package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudNetworkAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNetworkAclCreate,
		Read:   resourceAlicloudNetworkAclRead,
		Update: resourceAlicloudNetworkAclUpdate,
		Delete: resourceAlicloudNetworkAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"egress_acl_entries": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_acl_entry_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"icmp", "gre", "tcp", "udp", "all"}, false),
						},
					},
				},
			},
			"ingress_acl_entries": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_acl_entry_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"icmp", "gre", "tcp", "udp", "all"}, false),
						},
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"network_acl_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringLenBetween(2, 128),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.122.0. New field 'network_acl_name' instead",
				ConflictsWith: []string{"network_acl_name"},
				ValidateFunc:  validation.StringLenBetween(2, 128),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudNetworkAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateNetworkAcl"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("network_acl_name"); ok {
		request["NetworkAclName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["NetworkAclName"] = v
	}

	request["RegionId"] = client.RegionId
	request["VpcId"] = d.Get("vpc_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_network_acl", action, AlibabaCloudSdkGoERROR)
	}
	responseNetworkAclAttribute := response["NetworkAclAttribute"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseNetworkAclAttribute["NetworkAclId"]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudNetworkAclUpdate(d, meta)
}
func resourceAlicloudNetworkAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeNetworkAcl(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_network_acl vpcService.DescribeNetworkAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])

	egressAclEntry := make([]map[string]interface{}, 0)
	if egressAclEntryList, ok := object["EgressAclEntries"].(map[string]interface{})["EgressAclEntry"].([]interface{}); ok {
		for _, v := range egressAclEntryList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"description":            m1["Description"],
					"destination_cidr_ip":    m1["DestinationCidrIp"],
					"network_acl_entry_name": m1["NetworkAclEntryName"],
					"policy":                 m1["Policy"],
					"port":                   m1["Port"],
					"protocol":               m1["Protocol"],
				}
				egressAclEntry = append(egressAclEntry, temp1)

			}
		}
	}
	if err := d.Set("egress_acl_entries", egressAclEntry); err != nil {
		return WrapError(err)
	}

	ingressAclEntry := make([]map[string]interface{}, 0)
	if ingressAclEntryList, ok := object["IngressAclEntries"].(map[string]interface{})["IngressAclEntry"].([]interface{}); ok {
		for _, v := range ingressAclEntryList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"description":            m1["Description"],
					"network_acl_entry_name": m1["NetworkAclEntryName"],
					"policy":                 m1["Policy"],
					"port":                   m1["Port"],
					"protocol":               m1["Protocol"],
					"source_cidr_ip":         m1["SourceCidrIp"],
				}
				ingressAclEntry = append(ingressAclEntry, temp1)

			}
		}
	}
	if err := d.Set("ingress_acl_entries", ingressAclEntry); err != nil {
		return WrapError(err)
	}
	d.Set("network_acl_name", object["NetworkAclName"])
	d.Set("name", object["NetworkAclName"])
	d.Set("status", object["Status"])
	d.Set("vpc_id", object["VpcId"])
	return nil
}
func resourceAlicloudNetworkAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"NetworkAclId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("network_acl_name") {
		update = true
		request["NetworkAclName"] = d.Get("network_acl_name")
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["NetworkAclName"] = d.Get("name")
	}
	if update {
		action := "ModifyNetworkAclAttributes"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("description")
		d.SetPartial("name")
		d.SetPartial("network_acl_name")
	}
	update = false
	updateNetworkAclEntriesReq := map[string]interface{}{
		"NetworkAclId": d.Id(),
	}
	updateNetworkAclEntriesReq["RegionId"] = client.RegionId
	if d.HasChange("egress_acl_entries") {
		updateNetworkAclEntriesReq["UpdateEgressAclEntries"] = true
		update = true
		EgressAclEntries := make([]map[string]interface{}, len(d.Get("egress_acl_entries").(*schema.Set).List()))
		for i, EgressAclEntriesValue := range d.Get("egress_acl_entries").(*schema.Set).List() {
			EgressAclEntriesMap := EgressAclEntriesValue.(map[string]interface{})
			EgressAclEntries[i] = make(map[string]interface{})
			EgressAclEntries[i]["Description"] = EgressAclEntriesMap["description"]
			EgressAclEntries[i]["DestinationCidrIp"] = EgressAclEntriesMap["destination_cidr_ip"]
			EgressAclEntries[i]["NetworkAclEntryName"] = EgressAclEntriesMap["network_acl_entry_name"]
			EgressAclEntries[i]["Policy"] = EgressAclEntriesMap["policy"]
			EgressAclEntries[i]["Port"] = EgressAclEntriesMap["port"]
			EgressAclEntries[i]["Protocol"] = EgressAclEntriesMap["protocol"]
		}
		updateNetworkAclEntriesReq["EgressAclEntries"] = EgressAclEntries

	}
	if d.HasChange("ingress_acl_entries") {
		updateNetworkAclEntriesReq["UpdateIngressAclEntries"] = true
		update = true
		IngressAclEntries := make([]map[string]interface{}, len(d.Get("ingress_acl_entries").(*schema.Set).List()))
		for i, IngressAclEntriesValue := range d.Get("ingress_acl_entries").(*schema.Set).List() {
			IngressAclEntriesMap := IngressAclEntriesValue.(map[string]interface{})
			IngressAclEntries[i] = make(map[string]interface{})
			IngressAclEntries[i]["Description"] = IngressAclEntriesMap["description"]
			IngressAclEntries[i]["NetworkAclEntryName"] = IngressAclEntriesMap["network_acl_entry_name"]
			IngressAclEntries[i]["Policy"] = IngressAclEntriesMap["policy"]
			IngressAclEntries[i]["Port"] = IngressAclEntriesMap["port"]
			IngressAclEntries[i]["Protocol"] = IngressAclEntriesMap["protocol"]
			IngressAclEntries[i]["SourceCidrIp"] = IngressAclEntriesMap["source_cidr_ip"]
		}
		updateNetworkAclEntriesReq["IngressAclEntries"] = IngressAclEntries

	}
	if update {
		action := "UpdateNetworkAclEntries"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, updateNetworkAclEntriesReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"TaskConflict"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, updateNetworkAclEntriesReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("egress_acl_entries")
		d.SetPartial("ingress_acl_entries")
	}
	d.Partial(false)
	return resourceAlicloudNetworkAclRead(d, meta)
}
func resourceAlicloudNetworkAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteNetworkAcl"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"NetworkAclId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
