package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)
	input := &ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("volume-id"),
				Values: []*string{
					// Volume ID of the EBS volume these snapshots are of.
					aws.String("vol-0d2d00e4b56405c25"),
				},
			},
		},
	}
	result, err := svc.DescribeSnapshots(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	now := time.Now()
	for _, s := range result.Snapshots {
		// Delete the snapshot if it is more than 5 days old
		if now.Sub(*s.StartTime).Hours() > 120 {
			input := &ec2.DeleteSnapshotInput{
				SnapshotId: s.SnapshotId,
			}
			del, err := svc.DeleteSnapshot(input)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(del)
		}
	}
}
