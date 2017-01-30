
package main

import (
    "flag"
    "fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

func main() {
    flag.Usage = func() {
        fmt.Printf("usage: %s <command> [<arg>..]\n", os.Args[0])
        fmt.Println()
        fmt.Println("Take all enviroment variables and expand \"secret:<name>\"")
        fmt.Println("with the content of '/run/secrets/<name>' file. The path")
        fmt.Println("may be overriden by setting $SECRETS_PATH.")
        fmt.Println()
        fmt.Println("Then it will executed <command> with <args>, in such environment.")
    }
    flag.Parse()

	if len(os.Args) < 2 {
		log.Fatal("expecting command as an argument")
	}

	SECRETS_PATH := os.Getenv("SECRETS_PATH")
	if SECRETS_PATH == "" {
		SECRETS_PATH = "/run/secrets"
	}

	secret_not_found := false
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key, value := pair[0], pair[1]
		if strings.HasPrefix(value, "secret:") {
			secret_name := strings.TrimPrefix(value, "secret:")
			secret_filename := path.Join(SECRETS_PATH, secret_name)
			secret_value, err := ioutil.ReadFile(secret_filename)
			if err != nil {
				log.Printf("Unable to read \"%s\" file", secret_filename)
				secret_not_found = true
			} else {
				os.Setenv(key, string(secret_value))
			}
		}
	}
	if secret_not_found {
		log.Fatal("One or more secrets were not found, exiting")
	}

	command := os.Args[1]
	binary, err := exec.LookPath(command)
	if err != nil {
		log.Fatalf("Command \"%s\" not found", command)
	}

	execErr := syscall.Exec(binary, os.Args[1:], os.Environ())
	if execErr != nil {
		panic(execErr)
	}
}
