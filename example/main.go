package main

import (
	"fmt"
	//er "github.com/skilld-lab/entity-rules-engine"
	er ".."
)

func main() {
	entityRules, err := er.LoadFromYAML("rules.yaml", er.WithEntityFuncs(Contribution{}))
	if err != nil {
		fmt.Println(err)
	} else {
		c := LoadContribution()
		if err := entityRules.ApplyOn(&c); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(c)
		}
	}
}
