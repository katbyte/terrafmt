#!/usr/bin/ruby

require 'thor'
require 'colorize'
require 'open3'

BlockPair = Struct.new(:start, :finish, :desc) do

  def starts?(line)
    return line.start_with?(start)
  end

  def finishes?(line)
    return line.start_with?(finish)
  end

end

class BlkFmt

  @@pairs = [
      BlockPair.new('```hcl',               '```', 'markdown'),
      BlockPair.new('return fmt.Sprintf(`', '`, ', 'markdown')
  ]

  @count = 0

  def initialize(mode, file=nil)
    raise Thor::Error, 'unknown BlkFmt mode'.red unless  [:fmt, :diff, :count].include?(mode)
    @mode = mode
  end

  def go

    #current block pair we are in (not nil == buffering)
    pair = nil

    #the current working block
    buffer = []

    #for each line in file/stdin
    STDIN.each_line do |line|

      #if we are capturing data to format
      if pair != nil

        #check to see if we should stop
        unless pair.finishes? line
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
          pair = nil

          #and output closing line
          puts line
          next

        end
      end


      #see if any pairs start here
      @@pairs.each do |p|
        if p.starts? line
          pair = p
          break
        end
      end

      #put starting line
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
