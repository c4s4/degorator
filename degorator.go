package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Terminate(message string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Println(message)
	} else {
		fmt.Printf(message+"\n", args...)
	}
	os.Exit(1)
}

func WriteResponse(writer http.ResponseWriter, status int, message string, args ...interface{}) {
	writer.WriteHeader(status)
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	writer.Write([]byte(message))
}

func main() {
	if len(os.Args) != 2 {
		Terminate("You must pass service configuration file on command line")
	}
	source, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		Terminate("Error reading configuration file: %v", err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		Terminate("Error parsing YAML configuration file: %v", err)
	}
	if err := (&config).Compile(); err != nil {
		Terminate("Error loading configuration: %v", err)
	}
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
