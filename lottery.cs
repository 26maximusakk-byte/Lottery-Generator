// lottery.cs
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;

class LotteryGenerator
{
    static string Colorize(string text, string color)
    {
        string col = color switch
        {
            "green" => "\x1b[92m",
            "red" => "\x1b[91m",
            "yellow" => "\x1b[93m",
            "blue" => "\x1b[94m",
            "bold" => "\x1b[1m",
            _ => "\x1b[0m"
        };
        return col + text + "\x1b[0m";
    }

    private int numNumbers;
    private int maxNumber;
    private List<List<int>> history;
    private string historyFile;

    public LotteryGenerator(int num, int max)
    {
        numNumbers = num;
        maxNumber = max;
        historyFile = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), ".lottery_history.txt");
        LoadHistory();
    }

    void LoadHistory()
    {
        history = new List<List<int>>();
        if (File.Exists(historyFile))
        {
            foreach (var line in File.ReadLines(historyFile))
            {
                var nums = line.Split(' ', StringSplitOptions.RemoveEmptyEntries)
                               .Select(int.Parse).ToList();
                if (nums.Count > 0) history.Add(nums);
            }
        }
    }

    void SaveHistory(List<int> ticket)
    {
        history.Add(ticket);
        using var sw = new StreamWriter(historyFile, true);
        sw.WriteLine(string.Join(" ", ticket));
    }

    public List<int> GenerateTicket()
    {
        var rnd = new Random();
        var nums = new HashSet<int>();
        while (nums.Count < numNumbers)
            nums.Add(rnd.Next(1, maxNumber + 1));
        return nums.OrderBy(n => n).ToList();
    }

    public (int matches, List<int> matched) CheckTicket(List<int> ticket, List<int> winning)
    {
        var matched = ticket.Intersect(winning).ToList();
        return (matched.Count, matched);
    }

    public void ShowStats()
    {
        if (history.Count == 0)
        {
            Console.WriteLine(Colorize("Нет истории для статистики.", "yellow"));
            return;
        }
        var freq = new Dictionary<int, int>();
        foreach (var t in history)
            foreach (int n in t)
                freq[n] = freq.ContainsKey(n) ? freq[n] + 1 : 1;
        Console.WriteLine(Colorize("📊 Статистика выпадений:", "bold"));
        foreach (var kv in freq.OrderBy(kv => kv.Key))
            Console.WriteLine($"  {kv.Key,2}: {kv.Value} раз");
    }

    static void Main(string[] args)
    {
        int num = 0, max = 0, count = 1;
        string winningStr = null, outputFile = null;
        bool stats = false, verbose = false;

        for (int i = 0; i < args.Length; i++)
        {
            string arg = args[i];
            if (arg == "-h" || arg == "--help")
            {
                Console.WriteLine("Usage: lottery <num> <max> [-c N] [-w nums] [-s] [-o file] [-v]");
                return;
            }
            else if (arg == "-c" && i + 1 < args.Length)
                count = int.Parse(args[++i]);
            else if (arg == "-w" && i + 1 < args.Length)
                winningStr = args[++i];
            else if (arg == "-s")
                stats = true;
            else if (arg == "-o" && i + 1 < args.Length)
                outputFile = args[++i];
            else if (arg == "-v")
                verbose = true;
            else if (num == 0)
                num = int.Parse(arg);
            else if (max == 0)
                max = int.Parse(arg);
        }
        if (num <= 0 || max <= 0 || num > max)
        {
            Console.WriteLine(Colorize("Неверные параметры. Укажите num и max (num <= max).", "red"));
            return;
        }

        var game = new LotteryGenerator(num, max);

        List<int> winningNumbers = null;
        if (!string.IsNullOrEmpty(winningStr))
        {
            winningNumbers = winningStr.Split(',').Select(s => int.Parse(s.Trim())).ToList();
            if (winningNumbers.Count != num)
            {
                Console.WriteLine(Colorize("Количество выигрышных номеров должно совпадать с num.", "red"));
                return;
            }
        }

        if (stats)
        {
            game.ShowStats();
            return;
        }

        var tickets = new List<List<int>>();
        for (int i = 0; i < count; i++)
        {
            var t = game.GenerateTicket();
            tickets.Add(t);
            game.SaveHistory(t);
        }

        var outputLines = new List<string>();
        if (verbose)
        {
            for (int i = 0; i < tickets.Count; i++)
            {
                var t = tickets[i];
                string line = $"Билет {i+1}: ";
                if (winningNumbers != null)
                {
                    var res = game.CheckTicket(t, winningNumbers);
                    var colored = t.Select(n => res.matched.Contains(n) ? Colorize(n.ToString(), "green") : n.ToString());
                    line += string.Join(" ", colored) + $" (совпадений: {res.matches})";
                }
                else
                {
                    line += string.Join(" ", t);
                }
                outputLines.Add(line);
            }
        }
        else
        {
            foreach (var t in tickets)
                outputLines.Add(string.Join(" ", t));
        }

        string output = string.Join("\n", outputLines);
        if (!string.IsNullOrEmpty(outputFile))
        {
            File.WriteAllText(outputFile, output);
            Console.WriteLine(Colorize($"Результат сохранён в {outputFile}", "green"));
        }
        else
        {
            Console.WriteLine(output);
        }
    }
}
