package main

import (
	"flag"
	"reflect"
	"strings"
)

func (c *Contribution) LabelsLen(l int) bool {
	return len(c.Labels) == l
}

func (c *Contribution) HasLabel(l string) bool {
	return true
}

func (c *Contribution) HasRevision() bool {
	return true
}
func (c *Contribution) LabelsHasChanged() bool {
	return true
}
func (c *Contribution) AddLabel(l string) {
	c.Labels = append(c.Labels, l)
}
func (c *Contribution) HasAssignee() bool {
	return c.Assignee != ""
}

func (c *Contribution) SetAssignee(l string) {
	c.Assignee = l
}

type Condition func(reflect.Value) bool

type Contribution struct {
	Assignee string
	Author   string
	Labels   []string
}

func LoadContribution() Contribution {
	assignee := flag.String("assignee", "Tommy ", "")
	labels := flag.String("labels", "kind/story", "plop")
	author := flag.String("author", "I'm the author", "")
	flag.Parse()
	c := Contribution{Assignee: *assignee, Labels: strings.Split(*labels, ","), Author: *author}
	return c
}
