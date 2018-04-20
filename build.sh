#!/bin/bash

set -euo pipefail

region=eu-west-2
account=430933187925
role=lambda_basic_execution
release=12345
function_name=ApiLambda-${release}

GOOS=linux go build -o main bin/lambda.go

if [ -f main ]; then

    zip deployment.zip main

    if aws lambda get-function --function-name $function_name 2>/dev/null; then
        aws lambda update-function-code --region $region --function-name $function_name --zip-file fileb://./deployment.zip
    else
        echo $function_name not found
#        aws lambda create-function --region $region --function-name $function_name --zip-file fileb://./deployment.zip --runtime go1.x --role arn:aws:iam::$account:role/$role --handler main
    fi

fi

rm deployment.zip
rm main
