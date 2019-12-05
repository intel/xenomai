This project contains the scripts to setup a Xenomai(xenomai.org) powered real-time co-kernel Linux distribution.

## Overview: 

The entire project is based on Yocto, Xenomai's code and ipipe pathces are organized into meta-xenomai. 'download.sh' will download a code snapshot of Yocto poky and other meta layers, based on the branch&revision in the manifest which you selected.

## Build guide:

- ### Prerequisite:
The scripts are written in Go, please setup Golang environment on your host machine:
https://golang.org/doc/install  or  https://golang.google.cn/doc/install  for China.
- ### download:
```
$./download.sh
Pls select a snapshot of code to download:
[1] ecs_b_1.0
[2] ecs_b_2.0
```
Type "1" to choose "ecs_b_1.0", it combines kernel 4.14.68 and xenomai 3.0.7; detailed branch and revision pls see: setup/manifest_ecs_b_1.0_.go;
depends on your network, download code snapshot will take some time.
- ### build:
```
$cd snapshots/ecs_b_1.0/
$source poky/oe-init-build-env build  ### will jump to build/ automaticly
$bitbake -k core-image-xfce-sdk
# or target without build facilities:
$bitbake -k core-image-xfce
# or kernel only:
$bitbake virtual/kernel
```
- ### flash image:
Output images under:  build/tmp/deploy/images/intel-corei7-64/

Liveboot USB disk creation for UEFI BIOS:
```
   $ sudo dd if=core-image-xfce-sdk-intel-corei7-64.rootfs.wic of=/dev/sdb bs=4M status=progress conv=fdatasync
```
Liveboot USB disk creation for legacy BIOS:
```
   $ sudo dd if=core-image-xfce-sdk-intel-corei7-64.hddimg of=/dev/sdb bs=4M status=progress conv=fdatasync
```