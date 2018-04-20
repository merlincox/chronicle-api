#!/usr/bin/env bash

set -euo pipefail

api_id=0si9vxu19d
stage_name=1804201241
package=models
package_dir=models

# Export a JSON-format Swagger API definition from the AWS Gateway API

aws apigateway get-export --rest-api-id $api_id  --stage-name $stage_name --export-type swagger export/$stage_name.json

# The schema-generator executable can be created from here: https://github.com/merlincox/generate
# It generates a Go source file of struct declarations from the Swagger API definition file

schema_generator=$(which schema-generator)

if [ ! -z $schema_generator ]; then

    $schema_generator -p $package -nsk export/$stage_name.json > export/$stage_name.go

fi
