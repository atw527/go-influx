Write points to Influx in bulk when needed and also trickle in the points when the flow slows down.

# The General Idea

Flow new points through a channel.  Receive items on the channel, but also have
a timeout option.  So any point will only sit in the batch for the maximum time
of the `delayChan`.

```go
delayChan := time.After(10 * time.Second)
pointsChan := make(chan *client.Point, 10000)

for {
    select {
    case point := <-pointsChan:
        /* write points */
    case <-delayChan:
        /* write points */
    }
}
```

# Tests on Tests

This whole project is a giant test, but there are integration tests to make
sure the environment is sane.

```bash
#[user]$
./scripts/test.sh
```

# Running Benchmark

Running `benchmark.sh` will take down Influx if running and restart with a clean
graph and run through the benchmarks.

```bash
#[user]$
./scripts/benchmark.sh
```

Watch the test results: http://localhost:4401/dashboard/db/playground?refresh=5s&orgId=1

Both single-stat modules should go green with the expected number of points when the test completes.

When done, run:

```bash
#[user]$
docker-compose down
```

# Known Issues

The best way I know of to stopping the flow of points is to
shove a nil pointer into the `points` channel.  I tried using a separate stop
channel, but select statements don't give priority to one case over another.  
The goal was to drain the `points` channel first, then wait on three conditions:

- points - write to `client.BatchPoints`
- sthap signal - flush `client.BatchPoints` to Influx and break
- timeout - flush `client.BatchPoints` to Influx and continue

But Go's select structure will select a case at random if more than one is
ready.  So when the `sthap` channel had a value, it would close with items still
in the `points` channel.
