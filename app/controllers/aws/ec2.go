package aws

import (
	"github.com/robfig/revel"
	"github.com/ahmad972/goamz/aws"
	"github.com/ahmad972/goamz/ec2"
	"awsgraph/app/models"
	"time"
)

func ListInstances() []models.Instance {

	var instances []models.Instance

	results := asyncApiCalls()
	for _, result := range results {

		for _,res := range result.Reservations {
			for _,instance := range res.Instances {

				revel.INFO.Printf("Instance: %s", instance.InstanceId)
				instances = append(instances,models.Instance{
					Name:"", //TODO
					InstanceId:             instance.InstanceId,
					PrivateIPAddress:       instance.PrivateIPAddress,
					PublicIPAddress:        instance.IPAddress,
					VpcId:                  instance.VpcId,
					SubnetId:               instance.SubnetId,
					ImageId:                instance.ImageId,
					RootDeviceType:         instance.RootDeviceType,
					InstanceType:           instance.InstanceType,
					State:                  instance.State.Name,
					KeyName:                instance.KeyName,
					AvailZone:              instance.AvailZone,
				})
			}
		}
	}
	return instances
}

func asyncApiCalls() []ec2.InstancesResp {
	ch := make(chan *ec2.InstancesResp)
	responses := []ec2.InstancesResp{}

	for name,region := range aws.Regions {
		go func(name string, region aws.Region) {
			revel.INFO.Printf("Fetching region: %s", name)
			auth, err := aws.GetAuth("", "", "", time.Time{})

			if err != nil {
				panic(err)
			}
			e := ec2.New(auth, region)
			resp, err := e.Instances(nil, nil)

			if err != nil {
				panic(err)
			}

			ch <- resp
		}(name, region)
	}

	for {
		select {
		case r := <-ch:
			revel.INFO.Printf("Region fetched with requestId: %s", r.RequestId)
			responses = append(responses, *r)
			if len(responses) == len(aws.Regions) {
				return responses
			}
		case <-time.After(15 * time.Second):
			revel.INFO.Print("waiting...")
		}
	}
	return responses
}
