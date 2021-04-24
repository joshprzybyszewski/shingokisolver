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
|2x2|1|45.4µs|13.609|0.00|0s|
|5x5|11|551.354µs|96.494|0.00|0s|
|7x7|11|1.430845ms|243.632|0.00|0s|
|10x10|11|5.328029ms|943.896|0.09|2.007µs|
|15x15|11|38.655322ms|7006.863|2.00|74.93µs|
|20x20|11|269.206666ms|48859.119|15.45|1.458867ms|
|25x25|11|690.948423ms|145390.564|36.45|1.435432ms|
