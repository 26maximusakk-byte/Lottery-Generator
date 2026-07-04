// lottery.cpp
#include <iostream>
#include <vector>
#include <set>
#include <string>
#include <random>
#include <algorithm>
#include <fstream>
#include <sstream>
#include <filesystem>
#include <map>

using namespace std;
namespace fs = std::filesystem;

const string RESET = "\033[0m";
const string GREEN = "\033[92m";
const string RED = "\033[91m";
const string YELLOW = "\033[93m";
const string BLUE = "\033[94m";
const string BOLD = "\033[1m";

string colorize(const string& text, const string& color) {
    return color + text + RESET;
}

string getHomeDir() {
    const char* home = getenv("HOME");
    if (!home) home = getenv("USERPROFILE");
    return string(home);
}

vector<int> generateTicket(int num, int maxVal) {
    random_device rd;
    mt19937 gen(rd());
    set<int> nums;
    while (nums.size() < num) {
        nums.insert(gen() % maxVal + 1);
    }
    return vector<int>(nums.begin(), nums.end());
}

pair<int, vector<int>> checkTicket(const vector<int>& ticket, const vector<int>& winning) {
    vector<int> matched;
    for (int n : ticket) {
        if (find(winning.begin(), winning.end(), n) != winning.end()) {
            matched.push_back(n);
        }
    }
    return {matched.size(), matched};
}

void showStats(const vector<vector<int>>& history) {
    map<int, int> freq;
    for (auto& t : history) {
        for (int n : t) freq[n]++;
    }
    cout << colorize("📊 Статистика выпадений:", BOLD) << endl;
    for (auto& p : freq) {
        cout << "  " << p.first << ": " << p.second << " раз" << endl;
    }
}

int main(int argc, char* argv[]) {
    int numNumbers = 0, maxNumber = 0;
    int count = 1;
    string winningStr, outputFile;
    bool showStatsFlag = false, verbose = false;

    for (int i=1; i<argc; ++i) {
        string arg = argv[i];
        if (arg == "-h" || arg == "--help") {
            cout << "Usage: lottery <num> <max> [-c N] [-w nums] [-s] [-o file] [-v]" << endl;
            return 0;
        } else if (arg == "-c" && i+1 < argc) {
            count = stoi(argv[++i]);
        } else if (arg == "-w" && i+1 < argc) {
            winningStr = argv[++i];
        } else if (arg == "-s") {
            showStatsFlag = true;
        } else if (arg == "-o" && i+1 < argc) {
            outputFile = argv[++i];
        } else if (arg == "-v") {
            verbose = true;
        } else if (numNumbers == 0) {
            numNumbers = stoi(arg);
        } else if (maxNumber == 0) {
            maxNumber = stoi(arg);
        }
    }
    if (numNumbers <= 0 || maxNumber <= 0 || numNumbers > maxNumber) {
        cerr << colorize("Неверные параметры. Укажите num и max (num <= max).", RED) << endl;
        return 1;
    }

    vector<int> winningNumbers;
    if (!winningStr.empty()) {
        stringstream ss(winningStr);
        string token;
        while (getline(ss, token, ',')) {
            winningNumbers.push_back(stoi(token));
        }
        if ((int)winningNumbers.size() != numNumbers) {
            cerr << colorize("Количество выигрышных номеров должно совпадать с num.", RED) << endl;
            return 1;
        }
    }

    // Загружаем историю
    string histFile = getHomeDir() + "/.lottery_history.txt";
    vector<vector<int>> history;
    if (fs::exists(histFile)) {
        ifstream f(histFile);
        string line;
        while (getline(f, line)) {
            vector<int> ticket;
            stringstream ss(line);
            int n;
            while (ss >> n) ticket.push_back(n);
            if (!ticket.empty()) history.push_back(ticket);
        }
    }

    if (showStatsFlag) {
        showStats(history);
        return 0;
    }

    vector<vector<int>> tickets;
    for (int i=0; i<count; ++i) {
        auto ticket = generateTicket(numNumbers, maxNumber);
        tickets.push_back(ticket);
        history.push_back(ticket);
    }

    // Сохраняем историю
    ofstream f(histFile);
    for (auto& t : history) {
        for (int n : t) f << n << " ";
        f << "\n";
    }

    // Вывод
    vector<string> lines;
    if (verbose) {
        for (size_t i=0; i<tickets.size(); ++i) {
            string line = "Билет " + to_string(i+1) + ": ";
            if (!winningNumbers.empty()) {
                auto res = checkTicket(tickets[i], winningNumbers);
                string colored;
                for (int n : tickets[i]) {
                    if (find(res.second.begin(), res.second.end(), n) != res.second.end()) {
                        colored += colorize(to_string(n), GREEN) + " ";
                    } else {
                        colored += to_string(n) + " ";
                    }
                }
                line += colored + " (совпадений: " + to_string(res.first) + ")";
            } else {
                for (int n : tickets[i]) line += to_string(n) + " ";
            }
            lines.push_back(line);
        }
    } else {
        for (auto& t : tickets) {
            string line;
            for (int n : t) line += to_string(n) + " ";
            lines.push_back(line);
        }
    }

    string output;
    for (auto& l : lines) output += l + "\n";

    if (!outputFile.empty()) {
        ofstream fout(outputFile);
        if (fout) {
            fout << output;
            cout << colorize("Результат сохранён в " + outputFile, GREEN) << endl;
        } else {
            cerr << colorize("Ошибка записи файла.", RED) << endl;
        }
    } else {
        cout << output;
    }
    return 0;
}
