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
|2x2|1|71.471µs|9.461|0.00|0s|
|5x5|59|329.933µs|52.116|0.00|0s|
|7x7|15|866.105µs|166.615|0.00|0s|
|10x10|20|8.751219ms|2122.089|0.30|18.125µs|
|15x15|12|87.148893ms|38709.186|14.50|1.107551ms|
