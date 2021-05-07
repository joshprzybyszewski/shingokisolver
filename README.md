# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 05-06-2021

_Intel(R) Core(TM) i5-3570 CPU @ 3.40GHz_

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|200|232.923µs|53.895|0.00|0s|
|5x5|medium|200|266.342µs|57.998|0.00|0s|
|5x5|hard|200|272.302µs|59.696|0.00|0s|
|7x7|easy|124|464.218µs|105.536|0.00|0s|
|7x7|medium|123|460.953µs|106.546|0.00|0s|
|7x7|hard|123|672.306µs|137.797|0.00|0s|
|10x10|easy|50|845.335µs|213.771|0.00|0s|
|10x10|medium|50|865.555µs|217.762|0.00|0s|
|10x10|hard|50|1.318388ms|354.198|0.00|0s|
|15x15|easy|97|1.862215ms|602.772|0.00|0s|
|15x15|medium|50|1.881001ms|627.623|0.00|0s|
|15x15|hard|50|9.715581ms|8658.567|0.72|144.758µs|
|20x20|easy|51|7.41151ms|11983.259|0.76|373.931µs|
|20x20|medium|51|7.602695ms|15829.578|1.20|145.486µs|
|20x20|hard|51|300.474535ms|428391.528|33.24|4.142838ms|
|25x25|easy|51|595.772571ms|899339.888|10.04|2.778623ms|
|25x25|medium|51|1.04394126s|1404577.270|15.53|3.296313ms|
|25x25|hard|52|14.200703531s|18268302.884|360.08|54.488474ms|
