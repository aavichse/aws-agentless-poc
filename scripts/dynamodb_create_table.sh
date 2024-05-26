#!/bin/bash

set -e

# DynamoDB table details
TABLE_NAME="POCFinancialData"
PRIMARY_KEY="Date"  # Primary key for the table
PRIMARY_KEY_TYPE="S"  # N for Number, S for String


USE_INTERFACE=false
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --interface) 
                USE_INTERFACE=true 
                shift 
                ENDPOINT_URL="https://$1"
                ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

# Create the DynamoDB table
echo "Creating DynamoDB table: $TABLE_NAME"

if $USE_INTERFACE; then
    aws dynamodb create-table \
        --table-name $TABLE_NAME \
        --region us-east-1 \
        --attribute-definitions AttributeName=$PRIMARY_KEY,AttributeType=$PRIMARY_KEY_TYPE \
        --key-schema AttributeName=$PRIMARY_KEY,KeyType=HASH \
        --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
        --endpoint "$ENDPOINT_URL"
else
    aws dynamodb create-table \
        --table-name $TABLE_NAME \
        --attribute-definitions AttributeName=$PRIMARY_KEY,AttributeType=$PRIMARY_KEY_TYPE \
        --key-schema AttributeName=$PRIMARY_KEY,KeyType=HASH \
        --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
fi
# Check if the table creation was successful
if [ $? -ne 0 ]; then
    echo "Failed to create the table."
    exit 1
else
    echo "Table created successfully."
fi

