package main


import (
    "fmt"
    "time"
)

func work(ch chan int) {
    time.Sleep(1e9)
    ch <- 1
}

func main() {
    ch := make(chan int)  
    timeout := make(chan bool, 1)

    go work(ch)

    go func() {
        for {
            time.Sleep(1e9)
            timeout <- true
        }
    }()

    var timeoutCount int = 0
    for {
        select {
            case i := <-ch:
                fmt.Println("Got value:", i)
            case <-timeout:
                fmt.Println("timeout!")
                timeoutCount++
        }

        if timeoutCount > 3 {
            break
        }
    }
  
}
