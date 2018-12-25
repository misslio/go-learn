package main


import (
    "fmt"
)


func main() {
    ch := make(chan int, 1)  
    timeout := make(chan bool, 1)

    ch <- 1
 
    timeout <- true
   

    var timeoutCount int = 0
   
    select {
        case i:= <-ch:
            fmt.Println("Got value:", i)
        case <-timeout:
            fmt.Println("timeout!")
            timeoutCount++
    }

}
