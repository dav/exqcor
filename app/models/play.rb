class Play < ActiveRecord::Base
  attr_accessible :description, :title
  has_many :characters, :inverse_of => :play, :dependent => :delete_all
  has_many :sections, :inverse_of => :play, :dependent => :delete_all
  
  def lines
    lines = []
    sections.each do |s|
      s.sub_sections.each do |ss|
        lines += ss.lines
      end
    end
    lines
  end
  
end
