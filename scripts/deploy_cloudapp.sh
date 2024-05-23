#!/bin/bash

# Description: Deployment of cloud app in GC Management. 

echo "Deploy AWS cloudapp"

gc-mgmtctl onboard_cloud_app --path_to_manifest  /storage/manifest.json \
    --cloud_app_type aws \
    --external_cloud_app_id cf178df3-1d8c-46a5-86b7-974a941c4d80 \
    --name arikCloudApp   \
    --coverage_mode accounts \
    --auto_discover_subscriptions False \
    --operation_mode reveal \
    --cluster_id cloud \
    --accounts 324264561773
