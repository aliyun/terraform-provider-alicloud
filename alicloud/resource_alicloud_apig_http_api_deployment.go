package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApigHttpApiDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigHttpApiDeploymentCreate,
		Read:   resourceAliCloudApigHttpApiDeploymentRead,
		Delete: resourceAliCloudApigHttpApiDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAliCloudApigHttpApiDeploymentImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"http_api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigHttpApiDeploymentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts, err := apigParseCompositeID(d.Id(), 4)
	if err != nil {
		return nil, err
	}
	if err := d.Set("http_api_id", parts[0]); err != nil {
		return nil, err
	}
	if err := d.Set("route_id", parts[1]); err != nil {
		return nil, err
	}
	if err := d.Set("environment_id", parts[2]); err != nil {
		return nil, err
	}
	if err := d.Set("gateway_id", parts[3]); err != nil {
		return nil, err
	}
	d.SetId(strings.Join(parts, ":"))
	return []*schema.ResourceData{d}, nil
}

func resourceAliCloudApigHttpApiDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/v1/http-apis/%s/deploy", d.Get("http_api_id"))
	body := apigHttpApiDeploymentRequest(d)
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, nil, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_http_api_deployment", action, AlibabaCloudSdkGoERROR)
	}
	data, err := apigResponseData(response)
	if err != nil {
		return err
	}
	if httpAPIID, ok := data["httpApiId"].(string); !ok || httpAPIID != d.Get("http_api_id") {
		return fmt.Errorf("APIG DeployHttpApi response does not match http_api_id")
	}
	d.SetId(fmt.Sprintf("%s:%s:%s:%s", d.Get("http_api_id"), d.Get("route_id"), d.Get("environment_id"), d.Get("gateway_id")))
	return apigWaitForHttpApiDeployment(d, client, "Deployed", d.Timeout(schema.TimeoutCreate))
}

func apigHttpApiDeploymentRequest(d *schema.ResourceData) map[string]interface{} {
	// APIG's current HTTP-route deployment contract uses routeId together with
	// isInternalRoute=false. Do not use the deprecated httpApiConfig shape here.
	return map[string]interface{}{
		"routeId":         d.Get("route_id"),
		"isInternalRoute": false,
	}
}

func resourceAliCloudApigHttpApiDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	route, err := apigDescribeHttpApiRoute(client, d.Get("http_api_id").(string), d.Get("route_id").(string))
	if err != nil {
		if !d.IsNewResource() && (IsExpectedErrors(err, []string{"NotFound.RouteNotFound"}) || NotFoundError(err)) {
			log.Printf("[DEBUG] APIG HTTP API route %s was not found", d.Get("route_id"))
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	status, _ := route["deployStatus"].(string)
	if status == "NotDeployed" {
		d.SetId("")
		return nil
	}
	if !apigRouteMatchesDeployment(route, d.Get("environment_id").(string), d.Get("gateway_id").(string)) {
		return fmt.Errorf("APIG route %s deployment target does not match Terraform state", d.Get("route_id"))
	}
	return d.Set("status", status)
}

func apigDescribeHttpApiRoute(client *connectivity.AliyunClient, httpAPIID, routeID string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/v1/http-apis/%s/routes/%s", httpAPIID, routeID)
	response, err := client.RoaGet("APIG", "2024-03-27", action, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return apigResponseData(response)
}

func apigRouteMatchesDeployment(route map[string]interface{}, environmentID, gatewayID string) bool {
	environment, ok := route["environmentInfo"].(map[string]interface{})
	if !ok || environment["environmentId"] != environmentID {
		return false
	}
	gateway, ok := environment["gatewayInfo"].(map[string]interface{})
	return ok && gateway["gatewayId"] == gatewayID
}

func apigWaitForHttpApiDeployment(d *schema.ResourceData, client *connectivity.AliyunClient, target string, timeout time.Duration) error {
	err := resource.Retry(timeout, func() *resource.RetryError {
		route, err := apigDescribeHttpApiRoute(client, d.Get("http_api_id").(string), d.Get("route_id").(string))
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if !apigRouteMatchesDeployment(route, d.Get("environment_id").(string), d.Get("gateway_id").(string)) {
			return resource.NonRetryableError(fmt.Errorf("APIG route deployment target does not match the requested environment and gateway"))
		}
		status, _ := route["deployStatus"].(string)
		if status != target {
			return resource.RetryableError(fmt.Errorf("APIG route deployment status is %q, waiting for %q", status, target))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return d.Set("status", target)
}

func resourceAliCloudApigHttpApiDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/v1/http-apis/%s/undeploy", d.Get("http_api_id"))
	body := map[string]interface{}{
		"routeId":       d.Get("route_id"),
		"environmentId": d.Get("environment_id"),
		"gatewayId":     d.Get("gateway_id"),
	}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = client.RoaPost("APIG", "2024-03-27", action, nil, nil, body, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFound.RouteNotFound"}) || NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return apigWaitForHttpApiDeployment(d, client, "NotDeployed", d.Timeout(schema.TimeoutDelete))
}
