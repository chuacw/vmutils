package vmutils

import (
	"bytes"
	"net"
	"os/exec"
	"strings"
	"time"
)

type (
	// VM is the structure for a VM configuration
	VM struct {
		controller string
		vmx        string
		vmStarted  bool
	}
)

// NewVM creates a new VM, with controller being the location of the vmrun utility and vmx being the location of the vmx to start.
func NewVM(controller, vmx string) (vm *VM) {
	vm = &VM{controller: controller, vmx: vmx}
	return
}

// IPTimeoutInMS waits for the IP for the VM, timing out in timeoutInMS, returning an empty string, or the IP if the timeout isn't reached.
func (vm *VM) IPTimeoutInMS(timeoutInMS uint64) (result string) {
	timeout := timeoutInMS
	result = vm.IPTimeout(time.Duration(timeout) * time.Millisecond)
	return
}

// IPTimeout waits for the specified timeout, or returns the IP for the VM
func (vm *VM) IPTimeout(timeout time.Duration) (result string) {
	var (
		ipAddr  string
		elapsed time.Duration
	)
	start := time.Now()
	timeoutInMS := uint64(timeout) / uint64(time.Millisecond)
	for ok := true; ok; ok = result == "" {
		cmd := exec.Command(vm.controller, "-T", "ws", "readVariable", vm.vmx, "guestVar", "ip")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		ipAddr = strings.TrimSuffix(out.String(), "\n")
		ipAddr = strings.TrimSuffix(ipAddr, "\r")
		if ipAddr != "" {
			ip := net.ParseIP(ipAddr)
			if ip != nil {
				result = ipAddr
			}
		}
		elapsed = time.Since(start)
		elapsedMS := uint64(elapsed / time.Millisecond)
		if elapsedMS > timeoutInMS {
			return
		}
	}
	return
}

// IP returns the IP for a VM, waiting infinitely if it doesn't get an IP. Requires the VM to be started.
func (vm *VM) IP() (result string) {
	if !vm.vmStarted {
		vm.Start()
	}
	for ok := true; ok; ok = result == "" {
		cmd := exec.Command(vm.controller, "-T", "ws", "readVariable", vm.vmx, "guestVar", "ip")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		ipAddr := strings.TrimSuffix(out.String(), "\n")
		ipAddr = strings.TrimSuffix(ipAddr, "\r")
		if ipAddr != "" {
			ip := net.ParseIP(ipAddr)
			if ip != nil {
				result = ipAddr
			}
		}
	}
	return
}

func (vm *VM) getWSCmd(cmd string) (result *exec.Cmd) {
	result = exec.Command(vm.controller, "-T", "ws", cmd, vm.vmx)
	return
}

// Start starts the VM
func (vm *VM) Start() {
	cmd := vm.getWSCmd("start")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	vm.vmStarted = true
}

// Stop stops the currently running VM
func (vm *VM) Stop() {
	cmd := vm.getWSCmd("stop")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	vm.vmStarted = false
}
