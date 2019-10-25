package main

import (
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/qemu"
	"github.com/JulienVdG/tastevin/pkg/testsuite"
	expect "github.com/google/goexpect"
)

var BuildrootBatcher []expect.Batcher = []expect.Batcher{
	&testsuite.BExpTLog{
		L: "Matched Buildroot banner",
		R: "Welcome to Buildroot",
		T: 40,
	},
	&testsuite.BExpTLog{
		L: "Matched Buildroot login prompt",
		R: "buildroot login:",
		T: 5,
	},
	&expect.BSnd{S: "root\r\n"},
	&expect.BExpT{
		R: "Password:",
		T: 5,
	},
	&expect.BSnd{S: "r\r\n"},
	&testsuite.BExpTLog{
		L: "Matched Buildroot prompt",
		R: "# ",
		T: 5,
	},
	&expect.BSnd{S: "poweroff\r\n"},
	&testsuite.BExpTLog{
		L: "Matched Power down",
		R: "reboot: Power down",
		T: 10,
	},
}

func TestBootGrub(t *testing.T) {
	vm, err := qemu.NewVM("",
		"-hda", "grub/output/images/disk.img")
	if err != nil {
		t.Fatal(err)
	}

	opts, warn := testsuite.ExpectOptions("")
	if warn != nil {
		t.Log(warn)
	}

	e, _, err := vm.Spawn(1*time.Second, opts...)
	if err != nil {
		t.Fatalf("Spawn failed: %v", err)
	}
	err = vm.PowerUp()
	if err != nil {
		t.Error(err)
	}

	batcher := append([]expect.Batcher{
		&testsuite.BExpTLog{
			L: "Matched Linux starting",
			R: "Linux version",
			T: 50,
		}}, BuildrootBatcher...)
	res, err := e.ExpectBatch(batcher, 0)
	if err != nil {
		t.Errorf("booting trough grub: %v", testsuite.DescribeBatcherErr(batcher, res, err))

	}

	err = vm.PowerDown()
	if err != nil {
		t.Error(err)
	}

	err = vm.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestURootBootGrub(t *testing.T) {
	vm, err := qemu.NewVM("",
		"-kernel", "output/images/bzImage",
		"-initrd", "output/images/uroot.cpio",
		"-append", "console=ttyS0",
		"-hda", "grub/output/images/disk.img")
	if err != nil {
		t.Fatal(err)
	}

	opts, warn := testsuite.ExpectOptions("")
	if warn != nil {
		t.Log(warn)
	}

	e, _, err := vm.Spawn(1*time.Second, opts...)
	if err != nil {
		t.Fatalf("Spawn failed: %v", err)
	}
	err = vm.PowerUp()
	if err != nil {
		t.Error(err)
	}

	batcher := append(testsuite.Linuxboot2urootBatcher,
		[]expect.Batcher{
			&expect.BSnd{S: "boot\r\n"},
			&testsuite.BExpTLog{
				L: "kexec done",
				R: "kexec_core: Starting new kernel",
				T: 10,
			}}...)
	batcher = append(batcher, BuildrootBatcher...)

	res, err := e.ExpectBatch(batcher, 0)
	if err != nil {
		t.Errorf("u-root 'boot' grub config: %v", testsuite.DescribeBatcherErr(batcher, res, err))

	}

	err = vm.PowerDown()
	if err != nil {
		t.Error(err)
	}

	err = vm.Close()
	if err != nil {
		t.Error(err)
	}
}
