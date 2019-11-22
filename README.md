# Buildroot For Testing bootloaders

This project holds several buildroot configurations to test bootloaders
It should be cloned with the buildroot folder as git submodule (ie use recursive)

It currently covers:
 - grub
 - syslinux on cd aka isolinux
 - linuxboot u-root 'boot', 'boot2' and 'localboot' commands on a grub
   installed partition
 - linuxboot u-root 'boot', 'boot2' and 'localboot' commands on an isolinux cd

## Test

```shell
./test.sh
```

## Building

### Requirements

#### buildroot

On debian
```shell
sudo apt install build-essential bc
```

Else see [requirements](https://buildroot.org/downloads/manual/manual.html#requirement).

#### u-root

Make sure your Go version is 1.12. Make sure your `GOPATH` is set up correctly.

Download and install u-root:

```shell
go get github.com/u-root/u-root
```


### Build

```shell
make

make PROJECT_NAME=grub
make PROJECT_NAME=syslinux
```

## (buildroot-submodule) Licence

buildroot-submodule is provided under the GPLv3 or later. The licence is provided in the _LICENCE_ file. Note that this licence only covers the files provided by buildroot-submodule. It does not cover buildroot (which is GPLv2 or later) nor any software installed by buildroot (they have their own licences) nor your own code (which you are free to licence as you want).
