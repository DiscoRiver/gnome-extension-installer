package toadstool

import (
	"fmt"
	"github.com/discoriver/toadstool/util"
	"github.com/lithammer/shortuuid"
)

func InstallExtension(extensionZip string) {
	initToadstool()

	if v := isValidZipFile(extensionZip); v != true {
		fmt.Printf("SKIPPING %s, not a valid .zip directory. \n", extensionZip)
		return
	}
	// Protection for concurrency
	tmpDir := fmt.Sprintf("%s%s", extensionsDirectory, shortuuid.New())

	// Make temp directory
	err := makeDirectory(tmpDir, 0775)
	if err != nil {
		fmt.Printf("SKIPPING %s, couldn't make directory: %s\n", extensionZip, tmpDir)
		installFailureTeardown(tmpDir)
		return
	}

	// Unzip extension
	err = util.Unzip(extensionZip, tmpDir)
	if err != nil {
		fmt.Printf("SKIPPING %s, couldn't unzip: %s\n", extensionZip, err)
		installFailureTeardown(tmpDir)
		return
	}

	extensionUUID, err := getMetaFromExtention(tmpDir)
	if err != nil {
		fmt.Printf("SKIPPING %s, couldn't get metadata UUID: %s\n", extensionZip, err)
		installFailureTeardown(tmpDir)
		return
	}

	finalExtensionDirectory := fmt.Sprintf("%s%s%s", extensionsDirectory, "/", extensionUUID)

	if !doesDirectoryExist(finalExtensionDirectory) {
		if err := renameDirectory(tmpDir, finalExtensionDirectory); err != nil {
			fmt.Printf("SKIPPING %s, couldn't rename directory: %s\n", extensionZip, err)
			installFailureTeardown(tmpDir)
		}
		fmt.Printf("INSTALLATION COMPLETE: %s", extensionUUID)
	} else {
		if r := askYesNo(fmt.Sprintf("WARNING: Attempting to install extension at %v, but directory already exists. Continuing will overwrite this directory. Would you like to continue? (yes/no): ", finalExtensionDirectory)); r == "yes" {
			err := removeDirectory(finalExtensionDirectory)
			if err != nil {
				fmt.Printf("SKIPPING %s, couldn't remove directory: %s\n", extensionZip, finalExtensionDirectory)
				installFailureTeardown(tmpDir)
				return
			}

			err = renameDirectory(tmpDir, finalExtensionDirectory)
			if err != nil {
				fmt.Printf("SKIPPING %s, couldn't rename directory: %s\n", extensionZip, tmpDir)
				installFailureTeardown(tmpDir)
				return
			}
		}
	}
}

func installFailureTeardown(dir string) {
	err := removeDirectory(dir)
	if err != nil {
		fmt.Printf("Failure in teardown, couldn't remove directory %s\n", dir)
		return
	}
}

