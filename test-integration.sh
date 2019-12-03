#!/bin/bash
bz=$(realpath output/images/bzImage)
mbd=$(realpath grub/output/target/boot/)
cd $GOPATH/src/github.com/u-root/u-root/integration
UROOT_QEMU_TIMEOUT_X=4 UROOT_QEMU="qemu-system-x86_64 -m 1024" UROOT_KERNEL="$bz" UROOT_MULTIBOOT_TEST_KERNEL_DIR="$mbd" go test "$@"
