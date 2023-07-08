package main

import (
	"log"
	"os"

	"github.com/Expand-My-Business/go_windows_agent/ExecuBot/instruction"
	"github.com/kardianos/service"
)

type myService struct{}

func (m *myService) Start(s service.Service) error {
	go m.run()
	return nil
}

func (m *myService) run() {
	instruction.GetInstructions()
}

func (m *myService) Stop(s service.Service) error {
	// Add any necessary cleanup or shutdown logic here
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "ExecuBot",
		DisplayName: "ExecuBot Service",
		Description: "Background service for executing instructions",
	}

	prg := &myService{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
