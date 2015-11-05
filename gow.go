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
)

// Disable vcs folder by default
var filePath = []string{}

var noteDescription = map[fsevents.EventFlags]string{
	fsevents.MustScanSubDirs: "MustScanSubdirs",
	fsevents.UserDropped:     "UserDropped",
	fsevents.KernelDropped:   "KernelDropped",
	fsevents.EventIDsWrapped: "EventIDsWrapped",
	fsevents.HistoryDone:     "HistoryDone",
	fsevents.RootChanged:     "RootChanged",
	fsevents.Mount:           "Mount",
	fsevents.Unmount:         "Unmount",

	// fsevents.ItemCreated:       "Created",
	// fsevents.ItemRemoved:       "Removed",
	// fsevents.ItemRenamed:       "Renamed",
	// fsevents.ItemModified:      "Modified",
	fsevents.ItemCreated:       "c",
	fsevents.ItemRemoved:       "r",
	fsevents.ItemInodeMetaMod:  "InodeMetaMod",
	fsevents.ItemRenamed:       "rn",
	fsevents.ItemModified:      "m",
	fsevents.ItemFinderInfoMod: "FinderInfoMod",
	fsevents.ItemChangeOwner:   "ChangeOwner",
	fsevents.ItemXattrMod:      "XAttrMod",
	fsevents.ItemIsFile:        "IsFile",
	fsevents.ItemIsDir:         "IsDir",
	fsevents.ItemIsSymlink:     "IsSymLink",
}

var (
	path    = flag.String("path", CurrentDir(), "Watch directory path")
	command = flag.String("command", "", "Run command after flag 'event'. Require -event flag")
	event   = flag.String("event", "", "Watch event type")
	file    = flag.String("file", "", "Watch file type")
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
	args := strings.Split(cmd, ",")
	for _, i := range args {
		split := strings.Split(i, " ")
		argc := split[0]
		var argv = strings.Fields(split[1])
		for i := 2; i < len(split); i++ {
			argv = append(argv, split[i])
		}

		c := exec.Command(argc, argv...)

		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		c.Run()
		log.Infoln("Finished")

	}

	return
}

func logEvent(events fsevents.Event) {
	// Split "file" flag separate "."
	// Detection file extension
	f := strings.Split(events.Path, ".")
	if StringInSlice(f[len(f)-1], strings.Split(*file, ",")) {
		note := ""
		for bit, description := range noteDescription {
			if events.Flags&bit == bit {
				note += description + " "
			}
		}

		log.Infoln("EventID:", events.ID)
		log.Infoln("Path:", events.Path)
		log.Infoln("Flags:", note)

		for _, events := range noteDescription {
			// switch event.Flags & eventType {
			switch true {

			case events == *event:
				go execCommand(*command)
			}
		}
	}
}
