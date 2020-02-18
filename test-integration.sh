#!/bin/bash
bz=$(realpath output/images/bzImage)
mbd=$(realpath grub/output/target/boot/)
cd $GOPATH/src/github.com/u-root/u-root/integration
export UROOT_QEMU_TIMEOUT_X=4
export UROOT_QEMU="qemu-system-x86_64 -m 1024"
export UROOT_KERNEL="$bz"
export UROOT_MULTIBOOT_TEST_KERNEL_DIR="$mbd"

go test "$@" ./...
