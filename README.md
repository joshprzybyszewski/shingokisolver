# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-20-2021

|Num Edges|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause (ns)|
|-:|-:|-:|-:|-:|-:|
|2x2|1|84.545µs|7.898|0.00|0s|
|5x5|11|404.295µs|44.504|0.00|0s|
|7x7|11|1.133645ms|126.812|0.00|0s|
|10x10|11|13.871393ms|2121.141|0.36|14.508µs|
|15x15|11|195.375509ms|55003.029|25.73|1.998706ms|
|20x20|11|14.649298779s|3458576.045|1380.91|48.177693ms|
