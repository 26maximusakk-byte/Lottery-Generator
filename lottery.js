// lottery.js
#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const os = require('os');
const readline = require('readline');

const COLORS = {
    reset: '\x1b[0m',
    green: '\x1b[92m',
    red: '\x1b[91m',
    yellow: '\x1b[93m',
    blue: '\x1b[94m',
    bold: '\x1b[1m'
};

function colorize(text, color) {
    return COLORS[color] + text + COLORS.reset;
}

class Lottery {
    constructor(num, max) {
        this.num = num;
        this.max = max;
        this.historyFile = path.join(os.homedir(), '.lottery_history.txt');
        this.history = this.loadHistory();
    }

    loadHistory() {
        try {
            const data = fs.readFileSync(this.historyFile, 'utf8');
            const lines = data.split('\n').filter(line => line.trim() !== '');
            return lines.map(line => line.trim().split(/\s+/).map(Number));
        } catch {
            return [];
        }
    }

    saveHistory(ticket) {
        this.history.push(ticket);
        const line = ticket.join(' ') + '\n';
        fs.appendFileSync(this.historyFile, line, 'utf8');
    }

    generateTicket() {
        const nums = new Set();
        while (nums.size < this.num) {
            nums.add(Math.floor(Math.random() * this.max) + 1);
        }
        return Array.from(nums).sort((a, b) => a - b);
    }

    checkTicket(ticket, winning) {
        const matched = ticket.filter(n => winning.includes(n));
        return { matches: matched.length, matched };
    }

    showStats() {
        if (this.history.length === 0) {
            console.log(colorize('Нет истории для статистики.', 'yellow'));
            return;
        }
        const freq = {};
        for (const ticket of this.history) {
            for (const n of ticket) {
                freq[n] = (freq[n] || 0) + 1;
            }
        }
        console.log(colorize('📊 Статистика выпадений:', 'bold'));
        for (let i = 1; i <= this.max; i++) {
            if (freq[i]) {
                console.log(`  ${i.toString().padStart(2)}: ${freq[i]} раз`);
            }
        }
    }
}

function main() {
    const args = process.argv.slice(2);
    if (args.length === 0 || args[0] === '-h' || args[0] === '--help') {
        console.log('Usage: node lottery.js <num> <max> [-c N] [-w nums] [-s] [-o file] [-v]');
        process.exit(0);
    }
    let num = 0, max = 0, count = 1;
    let winningStr = null, outputFile = null;
    let statsFlag = false, verbose = false;

    for (let i = 0; i < args.length; i++) {
        const arg = args[i];
        if (arg === '-c' && i+1 < args.length) {
            count = parseInt(args[++i]);
        } else if (arg === '-w' && i+1 < args.length) {
            winningStr = args[++i];
        } else if (arg === '-s') {
            statsFlag = true;
        } else if (arg === '-o' && i+1 < args.length) {
            outputFile = args[++i];
        } else if (arg === '-v') {
            verbose = true;
        } else if (num === 0) {
            num = parseInt(arg);
        } else if (max === 0) {
            max = parseInt(arg);
        }
    }
    if (num <= 0 || max <= 0 || num > max) {
        console.log(colorize('Неверные параметры. Укажите num и max (num <= max).', 'red'));
        process.exit(1);
    }

    const game = new Lottery(num, max);

    let winningNumbers = null;
    if (winningStr) {
        winningNumbers = winningStr.split(',').map(s => parseInt(s.trim()));
        if (winningNumbers.length !== num) {
            console.log(colorize('Количество выигрышных номеров должно совпадать с num.', 'red'));
            process.exit(1);
        }
    }

    if (statsFlag) {
        game.showStats();
        return;
    }

    const tickets = [];
    for (let i = 0; i < count; i++) {
        const ticket = game.generateTicket();
        tickets.push(ticket);
        game.saveHistory(ticket);
    }

    const outputLines = [];
    if (verbose) {
        tickets.forEach((ticket, idx) => {
            let line = `Билет ${idx+1}: `;
            if (winningNumbers) {
                const { matches, matched } = game.checkTicket(ticket, winningNumbers);
                const colored = ticket.map(n => matched.includes(n) ? colorize(String(n), 'green') : String(n));
                line += colored.join(' ') + ` (совпадений: ${matches})`;
            } else {
                line += ticket.join(' ');
            }
            outputLines.push(line);
        });
    } else {
        tickets.forEach(ticket => outputLines.push(ticket.join(' ')));
    }

    const output = outputLines.join('\n');
    if (outputFile) {
        fs.writeFileSync(outputFile, output, 'utf8');
        console.log(colorize(`Результат сохранён в ${outputFile}`, 'green'));
    } else {
        console.log(output);
    }
}

main();
