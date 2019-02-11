package vmutils

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestVM_IP(t *testing.T) {
	vmcontroller := os.Getenv("vmcontroller") 
	vmpath := os.Getenv("vmx")                
	vm := NewVM(vmcontroller, vmpath)
	start := time.Now()
	vm.Start()
	defer vm.Stop()
	fmt.Println("Waiting for IP at ", time.Now())
	if IP := vm.IPTimeout(1 * time.Hour); IP == "" {
		t.Fatal("Failed to get IP")
	} else {
		elapsed := time.Since(start)
		fmt.Println("IP address is: ", IP)
		fmt.Println("Time taken is: ", elapsed)
	}
}
