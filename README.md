# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-24-2021

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|25|448.649µs|89.713|0.00|0s|
|5x5|medium|25|447.675µs|93.267|0.00|0s|
|5x5|hard|25|468.529µs|93.660|0.00|0s|
|7x7|easy|25|1.224308ms|231.923|0.00|0s|
|7x7|medium|25|1.24721ms|220.172|0.00|0s|
|7x7|hard|25|1.6889ms|261.374|0.00|0s|
|10x10|easy|25|3.095097ms|608.803|0.00|0s|
|10x10|medium|25|3.596642ms|625.175|0.00|0s|
|10x10|hard|25|6.822758ms|1184.120|0.04|1.414µs|
|15x15|easy|25|22.296249ms|4600.775|0.68|26.189µs|
|15x15|medium|25|9.706095ms|2303.052|0.16|5.2µs|
|15x15|hard|25|282.517351ms|54025.816|10.88|441.344µs|
|20x20|easy|25|31.915746ms|5650.408|1.16|59.714µs|
|20x20|medium|25|45.409152ms|8436.734|2.04|76.433µs|
|20x20|hard|25|28.04055632s|5069801.697|958.32|38.920586ms|
|25x25|easy|25|183.384661ms|38193.643|10.88|373.922µs|
|25x25|medium|25|9.108905517s|1858494.769|251.16|11.1625ms|
|25x25|hard|3|15m12.884375197s|169302135.792|51234.67|2.047683586s|
