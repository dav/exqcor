class Play < ActiveRecord::Base
  attr_accessible :description, :title
  has_many :characters, :inverse_of => :play, :dependent => :delete_all
  has_many :sections, :inverse_of => :play, :dependent => :delete_all
  after_create :setup_play_after_create
  
  def setup_play_after_create
    self.VOSD # side-effect TODO less ugh
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
    existing = characters.find_by_name('VOSD')
    existing || create_vosd
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

  def duplicate(options = {})
    new_play = Play.new
    new_play.title = options[:title] || self.title
    new_play.description = options[:description] || self.description
    new_play.save!

    self.characters.each do |character|
      next if character == self.VOSD

      new_character = character.dup
      new_character.actor = character.actor
      new_character.play = new_play
      new_character.save
      new_play.characters <<  new_character
    end

    self.sections.each do |section|
      new_section = Section.new
      new_section.play = new_play
      new_section.ordering = section.ordering
      new_section.description = section.description
      new_section.name= section.name
      new_section.save!
      new_section.reload

      section.character_sections.each do |cs|
        next if cs.character == self.VOSD
        character = Character.find(cs.character_id)
        new_character = new_play.characters.find_by_name(character.name)
        new_cs = CharacterSection.where(character_id: new_character.id, section_id: new_section.id)
        if new_cs.empty?
          new_cs = CharacterSection.new
          new_cs.character = new_character
          new_cs.section = new_section
          new_cs.on_stage = cs.on_stage
          new_cs.save!
        end
      end

      first_ss = section.sub_sections.first
      if first_ss
        first_line = first_ss.lines.first
        if first_line
          new_ss = new_section.sub_sections.first # this is auto created so already exists
          new_first_line = first_line.dup
          if new_first_line
            new_first_line.sub_section = new_ss
            new_first_line.character = new_play.VOSD
            new_first_line.ordering = 1
            new_ss.lines << new_first_line
          end
        end
      end
      new_play.sections << new_section
    end

    new_play.save!
    new_play
  end

  def sorted_sections
    sections.sort do |a, b|
      a.sort_name <=> b.sort_name
    end
  end
end
