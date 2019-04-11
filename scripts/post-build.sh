#!/bin/sh

set -e

cp -f "$BR2_EXTERNAL_BUILDROOT_SUBMODULE_PATH/conf/grub-bios.cfg" "$TARGET_DIR/boot/grub/grub.cfg"

# Copy grub 1st stage to binaries, required for genimage
cp -f "$HOST_DIR/lib/grub/i386-pc/boot.img" "$BINARIES_DIR"
