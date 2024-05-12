#!/bin/bash

CONTRACTS_SRC_DIR="$HOME/work/contracts/guardicore/connector/v1"

# Find all YAML files and generate Go structs
find "$CONTRACTS_SRC_DIR" -name '*.yaml' -exec bash -c '
  FILE="$1"
  parentdir=$(basename "$(dirname "$FILE")")
  outdir=$(pwd)/agent/go/infra/model/$parentdir
  mkdir -p $outdir 
  OUTPUT_FILE="$outdir/$(basename ${FILE%.yaml}).go"
  echo "Generating Go code for $FILE -> $OUTPUT_FILE"
  oapi-codegen --package models \
   --import-mapping=../common/common.yaml:common/common.yaml,../common/inventory.yaml:common/inventory.yaml \
   --generate types "$FILE" > "$OUTPUT_FILE"
' sh {} \;
