class Section < ActiveRecord::Base
  attr_accessible :description, :name, :ordering, :play_id
  belongs_to :play, :inverse_of => :sections
  has_many :props, :inverse_of => :section, :dependent => :delete_all
  has_many :sub_sections, :inverse_of => :section, :order => 'ordering ASC', :dependent => :delete_all
  has_many :character_for_sections
  has_many :characters, :through => :character_for_sections
  
  after_create :build_first_sub_section, :add_vosd
  after_find :ensure_characters_include_vosd
  
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
  
  def ensure_characters_include_vosd
    add_vosd unless self.characters.include? self.play.VOSD
  end
  
  def add_vosd
    self.characters << self.play.VOSD
  end
  
  def build_first_sub_section
    ss = SubSection.new(:ordering => next_ordering_index)
    ss.section = self
    return ss.save
  end
  
  def sorted_characters
    [self.play.VOSD] + self.characters.reject {|c| c==self.play.VOSD}
  end
  
  def writing_duration
    return 0 if self.sub_sections.size == 0
    total = 0
    self.sub_sections.each do |ss|
      dur = ss.writing_duration
      next if dur.nil?
      total += dur
    end
    total
  end
  
end
