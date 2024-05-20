#!/bin/bash

CONTRACTS_SRC_DIR="$HOME/work/contracts/guardicore/connector/"

oapi-codegen \
  --package api \
  --import-mapping=../common/common.yaml:common/common.yaml,../common/inventory.yaml:common/inventory.yaml \
  --generate gin  "$CONTRACTS_SRC_DIR/version.yaml" > openapi_server.version.gen.go


#!/bin/bash

CONTRACTS_SRC_DIR="$HOME/work/contracts/guardicore/connector/v1/operations"

# Find all YAML files and generate Go structs
find "$CONTRACTS_SRC_DIR" -name '*.yaml' -exec bash -c '
  FILE="$1"
  name=$(basename ${FILE%.yaml})
  OUTPUT_FILE="openapi_server.$name.gen.go"
  echo "Generating Go code for $FILE -> $OUTPUT_FILE"
  oapi-codegen --package api \
   --generate gin "$FILE" > "$OUTPUT_FILE"
' sh {} \;


