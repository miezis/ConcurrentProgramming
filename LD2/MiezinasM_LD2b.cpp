#include <iostream>
#include <thread>
#include <string>
#include <fstream>
#include <vector>
#include <iomanip>
#include <omp.h>
#include <stdlib.h>
#include <sstream>
#include <algorithm>

using namespace std;

//kintamieji parodantys darbu pabaiga
volatile bool doneMaking = false;
volatile bool doneUsing = false;

//konstanta duomenu failo pavadinimui
const char DataFile[] = "MiezinasM_L2.txt";
const int MAX_THREADS = 5;
const int MAX_ARRAY_SIZE = 10;

//bazine struktura saugoti vienam modelio irasui
struct model {
	string name;
	int quantity;
	double price;
};

//paprastas skaitliukas, aprasyti operatoriai norint lengvai su juo dirbti
struct Counter {
	double price;
	int count;
public:
	//Counter(double price, int count):price(price), count(count){}
	int operator++(){return ++count;}
	int operator--(){return --count;}
	bool operator==(const Counter &other){return price == other.price;}
	bool operator<(const Counter &other){return price < other.price;}
};

//konteinerine klase, skirta saugoti modeliams (vektoriu)
class Manufacturer {
private:
	string name;
	int quantity;
	vector<model> models;
public:
	Manufacturer(string name, int quantity, vector<model> models) {
		this->name = name;
		this->quantity = quantity;
		this->models = models;
	}
	string getName() {
		return this->name;
	}
	int getQuantity() {
		return this->quantity;
	}
	vector<model> getModels(){
		return this->models;
	}
};


//buferis apsaugotas nuo maigymo is keliu giju
class Buffer {
	vector<Counter> buffer;//duomenys
	volatile bool empty;
public:
	Buffer() : empty(true){}
	bool Add(Counter c) {
		if(doneUsing)
			return false;
		#pragma omp critical(access)
		{
			//randamas atitinkamo pavadinimo skaitliukas
			auto i = find(buffer.begin(), buffer.end(), c);
			if(i != buffer.end()) {
				(*i).count += c.count;
			} else { //jei skaitliuko neradome, kuriame nauja reikiamoje vietoje
				auto size = buffer.size();
				for(auto i = buffer.begin(); i != buffer.end(); i++) {
					if(c < (*i)) {
						buffer.insert(i, c);
						break;
					}
				}
				if(buffer.size() == size)
					buffer.push_back(c);
			}
		}
		return true;
	}

	int Take(Counter c) {
		//unique_lock<mutex> lock(mtx);
		//laukiama jei buferis yra tuscias ir nebaigta rasyti
		//sis tikrinimas ar nebaigta rasyti nera butinas, bet turi galimybe sukelti deadlock
		while(empty);
		int taken = 0;
		#pragma omp critical(access)
		{
			//randamas atitinkamas skaitliukas
			auto i = find(buffer.begin(), buffer.end(), c);
			if(i != buffer.end()) { //paimame kiek imanoma
				if((*i).count >= c.count)
					taken = c.count;
				else
					taken = (*i).count;
				(*i).count -= taken;
				//triname tuscia skaitliuka
				if((*i).count <= 0)
				buffer.erase(i);
			}
		}	
		return taken;
	}
	int Size() {
		int size;
		#pragma omp critical(access)
		{
			size = buffer.size();
		}
		return size;
	}
	string Print() {
		stringstream ss;
		for(auto &c : buffer)
			ss << c.price << " " << c.count << endl;
		return ss.str();
		}
	void Done() {
		empty = false;
	}

};

//Funkciju prototipai
vector<Manufacturer> ReadFile(string filename, vector<vector<Counter>> &users);
void PrintTable(vector<Manufacturer> printOut, vector<vector<Counter>> users);
void PrintManufacturerModels(Manufacturer printOut, string procNum);
void Make(Manufacturer stuff, Buffer & buffer);
void Use(vector<Counter> stuff, Buffer & buffer);




int main() {
	Buffer buffer;
	vector<vector<Counter>> Users;
	vector<Manufacturer> AllModels;

	AllModels = ReadFile(DataFile, Users);
	PrintTable(AllModels, Users);

	int gijosNr;
	int done = 0;
	//omp_set_nested(true);
	//omp_set_num_threads(AllModels.size() + Users.size());

	#pragma omp parallel shared(buffer) private(gijosNr)
	{
		gijosNr = omp_get_thread_num();
		cout << gijosNr << endl;
		#pragma omp for
		for (int i = 0; i < AllModels.size(); i++) {
			Make(AllModels.at(i), buffer);
		}
		#pragma omp barrier
		doneMaking = true;
		buffer.Done();
		#pragma omp for
		for(int i = 0; i < Users.size(); i++)
			Use(Users.at(i), buffer);
		#pragma omp barrier
		doneUsing = true;
	} 
	cout << "Nesuvartota liko:\n\n" << buffer.Print();
	return 0;
}

