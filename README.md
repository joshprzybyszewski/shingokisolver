# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-20-2021

|Num Edges|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause (ns)|
|-:|-:|-:|-:|-:|-:|
|2x2|1|48.159µs|8.086|0.00|0s|
|5x5|11|402.845µs|45.724|0.00|0s|
|7x7|11|922.258µs|104.684|0.00|0s|
|10x10|11|9.387663ms|1354.891|0.18|4.39µs|
|15x15|11|76.691162ms|20356.461|7.64|381.611µs|
|20x20|11|5.553383297s|1559862.040|496.09|25.322102ms|
