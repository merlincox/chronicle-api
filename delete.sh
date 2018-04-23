#!/usr/bin/env bash

set -euo pipefail

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

cf_stack=api-stack-${release}
api_bucket=api-mock-data-${release}

aws cloudformation describe-stacks --stack-name $cf_stack >/dev/null

echo $cf_stack exists and will be deleted

if aws s3 ls "s3://${api_bucket}" 2>&1 | grep -v -q 'NoSuchBucket' ;then

    echo $api_bucket exists and will be deleted
    aws s3 rm s3://$api_bucket --recursive
    aws s3 rb s3://$api_bucket
fi

aws cloudformation delete-stack --stack-name $cf_stack

