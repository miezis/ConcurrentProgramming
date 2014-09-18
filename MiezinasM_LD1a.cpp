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
				cin >> modelis.name >> modelis.quantity >> modelis.price;
				models.push_back(modelis);
			}
			AllModels.push_back(Manufacturer(title, count, models));
		}
		fin.close();
	}	
	return AllModels;
}

void PrintTable(vector<Manufacturer> printOut){
	for(Manufacturer & manu : printOut){
		cout << setw(70) << manu.getName() << "\n";
		for(int i = 0; i < manu.getQuantity(); i++){
			model forPrinting = manu.getModels().at(i);
			cout << setw(60) << forPrinting.name << setw(8) << forPrinting.quantity << setw(5) << setprecision(2) << forPrinting.price << "\n";
		}
	}
}

int main() {
	vector<Manufacturer> AllModels = ReadFile(DataFile);
	PrintTable(AllModels);
}