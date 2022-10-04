package aws

import (
	"context"
	"fmt"
	"kamogawa/types"
	"kamogawa/types/aws/ec2types"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awsec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func AWSListEC2Instances(db *gorm.DB, user types.User, useCache bool) ec2types.EC2AggregatedInstances {
	return mockData()

	//if config.CacheEnabled && useCache {
	//	return mockData()
	//}
	//
	//return AWSListEC2InstancesMain(db, user)
}

func AWSListEC2InstancesMain(db *gorm.DB, user types.User) ec2types.EC2AggregatedInstances {
	cfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "AKIAQ6ESS3N7RZYDUZVA", SecretAccessKey: "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",
				Source: "DiceDuckMonk",
			},
		}),
		awsconfig.WithRegion("us-west-2"),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := ec2.NewFromConfig(cfg)

	describeRegionsInput := ec2.DescribeRegionsInput{Filters: []awsec2types.Filter{{Name: aws.String("opt-in-status"), Values: []string{"opt-in-not-required", "opted-in"}}}}
	describeRegionsResult, err := client.DescribeRegions(context.TODO(), &describeRegionsInput)

	regions := lo.Map[awsec2types.Region, string](describeRegionsResult.Regions, func(r awsec2types.Region, _ int) string {
		return *r.RegionName
	})

	ec2Instances := lo.Flatten[awsec2types.Instance](lo.Flatten[[]awsec2types.Instance](lo.FilterMap[string, [][]awsec2types.Instance](regions, func(region string, _ int) ([][]awsec2types.Instance, bool) {
		cfg, err := awsconfig.LoadDefaultConfig(
			context.TODO(),
			awsconfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID: "AKIAQ6ESS3N7RZYDUZVA", SecretAccessKey: "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",
					Source: "DiceDuckMonk",
				},
			}),
			awsconfig.WithRegion(region),
		)
		if err != nil {
			panic("configuration error, " + err.Error())
		}
		client := ec2.NewFromConfig(cfg)

		input := &ec2.DescribeInstancesInput{}
		result, err := GetInstances(context.TODO(), client, input)
		if err != nil {
			fmt.Println("Got an error retrieving information about your Amazon EC2 instances:")
			fmt.Println(err)
			return [][]awsec2types.Instance{}, false
		}

		instances := lo.Map[awsec2types.Reservation, []awsec2types.Instance](result.Reservations, func(r awsec2types.Reservation, _ int) []awsec2types.Instance {
			return r.Instances
		})
		return instances, len(instances) > 0
	})))

	zoneMap := make(map[string][]ec2types.EC2Instance)
	for _, i := range ec2Instances {
		nameTag := lo.Filter[awsec2types.Tag](i.Tags, func(t awsec2types.Tag, _ int) bool {
			return *t.Key == "Name"
		})[0]

		zoneMap[*i.Placement.AvailabilityZone] = append(zoneMap[*i.Placement.AvailabilityZone], ec2types.EC2Instance{Id: *i.InstanceId, Name: *nameTag.Value})
	}

	var ec2AggregatedInstances ec2types.EC2AggregatedInstances
	for zone, instances := range zoneMap {
		ec2AggregatedInstances.Zones = append(ec2AggregatedInstances.Zones, ec2types.EC2Zone{Zone: zone, Instances: instances})
	}

	return ec2AggregatedInstances
}

// EC2DescribeInstancesAPI defines the interface for the DescribeInstances function.
// We use this interface to test the function using a mocked service.
type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// GetInstances retrieves information about your Amazon Elastic Compute Cloud (Amazon EC2) instances.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a DescribeInstancesOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to DescribeInstances.
func GetInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

func mockData() ec2types.EC2AggregatedInstances {
	return ec2types.EC2AggregatedInstances{
		Zones: []ec2types.EC2Zone{
			ec2types.EC2Zone{
				Zone: "us-west-1c ",
				Instances: []ec2types.EC2Instance{
					ec2types.EC2Instance{Id: "shimogawa", Name: "shimogawa"},
					ec2types.EC2Instance{Id: "akari", Name: "akari"},
					ec2types.EC2Instance{Id: "ichiban", Name: "ichiban"},
					ec2types.EC2Instance{Id: "ichiro", Name: "ichiro"},
				},
			},
			ec2types.EC2Zone{
				Zone: "us-west-2c ",
				Instances: []ec2types.EC2Instance{
					ec2types.EC2Instance{Id: "kaze", Name: "kaze"},
					ec2types.EC2Instance{Id: "oni", Name: "oni"},
					ec2types.EC2Instance{Id: "moku", Name: "moku"},
					ec2types.EC2Instance{Id: "mizu", Name: "mizu"},
					ec2types.EC2Instance{Id: "oto", Name: "oto"},
					ec2types.EC2Instance{Id: "kumo", Name: "kumo"},
				},
			},
		},
	}
}
