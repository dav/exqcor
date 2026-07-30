class CreateCharactersSections < ActiveRecord::Migration
  def self.up
    create_table :characters_sections, :id => false do |t|
        t.references :character
        t.references :section
    end
    add_index :characters_sections, [:character_id, :section_id]
    add_index :characters_sections, :section_id
  end

  def self.down
    drop_table :characters_sections
  end
end
