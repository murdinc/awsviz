package aws

import (
	"github.com/robfig/revel"
	"github.com/ahmad972/goamz/aws"
	"github.com/ahmad972/goamz/ec2"
	"awsviz/app/models"
	"time"
)

func ListInstances() models.Flare {

	var flare models.Flare
	flare.Name = "Everything"

	temp := make(map[string]map[string][]models.Child)

	// Build the instances into their placeholder
	results := asyncApiCalls()
	for _, result := range results {

		// classes placeholder
		classes := make(map[string][]models.Child)
		var availZone string

		for _,res := range result.Reservations {
			for _,instance := range res.Instances {

				var instName string
				var instClass string

				for _,tag := range instance.Tags {
					if tag.Key == "Name" {
						instName = tag.Value
					}
					if tag.Key == "Class" {
						instClass = tag.Value
					}
				}

				var inst models.Child
				inst.Name = instName
				inst.Class = instClass
				inst.AvailZone = instance.AvailZone
				inst.InstanceType = instance.InstanceType

				classes[instClass] =  append(classes[instClass], inst)
				availZone = instance.AvailZone


			}
			temp[availZone] = classes
		}
	}

	// Combine everything in the correct order
	for tempRegion,tempClasses := range temp {
		var regionChild models.Child
		regionChild.Name = tempRegion

		for tempClass,tempInstances := range tempClasses {
			var classChild models.Child
			classChild.Name = tempClass
			for _, child := range tempInstances {

				classChild.Children = append(classChild.Children, child)

			}
			regionChild.Children = append(regionChild.Children, classChild)
		}
		flare.Children = append(flare.Children, regionChild)

	}

	return flare
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
