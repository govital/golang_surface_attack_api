package models

type InputJson struct {
	Vms      []Vm       `json:"vms"`
	Fw_rules []Fw_rules `json:"fw_rules"`
}

type Vm struct {
	Vm_id string   `json:"vm_id"`
	Name  string   `json:"name"`
	Tags  []string `json:"tags"`
}

type Fw_rules struct {
	Fw_id      string `json:"fw_id"`
	Source_tag string `json:"source_tag"`
	Dest_tag   string `json:"dest_tag"`
}
