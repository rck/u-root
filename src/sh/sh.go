// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
sh reads in a line at a time and runs it. 
prompt is '% '
*/

package main

import (
	"os/exec"
	"fmt"
	"os"
	"strings"
	"bufio"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Println("no scripts/args yet")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%% ")
	for scanner.Scan() {
		cmd := scanner.Text()
		argv := strings.Split(cmd, " ")
		e := os.Environ()
		// sudo sensibly doesn't inherit the path if you are root.
		// and probably doesn't in general so we need to do this much.
		// interestingly, on arch, I did not need to do this. Sounds bad.
		np := strings.NewReplacer("PATH=", "PATH=/go/bin:/buildbin:/bin:")
		for i := range(e) {
			e[i] = np.Replace(e[i])
		}
		e = append(e, "GOROOT=/go")
		e = append(e, "GOPATH=/")
		e = append(e, "GOBIN=/bin")
		// oh, and, Go looks in the environment, NOT the env in the cmd.
		p := os.Getenv("PATH")
		if err := os.Setenv("PATH", "/go/bin:/buildbin:/bin:/usr/local/bin:" + p); err != nil {
			fmt.Printf("Couldn't set path; %v\n", err)
			continue
		}
		p = os.Getenv("LD_LIBRARY_PATH")
		if err := os.Setenv("LD_LIBRARY_PATH", p + ":/usr/local/lib"); err != nil {
			fmt.Printf("Couldn't set LD_LIBRARY_PATH; %v\n", err)
			continue
		}
		run := exec.Command(argv[0], argv[1:]...)
		run.Env = e
		out, err := run.CombinedOutput()
		if err != nil {
			fmt.Printf("%v: Path %v\n", err, os.Getenv("PATH"))
		}
		fmt.Printf("%s", out)
		fmt.Printf("%% ")
	}
}