#!/usr/bin/ruby

require 'colorize'
require 'diffy'
require 'fileutils'
require 'open3'
require 'thor'

#todo a better way to output message with filename line number etc
# funtion: message("text") handles rest, outputs as white

BlockPair = Struct.new(:start, :finish, :desc) do

  def starts?(line)
    return line.strip.start_with?(start)
  end

  def finishes?(line)
    return line.start_with?(finish)
  end

end

class BlkReader

  @@pairs = [
      BlockPair.new('```hcl', '```', 'markdown'),
      BlockPair.new('return fmt.Sprintf(`', '`,', 'acctest')
  ]

  def initialize(mode, file = nil, context=5)
    raise Thor::Error, 'ERROR unknown BlkFmt mode'.red unless [:fmt, :diff, :count].include?(mode)
    @mode = mode
    @file = file
    @contex = context

    @lines = 0
    @lines_block = 0
    @count = 0
  end

  def go

    if @file == nil
      io = $stdin
    else
      io  = File.open(@file, 'r+')
    end

    buffer = [] #the current block
    pair = nil  #current block pair we are in (not nil == buffering)

    io.each_line do |line|
      @lines += 1

      if pair != nil #if we have started a pair and should buffer
        @lines_block += 1
        unless pair.finishes? line #check to see if we are at the end of a block
          buffer << line #if not buffer line and goto next
          next
        else

          block = buffer.join("")
          block_fmt, error, status = Open3.capture3("terraform fmt -", stdin_data: block)

          #common error handling
          if status.exitstatus != 0
            if @file == nil
              STDERR.puts "STDIN@#{@line_block_start}:".white.bold + " #{error}"
            else
              STDERR.puts "#{@file}@#{@line_block_start}:".white.bold + " #{error}"
            end
          end

          block_read(line, block, block_fmt, status)

          #noe reset the buffer/pair
          buffer = []
          pair = nil
          next #skip to next line
        end
      end

      #see if any pairs start here
      @@pairs.each do |p|
        if p.starts? line
          @count += 1
          @line_block_start = @lines
          pair = p
          break
        end
      end

      #put starting line
      line_read(line)
    end

    done(io)

    if @file != nil
      io.close
    end
  end

  #after each line is read, default to output it (passthrough)
  def line_read(line)
    puts line
  end

  #block has been read ito buffer, line that finished the block is passed in
  def block_read(line, block, block_fmt, status)
    puts buffer
    puts line
  end


  def done(io)

  end
end


#format each block
class BlkFmt < BlkReader

  #todo blocks_err, blocks_found, blocks_formatted

  def initialize(mode, file)
    super(mode, file)
    @output = []
  end

  def line_read(line)
    @output  << line
  end

  def block_read(line, block, block_fmt, status)
    if status.exitstatus == 0
      @output << block_fmt
    else
      @output  << block
    end

    @output  << line
  end

  def done(io)
    if @file != nil #read from a file, so lets rewind it and write it back

      io.close

      tmp = Tempfile.new('terrafmt-blocks')
      tmp.write @output.join("")
      tmp.flush
      tmp.close
      FileUtils.mv(tmp.path, @file)

      #this should work but there are stange IO errors that occue, TODO investigate
      #io.rewind
      #io.puts @output
      #io.flush
      #io.close


      if @count == 0
        puts "#{@file}:".white + " no blocks found!".yellow
      end
      puts "#{@file}:".white + " formatted #{@count} blocks".green
    else
      STDOUT.puts @output
    end
  end

end

class BlkDiff < BlkReader
  def line_read(line)
    #prevent any non block lines
  end

  def block_read(line, block, block_fmt, status)
    d = Diffy::Diff.new(block, block_fmt)
    dstr = d.to_s(:color).strip
    unless dstr.empty?
      puts "#{@file}@#{@line_block_start}:".white.bold + " block ##{@count}".magenta
      puts dstr
    end
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

  def block_read(line, block, block_fmt, status)
    if block_fmt != block
      @count_diff += 1
    end
  end

  def done(io)
    if !@quiet || @count_diff > 0
      #todo add file name
      puts "#{@count_diff} blocks require formatting (of #{@count})"
    end
  end
end