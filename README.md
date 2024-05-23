# databaseChunkerWithThreadSearcher
we chunk the "database" into a number of groups, we then spread those groups over a worker pool looking for a match, if we dont find a match we gracefully close down.
