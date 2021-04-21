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
|2x2|1|90.951µs|9.055|0.00|0s|
|5x5|11|445.652µs|53.676|0.00|0s|
|7x7|11|1.102149ms|127.753|0.00|0s|
|10x10|11|8.711558ms|1679.332|0.55|14.658µs|
|15x15|11|74.781044ms|22324.729|9.64|799.808µs|
|20x20|11|6.200534405s|1742138.094|736.27|34.851937ms|
