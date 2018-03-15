#!/bin/bash

set -e

# Prepare environment for quota support
apt-get update
apt-get install -y quota
truncate -s1G /tmp/test.ext4
mkfs.ext4 -F /tmp/test.ext4
mkdir -p /mnt/quota_test
mount -o usrjquota=aquota.user,grpjquota=aquota.group,jqfmt=vfsv1 /tmp/test.ext4 /mnt/quota_test
quotacheck -vucm /mnt/quota_test
quotacheck -vugm /mnt/quota_test
quotaon -v /mnt/quota_test
