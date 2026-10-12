---
subcategory: "Container Registry"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_authorization_token"
sidebar_current: "docs-alicloud-ephemeral-cr-ee-authorization-token"
description: |-
    Fetches a temporary Container Registry Enterprise Edition login credential without persisting it to the state or plan.
---

# Ephemeral: alicloud_cr_ee_authorization_token

This ephemeral resource fetches a temporary login credential (username and password) for a Container Registry Enterprise Edition instance through the `GetAuthorizationToken` API without persisting the credential to the state or plan. The credential is available for the duration of a single Terraform operation and can be referenced in provider configuration, locals, provisioners, and other ephemeral resources, but never appears in any Terraform artifact.

-> **NOTE:** Available since v2.0.0-beta5. Ephemeral resources require Terraform 1.10 or later.

-> **NOTE:** The token is valid for one hour by default. The validity period can be extended to at most 24 hours through `expires_in_hours`; when the provider is configured with STS, the validity period is capped at the remaining validity of the STS token. The token is re-fetched on every operation that opens it and expires at the reported `expire_time`.

-> **NOTE:** The ephemeral resource is opened during the plan phase when `instance_id` is known at plan time, so the instance must already exist. When `instance_id` references another resource in the same configuration, its value is unknown during plan and the open is deferred to the apply phase automatically.

-> **NOTE:** The [`alicloud_cr_ee_instances`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cr_ee_instances) data source can also expose the authorization token and temporary username of an instance through its `enable_details` argument, persisting them into the state; use this ephemeral resource when the credential must never appear in any Terraform artifact.

## Example Usage

Fetch a short-lived credential for logging on to the Registry of an instance:

```terraform
ephemeral "alicloud_cr_ee_authorization_token" "example" {
  instance_id = "cri-xxxxxxxxxxxxxxxx"
}
```

Fetch a credential with a longer validity period:

```terraform
ephemeral "alicloud_cr_ee_authorization_token" "example" {
  instance_id      = "cri-xxxxxxxxxxxxxxxx"
  expires_in_hours = 24
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the Container Registry Enterprise Edition instance.
* `expires_in_hours` - (Optional) The validity period of the token, in hours. Valid values are `1` to `24`. When this field is omitted, the token is valid for one hour. When the provider is configured with STS, the validity period is capped at the remaining validity of the STS token.

## Attributes Reference

The following attributes are exported in the ephemeral result:

* `temp_username` - The temporary username used to log on to the Registry of the instance.
* `authorization_token` - (Sensitive) The temporary password used to log on to the Registry of the instance.
* `expire_time` - The expiration time of the token, in milliseconds since the Unix epoch.
