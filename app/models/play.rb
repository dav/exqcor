class Play < ActiveRecord::Base
  attr_accessible :description, :title
  has_many :characters, :inverse_of => :play, :dependent => :delete_all
  has_many :sections, :inverse_of => :play, :dependent => :delete_all
  after_create :setup_play_after_create
  
  def setup_play_after_create
    create_vosd
  end

  def sub_sections
    ss = []
    sections.each do |s|
      ss += s.sub_sections
    end
    ss
  end
      
  def lines
    lines = []
    self.sub_sections.each do |ss|
        lines += ss.lines
    end
    lines
  end
  
  def VOSD
    characters.find_by_name('VOSD') || create_vosd
  end
  
  def create_vosd
    vosd = Character.new
    vosd.play = self
    vosd.name = 'VOSD'
    vosd.description = 'The scenic narration: all stage directions are spoken by Voice of Stage Directions.'
    vosd.save!
    vosd
  end
  
  def real_characters
    self.characters.reject {|c| c == self.VOSD}
  end
  
  def average_section_writing_duration
    return 0 if self.sections.size == 0
    durs = []
    self.sections.each do |s|
      durs << s.writing_duration
    end
    
    durs.inject{ |sum, s| sum + s }.to_f / durs.size
  end
  
  
end
