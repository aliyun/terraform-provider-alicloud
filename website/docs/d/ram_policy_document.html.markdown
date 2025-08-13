---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_policy_document"
description: |-
  Generates a RAM policy document to the user.
---

# alicloud_ram_policy_document

This data source Generates a RAM policy document of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.184.0.

## Example Usage

### Basic Example

```terraform
data "alicloud_ram_policy_document" "basic_example" {
  version = "1"
  statement {
    effect   = "Allow"
    action   = ["oss:*"]
    resource = ["acs:oss:*:*:myphotos", "acs:oss:*:*:myphotos/*"]
  }
}

resource "alicloud_ram_policy" "default" {
  policy_name     = "tf-example"
  policy_document = data.alicloud_ram_policy_document.basic_example.document
  force           = true
}
```

`data.alicloud_ram_policy_document.basic_example.document` will evaluate to:

```json
{
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "oss:*",
      "Resource": [
        "acs:oss:*:*:myphotos",
        "acs:oss:*:*:myphotos/*"
      ]
    }
  ],
  "Version": "1"
}
```

### Example Multiple Condition Keys and Values

```terraform
data "alicloud_ram_policy_document" "multiple_condition" {
  version = "1"
  statement {
    effect   = "Allow"
    action   = ["oss:ListBuckets", "oss:GetBucketStat", "oss:GetBucketInfo", "oss:GetBucketTagging", "oss:GetBucketAcl"]
    resource = ["acs:oss:*:*:*"]
  }
  statement {
    effect   = "Allow"
    action   = ["oss:GetObject", "oss:GetObjectAcl"]
    resource = ["acs:oss:*:*:myphotos/hangzhou/2015/*"]
  }
  statement {
    effect   = "Allow"
    action   = ["oss:ListObjects"]
    resource = ["acs:oss:*:*:myphotos"]
    condition {
      operator = "StringLike"
      variable = "oss:Delimiter"
      values   = ["/"]
    }
    condition {
      operator = "StringLike"
      variable = "oss:Prefix"
      values   = ["", "hangzhou/", "hangzhou/2015/*"]
    }
  }
}

resource "alicloud_ram_policy" "policy" {
  policy_name     = "tf-example-condition"
  policy_document = data.alicloud_ram_policy_document.multiple_condition.document
  force           = true
}
```

`data.alicloud_ram_policy_document.multiple_condition.document` will evaluate to:

```json
{
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "oss:ListBuckets",
        "oss:GetBucketStat",
        "oss:GetBucketInfo",
        "oss:GetBucketTagging",
        "oss:GetBucketAcl"
      ],
      "Resource": "acs:oss:*:*:*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "oss:GetObject",
        "oss:GetObjectAcl"
      ],
      "Resource": "acs:oss:*:*:myphotos/hangzhou/2015/*"
    },
    {
      "Effect": "Allow",
      "Action": "oss:ListObjects",
      "Resource": "acs:oss:*:*:myphotos",
      "Condition": {
        "StringLike": {
          "oss:Delimiter": "/",
          "oss:Prefix": [
            "",
            "hangzhou/",
            "hangzhou/2015/*"
          ]
        }
      }
    }
  ],
  "Version": "1"
}
```

### Example Assume-Role Policy with RAM Principal

```terraform
data "alicloud_ram_policy_document" "ram_example" {
  statement {
    effect = "Allow"
    action = ["sts:AssumeRole"]
    principal {
      entity      = "RAM"
      identifiers = ["acs:ram::123456789012****:root"]
    }
  }
}

resource "alicloud_ram_role" "role" {
  name     = "tf-example-role-ram"
  document = data.alicloud_ram_policy_document.ram_example.document
  force    = true
}
```

`data.alicloud_ram_policy_document.ram_example.document` will evaluate to:

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

### Example Assume-Role Policy with Service Principal

```terraform
data "alicloud_ram_policy_document" "service_example" {
  statement {
    effect = "Allow"
    action = ["sts:AssumeRole"]
    principal {
      entity      = "Service"
      identifiers = ["ecs.aliyuncs.com"]
    }
  }
}

resource "alicloud_ram_role" "role" {
  name     = "tf-example-role-service"
  document = data.alicloud_ram_policy_document.service_example.document
  force    = true
}
```

`data.alicloud_ram_policy_document.service_example.document` will evaluate to:

```json
{
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": [
          "ecs.aliyuncs.com"
        ]
      }
    }
  ],
  "Version": "1"
}
```

### Example Assume-Role Policy with Federated Principal

```terraform
data "alicloud_ram_policy_document" "federated_example" {
  statement {
    effect = "Allow"
    action = ["sts:AssumeRole"]
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
  name     = "tf-example-role-federated"
  document = data.alicloud_ram_policy_document.federated_example.document
  force    = true
}
```

`data.alicloud_ram_policy_document.federated_example.document` will evaluate to:

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

* `version` - (Optional) Version of the RAM policy document. Valid value is `1`. Default value is `1`.
* `statement` - (Optional) Statement of the RAM policy document. See the following `Block statement`. See [`statement`](#statement) below.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

### `statement`

The statement supports the following:

* `effect` - (Optional) This parameter indicates whether or not the `action` is allowed. Valid values are `Allow` and `Deny`. Default value is `Allow`. If you want to create a RAM role policy document, it must be `Allow`.
* `action` - (Required) Action of the RAM policy document. If you want to create a RAM role policy document, it must be `["sts:AssumeRole"]`.
* `resource` - (Optional) List of specific objects which will be authorized. If you want to create a RAM policy document, it must be set.
* `principal` - (Optional) Principal of the RAM policy document. If you want to create a RAM role policy document, it must be set. See [`principal`](#statement-principal) below.
* `condition` - (Optional) Specifies the condition that are required for a policy to take effect. See [`condition`](#statement-condition) below.

### `statement-principal`

The principal supports the following:

* `entity` - (Required) The trusted entity. Valid values: `RAM`, `Service` and `Federated`.
* `identifiers` - (Required) The identifiers of the principal.

### `statement-condition`

The condition supports the following:

* `operator` - (Required) The operator of the condition.
* `variable` - (Required) The variable of the condition.
* `values` - (Required) The values of the condition.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `document` - Standard policy document rendered based on the arguments above.
