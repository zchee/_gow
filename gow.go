// +build darwin

package main

import (
	"bufio"
	"flag"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/go-fsnotify/fsevents"
	log "github.com/sirupsen/logrus"
	"github.com/tcnksm/go-gitconfig"
)

// Disable vcs folder by default
var filePath = []string{
	".git",
	".hg",
	".svn",
}

var noteDescription = map[fsevents.EventFlags]string{
	fsevents.MustScanSubDirs: "MustScanSubdirs",
	fsevents.UserDropped:     "UserDropped",
	fsevents.KernelDropped:   "KernelDropped",
	fsevents.EventIDsWrapped: "EventIDsWrapped",
	fsevents.HistoryDone:     "HistoryDone",
	fsevents.RootChanged:     "RootChanged",
	fsevents.Mount:           "Mount",
	fsevents.Unmount:         "Unmount",

	fsevents.ItemCreated:       "Created",
	fsevents.ItemRemoved:       "Removed",
	fsevents.ItemInodeMetaMod:  "InodeMetaMod",
	fsevents.ItemRenamed:       "Renamed",
	fsevents.ItemModified:      "Modified",
	fsevents.ItemFinderInfoMod: "FinderInfoMod",
	fsevents.ItemChangeOwner:   "ChangeOwner",
	fsevents.ItemXattrMod:      "XAttrMod",
	fsevents.ItemIsFile:        "IsFile",
	fsevents.ItemIsDir:         "IsDir",
	fsevents.ItemIsSymlink:     "IsSymLink",
}

var (
	path    = flag.String("path", CurrentDir(), "Watch directory path")
	command = flag.String("command", "", "Run command name after any event. Require -event flag")
	event   = flag.String("event", "", "Watch directory path")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// TODO: Support Windows ansi color
	// It enable when support cross platform
	// "github.com/mattn/go-colorable"
	// log.SetOutput(colorable.NewColorableStdout())

	flag.Parse()

	es := &fsevents.EventStream{
		Paths:   []string{*path},
		Latency: 0 * time.Millisecond,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot}
	es.Start()
	ec := es.Events

	go func() {
		for msg := range ec {
			for _, event := range msg {
				logEvent(event)
			}
		}
	}()

	in := bufio.NewReader(os.Stdin)

	log.Warnln("Watch started.")
	for {
		log.Warnln("press C-c to quit. press enter to pause.")
		in.ReadString('\n')
		runtime.GC()
		es.Stop()

		log.Warnln("Stopped, press enter to restart.")
		in.ReadString('\n')
		es.Resume = true
		es.Start()
	}
}

// Support space separate cmd args for exec.Command()
// Can be specified in space separate
func execCommand(cmd string) {
	// Skip when not set command flag
	if *command == "" {
		return
	}

	// Split & convert args "c" to string slice
	split := strings.Split(cmd, " ")
	argc := split[0]
	var argv = strings.Fields(split[1])
	for i := 2; i < len(split); i++ {
		argv = append(argv, split[i])
	}

	c := exec.Command(argc, argv...)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	c.Run()

	return
}

func logEvent(event fsevents.Event) {
	gitignoreGloabalPath, err := gitconfig.Global("core.excludesfile")
	if err != nil {
		log.Fatal(err)
	}

	gitignoreGlobal, err := os.Open(gitignoreGloabalPath)
	if err != nil {
		log.Fatal(err)
	}
	defer gitignoreGlobal.Close()

	scanner := bufio.NewScanner(gitignoreGlobal)
	for scanner.Scan() {
		filePath = append(filePath, scanner.Text())
	}

	f := strings.Split(event.Path, "/")
	if !StringInSlice(f[len(f)-1], filePath) {
		note := ""
		for bit, description := range noteDescription {
			if event.Flags&bit == bit {
				note += description + " "
			}
		}

		log.Infoln("EventID:", event.ID)
		log.Infoln("Path:", event.Path)
		log.Infoln("Flags:", note)

		for eventType, _ := range noteDescription {
			switch event.Flags & eventType {

			case fsevents.ItemModified:
				go execCommand(*command)

			}
		}
	}
}
