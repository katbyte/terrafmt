#!/usr/bin/ruby
#take code/markdown as input on stdin, pull out tf blocks, format it with `terraform fmt` and insert pretty tf back in

require 'thor'
require 'colorize'
require 'open3'

#load class that does all the work
require_relative 'blkfmt.rb'

#define the program
class TerraFmtBlocks < Thor

  desc "fmt FILE", "format blocks of terraform found in FILE"
  long_desc <<-LONGDESC
  `terrafmt-blocks |fmt| FILE` will format blocks of terraform found in FILE

      use --i FILE to read from stdin

      > $ cat file | ./terrafmt-blocks -i
  LONGDESC

  option :stdin, type: :boolean, aliases: 'i'
  def fmt(file=nil)
    BlkFmt.new(:fmt, file).go
  end

  desc "diff", "show diff"
  def count(name)
    puts "Hello #{name}"
  end

  desc "count", "counts the number of blocks # (# requiring format)"
  def count(name)
    puts "Hello #{name}"
  end



  default_task :fmt
end

#run the program
TerraFmtBlocks.start(ARGV)
