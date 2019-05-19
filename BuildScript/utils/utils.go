package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

/*Copy Copies files or folder based on The OS*/
func Copy(inSrc string, inDest string) {
	if runtime.GOOS == "windows" {
		copyWindows(inSrc, inDest)
	} else {
		copyLinux(inSrc, inDest)
	}
}

func copyWindows(inSrc string, inDest string) {
	winpathinDest := strings.ReplaceAll(inDest, "/", "\\")
	var cmd *exec.Cmd
	if isADirectory(inSrc) {
		err := os.MkdirAll(winpathinDest+"\\"+inSrc, 0755)
		PrintError("Error creating Directory :"+inSrc, err)
		finalDest := winpathinDest + "\\" + inSrc
		log.Println(finalDest)
		cmd = exec.Command("xcopy", inSrc, finalDest, "/E")
	} else {
		cmd = exec.Command("xcopy", inSrc, winpathinDest)
	}
	_, err := cmd.CombinedOutput()
	msg := "copying `%s` failed with :"
	PrintError(fmt.Sprintf(msg, inSrc), err)
}

func copyLinux(inSrc string, inDest string) {
	cmd := exec.Command("cp", "-r", inSrc, inDest)
	out, err := cmd.CombinedOutput()
	log.Println("cmd Output :", out)
	msg := "copying `%s` failed with :"
	PrintError(fmt.Sprintf(msg, inSrc), err)
}

func isADirectory(path string) bool {
	f, err := os.Open(path)
	defer f.Close()
	PrintError("isADirectory :", err)
	stat, err := f.Stat()
	PrintError("isADirectory :", err)
	PrintError("isADirectory :", err)
	return stat.IsDir()
}

/*Compress Creates zip or tar package based on OS*/
func Compress(compressedFilename string, folderName string) {
	if runtime.GOOS == "windows" {
		compressWindows(compressedFilename, folderName)
	} else {
		compressLinux(compressedFilename, folderName)
	}
}

func compressWindows(compressedFilename string, folderName string) {
	log.Println("Start: Building zip")
	cmd := exec.Command("7z", "a", compressedFilename+".zip", folderName)
	_, err := cmd.CombinedOutput()
	PrintError("Building zip failed with :", err)
	log.Println("End: Building zip")
}

func compressLinux(compressedFilename string, folderName string) {
	log.Println("Start: Building tar")
	cmd := exec.Command("tar", "-cf", compressedFilename+".tar", folderName)
	out, err := cmd.CombinedOutput()
	log.Println("cmd Output :", out)
	PrintError("Building tar failed with :", err)
	log.Println("End: Building tar")
}

/*PrintError Prints an error and mesage on the screen*/
func PrintError(message string, err error) {
	if err != nil {
		log.Fatalln(message+" : ", err)
	}
}
