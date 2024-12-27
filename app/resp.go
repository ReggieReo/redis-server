package main

import (
	"fmt"
	"strconv"
)

type Type byte

const (
	Integer = ':'
	String  = '+'
	Bulk    = '$'
	Array   = '*'
	Error   = '-'
)

type RESP struct {
	Type      Type
	Raw       []byte
	Data      []byte
	ArrayData []byte
	Count     int
}

func formatReturnString(returnString string) string {
	response := fmt.Sprintf("$%d\r\n%s\r\n", len(returnString), returnString)
	return response
}

// ReadCommand the number of byte and resp command
func ReadCommand(packet []byte) (n int, resp RESP) {
	if len(packet) == 0 {
		return 0, RESP{}
	}
	// get type from first byte
	resp.Type = Type(packet[0])
	switch resp.Type {
	case Integer, String, Bulk, Array, Error:
	default:
		return 0, RESP{} // invalid type
	}
	//read to first end of line
	i := 1
	for ; ; i++ {
		if len(packet) == i {
			return 0, RESP{} // not enough data
		}
		if packet[i] == '\n' {
			if packet[i-1] != '\r' {
				return 0, RESP{} // missing CR
			}
			i++
			break
		}
	}
	resp.Raw = packet[0:i]      // all data
	resp.Data = packet[1 : i-2] // data except CR
	// optional - and integer from 0 to 9
	if resp.Type == Integer {
		var j int
		if resp.Data[0] == '-' {
			if len(resp.Data) == 1 {
				return 0, RESP{} // invalid integer
			}
			j++
		}
		for ; j < len(resp.Data); j++ {
			if resp.Data[j] < '0' || resp.Data[j] > '9' {
				return 0, RESP{} // invalid integer
			}
		}
		return len(resp.Raw), resp
	}
	if resp.Type == String || resp.Type == Error {
		return len(resp.Raw), resp
	}
	var err error
	resp.Count, err = strconv.Atoi(string(resp.Data)) // convert length to integer
	if err != nil {
		fmt.Println("Error getting count")
		return 0, RESP{}
	}
	if resp.Type == Bulk {
		if resp.Count < 0 {
			resp.Data = nil
			resp.Count = 0
			return len(resp.Raw), resp
		}
		if len(packet) < i+resp.Count+2 {
			// i to first cr and count + 2 -> string + cr
			return 0, RESP{} // not enough data
		}
		if packet[i+resp.Count] != '\r' || packet[i+resp.Count+1] != '\n' {
			return 0, RESP{} // missing CR
		}
		resp.Data = packet[i : i+resp.Count]
		resp.Raw = packet[0 : i+resp.Count+2]
		resp.Count = 0
		return len(resp.Raw), resp
	}
	// last type array
	var tn int
	sdata := packet[i:]               // the remaining data
	for j := 0; j < resp.Count; j++ { // loop through the remaining in array
		rn, rresp := ReadCommand(sdata) // any next data in array
		if rresp.Type == 0 {
			fmt.Println("error")
			return 0, RESP{} // invalid type in next command in array
		}
		tn += rn           // add number of byte of currnet array
		sdata = sdata[rn:] // read next data
	}
	resp.Data = packet[i : i+tn] // Data is all command except type array at the start
	resp.Raw = packet[0 : i+tn]
	// fmt.Println("packet before return")
	// fmt.Println((resp.Data))
	// fmt.Println(resp.Raw)
	// fmt.Println(string(resp.Count))
	// fmt.Println(string(resp.Type))
	return len(resp.Raw), resp
}
