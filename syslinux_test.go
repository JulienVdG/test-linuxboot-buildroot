package main

import (
	"testing"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
	expect "github.com/google/goexpect"
)

func TestBootIsolinux(t *testing.T) {
	e, cleanup := qemuTest(t,
		"-cdrom", "syslinux/output/images/rootfs.iso9660")
	defer cleanup()

	batcher := append([]expect.Batcher{
		&testsuite.BExpTLog{
			L: "Matched Linux starting",
			R: "Linux version",
			T: 50,
		}}, BuildrootBatcher...)
	res, err := e.ExpectBatch(batcher, 0)
	if err != nil {
		t.Errorf("booting trough isolinux: %v", testsuite.DescribeBatcherErr(batcher, res, err))
	}
}

func TestURootBootIsolinux(t *testing.T) {
	e, cleanup := qemuTest(t,
		"-kernel", "output/images/bzImage",
		"-initrd", "output/images/uroot.cpio",
		"-append", "console=ttyS0",
		"-cdrom", "syslinux/output/images/rootfs.iso9660")
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
		t.Errorf("u-root 'boot' isolinux config: %v", testsuite.DescribeBatcherErr(batcher, res, err))
	}
}

func TestURootBoot2Isolinux(t *testing.T) {
	e, cleanup := qemuTest(t,
		"-kernel", "output/images/bzImage",
		"-initrd", "output/images/uroot.cpio",
		"-append", "console=ttyS0",
		"-cdrom", "syslinux/output/images/rootfs.iso9660")
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
		t.Errorf("u-root 'boot2' isolinux config: %v", testsuite.DescribeBatcherErr(batcher, res, err))
	}
}

func TestURootLocalbootIsolinux(t *testing.T) {
	e, cleanup := qemuTest(t,
		"-kernel", "output/images/bzImage",
		"-initrd", "output/images/uroot.cpio",
		"-append", "console=ttyS0",
		"-cdrom", "syslinux/output/images/rootfs.iso9660")
	defer cleanup()

	batcher := append(testsuite.Linuxboot2urootBatcher,
		[]expect.Batcher{
			&expect.BSnd{S: "localboot -isolinux\r\n"},
			&testsuite.BExpTLog{
				L: "kexec done",
				R: "kexec_core: Starting new kernel",
				T: 10,
			}}...)
	batcher = append(batcher, BuildrootBatcher...)

	res, err := e.ExpectBatch(batcher, 0)
	if err != nil {
		t.Errorf("u-root 'localboot' isolinux config: %v", testsuite.DescribeBatcherErr(batcher, res, err))
	}
}
