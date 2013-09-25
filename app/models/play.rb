class Play < ActiveRecord::Base
  attr_accessible :description, :title
  has_many :characters, :inverse_of => :play, :dependent => :delete_all
  has_many :sections, :inverse_of => :play, :dependent => :delete_all
  after_create :setup_play_after_create
  
  def setup_play_after_create
    create_vosd
  end
  
  def lines
    lines = []
    sections.each do |s|
      s.sub_sections.each do |ss|
        lines += ss.lines
      end
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
    vosd.description = 'Voice of Stage Directions. The scenic narration: all Stage Directions are spoken by VOSD'
    vosd.save!
    vosd
  end
  
end
