package main

// Done till pre git pull step(look into og-git for this)
// must do pull, build and the tar the the binary and copy it to final dest
// the use the compiled binary to build a docker image and push it docker.io
import (
	"log"
	"os"
	"os/exec"
)

var goPath = ""

func main() {

	log.Println("Build Started")
	finalDest, err := os.Getwd()
	printError("Error in getting Current Working Dir: ", err)
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
		printError("Error in setting GOPATH env variable", err)
	}

	// Clean up before build
	cleanUp(goPath)

	// Clean up after build
	defer cleanUp(goPath)

	// Make src dir for go code and change to it.
	log.Println("Start: creating /src dir and changing into it.")
	workingDir := goPath + "/src"
	err = os.MkdirAll(workingDir, 0755)
	printError("Error creating Directory `src`", err)
	err = os.Chdir(workingDir)
	printError("Error changing working directory to  \""+workingDir+"\" : ", err)
	log.Println("End: creating /src dir and changing into it.")

	// git clone the repo
	log.Println("Start: git cloning")
	repo := "https://github.com/Abhinkop/UrlShortnerGoLang.git"
	cmd := exec.Command("git", "clone", repo)
	err = cmd.Run()
	printError("git clone failed with :", err)
	log.Println("End: git cloning")

	// change to UrlShortnerGoLang dir to build
	log.Println("Start: change to `UrlShortnerGoLang` dir to build")
	err = os.Chdir(workingDir + "/UrlShortnerGoLang")
	printError("changing dir to `"+workingDir+"` UrlShortnerGoLang failed :", err)
	log.Println("End: change to `UrlShortnerGoLang` dir to build")

	// dep ensue to pull dependecies
	log.Println("Start: fetching dependecies (dep ensure)")
	cmd = exec.Command("dep", "ensure")
	out, err := cmd.CombinedOutput()
	log.Println("cmd Output :",out)
	printError("`dep ensure` failed with :", err)
	log.Println("End: fetching dependencies (dep ensur)e")

	// Build main.go
	log.Println("Start: go build main.go")
	cmd = exec.Command("go", "build", "main.go")
	out, err = cmd.CombinedOutput()
	log.Println("cmd Output :",out)
	printError("`go build maing.go` failed with :", err)
	log.Println("End: go build main.go")

	// create and copy the resouces and binary to build_output dir
	log.Println("Start: copy `resources` to `build_output")
	buildOutput := goPath + "/build_output"
	err = os.MkdirAll(buildOutput, 0755)
	printError("Error creating Directory :"+buildOutput, err)
	cmd = exec.Command("cp", "-r", "resources", buildOutput)
	out, err = cmd.CombinedOutput()
	log.Println("cmd Output :",out)
	printError("copying `resources` folder failed with :", err)
	log.Println("Start: copy `main` to `build_output")
	cmd = exec.Command("cp", "-r", "main", buildOutput)
	out, err = cmd.CombinedOutput()
	log.Println("cmd Output :",out)
	printError("copying `main` binary failed with :", err)
	log.Println("End: finished copying")

	// change to $GOPATH dir to build tar
	log.Println("Start: change dir to " + goPath)
	err = os.Chdir(goPath)
	printError("changing dir to "+goPath+"UrlShortnerGoLang failed :", err)
	log.Println("End: change dir to " + goPath)

	tarFileName := "Url.tar"
	// tarball the build
	log.Println("Start: Building tar")
	cmd = exec.Command("tar", "-cf", tarFileName, "build_output")
	out, err = cmd.CombinedOutput()
	log.Println("cmd Output :",out)
	printError("Building tar failed with :", err)
	log.Println("End: Building tar")

	// copy to final destination
	log.Println("Start: Copy to " + finalDest)
	cmd = exec.Command("cp", tarFileName, finalDest)
	out, err = cmd.CombinedOutput()
	log.Println("cmd Output :",out)
	printError("copy to final destination("+finalDest+") failed with :", err)
	log.Println("End: Copy to " + finalDest)

	log.Println("Build Completed. Binaries copied to " + finalDest + "/"+tarFileName)
}

func printError(message string, err error) {
	if err != nil {
		log.Fatalln(message+" : ", err)
	}
}

func cleanUp(goPath string) {
	log.Println("Start: removing "+ goPath +" dir.")
	err := os.RemoveAll(goPath)
	printError("Error cleaning up  Directory \""+goPath+"\"", err)
	log.Println("End: removing "+ goPath +" dir.")
}
