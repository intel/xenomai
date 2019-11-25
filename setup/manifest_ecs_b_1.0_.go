/******************************************************************************
# Copyright (c) 2019 Intel Corporation
# Authored-by:  Fino Meng <fino.meng@intel.com>
# Licensed under the MIT license. See LICENSE.MIT in the project root.
******************************************************************************/
package main

var manifest_ecs_b_1_0 = []typeRepoGit{
	typeRepoGit{
		Name:         "poky",
		SrcURI:       "git://git.yoctoproject.org/poky",
		Branch:       "sumo",
		LastRevision: "d240b885f26e9b05c8db0364ab2ace9796709aad",
	},
	typeRepoGit{
		Name:         "meta-openembedded",
		SrcURI:       "git://git.openembedded.org/meta-openembedded",
		Branch:       "sumo",
		LastRevision: "a19aa29f7fa336cd075b72c496fe1102e6e5422b",
	},
	typeRepoGit{
		Name:         "meta-intel",
		SrcURI:       "git://git.yoctoproject.org/meta-intel",
		Branch:       "sumo",
		LastRevision: "aa8f5fad12ed5cda9b3779dd380e29755ab24c79",
	},
	typeRepoGit{
		Name:         "meta-intel-iotg-bsp",
		LocalRepoDir: "meta-intel-iotg",
		SrcURI:       "https://github.com/intel/iotg-yocto-bsp-public.git",
		Branch:       "e3900/master",
		LastRevision: "a9e85b077160bfa6156ca266599efa9c3eaf5837",
	},
	typeRepoGit{
		Name:         "meta-measured",
		SrcURI:       "https://github.com/flihp/meta-measured.git",
		Branch:       "rocko",
		LastRevision: "b7d667e4796623812ab22fdde99d7ce68622de19",
	},
	typeRepoGit{
		Name:   "meta-xenomai",
		SrcURI: "https://github.com/intel/meta-xenomai.git",
		Branch: "4.14.68/base/3.0.7",
	},
}

///////////////////////////////////////////////////////////////////////////////

const bblayers_conf_ecs_b_1_0_template = `
# POKY_BBLAYERS_CONF_VERSION is increased each time build/conf/bblayers.conf
# changes incompatibly
POKY_BBLAYERS_CONF_VERSION = "2"

BBPATH = "${TOPDIR}"
BBFILES ?= ""

BBLAYERS ?= " \
  {{.SnapshotDir}}poky/meta \
  {{.SnapshotDir}}poky/meta-poky \
  {{.SnapshotDir}}poky/meta-yocto-bsp \
  {{.SnapshotDir}}meta-openembedded/meta-filesystems \
  {{.SnapshotDir}}meta-openembedded/meta-networking \
  {{.SnapshotDir}}meta-openembedded/meta-oe \
  {{.SnapshotDir}}meta-openembedded/meta-python \
  {{.SnapshotDir}}meta-openembedded/meta-multimedia \
  {{.SnapshotDir}}meta-openembedded/meta-gnome \
  {{.SnapshotDir}}meta-openembedded/meta-xfce \
  {{.SnapshotDir}}meta-intel \
  {{.SnapshotDir}}meta-measured \
  {{.SnapshotDir}}meta-intel-iotg/meta-intel-middleware \
  {{.SnapshotDir}}meta-xenomai \
  "
`

///////////////////////////////////////////////////////////////////////////////

