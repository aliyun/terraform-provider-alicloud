// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCSKubernetesPolicyInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCSKubernetesPolicyInstanceCreate,
		Read:   resourceAliCloudCSKubernetesPolicyInstanceRead,
		Update: resourceAliCloudCSKubernetesPolicyInstanceUpdate,
		Delete: resourceAliCloudCSKubernetesPolicyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespaces": {
				Type:     schema.TypeList,
				Optional: true,
				//Computed: true,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"warn", "deny"}, false),
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCSKubernetesPolicyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return WrapError(clientErr)
	}
	csClient := CsClient{client}

	action := "DeployPolicyInstance"
	cluster_id := d.Get("cluster_id").(string)
	policy_name := d.Get("policy_name").(string)

	// Wait for required addons to be ready
	if err := csClient.WaitForRequiredAddons(cluster_id, 10*time.Minute); err != nil {
		return WrapErrorf(err, "waiting for required addons failed")
	}

	createReq := &cs.DeployPolicyInstanceRequest{
		Action: tea.String(d.Get("action").(string)),
	}

	if v, ok := d.GetOk("namespaces"); ok {
		namespacesMapsArray := convertToInterfaceArray(v)
		stringList := make([]*string, len(namespacesMapsArray))
		for i, v := range namespacesMapsArray {
			stringList[i] = tea.String(v.(string))
		}

		createReq.Namespaces = stringList
	}

	if v, ok := d.GetOk("parameters"); ok {
		parametersMap := v.(map[string]interface{})
		createReq.Parameters = NormalizeMap(parametersMap)
	}

	var response *cs.DeployPolicyInstanceResponse
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.DeployPolicyInstance(tea.String(cluster_id), tea.String(policy_name), createReq)
		if err != nil {
			if NeedRetry(err) || strings.Contains(err.Error(), "the object has been modified") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_policy_instance", action, AlibabaCloudSdkGoERROR)
	}

	var instance_name string
	if response != nil && response.Body != nil && len(response.Body.Instances) == 1 {
		instance_name = tea.StringValue(response.Body.Instances[0])
	} else {
		return WrapErrorf(fmt.Errorf("no instance_name returned"), DefaultErrorMsg, "alicloud_cs_kubernetes_policy_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", cluster_id, policy_name, instance_name))

	return resourceAliCloudCSKubernetesPolicyInstanceRead(d, meta)
}

func resourceAliCloudCSKubernetesPolicyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ackServiceV2 := AckServiceV2{client}
	csClient, clientErr := client.NewRoaCsClient()
	if clientErr != nil {
		return WrapError(clientErr)
	}
	csClientWrapper := CsClient{csClient}

	object, err := ackServiceV2.DescribeAckPolicyInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cs_kubernetes_policy_instance DescribeAckPolicyInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	clusterId := object["cluster_id"].(string)
	// Wait for required addons to be ready
	if err := csClientWrapper.WaitForRequiredAddons(clusterId, 10*time.Minute); err != nil {
		return WrapErrorf(err, "waiting for required addons failed")
	}

	d.Set("action", object["policy_action"])
	d.Set("cluster_id", object["cluster_id"])
	d.Set("instance_name", object["instance_name"])
	d.Set("policy_name", object["policy_name"])

	// Process namespaces - convert string to array
	if policyScope, ok := object["policy_scope"].(string); ok {
		if policyScope == "*" {
			// Convert "*" to empty array
			d.Set("namespaces", []interface{}{})
		} else {
			// Split comma-separated values into array
			namespaces := strings.Split(policyScope, ",")
			result := make([]interface{}, len(namespaces))
			for i, ns := range namespaces {
				result[i] = strings.TrimSpace(ns) // Trim whitespace
			}
			d.Set("namespaces", result)
		}
	}

	// Process parameters - convert YAML string to map
	if policyParameters, ok := object["policy_parameters"].(string); ok {
		paramsMap, err := convertYamlToObject(policyParameters)
		if err == nil {
			// Apply type conversion to the parameters
			convertedParams := NormalizeMap(paramsMap)
			d.Set("parameters", convertedParams)
		} else {
			log.Printf("[WARN] Failed to parse policy_parameters as YAML: %v", err)
		}
	}

	return nil
}

func resourceAliCloudCSKubernetesPolicyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return WrapError(clientErr)
	}
	csClient := CsClient{client}

	action := "ModifyPolicyInstance"
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 3, len(parts)))
	}

	cluster_id := parts[0]
	policy_name := parts[1]
	instance_name := parts[2]

	// Wait for required addons to be ready
	if err := csClient.WaitForRequiredAddons(cluster_id, 10*time.Minute); err != nil {
		return WrapErrorf(err, "waiting for required addons failed")
	}

	updateReq := &cs.ModifyPolicyInstanceRequest{
		InstanceName: tea.String(instance_name),
	}

	if v, ok := d.GetOk("namespaces"); ok {
		namespacesMapsArray := convertToInterfaceArray(v)
		stringList := make([]*string, len(namespacesMapsArray))
		for i, v := range namespacesMapsArray {
			stringList[i] = tea.String(v.(string))
		}
		updateReq.Namespaces = stringList
	}

	if v, ok := d.GetOk("action"); ok {
		updateReq.Action = tea.String(v.(string))
	}

	if v, ok := d.GetOk("parameters"); ok {
		parametersMap := v.(map[string]interface{})
		updateReq.Parameters = NormalizeMap(parametersMap)
	}

	var response *cs.ModifyPolicyInstanceResponse
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.ModifyPolicyInstance(tea.String(cluster_id), tea.String(policy_name), updateReq)
		if err != nil {
			if NeedRetry(err) || strings.Contains(err.Error(), "the object has been modified") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudCSKubernetesPolicyInstanceRead(d, meta)
}
func resourceAliCloudCSKubernetesPolicyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		WrapError(clientErr)
	}
	csClient := CsClient{client}

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		err := WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 3, len(parts)))
		return err
	}

	cluster_id := parts[0]
	policy_name := parts[1]
	action := "DeletePolicyInstance"
	deleteReq := &cs.DeletePolicyInstanceRequest{
		InstanceName: tea.String(parts[2]),
	}

	// Wait for required addons to be ready
	if err := csClient.WaitForRequiredAddons(cluster_id, 10*time.Minute); err != nil {
		return WrapErrorf(err, "waiting for required addons failed")
	}

	var response *cs.DeletePolicyInstanceResponse
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.DeletePolicyInstance(tea.String(cluster_id), tea.String(policy_name), deleteReq)
		if err != nil {
			if NeedRetry(err) || strings.Contains(err.Error(), "the object has been modified") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, deleteReq)

	if err != nil {
		if IsExpectedErrors(err, []string{"404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func (s *CsClient) WaitForRequiredAddons(clusterId string, timeout time.Duration) error {
	requiredAddons := []string{"gatekeeper", "policy-template-controller"}
	for _, addonName := range requiredAddons {
		addon, err := s.GetCsKubernetesAddonInstance(clusterId, addonName)
		if err != nil {
			if !NotFoundError(err) {
				return WrapErrorf(err, "checking addon %s status failed", addonName)
			}
			// Addon not found, need to wait for installation
		} else if addon.Status == "running" || addon.Status == "active" {
			// Addon is already active, no need to wait
			continue
		}

		// Addon exists but not active, wait for it to become active
		stateConf := BuildStateConf([]string{}, []string{"running", "active"}, timeout, 10*time.Second, s.CsKubernetesAddonStateRefreshFunc(clusterId, addonName, []string{"failed", "error"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, "waiting for addon %s to be ready failed", addonName)
		}
	}
	return nil
}
