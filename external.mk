# Makefiles used by all subprojects

include $(sort $(wildcard $(BR2_EXTERNAL)/package/*/*.mk))
include $(BR2_EXTERNAL)/busybox_extra_getty.mk
