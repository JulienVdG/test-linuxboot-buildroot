# Buildroot For Testing bootloaders

This project holds several buildroot configurations to test bootloaders
It should be cloned with the buildroot folder as git submodule (ie use recursive)

## Test

TODO go tests with tastevin ;)

## Building

### Requirements

On debian
```
sudo apt install build-essential bc
```

Else see [requirements](https://buildroot.org/downloads/manual/manual.html#requirement).

### Build
```
make
```

The images will be in `output/images`
 - `disk.img` can be dd'ed to an usbstick
 - `bzImage` and `rootfs.cpio.xz` can be used for tftp

## (buildroot-submodule) Licence

buildroot-submodule is provided under the GPLv3 or later. The licence is provided in the _LICENCE_ file. Note that this licence only covers the files provided by buildroot-submodule. It does not cover buildroot (which is GPLv2 or later) nor any software installed by buildroot (they have their own licences) nor your own code (which you are free to licence as you want).
