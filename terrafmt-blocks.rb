#!/usr/bin/ruby
#take code/markdown as input on stdin, pull out tf blocks, format it with `terraform fmt` and insert pretty tf back in

require 'thor'
require 'colorize'
require 'open3'

#define program
# terrafmt-blocks file
# # cat f | terrafmt-blocks -i

BlockPair = Struct.new(:start, :end, :desc)

block_pairs = [
  BlockPair.new('```hcl',               '```', 'markdown'),
  BlockPair.new('return fmt.Sprintf(`', '`, ', 'markdown')

            ]

class BlkFmt

  @count = 0

  def initialize(mode, file=nil)
    raise Thor::Error, 'unknown BlkFmt mode'.red unless  [:fmt, :diff, :count].include?(mode)
    @mode = mode
  end

  def go

    starts = [
        "```hcl",
    ]
    ends = [
        "```",
    ]

    def line_has_patterns(line, patterns)
      patterns.each do |p|
        #puts "checking #{line} for #{p}"
        if line.match(Regexp.escape(p))
          # puts "!!!!!"
          return true
        end
      end

      return false
    end

    #true if we are buffering a block
    buffering = false

    #the current working block
    buffer = []

    STDIN.each_line do |line|

      #if we are capturing data to format
      if buffering

        #check to see if we should stop
        unless line_has_patterns(line, ends)
          #nope, add line to buffer and goto next line
          buffer << line
          next
        else
          #we are finish capturing data, format it!
          o, e, s = Open3.capture3("terraform fmt -", stdin_data: buffer.join(""))


          #check exit status
          if s.exitstatus == 0
            #success! output it and the closing line
            puts o
          else
            #have error, log it to stderr & output unformatted buffer
            STDERR.puts e
            puts buffer
          end

          #either way reset buffer and buffering value
          buffer = []
          buffering = false

          #and output closing line
          puts line
          next

        end
      end

      #if line includes a start match, start buffering
      if line_has_patterns(line, starts)
        buffering = true
      end

      puts line

    end
  end
end


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



TerraFmtBlocks.start(ARGV)
