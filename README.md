# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-25-2021

_Intel(R) Core(TM) i5-3570 CPU @ 3.40GHz_

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|50|514.39µs|95.557|0.00|0s|
|5x5|medium|50|586.767µs|104.553|0.00|0s|
|5x5|hard|50|613.405µs|109.685|0.00|0s|
|7x7|easy|50|1.43465ms|239.585|0.00|0s|
|7x7|medium|50|1.341624ms|229.398|0.00|0s|
|7x7|hard|50|1.637174ms|320.160|0.00|0s|
|10x10|easy|50|2.406236ms|567.274|0.00|0s|
|10x10|medium|50|2.34822ms|594.756|0.00|0s|
|10x10|hard|50|3.399003ms|1341.691|0.00|0s|
|15x15|easy|50|7.674662ms|2506.502|0.00|0s|
|15x15|medium|50|6.973568ms|2017.902|0.00|0s|
|15x15|hard|50|43.339062ms|27397.388|0.48|189.735µs|
|20x20|easy|50|17.992177ms|5679.180|0.12|6.301µs|
|20x20|medium|50|21.858473ms|8656.406|0.50|38.38µs|
|20x20|hard|50|2.275379931s|1390197.931|32.20|3.552498ms|
|25x25|easy|50|48.626982ms|22290.581|0.40|22.854µs|
|25x25|medium|50|758.824275ms|438928.311|4.64|357.032µs|
|25x25|hard|3|4.040006905s|2649061.138|322.67|37.343271ms|
