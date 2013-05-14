class MakeIdRefsInteger < ActiveRecord::Migration
  def up
    change_column :characters, :play_id, :integer
    change_column :props, :section_id, :integer
    change_column :sections, :play_id, :integer
    change_column :sub_sections, :section_id, :integer
    change_column :lines, :sub_section_id, :integer
    change_column :lines, :character_id, :integer
  end

  def down
    change_column :characters, :play_id, :string
    change_column :props, :section_id, :string
    change_column :sections, :play_id, :string
    change_column :sub_sections, :section_id, :string
    change_column :lines, :sub_section_id, :string
    change_column :lines, :character_id, :string
  end
end
