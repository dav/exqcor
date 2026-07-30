class Character < ActiveRecord::Base
  before_destroy :confirm_not_vosd

  belongs_to :play, :inverse_of => :characters
  belongs_to :actor, :inverse_of => :characters

  attr_accessible :description, :name, :actor_id

  has_many :lines, :inverse_of => :sub_section
  has_many :character_sections
  has_many :sections, :through => :character_sections
  validates_uniqueness_of :name, scope: :play_id
  
  def confirm_not_vosd
    return false if self.name == 'VOSD'
    true
  end
end
