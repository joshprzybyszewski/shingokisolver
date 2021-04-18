# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-18-2021

|Num Edges|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause (ns)|
|-:|-:|-:|-:|-:|-:|
|2x2|1|90.494µs|9.461|0.00|0s|
|5x5|101|382.187µs|52.779|0.00|0s|
|7x7|102|1.285391ms|200.680|0.00|0s|
|10x10|101|18.338776ms|3709.086|1.30|42.122µs|
|15x15|102|8.159503524s|3222823.863|1095.00|54.463891ms|
