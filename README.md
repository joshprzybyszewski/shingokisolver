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
|2x2|1|38.022µs|13.609|0.00|0s|
|5x5|11|571.608µs|95.036|0.00|0s|
|7x7|11|1.368588ms|227.580|0.00|0s|
|10x10|11|5.131716ms|897.665|0.09|2.008µs|
|15x15|11|35.69262ms|6737.759|1.82|85.01µs|
|20x20|11|254.156415ms|47263.849|14.82|497.808µs|
|25x25|11|1.508781847s|297822.812|81.45|2.966065ms|
