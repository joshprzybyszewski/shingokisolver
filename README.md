# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 05-05-2021

_Intel(R) Core(TM) i5-3570 CPU @ 3.40GHz_

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|50|203.39µs|52.430|0.00|0s|
|5x5|medium|50|228.771µs|56.369|0.00|0s|
|5x5|hard|50|277.532µs|56.589|0.00|0s|
|7x7|easy|50|475.106µs|102.712|0.00|0s|
|7x7|medium|50|453.839µs|98.315|0.00|0s|
|7x7|hard|50|636.527µs|122.810|0.00|0s|
|10x10|easy|50|648.602µs|595.046|0.04|8.865µs|
|10x10|medium|50|621.769µs|567.789|0.02|2.842µs|
|10x10|hard|50|2.366452ms|2552.017|0.26|62.789µs|
|15x15|easy|50|4.140292ms|5085.206|0.52|141.424µs|
|15x15|medium|50|4.739465ms|7936.229|0.64|140.517µs|
|15x15|hard|50|13.232255ms|17013.244|1.26|233.901µs|
|20x20|easy|50|5.438644ms|10278.225|0.78|222.4µs|
|20x20|medium|50|5.382406ms|11103.297|0.92|221.826µs|
|20x20|hard|50|322.722698ms|435592.139|29.22|3.038454ms|
|25x25|easy|50|22.72108ms|41470.013|2.86|516.343µs|
|25x25|medium|50|49.525739ms|58735.543|4.24|1.277834ms|
|25x25|hard|24|24.491601749s|34779108.233|2113.96|250.568288ms|
