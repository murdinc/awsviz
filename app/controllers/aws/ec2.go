package aws

import (
	"github.com/robfig/revel"
	"github.com/ahmad972/goamz/aws"
	"github.com/ahmad972/goamz/ec2"
	"awsgraph/app/models"
	"time"
)

func ListInstances() models.Arbor {

	var arbor models.Arbor

	nodes := map[string]models.Node{}
	//groups := map[string]map[string][]string{}
	regions := []string{}
	classes := map[string][]string{}
	edges := map[string]map[string][]string{}

	results := asyncApiCalls()
	for _, result := range results {

		if ( result.Reservations != nil ) {
			//revel.INFO.Printf("Region with instances: %s", result)
		}

		for _,res := range result.Reservations {
			for _,instance := range res.Instances {

				// Build our region "nodes" when we come across one we haven't seen before.
				if _,ok := nodes[instance.AvailZone]; ok {
					//revel.INFO.Printf("Region node already found: %s", instance.AvailZone)
				} else {
					nodes[instance.AvailZone] = models.Node{
						Name:		instance.AvailZone,
						Class:          "Availability Zone",
						Region:         instance.AvailZone,
					}
					// Add region to our regions
					regions = append(regions, instance.AvailZone)
				}

				var name, class string
				for _,tag := range instance.Tags {
					if tag.Key == "Name" {
						name = tag.Value
					}
					if tag.Key == "Class" {
						class = tag.Value
					}
				}

				// Add the instance to our nodes
				nodes[instance.InstanceId] = models.Node{
					Name:           name,
					Class:		class,
					Region:		instance.AvailZone,
				}

				// Build our class "nodes" when we come across one we haven't seen before.
				nodename := class + "." + instance.AvailZone
				if _,ok := nodes[nodename]; ok {
					//revel.INFO.Printf("CLASS FOUND: %s", class)                                                                                           
				} else {
					nodes[nodename] = models.Node{
						Name:           class,
						Class:          "Class",
						Region:         instance.AvailZone,
					}

				}
				// Add classes to our groups
				classes[instance.AvailZone] = append(classes[instance.AvailZone], class)

			}
		}
	}

	// Build our edges for regions and classes
	for _,region := range regions {
		revel.INFO.Printf("Building region: %s", region)
		tempedge := map[string][]string{}

		// Connect our regions together
		for _,targetregion := range regions {
			if targetregion != region {
				tempedge[targetregion] = nil
			}
		}

		// Connect our classes to their regions
		for classregion,classlist := range classes {
			if classregion == region {
				for _,regionclass := range classlist {
					edgename := regionclass + "." + region
					tempedge[edgename] = nil
				}
			}
		}

		edges[region] = tempedge
	}


	// Build our edges for instances
	for classregion, classlist := range classes {
		for _,regionclass := range classlist {

			revel.INFO.Printf("Building class: %s in: %s", classregion, regionclass)
			edgename := regionclass + "." + classregion
			tempedge := map[string][]string{}
			for nodename,val := range nodes {
				if (val.Class == regionclass && val.Region == classregion) {
					//tempedge := map[string][]string{}
					revel.INFO.Printf("Found class: %s and including: %s", regionclass, nodename)
					tempedge[nodename] = nil
				//	edges[edgename] = tempedge
				}
				edges[edgename] = tempedge

			}
		}
	}

	arbor.Nodes = nodes
	arbor.Regions = regions
	arbor.Classes = classes
	arbor.Edges = edges

	return arbor
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
