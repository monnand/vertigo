package vm

type VirtualMachineSpec struct {
	Name        string
	CpuLevel    int
	MemoryLevel int
	Image       string
}

type VirtualMachineInfo struct {
	Name    string
	Address string
}

type VirtualMachineManager interface {
	NewMachine(spec *VirtualMachineSpec) (*VirtualMachineInfo, error)
	DelMachine(info *VirtualMachineInfo) error
}
