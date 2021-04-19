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
|2x2|1|99.14µs|9.461|0.00|0s|
|5x5|107|389.19µs|52.525|0.00|0s|
|7x7|104|1.29097ms|200.635|0.00|0s|
|10x10|101|18.267574ms|3706.396|1.32|49.869µs|
|15x15|102|8.828381456s|3223315.572|997.75|41.120258ms|
|20x20|12|15.036787502s|5144986.529|2237.08|80.910821ms|
