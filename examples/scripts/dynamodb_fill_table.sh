#!/bin/bash

set -e

# DynamoDB table details
TABLE_NAME="POCFinancialData"
PRIMARY_KEY="Date"  # Primary key for the table

# Date range for data generation: past 5 years
start_date=$(date -d "5 years ago" +%Y-%m-%d)
end_date=$(date +%Y-%m-%d)

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

# Generate data for each day
current_date="$start_date"
while [[ "$current_date" != "$end_date" ]]; do
    # Generate random financial data
    open=$(printf "%.2f" "$(echo "$RANDOM%100+100" | bc)")
    high=$(printf "%.2f" "$(echo "$open + $RANDOM%10" | bc)")
    low=$(printf "%.2f" "$(echo "$open - $RANDOM%10" | bc)")
    close=$(printf "%.2f" "$(echo "$low + $RANDOM%($high - $low+0.01)" | bc)")
    adj_close=$(printf "%.2f" "$(echo "$close" | bc)")
    volume=$(printf "%d" "$((RANDOM%10000+100000))")

    # Insert data into DynamoDB
    echo "Inserting data for Date: $current_date, $open, $high, $low, $close, $adj_close, $volume"
    
    
    if $USE_INTERFACE; then
        aws dynamodb put-item \
            --table-name $TABLE_NAME \
            --region us-east-1 \
            --item "{
                \"$PRIMARY_KEY\": {\"S\": \"$current_date\"},
                \"Open\": {\"N\": \"$open\"},
                \"High\": {\"N\": \"$high\"},
                \"Low\": {\"N\": \"$low\"},
                \"Close\": {\"N\": \"$close\"},
                \"Adj_Close\": {\"N\": \"$adj_close\"},
                \"Volume\": {\"N\": \"$volume\"}
            }" \
            --endpoint "$ENDPOINT_URL"
    else
        aws dynamodb put-item \
            --table-name $TABLE_NAME \
            --item "{
                \"$PRIMARY_KEY\": {\"S\": \"$current_date\"},
                \"Open\": {\"N\": \"$open\"},
                \"High\": {\"N\": \"$high\"},
                \"Low\": {\"N\": \"$low\"},
                \"Close\": {\"N\": \"$close\"},
                \"Adj_Close\": {\"N\": \"$adj_close\"},
                \"Volume\": {\"N\": \"$volume\"}
            }" 
    fi

    # Increment date by one day
    current_date=$(date -I -d "$current_date + 1 day")
done

