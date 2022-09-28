package ec2types

type EC2AggregatedInstances struct {
	Zones []EC2Zone
}

type EC2Zone struct {
	Zone      string
	Instances []EC2Instance
}

type EC2Instance struct {
	Id   string
	Name string
}
