package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	//er "github.com/skilld-lab/entity-rules-engine"
	er ".."
	"log"
)

func main() {
	c := LoadContribution()
	entity := Contribution{}
	entityRules, err := er.LoadFromYAML("/home/bgrossmann/Sources/skilld-machine/Test2/entity-rules-engine/example/config.yaml", er.WithEntityFuncs(entity))
	if err != nil {
		log.Fatal(err)
	} else {
		if err := entityRules.ApplyOn(&c); err != nil {
			fmt.Println(err)
		} else {
			spew.Dump(c)
		}
	}
}
