package main

import(
	"fmt"
	"os"
	"bufio"
	"time"
	"flag"
	"io"
	// "compress/gzip"
)

var filepath string
var toFile bool
var toStdout bool
var toInput string
var toHex bool
var toBin bool
var nb []byte
var bsize uint64
var gzip bool
var gzipCheck bool
var del bool
var nullByte = byte(0)
var payload []byte

func check(e error){
	if e!=nil{
		panic(e)
	}
}

func checkEOF(e error)bool{
	if e != nil{
		if e == io.EOF {
	        return true
	    }else{
			panic(e)
		}
	}
	return false
}

/**
 * @brief      Parse the scritp call
 *
 * @return     void
 */
func parseFlags(){

	// Get the filepath flag
	flag.StringVar(&filepath, "file", "./deflibe.txt", "Target file path.")
	// Get the toStdout flag. If true, the output will be sent to stdout
	flag.BoolVar(&toStdout, "o", true, "If true, the program output will be sent to stdout.")
	//Get the input flag
	flag.StringVar(&toInput, "in", "", "Set the input file path.")
	// Get the stdout to file flag
	flag.BoolVar(&toFile, "f", false, "If true write the stdout to a file.")
	// Get the stdout in hex form
	flag.BoolVar(&toHex, "ox", false, "If true write the stdout in hex form.")
	// Get the stdout in bin form
	flag.BoolVar(&toBin, "ob", false, "If true write the stdout in bin form.")
	// Get how many byte to write
	flag.Uint64Var(&bsize, "s", 100, "How many bytes to write. For example, set this to 1000 to write a KByte.")
	// Get the gzip flag
	flag.BoolVar(&gzip, "gc", false, "If true the output file will be compressed with gzip.")
	// Get the delete flag
	flag.BoolVar(&del, "d", false, "If true, delete the file at the end of the script execution.")

	// Flag parsing
	flag.Parse()
}


func doGzip(gfp *os.File)(bool, *os.File){
	fmt.Printf("%x", gfp)
	return true,gfp
}

/**
 * @brief      Handle the doDelete flag
 *
 * @param      path  the file path
 *
 * @return     true if file has been correctly deleted
 */
func doDelete(path string)bool{
	var err = os.Remove(path)
    check(err)
    return true
}

/**
 * @brief      Build the byte payload. By default it will be filled with null bytes
 *
 * @param      bsize  The size in bytes of the payload
 *
 * @return     The payload.
 * 
 * WARNING!! It can cause out of memory error
 */
func buildPayload(bsize uint64)[]byte{
	var payload =  make([]byte, bsize)
	for i := range payload{
		payload[i] = nullByte
	}
	fmt.Printf("The built payload is %d byte long.\n", bsize)
	return payload
}

/**
 * @brief      Write the payload to a file
 *
 * @return     the file pointer, the number of bytes written
 */
func doToFile() (*os.File, int){

	// Create file
	filepointer, err := os.Create(filepath)
	check(err)

    w := bufio.NewWriter(filepointer)
    byteCounter, err := w.Write(payload)
    check(err)

    w.Flush()

    return filepointer, byteCounter
}

/**
 * @brief      Send the payload to stdout
 *
 * @return
 */
func sendToStdout(){
	if toHex == true{
		fmt.Printf("%x",payload)
		fmt.Println()
	}else if toBin == true{
		fmt.Printf("%b", payload)
		fmt.Println()
	}else{
		fmt.Printf("%s", payload)
		fmt.Println()
	}
}

func handleInput(){

	// Get the input file dereferencing the flag pointer
	file, err := os.Open(toInput)
    check(err)

    // Create a slice to handle the file byte
    // It will be read byte by byte
    data := make([]byte, 2)

    // Point the reader to the file
    reader := bufio.NewReader(file)

    defer file.Close()

    var byteCounter int = 0
    for{
    	n, err := reader.Read(data)
    	byteCounter += n
    	eofCheck := checkEOF(err)
    	if eofCheck == true{
    		fmt.Printf("The input file is %d bytes long.\n", byteCounter)
    		break
    	}
        fmt.Printf("%X  --> %s\n", data[:n], data[:n])
    }

    
}


func main(){

	//Start timer
	t1 := time.Now()

	parseFlags()

	if toInput != ""{
		handleInput()
		return
	}

	// Build the payload
	payload = buildPayload(bsize)

	/**
	 * Write stdout to file
	 */
	if toFile == true {
		filepointer, byteCounter := doToFile()
		fmt.Printf("wrote %d bytes in %s\n", byteCounter, time.Since(t1))

		/**
	     * Delete the file afterwards
	     */
	    if del == true{
	    	doDelete(filepath)
	    }

	    defer filepointer.Close()
	}

	/**
	 * Handle the toStdou flag
	 */
	if toStdout == true{
		sendToStdout()
	}
	

    // if gzip == true{
    // 	gzipCheck,filepointer := doGzip(filepointer)
    // 	fmt.Printf("%s", gzipCheck)
    // 	fmt.Printf("%x", filepointer)
    // }

}