//Mantas Miezinas 10
//kartais buna deadlock
#include <mutex>
#include <condition_variable>
#include <iostream>
#include <thread>
#include <vector>

using namespace std;

volatile bool doneUsing = false;

class semaphore {
private:
	int count;
	mutex mtx;
	condition_variable cv;
public:
	semaphore(int count_):count(count_){;}
	semaphore() { count = 1; }
	void acquire(){
		unique_lock<mutex> lck(mtx);
		while(count == 0){
			cv.wait(lck);
		}
		count--;
	}
	void release(){
		unique_lock<mutex> lck(mtx);
        ++count;
        cv.notify_all();
	}
};

void Make(int add);
void Use(int procNr);

int common = 10;
int howMany = 0;

semaphore S;

int main() {
	vector<thread> Making;
	vector<thread> Using;

	for(int i = 0; i < 3; i++){
		Using.emplace_back(Use, (i + 3));
	}
	
	for(int i = 0; i < 2; i++){
		Making.emplace_back(Make, (i + 2));
	}

	for(auto &v : Making)
		v.join();


	for(auto &v : Using)
		v.join();


	return 0;
}

void Make(int add) {
	int used;
	while(!doneUsing){
		S.acquire();
		if((howMany % 2 == 0 && howMany != used) || (howMany % 3 == 0 && howMany != used)){
			common += add;
			used = howMany;
		}
		S.release();
	}
	
}

void Use(int procNr) {
	int val;
	while(!doneUsing){
		S.acquire();
		if(howMany < 10 && val != common){
			cout << procNr << ") reiksme: " << common << endl;
			val = common;
			howMany++;
		} else if (howMany == 10) {
			doneUsing = true;
		}
		S.release();
	}
	
}