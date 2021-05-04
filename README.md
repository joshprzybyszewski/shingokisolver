# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 05-04-2021

_Intel(R) Core(TM) i5-3570 CPU @ 3.40GHz_

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|50|654.717µs|109.587|0.00|0s|
|5x5|medium|50|698.831µs|113.890|0.00|0s|
|5x5|hard|50|711.753µs|128.348|0.00|0s|
|7x7|easy|50|1.661434ms|263.243|0.00|0s|
|7x7|medium|50|1.749875ms|255.223|0.00|0s|
|7x7|hard|50|2.071389ms|327.707|0.00|0s|
|10x10|easy|50|3.288001ms|614.220|0.00|0s|
|10x10|medium|50|3.235757ms|616.755|0.00|0s|
|10x10|hard|50|3.837085ms|804.955|0.00|0s|
|15x15|easy|50|9.275434ms|1940.863|0.00|0s|
|15x15|medium|50|8.711917ms|1813.346|0.00|0s|
|15x15|hard|50|14.543982ms|5194.380|0.06|3.179µs|
|20x20|easy|50|16.933197ms|3905.135|0.00|0s|
|20x20|medium|50|17.485525ms|4101.870|0.00|0s|
|20x20|hard|50|290.047802ms|163579.378|2.84|487.52µs|
|25x25|easy|50|30.675832ms|8558.995|0.02|1.048µs|
|25x25|medium|50|34.994685ms|11673.703|0.06|3.156µs|
|25x25|hard|25|11.051573748s|5707426.137|40.68|6.021915ms|
