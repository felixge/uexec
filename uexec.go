package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

var userName = flag.String("user", "", "user name")

func main() {
	flag.Usage = func() {
		fmt.Println("usage: uexec [-user=user] command")
		fmt.Println("")
		fmt.Println("Works like sudo- u, but does not messs with your environment.")
		fmt.Println("")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		panic(err)
	}
	args[0] = path
	env := os.Environ()
	if *userName != "" {
		if err := changeUser(*userName); err != nil {
			panic(err)
		}
	}
	if err := syscall.Exec(path, args, env); err != nil {
		panic(err)
	}
}

func changeUser(name string) error {
	user, err := user.Lookup(name)
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(user.Gid)
	if err = syscall.Setgroups([]int{gid}); err != nil {
		return err
	}
	if err = syscall.Setuid(uid); err != nil {
		return err
	}
	return nil
}
