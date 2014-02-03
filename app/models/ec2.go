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
	Edges                   map[string]string `json:"edges"`
}
