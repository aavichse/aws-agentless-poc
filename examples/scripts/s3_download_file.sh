#!/bin/bash

# Set the S3 bucket name and file name
BUCKET_NAME="poc-examples"
FILE_NAME="awscliv2.zip"

# Path where you want to download the file
DOWNLOAD_PATH="./$FILE_NAME"  # Downloads to the current directory

VPC_ENDPOINT_NAME="poc-examples-vpcep-interface-s3"

USE_INTERFACE=false
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --interface) 
                USE_INTERFACE=true 
                shift 
                ENDPOINT_URL="https://bucket.$1"
                ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

# Download the file from S3
echo "Downloading $FILE_NAME from bucket $BUCKET_NAME..."
if $USE_INTERFACE; then
        echo "Uploading $FILE_NAME to $BUCKET_NAME using VPC endpoint"
        aws s3 cp "$FILE_NAME" "s3://$BUCKET_NAME/" --endpoint-url "$ENDPOINT_URL"
else
        aws s3 cp "s3://$BUCKET_NAME/$FILE_NAME" "$DOWNLOAD_PATH"
fi

# Check if the download was successful
if [ $? -eq 0 ]; then
    echo "File downloaded successfully to $DOWNLOAD_PATH."
else
    echo "Failed to download the file."
fi
