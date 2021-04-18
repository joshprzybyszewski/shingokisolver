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
|2x2|1|59.711µs|9.461|0.00|0s|
|5x5|107|382.703µs|52.546|0.00|0s|
|7x7|104|1.283603ms|200.537|0.00|0s|
|10x10|101|18.343723ms|3707.658|1.18|35.548µs|
