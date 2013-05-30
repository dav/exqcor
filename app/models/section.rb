class Section < ActiveRecord::Base
  attr_accessible :description, :name, :ordering, :play_id
  belongs_to :play, :inverse_of => :sections
  has_many :props, :inverse_of => :section, :dependent => :delete_all
  has_many :sub_sections, :inverse_of => :section, :order => 'ordering ASC', :dependent => :delete_all
  
  after_create :build_first_sub_section
  
  def next_ordering_index
    return 0 if self.sub_sections.size == 0
    
    max = -1
    self.sub_sections.each do |ss|
      if ss.ordering > max
        max = ss.ordering
      end
    end
    return max+1
  end
  
  def build_first_sub_section
    ss = SubSection.new(:ordering => next_ordering_index)
    ss.section = self
    return ss.save
  end
end
