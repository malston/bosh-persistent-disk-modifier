package vc

import (
	"context"
	"strings"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
)

func GetPersistentDiskName(ctx context.Context, c *vim25.Client, vmName string) (error, string) {
	finder := find.NewFinder(c)
	vm, err := finder.VirtualMachine(ctx, vmName)
	devices, err := vm.Device(ctx) // get the VM's virtual device list
	if err != nil {
		return err, ""
	}

	for i := range devices {
		switch d := devices[i].(type) {
		case *types.VirtualDisk:
			switch d.Backing.(type) {
			case *types.VirtualDiskFlatVer2BackingInfo:
				disk := d.Backing.(*types.VirtualDiskFlatVer2BackingInfo)
				if disk.DiskMode == "independent_persistent" {
					return nil, removePath(disk.FileName)
				}
			}
		}
	}

	return nil, ""
}

func removePath(filename string) string {
	if filename == "" {
		return ""
	}

	return strings.Split(filename[strings.LastIndex(filename, "/")+1:], ".")[0]
}
