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
|5x5|easy|50|471.315µs|107.095|0.00|0s|
|5x5|medium|50|531.737µs|111.576|0.00|0s|
|5x5|hard|50|560.506µs|124.316|0.00|0s|
|7x7|easy|50|1.212409ms|260.533|0.00|0s|
|7x7|medium|50|1.19916ms|251.861|0.00|0s|
|7x7|hard|50|1.452728ms|324.408|0.00|0s|
|10x10|easy|50|3.025905ms|605.543|0.00|0s|
|10x10|medium|50|2.994904ms|605.736|0.00|0s|
|10x10|hard|50|3.522911ms|806.998|0.00|0s|
|15x15|easy|50|8.508744ms|1944.766|0.00|0s|
|15x15|medium|50|8.394647ms|1794.614|0.00|0s|
|15x15|hard|50|16.684623ms|6312.968|0.02|1.059µs|
|20x20|easy|50|16.381779ms|3959.535|0.00|0s|
|20x20|medium|50|17.27952ms|4221.413|0.00|0s|
|20x20|hard|50|374.941735ms|189206.937|0.92|143.742µs|
|25x25|easy|50|31.231593ms|8837.741|0.00|0s|
|25x25|medium|50|37.915497ms|13387.654|0.06|3.175µs|
|25x25|hard|10|9.074665341s|3912919.274|36.00|4.98836ms|
