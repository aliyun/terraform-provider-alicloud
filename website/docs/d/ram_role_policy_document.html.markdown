---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_role_policy_document"
sidebar_current: "docs-alicloud-datasource-ram-role-policy-document"
description: |-
    Generates a RAM role policy document to the user.
---

# alicloud\_ram\_role\_policy\_document

This data source Generates a RAM role policy document of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.183.0+.

## Example Usage

### Basic Example

```terraform
data "alicloud_ram_role_policy_document" "basic_example" {
  statement {
    effect = "Allow"
    action = "sts:AssumeRole"
    principal {
      entity      = "RAM"
      identifiers = ["acs:ram::123456789012****:root"]
    }
  }
}

resource "alicloud_ram_role" "role" {
  name     = "test-role-ram"
  document = data.alicloud_ram_role_policy_document.basic_example.document
  force    = true
}
```

`data.alicloud_ram_role_policy_document.basic_example.document` will evaluate to:

```json
{
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "RAM": [
          "acs:ram::123456789012****:root"
        ]
      }
    }
  ],
  "Version": "1"
}
```

### Example Condition Keys and Values

```terraform
data "alicloud_ram_role_policy_document" "multiple_condition" {
  statement {
    effect = "Allow"
    action = "sts:AssumeRole"
    principal {
      entity      = "Federated"
      identifiers = ["acs:ram::123456789012****:saml-provider/testprovider"]
    }
    condition {
      operator = "StringEquals"
      variable = "saml:recipient"
      values   = ["https://signin.aliyun.com/saml-role/sso"]
    }
  }
}

resource "alicloud_ram_role" "role" {
  name     = "test-role-federated"
  document = data.alicloud_ram_role_policy_document.multiple_condition.document
  force    = true
}
```

`data.alicloud_ram_role_policy_document.multiple_condition.document` will evaluate to:

```json
{
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "Federated": [
          "acs:ram::123456789012****:saml-provider/testprovider"
        ]
      },
      "Condition": {
        "StringEquals": {
          "saml:recipient": "https://signin.aliyun.com/saml-role/sso"
        }
      }
    }
  ],
  "Version": "1"
}
```

## Argument Reference

The following arguments are supported:

* `version` - (Optional) Version of the RAM role policy document. Valid value is `1`. Default value is `1`.
* `statement` - (Optional) Statement of the RAM role policy document. See the following `Block statement`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

#### Block statement

The statement supports the following:

* `effect` - (Optional) This parameter indicates whether or not the `action` is allowed. Valid values: `Allow`. Default value is `Allow`.
* `action` - (Optional) Action of the RAM role policy document. Valid values: `sts:AssumeRole`. Default value is `sts:AssumeRole`.
* `principal` - (Optional) Principal of the RAM role policy document. See the following `Block principal`.
* `condition` - (Optional) Specifies the condition that are required for a policy to take effect. See the following `Block condition`.

#### Block principal

The principal supports the following:

* `entity` - (Required) The trusted entity . Valid values: `RAM`, `Service` and `Federated`.
* `identifiers` - (Required) The identifiers of the principal.

#### Block condition

The condition supports the following:

* `operator` - (Required) The operator of the condition.
* `variable` - (Required) The variable of the condition.
* `values` - (Required) The values of the condition.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `document` - Standard RAM role policy document rendered based on the arguments above.