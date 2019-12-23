Go Worker Example
===

For episode [14](https://codeandship.rocks/podcast/14).

```
go run main.go
```

# Challenges:

1. Remove or comment out line 108 and observe the changes in the log output.
2. Remove or comment out line 108 and 112 and observe the changes in program execution.
3. Think about why we do need a dedicated interruped handler?
4. Add another goroutine to produce more work.
5. Use a buffered job queue and print the length of `jobQueue` when exiting.
6. Exit after a given timeout using another goroutine. `os.Exit(0)` is not allowed, the `time` package is your friend here.
7. Use a context instead of the `stop` channel.