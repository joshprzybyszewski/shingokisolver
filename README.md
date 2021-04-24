# Shingoki Solver

## What?
[Shingoki](https://www.puzzle-shingoki.com) is a fun logic puzzle. Check out the website to learn and play it.

## Why?

![Travis challenged me](https://user-images.githubusercontent.com/23204038/112846696-f1f1fb00-906b-11eb-9693-3130ce4e78d7.png)

## How?

Using golang, I've built a solver. You can see it execute on cached puzzles with `make run`. Below are the latest results from my machine.

</startResults>

#### Results from 04-24-2021

|Num Edges|Difficulty|Sample Size|Average Duration|Average Allocations (KB)|Average Garbage Collections|Average GC Pause (ns)|
|-:|-|-:|-:|-:|-:|-:|
|5x5|easy|25|493.359µs|89.169|0.00|0s|
|5x5|medium|25|560.54µs|94.296|0.00|0s|
|5x5|hard|25|546.944µs|93.827|0.00|0s|
|7x7|easy|25|1.339641ms|227.682|0.00|0s|
|7x7|medium|25|1.288164ms|219.669|0.00|0s|
|7x7|hard|25|1.82026ms|264.583|0.00|0s|
|10x10|easy|25|3.337271ms|605.910|0.00|0s|
|10x10|medium|25|3.660016ms|633.897|0.00|0s|
|10x10|hard|25|6.754915ms|1138.570|0.04|1.354µs|
|15x15|easy|25|20.636356ms|4301.091|0.64|27.167µs|
|15x15|medium|25|10.791405ms|2305.660|0.08|3.112µs|
|15x15|hard|25|267.467497ms|54134.801|12.36|533.738µs|
|20x20|easy|18|33.664303ms|5911.913|1.67|68.978µs|
|20x20|medium|17|46.336826ms|8360.928|2.53|112.336µs|
|20x20|hard|17|5.765709509s|976857.872|336.35|12.890081ms|
|25x25|easy|3|847.576276ms|166312.586|78.67|2.917524ms|
|25x25|medium|3|206.32982ms|38184.898|14.00|631.83µs|
|25x25|hard|1|3m28.732317054s|45311328.836|17060.00|713.204353ms|
