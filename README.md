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
|2x2|1|43.663µs|9.461|0.00|0s|
|5x5|59|371.087µs|51.914|0.00|0s|
|7x7|15|1.154762ms|167.512|0.00|0s|
|10x10|20|9.379648ms|2126.325|0.30|11.733µs|
|15x15|12|89.361461ms|38821.718|14.75|505.228µs|
