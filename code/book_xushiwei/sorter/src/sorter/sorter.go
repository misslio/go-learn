package main

import (
    "flag"
    "fmt"
    "os"
    "bufio"
    "strconv"
    "algorithms/bubblesort"
)


var infile *string = flag.String("i", "infile", "File contains values for sorting")
var outfile *string = flag.String("o", "outfile", "File to receive sorted values")
var algorithm *string = flag.String("a", "qsort", "Sort algorithm")


func readValues(infile string)(values []int, err error) {  
    //只读模式打开文件  
    file, err := os.Open(infile)
    if err != nil {
        fmt.Println("Failed to open input file ", infile)
        return    
    }

    defer file.Close()

    //也可用bufio.Scanner来读取
    br := bufio.NewReader(file)
    values = make([]int, 0)
    
    for {
        line, isPrefix, err1 := br.ReadLine()
        if err1 != nil {
            break;
        }

        if isPrefix {
            return
        }
        str := string(line)
        value, err1 := strconv.Atoi(str)
        if err1 != nil {
            err = err1
            return
        }

        values = append(values, value)
    }

    return
}

func writeValues(values []int, outfile string) error {
    //以O_RDWR, O_CREAT, O_TRUNC模式创建文件
    file, err := os.Create(outfile)
    if err != nil {
        fmt.Println("Failed to createa outfile ", outfile)
        return err
    }

    defer file.Close()

    for _, value := range values {
        str := strconv.Itoa(value)
        file.WriteString(str + "\n") //类似file.Write,但接受一个字符串参数，方便使用
    }

    return nil
}

func main() {
    flag.Parse()

    if infile != nil {
        fmt.Println("infile =", *infile, "outfile =", *outfile, "algorithm =",
            *algorithm)
    }

    values, _:= readValues(*infile)

    fmt.Println(values)
    bubblesort.BubbleSort(values)
    fmt.Println(values)
    writeValues(values, *outfile)
}
