package main

import (
    "os"
    "fmt"
    "simplemath"
    "strconv"
)

var Usage = func() {
    fmt.Println("USAGE: calc command [arguments] ...")
    fmt.Println("\nThe commands are:\n\tadd\tAddition of two values.")
    fmt.Println("\tsqrt\tSquare root of a non-negative value")
}

func main() {
    //Args保管了命令行参数，第一个是程序名
    args := os.Args
    if args == nil || len(args) < 3 {
        Usage()
        return
    }
    fmt.Println(args)
    switch args[1] {
        case "add":
            if len(args) != 4 {
                fmt.Println("USAGE: calc add <integer1> <integer2>")
                return
            }
            v1, err1 := strconv.Atoi(args[2])
            v2, err2 := strconv.Atoi(args[3])
            if err1 != nil || err2 != nil {
                
            }
            ret := simplepath.Add(v1, v2)
            fmt.Println("Result: ", ret)
        case "sqrt":
            if len(args) != 3 {
                fmt.Println("USAGE: calc add <integer1> <integer2>")
                return
            }
            v, err := strconv.Atoi(args[2])
            if err != nil {
                fmt.Println("USAGE: calc sqrt <integer>")
                return
            }
            ret := simplepath.Sqrt(v)
            fmt.Println("Result: ", ret)
        default:
            Usage()
    }
}
