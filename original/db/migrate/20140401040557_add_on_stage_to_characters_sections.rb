class AddOnStageToCharactersSections < ActiveRecord::Migration
  def change
    add_column :characters_sections, :on_stage, :boolean
    add_column :characters_sections, :id, :primary_key
    rename_table :characters_sections, :character_sections
  end
end
