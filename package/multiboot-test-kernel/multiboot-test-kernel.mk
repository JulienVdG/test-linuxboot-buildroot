################################################################################
#
# multiboot-test-kernel
#
################################################################################

MULTIBOOT_TEST_KERNEL_VERSION = 8d785e22f07cef36cf02dc173ba4f6d7901f7ce9
MULTIBOOT_TEST_KERNEL_SITE = $(call github,u-root,multiboot-test-kernel,$(MULTIBOOT_TEST_KERNEL_VERSION))
MULTIBOOT_TEST_KERNEL_LICENSE = GPL-3.0+
MULTIBOOT_TEST_KERNEL_FILES = LICENSE

define MULTIBOOT_TEST_KERNEL_BUILD_CMDS
    # redefine CFLAGS and LDFLAGS as cmdline overrides the Makefile content
    $(MAKE) $(TARGET_CONFIGURE_OPTS) CFLAGS="-c -m32" LDFLAGS="-m elf_i386" -C $(@D) all
    gzip -k $(@D)/kernel
endef

define MULTIBOOT_TEST_KERNEL_INSTALL_TARGET_CMDS
    $(INSTALL) -D -m 0755 $(@D)/kernel $(TARGET_DIR)/boot
    $(INSTALL) -D -m 0755 $(@D)/kernel.gz $(TARGET_DIR)/boot
endef

MULTIBOOT_TEST_KERNEL_BIN_ARCH_EXCLUDE = "/boot"

$(eval $(generic-package))
