// lottery.java
import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.util.stream.*;

public class lottery {
    private static final String RESET = "\u001B[0m";
    private static final String GREEN = "\u001B[92m";
    private static final String RED = "\u001B[91m";
    private static final String YELLOW = "\u001B[93m";
    private static final String BLUE = "\u001B[94m";
    private static final String BOLD = "\u001B[1m";

    private static String colorize(String text, String color) {
        return color + text + RESET;
    }

    private int numNumbers;
    private int maxNumber;
    private List<List<Integer>> history = new ArrayList<>();
    private String historyFile;

    public lottery(int num, int max) {
        numNumbers = num;
        maxNumber = max;
        historyFile = System.getProperty("user.home") + "/.lottery_history.txt";
        loadHistory();
    }

    private void loadHistory() {
        try {
            List<String> lines = Files.readAllLines(Paths.get(historyFile));
            for (String line : lines) {
                if (line.trim().isEmpty()) continue;
                List<Integer> ticket = Arrays.stream(line.trim().split("\\s+"))
                        .map(Integer::parseInt)
                        .collect(Collectors.toList());
                if (!ticket.isEmpty()) history.add(ticket);
            }
        } catch (IOException e) {
            // file not exists
        }
    }

    private void saveHistory(List<Integer> ticket) throws IOException {
        history.add(ticket);
        try (FileWriter fw = new FileWriter(historyFile, true);
             BufferedWriter bw = new BufferedWriter(fw)) {
            for (int n : ticket) bw.write(n + " ");
            bw.newLine();
        }
    }

    public List<Integer> generateTicket() {
        Random rnd = new Random();
        Set<Integer> nums = new HashSet<>();
        while (nums.size() < numNumbers) {
            nums.add(rnd.nextInt(maxNumber) + 1);
        }
        return nums.stream().sorted().collect(Collectors.toList());
    }

    public int[] checkTicket(List<Integer> ticket, List<Integer> winning) {
        List<Integer> matched = ticket.stream()
                .filter(winning::contains)
                .collect(Collectors.toList());
        return new int[]{matched.size(), 0}; // second not used
    }

    public void showStats() {
        if (history.isEmpty()) {
            System.out.println(colorize("Нет истории для статистики.", YELLOW));
            return;
        }
        Map<Integer, Integer> freq = new HashMap<>();
        for (List<Integer> t : history) {
            for (int n : t) freq.put(n, freq.getOrDefault(n, 0) + 1);
        }
        System.out.println(colorize("📊 Статистика выпадений:", BOLD));
        for (int i = 1; i <= maxNumber; i++) {
            if (freq.containsKey(i)) {
                System.out.printf("  %2d: %d раз\n", i, freq.get(i));
            }
        }
    }

    public static void main(String[] args) throws IOException {
        int num = 0, max = 0, count = 1;
        String winningStr = null, outputFile = null;
        boolean statsFlag = false, verbose = false;

        for (int i = 0; i < args.length; i++) {
            String arg = args[i];
            if (arg.equals("-h") || arg.equals("--help")) {
                System.out.println("Usage: lottery <num> <max> [-c N] [-w nums] [-s] [-o file] [-v]");
                return;
            } else if (arg.equals("-c") && i+1 < args.length) {
                count = Integer.parseInt(args[++i]);
            } else if (arg.equals("-w") && i+1 < args.length) {
                winningStr = args[++i];
            } else if (arg.equals("-s")) {
                statsFlag = true;
            } else if (arg.equals("-o") && i+1 < args.length) {
                outputFile = args[++i];
            } else if (arg.equals("-v")) {
                verbose = true;
            } else if (num == 0) {
                num = Integer.parseInt(arg);
            } else if (max == 0) {
                max = Integer.parseInt(arg);
            }
        }
        if (num <= 0 || max <= 0 || num > max) {
            System.err.println(colorize("Неверные параметры. Укажите num и max (num <= max).", RED));
            return;
        }

        lottery game = new lottery(num, max);

        List<Integer> winningNumbers = null;
        if (winningStr != null) {
            winningNumbers = Arrays.stream(winningStr.split(","))
                    .map(s -> Integer.parseInt(s.trim()))
                    .collect(Collectors.toList());
            if (winningNumbers.size() != num) {
                System.err.println(colorize("Количество выигрышных номеров должно совпадать с num.", RED));
                return;
            }
        }

        if (statsFlag) {
            game.showStats();
            return;
        }

        List<List<Integer>> tickets = new ArrayList<>();
        for (int i = 0; i < count; i++) {
            List<Integer> ticket = game.generateTicket();
            tickets.add(ticket);
            game.saveHistory(ticket);
        }

        List<String> outputLines = new ArrayList<>();
        if (verbose) {
            for (int i = 0; i < tickets.size(); i++) {
                List<Integer> ticket = tickets.get(i);
                StringBuilder line = new StringBuilder("Билет " + (i+1) + ": ");
                if (winningNumbers != null) {
                    int matches = (int) ticket.stream().filter(winningNumbers::contains).count();
                    String colored = ticket.stream()
                            .map(n -> winningNumbers.contains(n) ? colorize(n.toString(), GREEN) : n.toString())
                            .collect(Collectors.joining(" "));
                    line.append(colored).append(" (совпадений: ").append(matches).append(")");
                } else {
                    line.append(ticket.stream().map(String::valueOf).collect(Collectors.joining(" ")));
                }
                outputLines.add(line.toString());
            }
        } else {
            for (List<Integer> ticket : tickets) {
                outputLines.add(ticket.stream().map(String::valueOf).collect(Collectors.joining(" ")));
            }
        }

        String output = String.join("\n", outputLines);
        if (outputFile != null) {
            Files.write(Paths.get(outputFile), output.getBytes());
            System.out.println(colorize("Результат сохранён в " + outputFile, GREEN));
        } else {
            System.out.println(output);
        }
    }
}
