# Unpredictable Data Flows

## Problems

For a test app, I had the point batch set to 1000 points.  This worked well for
when large import operations were processing.  But if the data stream slowed
down to a few per second, several minutes can go before the batch is large
enough to write.

Also when writing the point batch, how can I be sure another go routine isn't
trying to use it in `BatchPoints.AddPoint()`?

## Solution

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
