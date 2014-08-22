package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	cmd_cc   = flag.String("CC", "gcc", "C compiler, support gcc and clang")
	cmd_cxx  = flag.String("CXX", "g++", "C++ compiler, support g++ and clang++")
	cmd_make = flag.String("MAKE", "make", "support GNU make")
)

func main() {
	flag.Parse()

	fmt.Println("CC: " + *cmd_cc)
	fmt.Println("CXX: " + *cmd_cxx)
	fmt.Println("Make: " + *cmd_make)

	cmd_rm := "rm"
	switch runtime.GOOS {
	case "darwin":
	case "linux":
	case "windows":
		cmd_rm = "del"
	}
	os.Setenv("CC", *cmd_cc)
	os.Setenv("CXX", *cmd_cxx)
	os.Setenv("RM", cmd_rm)

	log.Println("Compiling")
	cmd := exec.Command(*cmd_make, "-C", "liblinear_src")
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		log.Println(err.Error())
		return
	}

	os.Mkdir("lib", os.ModeDir|0777)

	log.Println("Installing libs")
	if runtime.GOOS == "windows" {
		err = copyFile("liblinear_src/linear.dll", "lib/linear.dll")
		if err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		err = copyFile("liblinear_src/liblinear.so", "lib/liblinear.so")
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	log.Println("Cleaning")
	exec.Command(*cmd_make, "-C", "liblinear_src", "clean").Run()

	lib_path, err := filepath.Abs("lib")
	var target_envir string
	switch runtime.GOOS {
	case "darwin":
		target_envir = "DYLD_LIBRARY_PATH"
	case "linux":
		target_envir = "LD_LIBRARY_PATH"
	case "windows":
		target_envir = "PATH"
	}

	fmt.Println("IMPORTANT:")
	fmt.Println("Add \"" + lib_path + "\" to your " + target_envir + " environment manually")
}

func copyFile(src, dst string) error {
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if sfi.Mode()&os.ModeType != 0 {
		log.Fatalf("mirrorFile can't deal with non-regular file %s", src)
	}
	dfi, err := os.Stat(dst)
	if err == nil &&
		isExecMode(sfi.Mode()) == isExecMode(dfi.Mode()) &&
		(dfi.Mode()&os.ModeType == 0) &&
		dfi.Size() == sfi.Size() &&
		dfi.ModTime().Unix() == sfi.ModTime().Unix() {
		// Seems to not be modified.
		return nil
	}

	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	n, err := io.Copy(df, sf)
	if err == nil && n != sfi.Size() {
		err = fmt.Errorf("copied wrong size for %s -> %s: copied %d; want %d", src, dst, n, sfi.Size())
	}
	cerr := df.Close()
	if err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Chmod(dst, sfi.Mode())
	}
	if err == nil {
		err = os.Chtimes(dst, sfi.ModTime(), sfi.ModTime())
	}
	return err
}
func isExecMode(mode os.FileMode) bool {
	return (mode & 0111) != 0
}