vector<Manufacturer> ReadFile(string filename, vector<vector<Counter>> &users){
	bool readUsers = false;
	vector<Manufacturer> AllModels;
	vector<Counter> temp;
	string title;
	int count;
	ifstream fin(filename);

	if(!fin) {
		cerr << "Couldn't open file!\n";
	} else {
		while(!fin.eof()){
			string line;
			int start, end;
			start = 0;
			decide:
			getline(fin, line);
			if(!readUsers) {
				if(line == "UsersWant"){
					readUsers = true;
					goto decide;
				}
				end = line.find(' ');
				title = line.substr(start, end);
				start = end + 1;
				end = line.find('\n', start);
				count = atoi(line.substr(start, end - start).c_str());
				//count = stoi(count);
				//line >> title >> count;
				vector<model> models;
				for(int i = 0; i < count; i++){
					model modelis;
					getline(fin, line);
					start = 0;
					end = line.find('\t');
					modelis.name = line.substr(start, end);
					start = end + 1;
					end = line.find('\t', start);
					modelis.quantity = atoi(line.substr(start, end - start).c_str());
					start = end + 1;
					end = line.find('\n', start);
					modelis.price = atof(line.substr(start, end - start).c_str());
					//fin >> modelis.name >> modelis.quantity >> modelis.price;
					models.push_back(modelis);
				}
				AllModels.push_back(Manufacturer(title, count, models));
			} else {
				if(line == " "){
					users.push_back(temp);
					temp.clear();
				} else {
					Counter tmp;
					start = 0;
					end = line.find('\t');
					tmp.price = atof(line.substr(start, end).c_str());
					start = end + 1;
					end = line.find('\n', start);
					tmp.count = atoi(line.substr(start, end - start).c_str());
					temp.push_back(tmp);
				}
			}
		}
		fin.close();
	}	
	return AllModels;
}

void PrintTable(vector<Manufacturer> printOut, vector<vector<Counter>> users){
	cout << "-----------------------------------------------------------------------------\n";
	for(Manufacturer & manu : printOut){
		cout << right << setw(35) << manu.getName() << "\n";
		cout << "-----------------------------------------------------------------------------\n";
		cout << left << setw(63) << "Modelio Pavadinimas"
			 << setw(8) << "Kiekis"
			 << setw(5) << "Kaina" << "\n";
		for(int i = 0; i < manu.getQuantity(); i++){
			model forPrinting = manu.getModels().at(i);
			cout << left << setw(3) << to_string(i+1) + ")"
				 << setw(60) << forPrinting.name
				 << setw(8) << forPrinting.quantity
				 << setw(5) << setprecision(4) << forPrinting.price
				 << "\n";
		}
		cout << "-----------------------------------------------------------------------------\n";
	}
	int i = 1;
	for(auto & user : users) {
		cout << "User Nr. " << i << endl;
		cout << "------------\n";
		cout << "Price  Count\n";
		for(Counter & one : user){
			cout << one.price << "\t" << one.count << endl;
		}
		cout << "------------\n";
		i++;
	}
}

void PrintManufacturerModels(Manufacturer printOut, string procNum){
	int i = 1;
	for(model &mod : printOut.getModels()){
		cout << left << setw(20) << procNum
			 << setw(5) << i
			 << setw(30) << mod.name
			 << setw(5) << mod.quantity
			 << setw(5) << setprecision(4) << mod.price << "\n";
		i++;
	}
}

void Make(Manufacturer stuff, Buffer &buffer) {
	for(auto &s : stuff.getModels()) {
		Counter tmp;
		tmp.price = s.price;
		tmp.count = s.quantity;
		buffer.Add(tmp);
	}
}
//vartojimo funkcija
void Use(vector<Counter> stuff, Buffer &buffer) {
	auto i = stuff.begin();

	while((!doneMaking || buffer.Size() > 0) && stuff.size() > 0) {
		i++;
		if(i == stuff.end())
			i = stuff.begin();
		int taken = buffer.Take(*i);
		(*i).count -= taken;

		if((*i).count <= 0 || (taken == 0 && doneMaking)) {
			stuff.erase(i);
			i = stuff.begin();
		}
	}
}