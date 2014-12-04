/*
Mantas Miežinas, IFF-2
LD3a

Deja, nemanau, kad tokią užduotį neįmanoma realizuoti pasitelkiant go sinchroninius kanalus,
viskas yra gražu, kol reikia pasirašyti alternatyvą, o tam būdo su sinchroninių kanalų masyvais
neradau, užtat naudojant asinchroninius kanalus viskas gaunasi natūraliai, nes užduoties furmoluotė
natūraliai turėtų būti sprendžiama su buferizuotais kanalais (mano nuomonė). Tad ši programa nėra lygiagreti,
nes iteruojama per siuntėjų kanalus bei vartotojų kanalus ir galiausiai gaunama deadlock. Jei aš klystu ir spręndimas
su sinchroniniais kanalais yra galimas go kalboje, prašau paaiškinkite?

Neesu pirmas susidūręs su tokia problema:

https://groups.google.com/forum/#!topic/golang-nuts/WDdWXO07Lj0
*/

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

/*
============================================================================

	Esamo skaitliuko pozicijos radimas, jei tokio nera, grazinama -1

============================================================================
*/
func (slice *counterSlice) find(value float32) int {
	for pos, v := range *slice {
		if v.price == value {
			return pos
		}
	}
	return -1
}

/*
============================================================================

	Funkcija rasti pozicija, i kuria reikia iterpti skaitliuka,
	jei pozicija pabaigoje slice'o, tuomet graziname -1

============================================================================
*/
func (slice *counterSlice) getPos(value float32) int {
	for pos, v := range *slice {
		if v.price > value {
			return pos
		}
	}
	return -1
}

/*
============================================================================

	konteinerine struktura, skirta saugoti modeliams

============================================================================
*/
type manufacturer struct {
	name     string
	quantity uint
	models   []model
}

/*
============================================================================

	struktura, skirta bendrai atminciai, su pagalbiniais metodais:
	Add, Take ir Print

============================================================================
*/
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
				//b.buff[j] = *c

				b.buff = append(b.buff, counter{})
				copy(b.buff[j+1:], b.buff[j:])
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
				copy(b.buff[i:], b.buff[i+1:])
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

/*
============================================================================

	Pagalbine funkcija, tikrini klaidoms

============================================================================
*/
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

/*
============================================================================
PrintTable

	Atspausdina pradinius duomenis lentelemis
============================================================================
*/
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
	for i := range user {
		channel <- user[i]
	consume:
		select {
		case c := <-taken:
			{
				channel <- c
				goto consume
			}
		default:
			break
		}
		//<-channel
	}
}

func Make(channel chan counter, models manufacturer) {
	for _, element := range models.models {
		channel <- counter{element.price, element.quantity}
		//<-channel
	}
}

func Valdytojas(userChannels []chan counter, producerChannels []chan counter, taken []chan counter, done chan bool, buff *buffer, howMuch int, removes int) {

	for i := 0; i < howMuch+removes; {
		select { //Alternatyvos pasirinkimas
		case counter := <-producerChannels[0]: //Siuo atveju priimame duomenis is gamintoju ir dedame i bendraja atmint
			{
				buff.Add(&counter)
				i++
			}
		case counter := <-producerChannels[1]: //Siuo atveju priimame duomenis is gamintoju ir dedame i bendraja atmint
			{
				buff.Add(&counter)
				i++
			}
		case counter := <-producerChannels[2]: //Siuo atveju priimame duomenis is gamintoju ir dedame i bendraja atmint
			{
				buff.Add(&counter)
				i++
			}
		case counter := <-producerChannels[3]: //Siuo atveju priimame duomenis is gamintoju ir dedame i bendraja atmint
			{
				buff.Add(&counter)
				i++
			}
		case counter := <-producerChannels[4]: //Siuo atveju priimame duomenis is gamintoju ir dedame i bendraja atmint
			{
				buff.Add(&counter)
				i++
			}
		case counter := <-userChannels[0]: //Siuo atveju priimame duomenis is vartotoju ir isimame imanoma kieki
			{ //elemento is bendros atminties, jei lieka neisimtu, siunciame counter objekta atgal (su likuciu)
				counter.count -= buff.Take(&counter)
				if counter.count == 0 {
					i++
				} else {
					taken[0] <- counter
				}
			}
		case counter := <-userChannels[1]: //Siuo atveju priimame duomenis is vartotoju ir isimame imanoma kieki
			{ //elemento is bendros atminties, jei lieka neisimtu, siunciame counter objekta atgal (su likuciu)
				counter.count -= buff.Take(&counter)
				if counter.count == 0 {
					i++
				} else {
					taken[1] <- counter
				}
			}
		case counter := <-userChannels[2]: //Siuo atveju priimame duomenis is vartotoju ir isimame imanoma kieki
			{ //elemento is bendros atminties, jei lieka neisimtu, siunciame counter objekta atgal (su likuciu)
				counter.count -= buff.Take(&counter)
				if counter.count == 0 {
					i++
				} else {
					taken[2] <- counter
				}
			}
		case counter := <-userChannels[3]: //Siuo atveju priimame duomenis is vartotoju ir isimame imanoma kieki
			{ //elemento is bendros atminties, jei lieka neisimtu, siunciame counter objekta atgal (su likuciu)
				counter.count -= buff.Take(&counter)
				if counter.count == 0 {
					i++
				} else {
					taken[3] <- counter
				}
			}

		default:
			break
		}
	}

	done <- true

}

func main() {
	done := make(chan bool, 1)

	var howMuch int = 0
	var removes int = 0

	var buff buffer
	var AllModels []manufacturer
	var Users []counterSlice

	AllModels, Users = ReadFile()

	for _, el := range AllModels {
		howMuch += len(el.models)
	}

	for _, el := range Users {
		removes += len(el)
	}
	fmt.Println(removes)

	PrintTable(AllModels, Users)

	var producerMessages []chan counter
	var consumerMessages []chan counter
	var taken []chan counter

	for i := 0; i < len(AllModels); i++ {
		producerMessages = append(producerMessages, make(chan counter))
	}

	for i := 0; i < len(Users); i++ {
		consumerMessages = append(consumerMessages, make(chan counter))
		taken = append(taken, make(chan counter))
	}

	for i, element := range AllModels {
		go Make(producerMessages[i], element)
	}

	for i, element := range Users {
		go Use(consumerMessages[i], element, taken[i])
	}
	go Valdytojas(consumerMessages, producerMessages, taken, done, &buff, howMuch, removes)
	<-done
	buff.Print()
}
