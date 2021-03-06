#!/usr/bin/env bash

set -euo pipefail

timestamp=$(date +"%y%m%d%H%M")
release_pattern="^[a-z0-9]+$"

if [ $# -ne 1 ]; then

    echo "USAGE $0 release-name"
    exit 1

fi

if [[ ! $1 =~ $release_pattern ]] ; then

    echo "Invalid release-name: $1 does not match $release_pattern regex"
    exit 1

fi

release=$1

cf_bucket=cf-api-import-${timestamp}
cf_stack=api-stack-${release}

api_bucket=api-mock-data-${release}
api_filename=data.json

cd $( dirname $0 )

run_dir=$(pwd)

package_yml=$run_dir/export/package.yml

bucket_created=0

function cleanup {

    if [ $bucket_created -eq 1 ] ; then
        aws s3 rm s3://$cf_bucket --recursive
        aws s3 rb s3://$cf_bucket
    fi
}

trap cleanup EXIT

aws s3api create-bucket --bucket $cf_bucket --create-bucket-configuration LocationConstraint=eu-west-2

bucket_created=1

aws-sam-local package \
       --template-file $run_dir/template.yml \
       --s3-bucket $cf_bucket \
       --output-template-file $package_yml

aws cloudformation deploy \
       --template-file $package_yml \
       --stack-name $cf_stack \
       --capabilities CAPABILITY_IAM \
       --parameter-overrides Release=${release} S3Bucket=${api_bucket} S3Filename=${api_filename}

aws s3 cp $run_dir/mock/${api_filename} s3://${api_bucket}/${api_filename}




