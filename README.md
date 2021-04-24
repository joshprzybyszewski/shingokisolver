# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-24-2021

|Num Edges|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause (ns)|
|-:|-:|-:|-:|-:|-:|
|2x2|1|53.01µs|14.148|0.00|0s|
|5x5|11|739.242µs|148.960|0.00|0s|
|7x7|11|2.221061ms|439.010|0.00|0s|
|10x10|11|9.749116ms|2161.961|0.36|13.324µs|
|15x15|11|70.286356ms|16823.256|5.64|230.282µs|
|20x20|11|604.25471ms|137766.121|46.00|1.550397ms|
