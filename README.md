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

# Testing

Running `reset.sh` will take down Influx if running and restart with a clean
graph and run through the benchmarks.

```bash
#[user]$
./reset.sh
```

Watch the test results: http://localhost:4401/dashboard/db/playground?refresh=5s&orgId=1

Both single-stat modules should go green with the expected number of points when the test completes.

When done testing, run:

```bash
#[user]$
docker-compose down
```

# Known Issues

Currently it's not shutting down correctly.  In the output, the last point batch will be less than what's sent to the channel.

Putting in a manual delay between the last `AddPoint()` and `Stahp()` will fix the problem.  However, I want the channels to block each other while things finish.
