#!/bin/sh

set -e

u-root -build=bb -o $BINARIES_DIR/uroot.cpio -files $TARGET_DIR/lib/modules:lib/modules core boot github.com/u-root/u-root/cmds/exp/modprobe

