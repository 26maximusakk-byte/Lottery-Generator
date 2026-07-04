# lottery.py
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import random
import argparse
import json
from pathlib import Path
from collections import Counter

# ANSI-цвета
COLORS = {
    'reset': '\033[0m',
    'green': '\033[92m',
    'red': '\033[91m',
    'yellow': '\033[93m',
    'blue': '\033[94m',
    'bold': '\033[1m'
}

def colorize(text, color):
    return f"{COLORS.get(color, '')}{text}{COLORS['reset']}"

class LotteryGenerator:
    def __init__(self, num_numbers, max_number):
        self.num_numbers = num_numbers
        self.max_number = max_number
        self.history_file = Path.home() / '.lottery_history.json'
        self.history = self.load_history()

    def load_history(self):
        if self.history_file.exists():
            with open(self.history_file, 'r') as f:
                return json.load(f)
        return []

    def save_history(self, ticket):
        self.history.append(ticket)
        with open(self.history_file, 'w') as f:
            json.dump(self.history, f)

    def generate_ticket(self):
        numbers = sorted(random.sample(range(1, self.max_number+1), self.num_numbers))
        return numbers

    def check_ticket(self, ticket, winning_numbers):
        matches = set(ticket) & set(winning_numbers)
        return len(matches), sorted(matches)

    def stats(self):
        if not self.history:
            print(colorize("Нет истории для статистики.", 'yellow'))
            return
        all_numbers = [num for ticket in self.history for num in ticket]
        counter = Counter(all_numbers)
        print(colorize("📊 Статистика выпадений:", 'bold'))
        for num, count in sorted(counter.items()):
            print(f"  {num:2}: {count} раз")

def main():
    parser = argparse.ArgumentParser(description="Lottery Generator")
    parser.add_argument('num_numbers', type=int, help='Количество номеров в билете')
    parser.add_argument('max_number', type=int, help='Максимальное значение номера')
    parser.add_argument('-c', '--count', type=int, default=1, help='Количество билетов')
    parser.add_argument('-w', '--winning', help='Выигрышные номера (через запятую)')
    parser.add_argument('-s', '--stats', action='store_true', help='Показать статистику')
    parser.add_argument('-o', '--output', help='Файл для сохранения')
    parser.add_argument('-v', '--verbose', action='store_true', help='Подробный вывод')
    args = parser.parse_args()

    if args.num_numbers > args.max_number:
        print(colorize("Количество номеров не может превышать максимальное число.", 'red'))
        sys.exit(1)

    game = LotteryGenerator(args.num_numbers, args.max_number)

    if args.stats:
        game.stats()
        return

    winning_numbers = None
    if args.winning:
        winning_numbers = [int(x.strip()) for x in args.winning.split(',')]
        if len(winning_numbers) != args.num_numbers:
            print(colorize("Количество выигрышных номеров должно совпадать с num_numbers.", 'red'))
            sys.exit(1)

    tickets = []
    for _ in range(args.count):
        ticket = game.generate_ticket()
        tickets.append(ticket)
        game.save_history(ticket)

    # Вывод
    output_lines = []
    if args.verbose:
        for i, ticket in enumerate(tickets):
            line = f"Билет {i+1}: {' '.join(map(str, ticket))}"
            if winning_numbers:
                matches, matched = game.check_ticket(ticket, winning_numbers)
                if matches:
                    line += f"  Совпадений: {matches} ({', '.join(map(str, matched))})"
                    # Цветные совпадения
                    colored = []
                    for n in ticket:
                        if n in matched:
                            colored.append(colorize(str(n), 'green'))
                        else:
                            colored.append(str(n))
                    line = f"Билет {i+1}: {' '.join(colored)}  (совпадений: {matches})"
            output_lines.append(line)
    else:
        for ticket in tickets:
            output_lines.append(' '.join(map(str, ticket)))

    output = '\n'.join(output_lines)

    if args.output:
        with open(args.output, 'w') as f:
            f.write(output)
        print(colorize(f"Результат сохранён в {args.output}", 'green'))
    else:
        print(output)

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print(colorize("\nПрервано.", 'yellow'))
        sys.exit(0)
