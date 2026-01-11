class CreateCharacterSections < ActiveRecord::Migration[8.0]
  def change
    create_table :character_sections do |t|
      t.references :character, null: false, foreign_key: true
      t.references :section, null: false, foreign_key: true
      t.boolean :on_stage

      t.timestamps
    end

    add_index :character_sections, [:character_id, :section_id], name: "index_character_sections_on_character_id_and_section_id"
  end
end
