# frozen_string_literal: true

require 'colorize'
require 'fileutils'

require_relative 'blkreader.rb'

# format each block
class BlkFmt < BlkReader
  def initialize(file, quiet)
    super(file)
    @quiet = quiet
    @output = []
  end

  def line_read(line)
    @output << line
  end

  def formatted_block(block, block_fmt, status)
    @output << if status.exitstatus.zero?
                 block_fmt
               else
                 block
               end
  end

  def done(io)
    if @is_stdin
      STDOUT.puts @output
    else
      io.close # close read file

      # HACK: because rewinding and writing out the file is buggy
      tmp = Tempfile.new('terrafmt-blocks')
      tmp.write @output.join('')
      tmp.flush
      tmp.close
      FileUtils.mv(tmp.path, @file)

      # io.rewind
      # io.puts @output
      # io.flush
      # io.close

      if !@quiet && @blocks_found.zero?
        puts "#{@file}:".light_white + ' no blocks found!'.yellow
      elsif !@quiet && @blocks_diff.zero?
        puts "#{@file}:".light_white + " #{@blocks_ok} already formatted".light_blue
      else
        puts "#{@file}:".light_white + " formatted #{@blocks_diff} blocks".green
      end
    end
  end
end
