# simulationgame
Welcome to the simulation game where you can simulate creatures and how they evolve!

The thought behind this repo so to create a very simple evolution simulator, where creatures wake up, try to stay alive by competing for food, reproducing and mutating.

The aim is mostly to better myself in the language of Go, and programming in general rather than a perfect evolution simulator.

# Rules and how it works
A number of creatures compete for food each round, and if they manage to get some food, they will survive and maybe reproduce depending on their energy levels.

Creatures can perform actions, including:
* Moving
* Scanning for food
* Interacting with other creatures

Moving and scanning for food costs energy, and eating foods gives the creature energy. 

If a creature has enough energy at the end of a round, it will reproduce.

The offsprings will inherit their parent's attributes, e.g. speed and scan proc chance.

There is also possibility that the creature will mutate and thus change their attributes. But mutations can come with a cost, e.g. higher speed also equals higher energy consumption per step taken, and it's similiar with the scan proc chance - the higher the chance, the higher the cost.

Currently, there are two different types of creatures, passive and aggressive ones.

If an aggressive creature faces a passive one with food, the aggresive one will attack and kill the passive, and take its food.

If an aggressive creature faces an aggressive one with food, it will attach and try to get the food. The creatures will fight, one of them will die and the other will take some damage.

If a passive creaures faces another passive one with food, they will share the food with each other.

# How to run - quick run
Import the module:

```
import (
  "github.com/sebastianring/simulationgame"
)
```

And run it using this code:

```
sc := simulationgame.GetStandardSimulationConfig()
	_, err := simulationgame.RunSimulation(sc)

	if err != nil {
		fmt.Println("Error... : " + err.Error())
	}
```
