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
|2x2|1|80.408µs|9.055|0.00|0s|
|5x5|11|488.442µs|57.092|0.00|0s|
|7x7|11|872.094µs|132.180|0.00|0s|
|10x10|11|7.723614ms|1349.626|0.36|13.98µs|
|15x15|11|52.986534ms|13966.004|5.91|381.823µs|
|20x20|11|4.539295972s|1169188.070|478.45|29.174743ms|
