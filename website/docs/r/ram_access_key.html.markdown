---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_access_key"
description: |-
  Provides a Alicloud RAM Access Key resource.
---

# alicloud_ram_access_key

Provides a RAM Access Key resource.



For information about RAM Access Key and how to use it, see [What is Access Key](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ram-2015-05-01-createaccesskey).

-> **NOTE:** Available since v1.0.0.

-> **NOTE:**  You should set the `secret_file` if you want to get the access key.  

-> **NOTE:**  From version 1.98.0, if not set `pgp_key`, the resource will output the access key secret to field `secret` and please protect your backend state file judiciously


## Example Usage

Output the secret to a file.
<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_access_key&exampleId=215c87f4-930d-fb72-356c-fe0010b412841e32d6f9&activeTab=example&spm=docs.r.ram_access_key.0.215c87f493&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

# Create a new RAM access key for user.
resource "alicloud_ram_user" "user" {
  name         = "terraform-example-${random_integer.default.result}"
  display_name = "user_display_name"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
  force        = true
}

resource "alicloud_ram_access_key" "ak" {
  user_name   = alicloud_ram_user.user.name
  secret_file = "/xxx/xxx/xxx.txt"
}
```

Using `pgp_key` to encrypt the secret.
<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_access_key&exampleId=2a73ed24-9d33-6847-6269-0087c126dba0828f3649&activeTab=example&spm=docs.r.ram_access_key.1.2a73ed249d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

# Create a new RAM access key for user.
resource "alicloud_ram_user" "user" {
  name         = "terraform-example-${random_integer.default.result}"
  display_name = "user_display_name"
  mobile       = "86-18688888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
  force        = true
}

resource "alicloud_ram_access_key" "encrypt" {
  user_name = alicloud_ram_user.user.name
  pgp_key   = <<EOF
mQENBFXbjPUBCADjNjCUQwfxKL+RR2GA6pv/1K+zJZ8UWIF9S0lk7cVIEfJiprzzwiMwBS5cD0da
rGin1FHvIWOZxujA7oW0O2TUuatqI3aAYDTfRYurh6iKLC+VS+F7H+/mhfFvKmgr0Y5kDCF1j0T/
063QZ84IRGucR/X43IY7kAtmxGXH0dYOCzOe5UBX1fTn3mXGe2ImCDWBH7gOViynXmb6XNvXkP0f
sF5St9jhO7mbZU9EFkv9O3t3EaURfHopsCVDOlCkFCw5ArY+DUORHRzoMX0PnkyQb5OzibkChzpg
8hQssKeVGpuskTdz5Q7PtdW71jXd4fFVzoNH8fYwRpziD2xNvi6HABEBAAG0EFZhdWx0IFRlc3Qg
S2V5IDGJATgEEwECACIFAlXbjPUCGy8GCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJEOfLr44B
HbeTo+sH/i7bapIgPnZsJ81hmxPj4W12uvunksGJiC7d4hIHsG7kmJRTJfjECi+AuTGeDwBy84TD
cRaOB6e79fj65Fg6HgSahDUtKJbGxj/lWzmaBuTzlN3CEe8cMwIPqPT2kajJVdOyrvkyuFOdPFOE
A7bdCH0MqgIdM2SdF8t40k/ATfuD2K1ZmumJ508I3gF39jgTnPzD4C8quswrMQ3bzfvKC3klXRlB
C0yoArn+0QA3cf2B9T4zJ2qnvgotVbeK/b1OJRNj6Poeo+SsWNc/A5mw7lGScnDgL3yfwCm1gQXa
QKfOt5x+7GqhWDw10q+bJpJlI10FfzAnhMF9etSqSeURBRW5AQ0EVduM9QEIAL53hJ5bZJ7oEDCn
aY+SCzt9QsAfnFTAnZJQrvkvusJzrTQ088eUQmAjvxkfRqnv981fFwGnh2+I1Ktm698UAZS9Jt8y
jak9wWUICKQO5QUt5k8cHwldQXNXVXFa+TpQWQR5yW1a9okjh5o/3d4cBt1yZPUJJyLKY43Wvptb
6EuEsScO2DnRkh5wSMDQ7dTooddJCmaq3LTjOleRFQbu9ij386Do6jzK69mJU56TfdcydkxkWF5N
ZLGnED3lq+hQNbe+8UI5tD2oP/3r5tXKgMy1R/XPvR/zbfwvx4FAKFOP01awLq4P3d/2xOkMu4Lu
9p315E87DOleYwxk+FoTqXEAEQEAAYkCPgQYAQIACQUCVduM9QIbLgEpCRDny6+OAR23k8BdIAQZ
AQIABgUCVduM9QAKCRAID0JGyHtSGmqYB/4m4rJbbWa7dBJ8VqRU7ZKnNRDR9CVhEGipBmpDGRYu
lEimOPzLUX/ZXZmTZzgemeXLBaJJlWnopVUWuAsyjQuZAfdd8nHkGRHG0/DGum0l4sKTta3OPGHN
C1z1dAcQ1RCr9bTD3PxjLBczdGqhzw71trkQRBRdtPiUchltPMIyjUHqVJ0xmg0hPqFic0fICsr0
YwKoz3h9+QEcZHvsjSZjgydKvfLYcm+4DDMCCqcHuJrbXJKUWmJcXR0y/+HQONGrGJ5xWdO+6eJi
oPn2jVMnXCm4EKc7fcLFrz/LKmJ8seXhxjM3EdFtylBGCrx3xdK0f+JDNQaC/rhUb5V2XuX6VwoH
/AtY+XsKVYRfNIupLOUcf/srsm3IXT4SXWVomOc9hjGQiJ3rraIbADsc+6bCAr4XNZS7moViAAcI
PXFv3m3WfUlnG/om78UjQqyVACRZqqAGmuPq+TSkRUCpt9h+A39LQWkojHqyob3cyLgy6z9Q557O
9uK3lQozbw2gH9zC0RqnePl+rsWIUU/ga16fH6pWc1uJiEBt8UZGypQ/E56/343epmYAe0a87sHx
8iDV+dNtDVKfPRENiLOOc19MmS+phmUyrbHqI91c0pmysYcJZCD3a502X1gpjFbPZcRtiTmGnUKd
OIu60YPNE4+h7u2CfYyFPu3AlUaGNMBlvy6PEpU=
	  EOF
}

output "encrypted_secret" {
  value = alicloud_ram_access_key.encrypt.encrypted_secret
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ram_access_key&spm=docs.r.ram_access_key.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `status` - (Optional, Computed) The status of the AccessKey. Value:
  - Active: Activated.
  - Inactive: Disabled.
* `user_name` - (Optional, ForceNew) The RAM user name.
* `secret_file` - (Optional, ForceNew) The name of file that can save access key id and access key secret. Strongly suggest you to specified it when you creating access key, otherwise, you wouldn't get its secret ever.
* `pgp_key` - (Optional, ForceNew, Available since v1.47.0) Either a base-64 encoded PGP public key, or a keybase username in the form `keybase:some_person_that_exists`

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - (Available since v1.246.0) The create time of the AccessKey.
* `secret` - (Available since v1.98.0) The secret access key. Note that this will be written to the state file. 
If you use this, please protect your backend state file judiciously. 
Alternatively, you may supply a `pgp_key` instead, which will prevent the secret from being stored in plaintext, 
at the cost of preventing the use of the secret key in automation.
* `key_fingerprint` - (Available since v1.47.0) The fingerprint of the PGP key used to encrypt the secret
* `encrypted_secret` - (Available since v1.47.0) The encrypted secret, base64 encoded. ~> NOTE: The encrypted secret may be decrypted using the command line, for example: `terraform output encrypted_secret | base64 --decode | keybase pgp decrypt`.

## Timeouts

-> **NOTE:** Available since v1.246.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Key.
* `delete` - (Defaults to 5 mins) Used when delete the Access Key.
* `update` - (Defaults to 5 mins) Used when update the Access Key.