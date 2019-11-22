package main

import (
	"testing"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
	expect "github.com/google/goexpect"
)

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

func TestURootBoot2Grub(t *testing.T) {
	e, cleanup := qemuTest(t,
		"-kernel", "output/images/bzImage",
		"-initrd", "output/images/uroot.cpio",
		"-append", "console=ttyS0",
		"-hda", "grub/output/images/disk.img")
	defer cleanup()

	batcher := append(testsuite.Linuxboot2urootBatcher,
		[]expect.Batcher{
			&expect.BSnd{S: "boot2\r\n"},
			&testsuite.BExpTLog{
				L: "kexec done",
				R: "kexec_core: Starting new kernel",
				T: 10,
			}}...)
	batcher = append(batcher, BuildrootBatcher...)

	res, err := e.ExpectBatch(batcher, 0)
	if err != nil {
		t.Errorf("u-root 'boot2' grub config: %v", testsuite.DescribeBatcherErr(batcher, res, err))
	}
}

func TestURootLocalbootGrub(t *testing.T) {
	e, cleanup := qemuTest(t,
		"-kernel", "output/images/bzImage",
		"-initrd", "output/images/uroot.cpio",
		"-append", "console=ttyS0",
		"-hda", "grub/output/images/disk.img")
	defer cleanup()

	batcher := append(testsuite.Linuxboot2urootBatcher,
		[]expect.Batcher{
			&expect.BSnd{S: "localboot -grub\r\n"},
			&testsuite.BExpTLog{
				L: "kexec done",
				R: "kexec_core: Starting new kernel",
				T: 10,
			}}...)
	batcher = append(batcher, BuildrootBatcher...)

	res, err := e.ExpectBatch(batcher, 0)
	if err != nil {
		t.Errorf("u-root 'localboot' grub config: %v", testsuite.DescribeBatcherErr(batcher, res, err))
	}
}
