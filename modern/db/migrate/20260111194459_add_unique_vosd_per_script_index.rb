class AddUniqueVosdPerScriptIndex < ActiveRecord::Migration[8.0]
  def change
    add_index :characters, :script_id,
              unique: true,
              where: "role = 'vosd'",
              name: "index_characters_on_script_id_unique_vosd"
  end
end
