# frozen_string_literal: true

require 'colorize'
require 'diffy'

require_relative 'blkreader.rb'

# shows a fmt diff for blocks
class BlkDiff < BlkReader
  def initialize(file, context)
    super(file)
    @context = context
    @exit_code = 0
  end

  def line_read(line)
    # prevent any non block lines for being output
  end

  def formatted_block(block, block_fmt, status)
    return unless status.exitstatus.zero?

    d = Diffy::Diff.new(block, block_fmt)
    dstr = d.to_s(:color).strip

    return if dstr.empty?

    puts "#{@file}@#{@line_block_start}:".white.bold + " block ##{@blocks_found}".magenta
    puts dstr
    puts
    @exit_code = 1
  end

  def go
    super
    @exit_code
  end
end
