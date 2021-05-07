# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 05-07-2021

_Intel(R) Core(TM) i5-3570 CPU @ 3.40GHz_
_4 logical cores_
_4 physical cores_
_1 threads per core_

|Num Edges|Difficulty|Sample Size|Best Duration|Average Duration|Best Allocations (KB)|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|-:|-:|
|5x5|easy|200|80.223µs|223.349µs|10849.867|54.249|0.00|0s|
|5x5|medium|200|77.263µs|255.415µs|11572.008|57.860|0.00|0s|
|5x5|hard|200|78.741µs|267.539µs|11995.195|59.976|0.00|0s|
|7x7|easy|124|145.374µs|436.356µs|13138.203|105.953|0.00|0s|
|7x7|medium|123|153.756µs|453.792µs|13110.430|106.589|0.00|0s|
|7x7|hard|123|219.805µs|645.897µs|16900.695|137.404|0.00|0s|
|10x10|easy|50|368.215µs|817.693µs|10736.375|214.727|0.00|0s|
|10x10|medium|50|372.447µs|851.782µs|10937.023|218.740|0.00|0s|
|10x10|hard|50|614.652µs|1.325512ms|18036.516|360.730|0.00|0s|
|15x15|easy|97|666.873µs|1.812792ms|58897.367|607.189|0.00|0s|
|15x15|medium|50|821.753µs|1.776363ms|28943.969|578.879|0.00|0s|
|15x15|hard|50|1.467422ms|7.716184ms|329705.297|6594.106|0.44|83.881µs|
|20x20|easy|51|1.646422ms|3.678649ms|123714.734|2425.779|0.10|18.232µs|
|20x20|medium|51|1.638269ms|3.870252ms|124037.258|2432.103|0.10|23.216µs|
|20x20|hard|51|2.755005ms|162.31287ms|11509386.883|225674.253|17.96|2.425415ms|
|25x25|easy|51|2.057323ms|10.451354ms|629709.906|12347.253|0.73|161.615µs|
|25x25|medium|51|2.104661ms|26.664096ms|1555579.672|30501.562|1.82|681.66µs|
|25x25|hard|76|14.039878ms|5.926596704s|715952173.430|9420423.335|540.29|57.181294ms|
