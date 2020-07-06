This package allows you to create an entity rules engine.

Let's imagine you have a type object with getter methods (returning boolean) and setter methods, and that you want to establish rules to change this object accordingly to events.

# Rule 

You will have to define rules, here is the structure of one rule:
```go
type Rule struct {
	Name 	    string     // Id
 	Description string     // Describes the rule
	When        []string   // Rule's triggers
	If          string     // Rule's conditions
	Do          []string   // Rule's actions
}
```
A rule provides the object with a reaction to some triggers (When) and under some conditions (If) by performing actions (Do).

 ### When
Is a string list of the conditions we want to have as triggers.
The logical link between them is an OR. Meaning that if one trigger is valid from the selected triggers then When is validated.
 ### If
Is a string representing the logical link between the conditions.
write the name of the conditions in a logical sentence like that for instance:
cond_1 and ( cond_2 or cond_3 )
 ### Do
Is a string list of the actions we want to perform.

# Listing conditions and actions

This part gives the code information related to method execution.
It needs to know what method to call and with which arguments.
```go
type MethodExecution struct {
        Name      string    //Id
        Method    string    //Name of the method related to the object
        Arguments Arguments //Arguments to evaluate the Method with
}
type Arguments []Argument
type Argument interface{}
```
We put apart conditions from actions methods:
```go
type MethodsExecution []MethodExecution
Conditions MethodsExecution // List of methods for the When and If statements (getters)
Actions    MethodsExecution // List of methods for the Do statement (setters)
```

type EntityRules struct {

        Conditions MethodsExecution
        Actions    MethodsExecution
        Conditions *MethodsExecution
        Actions    *MethodsExecution
        Rules      Rules
}

# Usage

```go
import ebr "github.com/skilld-labs/entity-rules-engine"
```
# Generate your rules

```go 
c := LoadContribution()
entity := Contribution{}
entityRules, err := ebr.LoadFromYAML("~.yaml")
if err != nil {
	log.Fatal(err)
	} else {
        if err := entityRules.ApplyOn(&c); err != nil{  
		log.Fatal(err)
                } 
        fmt.Println(c)        
        }
```

Load a New EntityRules from either :
- Json, yaml file
- interface
- map[string]interface{}

There are a few With... option functions that can be used to customize the entityRules loading. 

For example, to set a custom template.FuncMap parsing:
-> WithEntityFuncs(entity)
-> WithFuncMaps(fm) 





for instance:
rule:
- description: Hugging Tommy as welcome
- name: welcomingTommy
- when: someone enters the room
- on: this person is Tommy
- do: hug

let's say we have 

conditions:

- name: EnteringRoom 
  method: IsNewInRoom

- name: IsTommy
  method: IsName
  params: Tommy
 
actions:

- name: HugTommy
  method: Hug
  params: Tommy
  
then the rule will be defined like this :

rule:

- description: Hugging Tommy as welcome
- name: WelcomingTommy
- when: EnteringRoom
- on: IsTommy
- do: HugTommy
 


# Contribution
written by Brunelle Grossmann
with help from skilld-labs
