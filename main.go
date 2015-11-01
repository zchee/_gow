// +build darwin

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/go-fsnotify/fsevents"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/tcnksm/go-gitconfig"
)

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

func main() {
	log.SetOutput(colorable.NewColorableStdout())
	p, _ := os.Getwd()
	path := flag.String("path", p, "Watch directory path")

	flag.Parse()

	// dev, _ := fsevents.DeviceForPath(*path)
	// log.Print(dev)
	// log.Println(fsevents.EventIDForDeviceBeforeTime(dev, time.Now()))

	es := &fsevents.EventStream{
		Paths:   []string{*path},
		Latency: 0 * time.Millisecond,
		// Device:  dev,
		Flags: fsevents.FileEvents | fsevents.WatchRoot}
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
func execCommand(command string) {
	split := strings.Split(command, " ")
	argc := split[0]

	var argv = strings.Fields(split[1])
	for i := 2; i < len(split); i++ {
		argv = append(argv, split[i])
		fmt.Println(argv)
	}

	cmd := exec.Command(argc, argv...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v")
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
	// cmd.Run()
	return
}

func logEvent(event fsevents.Event) {
	// p, _ := os.Getwd()
	// path := flag.String("path", p, "Watch directory path")
	var filePath = []string{}

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

	note := ""
	for bit, description := range noteDescription {
		if event.Flags&bit == bit {
			note += description + " "
		}
	}

	f := strings.Split(event.Path, "/")
	if !StringInSlice(f[len(f)-1], filePath) {
		log.Infoln("EventID:", event.ID)
		log.Infoln("Path:", event.Path)
		log.Infoln("Flags:", note)

		for eventType, _ := range noteDescription {
			if event.Flags&eventType == fsevents.ItemModified {
				execCommand("echo hello")
			}
		}
	}
}
