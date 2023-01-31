package main

import (
	"crypto/sha256"
	"fmt"
	spaolacciMurmur "github.com/spaolacci/murmur3"
	"hash/fnv"
	"math"
	"math/big"
	"strconv"
	"time"
)

func main() {
	count := 100000

	//Hash(count, 1)
	//Hash(count, 2)
	Hash(count, 3)
}

func Hash(count int, x int) {

	mymap := make(map[string]string)

	if x == 1 {
		fmt.Println("MURMUR")
	} else if x == 2 {
		fmt.Println("FNV")
	} else {
		fmt.Println("sha256")
	}

	start := time.Now()
	invalidCount := 0
	for i := 0; i < count; i++ {
		v := "pk" + strconv.Itoa(i) + "psid" + strconv.Itoa(153)

		vByte := []byte(v)

		var hashedProductID uint32
		if x == 1 {
			hashedProductID = spaolacciMurmur.Sum32(vByte)
		} else if x == 2 {
			algorithm := fnv.New32a()

			algorithm.Write(vByte)
			hashedProductID = algorithm.Sum32()
		} else {
			var minProductID int64 = 50000000                 //highest value already persisted in the db
			var maxProductID int64 = math.MaxInt32            //maximum allowed value for int32
			idList := big.NewInt(maxProductID - minProductID) //number of allowable ids available
			//Range of potential ID's 50M - 2.1B

			//hash productKey and primaryStoreId []byte from above
			hash := sha256.New()
			hash.Write(vByte)
			sum := hash.Sum(nil)

			//new int from the Hashed value and ensure it's a positive value
			intFromHash := new(big.Int).SetBytes(sum[0:8])
			intFromHash.Abs(intFromHash)

			//get remainder of intFromHash/idList
			//mod will be guaranteed to be less than idList - which is the range of valid id's
			mod := intFromHash.Mod(intFromHash, idList)

			hashedProductID = uint32(mod.Int64() + 50000000)
		}

		if hashedProductID > 0 {
			productId := int(hashedProductID)
			mymap[v] = strconv.Itoa(productId)
			if productId < 50000000 {
				//fmt.Printf("%d: 			%d\n", i, productId)
				invalidCount++
			} else if productId > 2147483647 {
				//fmt.Printf("%d: 			%d\n", i, productId)
				invalidCount++
			} else {
				//fmt.Printf("%d: %d\n", i, productId)

				if productId == 731575512 || productId == 562221661 || productId == 1351126706 || productId == 379937701 {
					fmt.Println(v)
					fmt.Println(productId)
				}

			}

			//fmt.Println()
		}
	}

	fmt.Printf("		Invalid: %d/%d\n", invalidCount, count)
	duration := time.Since(start)
	fmt.Printf("		Duration: %d ns\n", duration.Nanoseconds())
	fmt.Printf("		Duration: %d ms\n", duration.Milliseconds())

	fmt.Println("DUP CHECK:")
	fmt.Println(hasDupes(mymap))
}

func hasDupes(m map[string]string) bool {
	x := make(map[string]struct{})

	for _, v := range m {
		if _, has := x[v]; has {
			//fmt.Println(v)
			//fmt.Println(x[v])
			//return true
		}
		x[v] = struct{}{}
	}

	return false
}
