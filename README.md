# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-25-2021

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|25|414.569µs|81.897|0.00|0s|
|5x5|medium|25|454.757µs|87.076|0.00|0s|
|5x5|hard|25|448.433µs|86.037|0.00|0s|
|7x7|easy|25|1.116435ms|216.539|0.00|0s|
|7x7|medium|25|1.126007ms|202.199|0.00|0s|
|7x7|hard|25|1.325636ms|244.767|0.00|0s|
|10x10|easy|25|2.622265ms|535.524|0.00|0s|
|10x10|medium|25|2.886998ms|557.241|0.00|0s|
|10x10|hard|25|5.700566ms|1021.832|0.04|1.831µs|
|15x15|easy|25|12.093469ms|2568.308|0.16|8.598µs|
|15x15|medium|25|8.370789ms|1828.152|0.04|1.674µs|
|15x15|hard|25|134.348865ms|22414.195|4.08|164.206µs|
|20x20|easy|25|23.682627ms|4828.108|0.84|35.441µs|
|20x20|medium|25|30.62885ms|6176.988|1.24|48.989µs|
|20x20|hard|25|8.716190585s|1080382.218|225.68|10.315545ms|
|25x25|easy|25|146.889727ms|25560.890|4.04|169.272µs|
|25x25|medium|25|3.842527598s|654892.437|87.40|3.872958ms|
|25x25|hard|3|5m58.592184806s|53474395.708|15305.00|597.432111ms|
