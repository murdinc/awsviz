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
