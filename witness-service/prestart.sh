#!/bin/bash

set -e

function parse_yaml {
   local prefix=$2
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
   sed -ne "s|^\($s\):|\1|" \
        -e "s|^\($s\)\($w\)$s:$s[\"']\(.*\)[\"']$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
   awk -F$fs '{
      indent = length($1)/2;
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
      if (length($3) > 0) {
         vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
         printf("%s%s%s=%s\n", "'$prefix'",vn, $2, $3);
      }
   }'
}

yml_cfg=$1
uri=$(parse_yaml $yml_cfg | grep "system__dbURI" | awk -F "=" '{print $2}' )
user=$(echo $uri | awk -F ":" '{print $1}')
passwd=$(echo $uri | awk -F ":" '{print $2}' | awk -F "@" '{print $1}')
host=$(echo $uri | awk -F ":" '{print $2}' | awk -F "@" '{print $2}' | awk -F "/" '{print $1}')
if [ -z $host ]
then
    host='localhost'
fi

mysql -h$host -u$user -p$passwd << EOF
CREATE DATABASE IF NOT EXISTS tubedb;
USE tubedb;
CREATE TABLE IF NOT EXISTS request (
  token varchar(45) NOT NULL,
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  uuid varchar(45) NOT NULL,
  fromAssetType varchar(45) NOT NULL DEFAULT '',
  toAssetType varchar(45) NOT NULL,
  tokenAmount int(10) unsigned NOT NULL,
  toAddr varchar(45) NOT NULL,
  fromAddr varchar(45) NOT NULL,
  redeemCode varchar(45) NOT NULL DEFAULT '',
  status varchar(45) NOT NULL DEFAULT 'requested',
  createdAt datetime DEFAULT NULL,
  fullfilledAt datetime DEFAULT NULL,
  finalizedAt datetime DEFAULT NULL,
  finalizedOn varchar(45) NOT NULL DEFAULT '',
  txHash varchar(128) DEFAULT '',
  PRIMARY KEY (id, token)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8;
EOF
