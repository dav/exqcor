class Script < ApplicationRecord
  has_many :characters, inverse_of: :script, dependent: :delete_all
  has_many :sections, inverse_of: :script, dependent: :delete_all

  def sub_sections
    sections.flat_map(&:sub_sections)
  end

  def lines
    sub_sections.flat_map(&:lines)
  end

  def VOSD
    characters.find_by(role: "vosd")
  end

  def ensure_vosd!
    self.VOSD || create_vosd
  end

  def create_vosd
    Character.create!(
      script: self,
      name: "VOSD",
      role: "vosd",
      description: "The scenic narration: all stage directions are spoken by Voice of Stage Directions."
    )
  end

  def real_characters
    characters.reject { |c| c == self.VOSD }
  end

  def average_section_writing_duration
    return 0 if sections.empty?

    durs = sections.map(&:writing_duration)
    durs.sum.to_f / durs.size
  end

  def duplicate(options = {})
    new_script = Script.new
    new_script.title = options[:title] || title
    new_script.description = options[:description] || description
    new_script.save!

    characters.each do |character|
      next if character == self.VOSD

      new_character = character.dup
      new_character.actor = character.actor
      new_character.script = new_script
      new_character.save
      new_script.characters << new_character
    end

    sections.each do |section|
      new_section = Section.new
      new_section.script = new_script
      new_section.ordering = section.ordering
      new_section.description = section.description
      new_section.name = section.name
      new_section.save!
      new_section.reload

      section.character_sections.each do |cs|
        next if cs.character == self.VOSD
        character = Character.find(cs.character_id)
        new_character = new_script.characters.find_by(name: character.name)
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
          new_ss = new_section.sub_sections.first
          new_first_line = first_line.dup
          if new_first_line
            new_first_line.sub_section = new_ss
            new_first_line.character = new_script.VOSD
            new_first_line.ordering = 1
            new_ss.lines << new_first_line
          end
        end
      end
      new_script.sections << new_section
    end

    new_script.save!
    new_script
  end

  def sorted_sections
    sections.sort { |a, b| a.sort_name <=> b.sort_name }
  end
end
