#!/bin/bash
terraform init
default_provider_root_dir='.terraform/providers/registry.terraform.io/hashicorp/alicloud/'
provide_version=$(ls $default_provider_root_dir)

default_provider_dir="${default_provider_root_dir}${provide_version}/"
platform_path=$(ls "${default_provider_dir}")
platform=$(echo "${platform_path}" | tr '_' '-')

local_provider_tgz=$(ls ../../bin/ | grep "$platform")
filename=$(ls "${default_provider_dir}${platform_path}")


provider_root_dir='/usr/local/lib/terraform/registry.terraform.io/hashicorp/alicloud/'
cat <<EOT > ~/.terraformrc
provider_installation {
  filesystem_mirror {
    path = "/usr/local/lib/terraform"
  }
}
EOT
tar -zxvf "../../bin/${local_provider_tgz}" \
  && mv -f bin/terraform-provider-alicloud "${provider_root_dir}${provide_version}/${platform_path}/${filename}" \
  && rm -rf bin/

rm -rf .terraform.lock.hcl && terraform init
source env.sh && echo "ots account has set to env."


