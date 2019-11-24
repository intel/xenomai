/******************************************************************************
# Copyright (c) 2019 Intel Corporation
# Authored-by:  Fino Meng <fino.meng@intel.com>
# Licensed under the MIT license. See LICENSE.MIT in the project root.
******************************************************************************/
package main

const ecs_b_1_0 = "ecs_b_1.0"
const ecs_b_2_0 = "ecs_b_2.0"

var manifestNames = []string{ecs_b_1_0, ecs_b_2_0}

// branch is mandatory, tag and commit-id is optional;
// it's due to a tag and a branch may set to same name accidentally.
var manifests = map[string][]typeRepoGit{

	ecs_b_1_0: manifest_ecs_b_1_0,
	ecs_b_2_0: manifest_ecs_b_2_0,
}

var bblayersConfTemplates = map[string]string{

	ecs_b_1_0: bblayers_conf_ecs_b_1_0_template,
	ecs_b_2_0: bblayers_conf_ecs_b_2_0_template,
}

var localConfTemplates = map[string]string{

	ecs_b_1_0: local_conf_ecs_b_1_0_template,
	ecs_b_2_0: local_conf_ecs_b_2_0_template,
}
