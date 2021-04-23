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
|2x2|1|79.042µs|10.219|0.00|0s|
|5x5|11|527.843µs|72.526|0.00|0s|
|7x7|11|1.502013ms|318.414|0.00|0s|
|10x10|11|26.349889ms|7169.555|3.00|113.935µs|
|15x15|11|2.578114988s|743048.937|315.82|12.533358ms|
|20x20|11|4m24.825507892s|55171348.145|21603.18|1.000443931s|
