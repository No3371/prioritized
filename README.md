# prioritized
A simple work priority implementation.

## Why
I was making a ws server and was stress testing, it turned out can't properly serve clients connected before clients connected later.

This simple implementation allows users to queue jobs into high/low priority channel. Jobs in high priority channels will always be dequeued before low priority jobs.