const local_conf_ecs_b_1_0_template = `
#
# This file is your local configuration file and is where all local user settings
# are placed. The comments in this file give some guide to the options a new user
# to the system might want to change but pretty much any configuration option can
# be set in this file. More adventurous users can look at local.conf.extended
# which contains other examples of configuration which can be placed in this file
# but new users likely won't need any of them initially.
#
# Lines starting with the '#' character are commented out and in some cases the
# default values are provided as comments to show people example syntax. Enabling
# the option is a question of removing the # character and making any change to the
# variable as required.

#
# Machine Selection
#
MACHINE = "intel-corei7-64"

# Set preferred version for kernel
PREFERRED_PROVIDER_virtual/kernel = "linux-intel"
PREFERRED_VERSION_linux-intel-rt = "4.14%"

# Add serial console and boot parameters for RT
APPEND += "3 \
scsi_mod.scan=async \
reboot=efi \
console=ttyS2,115200n8 \
processor.max_cstate=0 \
intel.max_cstate=0 \
processor_idle.max_cstate=0 \
intel_idle.max_cstate=0 \
clocksource=tsc \
tsc=reliable \
nmi_watchdog=0 \
nosoftlockup \
idle=poll \
noht \
isolcpus=1-3 \
rcu_nocbs=1-3 \
nohz_full=1-3 \
i915.enable_guc_loading=1 \
i915.enable_guc_submission=1 \
i915.enable_rc6=0 \
i915.enable_dc=0 \
i915.disable_power_well=0 \
"

# Take note that as we are building 3rd party ingredient.
# We need the LICENSE_FLAGS below.
LICENSE_FLAGS_WHITELIST += "commercial"

# Enable the security flags as per SDL 4.0 requirement
require conf/distro/include/security_flags.inc
SECURITY_CFLAGS = "-fstack-protector-strong -pie -fpie -D_FORTIFY_SOURCE=2 -O2 -Wformat -Wformat-security"
SECURITY_NO_PIE_CFLAGS = "-fstack-protector-strong -D_FORTIFY_SOURCE=2 -O2 -Wformat -Wformat-security"
SECURITY_LDFLAGS = "-Wl,-z,relro,-z,now,-z,noexecstack"

# GCC Sanitizers does not compiles with PIE
SECURITY_CFLAGS_pn-gcc-sanitizers = "${SECURITY_NO_PIE_CFLAGS}"

#
# User Space Configuration Override for Intel Platforms
#
# Selection of jpeg package provider
PREFERRED_PROVIDER_jpeg = "jpeg"
PREFERRED_PROVIDER_jpeg-native = "jpeg-native"

# Include rt test in image
IMAGE_INSTALL_append = " rt-tests-ptest"

# Add WKS file with updated boot parameters
WKS_FILE = "systemd-bootdisk-microcode-custom.wks"

#Exclude piglit, ltp and mesa-demos packages
PACKAGE_EXCLUDE = "packagegroup-core-apl-extra"

# Use samba
DISTRO_FEATURES_append = " pam"

# Enable xserver-xorg
IMAGE_INSTALL_append = " xserver-xorg"

# Enable vp8dec for gstreamer1.0-plugins-good
PACKAGECONFIG_append_pn-gstreamer1.0-plugins-good = "vpx"

# Multi-libraries support is by default "OFF"
# Please uncomment the 4 lines below to enable multilib support.
#require conf/multilib.conf
#DEFAULTTUNE = "corei7-64"
#MULTILIBS = "multilib:lib32"
#DEFAULTTUNE_virtclass-multilib-lib32 = "corei7-32"

# Install autoconf-archive
IMAGE_INSTALL_append = " autoconf-archive"

# Install libva
IMAGE_INSTALL_append = " libva"

# Install Wayland in image
DISTRO_FEATURES_append = " wayland pam"
CORE_IMAGE_EXTRA_INSTALL += "wayland weston weston-examples"

# Install mesa glxinfo
IMAGE_INSTALL_append = " mesa-glxinfo"

# Install USB-modeswitch and USB-modeswitch-data in image
IMAGE_INSTALL_append = " usb-modeswitch usb-modeswitch-data"

# Install JHI
IMAGE_INSTALL_append = " jhi"

# Install IQV driver
IMAGE_INSTALL_append = " iqvlinux"

# Disable lttng modules
LTTNGMODULES_corei7-64-intel-common = ""

# Install xinitrc environment file
IMAGE_INSTALL_append = " xinit-env"

# By default, we want our OS to includes all kernel modules.
IMAGE_INSTALL_append = " kernel-modules"

# Use systemd init instead of sysV init
DISTRO_FEATURES_append = " systemd"
VIRTUAL-RUNTIME_init_manager = "systemd"
DISTRO_FEATURES_BACKFILL_CONSIDERED = "sysvinit"

#
# Where to place downloads
#
# During a first build the system will download many different source code tarballs
# from various upstream projects. This can take a while, particularly if your network
# connection is slow. These are all stored in DL_DIR. When wiping and rebuilding you
# can preserve this directory to speed up this part of subsequent builds. This directory
# is safe to share between multiple builds on the same machine too.
#
# The default is a downloads directory under TOPDIR which is the build directory.
#
#DL_DIR ?= "${TOPDIR}/downloads"

#
# Where to place shared-state files
#
# BitBake has the capability to accelerate builds based on previously built output.
# This is done using "shared state" files which can be thought of as cache objects
# and this option determines where those files are placed.
#
# You can wipe out TMPDIR leaving this directory intact and the build would regenerate
# from these files if no changes were made to the configuration. If changes were made
# to the configuration, only shared state files where the state was still valid would
# be used (done using checksums).
#
# The default is a sstate-cache directory under TOPDIR.
#
#SSTATE_DIR ?= "${TOPDIR}/sstate-cache"

#
# Where to place the build output
#
# This option specifies where the bulk of the building work should be done and
# where BitBake should place its temporary files and output. Keep in mind that
# this includes the extraction and compilation of many applications and the toolchain
# which can use Gigabytes of hard disk space.
#
# The default is a tmp directory under TOPDIR.
#
#TMPDIR = "${TOPDIR}/tmp"

#
# Default policy config
#
# The distribution setting controls which policy settings are used as defaults.
# The default value is fine for general Yocto project use, at least initially.
# Ultimately when creating custom policy, people will likely end up subclassing 
# these defaults.
#
DISTRO ?= "poky"
# As an example of a subclass there is a "bleeding" edge policy configuration
# where many versions are set to the absolute latest code from the upstream 
# source control systems. This is just mentioned here as an example, its not
# useful to most new users.
# DISTRO ?= "poky-bleeding"

#
# Package Management configuration
#
# This variable lists which packaging formats to enable. Multiple package backends
# can be enabled at once and the first item listed in the variable will be used
# to generate the root filesystems.
# Options are:
#  - 'package_deb' for debian style deb files
#  - 'package_ipk' for ipk files are used by opkg (a debian style embedded package manager)
#  - 'package_rpm' for rpm style packages
# E.g.: PACKAGE_CLASSES ?= "package_rpm package_deb package_ipk"
# We default to rpm:
PACKAGE_CLASSES ?= "package_rpm package_deb"

#
# SDK target architecture
#
# This variable specifies the architecture to build SDK items for and means
# you can build the SDK packages for architectures other than the machine you are
# running the build on (i.e. building i686 packages on an x86_64 host).
# Supported values are i686 and x86_64
#SDKMACHINE ?= "i686"

#
# Extra image configuration defaults
#
# The EXTRA_IMAGE_FEATURES variable allows extra packages to be added to the generated
# images. Some of these options are added to certain image types automatically. The
# variable can contain the following options:
#  "dbg-pkgs"       - add -dbg packages for all installed packages
#                     (adds symbol information for debugging/profiling)
#  "dev-pkgs"       - add -dev packages for all installed packages
#                     (useful if you want to develop against libs in the image)
#  "ptest-pkgs"     - add -ptest packages for all ptest-enabled packages
#                     (useful if you want to run the package test suites)
#  "tools-sdk"      - add development tools (gcc, make, pkgconfig etc.)
#  "tools-debug"    - add debugging tools (gdb, strace)
#  "eclipse-debug"  - add Eclipse remote debugging support
#  "tools-profile"  - add profiling tools (oprofile, lttng, valgrind)
#  "tools-testapps" - add useful testing tools (ts_print, aplay, arecord etc.)
#  "debug-tweaks"   - make an image suitable for development
#                     e.g. ssh root access has a blank password
# There are other application targets that can be used here too, see
# meta/classes/image.bbclass and meta/classes/core-image.bbclass for more details.
# We default to enabling the debugging tweaks.
EXTRA_IMAGE_FEATURES ?= "debug-tweaks"

#
# Additional image features
#
# The following is a list of additional classes to use when building images which
# enable extra features. Some available options which can be included in this variable
# are:
#   - 'buildstats' collect build statistics
#   - 'image-mklibs' to reduce shared library files size for an image
#   - 'image-prelink' in order to prelink the filesystem image
# NOTE: if listing mklibs & prelink both, then make sure mklibs is before prelink
# NOTE: mklibs also needs to be explicitly enabled for a given image, see local.conf.extended
USER_CLASSES ?= "buildstats image-mklibs image-prelink"

#
# Runtime testing of images
#
# The build system can test booting virtual machine images under qemu (an emulator)
# after any root filesystems are created and run tests against those images. To
# enable this uncomment this line. See classes/testimage(-auto).bbclass for
# further details.
#TEST_IMAGE = "1"
#
# Interactive shell configuration
#
# Under certain circumstances the system may need input from you and to do this it
# can launch an interactive shell. It needs to do this since the build is
# multithreaded and needs to be able to handle the case where more than one parallel
# process may require the user's attention. The default is iterate over the available
# terminal types to find one that works.
#
# Examples of the occasions this may happen are when resolving patches which cannot
# be applied, to use the devshell or the kernel menuconfig
#
# Supported values are auto, gnome, xfce, rxvt, screen, konsole (KDE 3.x only), none
# Note: currently, Konsole support only works for KDE 3.x due to the way
# newer Konsole versions behave
#OE_TERMINAL = "auto"
# By default disable interactive patch resolution (tasks will just fail instead):
PATCHRESOLVE = "noop"

#
# Disk Space Monitoring during the build
#
# Monitor the disk space during the build. If there is less that 1GB of space or less
# than 100K inodes in any key build location (TMPDIR, DL_DIR, SSTATE_DIR), gracefully
# shutdown the build. If there is less that 100MB or 1K inodes, perform a hard abort
# of the build. The reason for this is that running completely out of space can corrupt
# files and damages the build in ways which may not be easily recoverable.
# It's necesary to monitor /tmp, if there is no space left the build will fail
# with very exotic errors.
BB_DISKMON_DIRS ??= "\
    STOPTASKS,${TMPDIR},1G,100K \
    STOPTASKS,${DL_DIR},1G,100K \
    STOPTASKS,${SSTATE_DIR},1G,100K \
    STOPTASKS,/tmp,100M,100K \
    ABORT,${TMPDIR},100M,1K \
    ABORT,${DL_DIR},100M,1K \
    ABORT,${SSTATE_DIR},100M,1K \
    ABORT,/tmp,10M,1K"

#
# Shared-state files from other locations
#
# As mentioned above, shared state files are prebuilt cache data objects which can
# used to accelerate build time. This variable can be used to configure the system
# to search other mirror locations for these objects before it builds the data itself.
#
# This can be a filesystem directory, or a remote url such as http or ftp. These
# would contain the sstate-cache results from previous builds (possibly from other
# machines). This variable works like fetcher MIRRORS/PREMIRRORS and points to the
# cache locations to check for the shared objects.
# NOTE: if the mirror uses the same structure as SSTATE_DIR, you need to add PATH
# at the end as shown in the examples below. This will be substituted with the
# correct path within the directory structure.
#SSTATE_MIRRORS ?= "\
#file://.* http://someserver.tld/share/sstate/PATH;downloadfilename=PATH \n \
#file://.* file:///some/local/dir/sstate/PATH"

#
# Yocto Project SState Mirror
#
# The Yocto Project has prebuilt artefacts available for its releases, you can enable
# use of these by uncommenting the following line. This will mean the build uses
# the network to check for artefacts at the start of builds, which does slow it down
# equally, it will also speed up the builds by not having to build things if they are
# present in the cache. It assumes you can download something faster than you can build it
# which will depend on your network.
#
#SSTATE_MIRRORS ?= "file://.* http://sstate.yoctoproject.org/2.5/PATH;downloadfilename=PATH"

#
# Qemu configuration
#
# By default qemu will build with a builtin VNC server where graphical output can be
# seen. The two lines below enable the SDL backend too. By default libsdl-native will
# be built, if you want to use your host's libSDL instead of the minimal libsdl built
# by libsdl-native then uncomment the ASSUME_PROVIDED line below.
PACKAGECONFIG_append_pn-qemu-native = " sdl"
PACKAGECONFIG_append_pn-nativesdk-qemu = " sdl"
#ASSUME_PROVIDED += "libsdl-native"

# CONF_VERSION is increased each time build/conf/ changes incompatibly and is used to
# track the version of this file when it was generated. This can safely be ignored if
# this doesn't mean anything to you.
CONF_VERSION = "1"
`
