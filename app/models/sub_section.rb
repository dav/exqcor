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
end
