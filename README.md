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
|5x5|easy|50|206.645µs|67.373|0.00|0s|
|5x5|medium|50|204.306µs|63.438|0.00|0s|
|5x5|hard|50|248.236µs|70.711|0.00|0s|
|7x7|easy|50|454.456µs|102.257|0.00|0s|
|7x7|medium|50|434.366µs|97.791|0.00|0s|
|7x7|hard|50|616.952µs|122.861|0.00|0s|
|10x10|easy|50|3.654551ms|6972.158|0.46|119.474µs|
|10x10|medium|50|1.307743ms|1592.587|0.08|16.323µs|
|10x10|hard|50|2.276954ms|3044.255|0.32|62.755µs|
|15x15|easy|50|3.63875ms|5844.970|0.60|164.546µs|
|15x15|medium|50|2.628653ms|4413.470|0.56|146.227µs|
|15x15|hard|50|12.821558ms|16170.723|1.32|298.773µs|
|20x20|easy|50|7.607079ms|14249.280|1.00|239.754µs|
|20x20|medium|50|4.169592ms|7757.678|0.60|155.862µs|
|20x20|hard|50|338.495615ms|473359.284|32.14|4.33315ms|
|25x25|easy|50|28.133478ms|45938.686|3.02|877.332µs|
|25x25|medium|50|39.588236ms|68162.789|4.64|809.405µs|
|25x25|hard|23|13.7624251s|19940972.607|1204.74|154.483626ms|
