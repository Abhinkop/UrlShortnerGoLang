package main

// Done till pre git pull step(look into og-git for this)
// must do pull, build and the tar the the binary and copy it to final dest
// the use the compiled binary to build a docker image and push it docker.io
import (
	"fmt"
	"os"
)

var goPath = ""

func main() {
	finalDest, err := os.Getwd()
	printError("Error in getting Current Working Dir env variable", err)
	fmt.Println(finalDest)
	goPath, isPresent := os.LookupEnv("GOPATH")
	if !isPresent {
		homeDir, _ := os.LookupEnv("HOME")
		fmt.Println(homeDir)
		goPath = homeDir + "/staging_dir/go"
		fmt.Println("Go Path not set")
		fmt.Println("Using \"" + goPath + "\" as default value for GOPATH")
		err = os.Setenv("GOPATH", goPath)
		printError("Error in setting GOPATH env variable", err)
	}
	workingDir := goPath + "/src"
	err = os.RemoveAll(goPath)
	printError("Error pre-cleaning up  Directory \""+goPath+"\"", err)
	err = os.MkdirAll(workingDir, 0755)
	printError("Error creating Directory \"src\"", err)

	err = os.Chdir(workingDir)
	printError("Error changing working directory to  \""+workingDir+"\" : ", err)
	str, _ := os.Getwd()
	err = os.RemoveAll(goPath)
	printError("Error pre-cleaning up  Directory \""+goPath+"\"", err)
	fmt.Println(str)
	fmt.Println(isPresent)
}

func printError(message string, err error) {
	if err != nil {
		fmt.Println(message+" : ", err)
		os.Exit(-1)
	}
}
