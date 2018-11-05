# frozen_string_literal: true

#take code/markdown as input on stdin, pull out tf blocks, format it with `terraform fmt` and insert pretty tf back in

require 'thor'
require 'colorize'

#todo diff that only shows a couple lines before and after changes

#load class that does all the work
require_relative 'tfformatter_blkreader.rb'

#define the program
class TerraFmtBlocks < Thor

  def self.exit_on_failure?
    true
  end

  desc "fmt FILE", "format blocks of terraform found in FILE"
  long_desc <<-LONGDESC
  `terrafmt-blocks |fmt| |FILE|` will format blocks of terraform found in FILE or stdin if no file is specified
  LONGDESC

  option :stdin, type: :boolean, aliases: 'i'
  def fmt(file=nil)
    exit BlkFmt.new(file).go

  end

  desc "diff FILE", "will show a diff of what will be changed in the file"
  option :contex, type: :numeric, aliases: 'c'
  def diff(file=nil)
    exit BlkDiff.new(file, options[:context]).go
  end

  default_task :fmt
end

#run the program
TerraFmtBlocks.start(ARGV)
