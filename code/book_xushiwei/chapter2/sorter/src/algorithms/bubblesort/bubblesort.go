package bubblesort

func BubbleSort(values []int){
 //   flag := true
    var tmp int;

    for i := len(values) -1 ; i > 0; i-- {
        for j := 0; j < i; j++ {
         if(values[j] > values[j+1]){
                tmp = values[j]
                values[j] = values[j+1]
                values[j + 1] = tmp
            }   
        }        
    }
}
