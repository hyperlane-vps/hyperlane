package control

type DesiredVM struct {
	Name   string
	Cpu    int
	Ram    int
	Image  string
	NodeID string
}

type ObservedVM struct {
	Name  string
	State string
	Node  string
}

type ClusterState struct {
	Desired  map[string]DesiredVM
	Observed map[string]ObservedVM
}
