# frozen_string_literal: true

require 'colorize'
require 'open3'

require_relative 'blkpairs.rb'

# reads a file and finds blocks to work on
class BlkReader
  def initialize(file = nil)
    @file = file

    # line counters
    @lines = 0
    @lines_block = 0

    # stats
    @blocks_found = 0
    @blocks_ok = 0
    @blocks_err = 0
    @blocks_diff = 0

    @is_stdin = @file.nil?
    @file = '$STDIN' if @is_stdin
  end

  # common logging
  def print_msg(file, line, msg)
    STDERR.puts file.light_white + '@'.white + line.to_s.light_white + ' ' + msg
  end

  def go
    io = if @is_stdin
           $stdin
         else
           File.open(@file, 'r+')
         end

    buffer = [] # the current block
    pair = nil  # current block pair we are in (not nil == buffering)

    io.each_line do |line|
      @lines += 1 # count a line

      # are we in a block
      if pair.nil?
        # put starting line
        notblock_line_read(line)

        # see if any pairs start here (returns nil for no)
        pair = BlkPairs.check(line)
        unless pair.nil?
          @blocks_found += 1
          @line_block_start = @lines
        end
      else
        @lines_block += 1 # count a block line

        if pair.finishes? line

          # handle a fully read block
          block_read(buffer.join(''))

          # put closing line
          notblock_line_read(line)

          # now reset the buffer/pair
          buffer = []
          pair = nil
        else # check to see if we are at the end of a block
          buffer << line # if not buffer line and goto next
        end

      end
    end

    # if we get here still buffering there is a malformed block
    unless pair.nil?
      print_msg(@file, @line_block_start, "MALFORMED BLOCK: `#{pair.start}` missing `#{pair.finish}`".red)
      @blocks_err += 1
      # processed_block(block, block, status) #todo
    end

    done(io)

    io.close unless @is_stdin

    0
  end

  # after each  nonbock line is read, default to output it (passthrough)
  def notblock_line_read(line)
    puts line
  end

  def block_read(block)
    block_fmt, error, status = Open3.capture3('terraform fmt -', stdin_data: block)

    # common error handling
    if !status.exitstatus.zero?
      print_msg(@file, @line_block_start, error)
      @blocks_err += 1
    else
      @blocks_ok += 1

      # see if different
      @blocks_diff += 1 if block_fmt != block
    end

    processed_block(block, block_fmt, status)
  end

  # block after it has been formatted
  def processed_block(block, _block_fmt, _status)
    puts block
  end

  def done(_io); end
end
