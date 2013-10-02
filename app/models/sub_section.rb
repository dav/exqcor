class SubSection < ActiveRecord::Base
  attr_accessible :ordering, :play_id
  belongs_to :section, :inverse_of => :sub_sections
  has_many :lines, :inverse_of => :sub_section, :dependent => :delete_all
  
  def next_ordering_index
    return 0 if self.lines.size == 0
    
    max = -1
    self.lines.each do |ss|
      if ss.ordering > max
        max = ss.ordering
      end
    end
    return max+1
  end
  
  def next_section
    sss = self.section.sub_sections
    sss.each_with_index do |sub_section, i|
      if sub_section==self
        if i+1 < sss.size
          return sss[i+1]
        end
      end
    end
    nil
  end
  
  def writing_duration
    # because currently we may prime the pump with an initial VOSD line, just toss out
    # the very first sub_section in a scene. It shouldn't affect the math much.
    return nil if self.lines.size < 2
    start_line = self.lines.min_by(&:created_at)
    showtime_lines = self.lines.reject {|line| line == start_line }
    first_line = showtime_lines.min_by(&:created_at)
    last_line = showtime_lines.max_by(&:created_at)
    return nil if last_line.nil? || first_line.nil?
    
    seconds = (last_line.created_at - first_line.created_at)
    seconds
  end
  
end
