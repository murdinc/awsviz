package models

import (
)

type    Instance                struct {
	Name                    string
	InstanceId              string
	PrivateIPAddress        string
	PublicIPAddress         string
	VpcId                   string
	SubnetId                string
	ImageId                 string
	RootDeviceType          string
	InstanceType            string
	State                   string
	KeyName                 string
	AvailZone               string
}

type	Node			struct {
	Name			string `json:"name"`
	Class			string `json:"class"`
	Region                  string `json:"region"`
}

type	Arbor			struct {
	Nodes			map[string]Node `json:"nodes"`
	Edges                   map[string]map[string][]string `json:"edges"`
	Regions			[]string `json:"regions"`
	Classes			map[string][]string `json:"classes"`
}


type	Flare		struct {
        Name            string          `json:"name"`
	Children	[]Child		`json:"children"`
	Class           string          `json:"class"`
	AvailZone       string          `json:"availZone"`
}

type	Child		struct {
        Name            string          `json:"name"`
	InstanceType    string          `json:"instanceType"`
	Children	[]Child		`json:"children"`
	Class		string		`json:"class"`
	AvailZone	string		`json:"availZone"`
}
