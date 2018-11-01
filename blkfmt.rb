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

class BlkReader

  @@pairs = [
      BlockPair.new('```hcl', '```', 'markdown'),
      BlockPair.new('return fmt.Sprintf(`', '`, ', 'markdown')
  ]

  def initialize(mode, file = nil)
    raise Thor::Error, 'ERROR unknown BlkFmt mode'.red unless [:fmt, :diff, :count].include?(mode)
    @mode = mode

    @lines = 0
    @lines_block = 0
    @count = 0
  end

  def go

    @buffer = [] #the current block
    pair = nil  #current block pair we are in (not nil == buffering)

    STDIN.each_line do |line|
      @lines += 1

      if pair != nil #if we have started a pair and should buffer
        @lines_block += 1
        unless pair.finishes? line #check to see if we are at the end of a block
          @buffer << line #if not buffer line and goto next
          next
        else

          #block done! so call block_read function
          block_read(line)

          #no reset the buffer/pair
          @buffer = []
          pair = nil
          next #skip to next line
        end
      end

      #see if any pairs start here
      @@pairs.each do |p|
        if p.starts? line
          @count += 1
          pair = p
          break
        end
      end

      #put starting line
      line_read(line)
    end

    done
  end

  def run_format
    return Open3.capture3("terraform fmt -", stdin_data: @buffer.join(""))
  end

  #after each line is read, default to output it (passthrough)
  def line_read(line)
    puts line
  end

  #block has been read ito buffer, line that finished the block is passed in
  def block_read(line)
    puts @buffer
    puts line
  end

  def done

  end
end


#format each block
class BlkFmt < BlkReader
  def line_read(line)
    puts line
  end

  def block_read(line)
    o, e, s = run_format

    #check exit status
    if s.exitstatus == 0
      #success! output it and the closing line
      puts o
    else
      #have error, log it to stderr & output unformatted buffer
      STDERR.puts e
      puts @buffer
    end

    puts line
  end
end


class BlkCount < BlkReader

  def initialize(mode, quiet, file = nil)
    super(mode, file)
    @quiet = quiet
    @count_diff = 0
  end

  def line_read(line)

  end

  def block_read(line)
    o, e, s = run_format

    if s.exitstatus == 0
      if o != @buffer.join("")
        @count_diff += 1
      end
    else
      STDERR.puts e
    end
  end

  def done
    if !@quiet || @count_diff > 0
      puts "#{@count_diff} blocks require formatting (of #{@count})"
    end
  end
end