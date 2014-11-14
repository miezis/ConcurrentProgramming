package main

import "fmt"
import "os"
import "bufio"
import "io"
import "strings"
import "strconv"

//konstanta duomenu failo pavadinimui
const DataFile string = "MiezinasM_L2.txt"

/*
============================================================================

	Bazine struktura saugoti vienam modelio irasui

============================================================================
*/
type model struct {
	name     string
	quantity uint
	price    float32
}

/*
============================================================================

	Skaitliukas su metodais, kad butu lengviau dirbti

============================================================================
*/
type counter struct {
	price float32
	count uint
}

type counterSlice []counter

func (slice *counterSlice) find(value float32) int {
	for pos, v := range *slice {
		if v.price == value {
			return pos
		}
	}
	return -1
}

func (slice *counterSlice) getPos(value float32) int {
	curPos := 0
	for pos, v := range *slice {
		if v.price < value {
			curPos = pos
		} else {
			return curPos
		}
	}
	return -1
}

type manufacturer struct {
	name     string
	quantity uint
	models   []model
}

type buffer struct {
	buff counterSlice
}

func (b *buffer) Add(c *counter) {
	i := (b.buff).find(c.price)
	switch i {
	case -1:
		{
			j := (b.buff).getPos(c.price)
			if j == -1 {
				b.buff = append(b.buff, *c)
			} else {
				b.buff[j] = *c
			}
		}
	default:
		{
			b.buff[i].count += c.count
		}
	}
}

func (b *buffer) Take(c *counter) uint {
	if len(b.buff) == 0 {
		return 0
	}
	i := (b.buff).find(c.price)
	switch i {
	case -1:
		{
			return 0
		}
	default:
		{
			var j uint
			if b.buff[i].count >= c.count {
				j = c.count
			} else {
				j = b.buff[i].count

				//b.buff[] = append(b.buff[:i-1], b.buff[i:])
			}
			b.buff[i].count -= j
			if b.buff[i].count <= 0 {
				copy(b.buff[:i-1], b.buff[i:])
				b.buff[len(b.buff)-1] = counter{}
				b.buff = b.buff[:len(b.buff)-1]
			}

			return j
		}
	}
}

func (b *buffer) Print() {
	fmt.Println("Nesuvartota liko: \n")
	for _, v := range b.buff {
		fmt.Printf("%v\n", v)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
============================================================================
ReadFile

	Pradiniu duomenu nuskaitymo funkcija, grazina gamintoju ir vartotoju
	sarasa.

============================================================================
*/

func ReadFile() ([]manufacturer, []counterSlice) {
	var Allmodels []manufacturer
	var users []counterSlice
	var readUsers bool = false
	var temp counterSlice

	file, err := os.Open(DataFile)
	check(err)
	read := bufio.NewReader(file)

	for j := 0; ; j++ {
	decide:
		line, _, err := read.ReadLine()
		if err == io.EOF {
			break
		}
		if !readUsers {
			if "UsersWant" == string(line) {
				readUsers = true
				goto decide
			}
			vars := strings.Split(string(line), " ")
			title := vars[0]
			count, _ := strconv.Atoi(vars[1])
			models := make([]model, count)
			for i := 0; i < count; i++ {
				line, _, err := read.ReadLine()
				check(err)
				vars := strings.Split(string(line), "\t")
				name := vars[0]
				quantity, err := strconv.Atoi(vars[1])
				check(err)
				price, err := strconv.ParseFloat(vars[2], 32)
				check(err)
				models[i] = model{name, uint(quantity), float32(price)}
			}
			Allmodels = append(Allmodels, manufacturer{title, uint(count), models})
		} else {
			if " " == string(line) {
				users = append(users, temp)
				temp = make(counterSlice, 0)
			} else {
				vars := strings.Split(string(line), "\t")
				price, err := strconv.ParseFloat(vars[0], 32)
				check(err)
				count, err := strconv.Atoi(vars[1])
				check(err)
				temp = append(temp, counter{float32(price), uint(count)})
			}
		}
	}
	file.Close()
	return Allmodels, users
}

func PrintTable(printOut []manufacturer, users []counterSlice) {
	fmt.Println("-----------------------------Pradiniai Duomenys------------------------------")
	for _, element := range printOut {
		fmt.Println(element.name)
		fmt.Println("-----------------------------------------------------------------------------")
		fmt.Printf("%20v %45v %10v\n", "Modelio Pavadinimas", "Kiekis", "Kaina")
		for i, modelis := range element.models {
			fmt.Printf("%2s %17s %45d %10.2f\n", strconv.Itoa(i+1)+")", modelis.name, modelis.quantity, modelis.price)
		}
		fmt.Println("-----------------------------------------------------------------------------")
	}
	fmt.Println("------------------------------------Users------------------------------------")
	for i, element := range users {
		fmt.Printf("%s %d\n", "User Nr.", i)
		fmt.Println("------------")
		fmt.Println("Price  Count")
		for _, user := range element {
			fmt.Printf("%5.2f  %5d\n", user.price, user.count)
		}
		fmt.Println("------------")
	}
}

func Use(channel chan counter, user counterSlice, taken chan counter) {
	for _, el := range user {
		channel <- el
		c := <-taken
		for c.count > 0 {
			channel <- c
			c = <-taken
		}
	}
}

func Make(channel chan counter, models manufacturer) {
	for _, element := range models.models {
		channel <- counter{element.price, element.quantity}
	}
}

func Valdytojas(userChannels chan counter, producerChannels chan counter, done chan bool, buff *buffer, taken chan counter) {

	for i := 0; i < 50000; i++ {
		select {
		case counter := <-producerChannels:
			buff.Add(&counter)
		case counter := <-userChannels:
			{
				counter.count -= buff.Take(&counter)
				taken <- counter
			}

		default:
		}
	}

	done <- true

}

func main() {
	done := make(chan bool, 1)

	var buff buffer
	var AllModels []manufacturer
	var Users []counterSlice

	AllModels, Users = ReadFile()

	PrintTable(AllModels, Users)

	var producerMessages chan counter
	var consumerMessages chan counter
	var taken chan counter
	//	var isDone []chan bool

	producerMessages = make(chan counter, len(AllModels))
	consumerMessages = make(chan counter, len(Users))
	taken = make(chan counter, len(Users))
	for _, element := range AllModels {
		go Make(producerMessages, element)
	}

	for _, element := range Users {
		go Use(consumerMessages, element, taken)
	}

	go Valdytojas(consumerMessages, producerMessages, done, &buff, taken)

	<-done
	buff.Print()
}
