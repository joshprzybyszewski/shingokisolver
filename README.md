# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-23-2021

|Num Edges|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause (ns)|
|-:|-:|-:|-:|-:|-:|
|2x2|1|76.265µs|15.875|0.00|0s|
|5x5|11|719.973µs|112.722|0.00|0s|
|7x7|11|2.090608ms|377.904|0.00|0s|
|10x10|11|13.18223ms|1990.376|0.36|14.041µs|
|15x15|11|229.686921ms|41195.443|14.45|1.055876ms|
|20x20|11|44.921049914s|7530302.848|2485.18|138.255561ms|
