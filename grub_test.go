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

func qemuTest(t *testing.T, extraArgs ...string) (*expect.GExpect, func()) {
	vm, err := qemu.NewVM("", extraArgs...)
	if err != nil {
		t.Fatal(err)
	}

	// Use test name as log filename
	opts, warn := testsuite.ExpectOptions(t.Name())
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

	return e, func() {
		err = vm.PowerDown()
		if err != nil {
			t.Error(err)
		}

		err = vm.Close()
		if err != nil {
			t.Error(err)
		}
	}
}

func TestBootGrub(t *testing.T) {
	e, cleanup := qemuTest(t,
		"-hda", "grub/output/images/disk.img")
	defer cleanup()

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
}

func TestURootBootGrub(t *testing.T) {
	e, cleanup := qemuTest(t,
		"-kernel", "output/images/bzImage",
		"-initrd", "output/images/uroot.cpio",
		"-append", "console=ttyS0",
		"-hda", "grub/output/images/disk.img")
	defer cleanup()

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
}
