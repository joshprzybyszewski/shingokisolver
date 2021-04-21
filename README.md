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
|2x2|1|77.513µs|8.086|0.00|0s|
|5x5|11|1.166062ms|46.058|0.00|0s|
|7x7|11|665.041µs|105.496|0.00|0s|
|10x10|11|7.595101ms|1356.197|0.18|9.228µs|
|15x15|11|72.774476ms|20376.112|7.55|298.834µs|
|20x20|11|5.605869917s|1558540.170|498.45|25.03296ms|
