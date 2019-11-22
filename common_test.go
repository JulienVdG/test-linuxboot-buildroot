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
