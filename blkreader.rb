#!/usr/bin/ruby

require 'colorize'
require 'diffy'
require 'fileutils'
require 'open3'
require 'thor'

#defines the start and end of a block
BlockPair = Struct.new(:start, :finish, :desc) do

  def starts?(line)
    return line.strip.start_with?(start)
  end

  def finishes?(line)
    return line.start_with?(finish)
  end

end

#reads a file and finds blocks to work on
class BlkReader

  @@pairs = [
      BlockPair.new('```hcl', '```', 'markdown'),
      BlockPair.new('return fmt.Sprintf(`', '`,', 'acctest')
  ]

  def initialize(file = nil, context=5)
    @file = file
    @contex = context

    #line counters
    @lines = 0
    @lines_block = 0

    #stats
    @blocks_found = 0
    @blocks_ok = 0
    @blocks_err = 0
    @blocks_diff = 0
    @blocks_formatted = 0

    @is_stdin = @file.nil?
    if @is_stdin
       @file = 'STDIN'
    end
  end

  #common logging
  def print_msg(file, line, msg)
    STDERR.puts file.white.bold + '@'.white + line.to_s.white.bold + ' ' + msg
  end

  def go

    if @is_stdin
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
          block_fmt, error, status = Open3.capture3('terraform fmt -', stdin_data: block)

          #common error handling
          if status.exitstatus != 0
            print_msg(@file, @line_block_start, error)
            @blocks_err += 1
          else
            @blocks_ok += 1
          end

          #see if different
          if block_fmt != block
            @blocks_diff += 1
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
          @blocks_found += 1
          @line_block_start = @lines
          pair = p
          break
        end
      end

      #put starting line
      line_read(line)
    end

    #if we get here still buffering there is a malformed block
    if pair != nil
      print_msg(@file, @line_block_start, "MALFORMED BLOCK: `#{pair.start}` missing `#{pair.finish}`".red)
      @blocks_err += 1
    end

    r = done(io)

    if !@is_stdin
      io.close
    end

    return r
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
    return 0
  end
end


#format each block
class BlkFmt < BlkReader

  #todo blocks_err, blocks_found, blocks_formatted

  def initialize(file)
    super(file)
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
      else
        puts "#{@file}:".white + " formatted #{@count} blocks".green
      end
    else
      STDOUT.puts @output
    end

    return 0
  end

end


#shows a fmt diff for blocks
class BlkDiff < BlkReader
  def line_read(line)
    #prevent any non block lines
  end

  def block_read(line, block, block_fmt, status)
    if status.exitstatus == 0
      d = Diffy::Diff.new(block, block_fmt)
      dstr = d.to_s(:color).strip
      unless dstr.empty?
        puts "#{@file}@#{@line_block_start}:".white.bold + " block ##{@blocks_found}".magenta
        puts dstr
        return 1
      else
        return 0
      end
    end
  end
end


# counts blocks that need formatting
# is there any point in having this? maybe remove?
class BlkCount < BlkReader

  def initialize(quiet, file = nil)
    super(file)
    @quiet = quiet
  end

  #override to do nothing
  def line_read(line)
  end

  #override to do nothing
  def block_read(line, block, block_fmt, status)
  end

  def done(io)
    if !@quiet && @count_diff != 0
      if @is_stdin
        puts "#{@blocks_diff} blocks require formatting (of #{@blocks_found})"
        return 0
      else
        if @blocks_diff == 0
          msg = "no blocks found".yellow
        elsif @blocks_diff == 0
          msg = "#{@blocks_found} already formatted".green
        else
          msg = "#{@blocks_diff} of #{@blocks_found} require formatting".yellow
        end
        puts "#{@file}: ".white + msg
      end

      return 1
    end
  end
end