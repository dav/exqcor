class AddStartOnStageToCharactersSections < ActiveRecord::Migration
  def change
    add_column :characters_sections, :start_on_stage, :boolean
  end
end
