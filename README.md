# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 05-03-2021

_Intel(R) Core(TM) i5-3570 CPU @ 3.40GHz_

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|50|814.361µs|109.986|0.00|0s|
|5x5|medium|50|845.055µs|113.796|0.00|0s|
|5x5|hard|50|887.523µs|127.552|0.00|0s|
|7x7|easy|50|1.93479ms|266.605|0.00|0s|
|7x7|medium|50|1.876116ms|255.692|0.00|0s|
|7x7|hard|50|2.262088ms|326.874|0.00|0s|
|10x10|easy|50|3.551191ms|614.276|0.00|0s|
|10x10|medium|50|3.601201ms|618.013|0.00|0s|
|10x10|hard|50|4.031859ms|812.658|0.00|0s|
|15x15|easy|50|8.874533ms|1938.941|0.00|0s|
|15x15|medium|50|8.789547ms|1806.292|0.00|0s|
|15x15|hard|50|15.383141ms|5698.622|0.08|3.814µs|
|20x20|easy|50|16.911164ms|3918.270|0.00|0s|
|20x20|medium|50|17.871035ms|4279.802|0.02|1.445µs|
|20x20|hard|50|300.208106ms|164565.506|3.78|319.153µs|
|25x25|easy|50|30.649999ms|8662.978|0.02|990ns|
|25x25|medium|50|35.473427ms|11903.493|0.08|4.214µs|
|25x25|hard|15|3.535354902s|1982143.208|43.33|3.183393ms|
