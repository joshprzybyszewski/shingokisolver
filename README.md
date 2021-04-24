# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-24-2021

|Num Edges|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause (ns)|
|-:|-:|-:|-:|-:|-:|
|2x2|1|48.158µs|13.211|0.00|0s|
|5x5|11|505.08µs|110.144|0.00|0s|
|7x7|11|1.680188ms|305.159|0.00|0s|
|10x10|11|6.856455ms|1447.711|0.18|6.418µs|
|15x15|11|50.683794ms|11503.874|3.82|402.087µs|
|20x20|11|398.806366ms|85852.584|28.55|1.943443ms|
