
require 'rubocop/rake_task'

RuboCop::RakeTask.new(:rubocop) do |t|
  t.options = ['--display-cop-names']
end

task :style   => %I(rubocop )
task :test    => %I(style)
task :default => :test