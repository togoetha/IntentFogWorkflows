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
|:-----------------------|:-------------------------|:-------------------|
|Number of workflows| 25 | 1|
|Workflow popularity distribution |  Zipf(0.8) | Normal ($\mu$=0.025, $\sigma$=0.025)|
|Number of microservices per workflow | 5 | 5|
|Workflow graph | Hub & Spoke (H&S), Chain | H&S, Chain|
|Average task size per microservice [mCPU] | [100 - 200] | 1|
|Average input/output data per microservice [Mb] | [0.4, 8.0] | 0.03 |
|Average response time tolerance [msec] | [30 - 150] | 35 |
|Fog nodes per tier | { $t_1:$ 2, $t_2:$ 3, $t_3:$ 4} | { $t_1:$ 1-2, $t_2:$ 2, $t_3:$ 3} |
|Number of users per access node | [500, 2000] | 1 |
|Average request rate per access node [request/s]| 1000 | 40 |
|Average CPU capacity per node [mCPUs] | { $t_1:$ [10^6 - 10$^7$], $t_2:$ [10$^5$ - 10$^6$], $t_3:$ [10$^4$ - 10$^5$] } | { $t_1:$ 5000, $t_2:$ 3000, $t_3:$ 1000}|
|Average CPU energy price per node [PpmCPU] | $\{t_1:$ [10$^{-5}$ - 10$^{-4}$], $t_2:$ [10$^{-4}$ - 10$^{-3}$], $t_3:$[10$^{-3}$ - 10$^{-2}$] $\}$ |-- |
|Average fraction of supply of green energy | $\{t_1:$ [0.8-0.9], $t_2:$ [0.7-0.8], $t_3:$ [0.6-0.7] $\}$ | $\{t_1:$ [0.8-0.9], $t_2:$ [0.7-0.8], $t_3:$ [0.6-0.7] $\}$|
|Average bandwidth capacity per path [Mb/s] | $\{t_1:$ [10$^5$ - 10$^6$], $t_2:$ [10$^4$ - 10$^5$], $t_3:$ [10$^3$ - 10$^4$] $\}$ | $\{t_1:$ 1000, $t_2:$ 1000, $t_3:$ 100 $\}$ |
|Average link length [Km] | $\{t_1:$ [10-100], $t_2:$ [1-10], $t_3:$ [0.5-1] $\}$ | $\{t_1:$ 10, $t_2:$ 10, $t_3:$ 10 $\}$|



