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
|5x5|easy|50|745.827µs|106.583|0.00|0s|
|5x5|medium|50|757.286µs|111.236|0.00|0s|
|5x5|hard|50|852.16µs|122.505|0.00|0s|
|7x7|easy|50|1.84966ms|258.291|0.00|0s|
|7x7|medium|50|1.807912ms|251.583|0.00|0s|
|7x7|hard|50|2.131409ms|319.977|0.00|0s|
|10x10|easy|50|3.419776ms|607.520|0.00|0s|
|10x10|medium|50|3.499216ms|611.279|0.00|0s|
|10x10|hard|50|3.995545ms|794.031|0.00|0s|
|15x15|easy|50|8.831534ms|1919.767|0.00|0s|
|15x15|medium|50|8.918229ms|1792.963|0.00|0s|
|15x15|hard|50|14.853537ms|5485.020|0.08|8.581µs|
|20x20|easy|50|16.985063ms|3907.831|0.00|0s|
|20x20|medium|50|17.407853ms|4200.510|0.02|2.467µs|
|20x20|hard|50|179.851807ms|112204.417|1.70|149.37µs|
|25x25|easy|50|30.804949ms|8546.195|0.02|1.003µs|
|25x25|medium|50|34.324436ms|11432.252|0.10|6.218µs|
|25x25|hard|25|10.108149463s|5405117.675|39.64|5.117722ms|
