This package allows you to create an entity rules engine.
This library executes different actions on an object when the right events occur and if some conditions are valid.

Let's imagine you have a type object with getters and setters methods.

the getter methods return boolean
setter methods perform action on the object

# Listing conditions and actions
Let's have conditions and actions following this structure:
```
conditions map[string]methodExecution
actions    map[string]methodExecution

type methodExecution struct {
	name   string
	method string
	params []interface{}
}
```
conditions being a map of methodExecution representing the getter methods.
actions being a map of methodExecution representing the setter methods

a methodExecution type has:
- name:    It stands to be an id
- method:  The name of the method of the type object
- params: The arguments that you want to evaluate the method with


# Defining a rule
```
type rule struct {
	description string
	name        string
	when        []string
	on          string
	do          []string
}
```

- description: describes the rule.
- name: rule's id.
- when: defines the event on which the rule is called.
- on: defines the conditions that are requiered to execute the rule.
- do: defines the actions performed after when and on statements have been validated.

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
 
 ### When
Is a string list of the conditions we want to have as triggers.
The logical link between them is an OR. Meaning that if one trigger is valid from the selected triggers then When is validated
 
 ### On
Is a string representing the logical link between the conditions.
write the name of the conditions in a logical sentence like that for instance:
 
cond_1 and ( cond_2 or cond_3 )
 
*(beware of the blank spaces around "(" and ")" )*
 
 ### Do
Is a string list of the actions we want to perform.
 
# Generate your models
- Lists conditions and actions methods you want to call.
- Create your rules.

# Contribution
written by Brunelle Grossmann
with help from skilld-labs
