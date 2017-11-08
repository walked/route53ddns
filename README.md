# route53ddns
Route53 Dynamic DNS Client

This is a simple DDNS client for Route 53. It is meant to be lightweight, simple, and effective. In my testing it utilizes approx 4mb of memory and virtually zero CPU in Windows 10.


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
- Extract and place route53ddns.exe and conf.toml side by side in a directory 
- Run `route53ddns.exe --configure` 
- Give it the information referenced in the `information required` section.
- Run `route53ddns.exe --install` to install as a Windows Service
- Start the service via powershell (or in services.msc) `Start-Service -Name "R53DDNSSRV"` (it will now automatically start on reboots)

## Uninstall
- Simply run `route53ddns.exe --uninstall` 


---

## Known Issues

- Currently only checks against the DNS record that exists in Route53 at Service Start. This means that if your Route53 DNS record drifts, this will not catch it until the service is restarted. This will be fixed.

- Cannot run manually by hand at this time. This will be fixed so that the executable can also be triggered by external scripts.