#include <iostream>
#include <thread>
#include <string>
#include <fstream>
#include <vector>
#include <iomanip>

using namespace std;

const char DataFile[] = "MiezinasM.txt";

struct model {
	string name;
	int quantity;
	double price;
};

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

vector<Manufacturer> ReadFile(string filename){
	vector<Manufacturer> AllModels;
	string title;
	int count;
	ifstream fin(filename);
	if(!fin) {
		cerr << "Couldn't open file!\n";
	} else {
		while(!fin.eof()){
			fin >> title >> count;
			vector<model> models;
			for(int i = 0; i < count; i++){
				model modelis;
				fin >> modelis.name >> modelis.quantity >> modelis.price;
				models.push_back(modelis);
			}
			AllModels.push_back(Manufacturer(title, count, models));
		}
		fin.close();
	}	
	return AllModels;
}

void PrintTable(vector<Manufacturer> printOut){
	cout << "-----------------------------------------------------------------------\n";
	for(Manufacturer & manu : printOut){
		cout << right << setw(35) << manu.getName() << "\n";
		cout << "-----------------------------------------------------------------------\n";
		for(int i = 0; i < manu.getQuantity(); i++){
			model forPrinting = manu.getModels().at(i);
			cout 	<< left << setw(60) << forPrinting.name 
					<< setw(8) << forPrinting.quantity 
					<< setw(5) << setprecision(2) << forPrinting.price 
					<< "\n";
		}
		cout << "-----------------------------------------------------------------------\n";
	}
}

void PrintManufacturerModels(Manufacturer printOut, string procNum){
	int i = 0;
	for(model &mod : printOut.getModels()){
		cout 	<< left << setw(20) << procNum 
				<< setw(5) << i 
				<< setw(30) << mod.name 
				<< setw(5) << mod.quantity 
				<< setw(5) << mod.price << "\n";
		i++;
	}
}

int main() {
	vector<Manufacturer> AllModels = ReadFile(DataFile);
	vector<thread> threads;
	int i = 0;

	PrintTable(AllModels);

	for(Manufacturer &manu : AllModels) {
		threads.push_back(thread(PrintManufacturerModels, manu, "procesas_" + to_string(i)));
		i++;
	}

	for (thread &t: threads) {
        t.join();
    }

	return 0;
}