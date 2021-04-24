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
|2x2|1|64.699µs|12.648|0.00|0s|
|5x5|11|582.926µs|97.188|0.00|0s|
|7x7|11|1.219036ms|268.479|0.00|0s|
|10x10|11|7.175004ms|1526.866|0.18|6.836µs|
|15x15|11|40.571497ms|9448.500|3.00|515.6µs|
|20x20|11|2.596785303s|538339.442|171.36|11.442594ms|
