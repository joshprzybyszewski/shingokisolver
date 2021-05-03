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
|5x5|easy|50|485.327µs|105.036|0.00|0s|
|5x5|medium|50|497.633µs|113.546|0.00|0s|
|5x5|hard|50|536.363µs|124.052|0.00|0s|
|7x7|easy|50|1.281013ms|262.215|0.00|0s|
|7x7|medium|50|1.244136ms|257.051|0.00|0s|
|7x7|hard|50|1.512787ms|366.297|0.00|0s|
|10x10|easy|50|2.403451ms|640.644|0.00|0s|
|10x10|medium|50|2.43429ms|666.116|0.00|0s|
|10x10|hard|50|3.378052ms|1330.407|0.00|0s|
|15x15|easy|50|8.369696ms|2745.453|0.00|0s|
|15x15|medium|50|7.297772ms|2133.993|0.00|0s|
|15x15|hard|50|32.757112ms|21538.512|0.40|25.25µs|
|20x20|easy|50|17.559365ms|5852.889|0.00|0s|
|20x20|medium|50|21.572674ms|8529.683|0.10|69.084µs|
|20x20|hard|50|2.032246175s|1231065.337|16.42|2.053403ms|
|25x25|easy|50|46.77189ms|21211.021|0.12|7.241µs|
|25x25|medium|50|433.944644ms|256517.041|3.60|243.662µs|
|25x25|hard|5|7.719202805s|4934137.653|204.60|20.331256ms|
