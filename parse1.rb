#!/usr/bin/env ruby

require 'time'

class Array
  def sum
    inject(0.0) { |result, el| result + el }
  end

  def mean 
    sum / size
  end

  def percentile(percentile)
    case self.size
    when 0
      return 0
    when 1
      return first
    end

    values_sorted = self.sort
    k = (percentile*(values_sorted.length-1)+1).floor - 1
    f = (percentile*(values_sorted.length-1)+1).modulo(1)

    return values_sorted[k] + (f * (values_sorted[k+1] - values_sorted[k]))
  rescue
    puts self.inspect
    raise
  end
end

def initialize_bucket
  $counters = Hash.new(0)
  $measures = Hash.new {|h,k| h[k] = [] }
end

CHEADERS = %w(end_SendMessage start_SendMessage received_connection end_connection)
DHEADERS = %w(end_SendMessage_duration end_connection_duration goroutines)

puts(['time', CHEADERS, DHEADERS.map {|n| "#{n}_95" } ].flatten.join(','))

def print(microseconds)
  list = [
    microseconds,
    CHEADERS.map {|name| $counters[name] },
    DHEADERS.map {|name|
      if list = $measures[name]
        list.percentile(0.95)
      else
        []
      end
    }
  ]

  puts list.flatten.join(',')
  initialize_bucket
end

def bucket(parts)
  at = nil
  parts.each do |p|
    k, v = p.split '='
    case k
    when 'at'
      at = v
      $counters[v] += 1
    when 'duration'
      $measures["#{at}_duration"] << v.to_i
    when 'goroutines'
      $measures["goroutines"] << v.to_i
    end
  end
end

previous_us = nil
initialize_bucket

File.open(ARGV[0]).each do |line|
  parts = line.split(' ')
  t = Time.parse(parts[1])

  line_us = t.nsec/10_000_000 + (t.sec * 100)
  if line_us != previous_us
    print(previous_us)
    previous_us = line_us
  end
  
  bucket(parts[2..-1])
end

print(previous_us)


