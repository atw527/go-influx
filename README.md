Streaming points into Influx can go at a faster throughput if the point batch is larger.

For a test app, I had the point batch set to 1000 points, but this app had an inconsistent data flow.  Idle chatter can be as low as a few points per second.
