package main

import "fmt"
import "math"
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
	buff    counterSlice
	channel chan counter
}

func (b buffer) Add(c *counter) {
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

func (b buffer) Take(c *counter) uint {
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
				b.buff[i].count -= c.count
			} else {
				j = b.buff[i].count
				b.buff[i].count -= j

				copy(b.buff[:i-1], b.buff[i:])
				b.buff[len(b.buff)-1] = counter{}
				b.buff = b.buff[:len(b.buff)-1]
				//b.buff[] = append(b.buff[:i-1], b.buff[i:])
			}
			return j
		}
	}
}

func (b buffer) Print() {
	fmt.Println("Nesuvartota liko: \n")
	for v := range b.buff {
		fmt.Println("%+v", v)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
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

func main() {

	var AllModels []manufacturer
	var Users []counterSlice

	AllModels, Users = ReadFile()

	PrintTable(AllModels, Users)

	// Go will infer the type of initialized variables.
	var d = true
	fmt.Println(d)

	// The `:=` syntax is shorthand for declaring and
	// initializing a variable, e.g. for
	// `var f string = "short"` in this case.
	f := "short"
	fmt.Println(f)

	fmt.Println(DataFile)
	const n = 500000000
	fmt.Println(math.Sin(n))

	fmt.Println(int64(n))

	// The most basic type, with a single condition.
	i := 1
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}

	// A classic initial/condition/after `for` loop.
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}

	// `for` without a condition will loop repeatedly
	// until you `break` out of the loop or `return` from
	// the enclosing function.
	for {
		fmt.Println("loop")
		break
	}

	if 7%2 == 0 {
		fmt.Println("7 is even")
	} else {
		fmt.Println("7 is odd")
	}

	var az [5]int
	fmt.Println("emp:", az)

	fmt.Println("len:", len(az))
}
