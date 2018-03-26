# get-ssm-params
Grab values out of the AWS SSM Parameter Store

## Installation

    sudo curl -L https://github.com/justmiles/go-get-ssm-params/releases/download/v1.0.0/get-ssm-params.v1.0.0.linux-amd64 -o /usr/local/bin/get-ssm-params
    sudo chmod +x /usr/local/bin/get-ssm-params

## Usage
    
    # as JSON
    > get-ssm-params -path /dev/default -path /dev/myapp -path /dev/database
    
    {
      "Parameters": [
        {
          "key": "SOME_PARAMETER",
          "value": "some parameter value"
        }
      ]
    }
    
    # as shell
    > get-ssm-params -shell -path /dev/default -path /dev/myapp -path /dev/database
    
    export SOME_PARAMETER="some parameter value"



Usage of get-ssm-params:
    
    -path value
        SSM Parameters path to source (can be passed multiple times)
    
    -region string
        aws region (default "us-east-1")
    
    -shell
        optionally output shell command to export variables. otherwise, output as json

