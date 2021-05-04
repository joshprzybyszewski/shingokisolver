# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 05-04-2021

_Intel(R) Core(TM) i5-3570 CPU @ 3.40GHz_

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|50|732.515µs|107.095|0.00|0s|
|5x5|medium|50|799.011µs|111.242|0.00|0s|
|5x5|hard|50|823.589µs|122.047|0.00|0s|
|7x7|easy|50|1.799652ms|259.220|0.00|0s|
|7x7|medium|50|1.83664ms|249.057|0.00|0s|
|7x7|hard|50|2.077105ms|320.196|0.00|0s|
|10x10|easy|50|3.204782ms|602.794|0.00|0s|
|10x10|medium|50|3.400402ms|607.065|0.00|0s|
|10x10|hard|50|3.962109ms|792.933|0.00|0s|
|15x15|easy|50|8.693696ms|1919.200|0.00|0s|
|15x15|medium|50|8.572092ms|1799.610|0.00|0s|
|15x15|hard|50|15.298411ms|5940.707|0.08|10.27µs|
|20x20|easy|50|16.610619ms|3884.668|0.00|0s|
|20x20|medium|50|17.173852ms|4202.123|0.02|981ns|
|20x20|hard|50|161.365736ms|109921.626|1.58|175.903µs|
|25x25|easy|50|29.939173ms|8619.878|0.02|1.225µs|
|25x25|medium|50|32.955742ms|11639.914|0.12|19.286µs|
|25x25|hard|25|8.841767888s|5234772.060|35.96|4.538768ms|
