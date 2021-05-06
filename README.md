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
|5x5|easy|25|216.684µs|54.338|0.00|0s|
|5x5|medium|25|213.054µs|51.767|0.00|0s|
|5x5|hard|25|286.194µs|61.592|0.00|0s|
|7x7|easy|25|415.757µs|105.824|0.00|0s|
|7x7|medium|25|381.553µs|102.386|0.00|0s|
|7x7|hard|25|544.189µs|122.085|0.00|0s|
|10x10|easy|25|810.584µs|878.948|0.08|11.599µs|
|10x10|medium|25|734.204µs|567.983|0.04|9.726µs|
|10x10|hard|25|1.184387ms|436.592|0.04|6.567µs|
|15x15|easy|25|1.640459ms|646.724|0.00|0s|
|15x15|medium|25|1.778146ms|595.574|0.00|0s|
|15x15|hard|25|9.61487ms|8657.066|0.68|114.539µs|
|20x20|easy|25|2.859009ms|1997.595|0.12|18.979µs|
|20x20|medium|25|5.655462ms|3913.640|0.32|56.266µs|
|20x20|hard|25|236.659492ms|253960.126|25.52|9.165234ms|
|25x25|easy|25|18.489244ms|22165.468|2.20|870.821µs|
|25x25|medium|25|46.802132ms|65801.131|5.16|807.294µs|
|25x25|hard|11|4.284155646s|6146465.084|515.64|67.854976ms|
