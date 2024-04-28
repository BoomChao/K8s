#!/bin/bash

# 下面脚本的参数通过pod启动的时候传入，一个是挂载的主机的目录，一个是主机的ip
set -e
host_root_dir=$1"/root"
linux_os_relese=$1"/etc/os-release"
host_ip=$2

# 检查系统配置
echo "root dir:" $host_root_dir  " host ip:" $host_ip  " linux release:" $linux_os_relese
if [ ! -f ${linux_os_relese} ]; then
  echo "unsupported os"
  exit 255
fi

if [ ! -d "$host_root_dir" ]; then
  echo "root dir $host_root_dir missing"
  exit 255
fi

# 检查文件是否存在
authorized_keys=$host_root_dir/.ssh/authorized_keys
if [ ! -f ${authorized_keys} ]; then
  mkdir -p $host_root_dir/.ssh/
  touch $authorized_keys
fi

# 检查公钥是否已经存在
echo "update host public key"
empty=false
grep "^root@zxc" $authorized_keys || empty=true
if [ $empty = true ]; then
  echo "write pubkey in authorized_keys"
  cat /root/.ssh/id_rsa.pub >> $authorized_keys
fi

sshd_config=$1/etc/ssh/sshd_config
echo "check sshd config"

# 是否允许公钥登录
pubkey_disabled=false
grep -E "^PubkeyAuthentication[[:blank:]]+yes" $sshd_config || pubkey_disabled=true
grep "^PubkeyAuthentication " $sshd_config || pubkey_disabled=false
if [ $pubkey_disabled = true ]; then
    echo "PubkeyAuthentication disabled"
    exit 255
fi

# 是否允许root登录
permit_disabled=false
grep -E "^PermitRootLogin[[:blank:]]+yes" $sshd_config || grep -E "^PermitRootLogin[[:blank:]]+prohibit-password" $sshd_config || permit_disabled=true
grep "^PermitRootLogin " $sshd_config || permit_disabled=false
if [ $permit_disabled = true ]; then
    echo "PermitRootLogin disabled"
    exit 255
fi

echo "ssh and excute sh remotely"
ssh root@${host_ip} 'bash -s' < /usr/local/bin/cert.sh