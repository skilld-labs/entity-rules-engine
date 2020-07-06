This package allows you to create an entity rules engine.

Let's imagine you have a type object with getter methods (returning boolean) and setter methods, and that you want to establish rules to change this object accordingly to events.
You will have to define rules, they are composed of:

```go
type Rules []Rule
type Rule struct {
	Name 	    string     // Id
 	Description string     // Describes the rule
	When        []string   // Rule's triggers
	If          string     // Rule's conditions
	Do          []string   // Rule's actions
}
```
Imagine you need your object to react to some triggers (When) and then check some conditions (If) before performing actions (Do).

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
Let's have conditions and actions following this structure:
```go
Conditions MethodsExecution // List of methods for the When and If statements (getters)
Actions    MethodsExecution // List of methods for the Do statement (setters)

type MethodsExecution []MethodExecution
type MethodExecution struct {
        Name      string    //Id
        Method    string    //Name of the method related to the object
        Arguments Arguments //Arguments to evaluate the Method with
}

type Arguments []Argument
type Argument interface{}
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
