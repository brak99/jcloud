## Jumpcloud Interview Assignment


----------


### Setup
Each step is in it's own directory and has a main.go.  To run each step one simply has to run:

    > go run main.go

Each step except the first runs a server on port 8088 so step 2 can be reached at http://localhost:8088/hash?password=angryMonkey.

### Notes
If this were a real service there are things I would do different.  However in the interests of time and since this is just a sample there were some tradeoffs made.

- The idStore is "ok" for an example.  A real data store would be much better instead of the simple and limited one here.
- The stats could be better.  When in a production environment, more stats the better for monitoring.  There are better stat packages out there that do a great job of monitoring and alerting.
- The workQueue while adequate but could be better.
- There's no security on any of these endpoints.  Especially the shutdown one.  
- Since we are talking passwords, encryption should be used.
- Etc