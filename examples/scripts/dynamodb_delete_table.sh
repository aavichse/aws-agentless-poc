#!/bin/bash

# DynamoDB table details
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

echo "Deleting DynamoDB table: $TABLE_NAME"
if $USE_INTERFACE; then
    aws dynamodb delete-table --table-name $TABLE_NAME  --endpoint-url "$ENDPOINT_URL"

else
    aws dynamodb delete-table --table-name $TABLE_NAME
fi

# Check if the table deletion was successful
if [ $? -eq 0 ]; then
    echo "Table deleted successfully."
else
    echo "Failed to delete the table."
fi

