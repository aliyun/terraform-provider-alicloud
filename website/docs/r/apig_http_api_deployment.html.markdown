---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_http_api_deployment"
description: |-
  Provides an Alicloud APIG HTTP API route deployment resource.
---

# alicloud_apig_http_api_deployment

Deploys one APIG HTTP API route to an environment and gateway. Destroying the resource undeploys the route but does not delete the route or HTTP API.

For information about route deployment, see [DeployHttpApi](https://next.api.alibabacloud.com/document/APIG/2024-03-27/DeployHttpApi).

-> **NOTE:** Available since v1.291.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_apig_http_api_deployment&exampleId=55c0db89-65c4-c22d-1c3c-031a3035ceb04704350e&activeTab=example&spm=docs.r.apig_http_api_deployment.0.55c0db8965&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_apig_http_api_deployment" "default" {
  http_api_id    = alicloud_apig_http_api.default.id
  route_id       = alicloud_apig_route.default.route_id
  environment_id = alicloud_apig_gateway.default.environments[0].environment_id
  gateway_id     = alicloud_apig_gateway.default.id
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_apig_http_api_deployment&spm=docs.r.apig_http_api_deployment.example&intl_lang=EN_US)


## Argument Reference

* `environment_id` - (Required, ForceNew) The deployment environment ID.
* `gateway_id` - (Required, ForceNew) The deployment gateway ID.
* `http_api_id` - (Required, ForceNew) The HTTP API ID.
* `route_id` - (Required, ForceNew) The HTTP API route ID.

## Attributes Reference

* `id` - A composite ID containing the HTTP API, route, environment, and gateway IDs.
* `status` - The observed route deployment status.

## Timeouts

* `create` - (Defaults to 5 mins) Used when deploying and waiting for `Deployed`.
* `delete` - (Defaults to 5 mins) Used when undeploying and waiting for `NotDeployed`.

## Import

Existing deployed routes can be imported without redeployment:

```shell
$ terraform import alicloud_apig_http_api_deployment.example '<http_api_id>:<route_id>:<environment_id>:<gateway_id>'
```
