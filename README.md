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
|5x5|easy|25|399.343µs|96.172|0.00|0s|
|5x5|medium|25|444.705µs|102.162|0.00|0s|
|5x5|hard|25|478.845µs|111.345|0.00|0s|
|7x7|easy|25|1.109209ms|232.930|0.00|0s|
|7x7|medium|25|1.052149ms|226.669|0.00|0s|
|7x7|hard|25|1.225188ms|281.271|0.00|0s|
|10x10|easy|25|2.380651ms|561.410|0.00|0s|
|10x10|medium|25|2.372194ms|601.148|0.00|0s|
|10x10|hard|25|3.378229ms|1234.306|0.00|0s|
|15x15|easy|25|7.52693ms|2261.206|0.04|979ns|
|15x15|medium|25|7.465625ms|1966.047|0.00|0s|
|15x15|hard|25|41.048529ms|26146.892|0.80|44.202µs|
|20x20|easy|25|17.530284ms|5236.225|0.28|100.149µs|
|20x20|medium|25|24.840437ms|11134.965|1.36|231.804µs|
|20x20|hard|25|1.558243388s|903780.442|61.68|14.11115ms|
|25x25|easy|25|50.544941ms|19898.019|1.12|98.189µs|
|25x25|medium|25|1.454518716s|836497.088|9.72|899.171µs|
