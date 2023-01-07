DISCONTINUATION OF PROJECT

This project will no longer be maintained by Intel.

Intel has ceased development and contributions including, but not limited to, maintenance, bug fixes, new releases, or updates, to this project.  

Intel no longer accepts patches to this project.

If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the open source software community, please create your own fork of this project.  

Contact: webadmin@linux.intel.com
# Xenomai

This project contains the scripts to setup a Xenomai(xenomai.org) powered real-time co-kernel Linux distribution in Yocto/Bitbake way, it has fewest code for education purpose but not for production.

## Overview: 

The entire project is based on Yocto, Xenomai's code and ipipe pathces are organized into meta-xenomai-intel. 'download.sh' will download a code snapshot of Yocto poky and other meta layers, based on the branch&revision in the manifest which you selected.

## Build guide:

- ### Prerequisite:
The scripts are written in Go, please setup Golang environment on your host machine:
https://golang.org/doc/install  or  https://golang.google.cn/doc/install  for China.
- ### Download:
```
$./download.sh
Pls select a snapshot of code to download:
[1] manifest_1
[2] manifest_2
```
Type "1" to choose "manifest_1", it combines kernel 4.14.68 and xenomai 3.0.7; detailed branch and revision pls see: setup/manifest_1_.go;
depends on your network, download code snapshot will take some time.
- ### Build:
```
$cd snapshots/manifest_1/
$source poky/oe-init-build-env build  ### will jump to build/ automaticly
$bitbake -k core-image-xfce-sdk
# or target without build facilities:
$bitbake -k core-image-xfce
# or kernel only:
$bitbake virtual/kernel
```
- ### Make bootable USB disk:
Output images under:  build/tmp/deploy/images/intel-corei7-64/

Assume a USB disk is plugged in and enum as: /dev/sdb

Liveboot USB disk creation for UEFI BIOS:
```
   $ sudo dd if=core-image-xfce-sdk-intel-corei7-64.rootfs.wic of=/dev/sdb bs=4M status=progress oflag=sync
```
Liveboot USB disk creation for legacy BIOS:
```
   $ sudo dd if=core-image-xfce-sdk-intel-corei7-64.hddimg of=/dev/sdb bs=4M status=progress oflag=sync
```
