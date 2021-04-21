# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-21-2021

|Num Edges|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause (ns)|
|-:|-:|-:|-:|-:|-:|
|2x2|1|68.822µs|9.055|0.00|0s|
|5x5|11|334.724µs|53.991|0.00|0s|
|7x7|11|957.646µs|128.607|0.00|0s|
|10x10|11|12.926302ms|1678.897|0.45|734.516µs|
|15x15|11|75.847298ms|22362.230|9.55|555.165µs|
|20x20|11|6.081758628s|1743986.348|716.82|42.963817ms|
