# frozen_string_literal: true

# TODO multiple start and finish strings

# defines the start and end of a block
BlkPair = Struct.new(:start, :finish, :desc) do
  def starts?(line)
    line.strip.start_with?(start)
  end

  def finishes?(line)
    line.start_with?(finish)
  end
end

# collection of pairs and helper to find
class BlkPairs
  PAIRS = [
    BlkPair.new('```hcl', '```', 'markdown'),
    BlkPair.new('return fmt.Sprintf(`', '`,', 'acctest')
  ].freeze

  def self.check(line)
    PAIRS.each do |p|
      return p if p.starts? line
    end

    nil
  end
end
