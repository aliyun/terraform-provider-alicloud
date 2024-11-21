---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_custom_domain"
description: |-
  Provides a Alicloud FCV3 Custom Domain resource.
---

# alicloud_fcv3_custom_domain

Provides a FCV3 Custom Domain resource.

Custom Domain names allow users to access FC functions through custom domain names, providing convenience for building Web services using function compute.
You can bind a custom domain name to Function Compute and set different paths to different functions of different services.

For information about FCV3 Custom Domain and how to use it, see [What is Custom Domain](https://www.alibabacloud.com/help/en/functioncompute/developer-reference/api-fc-2023-03-30-getcustomdomain).

~> **NOTE:** This content is a technical preview, and should not be relied on in a production environment.

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fcv3_custom_domain&exampleId=d8ff329b-9334-7860-59b2-fd3b11bdfd0abc9c8568&activeTab=example&spm=docs.r.fcv3_custom_domain.0.d8ff329b93&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

variable "name" {
  default = "flask-07ap.fcv3.1511928242963727.cn-shanghai.fc.devsapp.net"
}

variable "function_name1" {
  default = "terraform-custom-domain-t1"
}

variable "auth_config" {
  default = <<EOF
{
    "jwks": {
        "keys": [
            {
                "p": "8AdUVeldoE4LueFuzEF_C8tvJ7NhlkzS58Gz9KJTPXPr5DADSUVLWJCr5OdFE79q513SneT0UhGo-JfQ1lNMoNv5-YZ1AxIo9fZUEPIe-KyX9ttaglpzCAUE3TeKdm5J-_HZQzBPKbyUwJHAILNgB2-4IBZZwK7LAfbmfi9TmFM",
                "kty": "RSA",
                "q": "x8m5ydXwC8AAp9I-hOnUAx6yQJz1Nx-jXPCfn--XdHpJuNcuwRQsuUCSRQs_h3SoCI3qZZdzswQnPrtHFxgUJtQFuMj-QZpyMnebDb81rmczl2KPVUtaVDVagJEF6U9Ov3PfrLhvHUEv5u7p6s4Z6maBUaByfFlhEVPv4_ao8Us",
                "d": "bjIQAKD2e65gwJ38_Sqq_EmLFuMMey3gjDv1bSCHFH8fyONJTq-utrZfvspz6EegRwW2mSHW9kq87hRwIBW9y7ED5N4KG5gHDjyh57BRE0SKv0Dz1igtKLyp-nl8-aHc1DbONwr1d7tZfFt255TxIN8cPTakXOp2Av_ztql_JotVUGK8eHmXNJFlvq5tc180sKWMHNSNsCUhQgcB1TWb_gwcqxdsIWPsLZI491XKeTGQ98J7z5h6R1cTC97lfJZ0vNtJahd2jHd3WfTUDj5-untMKyZpYYak2Vr8xtFz8H6Q5Rsz8uX_7gtEqYH2CMjPdbXcebrnD1igRSJMYiP0lQ",
                "e": "AQAB",
                "use": "sig",
                "qi": "MTCCRu8AcvvjbHms7V_sDFO7wX0YNyvOJAAbuTmHvQbJ0NDeDta-f-hi8cjkMk7Fpk2hej158E5gDyO62UG99wHZSbmHT34MvIdmhQ5mnbL-5KK9rxde0nayO1ebGepD_GJThPAg9iskzeWpCg5X2etNo2bHoG_ZLQGXj2BQ1VM",
                "dp": "J4_ttKNcTTnP8PlZO81n1VfYoGCOqylKceyZbq76rVxX-yp2wDLtslFWI8qCtjiMtEnglynPo19JzH-pakocjT70us4Qp0rs-W16ebiOpko8WfHZvzaNUzsQjC3FYrPW-fHo74wc4DI3Cm57jmhCYbdmT9OfQ4UL7Oz3HMFMNAU",
                "alg": "RS256",
                "dq": "H4-VgvYB-sk1EU3cRIDv1iJWRHDHKBMeaoM0pD5kLalX1hRgNW4rdoRl1vRk79AU720D11Kqm2APlxBctaA_JrcdxEg0KkbsvV45p11KbKeu9b5DKFVECsN27ZJ7XZUCuqnibtWf7_4pRBD_8PDoFShmS2_ORiiUdflNjzSbEas",
                "n": "u1LWgoomekdOMfB1lEe96OHehd4XRNCbZRm96RqwOYTTc28Sc_U5wKV2umDzolfoI682ct2BNnRRahYgZPhbOCzHYM6i8sRXjz9Ghx3QHw9zrYACtArwQxrTFiejbfzDPGdPrMQg7T8wjtLtkSyDmCzeXpbIdwmxuLyt_ahLfHelr94kEksMDa42V4Fi5bMW4cCLjlEKzBEHGmFdT8UbLPCvpgsM84JK63e5ifdeI9NdadbC8ZMiR--dFCujT7AgRRyMzxgdn2l-nZJ2ZaYzbLUtAW5_U2kfRVkDNa8d1g__2V5zjU6nfLJ1S2MoXMgRgDPeHpEehZVu2kNaSFvDUQ"
            }
        ]
    },
    "tokenLookup": "header:auth",
    "claimPassBy": "header:name:name"
}
EOF
}

variable "certificate" {
  default = <<EOF
-----BEGIN CERTIFICATE-----
MIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV
BAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP
MA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0
ZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow
djELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE
ChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG
9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ
AoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB
coG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook
KOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw
HQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy
+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC
QkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN
MAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ
AJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT
cQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1
Ofi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd
DUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=
-----END CERTIFICATE-----
EOF
}

variable "private_key" {
  default = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV
kg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM
ywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB
AoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd
6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP
hwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4
MdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz
71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm
Ev9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE
qygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8
9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM
zWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe
DrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=
-----END RSA PRIVATE KEY-----
EOF
}

resource "alicloud_fcv3_custom_domain" "default" {
  custom_domain_name = var.name
  route_config {
    routes {
      function_name = var.function_name1
      rewrite_config {
        regex_rules {
          match       = "/api/*"
          replacement = "$1"
        }
        regex_rules {
          match       = "/api1/*"
          replacement = "$1"
        }
        regex_rules {
          match       = "/api2/*"
          replacement = "$1"
        }

        wildcard_rules {
          match       = "^/api1/.+?/(.*)"
          replacement = "/api/v1/$1"
        }
        wildcard_rules {
          match       = "^/api2/.+?/(.*)"
          replacement = "/api/v2/$1"
        }
        wildcard_rules {
          match       = "^/api2/.+?/(.*)"
          replacement = "/api/v3/$1"
        }

        equal_rules {
          match       = "/old"
          replacement = "/new"
        }
        equal_rules {
          replacement = "/new1"
          match       = "/old1"
        }
        equal_rules {
          match       = "/old2"
          replacement = "/new2"
        }

      }

      methods = [
        "GET",
        "POST",
        "DELETE",
        "HEAD"
      ]
      path      = "/a"
      qualifier = "LATEST"
    }
    routes {
      function_name = var.function_name1
      methods = [
        "GET"
      ]
      path      = "/b"
      qualifier = "LATEST"
    }
    routes {
      function_name = var.function_name1
      methods = [
        "POST"
      ]
      path      = "/c"
      qualifier = "1"
    }

  }

  auth_config {
    auth_type = "jwt"
    auth_info = var.auth_config
  }

  protocol = "HTTP,HTTPS"
  cert_config {
    cert_name   = "cert-name"
    certificate = var.certificate
    private_key = var.private_key
  }

  tls_config {
    cipher_suites = [
      "TLS_RSA_WITH_AES_128_CBC_SHA",
      "TLS_RSA_WITH_AES_256_CBC_SHA",
      "TLS_RSA_WITH_AES_128_GCM_SHA256",
      "TLS_RSA_WITH_AES_256_GCM_SHA384"
    ]
    max_version = "TLSv1.3"
    min_version = "TLSv1.0"
  }

  waf_config {
    enable_waf = "false"
  }

}
```

## Argument Reference

The following arguments are supported:
* `auth_config` - (Optional, List) Permission authentication configuration See [`auth_config`](#auth_config) below.
* `cert_config` - (Optional, Computed, List) HTTPS certificate information See [`cert_config`](#cert_config) below.
* `custom_domain_name` - (Optional, ForceNew, Computed) The name of the resource
* `protocol` - (Optional) The protocol type supported by the domain name. HTTP: only HTTP protocol is supported. HTTPS: only HTTPS is supported. HTTP,HTTPS: Supports HTTP and HTTPS protocols.
* `route_config` - (Optional, List) Route matching rule configuration See [`route_config`](#route_config) below.
* `tls_config` - (Optional, Computed, List) TLS configuration information See [`tls_config`](#tls_config) below.
* `waf_config` - (Optional, List) Web application firewall configuration information See [`waf_config`](#waf_config) below.

### `auth_config`

The auth_config supports the following:
* `auth_info` - (Optional) Authentication Information
* `auth_type` - (Optional) Authentication type. anonymous, function, or jwt.

### `cert_config`

The cert_config supports the following:
* `cert_name` - (Optional) Certificate Name
* `certificate` - (Optional) PEM format certificate
* `private_key` - (Optional) Private Key in PEM format

### `route_config`

The route_config supports the following:
* `routes` - (Optional, List) Routing Configuration List See [`routes`](#route_config-routes) below.

### `route_config-routes`

The route_config-routes supports the following:
* `function_name` - (Optional) Function name
* `methods` - (Optional, List) List of supported HTTP methods
* `path` - (Optional) Route matching rule
* `qualifier` - (Optional) Version or Alias
* `rewrite_config` - (Optional, List) Override Configuration See [`rewrite_config`](#route_config-routes-rewrite_config) below.

### `route_config-routes-rewrite_config`

The route_config-routes-rewrite_config supports the following:
* `equal_rules` - (Optional, List) Exact Match Rule List See [`equal_rules`](#route_config-routes-rewrite_config-equal_rules) below.
* `regex_rules` - (Optional, List) Regular match rule list See [`regex_rules`](#route_config-routes-rewrite_config-regex_rules) below.
* `wildcard_rules` - (Optional, List) List of wildcard matching rules See [`wildcard_rules`](#route_config-routes-rewrite_config-wildcard_rules) below.

### `route_config-routes-rewrite_config-equal_rules`

The route_config-routes-rewrite_config-equal_rules supports the following:
* `match` - (Optional) Matching Rules
* `replacement` - (Optional) Replace Rules

### `route_config-routes-rewrite_config-regex_rules`

The route_config-routes-rewrite_config-regex_rules supports the following:
* `match` - (Optional) Matching Rules
* `replacement` - (Optional) Replace Rules

### `route_config-routes-rewrite_config-wildcard_rules`

The route_config-routes-rewrite_config-wildcard_rules supports the following:
* `match` - (Optional) Matching Rules
* `replacement` - (Optional) Replace Rules

### `tls_config`

The tls_config supports the following:
* `cipher_suites` - (Optional, List) List of TLS cipher suites
* `max_version` - (Optional) The maximum version of TLS. Enumeration values: TLSv1.3, TLSv1.2, TLSv1.1, TLSv1.0
* `min_version` - (Optional) TLS minimum version number. Enumeration values: TLSv1.3, TLSv1.2, TLSv1.1, TLSv1.0

### `waf_config`

The waf_config supports the following:
* `enable_waf` - (Optional) Enable WAF protection

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `account_id` - (Available since v1.234.0) The ID of your Alibaba Cloud account (primary account).
* `api_version` - (Available since v1.234.0) API version of Function Compute
* `create_time` - The creation time of the resource
* `last_modified_time` - (Available since v1.234.0) The last time the custom domain name was Updated
* `subdomain_count` - (Available since v1.234.0) Number of subdomains

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Custom Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Domain.
* `update` - (Defaults to 5 mins) Used when update the Custom Domain.

## Import

FCV3 Custom Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_custom_domain.example <id>
```