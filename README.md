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
|2x2|1|95.378µs|8.133|0.00|0s|
|5x5|11|453.743µs|45.989|0.00|0s|
|7x7|11|811.382µs|106.673|0.00|0s|
|10x10|11|8.720509ms|1376.436|0.27|6.801µs|
|15x15|11|72.088909ms|20339.393|8.91|305.406µs|
|20x20|11|5.947619504s|1571765.308|669.64|26.741651ms|
