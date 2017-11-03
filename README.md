# route53ddns
Route53 Dynamic DNS Client

This is a simple DDNS client for Route 53.


## Prerequisites:
- AWS Account
- Hosted Zone on the AWS account
- IAM user with programmatic access to modify the hosted zone
- Note: **_aws cli is not strictly required_**

## Information Required:

You will need 4 pieces of information:
- AWS IAM Access Key ID
- AWS IAM Secret Key 
- The Zone ID you are updating
- The desired FQDN e.g. `test.test.tld`

## Installation
- Install and place route53ddns.exe and conf.toml side by side in a directory
- Run `route53ddns.exe --configure` 
- Give it the information referenced in the `information required` section.

## Running
- Simply run `route53ddns.exe`

## Scheduled Tasking:
This is pretty easy to do:

- Create a task in Windows Task Scheduler to whatever schedule you want
- For action; configure to start a program pointed at the exe
- Under the action properties, insert the full path including the trailing \ where the conf.toml lives.
