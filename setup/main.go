/******************************************************************************
# Copyright (c) 2019 Intel Corporation
# Authored-by:  Fino Meng <fino.meng@intel.com>
# Licensed under the MIT license. See LICENSE.MIT in the project root.
******************************************************************************/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"text/template"
	"time"
)

// G collecting all global var, no exception.
var G = struct {
	DirDownload string
	DirBuild    string // inside each code snapshot
	DirConf     string
	PWD         string
}{
	"snapshots/",
	"build/",
	"build/conf/",
	"",
}

func printHelp() {
	fmt.Printf("\nusage:\n")
	fmt.Printf("	supported cmd: down, check, pull\n\n")
}

///////////////////////////////////////////////////////////////////////////////
func main() {

	flag.Parse()
	lenOfArgs := len(flag.Args())
	flagCmd := ""
	if lenOfArgs == 0 {
		printHelp()
		os.Exit(0)
	} else {
		flagCmd = flag.Arg(0)
	}

	G.PWD, _ = os.Getwd()

	switch flagCmd {
	case "down":

		idxSelect := 0
		fmt.Printf("Pls select a snapshot of code to download:\n")
		for k, v := range manifestNames {
			fmt.Printf("[%d] %s\n", k+1, v)
		}
		fmt.Scanln(&idxSelect)
		dirSnapshot := G.PWD + "/" + G.DirDownload + manifestNames[idxSelect-1] + "/"
		fmt.Printf("\nSelected code snapshot: %s\n\n", dirSnapshot)
		time.Sleep(1)

		generateConf(dirSnapshot, manifestNames[idxSelect-1])
		downloadRepos(dirSnapshot, manifestNames[idxSelect-1])

	case "check":
		//compare commit diff between local branch and upsteam branch

	case "pull":

	default:
		println("unknown cmd.")
		os.Exit(0)
	}

}

///////////////////////////////////////////////////////////////////////////////

type typeRepoGit struct {
	Name         string
	SrcURI       string
	Branch       string
	Tag          string
	LastRevision string
	LocalRepoDir string
	ReDown       bool // record user selection if re-download the repo
}

func execCommand(commandName, runDir string, params []string) bool {

	cmd := exec.Command(commandName, params...)
	cmd.Dir = runDir
	fmt.Println(cmd.Args)
	fmt.Println("Run dir: ", runDir)
	//show stdout&stderr instantly
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return true
}

func worker(wg *sync.WaitGroup, id int, name string, commandName, runDir string, params []string) {
	defer wg.Done()

	fmt.Printf("\nWorker %v %v: Started\n", id, name)
	execCommand(commandName, runDir, params)
	fmt.Printf("\nWorker %v %v: Finished\n", id, name)
}

type typeBBlayersConfTemplate struct {
	SnapshotDir string
}

type typeLocalConfTemplate struct {
	SnapshotDir string
}

func generateConf(dirSnapshot string, manifest string) {

	dirConf := dirSnapshot + G.DirConf

	if _, err := os.Stat(dirConf); !os.IsNotExist(err) {
		// path/to/whatever already exist
		println("conf dir " + dirConf + " already exist.")
		// to do: press y for delete, N for keep
		fmt.Printf("do you want to delete existing conf dir and re-generate it? [y/N]\n")
		var input string
		fmt.Scanln(&input)
		if input == "y" {
			os.RemoveAll(dirConf)
			time.Sleep(3)
			println(dirConf + " deleted.\n")
		} else {
			println(dirConf + " keeped.\n")
			return
		}
	}
	os.MkdirAll(dirConf, os.ModePerm)
	fmt.Printf("created conf dir for bitbake project: %s\n", dirConf)

	filenameBBlayers := dirConf + "/bblayers.conf"
	fileBBlayers, err := os.Create(filenameBBlayers)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(" Write to file : " + filenameBBlayers)
	bblayers := typeBBlayersConfTemplate{SnapshotDir: dirSnapshot}
	t := template.Must(template.New("bblayers").Parse(bblayersConfTemplates[manifest]))

	//err = t.Execute(os.Stdout, bblayers)
	err = t.Execute(fileBBlayers, bblayers)
	if err != nil {
		log.Println("executing template:", err)
	}

	fileBBlayers.Close()

	filenameLocalConf := dirConf + "/local.conf"
	fileLocalConf, err := os.Create(filenameLocalConf)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(" Write to file : " + filenameLocalConf)
	localConf := typeLocalConfTemplate{SnapshotDir: dirSnapshot}
	t = template.Must(template.New("localConf").Parse(localConfTemplates[manifest]))

	//err = t.Execute(os.Stdout, bblayers)
	err = t.Execute(fileLocalConf, localConf)
	if err != nil {
		log.Println("executing template:", err)
	}

	fileLocalConf.Close()

	println("")
}

