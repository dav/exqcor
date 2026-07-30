class ChangeLineTextLength < ActiveRecord::Migration
  def up
    change_column :lines, :text, :text, :limit => 4000
  end

  def down
    change_column :lines, :text, :string
  end
end
