#!/bin/bash

# Set the bucket name
BUCKET_NAME="poc-examples"

# Set the file name to be uploaded
FILE_NAME="awscliv2.zip"


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

if [ -f "$FILE_NAME" ]; then

    if $USE_INTERFACE; then
        echo "Uploading $FILE_NAME to $BUCKET_NAME using VPC endpoint"
        aws s3 cp "$FILE_NAME" "s3://$BUCKET_NAME/" --endpoint-url "$ENDPOINT_URL"
    else
        # Upload the file to the bucket without the endpoint URL
        echo "Uploading $FILE_NAME to $BUCKET_NAME"
        aws s3 cp "$FILE_NAME" "s3://$BUCKET_NAME/"
    fi

    if [ $? -eq 0 ]; then
        echo "File uploaded successfully."
    else
        echo "Failed to upload the file."
    fi
else
    echo "Error: File $FILE_NAME does not exist."
fi
