This package allows you to create an entity rules engine.

Let's imagine you have an object with getter methods (returning boolean) and setter methods, and that you want to establish rules to change this object accordingly to events.

# Rules

You will have to define rules. A rule provides the object with a reaction to some triggers (When) and under some conditions (If) by performing actions (Do).
here is the structure of one rule:
```go
type Rule struct {
	Name 	    string     // Id
 	Description string     // Describes the rule
	When        []string   // Rule's triggers
	If          string     // Rule's conditions
	Do          []string   // Rule's actions
}
```
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
It needs to know what method of the object to call and with which arguments.
```go
type MethodExecution struct {
        Name      string    //Id
        Method    string    //Name of the method related to the object
        Arguments Arguments //Arguments to evaluate the Method with
}
type Arguments []Argument
type Argument interface{}
```
We divide method execution information into: conditions and actions methods
```go
type MethodsExecution []MethodExecution
Conditions MethodsExecution // List of methods for the When and If statements (getters)
Actions    MethodsExecution // List of methods for the Do statement (setters)
```

# EntityRules

The entity EntityRules combine all the previous information.

```go
type EntityRules struct {
        Conditions MethodsExecution
        Actions    MethodsExecution
        Rules      Rules
}
```
While Conditions and Actions contain method execution information, Rules describe the logic between those evaluations.
Then it is mandatory to have the methods called in Rules be listed in Conditions and Actions.

# Usage

```go
import ebr "github.com/skilld-labs/entity-rules-engine"
```
## Generate your rules

Let's load the object we want to apply rules on.
```go 
c := LoadContribution()
```
Then, let's load our EntityRules:

You can either apply your own loader or use our loaders from either :
```go
func LoadFromJSON(filePath string, opts ...LoadOption) (entityRules *EntityRules, err error) 
func LoadFromYAML(filePath string, opts ...LoadOption) (entityRules *EntityRules, err error) 
func LoadFromMap(m map[string]interface{}, opts ...LoadOption) (*EntityRules, error) 
func LoadFromInterface(er interface{}, opts ...LoadOption) (*EntityRules, error) 
```
Knowing there are a few With... option functions that can be used to customize the entityRules loading.
For example, to set a custom template.FuncMap parsing, using:
```go
WithEntityFuncs(entity) //Entity being an empty variable of the object on which we want to apply rules
WithFuncMaps(fm) //fm being a funcMap we want to use for parsing
```
Please note that by default there is no template parsing.

Finally you can apply on the object the loaded EntityRules using: 
```go
entityRules.ApplyOn(&c)
```
And even of multiple variables of the same object.
```go
entityRules.ApplyOn(&c1, &c2, ...)
```

# Examples

The example directory provides the reader with a clear examples, but here is another one to understand better this library's logic.

for instance let's say we want to hug Tommy when he enters the room where are in.

The object is the room we are in, having attributes about who is in it and what is going on.

To define the EntityRules we will define a rule and list Actions and Conditions.

The EntityRule will be defined like this :
```
rule:
- description: Hugging Tommy as welcome
- name: WelcomingTommy
- when: EnteringRoom
- on: IsTommy
- do: HugTommy

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
 ```


# Contribution
written by Brunelle Grossmann
with help from skilld-labs
