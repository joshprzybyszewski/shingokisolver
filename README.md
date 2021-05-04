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
|5x5|easy|50|594.679µs|101.906|0.00|0s|
|5x5|medium|50|633.913µs|108.600|0.00|0s|
|5x5|hard|50|629.96µs|116.638|0.00|0s|
|7x7|easy|50|1.340851ms|235.837|0.00|0s|
|7x7|medium|50|1.345164ms|225.610|0.00|0s|
|7x7|hard|50|1.634597ms|299.770|0.00|0s|
|10x10|easy|50|2.976175ms|538.148|0.00|0s|
|10x10|medium|50|3.153786ms|548.700|0.00|0s|
|10x10|hard|50|3.666727ms|737.621|0.00|0s|
|15x15|easy|50|7.341444ms|1810.236|0.00|0s|
|15x15|medium|50|7.097797ms|1596.448|0.00|0s|
|15x15|hard|50|13.868836ms|6284.446|0.12|6.241µs|
|20x20|easy|50|12.71696ms|3550.214|0.00|0s|
|20x20|medium|50|14.227868ms|4520.928|0.06|2.999µs|
|20x20|hard|50|527.35525ms|283294.391|6.18|611.466µs|
|25x25|easy|50|23.610058ms|8657.903|0.04|2.906µs|
|25x25|medium|50|38.408075ms|18542.389|0.40|26.195µs|
|25x25|hard|5|1.251973295s|785849.655|53.60|4.208929ms|
