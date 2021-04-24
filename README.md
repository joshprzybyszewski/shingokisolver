# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-24-2021

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause (ns)|
|-:|-:|-:|-:|-:|-:|
|5x5|easy|26|395.174µs|89.107|0.00|0s|
|5x5|medium|26|405.068µs|93.618|0.00|0s|
|5x5|hard|26|405.183µs|93.717|0.00|0s|
|7x7|easy|26|1.059511ms|227.793|0.00|0s|
|7x7|medium|26|1.166019ms|223.474|0.00|0s|
|7x7|hard|26|1.382657ms|259.949|0.00|0s|
|10x10|easy|7|3.542988ms|639.730|0.00|0s|
|10x10|medium|5|3.865333ms|620.036|0.00|0s|
|10x10|hard|5|5.048951ms|977.636|0.00|0s|
|15x15|easy|4|12.548349ms|2396.385|0.25|10.918µs|
|15x15|medium|4|14.715121ms|3028.678|0.25|10.046µs|
|15x15|hard|4|46.012142ms|9566.773|3.25|139.971µs|
