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
|5x5|easy|25|412.107µs|83.563|0.00|0s|
|5x5|medium|25|421.049µs|87.843|0.00|0s|
|5x5|hard|25|463.072µs|88.623|0.00|0s|
|7x7|easy|25|1.076272ms|217.632|0.00|0s|
|7x7|medium|25|1.069027ms|206.059|0.00|0s|
|7x7|hard|25|1.388225ms|251.638|0.00|0s|
|10x10|easy|25|2.530788ms|559.705|0.00|0s|
|10x10|medium|25|3.045558ms|578.148|0.00|0s|
|10x10|hard|25|5.976311ms|1072.859|0.04|1.833µs|
|15x15|easy|25|12.62988ms|2673.391|0.16|6.626µs|
|15x15|medium|25|8.492091ms|1870.416|0.04|2.091µs|
|15x15|hard|25|141.041739ms|24264.430|4.56|209.612µs|
|20x20|easy|25|23.989072ms|5004.065|0.96|47.749µs|
|20x20|medium|25|31.38743ms|6494.407|1.36|57.2µs|
|20x20|hard|25|11.573641194s|1702298.609|343.96|13.930566ms|
|25x25|easy|25|158.83677ms|27491.232|4.76|209.457µs|
|25x25|medium|25|4.199782152s|740468.723|100.56|4.486642ms|
|25x25|hard|3|6m24.83969554s|60906591.448|18141.33|678.140397ms|
