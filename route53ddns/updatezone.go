package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func createRecord(c Config, ipaddr string) {

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := route53.New(sess)

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(c.Record),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ipaddr),
							},
						},
						TTL: aws.Int64(300),
						// Weight:        aws.Int64(1),
						// SetIdentifier: aws.String("Primary Record"),
					},
				},
			},
			Comment: aws.String("Updated automatically via route53ddns"),
		},
		HostedZoneId: aws.String(c.ZoneID),
	}
	resp, err := svc.ChangeResourceRecordSets(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Change Response:")
	fmt.Println(resp)
}
