#!/usr/bin/ruby
#take code/markdown as input on stdin, pull out tf blocks, format it with `terraform fmt` and insert pretty tf back in

require 'open3'

starts = [
    "```hcl",
   # "return fmt.Sprintf(`"
]
ends = [
    "```",
    #"`,",
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


buffering = false
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


