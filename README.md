# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-25-2021

_Intel(R) Core(TM) i5-3570 CPU @ 3.40GHz_

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|50|345.081µs|95.294|0.00|0s|
|5x5|medium|50|390.415µs|104.934|0.00|0s|
|5x5|hard|50|406.752µs|109.612|0.00|0s|
|7x7|easy|50|957.262µs|237.313|0.00|0s|
|7x7|medium|50|878.029µs|230.280|0.00|0s|
|7x7|hard|50|1.132ms|320.329|0.00|0s|
|10x10|easy|25|2.292988ms|557.699|0.00|0s|
|10x10|medium|44|2.400898ms|605.510|0.00|0s|
|10x10|hard|25|3.445518ms|1284.954|0.00|0s|
|15x15|easy|50|8.075881ms|2527.244|0.00|0s|
|15x15|medium|25|7.244041ms|1973.797|0.00|0s|
|15x15|hard|25|35.298641ms|24656.237|0.00|0s|
|20x20|easy|28|17.375506ms|5475.946|0.00|0s|
|20x20|medium|27|22.060404ms|10089.351|0.00|0s|
|20x20|hard|26|1.43115375s|990118.547|2.27|148.945µs|
|25x25|easy|27|46.353767ms|18876.300|0.00|0s|
|25x25|medium|27|1.172287398s|782207.533|1.52|105.348µs|
|25x25|hard|5|3m34.908844939s|112604444.645|421.00|59.570361ms|
