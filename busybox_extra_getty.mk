# Rules to extend the busybox getty setup in inittab

# List of extra configuration (name used as both comment and variable prefix)
BR2_TARGET_EXTRA_GETTY = VGA IPMI

# Variables similar to SYSTEM_GETTY_* for each extra conf
VGA_GETTY_PORT = tty1
VGA_GETTY_BAUDRATE = 0
VGA_GETTY_TERM = vt100
VGA_GETTY_OPTIONS = 

IPMI_GETTY_PORT = ttyS1
IPMI_GETTY_BAUDRATE = 0
IPMI_GETTY_TERM = vt100
IPMI_GETTY_OPTIONS = 

# The code

ifneq ($(BR2_TARGET_EXTRA_GETTY),)
define BUSYBOX_SET_EXTRA_GETTY_template
	grep -q '# $(1)$$' $(TARGET_DIR)/etc/inittab || \
	$(SED) '/# GENERIC_SERIAL$$/a# $(1)' \
		$(TARGET_DIR)/etc/inittab
	$(SED) '/# $(1)$$/s~^.*#~$($(1)_GETTY_PORT)::respawn:/sbin/getty -L $($(1)_GETTY_OPTIONS) $($(1)_GETTY_PORT) $($(1)_GETTY_BAUDRATE) $($(1)_GETTY_TERM) #~' \
		$(TARGET_DIR)/etc/inittab
endef
BUSYBOX_SET_EXTRA_GETTY = $(foreach name,$(BR2_TARGET_EXTRA_GETTY),$(call BUSYBOX_SET_EXTRA_GETTY_template,$(name))$(sep))
BUSYBOX_TARGET_FINALIZE_HOOKS += BUSYBOX_SET_EXTRA_GETTY
endif # BR2_TARGET_EXTRA_GETTY

