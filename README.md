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
|5x5|easy|50|772.94µs|109.906|0.00|0s|
|5x5|medium|50|862.92µs|114.207|0.00|0s|
|5x5|hard|50|868.833µs|128.105|0.00|0s|
|7x7|easy|50|1.842534ms|264.535|0.00|0s|
|7x7|medium|50|1.851079ms|257.137|0.00|0s|
|7x7|hard|50|2.087462ms|328.686|0.00|0s|
|10x10|easy|50|3.228031ms|613.231|0.00|0s|
|10x10|medium|50|3.509816ms|622.418|0.00|0s|
|10x10|hard|50|3.984722ms|808.854|0.00|0s|
|15x15|easy|50|8.967627ms|1931.356|0.00|0s|
|15x15|medium|50|8.600977ms|1821.697|0.00|0s|
|15x15|hard|50|14.764542ms|5510.361|0.08|6.172µs|
|20x20|easy|50|16.932936ms|3928.870|0.00|0s|
|20x20|medium|50|17.484775ms|4187.171|0.00|0s|
|20x20|hard|50|284.030989ms|166628.684|2.66|298.241µs|
|25x25|easy|50|30.88171ms|8724.454|0.00|0s|
|25x25|medium|50|34.301038ms|11702.271|0.04|2.004µs|
|25x25|hard|25|9.988375741s|5427699.240|40.76|5.324372ms|
