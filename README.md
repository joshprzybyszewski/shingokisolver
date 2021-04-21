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
|2x2|1|61.373µs|7.961|0.00|0s|
|5x5|11|352.489µs|45.521|0.00|0s|
|7x7|11|836.909µs|105.219|0.00|0s|
|10x10|11|8.453513ms|1346.837|0.18|5.658µs|
|15x15|11|73.302242ms|20319.771|8.64|626.212µs|
|20x20|11|5.864814585s|1556511.207|664.00|26.700185ms|
