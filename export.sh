#!/usr/bin/env bash

set -euo pipefail

release_pattern="^[a-z0-9]+$"

if [ $# -lt 1 ] ; then

    echo "USAGE $0 release-name [stage-name]"
    exit 1
fi

if [[ ! $1 =~ $release_pattern ]] ; then

    echo "Invalid release-name: $1 does not match $release_pattern regex"
    exit 1
fi

release=$1
package=models
package_dir=models
export_dir=export

# Get the API id from the API name using jq command-line json tool. See https://stedolan.github.io/jq

api_id=$( aws apigateway get-rest-apis | jq  -r '.items[] | select(.name == "Chronicle-API-'${release}'") | .id' )

# Export a JSON-format Swagger API definition from the AWS Gateway API

aws apigateway get-export --rest-api-id $api_id  --stage-name $release --export-type swagger $export_dir/$release.json

# The schema-generator executable can be created from here: https://github.com/merlincox/generate
# It generates a Go source file of struct declarations from the Swagger API definition file

schema_generator=$(which schema-generator)

if [ ! -z "$schema_generator" ]; then

    $schema_generator -p $package -nsk $export_dir/$release.json > $export_dir/$release.go

fi
