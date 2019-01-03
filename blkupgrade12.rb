# frozen_string_literal: true

require 'colorize'
require 'fileutils'

require_relative 'blkreader.rb'

# format each block
class BlkUpgrade12 < BlkReader
  def initialize(file, quiet, diff, diff_context)
    super(file)
    @quiet = quiet
    @diff = diff
    @context = diff_context
    @output = []
  end

  def notblock_line_read(line)
    # only output non block lines if we are not in 'diff mode'
    @output << line unless @diff
  end

  def block_read(block)
    # upgrade command works on directories so
    Dir.mktmpdir('terrafmt') do |dir|
      File.open("#{dir}/main.tf", 'w') do |io|
        io.write block
      end

      _, error, status = Open3.capture3("terraform 0.12upgrade #{dir}", stdin_data: "yes\n")

      result = File.read("#{dir}/main.tf")

      # common error handling
      if !status.exitstatus.zero?
        print_msg(@file, @line_block_start, error)
        @blocks_err += 1
      else
        @blocks_ok += 1

        # see if different
        @blocks_diff += 1 if result != block
      end

      processed_block(block, result, status)
    end
  end

  # these duplicate other files, figure out a better way to handle this
  def processed_block(block, block_fmt, status)
    # output
    if @diff
      return unless status.exitstatus.zero?

      d = Diffy::Diff.new(block, block_fmt, context: @context)
      dstr = d.to_s(:color).strip

      return if dstr.empty?

      puts "#{@file}@#{@line_block_start}:".white.bold + " block ##{@blocks_found}".magenta
      puts dstr
      puts
      @exit_code = 1
    else
      @output << if status.exitstatus.zero?
                   block_fmt
                 else
                   block
                 end
    end
  end

  def done(io)
    return if @diff

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

      if @blocks_found.zero?
        puts "#{@file}:".light_white + ' no blocks found!'.yellow unless @quiet
      elsif @blocks_diff.zero?
        puts "#{@file}:".light_white + " #{@blocks_ok} already formatted".light_blue unless @quiet
      else
        puts "#{@file}:".light_white + " formatted #{@blocks_diff} blocks".green
      end
    end
  end
end
