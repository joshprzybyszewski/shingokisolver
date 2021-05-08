# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 05-08-2021

_Intel(R) Core(TM) i5-3570 CPU @ 3.40GHz, 4 logical cores, 4 physical cores, 1 threads per core_

|Num Edges|Difficulty|Sample Size|Best Duration|Average Duration|Best Allocations (KB)|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|-:|-:|
|5x5|easy|200|65.071µs|232.947µs|41.312|54.230|0.00|0s|
|5x5|medium|200|89.624µs|265.674µs|41.711|57.943|0.00|0s|
|5x5|hard|200|77.948µs|276.698µs|42.289|60.012|0.00|0s|
|7x7|easy|124|135.803µs|464.362µs|79.414|106.039|0.00|0s|
|7x7|medium|123|168.497µs|479.832µs|79.688|106.732|0.00|0s|
|7x7|hard|123|152.33µs|678.633µs|79.008|137.626|0.00|0s|
|10x10|easy|50|536.419µs|826.389µs|162.562|214.906|0.00|0s|
|10x10|medium|50|272.839µs|857.075µs|160.586|217.928|0.00|0s|
|10x10|hard|50|610.186µs|1.255745ms|160.070|358.020|0.00|0s|
|15x15|easy|97|636.893µs|1.868746ms|376.953|611.001|0.00|0s|
|15x15|medium|50|718.798µs|1.887606ms|386.406|578.767|0.00|0s|
|15x15|hard|50|1.688295ms|8.652569ms|441.469|7333.006|0.54|85.34µs|
|20x20|easy|51|1.605868ms|3.775921ms|726.883|2403.278|0.10|29.891µs|
|20x20|medium|51|1.581474ms|3.794539ms|701.445|2428.588|0.10|29.271µs|
|20x20|hard|51|3.597126ms|122.150281ms|948.297|167799.785|13.10|1.511727ms|
|25x25|easy|52|2.645406ms|10.714063ms|1106.594|12684.376|0.69|64.906µs|
|25x25|medium|52|2.353704ms|23.485657ms|1129.531|32276.131|2.06|194.031µs|
|25x25|hard|77|13.176176ms|5.553725581s|13773.039|9089625.560|516.87|51.941972ms|
|30x30|easy|2|5.600254ms|23.810177ms|8400.141|32301.168|1.00|44.981µs|
|35x35|easy|1|293.082068ms|293.082068ms|512113.406|512113.406|62.00|13.766628ms|
|40x40|easy|1|984.363276ms|984.363276ms|1750759.586|1750759.586|176.00|38.874428ms|