// lottery.go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	reset  = "\033[0m"
	green  = "\033[92m"
	red    = "\033[91m"
	yellow = "\033[93m"
	blue   = "\033[94m"
	bold   = "\033[1m"
)

func colorize(text, color string) string {
	return color + text + reset
}

type Lottery struct {
	numNumbers int
	maxNumber  int
	history    [][]int
	historyFile string
}

func NewLottery(num, max int) *Lottery {
	l := &Lottery{
		numNumbers: num,
		maxNumber:  max,
		historyFile: filepath.Join(os.Getenv("HOME"), ".lottery_history.txt"),
	}
	l.loadHistory()
	return l
}

func (l *Lottery) loadHistory() {
	data, err := os.ReadFile(l.historyFile)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		var ticket []int
		for _, p := range parts {
			if n, err := strconv.Atoi(p); err == nil {
				ticket = append(ticket, n)
			}
		}
		if len(ticket) > 0 {
			l.history = append(l.history, ticket)
		}
	}
}

func (l *Lottery) saveHistory(ticket []int) {
	l.history = append(l.history, ticket)
	f, err := os.OpenFile(l.historyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	for _, n := range ticket {
		fmt.Fprintf(f, "%d ", n)
	}
	fmt.Fprintln(f)
}

func (l *Lottery) generateTicket() []int {
	rand.Seed(time.Now().UnixNano())
	m := make(map[int]bool)
	for len(m) < l.numNumbers {
		m[rand.Intn(l.maxNumber)+1] = true
	}
	var ticket []int
	for n := range m {
		ticket = append(ticket, n)
	}
	sort.Ints(ticket)
	return ticket
}

func (l *Lottery) checkTicket(ticket, winning []int) (int, []int) {
	var matched []int
	for _, n := range ticket {
		for _, w := range winning {
			if n == w {
				matched = append(matched, n)
				break
			}
		}
	}
	return len(matched), matched
}

func (l *Lottery) showStats() {
	if len(l.history) == 0 {
		fmt.Println(colorize("Нет истории для статистики.", yellow))
		return
	}
	freq := make(map[int]int)
	for _, t := range l.history {
		for _, n := range t {
			freq[n]++
		}
	}
	fmt.Println(colorize("📊 Статистика выпадений:", bold))
	for n := 1; n <= l.maxNumber; n++ {
		if freq[n] > 0 {
			fmt.Printf("  %2d: %d раз\n", n, freq[n])
		}
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Println("Usage: lottery <num> <max> [-c N] [-w nums] [-s] [-o file] [-v]")
		return
	}
	var num, max, count int
	var winningStr, outputFile string
	statsFlag := false
	verbose := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-c":
			if i+1 < len(args) {
				count, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "-w":
			if i+1 < len(args) {
				winningStr = args[i+1]
				i++
			}
		case "-s":
			statsFlag = true
		case "-o":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		case "-v":
			verbose = true
		default:
			if num == 0 {
				num, _ = strconv.Atoi(args[i])
			} else if max == 0 {
				max, _ = strconv.Atoi(args[i])
			}
		}
	}
	if num <= 0 || max <= 0 || num > max {
		fmt.Println(colorize("Неверные параметры. Укажите num и max (num <= max).", red))
		return
	}
	if count == 0 {
		count = 1
	}

	game := NewLottery(num, max)

	var winningNumbers []int
	if winningStr != "" {
		parts := strings.Split(winningStr, ",")
		for _, p := range parts {
			n, _ := strconv.Atoi(strings.TrimSpace(p))
			winningNumbers = append(winningNumbers, n)
		}
		if len(winningNumbers) != num {
			fmt.Println(colorize("Количество выигрышных номеров должно совпадать с num.", red))
			return
		}
	}

	if statsFlag {
		game.showStats()
		return
	}

	tickets := make([][]int, count)
	for i := 0; i < count; i++ {
		ticket := game.generateTicket()
		tickets[i] = ticket
		game.saveHistory(ticket)
	}

	outputLines := []string{}
	if verbose {
		for i, ticket := range tickets {
			line := fmt.Sprintf("Билет %d: ", i+1)
			if winningNumbers != nil {
				matches, matched := game.checkTicket(ticket, winningNumbers)
				colored := []string{}
				for _, n := range ticket {
					found := false
					for _, m := range matched {
						if n == m {
							found = true
							break
						}
					}
					if found {
						colored = append(colored, colorize(strconv.Itoa(n), green))
					} else {
						colored = append(colored, strconv.Itoa(n))
					}
				}
				line += strings.Join(colored, " ") + fmt.Sprintf(" (совпадений: %d)", matches)
			} else {
				line += strings.Join(strings.Fields(fmt.Sprint(ticket)), " ")
			}
			outputLines = append(outputLines, line)
		}
	} else {
		for _, ticket := range tickets {
			outputLines = append(outputLines, strings.Join(strings.Fields(fmt.Sprint(ticket)), " "))
		}
	}

	output := strings.Join(outputLines, "\n")
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(output), 0644)
		if err == nil {
			fmt.Println(colorize("Результат сохранён в "+outputFile, green))
		} else {
			fmt.Println(colorize("Ошибка записи файла.", red))
		}
	} else {
		fmt.Println(output)
	}
}

func sort.Ints(a []int) {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i] > a[j] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}
