package main

import (
	"fmt"
	"bytes"
	"strconv"
	"strings"
)

func GetFirstN(n uint) (byte){
	o := byte(0)
	for j := uint(0); j < n; j++{
		o = o | 1<<(7-j)
	}
	return o
}

func main(){

	var filled = uint(0)
	var all_used_characters []byte
	var output = []byte{0}
	var input = "The way this script works is that it first creates a reference array containing all of the used characters in the string at most once. Then, it just creates a stream of 1s and 0s which tell the position of the encrypted character in that reference array. It finally outputs that stream of 1s and 0s as string with some more additional information. Depending on how diverse the characters used are, this script can save from 10% to 70% in terms of storage for the information"
	for  i := 0; i < len(input); i++ {
		if bytes.IndexByte(all_used_characters, input[i]) == -1 {
			all_used_characters = append(all_used_characters, input[i])
		}
	}

	push_up := 30
	size_of_char := byte(len(all_used_characters))
	size_of_block := uint(0)
	for i := size_of_char; i != 0; i = i >> 1{
		size_of_block++
	}


	diff := uint(8-size_of_block)


	for  i := 0; i < len(input); i++ {
		n := byte(bytes.IndexByte(all_used_characters, input[i]))

		n = n<<diff


		free := uint(8-filled)
		o := GetFirstN(free)

		to_append := o & n
		to_append = to_append>>filled

		length := len(output)

		output[length-1] = output[length-1] | to_append

		filled += uint(size_of_block)
		n = n<<free

		if filled > 7 {
			output = append(output, n)
			filled = filled - 8
		}
	}

	toprint := ""
	for j := 0; j < len(output); j++ {
		_ = toprint + strconv.FormatInt(int64(output[j]), 2) + " "
		toprint = toprint + string(byte(push_up)+output[j])
	}

	key := ""
	for j := 0; j < len(all_used_characters); j++ {
		key = key + string(all_used_characters[j])
	}

	result := "["+strconv.Itoa(int(size_of_block))+"]{"+strconv.Itoa(int(push_up))+"}"+key+"..."+toprint
	fmt.Println(result)
	fmt.Println(" ")


	/*
	Do we want to reorder the key after scanning the whole text?
		In that case, it wont necessarily be possible to have the first few "encrypted" couple of characters to be the first few characters in the key.
			Thus, we will need to indicate what is the push_up that we used.
		However, if we make sure that the first few characters in the "encrypted" string be the same as the first few characters of the key, we can figure out what is the push_up programatically
	*/

	start := strings.Index(result, "...") + 3
	size_of_block_temp, _ := strconv.Atoi( string(result[1: strings.Index(result, "]")]) )
	size_of_block = uint(size_of_block_temp)
	to_get := uint(size_of_block)
	pos := uint(0)


	push_up, _ = strconv.Atoi( result[strings.Index(result, "{")+1: strings.Index(result, "}")] )

	alphabet := string(result[strings.Index(result, "}")+1:strings.Index(result, "...") ])

	var decrypted string
	decrypted = ""
	to_decrypt_pre := []rune(result[start:])
	to_decrypt := make([]byte, len(to_decrypt_pre) )

	for j := 0; j < len(to_decrypt) ; j++ {
		to_decrypt[j] = byte(byte(to_decrypt_pre[j])-uint8(push_up))
	}
	// Now, the first size_of_block bits should all be 0.

	for j := 0; j < len(to_decrypt); {

		o := GetFirstN(to_get)
		n := o & (to_decrypt[j]<<pos)
		if to_get < size_of_block { // We have gotten less bits than we need, so we will look at the next byte
			if j+1 < len(to_decrypt) {
				next_bits := to_decrypt[j+1]
				next_selected_bits := GetFirstN(size_of_block-to_get) & next_bits
				appending_bits := next_selected_bits>>to_get
				n = n | appending_bits
			}
		}

		n = n>>(8-size_of_block)
		decrypted = decrypted + string(alphabet[n])



		pos += size_of_block
		if pos > 7 {
			to_get = size_of_block
			pos = pos - 8
			j++
		}

		if 8-pos < size_of_block && pos != 0 {
			to_get = 8-pos
		} else {
			to_get = size_of_block
		}

	}
	fmt.Println(decrypted)
}
