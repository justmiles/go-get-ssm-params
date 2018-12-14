# get-ssm-params
Grab values out of the AWS SSM Parameter Store

## Installation

    sudo curl -L https://github.com/justmiles/go-get-ssm-params/releases/download/v1.6.1/get-ssm-params.v1.6.1.linux-amd64 -o /usr/local/bin/get-ssm-params
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
    
as terraform-ecs

    > get-ssm-params -output terraform-ecs -path /dev/default -path /dev/myapp
    
    {
      "JSONString": "[{\"name\":\"MY_CONFIG_KEY\",\"value\":\"overridden\"}]"
    }
    
as a custom template
    
    > echo "The database hostname is {{.DB_HOSTNAME}} and the password is {{.DB_PASSWORD}}." > example-template.tpl
    > get-ssm-params -template example-template.tpl -path /dev/default -path /dev/myapp
    
    The database hostname is db.dev.mycompany.com and the password is password.

Usage of get-ssm-params:

    -key string
        if specified, gets a single key and sends value to stdout
    -output string
        set the desired output (default "json")
    -path value
        SSM Parameter path
    -region string
        aws region (default "us-east-1")
    -template string
        if specified, renders custom output from a template file
    -version
        show current version


