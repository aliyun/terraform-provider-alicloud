#!/usr/bin/env bash

set -e -o pipefail

my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"

source ${release_dir}/ci/tasks/utils.sh

: ${terraform_provider_bucket_name:?}
: ${terraform_provider_bucket_region:?}
: ${terraform_provider_access_key:?}
: ${terraform_provider_secret_key:?}

ls -la
provider="terraform-provider-alicloud"
echo "tar ${provider} ..."
tar -czvf ${provider}.tgz ${provider}

echo -e "Uploading ${provider}.tgz ..."
aliyun oss cp ${provider}.tgz oss://${terraform_provider_bucket_name}/${provider}.tgz -f --access-key-id ${terraform_provider_access_key} --access-key-secret ${terraform_provider_secret_key} --region ${terraform_provider_bucket_region}

echo -e "Upload Finished!"
ls -l ${output_path}
rm -rf ${provider}.tgz