// [branch, tag, commit-id] can have multi combinations
// same tag name can exist in multi branch
// dirDownload should be created by generateConf() already.
func downloadRepos(dirDownload string, manifest string) {

	mani := manifests[manifest]
	amountOfRepos := len(mani)

	for i := 0; i < amountOfRepos; i++ {

		var localRepoDir string = dirDownload + mani[i].Name
		if mani[i].LocalRepoDir != "" {
			localRepoDir = dirDownload + mani[i].LocalRepoDir
		}

		if _, err := os.Stat(localRepoDir); !os.IsNotExist(err) {
			println("repo " + localRepoDir + " already exist.")
			fmt.Println("do you want to delete existing repo and re-download again? [y/N] ")
			var input string
			fmt.Scanln(&input)
			if input == "y" {
				os.RemoveAll(localRepoDir)
				time.Sleep(3)
				println("repo " + localRepoDir + " deleted.\n")
				mani[i].ReDown = true
			} else {
				println("repo " + localRepoDir + " keeped.\n")
				mani[i].ReDown = false
				continue
			}
		} else {
			println("repo " + localRepoDir + " do not exist, will download it.\n")
			mani[i].ReDown = true // folder not exist
		}

	}

	var wg sync.WaitGroup

	fmt.Println("Main: start downloading repos...\n")
	time.Sleep(3)
	timeStart := time.Now().UTC()

	for i := 0; i < amountOfRepos; i++ {
		if mani[i].ReDown == true {
			var gitParams []string
			var localRepoDir string = dirDownload + mani[i].Name
			if mani[i].LocalRepoDir != "" {
				localRepoDir = dirDownload + mani[i].LocalRepoDir
			}

			if mani[i].Branch != "" {
				gitParams = []string{"clone", "-b", mani[i].Branch, mani[i].SrcURI, localRepoDir}
			} else {
				fmt.Println("Branch name is mandatory for a git repo, pls add branch in manifest.")
				os.Exit(1)
			}

			wg.Add(1)
			dir, _ := os.Getwd()
			go worker(&wg, i, mani[i].Name, "git", dir, gitParams)
		}
	}

	wg.Wait() //wait for all the goroutine return
	timeEnd := time.Now().UTC()
	var dura time.Duration = timeEnd.Sub(timeStart)
	dura = dura.Round(time.Second)

	fmt.Printf("Main: download completed, cost time: %v\n\n", dura)
	fmt.Println("Main: checking out by tag or commit-id.\n")

	for i := 0; i < amountOfRepos; i++ {
		if mani[i].ReDown == true {

			var gitParams []string
			var localRepoDir string = dirDownload + mani[i].Name

			if mani[i].LocalRepoDir != "" {
				localRepoDir = dirDownload + mani[i].LocalRepoDir
			}

			//Tag and commit-id should not co-exist
			if mani[i].Tag != "" {
				gitParams = []string{"checkout", mani[i].Tag}
			} else if mani[i].LastRevision != "" {
				gitParams = []string{"checkout", mani[i].LastRevision}
			} else {
				fmt.Println(localRepoDir, ": last revision is the head of the branch.")
				continue
			}
			execCommand("git", localRepoDir, gitParams)
		}
	}
}
