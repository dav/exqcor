class CreateLines < ActiveRecord::Migration
  def change
    create_table :lines do |t|
      t.string :text
      t.integer :ordering
      t.integer :sub_section_id
      t.integer :character_id

      t.timestamps
    end
  end
end
