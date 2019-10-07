// acmecrystal watches for .cr files being written then runs the
// crystal formatting tool on that written file and then reloads the
// acme window with the formatted file
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

// web https://9fans.github.io/plan9port/man/man4/acme.html
func main() {
	flag.Parse()
	//log   reports a log of window operations since the opening of the
	//log file. Each line describes a single operation using three fields
	//separated by single spaces: the decimal window ID, the operation,
	//and the window name. Reading from log blocks until there is an
	//operation to report, so reading the file can be used to monitor
	//editor activity and react to changes. The reported operations are
	//new (window creation), zerox (window creation via zerox), get, put,
	//and del (window deletion). The window name can be the empty string;
	//in particular it is empty in new log entries corresponding to windows
	//created by external programs.
	l, err := acme.Log()

	if err != nil {
		log.Fatal(err)
	}

	for {
		event, err := l.Read()
		if err != nil {
			log.Fatal(err)
		}

		//If log reports that a crystal file is being "put", then format it
		if event.Name != "" && event.Op == "put" && strings.HasSuffix(event.Name, ".cr") {
			crystalFormat(event.ID, event.Name)
		}
	}
}

func crystalFormat(id int, name string) {
	// When a command is run under acme, a directory holding these files
	// is posted as the 9P service acme
	w, err := acme.Open(id, nil)
	if err != nil {
		log.Print(err)
		return
	}

	// remember dot
	w.ReadAddr() // not sure why you need to open this to init
	w.Ctl("addr=dot")
	p1, p2, _ := w.ReadAddr()

	defer w.CloseFiles()
	old, err := ioutil.ReadFile(name)
	if err != nil {
		return
	}

	// `crystal tool format <filename>` reformats the file in place.  So
	// after running this the file on disk will be reformatted, but not
	// the file in the Acme window
	// os.Setenv("TERM", "dumb") when my pull request gets accepted this should work
	// https://github.com/crystal-lang/crystal/pull/8271
	out, err := exec.Command("crystal", "tool", "format", "--no-color", name).CombinedOutput()
	if err != nil {
		log.Printf("%s", out)
		return
	}

	new, err := ioutil.ReadFile(name)
	if err != nil {
		return
	}
	if bytes.Equal(old, new) {
		return
	}

	w.Write("ctl", []byte("mark"))
	w.Write("ctl", []byte("nomark"))
	err = w.Addr("0,$")
	if err != nil {
		log.Print(err)
		return
	}
	w.Write("data", new)

	// put dot back where it was
	w.Addr("#%d,#%d", p1, p2)
	w.Ctl("dot=addr")
	w.Ctl("show")
}
