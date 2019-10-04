// acmecrystal watches for .cr files being written
// it then runs the crystal formatting tool on that written file
// and the reloads the acme window with the formatted file
package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"9fans.net/go/acme"
)

func main() {
	flag.Parse()
	l, err := acme.Log()
	if err != nil {
		log.Fatal(err)
	}

	for {
		event, err := l.Read()
		if err != nil {
			log.Fatal(err)
		}

		if event.Name != "" && event.Op == "put" && strings.HasSuffix(event.Name, ".cr") {
			crystalFormat(event.ID, event.Name)
		}
	}
}

func crystalFormat(id int, name string) {
	w, err := acme.Open(id, nil)
	if err != nil {
		log.Print(err)
		return
	}

	defer w.CloseFiles()
	old, err := ioutil.ReadFile(name)
	if err != nil {
		return
	}

	new, _ := exec.Command("crystal", "tool", "format", name).CombinedOutput()
	if bytes.Equal(old, new) {
		return
	}

	w.Write("ctl", []byte("mark"))
	w.Write("ctl", []byte("nomark"))
	w.Write("ctl", []byte("get"))
}
