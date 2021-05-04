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
|5x5|easy|50|560.207µs|104.025|0.00|0s|
|5x5|medium|50|596.995µs|111.225|0.00|0s|
|5x5|hard|50|648.965µs|120.587|0.00|0s|
|7x7|easy|50|1.251879ms|249.668|0.00|0s|
|7x7|medium|50|1.242455ms|241.650|0.00|0s|
|7x7|hard|50|1.668787ms|342.766|0.00|0s|
|10x10|easy|50|2.574934ms|603.287|0.00|0s|
|10x10|medium|50|2.552489ms|630.071|0.00|0s|
|10x10|hard|50|3.292302ms|1106.355|0.00|0s|
|15x15|easy|50|7.9468ms|2540.750|0.00|0s|
|15x15|medium|50|7.202436ms|2046.831|0.00|0s|
|15x15|hard|50|27.311807ms|16789.788|0.44|24.005µs|
|20x20|easy|50|17.131594ms|5371.800|0.00|0s|
|20x20|medium|50|20.63857ms|7944.371|0.14|7.173µs|
|20x20|hard|50|1.914540392s|1061908.135|14.30|1.164864ms|
|25x25|easy|50|41.597037ms|16795.014|0.12|7.799µs|
|25x25|medium|50|281.292699ms|169381.616|3.12|182.855µs|
|25x25|hard|5|5.528196895s|3576633.631|157.20|12.022983ms|
