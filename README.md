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
|2x2|1|66.955µs|13.211|0.00|0s|
|5x5|11|461.335µs|106.238|0.00|0s|
|7x7|11|1.636297ms|290.864|0.00|0s|
|10x10|11|7.230168ms|1461.398|0.18|7.99µs|
|15x15|11|51.32099ms|11418.047|3.36|239.116µs|
|20x20|11|404.828038ms|86138.817|27.36|2.343069ms|
|25x25|11|1.185418182s|304917.551|75.73|4.395992ms|
