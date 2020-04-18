/******************************************************************************
# Copyright (c) 2019 Intel Corporation
# Authored-by:  Fino Meng <fino.meng@intel.com>
# Licensed under the MIT license. See LICENSE.MIT in the project root.
******************************************************************************/
package main

const key_1 = "manifest_1"
const key_2 = "manifest_2"

var manifestNames = []string{key_1, key_2}

// branch is mandatory, tag and commit-id is optional;
// it's due to a tag and a branch may set to same name accidentally.
var manifests = map[string][]typeRepoGit{

	key_1: manifest_1,
	key_2: manifest_2,
}

var bblayersConfTemplates = map[string]string{

	key_1: bblayers_conf_manifest_1_template,
	key_2: bblayers_conf_manifest_2_template,
}

var localConfTemplates = map[string]string{

	key_1: local_conf_manifest_1_template,
	key_2: local_conf_manifest_2_template,
}
