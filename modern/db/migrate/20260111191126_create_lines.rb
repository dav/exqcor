class CreateLines < ActiveRecord::Migration[8.0]
  def change
    create_table :lines do |t|
      t.text :text
      t.integer :ordering
      t.references :sub_section, null: false, foreign_key: true
      t.references :character, null: false, foreign_key: true

      t.timestamps
    end
  end
end
