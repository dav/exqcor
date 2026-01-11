class AddRoleToCharacters < ActiveRecord::Migration[8.0]
  def up
    add_column :characters, :role, :string, default: "character"
    execute <<~SQL.squish
      UPDATE characters
      SET role = CASE WHEN name = 'VOSD' THEN 'vosd' ELSE 'character' END
    SQL
    change_column_null :characters, :role, false
  end

  def down
    remove_column :characters, :role
  end
end
