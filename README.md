# Intent-based decentralized orchestration for green energy-aware provisioning of fog-native workflows

Mays AL-Naday, Tom Goethals, Bruno Volckaert


## Evaluation parameters and their settings

---

### Table 1: Algorithm settings
| Parameter   | Setting     |
|:-----------|:-----------|
|error gaps   | $\epsilon_{pri}:10^{-2}, \epsilon_{dual}:10^{-3}$   |
| $\rho$, $\zeta$, $\kappa_{inc}$, $\kappa_{dec}$ | 1, 10, 2, 2|


### Table 2: Analytical and experimental settings

| Parameter   | Setting (Analytical)    | Setting (Experimental)    |
|:-----------------------|:-------------------|:-------------------------|
|Number of workflows| 25 | 1|
|Workflow popularity distribution |  Zipf(0.8) | Normal ($\mu$=0.025, $\sigma$=0.025)|
|Number of microservices per workflow | 5 | 5|
|Workflow graph | Hub & Spoke (H&S), Chain | H&S, Chain|
|Average task size per microservice [mCPU] | [100 - 200] | 1|
|Average input/output data per microservice [Mb] | [0.4, 8.0] | 0.03 |
|Average response time tolerance [msec] | [30 - 150] | 35 |
|Fog nodes per tier | {t<sub>1</sub>:2, t<sub>2</sub>:3, t<sub>3</sub>:4} | {t<sub>1</sub>:[1-2], t<sub>2</sub>:2, t<sub>3</sub>:3} |
|Number of users per access node | [500, 2000] | 1 |
|Average request rate per access node [request/s]| 1000 | 40 |
|Average CPU capacity per node [mCPUs] | {t<sub>1</sub>:[10<sup>6</sup> - 10<sup>7</sup>], t<sub>2</sub>:[10<sup>5</sup> - 10<sup>6</sup>], t<sub>3</sub>:[10<sup>4</sup> - 10<sup>5</sup>]} | {t<sub>1</sub>:5000, t<sub>2</sub>:3000, t<sub>3</sub>:1000}|
|Average CPU energy price per node [PpmCPU] | {t<sub>1</sub>:[10<sup>-5</sup> - 10<sup>-4</sup>], t<sub>2</sub>:[10<sup>-4</sup> - 10<sup>-3</sup>], t<sub>3</sub>:[10<sup>-3</sup> - 10<sup>-2</sup>]} |-- |
|Average fraction of supply of green energy | {t<sub>1</sub>:[0.8-0.9], t<sub>2</sub>:[0.7-0.8], t<sub>3</sub>:[0.6-0.7]} | {t<sub>1</sub>:[0.8-0.9], t<sub>2</sub>:[0.7-0.8], t<sub>3</sub>:[0.6-0.7]}|
|Average bandwidth capacity per path [Mb/s] | {t<sub>1</sub>:[10<sup>5</sup> - 10<sup>6</sup>], t<sub>2</sub>:[10<sup>4</sup> - 10<sup>5</sup>], t<sub>3</sub>:[10<sup>3</sup> - 10<sup>4</sup>]} | {t<sub>1</sub>:1000, t<sub>2</sub>:1000, t<sub>3</sub>:100} |
|Average link length [Km] | {t<sub>1</sub>:[10-100], t<sub>2</sub>:[1-10], t<sub>3</sub>:[0.5-1]} | {t<sub>1</sub>:10, t<sub>2</sub>:10, t<sub>3</sub>:10} |

## Experimental Evaluation

---

Experimental evaluations are run on fog nodes with a Gigabit Ethernet connection, an Intel i5 9400 processor, and 32GiB RAM. A unit of processing work (mCPU) is defined as the CPU time required to bubble sort 1000 integers in Golang, which equals $6/5$ physical milliCPU on the evaluation nodes. A tiered architecture is simulated by limiting the total CPU budget of all microservices on each node.


