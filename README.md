# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 05-06-2021

_Intel(R) Core(TM) i5-3570 CPU @ 3.40GHz_

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|200|236.142µs|54.337|0.00|0s|
|5x5|medium|200|262.642µs|57.902|0.00|0s|
|5x5|hard|200|279.876µs|59.935|0.00|0s|
|7x7|easy|124|471.996µs|106.395|0.00|0s|
|7x7|medium|123|485.003µs|106.677|0.00|0s|
|7x7|hard|123|687.332µs|137.421|0.00|0s|
|10x10|easy|50|848.354µs|214.655|0.00|0s|
|10x10|medium|50|901.537µs|218.482|0.00|0s|
|10x10|hard|50|1.308579ms|356.737|0.00|0s|
|15x15|easy|97|1.938959ms|610.587|0.00|0s|
|15x15|medium|50|1.943677ms|602.755|0.00|0s|
|15x15|hard|50|7.688788ms|6537.212|0.46|73.995µs|
|20x20|easy|51|3.78406ms|2467.855|0.12|16.284µs|
|20x20|medium|51|3.998614ms|2259.126|0.08|18.25µs|
|20x20|hard|51|107.199682ms|151189.124|11.82|1.196713ms|
|25x25|easy|51|11.318101ms|13989.538|0.84|86.048µs|
|25x25|medium|51|7.610279ms|7616.098|0.31|38.604µs|
|25x25|hard|75|5.751988479s|8916040.409|504.27|63.06952ms|
