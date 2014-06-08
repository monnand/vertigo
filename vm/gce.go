package vm

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// XXX(monnand): We know we should use oauth. But this is a hackathon.
type gceVmManager struct {
	gcutilPath string
}

const (
	gcutil_OP_ADD = iota
	gcutil_OP_DEL
)

type gcutilOperation int

func (self gcutilOperation) String() string {
	switch self {
	case gcutil_OP_ADD:
		return "addinstance"
	case gcutil_OP_DEL:
		return "deleteinstance"
	}
	return ""
}

type gcutlCmdParams struct {
	Name        string
	Op          gcutilOperation
	GcutilPath  string
	Zone        string
	MachineType string
	Image       string
}

func randomUniqString() string {
	var d [8]byte
	io.ReadFull(rand.Reader, d[:])
	str := hex.EncodeToString(d[:])
	return fmt.Sprintf("%x-%v", time.Now().Unix(), str)
}

func (self *gcutlCmdParams) fillDefault() {
	if self.Name == "" {
		self.Name = fmt.Sprintf("virtigo-%v", randomUniqString())
	}
	if self.Zone == "" {
		self.Zone = "us-central1-a"
	}
	if self.MachineType == "" {
		self.MachineType = "n1-standard-2"
	}
	if self.Image == "" {
		self.Image = "ubuntu-trusty"
	}
}

func (self *gcutlCmdParams) ToParamList() []string {
	self.fillDefault()
	ret := make([]string, 0, 6)
	ret = append(ret, self.Op.String())
	switch self.Op {
	case gcutil_OP_ADD:
		ret = append(ret, fmt.Sprintf("--zone=%v", self.Zone))
		ret = append(ret, fmt.Sprintf("--machine_type=%v", self.MachineType))
		ret = append(ret, fmt.Sprintf("--image=%v", self.Image))
	case gcutil_OP_DEL:
		ret = append(ret, "-f")
		ret = append(ret, "--nodelete_boot_pd")
	}
	ret = append(ret, self.Name)
	return ret
}

func specToAddInstanceCmd(spec *VirtualMachineSpec) *gcutlCmdParams {
	ret := &gcutlCmdParams{
		Op:    gcutil_OP_ADD,
		Name:  spec.Name,
		Image: spec.Image,
	}
	return ret
}

func infoToDelInstanceCmd(info *VirtualMachineInfo) *gcutlCmdParams {
	return &gcutlCmdParams{
		Op:   gcutil_OP_DEL,
		Name: info.Name,
	}
}

func (self *gceVmManager) runGcutil(params *gcutlCmdParams) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	var err error
	go func() {
		defer wg.Done()
		paramList := params.ToParamList()
		gcutil := self.gcutilPath
		if self.gcutilPath == "" {
			gcutil = "gcutil"
		}
		fmt.Printf("%v %v\n", gcutil, strings.Join(paramList, " "))
		cmd := exec.Command(gcutil, paramList...)
		err = cmd.Run()
	}()
	wg.Wait()
	if err != nil {
		return fmt.Errorf("unable to use gcutil: %v", err)
	}
	return nil
}

func (self *gceVmManager) NewMachine(spec *VirtualMachineSpec) (*VirtualMachineInfo, error) {
	params := specToAddInstanceCmd(spec)
	err := self.runGcutil(params)
	if err != nil {
		return nil, err
	}
	info := &VirtualMachineInfo{
		Name: params.Name,
	}
	return info, nil
}

func (self *gceVmManager) DelMachine(info *VirtualMachineInfo) error {
	params := infoToDelInstanceCmd(info)
	err := self.runGcutil(params)
	if err != nil {
		return err
	}
	return nil
}

func NewGceManager(gcutilPath string) (VirtualMachineManager, error) {
	ret := &gceVmManager{
		gcutilPath: "gcutil",
	}
	return ret, nil
}
