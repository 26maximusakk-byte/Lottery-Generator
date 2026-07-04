#!/usr/bin/env ruby
# lottery.rb
# encoding: UTF-8

require 'json'
require 'fileutils'

COLORS = {
  reset: "\e[0m",
  green: "\e[92m",
  red: "\e[91m",
  yellow: "\e[93m",
  blue: "\e[94m",
  bold: "\e[1m"
}

def colorize(text, color)
  "#{COLORS[color]}#{text}#{COLORS[:reset]}"
end

class Lottery
  attr_reader :num, :max, :history, :history_file

  def initialize(num, max)
    @num = num
    @max = max
    @history_file = File.join(Dir.home, '.lottery_history.txt')
    load_history
  end

  def load_history
    @history = []
    return unless File.exist?(@history_file)
    File.readlines(@history_file).each do |line|
      ticket = line.strip.split.map(&:to_i)
      @history << ticket if ticket.any?
    end
  end

  def save_history(ticket)
    @history << ticket
    File.open(@history_file, 'a') { |f| f.puts(ticket.join(' ')) }
  end

  def generate_ticket
    nums = Set.new
    while nums.size < @num
      nums.add(rand(1..@max))
    end
    nums.to_a.sort
  end

  def check_ticket(ticket, winning)
    matched = ticket & winning
    { matches: matched.size, matched: matched }
  end

  def show_stats
    if @history.empty?
      puts colorize('Нет истории для статистики.', :yellow)
      return
    end
    freq = Hash.new(0)
    @history.each { |t| t.each { |n| freq[n] += 1 } }
    puts colorize('📊 Статистика выпадений:', :bold)
    (1..@max).each do |n|
      puts "  #{n.to_s.rjust(2)}: #{freq[n]} раз" if freq[n] > 0
    end
  end
end

def main
  args = ARGV
  if args.empty? || args[0] == '-h' || args[0] == '--help'
    puts "Usage: ruby lottery.rb <num> <max> [-c N] [-w nums] [-s] [-o file] [-v]"
    exit
  end

  num = 0
  max = 0
  count = 1
  winning_str = nil
  output_file = nil
  stats_flag = false
  verbose = false

  i = 0
  while i < args.size
    case args[i]
    when '-c'
      count = args[i+1].to_i if i+1 < args.size
      i += 1
    when '-w'
      winning_str = args[i+1] if i+1 < args.size
      i += 1
    when '-s'
      stats_flag = true
    when '-o'
      output_file = args[i+1] if i+1 < args.size
      i += 1
    when '-v'
      verbose = true
    else
      if num == 0
        num = args[i].to_i
      elsif max == 0
        max = args[i].to_i
      end
    end
    i += 1
  end

  if num <= 0 || max <= 0 || num > max
    puts colorize('Неверные параметры. Укажите num и max (num <= max).', :red)
    exit 1
  end

  game = Lottery.new(num, max)

  winning_numbers = nil
  if winning_str
    winning_numbers = winning_str.split(',').map(&:to_i)
    if winning_numbers.size != num
      puts colorize('Количество выигрышных номеров должно совпадать с num.', :red)
      exit 1
    end
  end

  if stats_flag
    game.show_stats
    exit
  end

  tickets = []
  count.times do
    ticket = game.generate_ticket
    tickets << ticket
    game.save_history(ticket)
  end

  output_lines = []
  if verbose
    tickets.each_with_index do |ticket, idx|
      line = "Билет #{idx+1}: "
      if winning_numbers
        res = game.check_ticket(ticket, winning_numbers)
        colored = ticket.map { |n| res[:matched].include?(n) ? colorize(n.to_s, :green) : n.to_s }
        line += colored.join(' ') + " (совпадений: #{res[:matches]})"
      else
        line += ticket.join(' ')
      end
      output_lines << line
    end
  else
    tickets.each { |t| output_lines << t.join(' ') }
  end

  output = output_lines.join("\n")
  if output_file
    File.write(output_file, output)
    puts colorize("Результат сохранён в #{output_file}", :green)
  else
    puts output
  end
end

main if __FILE__ == $0
