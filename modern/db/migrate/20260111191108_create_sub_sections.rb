class CreateSubSections < ActiveRecord::Migration[8.0]
  def change
    create_table :sub_sections do |t|
      t.integer :ordering
      t.references :section, null: false, foreign_key: true
      t.references :writer, null: true, foreign_key: true

      t.timestamps
    end
  end
end
