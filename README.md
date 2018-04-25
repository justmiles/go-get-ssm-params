# get-ssm-params
Grab values out of the AWS SSM Parameter Store

## Installation

    sudo curl -L https://github.com/justmiles/go-get-ssm-params/releases/download/v1.2.0/get-ssm-params.v1.2.0.linux-amd64 -o /usr/local/bin/get-ssm-params
    sudo chmod +x /usr/local/bin/get-ssm-params

## Usage
Group your parameters in SSM by path. When you retrieve them with get-ssm-params, parameters in latest path you provide will overwrite any previous ones. For the examples below, this is what is in SSM:

    /dev/default/MY_CONFIG_KEY=myconfigvalue
    /dev/default/DB_HOSTNAME=db.dev.mycompany.com
    /dev/default/DB_PASSWORD=password
    
    /dev/myapp/MY_CONFIG_KEY=overridden

as JSON

    > get-ssm-params -output json -path /dev/default -path /dev/myapp
    
    {
      "MY_CONFIG_KEY": "overridden",
      "DB_HOSTNAME": "db.dev.mycompany.com"
      "DB_PASSWORD": "password"
    }

    
as shell

    > get-ssm-params -output shell -path /dev/default -path /dev/myapp
    
    export MY_CONFIG_KEY="overridden"
    export DB_HOSTNAME="db.dev.mycompany.com"
    export DB_PASSWORD="password"
    
as text

    > get-ssm-params -output text -path /dev/default -path /dev/myapp
    
    MY_CONFIG_KEY=overridden
    DB_HOSTNAME=db.dev.mycompany.com
    DB_PASSWORD=password
    
as ECS

    > get-ssm-params -output ecs -path /dev/default -path /dev/myapp
    
    [
      {
        "name":"MY_CONFIG_KEY",
        "value":"overridden"
      },
      {
        "name":"DB_HOSTNAME",
        "value":"db.dev.mycompany.com"
      },
      {
        "name":"DB_PASSWORD",
        "value":"password"
      }
    ]

Usage of get-ssm-params:

    -path value
        SSM Parameter path

    -output string
        set the desired output (default "json")
        
    -region string
        aws region (default "us-east-1")
        
    -version
        show current version


