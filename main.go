package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	//read config
	if len(os.Args) < 2 {
		log.Fatal("Usage: ", "go-fswatcher config.json")
		os.Exit(1)
	}
	configPath := os.Args[1]

	file, e := ioutil.ReadFile(configPath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var watches []Watch
	json.Unmarshal(file, &watches)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	commands := make(map[string][]string)

	for _, watch := range watches {
		log.Println("Run watcher for ", watch.Path)

		if _, err := os.Stat(watch.Path); os.IsNotExist(err) {
			log.Println("File ", watch.Path, " doesn't exists. Skip")
			continue
		}

		err = watcher.Add(watch.Path)
		if err != nil {
			log.Fatal(err)
		}

		for _, command := range watch.Commands {
			log.Println("Register command: ", command)
			commands[watch.Path] = append(commands[watch.Path], command)
		}
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file: ", event.Name)
					for _, command := range commands[event.Name] {
						log.Println("Run command: ", command)
						exeCmd(command)
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	<-done

}

func exeCmd(cmd string) {
	fmt.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
}

type Watch struct {
	Path     string
	Commands []string
}
