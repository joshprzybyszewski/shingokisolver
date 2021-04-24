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
|2x2|1|112.649µs|14.148|0.00|0s|
|5x5|11|757.575µs|141.564|0.00|0s|
|7x7|11|1.86995ms|409.808|0.00|0s|
|10x10|11|9.976014ms|2165.457|0.45|13.266µs|
|15x15|11|71.874175ms|16980.570|5.73|232.07µs|
|20x20|11|596.526623ms|137849.810|45.91|1.866981ms|
