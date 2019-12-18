package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
	expect "github.com/google/goexpect"
	"github.com/u-root/u-root/pkg/multiboot"
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

const MultibootStarting = "Starting multiboot kernel"

var MultibootBatcher []expect.Batcher = []expect.Batcher{
	&testsuite.BExpTLog{
		L: "Matched multiboot test kernel starting",
		R: MultibootStarting,
		T: 5,
	},
	&expect.BExp{R: `"status": "ok"`},
	&expect.BExp{R: "}"},
}

func TestBootGrubMultiboot(t *testing.T) {
	e, cleanup := qemuTest(t,
		"-hda", "grub/output/images/disk.img")
	defer cleanup()

	batcher := []expect.Batcher{
		&testsuite.BExpTLog{
			L: "Matched Grub starting",
			R: "GNU GRUB  version 2.02",
			T: 50,
		},
		&expect.BExp{R: "multiboot-test-kernel"},
		&expect.BSnd{S: "v"},
		&expect.BExp{R: "\\*multiboot-test-kernel"},
		&expect.BSnd{S: "\r\n"},
	}
	batcher = append(batcher, MultibootBatcher...)
	res, err := e.ExpectBatch(batcher, 5*time.Second)
	if err != nil {
		t.Errorf("booting trough grub: %v", testsuite.DescribeBatcherErr(batcher, res, err))
	}
	testMultibootResults(t, res, wantGrubMultibootJSON)
}

func testMultibootResults(t *testing.T, res []expect.BatchRes, wantJSON []byte) {
	var output []byte
	for _, br := range res {
		output = append(output, []byte(br.Output)...)
	}

	var want multiboot.Description
	if err := json.Unmarshal(wantJSON, &want); err != nil {
		t.Fatalf("Cannot unmarshal multiboot debug information: %v", err)
	}

	i := bytes.Index(output, []byte(MultibootStarting))
	if i == -1 {
		t.Fatalf("Multiboot kernel was not executed")
	}
	output = output[i+len(MultibootStarting):]

	var got multiboot.Description
	if err := json.Unmarshal(output, &got); err != nil {
		t.Fatalf("Cannot unmarshal multiboot information from executed kernel: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Logf("var wantJSON []byte = []byte(%q)\n", output)
		t.Errorf("multiboot info failed: got\n%#v, want\n%#v", got, want)
	}
}

func TestURootBootGrubMultiboot(t *testing.T) {
	e, cleanup := qemuTest(t,
		"-kernel", "output/images/bzImage",
		"-initrd", "output/images/uroot.cpio",
		"-append", "console=ttyS0",
		"-hda", "grub/output/images/disk.img")
	defer cleanup()

	batcher := append(testsuite.Linuxboot2urootBatcher,
		[]expect.Batcher{
			&expect.BSnd{S: "boot -boot=multiboot-test-kernel\r\n"},
			&testsuite.BExpTLog{
				L: "kexec done",
				R: "kexec_core: Starting new kernel",
				T: 10,
			}}...)
	batcher = append(batcher, MultibootBatcher...)

	res, err := e.ExpectBatch(batcher, 5*time.Second)
	if err != nil {
		t.Errorf("u-root 'boot' grub config: %v", testsuite.DescribeBatcherErr(batcher, res, err))
	}
	testMultibootResults(t, res, wantUBootGrubMultibootJSON)
}

// printed by test or error, use as golden
var wantGrubMultibootJSON []byte = []byte("\n{\n\"flags\": 6767,\n\"mem_lower\": 639, \"mem_upper\": 1047424,\n\"boot_device\": 2147549183,\n\"cmdline\": \"\",\n\"mods_count\": 2, \"mods_addr\": 65692,\n\"modules\": [\n{\"start\": 1052672, \"end\": 1057988, \"cmdline\": \"foo=bar\", \"sha256\": \"7e28d3515e28dda2d9db617a10bf1843b8673a366fdd32eec3bfb5e6fe30b273\"},\n{\"start\": 1060864, \"end\": 6156848, \"cmdline\": \"\", \"sha256\": \"26c13d10d7038e70979206f3175d86bfaba3e23fc9617e6570116b796fc5cd6f\"}\n],\n\"multiboot_elf_sec\": {\"num\": 6, \"size\": 40, \"addr\": 65880, \"shndx\": 5},\n\"mmap_addr\": 65736, \"mmap_length\": 144,\n\"mmap\": [\n{\"size\": 20, \"base_addr\": \"0x000000000\", \"length\": \"0x00009fc00\", \"type\": 1},\n{\"size\": 20, \"base_addr\": \"0x00009fc00\", \"length\": \"0x000000400\", \"type\": 2},\n{\"size\": 20, \"base_addr\": \"0x0000f0000\", \"length\": \"0x000010000\", \"type\": 2},\n{\"size\": 20, \"base_addr\": \"0x000100000\", \"length\": \"0x03fee0000\", \"type\": 1},\n{\"size\": 20, \"base_addr\": \"0x03ffe0000\", \"length\": \"0x000020000\", \"type\": 2},\n{\"size\": 20, \"base_addr\": \"0x0fffc0000\", \"length\": \"0x000040000\", \"type\": 2}\n],\n\"bootloader\": \"GRUB 2.02\",\n\"status\": \"ok\"\n}")
var wantUBootGrubMultibootJSON []byte = []byte("\n{\n\"flags\": 589,\n\"mem_lower\": 639, \"mem_upper\": 1047424,\n\"cmdline\": \"\",\n\"mods_count\": 2, \"mods_addr\": 6164480,\n\"modules\": [\n{\"start\": 1056768, \"end\": 1062084, \"cmdline\": \"foo=bar\", \"sha256\": \"7e28d3515e28dda2d9db617a10bf1843b8673a366fdd32eec3bfb5e6fe30b273\"},\n{\"start\": 1064960, \"end\": 6160944, \"cmdline\": \"\", \"sha256\": \"26c13d10d7038e70979206f3175d86bfaba3e23fc9617e6570116b796fc5cd6f\"}\n],\n\"mmap_addr\": 1048576, \"mmap_length\": 144,\n\"mmap\": [\n{\"size\": 20, \"base_addr\": \"0x000000000\", \"length\": \"0x00009fc00\", \"type\": 1},\n{\"size\": 20, \"base_addr\": \"0x00009fc00\", \"length\": \"0x000000400\", \"type\": 2},\n{\"size\": 20, \"base_addr\": \"0x0000f0000\", \"length\": \"0x000010000\", \"type\": 2},\n{\"size\": 20, \"base_addr\": \"0x000100000\", \"length\": \"0x03fee0000\", \"type\": 1},\n{\"size\": 20, \"base_addr\": \"0x03ffe0000\", \"length\": \"0x000020000\", \"type\": 2},\n{\"size\": 20, \"base_addr\": \"0x0fffc0000\", \"length\": \"0x000040000\", \"type\": 2}\n],\n\"bootloader\": \"u-root kexec\",\n\"status\": \"ok\"\n}")
