package main

// Done till pre git pull step(look into og-git for this)
// must do pull, build and the tar the the binary and copy it to final dest
// the use the compiled binary to build a docker image and push it docker.io
import (
	"log"
	"os"
	"os/exec"
	"runtime"

	"UrlShortnerGoLang/BuildScript/utils"
)

var goPath = ""

func main() {

	log.Println("Build Started")
	finalDest, err := os.Getwd()
	utils.PrintError("Error in getting Current Working Dir: ", err)
	goPath, isPresent := os.LookupEnv("GOPATH")
	if !isPresent {
		homeDir, isPresent := os.LookupEnv("HOME")
		if !isPresent {
			log.Fatalln("$HOME env var not found")
		}
		goPath = homeDir + "/staging_dir/go"
		log.Println("Go Path not set")
		log.Println("Using `" + goPath + "` as default value for GOPATH")
		err = os.Setenv("GOPATH", goPath)
		utils.PrintError("Error in setting GOPATH env variable", err)
	}

	// Clean up before build
	cleanUp(goPath)

	// Make src dir for go code and change to it.
	log.Println("Start: creating /src dir and changing into it.")
	workingDir := goPath + "/src"
	err = os.MkdirAll(workingDir, 0755)
	utils.PrintError("Error creating Directory `src`", err)
	err = os.Chdir(workingDir)
	utils.PrintError("Error changing working directory to  \""+workingDir+"\" : ", err)
	log.Println("End: creating /src dir and changing into it.")

	// git clone the repo
	log.Println("Start: git cloning")
	repo := "https://github.com/Abhinkop/UrlShortnerGoLang.git"
	cmd := exec.Command("git", "clone", repo)
	err = cmd.Run()
	utils.PrintError("git clone failed with :", err)
	log.Println("End: git cloning")

	// change to UrlShortnerGoLang dir to build
	log.Println("Start: change to `UrlShortnerGoLang` dir to build")
	err = os.Chdir(workingDir + "/UrlShortnerGoLang")
	utils.PrintError("changing dir to `"+workingDir+"` UrlShortnerGoLang failed :", err)
	log.Println("End: change to `UrlShortnerGoLang` dir to build")

	// dep ensue to pull dependecies
	log.Println("Start: fetching dependecies (dep ensure)")
	cmd = exec.Command("dep", "ensure")
	out, err := cmd.CombinedOutput()
	log.Println("cmd Output :", out)
	utils.PrintError("`dep ensure` failed with :", err)
	log.Println("End: fetching dependencies (dep ensur)e")

	// Build main.go
	log.Println("Start: go build main.go")
	cmd = exec.Command("go", "build", "main.go")
	out, err = cmd.CombinedOutput()
	log.Println("cmd Output :", out)
	utils.PrintError("`go build maing.go` failed with :", err)
	log.Println("End: go build main.go")

	// create and copy the resouces and binary to build_output dir
	log.Println("Start: copy `resources` to `build_output")
	buildOutput := goPath + "/build_output"
	err = os.MkdirAll(buildOutput, 0755)
	utils.PrintError("Error creating Directory :"+buildOutput, err)

	utils.Copy("resources", buildOutput)
	utils.Copy("Config", buildOutput)
	switch os := runtime.GOOS; os {
	case "windows":
		utils.Copy("main.exe", buildOutput)
	case "linux":
		utils.Copy("main", buildOutput)
	default:
		log.Printf("%s not yet supported.\n", os)
	}
	log.Println("End: finished copying")

	// change to $GOPATH dir to build tar
	log.Println("Start: change dir to " + goPath)
	err = os.Chdir(goPath)
	utils.PrintError("changing dir to "+goPath+"UrlShortnerGoLang failed :", err)
	log.Println("End: change dir to " + goPath)

	compressedFileName := "Url"
	utils.Compress(compressedFileName, "build_output")

	switch os := runtime.GOOS; os {
	case "windows":
		utils.Copy(compressedFileName+".zip", finalDest)
	case "linux":
		utils.Copy(compressedFileName+".tar", finalDest)
	default:
		log.Printf("%s not yet supported.\n", os)
	}

	log.Println("Build Completed. Binaries copied to " + finalDest)
	cleanUp(goPath)
}

func cleanUp(goPath string) {
	log.Println("Start: removing " + goPath + " dir.")
	err := os.RemoveAll(goPath)
	utils.PrintError("Error cleaning up  Directory \""+goPath+"\"", err)
	log.Println("End: removing " + goPath + " dir.")
}
