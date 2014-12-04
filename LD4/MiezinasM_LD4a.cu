/*
============================================================================
Mantas Miežinas, IFF-2
LD4a
============================================================================
*/

#include "cuda_runtime.h"
#include "device_launch_parameters.h"
#include <cuda.h>

#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <string>
#include <fstream>
#include <iomanip>

using namespace std;

//Konstantos
const char DataFile[] = "MiezinasM_L4.txt";
const int MAX_THREADS = 5;
const int MAX_ARRAY_SIZE = 5;

//Struktura gamintojo modeliams saugoti
struct model {
	char name[75];
	int quantity;
	double price;
};

//Konteinerine struktura saugoti modeliams
struct manufacturer {
	char name[15];
	int quantity;
	model models[MAX_ARRAY_SIZE];
};

//Funkciju prototipai
void ReadFile(string filename, manufacturer (&AllModels)[MAX_THREADS]);
void PrintTable(manufacturer(&printOut)[MAX_THREADS]);
void PrintResults(model (&printOut)[MAX_ARRAY_SIZE]);

__global__ void RunOnGPU(manufacturer *printOut, model *resultsArray);
__device__ char * my_strcpy(char *dest, const char *src);
__device__ char * my_strcat(char *dest, const char *src);

int main() {
	manufacturer AllModels[MAX_THREADS];
	model results[MAX_ARRAY_SIZE];

	ReadFile(DataFile, AllModels);
	PrintTable(AllModels);
	//Nusinuliname rezultatu masyvo elementus
	for (int i = 0; i < MAX_ARRAY_SIZE; i++){
		model data;
		strcpy(data.name, "");
		data.price = 0.0;
		data.quantity = 0;
		results[i] = data;
	}

	//Paruosiame GPU, t.y. perkeliame duomenis is RAM i VRAM
	manufacturer *manu;
	model *resultsArray;
	
	int i = MAX_THREADS * sizeof(manufacturer);
	int j = MAX_ARRAY_SIZE * sizeof(model);

	cudaMalloc((void**)&resultsArray , j);
	cudaMalloc((void**)&manu, i);

	cudaMemcpy(resultsArray, results, j, cudaMemcpyHostToDevice);
	cudaMemcpy(manu, AllModels, i, cudaMemcpyHostToDevice);

	//Iskvieciame GPU metoda, sumuojanti masyvu elementus
	RunOnGPU << <1, MAX_THREADS >> >(manu, resultsArray);
	cudaDeviceSynchronize();

	//Perkeliame rezultatus is VRAM (GPU) i RAM (CPU) atminti
	cudaMemcpy(results, resultsArray, j, cudaMemcpyDeviceToHost);

	//Atlaisviname VRAM atminti
	cudaFree(manu);
	cudaFree(resultsArray);

	//Atspausdiname rezultatus
	PrintResults(results);
	system("Pause");
	return 0;
}

/*
============================================================================
ReadFile

Pradiniu duomenu nuskaitymo funkcija, per nuoroda grazina manufacturer
strukturos masyva
============================================================================
*/
void ReadFile(string filename, manufacturer (&AllModels)[MAX_THREADS]) {
	string title;
	int count, j;
	j = 0;
	ifstream fin(filename);
	if (!fin) {
		cerr << "Couldn't open file!\n";
	}
	else {
		while (!fin.eof()){
			fin >> title >> count;
			strcpy(AllModels[j].name, title.c_str());
			//AllModels[j].name = title;
			AllModels[j].quantity = count;
			model models[MAX_ARRAY_SIZE];
			for (int i = 0; i < count; i++){
				model modelis;
				fin >> modelis.name >> modelis.quantity >> modelis.price;
				AllModels[j].models[i] = modelis;
			}
			j++;
		}
		fin.close();
	}
}

/*
============================================================================
PrintTable

Atspausdina pradinius duomenis lentelemis
============================================================================
*/
void PrintTable(manufacturer(&printOut)[MAX_THREADS]){
	cout << "-----------------------------------------------------------------------------\n";
	for (manufacturer & manu : printOut){
		cout << right << setw(35) << manu.name << "\n";
		cout << "-----------------------------------------------------------------------------\n";
		cout << left << setw(63) << "Modelio Pavadinimas"
			<< setw(8) << "Kiekis"
			<< setw(5) << "Kaina" << "\n";
		cout << "-----------------------------------------------------------------------------\n";
		for (int i = 0; i < manu.quantity; i++){
			model forPrinting = manu.models[i];
			cout << left << setw(3) << to_string(i + 1) + ")"
				<< setw(60) << forPrinting.name
				<< setw(8) << forPrinting.quantity
				<< setw(5) << setprecision(4) << forPrinting.price
				<< "\n";
		}
		cout << "-----------------------------------------------------------------------------\n";
	}
}

/*
============================================================================
PrintResults

Atspausdina rezultatu masyva lenteleje
============================================================================
*/
void PrintResults(model(&printOut)[MAX_ARRAY_SIZE]) {
	cout << "-----------------------------------------------------------------------------\n";
	cout << left << setw(63) << "Modelio Pavadinimas"
		<< setw(8) << "Kiekis"
		<< setw(5) << "Kaina" << "\n";
	cout << "-----------------------------------------------------------------------------\n";
	for (int i = 0; i < MAX_ARRAY_SIZE; i++) {
		model forPrinting = printOut[i];
		cout << left << setw(3) << to_string(i + 1) + ")"
			<< setw(60) << forPrinting.name
			<< setw(8) << forPrinting.quantity
			<< setw(5) << setprecision(4) << forPrinting.price
			<< "\n";
	}
}

/*
============================================================================
RunOnGPU

Susumuoja pradiniu duomenu masyvus i viena masyva. Rezultatas lieka GPU
atmintyje, is kur veliau ji paimsime.
============================================================================
*/
__global__ void RunOnGPU(manufacturer *printOut, model *resultsArray) {
	int gija = threadIdx.x;
	for (int i = 0; i < MAX_THREADS; i++) {
		my_strcat(resultsArray[gija].name, printOut[i].models[gija].name);
		resultsArray[gija].price += printOut[i].models[gija].price;
		resultsArray[gija].quantity += printOut[i].models[gija].quantity;
	}
}

/*
============================================================================
Pagalbines funkcijos, my_strcat skirta apjungti dviems simboliu eilutems,
ji naudoja my_strcpy. Ju reikia tam, nes CUDA nepalaiko string.h
bibliotekos, kurioje yra strcat() funkcija.
============================================================================
*/
__device__ char * my_strcpy(char *dest, const char *src){
	int i = 0;
	do {
		dest[i] = src[i];
	} while (src[i++] != 0);
	return dest;
}

__device__ char * my_strcat(char *dest, const char *src){
	int i = 0;
	while (dest[i] != 0) i++;
	my_strcpy(dest + i, src);
	return dest;
}