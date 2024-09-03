# go-node-server-benchmark-test

## Overview

This is a simple test to compare a the speed and memory consumption of a go and js server endpoint. For work we had something similar implemented with Node.js and I wanted to see the performance gains when using Go.

The API simply takes in an array in json format, filters the items and returns the filterd items. The goal was to minimize the memory footprint of the applications, to allow filtering for huge json data. For this I implemented a steaming approch, so each item is processed one by one.
I also included a `/status` endpoint, since we had some problems in the past dealing with health check failures under heavy load and I wanted to gain some insights here too.

To be fair with Node.js, since Node.Js cannot run code in parallel I used the "cluster" module to run the API on all cores.

The load test was done with 100 parallel requests and 50 iterations on a Macbook Pro 2019 with an Intel i7 6-Core 2.6GHz CPU. The heartbeat request was made every 5 seconds.

## Results

With 12 threads (utilizing all logical processors):

Go:

Filter request: Mean: 1.86734146s - Median: 1.81635645s - Max: 3.443674941s - Min: 386.940174ms
Heartbeat request: Mean: 80.81609ms - Median: 10.963387ms - Max: 352.184537ms - Min: 1.473372ms

Node:

(Heartbeat) Final: Mean: 3.110083141s - Median: 1.376866903s - Max: 16.185833569s - Min: 4.93769ms
Final: Mean: 8.340175742s - Median: 8.027013732s - Max: 18.5924755s - Min: 1.033271597s

With 6 threads:

Go:

(Heartbeat) Final: Mean: 40.562995ms - Median: 12.623554ms - Max: 217.326776ms - Min: 2.778162ms
Final: Mean: 2.305342424s - Median: 2.074621453s - Max: 4.373371328s - Min: 464.529213ms

Node:

(Heartbeat) Final: Mean: 1.193961734s - Median: 951.717178ms - Max: 15.531011697s - Min: 1.879371ms
Final: Mean: 9.591364943s - Median: 8.920234559s - Max: 22.272940865s - Min: 600.678025ms

With 2 threads:

Go:

(Heartbeat) Final: Mean: 597.767714ms - Median: 382.351985ms - Max: 2.252416645s - Min: 11.297192ms
Final: Mean: 5.918676677s - Median: 5.118359501s - Max: 14.060162721s - Min: 507.217983ms

With 1 thread:
