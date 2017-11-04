go build ../r53conf
go build -ldflags -H=windowsgui ../route53ddns
Copy-Item ../r53conf/conf.toml ./