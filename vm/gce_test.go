package vm

import "testing"

func TestNewMachine(t *testing.T) {
	mngr, _ := NewGceManager("")
	spec := &VirtualMachineSpec{}
	info, err := mngr.NewMachine(spec)
	if err != nil {
		t.Error(err)
	}
	err = mngr.DelMachine(info)
	if err != nil {
		t.Error(err)
	}
}
