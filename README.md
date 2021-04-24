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
|2x2|1|108.292µs|13.883|0.00|0s|
|5x5|11|757.685µs|139.976|0.00|0s|
|7x7|11|1.648815ms|388.660|0.00|0s|
|10x10|11|9.715ms|1952.685|0.36|10.778µs|
|15x15|11|69.737382ms|15239.889|5.36|203.907µs|
|20x20|11|576.287238ms|122749.622|41.27|1.411222ms|
