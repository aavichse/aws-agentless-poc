#!/bin/bash

# DynamoDB table name
TABLE_NAME="POCFinancialData"

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

# Scan the table and format output as CSV
echo "Scanning all items from $TABLE_NAME and printing in CSV format"
if $USE_INTERFACE; then
    aws dynamodb scan \
        --table-name $TABLE_NAME \
        --region us-east-1 \
        --endpoint "$ENDPOINT_URL" \
        --return-consumed-capacity TOTAL | \
    jq -r '.Items[] | [.Date.S, .Open.N, .High.N, .Low.N, .Close.N, .Adj_Close.N, .Volume.N] | @csv'
else
    aws dynamodb scan \
        --table-name $TABLE_NAME \
        --return-consumed-capacity TOTAL | \
    jq -r '.Items[] | [.Date.S, .Open.N, .High.N, .Low.N, .Close.N, .Adj_Close.N, .Volume.N] | @csv'
fi

# Check if the scan was successful
if [ $? -eq 0 ]; then
    echo "Scan executed and output formatted successfully."
else
    echo "Failed to execute scan or format output."
fi

