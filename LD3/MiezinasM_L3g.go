//Miezinas Mantas IFF-2, sedejau prie 23 komp.

package main

import "fmt"

func Valdytojas(done chan bool, kanalasIValdytoja chan int, kanalasISpausdintoja chan [100]int, isValdytojasDone chan bool, doneToSpausdintojas chan bool) {
	var masyvas [100]int
	//nesudeta logika, kad vienu atveju detu i prieki, kitu i gala
	for i := 0; i < 100; {
		var didInc bool
		select {
		case masyvas[i] = <-kanalasIValdytoja:
			{
				done <- false
				i++
				didInc = true
			}
		default:
			didInc = false
		}
		if (i%10 == 0 && didInc) || i == 100 {
			kanalasISpausdintoja <- masyvas
		}
	}

	done <- true
	doneToSpausdintojas <- true
	isValdytojasDone <- true
}

func Numeriai(done chan bool, iValdytoja chan int, procNr int) {
	salyga := <-done
	for salyga != true {
		iValdytoja <- procNr
		procNr += 2
		salyga = <-done
	}
}

func Spausdintojas(kanalas chan [100]int, done chan bool) {
	for i := 0; i < 10; {
		select {
		case masyvas := <-kanalas:
			{
				fmt.Println(masyvas)
				i++
			}
		default:
		}
	}
	<-done
}

func main() {

	var kanalasIValdytoja chan int
	var kanalasISpausdintoja chan [100]int
	var done chan bool
	var isValdytojasDone chan bool
	var doneToSpausdintojas chan bool

	kanalasIValdytoja = make(chan int, 10)        //Many2One
	kanalasISpausdintoja = make(chan [100]int, 1) //One2One
	done = make(chan bool, 1)
	isValdytojasDone = make(chan bool, 1)
	doneToSpausdintojas = make(chan bool, 1)

	done <- false

	for i := 1; i < 3; i++ {
		go Numeriai(done, kanalasIValdytoja, i)
	}

	go Spausdintojas(kanalasISpausdintoja, doneToSpausdintojas)
	go Valdytojas(done, kanalasIValdytoja, kanalasISpausdintoja, isValdytojasDone, doneToSpausdintojas)
	<-isValdytojasDone
}
