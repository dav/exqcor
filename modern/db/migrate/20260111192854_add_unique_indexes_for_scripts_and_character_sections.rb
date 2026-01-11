class AddUniqueIndexesForScriptsAndCharacterSections < ActiveRecord::Migration[8.0]
  def change
    add_index :characters, [:script_id, :name], unique: true
    add_index :character_sections, [:character_id, :section_id], unique: true,
              name: "index_character_sections_on_character_id_and_section_id", if_not_exists: true
  end
end
