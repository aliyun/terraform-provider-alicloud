#!/usr/bin/env bash

set -e

: ${GITHUB_TOKEN:=?}
: ${GPG_FINGERPRINT=?}

export GITHUB_TOKEN=${GITHUB_TOKEN}
export GPG_FINGERPRINT=${GPG_FINGERPRINT}
# https://www.terraform.io/docs/registry/providers/publishing.html

CURRENT_PATH=$(pwd)

go version

cd $GOPATH
mkdir -p src/github.com/aliyun
cp -rf $CURRENT_PATH/terraform-provider-alicloud src/github.com/aliyun/
# 每次下载 goreleaser 太慢了，所以采用本地编译
mkdir -p src/github.com/goreleaser
cp -rf $CURRENT_PATH/goreleaser src/github.com/goreleaser/
pushd src/github.com/goreleaser/goreleaser
go mod tidy
go mod vendor
go build -o /usr/bin/goreleaser .
goreleaser --version
popd

echo "Building goreleaser finished."

pushd src/github.com/aliyun/terraform-provider-alicloud

# update the changelog
# 更新 Changelog 总共有一部分：
# 1. 更新 PR 的URL
# 2. 更新版本号。当前版本号的计算原则是：具体的做法是，查找上一个版本号，根据然后计算当前版本号和下个待发布的版本号。
#    如果有新的resource或者datasource，发布大版本，否则发布小版本。最后再把日期更改了。
bash scripts/changelog-links.sh

changelog="CHANGELOG.md"
unreleased_line=`grep -n "(Unreleased)" $changelog | head -1 | cut -d ":" -f 1`
last_version_line=`grep -n "##" $changelog | head -2 | tail -1 | cut -d ":" -f 1`
new_feature_line=`grep -n "\*\*New Resource" $changelog | head -1 | head -$last_version_line | cut -d ":" -f 1`
if [[ ${new_feature_line} -gt ${last_version_line} ]]; then
  new_feature_line=`grep -n "\*\*Data" $changelog | head -1 | head -$last_version_line | cut -d ":" -f 1`
fi
echo "new_feature_line: ${new_feature_line} ang last version line: ${last_version_line}"

last_version=`grep -n "##" $changelog | head -2 | tail -1 | cut -d " " -f 2`
this_version=$last_version
next_version=$this_version
if [[ ${new_feature_line} -gt ${last_version_line} ]]; then
  arr=(${this_version//./ })
  this_version="$((${arr[0]})).$((${arr[1]})).$((${arr[2]} + 1))"
  next_version="$((${arr[0]})).$((${arr[1]})).$((${arr[2]} + 2))"
else
  arr=(${this_version//./ })
  this_version="$((${arr[0]})).$((${arr[1]} + 1)).0"
  next_version="$((${arr[0]})).$((${arr[1]} + 2)).0"
fi

echo "The last version is ${last_version}; this version is ${this_version}; next verison is ${next_version}"

bump_date=`env LANG=en_US.UTF-8 date '+%B %d, %Y'`
sed -i "/(Unreleased)/d" $changelog
sed -i -e "${unreleased_line}i \#\# ${this_version} ($bump_date)" $changelog
sed -i -e "${unreleased_line}i \#\# ${next_version} (Unreleased)" $changelog

git diff | cat
git add .

git config --global user.email guimin.hgm@alibaba-inc.com
git config --global user.name xiaozhu36
git commit -m "Cleanup after release $this_version"

git tag v$this_version

goreleaser release --rm-dist

echo "Building providers finished."

ls -l
rm -rf dist

echo "Sync the changelog to the output repo"
cp $changelog $CURRENT_PATH/terraform-provider-alicloud
popd

cd $CURRENT_PATH
pushd terraform-provider-alicloud
git diff | cat
git add .

git config --global user.email guimin.hgm@alibaba-inc.com
git config --global user.name xiaozhu36
git commit -m "Cleanup after release $this_version"
popd
}