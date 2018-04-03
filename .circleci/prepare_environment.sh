#!/bin/bash

set -e

# Prepare environment for quota support
apt-get update
apt-get install -y quota
truncate -s1G /tmp/quotas_enabled.ext4
truncate -s1G /tmp/quotas_disabled.ext4
mkfs.ext4 -F /tmp/quotas_enabled.ext4
mkfs.ext4 -F /tmp/quotas_disabled.ext4
mkdir -p /mnt/quota_test /mnt/noquota_test
mount -o usrjquota=aquota.user,grpjquota=aquota.group,jqfmt=vfsv1 /tmp/quotas_enabled.ext4 /mnt/quota_test
mount /tmp/quotas_disabled.ext4 /mnt/noquota_test
quotacheck -vucm /mnt/quota_test
quotacheck -vugm /mnt/quota_test
quotaon -v /mnt/quota_test
for i in {10000..10009}
do
    addgroup --gid $i test$i
    adduser --system --shell /bin/false --no-create-home --uid $i --gid $i --disabled-login --disabled-password test$i
done
